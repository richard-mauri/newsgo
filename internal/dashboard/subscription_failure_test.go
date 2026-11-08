package dashboard

import (
 "strings"
 "testing"
)
func TestSubscribeFailureDialogContract(t *testing.T) {
 for _, s := range []string{
  "function showSubscriptionFailure(trigger,f)",
  "trigger.dataset.subscribeError=error",
  "trigger.textContent='Failure'",
  "function subscribeClick(f,trigger)",
  "button('Retry'",
  "dialog.setAttribute('aria-modal','true')",
  "function safeSubscriptionFeedURL(raw)",
  "u.username='';u.password='';u.hash='';",
  "u.searchParams.set(key,'[redacted]')",
  "Attempted RSS/Atom/JSON feed URL:",
  "Validation failure:",
  "The address above is the attempted feed endpoint",
  "dialog.append(heading,endpoint,detail,note,actions)",
  "delete trigger.dataset.subscribeError;trigger.disabled=true",

 } {
  if !strings.Contains(personalJS,s) { t.Errorf("missing subscription failure behavior: %q",s) }
 }
}
