# Current work — NewsGo

**Merged main verified:** `cc86480caafbf3299b1064c380e0d7e5bd17ff3c`. **Active scope:** [Issue #114](https://github.com/richard-mauri/newsgo/issues/114), draft [PR #115](https://github.com/richard-mauri/newsgo/pull/115), branch `fix/114-authenticated-thumbnails` based on `main`. Coordinated ScraperGo [Issue #147](https://github.com/richard-mauri/scrapergo/issues/147), draft [PR #148](https://github.com/richard-mauri/scrapergo/pull/148), branch `fix/147-newsgo-thumbnail-auth` based on ScraperGo merged `main` `e22b2e7427a6884074d10e72b24b47fe544aded8`. Both historical authentication PRs (#111, #146) are superseded by merged main and are not active working branches.

## Problem and resolution

Owner Safari observed 196 NewsGo articles, 190 normal thumbnail failures (0 loads) and an unauthorized ScraperGo `GET /api/thumbnail` (HTTP 401) because the browser's direct `<img>` request lacked bearer authentication. RSS metadata-only publisher images worked. Draft PR #115 routes normal thumbnail loads through a NewsGo same-origin, Homepage-authenticated `GET /api/thumbnail` proxy, which bounds upstream URL, method, width, image type/size, redirects, timeout and concurrent requests. Browser uses `authorizedNewsgoFetch` and a temporary blob URL; RSS mode is unchanged.

Owner subsequently observed a valid Homepage personal session briefly displaying articles then disappearing after the new proxy returned HTTP 503 for missing service configuration. Fixed the proxy's own unavailable/configuration failure to HTTP **502**, retaining the outer Homepage verification gate's 503 behavior and CI-only browser regression. Homepage personal/Guest auth, per-user settings and explicit Save Changes remain unchanged.

The owner confirmed the original private `SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN` exists in `/Users/richardmauri/projects/homepage/private/service-auth.env` and instructed reuse of that **one existing 64-lowercase-hex machine credential for all allowlisted NewsGo→ScraperGo APIs**, instead of creating a per-endpoint thumbnail token. This draft now uses that existing variable to forward `GET /api/thumbnail`; ScraperGo PR #148 permits the same credential on `POST /api/classify` and `GET /api/thumbnail` only. Other endpoints remain strictly protected; no token value or new variable is introduced. Do not add or source `SUITE_NEWSGO_SCRAPERGO_THUMBNAIL_TOKEN`.

## Security, tests, documentation and operations

The outer NewsGo Homepage gate requires a verified user/Guest bearer for every app API request. The machine token stays server-side and is forwarded only to the fixed ScraperGo HTTPS/explicit-loopback endpoint, never in query params or HTML. Existing classification authorization and ScraperGo SSRF/media policies remain intact. Added Go and CI-only browser tests for proxy image success/errors, missing credential 502 without sign-out, real 503 fail-closed, shared token scope and other-route denial.

Reviewed current Homepage `SUITE_DEVELOPMENT_POLICY.md` v1.1.0 from `main` plus NewsGo `AGENTS.md`, `README.md`, `ARCHITECTURE.md`, `DEVELOPMENT.md`, `CURRENT_WORK.md`. Updated README/ARCHITECTURE/DEVELOPMENT and this handoff; AGENTS and Homepage policy unchanged. **Exact-head CI must be verified on the latest commit of each PR before handoff.** Previous CI #741 passed on the earlier NewsGo head and is not proof of this change. Owner's Mac local working tree, private credential value, local settings/browser state and tunnel are inaccessible; no private data, credentials, launchers, account settings, subscriptions, counters, processes or tunnels modified. No production restart/deployment/merge authorized. Owner Mac/Safari/live cross-service acceptance remains unverified.

## Owner Mac handoff (only after final CI)

```bash
cd ~/projects/newsgo
git status --short --branch
# Stop if there are uncommitted changes.
git fetch origin
git checkout fix/114-authenticated-thumbnails
git pull --ff-only origin fix/114-authenticated-thumbnails
./devbuild.sh
./dist/newsgo-darwin-arm64 -version
# Run ./newsgo.sh only after separate restart authorization.
```

Use the owner's **2026 M6 Apple Silicon MacBook Pro**, never the obsolete Intel hardware. `./devbuild.sh` already runs Go tests and builds ARM64 and compatibility AMD64 binaries. The version check is optional, not a startup prerequisite. The owner's private file remains as it is; verify that both running services load the existing credential, then retest Safari personal/Guest, initial load/tab focus, normal thumbnails, RSS-only mode and other API denials.

**Promotion authority:** Both PRs remain draft; no `FINISH`, merge, issue closure, branch deletion, restart, deployment, tag, release or semantic version change authorized.
