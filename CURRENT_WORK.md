# Current work — NewsGo

## Active draft — subscription failure details (Issue #55, November 8, 2026)

- **Issue:** [#55](https://github.com/richard-mauri/newsgo/issues/55); branch `fix/55-subscription-failure-details`; draft [PR #73](https://github.com/richard-mauri/newsgo/pull/73) against authoritative merged `main`.
- **Changes:** The existing Failure dialog shows the actual attempted feed endpoint and the validation error, with sensitive URL query values redacted and Retry preserved. **Find Subscriptions now prevalidates candidate feeds**, using the SSRF-hardened fetch path and existing RSS/Atom/RDF/JSON parser. The directory response labels candidates **Available** (readable entries), **Unavailable** (e.g., HTTP 400, 404, 410, unparseable/empty responses; Subscribe disabled), or **Not verified** (timeouts, 429, 503 and other uncertain results; Try Subscribe remains possible). The check is bounded to four simultaneous fetches, four seconds per request, nine seconds per search, twelve candidates, and two minutes of result caching. No automatic substitution of feed URLs or permanent blacklist.

- **Scope preserved:** The existing Failure button, Close/Retry controls, independent category subscriptions, SSRF/feed validation, feed directory, per-user subscriptions, exclusions, preferences, credentials, server settings, cache, build counter and Cloudflare are unchanged. Stale error state is cleared on retry; failures are not automatically retried.
- **Evidence:** Screenshot on November 8 shows NBC News feeds with per-card Failure and status `Feed not added: HTTP 400`; Issue #55 also has an observed HTTP 503. Exact upstream cause not yet established; this PR improves diagnosis, not feed availability.
- **Tests:** Dialog/URL redaction and retry contracts, feed health classification (RSS success, HTTP 400/404/429/503, invalid markup, empty feed), concurrency bound, and existing catalog/exclusion behavior. Exact PR-head Go CI status to be recorded after run. Mac/Safari live publisher and UI acceptance remain outstanding.
- **Owner Mac acceptance:** `cd ~/projects/newsgo`, `git fetch origin`, `git checkout fix/55-subscription-failure-details`, `git pull --ff-only origin fix/55-subscription-failure-details`, `go test -count=1 ./...`, `./devbuild.sh`, `./newsgo.sh`. In Safari search NBC and verify HTTP 400 entries are explicitly Unavailable (not selectable), valid feeds show Subscribe, uncertain 503s show Try Subscribe, and an attempted temporary failure gives a Failure dialog with the attempted endpoint, HTTP status and Retry. Confirm existing subscriptions/exclusions remain unchanged.
- **Merge boundary:** PR #73 remains draft until explicit owner approval; do not merge, tag, release or change `main`.

## Other work (independent)

- NewsGo [PR #69](https://github.com/richard-mauri/newsgo/pull/69) for Issue #57 (Homepage session renewal/countdown) is a separate draft requiring coordinated [Homepage PR #21](https://github.com/richard-mauri/homepage/pull/21), Mac/Safari acceptance and explicit owner promotion approval. This issue does not modify it.
- [PR #72](https://github.com/richard-mauri/newsgo/pull/72) for NBC dates and raw feed responses was merged into `main` on November 8 (merge commit `7ca9f5840fa9ac611a1e241113001ff09e16e8e3`), and Issue #71 was closed. The merged raw feed viewer remains separate from this subscription Failure dialog.
