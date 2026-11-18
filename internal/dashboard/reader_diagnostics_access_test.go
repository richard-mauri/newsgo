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
 selected := `.news-nav button.ui-button[aria-pressed="true"]{background:var(--soft);color:var(--text);border-color:var(--accent)}`
 if !strings.Contains(markup, selected) {
  t.Fatal("selected navigation buttons must retain shared dark fill with accent border")
 }
 if !strings.Contains(personalJS, "document.querySelectorAll('.news-nav button,.panel-close').forEach(b=>b.classList.add('ui-button-large'))") {
  t.Fatal("navigation must use shared large button styling")
 }
}

func TestAllActionButtonsUseSharedDarkAppearance(t *testing.T) {
 markup := page.Tree.Root.String()
 for _,needle := range []string{
  `button.ui-button.ui-button-primary{background:var(--soft);border-color:var(--accent);color:var(--text)}`,
  `.admin-settings-actions button.ui-button[type=submit]`,
  `.subscription-actions button.ui-button:first-child`,
  `.filter-btn.ui-button`,
  `.filtered-tab.ui-button[aria-selected=true]`,
  `.exclusion-confirm-action.ui-button`,
  `background:var(--soft)!important;color:var(--text)!important;border-color:var(--accent)!important`,
 } {
  if !strings.Contains(markup,needle) {t.Errorf("shared dark action style missing: %q",needle)}
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
