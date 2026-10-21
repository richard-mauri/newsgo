package dashboard

import (
	"regexp"
)

// Version returns the dashboard's displayed semantic version. The dashboard
// header is the authoritative version until its existing template is refactored.
// Keeping the CLI coupled to the displayed header avoids another hardcoded
// release number and protects the existing cache-buster contract.
func Version() string {
	match := regexp.MustCompile(`NEWS READING WITHOUT PAGE CLUTTER USING SCRAPERGO · V([0-9]+\.[0-9]+\.[0-9]+)`).FindStringSubmatch(page.Tree.Root.String())
	if len(match) != 2 {
		panic("NewsGo: dashboard semantic version missing")
	}
	return match[1]
}
