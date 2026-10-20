package dashboard

import (
	"strings"
	"testing"
)

func TestPersonalRefreshBusyIndicator(t *testing.T) {
	for _, fragment := range []string{"role=\"status\" aria-live=\"polite\"", "@keyframes newsgo-spin", "prefers-reduced-motion", "is-refreshing"} {
		if !strings.Contains(page.Tree.Root.String(), fragment) {
			t.Errorf("missing accessible busy indicator: %s", fragment)
		}
	}
	for _, fragment := range []string{"loading.classList.add('is-refreshing')", "loading.classList.remove('is-refreshing')", "snapshot.slice(offset,offset+8)", "sequence!==refreshSequence", "Could not refresh feeds:", "completed+=batch.length", "const succeeded=data.filter(group=>!group.error).length"} {
		if !strings.Contains(personalJS, fragment) {
			t.Errorf("missing refresh state handling: %s", fragment)
		}
	}
}
