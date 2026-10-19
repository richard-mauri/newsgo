package dashboard
import("strings";"testing")
func TestSharedSubscriptionPromotesItsOwnLeadStories(t *testing.T){
 for _,s:=range []string{
  "function arrangeSharedLeads(source)",
  "source==='shared:'+a.dataset.source",
  "included.includes(a.dataset.source)",
  "let spaces=Math.max(0,3-originals.filter(matches).length);",
  "if(!spaces||!matches(row))continue;",
  "lead.append(card);row.hidden=true;spaces--;",
  "lead.querySelectorAll('.subscription-promoted-lead').forEach(a=>a.remove());",
  "arrangeSharedLeads(source);",
 }{if !strings.Contains(personalJS,s){t.Errorf("missing shared lead behavior: %s",s)}}
}
