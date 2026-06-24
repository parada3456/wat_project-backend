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
├── migrations/
│   ├── 000001_users.down.sql
│   ├── 000001_users.up.sql
│   └── [Remaining 000002–000020 migration SQL files]
├── pkg/
│   ├── apperror/
│   │   ├── apperror.go
│   │   └── apperror_test.go
│   ├── response/
│   │   └── response.go
│   ├── timeutil/
│   │   ├── clock.go
│   │   └── clock_test.go
│   └── uid/
│       ├── uid.go
│       └── uid_test.go
├── job_housing.json
├── job_posting.json
└── internal/
    ├── admin/
    │   ├── adapter/
    │   │   ├── http/
    │   │   │   ├── dto/
    │   │   │   │   ├── admin_request.go
    │   │   │   │   └── admin_response.go
    │   │   │   └── admin_handler.go
    │   │   └── postgres/
    │   │       └── admin_repo.go
    │   ├── domain/
    │   ├── port/
    │   └── usecase/
    │       ├── admin_usecase.go
    │       └── mocks_test.go
    │
    ├── auth/
    │   ├── adapter/
    │   │   ├── http/
    │   │   │   ├── dto/
    │   │   │   │   ├── auth_request.go
    │   │   │   │   └── auth_response.go
    │   │   │   └── auth_handler.go
    │   │   └── postgres/
    │   │       └── auth_repo.go
    │   ├── domain/
    │   ├── port/
    │   └── usecase/
    │       ├── login_user.go
    │       ├── login_user_test.go
    │       ├── mocks_test.go
    │       ├── register_user.go
    │       └── register_user_test.go
    │
    ├── domain/               # Global/Shared Core Domain Errors only
    │   ├── domain_test.go
    │   └── errors.go
    │
    ├── expense/
    │   ├── adapter/
    │   │   ├── http/
    │   │   │   ├── dto/
    │   │   │   │   ├── expense_request.go
    │   │   │   │   └── expense_response.go
    │   │   │   └── expense_handler.go
    │   │   └── postgres/
    │   │       └── expense_repo.go
    │   ├── domain/
    │   │   └── expense.go
    │   ├── port/
    │   └── usecase/
    │       ├── manage_expense.go
    │       └── manage_expense_test.go
    │
    ├── friend/
    │   ├── adapter/
    │   │   ├── http/
    │   │   │   ├── dto/
    │   │   │   │   ├── friend_request.go
    │   │   │   │   └── friend_response.go
    │   │   │   └── friend_handler.go
    │   │   └── postgres/
    │   │       └── friendship_repo.go
    │   ├── domain/
    │   │   └── friendship.go
    │   ├── port/
    │   └── usecase/
    │       ├── manage_friendship.go
    │       ├── manage_friendship_test.go
    │       └── mocks_test.go
    │
    ├── gamification/
    │   ├── adapter/
    │   │   ├── http/
    │   │   │   ├── dto/
    │   │   │   │   └── journey_response.go
    │   │   │   └── journey_handler.go
    │   │   └── postgres/
    │   │       ├── gamification_repo.go
    │   │       ├── leaderboard_repo.go
    │   │       └── radar_repo.go
    │   ├── domain/
    │   │   └── gamification.go
    │   ├── port/
    │   └── usecase/
    │       ├── advance_phase.go
    │       ├── advance_phase_test.go
    │       ├── leaderboard.go
    │       ├── leaderboard_test.go
    │       ├── manage_journey.go
    │       ├── manage_journey_test.go
    │       ├── mocks_test.go
    │       ├── radar.go
    │       ├── radar_test.go
    │       ├── reward_engine.go
    │       └── reward_engine_test.go
    │
    ├── job/
    │   ├── adapter/
    │   │   ├── http/
    │   │   │   ├── dto/
    │   │   │   │   ├── job_request.go
    │   │   │   │   └── job_response.go
    │   │   │   └── job_handler.go
    │   │   └── postgres/
    │   │       └── job_repo.go
    │   ├── domain/
    │   │   └── job.go
    │   ├── port/
    │   └── usecase/
    │       ├── manage_job.go
    │       ├── manage_job_test.go
    │       └── mocks_test.go
    │
    ├── mission/
    │   ├── adapter/
    │   │   ├── http/
    │   │   │   ├── dto/
    │   │   │   │   ├── mission_request.go
    │   │   │   │   └── mission_response.go
    │   │   │   └── mission_handler.go
    │   │   └── postgres/
    │   │       └── mission_repo.go
    │   ├── domain/
    │   │   └── mission.go
    │   ├── port/
    │   └── usecase/
    │       ├── complete_mission.go
    │       ├── complete_mission_test.go
    │       ├── manage_mission.go
    │       ├── manage_mission_test.go
    │       └── mocks_test.go
    │
    ├── notification/
    │   ├── adapter/
    │   │   ├── http/
    │   │   │   └── notification_handler.go
    │   │   └── postgres/
    │   │       └── notification_repo.go
    │   ├── domain/
    │   │   └── notification.go
    │   ├── port/
    │   └── usecase/
    │       ├── manage_notification.go
    │       ├── manage_notification_test.go
    │       └── mocks_test.go
    │
    ├── user/
    │   ├── adapter/
    │   │   ├── http/
    │   │   │   ├── dto/
    │   │   │   │   ├── user_request.go
    │   │   │   │   └── user_response.go
    │   │   │   └── user_handler.go
    │   │   └── postgres/
    │   │       ├── profile_repo.go
    │   │       └── user_repo.go
    │   ├── domain/
    │   │   └── user.go
    │   ├── port/
    │   └── usecase/
    │       ├── manage_user.go
    │       ├── manage_user_test.go
    │       └── mocks_test.go
    │
    ├── infrastructure/       # Cross-cutting tech adapters/engines (Secondary Adapters)
    │   ├── config/
    │   │   ├── config.go
    │   │   └── config_test.go
    │   ├── db/
    │   │   ├── postgres.go
    │   │   └── postgres_test.go
    │   ├── security/         # Moved from internal/adapter/auth
    │   │   ├── argon2_hasher.go
    │   │   ├── argon2_hasher_test.go
    │   │   ├── jwt_issuer.go
    │   │   └── jwt_issuer_test.go
    │   ├── notification/     # Base client integrations
    │   │   ├── fcm_notifier.go
    │   │   └── fcm_notifier_test.go
    │   ├── storage/          # External object storage clients
    │   │   ├── supabase_storage.go
    │   │   └── supabase_storage_test.go
    │   ├── scheduler/        # Internal async background engines
    │   │   ├── cron.go
    │   │   ├── cron_test.go
    │   │   ├── overdue_expense_job.go
    │   │   ├── overdue_jobs_test.go
    │   │   └── overdue_mission_job.go
    │   └── outbound/         # External scrapers and integrations
    │       └── scraper/
    │           ├── acadex/
    │           │   ├── acadex.go
    │           │   └── acadex_test.go
    │           ├── iee/
    │           │   ├── iee.go
    │           │   └── iee_test.go
    │           ├── ihappy/
    │           │   ├── ihappy.go
    │           │   └── ihappy_test.go
    │           ├── scraper.go
    │           ├── scrape_jobs.go
    │           └── scrape_jobs_test.go
    │
    └── transport/            # HTTP Entry point and routing (Primary Adapter)
        └── http/
            ├── middleware/
            │   ├── auth.go
            │   ├── cors.go
            │   ├── logger.go
            │   └── middleware_test.go
            ├── test/         # End-to-End/Integration Router tests
            │   ├── auth_handler_test.go
            │   ├── expense_handler_test.go
            │   ├── friend_handler_test.go
            │   ├── job_handler_test.go
            │   ├── journey_handler_test.go
            │   ├── mission_handler_test.go
            │   ├── mocks_test.go
            │   ├── notification_handler_test.go
            │   ├── router_test.go
            │   └── user_handler_test.go
            └── router.go     # Wireframes everything together