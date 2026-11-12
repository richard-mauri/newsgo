# NewsGo — subscription-based news reader

NewsGo is a Go web service with a Safari-friendly dashboard for following RSS/Atom/JSON news feeds. It shows selected subscriptions and links articles to ScraperGo for clean reading. NewsGo does not replace ScraperGo's article extraction.

## Using the dashboard
Open `http://127.0.0.1:8084/` on the Mac running NewsGo, or use the separately configured NewsGo public hostname. Use **Find Subscriptions** to discover or add feeds, **My Subscriptions** to select a source or all sources, **Reading Preferences** to choose maximum article age, and **Diagnostics → Excluded Domains** (expandable panel) to manage publisher and article-destination exclusions. A selection starts fetching headlines; initial landing does not fetch every subscription automatically.

NewsGo shows only articles classified `article` by ScraperGo. Within **Diagnostics**, the collapsed **Filtered / unprocessed articles** section lists rejected or out-of-range entries and entries beyond the per-feed processing cap. The main briefing shows only lead stories and latest headlines, without filtering controls. Selecting an article opens the ScraperGo reader; original links open the publisher.

Subscriptions and reading preferences remain available locally in the browser. For secure per-user synchronization, sign in through Homepage and launch NewsGo from its card. NewsGo automatically loads your verified Homepage account profile after sign-in, then saves subsequent subscription, exclusion, and reading-preference changes to your per-user account with revision conflict protection. No manual Account Sync action is needed. If a network error prevents reading the account, browser settings are retained and no account data is overwritten. NewsGo honors Homepage's twelve-hour session without a separate login. The NewsGo header shows the Homepage-authenticated reader identity, without an app-specific countdown or Renew button. Manage the twelve-hour login, renewal timer, and Disconnect through Homepage's **Manage tokens → Current session** panel. NewsGo retains only a tab-scoped Homepage handoff. If that temporary bearer is absent on a direct NewsGo visit (for example, after closing and reopening Safari), NewsGo makes a top-level navigation to Homepage's fixed `/launch/newsgo` route. A bridge running on Homepage's origin verifies the existing approved 12-hour login using its retained opaque session ID, falling back to the host-only cookie, then delivers a fresh bound bearer in the NewsGo URL fragment. If the Homepage session is absent, revoked or expired, NewsGo resumes browser-only mode without changing local preferences or looping. NewsGo also rechecks Homepage's profile endpoint at the cached bearer expiry, so a still-valid bearer extended by Homepage is not discarded prematurely. This recovery does not start or extend a Homepage login. Expired or revoked sessions must sign in through Homepage again. Homepage administrators can open Admin Settings and request fresh filtered-article diagnostics; the NewsGo server independently verifies their Homepage administrator role. There is no Server Sync token login or legacy synchronization endpoint.

## Build and run on the owner's Mac

From the project folder:

```bash
cd ~/projects/newsgo
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

`devbuild.sh` builds `dist/newsgo-darwin-amd64` and `dist/newsgo-darwin-arm64`; `newsgo.sh` launches the binary for the Mac's architecture on `127.0.0.1:8084` with NewsGo's dedicated Cloudflare tunnel. The launcher expects its established local `.cloudflared/token` and cloudflared executable; do not store tokens in Git. The ScraperGo reader defaults to `https://scrapergo.themauris.org` (override with `-scrapergo` when intentionally developing).

```bash
./dist/newsgo-darwin-amd64 -version
```

The CLI, launcher, server log, and dashboard display `NewsGo v0.2.58 · Build NNNNNN` on the present development candidate. The six-digit build counter is local to the Mac's Git metadata and increments once after successful `./devbuild.sh` completion; it is not a semantic release version or a Git-tracked file. Builds without injected metadata show `000000`.

`./devbuild.sh --release vX.Y.Z` has additional release safeguards and is **not** part of ordinary development; only run it after explicit owner approval on the release branch.

## Configuration and care
- `feeds.json` supplies the discovery catalog, not the user's saved subscriptions.
- `data/server-settings.json` contains operational settings and is the only NewsGo server-side JSON data store. Homepage holds account preferences; browser `localStorage` holds any unsaved preferences. The old `approved-tokens.json` and `user-preferences.json` files are unused and can be deleted from your Mac after you stop NewsGo, if present.
- `newsgo.log`, `dist/` binaries, and `.cloudflared/token` are machine-local runtime artifacts. Never publish secrets or overwrite data during updates.
- Feed and article availability, thumbnails, and publisher blocking can vary. Live Safari checks are required even when Go CI passes.

## Documentation map
[AGENTS.md](AGENTS.md) specifies development rules. [ARCHITECTURE.md](ARCHITECTURE.md) describes implemented internals. [CURRENT_WORK.md](CURRENT_WORK.md) records the active candidate and acceptance status. Historical development is available through Git history rather than standalone version-note files.

## Reading diagnostics

All readers can open **Diagnostics** after reading one feed or All Subscriptions. Expand the collapsed **Latest read** panel to inspect the latest completed browser read: fetched feed entries, accepted articles, filtered entries, entries beyond the processing cap, feed errors and actual thumbnail load outcomes (loaded, failed, pending or not requested). Thumbnail counts update as image events occur. These browser-local observations do not require additional publisher requests. The existing **Filtered / unprocessed articles** section groups rejected items by classification or exclusion reason; entries beyond the cap remain separately labeled as unprocessed.

Timing diagnostics show browser-observed refresh and batch durations, browser render duration, sum of per-feed fetch/wait and classification durations, and mean thumbnail completion-event latency. Summed server per-feed durations are not additive wall-clock stages. The panel additionally reports NewsGo classification-cache hits, misses, bypasses, requests disabled by the server, configured TTL and cache hit rate for each refresh. It does not report ScraperGo's separate internal cache statistics.

### Optional RSS metadata-only mode

In **Reading Preferences**, enable **Use RSS metadata only** to display RSS/Atom/JSON feed titles, dates and safe images while still requiring ScraperGo's `article` classification for every eligible link. This mode changes presentation, not article eligibility; it does not bypass classification. The preference is browser-local. Opening a headline uses ScraperGo.

Reading Preferences also offer browser-local Maximum Lead Stories (1–20, default 3) and Lead Stories only in All Subscriptions (default off). The latter skips Latest Articles rendering, not feed fetching or classification. These coexist with RSS metadata only and admin Diagnostics.

## Verified article filtering and shared cache

Both normal reading and **Use RSS metadata only** classify eligible feed links with ScraperGo before presenting headlines; only confirmed `article` classifications are displayed. Non-article classifications, including `homepage`, `unreadable`, `unavailable`, and `unknown`, are shown in the filtered-articles audit. RSS metadata-only mode still uses feed-provided headlines, publication dates, and safe RSS images rather than ScraperGo thumbnails.

Reading Preferences includes **Use classification cache**, enabled by default. Cache entries are shared server-side across reading modes and readers. Disabling this preference bypasses both lookup and insertion for that reader's requests without clearing the shared cache; ScraperGo itself may implement separate internal caching. Confirmed stable classifications `article`, `homepage`, and `unreadable` are eligible for caching, while `unavailable` and `unknown` are not. New server settings default classification TTL to 1440 minutes (24 hours), configurable by administrators from 0–1440 minutes. Existing saved explicit TTL settings are preserved; they are not silently overwritten.

## Administrator article URL export (diagnostic)

After a completed **Read All Articles** refresh, the **Download Article URLs** button appears in the briefing **only for a verified Homepage administrator signed in through NewsGo's Homepage handoff**. It is hidden for ordinary users, signed-out browsers and expired sessions, and its click handler rechecks the admin role before exporting. The control uses the normal NewsGo action styling. On localhost without a signed-in Homepage session, the export action remains hidden. This is a client-side export of already visible article URLs, not a server-side privileged download API.

Click to save `newsgo-scrapergo-urls.json` (which may contain publisher URL query strings; keep it private). With ScraperGo running, use the bundled sequential diagnostic from `~/projects/newsgo`: `go run ./cmd/scrapertrace -input ~/Downloads/newsgo-scrapergo-urls.json -limit 10` for a small initial sample; to replay the complete filtered export, specify its actual URL count with `-limit N` (maximum 500), and retain the default `-delay 1s` pacing. The tool calls `/api/classify` and `/api/read` separately for each selected URL and reports classification, both HTTP statuses, bounded failure kind, request wall times, and aggregate latency. It defaults to 10 URLs, one sequential worker, and a 1-second inter-URL delay for system safety. On a sufficiently capable machine, optionally use `-workers 4 -delay 0s` for a bounded four-way batch; classification and read remain sequential within each worker. Concurrent output arrives in completion order. Avoid running this on underpowered machines or when ScraperGo is already busy. The export contains only NewsGo-accepted headlines, not rejected entries. From `~/projects/scrapergo`, run `go run ./cmd/recoveryreport` for cumulative results. Exporting does not automatically submit article requests.

To identify **candidate duplicate text** in ScraperGo's extracted article blocks, use the separate, version-controlled `cmd/duplicateaudit` tool. From `~/projects/newsgo`, run `go run ./cmd/duplicateaudit -input ~/Downloads/newsgo-scrapergo-urls.json -limit 10 -output ~/Downloads/duplicate-audit-sample.csv` against local ScraperGo. For the approved public endpoint, add `-scrapergo https://scrapergo.themauris.org/api/read -public`. Its sequential, paced requests produce a private CSV of exact or substantial overlapping text, article URL, block indices and kind. Failed reads are reported separately; candidates may include legitimate repetition. See [duplicate-audit usage and limitations](cmd/duplicateaudit/README.md). This is development tooling only and does not change the NewsGo service.

## Homepage-only identity and account storage (draft PR #68)

Issue #67 / draft PR #68 removes the old Server Sync runtime, the legacy browser subscription format converters, the token authentication and issuing command, and old `/api/sync` route. Homepage handles user identity and per-user account save/load. NewsGo's administrator endpoints verify Homepage administrator credentials server-side. Server settings use the standalone `data/server-settings.json` store, configurable with `-settings-dir` (default `data`). There is no migration or legacy compatibility path.

A zero-minute saved classification TTL disables NewsGo's classification cache, and the Diagnostics panel identifies disabled requests separately from misses. The classification cache is in memory and is cleared when the NewsGo process restarts; restarting just a Safari tab does not clear it.

### Raw feed responses in Diagnostics

After a completed read, open **Diagnostics → Raw feed responses** to inspect the original response body for each selected subscription. The inspector retains the exact bytes NewsGo received from its normal feed HTTP request (after HTTP transfer/content encoding handled by Go's HTTP client), before RSS/Atom/RDF/JSON parsing. The source link opens that publisher's actual feed URL in a separate tab. The raw feed response area is a collapsed, expandable panel in Diagnostics, matching the Excluded Domains and Filtered / unprocessed articles panels. Each subscription has an individual expandable plain-text view, **Copy Raw Feed** and **Download Raw Feed**; download uses the preserved bytes without changing line endings or encoding. This includes NewsGo's previously cached feed body if its normal feed cache was used, with a `(cached)` label. No additional fetch is made. Browser data is replaced on the next completed read and is not persisted by this feature. Feed URLs and source text may contain private query values; handle downloads or copied text accordingly. This replaces the prior Article Date Diagnostics panel. The existing general Diagnostics counts and timing remain available.

### Grouped filtered-article audit

In Diagnostics → Filtered / unprocessed articles, each classification tab groups filtered article rows by their normalized publisher URL domain and shows a count per domain. A single **Exclude domain** control appears per non-excluded valid domain, rather than repeated controls for every article. Choosing it uses the existing exclusion workflow, including confirmation before removing affected subscriptions, saves the reader's exclusion settings, and refreshes the reading results. No domain is excluded based solely on high filtered counts; the individual article links and optional administrator fresh-classification diagnostics remain available. Unknown or invalid domains have no exclusion button.

### Understanding classification cache misses

Diagnostics → Latest read now separates **cold / absent** and **expired** classification-cache entries, and reports misses that ended in an **unknown** or **unavailable** classification or an **upstream/protocol failure**. Only successful `article`, `homepage` and `unreadable` results are cached. As a result, recurring unknown/unavailable classifications or transient upstream failures can keep the overall hit rate below 100%, even with unchanged feeds and a 720-minute cache TTL. The counters identify categories, not article URLs, and do not imply every miss is a defect. This diagnostic change does not negatively cache ambiguous results or alter classification eligibility.

### Excluded domains bypass classification

When a domain is excluded under Diagnostics → Excluded Domains, NewsGo passes its normalized browser/authorized Homepage preference exclusions to the server with the reading request. After fetching a feed, but **before** classifying each article with ScraperGo, the server marks matching publisher URLs as `publisher exclusion`. Those articles do not make classification requests and cannot generate cache misses. The filtered audit still lists them as excluded. Article host and subdomain matching is boundary-safe (excluding `investors.com` does not exclude `notinvestors.com`). The RSS feed itself may still be fetched; this feature avoids unnecessary *article classification* rather than automatically unsubscribing from a feed.

## Homepage suite activity reporting (Issue #89, draft)

NewsGo can optionally send **aggregate** visitor/processing activity to Homepage's pending admin activity collector. When both environment variables `NEWSGO_ACTIVITY_URL` (for example `http://127.0.0.1:8788/api/admin/activity/ingest` in a separate local test) and `HOMEPAGE_ACTIVITY_KEY_NEWSGO` (a unique secret of at least 32 characters, **server environment only**) are configured, NewsGo serves an ephemeral per-tab browser presence script, observes same-origin `POST /api/activity/ping` requests, and sends one aggregate report to Homepage approximately every minute. Without both, reporting is off and ordinary NewsGo usage is unchanged. The Homepage collector must have the **same NewsGo key** in its private environment, with `SUITE_DATA_DIR` enabled, for the reports to be accepted. Never put the key into browser JavaScript, localStorage, settings or source control. Only HTTPS and explicit loopback HTTP collector URLs are accepted; collector redirects are refused.

Presence reports include 2/5/15/60-minute tab counts, foreground counts and coarse anonymous/unverified counts. The NewsGo server also observes concurrent in-flight `/api/personal` feed-reading requests, which can trigger outbound feed and ScraperGo work; however, detached or separately triggered ScraperGo tasks will require separate ScraperGo instrumentation. **Unverified** NewsGo tabs currently contribute to the `anonymous` field even if their user is signed in to Homepage; a later change should add server-verified account classification without trusting a browser role claim. No article URLs, user IDs, session IDs or browser IDs are transmitted to Homepage. This is advisory telemetry, not accurate hourly unique-visitor counts or a restart permission signal. No service is restarted or released by this draft PR.
