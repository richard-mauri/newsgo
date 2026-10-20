package dashboard
import("strings";"testing")
func TestUnifiedSubscriptionReadingSelection(t *testing.T){
 for _,x:=range []string{"function readingSelection(source)","showPersonal:true", "renderPersonalArticles();", "feeds.find(f=>f.url===url)","registry.length>0"}{if !strings.Contains(personalJS,x){t.Errorf("missing unified reading contract %q",x)}}
 for _,legacy:=range []string{"arrangeSharedLeads(selection)","visibleShared("}{if strings.Contains(personalJS,legacy){t.Errorf("legacy shared headline path remains: %q",legacy)}}
}
