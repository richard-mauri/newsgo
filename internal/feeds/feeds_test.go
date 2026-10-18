package feeds

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	cases := []struct{ name, data, format string }{{"rss", `<rss version="2.0"><channel><item><title>One</title><link>https://example.org/one</link><pubDate>Sun, 18 Oct 2026 10:00:00 GMT</pubDate></item></channel></rss>`, "RSS"}, {"atom", `<feed xmlns="http://www.w3.org/2005/Atom"><entry><title>Two</title><link href="https://example.org/two" rel="alternate"/><updated>2026-10-18T10:00:00Z</updated></entry></feed>`, "Atom"}, {"json", `{"version":"https://jsonfeed.org/version/1.1","items":[{"title":"Three","url":"https://example.org/three","date_published":"2026-10-18T10:00:00Z"}]}`, "JSON Feed"}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			format, items, err := Parse([]byte(tc.data))
			if err != nil || format != tc.format || len(items) != 1 || items[0].date.IsZero() {
				t.Fatalf("%s %+v %v", format, items, err)
			}
		})
	}
}
func TestCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<rss><channel><item><title>Headline</title><link>https://example.org/story</link><pubDate>Sun, 18 Oct 2026 10:00:00 GMT</pubDate></item></channel></rss>`))
	}))
	defer server.Close()
	r := Check(context.Background(), server.Client(), Source{"Test", "World", server.URL}, time.Date(2026, 10, 18, 11, 0, 0, 0, time.UTC))
	if r.Status != "healthy" || r.WithLink != 1 {
		t.Fatalf("%+v", r)
	}
}
func TestInvalid(t *testing.T) {
	_, _, err := Parse([]byte(`<html></html>`))
	if err == nil {
		t.Fatal("expected invalid feed")
	}
	if ValidURL("file:///tmp/file") {
		t.Fatal("file URLs not allowed")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Fatal(err)
	}
}
