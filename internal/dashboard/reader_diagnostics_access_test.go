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
 if !strings.Contains(personalJS,"el('manage-admin').hidden=!mayAdminister()"){t.Fatal("administrator UI role gate missing")}
 if !strings.Contains(personalJS,"function homepageAdminActive()"){t.Fatal("Homepage admin handoff gate missing")}
}
