package dashboard

import (
	"encoding/json"
	"github.com/richard-mauri/newsgo/internal/feeds"
	"html"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"

	"strings"
)

type discoverRequest struct {
	URL string `json:"url"`
}
type discoverResponse struct {
	Feeds  []feeds.Source `json:"feeds,omitempty"`
	Counts map[string]int `json:"counts,omitempty"`
	Error  string         `json:"error,omitempty"`
}

var titleTagRE = regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)
var titleMarkupRE = regexp.MustCompile(`(?s)<[^>]+>`)

// A readable display title is derived from ordinary page metadata and the
// requested section URL. It never changes the actual feed endpoint.
func discoveryName(pageHTML []byte, website *url.URL, feedURL string) string {
	host := strings.TrimPrefix(strings.ToLower(website.Hostname()), "www.")
	publisher := strings.TrimSuffix(host, ".com")
	if publisher == host {
		publisher = strings.Split(host, ".")[0]
	}
	if publisher == "" {
		publisher = host
	}
	publisher = strings.ToUpper(publisher[:1]) + publisher[1:]
	var title string
	if m := titleTagRE.FindSubmatch(pageHTML); len(m) > 1 {
		title = strings.TrimSpace(html.UnescapeString(titleMarkupRE.ReplaceAllString(string(m[1]), "")))
	}
	// Prefer concise publisher/section title components over long SEO headlines.
	for _, separator := range []string{" | ", " — ", " - "} {
		if strings.Contains(title, separator) {
			parts := strings.Split(title, separator)
			if len(parts) > 1 && len(strings.TrimSpace(parts[0])) > 2 && len(parts[0]) < 65 {
				title = strings.TrimSpace(parts[0])
				break
			}
		}
	}
	section := strings.Trim(path.Base(strings.TrimRight(website.Path, "/")), " ")
	if section != "" && section != "." && section != "/" && section != "rss-feeds" {
		section = strings.ReplaceAll(strings.ReplaceAll(section, "-", " "), "_", " ")
		section = strings.TrimSpace(section)
		section = strings.Title(section)
	}
	if title == "" || len(title) > 72 || strings.EqualFold(title, publisher) || strings.EqualFold(title, host) {
		title = section
	}
	if title == "" {
		title = publisher
	}
	// If the title already names the publisher, don't repeat it.
	if strings.Contains(strings.ToLower(title), strings.ToLower(publisher)) {
		return title
	}
	return title + " · " + publisher
}
func discoveryCategory(name, websiteURL string) string {
	text := strings.ToLower(name + " " + websiteURL)
	for _, word := range []string{"sports", "football", "soccer", "basketball", "baseball", "hockey", "tennis", "golf", "nfl", "nba", "mlb", "ncaa"} {
		if strings.Contains(text, word) {
			return "Sports"
		}
	}
	return "Other"
}

var linkTagRE = regexp.MustCompile(`(?is)<link\s+[^>]*>`)
var anchorTagRE = regexp.MustCompile(`(?is)<a\s+[^>]*>`)
var htmlAttrRE = regexp.MustCompile(`(?is)([a-z][a-z0-9_-]*)\s*=\s*(?:"([^"]*)"|'([^']*)')`)

func advertisedFeeds(body []byte, base *url.URL) []string {
	out := []string{}
	for _, tag := range linkTagRE.FindAllString(string(body), -1) {
		attrs := map[string]string{}
		for _, m := range htmlAttrRE.FindAllStringSubmatch(tag, -1) {
			v := m[2]
			if v == "" {
				v = m[3]
			}
			attrs[strings.ToLower(m[1])] = html.UnescapeString(v)
		}
		rel := strings.Fields(strings.ToLower(attrs["rel"]))
		alternate := false
		for _, v := range rel {
			if v == "alternate" {
				alternate = true
			}
		}
		typ := strings.ToLower(attrs["type"])
		if !alternate || !(strings.Contains(typ, "rss") || strings.Contains(typ, "atom") || strings.Contains(typ, "feed+json")) {
			continue
		}
		u, e := base.Parse(attrs["href"])
		if e == nil && safeURL(u.String()) {
			out = append(out, u.String())
		}
	}
	return out
}

// linkedFeeds discovers ordinary feed links from a publisher's RSS directory.
// It deliberately requires feed-related labels or URL hints before probing.
func linkedFeeds(body []byte, base *url.URL) []string {
	out := []string{}
	for _, tag := range anchorTagRE.FindAllString(string(body), -1) {
		attrs := map[string]string{}
		for _, m := range htmlAttrRE.FindAllStringSubmatch(tag, -1) {
			v := m[2]
			if v == "" {
				v = m[3]
			}
			attrs[strings.ToLower(m[1])] = html.UnescapeString(v)
		}
		href := attrs["href"]
		hint := strings.ToLower(href + " " + attrs["title"] + " " + attrs["type"] + " " + attrs["aria-label"])
		if !(strings.Contains(hint, "rss") || strings.Contains(hint, "atom") || strings.Contains(hint, "feed") || strings.Contains(hint, "optimized-rss")) {
			continue
		}
		u, err := base.Parse(href)
		if err == nil && safeURL(u.String()) {
			out = append(out, u.String())
		}
	}
	return out
}
func (s *State) serveDiscover(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		_ = json.NewEncoder(w).Encode(discoverResponse{Error: "POST required"})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var req discoverRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil || len(req.URL) > 2048 || !safeURL(req.URL) {
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(discoverResponse{Error: "Enter a public HTTP(S) website URL"})
		return
	}
	client := safeFeedClient()
	base, _ := url.Parse(req.URL)
	candidates := []string{req.URL}
	var pageHTML []byte
	// Prefer explicitly advertised feeds. Probe conventional endpoints only when needed.
	resp, e := client.Get(req.URL)
	if e == nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 && strings.Contains(strings.ToLower(resp.Header.Get("Content-Type")), "html") {
			b, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
			pageHTML = b
			candidates = append(candidates, advertisedFeeds(b, resp.Request.URL)...)
			candidates = append(candidates, linkedFeeds(b, resp.Request.URL)...)
		}
	}
	for _, suffix := range []string{"/feed", "/rss", "/rss.xml", "/atom.xml"} {
		candidates = append(candidates, base.Scheme+"://"+base.Host+suffix)
	}
	seen := map[string]bool{}
	results := []feeds.Source{}
	counts := map[string]int{}
	probes := 0
	for _, raw := range candidates {
		if seen[raw] || !safeURL(raw) {
			continue
		}
		seen[raw] = true
		probes++
		if probes > 24 {
			break
		}
		name := discoveryName(pageHTML, base, raw)
		source := feeds.Source{Name: name, Category: discoveryCategory(name, req.URL), URL: raw}
		items, err := feeds.FetchArticles(r.Context(), client, source)
		if err != nil || len(items) == 0 {
			continue
		}
		results = append(results, source)
		counts[source.URL] = len(items)
		if len(results) >= 8 {
			break
		}
	}
	_ = json.NewEncoder(w).Encode(discoverResponse{Feeds: results, Counts: counts})
}
