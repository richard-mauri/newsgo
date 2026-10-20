package dashboard

import (
	"strings"
	"testing"
)

func TestExcludedSubscriptionRetainsActionToAllowPublisher(t *testing.T) {
	// Both shared subscriptions and individually followed personal feeds must
	// allow reversal without removing subscription membership.
	for _, marker := range []string{
		"function excludeButton(domain)",
		"excludedURL(d)?'Allow this site':'Exclude this site'",
		"const exclude=excludeButton(domain);if(exclude)actions.append(exclude);row.append(info,actions);root.append(row)",
		"if(!domain||!excludedURL(domain))actions.append(button('Read Articles'",
	} {
		if !strings.Contains(personalJS, marker) {
			t.Errorf("missing publisher exclusion action contract: %s", marker)
		}
	}
}
