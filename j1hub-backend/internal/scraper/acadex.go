package scraper

import (
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/j1hub/backend/internal/domain"
)

func (s *Scraper) ScrapeAcadex(url string, results chan<- Result) {
	c := s.Collector.Clone()

	// If this URL is already a detail page, parse it directly
	if strings.Contains(url, "/location/") {
		c.OnHTML("body", func(de *colly.HTMLElement) {
			s.parseAcadexDetail(url, "Full-Time", de, results)
		})
		c.Visit(url)
		return
	}

	// List page parsing
	c.OnHTML(".job-item", func(e *colly.HTMLElement) {
		// Capture position type from list level
		positionType := e.ChildText(".category-tag")
		detailURL := e.ChildAttr("a.detail-link", "href")
		absoluteURL := e.Request.AbsoluteURL(detailURL)

		detailCollector := s.Collector.Clone()
		detailCollector.OnHTML("body", func(de *colly.HTMLElement) {
			s.parseAcadexDetail(absoluteURL, positionType, de, results)
		})

		detailCollector.Visit(absoluteURL)
	})

	c.Visit(url)
}

func (s *Scraper) parseAcadexDetail(absoluteURL string, positionType string, de *colly.HTMLElement, results chan<- Result) {
	jobID := GenerateJobID(absoluteURL)
	
	titleText := de.ChildText(".location_title")
	if titleText == "" {
		titleText = de.ChildText("h1")
	}
	titleText = strings.TrimSpace(titleText)

	employerTitle := titleText
	if idx := strings.Index(titleText, ","); idx != -1 {
		employerTitle = strings.TrimSpace(titleText[:idx])
	}

	locationText := de.ChildText(".location_txt span")
	if locationText == "" {
		locationText = de.ChildText(".location_txt")
	}
	locationText = strings.TrimSpace(locationText)
	city, state := parseLocation(locationText)

	var rank string
	lowerTitle := strings.ToLower(titleText)
	lowerURL := strings.ToLower(absoluteURL)
	if strings.Contains(lowerTitle, "group a") || strings.Contains(lowerURL, "group-a") {
		rank = "A"
	} else if strings.Contains(lowerTitle, "group x") || strings.Contains(lowerURL, "group-x") {
		rank = "X"
	} else if strings.Contains(lowerTitle, "group y") || strings.Contains(lowerURL, "group-y") {
		rank = "Y"
	} else if strings.Contains(lowerTitle, "group z") || strings.Contains(lowerURL, "group-z") {
		rank = "Z"
	}

	var positions []string
	de.ForEach(".main_subtitle", func(_ int, se *colly.HTMLElement) {
		sub := se.ChildText(".subtitle")
		if strings.TrimSpace(sub) == "Position" {
			se.ForEach(".subtitlelist", func(_ int, sle *colly.HTMLElement) {
				pText := strings.TrimPrefix(strings.TrimSpace(sle.Text), "- ")
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

	jobPosting := &domain.JobPosting{
		JobID:         jobID,
		AgencyName:    "Acadex",
		EmployerTitle: employerTitle,
		Position:      position,
		PositionType:  positionType,
		LocationCity:  city,
		LocationState: state,
		GroupLocation: MapLocationRank(rank, "Acadex"),
		USSponsor:     true,
		SourceURL:     absoluteURL,
		ScrapeAt:      time.Now(),
		PostedAt:      time.Now(),
		UpdatedAt:     time.Now(),
	}

	housing := &domain.JobHousing{
		HousingID: GenerateJobID(absoluteURL + "/housing"),
		JobID:     jobID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	results <- Result{Job: jobPosting, Housing: housing}
}

func parseLocation(loc string) (string, string) {
	parts := strings.Split(loc, ",")
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	return loc, ""
}
