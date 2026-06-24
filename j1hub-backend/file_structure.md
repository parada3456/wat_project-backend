# j1hub-backend Directory Structure

This document shows the complete directory structure of the `j1hub-backend` Go service.

```text
.
├── .env.example
├── Makefile
├── clean_architecture_refactor_plan.md
├── cmd/
│   ├── mock_server.go
│   ├── scraper/
│   │   ├── main.go
│   │   └── main_test.go
│   └── server/
│       └── main.go
├── configs/
│   └── fcm-credentials.json
├── docker-compose.yml
├── file_structure.md
├── go.mod
├── go.sum
├── internal/
│   ├── adapter/
│   │   ├── auth/
│   │   │   ├── argon2_hasher.go
│   │   │   ├── argon2_hasher_test.go
│   │   │   ├── jwt_issuer.go
│   │   │   └── jwt_issuer_test.go
│   │   ├── http/
│   │   │   ├── handler/
│   │   │   │   ├── dto/
│   │   │   │   │   ├── admin_request.go
│   │   │   │   │   ├── admin_response.go
│   │   │   │   │   ├── auth_request.go
│   │   │   │   │   ├── auth_response.go
│   │   │   │   │   ├── expense_request.go
│   │   │   │   │   ├── expense_response.go
│   │   │   │   │   ├── friend_request.go
│   │   │   │   │   ├── friend_response.go
│   │   │   │   │   ├── job_request.go
│   │   │   │   │   ├── job_response.go
│   │   │   │   │   ├── journey_response.go
│   │   │   │   │   ├── mission_request.go
│   │   │   │   │   ├── mission_response.go
│   │   │   │   │   ├── user_request.go
│   │   │   │   │   └── user_response.go
│   │   │   │   ├── job_handler.go
│   │   │   │   ├── router.go
│   │   │   │   └── test/
│   │   │   │       ├── auth_handler_test.go
│   │   │   │       ├── expense_handler_test.go
│   │   │   │       ├── friend_handler_test.go
│   │   │   │       ├── job_handler_test.go
│   │   │   │       ├── journey_handler_test.go
│   │   │   │       ├── mission_handler_test.go
│   │   │   │       ├── mocks_test.go
│   │   │   │       ├── notification_handler_test.go
│   │   │   │       ├── router_test.go
│   │   │   │       └── user_handler_test.go
│   │   │   └── middleware/
│   │   │       ├── auth.go
│   │   │       ├── cors.go
│   │   │       ├── logger.go
│   │   │       └── middleware_test.go
│   │   ├── notification/
│   │   │   ├── fcm_notifier.go
│   │   │   └── fcm_notifier_test.go
│   │   ├── outbound/
│   │   │   └── scraper/
│   │   │       ├── acadex/
│   │   │       │   ├── acadex.go
│   │   │       │   └── acadex_test.go
│   │   │       ├── iee/
│   │   │       │   ├── iee.go
│   │   │       │   └── iee_test.go
│   │   │       ├── ihappy/
│   │   │       │   ├── ihappy.go
│   │   │       │   └── ihappy_test.go
│   │   │       └── scraper.go
│   │   ├── postgres/
│   │   │   ├── admin_repo.go
│   │   │   ├── expense_repo.go
│   │   │   ├── friendship_repo.go
│   │   │   ├── gamification_repo.go
│   │   │   ├── job_repo.go
│   │   │   ├── leaderboard_repo.go
│   │   │   ├── mission_repo.go
│   │   │   ├── notification_repo.go
│   │   │   ├── postgres_test.go
│   │   │   ├── profile_repo.go
│   │   │   ├── radar_repo.go
│   │   │   └── user_repo.go
│   │   └── storage/
│   │       ├── supabase_storage.go
│   │       └── supabase_storage_test.go
│   ├── admin/
│   │   ├── adapter/
│   │   │   ├── http/
│   │   │   │   └── handler/
│   │   │   │       └── admin_handler.go
│   │   │   └── postgres/
│   │   ├── domain/
│   │   ├── port/
│   │   └── usecase/
│   │       ├── admin_usecase.go
│   │       └── mocks_test.go
│   ├── auth/
│   │   ├── adapter/
│   │   │   ├── http/
│   │   │   │   └── handler/
│   │   │   │       └── auth_handler.go
│   │   │   └── postgres/
│   │   ├── domain/
│   │   ├── port/
│   │   └── usecase/
│   │       ├── login_user.go
│   │       ├── login_user_test.go
│   │       ├── mocks_test.go
│   │       ├── register_user.go
│   │       └── register_user_test.go
│   ├── domain/
│   │   ├── domain_test.go
│   │   └── errors.go
│   ├── expense/
│   │   ├── adapter/
│   │   │   ├── http/
│   │   │   │   └── handler/
│   │   │   │       └── expense_handler.go
│   │   │   └── postgres/
│   │   ├── domain/
│   │   │   └── expense.go
│   │   ├── port/
│   │   └── usecase/
│   │       ├── manage_expense.go
│   │       └── manage_expense_test.go
│   ├── friend/
│   │   ├── adapter/
│   │   │   ├── http/
│   │   │   │   └── handler/
│   │   │   │       └── friend_handler.go
│   │   │   └── postgres/
│   │   ├── domain/
│   │   │   └── friendship.go
│   │   ├── port/
│   │   └── usecase/
│   │       ├── manage_friendship.go
│   │       ├── manage_friendship_test.go
│   │       └── mocks_test.go
│   ├── gamification/
│   │   ├── adapter/
│   │   │   ├── http/
│   │   │   │   └── handler/
│   │   │   │       └── journey_handler.go
│   │   │   └── postgres/
│   │   ├── domain/
│   │   │   └── gamification.go
│   │   ├── port/
│   │   └── usecase/
│   │       ├── leaderboard.go
│   │       ├── leaderboard_test.go
│   │       ├── mocks_test.go
│   │       ├── radar.go
│   │       ├── radar_test.go
│   │       ├── reward_engine.go
│   │       └── reward_engine_test.go
│   ├── infrastructure/
│   │   ├── config/
│   │   │   ├── config.go
│   │   │   └── config_test.go
│   │   ├── db/
│   │   │   ├── postgres.go
│   │   │   └── postgres_test.go
│   │   └── scheduler/
│   │       ├── cron.go
│   │       └── cron_test.go
│   ├── job/
│   │   ├── domain/
│   │   │   └── job.go
│   │   └── usecase/
│   │       ├── manage_job.go
│   │       ├── manage_job_test.go
│   │       └── mocks_test.go
│   ├── mission/
│   │   ├── adapter/
│   │   │   ├── http/
│   │   │   │   └── handler/
│   │   │   │       └── mission_handler.go
│   │   │   └── postgres/
│   │   ├── domain/
│   │   │   └── mission.go
│   │   ├── port/
│   │   └── usecase/
│   │       ├── complete_mission.go
│   │       ├── complete_mission_test.go
│   │       ├── manage_mission.go
│   │       ├── manage_mission_test.go
│   │       └── mocks_test.go
│   ├── notification/
│   │   ├── adapter/
│   │   │   ├── http/
│   │   │   │   └── handler/
│   │   │   │       └── notification_handler.go
│   │   │   └── postgres/
│   │   ├── domain/
│   │   │   └── notification.go
│   │   ├── port/
│   │   └── usecase/
│   │       ├── manage_notification.go
│   │       ├── manage_notification_test.go
│   │       └── mocks_test.go
│   ├── port/
│   │   ├── repository.go
│   │   └── service.go
│   ├── usecase/
│   │   ├── advance_phase.go
│   │   ├── advance_phase_test.go
│   │   ├── dto.go
│   │   ├── manage_journey.go
│   │   ├── manage_journey_test.go
│   │   ├── mocks_test.go
│   │   ├── overdue_expense_job.go
│   │   ├── overdue_jobs_test.go
│   │   ├── overdue_mission_job.go
│   │   ├── scrape_jobs.go
│   │   └── scrape_jobs_test.go
│   └── user/
│       ├── adapter/
│       │   ├── http/
│       │   │   └── handler/
│       │   │       └── user_handler.go
│       │   └── postgres/
│       ├── domain/
│       │   └── user.go
│       ├── port/
│       └── usecase/
│           ├── manage_user.go
│           ├── manage_user_test.go
│           └── mocks_test.go
├── job_housing.json
├── job_posting.json
├── migrations/
│   ├── 000001_users.down.sql
│   ├── 000001_users.up.sql
│   ├── 000002_profiles.down.sql
│   ├── 000002_profiles.up.sql
│   ├── 000003_friendships.down.sql
│   ├── 000003_friendships.up.sql
│   ├── 000004_journey_phases.down.sql
│   ├── 000004_journey_phases.up.sql
│   ├── 000005_user_phase_history.down.sql
│   ├── 000005_user_phase_history.up.sql
│   ├── 000006_missions.down.sql
│   ├── 000006_missions.up.sql
│   ├── 000007_user_missions.down.sql
│   ├── 000007_user_missions.up.sql
│   ├── 000008_tasks.down.sql
│   ├── 000008_tasks.up.sql
│   ├── 000009_user_tasks.down.sql
│   ├── 000009_user_tasks.up.sql
│   ├── 000010_point_ledger.down.sql
│   ├── 000010_point_ledger.up.sql
│   ├── 000011_badges.down.sql
│   ├── 000011_badges.up.sql
│   ├── 000012_user_badges.down.sql
│   ├── 000012_user_badges.up.sql
│   ├── 000013_credit_scores.down.sql
│   ├── 000013_credit_scores.up.sql
│   ├── 000014_expense_transactions.down.sql
│   ├── 000014_expense_transactions.up.sql
│   ├── 000015_expense_splits.down.sql
│   ├── 000015_expense_splits.up.sql
│   ├── 000016_job_postings.down.sql
│   ├── 000016_job_postings.up.sql
│   ├── 000017_job_housings.down.sql
│   ├── 000017_job_housings.up.sql
│   ├── 000018_user_carts.down.sql
│   ├── 000018_user_carts.up.sql
│   ├── 000019_job_overall_ratings.down.sql
│   ├── 000019_job_overall_ratings.up.sql
│   ├── 000020_job_reviews.down.sql
│   └── 000020_job_reviews.up.sql
└── pkg/
    ├── apperror/
    │   ├── apperror.go
    │   └── apperror_test.go
    ├── response/
    │   └── response.go
    ├── timeutil/
    │   ├── clock.go
    │   └── clock_test.go
    └── uid/
        ├── uid.go
        └── uid_test.go
```

## Key Directories and Architecture
- **`cmd/`**: Entry points for compiling the binaries.
  - `server/`: The main HTTP REST server.
  - `scraper/`: The worker process for scraping jobs and housing info.
- **`configs/`**: Application configuration files (e.g., Firebase Cloud Messaging credentials).
- **`migrations/`**: Postgres database migration SQL files schema changes.
- **`pkg/`**: Utility packages shared across the application (errors, HTTP responses, UUID helpers, custom clocks).
- **`internal/`**: Core Clean Architecture business logic separated by domains (auth, user, job, expense, friend, gamification, mission, notification, admin).
  - **`domain/`**: Enterprise-level entities and business rules (innermost layer).
  - **`usecase/`**: Application-specific business rules/interactors.
  - **`port/`**: Interfaces defining input ports (interfaces called by controllers/handlers) and output ports (repository interfaces).
  - **`adapter/`**: Secondary adapters implementing the port interfaces (e.g., HTTP handlers under `adapter/http`, database repositories under `adapter/postgres`).
