package dashboard

import (
 "context"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "testing"
)

func TestNewsClassifyDetailRetainsPreviewStatus(t *testing.T) {
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path!="/api/classify" {http.NotFound(w,r);return}
  _=json.NewEncoder(w).Encode(classificationOutcome{Classification:"article",ContentStatus:"partial",ContentReason:"publisher_access_restricted"})
 }))
 defer server.Close()
 result:=newsClassifyDetail(context.Background(),server.URL,"https://example.org/story/preview-status")
 if result.Classification!="article" || result.ContentStatus!="partial" {t.Fatalf("unexpected result %+v",result)}
}
func TestNewsClassifyDetailLegacyStatusUnknown(t *testing.T) {
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  _=json.NewEncoder(w).Encode(map[string]any{"classification":"article"})
 }))
 defer server.Close()
 result:=newsClassifyDetail(context.Background(),server.URL,"https://example.org/story/legacy-status")
 if result.Classification!="article" || result.ContentStatus!="unknown" {t.Fatalf("legacy response %+v",result)}
}
