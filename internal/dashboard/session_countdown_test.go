package dashboard

import (
 "strings"
 "testing"
)

// Homepage alone exposes session renewal and countdown. NewsGo consumes
// the parent-bound handoff, verifies account identity, and remains local-only
// without one.
func TestHomepageOwnsSessionManagement(t *testing.T) {
 for _,snippet:=range []string{
  "const suiteHandoffKey='newsgo.homepage.handoff.v1'",
  "const browserOnly=fragment.get('suite_mode')==='local'",
  "sessionStorage.removeItem(suiteHandoffKey)",
  "const homepageAccountHeaders=()=>suiteHandoff?{Authorization:'Bearer '+suiteHandoff}:{}",
  "async function checkHomepageSession()",
  "'/api/session/profile'",
  "suiteExpiresAt=Date.parse(profile.expiresAt)",
  "Signed in as '+homepageIdentity.user+' (Homepage)",
  "Homepage session expired",
  "deadlineCheckPending=true",
  "checkHomepageSession().finally",
  "promptSuiteSignIn()",
  "temporary authenticated guest session valid for one hour",
  "homepageIdentity.role==='guest'",
  "if(homepageIdentity.role!=='guest')void initializeHomepageAccount",
  "explicitBrowserOnly=true",
  "if(browserOnly)",
  "if(!suiteHandoff && !explicitBrowserOnly)",
 } {
  if !strings.Contains(personalJS,snippet){t.Errorf("missing Homepage handoff/expiry behavior: %s",snippet)}
 }
 for _,obsolete:=range []string{
  "location.replace(homepageOrigin+'/launch/newsgo')",
  "renewSuiteSession", "suiteRenewLink", "suiteRenewStatus",
  "sessionRenewFeedback", "sessionTimeRemaining",
  "showSuiteCountdown", "suite-renew-button",
  "/api/handoff/extend-session", "/api/handoff/renew", "/api/session/renew",
  "setInterval(showSuiteCountdown,1000)",
  "Current session: '+sessionTimeRemaining",
 } {
  if strings.Contains(personalJS,obsolete){t.Errorf("NewsGo still owns session renewal: %s",obsolete)}
 }
}

func TestNoNewsGoRenewButtonOrCountdownStyles(t *testing.T) {
 markup:=page.Tree.Root.String()
 for _,obsolete:=range []string{
  ".suite-renew-button", "#suite-renew-status", "Renew current session",
 } {
  if strings.Contains(markup,obsolete){t.Errorf("obsolete session renewal UI in markup: %s",obsolete)}
 }
}
