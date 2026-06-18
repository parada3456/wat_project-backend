package usecase

import (
	"context"
	"log"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
)

type OverdueMissionJob struct {
	umRepo      port.UserMissionRepository
	missionRepo port.MissionRepository
	userRepo    port.UserRepository
	notifier    port.NotifierPort
}

func NewOverdueMissionJob(umRepo port.UserMissionRepository, missionRepo port.MissionRepository, userRepo port.UserRepository, notifier port.NotifierPort) *OverdueMissionJob {
	return &OverdueMissionJob{umRepo: umRepo, missionRepo: missionRepo, userRepo: userRepo, notifier: notifier}
}

func (j *OverdueMissionJob) Run(ctx context.Context) error {
	ums, err := j.umRepo.FindOverdue(ctx)
	if err != nil {
		return err
	}

	count := 0
	for _, um := range ums {
		if err := j.umRepo.UpdateStatus(ctx, um.UserMissionID, domain.StatusOverdue); err != nil {
			continue
		}

		m, err := j.missionRepo.FindByID(ctx, um.MissionID)
		if err == nil && m.IsMandatory {
			j.userRepo.ResetStreak(ctx, um.UserID)
		}

		j.notifier.Send(ctx, um.UserID, "Mission overdue", "A mission is past its due date!")
		count++
	}

	log.Printf("Processed %d overdue missions", count)
	return nil
}
