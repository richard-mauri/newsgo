package feeds

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Article is a normalized headline; its content remains on the publisher site.
type Article struct {
	Title     string
	URL       string
	Source    string
	Category  string
	Published time.Time
	Image     string
}

// FetchArticles reuses FeedCheck's parser and conservative 4 MiB response limit.
func FetchArticles(ctx context.Context, client *http.Client, src Source) ([]Article, error) {
	if !ValidURL(src.URL) {
		return nil, fmt.Errorf("invalid feed URL")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, src.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "NewsGo/0.2 (+personal RSS reader)")
	req.Header.Set("Accept", "application/rss+xml, application/atom+xml, application/feed+json, application/xml, text/xml, application/json;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, (4<<20)+1))
	if err != nil {
		return nil, err
	}
	if len(b) > 4<<20 {
		return nil, fmt.Errorf("feed exceeds 4 MiB")
	}
	_, entries, err := Parse(b)
	if err != nil {
		return nil, err
	}
	out := make([]Article, 0, len(entries))
	for _, e := range entries {
		if strings.TrimSpace(e.title) == "" || !ValidURL(e.link) {
			continue
		}
		out = append(out, Article{Title: strings.TrimSpace(e.title), URL: strings.TrimSpace(e.link), Source: src.Name, Category: src.Category, Published: e.date, Image: e.image})
	}
	return out, nil
}
