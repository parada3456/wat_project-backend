package postgres

import (
	"context"
	"fmt"
	"log"

	notificationdomain "github.com/parada3456/wat_project-backend/internal/notification/domain"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	port "github.com/parada3456/wat_project-backend/internal/notification/port"
)

type notificationRepository struct {
	pool *pgxpool.Pool
}

func NewNotificationRepository(pool *pgxpool.Pool) port.NotificationRepository {
	log.Println("debugprint: entering NewNotificationRepository")
	return &notificationRepository{pool: pool}
}

func (r *notificationRepository) Insert(ctx context.Context, n *notificationdomain.Notification) error {
	log.Println("debugprint: entering (*notificationRepository).Insert")
	query := `
		INSERT INTO notifications (id, user_id, title, body, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, query, n.NotificationID, n.UserID, n.Title, n.Body, n.IsRead, n.CreatedAt)
	return err
}

func (r *notificationRepository) FindByUser(
	ctx context.Context,
	userID string,
	isRead *bool,
	limit,
	offset int,
) ([]notificationdomain.Notification, int, error) {
	log.Println("debugprint: entering (*notificationRepository).FindByUser")

	// ----------------------------------------------------
	// STEP 1: Build & Execute Count Query
	// ----------------------------------------------------
	countBuilder := sq.Select("COUNT(*)").
		From("notifications").
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar) // Essential for PostgreSQL ($1, $2)

	if isRead != nil {
		countBuilder = countBuilder.Where(sq.Eq{"is_read": *isRead})
	}

	countSQL, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build count query: %w", err)
	}

	var totalCount int
	// Use r.pool directly with pgx native QueryRow method
	err = r.pool.QueryRow(ctx, countSQL, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute count query: %w", err)
	}

	// Early exit optimization if no rows match criteria
	if totalCount == 0 {
		return []notificationdomain.Notification{}, 0, nil
	}

	// ----------------------------------------------------
	// STEP 2: Build & Execute Paginated Data Query
	// ----------------------------------------------------
	dataBuilder := sq.Select("id", "user_id", "title", "body", "is_read", "created_at").
		From("notifications").
		Where(sq.Eq{"user_id": userID}).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		OrderBy("created_at DESC"). // Good practice for lists
		PlaceholderFormat(sq.Dollar)

	if isRead != nil {
		dataBuilder = dataBuilder.Where(sq.Eq{"is_read": *isRead})
	}

	dataSQL, dataArgs, err := dataBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build data query: %w", err)
	}

	// Use r.pool directly with pgx native Query method
	rows, err := r.pool.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute data query: %w", err)
	}
	defer rows.Close() // pgx rows MUST be closed to return the connection back to the pool

	var notifications []notificationdomain.Notification
	for rows.Next() {
		var n notificationdomain.Notification
		// pgx automatically converts underlying PG types to Go types elegantly
		err := rows.Scan(&n.NotificationID, &n.UserID, &n.Title, &n.Body, &n.IsRead, &n.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan notification row: %w", err)
		}
		notifications = append(notifications, n)
	}

	// Check for any async errors during row iteration
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during row iteration: %w", err)
	}

	return notifications, totalCount, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id string) error {
	log.Println("debugprint: entering (*notificationRepository).MarkAsRead")
	query := `UPDATE notifications SET is_read = TRUE WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	log.Println("debugprint: entering (*notificationRepository).MarkAllAsRead")
	query := `UPDATE notifications SET is_read = TRUE WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}

func (r *notificationRepository) Delete(ctx context.Context, id string) error {
	log.Println("debugprint: entering (*notificationRepository).Delete")
	query := `DELETE FROM notifications WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
