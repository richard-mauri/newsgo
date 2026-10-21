package dashboard

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewsClassifyHomepage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/classify" {
			t.Error("wrong path")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"classification":"homepage","readableArticle":false}`))
	}))
	defer srv.Close()
	if got := newsClassify(context.Background(), srv.URL, "https://example.com/test-homepage"); got != "homepage" {
		t.Fatalf("got %q", got)
	}
}
func TestNewsClassifyUnknownPreserved(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "unavailable", 502) }))
	defer srv.Close()
	if got := newsClassify(context.Background(), srv.URL, "https://example.com/unknown"); got != "unknown" {
		t.Fatal(got)
	}
	_ = strings.TrimSpace
}
