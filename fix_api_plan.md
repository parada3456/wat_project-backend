# Implementation Plan: User Jobs, Job CRUD, Password updates & Media APIs

## Overview
This plan implements the required backend routes and features to allow multiple jobs per user (with one main job), password updates, job CRUD operations, and media upload/delete APIs, ensuring the backend compiles and matches the routes in `router.go`.

## Architectural Decisions
- **Multiple User Jobs Database Migration**: We will create migration `000022_alter_user_jobs.up.sql` to add `is_main` (boolean), `start_date` (timestamp), and `end_date` (timestamp) to the `user_jobs` table.
- **Main Job Constraint**: In the `userRepo` implementation of `AssignJob`, we will implement a constraint: if `is_main` is set to `true`, we will automatically update all other jobs of the user to be `is_main = false`.
- **Dynamic File Storage**: Extend `StoragePort` in `internal/infrastructure/storage` to support `DeleteFile(ctx context.Context, bucket, key string) error` via Supabase REST storage API.
- **Media Module**: Create a new `internal/media` package to cleanly manage the media upload/deletion adapter, port, and use case, keeping details separate from user/job domains.

---

## Task List

### Phase 1: User Jobs & Schema Adjustments
#### Task 1: Create Database Migration
- **Description:** Add columns `is_main`, `start_date`, and `end_date` to `user_jobs` table.
- **Acceptance criteria:**
  - Database schema contains the new columns.
  - Migration up and down scripts successfully run.
- **Verification:**
  - Build successfully and migrations run on application startup.
- **Files likely touched:**
  - `migrations/000022_alter_user_jobs.up.sql`
  - `migrations/000022_alter_user_jobs.down.sql`
- **Estimated scope:** XS

#### Task 2: Domain, Port, and Repo Updates for User Jobs
- **Description:** Update `userdomain.UserJob` struct, `UserRepository` interface and postgres implementation. Add `FindUserJobs` to retrieve all jobs, and adjust `FindUserJob` to return the main job.
- **Acceptance criteria:**
  - `UserJob` domain struct contains `IsMain`, `StartDate`, and `EndDate`.
  - Repository methods read/write the new fields.
  - Mocks updated.
- **Verification:**
  - Go tests pass.
- **Files likely touched:**
  - `internal/user/domain/user.go`
  - `internal/user/port/user.go`
  - `internal/user/adapter/postgres/user_repo.go`
  - `internal/user/usecase/mocks_test.go`
  - `internal/infrastructure/notification/fcm_notifier_test.go`
- **Estimated scope:** M

#### Task 3: Usecase and Response DTO Updates for User Jobs
- **Description:** Update `GetProfile` usecase, `UserProfileResponse`, `GetProfileResponse` DTO to return `user_job` (main) and `user_jobs` (all). Update `AssignJob` to handle `is_main`, `start_date`, `end_date`.
- **Acceptance criteria:**
  - Endpoint `GET /api/v1/users/me` returns `user_job` and `user_jobs`.
  - `POST /api/v1/users/me/job` accepts parameters for `is_main`, `start_date`, `end_date`.
- **Verification:**
  - Usecase and handler tests pass.
- **Files likely touched:**
  - `internal/user/usecase/manage_user.go`
  - `internal/user/usecase/manage_user_test.go`
  - `internal/user/adapter/http/user_handler.go`
  - `internal/user/adapter/http/dto/user_response.go`
  - `internal/transport/http/test/user_handler_test.go`
- **Estimated scope:** M

---

### Checkpoint 1: Multiple Jobs Feature
- [ ] Database migrations execute successfully.
- [ ] User profile loads main job and all jobs list.
- [ ] Usecase unit tests and handler integration tests pass.

---

### Phase 2: User Password & Job CRUD APIs
#### Task 4: Implement Update Password API
- **Description:** Implement handler `userH.UpdatePassword` to let users update their password. Includes usecase, repository support if needed (or reusing `Update`), password verification via hasher.
- **Acceptance criteria:**
  - Endpoint `PUT /api/v1/users/me/password` validates current password and updates to new password.
- **Verification:**
  - Handler tests pass.
- **Files likely touched:**
  - `internal/user/usecase/manage_user.go`
  - `internal/user/adapter/http/user_handler.go`
  - `internal/transport/http/test/user_handler_test.go`
- **Estimated scope:** S

#### Task 5: Implement Job CRUD Handlers
- **Description:** Add `CreateJob`, `UpdateJob`, `PatchJob`, `DeleteJob` CRUD handlers to `JobHandler` and map their logic to `ManageJobUseCase` and repository layer.
- **Acceptance criteria:**
  - Create, fully update (PUT), partially update (PATCH), and delete (DELETE) job postings work.
- **Verification:**
  - `go build ./...` succeeds and compilation errors are resolved.
- **Files likely touched:**
  - `internal/job/port/job.go`
  - `internal/job/adapter/postgres/job_repo.go`
  - `internal/job/usecase/manage_job.go`
  - `internal/job/adapter/http/job_handler.go`
- **Estimated scope:** M

---

### Checkpoint 2: Core APIs
- [ ] Usecase tests for password updates and job CRUD pass.
- [ ] Compilation errors for `userH` and `jobH` resolved.

---

### Phase 3: Media Storage & Documentation
#### Task 6: Implement Media Upload & Delete Handler
- **Description:** Create media package to handle `/media/upload` and `/media//{key}` requests. Extend Supabase storage integration to implement file deletion. Register `mediaH` in `main.go` and `router.go`.
- **Acceptance criteria:**
  - Endpoint `POST /api/v1/media/upload` uploads a file and returns its public URL.
  - Endpoint `DELETE /api/v1/media/{key}` deletes the file.
- **Verification:**
  - `go test ./...` passes.
- **Files likely touched:**
  - `internal/infrastructure/storage/port.go`
  - `internal/infrastructure/storage/supabase_storage.go`
  - `internal/media/port/media.go`
  - `internal/media/adapter/http/media_handler.go`
  - `internal/media/usecase/manage_media.go`
  - `internal/transport/http/router.go`
  - `cmd/server/main.go`
- **Estimated scope:** L

#### Task 7: Update API list documentation
- **Description:** Update `api_list_v3.md` to document the new and updated endpoints, including parameters, request bodies, and expected responses.
- **Acceptance criteria:**
  - All endpoints match `router.go` exactly.
- **Verification:**
  - Document matches implementation.
- **Files likely touched:**
  - `api_list_v3.md`
- **Estimated scope:** S

---

### Checkpoint 3: Complete Project Delivery
- [ ] Whole backend builds successfully.
- [ ] All tests pass.
- [ ] API documentation matches code.
- [ ] Ready for final validation.
