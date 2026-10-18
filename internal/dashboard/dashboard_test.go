package dashboard

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/richard-mauri/newsgo/internal/feeds"
)

func TestRefreshAndFiltering(t *testing.T) {
	feed := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<rss><channel><item><title>Alpha &amp; beta</title><link>https://example.org/article</link><pubDate>Sun, 18 Oct 2026 18:00:00 GMT</pubDate></item></channel></rss>`)
	}))
	defer feed.Close()
	s, err := New([]feeds.Source{{Name: "One", Category: "World", URL: feed.URL}, {Name: "Two", Category: "Tech", URL: feed.URL + "/other"}}, "http://127.0.0.1:8083", &http.Client{Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	s.Refresh(context.Background())
	for _, tc := range []struct {
		path   string
		want   string
		absent string
	}{{"/", "Alpha &amp; beta", "xyz"}, {"/?source=One", "Alpha &amp; beta", "No headlines match"}, {"/?category=Tech", "No headlines match", "Alpha &amp; beta"}, {"/health", "ok", "Alpha"}} {
		// The duplicate URL is retained once, from the first configured source.
		w := httptest.NewRecorder()
		s.ServeHTTP(w, httptest.NewRequest("GET", tc.path, nil))
		if !strings.Contains(w.Body.String(), tc.want) {
			t.Fatalf("%s missing %q: %s", tc.path, tc.want, w.Body.String())
		}
		if strings.Contains(w.Body.String(), tc.absent) {
			t.Fatalf("%s contains %q", tc.path, tc.absent)
		}
	}
	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if strings.Count(w.Body.String(), "class=\"story-title\"") != 1 {
		t.Fatalf("expected deduplicated entry: %s", w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "target=\"newsgo-scrapergo\"") {
		t.Fatal("article link missing reusable ScraperGo tab target")
	}
	if !strings.Contains(w.Body.String(), "?url=https%3A%2F%2Fexample.org%2Farticle") {
		t.Fatalf("no encoded reader link: %s", w.Body.String())
	}
}
func TestReaderURLValidation(t *testing.T) {
	src := []feeds.Source{{Name: "One", URL: "https://example.org/rss"}}
	for _, v := range []string{"javascript:alert(1)", "http://", "http://name:pass@host"} {
		if _, e := New(src, v, http.DefaultClient); e == nil {
			t.Errorf("accepted %q", v)
		}
	}
}
