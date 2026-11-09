package dashboard

import (
	"strings"
	"testing"
)

func TestMainNavigationHoverAndKeyboardFocusHighlight(t *testing.T) {
	html := page.Tree.Root.String()
	// Match the Homepage Manage Tokens interaction without changing the
	// selected navigation button's contrasting text and background.
	for _, css := range []string{
		".news-nav button:hover:not([aria-pressed=\"true\"]),.news-nav button:focus-visible:not([aria-pressed=\"true\"]){border-color:var(--accent);color:var(--accent)}",
		".news-nav button:focus-visible{outline:2px solid var(--accent);outline-offset:2px}",
		".news-nav button[aria-pressed=\"true\"]{background:var(--accent);color:light-dark(#fff,#082028)}",
	} {
		if !strings.Contains(html, css) {
			t.Errorf("missing main navigation interaction style %q", css)
		}
	}
	if !strings.Contains(html, `<nav class="news-nav" aria-label="NewsGo sections">`) {
		t.Fatal("main navigation missing")
	}
}
