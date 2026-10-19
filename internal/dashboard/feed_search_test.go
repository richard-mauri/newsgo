package dashboard

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseFeedlyResults(t *testing.T) {
	body := []byte(`{"results":[{"title":"Sailing","feedId":"feed/https://example.org/rss","website":"https://example.org"},{"title":"Duplicate","id":"feed/https://example.org/rss"},{"title":"Private","feedId":"feed/http://127.0.0.1/private"}]}`)
	got, err := feedlyResults(body)
	if err != nil || len(got) != 1 || got[0].Name != "Sailing" {
		t.Fatalf("unexpected %#v %v", got, err)
	}
}
func TestSearchRejectsBadQuery(t *testing.T) {
	s := &State{}
	w := httptest.NewRecorder()
	s.serveFeedSearch(w, httptest.NewRequest("GET", "/api/search-feeds?q=x", nil))
	if w.Code != 400 {
		t.Fatalf("HTTP %d", w.Code)
	}
	var r feedSearchResponse
	if err := json.Unmarshal(w.Body.Bytes(), &r); err != nil || !strings.Contains(r.Error, "topic") {
		t.Fatalf("%v %+v", err, r)
	}
}
