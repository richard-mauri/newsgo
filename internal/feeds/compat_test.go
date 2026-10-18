package feeds

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRSSLinkFallback(t *testing.T) {
	cases := []struct{ name, xml, expected string }{
		{"guid", `<rss><channel><item><title>A</title><guid>https://example.com/a</guid></item></channel></rss>`, "https://example.com/a"},
		{"explicit false", `<rss><channel><item><title>A</title><guid isPermaLink="false">https://example.com/a</guid></item></channel></rss>`, ""},
		{"atom style", `<rss><channel><item><title>A</title><link href="https://example.com/b" rel="alternate"/></item></channel></rss>`, "https://example.com/b"},
		{"regular", `<rss><channel><item><title>A</title><link>https://example.com/c</link></item></channel></rss>`, "https://example.com/c"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, items, err := Parse([]byte(tc.xml))
			if err != nil || len(items) != 1 || items[0].link != tc.expected {
				t.Fatalf("items=%+v err=%v", items, err)
			}
		})
	}
}
func TestScraperGoSample(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/read" || r.Method != "POST" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL)
			return
		}
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Error(err)
			return
		}
		switch payload["url"] {
		case "https://example.com/ok":
			json.NewEncoder(w).Encode(map[string]string{"title": "A story"})
		case "https://example.com/blocked":
			w.WriteHeader(403)
			json.NewEncoder(w).Encode(map[string]string{"kind": "publisher_access_denied"})
		default:
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(map[string]string{"kind": "extraction_error"})
		}
	}))
	defer srv.Close()
	endpoint, err := ScraperGoEndpoint(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	result := SampleScraperGo(context.Background(), srv.Client(), endpoint, []string{"https://example.com/ok", "https://example.com/blocked", "https://example.com/fail", "https://example.com/fail"}, 3)
	if result.Tested != 3 || result.Readable != 1 || result.Restricted != 1 || result.Other != 1 {
		t.Fatalf("%+v", result)
	}
}
