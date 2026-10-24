package dashboard

import (
 "strings"
 "testing"
)

// An up-to-date Server Sync connection must refresh the subscription navigation.
// Previously its reading view was restored but My Subscriptions could remain (0).
func TestConnectedSyncRefreshesSubscriptionNavigation(t *testing.T) {
 if !strings.Contains(personalJS,"if(syncConnected)renderManager();") {
  t.Fatal("successful Server Sync connection must rerender subscription manager and counts")
 }
}
