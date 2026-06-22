package acadex_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/j1hub/backend/internal/adapter/outbound/scraper/acadex"
)

func TestAcadexScraper_GetJobLinks(t *testing.T) {
	html := `<html><body>
		<div class="job-item"><a class="detail-link" href="/job1">Job 1</a></div>
		<div class="job-item"><a class="detail-link" href="/job2">Job 2</a></div>
	</body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := acadex.NewAcadexScraper(server.Client())
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

func TestAcadexScraper_GetJobDetails(t *testing.T) {
	html := `<html><body>
		<h1 class="location_title">Rosauers Supermarkets, Kalispell - Group A</h1>
		<div class="location_txt"><span>Kalispell, MT</span></div>
		<div class="main_subtitle">
			<span class="subtitle">Position</span>
			<span class="subtitlelist">- Cashier</span>
			<span class="subtitlelist">- Stocker</span>
		</div>
	</body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := acadex.NewAcadexScraper(server.Client())
	job, err := scraper.GetJobDetails(context.Background(), server.URL)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if job.EmployerTitle != "Rosauers Supermarkets" {
		t.Errorf("expected employer title 'Rosauers Supermarkets', got '%s'", job.EmployerTitle)
	}
	if job.LocationCity != "Kalispell" {
		t.Errorf("expected city 'Kalispell', got '%s'", job.LocationCity)
	}
	if job.LocationState != "MT" {
		t.Errorf("expected state 'MT', got '%s'", job.LocationState)
	}
	if job.GroupLocation != "Rank S" {
		t.Errorf("expected rank 'Rank S', got '%s'", job.GroupLocation)
	}
	if job.Position != "Cashier / Stocker" {
		t.Errorf("expected position 'Cashier / Stocker', got '%s'", job.Position)
	}
}
