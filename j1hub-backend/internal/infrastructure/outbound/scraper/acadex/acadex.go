package acadex

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	jobdomain "github.com/parada3456/wat_project-backend/internal/job/domain"
)

type AcadexScraper struct {
	client *http.Client
}

func NewAcadexScraper(client *http.Client) *AcadexScraper {
	if client == nil {
		client = http.DefaultClient
	}
	return &AcadexScraper{
		client: client,
	}
}

func (s *AcadexScraper) GetJobLinks(ctx context.Context, listURL string) ([]string, error) {
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
	doc.Find(".job-item a.detail-link").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if exists {
			if u, err := url.Parse(href); err == nil {
				links = append(links, baseURL.ResolveReference(u).String())
			}
		}
	})

	return links, nil
}

func (s *AcadexScraper) GetJobDetails(ctx context.Context, detailURL string) (*jobdomain.JobPosting, error) {
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

	titleText := doc.Find(".location_title").Text()
	if titleText == "" {
		titleText = doc.Find("h1").Text()
	}
	titleText = strings.TrimSpace(titleText)

	employerTitle := titleText
	if idx := strings.Index(titleText, ","); idx != -1 {
		employerTitle = strings.TrimSpace(titleText[:idx])
	}

	locationText := doc.Find(".location_txt span").Text()
	if locationText == "" {
		locationText = doc.Find(".location_txt").Text()
	}
	locationText = strings.TrimSpace(locationText)

	city, state := locationText, ""
	parts := strings.Split(locationText, ",")
	if len(parts) >= 2 {
		city = strings.TrimSpace(parts[0])
		state = strings.TrimSpace(parts[1])
	}

	var rank string
	lowerTitle := strings.ToLower(titleText)
	lowerURL := strings.ToLower(detailURL)
	if strings.Contains(lowerTitle, "group a") || strings.Contains(lowerURL, "group-a") {
		rank = "Rank S"
	} else if strings.Contains(lowerTitle, "group x") || strings.Contains(lowerURL, "group-x") {
		rank = "Rank A"
	} else if strings.Contains(lowerTitle, "group y") || strings.Contains(lowerURL, "group-y") {
		rank = "Rank B"
	} else if strings.Contains(lowerTitle, "group z") || strings.Contains(lowerURL, "group-z") {
		rank = "Rank C"
	} else {
		rank = "Rank D"
	}

	var positions []string
	doc.Find(".main_subtitle").Each(func(i int, se *goquery.Selection) {
		sub := se.Find(".subtitle").Text()
		if strings.TrimSpace(sub) == "Position" {
			se.Find(".subtitlelist").Each(func(j int, sle *goquery.Selection) {
				pText := strings.TrimPrefix(strings.TrimSpace(sle.Text()), "- ")
				if pText != "" && !strings.Contains(pText, "ผู้พิจารณา") {
					positions = append(positions, pText)
				}
			})
		}
	})
	position := strings.Join(positions, " / ")
	if position == "" {
		position = "Server / Line cook"
	}

	jobID := fmt.Sprintf("%x", md5.Sum([]byte(detailURL)))

	job := &jobdomain.JobPosting{
		JobID:         jobID,
		AgencyName:    "Acadex",
		EmployerTitle: employerTitle,
		Position:      position,
		PositionType:  "Full-Time", // Defaults to Full-Time
		LocationCity:  city,
		LocationState: state,
		GroupLocation: rank,
		USSponsor:     true,
		SourceURL:     detailURL,
		ScrapeAt:      time.Now(),
		PostedAt:      time.Now(),
		UpdatedAt:     time.Now(),
	}

	return job, nil
}
