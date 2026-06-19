package scheduler

import (
	"context"
	"log"

	"github.com/j1hub/backend/internal/infrastructure/config"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/robfig/cron/v3"
)

func NewScheduler(
	cfg *config.Config,
	overdueExpenseJob *usecase.OverdueExpenseJob,
	overdueMissionJob *usecase.OverdueMissionJob,
	scrapeJobsJob *usecase.ScrapeJobsUseCase,
) *cron.Cron {
	log.Println("debugprint: entering NewScheduler")
	c := cron.New()

	c.AddFunc(cfg.CronOverdueExpense, func() {
		log.Println("Running overdue expense job...")
		if err := overdueExpenseJob.Run(context.Background()); err != nil {
			log.Printf("overdue expense job failed: %v", err)
		}
	})

	c.AddFunc(cfg.CronOverdueMission, func() {
		log.Println("Running overdue mission job...")
		if err := overdueMissionJob.Run(context.Background()); err != nil {
			log.Printf("overdue mission job failed: %v", err)
		}
	})

	c.AddFunc(cfg.CronScraper, func() {
		log.Println("Running weekly jobs scraping cron...")
		if err := scrapeJobsJob.Run(context.Background()); err != nil {
			log.Printf("scraping cron failed: %v", err)
		}
	})

	return c
}
