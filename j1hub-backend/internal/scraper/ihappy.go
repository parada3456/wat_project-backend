package scraper

import (
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/j1hub/backend/internal/domain"
)

func (s *Scraper) ScrapeIHappy(url string, results chan<- Result) {
	c := s.Collector.Clone()

	// If this URL is already a detail page, parse it directly
	if !strings.HasSuffix(strings.TrimRight(url, "/"), "job-location-summer") && !strings.HasSuffix(strings.TrimRight(url, "/"), "job-location-spring") {
		c.OnHTML("body", func(de *colly.HTMLElement) {
			s.parseIHappyDetail(url, de, results)
		})
		c.Visit(url)
		return
	}

	c.OnHTML(".job-card", func(e *colly.HTMLElement) {
		detailURL := e.ChildAttr("a", "href")
		absoluteURL := e.Request.AbsoluteURL(detailURL)

		detailCollector := s.Collector.Clone()
		detailCollector.OnHTML("body", func(de *colly.HTMLElement) {
			s.parseIHappyDetail(absoluteURL, de, results)
		})
		detailCollector.Visit(absoluteURL)
	})

	c.Visit(url)
}

func (s *Scraper) parseIHappyDetail(absoluteURL string, de *colly.HTMLElement, results chan<- Result) {
	jobID := GenerateJobID(absoluteURL)

	employerTitle := ""
	de.ForEach("p, span, h1", func(_ int, element *colly.HTMLElement) {
		text := element.Text
		if strings.Contains(text, "Employer Name:") {
			employerTitle = strings.TrimSpace(strings.TrimPrefix(text, "Employer Name:"))
		}
	})
	if employerTitle == "" {
		h1Text := de.ChildText(".sc_layouts_title_caption")
		if h1Text == "" {
			h1Text = de.ChildText("h1")
		}
		if idx := strings.Index(h1Text, "/"); idx != -1 {
			employerTitle = strings.TrimSpace(h1Text[:idx])
		} else {
			employerTitle = h1Text
		}
	}
	employerTitle = strings.TrimSpace(employerTitle)

	position := ""
	de.ForEach("p, span", func(_ int, element *colly.HTMLElement) {
		text := element.Text
		if strings.Contains(text, "Position:") && !strings.Contains(text, "Available") {
			position = strings.TrimSpace(strings.TrimPrefix(text, "Position:"))
		}
	})
	position = strings.TrimSpace(position)

	positionType := ""
	de.ForEach("p, span, tr", func(_ int, element *colly.HTMLElement) {
		text := element.Text
		if strings.Contains(text, "Position Type:") {
			positionType = strings.TrimSpace(strings.TrimPrefix(text, "Position Type:"))
		} else if element.Name == "tr" && strings.Contains(element.ChildText("td:first-child"), "Position Type") {
			positionType = element.ChildText("td:last-child")
		}
	})
	if positionType == "" {
		if strings.Contains(strings.ToLower(position), "cook") || strings.Contains(strings.ToLower(position), "busser") || strings.Contains(strings.ToLower(position), "server") {
			positionType = "Restaurant"
		} else {
			positionType = "Resort Worker"
		}
	}
	positionType = strings.TrimSpace(positionType)

	city := ""
	state := ""
	de.ForEach("p, span", func(_ int, element *colly.HTMLElement) {
		text := element.Text
		if strings.Contains(text, "City:") {
			city = strings.TrimSpace(strings.TrimPrefix(text, "City:"))
		}
		if strings.Contains(text, "States:") || strings.Contains(text, "State:") {
			state = strings.TrimSpace(strings.TrimPrefix(text, "States:"))
			state = strings.TrimSpace(strings.TrimPrefix(state, "State:"))
		}
	})
	if city == "" || state == "" {
		h1Text := de.ChildText(".sc_layouts_title_caption")
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

	rank := "premium" // default
	bodyText := de.Text
	if strings.Contains(strings.ToLower(bodyText), "super premium") {
		rank = "super premium"
	} else if strings.Contains(strings.ToLower(bodyText), "signature") {
		rank = "signature"
	} else if strings.Contains(strings.ToLower(bodyText), "prestige") {
		rank = "prestige"
	} else if strings.Contains(strings.ToLower(bodyText), "premium") {
		rank = "premium"
	} else if strings.Contains(strings.ToLower(bodyText), "promotion") {
		rank = "promotion"
	}

	jobPosting := &domain.JobPosting{
		JobID:         jobID,
		AgencyName:    "iHappy",
		EmployerTitle: employerTitle,
		Position:      position,
		PositionType:  positionType,
		LocationCity:  city,
		LocationState: state,
		GroupLocation: MapLocationRank(rank, "IHappy"),
		USSponsor:     true,
		SourceURL:     absoluteURL,
		ScrapeAt:      time.Now(),
		PostedAt:      time.Now(),
		UpdatedAt:     time.Now(),
	}

	var weeklyRate float64
	de.ForEach("p, span, tr", func(_ int, element *colly.HTMLElement) {
		text := element.Text
		if strings.Contains(text, "Housing Cost:") {
			costStr := strings.TrimSpace(strings.TrimPrefix(text, "Housing Cost:"))
			costStr = strings.TrimPrefix(costStr, "$")
			if idx := strings.Index(costStr, " "); idx != -1 {
				costStr = costStr[:idx]
			}
			if val, err := strconv.ParseFloat(costStr, 64); err == nil {
				weeklyRate = val
			}
		} else if element.Name == "tr" && strings.Contains(element.ChildText("td:first-child"), "Housing Weekly Rate") {
			rateStr := element.ChildText("td:last-child")
			rateStr = strings.TrimPrefix(rateStr, "$")
			if val, err := strconv.ParseFloat(rateStr, 64); err == nil {
				weeklyRate = val
			}
		}
	})

	housing := &domain.JobHousing{
		HousingID:  GenerateJobID(absoluteURL + "/housing"),
		JobID:      jobID,
		WeeklyRate: weeklyRate,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	results <- Result{Job: jobPosting, Housing: housing}
}
