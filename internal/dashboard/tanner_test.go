package dashboard

import (
	"github.com/richard-mauri/newsgo/internal/feeds"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTannerHeaderAndImage(t *testing.T) {
	s, err := New([]feeds.Source{{Name: "Test", Category: "Science", URL: "https://example.com/rss"}}, "http://127.0.0.1:8083", http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRecorder()
	s.ServeHTTP(r, httptest.NewRequest("GET", "/", nil))
	if r.Code != http.StatusOK || !strings.Contains(r.Body.String(), `alt="Tanner, the puppy"`) || !strings.Contains(r.Body.String(), "V0.2.7") {
		t.Fatalf("header unavailable: %d", r.Code)
	}
	r = httptest.NewRecorder()
	s.ServeHTTP(r, httptest.NewRequest("GET", "/tanner.jpg", nil))
	if r.Code != http.StatusOK || r.Header().Get("Content-Type") != "image/jpeg" || r.Body.Len() < 1000 {
		t.Fatalf("Tanner image unavailable: status=%d bytes=%d", r.Code, r.Body.Len())
	}
}
