package iee_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/j1hub/backend/internal/infrastructure/outbound/scraper/iee"
)

func TestIEEScraper_GetJobLinks(t *testing.T) {
	html := `<html><body>
		<div><a href="/work_and_travel_2/mcdonalds/">McDonald's</a></div>
		<div><a href="/other">Other</a></div>
	</body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := iee.NewIEEScraper(server.Client())
	links, err := scraper.GetJobLinks(context.Background(), server.URL)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(links) != 1 {
		t.Errorf("expected 1 link, got %d", len(links))
	}
	if links[0] != server.URL+"/work_and_travel_2/mcdonalds/" {
		t.Errorf("expected %s/work_and_travel_2/mcdonalds/, got %s", server.URL, links[0])
	}
}

func TestIEEScraper_GetJobDetails(t *testing.T) {
	html := `<html><body>
		<h1>McDonald's - Florida</h1>
		<p>Position: Cashier</p>
		<ul>
			<li>City: Orlando</li>
			<li>State: FL</li>
		</ul>
	</body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := iee.NewIEEScraper(server.Client())
	job, err := scraper.GetJobDetails(context.Background(), server.URL)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if job.EmployerTitle != "McDonald's - Florida" {
		t.Errorf("expected employer title 'McDonald's - Florida', got '%s'", job.EmployerTitle)
	}
	if job.LocationCity != "Orlando" {
		t.Errorf("expected city 'Orlando', got '%s'", job.LocationCity)
	}
	if job.LocationState != "FL" {
		t.Errorf("expected state 'FL', got '%s'", job.LocationState)
	}
	if job.Position != "Cashier" {
		t.Errorf("expected position 'Cashier', got '%s'", job.Position)
	}
}
