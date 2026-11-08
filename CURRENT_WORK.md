# Current work — NewsGo

## Parallel draft — NBC date format (November 8, 2026)

- **Issue:** [#71](https://github.com/richard-mauri/newsgo/issues/71), branch `fix/71-nbc-feed-dates`, draft [PR #72](https://github.com/richard-mauri/newsgo/pull/72) targeting `main`.
- **Diagnostics scope:** Replaced the redundant Article Date Diagnostics with a format-agnostic **Raw Feed Inspector**. It exposes the publisher's actual feed response bytes (XML, JSON, etc.) from the existing or cached fetch, a direct feed-URL link, plain source text, and Copy/Download controls. Download preserves the recorded bytes; the source view is text-only. No extra HTTP requests; latest read remains browser-local.
- **Scope:** Recognize NBC Bay Area feed dates such as `Fri, Nov 06 2026 03:05:45 PM` in the existing Go RSS parser; retain the raw date fallback for formats that still cannot be parsed, and add regression tests. No publisher dates are rewritten.
- **Evidence and outstanding checks:** NBC parser, raw-byte transport and raw-feed tests added. GitHub Actions on the final branch head is being checked. Mac/Safari acceptance remains pending; NBC's 49 identical dates are not changed or guessed.
- **Acceptance:** From `~/projects/newsgo`, fetch and check out `fix/71-nbc-feed-dates`, run `go test -count=1 ./...`, `./devbuild.sh`, and `./newsgo.sh`. In Safari, read NBC Bay Area, open Diagnostics → Raw Feed Inspector, inspect the unchanged `pubDate` strings in XML, verify the source feed URL, Copy Raw Feed, and Download Raw Feed byte fidelity, plus normal date display and filters.
- **Merge boundary:** PR #72 remains draft. Never merge, tag, release, or modify `main` without explicit owner approval; this work does not supersede the separate session-renewal PR #69.


## Verified state (November 8, 2026)

- **Merged baseline:** GitHub `main` at `3c054600a3200fab059fda4fef721295735b7d54`, after owner-merged PR #68 for Issue #67.
- **Active task:** [Issue #57](https://github.com/richard-mauri/newsgo/issues/57), branch `fix/57-session-countdown-renewal`, draft [PR #69](https://github.com/richard-mauri/newsgo/pull/69) targeting `main`.
- **Asset delivery:** HTML now uses a content-fingerprinted `/personal.js` URL so Safari and intermediary caches cannot reuse a script from a prior build with the same semantic version.
- **UX fix (November 8):** Renew now highlights on hover and focus, keeps the button caption stable and places only failures beside the header button rather than inside the hidden Find Subscriptions view. The reset countdown is the sole success indicator, without flashing intermediate text or changing cursor/hover state: duplicate clicks are guarded in JavaScript without disabling the button or setting transient ARIA busy attributes. The remaining browser-level cursor behavior must be verified by the owner in Safari. A missing Homepage endpoint reports HTTP 404 and instructs installing Homepage PR #21.
- **Scope:** The header displays the authenticated Homepage reader and a per-second HH:MM:SS remaining counter based on server-asserted `expiresAt`. It shows expiration, clears the local handoff on expiry, and offers one-click Renew through Homepage's new authenticated `POST /api/handoff/extend-session` endpoint. Renewal atomically extends the still-active parent login twelve hours, rotates the NewsGo bearer, and updates local expiration without navigation or approved-token reentry. Requires [Homepage draft PR #21](https://github.com/richard-mauri/homepage/pull/21), Issue #20, installed first. Expired or revoked sessions still require Homepage sign-in. No Homepage or other app changes.
- **Security/data:** Homepage remains the sole identity and storage authority. No new token issuance, no server-side entitlement changes, no credentials logged or persisted beyond the existing tab handoff. Existing per-user account data, local preferences, cache, server settings, launchers, Cloudflare and build labels remain untouched.
- **Validation:** Added `session_countdown_test.go`; GitHub Actions exact-head validation is tracked in PR #69; Mac/Safari acceptance remains owner-controlled. Live Mac/Safari acceptance remains required. This environment cannot verify Safari, live Homepage login, countdown over real elapsed time or successful reauthentication.
- **Owner acceptance:** After building and launching Homepage PR #21 and this NewsGo candidate, launch NewsGo while signed in; confirm HH:MM:SS decreases and Renew resets the timer near twelve hours **without navigation or entering an approved token**. Return to Homepage and confirm its recovered login recognizes the extension. Verify expired/revoked sessions fail, old app bearer cannot be replayed, expiration disables account/admin controls, and browser reader settings remain intact.
- **Merge boundary:** Keep the PR draft. Do not merge, release, tag or modify `main` without explicit owner approval.

## Mac acceptance commands

```bash
cd ~/projects/newsgo
git fetch origin
git checkout fix/57-session-countdown-renewal
git pull --ff-only origin fix/57-session-countdown-renewal
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

After CI and owner Safari checks, update PR evidence and request explicit owner approval before merge.
