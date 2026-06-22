# Implementation Plan: Sentinel Errors & Clean Architecture DTOs

## Overview
We need to refactor all API handlers to strictly adhere to Clean Architecture principles. This involves two major transformations. First, handlers must stop constructing explicit `apperror.AppError` structs (which leak HTTP status codes and error serialization logic into the delivery layer) and instead return standard sentinel errors defined in the `domain` package. Second, we must ensure that endpoints returning arrays/lists of domain models no longer leak those raw models; they must be mapped to array of Response DTOs before being handed off to `RespondList`.

## Architecture Decisions
- **Sentinel Errors via `domain` package**: Handlers will return `domain.ErrUnauthorized`, `domain.ErrInvalidInput`, etc. If custom messages are needed, they will wrap the sentinel error (e.g., `fmt.Errorf("%w: missing user_id", domain.ErrInvalidInput)`). `pkg/apperror` already knows how to map these to the correct HTTP status codes.
- **DTO Array Mapping**: Endpoints utilizing `apperror.RespondList(w, items...)` currently pass slices like `[]domain.JobPosting`. We will create explicit mappers (e.g. `NewJobPostingListResponse`) to map `[]domain.Entity` to `[]dto.EntityResponse`.
- **Vertical Slicing**: We will slice the work by handler groupings rather than horizontal layers to ensure each domain remains fully functional and testable at checkpoints.

## Task List

### Phase 1: Authentication & User Domain
- [ ] **Task 1: Add Missing Sentinel Errors**
  - **Description**: Add any missing sentinel errors to `domain/errors.go` (e.g., `ErrMalformedRequest`, `ErrMissingParameter`) and update `pkg/apperror/apperror.go` to map them to 400 Bad Request if they aren't already covered by `ErrInvalidInput`.
  - **Acceptance criteria**:
    - `domain/errors.go` contains necessary sentinel errors.
    - `pkg/apperror/apperror.go` handles mapping to HTTP status codes.
  - **Verification**: `go build ./...` succeeds.
  - **Dependencies**: None
  - **Files likely touched**: `domain/errors.go`, `pkg/apperror/apperror.go`
  - **Estimated scope**: XS
- [ ] **Task 2: Refactor Auth & User Handlers**
  - **Description**: Replace all `&apperror.AppError` instantiations with `domain.Err...`. Ensure any lists (if applicable) are mapped to DTO arrays.
  - **Acceptance criteria**:
    - No explicit `&apperror.AppError` usage remains in `auth_handler.go` or `user_handler.go`.
  - **Verification**: Tests pass, builds clean.
  - **Dependencies**: Task 1
  - **Files likely touched**: `auth_handler.go`, `user_handler.go`
  - **Estimated scope**: Small

### Checkpoint: Phase 1
- [ ] All tests pass
- [ ] Application builds without errors
- [ ] Auth and User domains are cleanly decoupled from `apperror.AppError` structs

### Phase 2: Core Business Logic (Missions, Jobs, Journeys)
- [ ] **Task 3: Refactor Job & Mission Handlers**
  - **Description**: Replace explicit errors with sentinel errors. Map all domain slices (e.g. `ListJobs`, `ListMissions`) to DTO slices before `RespondList`.
  - **Acceptance criteria**:
    - No `&apperror.AppError` usage.
    - `RespondList` only accepts slices of `dto.[Name]Response`.
  - **Verification**: `go build ./...`
  - **Dependencies**: Task 1
  - **Files likely touched**: `job_handler.go`, `mission_handler.go`, `dto/job_response.go`, `dto/mission_response.go`
  - **Estimated scope**: Medium
- [ ] **Task 4: Refactor Journey Handler**
  - **Description**: Replace explicit errors and implement DTO slice mappers for phases, leaderboards, and history.
  - **Acceptance criteria**:
    - No `&apperror.AppError` usage.
    - `RespondList` uses DTO slices.
  - **Verification**: `go build ./...`
  - **Dependencies**: Task 1
  - **Files likely touched**: `journey_handler.go`, `dto/journey_response.go`
  - **Estimated scope**: Medium

### Checkpoint: Phase 2
- [ ] Core feature compilation succeeds
- [ ] List endpoints successfully map to DTOs

### Phase 3: Social, Finance, & Administration
- [ ] **Task 5: Refactor Friend & Expense Handlers**
  - **Description**: Replace explicit errors with sentinel errors. Map lists (e.g. pending requests, expense lists) to DTOs.
  - **Acceptance criteria**:
    - No `&apperror.AppError` usage.
    - `RespondList` uses DTO slices.
  - **Verification**: `go build ./...`
  - **Dependencies**: Task 1
  - **Files likely touched**: `friend_handler.go`, `expense_handler.go`, `dto/friend_response.go`, `dto/expense_response.go`
  - **Estimated scope**: Medium
- [ ] **Task 6: Refactor Admin & Notification Handlers**
  - **Description**: Replace explicit errors with sentinel errors. Map lists to DTOs.
  - **Acceptance criteria**:
    - No `&apperror.AppError` usage.
    - `RespondList` uses DTO slices.
  - **Verification**: `go build ./...`
  - **Dependencies**: Task 1
  - **Files likely touched**: `admin_handler.go`, `notification_handler.go`, `dto/admin_response.go`, `dto/notification_response.go`
  - **Estimated scope**: Medium

### Checkpoint: Complete
- [ ] All acceptance criteria met
- [ ] Project builds cleanly (`go build ./...`)
- [ ] Existing tests run cleanly (`go test ./...`)
- [ ] Ready for review

## Risks and Mitigations
| Risk | Impact | Mitigation |
|------|--------|------------|
| Sentinel error wrapping alters JSON response format | High | We must ensure `pkg/apperror` properly unwraps or formats `err.Error()` without breaking existing frontend expectations for error messaging. |
| Redundant mapping overhead | Low | Mapping slices iteratively adds marginal CPU cost; negligible for standard pagination sizes. |

## Open Questions
- Should validation errors (like "Invalid request body") wrap `domain.ErrInvalidInput` or use a new sentinel error like `domain.ErrMalformedRequest`? I will default to wrapping `domain.ErrInvalidInput` to keep the domain error surface minimal and automatically map to 400 Bad Request.
