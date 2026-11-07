# Current work

**Merged baseline:** GitHub `main` at `046783f1c020938a4d431711eb5bb905ea179727` after the owner-approved PR #66 merge and closure of Issue #65, verified November 7, 2026. This baseline is unchanged.

**Active work:** [Issue #67](https://github.com/richard-mauri/newsgo/issues/67), branch `feat/67-homepage-auth-diagnostics`, [draft PR #68](https://github.com/richard-mauri/newsgo/pull/68) targeting main. The owner explicitly authorized removing the legacy Server Sync implementation on November 7, 2026. Do not merge, tag or release without separate owner authorization.

**Implemented:** Homepage is the only live account and administrator identity. NewsGo's Go server verifies active Homepage admin roles for both privileged endpoints using the fixed Homepage profile endpoint, no redirects, a short timeout, and a dedicated auth header. Browser Account Sync supports Homepage's revisioned per-user save/load and local browser preferences. Ordinary readers can access latest-read Diagnostics showing cache hits, misses, bypasses, disabled counts and TTL. Owner Safari screenshots established that a saved TTL of zero disabled caching; after changing it to 1440 minutes, repeated readings demonstrated 329 hits/6 misses (98%) and a 245 ms refresh versus 11.71 s cold.

**Complete cleanup in candidate:** Removed `sync.go`, its token-backed profile logic, former token CLI, legacy-only tests, browser token login, JavaScript legacy synchronization, `/api/sync`, administrator token fallback, old subscription-format migration and shared-record placeholders, and old server-settings schema fallbacks. A standalone `SettingsStore` now owns only operational settings. The owner confirmed that migration is unnecessary and is willing to start over; historical token/profile files, if left on disk, are unused and may be removed locally after stopping the service. Homepage is the sole active account data store.

**Validation:** GitHub CI checks tests and build after each candidate commit. Exact-head CI result must be recorded in draft PR #68 after the final documentation changes. Live Mac/Safari acceptance of the cleaned-up account UI and Homepage-only admin APIs remains outstanding. CI does not prove real Homepage network access, session expiry, Cloudflare routing or persistence of personal account data.

**Owner Mac acceptance commands:**
```sh
cd ~/projects/newsgo
git fetch origin
git checkout feat/67-homepage-auth-diagnostics
git pull --ff-only origin feat/67-homepage-auth-diagnostics
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

**Safari acceptance:** Launch from Homepage. Confirm only Homepage Account Sync save/load appears, with no Server Sync token dialog; account revisioned save/load works with the current subscription and exclusion schema; Homepage admin can view and save Admin Settings and run fresh article diagnostics; ordinary Homepage users cannot; Diagnostics works for everyone and classification cache remains effective. Verify `/api/sync` is no longer served. Current operational server settings must persist across restart. No change to ScraperGo, Cloudflare settings, user private files, builds or semantic version.
