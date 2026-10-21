**Development correction:** Fixed missing date check in the classification pipeline: articles older than the selected maximum age are now marked `too-old` before invoking ScraperGo. Added `article_age_pipeline_test.go` to prevent recurrence. Pending CI/Mac verification.

## User reading preference: Maximum Article Age (development)

The Reading Preferences panel offers All (default), 1, 3, 7, 14, 30, or 60 days. Selected age is stored per browser and synchronized with the authenticated user's existing Server Sync profile. Older dated feed items are excluded before expensive ScraperGo classification and appear in the Filtered / unprocessed articles card with reason `too-old`. Undated entries are retained. Existing subscriptions, exclusions, and Admin Settings are unchanged. CI and live Safari acceptance are required before merging to main.

## NewsGo v0.2.53 development: filtered article audit and operational controls

The briefing now offers a collapsible **Filtered / unprocessed articles** card. It reports headlines rejected by ScraperGo classification with their classification codes, items hidden by saved exclusion domains, and a separate count of feed entries not classified because they exceed the configured per-feed cap. It does not claim to diagnose entries absent from the RSS feed. Admin Settings now exposes four **active** persisted controls: maximum articles per feed (default 60, range 1–250), classification timeout (7 seconds, 2–30), parallel classifications (8, 1–16), and successful-feed cache TTL (10 minutes, 1–60). Older saved settings migrate transparently. Historic image enrichment settings are clearly labeled inactive because thumbnail delivery is now managed by ScraperGo. The API limit of eight feeds per request and the 4 MiB feed payload safety cap remain fixed. Changes have passed GitHub Go tests/build, but browser acceptance is pending; no production promotion intended.

# NewsGo v0.2.52 — strict classification gate (development)

Requires ScraperGo v0.2.65. Shows headlines and eager thumbnails only when `/api/classify` confirms `article`. Suppresses homepage, unreadable, unavailable, and unknown during the current reading session. Transient failures can be retried on subsequent selections; nothing is unsubscribed. Preserves v0.2.51 eager thumbnail loading and on-demand welcome/startup behavior. Does not repair backend image-fetch failures. Test full local suite and browser before commit.

# NewsGo v0.2.51 — eager thumbnail loading (development candidate)

Changes Safari card thumbnail HTML from lazy to eager loading, so all rendered images request immediately rather than appearing only when scrolled into view. Preserves ScraperGo URLs, the no-auto-fetch welcome screen, single-feed selection and v0.2.50 classification filtering. This will cause more simultaneous image requests and does not guarantee availability where ScraperGo returns an error or cannot select an image. Install direct-unzip, test/build/restart on Mac, and verify in Safari before commit.

# 0.2.50 development candidate — 2026-10-21

Startup waits for explicit subscription or All selection. Only selected feeds fetched; suppress confirmed homepage/unreadable classifications.

Install ScraperGo before NewsGo. Run full Mac tests/build and browser acceptance prior to commit. Preserve launcher/tunnel and local runtime state.


## Development candidate v0.2.49 — 2026-10-21

NewsGo POST /api/personal now classifies up to 60 feed entries per source using ScraperGo /api/classify before response. Confirmed homepage entries are suppressed. Unknown entries are retained. Eight concurrent classification requests and a 7-second per-request timeout; classification caching 30 minutes. This adds latency on cold feeds and awaits live load testing.

## v0.2.48 — ScraperGo thumbnail integration (development candidate)

NewsGo delegates article-thumbnail discovery, downloading, resizing and caching to ScraperGo `/api/thumbnail`. The NewsGo `/api/personal` response now supplies a ScraperGo thumbnail URL for each valid article URL instead of inspecting publisher HTML or returning publisher image URLs. Dashboard thumbnails use lazy loading, with the existing image-error fallback. Feed fetching, subscriptions, exclusions, and server sync are unchanged. ScraperGo must have a reachable HTTPS `/api/thumbnail` route through its existing tunnel. Server image-retrieval tuning fields from v0.2.47 are retained for backward-compatible server-settings JSON, but no longer govern NewsGo thumbnail lookups; future settings cleanup should be handled separately.

**Development verification:** run `go test ./...`, `./devbuild.sh`, and live Safari tests on the user's Go 1.27.1 Mac before commit or deployment. Test CNBC, Seahorse, no-image stories, and publisher exclusions. Image requests are now browser-to-ScraperGo; NewsGo does not fetch thumbnail bytes. ScraperGo cache and rate limits require attention for high headline volumes.

# CURRENT CANDIDATE — NewsGo v0.2.47 admin server settings (2026-10-21)

Based on authoritative user-provided source matching GitHub commit `598ec5b`. This direct-unzip candidate adds local CLI `-set-role admin` / `-set-role user` to designate an existing hashed token without rotating it, admin-guarded `/api/admin/settings` GET/PUT, separately revisioned atomic `data/server-settings.json`, and a browser Admin Settings panel revealed to the connected admin token. Image metadata timeout, maximum HTML response, concurrent lookups, enrichment deadline, and missing-image cache TTL now use bounded live server settings. No changes to user preferences or secrets. The default values remain 2500 ms, 512 KiB, 12, 25 seconds, 90 seconds. Tests passed in an isolated Go 1.23.2 compatibility run (temporary go.mod only, original go.mod retained); the required Go 1.27.1 project tests, ./devbuild.sh and live Safari acceptance remain unverified. No Git commit, push, restart or release authorized. ScraperGo and NewsValidate will be handled separately after the CNBC investigation. The existing historical handoff below is retained for context; old "next step" statements are superseded by this paragraph.

## NewsGo administrator settings candidate — October 21, 2026

Based on the user's clean commit `598ec5b` and source archive. Server tokens support optional `role: "admin"`, assigned locally with `go run ./cmd/newsgo-token -dir data -user richard -set-role admin` (use the existing user ID; no token rotation). `GET/PUT /api/admin/settings` require the token's admin role, use revision checks, JSON atomic writes, input bounds, and server-wide image lookup parameters. The Admin Settings browser panel is available only after admin Server Sync connection; authorization is enforced server-side. Existing subscriptions, exclusions, tokens, data files and tunnels are not replaced. Candidate requires full local tests and browser acceptance prior to commit or deployment. ScraperGo and NewsValidate remain separate projects in the planned shared extraction architecture.

---

## v0.2.46 — Executable version reporting (development patch, October 21, 2026)

The NewsGo executable now supports `-version` and prints `NewsGo v0.2.46` without loading feeds or starting a server. The launcher `newsgo.sh` reports the exact `dist/newsgo-darwin-amd64` path and checks its executable version **before** terminating any running NewsGo process. The startup log uses the dashboard's current version instead of its previous stale v0.2.34 message. This is a narrowly scoped v0.2.46 follow-up and does not alter subscriptions, sync, exclusions, or article rendering. Local full-project Go tests, Intel builds, and runtime behavior remain to be verified on the user's Mac. Use the established `./devbuild.sh` build command; its output goes to `dist/`.

## v0.2.44 — Missing-thumbnail retry and lookup budget

The Seahorse Magazine Finn article had no thumbnail even after the v0.2.44 body-image fallback. Preserve feed/metadata/body extraction; allow 12 parallel metadata lookups within a 25-second shared deadline so late article lookups are less likely to be abandoned. Cache successful feeds containing image-less articles for 90 seconds rather than ten minutes to recover from transient metadata failures. Regression covers the Finn lead-image HTML and retry TTL. Live publisher image retrieval remains a browser acceptance test; no ScraperGo changes or data migration.

## v0.2.44 — Article image fallback

NewsGo now inspects article-body image tags when the feed and social metadata have no thumbnail. This addresses the Seahorse Magazine image in Sailing Anarchy’s “Control from the cockpit” link. Existing image metadata, safe network fetching, subscriptions, and exclusions are unchanged.

## v0.2.42 — exclude article destination domains

**Excluded Sites** now manages two independent lists: excluded publisher domains (the existing behavior) and excluded article domains. Add `youtube.com` to **Excluded Article Domains** to hide linked YouTube stories from any subscribed RSS feed, including KCRA, without removing the subscription or hiding non-YouTube articles. `youtu.be` short links are treated as YouTube. Use **Allow articles** to undo the exclusion. The button count includes entries from both exclusion lists. Both exclusion lists are retained in local browser storage, preference export/import, and token-authenticated server-side JSON synchronization. Existing JSON user profiles require no migration; omitted `excludedArticleDomains` means no article exclusions. Automatic saves, revision/conflict checks, and safe reconnect/merge now include article domains.

**v0.2.42:** Corrected the visible Server Sync token-entry placeholder to “Enter the 64-character approved token distributed by Richard.” No change to token handling, sync, subscriptions, or exclusions.

**v0.2.42:** Presentation-only changes: header tagline is “NEWS READING WITHOUT PAGE CLUTTER USING SCRAPERGO” and invalid/absent sync token guidance reads “Enter the 64-character approved token distributed by Richard.” No change to synchronization, subscription data, exclusions, or ScraperGo.

**v0.2.38:** Removes the redundant upper-right “Updated subscribed feeds” status. The detailed line showing updated feed counts and visible headlines remains the single source of refresh progress. No changes to subscriptions, exclusions, or server sync.

**v0.2.37:** Removes the redundant standalone ← All Subscriptions button from the selected-feed reading view. Use **My Subscriptions → Read All Articles** to return to the combined feed. Existing feed selection and exclusion fallback still use `chooseAllSubscriptions()` internally. No changes to browser/server persistence, sync, or ScraperGo.

**v0.2.36:** Excluded Sites now shows the saved publisher name(s) (for example, Hacker News) alongside the underlying excluded domain (`news.ycombinator.com`), and refreshes the panel when opened. Domain keys, counts, preferences, and sync APIs are unchanged.

# NewsGo — Subscription-Based News Reader

NewsGo brings together headlines from news sources **you choose to subscribe to**. The dashboard presents one unified briefing: there is no separate built-in news stream, and no unsolicited headlines from the source catalog.

To begin, open **Find Subscriptions**, browse the available sources or discover a publication's RSS feed, and choose **Subscribe**. Your selections appear in **My Subscriptions** and all their non-excluded headlines appear together in **Today's briefing**. With no subscriptions, the briefing stays empty until you choose sources.

**Excluded Sites** lets you hide publishers by domain without deleting subscriptions. An excluded subscription can be restored with **Allow this site**. Preferences are kept in browser storage and may be exported and imported with the dashboard controls. They are synchronized across devices only when optional Server Sync is connected.

Headlines open in [ScraperGo](https://scrapergo.themauris.org/), a separate article reader. NewsGo does not alter ScraperGo's extraction behavior.

The service is written in Go, uses an RSS/Atom/JSON feed catalog (`feeds.json`) for discoverability, and presents article headlines from the feeds you have subscribed to. Provider access and publication frequencies vary. NewsGo is an aggregator, not a subscription/paywall bypass.

**Development baseline:** v0.2.46 (October 21, 2026), committed on GitHub in `59322fa`, with handoff update `9766238`. Live executable version and manual-search acceptance were not established by the handoff. This patch adds version introspection; it has not yet been installed, built, or browser-tested on the Mac.

**Developers and future AI agents:** See [`AGENTS.md`](AGENTS.md) for the authoritative current engineering handoff, development procedure, version synchronization rules, and Git safeguards.

### v0.2.33 discovery and cleanup recap

Removed the obsolete server-side catalog headline fetching and image enrichment path. Improved Feedly publisher-name lookup with a prefix fallback when no results are returned, clearer external-feed HTTP error messages, and discovery result card layout. The feed catalog remains a source of discoverable subscriptions, not an automatic headline stream.

## Optional server synchronization (v0.2.34)

The NewsGo server can persist **both subscriptions and excluded sites** per approved user token.
The same token works across MacBook Safari, Docked Safari, and iPhone Safari.
Token entry is optional; browser-only operation remains available.

The server defaults to `-sync-dir data` relative to its working directory. Do not commit this directory; it holds user preferences and token digests. Back it up securely.

Issue the first token on the MacBook running NewsGo:

```sh
go run ./cmd/newsgo-token -dir data -user richard -label 'Primary user'
```

Copy the printed secret token into **Server Sync** on each trusted browser once. The first connection automatically uploads both local lists if no server profile exists; existing server profiles load on Connect. Both subscriptions and excluded sites save automatically after edits. Other devices check for updates every 30 seconds while visible and when returning to the tab. **Advanced** contains optional Load from Server, Save Now and Disconnect controls. Stale revisions cannot overwrite newer server data; after a conflict use Advanced > Load from Server (with confirmation when unsaved local changes exist).

Revoke a user token: `go run ./cmd/newsgo-token -dir data -user richard -revoke`. Reissue after revocation if necessary. Tokens are stored as SHA-256 digests; each token should be treated as a secret. The browser remembers the token in localStorage for automatic reconnection. Only enter tokens on trusted devices: browser storage is accessible to scripts running on this origin. Disconnect erases the locally remembered token, not the server account. The server must be reachable over HTTPS through your existing Cloudflare Tunnel. No additional hosted database is required.

### v0.2.42: safe reconnection

When a valid token already has preferences on the server and the local browser lists differ, NewsGo offers **Use Server Preferences**, **Merge Preferences**, and **Cancel**. The decision applies to both subscriptions and excluded sites. Merge uses subscription IDs and unique excluded domains, and saves against the server revision; stale writes still receive HTTP 409. A new server profile is initialized from browser preferences as before. This applies to remembered-token reconnects as well as manual Connect.

## v0.2.45 — Subscription search feedback

Subscription search preserves old results during new searches and shows pending/waiting/retry messages. Browser pacing reflects the existing three-second per-IP uncached-search limit, with bounded automatic retries for HTTP 429. Obsolete responses cannot replace results for a newer query. No changes to the feed-search API, subscriptions, or exclusions.

### v0.2.46 corrected release gate
Version-sensitive tests must derive expectations from the displayed application version rather than encode a release number. Tests for changed behaviors must be updated with the implementation, not suppressed. The explicit subscription search requires a button or Enter; keystrokes and initial load must not trigger discovery requests. Run the complete `go test ./...` using the actual project toolchain; build through `./devbuild.sh` so the launcher uses the new executable. Do not describe a ZIP as verified unless tests have passed.
