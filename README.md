# NewsGo — RSS News Aggregator

**Current development candidate:** v0.2.24 (not yet tested on the user's Mac or committed). **Last documented pushed baseline:** v0.2.20; local v0.2.23 changes were previously reported pending commit. Verify `git status` before replacing any locally modified files.

**Repository:** https://github.com/richard-mauri/newsgo  
**Service:** https://newsgo.themauris.org  
**Local directory:** `~/projects/newsgo` · **Port:** `8084`

## Purpose and architecture

NewsGo is a lightweight self-hosted Go RSS news aggregator with article thumbnails and a Safari-compatible browser dashboard. The companion ScraperGo application opens the full articles and follows a separate release lifecycle. NewsGo uses server-cached publications plus individually requested RSS feeds; the different retrieval mechanisms are implementation details, not separate user-facing subscription systems.

Source layout: `cmd/newsgo`, `internal/dashboard` (HTML, JavaScript, subscription UX), `internal/feeds`. Diagnostic tools: `cmd/feedcheck` and `cmd/discoverycheck`. Go standard library is preferred. The existing `go.mod` specifies Go 1.27.1. The app runs on a 2020 Intel MacBook Air behind Cloudflare Tunnel.

## v0.2.24 — Read All Articles correction

The v0.2.23 **Read All Articles** action cleared the selected publication but then hid the server-cached article section unconditionally; if no separately fetched RSS subscriptions existed, *both* article sections were hidden. This affected fresh Safari Dock instances with only built-in publisher subscriptions. The article markup was present in the DOM, but CSS hid its parent container.

v0.2.24 corrects the selection logic. When no specific source is selected, NewsGo displays articles from subscribed server-cached publishers and subscribed individually fetched RSS feeds; unsubscribed server publishers are filtered out. Selecting an individual publication still scopes the article view to that publication. The existing per-publication lead promotion and personal RSS lead/latest grouping are retained. With zero subscriptions, the page provides a Find Subscriptions prompt.

The two feed retrieval and rendering sections remain separate internally; this patch does **not** create one globally timestamp-sorted mixed-source timeline. The goal is to fix the known defect without expanding the scope of the change.

The v0.2.22 browser subscription registry at `localStorage['newsgo.subscriptions.v2']` is untouched, including its v1 migration and import/export. Do not clear website data to deploy this patch, as that would discard local subscriptions. Built-in feeds remain suggestions in Find Subscriptions; personal and built-in subscriptions behave identically for read, unsubscribe, and Read All.

## Regression tests and handoff

The first v0.2.24 patch failed six older tests because those tests searched for exact v0.2.23 version strings or the old shared-only selection behavior. This corrected patch updates the six test assertions to v0.2.24 and preserves the original checks for discovery, navigation, lead promotion, initial hidden state, and Tanner. These are compatibility checks, not a substitute for browser tests.

## Handoff: minimal ZIP patch

This ZIP contains only:

- `internal/dashboard/personal_script.go`
- `internal/dashboard/dashboard.go`
- `internal/dashboard/combined_subscriptions_v024_test.go`
- Six updated existing regression tests (`browse_test.go`, `catalog_test.go`, `discovery_presentation_test.go`, `shared_lead_selection_test.go`, `startup_view_test.go`, `tanner_test.go`) so their expectations track v0.2.24
- `README.md`

No backend feed code, `.git`, `go.mod`, tunnel settings, shell scripts, or ScraperGo files are included. The two Go source files were derived directly from the user's uploaded v0.2.23 files, not from an assumed older GitHub baseline.

After protecting any local changes, install and verify on the Mac:

```sh
cd ~/projects/newsgo
git status
unzip -l ~/Downloads/newsgo-v0.2.24-patch.zip
unzip -o ~/Downloads/newsgo-v0.2.24-patch.zip
gofmt -w internal/dashboard/*.go
go test ./...
./devbuild.sh
./newsgo.sh
```

First verify with a fresh Dock app with two built-in publisher subscriptions: the **Read All Articles** view must show only those publishers, and individual Read Articles must continue working. Also test two RSS feeds, mixed built-in/RSS subscriptions, no subscriptions, unsubscribe behavior, thumbnails, refresh status, and the regular Safari version. Confirm local v0.2.24 code is actually running before evaluating Dock behavior. Existing Safari and Dock web apps have separate local storage.

Once verified, inspect `git status`, `git diff --check`, `git diff --cached --stat`, and `git diff --cached --check` before committing and pushing. Do not tag until after push. Keep the rollback copy or the last good git commit available. **A test-source contract check is not a substitute for a full browser regression test.**

## Earlier development decisions

- **v0.2.20:** Last documented pushed baseline (verify repository before release).
- **v0.2.21:** Generic article-level image fallback for cached publishers; up to eight recent missing images per publisher, shared 12-second deadline, four-hour positive/negative cache.
- **v0.2.22:** One authoritative subscription registry, valid import/export and legacy migration; avoid redesigning storage.
- **v0.2.23:** Per-publication lead selection, Lead Stories and Latest Headlines for RSS subscriptions, refresh spinner and progress status, All Subscriptions navigation; pending local commit at earlier handoff.
- **v0.2.24:** Combined reading view works for built-in-only, RSS-only, and mixed subscriptions, preserving existing infrastructure.

Historical README text in GitHub remains available at https://github.com/richard-mauri/newsgo/blob/main/README.md. This refreshed README replaces the long early FeedCheck transcript with the current operational handoff and architecture.
