package feeds

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Source struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	URL      string `json:"url"`
}
type Report struct {
	Name          string         `json:"name"`
	Category      string         `json:"category"`
	URL           string         `json:"url"`
	Status        string         `json:"status"`
	HTTPStatus    int            `json:"http_status,omitempty"`
	Format        string         `json:"format,omitempty"`
	Items         int            `json:"items"`
	Latest        string         `json:"latest,omitempty"`
	AgeHours      *float64       `json:"age_hours,omitempty"`
	WithTitle     int            `json:"with_title"`
	WithLink      int            `json:"with_link"`
	WithDate      int            `json:"with_date"`
	Error         string         `json:"error,omitempty"`
	Compatibility *Compatibility `json:"scrapergo,omitempty"`
}
type item struct {
	title, link, image string
	date               time.Time
}
type rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Items []struct {
			Title string `xml:"title"`
			Link  struct {
				Value string `xml:",chardata"`
				Href  string `xml:"href,attr"`
				Rel   string `xml:"rel,attr"`
			} `xml:"link"`
			GUID struct {
				Value       string `xml:",chardata"`
				IsPermaLink string `xml:"isPermaLink,attr"`
			} `xml:"guid"`
			Description    string `xml:"description"`
			ContentEncoded string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
			PubDate        string `xml:"pubDate"`
			DCDate         string `xml:"date"`
			Enclosure      struct {
				URL  string `xml:"url,attr"`
				Type string `xml:"type,attr"`
			} `xml:"enclosure"`
			MediaThumbnail []struct {
				URL string `xml:"url,attr"`
			} `xml:"http://search.yahoo.com/mrss/ thumbnail"`
			MediaContent []struct {
				URL    string `xml:"url,attr"`
				Type   string `xml:"type,attr"`
				Medium string `xml:"medium,attr"`
			} `xml:"http://search.yahoo.com/mrss/ content"`
			ITunesImage struct {
				Href string `xml:"href,attr"`
			} `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd image"`
		} `xml:"item"`
	} `xml:"channel"`
}
type atom struct {
	XMLName xml.Name `xml:"feed"`
	Entries []struct {
		Title string `xml:"title"`
		Links []struct {
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
		} `xml:"link"`
		Summary        string `xml:"summary"`
		Content        string `xml:"http://www.w3.org/2005/Atom content"`
		Updated        string `xml:"updated"`
		Published      string `xml:"published"`
		MediaThumbnail []struct {
			URL string `xml:"url,attr"`
		} `xml:"http://search.yahoo.com/mrss/ thumbnail"`
		MediaContent []struct {
			URL    string `xml:"url,attr"`
			Type   string `xml:"type,attr"`
			Medium string `xml:"medium,attr"`
		} `xml:"http://search.yahoo.com/mrss/ content"`
	} `xml:"entry"`
}
type rdf struct {
	XMLName xml.Name `xml:"RDF"`
	Items   []struct {
		Title string `xml:"title"`
		Link  string `xml:"link"`
		Date  string `xml:"date"`
	} `xml:"item"`
}
type jsonFeed struct {
	Version string `json:"version"`
	Items   []struct {
		Title       string `json:"title"`
		URL         string `json:"url"`
		ExternalURL string `json:"external_url"`
		Published   string `json:"date_published"`
		Modified    string `json:"date_modified"`
		Image       string `json:"image"`
		BannerImage string `json:"banner_image"`
	} `json:"items"`
}

func Parse(data []byte) (string, []item, error) {
	var j jsonFeed
	if json.Unmarshal(data, &j) == nil && strings.Contains(j.Version, "jsonfeed.org/version/") {
		out := make([]item, 0, len(j.Items))
		for _, v := range j.Items {
			l := v.URL
			if l == "" {
				l = v.ExternalURL
			}
			d := v.Published
			if d == "" {
				d = v.Modified
			}
			out = append(out, item{title: v.Title, link: l, image: firstImage(v.Image, v.BannerImage), date: parseDate(d)})
		}
		return "JSON Feed", out, nil
	}
	var root struct{ XMLName xml.Name }
	if err := xml.Unmarshal(data, &root); err != nil {
		return "", nil, fmt.Errorf("invalid feed: %w", err)
	}
	switch strings.ToLower(root.XMLName.Local) {
	case "rss":
		var f rss
		if err := xml.Unmarshal(data, &f); err != nil {
			return "", nil, err
		}
		out := make([]item, 0, len(f.Channel.Items))
		for _, v := range f.Channel.Items {
			d := v.PubDate
			if d == "" {
				d = v.DCDate
			}
			l := strings.TrimSpace(v.Link.Value)
			if !ValidURL(l) && (v.Link.Rel == "" || v.Link.Rel == "alternate") && ValidURL(v.Link.Href) {
				l = v.Link.Href
			}
			if !ValidURL(l) && !strings.EqualFold(v.GUID.IsPermaLink, "false") && ValidURL(strings.TrimSpace(v.GUID.Value)) {
				l = strings.TrimSpace(v.GUID.Value)
			}
			image := ""
			for _, m := range v.MediaThumbnail {
				image = firstImage(image, m.URL)
			}
			for _, m := range v.MediaContent {
				if m.Medium == "image" || strings.HasPrefix(m.Type, "image/") {
					image = firstImage(image, m.URL)
				}
			}
			if strings.HasPrefix(v.Enclosure.Type, "image/") {
				image = firstImage(image, v.Enclosure.URL)
			}
			image = firstImage(image, v.ITunesImage.Href, imageFromHTML(v.ContentEncoded, l), imageFromHTML(v.Description, l))
			out = append(out, item{title: v.Title, link: l, image: image, date: parseDate(d)})
		}
		return "RSS", out, nil
	case "feed":
		var f atom
		if err := xml.Unmarshal(data, &f); err != nil {
			return "", nil, err
		}
		out := make([]item, 0, len(f.Entries))
		for _, v := range f.Entries {
			l := ""
			for _, link := range v.Links {
				if link.Rel == "" || link.Rel == "alternate" {
					l = link.Href
					break
				}
			}
			d := v.Published
			if d == "" {
				d = v.Updated
			}
			image := ""
			for _, m := range v.MediaThumbnail {
				image = firstImage(image, m.URL)
			}
			for _, m := range v.MediaContent {
				if m.Medium == "image" || strings.HasPrefix(m.Type, "image/") {
					image = firstImage(image, m.URL)
				}
			}
			image = firstImage(image, imageFromHTML(v.Content, l), imageFromHTML(v.Summary, l))
			out = append(out, item{title: v.Title, link: l, image: image, date: parseDate(d)})
		}
		return "Atom", out, nil
	case "rdf":
		var f rdf
		if err := xml.Unmarshal(data, &f); err != nil {
			return "", nil, err
		}
		out := make([]item, 0, len(f.Items))
		for _, v := range f.Items {
			out = append(out, item{title: v.Title, link: v.Link, date: parseDate(v.Date)})
		}
		return "RSS RDF", out, nil
	default:
		return "", nil, fmt.Errorf("unsupported XML root <%s>", root.XMLName.Local)
	}
}

// Embedded feed HTML often contains a lead photograph without Media RSS tags.
// Extract the first HTTPS image without downloading an article or guessing a replacement.
var imgTagRE = regexp.MustCompile(`(?is)<img\b[^>]*>`)
var imgAttributeRE = regexp.MustCompile(`(?is)\b(src|data-src|data-original|data-lazy-src)\s*=\s*["']([^"']+)["']`)

func imageFromHTML(content, article string) string {
	base, _ := url.Parse(article)
	for _, tag := range imgTagRE.FindAllString(content, 25) {
		for _, a := range imgAttributeRE.FindAllStringSubmatch(tag, -1) {
			raw := strings.TrimSpace(html.UnescapeString(a[2]))
			u, err := url.Parse(raw)
			if err != nil || base == nil {
				continue
			}
			resolved := base.ResolveReference(u)
			if out := firstImage(resolved.String()); out != "" {
				return out
			}
		}
	}
	return ""
}

// Only absolute HTTPS image links are used: mixed-content HTTP images would not load
// on the public HTTPS dashboard. Never guess an image from article content.
func firstImage(values ...string) string {
	for _, raw := range values {
		raw = strings.TrimSpace(raw)
		if !strings.HasPrefix(raw, "https://") {
			continue
		}
		if u, err := url.Parse(raw); err == nil && u.Hostname() != "" && u.User == nil {
			return raw
		}
	}
	return ""
}
func parseDate(s string) time.Time {
	for _, layout := range []string{time.RFC3339, time.RFC3339Nano, time.RFC1123Z, time.RFC1123, time.RFC822Z, time.RFC822, "Mon, 2 Jan 2006 15:04:05 MST", "2006-01-02T15:04:05Z0700", "2006-01-02 15:04:05", "2006-01-02"} {
		if t, e := time.Parse(layout, strings.TrimSpace(s)); e == nil {
			return t
		}
	}
	return time.Time{}
}
func ValidURL(raw string) bool {
	u, e := url.Parse(raw)
	return e == nil && (u.Scheme == "https" || u.Scheme == "http") && u.Hostname() != "" && u.User == nil
}
func Check(ctx context.Context, client *http.Client, s Source, now time.Time) Report {
	r := Report{Name: s.Name, Category: s.Category, URL: s.URL, Status: "error"}
	if !ValidURL(s.URL) {
		r.Error = "invalid HTTP(S) feed URL"
		return r
	}
	req, e := http.NewRequestWithContext(ctx, http.MethodGet, s.URL, nil)
	if e != nil {
		r.Error = e.Error()
		return r
	}
	req.Header.Set("User-Agent", "NewsGo-FeedCheck/0.1 (+personal feed validation)")
	req.Header.Set("Accept", "application/rss+xml, application/atom+xml, application/feed+json, application/xml, text/xml, application/json;q=0.8")
	resp, e := client.Do(req)
	if e != nil {
		r.Error = e.Error()
		return r
	}
	defer resp.Body.Close()
	r.HTTPStatus = resp.StatusCode
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		r.Status = "http_error"
		r.Error = resp.Status
		return r
	}
	data, e := io.ReadAll(io.LimitReader(resp.Body, 4<<20+1))
	if e != nil {
		r.Error = e.Error()
		return r
	}
	if len(data) > 4<<20 {
		r.Error = "feed exceeds 4 MiB limit"
		return r
	}
	format, items, e := Parse(data)
	if e != nil {
		r.Status = "invalid_feed"
		r.Error = e.Error()
		return r
	}
	r.Format = format
	r.Items = len(items)
	var latest time.Time
	for _, v := range items {
		if strings.TrimSpace(v.title) != "" {
			r.WithTitle++
		}
		if ValidURL(v.link) {
			r.WithLink++
		}
		if !v.date.IsZero() {
			r.WithDate++
			if v.date.After(latest) {
				latest = v.date
			}
		}
	}
	if !latest.IsZero() {
		r.Latest = latest.UTC().Format(time.RFC3339)
		h := now.Sub(latest).Hours()
		r.AgeHours = &h
	}
	switch {
	case len(items) == 0:
		r.Status = "empty"
	case r.WithTitle == 0 || r.WithLink == 0:
		r.Status = "poor_metadata"
	case latest.IsZero():
		r.Status = "undated"
	case now.Sub(latest) > 72*time.Hour:
		r.Status = "stale"
	default:
		r.Status = "healthy"
	}
	return r
}
func ValidateSources(sources []Source) error {
	if len(sources) == 0 {
		return errors.New("no feeds configured")
	}
	seen := map[string]bool{}
	for _, s := range sources {
		if s.Name == "" || !ValidURL(s.URL) {
			return fmt.Errorf("invalid source %q", s.Name)
		}
		if seen[s.URL] {
			return fmt.Errorf("duplicate feed URL: %s", s.URL)
		}
		seen[s.URL] = true
	}
	return nil
}
func SortReports(rs []Report) { sort.Slice(rs, func(i, j int) bool { return rs[i].Name < rs[j].Name }) }

// ArticleURLs fetches one feed and returns link candidates in feed order.
func ArticleURLs(ctx context.Context, client *http.Client, feedURL string) ([]string, error) {
	if !ValidURL(feedURL) {
		return nil, fmt.Errorf("invalid feed URL")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "NewsGo-FeedCheck/0.1 (+personal feed validation)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("feed returned HTTP %d", resp.StatusCode)
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, (4<<20)+1))
	if err != nil {
		return nil, err
	}
	if len(b) > 4<<20 {
		return nil, fmt.Errorf("feed exceeds 4 MiB limit")
	}
	_, items, err := Parse(b)
	if err != nil {
		return nil, err
	}
	var urls []string
	for _, v := range items {
		if ValidURL(v.link) {
			urls = append(urls, v.link)
		}
	}
	return urls, nil
}
