package dashboard

import (
	"strings"
	"testing"
)

// Homepage owns login lifetime. NewsGo must display the server-provided
// expiration and direct users back to Homepage to obtain a new login.
func TestHomepageSessionCountdownAndRenewal(t *testing.T) {
	for _, snippet := range []string{
		"sessionTimeRemaining(seconds)",
		"suiteExpiresAt-Date.now()",
		"Current session: '+sessionTimeRemaining(remaining)+' remaining",
		"Current session expired",
		"setInterval(showSuiteCountdown,1000)",
		"suiteRenewLink.addEventListener('click',renewSuiteSession)",
		"async function renewSuiteSession()",
		"'/api/handoff/extend-session'",
		"sessionStorage.setItem(suiteHandoffKey,data.token)",
		"verifiedToken!==suiteHandoff",
		"suiteAccountLine.setAttribute('aria-live','off')",
		"if(!Number.isFinite(suiteExpiresAt)||suiteExpiresAt<=Date.now())",
	} {
		if !strings.Contains(personalJS, snippet) {
			t.Errorf("missing session display/renewal behavior: %s", snippet)
		}
	}
	for _, forbidden := range []string{"/api/handoff/renew'", "/api/session/renew'"} {
		if strings.Contains(personalJS, forbidden) {
			t.Errorf("NewsGo must not attempt to extend the fixed Homepage login: %s", forbidden)
		}
	}
}
