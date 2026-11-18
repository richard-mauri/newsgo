package dashboard

import (
	"strings"
	"testing"
)

func TestMainNavigationHoverAndKeyboardFocusHighlight(t *testing.T) {
	html := page.Tree.Root.String()
	// Navigation buttons share the same resting appearance; only hover
	// and keyboard focus add a visible highlight.
	for _, css := range []string{
		".news-nav button:focus-visible{outline:2px solid var(--accent);outline-offset:2px}",
	} {
		if !strings.Contains(html, css) {
			t.Errorf("missing main navigation interaction style %q", css)
		}
	}
	if !strings.Contains(suiteButtonCSS,"button:not(:disabled):hover") || !strings.Contains(suiteButtonCSS,"color: var(--accent)") {
		t.Error("Homepage canonical button hover must change label color")
	}
	if strings.Contains(html, `.news-nav button[aria-pressed="true"]{`) {
		t.Fatal("selected navigation state must not add persistent styling")
	}
	if !strings.Contains(html, `<nav class="news-nav" aria-label="NewsGo sections">`) {
		t.Fatal("main navigation missing")
	}
}
