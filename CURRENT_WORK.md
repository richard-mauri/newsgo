# Current work — NewsGo

**Active Issue:** [#116](https://github.com/richard-mauri/newsgo/issues/116) — follow-up thumbnail diagnostics and administrator settings after merged routing PR #117. **New branch:** `fix/116-thumbnail-diagnostics-settings`, based on verified merged `main` `4134c7bf8b546a56a4a88ce7e96ec7008b9d7030`. Previous routing branch `fix/116-thumbnail-loopback-backend` is merged and must not be reused.

## Scoped implementation

The owner issued NEXT after the explicit FINISH of PR #117, authorizing a fresh branch and draft PR. This phase separates semaphore queue admission from the upstream timeout; configurable bounded settings default to concurrency 6, queue wait 12 seconds and upstream timeout 12 seconds. A new same-origin redacted response-header summary powers Diagnostics → Thumbnail Performance with observed queue/upstream mean and p95 and failure categories. Admin Settings has a new server-side validated thumbnail section, following the existing authorized, revisioned save path and remaining compatible with old settings files without writing a migration on read. Fixed media limits (3 MiB, allowed types, width, no redirects and URL restrictions) remain unchanged. No browser scheduler, ScraperGo change, credential rotation, classification cache namespace change, app version bump or production launcher modification.

## Safety and verification

Read authoritative Homepage `SUITE_DEVELOPMENT_POLICY.md` v1.1.0 and NewsGo's five required documents. No owner-local files or processes are accessible to this GitHub-connected agent; owner-local dirty or ahead state remains **unverified** and must be checked before fetching. The original service token, Homepage bearer verification, admin-only settings, per-user storage, subscriptions, exclusions, browser histories, classification caches, tunnels and private configuration must be preserved. The branch adds deterministic Go tests for admission timeouts, diagnostics headers, legacy settings and form layout; CI-only browser checks still apply. GitHub Actions exact-head CI and Mac/Safari/live checks remain unverified until observed.

## Owner Mac handoff and promotion authority

Only after exact-head CI passes: `cd ~/projects/newsgo`, `git status --short --branch` (stop if dirty), `git fetch origin`, `git checkout fix/116-thumbnail-diagnostics-settings`, `git pull --ff-only origin fix/116-thumbnail-diagnostics-settings`, then `./devbuild.sh` **once** and `./dist/newsgo-darwin-arm64 -version`. `./newsgo.sh` restarts NewsGo and its dedicated Cloudflare tunnel and requires separate authorization; never execute without it. Owner Safari acceptance should be limited to a normal All Subscriptions run and examining the new Diagnostics panel, with optional administrator settings validation. No FINISH, merge, issue closure, branch deletion, deployment, restart, tag or release has been authorized.
