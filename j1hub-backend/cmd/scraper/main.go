package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/j1hub/backend/internal/infrastructure/outbound/scraper"
	"github.com/j1hub/backend/internal/infrastructure/outbound/scraper/acadex"
	"github.com/j1hub/backend/internal/infrastructure/outbound/scraper/iee"
	"github.com/j1hub/backend/internal/infrastructure/outbound/scraper/ihappy"
	jobdomain "github.com/j1hub/backend/internal/job/domain"
)

func main() {
	log.Println("debugprint: entering main")

	acadexScraper := acadex.NewAcadexScraper(nil)
	ihappyScraper := ihappy.NewIHappyScraper(nil)
	ieeScraper := iee.NewIEEScraper(nil)

	results := make(chan *jobdomain.JobPosting)
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scrapeSource := func(source scraper.JobSource, name string, url string) {
		defer wg.Done()
		log.Printf("Starting %s scraper for %s...", name, url)

		links, err := source.GetJobLinks(ctx, url)
		if err != nil {
			log.Printf("Failed to get job links for %s: %v", name, err)
			return
		}
		log.Printf("Found %d links for %s", len(links), name)

		for _, link := range links {
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

	// Phase 2: Scrape Acadex
	wg.Add(1)
	go scrapeSource(acadexScraper, "Acadex", "https://www.acadexthailand.com/program/work-and-travel-summer/")

	// Phase 3: Scrape iHappy
	wg.Add(1)
	go scrapeSource(ihappyScraper, "iHappy", "https://www.ihappyeducation.com/job-location-summer/")

	// Phase 4: Scrape IEE
	wg.Add(1)
	go scrapeSource(ieeScraper, "IEE", "https://www.ieethailand.com/work-and-travel-new/")

	var jobs []jobdomain.JobPosting

	// Phase 5: Output Generation
	done := make(chan bool)
	go func() {
		for job := range results {
			if job != nil {
				jobs = append(jobs, *job)
			}
		}
		done <- true
	}()

	wg.Wait()
	close(results)
	<-done

	// Export to JSON
	if err := exportToJSON("job_posting.json", jobs); err != nil {
		log.Fatalf("Failed to export jobs: %v", err)
	}
	log.Printf("Exported %d job postings to job_posting.json", len(jobs))

	// The old code exported housings, but we didn't specify JobHousing in the new interface.
	// You can add logic to extract housing from JobPosting or create them as separate entities.
	log.Println("Scraping and export completed successfully.")
}

func exportToJSON(filename string, data interface{}) error {
	log.Println("debugprint: entering exportToJSON")
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
