package usecase

import (
	"context"
	"log"
	"sync"

	"github.com/j1hub/backend/internal/port"
	"github.com/j1hub/backend/internal/scraper"
)

type ScrapeJobsUseCase struct {
	jobRepo     port.JobPostingRepository
	housingRepo port.JobHousingRepository
}

func NewScrapeJobsUseCase(jobRepo port.JobPostingRepository, housingRepo port.JobHousingRepository) *ScrapeJobsUseCase {
	return &ScrapeJobsUseCase{
		jobRepo:     jobRepo,
		housingRepo: housingRepo,
	}
}

func (uc *ScrapeJobsUseCase) Run(ctx context.Context) error {
	scr := scraper.NewScraper()
	results := make(chan scraper.Result, 100)
	var wg sync.WaitGroup

	// Scrape Acadex
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Starting Acadex scraper job...")
		scr.ScrapeAcadex("https://www.acadexthailand.com/program/work-and-travel-summer/", results)
	}()

	// Scrape iHappy
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Starting iHappy scraper job...")
		scr.ScrapeIHappy("https://www.ihappyeducation.com/job-location-summer/", results)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	// Process and Upsert results
	for res := range results {
		if res.Job != nil {
			log.Printf("Scraper found job: %s - %s", res.Job.EmployerTitle, res.Job.Position)
			if err := uc.jobRepo.Upsert(ctx, res.Job); err != nil {
				log.Printf("Failed to upsert job posting %s: %v", res.Job.JobID, err)
				continue
			}
			if res.Housing != nil {
				if err := uc.housingRepo.Upsert(ctx, res.Housing); err != nil {
					log.Printf("Failed to upsert job housing %s: %v", res.Housing.HousingID, err)
				}
			}
		}
	}

	return nil
}
