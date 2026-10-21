package dashboard

import (
	"strings"
	"testing"
	"time"
)

func TestScraperThumbnailURL(t *testing.T) {
	got := scraperThumbnailURL("https://scrapergo.themauris.org/", "https://www.cnbc.com/article.html#new_tab")
	if !strings.HasPrefix(got, "https://scrapergo.themauris.org/api/thumbnail?") || !strings.Contains(got, "width=256") || !strings.Contains(got, "url=https%3A%2F%2Fwww.cnbc.com%2Farticle.html") || strings.Contains(got, "new_tab") {
		t.Fatalf("bad thumbnail URL: %s", got)
	}
	for _, bad := range []string{"javascript:alert(1)", "https://user:password@example.org/x", "file:///etc/passwd"} {
		if v := scraperThumbnailURL("https://scrapergo.themauris.org", bad); v != "" {
			t.Fatalf("unsafe URL accepted: %s", v)
		}
	}
}
func TestFeedCacheNoLongerDependsOnImages(t *testing.T) {
	if v := personalCacheTTL(cacheEntry{}); v != 10*time.Minute {
		t.Fatalf("TTL = %s", v)
	}
}
