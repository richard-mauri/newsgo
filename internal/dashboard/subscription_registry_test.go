package dashboard

import (
	"strings"
	"testing"
)

// The browser registry is authoritative; the previous Safari keys are read-only
// migration inputs and must remain available for recovery.
func TestSingleSubscriptionRegistryMigration(t *testing.T) {
	for _, token := range []string{
		"newsgo.subscriptions.v2",
		"loadSubscriptions()",
		"localStorage.setItem(key,JSON.stringify(migrated))",
		"function saveSubscriptions(next)",
		"function subscribeRecord(record)",
		"function unsubscribeRecord(id)",
		"feedIdentity(f.url)",
		"sharedRecord(name)",
		"snapshot.slice(offset,offset+8)",
	} {
		if !strings.Contains(personalJS, token) {
			t.Errorf("registry feature missing: %s", token)
		}
	}
	for _, forbidden := range []string{
		"localStorage.removeItem(legacyPersonalKey)",
		"localStorage.removeItem(legacyIncludedKey)",
		"slice(0,8)",
		"Maximum 8 personal feeds",
		"saveIncluded()",
	} {
		if strings.Contains(personalJS, forbidden) {
			t.Errorf("old destructive subscription behavior: %s", forbidden)
		}
	}
}
