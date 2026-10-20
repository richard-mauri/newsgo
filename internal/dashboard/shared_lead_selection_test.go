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
		"const matches=a=>selection.visibleShared(a.dataset.source);",
		"const eligible=a=>matches(a)&&!excludedURL(a.querySelector('a.original')?.href||'');",
		"originals.forEach(a=>{a.hidden=!eligible(a)});",
		"rows.forEach(a=>{a.hidden=!eligible(a)});",
		"let spaces=Math.max(0,3-originals.filter(eligible).length);",
		"if(!spaces||!eligible(row))continue;",
		"lead.append(card);row.hidden=true;spaces--;",
		"lead.querySelectorAll('.subscription-promoted-lead').forEach(a=>a.remove());",
		"arrangeSharedLeads(selection);",
	} {
		if !strings.Contains(personalJS, s) {
			t.Errorf("missing shared lead behavior: %s", s)
		}
	}
}
