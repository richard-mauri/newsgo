// discoverycheck evaluates third-party keyword-to-RSS search services.
// It is diagnostic only: it makes no changes to NewsGo subscriptions.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type queries []string

func (q *queries) String() string { return strings.Join(*q, ", ") }
func (q *queries) Set(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("empty query")
	}
	*q = append(*q, s)
	return nil
}

type Result struct {
	Provider   string      `json:"provider"`
	Query      string      `json:"query"`
	Status     int         `json:"http_status"`
	DurationMS int64       `json:"duration_ms"`
	Count      int         `json:"returned"`
	Candidates []Candidate `json:"candidates,omitempty"`
	Error      string      `json:"error,omitempty"`
}
type Candidate struct {
	Title       string `json:"title"`
	URL         string `json:"feed_url"`
	Website     string `json:"website,omitempty"`
	Description string `json:"description,omitempty"`
}

type provider struct{ name, endpoint, param string }

var providers = []provider{
	{"Feedly", "https://cloud.feedly.com/v3/search/feeds", "query"},
	{"RSSAmplifier", "https://rssamplifier.com/api/search", "q"},
}

func candidateURL(s string) string {
	if strings.HasPrefix(s, "feed/") {
		s = strings.TrimPrefix(s, "feed/")
	}
	u, e := url.Parse(s)
	if e != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" {
		return ""
	}
	return s
}
func parseFeedly(body []byte) ([]Candidate, error) {
	var payload struct {
		Results []struct {
			Title       string `json:"title"`
			FeedID      string `json:"feedId"`
			ID          string `json:"id"`
			Website     string `json:"website"`
			Description string `json:"description"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	out := []Candidate{}
	for _, x := range payload.Results {
		fid := x.FeedID
		if fid == "" {
			fid = x.ID
		}
		if u := candidateURL(fid); u != "" {
			out = append(out, Candidate{x.Title, u, x.Website, x.Description})
		}
	}
	return out, nil
}
func parseAmplifier(body []byte) ([]Candidate, error) {
	// Response formats may evolve; accept a direct array or a documented collection wrapper.
	var raw any
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	var rows []any
	switch v := raw.(type) {
	case []any:
		rows = v
	case map[string]any:
		for _, key := range []string{"results", "feeds", "items", "data"} {
			if a, ok := v[key].([]any); ok {
				rows = a
				break
			}
		}
		if rows == nil {
			return nil, errors.New("unrecognized JSON response shape")
		}
	default:
		return nil, errors.New("unrecognized JSON response shape")
	}
	out := []Candidate{}
	for _, row := range rows {
		m, ok := row.(map[string]any)
		if !ok {
			continue
		}
		get := func(keys ...string) string {
			for _, k := range keys {
				if s, ok := m[k].(string); ok && s != "" {
					return s
				}
			}
			return ""
		}
		f := candidateURL(get("feed_url", "feedUrl", "rss_url", "rssUrl", "url"))
		if f == "" {
			continue
		}
		out = append(out, Candidate{get("title", "name"), f, get("website", "site_url", "siteUrl"), get("description")})
	}
	return out, nil
}
func probe(ctx context.Context, client *http.Client, p provider, query string, max int) Result {
	r := Result{Provider: p.name, Query: query}
	u, _ := url.Parse(p.endpoint)
	v := u.Query()
	v.Set(p.param, query)
	v.Set("count", fmt.Sprint(max))
	v.Set("limit", fmt.Sprint(max))
	u.RawQuery = v.Encode()
	req, e := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if e != nil {
		r.Error = e.Error()
		return r
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "NewsGo-DiscoveryCheck/0.1 (+https://github.com/richard-mauri/newsgo)")
	start := time.Now()
	res, e := client.Do(req)
	r.DurationMS = time.Since(start).Milliseconds()
	if e != nil {
		r.Error = e.Error()
		return r
	}
	defer res.Body.Close()
	r.Status = res.StatusCode
	body, e := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if e != nil {
		r.Error = e.Error()
		return r
	}
	if res.StatusCode != 200 {
		r.Error = fmt.Sprintf("HTTP %d: %s", res.StatusCode, strings.TrimSpace(string(body[:min(len(body), 160)])))
		return r
	}
	if p.name == "Feedly" {
		r.Candidates, e = parseFeedly(body)
	} else {
		r.Candidates, e = parseAmplifier(body)
	}
	if e != nil {
		r.Error = e.Error()
		return r
	}
	if len(r.Candidates) > max {
		r.Candidates = r.Candidates[:max]
	}
	r.Count = len(r.Candidates)
	return r
}
func main() {
	var qs queries
	flag.Var(&qs, "q", "topic or publisher to test; repeat for multiple searches")
	max := flag.Int("limit", 5, "maximum results per search (1–20)")
	timeout := flag.Duration("timeout", 12*time.Second, "per-request timeout")
	output := flag.String("output", "", "optional JSON report destination")
	flag.Parse()
	if len(qs) == 0 {
		qs = queries{"college football", "sailing", "gardening"}
	}
	if *max < 1 || *max > 20 {
		fmt.Fprintln(os.Stderr, "-limit must be 1–20")
		os.Exit(2)
	}
	client := &http.Client{Timeout: *timeout}
	results := make([]Result, 0, len(qs)*len(providers))
	for _, q := range qs {
		for _, p := range providers {
			ctx, cancel := context.WithTimeout(context.Background(), *timeout)
			r := probe(ctx, client, p, q, *max)
			cancel()
			results = append(results, r)
			fmt.Printf("\n%s | %q | HTTP %d | %d ms | %d usable feeds\n", r.Provider, r.Query, r.Status, r.DurationMS, r.Count)
			if r.Error != "" {
				fmt.Printf("  Error: %s\n", r.Error)
			}
			for _, c := range r.Candidates {
				fmt.Printf("  %s\n    %s\n", c.Title, c.URL)
			}
		}
	}
	if *output != "" {
		b, e := json.MarshalIndent(results, "", "  ")
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(1)
		}
		if e = os.WriteFile(*output, append(b, '\n'), 0600); e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(1)
		}
		fmt.Println("\nSaved:", *output)
	}
}
