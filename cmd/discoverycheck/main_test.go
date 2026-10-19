package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFeedlyParse(t *testing.T) {
	x, e := parseFeedly([]byte(`{"results":[{"title":"Sailing","feedId":"feed/https://example.com/rss.xml","website":"https://example.com"},{"title":"bad","feedId":"feed/javascript:bad"}]}`))
	if e != nil || len(x) != 1 || x[0].URL != "https://example.com/rss.xml" {
		t.Fatalf("%+v %v", x, e)
	}
}
func TestAmplifierParse(t *testing.T) {
	x, e := parseAmplifier([]byte(`{"results":[{"title":"Sailing","feed_url":"https://example.org/rss"}]}`))
	if e != nil || len(x) != 1 {
		t.Fatalf("%+v %v", x, e)
	}
}
func TestProbe(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("query") != "college football" {
			t.Errorf("query: %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{"results":[{"title":"CFB","feedId":"feed/https://example.com/rss"}]}`))
	}))
	defer srv.Close()
	r := probe(context.Background(), srv.Client(), provider{"Feedly", srv.URL, "query"}, "college football", 5)
	if r.Status != 200 || r.Count != 1 || r.Error != "" {
		t.Fatalf("%+v", r)
	}
}
func TestNoInvalidFeedURLs(t *testing.T) {
	for _, u := range []string{"javascript:alert(1)", "feed/file:///etc/passwd", "ftp://example.com/feed"} {
		if candidateURL(u) != "" {
			t.Fatal(u)
		}
	}
	if !strings.HasPrefix(candidateURL("feed/https://example.com"), "https://") {
		t.Fatal("unexpected URL")
	}
}
