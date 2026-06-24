package postgres

import (
	"context"
	"log"

	userdomain "github.com/j1hub/backend/internal/user/domain"

	port "github.com/j1hub/backend/internal/gamification/port"
	"github.com/jackc/pgx/v5/pgxpool"
)

type radarRepo struct {
	pool *pgxpool.Pool
}

func NewRadarRepository(pool *pgxpool.Pool) port.RadarRepository {
	log.Println("debugprint: entering NewRadarRepository")
	return &radarRepo{pool: pool}
}

func (r *radarRepo) FindNearby(ctx context.Context, lat, lng, radius float64, staleMinutes int) ([]userdomain.Profile, error) {
	log.Println("debugprint: entering (*radarRepo).FindNearby")
	query := `
		SELECT 
			profile_id, user_id, phone_number, bio, avatar_url, 
			radar_visibility, ST_X(current_coordinates), ST_Y(current_coordinates), 
			location_updated_at, updated_at
		FROM profiles 
		WHERE ST_DWithin(
			current_coordinates::geography, 
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, 
			$3
		)
		AND location_updated_at > NOW() - (interval '1 minute' * $4)`

	rows, err := r.pool.Query(ctx, query, lng, lat, radius, staleMinutes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var profiles []userdomain.Profile
	for rows.Next() {
		var p userdomain.Profile
		if err := rows.Scan(&p.ProfileID, &p.UserID, &p.PhoneNumber, &p.Bio, &p.AvatarURL, &p.RadarVisibility, &p.Lng, &p.Lat, &p.LocationUpdatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}
