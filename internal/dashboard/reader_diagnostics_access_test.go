package dashboard

import (
 "strings"
 "testing"
)

// Metrics about a reader's own briefing must not require a legacy
// Server Sync credential or the administrator role.
func TestReadingDiagnosticsAvailableWithoutLegacySession(t *testing.T) {
 for _,needle:=range []string{
  "el('manage-diagnostics').hidden=false",
  "el('manage-diagnostics').addEventListener('click',()=>{showView('diagnostics');updateDiagnostics()})",
 }{
  if !strings.Contains(personalJS,needle){t.Errorf("reader diagnostics not accessible: %q",needle)}
 }
 if strings.Contains(personalJS,"el('manage-diagnostics').hidden=!(syncConnected&&isAdmin)"){t.Fatal("diagnostics still gated by legacy credentials")}
 if !strings.Contains(page.Tree.Root.String(),`aria-label="Reading diagnostics"`){t.Fatal("reader diagnostics label missing")}
 // Admin endpoints and UI must remain privileged until server-verified
 // Homepage authorization replaces the independent legacy guard.
 if !strings.Contains(personalJS,"el('admin-panel').hidden=!homepageAdminActive()"){t.Fatal("administrator settings panel role gate missing")}
 if !strings.Contains(personalJS,"function homepageAdminActive()"){t.Fatal("Homepage admin handoff gate missing")}
}

func TestUnifiedSettingsNavigationPreservesRoleSeparation(t *testing.T) {
 for _,needle:=range []string{
  "settingsPanel.append(el('preferences-panel'),el('admin-panel'))",
  "el('manage-settings').addEventListener('click'",
  "el('settings-panel').classList.toggle('personal-hidden',view!=='settings')",
  "if(mayAdminister())loadAdmin()",
  "el('close-settings').addEventListener('click',()=>confirmAdminLeave(()=>showView('news')))",
 } {
  if !strings.Contains(personalJS,needle) {t.Errorf("missing unified settings contract %q",needle)}
 }
}

func TestNavigationSelectionKeepsSharedDarkButtonStyle(t *testing.T) {
 markup := page.Tree.Root.String()
 if strings.Contains(markup, `.news-nav button[aria-pressed="true"]{`) || strings.Contains(markup, `.news-nav button.ui-button[aria-pressed="true"]{`) {
  t.Fatal("selected navigation buttons must not have persistent border or fill highlights")
 }
 if !strings.Contains(personalJS, "document.querySelectorAll('.news-nav button,.panel-close').forEach(b=>b.classList.add('ui-button-large'))") {
  t.Fatal("navigation must use shared large button styling")
 }
}

func TestAllActionButtonsUseSharedDarkAppearance(t *testing.T) {
 markup := page.Tree.Root.String()
 for _,needle := range []string{
  `button.ui-button{appearance:none;background:var(--soft);color:var(--text);border:1px solid var(--line)`,
  `button.ui-button:not(:disabled):hover{color:var(--accent);border-color:var(--accent)`,
  `button.ui-button:not(:disabled):focus-visible{outline:2px solid var(--accent)`,
 } {
  if !strings.Contains(markup,needle) {t.Errorf("shared dark action style missing: %q",needle)}
 }
 for _,obsolete := range []string{
  `button.ui-button.ui-button-primary{`,
  `.admin-settings-actions button[type=submit]{`,
  `.subscription-actions button:first-child{`,
  `.filtered-tab[aria-selected=true]{`,
  `.exclusion-confirm-action{`,
  `background:var(--accent);color:light-dark(#fff,#082028);padding:11px`,
 } {
  if strings.Contains(markup,obsolete) {t.Errorf("obsolete button override remains: %q",obsolete)}
 }
 if strings.Contains(personalJS,"classList.add('ui-button-primary')") {
  t.Fatal("obsolete primary button class assignment remains")
 }
 for _,needle := range []string{
  "document.querySelectorAll('button').forEach(b=>b.classList.add('ui-button'))",
  "b.classList.add('ui-button');if(cls)b.classList.add(cls)",
  "tab.className='filtered-tab ui-button'",
  "el('close-settings').classList.add('ui-button','ui-button-large')",
 } {
  if !strings.Contains(personalJS,needle) {t.Errorf("shared button initialization missing: %q",needle)}
 }
}
