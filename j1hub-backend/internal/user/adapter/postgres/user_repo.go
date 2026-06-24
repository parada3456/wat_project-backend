package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	userdomain "github.com/j1hub/backend/internal/user/domain"

	"github.com/j1hub/backend/internal/domain"
	port "github.com/j1hub/backend/internal/user/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) port.UserRepository {
	log.Println("debugprint: entering NewUserRepository")
	return &userRepo{pool: pool}
}

func stringToNull(s string) interface{} {
	log.Println("debugprint: entering stringToNull")
	if s == "" {
		return nil
	}
	return s
}

func timeToNull(t time.Time) interface{} {
	log.Println("debugprint: entering timeToNull")
	if t.IsZero() {
		return nil
	}
	return t
}

func scanUser(row pgx.Row) (*userdomain.User, error) {
	log.Println("debugprint: entering scanUser")
	var u userdomain.User
	var currentPhaseID *string
	var arrivalDate *time.Time
	var jobStartDate *time.Time

	err := row.Scan(
		&u.UserID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
		&currentPhaseID, &u.TotalLifetimePoints, &u.CurrentPhasePoints,
		&u.MissionStreak, &arrivalDate, &jobStartDate, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
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

	return &u, nil
}

func (r *userRepo) Create(ctx context.Context, u *userdomain.User) error {
	log.Println("debugprint: entering (*userRepo).Create")
	query := `
		INSERT INTO users (
			user_id, email, password_hash, first_name, last_name, 
			current_phase_id, total_lifetime_points, current_phase_points, 
			mission_streak, arrival_date, job_start_date, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := r.pool.Exec(ctx, query,
		u.UserID, u.Email, u.PasswordHash, u.FirstName, u.LastName,
		stringToNull(u.CurrentPhaseID), u.TotalLifetimePoints, u.CurrentPhasePoints,
		u.MissionStreak, timeToNull(u.ArrivalDate), timeToNull(u.JobStartDate), u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*userdomain.User, error) {
	log.Println("debugprint: entering (*userRepo).FindByID")
	query := `
		SELECT 
			user_id, email, password_hash, first_name, last_name, 
			current_phase_id, total_lifetime_points, current_phase_points, 
			mission_streak, arrival_date, job_start_date, created_at, updated_at
		FROM users WHERE user_id = $1`

	row := r.pool.QueryRow(ctx, query, id)
	u, err := scanUser(row)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	return u, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	log.Println("debugprint: entering (*userRepo).FindByEmail")
	query := `
		SELECT 
			user_id, email, password_hash, first_name, last_name, 
			current_phase_id, total_lifetime_points, current_phase_points, 
			mission_streak, arrival_date, job_start_date, created_at, updated_at
		FROM users WHERE email = $1`

	row := r.pool.QueryRow(ctx, query, email)
	u, err := scanUser(row)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return u, nil
}

func (r *userRepo) Update(ctx context.Context, u *userdomain.User) error {
	log.Println("debugprint: entering (*userRepo).Update")
	query := `
		UPDATE users SET 
			email = $1, password_hash = $2, first_name = $3, last_name = $4, 
			current_phase_id = $5, total_lifetime_points = $6, current_phase_points = $7, 
			mission_streak = $8, arrival_date = $9, job_start_date = $10, updated_at = $11
		WHERE user_id = $12`

	_, err := r.pool.Exec(ctx, query,
		u.Email, u.PasswordHash, u.FirstName, u.LastName,
		stringToNull(u.CurrentPhaseID), u.TotalLifetimePoints, u.CurrentPhasePoints,
		u.MissionStreak, timeToNull(u.ArrivalDate), timeToNull(u.JobStartDate), u.UpdatedAt, u.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepo) IncrementPoints(ctx context.Context, userID string, lifetimeDelta, phaseDelta int) error {
	log.Println("debugprint: entering (*userRepo).IncrementPoints")
	query := `
		UPDATE users SET 
			total_lifetime_points = total_lifetime_points + $1,
			current_phase_points = current_phase_points + $2,
			updated_at = NOW()
		WHERE user_id = $3`

	_, err := r.pool.Exec(ctx, query, lifetimeDelta, phaseDelta, userID)
	if err != nil {
		return fmt.Errorf("failed to increment points: %w", err)
	}
	return nil
}

func (r *userRepo) ResetStreak(ctx context.Context, userID string) error {
	log.Println("debugprint: entering (*userRepo).ResetStreak")
	query := `UPDATE users SET mission_streak = 0, updated_at = NOW() WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to reset streak: %w", err)
	}
	return nil
}

func (r *userRepo) SetPhase(ctx context.Context, userID, phaseID string) error {
	log.Println("debugprint: entering (*userRepo).SetPhase")
	query := `UPDATE users SET current_phase_id = $1, updated_at = NOW() WHERE user_id = $2`
	_, err := r.pool.Exec(ctx, query, phaseID, userID)
	if err != nil {
		return fmt.Errorf("failed to set phase: %w", err)
	}
	return nil
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	log.Println("debugprint: entering (*userRepo).Delete")
	query := `DELETE FROM users WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *userRepo) FindUserJob(ctx context.Context, userID string) (*userdomain.UserJob, error) {
	log.Println("debugprint: entering (*userRepo).FindUserJob")
	query := `SELECT user_id, job_id, assigned_at, is_main, start_date, end_date FROM user_jobs WHERE user_id = $1 AND is_main = true`
	var uj userdomain.UserJob
	err := r.pool.QueryRow(ctx, query, userID).Scan(&uj.UserID, &uj.JobID, &uj.AssignedAt, &uj.IsMain, &uj.StartDate, &uj.EndDate)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user job: %w", err)
	}
	return &uj, nil
}

func (r *userRepo) FindUserJobs(ctx context.Context, userID string) ([]userdomain.UserJob, error) {
	log.Println("debugprint: entering (*userRepo).FindUserJobs")
	query := `SELECT user_id, job_id, assigned_at, is_main, start_date, end_date FROM user_jobs WHERE user_id = $1 ORDER BY assigned_at DESC`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user jobs: %w", err)
	}
	defer rows.Close()

	var ujs []userdomain.UserJob
	for rows.Next() {
		var uj userdomain.UserJob
		err := rows.Scan(&uj.UserID, &uj.JobID, &uj.AssignedAt, &uj.IsMain, &uj.StartDate, &uj.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user job: %w", err)
		}
		ujs = append(ujs, uj)
	}
	return ujs, nil
}

func (r *userRepo) AssignJob(ctx context.Context, userID, jobID string, isMain bool, startDate, endDate *time.Time) error {
	log.Println("debugprint: entering (*userRepo).AssignJob")
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if isMain {
		_, err = tx.Exec(ctx, `UPDATE user_jobs SET is_main = false WHERE user_id = $1`, userID)
		if err != nil {
			return fmt.Errorf("failed to update existing user jobs to non-main: %w", err)
		}
	} else {
		var hasMain bool
		err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM user_jobs WHERE user_id = $1 AND is_main = true)`, userID).Scan(&hasMain)
		if err != nil {
			return fmt.Errorf("failed to check existing main job: %w", err)
		}
		if !hasMain {
			isMain = true
		}
	}

	query := `
		INSERT INTO user_jobs (user_id, job_id, assigned_at, is_main, start_date, end_date)
		VALUES ($1, $2, NOW(), $3, $4, $5)
		ON CONFLICT (user_id, job_id) 
		DO UPDATE SET is_main = EXCLUDED.is_main, start_date = EXCLUDED.start_date, end_date = EXCLUDED.end_date, assigned_at = NOW()`
	
	_, err = tx.Exec(ctx, query, userID, jobID, isMain, startDate, endDate)
	if err != nil {
		return fmt.Errorf("failed to assign job to user: %w", err)
	}

	return tx.Commit(ctx)
}

