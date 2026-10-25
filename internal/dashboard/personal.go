package dashboard

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/richard-mauri/newsgo/internal/feeds"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type personalRequest struct {
	Feeds []feeds.Source `json:"feeds"`
 MaxArticleAgeDays int `json:"maxArticleAgeDays,omitempty"`
}
type personalResult struct {
	Source   feeds.Source      `json:"source"`
	Articles []personalArticle `json:"articles,omitempty"`
	Error    string            `json:"error,omitempty"`
 Filtered []filteredArticle `json:"filtered,omitempty"`
 NotProcessed int `json:"notProcessed,omitempty"`
 FeedEntries int `json:"feedEntries,omitempty"`
}
type filteredArticle struct {
 Title string `json:"title"`
 URL string `json:"url"`
 Source string `json:"source"`
 Reason string `json:"reason"`
 Published string `json:"published,omitempty"`
}
type personalArticle struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Open      string `json:"open"`
	Source    string `json:"source"`
	Category  string `json:"category"`
	Date      string `json:"date"`
 DateWarning string `json:"dateWarning,omitempty"`
	Published string `json:"published"`
	Image     string `json:"image,omitempty"`
 FeedImage string `json:"feedImage,omitempty"`
 ContentStatus string `json:"contentStatus,omitempty"`
 AvailabilityStatus string `json:"availabilityStatus,omitempty"`
}

// Personal feed fetching uses a dedicated transport to prevent SSRF, including DNS rebinding and redirects.
func safeFeedClient() *http.Client {
	dialer := &net.Dialer{Timeout: 5 * time.Second, KeepAlive: 10 * time.Second}
	transport := &http.Transport{Proxy: nil, DisableKeepAlives: true, MaxConnsPerHost: 2, ResponseHeaderTimeout: 8 * time.Second, DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return nil, err
		}
		ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
		if err != nil {
			return nil, err
		}
		var last error
		for _, a := range ips {
			if !publicIP(a.IP) {
				continue
			}
			conn, e := dialer.DialContext(ctx, "tcp", net.JoinHostPort(a.IP.String(), port))
			if e == nil {
				return conn, nil
			}
			last = e
		}
		if last != nil {
			return nil, last
		}
		return nil, errors.New("feed destination is not a public IP address")
	}}
	return &http.Client{Timeout: 12 * time.Second, Transport: transport, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 4 || !safeURL(req.URL.String()) {
			return http.ErrUseLastResponse
		}
		return nil
	}}
}
func publicIP(ip net.IP) bool {
	if ip == nil || !ip.IsGlobalUnicast() || ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return false
	}
	v4 := ip.To4()
	if v4 != nil {
		if v4[0] == 0 || v4[0] == 100 && v4[1] >= 64 && v4[1] <= 127 || v4[0] >= 224 || v4[0] == 192 && v4[1] == 0 && v4[2] == 0 || v4[0] == 198 && (v4[1] == 18 || v4[1] == 19) || v4[0] == 169 && v4[1] == 254 {
			return false
		}
		return true
	}
	// Only global IPv6 unicast (2000::/3); excludes ULA, multicast, documentation and IPv4 mappings.
	return len(ip) == 16 && ip[0]&0xe0 == 0x20 && !(ip[0] == 0x20 && ip[1] == 0x01 && ip[2] == 0x0d && ip[3] == 0xb8)
}
func safeURL(raw string) bool {
	u, e := url.Parse(raw)
	return e == nil && (u.Scheme == "https" || u.Scheme == "http") && u.Hostname() != "" && u.User == nil && u.Fragment == "" && u.Port() != "0" && !strings.HasSuffix(strings.ToLower(u.Hostname()), ".local") && !strings.HasSuffix(strings.ToLower(u.Hostname()), ".internal")
}

var personalSlots = make(chan struct{}, 4)
var personalCache = struct {
	sync.Mutex
	entries map[string]cacheEntry
}{entries: map[string]cacheEntry{}}

type cacheEntry struct {
	articles []feeds.Article
	err      string
	at       time.Time
}

// A successful fetch is reusable for ten minutes. Transient failures retry soon,
// and canceled requests are never cached.
func personalCacheTTL(entry cacheEntry) time.Duration {
	if entry.err != "" {
		return 20 * time.Second
	}
	return time.Duration(currentArticleSettings().FeedCacheMinutes) * time.Minute
}

func personalFetch(ctx context.Context, source feeds.Source) ([]feeds.Article, string) {
	if err := ctx.Err(); err != nil {
		return nil, "Request canceled"
	}
	personalCache.Lock()
	cached, ok := personalCache.entries[source.URL]
	personalCache.Unlock()
	if ok && time.Since(cached.at) < personalCacheTTL(cached) {
		return cached.articles, cached.err
	}
	select {
	case personalSlots <- struct{}{}:
		defer func() { <-personalSlots }()
	case <-ctx.Done():
		return nil, "Request canceled"
	}
	// Recheck after waiting for a slot.
	personalCache.Lock()
	cached, ok = personalCache.entries[source.URL]
	personalCache.Unlock()
	if ok && time.Since(cached.at) < personalCacheTTL(cached) {
		return cached.articles, cached.err
	}
	a, err := feeds.FetchArticles(ctx, safeFeedClient(), source)
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	// Do not poison the shared cache with a canceled browser request. A
	// timed-out feed request may be retried, but cannot mask later success.
	if ctx.Err() != nil || errors.Is(err, context.Canceled) {
		return a, msg
	}
	personalCache.Lock()
	if len(personalCache.entries) >= 64 {
		for k := range personalCache.entries {
			delete(personalCache.entries, k)
			break
		}
	}
	personalCache.entries[source.URL] = cacheEntry{a, msg, time.Now()}
	personalCache.Unlock()
	return a, msg
}
func (s *State) servePersonal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not allowed", 405)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	var req personalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Feeds) > 8 || !validArticleAgeDays(req.MaxArticleAgeDays) {
		http.Error(w, "Expected up to 8 feeds", 400)
		return
	}
	results := make([]personalResult, 0, len(req.Feeds))
	seen := map[string]bool{}
	for _, src := range req.Feeds {
		src.Name = strings.TrimSpace(src.Name)
		src.Category = strings.TrimSpace(src.Category)
		src.URL = strings.TrimSpace(src.URL)
		res := personalResult{Source: src}
		if len(src.Name) == 0 || len(src.Name) > 70 || len(src.Category) > 40 || len(src.URL) > 2048 || !safeURL(src.URL) {
			res.Error = "Invalid feed name or public HTTP(S) URL"
		} else if seen[src.URL] {
			res.Error = "Duplicate feed URL"
		} else {
			seen[src.URL] = true
			articles, msg := personalFetch(r.Context(), src)
			res.Error = msg
			cfg := currentArticleSettings()
   res.FeedEntries = len(articles)
   if len(articles) > cfg.MaxArticlesPerFeed { res.NotProcessed = len(articles) - cfg.MaxArticlesPerFeed }
   decisions := classifyFeedArticlesWithAgeDetailed(r.Context(), s.reader, articles, req.MaxArticleAgeDays)
			for i, a := range articles {
				if i >= cfg.MaxArticlesPerFeed {
					break
				}
				if decisions[i].Classification != "article" && classificationAvailability(decisions[i]) != "unavailable" {
      res.Filtered = append(res.Filtered, filteredArticle{Title:a.Title, URL:a.URL, Source:src.Name, Reason:decisions[i].Classification, Published:a.Published.Format(time.RFC3339)})
					continue
				}
				date := "Date unavailable"
                warning := ""
                if !a.Published.IsZero() {
                    date = a.Published.Local().Format("Jan 2, 3:04 PM")
                } else if strings.TrimSpace(a.RawDate) != "" {
                    date = a.RawDate
                    warning = "Unrecognized date format"
                }
				res.Articles = append(res.Articles, personalArticle{Title: a.Title, URL: a.URL, Open: s.reader + "/?url=" + url.QueryEscape(a.URL), Source: src.Name, Category: src.Category, Date: date, DateWarning: warning, Published: a.Published.Format(time.RFC3339), Image: scraperThumbnailURL(s.reader, a.URL), FeedImage: safeFeedImage(a.Image), ContentStatus: decisions[i].ContentStatus, AvailabilityStatus: classificationAvailability(decisions[i])})
			}
			if msg == "" && len(res.Articles) == 0 && res.FeedEntries == 0 {
				res.Error = "Feed contains no usable article links"
			}
		}
		results = append(results, res)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_ = json.NewEncoder(w).Encode(results)
}

// scraperThumbnailURL delegates discovery, retrieval, resizing and caching to ScraperGo.
// The browser never requests publisher image URLs from NewsGo's article cards.
func scraperThumbnailURL(reader, article string) string {
	article = strings.TrimSpace(article)
	if u, err := url.Parse(article); err == nil {
		u.Fragment = ""
		u.RawFragment = ""
		article = u.String()
	}
	if !safeURL(article) {
		return ""
	}
	return strings.TrimRight(reader, "/") + "/api/thumbnail?" + url.Values{"url": {article}, "width": {"256"}}.Encode()
}

// RSS-supplied public images remain available when the article's HTML is blocked.
// Restrict browser fallback to HTTPS image URLs with no credentials.
func safeFeedImage(raw string) string {
 u, err := url.Parse(strings.TrimSpace(raw))
 if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil {return ""}
 return u.String()
}

// Classification is authoritative only for a confirmed homepage/index. A failed
// lookup cannot delete a legitimate headline or prevent opening ScraperGo.
type classificationOutcome struct {
	Classification string `json:"classification"`
 ContentStatus string `json:"contentStatus,omitempty"`
 ContentReason string `json:"contentReason,omitempty"`
 Reason string `json:"reason,omitempty"`
}
type newsClassificationEntry struct {
	value  classificationOutcome
	expiry time.Time
}

var newsClassificationCache = struct {
	sync.Mutex
	entries map[string]newsClassificationEntry
}{entries: make(map[string]newsClassificationEntry)}
var newsClassificationSlots = make(chan struct{}, 16)

func validArticleAgeDays(n int) bool { return n==0 || n==1 || n==3 || n==7 || n==14 || n==30 || n==60 }

func articleTooOld(published time.Time, days int, now time.Time) bool {
 return days > 0 && !published.IsZero() && published.Before(now.AddDate(0,0,-days))
}
func classifyFeedArticles(ctx context.Context, reader string, articles []feeds.Article) []string {
 return classifyFeedArticlesWithAge(ctx, reader, articles, 0)
}
func classifyFeedArticlesWithAge(ctx context.Context, reader string, articles []feeds.Article, days int) []string {
 detail := classifyFeedArticlesWithAgeDetailed(ctx, reader, articles, days)
 out := make([]string,len(detail))
 for i, v := range detail { out[i]=v.Classification }
 return out
}
func classifyFeedArticlesWithAgeDetailed(ctx context.Context, reader string, articles []feeds.Article, days int) []classificationOutcome {
 out := make([]classificationOutcome,len(articles))
 limiter := make(chan struct{},currentArticleSettings().ClassificationConcurrency)
 var wg sync.WaitGroup
 for i,a := range articles {
  if i >= currentArticleSettings().MaxArticlesPerFeed {break}
  if articleTooOld(a.Published,days,time.Now()) {out[i].Classification="too-old";continue}
  wg.Add(1)
  go func(i int, article string) {
   defer wg.Done()
   select {case limiter<-struct{}{}: defer func(){<-limiter}();case <-ctx.Done():out[i].Classification="unknown";return}
   out[i]=newsClassifyDetail(ctx,reader,article)
  }(i,a.URL)
 }
 wg.Wait()
 return out
}
func newsClassify(ctx context.Context, reader, article string) string {
 return newsClassifyDetail(ctx,reader,article).Classification
}
func newsClassifyDetail(ctx context.Context, reader, article string) classificationOutcome {
 article = normalizeClassificationURL(article)
	if !safeURL(article) {
		return classificationOutcome{Classification:"unknown"}
	}
	key := reader + "|" + article
	newsClassificationCache.Lock()
	entry, ok := newsClassificationCache.entries[key]
	newsClassificationCache.Unlock()
	if ok && time.Now().Before(entry.expiry) {
		return entry.value
	}
	select {
	case newsClassificationSlots <- struct{}{}:
		defer func() { <-newsClassificationSlots }()
	case <-ctx.Done():
		return classificationOutcome{Classification:"unknown"}
	}
	requestCtx, cancel := context.WithTimeout(ctx, time.Duration(currentArticleSettings().ClassificationTimeoutSeconds)*time.Second)
	defer cancel()
	body, _ := json.Marshal(map[string]string{"url": article})
	req, err := http.NewRequestWithContext(requestCtx, http.MethodPost, strings.TrimRight(reader, "/")+"/api/classify", bytes.NewReader(body))
	if err != nil {
		return classificationOutcome{Classification:"unknown"}
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Duration(currentArticleSettings().ClassificationTimeoutSeconds)*time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return classificationOutcome{Classification:"unknown"}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return classificationOutcome{Classification:"unknown"}
	}
	var result classificationOutcome
	if json.NewDecoder(io.LimitReader(resp.Body, 2048)).Decode(&result) != nil {
		return classificationOutcome{Classification:"unknown"}
	}
	if result.ContentStatus != "partial" && result.ContentStatus != "complete" {result.ContentStatus="unknown"}
 if result.Classification != "homepage" && result.Classification != "unreadable" && result.Classification != "article" && result.Classification != "unavailable" {
		return classificationOutcome{Classification:"unknown"}
	}
	newsClassificationCache.Lock()
	if len(newsClassificationCache.entries) >= 2048 {
		for k := range newsClassificationCache.entries {
			delete(newsClassificationCache.entries, k)
			break
		}
	}
	ttl := 30 * time.Minute
	if result.Classification == "unavailable" {
		ttl = time.Minute
	}
	newsClassificationCache.entries[key] = newsClassificationEntry{result, time.Now().Add(ttl)}
	newsClassificationCache.Unlock()
	return result
}

// Confirmed publisher denial is distinct from unknown availability and completeness.
func classificationAvailability(result classificationOutcome) string {
 if result.Classification == "unavailable" && result.Reason == "publisher_access_denied" { return "unavailable" }
 return ""
}

// HTTP fragments are browser-only and must not invalidate otherwise safe article URLs.
// Keep the original feed URL for navigation; canonicalize only classification input.
func normalizeClassificationURL(raw string) string {
 u, err := url.Parse(raw)
 if err != nil { return raw }
 u.Fragment = ""
 u.RawFragment = ""
 return u.String()
}
