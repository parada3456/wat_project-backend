package postgres

import (
	"context"
	"log"
	"time"

	userdomain "github.com/j1hub/backend/internal/user/domain"

	port "github.com/j1hub/backend/internal/gamification/port"
	"github.com/jackc/pgx/v5/pgxpool"
)

type leaderboardRepo struct {
	pool *pgxpool.Pool
}

func NewLeaderboardRepository(pool *pgxpool.Pool) port.LeaderboardRepository {
	log.Println("debugprint: entering NewLeaderboardRepository")
	return &leaderboardRepo{pool: pool}
}

func (r *leaderboardRepo) FindByScope(ctx context.Context, scope, jobID string) ([]userdomain.User, error) {
	log.Println("debugprint: entering (*leaderboardRepo).FindByScope")
	query := `
		SELECT u.user_id, u.email, u.first_name, u.last_name, u.current_phase_id, u.total_lifetime_points, u.current_phase_points, u.mission_streak, u.arrival_date, u.job_start_date, u.created_at, u.updated_at
		FROM users u`

	var args []interface{}
	if scope == "employer" && jobID != "" {
		query += ` JOIN user_carts uc ON uc.user_id = u.user_id WHERE uc.job_id = $1 AND uc.status = 'Applied'`
		args = append(args, jobID)
	}

	query += ` ORDER BY u.current_phase_points DESC LIMIT 50`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []userdomain.User
	for rows.Next() {
		var u userdomain.User
		var currentPhaseID *string
		var arrivalDate *time.Time
		var jobStartDate *time.Time
		if err := rows.Scan(
			&u.UserID, &u.Email, &u.FirstName, &u.LastName,
			&currentPhaseID, &u.TotalLifetimePoints, &u.CurrentPhasePoints,
			&u.MissionStreak, &arrivalDate, &jobStartDate, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if currentPhaseID != nil {
			u.CurrentPhaseID = *currentPhaseID
		}
		if arrivalDate != nil {
			u.ArrivalDate = *arrivalDate
		}
		if jobStartDate != nil {
			u.JobStartDate = *jobStartDate
		}
		users = append(users, u)
	}
	return users, nil
}
