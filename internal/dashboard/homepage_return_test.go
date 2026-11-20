package dashboard

import (
 "net/http/httptest"
 "strings"
 "testing"

 "github.com/richard-mauri/newsgo/internal/feeds"
)

func TestBackToHomepageNativeLink(t *testing.T) {
 s,err:=New([]feeds.Source{{Name:"Example",URL:"https://example.org/feed"}},"http://127.0.0.1:8083",nil)
 if err!=nil {t.Fatal(err)}
 w:=httptest.NewRecorder()
 s.ServeHTTP(w,httptest.NewRequest("GET","/",nil))
 markup:=w.Body.String()
 if !strings.Contains(markup, `href="https://home.themauris.org/"`) {
  t.Fatal("NewsGo must offer a same-tab Homepage return in its header")
 }
 if strings.Contains(markup, `class="suite-banner__home suite-type-control" target="_blank"`) {
  t.Fatal("Homepage return must not open a new tab")
 }
 for _,want:=range []string{`class="suite-type-title">NewsGo`, `suite-banner__portrait-caption">The Mauris`, `class="suite-ambient suite-ambient-one"`, `class="suite-ambient suite-ambient-two"`, `id="suite-auth-status"`, `/suite-ui/banner.css?v=`} {
  if !strings.Contains(markup,want) {t.Errorf("shared banner missing %q",want)}
 }
 if strings.Contains(markup,`class="topbar"`) {t.Fatal("obsolete topbar retained")}
 if !strings.Contains(markup,`</div><header class="suite-ui suite-banner"`) || !strings.Contains(markup,"</header>\n<main>") {t.Fatal("NewsGo banner must mount outside its article-specific main container")}
 if !strings.Contains(markup,".suite-banner{margin-top:24px}") {t.Fatal("missing NewsGo-specific top spacing")}

}

func TestSuiteBannerStylesheetUsesExactHomepageAmbient(t *testing.T) {
 s,err:=New([]feeds.Source{{Name:"Example",URL:"https://example.org/feed"}},"http://127.0.0.1:8083",nil)
 if err!=nil {t.Fatal(err)}
 w:=httptest.NewRecorder()
 s.ServeHTTP(w,httptest.NewRequest("GET","/suite-ui/banner.css",nil))
 if w.Code!=200 {t.Fatalf("stylesheet HTTP %d",w.Code)}
 for _,want:=range []string{".suite-ambient-one", ".suite-ambient-two", "blur(70px)", "top: -300px", "bottom: 0", ".suite-banner"} {
  if !strings.Contains(w.Body.String(),want) {t.Errorf("missing shared stylesheet style %q",want)}
 }
 if strings.Contains(w.Body.String(),"radial-gradient(") {t.Fatal("Homepage ambient must not be approximated by gradients")}
 for _,want:=range []string{"width:min(1180px,calc(100vw - 96px))","width:min(720px,calc(100vw - 48px))","width:calc(100vw - 36px)","margin-inline:auto","border:0"} {if !strings.Contains(w.Body.String(),want) {t.Errorf("shared banner must own its gutters: %s",want)}}
 if strings.Contains(w.Body.String(),"border:1px solid var(--suite-line)") {t.Fatal("shared banner must not draw an outline")}
}

func TestNoRedundantHomepagePromptInNewsGo(t *testing.T) {
 script := personalJS
 if strings.Contains(script,"Not signed in · Open from Homepage") || strings.Contains(script,"suiteHomeLink") || strings.Contains(script,"textContent='Open Homepage'") {
  t.Fatal("NewsGo account status must not create another Homepage link")
 }
 if strings.Contains(script,"suiteRenewButton") || strings.Contains(script,"renewSuiteAccount(") {
  t.Fatal("separate one-hour renewal UI must be removed")
 }
 if !strings.Contains(script,"fragment.get('suite_mode')==='local'") || !strings.Contains(script,"sessionStorage.removeItem(suiteHandoffKey)") {
  t.Fatal("explicit signed-out launch must clear the previous app bearer")
 }
}

func TestHomepageAccountHeaderDefersSessionManagement(t *testing.T) {
 script := personalJS
 if !strings.Contains(script,"'Signed in as '+homepageIdentity.user+' (Homepage)'") {
  t.Fatal("NewsGo must show server-verified Homepage identity")
 }
 if !strings.Contains(script,"suiteExpiresAt>Date.now()") ||
    !strings.Contains(script,"checkHomepageSession().finally") {
  t.Fatal("NewsGo must revalidate expired cached handoffs against Homepage")
 }
 if strings.Contains(script,"sessionTimeRemaining(") ||
    strings.Contains(script,"renewSuiteSession") ||
    strings.Contains(script,"suiteRenewLink") {
  t.Fatal("NewsGo must not provide its own countdown or Renew control")
 }
}
