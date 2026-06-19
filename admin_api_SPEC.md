# Spec: J1 Hub Admin API Specification

## Objective
The J1 Hub Admin API provides administrators and moderators with the capability to manage the lifecycle of Thai J1 students, configure journey phases and compliance missions, audit points and streaks, manually adjust gamification rewards, manage the job directory, and moderate community-sourced job reviews.

Success looks like a clean, secure, and fully verified set of administrative APIs protected by JWT authentication with admin claim verification.

---

## Tech Stack
The Admin API aligns with the existing backend architecture:
- **Language**: Go 1.22+
- **HTTP Router**: `github.com/go-chi/chi/v5` via [router.go](file:///Users/user/development/work/WAT_project/backend/j1hub-backend/internal/adapter/http/handler/router.go)
- **Database Access**: PostgreSQL 16 via raw `pgx/v5` (no ORM)
- **Validation**: `github.com/go-playground/validator/v10`
- **Error Handling**: `pkg/apperror` package (mapping Go errors to standard JSON response envelopes)

---

## Commands
All operations are run inside the `backend/j1hub-backend` directory:
- **Build**: `make build`
- **Test**: `make test`
- **Lint**: `make lint`
- **Run Locally**: `make run` or `go run cmd/server/main.go`

---

## Project Structure
Administrative endpoints are integrated into the existing clean architecture layout:
- **Domain Models**: [internal/domain/](file:///Users/user/development/work/WAT_project/backend/j1hub-backend/internal/domain) (pure structs)
- **Ports (Interfaces)**: [internal/port/](file:///Users/user/development/work/WAT_project/backend/j1hub-backend/internal/port)
- **Usecases**: [internal/usecase/](file:///Users/user/development/work/WAT_project/backend/j1hub-backend/internal/usecase) (e.g., `admin_usecase.go`)
- **HTTP Handlers**: [internal/adapter/http/handler/](file:///Users/user/development/work/WAT_project/backend/j1hub-backend/internal/adapter/http/handler) (e.g., `admin_handler.go`)
- **Routing Setup**: [router.go](file:///Users/user/development/work/WAT_project/backend/j1hub-backend/internal/adapter/http/handler/router.go)

---

## Code Style
HTTP request validation is handled at the adapter layer, while authorization is checked using context values set by [auth.go](file:///Users/user/development/work/WAT_project/backend/j1hub-backend/internal/adapter/http/middleware/auth.go).

### Example Go Snippet:
```go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/j1hub/backend/internal/adapter/http/middleware"
	"github.com/j1hub/backend/internal/port"
	"github.com/j1hub/backend/pkg/apperror"
)

type AdminHandler struct {
	adminUseCase port.AdminUseCase
	validate     *validator.Validate
}

func NewAdminHandler(uc port.AdminUseCase, val *validator.Validate) *AdminHandler {
	return &AdminHandler{
		adminUseCase: uc,
		validate:     val,
	}
}

type VerifyMissionRequest struct {
	Approved        bool    `json:"approved"`
	RejectionReason *string `json:"rejectionReason" validate:"required_without=Approved"`
}

func (h *AdminHandler) VerifyMission(w http.ResponseWriter, r *http.Request) {
	userMissionID := chi.URLParam(r, "id")
	claims := middleware.GetClaims(r.Context())
	if claims == nil || !claims.IsAdmin {
		apperror.RespondError(w, apperror.ErrForbidden)
		return
	}

	var req VerifyMissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Malformed request body", Err: err})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Validation failed", Err: err})
		return
	}

	err := h.adminUseCase.VerifyMission(r.Context(), claims.UserID, userMissionID, req.Approved, req.RejectionReason)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"userMissionId": userMissionID,
		"status":        "Completed",
	})
}
```

---

## Testing Strategy
- **Unit Testing**: Implement `admin_usecase_test.go` using `testify/mock` to mock all repository ports. Validate state transitions (e.g. approving a mission triggers points engine calculations correctly).
- **Integration Testing**: Implement `admin_repo_integration_test.go` with `testcontainers-go` to run queries against a live PostgreSQL+PostGIS database. Verify correct schema behavior and table state.

---

## Admin API Endpoints

### Global Rules:
1. **Admin Middleware**: All endpoints specified below MUST reside within the router group protected by `middleware.Authenticate` and `middleware.RequireAdmin` filters.
2. **Response Envelope**: List endpoints use the standard data envelope:
   ```json
   {
     "data": [...],
     "pagination": {
       "page": 1,
       "pageSize": 10,
       "totalItems": 45,
       "totalPages": 5
     }
   }
   ```
3. **Standard Error Schema**: Follows the client-side error specification:
   ```json
   {
     "error": {
       "code": "FORBIDDEN",
       "message": "Admin access required",
       "details": null
     }
   }
   ```

---

### Endpoint Matrix

| Method | Endpoint | Description |
|---|---|---|
| **System Overview** | | |
| GET | `/api/v1/admin/dashboard/stats` | Retrieve aggregate system metrics |
| **Mission & Verification** | | |
| GET | `/api/v1/admin/user-missions` | List all user mission submissions pending review |
| PATCH | `/api/v1/admin/user-missions/{id}/verify` | Approve or reject a student's submitted mission proof |
| GET | `/api/v1/admin/missions` | Browse the complete catalog of missions |
| POST | `/api/v1/admin/missions` | Create a new catalog mission item |
| PUT | `/api/v1/admin/missions/{id}` | Edit catalog mission properties |
| DELETE | `/api/v1/admin/missions/{id}` | Delete a mission catalog item |
| POST | `/api/v1/admin/missions/{id}/tasks` | Append a task requirement to a mission |
| PUT | `/api/v1/admin/missions/{id}/tasks/{taskId}` | Edit details of a specific task requirement |
| DELETE | `/api/v1/admin/missions/{id}/tasks/{taskId}` | Remove a task requirement |
| **User & Gamification Control** | | |
| GET | `/api/v1/admin/users` | Search, filter, and list J1 users (paginated) |
| GET | `/api/v1/admin/users/{id}` | Retrieve comprehensive user record details |
| PATCH | `/api/v1/admin/users/{id}/status` | Block, suspend, or reactivate a user account |
| POST | `/api/v1/admin/users/{id}/adjust-points` | Manually adjust a user's points (creates `Admin_Adjust` ledger entry) |
| POST | `/api/v1/admin/users/{id}/badges` | Manually award an achievement badge |
| **Journey Phases** | | |
| GET | `/api/v1/admin/phases` | Retrieve all journey phases |
| POST | `/api/v1/admin/phases` | Configure a new journey phase |
| PUT | `/api/v1/admin/phases/{id}` | Update phase metadata |
| DELETE | `/api/v1/admin/phases/{id}` | Remove a journey phase |
| **Job Board & Moderation** | | |
| POST | `/api/v1/admin/jobs` | Manually publish a job posting |
| PUT | `/api/v1/admin/jobs/{id}` | Update job listing properties |
| DELETE | `/api/v1/admin/jobs/{id}` | Delete a job posting |
| POST | `/api/v1/admin/jobs/{id}/housing` | Configure housing options for a job listing |
| DELETE | `/api/v1/admin/reviews/{id}` | Moderate/delete an inappropriate job review |

---

### Endpoint Payloads & Examples

#### 1. System Stats (`GET /api/v1/admin/dashboard/stats`)
- **Response (200 OK)**:
  ```json
  {
    "totalUsers": 1250,
    "activeUsers": 1100,
    "pendingVerifications": 14,
    "activeJobs": 85,
    "averageCreditScore": 680,
    "totalPointsAwarded": 254000
  }
  ```

#### 2. Verify Submission (`PATCH /api/v1/admin/user-missions/{id}/verify`)
- **Request Body**:
  ```json
  {
    "approved": true,
    "rejectionReason": null
  }
  ```
- **Response (200 OK)**:
  ```json
  {
    "userMissionId": "ums_77102",
    "status": "Completed",
    "verifiedAt": "2026-06-18T13:25:00Z",
    "verifiedBy": "admin_001"
  }
  ```
- **Side Effects**: If `approved` is true, calculate streak bonuses, first-completer checks, write to `point_ledger`, issue earned badges, update `user` totals, and send push notification. If false, set status to `In_Progress` (or `Not_Started`) and notify user of rejection with `rejectionReason`.

#### 3. Adjust User Points (`POST /api/v1/admin/users/{id}/adjust-points`)
- **Request Body**:
  ```json
  {
    "pointsDelta": 150,
    "reason": "Outstanding contribution to community sharing"
  }
  ```
- **Response (200 OK)**:
  ```json
  {
    "userId": "usr_01ARZ3NDEKTSV4RRFFQ69G5FAV",
    "lifetimeBalanceAfter": 1400,
    "phaseBalanceAfter": 650,
    "ledgerId": "ldg_admin_adj_01"
  }
  ```

#### 4. Configure Mission (`POST /api/v1/admin/missions`)
- **Request Body**:
  ```json
  {
    "phaseId": "phs_001",
    "title": "Apply for Social Security Number (SSN)",
    "description": "Submit form SS-5 at the nearest Social Security Administration office.",
    "location": "SSA Office",
    "basePoints": 200,
    "isMandatory": true,
    "verification": "Upload",
    "dueDateType": "Relative",
    "relativeTriggerEvent": "arrival_date",
    "relativeDaysOffset": 7
  }
  ```
- **Response (201 Created)**:
  ```json
  {
    "missionId": "mis_01B982KDKW92KS83JSA29J2DKA",
    "phaseId": "phs_001",
    "title": "Apply for Social Security Number (SSN)",
    "basePoints": 200,
    "isMandatory": true,
    "verification": "Upload"
  }
  ```

---

## Boundaries

### Always Do:
- Validate that the calling user possesses the `is_admin` claim set to `true`.
- Write a record in [point_ledger](file:///Users/user/development/work/WAT_project/er_diagram/v5/schema_v5.sql#L180-L191) with type `Admin_Adjust` when adjusting points manually.
- Perform all double-entry points calculations inside database transactions (`pgx.Tx`).

### Ask First:
- Changes to existing database tables/schemas to accommodate audit logging.
- Introducing a soft-delete mechanism on missions or phases instead of raw cascades.

### Never Do:
- Delete users' earned points logs (`point_ledger`) permanently. All modifications must be additive adjustments.
- Bypass verification checks on target entities (e.g., awarding a badge that does not exist in the `badge` catalog).

---

## Success Criteria
1. Authorization validation blocks non-admin JWT payloads with a `403 Forbidden` response.
2. Direct CRUD endpoints for missions, phases, jobs, and badges function correctly and reflect changes in target database queries.
3. Verification endpoints trigger the reward calculation pipeline and correctly record points transactions in the ledger.
4. Manually updating a user's points modifies both target user balances and registers the audit trail in `point_ledger`.

---

## Open Questions
1. **Cascade Deletion Behavior**: Should deleting a catalog `mission` or `journey_phase` perform a cascade delete on user progression records (`user_mission`, `user_phase_history`), or should deletions be blocked if students have already initiated them?
2. **Audit Logging Log Table**: Should we implement a dedicated `admin_audit_log` table to capture non-financial actions like deactivating users, moderating reviews, or deleting catalog missions?
