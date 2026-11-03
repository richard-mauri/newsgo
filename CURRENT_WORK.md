# NewsGo current work

Date: 2026-11-03. Baseline: GitHub `main` after merged header accent PR. Issue [#46](https://github.com/richard-mauri/newsgo/issues/46), branch `feat/46-hourly-account-handoff` (draft PR in progress).

Implemented candidate: accept Homepage one-hour URL fragment handoff, strip fragment from URL, retain only in browser tab sessionStorage, use Bearer header to validate Homepage identity and read/write revisioned account settings. The NewsGo cookie proxy and route are removed; no fallback. Without a valid one-hour token, account settings are unavailable and all subscriptions/preferences remain browser-local. Browser-local subscriptions, exclusions and preferences unchanged. No automatic import, migration, or data deletion. Must deploy Homepage issue #6 first for full operation. Admin server endpoints remain on legacy authorization; this PR does not change them. No merge/restart authorized.

CI and Safari verification pending, including expired token and browser-only signed-out behavior.

Mac test:
```
cd ~/projects/newsgo
git fetch origin
git checkout feat/46-hourly-account-handoff
git pull --ff-only origin feat/46-hourly-account-handoff
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

Latest candidate: NewsGo header now presents signed-in user, server-based remaining time (mm:ss), and Renew. Renew makes a direct bearer-authenticated POST to Homepage `/api/handoff/renew`; it replaces the short-lived per-tab token. After expiry the button becomes Open Homepage. Countdown does not automatically renew, and browser-local settings remain intact. Verify iOS Safari and renewal against short-lifetime Homepage candidate.

## Expired-token Homepage return
Owner confirmed Homepage remains signed in after five-minute NewsGo token expiration when navigating by Safari address bar, but the JavaScript-driven Open Homepage button reportedly landed on a signed-out Homepage. The expired-state action is now a native same-tab `<a href="https://home.themauris.org/">` link, not `location.assign` from a button. It does not clear Homepage credentials or touch NewsGo browser data. This change is in the existing handoff PR, with Safari verification pending; it must not be presented as a confirmed browser fix before owner testing.
