# NewsGo — RSS News Aggregator

**Current development candidate:** v0.2.25 (Go tests pass on the user’s Mac; Guardian thumbnail coverage verified after generic fallback limit adjustment; release not yet tagged). **Last confirmed pushed and tagged baseline:** v0.2.24, commit `199f191503300a27fa7196d63fc01c2838197daf`.

**Repository:** https://github.com/richard-mauri/newsgo  
**Service:** https://newsgo.themauris.org  
**Local directory:** `~/projects/newsgo` · **Port:** `8084`

## Purpose and design

NewsGo is a self-hosted Go RSS news aggregator integrated with the separate ScraperGo reader. Cached server-side publishers and user-added RSS feeds remain different *retrieval* paths, but a single subscription registry and unified reading navigation provide one subscription experience. The browser-local registry is `newsgo.subscriptions.v2`; legacy migration, import/export and all saved subscriptions are preserved. Each Safari tab or Dock web app may maintain separate website data.

The app uses Go standard-library code. Project areas: `cmd/newsgo`, `internal/dashboard`, `internal/feeds`. Diagnostic tools: `cmd/feedcheck` and `cmd/discoverycheck`. Do not modify ScraperGo, Cloudflare Tunnel, `newsgo.sh`, or `devbuild.sh` as part of a dashboard cleanup.

## v0.2.25 — Focused frontend cleanup

The v0.2.24 fix for **Read All Articles** remains intact. This release extracts a small common `readingSelection()` model from repeated conditional logic in `renderReading()` and `arrangeSharedLeads()`. It identifies the active source, subscribed server publisher eligibility, and visibility for both article groups. Per-publisher promoted lead stories continue working; individually fetched RSS feeds retain their existing Lead Stories and Latest Headlines layout. It preserves both fetch mechanisms, the browser subscription registry, and the interface. The thumbnail fallback limit is changed separately below. This is **not** a global chronological merge of cached and separately fetched RSS stories.

Existing string-contract tests were updated for the factored selection helper and for the bumped script cache-busting version. Do not treat these static source checks as substitutes for actual browser tests.

## Local verification and release handoff

The working tree contains the v0.2.25 subscription-selection cleanup and eager-loading adjustment, plus the separate generic thumbnail fallback change to 45 probes per publisher and its regression test. The last confirmed pushed and tagged baseline is v0.2.24. Before changing files, inspect `git status` and preserve local modifications.

```sh
cd ~/projects/newsgo
git status
go test ./...
./devbuild.sh
./newsgo.sh
```

Browser-check in regular Safari and a NewsGo Dock app: **Read All Articles** with only two server-cached subscriptions, with only personal RSS feeds, and mixed; read a single publisher and return to All Subscriptions; ensure no unrelated publishers appear. Check article cards, thumbnails, source labels, refresh indicator, unsubscribe, and persistent subscription counts. Do **not** clear website data during installation.

After all tests and local checks pass, inspect `git diff --check`, `git status`, stage only intended files, review `git diff --cached --stat` and `git diff --cached --check`, commit, push, and only then tag v0.2.25. Roll back to the annotated `v0.2.24` tag if necessary.

## Prior releases

- **v0.2.21:** Generic thumbnail metadata fallback for server-cached feeds.
- **v0.2.22:** Unified authoritative local subscription registry with migration and import/export.
- **v0.2.23:** Refresher status, source-scoped leads, article sections, navigation.
- **v0.2.24:** Fixed Read All Articles for cached-only, RSS-only, and mixed subscriptions. Tagged at commit `199f191`.
- **v0.2.25:** Factored common selection logic; eagerly loads all provided thumbnail URLs; increases generic shared thumbnail fallback from 8 to 45 candidate articles per publisher per refresh. No subscription storage change.


## v0.2.25 corrected test handoff

The first v0.2.25 ZIP mistakenly added `combined_subscriptions_v025_test.go` with a test name already defined by `combined_subscriptions_v024_test.go`. This corrected ZIP does **not** add a duplicate test file. Instead, it updates the existing v0.2.24 regression test to verify the v0.2.25 `readingSelection()` helper. The implementation files are unchanged from the original v0.2.25 candidate.

If the previous v0.2.25 ZIP was extracted and the duplicate test file still exists, remove only `internal/dashboard/combined_subscriptions_v025_test.go` before running `go test ./...`. The full suite must pass on the Mac before rebuilding or pushing.


## v0.2.25 thumbnail loading adjustment (October 19, 2026)

All available article thumbnails are now requested using `loading="eager"` for server-rendered cached-publication cards and `img.loading='eager'` for dynamically created RSS cards. This is not a publisher-specific exception. Image extraction rules, metadata lookup implementation, caches, article selection, subscriptions and CSS layout are unchanged; the generic fallback lookup budget is increased as described below. Articles without an image URL still have no thumbnail. Eager loading may increase page load, memory and bandwidth on long lists; verify performance in Safari and standalone Dock web app.

## v0.2.25 generic thumbnail fallback adjustment (October 19, 2026)

In `internal/dashboard/shared_images.go`, `enrichSharedImages()` now uses `perPublisher = 45` rather than `8`. The independent `maxConcurrent = 8` setting remains unchanged. The fallback fetches article-page Open Graph or Twitter image metadata only when the feed did not provide a usable image URL. The limit is **per publisher per shared-feed refresh**, not a guaranteed number of successful images. It applies uniformly to all publishers, with no Guardian-specific branches.

The generic fallback retains its **12-second overall deadline**, **2.5-second per-article request timeout**, and **four-hour cache for both successful and unsuccessful lookups**. Enlarging the budget can increase metadata requests, and the deadline can still limit how many finish. No image bytes are fetched by this server-side metadata lookup. Browser eager loading is independent: it only changes when existing thumbnail URLs are requested.

`internal/dashboard/shared_images_test.go` now constructs 45 image-less articles and verifies that the fallback can dispatch 45 bounded lookups. The user reported `go test ./...` passing after this change and confirmed all Guardian thumbnails appeared in NewsGo, compared with 8 present and 37 missing before the budget increase. This is a local observation, not a guarantee that every article on every feed will have metadata.

## Development and release notes

The previous v0.2.25 handoff briefly included a duplicate `combined_subscriptions_v025_test.go` test; the corrected candidate updates `combined_subscriptions_v024_test.go` instead. If the obsolete duplicate file remains in the local tree, remove that file only. Do not clear Safari/Dock website data or localStorage. The current README update is standalone; no ZIP needs to be reapplied. Run `git diff --check`, inspect `git status` and `git diff`, verify Safari/Dock application behavior, then stage only intended files, commit, push, and tag v0.2.25 after the push.
