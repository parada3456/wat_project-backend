package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	adminport "github.com/parada3456/wat_project-backend/internal/admin/port"
	missiondomain "github.com/parada3456/wat_project-backend/internal/mission/domain"
	userdomain "github.com/parada3456/wat_project-backend/internal/user/domain"
)

type adminRepo struct {
	pool *pgxpool.Pool
}

func NewAdminRepository(pool *pgxpool.Pool) adminport.AdminRepository {
	log.Println("debugprint: entering NewAdminRepository")
	return &adminRepo{pool: pool}
}

type timeToNullWrapper struct {
	t time.Time
}

func (w *timeToNullWrapper) Scan(value interface{}) error {
	log.Println("debugprint: entering (*timeToNullWrapper).Scan")
	if value == nil {
		w.t = time.Time{}
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("invalid time value: %v", value)
	}
	w.t = t
	return nil
}

func (r *adminRepo) GetStats(ctx context.Context) (*adminport.AdminStats, error) {
	log.Println("debugprint: entering (*adminRepo).GetStats")
	var stats adminport.AdminStats

	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	if err != nil {
		return nil, err
	}

	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE updated_at > NOW() - INTERVAL '30 days'").Scan(&stats.ActiveUsers)
	if err != nil {
		return nil, err
	}

	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM user_missions WHERE status = 'Pending_Verification'").Scan(&stats.PendingVerifications)
	if err != nil {
		return nil, err
	}

	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM job_postings").Scan(&stats.ActiveJobs)
	if err != nil {
		return nil, err
	}

	var avgCredit sql.NullFloat64
	err = r.pool.QueryRow(ctx, "SELECT AVG(current_score) FROM credit_scores").Scan(&avgCredit)
	if err != nil {
		return nil, err
	}
	if avgCredit.Valid {
		stats.AverageCreditScore = int(avgCredit.Float64)
	} else {
		stats.AverageCreditScore = 0
	}

	var sumPoints sql.NullInt64
	err = r.pool.QueryRow(ctx, "SELECT SUM(points) FROM point_ledger").Scan(&sumPoints)
	if err != nil {
		return nil, err
	}
	if sumPoints.Valid {
		stats.TotalPointsAwarded = int(sumPoints.Int64)
	} else {
		stats.TotalPointsAwarded = 0
	}

	return &stats, nil
}

func (r *adminRepo) ListPendingVerifications(ctx context.Context, limit, offset int) ([]missiondomain.UserMission, int, error) {
	log.Println("debugprint: entering (*adminRepo).ListPendingVerifications")
	var totalCount int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM user_missions WHERE status = 'Pending_Verification'").Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pending verifications: %w", err)
	}

	if totalCount == 0 {
		return []missiondomain.UserMission{}, 0, nil
	}

	query := `
		SELECT 
			user_mission_id, user_id, mission_id, status, calculated_due_date, 
			proof_url, proof_submitted_at, verified_at, verified_by, rewarded_at
		FROM user_missions 
		WHERE status = 'Pending_Verification'
		ORDER BY proof_submitted_at ASC
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query pending verifications: %w", err)
	}
	defer rows.Close()

	var ums []missiondomain.UserMission
	for rows.Next() {
		var um missiondomain.UserMission
		var verifiedBy *string
		err := rows.Scan(
			&um.UserMissionID, &um.UserID, &um.MissionID, &um.Status, &um.CalculatedDueDate,
			&um.ProofURL, &um.ProofSubmittedAt, &um.VerifiedAt, &verifiedBy, &um.RewardedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user mission: %w", err)
		}
		if verifiedBy != nil {
			um.VerifiedBy = *verifiedBy
		}
		ums = append(ums, um)
	}

	return ums, totalCount, nil
}

func (r *adminRepo) SearchUsers(ctx context.Context, query string, limit, offset int) ([]adminport.UserWithProfile, int, error) {
	log.Println("debugprint: entering (*adminRepo).SearchUsers")
	var totalCount int
	var countQuery string
	var countArgs []interface{}

	if query == "" {
		countQuery = "SELECT COUNT(*) FROM users"
	} else {
		countQuery = `
			SELECT COUNT(*) FROM users u
			LEFT JOIN profiles p ON u.user_id = p.user_id
			WHERE p.first_name ILIKE $1 OR p.last_name ILIKE $1 OR u.email ILIKE $1`
		countArgs = append(countArgs, "%"+query+"%")
	}

	err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	if totalCount == 0 {
		return []adminport.UserWithProfile{}, 0, nil
	}

	var rows pgx.Rows
	if query == "" {
		sql := `
			SELECT 
				u.user_id, u.email, u.password_hash, 
				u.current_phase_id, u.total_lifetime_points, u.current_phase_points, 
				u.mission_streak, u.arrival_date, u.job_start_date, u.created_at, u.updated_at,
				p.profile_id, COALESCE(p.first_name, ''), COALESCE(p.last_name, ''), COALESCE(p.phone_number, ''),
				COALESCE(p.bio, ''), COALESCE(p.avatar_url, ''), COALESCE(p.radar_visibility, 'hidden')
			FROM users u
			LEFT JOIN profiles p ON u.user_id = p.user_id
			ORDER BY u.created_at DESC
			LIMIT $1 OFFSET $2`
		rows, err = r.pool.Query(ctx, sql, limit, offset)
	} else {
		sql := `
			SELECT 
				u.user_id, u.email, u.password_hash, 
				u.current_phase_id, u.total_lifetime_points, u.current_phase_points, 
				u.mission_streak, u.arrival_date, u.job_start_date, u.created_at, u.updated_at,
				p.profile_id, COALESCE(p.first_name, ''), COALESCE(p.last_name, ''), COALESCE(p.phone_number, ''),
				COALESCE(p.bio, ''), COALESCE(p.avatar_url, ''), COALESCE(p.radar_visibility, 'hidden')
			FROM users u
			LEFT JOIN profiles p ON u.user_id = p.user_id
			WHERE p.first_name ILIKE $1 OR p.last_name ILIKE $1 OR u.email ILIKE $1
			ORDER BY u.created_at DESC
			LIMIT $2 OFFSET $3`
		rows, err = r.pool.Query(ctx, sql, "%"+query+"%", limit, offset)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []adminport.UserWithProfile
	for rows.Next() {
		var u userdomain.User
		var p userdomain.Profile
		var currentPhaseID *string
		var arrivalDate *timeToNullWrapper
		var jobStartDate *timeToNullWrapper

		err := rows.Scan(
			&u.UserID, &u.Email, &u.PasswordHash,
			&currentPhaseID, &u.TotalLifetimePoints, &u.CurrentPhasePoints,
			&u.MissionStreak, &arrivalDate, &jobStartDate, &u.CreatedAt, &u.UpdatedAt,
			&p.ProfileID, &p.FirstName, &p.LastName, &p.PhoneNumber,
			&p.Bio, &p.AvatarURL, &p.RadarVisibility,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
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
		p.UserID = u.UserID
		users = append(users, adminport.UserWithProfile{User: u, Profile: p})
	}

	return users, totalCount, nil
}
