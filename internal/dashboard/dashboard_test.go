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

func TestCatalogDoesNotAutoRenderUnsubscribedHeadlines(t *testing.T) {
	feed := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<rss><channel><item><title>Unsubscribed article</title><link>https://example.org/a</link></item></channel></rss>`)
	}))
	defer feed.Close()
	s, err := New([]feeds.Source{{Name: "Catalog Source", Category: "World", URL: feed.URL}}, "http://127.0.0.1:8083", &http.Client{Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	// Even if an obsolete path refreshes the cached collection, the browser
	// receives catalog metadata only; subscribed feeds populate /api/personal.
	s.Refresh(context.Background())
	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	body := w.Body.String()
	if strings.Contains(body, "Unsubscribed article") || strings.Contains(body, `id="shared-headlines"`) {
		t.Fatal("unsubscribed articles rendered into dashboard")
	}
	if !strings.Contains(body, `data-feed-url="`+feed.URL+`"`) {
		t.Fatal("catalog feed URL missing")
	}
	if !strings.Contains(body, `id="personal-articles"`) {
		t.Fatal("unified article collection missing")
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
