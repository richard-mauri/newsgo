# Current work — NewsGo

## Active draft — subscription failure details (Issue #55, November 8, 2026)

- **Issue:** [#55](https://github.com/richard-mauri/newsgo/issues/55); branch `fix/55-subscription-failure-details`; draft [PR #73](https://github.com/richard-mauri/newsgo/pull/73) against authoritative merged `main`.
- **Change:** Existing Find Subscriptions **Failure** dialog now displays the *actual attempted feed URL* separately from the validation error (including HTTP 400/503 status, when returned). This URL comes from the candidate feed record, not the publisher homepage. Sensitive query values, username/password and fragment are redacted. Values of a small whitelist of ordinary feed-format flags are preserved for diagnosis. All content uses textContent.
- **Scope preserved:** The existing Failure button, Close/Retry controls, independent category subscriptions, SSRF/feed validation, feed directory, per-user subscriptions, exclusions, preferences, credentials, server settings, cache, build counter and Cloudflare are unchanged. Stale error state is cleared on retry; failures are not automatically retried.
- **Evidence:** Screenshot on November 8 shows NBC News feeds with per-card Failure and status `Feed not added: HTTP 400`; Issue #55 also has an observed HTTP 503. Exact upstream cause not yet established; this PR improves diagnosis, not feed availability.
- **Tests:** `internal/dashboard/subscription_failure_test.go` extends dialog/URL redaction and retry contract. Exact PR-head Go CI status to be recorded after run. Mac/Safari live feed and user acceptance are outstanding.
- **Owner Mac acceptance:** `cd ~/projects/newsgo`, `git fetch origin`, `git checkout fix/55-subscription-failure-details`, `git pull --ff-only origin fix/55-subscription-failure-details`, `go test -count=1 ./...`, `./devbuild.sh`, `./newsgo.sh`. In Safari search NBC, select a failing candidate, click Failure, verify attempted feed URL is not the shared publisher homepage and HTTP status is visible, Retry invokes a new attempt, and existing subscribed feeds remain unchanged.
- **Merge boundary:** PR #73 remains draft until explicit owner approval; do not merge, tag, release or change `main`.

## Other work (independent)

- NewsGo [PR #69](https://github.com/richard-mauri/newsgo/pull/69) for Issue #57 (Homepage session renewal/countdown) is a separate draft requiring coordinated [Homepage PR #21](https://github.com/richard-mauri/homepage/pull/21), Mac/Safari acceptance and explicit owner promotion approval. This issue does not modify it.
- [PR #72](https://github.com/richard-mauri/newsgo/pull/72) for NBC dates and raw feed responses was merged into `main` on November 8 (merge commit `7ca9f5840fa9ac611a1e241113001ff09e16e8e3`), and Issue #71 was closed. The merged raw feed viewer remains separate from this subscription Failure dialog.
