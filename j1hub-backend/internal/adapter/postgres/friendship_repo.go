package postgres

import (
	"context"
	"log"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type friendshipRepo struct {
	pool *pgxpool.Pool
}

func NewFriendshipRepository(pool *pgxpool.Pool) port.FriendshipRepository {
	log.Println("debugprint: entering NewFriendshipRepository")
	return &friendshipRepo{pool: pool}
}

func (r *friendshipRepo) Insert(ctx context.Context, f *domain.Friendship) error {
	log.Println("debugprint: entering (*friendshipRepo).Insert")
	_, err := r.pool.Exec(ctx, `INSERT INTO friendships (friendship_id, user_id_1, user_id_2, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		f.FriendshipID, f.UserID1, f.UserID2, f.Status, f.CreatedAt, f.UpdatedAt)
	return err
}

func (r *friendshipRepo) FindByCanonicalPair(ctx context.Context, u1, u2 string) (*domain.Friendship, error) {
	log.Println("debugprint: entering (*friendshipRepo).FindByCanonicalPair")
	var f domain.Friendship
	err := r.pool.QueryRow(ctx, `SELECT friendship_id, user_id_1, user_id_2, status, created_at, updated_at FROM friendships WHERE user_id_1 = $1 AND user_id_2 = $2`, u1, u2).Scan(&f.FriendshipID, &f.UserID1, &f.UserID2, &f.Status, &f.CreatedAt, &f.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &f, err
}

func (r *friendshipRepo) FindByID(ctx context.Context, id string) (*domain.Friendship, error) {
	log.Println("debugprint: entering (*friendshipRepo).FindByID")
	var f domain.Friendship
	err := r.pool.QueryRow(ctx, `SELECT friendship_id, user_id_1, user_id_2, status, created_at, updated_at FROM friendships WHERE friendship_id = $1`, id).Scan(&f.FriendshipID, &f.UserID1, &f.UserID2, &f.Status, &f.CreatedAt, &f.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &f, err
}

func (r *friendshipRepo) UpdateStatus(ctx context.Context, id string, status domain.FriendshipStatus) error {
	log.Println("debugprint: entering (*friendshipRepo).UpdateStatus")
	_, err := r.pool.Exec(ctx, `UPDATE friendships SET status = $1, updated_at = NOW() WHERE friendship_id = $2`, status, id)
	return err
}

func (r *friendshipRepo) FindFriendsOf(ctx context.Context, userID string) ([]domain.Friendship, error) {
	log.Println("debugprint: entering (*friendshipRepo).FindFriendsOf")
	rows, err := r.pool.Query(ctx, `SELECT friendship_id, user_id_1, user_id_2, status, created_at, updated_at FROM friendships WHERE (user_id_1 = $1 OR user_id_2 = $1) AND status = 'Accepted'`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var friends []domain.Friendship
	for rows.Next() {
		var f domain.Friendship
		if err := rows.Scan(&f.FriendshipID, &f.UserID1, &f.UserID2, &f.Status, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		friends = append(friends, f)
	}
	return friends, nil
}
