package dashboard

import (
	"strings"
	"testing"
)

func TestV40PresentationCopy(t *testing.T) {
	if !strings.Contains(page.Tree.Root.String(), "NEWS READING WITHOUT PAGE CLUTTER USING SCRAPERGO · V0.2.41") {
		t.Fatal("header must use requested tagline and version")
	}
	if strings.Contains(page.Tree.Root.String(), "NEWS DISCOVERY FOR SCRAPERGO") {
		t.Fatal("outdated header tagline")
	}
	if !strings.Contains(page.Tree.Root.String(), `id="sync-token" autocomplete="off" placeholder="Enter the 64-character approved token distributed by Richard."`) {
        t.Fatal("token input placeholder must contain the approved wording")
    }
    if strings.Contains(page.Tree.Root.String(), `placeholder="Paste your approved token"`) {
        t.Fatal("outdated token placeholder")
    }
    if !strings.Contains(personalJS, "Enter the 64-character approved token distributed by Richard.") {
		t.Fatal("server sync instructions must identify token distributor")
	}
}
