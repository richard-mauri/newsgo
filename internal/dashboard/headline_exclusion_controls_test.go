package dashboard

import (
	"strings"
	"testing"
)

// Domain-exclusion management belongs to My Subscriptions and Excluded Domains,
// not to lead-story and latest-headline cards.
func TestHeadlineCardsDoNotRenderExclusionButtons(t *testing.T) {
	for _, forbidden := range []string{
		"const exclude=excludeButton(a.url)",
		"body.append(button('Exclude articles from '+ad",
	} {
		if strings.Contains(personalJS, forbidden) {
			t.Errorf("headline exclusion button remains: %s", forbidden)
		}
	}
	for _, required := range []string{
		"function excludeButton(domain)",
		"function setUnifiedDomain(domain,allow)",
		"if(!excludedURL(domain)){const exclude=excludeButton(domain);if(exclude)actions.append(exclude)}",
	} {
		if !strings.Contains(personalJS, required) {
			t.Errorf("required subscription exclusion behavior missing: %s", required)
		}
	}
}
