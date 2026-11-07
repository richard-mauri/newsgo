# Current work

**Merged baseline:** GitHub `main` at `046783f1c020938a4d431711eb5bb905ea179727` after owner-authorized merge of PR #66, verified November 7, 2026. Feature Issue #65 closed completed. No change to ScraperGo.

**Active issue:** [#67](https://github.com/richard-mauri/newsgo/issues/67) — migrate NewsGo from legacy Server Sync administration to Homepage-first authorization, and show reader-owned diagnostic/cache metrics. Branch `feat/67-homepage-auth-diagnostics`, created from the verified merged main, one draft PR. Do not merge without owner's explicit authorization.

**Initial committed slice:** The existing latest-reading Diagnostics screen is accessible to any reader, without requiring a legacy Server Sync token. Admin Settings remains legacy-admin-only while the server-verifiable Homepage role-validation integration is designed; no browser-only privilege promotion, no legacy token or user data modification. Regression test checks diagnostics visibility and retains admin gating.

**Incomplete / next:** Add per-refresh classification cache hit, miss, bypass statistics (currently no counters). Implement NewsGo Go server-side verification of the Homepage bearer against fixed Homepage session/profile endpoint with bounded network timeouts and appropriate role checks for administrator APIs; preserve legacy credentials as compatibility while migrating. Avoid treating the browser-reported `role` as a server privilege. Add negative security tests and verify how Homepage's session handoff can be validated safely from NewsGo's deployment.

**Review and acceptance:** GitHub CI must pass at exact draft PR head; Safari on Mac remains untested. The owner should check that Diagnostics is visible after normal reading while Admin Settings still requires authorized legacy credentials until server authorization migration is completed. This draft is **not feature-complete** and must not be merged.

**Mac review commands:**
```sh
cd ~/projects/newsgo
git fetch origin
git checkout feat/67-homepage-auth-diagnostics
git pull --ff-only origin feat/67-homepage-auth-diagnostics
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

**Safety:** No changes to approved-token files, existing personal subscriptions/profile records, Cloudflare settings, local build counters, or production privilege validation. No tags, release, or merge without owner approval.
