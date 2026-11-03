# NewsGo current work

Date: 2026-11-03. Baseline: GitHub `main` after merged header accent PR. Issue [#46](https://github.com/richard-mauri/newsgo/issues/46), branch `feat/46-hourly-account-handoff` (draft PR in progress).

Implemented candidate: accept Homepage one-hour URL fragment handoff, strip fragment from URL, retain only in browser tab sessionStorage, use Bearer header to validate Homepage identity and read/write revisioned account settings. Existing NewsGo cookie proxy remains a fallback. Browser-local subscriptions, exclusions and preferences unchanged. No automatic import, migration, or data deletion. Must deploy Homepage issue #6 first for full operation. Admin server endpoints remain on legacy authorization; this PR does not change them. No merge/restart authorized.

CI and Safari verification pending, including expired token and signed-out fallback.

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
