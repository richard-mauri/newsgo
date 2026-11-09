# Current work — NewsGo

## Active draft — Homepage-only session management (Issue #79, November 9, 2026)

- **Merged baseline:** GitHub `main` at `4971142157d0c75b31406355f7ee967baf34d10c` (PR #78 merged). Previous Issue #57 / PR #69 is already merged.
- **Current scope:** [Issue #79](https://github.com/richard-mauri/newsgo/issues/79), `refactor/79-homepage-session-management`, draft [PR #80](https://github.com/richard-mauri/newsgo/pull/80). This removes NewsGo's redundant countdown, Renew button, and legacy cross-origin session extension calls, and fixes loss of the tab bearer after Safari/browser closure by delegating recovery to Homepage Issue #24, branch `fix/24-newsgo-session-recovery`. Homepage's merged Issue #22 is the single session authority.
- **Preserved:** Homepage transfer bearer and `/api/session/profile` validation; explicit browser-only launch; admin-only server-side controls; per-user account storage, subscriptions, exclusions and local preferences. On a direct HTTPS NewsGo visit without a bearer, the app navigates once via Homepage `/launch/newsgo`: a valid parent gets a fresh bound handoff; a signed-out/expired parent returns `#suite_mode=local`, recognized for the current visit to prevent a redirect loop; later direct visits recheck Homepage. At cached handoff expiry NewsGo revalidates with Homepage before expiring a bearer that Homepage may already have extended.
- **Other apps:** MauriStocks already consumes Homepage handoffs without a Renew button. SailingGo and ScraperGo are browser-only for account/session management. No edits to any of these repositories.
- **Validation:** Prior Go CI [run 42687686183](https://github.com/richard-mauri/newsgo/actions/runs/42687686183) passed at head `5b45f3479c94e0618ad775625397169519be066c`. New cross-app recovery changes require fresh CI on final head. Mac/Safari and live cross-origin behavior remain owner-only and unverified.
- **Owner acceptance:** After approving both Homepage #24 and NewsGo #80 candidate restarts, sign in to Homepage, launch NewsGo, then close Safari, restart NewsGo and open NewsGo directly. Confirm Homepage silently recovers the same parent login (not a new 12-hour session). Verify identity, account storage/admin role, and no local Renew/timer; navigate to Homepage Manage Tokens for Renew and Disconnect. Confirm invalid/expired Homepage login returns NewsGo in browser-only mode with intact local subscriptions, without loops.
- **Merge boundary:** Draft PR only; no merges, releases, version changes, deployment, process restart or Cloudflare changes without owner authorization.

## Mac commands

```bash
cd ~/projects/newsgo
git fetch origin
git checkout refactor/79-homepage-session-management
git pull --ff-only origin refactor/79-homepage-session-management
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```
