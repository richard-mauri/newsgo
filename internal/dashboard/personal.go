package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/richard-mauri/newsgo/internal/feeds"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type personalRequest struct {
	Feeds []feeds.Source `json:"feeds"`
}
type personalResult struct {
	Source   feeds.Source      `json:"source"`
	Articles []personalArticle `json:"articles,omitempty"`
	Error    string            `json:"error,omitempty"`
}
type personalArticle struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Open      string `json:"open"`
	Source    string `json:"source"`
	Category  string `json:"category"`
	Date      string `json:"date"`
	Published string `json:"published"`
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

func personalFetch(ctx context.Context, source feeds.Source) ([]feeds.Article, string) {
	personalCache.Lock()
	cached, ok := personalCache.entries[source.URL]
	personalCache.Unlock()
	if ok && time.Since(cached.at) < 10*time.Minute {
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
	if ok && time.Since(cached.at) < 10*time.Minute {
		return cached.articles, cached.err
	}
	a, err := feeds.FetchArticles(ctx, safeFeedClient(), source)
	msg := ""
	if err != nil {
		msg = err.Error()
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Feeds) > 8 {
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
			for i, a := range articles {
				if i >= 60 {
					break
				}
				date := "Date unavailable"
				if !a.Published.IsZero() {
					date = a.Published.Local().Format("Jan 2, 3:04 PM")
				}
				res.Articles = append(res.Articles, personalArticle{a.Title, a.URL, s.reader + "/?url=" + url.QueryEscape(a.URL), src.Name, src.Category, date, a.Published.Format(time.RFC3339)})
			}
			if msg == "" && len(res.Articles) == 0 {
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
