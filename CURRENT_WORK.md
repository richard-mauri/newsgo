# Current work

**Merged baseline:** GitHub `main` at `046783f1c020938a4d431711eb5bb905ea179727`, verified November 7, 2026 following approved PR #66 merge and Issue #65 closure.

**Active issue and candidate:** [Issue #67](https://github.com/richard-mauri/newsgo/issues/67), branch `feat/67-homepage-auth-diagnostics`, [draft PR #68](https://github.com/richard-mauri/newsgo/pull/68), targeting main. Homepage-first authorization and removal of the duplicate-login requirement, while preserving legacy Server Sync data and credentials until migration is accepted. Do not merge without explicit owner authorization.

**Implemented:** Latest-reading Diagnostics is available to everyone (without admin login), with latest-read classification cache hits, misses, bypasses, and hit rate. The NewsGo Go server now independently validates a dedicated `X-NewsGo-Homepage-Authorization` bearer against the fixed HTTPS `/api/session/profile` endpoint on each privileged request (three-second timeout, no redirects, strict token form, server-asserted admin role and nonexpired session); valid Homepage admin can use Admin Settings and per-article fresh classification diagnostics. The browser submits Homepage credentials for these endpoints via a distinct header, so legacy `Authorization` tokens are never forwarded to Homepage. Existing legacy Server Sync admin tokens remain a compatibility fallback. The Account Sync UI identifies Homepage as primary and folds the legacy token flow inside a collapsed compatibility section. No credentials are copied, deleted, rotated or persisted in new places.

**Pending/risks:** Exact-head CI; live Homepage connectivity and authorization in owner Mac/Safari; negative tests for full endpoint coverage (normal user, expiry, network fail, missing credentials), account data migration acceptance. Do not delete the `/api/sync` legacy backend, token records, user profiles or administrative data until the owner verifies Homepage account load/save and role-based administration. Classification cache remains in memory, so restarting NewsGo clears it; durable persistence is a separate decision.

**Mac acceptance commands:**
```sh
cd ~/projects/newsgo
git fetch origin
git checkout feat/67-homepage-auth-diagnostics
git pull --ff-only origin feat/67-homepage-auth-diagnostics
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

Safari: Homepage admin should see Admin Settings without legacy token; Homepage regular account should not. Diagnostics available to both and show per-refresh cache counts. Inspect Homepage account save/load and expired/invalid session handling. Verify legacy admin remains available if needed. No changes to ScraperGo, Cloudflare settings, local build counters or private profiles. PR remains draft and unmerged.
