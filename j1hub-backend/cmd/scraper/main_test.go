package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"
)

type mockTransport struct{}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var body string
	if req.URL.Path == "/program/work-and-travel-summer/" {
		body = `
			<html>
				<body>
					<div class="job-item">
						<span class="category-tag">Full-Time</span>
						<a class="detail-link" href="/program/work-and-travel-summer/detail">Detail</a>
					</div>
				</body>
			</html>
		`
	} else if req.URL.Path == "/program/work-and-travel-summer/detail" {
		body = `
			<html>
				<body>
					<div class="job-detail">
						<span class="employer-title">Acadex Employer</span>
						<span class="position-name">Developer</span>
						<span class="location">New York, NY</span>
						<span class="rank-tag">A</span>
					</div>
				</body>
			</html>
		`
	} else if req.URL.Path == "/job-location-summer/" {
		body = `
			<html>
				<body>
					<div class="job-card">
						<a href="/job-location-summer/detail">Detail</a>
					</div>
				</body>
			</html>
		`
	} else if req.URL.Path == "/job-location-summer/detail" {
		body = `
			<html>
				<body>
					<div class="job-details-container">
						<h1 class="employer-name">iHappy Employer</h1>
						<h2 class="job-title">Designer</h2>
						<span class="rank-badge">signature</span>
						<table>
							<tr>
								<td>Position Type</td>
								<td>Full-Time</td>
							</tr>
							<tr>
								<td>Housing Weekly Rate</td>
								<td>$150.00</td>
							</tr>
						</table>
					</div>
				</body>
			</html>
		`
	} else {
		body = `<html></html>`
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Request:    req,
	}, nil
}

func TestMainScraper(t *testing.T) {
	originalTransport := http.DefaultTransport
	http.DefaultTransport = &mockTransport{}
	defer func() {
		http.DefaultTransport = originalTransport
		os.Remove("job_posting.json")
		os.Remove("job_housing.json")
	}()

	main()
}
