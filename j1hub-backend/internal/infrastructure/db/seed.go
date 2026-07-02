package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedMockData(pool *pgxpool.Pool) error {
	ctx := context.Background()

	// Check if already seeded
	var exists bool
	err := pool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM journey_phases LIMIT 1)").Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if db is seeded: %w", err)
	}
	if exists {
		log.Println("Database already contains data. Skipping mock data seeding.")
		return nil
	}

	log.Println("Database is empty. Seeding mock data...")

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction for seeding: %w", err)
	}
	defer tx.Rollback(ctx)

	queries := []string{
		// 1. Journey Phases
		`INSERT INTO journey_phases (phase_id, phase_number, title, description) VALUES
		('phs_001', 1, 'Pre-Departure', 'Tasks to complete before leaving your home country.'),
		('phs_002', 2, 'First Week Arrival', 'Tasks for your first week in the United States.'),
		('phs_003', 3, 'Mid-Term J1 Experience', 'Tasks during your active working period.'),
		('phs_004', 4, 'Return & Graduation', 'Final program checkouts and tax filings.');`,

		// 2. Users (Password is "password123")
		`INSERT INTO users (user_id, email, password_hash, current_phase_id, total_lifetime_points, current_phase_points, mission_streak, arrival_date, job_start_date) VALUES
		('usr_admin', 'admin@j1hub.com', '$argon2id$v=19$m=65536,t=1,p=4$9o8p3gSSPWq2GK6884/oXQ$Y6y1Bry9hDlzuaZ9S/StpX4ddQqDwSyvQ/iATKuin1o', 'phs_001', 0, 0, 0, NOW(), NOW()),
		('usr_somchai', 'somchai.j1@kmitl.ac.th', '$argon2id$v=19$m=65536,t=1,p=4$9o8p3gSSPWq2GK6884/oXQ$Y6y1Bry9hDlzuaZ9S/StpX4ddQqDwSyvQ/iATKuin1o', 'phs_002', 270, 270, 3, NOW() - INTERVAL '30 days', NOW() - INTERVAL '25 days'),
		('usr_alice', 'alice@example.com', '$argon2id$v=19$m=65536,t=1,p=4$9o8p3gSSPWq2GK6884/oXQ$Y6y1Bry9hDlzuaZ9S/StpX4ddQqDwSyvQ/iATKuin1o', 'phs_003', 460, 150, 4, NOW() - INTERVAL '60 days', NOW() - INTERVAL '55 days');`,

		// 3. Profiles (using PostGIS points)
		`INSERT INTO profiles (profile_id, user_id, username, first_name, last_name, phone_number, bio, avatar_url, radar_visibility, current_coordinates, location_updated_at) VALUES
		('prf_admin', 'usr_admin', 'admin', 'admin', 'User', '+66812345678', 'I am the admin of J1 Hub.', 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde', 'hidden', ST_SetSRID(ST_MakePoint(100.5018, 13.7563), 4326), NOW()),
		('prf_somchai', 'usr_somchai', 'somchai', 'Somchai', 'Deejai', '+66898765432', 'J1 student from KMTIL University, Thailand.', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d', 'show_friends', ST_SetSRID(ST_MakePoint(100.5118, 13.7663), 4326), NOW()),
		('prf_alice', 'usr_alice', 'alice', 'Alice', 'Smith', '+15551234567', 'J1 student from London, UK. Excited to work in the US!', 'https://images.unsplash.com/photo-1494790108377-be9c29b29330', 'show_anonymous', ST_SetSRID(ST_MakePoint(100.5218, 13.7763), 4326), NOW());`,

		// 4. Friendships
		`INSERT INTO friendships (friendship_id, user_id_1, user_id_2, status) VALUES
		('frn_somchai_alice', 'usr_somchai', 'usr_alice', 'accepted');`,

		// 5. User Phase History
		`INSERT INTO user_phase_history (history_id, user_id, phase_id, phase_points_earned, entered_at, completed_at) VALUES
		('uph_somchai_phs_001', 'usr_somchai', 'phs_001', 270, NOW() - INTERVAL '35 days', NOW() - INTERVAL '30 days'),
		('uph_somchai_phs_002', 'usr_somchai', 'phs_002', 0, NOW() - INTERVAL '30 days', NULL),
		('uph_alice_phs_001', 'usr_alice', 'phs_001', 310, NOW() - INTERVAL '65 days', NOW() - INTERVAL '60 days'),
		('uph_alice_phs_002', 'usr_alice', 'phs_002', 150, NOW() - INTERVAL '60 days', NOW() - INTERVAL '55 days'),
		('uph_alice_phs_003', 'usr_alice', 'phs_003', 0, NOW() - INTERVAL '55 days', NULL);`,

		// 6. Missions
		`INSERT INTO missions (mission_id, phase_id, title, description, location, base_points, is_mandatory, verification_type, due_date_type, fixed_due_date, relative_trigger_event, relative_days_offset) VALUES
		('mis_001', 'phs_001', 'Submit DS-2019 Form', 'Upload a scanned copy of your signed DS-2019 form.', 'Home Country', 100, true, 'upload', 'relative', NULL, 'arrival_date', -15),
		('mis_002', 'phs_001', 'Pay SEVIS Fee', 'Provide proof of SEVIS fee payment.', 'Home Country', 50, true, 'none', 'relative', NULL, 'arrival_date', -10),
		('mis_003', 'phs_002', 'Obtain SSN', 'Go to the local Social Security office to apply for your SSN.', 'US Local Office', 200, true, 'admin', 'relative', NULL, 'arrival_date', 7),
		('mis_004', 'phs_002', 'Update SEVIS Address', 'Update your housing address in your sponsor portal.', 'Online Portal', 100, true, 'none', 'relative', NULL, 'arrival_date', 3),
		('mis_005', 'phs_003', 'Complete Midterm Survey', 'Complete the midterm J1 experience questionnaire.', 'Online Survey', 150, false, 'none', 'relative', NULL, 'job_start_date', 60);`,

		// 7. User Missions
		`INSERT INTO user_missions (user_mission_id, user_id, mission_id, status, calculated_due_date, proof_url, proof_submitted_at, verified_at, verified_by, base_points_earned, speed_bonus_points, streak_bonus_points, first_completer_bonus_points, total_points_earned, rewarded_at) VALUES
		('ums_somchai_mis_001', 'usr_somchai', 'mis_001', 'completed', NOW() - INTERVAL '20 days', 'https://supabase.co/storage/v1/object/public/media/ds2019_somchai.pdf', NOW() - INTERVAL '22 days', NOW() - INTERVAL '21 days', 'usr_admin', 100, 20, 0, 0, 120, NOW() - INTERVAL '21 days'),
		('ums_somchai_mis_002', 'usr_somchai', 'mis_002', 'completed', NOW() - INTERVAL '15 days', NULL, NULL, NULL, NULL, 50, 0, 0, 0, 50, NOW() - INTERVAL '15 days'),
		('ums_somchai_mis_003', 'usr_somchai', 'mis_003', 'in_progress', NOW() + INTERVAL '5 days', NULL, NULL, NULL, NULL, 0, 0, 0, 0, 0, NULL),
		('ums_somchai_mis_004', 'usr_somchai', 'mis_004', 'completed', NOW() - INTERVAL '10 days', NULL, NULL, NULL, NULL, 100, 0, 0, 0, 100, NOW() - INTERVAL '10 days'),
		('ums_alice_mis_001', 'usr_alice', 'mis_001', 'completed', NOW() - INTERVAL '45 days', 'https://supabase.co/storage/v1/object/public/media/ds2019_alice.pdf', NOW() - INTERVAL '48 days', NOW() - INTERVAL '47 days', 'usr_admin', 100, 10, 0, 0, 110, NOW() - INTERVAL '47 days'),
		('ums_alice_mis_002', 'usr_alice', 'mis_002', 'completed', NOW() - INTERVAL '40 days', NULL, NULL, NULL, NULL, 50, 0, 0, 0, 50, NOW() - INTERVAL '40 days'),
		('ums_alice_mis_003', 'usr_alice', 'mis_003', 'completed', NOW() - INTERVAL '35 days', NULL, NULL, NOW() - INTERVAL '35 days', 'usr_admin', 200, 0, 0, 0, 200, NOW() - INTERVAL '35 days'),
		('ums_alice_mis_004', 'usr_alice', 'mis_004', 'completed', NOW() - INTERVAL '38 days', NULL, NULL, NULL, NULL, 100, 0, 0, 0, 100, NOW() - INTERVAL '38 days');`,

		// 8. Tasks
		`INSERT INTO tasks (task_id, mission_id, title, description) VALUES
		('tsk_001', 'mis_003', 'Locate SSN Office', 'Find the closest Social Security Administration office to your workplace.'),
		('tsk_002', 'mis_003', 'Gather Documents', 'Print and bring your passport, DS-2019, J1 visa, and printed I-94 arrival record.'),
		('tsk_003', 'mis_003', 'Fill Form SS-5', 'Fill out the Application for a Social Security Card (Form SS-5).');`,

		// 9. Seed User Tasks
		`INSERT INTO user_tasks (user_task_id, user_id, task_id, user_mission_id, is_completed, completed_at) VALUES
		('utk_som_tsk_001', 'usr_somchai', 'tsk_001', 'ums_somchai_mis_003', true, NOW() - INTERVAL '2 days'),
		('utk_som_tsk_002', 'usr_somchai', 'tsk_002', 'ums_somchai_mis_003', true, NOW() - INTERVAL '2 days'),
		('utk_som_tsk_003', 'usr_somchai', 'tsk_003', 'ums_somchai_mis_003', false, NULL);`,

		// 10. Seed Point Ledger
		`INSERT INTO point_ledger (ledger_id, user_id, source_type, source_id, delta, lifetime_balance_after, phase_balance_after, note) VALUES
		('plg_som_1', 'usr_somchai', 'mission_base', 'ums_somchai_mis_001', 100, 100, 100, 'Completed Submit DS-2019 Form'),
		('plg_som_2', 'usr_somchai', 'speed_bonus', 'ums_somchai_mis_001', 20, 120, 120, 'Speed bonus: Submitted DS-2019 early'),
		('plg_som_3', 'usr_somchai', 'mission_base', 'ums_somchai_mis_002', 50, 170, 170, 'Completed Pay SEVIS Fee'),
		('plg_som_4', 'usr_somchai', 'mission_base', 'ums_somchai_mis_004', 100, 270, 270, 'Completed Update SEVIS Address'),
		('plg_ali_1', 'usr_alice', 'mission_base', 'ums_alice_mis_001', 100, 100, 100, 'Completed Submit DS-2019 Form'),
		('plg_ali_2', 'usr_alice', 'speed_bonus', 'ums_alice_mis_001', 10, 110, 110, 'Speed bonus: Submitted DS-2019 early'),
		('plg_ali_3', 'usr_alice', 'mission_base', 'ums_alice_mis_002', 50, 160, 160, 'Completed Pay SEVIS Fee'),
		('plg_ali_4', 'usr_alice', 'mission_base', 'ums_alice_mis_004', 100, 260, 260, 'Completed Update SEVIS Address'),
		('plg_ali_5', 'usr_alice', 'mission_base', 'ums_alice_mis_003', 200, 460, 460, 'Completed Obtain SSN');`,

		// 11. Seed Badges
		`INSERT INTO badges (badge_id, title, description, trigger_type, icon_url) VALUES
		('bdg_001', 'Early Bird', 'Complete a mission 7+ days before its calculated due date.', 'speed', 'https://example.com/badge_early_bird.png'),
		('bdg_002', 'Streak Master', 'Maintain a streak of 3 or more mission completions.', 'streak', 'https://example.com/badge_streak.png'),
		('bdg_003', 'Pioneer', 'Be the first student to submit proof for a mission.', 'first_completer', 'https://example.com/badge_pioneer.png');`,

		// 12. Seed User Badges
		`INSERT INTO user_badges (user_badge_id, user_id, badge_id, source_id, earned_at) VALUES
		('ubg_som_1', 'usr_somchai', 'bdg_001', 'ums_somchai_mis_001', NOW() - INTERVAL '20 days'),
		('ubg_som_2', 'usr_somchai', 'bdg_002', 'ums_somchai_mis_001', NOW() - INTERVAL '10 days'),
		('ubg_ali_1', 'usr_alice', 'bdg_001', 'ums_alice_mis_001', NOW() - INTERVAL '45 days'),
		('ubg_ali_2', 'usr_alice', 'bdg_002', 'ums_alice_mis_001', NOW() - INTERVAL '40 days'),
		('ubg_ali_3', 'usr_alice', 'bdg_003', 'ums_alice_mis_002', NOW() - INTERVAL '35 days');`,

		// 13. Seed Credit Scores
		`INSERT INTO credit_scores (credit_id, user_id, current_score, last_updated) VALUES
		('crd_admin', 'usr_admin', 100, NOW()),
		('crd_somchai', 'usr_somchai', 100, NOW()),
		('crd_alice', 'usr_alice', 98, NOW());`,

		// 14. Seed Expense Transactions
		`INSERT INTO expense_transactions (transaction_id, paid_by_user_id, title, total_amount, currency, memo, transaction_date, due_date) VALUES
		('tx_001', 'usr_alice', 'June Room Rent', 1200.00, 'USD', 'Split for shared accommodation.', NOW() - INTERVAL '5 days', NOW() + INTERVAL '10 days');`,

		// 15. Seed Expense Splits
		`INSERT INTO expense_splits (split_id, transaction_id, user_id, owe_amount, payment_status, payment_method, payslip_url, approval_status, settled_at) VALUES
		('spl_001', 'tx_001', 'usr_somchai', 600.00, 'pending', NULL, NULL, 'pending_approval', NULL);`,

		// 16. Seed Job Postings
		`INSERT INTO job_postings (job_id, agency_name, employer_title, position, position_type, location_city, location_state, group_location, us_sponsor, salary_range_min, salary_range_max, available_slots, description, source_url, scrape_at, posted_at) VALUES
		('job_001', 'Global Student Exchange', 'Grand Teton Lodge Company', 'Front Desk Agent', 'Hospitality', 'Jackson', 'WY', 'Colter Bay Area', true, 16.00, 18.00, 5, 'Responsible for welcoming guests, executing check-in and check-out procedures, and answering phone calls.', 'https://example.com/job_001', NOW(), NOW()),
		('job_002', 'InterExchange', 'Cedar Point Amusement Park', 'Ride Operator', 'Amusement Park', 'Sandusky', 'OH', 'Cedar Point', true, 14.50, 15.00, 12, 'Ensure guest safety on park rides, direct queues, and execute ride launch sequences.', 'https://example.com/job_002', NOW(), NOW()),
		('job_003', 'CIEE Exchange', 'Oceanfront Dining Corp', 'Server Assistant / Busser', 'Restaurant', 'Myrtle Beach', 'SC', 'Coastal Dining Area', true, 12.00, 14.00, 3, 'Clear dirty dishes, wipe down and sanitize tables, set tables, and assist servers with food deliveries.', 'https://example.com/job_003', NOW(), NOW());`,

		// 17. Seed Job Housings
		`INSERT INTO job_housings (housing_id, job_id, description, weekly_rate, deposit, transportation, range_min_start_date, range_max_start_date) VALUES
		('hsg_001', 'job_001', 'On-site employee housing, shared cabin-style dorms with dining hall options.', 125.00, 200.00, 'Walking distance', '2026-05-15', '2026-06-15'),
		('hsg_002', 'job_002', 'Off-site apartments, fully furnished, with free park shuttle bus access.', 95.00, 150.00, 'Free shuttle', '2026-05-20', '2026-06-20');`,

		// 18. Seed User Carts
		`INSERT INTO user_carts (cart_id, user_id, job_id, status, added_at) VALUES
		('crt_som_job_001', 'usr_somchai', 'job_001', 'saved', NOW() - INTERVAL '12 days'),
		('crt_som_job_002', 'usr_somchai', 'job_002', 'applied', NOW() - INTERVAL '8 days');`,

		// 19. Seed Job Overall Ratings
		`INSERT INTO job_overall_ratings (rating_summary_id, job_id, overall_rate, agency_rate, job_rate, coworkers_rate, town_rate, hours_rate, housing_rate, second_job_feasibility_rate, overtime_availability_rate, review_count) VALUES
		('rtg_job_001', 'job_001', 4.50, 4.00, 4.50, 4.80, 4.20, 4.00, 4.50, 5.00, 4.50, 1);`,

		// 20. Seed Job Reviews
		`INSERT INTO job_reviews (review_id, job_id, user_id, rating_stars, review_text, tips_for_next_generation, score_agency, score_job, score_coworkers, score_town, score_hours, score_housing, score_second_job_feasibility, score_overtime_availability) VALUES
		('rev_001', 'job_001', 'usr_alice', 4.50, 'Outstanding summer role! Co-workers were extremely friendly and housing was just a few minutes walking distance from Grand Teton.', 'Bring some warm clothing as early mornings in Wyoming can be very chilly.', 4.00, 4.50, 4.80, 4.20, 4.00, 4.50, 5.00, 4.50);`,

		// 21. Seed User Jobs
		`INSERT INTO user_jobs (user_id, job_id, assigned_at, is_main, start_date, end_date) VALUES
		('usr_somchai', 'job_001', NOW() - INTERVAL '25 days', true, '2026-05-20 09:00:00', '2026-11-20 17:00:00'),
		('usr_alice', 'job_002', NOW() - INTERVAL '55 days', true, '2026-05-10 09:00:00', '2026-09-10 17:00:00');`,
	}

	for i, query := range queries {
		_, err := tx.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to execute seed query %d: %w", i+1, err)
		}
	}

	// Generate 50 mock missions and corresponding user_missions
	for i := 1; i <= 50; i++ {
		missionID := fmt.Sprintf("mis_mock_%03d", i)
		phaseNum := (i % 4) + 1
		phaseID := fmt.Sprintf("phs_%03d", phaseNum)
		title := fmt.Sprintf("Mock Mission %d", i)
		description := fmt.Sprintf("This is the description for mock mission %d", i)
		location := "Office"
		if i%2 == 0 {
			location = "Online"
		}
		basePoints := 50 + (i*10)%200
		isMandatory := i%2 == 0

		verificationType := "none"
		if i%3 == 1 {
			verificationType = "upload"
		} else if i%3 == 2 {
			verificationType = "admin"
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO missions (mission_id, phase_id, title, description, location, base_points, is_mandatory, verification_type, due_date_type, relative_trigger_event, relative_days_offset)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'relative', 'arrival_date', $9)`,
			missionID, phaseID, title, description, location, basePoints, isMandatory, verificationType, i%15,
		)
		if err != nil {
			return fmt.Errorf("failed to seed mock mission %d: %w", i, err)
		}

		// Seed user missions for usr_somchai and usr_alice
		var status string
		var dueDateOffset string

		if i%3 == 0 {
			status = "in_progress"
			dueDateOffset = "-5 days" // Overdue
		} else if i%3 == 1 {
			status = "completed"
			dueDateOffset = "-2 days" // Completed
		} else {
			status = "in_progress"
			dueDateOffset = "5 days"  // In Progress
		}

		for _, userID := range []string{"usr_somchai", "usr_alice"} {
			umID := fmt.Sprintf("ums_%s_mock_%03d", userID, i)

			var verifiedAt interface{}
			var verifiedBy interface{}
			var rewardedAt interface{}
			var pointsEarned int

			if status == "completed" {
				now := time.Now()
				verifiedAt = &now
				verifiedBy = "usr_admin"
				rewardedAt = &now
				pointsEarned = basePoints
			}

			_, err = tx.Exec(ctx, fmt.Sprintf(`
				INSERT INTO user_missions (
					user_mission_id, user_id, mission_id, status, calculated_due_date, 
					verified_at, verified_by, base_points_earned, total_points_earned, rewarded_at
				) VALUES ($1, $2, $3, $4, NOW() + INTERVAL '%s', $5, $6, $7, $7, $8)`, dueDateOffset),
				umID, userID, missionID, status, verifiedAt, verifiedBy, pointsEarned, rewardedAt,
			)
			if err != nil {
				return fmt.Errorf("failed to seed user mission for mock mission %d: %w", i, err)
			}
		}
	}

	// 10 more user missions for completed, in progress, locked/not_started
	for i := 1; i <= 10; i++ {
		missionID := fmt.Sprintf("mis_mock_extra_%03d", i)
		phaseID := "phs_003"
		title := fmt.Sprintf("Extra Mock Mission %d", i)

		_, err = tx.Exec(ctx, `
			INSERT INTO missions (mission_id, phase_id, title, description, location, base_points, is_mandatory, verification_type, due_date_type, relative_trigger_event, relative_days_offset)
			VALUES ($1, $2, $3, 'Extra mission', 'Virtual', 100, false, 'none', 'relative', 'arrival_date', 10)`,
			missionID, phaseID, title,
		)
		if err != nil {
			return fmt.Errorf("failed to seed extra mock mission %d: %w", i, err)
		}

		var status string
		var dueDateOffset string
		if i%3 == 1 {
			status = "completed"
			dueDateOffset = "-1 days"
		} else if i%3 == 2 {
			status = "in_progress"
			dueDateOffset = "5 days"
		} else {
			status = "not_started"
			dueDateOffset = "10 days"
		}

		umID := fmt.Sprintf("ums_alice_extra_%03d", i)
		var pointsEarned int
		if status == "completed" {
			pointsEarned = 100
		}

		_, err = tx.Exec(ctx, fmt.Sprintf(`
			INSERT INTO user_missions (user_mission_id, user_id, mission_id, status, calculated_due_date, base_points_earned, total_points_earned)
			VALUES ($1, 'usr_alice', $2, $3, NOW() + INTERVAL '%s', $4, $4)`, dueDateOffset),
			umID, missionID, status, pointsEarned,
		)
		if err != nil {
			return fmt.Errorf("failed to seed extra user mission %d: %w", i, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit seed transaction: %w", err)
	}

	log.Println("Mock data seeded successfully!")
	return nil
}
