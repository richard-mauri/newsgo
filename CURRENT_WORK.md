# Current work — NewsGo

## Verified development state (November 7, 2026)

- **Merged baseline:** GitHub `main` at `046783f1c020938a4d431711eb5bb905ea179727` following owner-approved PR #66.
- **Active scope:** Issue [#67](https://github.com/richard-mauri/newsgo/issues/67), `feat/67-homepage-auth-diagnostics`, draft [PR #68](https://github.com/richard-mauri/newsgo/pull/68) targeting `main`.
- **Merge boundary:** No authorization to merge, tag, release, change `main`, modify ScraperGo/Cloudflare, delete historical user data, or bump the semantic version. Keep PR draft.
- **Repository process:** Read `AGENTS.md`, `DEVELOPMENT.md`, `ARCHITECTURE.md` and this file from the PR branch. Always inspect the actual latest PR head and its exact-head CI when resuming.

## Scope implemented in PR #68

Homepage is the only active NewsGo account/profile and admin authorization provider. Legacy Server Sync login, tokens, API, token CLI and related browser workflow were removed with owner approval. Homepage per-user revisioned save/load and server-side administrator verification remain. Operational settings use `data/server-settings.json`; historical token/user-profile files are not touched or read. Both reading modes require successful ScraperGo `article` classification. Latest-read Diagnostics reports cache hits/misses/bypasses/disabled counts, hit rate and classification timings.

**Persistent 60-minute classification cache has now been implemented on this feature branch**, replacing the earlier process-only cache. `internal/dashboard/classification_cache_disk.go` saves validated classifications as permission-restricted atomic snapshots in `<settings-dir>/classification-cache.json`. Startup restores only unexpired, valid entries, preserving absolute expiration across a restart/rebuild; corrupt/oversized data fails to a cold cache. The cache has a 2,048-entry bound, with expired/soonest-to-expire eviction. Failures, unknown and unavailable outcomes are never stored as durable classifications; reader cache bypass avoids both lookup and insertion and does not erase entries. An explicit administrator TTL of zero continues to disable caching. New default is **60 minutes**; historical saved values above 60 are effectively capped to 60 at runtime without overwriting `server-settings.json` or otherwise changing user configuration. Administrator settings updates clear the cache and its disk snapshot.

Regression tests added for cold population, simulated process restart/new activation and normalized URL hits, simulated expiry without sleeps, bypass, malformed disk, unsupported entries, capacity bounds and 0/60/1440 effective settings. Other prior authentication and diagnostics tests remain.

## Validation and remaining owner acceptance

- GitHub CI for the new cache implementation was green at commit `3c94db8a92d357bfa86d67f03673f6429ad519d7`, [Actions run 42400832112](https://github.com/richard-mauri/newsgo/actions/runs/42400832112). **The later documentation/capacity commits require another exact-head CI verification**.
- This session cannot run the owner Mac or Safari, and GitHub CI cannot validate the live Homepage account/session, real feeds, local rebuild, cloud tunnel or practical performance.
- Owner still needs to validate Homepage Save/Load, administrator settings GET/PUT, ordinary user refusal, session expiry, subscriptions/exclusions persistence, and `/api/sync` returning 404.
- **Restart acceptance:** Read an identical subscription batch, stop and restart the actual NewsGo **server process** within 60 minutes, read again, verify substantial cache hits and lower classification latency. Restarting only Safari does not prove process-level persistence. Check the effective displayed TTL is 60, even if an older server setting was saved as 1440.
- Evaluate on-Mac classification-cache disk write performance under real feed batches; writes currently use atomic snapshots on each successful cache insertion. Do not claim performance acceptance solely from green CI.
- Keep the PR draft and unmerged until owner approval.

## Owner Mac commands

```bash
cd ~/projects/newsgo
git status --short
git fetch origin
git checkout feat/67-homepage-auth-diagnostics
git pull --ff-only origin feat/67-homepage-auth-diagnostics
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

## Next handoff

Check actual PR head and exact-head GitHub Actions. Perform Mac/Safari restart acceptance and real-batch performance check, review any remaining regressions, update the PR evidence, and request owner approval before merging. Do not change `main` or release without explicit authorization.
