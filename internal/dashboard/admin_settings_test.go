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

func adminCall(t *testing.T, s *SyncStore, method, token string, payload any) *httptest.ResponseRecorder {
	t.Helper()
	var body []byte
	if payload != nil {
		var err error
		body, err = json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
	}
	req := httptest.NewRequest(method, "https://newsgo.themauris.org/api/admin/settings", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	if method == http.MethodPut {
		req.Header.Set("Origin", "https://newsgo.themauris.org")
	}
	w := httptest.NewRecorder()
	s.ServeAdminSettings(w, req)
	return w
}
func TestAdminSettingsAuthorizationAndPersistence(t *testing.T) {
	dir := t.TempDir()
	s, err := NewSyncStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	adminToken, err := s.IssueToken("richard", "admin candidate")
	if err != nil {
		t.Fatal(err)
	}
	guestToken, err := s.IssueToken("guest", "guest")
	if err != nil {
		t.Fatal(err)
	}
	if w := adminCall(t, s, "GET", adminToken, nil); w.Code != 403 {
		t.Fatalf("pre-promotion status %d", w.Code)
	}
	if err = s.SetRole("richard", "admin"); err != nil {
		t.Fatal(err)
	}
	if w := adminCall(t, s, "GET", guestToken, nil); w.Code != 403 {
		t.Fatalf("guest status %d", w.Code)
	}
	if w := adminCall(t, s, "GET", "invalid", nil); w.Code != 401 {
		t.Fatalf("invalid status %d", w.Code)
	}
	current := defaultServerSettings()
	current.Image.MetadataTimeoutMs = 4500
	w := adminCall(t, s, "PUT", adminToken, current)
	if w.Code != 200 {
		t.Fatalf("put %d %s", w.Code, w.Body.String())
	}
	w = adminCall(t, s, "PUT", adminToken, current)
	if w.Code != 409 {
		t.Fatalf("stale update %d", w.Code)
	}
	w = adminCall(t, s, "GET", adminToken, nil)
	if w.Code != 200 {
		t.Fatalf("get %d", w.Code)
	}
	var got serverSettings
	if err = json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Image.MetadataTimeoutMs != 4500 || got.Revision != 1 {
		t.Fatalf("settings %#v", got)
	}
	if strings.Contains(string(mustRead(t, filepath.Join(dir, "approved-tokens.json"))), adminToken) {
		t.Fatal("plain token persisted")
	}
	if _, err = os.Stat(filepath.Join(dir, "server-settings.json")); err != nil {
		t.Fatal(err)
	}
	if err = s.SetRole("richard", "user"); err != nil {
		t.Fatal(err)
	}
	if w = adminCall(t, s, "GET", adminToken, nil); w.Code != 403 {
		t.Fatalf("demotion status %d", w.Code)
	}
}
func TestAdminSettingsRangeValidation(t *testing.T) {
	s, _ := NewSyncStore(t.TempDir())
	token, _ := s.IssueToken("richard", "")
	_ = s.SetRole("richard", "admin")
	c := defaultServerSettings()
	c.Image.ParallelRequests = 100
	w := adminCall(t, s, "PUT", token, c)
	if w.Code != 400 {
		t.Fatalf("bad settings status %d", w.Code)
	}
}

// Confirm the 30-second classification preference survives a new store
// instance and is activated after a simulated process restart.
func TestClassificationTimeoutPersistsAcrossRestart(t *testing.T) {
 dir:=t.TempDir()
 store,err:=NewSyncStore(dir)
 if err!=nil {t.Fatal(err)}
 token,err:=store.IssueToken("timeout-admin","")
 if err!=nil {t.Fatal(err)}
 if err=store.SetRole("timeout-admin","admin");err!=nil {t.Fatal(err)}
 cfg:=defaultServerSettings()
 cfg.Articles.ClassificationTimeoutSeconds=30
 if w:=adminCall(t,store,http.MethodPut,token,cfg);w.Code!=http.StatusOK {
  t.Fatalf("save status %d: %s",w.Code,w.Body.String())
 }
 restarted,err:=NewSyncStore(dir)
 if err!=nil {t.Fatal(err)}
 if err=restarted.ActivateServerSettings();err!=nil {t.Fatal(err)}
 if got:=currentArticleSettings().ClassificationTimeoutSeconds;got!=30 {
  t.Fatalf("after restart classification timeout = %d, want 30",got)
 }
 if w:=adminCall(t,restarted,http.MethodGet,token,nil);w.Code!=http.StatusOK {
  t.Fatalf("reload status %d",w.Code)
 }
}
