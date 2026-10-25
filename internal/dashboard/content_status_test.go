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


func TestClassificationAvailabilityRequiresConfirmedDenial(t *testing.T) {
 cases:=[]struct{classification,reason,want string}{
  {"unavailable","publisher_access_denied","unavailable"},
  {"unavailable","publisher_unavailable",""},
  {"unavailable","timeout",""},
  {"unknown","extraction_unavailable",""},
  {"article","publisher_access_denied",""},
 }
 for _,tc:=range cases {
  got:=classificationAvailability(classificationOutcome{Classification:tc.classification,Reason:tc.reason})
  if got!=tc.want {t.Fatalf("%s/%s: got %q want %q",tc.classification,tc.reason,got,tc.want)}
 }
}
