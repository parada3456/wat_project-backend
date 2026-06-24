package gamificationusecase

import (
	"context"
	"log"

	frienddomain "github.com/j1hub/backend/internal/friend/domain"
	userdomain "github.com/j1hub/backend/internal/user/domain"

	"github.com/j1hub/backend/internal/infrastructure/config"
	port "github.com/j1hub/backend/internal/gamification/port"
)

type RadarUseCase struct {
	cfg         *config.Config
	profileRepo port.ProfileRepository
	radarRepo   port.RadarRepository
	friendRepo  port.FriendshipRepository
}

func NewRadarUseCase(cfg *config.Config, profileRepo port.ProfileRepository, radarRepo port.RadarRepository, friendRepo port.FriendshipRepository) *RadarUseCase {
	log.Println("debugprint: entering NewRadarUseCase")
	return &RadarUseCase{cfg: cfg, profileRepo: profileRepo, radarRepo: radarRepo, friendRepo: friendRepo}
}

type RadarEntry struct {
	UserID    string
	Name      string
	AvatarURL string
	Lat       float64
	Lng       float64
}

func (uc *RadarUseCase) GetRadar(ctx context.Context, requesterID string) ([]RadarEntry, error) {
	log.Println("debugprint: entering (*RadarUseCase).GetRadar")
	p, err := uc.profileRepo.FindByUserID(ctx, requesterID)
	if err != nil {
		return nil, err
	}

	nearby, err := uc.radarRepo.FindNearby(ctx, p.Lat, p.Lng, uc.cfg.RadarRadiusMeters, uc.cfg.RadarStaleMinutes)
	if err != nil {
		return nil, err
	}

	var results []RadarEntry
	for _, n := range nearby {
		if n.UserID == requesterID {
			continue
		}

		if n.RadarVisibility == userdomain.VisibilityHidden {
			continue
		}

		entry := RadarEntry{
			UserID:    n.UserID,
			Lat:       n.Lat,
			Lng:       n.Lng,
			AvatarURL: "anonymous_avatar.png",
			Name:      "J1 Student",
		}

		isFriend := false
		u1, u2 := frienddomain.CanonicalOrder(requesterID, n.UserID)
		f, err := uc.friendRepo.FindByCanonicalPair(ctx, u1, u2)
		if err == nil && f.Status == frienddomain.FriendshipAccepted {
			isFriend = true
		}

		if n.RadarVisibility == userdomain.VisibilityShowFriends && !isFriend {
			continue
		}

		if isFriend || n.RadarVisibility == userdomain.VisibilityShowAnonymous {
			// In real case we'd fetch user name
			entry.Name = "Real Name" // placeholder
			entry.AvatarURL = n.AvatarURL
		}

		results = append(results, entry)
	}

	return results, nil
}
