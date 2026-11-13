# Current work — NewsGo

- **Same-tab reader navigation:** NewsGo headline links now open ScraperGo in the current tab (no named `newsgo-scrapergo` target), preserving the existing URL and Homepage handoff. External original-publisher links remain separate. Await exact-head CI and owner Safari acceptance; no merge.

- **Heartbeat revocation follow-up:** NewsGo's normal activity ping now returns HTTP 401 on a supplied revoked/expired handoff (never reclassifying it as anonymous); the browser notifies its existing account controller, which revalidates against Homepage and opens the existing sign-in-required dialog. An unavailable Homepage returns 503. CI #606 passed before this note; Mac Safari acceptance pending.

## Issue #93 — Non-dismissible suite sign-in prompt

- **Baseline:** GitHub `main`, including merged NewsGo PR #92 (Homepage named/guest identity and activity).
- **Branch:** `fix/93-block-signin-dismiss` from `main`; focused draft PR pending. Homepage central authorization remains unchanged.
- **Owner-reproduced UX:** After Homepage admin revokes a one-hour guest, NewsGo recognizes sign-out but the sign-in-required dialog offered **Close**, leaving the app visible as an unidentified reader.
- **Change:** Remove the sign-in dialog's Close button and prevent native Escape/cancel dismissal. Keep **Go to Homepage** as explicit, user-triggered navigation in a separate tab; a successful Homepage profile verification still closes the dialog programmatically. No automatic redirect.
- **Boundaries:** No changes to NewsGo feeds, ScraperGo APIs, persistence, browser-only settings, private tokens, activity heartbeat, or user data. Homepage Issue #38 separately tracks a later protected-API authorization audit.
- **Regression:** Source-contract test in `internal/dashboard/session_countdown_test.go` covers modal blocking and verified dismissal. Exact-head Go CI and owner's Apple Silicon Safari acceptance are required; neither implies the other.
- **Owner acceptance:** Verify invalid/revoked session shows only **Go to Homepage**; click and Escape cannot dismiss; valid Homepage sign-in and handoff permit normal reader use. Do not merge without explicit owner approval.

Mac commands after draft PR head is confirmed:

```bash
cd ~/projects/newsgo
git fetch origin
git checkout fix/93-block-signin-dismiss
git pull --ff-only origin fix/93-block-signin-dismiss
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```
