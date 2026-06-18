# J1 Hub Backend — Go Implementation Skill

## Overview

J1 Hub is a mobile-first platform for Thai J1 visa students working in the USA. It guides users through compliance missions, gamifies progress with points and badges, enables expense splitting with a credit score system, and surfaces J1 job listings with community-sourced reviews.

This skill instructs an AI agent to implement the **complete Go backend** from a set of Mermaid sequence diagrams and an ER diagram (v5). The agent must produce production-ready code following **Clean Architecture**, **OOP via interfaces**, and **Go best practices**.

---

## Input Artifacts

The agent MUST read these files before writing any code:

```
er_diagram_v5.mmd          – full data model (all tables, fields, types, enums, constraints)
seq_01_user_registration.mmd
seq_02_mission_complete_upload.mmd
seq_03_phase_transition.mmd
seq_04_friendship.mmd
seq_05_expense_split.mmd
seq_06_expense_overdue_penalty.mmd
seq_07_job_browse_cart.mmd
seq_08_radar_location.mmd
seq_09_mission_overdue.mmd
seq_10_leaderboard.mmd
```

Each sequence diagram defines one complete user flow. Implement every step shown — every DB write, every state transition, every notification trigger — exactly as drawn.

---

## Tech Stack

| Concern | Choice |
|---|---|
| Language | Go 1.22+ |
| HTTP router | `github.com/go-chi/chi/v5` |
| Database | PostgreSQL 16 with PostGIS (for `geometry_point`) |
| DB driver | `github.com/jackc/pgx/v5` + `github.com/jackc/pgx/v5/pgxpool` |
| Migrations | `github.com/golang-migrate/migrate/v4` |
| Password hash | `golang.org/x/crypto/argon2` |
| JWT | `github.com/golang-jwt/jwt/v5` |
| File storage | Supabase Storage via REST (multipart upload) |
<!-- | Push notifications | FCM via `firebase.google.com/go/v4` | -->
| Scheduler (cron) | `github.com/robfig/cron/v3` |
| Config | `github.com/spf13/viper` + `.env` |
| Validation | `github.com/go-playground/validator/v10` |
| Logging | `log/slog` (structured, JSON in prod) |
| Testing | `testing` + `github.com/stretchr/testify` + `github.com/testcontainers/testcontainers-go` |

---

## Project Layout (Clean Architecture)

```
j1hub-backend/
├── cmd/
│   └── server/
│       └── main.go                  # wire everything, start HTTP + cron
├── internal/
│   ├── domain/                      # pure business types — NO framework imports
│   │   ├── user.go
│   │   ├── mission.go
│   │   ├── expense.go
│   │   ├── job.go
│   │   ├── friendship.go
│   │   ├── gamification.go
│   │   └── errors.go                # sentinel domain errors
│   ├── port/                        # interfaces (input & output ports)
│   │   ├── repository.go            # all Repository interfaces
│   │   └── service.go               # all Service interfaces
│   ├── usecase/                     # one file per sequence diagram flow
│   │   ├── register_user.go         # seq_01
│   │   ├── complete_mission.go      # seq_02
│   │   ├── advance_phase.go         # seq_03
│   │   ├── manage_friendship.go     # seq_04
│   │   ├── manage_expense.go        # seq_05
│   │   ├── overdue_expense_job.go   # seq_06 (cron)
│   │   ├── manage_job.go            # seq_07
│   │   ├── radar.go                 # seq_08
│   │   ├── overdue_mission_job.go   # seq_09 (cron)
│   │   └── leaderboard.go           # seq_10
│   ├── adapter/
│   │   ├── http/
│   │   │   ├── router.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go          # JWT validation
│   │   │   │   └── logger.go
│   │   │   └── handler/
│   │   │       ├── auth_handler.go
│   │   │       ├── user_handler.go
│   │   │       ├── mission_handler.go
│   │   │       ├── expense_handler.go
│   │   │       ├── job_handler.go
│   │   │       ├── friendship_handler.go
│   │   │       ├── radar_handler.go
│   │   │       ├── leaderboard_handler.go
│   │   │       └── admin_handler.go
│   │   ├── postgres/                # Repository implementations
│   │   │   ├── user_repo.go
│   │   │   ├── mission_repo.go
│   │   │   ├── expense_repo.go
│   │   │   ├── job_repo.go
│   │   │   ├── friendship_repo.go
│   │   │   ├── gamification_repo.go
│   │   │   └── radar_repo.go
│   │   ├── storage/
│   │   │   └── supabase_storage.go  # file upload adapter
│   │   └── notification/
│   │       └── fcm_notifier.go      # push notification adapter
│   └── infrastructure/
│       ├── config/
│       │   └── config.go
│       ├── db/
│       │   └── postgres.go          # pgxpool init, ping
│       └── scheduler/
│           └── cron.go              # register seq_06 & seq_09 jobs
├── migrations/
│   ├── 000001_create_users.up.sql
│   ├── 000002_create_profiles.up.sql
│   └── ...                          # one file per table in ER diagram order
├── pkg/
│   ├── apperror/                    # HTTP-friendly error wrapper
│   ├── uid/                         # ID generator (prefix + ulid: "usr_", "mis_", etc.)
│   └── timeutil/
│       └── now.go                   # injectable clock for testing
├── .env.example
├── docker-compose.yml               # postgres+postgis, app
├── Makefile
└── go.mod
```

---

## Architecture Rules

### Domain Layer (`internal/domain/`)
- Plain Go structs only. No DB tags, no JSON tags, no framework imports.
- Enums as typed `string` constants with a `Valid()` method.
- Business rules live here: `CanAdvancePhase()`, `CalculateDueDate()`, `CalculatePointReward()`.

```go
// Example
type UserMissionStatus string

const (
    StatusNotStarted          UserMissionStatus = "Not_Started"
    StatusInProgress          UserMissionStatus = "In_Progress"
    StatusPendingVerification UserMissionStatus = "Pending_Verification"
    StatusCompleted           UserMissionStatus = "Completed"
    StatusOverdue             UserMissionStatus = "Overdue"
)

func (s UserMissionStatus) Valid() bool { ... }
```

### Port Layer (`internal/port/`)
- One interface per repository, one interface per service.
- All repository methods accept `context.Context` as the first argument.
- Interfaces are defined in the `port` package and consumed by usecases — dependency inversion strictly enforced.

```go
type UserRepository interface {
    Create(ctx context.Context, u *domain.User) error
    FindByID(ctx context.Context, id string) (*domain.User, error)
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
    Update(ctx context.Context, u *domain.User) error
    IncrementPoints(ctx context.Context, userID string, lifetimeDelta, phaseDelta int) error
    ResetStreak(ctx context.Context, userID string) error
}

type StoragePort interface {
    UploadFile(ctx context.Context, bucket, key string, data io.Reader, contentType string) (url string, err error)
}

type NotifierPort interface {
    Send(ctx context.Context, userID, title, body string) error
}
```

### Usecase Layer (`internal/usecase/`)
- Structs with injected port interfaces — never concrete types.
- One public method per sequence diagram step-group (not one method per DB call).
- All DB mutations that span multiple tables MUST use `pgx` transactions.
- Return domain errors (from `internal/domain/errors.go`), not HTTP status codes.

```go
type RegisterUserUseCase struct {
    userRepo    port.UserRepository
    profileRepo port.ProfileRepository
    creditRepo  port.CreditScoreRepository
    phaseRepo   port.JourneyPhaseRepository
    missionRepo port.MissionRepository
    hasher      port.PasswordHasher
    notifier    port.NotifierPort
    clock       timeutil.Clock
}

func (uc *RegisterUserUseCase) Register(ctx context.Context, cmd RegisterCommand) (*domain.User, string, error) { ... }
func (uc *RegisterUserUseCase) InitializeJourney(ctx context.Context, userID string, cmd InitJourneyCommand) error { ... }
```

### Adapter / HTTP Layer (`internal/adapter/http/`)
- Handlers only: decode request → call usecase → encode response.
- Map domain errors to HTTP status codes in a single `respondError(w, err)` helper.
- All request bodies validated with `validator` tags before reaching usecases.

### Repository Layer (`internal/adapter/postgres/`)
- Implement port interfaces using raw `pgx` queries (no ORM).
- Struct fields mapped manually from `pgx.Row` / `pgx.Rows`.
- PostGIS coordinates stored as `geometry(Point, 4326)` and queried with `ST_DWithin`, `ST_AsText`.

---

## ID Generation

All primary keys use a prefixed ULID:

```go
// pkg/uid/uid.go
func New(prefix string) string {
    return prefix + ulid.Make().String()  // e.g. "usr_01ARZ3NDEKTSV4RRFFQ69G5FAV"
}
// prefixes: usr_ prf_ frn_ phs_ uph_ mis_ ums_ tsk_ utk_
//           ldg_ bdg_ ubdg_ crd_ txn_ spl_ job_ hsg_ crt_ smr_ rvw_
```

---

## Reward Engine (seq_02)

Implement as a standalone usecase struct `RewardEngine` with method:

```go
func (re *RewardEngine) Calculate(ctx context.Context, um *domain.UserMission, user *domain.User) (*domain.PointReward, error)
```

Rules (all configurable via Viper):
- **Base**: `mission.base_points`
- **Speed bonus**: if `proof_submitted_at` < `calculated_due_date - 7 days` → `+20%` of base; if < `due_date - 1 day` → `+10%`
- **Streak bonus**: streak 3–6 → `+10%`; 7+ → `+25%` of base
- **First completer**: query `USER_MISSION` count where `mission_id=X AND status=Completed` — if 0 rows before this one → `+200 flat bonus`
- Total is sum of all four. Write one `POINT_LEDGER` row per bonus source type.

---

## Cron Jobs

### seq_06 — Overdue Expense (daily 00:00 UTC)
1. Find all `EXPENSE_SPLIT` where `payment_status IN ('Pending','Submitted')` AND `EXPENSE_TRANSACTION.due_date < NOW()`.
2. Batch update `payment_status = 'Overdue'`.
3. For each: decrement `CREDIT_SCORE.current_score` by 10.
4. Insert `POINT_LEDGER` row (source_type=`Expense_Penalty`, delta=-10).
5. Send push notification to debtor.

### seq_09 — Overdue Mission (daily 00:00 UTC)
1. Find all `USER_MISSION` where `status IN ('Not_Started','In_Progress','Pending_Verification')` AND `calculated_due_date < NOW()`.
2. Batch update `status = 'Overdue'`.
3. For each where `MISSION.is_mandatory = true`: reset `USER.mission_streak = 0`, insert `POINT_LEDGER` audit row.
4. Send push notification.

Both jobs must be idempotent (safe to re-run). Use DB-level locks or `WHERE status != 'Overdue'` guards.

---

## PostGIS Radar (seq_08)

```sql
-- location update
UPDATE profile SET current_coordinates = ST_SetSRID(ST_MakePoint($1, $2), 4326),
                   location_updated_at = NOW()
WHERE user_id = $3;

-- nearby query (radius in meters, default 5000)
SELECT p.user_id, p.radar_visibility, p.avatar_url,
       u.first_name, u.last_name
FROM profile p
JOIN users u ON u.user_id = p.user_id
WHERE ST_DWithin(
    p.current_coordinates::geography,
    ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
    $3
)
AND p.location_updated_at > NOW() - INTERVAL '30 minutes'
AND p.user_id != $4;
```

Visibility masking applied in the usecase layer, NOT in SQL.

---

## API Routes

```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
PATCH  /api/v1/users/me
GET    /api/v1/users/me

GET    /api/v1/phases/current
GET    /api/v1/missions
PATCH  /api/v1/user-missions/:id/tasks/:task_id
POST   /api/v1/user-missions/:id/proof

GET    /api/v1/friends
POST   /api/v1/friends/request
PATCH  /api/v1/friends/:friendship_id

POST   /api/v1/expenses
POST   /api/v1/expenses/:id/splits/:split_id/slip
PATCH  /api/v1/expenses/:id/splits/:split_id/approve

GET    /api/v1/jobs
GET    /api/v1/jobs/:id
POST   /api/v1/cart
PATCH  /api/v1/cart/:cart_id
POST   /api/v1/jobs/:id/reviews

PATCH  /api/v1/profile/location
GET    /api/v1/radar

GET    /api/v1/leaderboard

# Admin (requires admin JWT claim)
PATCH  /api/v1/admin/user-missions/:id/verify
```

---

## Error Handling

Define sentinel errors in `internal/domain/errors.go`:

```go
var (
    ErrNotFound         = errors.New("not found")
    ErrUnauthorized     = errors.New("unauthorized")
    ErrConflict         = errors.New("conflict")
    ErrInvalidInput     = errors.New("invalid input")
    ErrAlreadyCompleted = errors.New("mission already completed")
    ErrSelfSplit        = errors.New("cannot split expense with yourself")
    ErrDuplicateFriend  = errors.New("friendship already exists")
)
```

HTTP handler maps these with `errors.Is()` to 400/401/403/404/409.

---

## Migration Files

Create one `.up.sql` + `.down.sql` per table, in dependency order matching the ER diagram:

```
000001_users.up.sql
000002_profiles.up.sql
000003_friendships.up.sql
000004_journey_phases.up.sql
000005_user_phase_history.up.sql
000006_missions.up.sql
000007_user_missions.up.sql
000008_tasks.up.sql
000009_user_tasks.up.sql
000010_point_ledger.up.sql
000011_badges.up.sql
000012_user_badges.up.sql
000013_credit_scores.up.sql
000014_expense_transactions.up.sql
000015_expense_splits.up.sql
000016_job_postings.up.sql
000017_job_housings.up.sql
000018_user_carts.up.sql
000019_job_overall_ratings.up.sql
000020_job_reviews.up.sql
```

All tables must include created_at/updated_at with `DEFAULT NOW()`. Enums implemented as `TEXT` with `CHECK` constraints. PostGIS extension enabled in migration 000001.

---

## Testing Requirements

- **Unit tests** for all usecase structs using mock repositories (implement port interfaces with `testify/mock`).
- **Integration tests** for all repository implementations using `testcontainers-go` with a real PostgreSQL+PostGIS container.
- **At minimum** test each sequence diagram flow end-to-end at the usecase level.
- Use `timeutil.Clock` interface to control time in tests (avoid `time.Now()` directly).

```go
type Clock interface {
    Now() time.Time
}
type RealClock struct{}
func (RealClock) Now() time.Time { return time.Now() }
```

---

## Output Checklist

The agent must produce:

- [ ] `go.mod` with all dependencies
- [ ] All domain structs with enum types and business rule methods
- [ ] All port interfaces
- [ ] All 10 usecase implementations (one per sequence diagram)
- [ ] All repository implementations with raw pgx queries
- [ ] HTTP router with all routes wired
- [ ] All HTTP handlers with request/response structs and validation
- [ ] JWT auth middleware
- [ ] Supabase Storage adapter
- [ ] FCM push notification adapter
- [ ] Cron scheduler setup (seq_06 + seq_09)
- [ ] Reward engine with all bonus types
- [ ] All 20 migration files (up + down)
- [ ] `docker-compose.yml` with postgres+postgis
- [ ] `Makefile` with targets: `run`, `migrate-up`, `migrate-down`, `test`, `lint`
- [ ] `.env.example`
- [ ] Unit tests for all usecases
- [ ] Integration tests for all repositories

---

## Key Constraints

1. **No ORM** — use raw `pgx` queries everywhere.
2. **No global state** — all dependencies injected via constructors.
3. **Transactions** — any usecase that writes to 2+ tables must use a `pgx.Tx`.
4. **Points are only awarded after `verified_at` is set** — never on proof submission alone.
5. **Friendship canonical order** — always sort `(user_id_1, user_id_2)` lexicographically before insert/query.
6. **Phase points reset to 0 on phase transition; lifetime points never reset.**
7. **Cron jobs must be idempotent** — guard with status checks before each mutation.
8. **Radar masking happens in usecase layer** — SQL returns raw data, business logic applies visibility rules in Go.
