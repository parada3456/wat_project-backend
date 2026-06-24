package scheduler

import (
	"context"
	"log"

	missiondomain "github.com/j1hub/backend/internal/mission/domain"
	missionport "github.com/j1hub/backend/internal/mission/port"
	notificationport "github.com/j1hub/backend/internal/notification/port"
	userport "github.com/j1hub/backend/internal/user/port"
)

type OverdueMissionJob struct {
	umRepo      missionport.UserMissionRepository
	missionRepo missionport.MissionRepository
	userRepo    userport.UserRepository
	notifier    notificationport.NotifierPort
}

func NewOverdueMissionJob(umRepo missionport.UserMissionRepository, missionRepo missionport.MissionRepository, userRepo userport.UserRepository, notifier notificationport.NotifierPort) *OverdueMissionJob {
	log.Println("debugprint: entering NewOverdueMissionJob")
	return &OverdueMissionJob{umRepo: umRepo, missionRepo: missionRepo, userRepo: userRepo, notifier: notifier}
}

func (j *OverdueMissionJob) Run(ctx context.Context) error {
	log.Println("debugprint: entering (*OverdueMissionJob).Run")
	ums, err := j.umRepo.FindOverdue(ctx)
	if err != nil {
		return err
	}

	count := 0
	for _, um := range ums {
		if err := j.umRepo.UpdateStatus(ctx, um.UserMissionID, missiondomain.StatusOverdue); err != nil {
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
