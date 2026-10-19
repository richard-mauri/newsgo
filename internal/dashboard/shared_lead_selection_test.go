package dashboard

import (
	"strings"
	"testing"
)

// The server selects lead stories globally; a subscription view needs to
// promote matching latest-headline rows when its global leads are from elsewhere.
func TestSharedSubscriptionPromotesItsOwnLeadStories(t *testing.T) {
	for _, needle := range []string{
		"function arrangeSharedLeads(source)",
		"const matches=a=>source==='shared:'+a.dataset.source;",
		"let spaces=Math.max(0,3-originals.filter(matches).length);",
		"if(!spaces||!matches(row))continue;",
		"lead.append(card);row.hidden=true;spaces--;",
		"lead.querySelectorAll('.subscription-promoted-lead').forEach(a=>a.remove());",
		"arrangeSharedLeads(source);",
	} {
		if !strings.Contains(personalJS, needle) {
			t.Errorf("missing shared lead behavior: %s", needle)
		}
	}
}
