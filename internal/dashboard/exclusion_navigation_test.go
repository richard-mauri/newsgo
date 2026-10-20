package dashboard

import (
	"strings"
	"testing"
)

// Guards the shared and personal subscription navigation contracts.
func TestExcludedPublishersHaveNoReadActionOrEmptyReadingView(t *testing.T) {
	for _, required := range []string{
		"function selectedPublisherExcluded()",
		"if(selectedPublisherExcluded()){el('read-source').value='';selectedSource='';showView('news')}",
		"if(selectedPublisherExcluded()){el('read-source').value='';selectedSource=''}",
		"if(!domain||!excludedURL(domain))actions.append(button('Read Articles',()=>chooseSource(f.url)))",
		"if(source&&excludedURL(publisherForFeed(source))){chooseAllSubscriptions();return}",
	} {
		if !strings.Contains(personalJS, required) {
			t.Errorf("missing exclusion navigation guard: %s", required)
		}
	}
}
