# NewsGo — subscription-based news reader

NewsGo is a Go web service with a Safari-friendly dashboard for following RSS/Atom/JSON news feeds. It shows selected subscriptions and links articles to ScraperGo for clean reading. NewsGo does not replace ScraperGo's article extraction.

## Using the dashboard
Open `http://127.0.0.1:8084/` on the Mac running NewsGo, or use the separately configured NewsGo public hostname. Use **Find Subscriptions** to discover or add feeds, **My Subscriptions** to select a source or all sources, **Reading Preferences** to choose maximum article age, and **Excluded Sites** to manage publisher and article-destination exclusions. A selection starts fetching headlines; initial landing does not fetch every subscription automatically.

NewsGo shows only articles classified `article` by ScraperGo. The collapsed **Filtered / unprocessed articles** section lists rejected or out-of-range entries and entries beyond the per-feed processing cap. Selecting an article opens the ScraperGo reader; original links open the publisher.

Subscriptions and reading preferences remain available locally in the browser. For secure per-user synchronization, sign in through Homepage and launch NewsGo from its card. Use Account Sync to explicitly save or load your Homepage profile; once established, changes save automatically with revision conflict protection. NewsGo honors Homepage's twelve-hour session without a separate login. The header displays the signed-in reader's remaining session time (HH:MM:SS) and a **Renew** button with hover/focus feedback and inline renewal status. If Homepage's new endpoint has not been deployed, the header reports HTTP 404 with instructions to update Homepage PR #21. While the Homepage login is valid, clicking Renew uses the current parent-bound NewsGo bearer to securely extend the Homepage session by twelve hours from renewal, rotate the bearer, and update the countdown without navigation or re-entering the personal approved token. Expired or revoked sessions must sign in through Homepage again. This requires the coordinated Homepage Issue #20 / draft PR #21 deployment. Homepage administrators can open Admin Settings and request fresh filtered-article diagnostics; the NewsGo server independently verifies their Homepage administrator role. There is no Server Sync token login or legacy synchronization endpoint.

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

All readers can open **Diagnostics** after reading one feed or All Subscriptions. The panel shows the latest completed browser read: fetched feed entries, accepted articles, filtered entries, entries beyond the processing cap, feed errors and actual thumbnail load outcomes (loaded, failed, pending or not requested). Thumbnail counts update as image events occur. These browser-local observations do not require additional publisher requests. The existing **Filtered / unprocessed articles** section groups rejected items by classification or exclusion reason; entries beyond the cap remain separately labeled as unprocessed.

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

## Homepage-only identity and account storage (draft PR #68)

Issue #67 / draft PR #68 removes the old Server Sync runtime, the legacy browser subscription format converters, the token authentication and issuing command, and old `/api/sync` route. Homepage handles user identity and per-user account save/load. NewsGo's administrator endpoints verify Homepage administrator credentials server-side. Server settings use the standalone `data/server-settings.json` store, configurable with `-settings-dir` (default `data`). There is no migration or legacy compatibility path.

A zero-minute saved classification TTL disables NewsGo's classification cache, and the Diagnostics panel identifies disabled requests separately from misses. The classification cache is in memory and is cleared when the NewsGo process restarts; restarting just a Safari tab does not clear it.
