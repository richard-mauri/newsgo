package dashboard

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type feedSearchItem struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Website string `json:"website,omitempty"`
}
type feedSearchResponse struct {
	Feeds []feedSearchItem `json:"feeds"`
	Error string           `json:"error,omitempty"`
}

var feedSearchClient = &http.Client{Timeout: 7 * time.Second}
var feedSearchState = struct {
	sync.Mutex
	last  map[string]time.Time
	cache map[string]struct {
		feeds []feedSearchItem
		at    time.Time
	}
}{last: map[string]time.Time{}, cache: map[string]struct {
	feeds []feedSearchItem
	at    time.Time
}{}}

func publicCandidate(raw string) bool {
	u, e := url.Parse(raw)
	if e != nil {
		return false
	}
	host := u.Hostname()
	if ip := net.ParseIP(host); ip != nil {
		return publicIP(ip)
	}
	return host != "localhost" && !strings.HasSuffix(host, ".localhost")
}
func feedlyResults(body []byte) ([]feedSearchItem, error) {
	var payload struct {
		Results []struct {
			Title   string `json:"title"`
			FeedID  string `json:"feedId"`
			ID      string `json:"id"`
			Website string `json:"website"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	out := []feedSearchItem{}
	for _, it := range payload.Results {
		raw := it.FeedID
		if raw == "" {
			raw = it.ID
		}
		raw = strings.TrimPrefix(raw, "feed/")
		if !safeURL(raw) || !publicCandidate(raw) || seen[raw] || len(out) >= 12 {
			continue
		}
		seen[raw] = true
		name := strings.TrimSpace(it.Title)
		if name == "" {
			u, _ := url.Parse(raw)
			name = u.Hostname()
		}
		if len(name) > 70 {
			name = name[:70]
		}
		website := ""
		if safeURL(it.Website) {
			website = it.Website
		}
		out = append(out, feedSearchItem{Name: name, URL: raw, Website: website})
	}
	return out, nil
}
func feedlySearch(query string) ([]feedSearchItem, error) {
	endpoint := "https://cloud.feedly.com/v3/search/feeds?query=" + url.QueryEscape(query) + "&count=12"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "NewsGo/0.2.9 (feed discovery)")
	resp, err := feedSearchClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("feed directory returned HTTP " + resp.Status)
	}
	limited := io.LimitReader(resp.Body, 256*1024+1)
	body, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if len(body) > 256*1024 {
		return nil, errors.New("feed directory response too large")
	}
	return feedlyResults(body)
}
func (s *State) serveFeedSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method not allowed", 405)
		return
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if len(q) < 2 || len(q) > 80 {
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(feedSearchResponse{Error: "Enter a topic or publisher (2–80 characters)."})
		return
	}
	normalized := strings.ToLower(q)
	feedSearchState.Lock()
	if entry, ok := feedSearchState.cache[normalized]; ok && time.Since(entry.at) < 10*time.Minute {
		feedSearchState.Unlock()
		_ = json.NewEncoder(w).Encode(feedSearchResponse{Feeds: entry.feeds})
		return
	}
	// Bound remote calls in a public deployment. Entries are shared by query, not by user.
	now := time.Now()
	for k, t := range feedSearchState.last {
		if now.Sub(t) > time.Minute {
			delete(feedSearchState.last, k)
		}
	}
	ip := r.RemoteAddr
	if h, _, err := net.SplitHostPort(ip); err == nil {
		ip = h
	}
	if now.Sub(feedSearchState.last[ip]) < 3*time.Second {
		feedSearchState.Unlock()
		w.WriteHeader(429)
		_ = json.NewEncoder(w).Encode(feedSearchResponse{Error: "Please wait a moment before searching again."})
		return
	}
	feedSearchState.last[ip] = now
	feedSearchState.Unlock()
	results, err := feedlySearch(q)
	if err != nil {
		w.WriteHeader(502)
		_ = json.NewEncoder(w).Encode(feedSearchResponse{Error: "Feed search is temporarily unavailable."})
		return
	}
	feedSearchState.Lock()
	if len(feedSearchState.cache) > 100 {
		feedSearchState.cache = map[string]struct {
			feeds []feedSearchItem
			at    time.Time
		}{}
	}
	feedSearchState.cache[normalized] = struct {
		feeds []feedSearchItem
		at    time.Time
	}{results, time.Now()}
	feedSearchState.Unlock()
	_ = json.NewEncoder(w).Encode(feedSearchResponse{Feeds: results})
}
