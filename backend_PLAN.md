# J1 Hub Backend — Agent Execution Plan

> **How to use this file**: Feed `SKILL.md` + `PLAN.md` + all `.mmd` files to the agent together.
> The agent MUST complete every task in the order listed. Each task is atomic — finish it fully before moving to the next.
> After completing a task, mark it `[x]` and state which files were created/modified.

---

## Pre-flight (agent reads these before writing a single line of code)

- [ ] **P-1** Read `SKILL.md` in full
- [ ] **P-2** Read `er_diagram_v5.mmd` — memorize every table, field, type, and constraint
- [ ] **P-3** Read all 10 `seq_*.mmd` files — map each flow to its usecase file name
- [ ] **P-4** Confirm understanding: output a one-sentence summary of each sequence diagram before proceeding

---

## Phase 1 — Project Scaffold

### Task 1.1 — Repository Init
- [ ] Create `go.mod` (`module github.com/j1hub/backend`, `go 1.22`)
- [ ] Add all dependencies from `SKILL.md` Tech Stack table with exact versions
- [ ] Create `Makefile` with targets: `run`, `build`, `migrate-up`, `migrate-down`, `test`, `lint`, `docker-up`
- [ ] Create `docker-compose.yml` with services: `postgres` (postgis/postgis:16-3.4), `app`
- [ ] Create `.env.example` with all required keys:
  ```
  DATABASE_URL, JWT_SECRET, JWT_EXPIRY_HOURS,
  SUPABASE_URL, SUPABASE_SERVICE_KEY, SUPABASE_BUCKET_PROOFS, SUPABASE_BUCKET_SLIPS,
  FCM_CREDENTIALS_PATH,
  RADAR_RADIUS_METERS, RADAR_STALE_MINUTES,
  REWARD_SPEED_BONUS_7D_PCT, REWARD_SPEED_BONUS_1D_PCT,
  REWARD_STREAK_3_PCT, REWARD_STREAK_7_PCT,
  REWARD_FIRST_COMPLETER_FLAT,
  CRON_OVERDUE_EXPENSE, CRON_OVERDUE_MISSION,
  PORT
  ```

**Output files**: `go.mod`, `go.sum`, `Makefile`, `docker-compose.yml`, `.env.example`

---

### Task 1.2 — Directory Skeleton
- [ ] Create all empty directories matching the layout in `SKILL.md`
- [ ] Add `.gitkeep` to empty leaf dirs
- [ ] Create `cmd/server/main.go` with `// TODO: wire` placeholder

**Output files**: full directory tree, `cmd/server/main.go`

---

## Phase 2 — Infrastructure & Shared Packages

### Task 2.1 — Config (`internal/infrastructure/config/config.go`)
- [ ] Struct with all fields matching `.env.example` keys
- [ ] Load via `viper.AutomaticEnv()` + `.env` file
- [ ] `MustLoad() *Config` panics on missing required fields
- [ ] Reward bonus values exposed as typed `RewardConfig` sub-struct

### Task 2.2 — Database (`internal/infrastructure/db/postgres.go`)
- [ ] `NewPool(cfg *config.Config) (*pgxpool.Pool, error)`
- [ ] Pings on startup; logs pool stats
- [ ] `RunMigrations(pool, migrationsPath)` using `golang-migrate`

### Task 2.3 — Shared Packages
- [ ] `pkg/uid/uid.go` — `New(prefix string) string` using ULID
- [ ] `pkg/timeutil/clock.go` — `Clock` interface + `RealClock` + `MockClock` (for tests)
- [ ] `pkg/apperror/apperror.go` — `AppError{Code int, Message string, Err error}` + `respondError(w, err)` mapping `errors.Is` to HTTP codes

**Output files**: `config.go`, `postgres.go`, `uid.go`, `clock.go`, `apperror.go`

---

## Phase 3 — Domain Layer

> Rule: zero imports outside stdlib in this package.

### Task 3.1 — `internal/domain/user.go`
Fields from ER: `UserID`, `Email`, `PasswordHash`, `FirstName`, `LastName`, `CurrentPhaseID`, `TotalLifetimePoints`, `CurrentPhasePoints`, `MissionStreak`, `ArrivalDate`, `JobStartDate`, `CreatedAt`, `UpdatedAt`

- [ ] Struct definition
- [ ] `RadarVisibility` enum (`ShowAnonymous | ShowFriends | Hidden`) + `Valid()`

### Task 3.2 — `internal/domain/mission.go`
- [ ] `JourneyPhase` struct
- [ ] `Mission` struct with `VerificationType` enum (`None | Upload | Admin`) + `Valid()`
- [ ] `UserMission` struct with `UserMissionStatus` enum (5 values) + `Valid()`
- [ ] `Task` + `UserTask` structs
- [ ] `UserPhaseHistory` struct
- [ ] Business method: `func (m *Mission) CalculateDueDate(triggerDate time.Time) time.Time`
- [ ] Business method: `func CanAdvancePhase(missions []UserMission) bool` — returns true if all mandatory missions are Completed

### Task 3.3 — `internal/domain/gamification.go`
- [ ] `PointLedger` struct with `SourceType` enum (6 values) + `Valid()`
- [ ] `Badge` struct with `TriggerType` enum (5 values) + `Valid()`
- [ ] `UserBadge` struct
- [ ] `CreditScore` struct
- [ ] `PointReward` struct: `Base, SpeedBonus, StreakBonus, FirstCompleterBonus, Total int`

### Task 3.4 — `internal/domain/expense.go`
- [ ] `ExpenseTransaction` struct
- [ ] `ExpenseSplit` struct with `PaymentStatus` enum (4 values) + `ApprovalStatus` enum (3 values)
- [ ] Business method: `func (s *ExpenseSplit) IsSettled() bool`

### Task 3.5 — `internal/domain/job.go`
- [ ] `JobPosting`, `JobHousing`, `JobOverallRating`, `JobReview` structs
- [ ] `UserCart` struct with `CartStatus` enum (4 values) + `Valid()`
- [ ] Business method: `func (r *JobReview) ScoreMap() map[string]float64` — returns all score fields as a map for aggregate calculation

### Task 3.6 — `internal/domain/friendship.go`
- [ ] `Friendship` struct with `FriendshipStatus` enum (`Pending | Accepted | Blocked`) + `Valid()`
- [ ] Business method: `func CanonicalOrder(a, b string) (string, string)` — returns lexicographic min, max

### Task 3.7 — `internal/domain/errors.go`
- [ ] All 7 sentinel errors from `SKILL.md` + `ErrForbidden`, `ErrPhaseNotComplete`, `ErrProofAlreadySubmitted`

**Output files**: 7 domain files

---

## Phase 4 — Port Interfaces

### Task 4.1 — `internal/port/repository.go`
Define one interface per repository. Every method takes `ctx context.Context` first.

- [ ] `UserRepository` (Create, FindByID, FindByEmail, Update, IncrementPoints, ResetStreak, SetPhase)
- [ ] `ProfileRepository` (Create, FindByUserID, UpdateLocation, UpdateVisibility)
- [ ] `CreditScoreRepository` (Create, FindByUserID, Decrement)
- [ ] `JourneyPhaseRepository` (FindByNumber, FindByID)
- [ ] `UserPhaseHistoryRepository` (Insert, CompleteCurrentPhase, FindByUserAndPhase)
- [ ] `MissionRepository` (FindByPhase, FindByID)
- [ ] `UserMissionRepository` (BulkInsert, FindByUserAndPhase, FindByID, UpdateStatus, UpdateVerification, UpdateReward, FindOverdue)
- [ ] `TaskRepository` (FindByMission)
- [ ] `UserTaskRepository` (Upsert, FindByUserMission)
- [ ] `PointLedgerRepository` (Insert, InsertBatch)
- [ ] `BadgeRepository` (FindByTriggerType, FindEligible)
- [ ] `UserBadgeRepository` (Insert, FindByUser)
- [ ] `FriendshipRepository` (Insert, FindByCanonicalPair, FindByID, UpdateStatus, FindFriendsOf)
- [ ] `ExpenseTransactionRepository` (Insert, FindByID, MarkSettled)
- [ ] `ExpenseSplitRepository` (BulkInsert, FindByID, UpdatePaymentStatus, UpdateApproval, FindOverdue, CountUnsettled)
- [ ] `JobPostingRepository` (FindWithFilters, FindByID)
- [ ] `JobHousingRepository` (FindByJobID)
- [ ] `JobOverallRatingRepository` (FindByJobID, Recalculate)
- [ ] `JobReviewRepository` (Insert, FindByJobID)
- [ ] `UserCartRepository` (Insert, FindByUserAndJob, FindByID, UpdateStatus)
- [ ] `RadarRepository` (FindNearby)
- [ ] `LeaderboardRepository` (FindByScope)

### Task 4.2 — `internal/port/service.go`
- [ ] `PasswordHasher` interface (`Hash(plain string) (string, error)`, `Verify(plain, hash string) bool`)
- [ ] `TokenIssuer` interface (`Issue(userID string, isAdmin bool) (string, error)`, `Verify(token string) (*Claims, error)`)
- [ ] `StoragePort` interface (`UploadFile(ctx, bucket, key string, data io.Reader, contentType string) (url string, err error)`)
- [ ] `NotifierPort` interface (`Send(ctx, userID, title, body string) error`)

**Output files**: `repository.go`, `service.go`

---

## Phase 5 — Database Migrations

### Task 5.1 — All 20 migration pairs
Write `.up.sql` and `.down.sql` for every table in the order specified in `SKILL.md`.

Rules per migration:
- [ ] Enable PostGIS in `000001` (`CREATE EXTENSION IF NOT EXISTS postgis;`)
- [ ] All PK columns: `TEXT PRIMARY KEY`
- [ ] Enums as `TEXT NOT NULL CHECK (col IN (...))`
- [ ] All timestamps: `TIMESTAMPTZ NOT NULL DEFAULT NOW()`
- [ ] Nullable timestamps: `TIMESTAMPTZ`
- [ ] `geometry(Point, 4326)` for `current_coordinates` in `profiles`
- [ ] Create spatial index: `CREATE INDEX ON profiles USING GIST (current_coordinates);`
- [ ] All FK references use `REFERENCES table(id) ON DELETE CASCADE` unless noted
- [ ] Friendship unique constraint: `UNIQUE (user_id_1, user_id_2)` with comment that app enforces canonical order
- [ ] Each `.down.sql` is the exact reverse (`DROP TABLE IF EXISTS ... CASCADE`)

**Output files**: `migrations/000001_*.up.sql` through `migrations/000020_*.up.sql` + matching `.down.sql` (40 files)

---

## Phase 6 — Repository Implementations

### Task 6.1 — `internal/adapter/postgres/user_repo.go`
- [ ] Implements `port.UserRepository`
- [ ] All 7 methods with raw `pgx` queries
- [ ] Map pgx rows → domain structs manually (no ORM)

### Task 6.2 — `internal/adapter/postgres/profile_repo.go`
- [ ] Implements `port.ProfileRepository`
- [ ] `UpdateLocation` uses `ST_SetSRID(ST_MakePoint($1,$2),4326)`

### Task 6.3 — `internal/adapter/postgres/mission_repo.go`
- [ ] Implements `MissionRepository`, `UserMissionRepository`, `TaskRepository`, `UserTaskRepository`, `UserPhaseHistoryRepository`, `JourneyPhaseRepository`

### Task 6.4 — `internal/adapter/postgres/gamification_repo.go`
- [ ] Implements `PointLedgerRepository`, `BadgeRepository`, `UserBadgeRepository`, `CreditScoreRepository`

### Task 6.5 — `internal/adapter/postgres/friendship_repo.go`
- [ ] Implements `FriendshipRepository`

### Task 6.6 — `internal/adapter/postgres/expense_repo.go`
- [ ] Implements `ExpenseTransactionRepository`, `ExpenseSplitRepository`
- [ ] `FindOverdue` joins `EXPENSE_SPLIT` + `EXPENSE_TRANSACTION` on due_date

### Task 6.7 — `internal/adapter/postgres/job_repo.go`
- [ ] Implements `JobPostingRepository`, `JobHousingRepository`, `JobOverallRatingRepository`, `JobReviewRepository`, `UserCartRepository`
- [ ] `Recalculate` does `SELECT AVG(score_*), COUNT(*) FROM job_reviews WHERE job_id=$1` and updates `job_overall_ratings`

### Task 6.8 — `internal/adapter/postgres/radar_repo.go`
- [ ] Implements `RadarRepository`
- [ ] `FindNearby` uses `ST_DWithin(current_coordinates::geography, ST_SetSRID(ST_MakePoint($1,$2),4326)::geography, $3)` with 30-min staleness filter

### Task 6.9 — `internal/adapter/postgres/leaderboard_repo.go`
- [ ] Implements `LeaderboardRepository`
- [ ] `FindByScope(scope, jobID string)` returns ranked list joined with USER_CART

**Output files**: 9 repo files

---

## Phase 7 — External Adapters

### Task 7.1 — `internal/adapter/storage/supabase_storage.go`
- [ ] Implements `port.StoragePort`
- [ ] Multipart POST to Supabase Storage REST API
- [ ] Returns public URL

<!-- ### Task 7.2 — `internal/adapter/notification/fcm_notifier.go`
- [ ] Implements `port.NotifierPort`
- [ ] Uses `firebase.google.com/go/v4/messaging`
- [ ] Looks up FCM token by `userID` (add `fcm_token TEXT` to USER table migration if not present — add as `000001b` patch)
- [ ] Gracefully no-ops if token is empty (user hasn't granted push permission) -->

### Task 7.3 — `internal/adapter/auth/argon2_hasher.go`
- [ ] Implements `port.PasswordHasher`
- [ ] `Hash` uses `argon2.IDKey` with recommended params (time=1, memory=64MB, threads=4, keyLen=32)
- [ ] `Verify` constant-time compare

### Task 7.4 — `internal/adapter/auth/jwt_issuer.go`
- [ ] Implements `port.TokenIssuer`
- [ ] Claims: `sub=userID`, `is_admin=bool`, `exp`
- [ ] `Verify` returns typed `*Claims` struct

**Output files**: `supabase_storage.go`, `fcm_notifier.go`, `argon2_hasher.go`, `jwt_issuer.go`

---

## Phase 8 — Usecases (one per sequence diagram)

> For every usecase: constructor accepts interfaces only, never concrete types.
> Wrap multi-table writes in `pgxpool.BeginTx`.

### Task 8.1 — `register_user.go` (seq_01)
- [ ] `Register(ctx, cmd) (*domain.User, token string, error)`
  1. Hash password
  2. `BEGIN TX` → INSERT user → INSERT profile → INSERT credit_score → `COMMIT`
  3. Return user + JWT
- [ ] `InitializeJourney(ctx, userID, cmd) error`
  1. UPDATE user (arrival_date, job_start_date)
  2. Find phase 1
  3. `BEGIN TX` → UPDATE user.current_phase_id → INSERT user_phase_history → bulk INSERT user_missions with calculated_due_date → `COMMIT`

### Task 8.2 — `complete_mission.go` (seq_02)
- [ ] `ListMissions(ctx, userID) ([]MissionWithTasks, error)`
- [ ] `CompleteTask(ctx, userID, userMissionID, taskID) error`
- [ ] `SubmitProof(ctx, userID, userMissionID string, file io.Reader, contentType string) error`
  1. Upload to storage
  2. UPDATE user_mission: status=Pending_Verification, proof_url, proof_submitted_at
- [ ] `VerifyMission(ctx, adminID, userMissionID string, approved bool) error` (admin only)
  1. UPDATE user_mission: verified_at, verified_by
  2. Call RewardEngine.Calculate
  3. `BEGIN TX` → UPDATE user_mission (points + status=Completed) → UPDATE user (points, streak) → INSERT point_ledger entries → check+INSERT user_badges → `COMMIT`
  4. Send push notification

### Task 8.3 — `advance_phase.go` (seq_03)
- [ ] `TryAdvancePhase(ctx, userID) (advanced bool, error)`
  1. Load all mandatory USER_MISSION for current phase
  2. Call `domain.CanAdvancePhase(missions)`
  3. If true: `BEGIN TX` → snapshot USER_PHASE_HISTORY → UPDATE user (next phase, reset current_phase_points=0) → INSERT new USER_PHASE_HISTORY → bulk INSERT new USER_MISSION → INSERT POINT_LEDGER audit row → `COMMIT`
  4. Send push notification

### Task 8.4 — `manage_friendship.go` (seq_04)
- [ ] `SendRequest(ctx, senderID, targetID string) error`
  1. Verify target exists
  2. Canonical order via `domain.CanonicalOrder`
  3. Check duplicate
  4. INSERT FRIENDSHIP (status=Pending)
  5. Send push to target
- [ ] `RespondToRequest(ctx, responderID, friendshipID, accept bool) error`
  1. Load FRIENDSHIP, verify responder is user_id_2
  2. UPDATE status → Accepted or Blocked
  3. If accepted: send push to requester

### Task 8.5 — `manage_expense.go` (seq_05)
- [ ] `CreateExpense(ctx, payerID string, cmd CreateExpenseCmd) error`
  1. Validate no split has user_id == payerID
  2. `BEGIN TX` → INSERT EXPENSE_TRANSACTION → bulk INSERT EXPENSE_SPLIT → `COMMIT`
  3. Send push to each debtor
- [ ] `SubmitSlip(ctx, debtorID, txnID, splitID string, file io.Reader, contentType string) error`
  1. Verify ownership
  2. Upload slip
  3. UPDATE EXPENSE_SPLIT (payment_status=Submitted, payslip_url)
  4. Send push to payer
- [ ] `ApproveSplit(ctx, payerID, txnID, splitID string) error`
  1. Verify ownership (payer)
  2. `BEGIN TX` → UPDATE EXPENSE_SPLIT (Approved, settled_at) → check if all splits settled → if yes UPDATE EXPENSE_TRANSACTION.updated_at → `COMMIT`
  3. Send push to debtor

### Task 8.6 — `overdue_expense_job.go` (seq_06)
- [ ] `RunOverdueExpenseJob(ctx) error`
  1. `FindOverdue` (batch SELECT)
  2. For each: UPDATE split → Overdue, decrement credit_score, INSERT POINT_LEDGER, send push
  3. Log total processed count

### Task 8.7 — `manage_job.go` (seq_07)
- [ ] `ListJobs(ctx, filter JobFilter) ([]JobWithRating, error)`
- [ ] `GetJob(ctx, jobID string) (*JobDetail, error)`
- [ ] `AddToCart(ctx, userID, jobID string) error`
  1. Check duplicate
  2. INSERT USER_CART (status=Saved)
- [ ] `UpdateCartStatus(ctx, userID, cartID string, status domain.CartStatus) error`
- [ ] `WriteReview(ctx, userID, jobID string, cmd ReviewCmd) error`
  1. INSERT JOB_REVIEW
  2. `Recalculate` JOB_OVERALL_RATING

### Task 8.8 — `radar.go` (seq_08)
- [ ] `UpdateLocation(ctx, userID string, lat, lng float64) error`
- [ ] `GetRadar(ctx, requesterID string) ([]RadarEntry, error)`
  1. Load requester location from PROFILE
  2. `FindNearby` from RadarRepository
  3. For each result: check friendship status, apply visibility masking rules in Go

### Task 8.9 — `overdue_mission_job.go` (seq_09)
- [ ] `RunOverdueMissionJob(ctx) error`
  1. `FindOverdue` USER_MISSION (batch SELECT)
  2. For each: UPDATE status=Overdue, if mandatory → reset streak + INSERT POINT_LEDGER, send push
  3. Log total processed

### Task 8.10 — `leaderboard.go` (seq_10)
- [ ] `GetLeaderboard(ctx, scope, jobID string) ([]LeaderboardEntry, error)`
  1. `FindByScope` from LeaderboardRepository
  2. Load USER_BADGE for ranked users
  3. Apply display name masking (Hidden → "J1 Student #N", else first_name + last_initial)
  4. Attach rank number

**Output files**: 10 usecase files + `reward_engine.go`

---

## Phase 9 — HTTP Layer

### Task 9.1 — Middleware
- [ ] `internal/adapter/http/middleware/auth.go`
  - `Authenticate` — validates JWT, injects `userID` into ctx
  - `RequireAdmin` — checks `is_admin` claim
- [ ] `internal/adapter/http/middleware/logger.go`
  - Request/response logger using `slog`, logs method, path, status, latency

### Task 9.2 — Handlers (one file per domain)
Each handler file must:
- [ ] Define request/response structs with `json` + `validate` tags
- [ ] Call `validator.Validate` on request before usecase
- [ ] Call usecase method
- [ ] Use `respondError(w, err)` from `pkg/apperror`
- [ ] Return JSON response with correct HTTP status

Files:
- [ ] `auth_handler.go` — POST /auth/register, POST /auth/login
- [ ] `user_handler.go` — GET /users/me, PATCH /users/me
- [ ] `mission_handler.go` — GET /missions, PATCH tasks, POST proof
- [ ] `admin_handler.go` — PATCH /admin/user-missions/:id/verify
- [ ] `friendship_handler.go` — GET/POST/PATCH /friends
- [ ] `expense_handler.go` — POST /expenses, POST slip, PATCH approve
- [ ] `job_handler.go` — GET /jobs, GET /jobs/:id, POST/PATCH /cart, POST /reviews
- [ ] `radar_handler.go` — PATCH /profile/location, GET /radar
- [ ] `leaderboard_handler.go` — GET /leaderboard

### Task 9.3 — `internal/adapter/http/router.go`
- [ ] Wire all routes from `SKILL.md` API Routes section
- [ ] Apply `Authenticate` middleware to all non-auth routes
- [ ] Apply `RequireAdmin` middleware to `/admin/*` routes
- [ ] Mount chi middleware: `RealIP`, `RequestID`, logger

**Output files**: `auth.go`, `logger.go`, 9 handler files, `router.go`

---

## Phase 10 — Scheduler & Wiring

### Task 10.1 — `internal/infrastructure/scheduler/cron.go`
- [ ] `NewScheduler(cfg *config.Config, overdueExpenseJob, overdueMissionJob usecases) *cron.Cron`
- [ ] Register seq_06 job at `cfg.CronOverdueExpense` (default `"0 0 * * *"`)
- [ ] Register seq_09 job at `cfg.CronOverdueMission` (default `"0 0 * * *"`)
- [ ] Recover from panics in job wrapper

### Task 10.2 — `cmd/server/main.go`
Wire everything in dependency order:
```
config → db pool → run migrations →
repos → adapters (storage, fcm, hasher, jwt) →
usecases → handlers → router →
scheduler.Start() → http.ListenAndServe()
```
- [ ] Graceful shutdown on SIGTERM/SIGINT (drain in-flight requests, stop cron)
- [ ] Log startup info (port, DB URL masked)

**Output files**: `cron.go`, `main.go`

---

## Phase 11 — Tests

### Task 11.1 — Unit Tests (usecase layer)
One `_test.go` per usecase file. Use `testify/mock` for all port interfaces.

- [ ] `register_user_test.go` — test Register + InitializeJourney happy path + duplicate email error
- [ ] `complete_mission_test.go` — test VerifyMission: speed bonus, streak bonus, first completer, badge award
- [ ] `advance_phase_test.go` — test phase advances when all mandatory complete; does NOT advance when optional incomplete
- [ ] `manage_friendship_test.go` — test canonical order enforced; duplicate rejected
- [ ] `manage_expense_test.go` — test self-split rejected; ApproveSplit closes transaction when all settled
- [ ] `overdue_expense_job_test.go` — mock FindOverdue returns 3 splits; verify 3 credit decrements + 3 ledger inserts
- [ ] `manage_job_test.go` — test AddToCart duplicate prevention; WriteReview triggers Recalculate
- [ ] `radar_test.go` — test masking: Hidden → excluded, ShowAnonymous → masked, ShowFriends + friend → full name
- [ ] `overdue_mission_job_test.go` — mandatory overdue resets streak; optional overdue does not
- [ ] `leaderboard_test.go` — Hidden users masked to "J1 Student #N"
- [ ] `reward_engine_test.go` — test all 4 bonus combinations (table-driven)

### Task 11.2 — Integration Tests (repository layer)
Use `testcontainers-go` with `postgis/postgis:16-3.4`.

- [ ] `user_repo_integration_test.go`
- [ ] `mission_repo_integration_test.go`
- [ ] `expense_repo_integration_test.go` — test FindOverdue with time manipulation
- [ ] `radar_repo_integration_test.go` — insert profiles with real ST_MakePoint, assert ST_DWithin results

**Output files**: 15 test files

---

## Phase 12 — Final Verification

- [ ] `make docker-up` — postgres+postgis starts clean
- [ ] `make migrate-up` — all 20 migrations apply without error
- [ ] `make build` — compiles with zero warnings
- [ ] `make test` — all unit tests pass
- [ ] `make lint` — zero `golangci-lint` errors (use `errcheck`, `govet`, `staticcheck`)
- [ ] Manually verify Output Checklist in `SKILL.md` — every item ticked

---

## Dependency Graph (build order enforced above)

```
SKILL.md + seq_*.mmd
        │
        ▼
Phase 1: Scaffold (go.mod, docker, Makefile)
        │
        ▼
Phase 2: Infrastructure (config, db, pkg/)
        │
        ▼
Phase 3: Domain (structs, enums, business methods)
        │
        ▼
Phase 4: Ports (interfaces)
        │
        ├──────────────────────┐
        ▼                      ▼
Phase 5: Migrations      Phase 7: External Adapters
        │                      │
        ▼                      │
Phase 6: Repo Impls ◄──────────┘
        │
        ▼
Phase 8: Usecases
        │
        ▼
Phase 9: HTTP Layer
        │
        ▼
Phase 10: Scheduler + main.go
        │
        ▼
Phase 11: Tests
        │
        ▼
Phase 12: Verify
```

---

## Agent Notes

- **Never skip a task** — each phase builds on the previous. Writing handlers before usecases will produce uncompilable code.
- **Never use `interface{}`** — all types must be concrete domain types or typed interfaces.
- **Never call `time.Now()` directly** — always use `clock.Now()` injected via constructor.
- **Multi-table writes = transaction** — if you find yourself calling two repo methods sequentially without a `pgx.Tx`, it is a bug.
- **After each phase**, do a quick compilation check (`go build ./...`) before proceeding.
- **Sequence diagrams are the source of truth** for business logic. If `SKILL.md` and a sequence diagram conflict, the sequence diagram wins.
