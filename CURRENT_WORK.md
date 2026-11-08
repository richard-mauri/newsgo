# Current work — NewsGo

## Verified state (November 8, 2026)

- **Merged baseline:** GitHub `main` at `3c054600a3200fab059fda4fef721295735b7d54`, after owner-merged PR #68 for Issue #67.
- **Active task:** [Issue #57](https://github.com/richard-mauri/newsgo/issues/57), branch `fix/57-session-countdown-renewal`, draft [PR #69](https://github.com/richard-mauri/newsgo/pull/69) targeting `main`.
- **Asset delivery:** HTML now uses a content-fingerprinted `/personal.js` URL so Safari and intermediary caches cannot reuse a script from a prior build with the same semantic version.
- **Scope:** The header displays the authenticated Homepage reader and a per-second HH:MM:SS remaining counter based on server-asserted `expiresAt`. It shows expiration, clears the local handoff on expiry, and offers a Renew link returning to Homepage for a fresh sign-in and relaunch. A fixed twelve-hour Homepage parent login cannot be extended by bearer rotation. No Homepage or other app changes.
- **Security/data:** Homepage remains the sole identity and storage authority. No new token issuance, no server-side entitlement changes, no credentials logged or persisted beyond the existing tab handoff. Existing per-user account data, local preferences, cache, server settings, launchers, Cloudflare and build labels remain untouched.
- **Validation:** Added `session_countdown_test.go`; exact-head CI must be inspected before review. This environment cannot verify Safari, live Homepage login, countdown over real elapsed time or successful reauthentication.
- **Owner acceptance:** In Safari, launch NewsGo while signed in to Homepage; confirm HH:MM:SS decreases, Renew navigates to Homepage and after a **fresh sign-in** relaunching NewsGo shows a new session. Confirm a still-active Homepage login does not get a falsely extended twelve-hour lifetime; verify expiration disables account/admin controls and ordinary reader settings remain intact.
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
