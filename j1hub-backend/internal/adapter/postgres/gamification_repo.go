package postgres

import (
	"context"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pointLedgerRepo struct {
	pool *pgxpool.Pool
}

func NewPointLedgerRepository(pool *pgxpool.Pool) port.PointLedgerRepository {
	return &pointLedgerRepo{pool: pool}
}

func (r *pointLedgerRepo) Insert(ctx context.Context, l *domain.PointLedger) error {
	query := `INSERT INTO point_ledger (ledger_id, user_id, source_type, source_id, delta, lifetime_balance_after, phase_balance_after, note, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pool.Exec(ctx, query, l.LedgerID, l.UserID, l.SourceType, l.SourceID, l.Delta, l.LifetimeBalanceAfter, l.PhaseBalanceAfter, l.Note, l.CreatedAt)
	return err
}

func (r *pointLedgerRepo) InsertBatch(ctx context.Context, ledgers []domain.PointLedger) error {
	batch := &pgx.Batch{}
	for _, l := range ledgers {
		batch.Queue(`INSERT INTO point_ledger (ledger_id, user_id, source_type, source_id, delta, lifetime_balance_after, phase_balance_after, note, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			l.LedgerID, l.UserID, l.SourceType, l.SourceID, l.Delta, l.LifetimeBalanceAfter, l.PhaseBalanceAfter, l.Note, l.CreatedAt)
	}
	return r.pool.SendBatch(ctx, batch).Close()
}

type badgeRepo struct {
	pool *pgxpool.Pool
}

func NewBadgeRepository(pool *pgxpool.Pool) port.BadgeRepository {
	return &badgeRepo{pool: pool}
}

func (r *badgeRepo) FindByTriggerType(ctx context.Context, triggerType domain.TriggerType) ([]domain.Badge, error) {
	rows, err := r.pool.Query(ctx, `SELECT badge_id, title, description, trigger_type, icon_url, created_at FROM badges WHERE trigger_type = $1`, triggerType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var badges []domain.Badge
	for rows.Next() {
		var b domain.Badge
		if err := rows.Scan(&b.BadgeID, &b.Title, &b.Description, &b.TriggerType, &b.IconURL, &b.CreatedAt); err != nil {
			return nil, err
		}
		badges = append(badges, b)
	}
	return badges, nil
}

func (r *badgeRepo) FindEligible(ctx context.Context, userID string, triggerType domain.TriggerType) ([]domain.Badge, error) {
	// Simple implementation: find all badges of this trigger type that the user doesn't have yet
	query := `
		SELECT b.badge_id, b.title, b.description, b.trigger_type, b.icon_url, b.created_at 
		FROM badges b 
		WHERE b.trigger_type = $1 
		AND NOT EXISTS (SELECT 1 FROM user_badges ub WHERE ub.user_id = $2 AND ub.badge_id = b.badge_id)`
	rows, err := r.pool.Query(ctx, query, triggerType, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var badges []domain.Badge
	for rows.Next() {
		var b domain.Badge
		if err := rows.Scan(&b.BadgeID, &b.Title, &b.Description, &b.TriggerType, &b.IconURL, &b.CreatedAt); err != nil {
			return nil, err
		}
		badges = append(badges, b)
	}
	return badges, nil
}

type userBadgeRepo struct {
	pool *pgxpool.Pool
}

func NewUserBadgeRepository(pool *pgxpool.Pool) port.UserBadgeRepository {
	return &userBadgeRepo{pool: pool}
}

func (r *userBadgeRepo) Insert(ctx context.Context, ub *domain.UserBadge) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO user_badges (user_badge_id, user_id, badge_id, source_id, earned_at) VALUES ($1, $2, $3, $4, $5)`,
		ub.UserBadgeID, ub.UserID, ub.BadgeID, ub.SourceID, ub.EarnedAt)
	return err
}

func (r *userBadgeRepo) FindByUser(ctx context.Context, userID string) ([]domain.UserBadge, error) {
	rows, err := r.pool.Query(ctx, `SELECT user_badge_id, user_id, badge_id, source_id, earned_at FROM user_badges WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ubs []domain.UserBadge
	for rows.Next() {
		var ub domain.UserBadge
		if err := rows.Scan(&ub.UserBadgeID, &ub.UserID, &ub.BadgeID, &ub.SourceID, &ub.EarnedAt); err != nil {
			return nil, err
		}
		ubs = append(ubs, ub)
	}
	return ubs, nil
}

type creditScoreRepo struct {
	pool *pgxpool.Pool
}

func NewCreditScoreRepository(pool *pgxpool.Pool) port.CreditScoreRepository {
	return &creditScoreRepo{pool: pool}
}

func (r *creditScoreRepo) Create(ctx context.Context, c *domain.CreditScore) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO credit_scores (credit_id, user_id, current_score, last_updated) VALUES ($1, $2, $3, $4)`,
		c.CreditID, c.UserID, c.CurrentScore, c.LastUpdated)
	return err
}

func (r *creditScoreRepo) FindByUserID(ctx context.Context, userID string) (*domain.CreditScore, error) {
	var c domain.CreditScore
	err := r.pool.QueryRow(ctx, `SELECT credit_id, user_id, current_score, last_updated FROM credit_scores WHERE user_id = $1`, userID).Scan(&c.CreditID, &c.UserID, &c.CurrentScore, &c.LastUpdated)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &c, err
}

func (r *creditScoreRepo) Decrement(ctx context.Context, userID string, delta int) error {
	_, err := r.pool.Exec(ctx, `UPDATE credit_scores SET current_score = current_score - $1, last_updated = NOW() WHERE user_id = $2`, delta, userID)
	return err
}
