package postgres

import (
	"context"
	"log"

	frienddomain "github.com/j1hub/backend/internal/friend/domain"

	"github.com/j1hub/backend/internal/domain"
	port "github.com/j1hub/backend/internal/friend/port"
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

func (r *friendshipRepo) Insert(ctx context.Context, f *frienddomain.Friendship) error {
	log.Println("debugprint: entering (*friendshipRepo).Insert")
	_, err := r.pool.Exec(ctx, `INSERT INTO friendships (friendship_id, user_id_1, user_id_2, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		f.FriendshipID, f.UserID1, f.UserID2, f.Status, f.CreatedAt, f.UpdatedAt)
	return err
}

func (r *friendshipRepo) FindByCanonicalPair(ctx context.Context, u1, u2 string) (*frienddomain.Friendship, error) {
	log.Println("debugprint: entering (*friendshipRepo).FindByCanonicalPair")
	var f frienddomain.Friendship
	err := r.pool.QueryRow(ctx, `SELECT friendship_id, user_id_1, user_id_2, status, created_at, updated_at FROM friendships WHERE user_id_1 = $1 AND user_id_2 = $2`, u1, u2).Scan(&f.FriendshipID, &f.UserID1, &f.UserID2, &f.Status, &f.CreatedAt, &f.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &f, err
}

func (r *friendshipRepo) FindByID(ctx context.Context, id string) (*frienddomain.Friendship, error) {
	log.Println("debugprint: entering (*friendshipRepo).FindByID")
	var f frienddomain.Friendship
	err := r.pool.QueryRow(ctx, `SELECT friendship_id, user_id_1, user_id_2, status, created_at, updated_at FROM friendships WHERE friendship_id = $1`, id).Scan(&f.FriendshipID, &f.UserID1, &f.UserID2, &f.Status, &f.CreatedAt, &f.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &f, err
}

func (r *friendshipRepo) UpdateStatus(ctx context.Context, id string, status frienddomain.FriendshipStatus) error {
	log.Println("debugprint: entering (*friendshipRepo).UpdateStatus")
	_, err := r.pool.Exec(ctx, `UPDATE friendships SET status = $1, updated_at = NOW() WHERE friendship_id = $2`, status, id)
	return err
}

func (r *friendshipRepo) FindFriendsOf(ctx context.Context, userID string, limit, offset int) ([]frienddomain.Friendship, int, error) {
	log.Println("debugprint: entering (*friendshipRepo).FindFriendsOf")
	
	var totalCount int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM friendships WHERE (user_id_1 = $1 OR user_id_2 = $1) AND status = 'accepted'`, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}
	
	if totalCount == 0 {
		return []frienddomain.Friendship{}, 0, nil
	}

	rows, err := r.pool.Query(ctx, `SELECT friendship_id, user_id_1, user_id_2, status, created_at, updated_at FROM friendships WHERE (user_id_1 = $1 OR user_id_2 = $1) AND status = 'accepted' LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var friends []frienddomain.Friendship
	for rows.Next() {
		var f frienddomain.Friendship
		if err := rows.Scan(&f.FriendshipID, &f.UserID1, &f.UserID2, &f.Status, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, 0, err
		}
		friends = append(friends, f)
	}
	return friends, totalCount, nil
}

func (r *friendshipRepo) FindPendingFor(ctx context.Context, userID string, limit, offset int) ([]frienddomain.Friendship, int, error) {
	log.Println("debugprint: entering (*friendshipRepo).FindPendingFor")
	
	var totalCount int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM friendships WHERE (user_id_1 = $1 OR user_id_2 = $1) AND status = 'pending'`, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	if totalCount == 0 {
		return []frienddomain.Friendship{}, 0, nil
	}

	rows, err := r.pool.Query(ctx, `SELECT friendship_id, user_id_1, user_id_2, status, created_at, updated_at FROM friendships WHERE (user_id_1 = $1 OR user_id_2 = $1) AND status = 'pending' LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var pending []frienddomain.Friendship
	for rows.Next() {
		var f frienddomain.Friendship
		if err := rows.Scan(&f.FriendshipID, &f.UserID1, &f.UserID2, &f.Status, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, 0, err
		}
		pending = append(pending, f)
	}
	return pending, totalCount, nil
}

func (r *friendshipRepo) Delete(ctx context.Context, id string) error {
	log.Println("debugprint: entering (*friendshipRepo).Delete")
	_, err := r.pool.Exec(ctx, `DELETE FROM friendships WHERE friendship_id = $1`, id)
	return err
}
