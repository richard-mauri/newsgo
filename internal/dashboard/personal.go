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
 "sync/atomic"
	"time"
)

type personalRequest struct {
	Feeds []feeds.Source `json:"feeds"`
 MaxArticleAgeDays int `json:"maxArticleAgeDays,omitempty"`
 RSSMetadataOnly bool `json:"rssMetadataOnly,omitempty"`
 UseClassificationCache *bool `json:"useClassificationCache,omitempty"`
 ExcludedDomains []string `json:"excludedDomains,omitempty"`
}
type personalResult struct {
	Source   feeds.Source      `json:"source"`
	Articles []personalArticle `json:"articles,omitempty"`
	Error    string            `json:"error,omitempty"`
 Filtered []filteredArticle `json:"filtered,omitempty"`
 NotProcessed int `json:"notProcessed,omitempty"`
 FeedEntries int `json:"feedEntries,omitempty"`
 RawFeed *rawFeedDiagnostic `json:"rawFeed,omitempty"`
 FeedDurationMs int64 `json:"feedDurationMs"`
 ClassificationDurationMs int64 `json:"classificationDurationMs"`
 ClassificationCacheHits int64 `json:"classificationCacheHits"`
 ClassificationCacheMisses int64 `json:"classificationCacheMisses"`
 ClassificationCacheBypasses int64 `json:"classificationCacheBypasses"`
 ClassificationCacheDisabled int64 `json:"classificationCacheDisabled"`
 ClassificationCacheTTLMinutes int `json:"classificationCacheTTLMinutes"`
 ClassificationCacheColdMisses int64 `json:"classificationCacheColdMisses"`
 ClassificationCacheExpiredMisses int64 `json:"classificationCacheExpiredMisses"`
 ClassificationNoncacheableUnknown int64 `json:"classificationNoncacheableUnknown"`
 ClassificationNoncacheableUnavailable int64 `json:"classificationNoncacheableUnavailable"`
 ClassificationUpstreamFailures int64 `json:"classificationUpstreamFailures"`
}
type filteredArticle struct {
 Title string `json:"title"`
 URL string `json:"url"`
 Source string `json:"source"`
 Reason string `json:"reason"`
 Published string `json:"published,omitempty"`
}
type rawFeedDiagnostic struct {
 URL string `json:"url"`
 Body []byte `json:"body,omitempty"`
 ContentType string `json:"contentType,omitempty"`
 HTTPStatus int `json:"httpStatus,omitempty"`
 Cached bool `json:"cached"`
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
 raw rawFeedDiagnostic
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

func personalFetch(ctx context.Context, source feeds.Source) ([]feeds.Article, rawFeedDiagnostic, string) {
	if err := ctx.Err(); err != nil {
		return nil, rawFeedDiagnostic{}, "Request canceled"
	}
	personalCache.Lock()
	cached, ok := personalCache.entries[source.URL]
	personalCache.Unlock()
	if ok && time.Since(cached.at) < personalCacheTTL(cached) {
		return cached.articles, rawFeedDiagnostic{URL:source.URL,Body:cached.raw.Body,ContentType:cached.raw.ContentType,HTTPStatus:cached.raw.HTTPStatus,Cached:true}, cached.err
	}
	select {
	case personalSlots <- struct{}{}:
		defer func() { <-personalSlots }()
	case <-ctx.Done():
		return nil, rawFeedDiagnostic{}, "Request canceled"
	}
	// Recheck after waiting for a slot.
	personalCache.Lock()
	cached, ok = personalCache.entries[source.URL]
	personalCache.Unlock()
	if ok && time.Since(cached.at) < personalCacheTTL(cached) {
		return cached.articles, rawFeedDiagnostic{URL:source.URL,Body:cached.raw.Body,ContentType:cached.raw.ContentType,HTTPStatus:cached.raw.HTTPStatus,Cached:true}, cached.err
	}
	a, body, contentType, status, err := feeds.FetchArticlesWithRaw(ctx, safeFeedClient(), source)
 raw := rawFeedDiagnostic{URL:source.URL,Body:body,ContentType:contentType,HTTPStatus:status}
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	// Do not poison the shared cache with a canceled browser request. A
	// timed-out feed request may be retried, but cannot mask later success.
	if ctx.Err() != nil || errors.Is(err, context.Canceled) {
		return a, raw, msg
	}
	personalCache.Lock()
	if len(personalCache.entries) >= 64 {
		for k := range personalCache.entries {
			delete(personalCache.entries, k)
			break
		}
	}
	personalCache.entries[source.URL] = cacheEntry{articles:a,raw:raw,err:msg,at:time.Now()}
	personalCache.Unlock()
	return a, raw, msg
}
func (s *State) servePersonal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not allowed", 405)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	var req personalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Feeds) > 8 || !validArticleAgeDays(req.MaxArticleAgeDays) || !validPersonalExcludedDomains(req.ExcludedDomains) {
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
			feedStarted:=time.Now()
            articles, raw, msg := personalFetch(r.Context(), src)
   res.RawFeed=&raw
            res.FeedDurationMs=time.Since(feedStarted).Milliseconds()
			res.Error = msg
			cfg := currentArticleSettings()
   res.FeedEntries = len(articles)
   if len(articles) > cfg.MaxArticlesPerFeed { res.NotProcessed = len(articles) - cfg.MaxArticlesPerFeed }
   classifyStarted:=time.Now()
   decisions := make([]string,len(articles))
   // RSS metadata-only controls presentation, not article eligibility.
   // Both reading modes require ScraperGo-confirmed articles.
   var cacheStats classificationCacheStats
   decisions = classifyFeedArticlesWithExclusionsStats(r.Context(), s.reader, articles, req.MaxArticleAgeDays, req.UseClassificationCache == nil || *req.UseClassificationCache, req.ExcludedDomains, &cacheStats)
   res.ClassificationCacheHits=cacheStats.hits.Load()
   res.ClassificationCacheMisses=cacheStats.misses.Load()
   res.ClassificationCacheBypasses=cacheStats.bypasses.Load()
   res.ClassificationCacheDisabled=cacheStats.disabled.Load()
   res.ClassificationCacheTTLMinutes=classificationTTLMinutes()
   res.ClassificationCacheColdMisses=cacheStats.cold.Load()
   res.ClassificationCacheExpiredMisses=cacheStats.expired.Load()
   res.ClassificationNoncacheableUnknown=cacheStats.unknown.Load()
   res.ClassificationNoncacheableUnavailable=cacheStats.unavailable.Load()
   res.ClassificationUpstreamFailures=cacheStats.failures.Load()
   res.ClassificationDurationMs=time.Since(classifyStarted).Milliseconds()
			for i, a := range articles {
				if i >= cfg.MaxArticlesPerFeed {
					break
				}
				if decisions[i] != "article" {
      res.Filtered = append(res.Filtered, filteredArticle{Title:a.Title, URL:a.URL, Source:src.Name, Reason:decisions[i], Published:a.Published.Format(time.RFC3339)})
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
				res.Articles = append(res.Articles, personalArticle{Title: a.Title, URL: a.URL, Open: s.reader + "/?url=" + url.QueryEscape(a.URL), Source: src.Name, Category: src.Category, Date: date, DateWarning: warning, Published: a.Published.Format(time.RFC3339), Image: func() string {if req.RSSMetadataOnly {if safeRSSImage(a.Image){return a.Image};return ""};return scraperThumbnailURL(s.reader, a.URL)}()})
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

// Classification is authoritative only for a confirmed homepage/index. A failed
// lookup cannot delete a legitimate headline or prevent opening ScraperGo.
type classificationOutcome struct {
	Classification string `json:"classification"`
}
type newsClassificationEntry struct {
	value  string
	expiry time.Time
}

var newsClassificationCache = struct {
	sync.Mutex
	entries map[string]newsClassificationEntry
 path string
}{entries: make(map[string]newsClassificationEntry)}
var newsClassificationSlots = make(chan struct{}, 16)

func clearNewsClassificationCache(){
 newsClassificationCache.Lock()
 newsClassificationCache.entries=make(map[string]newsClassificationEntry)
 _=persistClassificationCacheLocked()
 newsClassificationCache.Unlock()
}

func validArticleAgeDays(n int) bool { return n==0 || n==1 || n==3 || n==7 || n==14 || n==30 || n==60 }

func articleTooOld(published time.Time, days int, now time.Time) bool {
 return days > 0 && !published.IsZero() && published.Before(now.AddDate(0,0,-days))
}
func classifyFeedArticles(ctx context.Context, reader string, articles []feeds.Article) []string {
 return classifyFeedArticlesWithAge(ctx, reader, articles, 0)
}
func classifyFeedArticlesWithAge(ctx context.Context, reader string, articles []feeds.Article, days int) []string {
 return classifyFeedArticlesWithAgeCache(ctx,reader,articles,days,true)
}
func classifyFeedArticlesWithAgeCache(ctx context.Context, reader string, articles []feeds.Article, days int, useCache bool) []string {
 return classifyFeedArticlesWithAgeStats(ctx,reader,articles,days,useCache,nil)
}
type classificationCacheStats struct {hits,misses,bypasses,disabled,cold,expired,unknown,unavailable,failures atomic.Int64}

// Exclusions originate in the browser's local or Homepage-hydrated profile.
// Validate them before use; no user-supplied domain is used for network access.
func validPersonalExcludedDomains(domains []string) bool {
 if len(domains)>128 {return false}
 for _,domain:=range domains {
  if len(domain)==0||len(domain)>253||strings.ToLower(domain)!=domain||strings.ContainsAny(domain,"/:@?#%\\ ")||strings.HasPrefix(domain,".")||strings.HasSuffix(domain,".")||!strings.Contains(domain,".") {return false}
  for _,label:=range strings.Split(domain,".") {
   if len(label)==0||len(label)>63||label[0]=='-'||label[len(label)-1]=='-' {return false}
   for _,c:=range label {if !((c>='a'&&c<='z')||(c>='0'&&c<='9')||c=='-') {return false}}
  }
 }
 return true
}
func excludedArticleURL(raw string, excluded []string) bool {
 u,err:=url.Parse(raw)
 if err!=nil||(u.Scheme!="https"&&u.Scheme!="http")||u.User!=nil {return false}
 host:=strings.TrimSuffix(strings.ToLower(u.Hostname()),".")
 host=strings.TrimPrefix(host,"www.")
 for _,domain:=range excluded {
  if host==domain||strings.HasSuffix(host,"."+domain) {return true}
 }
 return false
}
func classifyFeedArticlesWithAgeStats(ctx context.Context, reader string, articles []feeds.Article, days int, useCache bool, stats *classificationCacheStats) []string {
 return classifyFeedArticlesWithExclusionsStats(ctx,reader,articles,days,useCache,nil,stats)
}
func classifyFeedArticlesWithExclusionsStats(ctx context.Context, reader string, articles []feeds.Article, days int, useCache bool, exclusions []string, stats *classificationCacheStats) []string {
	out := make([]string, len(articles))
 limiter := make(chan struct{}, currentArticleSettings().ClassificationConcurrency)
	var wg sync.WaitGroup
	for i, a := range articles {
		if i >= currentArticleSettings().MaxArticlesPerFeed {
			break
		}
        // Exclusion must precede age handling and every classification/cache
        // lookup. A skipped article is reported but never sent to ScraperGo.
        if excludedArticleURL(a.URL,exclusions) {
            out[i]="publisher exclusion"
            continue
        }
		if articleTooOld(a.Published, days, time.Now()) {
			out[i] = "too-old"
			continue
		}
		wg.Add(1)
		go func(i int, article string) { defer wg.Done(); select { case limiter <- struct{}{}: defer func(){ <-limiter }(); case <-ctx.Done(): out[i]="unknown"; return }; out[i] = newsClassifyWithCacheStats(ctx, reader, article, useCache, stats) }(i, a.URL)
	}
	wg.Wait()
	return out
}
func newsClassify(ctx context.Context, reader, article string) string {
 return newsClassifyWithCache(ctx,reader,article,true)
}
func newsClassifyWithCache(ctx context.Context, reader, article string, useCache bool) string {
 return newsClassifyWithCacheStats(ctx,reader,article,useCache,nil)
}
func newsClassifyWithCacheStats(ctx context.Context, reader, article string, useCache bool, stats *classificationCacheStats) (classification string) {
 miss := false
 upstreamCompleted := false
 defer func() {
  if !miss || stats==nil {return}
  if !upstreamCompleted {stats.failures.Add(1);return}
  switch classification {
  case "unknown": stats.unknown.Add(1)
  case "unavailable": stats.unavailable.Add(1)
  }
 }()
 article = normalizeClassificationURL(article)
	if !safeURL(article) {
		return "unknown"
	}
	key := reader + "|" + article
	newsClassificationCache.Lock()
	entry, ok := newsClassificationCache.entries[key]
	newsClassificationCache.Unlock()
	if useCache && classificationTTLMinutes()>0 && ok && time.Now().Before(entry.expiry) {
  if stats!=nil {stats.hits.Add(1)}
		return entry.value
	}
 if stats!=nil {
  if !useCache {stats.bypasses.Add(1)} else if classificationTTLMinutes()==0 {stats.disabled.Add(1)} else {
   stats.misses.Add(1)
   if ok {stats.expired.Add(1)} else {stats.cold.Add(1)}
   miss=true
  }
 }
	select {
	case newsClassificationSlots <- struct{}{}:
		defer func() { <-newsClassificationSlots }()
	case <-ctx.Done():
		return "unknown"
	}
	requestCtx, cancel := context.WithTimeout(ctx, time.Duration(currentArticleSettings().ClassificationTimeoutSeconds)*time.Second)
	defer cancel()
	body, _ := json.Marshal(map[string]string{"url": article})
	req, err := http.NewRequestWithContext(requestCtx, http.MethodPost, strings.TrimRight(reader, "/")+"/api/classify", bytes.NewReader(body))
	if err != nil {
		return "unknown"
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Duration(currentArticleSettings().ClassificationTimeoutSeconds)*time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "unknown"
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "unknown"
	}
	var result classificationOutcome
	if json.NewDecoder(io.LimitReader(resp.Body, 2048)).Decode(&result) != nil {
		return "unknown"
	}
	if result.Classification != "homepage" && result.Classification != "unreadable" && result.Classification != "article" && result.Classification != "unavailable" && result.Classification != "unknown" {
		return "unknown"
	}
 upstreamCompleted=true
 // Cache only successful classifications, when enabled. Never cache unavailable failures.
 if ttl:=classificationTTLMinutes();useCache && ttl>0 && validCachedClassification(result.Classification) && validClassificationKey(key) {
  newsClassificationCache.Lock()
  if len(newsClassificationCache.entries)>=classificationCacheCapacity {
   // Expired entries go first; otherwise evict the soonest-to-expire entry.
   victim:=""
   var earliest time.Time
   now:=time.Now()
   for k,v:=range newsClassificationCache.entries {
    if !now.Before(v.expiry) {victim=k;break}
    if victim=="" || v.expiry.Before(earliest) {victim=k;earliest=v.expiry}
   }
   if victim!="" {delete(newsClassificationCache.entries,victim)}
  }
  newsClassificationCache.entries[key]=newsClassificationEntry{result.Classification,time.Now().Add(time.Duration(ttl)*time.Minute)}
  _=persistClassificationCacheLocked()
  newsClassificationCache.Unlock()
 }
	return result.Classification
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

// safeRSSImage only exposes HTTPS feed-supplied images; no ScraperGo extraction is requested.
func safeRSSImage(raw string) bool {
 u,e:=url.Parse(strings.TrimSpace(raw));return e==nil&&u.Scheme=="https"&&u.Host!=""&&u.User==nil&&u.Fragment==""
}
