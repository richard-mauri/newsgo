# NewsGo — subscription-based news reader

NewsGo is a Go web service with a Safari-friendly dashboard for following RSS/Atom/JSON news feeds. It shows selected subscriptions and links articles to ScraperGo for clean reading. NewsGo does not replace ScraperGo's article extraction.

## Using the dashboard
Open `http://127.0.0.1:8084/` on the Mac running NewsGo, or use the separately configured NewsGo public hostname. Use **Find Subscriptions** to discover or add feeds, **My Subscriptions** to select a source or all sources, **Reading Preferences** to choose maximum article age, and **Excluded Sites** to manage publisher and article-destination exclusions. A selection starts fetching headlines; initial landing does not fetch every subscription automatically.

NewsGo shows only articles classified `article` by ScraperGo. The collapsed **Filtered / unprocessed articles** section lists rejected or out-of-range entries and entries beyond the per-feed processing cap. Selecting an article opens the ScraperGo reader; original links open the publisher.

Subscriptions and preferences are stored in the browser. For personal account settings, sign in through Homepage and launch NewsGo from its card; Homepage's single 12-hour session authorizes per-user save/load, with no app-specific countdown or Renew button. **Server Sync** can also be used in its existing administrative/token workflows. On a difference, choose whether to use server preferences or merge carefully; don't overwrite a desired profile. An authenticated administrator can manage server-wide article settings, inspect fresh per-article diagnostics under filtered entries, and **Exit Admin Mode** to disconnect the current browser. Disconnecting preserves locally held reading preferences and does not revoke the token. Reconnect through Server Sync to regain access.

## Build and run on the owner's Intel Mac

From the project folder:

```bash
cd ~/projects/newsgo
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

`devbuild.sh` builds `dist/newsgo-darwin-amd64` and `dist/newsgo-darwin-arm64`; `newsgo.sh` launches the Intel binary on `127.0.0.1:8084` with NewsGo's dedicated Cloudflare tunnel. The launcher expects its established local `.cloudflared/token` and cloudflared executable; do not store tokens in Git. The ScraperGo reader defaults to `https://scrapergo.themauris.org` (override with `-scrapergo` when intentionally developing).

```bash
./dist/newsgo-darwin-amd64 -version
```

The CLI, launcher, server log, and dashboard display `NewsGo v0.2.58 · Build NNNNNN` on the present development candidate. The six-digit build counter is local to the Mac's Git metadata and increments once after successful `./devbuild.sh` completion; it is not a semantic release version or a Git-tracked file. Builds without injected metadata show `000000`.

`./devbuild.sh --release vX.Y.Z` has additional release safeguards and is **not** part of ordinary development; only run it after explicit owner approval on the release branch.

## Configuration and care
- `feeds.json` supplies the discovery catalog, not the user's saved subscriptions.
- `data/` holds local tokens and synchronized profiles; treat as private persistent state.
- `newsgo.log`, `dist/` binaries, and `.cloudflared/token` are machine-local runtime artifacts. Never publish secrets or overwrite data during updates.
- Feed and article availability, thumbnails, and publisher blocking can vary. Live Safari checks are required even when Go CI passes.

## Documentation map
[AGENTS.md](AGENTS.md) specifies development rules. [ARCHITECTURE.md](ARCHITECTURE.md) describes implemented internals. [CURRENT_WORK.md](CURRENT_WORK.md) records the active candidate and acceptance status. Historical development is available through Git history rather than standalone version-note files.

## Admin diagnostics (candidate)

Authenticated administrators can open **Diagnostics** after reading one feed or All Subscriptions. The panel shows the latest completed browser read: fetched feed entries, accepted articles, filtered entries, entries beyond the processing cap, feed errors and actual thumbnail load outcomes (loaded, failed, pending or not requested). Thumbnail counts update as image events occur. These browser-local observations do not persist to Server Sync and do not run extra publisher requests. The existing **Filtered / unprocessed articles** section groups rejected items by classification or exclusion reason; entries beyond the cap remain separately labeled as unprocessed.

Timing diagnostics show browser-observed refresh and batch durations, browser render duration, sum of per-feed fetch/wait and classification durations, and mean thumbnail completion-event latency. Summed server per-feed durations are not additive wall-clock stages. These timings do not yet report cache hits or specific thumbnail failure stages.

### Optional RSS metadata-only mode

In **Reading Preferences**, enable **Use RSS metadata only** to bypass ScraperGo article classification and thumbnail extraction while constructing the headline briefing. NewsGo uses RSS/Atom/JSON feed titles, dates, links and HTTPS image metadata when provided; article-age and feed processing limits continue to apply. RSS-only results can include homepages or other non-articles, and some publisher images may be missing or blocked by hotlink protection. Opening a headline still uses ScraperGo. This new mode defaults off and is saved in this browser only, not in Server Sync.


Reading Preferences also offer browser-local Maximum Lead Stories (1–20, default 3) and Lead Stories only in All Subscriptions (default off). The latter skips Latest Articles rendering, not feed fetching or classification. These coexist with RSS metadata only and admin Diagnostics.

## Administrator article URL export (diagnostic)

After a completed **Read All Articles** refresh, the **Download Article URLs** button appears in the briefing **only for a verified Homepage administrator signed in through NewsGo's Homepage handoff**. It is hidden for ordinary users, signed-out browsers and expired sessions, and its click handler rechecks the admin role before exporting. The control uses the normal NewsGo action styling. On localhost without a signed-in Homepage session, the export action remains hidden. This is a client-side export of already visible article URLs, not a server-side privileged download API.

Click to save `newsgo-scrapergo-urls.json` (which may contain publisher URL query strings; keep it private). With ScraperGo running, run `go run ./cmd/scrapertrace -input ~/Downloads/newsgo-scrapergo-urls.json` from `~/projects/newsgo`. For each URL the runner makes sequential local `POST /api/classify` and `POST /api/read` requests and prints its article index, `classification`, `classify_http` (ScraperGo classification endpoint status), `read_http` (ScraperGo reader endpoint status), `read_kind` and `publisher_http` (the publisher's separate HTTP status) when provided. It flags cases classified `article` but failing the subsequent read. ScraperGo's classifier itself may invoke its reader; this workflow deliberately makes both calls to compare their answers, so requests can take longer than the original replay. Classification caches and changed upstream publisher responses may explain some disagreements. Outputs never include article URLs, query strings, raw error messages, or fetched content. The final successful/failed count refers to `/api/read` responses; a nonzero failure count exits with status 1. With ScraperGo recovery tracing enabled, run `go run ./cmd/recoveryreport` from `~/projects/scrapergo` for cumulative extraction results. Exporting does not automatically submit article requests.
