package dashboard

import (
 "strings"
 "testing"
)

func TestBackToHomepageNativeLink(t *testing.T) {
 markup := page.Tree.Root.String()
 if !strings.Contains(markup, `class="suite-home-link" href="https://home.themauris.org/"`) {
  t.Fatal("NewsGo must offer a same-tab Homepage return in its header")
 }
 if strings.Contains(markup, `class="suite-home-link" target="_blank"`) {
  t.Fatal("Homepage return must not open a new tab")
 }
 if !strings.Contains(markup,".suite-home-row{flex:0 0 100%;width:100%}") {t.Fatal("Homepage link must be in upper-left header row")}
 if !strings.Contains(markup,".suite-home-link:hover{background:rgba(96,151,172,.22);border-color:var(--accent,#69c6d3);text-decoration:none}") {t.Fatal("standardized hover border style missing")}

}
