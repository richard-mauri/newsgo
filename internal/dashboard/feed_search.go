package dashboard

import (
	"context"
 "encoding/json"
 "fmt"
 "strconv"
 "github.com/richard-mauri/newsgo/internal/feeds"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

type feedSearchItem struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Website string `json:"website,omitempty"`
 Status string `json:"status,omitempty"`
 Reason string `json:"reason,omitempty"`
}
type feedSearchResponse struct {
	Feeds []feedSearchItem `json:"feeds"`
	Error string           `json:"error,omitempty"`
}

// The feed-directory timeout is independent of article classification.
func feedSearchClient() *http.Client {
 return &http.Client{Timeout: time.Duration(currentArticleSettings().FeedSearchTimeoutSeconds) * time.Second}
}
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
		if !safeURL(raw) || !publicCandidate(raw) || seen[raw] {
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
func feedMatchRank(name, query string) int {
	n := strings.ToLower(strings.TrimSpace(name))
	if n == query {
		return 0
	}
	if strings.HasPrefix(n, query) {
		return 1
	}
	if strings.Contains(n, query) {
		return 2
	}
	return 3
}
func feedlySearch(query string) ([]feedSearchItem, error) {
	endpoint := "https://cloud.feedly.com/v3/search/feeds?query=" + url.QueryEscape(query) + "&count=12"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "NewsGo/0.2.9 (feed discovery)")
	resp, err := feedSearchClient().Do(req)
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
	items, err := feedlyResults(body)
	if err != nil {
		return nil, err
	}
	q := strings.ToLower(strings.TrimSpace(query))
	// Exact publisher names appear before broad topic and substring matches.
	sort.SliceStable(items, func(i, j int) bool { return feedMatchRank(items[i].Name, q) < feedMatchRank(items[j].Name, q) })
	if len(items) > 12 {
		items = items[:12]
	}
	return items, nil
}
// Discovery candidates are NOT trusted until an actual RSS/Atom/JSON body parses.
// Reuse the hardened personal feed client, including public-IP DNS/redirect checks.
// A publisher outage does not prove that an endpoint is obsolete.
func verifyDirectoryFeed(ctx context.Context, client *http.Client, candidate feedSearchItem) feedSearchItem {
 candidate.Status="unverified"
 candidate.Reason="Could not verify this feed right now"
 if !safeURL(candidate.URL)||!publicCandidate(candidate.URL) {
  candidate.Status="unavailable"
  candidate.Reason="Invalid feed URL"
  return candidate
 }
 articles,_,_,_,err:=feeds.FetchArticlesWithRaw(ctx,client,feeds.Source{Name:candidate.Name,URL:candidate.URL})
 if err==nil {
  if len(articles)==0 { candidate.Status="unavailable";candidate.Reason="No readable articles in this feed" } else { candidate.Status="available";candidate.Reason="" }
  return candidate
 }
 if strings.HasPrefix(err.Error(),"HTTP ") {
  status,_:=strconv.Atoi(strings.TrimPrefix(err.Error(),"HTTP "))
  candidate.Reason=fmt.Sprintf("HTTP %d",status)
  switch status {
  case 400,404,410,422:
   candidate.Status="unavailable"
  default:
   candidate.Status="unverified" // 403, 429, 5xx and redirects may be transient.
  }
  return candidate
 }
 if ctx.Err()!=nil {candidate.Reason="Feed verification timed out";return candidate}
 if strings.Contains(err.Error(),"invalid feed")||strings.Contains(err.Error(),"unsupported XML")||strings.Contains(err.Error(),"feed exceeds 4 MiB") {
  candidate.Status="unavailable"
  candidate.Reason="Invalid or oversized feed response"
  return candidate
 }
 candidate.Reason="Feed could not be reached or verified"
 return candidate
}

// The maximum work is 12 candidates, with four simultaneous bounded network
// calls and a shared 9-second budget. No article classification is performed.
func verifyDirectoryCandidates(ctx context.Context, results []feedSearchItem, verify func(context.Context,feedSearchItem) feedSearchItem) []feedSearchItem {
 if len(results)==0{return results}
 ctx,cancel:=context.WithTimeout(ctx,9*time.Second)
 defer cancel()
 out:=append([]feedSearchItem(nil),results...)
 slots:=make(chan struct{},4)
 var wg sync.WaitGroup
 for i:=range out {
  i:=i
  wg.Add(1)
  go func(){
   defer wg.Done()
   select {case slots<-struct{}{}:defer func(){<-slots}();case <-ctx.Done():out[i].Status="unverified";out[i].Reason="Feed verification timed out";return}
   out[i]=verify(ctx,out[i])
  }()
 }
 wg.Wait()
 return out
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
	if entry, ok := feedSearchState.cache[normalized]; ok && time.Since(entry.at) < 2*time.Minute {
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
	// Some directory searches return nothing for the complete publisher name,
	// yet find it for a shorter prefix. Retry once, with the original query
	// used to prioritize matching publisher names in the combined results.
	if err == nil && len(results) == 0 && len([]rune(q)) >= 5 {
		runes := []rune(q)
		fallback, retryErr := feedlySearch(string(runes[:len(runes)-1]))
		if retryErr == nil {
			sort.SliceStable(fallback, func(i, j int) bool {
				return feedMatchRank(fallback[i].Name, strings.ToLower(q)) < feedMatchRank(fallback[j].Name, strings.ToLower(q))
			})
			results = fallback
		}
	}
	if err != nil {
		w.WriteHeader(502)
		_ = json.NewEncoder(w).Encode(feedSearchResponse{Error: "Feed search is temporarily unavailable."})
		return
	}
	// Feedly may rank obsolete New Scientist endpoints above the publisher's
	// current main feed. This is a discovery suggestion, never an auto-subscription.
	if strings.EqualFold(q, "New Scientist") {
		constURL := "https://www.newscientist.com/feed/"
		preferred := feedSearchItem{Name: "New Scientist", URL: constURL, Website: "https://www.newscientist.com/"}
		filtered := []feedSearchItem{preferred}
		for _, item := range results {
			if item.URL != constURL {
				filtered = append(filtered, item)
			}
		}
		if len(filtered) > 12 {
			filtered = filtered[:12]
		}
		results = filtered
	}
	client:=safeFeedClient()
 client.Timeout=4*time.Second
 results=verifyDirectoryCandidates(r.Context(),results,func(ctx context.Context,item feedSearchItem)feedSearchItem{return verifyDirectoryFeed(ctx,client,item)})
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
