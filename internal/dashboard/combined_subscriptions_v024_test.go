package dashboard

import (
    "strings"
    "testing"
)

// Verify the consolidated selection model retains built-in and RSS subscriptions
// without reverting to the v0.2.23 shared-headline visibility bug.
func TestCombinedSubscriptionsVisibility(t *testing.T) {
    for _, expected := range []string{
        "function readingSelection(source)",
        "const all=source==='';",
        "all&&included.includes(name)",
        "showShared:!!sharedName||(all&&included.length>0)",
        "showPersonal:!!feedURL||(all&&feeds.length>0)",
        "const matches=a=>selection.visibleShared(a.dataset.source);",
        "arrangeSharedLeads(selection);",
        "el('subscriptions-empty-notice').hidden=!selection.all||registry.length>0",
    } {
        if !strings.Contains(personalJS, expected) {
            t.Fatalf("missing combined-view behavior: %s", expected)
        }
    }
}
