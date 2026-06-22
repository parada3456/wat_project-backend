package ihappy_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/j1hub/backend/internal/adapter/outbound/scraper/ihappy"
)

func TestIHappyScraper_GetJobLinks(t *testing.T) {
	html := `<html><body>
		<div class="job-card"><a href="/job1">Job 1</a></div>
		<div class="job-card"><a href="/job2">Job 2</a></div>
	</body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := ihappy.NewIHappyScraper(server.Client())
	links, err := scraper.GetJobLinks(context.Background(), server.URL)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(links) != 2 {
		t.Errorf("expected 2 links, got %d", len(links))
	}
	if links[0] != server.URL+"/job1" {
		t.Errorf("expected %s/job1, got %s", server.URL, links[0])
	}
}

func TestIHappyScraper_GetJobDetails(t *testing.T) {
	html := `<html><body>
		<h1 class="sc_layouts_title_caption">Yankee Rebel Tavern / Mackinac Island / MI</h1>
		<p>Employer Name: Yankee Rebel Tavern</p>
		<p>Position: Line Cook</p>
		<p>Position Type: Restaurant</p>
		<p>City: Mackinac Island</p>
		<p>State: MI</p>
		<div>premium</div>
	</body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := ihappy.NewIHappyScraper(server.Client())
	job, err := scraper.GetJobDetails(context.Background(), server.URL)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if job.EmployerTitle != "Yankee Rebel Tavern" {
		t.Errorf("expected employer title 'Yankee Rebel Tavern', got '%s'", job.EmployerTitle)
	}
	if job.LocationCity != "Mackinac Island" {
		t.Errorf("expected city 'Mackinac Island', got '%s'", job.LocationCity)
	}
	if job.LocationState != "MI" {
		t.Errorf("expected state 'MI', got '%s'", job.LocationState)
	}
	if job.GroupLocation != "Rank C" {
		t.Errorf("expected rank 'Rank C', got '%s'", job.GroupLocation)
	}
	if job.Position != "Line Cook" {
		t.Errorf("expected position 'Line Cook', got '%s'", job.Position)
	}
}
