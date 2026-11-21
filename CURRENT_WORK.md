# Current work — NewsGo

**Merged main baseline:** `3aeeb3426c1eb896dbda6c719c3986a2456a8ac6`, verified before starting. **Active settings scope:** [Issue #112](https://github.com/richard-mauri/newsgo/issues/112), branch `feat/112-explicit-settings-save`, draft PR to `main` (verify actual number/head). **Concurrent independent authentication work:** [Issue #110 / draft PR #111](https://github.com/richard-mauri/newsgo/pull/111), branch `auth/110-verified-suite-access`. Do not merge or modify that PR automatically. Both phases must be reconciled before production adoption.

**Owner requirement:** Match MauriStocks explicit Save Changes for NewsGo user settings and subscriptions, not NewsGo's legacy immediate browser writes and automatic 650ms Homepage save/30s retry. Current scoped branch stages subscriptions, exclusions, age, lead count/mode, RSS metadata-only, classification cache and reading selection in tab memory; gives one reader Save Changes action; performs transactional browser writes followed by one user-triggered personal-account PUT; Guest saves local preferences only. User gets Safari native unsaved-changes warning. Failed Homepage PUT is stored in an owner-specific browser snapshot for manual retry after reload, never associated with another user. Initial account GET must preserve in-progress edits. Recent search/website histories remain immediate local operational history, and Admin Save Server Settings stays separate. No Homepage framework adjustment required.

**Data safety:** No production records, private subscriptions, exclusions, browser history, personal settings, bearer credentials, local Mac files, build counters, listeners, process restarts, tunnels or DNS touched. The owner Mac's uncommitted/unpushed state is inaccessible. No change to account schema or server credentials. No new local build dependencies.

**Docs:** Read merged Homepage `SUITE_DEVELOPMENT_POLICY.md` v1.1.0 and NewsGo AGENTS, README, ARCHITECTURE, DEVELOPMENT, CURRENT_WORK. README, ARCHITECTURE, DEVELOPMENT and CURRENT_WORK updated. AGENTS and shared policy unchanged.

**Verification:** New Go static contract checks and Node.js behavioral tests execute in GitHub CI only; exact-head Go CI and executable build must pass. Owner Mac `./devbuild.sh` and Safari/Guest/personal/account save failures, reload, data isolation and final combined authentication integration remain unverified.

**Mac handoff:** `cd ~/projects/newsgo`; `git status --short --branch` and stop if dirty; `git fetch origin`; `git checkout feat/112-explicit-settings-save`; `git pull --ff-only origin feat/112-explicit-settings-save`; `./devbuild.sh`; `./dist/newsgo-darwin-arm64 -version`. No `./newsgo.sh` unless separately authorized.

**Promotion:** Issue open and PR draft; no merge, Issue closure, branch deletion, restart, deployment, tag, release, or version changes authorized.
