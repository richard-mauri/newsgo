package dashboard

import ("strings";"testing")

func TestV40PresentationCopy(t *testing.T){
 html:=renderedVersionPage(t)
 if !strings.Contains(html,"class=\"suite-ui suite-banner\"") {t.Fatal("shared banner missing")}
 if !strings.Contains(html,"Your news, without the noise") {t.Fatal("editorial tagline missing")}
 if strings.Contains(html,"NEWS DISCOVERY FOR SCRAPERGO"){t.Fatal("outdated tagline")}
 if strings.Contains(html,`id="sync-token"`)||strings.Contains(html,`id="legacy-sync-details"`){t.Fatal("legacy token interface still present")}
 if strings.Contains(html,`id="manage-sync"`){t.Fatal("obsolete Account Sync navigation remains")}
}
