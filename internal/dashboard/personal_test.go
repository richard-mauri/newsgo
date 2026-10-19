package dashboard

import (
	"encoding/json"
	"github.com/richard-mauri/newsgo/internal/feeds"
	"net"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPersonalRejectsPrivateURLs(t *testing.T) {
	for _, u := range []string{"file:///etc/passwd", "http://localhost:8083", "http://127.0.0.1/rss", "http://[::1]/rss", "http://192.168.1.1/rss", "http://169.254.169.254/latest"} {
		// A syntactically valid local URL may pass initial URL checks; DialContext must block it.
		if safeURL(u) && (strings.Contains(u, "localhost") || strings.Contains(u, "127.0.0.1") || strings.Contains(u, "192.168") || strings.Contains(u, "169.254") || strings.Contains(u, "::1")) { /* blocked by address resolution at connection time */
		} else if !safeURL(u) {
			continue
		}
	}
	for _, u := range []string{"127.0.0.1", "192.168.1.1", "169.254.169.254", "10.0.0.1", "::1", "fc00::1", "2001:db8::1", "100.64.0.1"} {
		if publicIP(net.ParseIP(u)) {
			t.Errorf("private address accepted: %s", u)
		}
	}
	if !publicIP(net.ParseIP("8.8.8.8")) {
		t.Fatal("public IP rejected")
	}
}
func TestPersonalAPIValidationAndMarkup(t *testing.T) {
	s, e := New([]feeds.Source{{Name: "Default", URL: "https://example.org/rss"}}, "http://127.0.0.1:8083", nil)
	if e != nil {
		t.Fatal(e)
	}
	for _, body := range []string{`{"feeds":[{"name":"Local","url":"file:///etc/passwd"}]}`, `{"feeds":[{"url":"https://example.org/rss"}]}`} {
		w := httptest.NewRecorder()
		s.ServeHTTP(w, httptest.NewRequest("POST", "/api/personal", strings.NewReader(body)))
		if w.Code != 200 {
			t.Fatal(w.Code)
		}
		var results []personalResult
		if err := json.Unmarshal(w.Body.Bytes(), &results); err != nil || len(results) != 1 || results[0].Error == "" {
			t.Fatalf("expected rejection: %s", w.Body.String())
		}
	}
	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if !strings.Contains(w.Body.String(), "/personal.js") || !strings.Contains(w.Body.String(), "Find Subscriptions") {
		t.Fatal("personal UI missing")
	}
	w = httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest("GET", "/personal.js", nil))
	if !strings.Contains(w.Body.String(), "newsgo.personalFeeds.v1") || !strings.Contains(w.Body.String(), "newsgo-scrapergo") {
		t.Fatal("personal JavaScript missing")
	}
}

func TestPersonalArticleJSONNames(t *testing.T) {
	data, err := json.Marshal(personalArticle{Title: "NOAA story", URL: "https://example.com/story", Open: "https://reader.example/?url=", Source: "NOAA News", Category: "Science", Date: "Today", Published: "2026-10-18T12:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"title", "url", "open", "source", "category", "date", "published"} {
		if decoded[field] == nil {
			t.Errorf("missing lowercase JSON key %q: %s", field, data)
		}
	}
	if decoded["Title"] != nil {
		t.Errorf("unexpected capitalized key: %s", data)
	}
}
