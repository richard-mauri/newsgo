package dashboard
import("strings";"testing")
func TestSubscriptionFirstStartupDoesNotFlashSharedHeadlines(t *testing.T){
 if !strings.Contains(page.Tree.Root.String(),`id="shared-headlines" class="personal-hidden"`){t.Fatal("shared headlines must be hidden until subscriptions load")}
 for _,s:=range []string{"const showShared=source.startsWith('shared:')||(source===''&&included.length>0);","included.includes(a.dataset.source)"}{if !strings.Contains(personalJS,s){t.Fatalf("shared headlines must appear only for subscribed publishers: %s",s)}}
 if !strings.Contains(page.Tree.Root.String(),`/personal.js?v=0.2.24`){t.Fatal("browser JS URL must be versioned")}
}
