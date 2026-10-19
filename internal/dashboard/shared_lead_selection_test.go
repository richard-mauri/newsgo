package dashboard

import (
	"strings"
	"testing"
)

func TestSharedSubscriptionPromotesItsOwnLeadStories(t *testing.T) {
	for _, s := range []string{
		"function arrangeSharedLeads(selection)",
		"name===sharedName",
		"included.includes(name)",
		"let spaces=Math.max(0,3-originals.filter(matches).length);",
		"if(!spaces||!matches(row))continue;",
		"lead.append(card);row.hidden=true;spaces--;",
		"lead.querySelectorAll('.subscription-promoted-lead').forEach(a=>a.remove());",
		"arrangeSharedLeads(selection);",
	} {
		if !strings.Contains(personalJS, s) {
			t.Errorf("missing shared lead behavior: %s", s)
		}
	}
}
