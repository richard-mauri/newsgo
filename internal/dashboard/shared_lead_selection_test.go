package dashboard
import("strings";"testing")
func TestSubscriptionOnlyLeads(t *testing.T){
 for _,x:=range []string{"function renderPersonalArticles()", "allPersonalArticles", "personal-lead-stories", "excludedURL("}{if !strings.Contains(personalJS,x){t.Errorf("missing article behavior %q",x)}}
 if strings.Contains(personalJS,"function arrangeSharedLeads("){t.Fatal("retired built-in lead pipeline restored")}
}
