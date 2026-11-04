package dashboard

import (
 "strings"
 "testing"
)

func TestSyncConnectInteractionFeedback(t *testing.T) {
 script := personalJS
 for _, expected := range []string{
  "if(syncConnectPending)return;",
  "connectButton.disabled=true;",
  "connectButton.textContent='Connecting…';",
  "connectButton.setAttribute('aria-busy','true');",
  "syncStatus('Checking approved token…');",
  "syncConnectPending=false;",
  "connectButton.disabled=false;",
  "connectButton.removeAttribute('aria-busy');",
  "token was rejected or has expired",
  "Browser preferences remain unchanged.",
 } {
  if !strings.Contains(script,expected) {t.Errorf("missing sign-in state: %s",expected)}
 }
 if !strings.Contains(script,"finally{") && !strings.Contains(script,"finally{\n") {t.Fatal("connection must restore controls in finally")}
}

func TestSyncConnectVisualFeedback(t *testing.T) {
 markup := page.Tree.Root.String()
 for _, expected := range []string{
  "#sync-connect:hover:not(:disabled)",
  "#sync-connect:focus-visible",
  "#sync-connect:active:not(:disabled)",
  "#sync-connect:disabled",
  "id=\"sync-status\" role=\"status\" aria-live=\"polite\"",
 } {
  if !strings.Contains(markup,expected) {t.Errorf("missing sign-in presentation: %s",expected)}
 }
}
