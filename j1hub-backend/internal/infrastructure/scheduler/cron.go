package scheduler

import (
	"context"
	"log"

	"github.com/parada3456/wat_project-backend/internal/infrastructure/config"
	scraper "github.com/parada3456/wat_project-backend/internal/infrastructure/outbound/scraper"
	"github.com/robfig/cron/v3"
)

func NewScheduler(
	cfg *config.Config,
	overdueExpenseJob *OverdueExpenseJob,
	overdueMissionJob *OverdueMissionJob,
	scrapeJobsJob *scraper.ScrapeJobsUseCase,
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
