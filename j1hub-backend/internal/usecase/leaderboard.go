package usecase

import (
	"context"
	"fmt"

	"github.com/j1hub/backend/internal/port"
)

type LeaderboardUseCase struct {
	leaderRepo  port.LeaderboardRepository
	profileRepo port.ProfileRepository
	ubRepo      port.UserBadgeRepository
}

func NewLeaderboardUseCase(leaderRepo port.LeaderboardRepository, profileRepo port.ProfileRepository, ubRepo port.UserBadgeRepository) *LeaderboardUseCase {
	return &LeaderboardUseCase{leaderRepo: leaderRepo, profileRepo: profileRepo, ubRepo: ubRepo}
}

type LeaderboardEntry struct {
	Rank   int
	UserID string
	Name   string
	Points int
	Streak int
	Badges []string
}

func (uc *LeaderboardUseCase) GetLeaderboard(ctx context.Context, scope, jobID string) ([]LeaderboardEntry, error) {
	users, err := uc.leaderRepo.FindByScope(ctx, scope, jobID)
	if err != nil {
		return nil, err
	}

	var results []LeaderboardEntry
	for i, u := range users {
		p, _ := uc.profileRepo.FindByUserID(ctx, u.UserID)

		lastNameInitial := ""
		if len(u.LastName) > 0 {
			lastNameInitial = " " + string(u.LastName[0]) + "."
		}
		name := fmt.Sprintf("%s%s", u.FirstName, lastNameInitial)
		if p != nil && p.RadarVisibility == "Hidden" {
			name = fmt.Sprintf("J1 Student #%d", i+1)
		}

		entry := LeaderboardEntry{
			Rank:   i + 1,
			UserID: u.UserID,
			Name:   name,
			Points: u.CurrentPhasePoints,
			Streak: u.MissionStreak,
		}

		// Load badges
		ubs, _ := uc.ubRepo.FindByUser(ctx, u.UserID)
		for _, ub := range ubs {
			entry.Badges = append(entry.Badges, ub.BadgeID)
		}

		results = append(results, entry)
	}

	return results, nil
}
