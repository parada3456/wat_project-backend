package scraper_test

import (
	"testing"
	"time"

	"github.com/j1hub/backend/internal/scraper"
	"github.com/stretchr/testify/assert"
)

func TestScraper_ScrapeAcadex(t *testing.T) {
	scr := scraper.NewScraper()
	results := make(chan scraper.Result, 10)

	listURL := "https://www.acadexthailand.com/location/thai-o-cha-ocean-city-summer-2027-group-a/"
	scr.ScrapeAcadex(listURL, results)

	select {
	case res := <-results:
		assert.Equal(t, "Thai O-Cha", res.Job.EmployerTitle)
		assert.Contains(t, res.Job.Position, "Server")
		assert.Contains(t, res.Job.Position, "Line cook")
		assert.Equal(t, "Ocean city", res.Job.LocationCity)
		assert.Equal(t, "Maryland", res.Job.LocationState)
		assert.Equal(t, "Rank S", res.Job.GroupLocation)
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for scraper results")
	}
}

func TestScraper_ScrapeIHappy(t *testing.T) {
	scr := scraper.NewScraper()
	results := make(chan scraper.Result, 10)

	listURL := "https://www.ihappyeducation.com/yankee-rebel-tavern-mackinac-island-michigan/"
	scr.ScrapeIHappy(listURL, results)

	select {
	case res := <-results:
		assert.Equal(t, "Yankee Rebel Tavern", res.Job.EmployerTitle)
		assert.Equal(t, "Host/Busser", res.Job.Position)
		assert.Equal(t, "Mackinac Island", res.Job.LocationCity)
		assert.Equal(t, "Michigan", res.Job.LocationState)
		assert.Equal(t, "Rank B", res.Job.GroupLocation)
		assert.Equal(t, 110.00, res.Housing.WeeklyRate)
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for scraper results")
	}
}

func TestScraper_ScrapeAcadexSilverDollarCity(t *testing.T) {
	scr := scraper.NewScraper()
	results := make(chan scraper.Result, 10)

	listURL := "https://www.acadexthailand.com/location/silver-dollar-city-branson-summer-2027-group-x/"
	scr.ScrapeAcadex(listURL, results)

	select {
	case res := <-results:
		assert.Equal(t, "Silver Dollar City", res.Job.EmployerTitle)
		assert.Contains(t, res.Job.Position, "Attractions Team Member")
		assert.Contains(t, res.Job.Position, "Food Team Member")
		assert.Equal(t, "Branson", res.Job.LocationCity)
		assert.Equal(t, "Missouri", res.Job.LocationState)
		assert.Equal(t, "Rank A", res.Job.GroupLocation)
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for scraper results")
	}
}

func TestScraper_ScrapeIHappyUnionRiver(t *testing.T) {
	scr := scraper.NewScraper()
	results := make(chan scraper.Result, 10)

	listURL := "https://www.ihappyeducation.com/union-river-lobster-pot-ellsworth-maine/"
	scr.ScrapeIHappy(listURL, results)

	select {
	case res := <-results:
		assert.Equal(t, "Union River Lobster Pot", res.Job.EmployerTitle)
		assert.Equal(t, "Back of House", res.Job.Position)
		assert.Equal(t, "Ellsworth", res.Job.LocationCity)
		assert.Equal(t, "Maine", res.Job.LocationState)
		assert.Equal(t, "Rank B", res.Job.GroupLocation)
		assert.Equal(t, 100.00, res.Housing.WeeklyRate)
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for scraper results")
	}
}
