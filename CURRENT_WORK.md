# Current work — NewsGo

## Issue #85 — diagnose repeated cache misses (November 10, 2026)

- **Merged baseline:** `main` at `83ca5081ca5f50e6651fc36cb84fe1048f0947b8`; Issue #83 / PR #84 merged and closed. New [Issue #85](https://github.com/richard-mauri/newsgo/issues/85), branch `fix/85-classification-cache-diagnostics`, draft [PR #86](https://github.com/richard-mauri/newsgo/pull/86).
- **Observed:** Repeated Safari reloads show 302 hits and 75 misses (80% hit rate) for 377 classifications with a 720-minute TTL; a prior screenshot showed 75 unknown article classifications. This suggests noncacheable outcomes, but exact per-URL match has not been confirmed.
- **Implementation:** Distinguish absent versus expired misses, explicit unknown versus unavailable returned classifications, and upstream/protocol failures in per-feed stats and Diagnostics. New screenshot confirms `investors.com` is in excluded domains but 75 unknown articles from that domain are still classified: exclusions were applied in browser after `/api/personal`, not before ScraperGo. The browser now sends the normalized excluded domains; Go validates them and short-circuits matching articles to `publisher exclusion` before classification and cache lookups. Feed fetch itself remains unchanged. Preserve no-cache policy for ambiguous/failure outcomes; no speculative negative caching.
- **Tests:** Added repeat-outcome, expired-entry, and excluded-domain no-ScraperGo-request regression tests. Go test and build passed on implementation head `6b040682444980eb9752151fdc07aa3fa73d0c2e` (GitHub Actions run 42866569601); reverify CI for this documentation update. Mac/Safari/ScraperGo live acceptance remains unverified; no release or runtime restart authorized.
- **Merge boundary:** Draft PR only; owner approval required for merge. No modifications to `main`, other apps, data, settings or Cloudflare configuration.

## Mac verification

```sh
cd ~/projects/newsgo
git status --short --branch
git fetch origin
git checkout fix/85-classification-cache-diagnostics
git pull --ff-only origin fix/85-classification-cache-diagnostics
go test -count=1 ./...
./devbuild.sh
# Restart only after separate owner approval:
./newsgo.sh
```

In Safari, repeat an All subscriptions read and reload; verify cache hits/misses and cold/expired breakdown add up, determine how many repeated misses are unknown, unavailable or upstream failures, confirm classification and filtered counts unchanged, then validate cache behavior after a deliberate service restart.
