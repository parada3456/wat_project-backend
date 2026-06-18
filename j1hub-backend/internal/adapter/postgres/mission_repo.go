package postgres

import (
	"context"
	"time"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type missionRepo struct {
	pool *pgxpool.Pool
}

func NewMissionRepository(pool *pgxpool.Pool) port.MissionRepository {
	return &missionRepo{pool: pool}
}

func (r *missionRepo) FindByPhase(ctx context.Context, phaseID string) ([]domain.Mission, error) {
	query := `SELECT mission_id, phase_id, title, description, location, base_points, is_mandatory, verification_type, due_date_type, fixed_due_date, relative_trigger_event, relative_days_offset, created_at, updated_at FROM missions WHERE phase_id = $1`
	rows, err := r.pool.Query(ctx, query, phaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var missions []domain.Mission
	for rows.Next() {
		var m domain.Mission
		var desc *string
		var loc *string
		var triggerEvent *string
		var daysOffset *int
		if err := rows.Scan(&m.MissionID, &m.PhaseID, &m.Title, &desc, &loc, &m.BasePoints, &m.IsMandatory, &m.VerificationType, &m.DueDateType, &m.FixedDueDate, &triggerEvent, &daysOffset, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		if desc != nil {
			m.Description = *desc
		}
		if loc != nil {
			m.Location = *loc
		}
		if triggerEvent != nil {
			m.RelativeTriggerEvent = *triggerEvent
		}
		if daysOffset != nil {
			m.RelativeDaysOffset = *daysOffset
		}
		missions = append(missions, m)
	}
	return missions, nil
}

func (r *missionRepo) FindByID(ctx context.Context, id string) (*domain.Mission, error) {
	query := `SELECT mission_id, phase_id, title, description, location, base_points, is_mandatory, verification_type, due_date_type, fixed_due_date, relative_trigger_event, relative_days_offset, created_at, updated_at FROM missions WHERE mission_id = $1`
	var m domain.Mission
	var desc *string
	var loc *string
	var triggerEvent *string
	var daysOffset *int
	err := r.pool.QueryRow(ctx, query, id).Scan(&m.MissionID, &m.PhaseID, &m.Title, &desc, &loc, &m.BasePoints, &m.IsMandatory, &m.VerificationType, &m.DueDateType, &m.FixedDueDate, &triggerEvent, &daysOffset, &m.CreatedAt, &m.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err == nil {
		if desc != nil {
			m.Description = *desc
		}
		if loc != nil {
			m.Location = *loc
		}
		if triggerEvent != nil {
			m.RelativeTriggerEvent = *triggerEvent
		}
		if daysOffset != nil {
			m.RelativeDaysOffset = *daysOffset
		}
	}
	return &m, err
}

type userMissionRepo struct {
	pool *pgxpool.Pool
}

func NewUserMissionRepository(pool *pgxpool.Pool) port.UserMissionRepository {
	return &userMissionRepo{pool: pool}
}

func (r *userMissionRepo) BulkInsert(ctx context.Context, ums []domain.UserMission) error {
	batch := &pgx.Batch{}
	for _, um := range ums {
		batch.Queue(`INSERT INTO user_missions (user_mission_id, user_id, mission_id, status, calculated_due_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			um.UserMissionID, um.UserID, um.MissionID, um.Status, um.CalculatedDueDate, um.CreatedAt, um.UpdatedAt)
	}
	return r.pool.SendBatch(ctx, batch).Close()
}

func (r *userMissionRepo) FindByUserAndPhase(ctx context.Context, userID, phaseID string) ([]domain.UserMission, error) {
	query := `SELECT um.user_mission_id, um.user_id, um.mission_id, um.status, um.calculated_due_date, um.proof_url, um.proof_submitted_at, um.verified_at, um.verified_by, um.base_points_earned, um.speed_bonus_points, um.streak_bonus_points, um.first_completer_bonus_points, um.total_points_earned, um.rewarded_at, um.created_at, um.updated_at 
	FROM user_missions um JOIN missions m ON um.mission_id = m.mission_id WHERE um.user_id = $1 AND m.phase_id = $2`
	rows, err := r.pool.Query(ctx, query, userID, phaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ums []domain.UserMission
	for rows.Next() {
		var um domain.UserMission
		var proofURL *string
		var verifiedBy *string
		if err := rows.Scan(&um.UserMissionID, &um.UserID, &um.MissionID, &um.Status, &um.CalculatedDueDate, &proofURL, &um.ProofSubmittedAt, &um.VerifiedAt, &verifiedBy, &um.BasePointsEarned, &um.SpeedBonusPoints, &um.StreakBonusPoints, &um.FirstCompleterBonusPoints, &um.TotalPointsEarned, &um.RewardedAt, &um.CreatedAt, &um.UpdatedAt); err != nil {
			return nil, err
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

func (r *userMissionRepo) FindByID(ctx context.Context, id string) (*domain.UserMission, error) {
	query := `SELECT user_mission_id, user_id, mission_id, status, calculated_due_date, proof_url, proof_submitted_at, verified_at, verified_by, base_points_earned, speed_bonus_points, streak_bonus_points, first_completer_bonus_points, total_points_earned, rewarded_at, created_at, updated_at FROM user_missions WHERE user_mission_id = $1`
	var um domain.UserMission
	var proofURL *string
	var verifiedBy *string
	err := r.pool.QueryRow(ctx, query, id).Scan(&um.UserMissionID, &um.UserID, &um.MissionID, &um.Status, &um.CalculatedDueDate, &proofURL, &um.ProofSubmittedAt, &um.VerifiedAt, &verifiedBy, &um.BasePointsEarned, &um.SpeedBonusPoints, &um.StreakBonusPoints, &um.FirstCompleterBonusPoints, &um.TotalPointsEarned, &um.RewardedAt, &um.CreatedAt, &um.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err == nil {
		if proofURL != nil {
			um.ProofURL = *proofURL
		}
		if verifiedBy != nil {
			um.VerifiedBy = *verifiedBy
		}
	}
	return &um, err
}

func (r *userMissionRepo) UpdateStatus(ctx context.Context, id string, status domain.UserMissionStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE user_missions SET status = $1, updated_at = NOW() WHERE user_mission_id = $2`, status, id)
	return err
}

func (r *userMissionRepo) UpdateVerification(ctx context.Context, id string, verifiedAt time.Time, verifiedBy string) error {
	_, err := r.pool.Exec(ctx, `UPDATE user_missions SET verified_at = $1, verified_by = $2, updated_at = NOW() WHERE user_mission_id = $3`, verifiedAt, verifiedBy, id)
	return err
}

func (r *userMissionRepo) UpdateReward(ctx context.Context, id string, reward *domain.PointReward, rewardedAt time.Time) error {
	_, err := r.pool.Exec(ctx, `UPDATE user_missions SET base_points_earned = $1, speed_bonus_points = $2, streak_bonus_points = $3, first_completer_bonus_points = $4, total_points_earned = $5, rewarded_at = $6, updated_at = NOW() WHERE user_mission_id = $7`,
		reward.Base, reward.SpeedBonus, reward.StreakBonus, reward.FirstCompleterBonus, reward.Total, rewardedAt, id)
	return err
}

func (r *userMissionRepo) FindOverdue(ctx context.Context) ([]domain.UserMission, error) {
	query := `SELECT user_mission_id, user_id, mission_id, status, calculated_due_date, created_at, updated_at FROM user_missions WHERE status IN ('Not_Started', 'In_Progress', 'Pending_Verification') AND calculated_due_date < NOW()`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ums []domain.UserMission
	for rows.Next() {
		var um domain.UserMission
		if err := rows.Scan(&um.UserMissionID, &um.UserID, &um.MissionID, &um.Status, &um.CalculatedDueDate, &um.CreatedAt, &um.UpdatedAt); err != nil {
			return nil, err
		}
		ums = append(ums, um)
	}
	return ums, nil
}

type taskRepo struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) port.TaskRepository {
	return &taskRepo{pool: pool}
}

func (r *taskRepo) FindByMission(ctx context.Context, missionID string) ([]domain.Task, error) {
	rows, err := r.pool.Query(ctx, `SELECT task_id, mission_id, title, description, created_at, updated_at FROM tasks WHERE mission_id = $1`, missionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		var desc *string
		if err := rows.Scan(&t.TaskID, &t.MissionID, &t.Title, &desc, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if desc != nil {
			t.Description = *desc
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

type userTaskRepo struct {
	pool *pgxpool.Pool
}

func NewUserTaskRepository(pool *pgxpool.Pool) port.UserTaskRepository {
	return &userTaskRepo{pool: pool}
}

func (r *userTaskRepo) Upsert(ctx context.Context, ut *domain.UserTask) error {
	query := `
		INSERT INTO user_tasks (user_task_id, user_id, task_id, user_mission_id, is_completed, completed_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_task_id) DO UPDATE SET 
			is_completed = EXCLUDED.is_completed,
			completed_at = EXCLUDED.completed_at,
			updated_at = EXCLUDED.updated_at`
	_, err := r.pool.Exec(ctx, query, ut.UserTaskID, ut.UserID, ut.TaskID, ut.UserMissionID, ut.IsCompleted, ut.CompletedAt, ut.UpdatedAt)
	return err
}

func (r *userTaskRepo) FindByUserMission(ctx context.Context, userMissionID string) ([]domain.UserTask, error) {
	rows, err := r.pool.Query(ctx, `SELECT user_task_id, user_id, task_id, user_mission_id, is_completed, completed_at, updated_at FROM user_tasks WHERE user_mission_id = $1`, userMissionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var uts []domain.UserTask
	for rows.Next() {
		var ut domain.UserTask
		if err := rows.Scan(&ut.UserTaskID, &ut.UserID, &ut.TaskID, &ut.UserMissionID, &ut.IsCompleted, &ut.CompletedAt, &ut.UpdatedAt); err != nil {
			return nil, err
		}
		uts = append(uts, ut)
	}
	return uts, nil
}

type journeyPhaseRepo struct {
	pool *pgxpool.Pool
}

func NewJourneyPhaseRepository(pool *pgxpool.Pool) port.JourneyPhaseRepository {
	return &journeyPhaseRepo{pool: pool}
}

func (r *journeyPhaseRepo) FindByNumber(ctx context.Context, number int) (*domain.JourneyPhase, error) {
	var jp domain.JourneyPhase
	var desc *string
	err := r.pool.QueryRow(ctx, `SELECT phase_id, phase_number, title, description, created_at, updated_at FROM journey_phases WHERE phase_number = $1`, number).Scan(&jp.PhaseID, &jp.PhaseNumber, &jp.Title, &desc, &jp.CreatedAt, &jp.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err == nil && desc != nil {
		jp.Description = *desc
	}
	return &jp, err
}

func (r *journeyPhaseRepo) FindByID(ctx context.Context, id string) (*domain.JourneyPhase, error) {
	var jp domain.JourneyPhase
	var desc *string
	err := r.pool.QueryRow(ctx, `SELECT phase_id, phase_number, title, description, created_at, updated_at FROM journey_phases WHERE phase_id = $1`, id).Scan(&jp.PhaseID, &jp.PhaseNumber, &jp.Title, &desc, &jp.CreatedAt, &jp.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err == nil && desc != nil {
		jp.Description = *desc
	}
	return &jp, err
}

type userPhaseHistoryRepo struct {
	pool *pgxpool.Pool
}

func NewUserPhaseHistoryRepository(pool *pgxpool.Pool) port.UserPhaseHistoryRepository {
	return &userPhaseHistoryRepo{pool: pool}
}

func (r *userPhaseHistoryRepo) Insert(ctx context.Context, h *domain.UserPhaseHistory) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO user_phase_history (history_id, user_id, phase_id, phase_points_earned, entered_at, completed_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		h.HistoryID, h.UserID, h.PhaseID, h.PhasePointsEarned, h.EnteredAt, h.CompletedAt)
	return err
}

func (r *userPhaseHistoryRepo) CompleteCurrentPhase(ctx context.Context, userID string, points int, completedAt time.Time) error {
	_, err := r.pool.Exec(ctx, `UPDATE user_phase_history SET phase_points_earned = $1, completed_at = $2 WHERE user_id = $3 AND completed_at IS NULL`, points, completedAt, userID)
	return err
}

func (r *userPhaseHistoryRepo) FindByUserAndPhase(ctx context.Context, userID, phaseID string) (*domain.UserPhaseHistory, error) {
	var h domain.UserPhaseHistory
	err := r.pool.QueryRow(ctx, `SELECT history_id, user_id, phase_id, phase_points_earned, entered_at, completed_at FROM user_phase_history WHERE user_id = $1 AND phase_id = $2`, userID, phaseID).Scan(&h.HistoryID, &h.UserID, &h.PhaseID, &h.PhasePointsEarned, &h.EnteredAt, &h.CompletedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &h, err
}
