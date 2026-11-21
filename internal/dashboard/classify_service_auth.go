package dashboard

import (
 "errors"
 "net/http"
 "net/url"
 "os"
 "regexp"
 "strings"
 "time"
)

// Separate, server-only classification credential. Never forward a Homepage
// account bearer, and never send this key to a non-TLS remote endpoint.
var serviceClassifyTokenPattern=regexp.MustCompile(`^[a-f0-9]{64}$`)

func authorizeScraperGoClassification(req *http.Request) error {
 token:=os.Getenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN")
 if token=="" {return nil} // Legacy local mock clients have no service identity; guarded ScraperGo rejects them.
 if !serviceClassifyTokenPattern.MatchString(token) {return errors.New("invalid classification service identity")}
 u:=req.URL
 host:=u.Hostname()
 loopback:=(host=="127.0.0.1"||host=="::1"||strings.EqualFold(host,"localhost"))
 if u.User!=nil || (u.Scheme!="https" && !(u.Scheme=="http"&&loopback)) {
  return errors.New("classification service identity requires HTTPS or loopback")
 }
 if req.Method!=http.MethodPost || u.Path!="/api/classify" {return errors.New("classification credential scope violation")}
 req.Header.Set("Authorization","Service "+token)
 return nil
}

// A redirect must never forward a server credential, even within the same host.
func classificationHTTPClient(timeoutSeconds int) *http.Client {
 return &http.Client{Timeout:time.Duration(timeoutSeconds)*time.Second,
  CheckRedirect:func(*http.Request,[]*http.Request)error{return http.ErrUseLastResponse}}
}
