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
 if !strings.Contains(markup,`<div class="brand">News<span>Go</span><div><a class="suite-home-link"`) {t.Fatal("return link must follow NewsGo heading")}
 if !strings.Contains(markup,".suite-home-link:hover{background:rgba(96,151,172,.16);color:var(--accent,#69c6d3);text-decoration:underline;opacity:1}") {t.Fatal("standardized hover border style missing")}

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
