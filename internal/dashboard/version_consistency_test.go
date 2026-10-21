package dashboard

import (
	"regexp"
	"strings"
	"testing"
)

// AssertVersionedScript checks that the actual header and cache-busting script
// reference advertise the same runtime version. No historical version encoded.
func assertVersionedScript(t *testing.T, html string) {
	t.Helper()
	header := regexp.MustCompile(`NEWS READING WITHOUT PAGE CLUTTER USING SCRAPERGO · V([0-9]+\.[0-9]+\.[0-9]+)`)
	match := header.FindStringSubmatch(html)
	if len(match) != 2 {
		t.Fatal("dashboard header missing semantic version")
	}
	if !strings.Contains(html, `/personal.js?v=`+match[1]+`"`) {
		t.Fatalf("script version does not match header version %q", match[1])
	}
}

func TestDashboardVersionConsistent(t *testing.T) { assertVersionedScript(t, page.Tree.Root.String()) }
