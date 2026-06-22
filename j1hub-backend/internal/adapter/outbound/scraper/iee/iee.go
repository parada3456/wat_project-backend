package iee

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/j1hub/backend/internal/adapter/outbound/scraper"
	jobdomain "github.com/j1hub/backend/internal/job/domain"
)

type IEEScraper struct {
	client *http.Client
}

func NewIEEScraper(client *http.Client) scraper.JobSource {
	if client == nil {
		client = http.DefaultClient
	}
	return &IEEScraper{
		client: client,
	}
}

func (s *IEEScraper) GetJobLinks(ctx context.Context, listURL string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, listURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(listURL)
	if err != nil {
		return nil, err
	}

	var links []string
	// Assuming links are inside standard grid/list on IEE
	doc.Find("a").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if exists && strings.Contains(href, "work_and_travel") {
			if u, err := url.Parse(href); err == nil {
				absoluteURL := baseURL.ResolveReference(u).String()
				// Basic deduplication & filter
				if absoluteURL != listURL {
					links = append(links, absoluteURL)
				}
			}
		}
	})

	return links, nil
}

func (s *IEEScraper) GetJobDetails(ctx context.Context, detailURL string) (*jobdomain.JobPosting, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, detailURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	employerTitle := doc.Find("h1").First().Text()
	employerTitle = strings.TrimSpace(employerTitle)
	if employerTitle == "" {
		employerTitle = "Unknown Employer"
	}

	// Basic parsing from generic tags since layout is not fully known
	position := "Team Member"
	city := "Unknown City"
	state := "Unknown State"

	doc.Find("p, li, div").Each(func(i int, sel *goquery.Selection) {
		text := strings.TrimSpace(sel.Text())
		lowerText := strings.ToLower(text)

		if strings.Contains(lowerText, "position:") && position == "Team Member" {
			parts := strings.SplitN(text, ":", 2)
			if len(parts) == 2 {
				position = strings.TrimSpace(parts[1])
			}
		}
		if strings.Contains(lowerText, "city:") && city == "Unknown City" {
			parts := strings.SplitN(text, ":", 2)
			if len(parts) == 2 {
				city = strings.TrimSpace(parts[1])
			}
		}
		if strings.Contains(lowerText, "state:") && state == "Unknown State" {
			parts := strings.SplitN(text, ":", 2)
			if len(parts) == 2 {
				state = strings.TrimSpace(parts[1])
			}
		}
	})

	jobID := fmt.Sprintf("%x", md5.Sum([]byte(detailURL)))

	job := &jobdomain.JobPosting{
		JobID:         jobID,
		AgencyName:    "IEE",
		EmployerTitle: employerTitle,
		Position:      position,
		PositionType:  "General",
		LocationCity:  city,
		LocationState: state,
		GroupLocation: "Rank D",
		USSponsor:     true,
		SourceURL:     detailURL,
		ScrapeAt:      time.Now(),
		PostedAt:      time.Now(),
		UpdatedAt:     time.Now(),
	}

	return job, nil
}
