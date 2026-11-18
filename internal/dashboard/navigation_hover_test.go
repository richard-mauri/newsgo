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
		"href=\"/suite-ui/buttons.css\"",
	} {
		if !strings.Contains(html, css) {
			t.Errorf("missing main navigation interaction style %q", css)
		}
	}
	if !strings.Contains(suiteButtonCSS,"button:not(:disabled):hover") || !strings.Contains(suiteButtonCSS,"color: var(--accent)") || !strings.Contains(suiteButtonCSS,":focus-visible") {
		t.Error("Homepage canonical button hover must change label color")
	}
	if strings.Contains(html, `.news-nav button[aria-pressed="true"]{`) {
		t.Fatal("selected navigation state must not add persistent styling")
	}
	if !strings.Contains(html, `<nav class="news-nav" aria-label="NewsGo sections">`) {
		t.Fatal("main navigation missing")
	}
}
