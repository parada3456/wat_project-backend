package postgres

import (
	"context"
	"log"

	gamificationdomain "github.com/parada3456/wat_project-backend/internal/gamification/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parada3456/wat_project-backend/internal/domain"
	port "github.com/parada3456/wat_project-backend/internal/gamification/port"
)

type pointLedgerRepo struct {
	pool *pgxpool.Pool
}

func NewPointLedgerRepository(pool *pgxpool.Pool) port.PointLedgerRepository {
	log.Println("debugprint: entering NewPointLedgerRepository")
	return &pointLedgerRepo{pool: pool}
}

func (r *pointLedgerRepo) Insert(ctx context.Context, l *gamificationdomain.PointLedger) error {
	log.Println("debugprint: entering (*pointLedgerRepo).Insert")
	query := `INSERT INTO point_ledger (ledger_id, user_id, source_type, source_id, delta, lifetime_balance_after, phase_balance_after, note, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pool.Exec(ctx, query, l.LedgerID, l.UserID, l.SourceType, l.SourceID, l.Delta, l.LifetimeBalanceAfter, l.PhaseBalanceAfter, l.Note, l.CreatedAt)
	return err
}

func (r *pointLedgerRepo) InsertBatch(ctx context.Context, ledgers []gamificationdomain.PointLedger) error {
	log.Println("debugprint: entering (*pointLedgerRepo).InsertBatch")
	batch := &pgx.Batch{}
	for _, l := range ledgers {
		batch.Queue(`INSERT INTO point_ledger (ledger_id, user_id, source_type, source_id, delta, lifetime_balance_after, phase_balance_after, note, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			l.LedgerID, l.UserID, l.SourceType, l.SourceID, l.Delta, l.LifetimeBalanceAfter, l.PhaseBalanceAfter, l.Note, l.CreatedAt)
	}
	return r.pool.SendBatch(ctx, batch).Close()
}

func (r *pointLedgerRepo) FindByUser(ctx context.Context, userID string) ([]gamificationdomain.PointLedger, error) {
	log.Println("debugprint: entering (*pointLedgerRepo).FindByUser")
	query := `SELECT ledger_id, user_id, source_type, source_id, delta, lifetime_balance_after, phase_balance_after, note, created_at FROM point_ledger WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ledgers []gamificationdomain.PointLedger
	for rows.Next() {
		var l gamificationdomain.PointLedger
		if err := rows.Scan(&l.LedgerID, &l.UserID, &l.SourceType, &l.SourceID, &l.Delta, &l.LifetimeBalanceAfter, &l.PhaseBalanceAfter, &l.Note, &l.CreatedAt); err != nil {
			return nil, err
		}
		ledgers = append(ledgers, l)
	}
	return ledgers, nil
}

func (r *pointLedgerRepo) FindByUserAndSourceType(ctx context.Context, userID string, sourceType gamificationdomain.SourceType) ([]gamificationdomain.PointLedger, error) {
	log.Println("debugprint: entering (*pointLedgerRepo).FindByUserAndSourceType")
	query := `SELECT ledger_id, user_id, source_type, source_id, delta, lifetime_balance_after, phase_balance_after, note, created_at FROM point_ledger WHERE user_id = $1 AND source_type = $2 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, userID, sourceType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ledgers []gamificationdomain.PointLedger
	for rows.Next() {
		var l gamificationdomain.PointLedger
		if err := rows.Scan(&l.LedgerID, &l.UserID, &l.SourceType, &l.SourceID, &l.Delta, &l.LifetimeBalanceAfter, &l.PhaseBalanceAfter, &l.Note, &l.CreatedAt); err != nil {
			return nil, err
		}
		ledgers = append(ledgers, l)
	}
	return ledgers, nil
}

type badgeRepo struct {
	pool *pgxpool.Pool
}

func NewBadgeRepository(pool *pgxpool.Pool) port.BadgeRepository {
	log.Println("debugprint: entering NewBadgeRepository")
	return &badgeRepo{pool: pool}
}

func (r *badgeRepo) FindByTriggerType(ctx context.Context, triggerType gamificationdomain.TriggerType) ([]gamificationdomain.Badge, error) {
	log.Println("debugprint: entering (*badgeRepo).FindByTriggerType")
	rows, err := r.pool.Query(ctx, `SELECT badge_id, title, description, trigger_type, icon_url, created_at FROM badges WHERE trigger_type = $1`, triggerType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var badges []gamificationdomain.Badge
	for rows.Next() {
		var b gamificationdomain.Badge
		if err := rows.Scan(&b.BadgeID, &b.Title, &b.Description, &b.TriggerType, &b.IconURL, &b.CreatedAt); err != nil {
			return nil, err
		}
		badges = append(badges, b)
	}
	return badges, nil
}

func (r *badgeRepo) FindEligible(ctx context.Context, userID string, triggerType gamificationdomain.TriggerType) ([]gamificationdomain.Badge, error) {
	log.
		// Simple implementation: find all badges of this trigger type that the user doesn't have yet
		Println("debugprint: entering (*badgeRepo).FindEligible")

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
	var badges []gamificationdomain.Badge
	for rows.Next() {
		var b gamificationdomain.Badge
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
	log.Println("debugprint: entering NewUserBadgeRepository")
	return &userBadgeRepo{pool: pool}
}

func (r *userBadgeRepo) Insert(ctx context.Context, ub *gamificationdomain.UserBadge) error {
	log.Println("debugprint: entering (*userBadgeRepo).Insert")
	_, err := r.pool.Exec(ctx, `INSERT INTO user_badges (user_badge_id, user_id, badge_id, source_id, earned_at) VALUES ($1, $2, $3, $4, $5)`,
		ub.UserBadgeID, ub.UserID, ub.BadgeID, ub.SourceID, ub.EarnedAt)
	return err
}

func (r *userBadgeRepo) FindByUser(ctx context.Context, userID string) ([]gamificationdomain.UserBadge, error) {
	log.Println("debugprint: entering (*userBadgeRepo).FindByUser")
	query := `
		SELECT ub.user_badge_id, ub.user_id, ub.badge_id, ub.source_id, ub.earned_at,
		       b.badge_id, b.title, b.description, b.trigger_type, b.icon_url, b.created_at
		FROM user_badges ub
		INNER JOIN badges b ON ub.badge_id = b.badge_id
		WHERE ub.user_id = $1
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ubs []gamificationdomain.UserBadge
	for rows.Next() {
		var ub gamificationdomain.UserBadge
		var b gamificationdomain.Badge
		if err := rows.Scan(
			&ub.UserBadgeID, &ub.UserID, &ub.BadgeID, &ub.SourceID, &ub.EarnedAt,
			&b.BadgeID, &b.Title, &b.Description, &b.TriggerType, &b.IconURL, &b.CreatedAt,
		); err != nil {
			return nil, err
		}
		ub.Badge = &b
		ubs = append(ubs, ub)
	}
	return ubs, nil
}

type creditScoreRepo struct {
	pool *pgxpool.Pool
}

func NewCreditScoreRepository(pool *pgxpool.Pool) port.CreditScoreRepository {
	log.Println("debugprint: entering NewCreditScoreRepository")
	return &creditScoreRepo{pool: pool}
}

func (r *creditScoreRepo) Create(ctx context.Context, c *gamificationdomain.CreditScore) error {
	log.Println("debugprint: entering (*creditScoreRepo).Create")
	_, err := r.pool.Exec(ctx, `INSERT INTO credit_scores (credit_id, user_id, current_score, last_updated) VALUES ($1, $2, $3, $4)`,
		c.CreditID, c.UserID, c.CurrentScore, c.LastUpdated)
	return err
}

func (r *creditScoreRepo) FindByUserID(ctx context.Context, userID string) (*gamificationdomain.CreditScore, error) {
	log.Println("debugprint: entering (*creditScoreRepo).FindByUserID")
	var c gamificationdomain.CreditScore
	err := r.pool.QueryRow(ctx, `SELECT credit_id, user_id, current_score, last_updated FROM credit_scores WHERE user_id = $1`, userID).Scan(&c.CreditID, &c.UserID, &c.CurrentScore, &c.LastUpdated)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &c, err
}

func (r *creditScoreRepo) Decrement(ctx context.Context, userID string, delta int) error {
	log.Println("debugprint: entering (*creditScoreRepo).Decrement")
	_, err := r.pool.Exec(ctx, `UPDATE credit_scores SET current_score = current_score - $1, last_updated = NOW() WHERE user_id = $2`, delta, userID)
	return err
}
