# NewsGo — implemented architecture

This document describes **how the current application works**, not ongoing task status or release history. For usage see [README.md](README.md); for agent process see [AGENTS.md](AGENTS.md); for the active candidate see [CURRENT_WORK.md](CURRENT_WORK.md).

## Runtime and components
- `cmd/newsgo/main.go`: CLI flags, source catalog loading, web-service initialization, version output, HTTP listen and graceful shutdown. NewsGo runs at `127.0.0.1:8084` under the standard launcher and uses a separate Cloudflare tunnel.
- `internal/dashboard/`: HTML template and embedded JavaScript, HTTP routing, personal feed orchestration, classification, administrator settings, diagnostics, Server Sync and token management.
- `internal/feeds/`: feed source validation, downloading, parsing, article normalization and discovery helpers.
- ScraperGo is a **separate service** for article classification, reading and thumbnails; NewsGo calls its HTTP API rather than performing publisher article extraction locally.

## Reader data flow
1. The browser holds a subscription registry, publisher exclusions, article-destination exclusions, reading age and selected publication. User selection initiates a batch request to `POST /api/personal`; a request includes up to eight selected feeds.
2. The service validates source URLs and fetches feeds through a bounded public-network-only client, with guarded redirects and feed caching. Feed results contain article metadata and original publisher links.
3. Entries beyond the configurable processing cap are counted as unprocessed. Dated entries older than the user's allowed age are marked `too-old`.
4. For eligible entries, NewsGo calls ScraperGo `POST /api/classify` with a normalized article URL. Browser-only fragments such as `#new_tab` are removed **only from classification input**; the original feed URL remains available for navigation. Unsafe URLs remain rejected.
5. Only `article` results pass the strict article gate. `unknown`, `homepage`, `unreadable`, `unavailable`, and age-filtered items are not silently promoted into headlines. The response includes filtered reasons and unprocessed counts.
6. The browser renders headline cards, ScraperGo thumbnail URLs, source selection and a collapsible filtered-entry audit. Publisher and article-domain exclusions are additionally reflected in display.

## Persistence and synchronization
- Browser `localStorage` maintains subscriptions, exclusions, selected publication and other local reading preferences; browser-origin differences matter.
- `data/` holds approved-token records and revisioned per-user Server Sync JSON profiles. `GET/PUT /api/sync` requires a valid approved bearer token and checks profile revisions to prevent silent overwrites. Connecting can load server preferences or merge when profiles differ.
- Administrator tokens have an explicit server-held `admin` role. `GET/PUT /api/admin/settings` and `POST /api/admin/classification-diagnostic` enforce the role on the server.
- **Exit Admin Mode** calls the existing browser-side Server Sync disconnect path: clears active and remembered token from that browser, hides admin controls, leaves local preferences available and does **not** revoke the server token.
- Server settings control max articles per feed, classification timeout/concurrency, and successful feed cache duration. Historic image-enrichment fields remain preserved for compatibility where applicable.

## Diagnostics and security
- `POST /api/admin/classification-diagnostic` runs a fresh on-demand classification of a selected filtered article and returns bounded classification/failure information, upstream HTTP status when applicable, elapsed milliseconds, and configured timeout. It does not change persisted article eligibility or subscription data.
- `POST /api/diagnostic/feed-trace` is localhost-restricted and traces a requested URL through a supplied feed's matching/processing decisions.
- `/api/discover` and `/api/search-feeds` support finding feeds; `/health` returns process health.
- URL validation, bounded response handling, safe network routing, per-request timeouts and token checks are security boundaries. Admin-only UI visibility is never a substitute for server authorization.
- Recognized article classification outcomes are cached for limited periods; transient `unknown` outcomes are not persisted by NewsGo's classification cache. Actual browser and publisher behavior is verified separately from Go tests.

## Build identification
The semantic version is obtained from the dashboard's header version and used by the CLI. `devbuild.sh` injects the Mac-local sequential build number into both binaries through Go linker flags. Build numbering is a debugging identifier, not an API or data-format migration.

Refer to code and tests for precise validation bounds and to [CURRENT_WORK.md](CURRENT_WORK.md) for known open bugs and acceptance results.
