package scraper

import (
	"log"

	"github.com/gocolly/colly/v2"
	"github.com/j1hub/backend/internal/domain"
)

type Scraper struct {
	Collector *colly.Collector
}

func NewScraper() *Scraper {
	log.Println("debugprint: entering NewScraper")
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       2,
	})
	return &Scraper{
		Collector: c,
	}
}

type Result struct {
	Job     *domain.JobPosting
	Housing *domain.JobHousing
}
