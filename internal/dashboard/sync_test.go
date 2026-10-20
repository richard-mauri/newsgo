package dashboard

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func syncCall(t *testing.T, s *SyncStore, method, token string, payload any) *httptest.ResponseRecorder {
	t.Helper()
	var b []byte
	if payload != nil {
		var e error
		b, e = json.Marshal(payload)
		if e != nil {
			t.Fatal(e)
		}
	}
	r := httptest.NewRequest(method, "https://newsgo.themauris.org/api/sync", bytes.NewReader(b))
	r.Header.Set("Authorization", "Bearer "+token)
	if method == http.MethodPut {
		r.Header.Set("Origin", "https://newsgo.themauris.org")
	}
	w := httptest.NewRecorder()
	s.ServeSync(w, r)
	return w
}
func TestSyncRoundTripAcrossInstances(t *testing.T) {
	dir := t.TempDir()
	a, e := NewSyncStore(dir)
	if e != nil {
		t.Fatal(e)
	}
	token, e := a.IssueToken("richard", "first")
	if e != nil {
		t.Fatal(e)
	}
	if strings.Contains(string(mustRead(t, filepath.Join(dir, "approved-tokens.json"))), token) {
		t.Fatal("plaintext token on disk")
	}
	if w := syncCall(t, a, "GET", token, nil); w.Code != 404 {
		t.Fatalf("initial read %d", w.Code)
	}
	in := storedPrefs{Subscriptions: []syncSubscription{{ID: "feed:https://www.newscientist.com/feed/", Kind: "feed", Name: "New Scientist", URL: "https://www.newscientist.com/feed/"}}, ExcludedDomains: []string{"example.com"}}
	w := syncCall(t, a, "PUT", token, in)
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	b, e := NewSyncStore(dir)
	if e != nil {
		t.Fatal(e)
	}
	w = syncCall(t, b, "GET", token, nil)
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	var got userPrefs
	if e = json.Unmarshal(w.Body.Bytes(), &got); e != nil {
		t.Fatal(e)
	}
	if got.User != "richard" || got.Revision != 1 || len(got.Subscriptions) != 1 || len(got.ExcludedDomains) != 1 || got.ExcludedDomains[0] != "example.com" {
		t.Fatalf("round trip %#v", got)
	}
	if w = syncCall(t, a, "PUT", token, in); w.Code != 409 {
		t.Fatalf("stale write %d", w.Code)
	}
	got.ExcludedDomains = []string{}
	w = syncCall(t, a, "PUT", token, got.storedPrefs)
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	w = syncCall(t, b, "GET", token, nil)
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if got.Revision != 2 || len(got.ExcludedDomains) != 0 {
		t.Fatalf("updated %#v", got)
	}
}
func TestSyncAuthorizationAndRevocation(t *testing.T) {
	s, _ := NewSyncStore(t.TempDir())
	token, _ := s.IssueToken("richard", "")
	if _, e := s.IssueToken("richard", ""); e == nil {
		t.Fatal("duplicate identity")
	}
	if w := syncCall(t, s, "GET", strings.Repeat("0", 64), nil); w.Code != 401 {
		t.Fatal(w.Code)
	}
	if e := s.Revoke("richard"); e != nil {
		t.Fatal(e)
	}
	if w := syncCall(t, s, "GET", token, nil); w.Code != 401 {
		t.Fatal(w.Code)
	}
}
func TestSyncRejectsBadProfilesAndOrigin(t *testing.T) {
	s, _ := NewSyncStore(t.TempDir())
	token, _ := s.IssueToken("first", "")
	if w := syncCall(t, s, "PUT", token, storedPrefs{ExcludedDomains: []string{"Example.COM"}}); w.Code != 400 {
		t.Fatal(w.Code)
	}
	r := httptest.NewRequest("PUT", "https://newsgo.themauris.org/api/sync", strings.NewReader(`{"revision":0}`))
	r.Header.Set("Authorization", "Bearer "+token)
	r.Header.Set("Origin", "https://attacker.example")
	w := httptest.NewRecorder()
	s.ServeSync(w, r)
	if w.Code != 403 {
		t.Fatal(w.Code)
	}
}
func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	b, e := os.ReadFile(path)
	if e != nil {
		t.Fatal(e)
	}
	return b
}
