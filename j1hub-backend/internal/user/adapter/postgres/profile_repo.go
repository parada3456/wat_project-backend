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

type profileRepo struct {
	pool *pgxpool.Pool
}

func NewProfileRepository(pool *pgxpool.Pool) port.ProfileRepository {
	log.Println("debugprint: entering NewProfileRepository")
	return &profileRepo{pool: pool}
}

func (r *profileRepo) Create(ctx context.Context, p *userdomain.Profile) error {
	log.Println("debugprint: entering (*profileRepo).Create")
	query := `
		INSERT INTO profiles (
			profile_id, user_id, phone_number, bio, avatar_url, 
			radar_visibility, current_coordinates, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, ST_SetSRID(ST_MakePoint($7, $8), 4326), $9)`

	_, err := r.pool.Exec(ctx, query,
		p.ProfileID, p.UserID, p.PhoneNumber, p.Bio, p.AvatarURL,
		p.RadarVisibility, p.Lng, p.Lat, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}
	return nil
}

func (r *profileRepo) FindByUserID(ctx context.Context, userID string) (*userdomain.Profile, error) {
	log.Println("debugprint: entering (*profileRepo).FindByUserID")
	query := `
		SELECT 
			profile_id, user_id, phone_number, bio, avatar_url, 
			radar_visibility, ST_X(current_coordinates), ST_Y(current_coordinates), 
			location_updated_at, updated_at
		FROM profiles WHERE user_id = $1`

	row := r.pool.QueryRow(ctx, query, userID)
	var p userdomain.Profile
	var locUpdated *time.Time
	err := row.Scan(
		&p.ProfileID, &p.UserID, &p.PhoneNumber, &p.Bio, &p.AvatarURL,
		&p.RadarVisibility, &p.Lng, &p.Lat, &locUpdated, &p.UpdatedAt,
	)
	if err == nil && locUpdated != nil {
		p.LocationUpdatedAt = *locUpdated
	}
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find profile by user id: %w", err)
	}
	return &p, nil
}

func (r *profileRepo) Update(ctx context.Context, p *userdomain.Profile) error {
	log.Println("debugprint: entering (*profileRepo).Update")
	query := `
		UPDATE profiles SET 
			phone_number = $1, bio = $2, avatar_url = $3, 
			radar_visibility = $4, updated_at = NOW()
		WHERE profile_id = $5`

	_, err := r.pool.Exec(ctx, query,
		p.PhoneNumber, p.Bio, p.AvatarURL, p.RadarVisibility, p.ProfileID,
	)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}
	return nil
}

func (r *profileRepo) UpdateLocation(ctx context.Context, userID string, lat, lng float64) error {
	log.Println("debugprint: entering (*profileRepo).UpdateLocation")
	query := `
		UPDATE profiles SET 
			current_coordinates = ST_SetSRID(ST_MakePoint($1, $2), 4326),
			location_updated_at = NOW(),
			updated_at = NOW()
		WHERE user_id = $3`

	_, err := r.pool.Exec(ctx, query, lng, lat, userID)
	if err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}
	return nil
}

func (r *profileRepo) UpdateVisibility(ctx context.Context, userID string, visibility userdomain.RadarVisibility) error {
	log.Println("debugprint: entering (*profileRepo).UpdateVisibility")
	query := `UPDATE profiles SET radar_visibility = $1, updated_at = NOW() WHERE user_id = $2`
	_, err := r.pool.Exec(ctx, query, visibility, userID)
	if err != nil {
		return fmt.Errorf("failed to update visibility: %w", err)
	}
	return nil
}
