package dashboard

import (
 "bytes"
 "context"
 "encoding/json"
 "errors"
 "io"
 "net/http"
 "strings"
 "time"
)

// The diagnostic is intentionally uncached so it can explain a current failure.
// It does not change the classification gate, article cache, or subscriptions.
type adminClassificationDiagnostic struct {
 Classification string `json:"classification"`
 Failure string `json:"failure,omitempty"`
 HTTPStatus int `json:"httpStatus,omitempty"`
 ElapsedMs int64 `json:"elapsedMs"`
 TimeoutSeconds int `json:"timeoutSeconds"`
}

func classifyDiagnostic(ctx context.Context, reader, article string) adminClassificationDiagnostic {
 article = normalizeClassificationURL(article)
 timeout := currentArticleSettings().ClassificationTimeoutSeconds
 result := adminClassificationDiagnostic{Classification:"unknown", TimeoutSeconds:timeout}
 started:=time.Now()
 if !safeURL(article) {result.Failure="invalid_url";return result}
 contextWithTimeout,cancel:=context.WithTimeout(ctx,time.Duration(timeout)*time.Second)
 defer cancel()
 body,_:=json.Marshal(map[string]string{"url":article})
 request,err:=http.NewRequestWithContext(contextWithTimeout,http.MethodPost,strings.TrimRight(reader,"/")+"/api/classify",bytes.NewReader(body))
 if err!=nil {result.Failure="request_creation";return result}
 request.Header.Set("Content-Type","application/json")
 response,err:=(&http.Client{Timeout:time.Duration(timeout)*time.Second}).Do(request)
 if err!=nil {
  if errors.Is(err,context.DeadlineExceeded)||errors.Is(contextWithTimeout.Err(),context.DeadlineExceeded) {result.Failure="timeout"} else if ctx.Err()!=nil {result.Failure="canceled"} else {result.Failure="transport_error"}
  result.ElapsedMs=time.Since(started).Milliseconds();return result
 }
 defer response.Body.Close()
 result.HTTPStatus=response.StatusCode
 if response.StatusCode!=http.StatusOK {result.Failure="http_error";result.ElapsedMs=time.Since(started).Milliseconds();return result}
 var parsed classificationOutcome
 err=json.NewDecoder(io.LimitReader(response.Body,2048)).Decode(&parsed)
 if err!=nil {result.Failure="invalid_json";result.ElapsedMs=time.Since(started).Milliseconds();return result}
 switch parsed.Classification {
 case "article","homepage","unreadable","unavailable":
  result.Classification=parsed.Classification
 default:
  result.Failure="unexpected_classification"
 }
 result.ElapsedMs=time.Since(started).Milliseconds()
 return result
}

func(s *State) serveAdminClassificationDiagnostic(w http.ResponseWriter,r *http.Request) {
 w.Header().Set("Cache-Control","no-store")
 w.Header().Set("Content-Type","application/json; charset=utf-8")
 w.Header().Set("X-Content-Type-Options","nosniff")
 if r.Method!=http.MethodPost {w.Header().Set("Allow","POST");http.Error(w,"POST required",http.StatusMethodNotAllowed);return}
 if s.syncStore==nil {http.Error(w,"sync unavailable",http.StatusServiceUnavailable);return}
 s.syncStore.mu.Lock()
 token,status:=s.syncStore.authenticate(r)
 s.syncStore.mu.Unlock()
 if status!=http.StatusOK {http.Error(w,"approved token required",status);return}
 if token.Role!="admin" {http.Error(w,"administrator required",http.StatusForbidden);return}
 r.Body=http.MaxBytesReader(w,r.Body,4096)
 var req struct {URL string `json:"url"`}
 if json.NewDecoder(r.Body).Decode(&req)!=nil || len(req.URL)>2048 || !safeURL(normalizeClassificationURL(req.URL)) {http.Error(w,"invalid article URL",http.StatusBadRequest);return}
 _=json.NewEncoder(w).Encode(classifyDiagnostic(r.Context(),s.reader,req.URL))
}
