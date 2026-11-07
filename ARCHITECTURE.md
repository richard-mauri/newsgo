# NewsGo — architecture

See [AGENTS.md](AGENTS.md) and [DEVELOPMENT.md](DEVELOPMENT.md) for workflow, [README.md](README.md) for operation, and [CURRENT_WORK.md](CURRENT_WORK.md) for the active draft PR.

## Runtime

- `cmd/newsgo/main.go`: service startup, catalog, server settings, listener, graceful shutdown; default `127.0.0.1:8084`.
- `internal/dashboard/`: UI, personal reading pipeline, Homepage account integration, administrator routes, diagnostics and operational settings.
- `internal/feeds/`: feed validation and fetching, URL safety, parsing and article normalization.
- ScraperGo is independent and provides article classifications, clean reading and thumbnails.

## Reading and classification

NewsGo fetches selected RSS/Atom/JSON subscriptions in batches of at most eight. Each eligible article is checked using ScraperGo `POST /api/classify`; only confirmed `article` URLs appear in either normal or RSS metadata-only reading. Homepages, unreadable/unavailable/unknown outcomes, age-filtered articles and entries beyond feed caps are reported in the filtered/unprocessed audit. RSS metadata-only changes the displayed metadata and thumbnail source, not the classifier gate.

NewsGo's shared classification cache is process-local, bounded and controlled by the server TTL (0–1440 minutes). Its default is 1440 minutes, but saved explicit values, including zero, are preserved. Readers can opt out of cache use in Reading Preferences without clearing shared entries. The per-reader Diagnostics view reports cache hits, misses, bypasses, disabled requests, configured TTL, feed durations, classification durations and render durations. ScraperGo has separate upstream caching not included in these counters.

## Homepage-only authentication and account persistence

Homepage provides **all live reader identity and account persistence**. The browser obtains a tab-scoped `#suite_token` handoff from Homepage and verifies it using `https://home.themauris.org/api/session/profile`. The token is removed from the visible URL and retained in `sessionStorage` while the original Homepage session is active, up to twelve hours. The authenticated reader explicitly loads or saves their account settings through Homepage's revisioned `/api/storage/newsgo` GET/PUT API. Profiles remain isolated by Homepage identity. Unsaved local preferences remain in browser `localStorage`. No separate NewsGo token issuance or profile-sync login exists.

NewsGo administrator endpoints `GET/PUT /api/admin/settings` and `POST /api/admin/classification-diagnostic` require server-side verification of the supplied Homepage bearer against Homepage's **fixed** HTTPS session-profile endpoint on every privileged request. The browser sends it in the dedicated `X-NewsGo-Homepage-Authorization` header. The Go verifier requires a valid bearer, a current expiry and the server-asserted `admin` role, with three-second timeout, no redirects and no persistent secret logging. An ordinary Homepage reader cannot administer NewsGo, regardless of browser UI manipulation. Reading Diagnostics is not privileged and remains browser-local.

The old Server Sync endpoint `/api/sync`, manually approved token authentication, legacy token CLI and token-backed browser flow have been **removed from the draft candidate**. Existing historical `data/approved-tokens.json` and `data/user-preferences.json` files are not deleted or migrated automatically, and are no longer accessed by the application.

## Operational server settings

`internal/dashboard/settings_store.go` implements a dedicated, mutex-protected `SettingsStore` that reads and writes the existing revisioned `data/server-settings.json` atomically with restrictive file permissions. The CLI flag `-settings-dir` selects the directory (default `data`). Classification timeout/concurrency, per-feed cap, feed cache, classification-cache TTL and historical image settings are preserved across migration. Saving server settings updates active settings and clears NewsGo's in-memory classification cache. Existing saved TTL zero is intentional and is not overwritten without an admin change.

## Other diagnostics and safety

The administrator's fresh classification diagnostic is deliberately uncached. The localhost-limited feed-trace endpoint remains available. The client-side Download Article URLs export is restricted by the browser to a live Homepage administrator session; it exports already displayed URLs and does not create a server-side privileged download API.

Publisher URLs and responses are untrusted and subject to existing validation, request limits, timeouts and network safety rules. A browser role label alone is never an authorization decision. `cmd/scrapertrace` independently replays exported URLs through ScraperGo with sequential default operation, optional bounded four-worker concurrency, separate classification/read HTTP status and latency reporting.

## Build and release

`devbuild.sh` produces architecture-specific Mac binaries with a local build label. GitHub `main` remains the authoritative merged baseline. Draft PR #68 is not merged; see [CURRENT_WORK.md](CURRENT_WORK.md) for owner acceptance and exact-head CI.
