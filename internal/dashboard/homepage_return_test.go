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
}
