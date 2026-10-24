package dashboard

import (
 "context"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
)

func TestClassificationDropsNavigationFragment(t *testing.T) {
 article := "https://www.seahorsemagazine.com/article/november-2026/control-from-the-cockpit#new_tab"
 canonical := "https://www.seahorsemagazine.com/article/november-2026/control-from-the-cockpit"
 if normalized := normalizeClassificationURL(article);normalized != canonical {t.Fatalf("normalization: %q",normalized)}
 if !safeURL(normalizeClassificationURL(article)) {t.Fatal("canonical article URL rejected")}
 if safeURL(article) {t.Fatal("fixture must test fragment rejection")}
 upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path != "/api/classify" {t.Errorf("unexpected path %s",r.URL.Path)}
  var body struct{ URL string `json:"url"`}
  if err:=json.NewDecoder(r.Body).Decode(&body);err!=nil {t.Error(err)}
  if body.URL!=canonical {t.Errorf("sent %q, expected canonical URL",body.URL)}
  w.Header().Set("Content-Type","application/json")
  _,_=w.Write([]byte(`{"classification":"article"}`))
 }))
 defer upstream.Close()
 result:=classifyDiagnostic(context.Background(),upstream.URL,article)
 if result.Classification!="article" ||result.Failure!="" {t.Fatalf("diagnostic result: %+v",result)}
 if got:=newsClassify(context.Background(),upstream.URL,article);got!="article" {t.Errorf("newsClassify=%s",got)}
}

func TestClassificationURLKeepsQueryAndRejectsInvalid(t *testing.T){
 raw:="https://example.org/story?utm_source=feed#new_tab"
 if got:=normalizeClassificationURL(raw);got!="https://example.org/story?utm_source=feed" {t.Fatalf("query lost: %s",got)}
 for _,raw:=range []string{"javascript:alert(1)#new_tab","http://user:secret@example.org/article#new_tab","https://example.local/article#new_tab"} {
  if safeURL(normalizeClassificationURL(raw)) {t.Errorf("unsafe URL accepted: %q",raw)}
 }
 if strings.Contains(normalizeClassificationURL("https://example.org/a#new_tab"),"#") {t.Fatal("fragment retained")}
}
