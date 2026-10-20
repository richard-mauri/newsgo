package dashboard

import (
 "strings"
 "testing"
)

func TestSyncReconnectRequiresExplicitDecision(t *testing.T) {
 for _, required := range []string{
  "if(!samePreferences(data))", "resolveConnectDifference()",
  "Use Server Preferences", "Merge Preferences", "Cancel",
  "decision==='cancel'", "decision==='merge'",
  "[...data.excludedDomains,...excludedDomains]",
  "[...convertLegacy(data.subscriptions),...registry]",
  "syncRevision=data.revision", "scheduleServerSave()",
  "if(remember){try{localStorage.setItem(syncCredentialKey,token)}",
 } {
  if !strings.Contains(personalJS, required) { t.Errorf("missing reconnection safeguard %q", required) }
 }
}
