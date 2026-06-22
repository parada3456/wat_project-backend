package usecase

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/j1hub/backend/internal/adapter/outbound/scraper"
	"github.com/j1hub/backend/internal/adapter/outbound/scraper/acadex"
	"github.com/j1hub/backend/internal/adapter/outbound/scraper/iee"
	"github.com/j1hub/backend/internal/adapter/outbound/scraper/ihappy"
	jobdomain "github.com/j1hub/backend/internal/job/domain"
	"github.com/j1hub/backend/internal/port"
)

type ScrapeJobsUseCase struct {
	jobRepo     port.JobPostingRepository
	housingRepo port.JobHousingRepository
}

func NewScrapeJobsUseCase(jobRepo port.JobPostingRepository, housingRepo port.JobHousingRepository) *ScrapeJobsUseCase {
	log.Println("debugprint: entering NewScrapeJobsUseCase")
	return &ScrapeJobsUseCase{
		jobRepo:     jobRepo,
		housingRepo: housingRepo,
	}
}

func (uc *ScrapeJobsUseCase) Run(ctx context.Context) error {
	log.Println("debugprint: entering (*ScrapeJobsUseCase).Run")

	acadexScraper := acadex.NewAcadexScraper(nil)
	ihappyScraper := ihappy.NewIHappyScraper(nil)
	ieeScraper := iee.NewIEEScraper(nil)

	results := make(chan *jobdomain.JobPosting, 100)
	var wg sync.WaitGroup

	scrapeSource := func(source scraper.JobSource, name string, url string) {
		defer wg.Done()
		log.Printf("Starting %s scraper job...", name)

		links, err := source.GetJobLinks(ctx, url)
		if err != nil {
			log.Printf("Failed to get job links for %s: %v", name, err)
			return
		}

		for _, link := range links {
			// small delay to prevent rate limit
			time.Sleep(200 * time.Millisecond)
			job, err := source.GetJobDetails(ctx, link)
			if err != nil {
				log.Printf("Failed to get job details for %s from %s: %v", name, link, err)
				continue
			}
			if job != nil {
				results <- job
			}
		}
	}

	// Scrape Acadex
	wg.Add(1)
	go scrapeSource(acadexScraper, "Acadex", "https://www.acadexthailand.com/program/work-and-travel-summer/")

	// Scrape iHappy
	wg.Add(1)
	go scrapeSource(ihappyScraper, "iHappy", "https://www.ihappyeducation.com/job-location-summer/")

	// Scrape IEE
	wg.Add(1)
	go scrapeSource(ieeScraper, "IEE", "https://www.ieethailand.com/work-and-travel-new/")

	go func() {
		wg.Wait()
		close(results)
	}()

	// Process and Upsert results
	for job := range results {
		if job != nil {
			log.Printf("Scraper found job: %s - %s", job.EmployerTitle, job.Position)
			if err := uc.jobRepo.Upsert(ctx, job); err != nil {
				log.Printf("Failed to upsert job posting %s: %v", job.JobID, err)
				continue
			}

			// If we had housing logic, we'd upsert it here.
			// Since our new interface focuses on JobPosting, we skip housing logic or create an empty one.
		}
	}

	return nil
}
