package dashboard

import (
	"strings"
	"testing"
)

func TestNewsGoAutomaticAccountPreferences(t *testing.T) {
	markup := page.Tree.Root.String()
	for _, obsolete := range []string{
		`id="manage-sync"`, `id="sync-panel"`, "Account Sync",
	} {
		if strings.Contains(markup, obsolete) {
			t.Errorf("legacy sync UI remains: %q", obsolete)
		}
	}
	for _, obsolete := range []string{
		"account-save-button", "account-load-button",
		"accountEstablish()", "accountLoad()", "mountAccountActions()",
		"Save to account", "Load from account",
	} {
		if strings.Contains(personalJS, obsolete) {
			t.Errorf("manual account sync code remains: %q", obsolete)
		}
	}
	for _, required := range []string{
		"initializeHomepageAccount(profile.user,verifiedToken)",
		"const record=await homepageReaderRecord()",
		"applyAccountPreferences(record.data)",
		"accountRevision=record.revision",
		"accountRevision=0",
		"if(!record)scheduleAccountSave()",
		"response.status===409",
		"accountOwner!==homepageIdentity.user",
		"newsgo.account.last-owner.v1",
		"homepageIdentity.user!==owner||accountOwner!==owner||suiteHandoff!==token",
		"accountStatus('Account preferences could not load:",
	} {
		if !strings.Contains(personalJS, required) {
			t.Errorf("missing automatic account persistence invariant: %q", required)
		}
	}
}
