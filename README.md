# NewsGo — RSS News Aggregator

**Current version:** v0.2.23 (release candidate; staged locally, not yet committed or pushed). **Last pushed baseline:** v0.2.20.
**Repository:** https://github.com/richard-mauri/newsgo
**Service:** https://newsgo.themauris.org
**Local directory:** `~/projects/newsgo` · **Port:** `8084`

## Purpose and architecture

NewsGo is a self-hosted Go RSS news aggregator with a browser dashboard. It aggregates shared server-side cached publications and individually requested RSS subscriptions, presenting article thumbnails, headlines and links. The companion **ScraperGo** application displays full article content; it is a separate project and release lifecycle.

The application runs on an Intel MacBook Air, with public access through a Cloudflare tunnel. The project is organized around `cmd/newsgo`, `internal/dashboard` (HTML and JavaScript presentation, subscription UX), and `internal/feeds` (feed processing). Diagnostics are available under `cmd/feedcheck` and `cmd/discoverycheck`.

## Build and local development

```sh
cd ~/projects/newsgo
go test ./...
./devbuild.sh
./newsgo.sh
```

`go.mod` requests Go **1.27.1**. `devbuild.sh` creates macOS Intel and Apple Silicon builds and checksums; see `./devbuild.sh --help` for release options. Avoid changing `newsgo.sh`, Cloudflare configuration, or unrelated files as part of ordinary UI patches. Local UI behavior must also be checked in Safari after a rebuild and reload.

## Current behavior and implementation decisions

### Unified subscriptions (v0.2.22)

There is **one authoritative browser subscription registry** under Safari `localStorage` key `newsgo.subscriptions.v2`. Built-in publications are suggestions under Find Subscriptions, not a distinct category of subscribed content. Existing records from `newsgo.personalFeeds.v1` and `newsgo.includedSubscriptions.v1` are migrated on first use; those earlier keys remain as backups. Once v2 is present it is authoritative: do not re-import stale unsubscribed entries on subsequent loads. Export uses v2; import validates before replacing the registry. The browser registry is scoped to the site origin. Do **not** redesign this already verified system.

Shared subscriptions are supplied through existing server-side cached feeds; individually discovered RSS subscriptions can require separate browser-initiated fetches. The personal-feed backend accepts batches of at most eight. Subscription totals and individually fetched feed totals therefore need not match.

### Articles and navigation (v0.2.23)

The landing page shows a combined article view. Selecting a single publication scopes the articles to that publication. **All Subscriptions** in the subscription pane and **← All Subscriptions** near an individual publication heading restore the combined view without adding a subscription or requiring a feed reload. The subscribed-publications count excludes this navigation option.

Both shared and individually fetched feeds use **The Lead Stories** (up to three selected articles) followed by **Latest Headlines**. On a publication-specific view, leading articles must be selected from *that publication*, not merely hidden from a global first-three set. Decorative `01`, `02`, `03` lead-story numbering was removed. The personal RSS path must not revert to a flat list.

### Refresh status (v0.2.23)

When individually fetched RSS feeds are refreshed, show a visible animated busy indicator and an accurate status; distinguish individually fetched feeds from already loaded shared publications. Count only real completed fetch operations, not projected or simulated progress. Indicate successes, failures and partial failures, retain previously loaded headlines on refresh errors when possible, prevent overlapping responses from overwriting newer ones, support an accessible live status and reduced-motion users. **Do not** display a bare `0 of 4` that looks inconsistent with six subscriptions.

### Thumbnail recovery (v0.2.21)

Shared built-in feeds, including Guardian Business, use the same article-level Open Graph / Twitter image fallback as individually fetched feeds. Feed-supplied images are preferred. Image metadata lookup is bounded to eight recent image-less articles per built-in publisher per 15-minute refresh, up to eight concurrent lookups under a twelve-second shared deadline. Positive and negative results are cached four hours. The approach is generic, not Guardian-specific, and some publishers will still lack thumbnails.

## Development history

- **v0.2.20** — Last known committed base at the October 19 handoff; subscriptions, discovery and thumbnails.
- **v0.2.21** — Generic thumbnail fallback for shared publications; user confirmed Guardian Business improvement.
- **v0.2.22** — Unified subscription registry with legacy migration and export/import; Safari persistence verified across restarts.
- **v0.2.23 (pending commit)** — Busy spinner and meaningful refresh progress; publication-specific lead selection; unified Lead Stories / Latest Headlines for personal RSS feeds; removal of decorative story numbers; visible return to All Subscriptions. User visually confirmed navigation and refresh indicator. Full final test pass and commit/tag still require local confirmation.

## Verification and release procedure

From the final local working tree run `go test ./...`, inspect `git status`, review staged files using `git diff --cached --stat` and `git diff --cached --check`, and exercise shared feeds, personal feeds, navigation, thumbnails and refresh in Safari. Do not infer a test pass from a test file or a syntax-only check. Once verified, commit and push before tagging a release. The README's baseline status must be updated in the same source regeneration or release patch.

The usual delivery is a **minimal ZIP patch** with paths relative to `~/projects/newsgo`. Every regeneration that modifies source must also include a freshly updated `README.md`. Preserve existing user modifications; avoid distributing outdated copies of concurrently edited files. Do not overwrite scripts, tunnel settings or ScraperGo source without explicit request.

## Next-development handoff

Treat the user’s local files as authoritative whenever `git status` shows uncommitted development beyond GitHub. At this handoff, the user intends to commit v0.2.23; do not assume the commit, tag, or final `go test ./...` completed until confirmed. Subsequent development should begin from the actual committed or uploaded v0.2.23 source, normally advancing to v0.2.24. Keep regressions out of subscription storage, thumbnails, article grouping and All Subscriptions navigation.

---

## Earlier README history (preserved)

The original README text below documents the initial FeedCheck milestone and its early diagnostics; its installation directions are **historical**, not the current NewsGo installation procedure.

# NewsGo v0.1.0 — FeedCheck

A standalone Go feed preflight utility. No third-party modules, Docker, or database are required. Supports RSS 2.0, Atom, RSS 1.0 (RDF), and JSON Feed. The initial ten feed URLs are candidates, **not yet verified**. Run the checker from your Mac to measure actual availability.

## Install into your existing repository

You have already created `~/projects/newsgo` and initialized `go.mod` with `github.com/richard-mauri/newsgo`. Extract the patch contents directly into that folder. The ZIP deliberately excludes `go.mod` and `.git`.

```sh
cd ~/projects/newsgo
unzip -o ~/Downloads/newsgo-v0.1.0-feedcheck.zip
go test ./...
go run ./cmd/feedcheck
```

The program prints feed status and metadata coverage, and writes `report.json` with per-feed HTTP status, format, item count, timestamps, and errors. `report.json` is local output, not part of the distributed ZIP. Use `-feeds custom.json`, `-out custom-report.json`, `-timeout 20s`, and `-parallel 2` as needed.

Statuses: `healthy` (at least one headline and link, newest known story <=72 hours old), `stale`, `undated`, `empty`, `poor_metadata`, `invalid_feed`, `http_error`, `error`. A successful feed fetch says nothing about ScraperGo article extraction: publisher access restrictions must be separately tested later.

## Scope

FeedCheck is the first milestone, not a running web server. It does not crawl arbitrary URLs, bypass paywalls, or modify ScraperGo. Next milestones are cached aggregation and an HTML dashboard opening articles via the existing ScraperGo URL interface.

## v0.1.1 — RSS link recovery and optional ScraperGo sampling

Some RSS publishers supply article links as a permalink `<guid>` instead of a normal `<link>`, or use an `href` attribute. FeedCheck now recognizes both formats generically, without publisher-specific logic. A GUID explicitly marked `isPermaLink="false"` is not interpreted as an article URL. FeedCheck continues to report missing links when no legitimate HTTP(S) URL exists.

Optionally ask your **running local** ScraperGo service to extract at most two articles per feed:

```sh
go run ./cmd/feedcheck -scrapergo http://127.0.0.1:8083 -samples 2
```

The tool makes POST requests to ScraperGo's existing `/api/read` endpoint and reports `readable`, `restricted`, or `other` for each sample in `report.json`, with a terminal summary. Publisher denial and subscription/registration notices are *restricted* rather than failed extraction. `other` includes network and format errors. This is diagnostic information, not a source eligibility decision. Sample tests make additional network requests to publishers via ScraperGo; they are **off by default**. Use `-samples 1` for minimal traffic. No settings, database, tunnel, or service code changes are required.

**Limitations:** a successful response must contain a nonempty article title to count as readable; this does not guarantee a complete clean extraction. Sampling is not a ranking score or comprehensive site certification. FeedCheck refetches a feed once to obtain sample URLs when `-scrapergo` is enabled. Avoid pointing the option to an untrusted service because it receives article URLs.
