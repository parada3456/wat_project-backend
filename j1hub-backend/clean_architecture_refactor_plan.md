# Implementation Plan: Clean Architecture Refactoring for j1hub-backend

## Overview
This plan outlines the refactoring necessary to bring the `j1hub-backend` into strict adherence with Clean Architecture and Go best practices. While the project already correctly structures layers (`domain`, `usecase`, `adapter`, `port`) and injects dependencies via `cmd/server/main.go`, we need to address leaking domain models, error handling coupling, and explicit DTO mapping boundaries.

## Architecture Decisions
- **Strict DTO Boundaries**: Handlers must not accept or return `domain` entities directly. The `usecase` layer should return UseCase Output DTOs, and the HTTP Handlers map these to HTTP JSON Responses.
- **Domain Sentinel Errors**: The `domain` package will define all business and state errors (`domain.ErrNotFound`, `domain.ErrInvalidInput`). The `usecase` and `repository` layers return these standard Go errors.
- **Centralized Error Translation**: Handlers (or an HTTP error middleware/wrapper) will translate `domain.Err*` to specific HTTP Status codes (e.g. `domain.ErrNotFound` -> `404 Not Found`). Explicit HTTP error structs (like `apperror.AppError`) will be removed from the core logic.
- **Payload Validation**: Input validation (e.g., using `go-playground/validator`) should happen at the HTTP adapter level before reaching the usecase.

## Task List

### Phase 1: Foundation (Errors & Validation)
- [ ] **Task 1: Define Domain Sentinel Errors**
  - **Description**: Add standard domain errors in `internal/domain/errors.go` (e.g., `ErrNotFound`, `ErrUnauthorized`, `ErrInvalidInput`, `ErrInternal`). 
  - **Acceptance criteria**:
    - [ ] `domain/errors.go` contains comprehensive sentinel errors.
    - [ ] Add an HTTP error translator in `internal/adapter/http/response` or `apperror` that maps these domain errors to HTTP responses.
  - **Verification**: 
    - [ ] Build succeeds: `go build ./...`
  - **Dependencies**: None
  - **Files likely touched**: `internal/domain/errors.go`, `pkg/apperror/apperror.go`
  - **Estimated scope**: Small (1-2 files)

### Checkpoint: Foundation
- [ ] Tests pass, builds clean
- [ ] Error definitions are available for all subsequent tasks.

### Phase 2: DTO Refactoring (Vertical Slices)

*We will slice the work vertically by domain to ensure testability at every step.*

- [ ] **Task 2: Refactor Authentication & User Domain**
  - **Description**: Update `LoginUseCase`, `RegisterUserUseCase`, and `UserUseCase` to return DTOs instead of `domain.User`. Update `AuthHandler` and `UserHandler` to use the centralized error translator instead of explicitly crafting `apperror.AppError`.
  - **Acceptance criteria**:
    - [ ] `Login`, `Register`, and Profile methods return DTOs.
    - [ ] Handlers do not import `domain` models directly for responses.
    - [ ] Handlers return errors which are translated automatically.
  - **Verification**: 
    - [ ] Tests pass: `go test ./internal/usecase/... -run "Test(Login|User|Register)"`
    - [ ] Build succeeds: `go build ./...`
  - **Dependencies**: Task 1
  - **Files likely touched**: `internal/usecase/login_user.go`, `internal/usecase/manage_user.go`, `internal/adapter/http/handler/auth_handler.go`, `internal/adapter/http/handler/user_handler.go`
  - **Estimated scope**: Medium (4-6 files)

- [ ] **Task 3: Refactor Jobs & Missions Domain**
  - **Description**: Update `ManageJobUseCase` and `MissionUseCase`. Ensure methods like `ListJobs` and `ListMissions` map arrays of domain models to array DTOs within the usecase or handler boundary strictly. Replace explicit handler errors with sentinel errors.
  - **Acceptance criteria**:
    - [ ] Usecases return DTO slices, preventing `domain.JobPosting` leaks to Handlers.
    - [ ] Handlers rely on centralized error translation.
  - **Verification**: 
    - [ ] Tests pass: `go test ./internal/usecase/... -run "Test(Job|Mission)"`
  - **Dependencies**: Task 1
  - **Files likely touched**: `internal/usecase/manage_job.go`, `internal/usecase/manage_mission.go`, `internal/adapter/http/handler/job_handler.go`, `internal/adapter/http/handler/mission_handler.go`
  - **Estimated scope**: Medium (4-6 files)

- [ ] **Task 4: Refactor Journeys, Friends & Expenses Domain**
  - **Description**: Update `JourneyUseCase`, `ManageFriendshipUseCase`, and `ManageExpenseUseCase` to return DTOs. Clean up error handling in their respective handlers.
  - **Acceptance criteria**:
    - [ ] Usecases return DTO slices.
    - [ ] Handlers rely on centralized error translation.
  - **Verification**: 
    - [ ] Build succeeds: `go build ./...`
  - **Dependencies**: Task 1
  - **Files likely touched**: `internal/usecase/manage_journey.go`, `internal/usecase/manage_friendship.go`, `internal/usecase/manage_expense.go`, corresponding handlers.
  - **Estimated scope**: Large (6-8 files)

### Checkpoint: Core Features
- [ ] End-to-end flows compile and all tests pass.
- [ ] No `apperror.AppError` explicit creation inside any handler.
- [ ] Handlers solely deal with parsing requests and formatting DTO responses.

### Phase 3: Polish & Repositories
- [ ] **Task 5: Refactor Repository Error Returns**
  - **Description**: Ensure Postgres repositories in `internal/adapter/postgres` do not leak `sql.ErrNoRows` or `pgx` errors. They must wrap or translate them to `domain.ErrNotFound` or `domain.ErrInternal`.
  - **Acceptance criteria**:
    - [ ] All `postgres` repositories return `domain.Err*` errors.
  - **Verification**: 
    - [ ] Tests pass.
  - **Dependencies**: Phase 2
  - **Files likely touched**: `internal/adapter/postgres/*.go`
  - **Estimated scope**: Medium (5-10 files)

### Checkpoint: Complete
- [ ] All acceptance criteria met
- [ ] Ready for human review

## Risks and Mitigations
| Risk | Impact | Mitigation |
|------|--------|------------|
| JSON Response Contract Breakage | High | Creating DTOs that exactly match the JSON struct tags of the current domain models will preserve API backwards compatibility. |
| Large PR Size | Medium | Slice refactoring into multiple PRs based on the task domains (User/Auth first, then Jobs/Missions, etc.). |

## Open Questions
- Is there a preferred validation library (e.g., `go-playground/validator`) already in use in `j1hub-backend` that we should standardize on for incoming requests?
