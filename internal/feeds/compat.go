package feeds

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Compatibility is diagnostic only. A failed sample never changes feed health.
type Compatibility struct {
	Tested     int      `json:"tested"`
	Readable   int      `json:"readable"`
	Restricted int      `json:"restricted"`
	Other      int      `json:"other"`
	Results    []Sample `json:"results"`
}
type Sample struct {
	URL    string `json:"url"`
	Status string `json:"status"`
	Kind   string `json:"kind,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func SampleScraperGo(ctx context.Context, client *http.Client, endpoint string, urls []string, limit int) *Compatibility {
	result := &Compatibility{Results: []Sample{}}
	seen := make(map[string]bool)
	for _, articleURL := range urls {
		if result.Tested >= limit {
			break
		}
		if !ValidURL(articleURL) || seen[articleURL] {
			continue
		}
		seen[articleURL] = true
		result.Tested++
		sample := Sample{URL: articleURL, Status: "other"}
		payload, _ := json.Marshal(map[string]string{"url": articleURL})
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			resp, e := client.Do(req)
			if e != nil {
				err = e
			} else {
				var body struct {
					Kind  string `json:"kind"`
					Title string `json:"title"`
					Error string `json:"error"`
				}
				e = json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&body)
				resp.Body.Close()
				if e != nil {
					err = fmt.Errorf("invalid ScraperGo JSON response: %w", e)
				} else {
					sample.Kind = body.Kind
					switch {
					case resp.StatusCode == http.StatusOK && strings.TrimSpace(body.Title) != "":
						sample.Status = "readable"
						result.Readable++
					case body.Kind == "publisher_access_denied" || body.Kind == "paywall" || body.Kind == "registration" || body.Kind == "restricted":
						sample.Status = "restricted"
						result.Restricted++
					default:
						sample.Detail = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, body.Error)
					}
				}
			}
		}
		if err != nil {
			sample.Detail = err.Error()
		}
		if sample.Status == "other" {
			result.Other++
		}
		result.Results = append(result.Results, sample)
	}
	return result
}
func ScraperGoEndpoint(base string) (string, error) {
	u, err := url.Parse(base)
	if err != nil || u == nil || u.Scheme != "http" && u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Path != "" && u.Path != "/" {
		return "", fmt.Errorf("invalid ScraperGo base URL")
	}
	u.Path = "/api/read"
	return u.String(), nil
}
