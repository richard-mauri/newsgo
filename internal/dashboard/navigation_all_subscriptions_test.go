package dashboard

import (
	"strings"
	"testing"
)

func TestAllSubscriptionsNavigation(t *testing.T) {
	for _, fragment := range []string{
		"function chooseAllSubscriptions()", "el('read-source').value=''", "renderReading();showView('news')", "Read All Articles", "el('all-subscriptions').addEventListener('click',chooseAllSubscriptions)",
	} {
		if !strings.Contains(personalJS, fragment) {
			t.Errorf("missing all-subscriptions navigation: %s", fragment)
		}
	}
	for _, fragment := range []string{"id=\"all-subscriptions\"", "← All Subscriptions"} {
		if !strings.Contains(page.Tree.Root.String(), fragment) {
			t.Errorf("missing navigation control: %s", fragment)
		}
	}
}
