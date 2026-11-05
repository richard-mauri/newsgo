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

## Homepage session-bound bearer handoff
NewsGo reads `#suite_token=...` from Homepage launcher links, immediately removes it from the visible URL, retains the bearer in tab-scoped `sessionStorage`, and verifies it with Homepage `/api/session/profile`. Homepage's revisioned `/api/storage/newsgo` GET/PUT uses an Authorization bearer header. There is no cookie fallback: an absent or invalid token leaves NewsGo in browser-only mode. The prior `/api/homepage/session` cookie proxy and its tests have been removed. No token is passed to NewsGo server endpoints, and normal reader subscriptions remain local until an explicit account load/save. An expired bearer requires navigating again through Homepage.

## Session status
The account status near the top shows only the authenticated identity (`Signed in as …`) or `Not signed in`, without a visible expiration countdown. The client still enforces the original twelve-hour Homepage session expiry using Homepage's `expiresAt` profile field. There is no independent app token renewal. Once the Homepage session expires, the user returns to Homepage and signs in again to restore account access.

## Return to Homepage navigation
NewsGo's browser UI exposes a visible native anchor `<a href="https://home.themauris.org/">← Back to Homepage</a>` with no `target` and no JavaScript history manipulation. It navigates in the current tab and can be used regardless of account sign-in. The link adds no new server endpoints, exposes no credentials, and does not modify app-local state. This is coordinated with Homepage Issue #11 / draft PR #12, where card clicks now navigate in the same tab.

The suite-wide return button uses the same `.suite-home-link` presentation in all four apps: 38px minimum height, 8px/12px padding, 10px radius, muted border at rest, accent border and darker fill on hover, and a visible keyboard focus outline. Place one link at the upper-left above the primary content/heading.

The shared suite navigation style is a quiet 12px text link below each page's primary heading, with subtle hover background and underline and keyboard focus outline. NewsGo has no duplicate Homepage link or Renew button.

## Twelve-hour Homepage sign-in handoff
NewsGo receives a one-time transfer bearer from Homepage for its current twelve-hour login session. The bearer is valid only while that Homepage parent session remains active and is rejected by Homepage after sign-out or session expiry. NewsGo does not offer independent token renewal. A Homepage `#suite_mode=local` launch explicitly discards any older cached bearer, preserving browser-only reading preferences. A launch with `#suite_token=...` replaces the cached bearer, and ordinary reloads may reuse it until server validation rejects it. Account save/load continues to require live server validation. This change is paired with Homepage draft PR #12; deploy the coordinated changes together only after owner Safari acceptance.
