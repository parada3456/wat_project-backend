package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type adminRepo struct {
	pool *pgxpool.Pool
}

func NewAdminRepository(pool *pgxpool.Pool) port.AdminRepository {
	log.Println("debugprint: entering NewAdminRepository")
	return &adminRepo{pool: pool}
}

func (r *adminRepo) GetStats(ctx context.Context) (*port.AdminStats, error) {
	log.Println("debugprint: entering (*adminRepo).GetStats")
	stats := &port.AdminStats{}

	// 1. Total Users
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to get total users: %w", err)
	}

	// 2. Active Users (using count for now or filtering where updated_at is recent)
	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&stats.ActiveUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}

	// 3. Pending Verifications
	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM user_missions WHERE status = $1", domain.StatusPendingVerification).Scan(&stats.PendingVerifications)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending verifications: %w", err)
	}

	// 4. Active Jobs
	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM job_postings").Scan(&stats.ActiveJobs)
	if err != nil {
		return nil, fmt.Errorf("failed to get active jobs: %w", err)
	}

	// 5. Average Credit Score
	var avgScore float64
	err = r.pool.QueryRow(ctx, "SELECT COALESCE(AVG(score), 0) FROM credit_scores").Scan(&avgScore)
	if err != nil {
		return nil, fmt.Errorf("failed to get average credit score: %w", err)
	}
	stats.AverageCreditScore = int(avgScore)

	// 6. Total Points Awarded
	err = r.pool.QueryRow(ctx, "SELECT COALESCE(SUM(points), 0) FROM point_ledger").Scan(&stats.TotalPointsAwarded)
	if err != nil {
		return nil, fmt.Errorf("failed to get total points: %w", err)
	}

	return stats, nil
}

func (r *adminRepo) ListPendingVerifications(ctx context.Context) ([]domain.UserMission, error) {
	log.Println("debugprint: entering (*adminRepo).ListPendingVerifications")
	query := `
		SELECT 
			user_mission_id, user_id, mission_id, status, calculated_due_date, 
			proof_url, proof_submitted_at, verified_at, verified_by, 
			base_points_earned, speed_bonus_points, streak_bonus_points, 
			first_completer_bonus_points, total_points_earned, rewarded_at, 
			created_at, updated_at 
		FROM user_missions 
		WHERE status = $1 
		ORDER BY proof_submitted_at ASC`

	rows, err := r.pool.Query(ctx, query, domain.StatusPendingVerification)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending verifications: %w", err)
	}
	defer rows.Close()

	var ums []domain.UserMission
	for rows.Next() {
		var um domain.UserMission
		var proofURL *string
		var verifiedBy *string
		err := rows.Scan(
			&um.UserMissionID, &um.UserID, &um.MissionID, &um.Status, &um.CalculatedDueDate,
			&proofURL, &um.ProofSubmittedAt, &um.VerifiedAt, &verifiedBy,
			&um.BasePointsEarned, &um.SpeedBonusPoints, &um.StreakBonusPoints,
			&um.FirstCompleterBonusPoints, &um.TotalPointsEarned, &um.RewardedAt,
			&um.CreatedAt, &um.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pending verification: %w", err)
		}
		if proofURL != nil {
			um.ProofURL = *proofURL
		}
		if verifiedBy != nil {
			um.VerifiedBy = *verifiedBy
		}
		ums = append(ums, um)
	}

	return ums, nil
}

func (r *adminRepo) SearchUsers(ctx context.Context, query string) ([]domain.User, error) {
	log.Println("debugprint: entering (*adminRepo).SearchUsers")
	var rows pgx.Rows
	var err error

	if query == "" {
		sql := `
			SELECT 
				user_id, email, password_hash, first_name, last_name, 
				current_phase_id, total_lifetime_points, current_phase_points, 
				mission_streak, arrival_date, job_start_date, created_at, updated_at
			FROM users 
			ORDER BY created_at DESC`
		rows, err = r.pool.Query(ctx, sql)
	} else {
		sql := `
			SELECT 
				user_id, email, password_hash, first_name, last_name, 
				current_phase_id, total_lifetime_points, current_phase_points, 
				mission_streak, arrival_date, job_start_date, created_at, updated_at
			FROM users 
			WHERE first_name ILIKE $1 OR last_name ILIKE $1 OR email ILIKE $1
			ORDER BY created_at DESC`
		rows, err = r.pool.Query(ctx, sql, "%"+query+"%")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		var currentPhaseID *string
		var arrivalDate *timeToNullWrapper
		var jobStartDate *timeToNullWrapper

		err := rows.Scan(
			&u.UserID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
			&currentPhaseID, &u.TotalLifetimePoints, &u.CurrentPhasePoints,
			&u.MissionStreak, &arrivalDate, &jobStartDate, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if currentPhaseID != nil {
			u.CurrentPhaseID = *currentPhaseID
		}
		if arrivalDate != nil && !arrivalDate.t.IsZero() {
			u.ArrivalDate = arrivalDate.t
		}
		if jobStartDate != nil && !jobStartDate.t.IsZero() {
			u.JobStartDate = jobStartDate.t
		}
		users = append(users, u)
	}

	return users, nil
}

type timeToNullWrapper struct {
	t javaTime
}

type javaTime = time.Time

func (w *timeToNullWrapper) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	t, ok := value.(javaTime)
	if !ok {
		return fmt.Errorf("expected time.Time, got %T", value)
	}
	w.t = t
	return nil
}
