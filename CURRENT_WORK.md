# Current work — NewsGo

## Active draft — Homepage-only session management (Issue #79, November 9, 2026)

- **Merged baseline:** GitHub `main` at `4971142157d0c75b31406355f7ee967baf34d10c` (PR #78 merged). Previous Issue #57 / PR #69 is already merged.
- **Current scope:** [Issue #79](https://github.com/richard-mauri/newsgo/issues/79), `refactor/79-homepage-session-management`, draft [PR #80](https://github.com/richard-mauri/newsgo/pull/80). This removes NewsGo's redundant countdown, Renew button, and legacy cross-origin session extension calls. Homepage's merged Issue #22 is the single session authority.
- **Preserved:** Homepage transfer bearer and `/api/session/profile` validation; explicit browser-only launch; admin-only server-side controls; per-user account storage, subscriptions, exclusions and local preferences. At its cached handoff expiry NewsGo revalidates with Homepage before expiring a bearer that Homepage may already have extended.
- **Other apps:** MauriStocks already consumes Homepage handoffs without a Renew button. SailingGo and ScraperGo are browser-only for account/session management. No edits to any of these repositories.
- **Validation:** Go CI [run 42687151371](https://github.com/richard-mauri/newsgo/actions/runs/42687151371) passed on code/test head `fa6c805b84ffd1cca2fa968036032b0ce1a0e22f` (Go tests and build). This documentation-only update requires a fresh exact-head CI check. Mac/Safari and live cross-origin behavior remain owner-only and unverified.
- **Owner acceptance:** Sign in to Homepage, launch NewsGo, confirm identity and no Renew/timer; navigate to Homepage Manage Tokens and renew, return to NewsGo and confirm account storage/admin role remain valid after the original child deadline; test Disconnect/revocation, expired credentials and signed-out browser-only mode without disturbing preferences.
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
