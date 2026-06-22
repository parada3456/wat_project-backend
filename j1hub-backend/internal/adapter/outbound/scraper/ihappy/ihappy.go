package ihappy

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

type IHappyScraper struct {
	client *http.Client
}

func NewIHappyScraper(client *http.Client) scraper.JobSource {
	if client == nil {
		client = http.DefaultClient
	}
	return &IHappyScraper{
		client: client,
	}
}

func (s *IHappyScraper) GetJobLinks(ctx context.Context, listURL string) ([]string, error) {
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
	doc.Find(".job-card a").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if exists {
			if u, err := url.Parse(href); err == nil {
				links = append(links, baseURL.ResolveReference(u).String())
			}
		}
	})

	return links, nil
}

func (s *IHappyScraper) GetJobDetails(ctx context.Context, detailURL string) (*jobdomain.JobPosting, error) {
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

	employerTitle := ""
	doc.Find("p, span, h1").Each(func(i int, sel *goquery.Selection) {
		text := sel.Text()
		if strings.Contains(text, "Employer Name:") {
			employerTitle = strings.TrimSpace(strings.TrimPrefix(text, "Employer Name:"))
		}
	})
	if employerTitle == "" {
		h1Text := doc.Find(".sc_layouts_title_caption").Text()
		if h1Text == "" {
			h1Text = doc.Find("h1").Text()
		}
		if idx := strings.Index(h1Text, "/"); idx != -1 {
			employerTitle = strings.TrimSpace(h1Text[:idx])
		} else {
			employerTitle = h1Text
		}
	}
	employerTitle = strings.TrimSpace(employerTitle)

	position := ""
	doc.Find("p, span").Each(func(i int, sel *goquery.Selection) {
		text := sel.Text()
		if strings.Contains(text, "Position:") && !strings.Contains(text, "Available") {
			position = strings.TrimSpace(strings.TrimPrefix(text, "Position:"))
		}
	})
	position = strings.TrimSpace(position)

	positionType := ""
	doc.Find("p, span, tr").Each(func(i int, sel *goquery.Selection) {
		text := sel.Text()
		if strings.Contains(text, "Position Type:") {
			positionType = strings.TrimSpace(strings.TrimPrefix(text, "Position Type:"))
		} else if goquery.NodeName(sel) == "tr" && strings.Contains(sel.Find("td:first-child").Text(), "Position Type") {
			positionType = sel.Find("td:last-child").Text()
		}
	})
	if positionType == "" {
		lowerPos := strings.ToLower(position)
		if strings.Contains(lowerPos, "cook") || strings.Contains(lowerPos, "busser") || strings.Contains(lowerPos, "server") {
			positionType = "Restaurant"
		} else {
			positionType = "Resort Worker"
		}
	}
	positionType = strings.TrimSpace(positionType)

	city := ""
	state := ""
	doc.Find("p, span").Each(func(i int, sel *goquery.Selection) {
		text := sel.Text()
		if strings.Contains(text, "City:") {
			city = strings.TrimSpace(strings.TrimPrefix(text, "City:"))
		}
		if strings.Contains(text, "States:") || strings.Contains(text, "State:") {
			state = strings.TrimSpace(strings.TrimPrefix(text, "States:"))
			state = strings.TrimSpace(strings.TrimPrefix(state, "State:"))
		}
	})
	if city == "" || state == "" {
		h1Text := doc.Find(".sc_layouts_title_caption").Text()
		if idx := strings.Index(h1Text, "/"); idx != -1 {
			parts := strings.Split(h1Text, "/")
			if len(parts) >= 3 {
				city = strings.TrimSpace(parts[1])
				state = strings.TrimSpace(parts[2])
			}
		}
	}
	city = strings.TrimSpace(city)
	state = strings.TrimSpace(state)

	rank := "Rank D"
	bodyText := doc.Text()
	lowerBody := strings.ToLower(bodyText)
	if strings.Contains(lowerBody, "signature") {
		rank = "Rank S"
	} else if strings.Contains(lowerBody, "prestige") {
		rank = "Rank A"
	} else if strings.Contains(lowerBody, "super premium") {
		rank = "Rank B"
	} else if strings.Contains(lowerBody, "premium") {
		rank = "Rank C"
	} else if strings.Contains(lowerBody, "promotion") {
		rank = "Rank D"
	}

	jobID := fmt.Sprintf("%x", md5.Sum([]byte(detailURL)))

	job := &jobdomain.JobPosting{
		JobID:         jobID,
		AgencyName:    "iHappy",
		EmployerTitle: employerTitle,
		Position:      position,
		PositionType:  positionType,
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
