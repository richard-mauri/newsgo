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
 } {
  if !strings.Contains(personalJS,s) { t.Errorf("missing subscription failure behavior: %q",s) }
 }
}
