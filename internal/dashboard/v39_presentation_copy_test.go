package dashboard

import ("strings";"testing")

func TestV40PresentationCopy(t *testing.T){
 html:=page.Tree.Root.String()
 if !strings.Contains(html,"NEWS READING WITHOUT PAGE CLUTTER USING SCRAPERGO · V"){t.Fatal("header tagline/version missing")}
 if strings.Contains(html,"NEWS DISCOVERY FOR SCRAPERGO"){t.Fatal("outdated tagline")}
 if strings.Contains(html,`id="sync-token"`)||strings.Contains(html,`id="legacy-sync-details"`){t.Fatal("legacy token interface still present")}
 if strings.Contains(html,`id="manage-sync"`){t.Fatal("obsolete Account Sync navigation remains")}
}
