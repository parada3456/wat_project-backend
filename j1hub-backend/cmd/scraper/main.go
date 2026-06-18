package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/scraper"
)

func main() {
	s := scraper.NewScraper()

	results := make(chan scraper.Result)
	var wg sync.WaitGroup

	// Phase 2: Scrape Acadex
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Starting Acadex scraper...")
		// Using placeholder URLs, in a real scenario we'd use the actual ones
		s.ScrapeAcadex("https://www.acadexthailand.com/program/work-and-travel-summer/", results)
	}()

	// Phase 3: Scrape iHappy
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Starting iHappy scraper...")
		s.ScrapeIHappy("https://www.ihappyeducation.com/job-location-summer/", results)
	}()

	var jobs []domain.JobPosting
	var housings []domain.JobHousing

	// Phase 4: Output Generation
	done := make(chan bool)
	go func() {
		for res := range results {
			if res.Job != nil {
				jobs = append(jobs, *res.Job)
			}
			if res.Housing != nil {
				housings = append(housings, *res.Housing)
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

	if err := exportToJSON("job_housing.json", housings); err != nil {
		log.Fatalf("Failed to export housings: %v", err)
	}
	log.Printf("Exported %d job housings to job_housing.json", len(housings))

	log.Println("Scraping and export completed successfully.")
}

func exportToJSON(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
