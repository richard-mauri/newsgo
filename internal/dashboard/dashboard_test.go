package dashboard

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/richard-mauri/newsgo/internal/feeds"
)

func TestCatalogDoesNotAutoRenderUnsubscribedHeadlines(t *testing.T) {
	s, err := New([]feeds.Source{{Name: "Catalog Source", Category: "World", URL: "https://example.org/rss"}}, "http://127.0.0.1:8083", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	body := w.Body.String()
	if strings.Contains(body, `id="shared-headlines"`) {
		t.Fatal("unsubscribed articles rendered into dashboard")
	}
	if !strings.Contains(body, `data-feed-url="https://example.org/rss"`) {
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
