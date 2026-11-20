package dashboard

import (
	"crypto/sha256"
	"fmt"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"github.com/richard-mauri/newsgo/internal/feeds"
	suiteui "github.com/richard-mauri/homepage/shared/ui"
	"github.com/richard-mauri/homepage/shared/ui/banner"
)

// AssertVersionedScript checks that the actual header and cache-busting script
// reference advertise the same runtime version. No historical version encoded.
func assertVersionedScript(t *testing.T, html string) {
	t.Helper()
	header := regexp.MustCompile(`class="suite-type-meta">v([0-9]+\.[0-9]+\.[0-9]+)`)
	match := header.FindStringSubmatch(html)
	if len(match) != 2 {
		t.Fatal("dashboard header missing semantic version")
	}
	if !strings.Contains(html, `/personal.js?v=`+match[1]+`-`) {
		t.Fatalf("script version does not match header version %q", match[1])
	}
}

func renderedVersionPage(t *testing.T) string {
 t.Helper()
 s,err:=New([]feeds.Source{{Name:"Test",URL:"https://example.org/rss"}},"http://127.0.0.1:8083",nil)
 if err!=nil {t.Fatal(err)}
 w:=httptest.NewRecorder()
 s.ServeHTTP(w,httptest.NewRequest("GET","/",nil))
 if w.Code!=200 {t.Fatalf("dashboard HTTP %d",w.Code)}
 return w.Body.String()
}

func TestDashboardVersionConsistent(t *testing.T) { assertVersionedScript(t, renderedVersionPage(t)) }

func TestPersonalScriptURLTracksContent(t *testing.T) {
	s, err := New([]feeds.Source{{Name:"Test",URL:"https://example.org/rss"}}, "https://scrapergo.themauris.org", nil)
	if err != nil { t.Fatal(err) }
	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	digest := sha256.Sum256([]byte(personalJS))
	expected := "/personal.js?v=" + Version() + "-" + fmt.Sprintf("%x", digest[:6])
	if !strings.Contains(w.Body.String(), expected) {
		t.Fatalf("rendered script URL must contain its content fingerprint %q", expected)
	}
}

func TestBannerStylesheetURLTracksEmbeddedCSS(t *testing.T) {
 html := renderedVersionPage(t)
 contents := banner.CSS(suiteui.DesignCSS)
 expected := "/suite-ui/banner.css?v=" + Version() + "-" + bannerCSSRevision(contents)
 if !strings.Contains(html, expected) {t.Fatalf("banner stylesheet URL must fingerprint its CSS, expected %q", expected)}
 if bannerCSSRevision([]byte("original")) == bannerCSSRevision([]byte("revised")) {t.Fatal("changing banner CSS must invalidate Safari's cached URL")}
}

func TestBannerGeometryCannotDriftInDashboardStyles(t *testing.T) {
 html:=renderedVersionPage(t)
 if !strings.Contains(string(banner.CSS(suiteui.DesignCSS)),"margin:24px auto 0") {t.Fatal("shared banner must own viewport-top spacing")}
 if strings.Contains(html,".suite-banner{margin-top:") {t.Fatal("NewsGo must not override the common banner offset")}
}
