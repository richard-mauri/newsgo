# NewsGo — AI Development Process and Chat Handoff

**Read this file first when coding or resuming NewsGo in a new chat.** `README.md` is the NewsGo user guide, not a ScraperGo document.

## Repository identity and critical documentation correction

- Repository: https://github.com/richard-mauri/newsgo
- Local directory: `~/projects/newsgo` on the established Mac setup.
- Local dashboard: `http://127.0.0.1:8084/`.
- Go project with `cmd/newsgo`, `cmd/feedcheck`, `cmd/discoverycheck`, `internal/dashboard`, `internal/feeds`, and `feeds.json`.
- Existing scripts: `./devbuild.sh` and `./newsgo.sh`. The startup script launches NewsGo's own Intel executable and **dedicated Cloudflare Tunnel**, keeping ScraperGo and other tunnels untouched.
- ScraperGo is **a separate article reader** (historically `http://127.0.0.1:8083/` or `https://scrapergo.themauris.org/`). NewsGo sends articles there but does not share its extraction code.

**Documentation corruption on GitHub, observed October 20, 2026:** Both `README.md` **and existing `AGENTS.md`** in the NewsGo repository currently contain **ScraperGo** documents. These replacement files deliberately restore NewsGo identity. Existing `NEWSGO-DASHBOARD.md`, `NEWSGO-v0.2.20.md`, `NEWSGO-DISCOVERY-EVALUATION.md`, and other NewsGo-specific historical files remain useful background; do not delete them just because these two documents are corrected.

## Current handoff / uncertainty — October 20, 2026

The GitHub `NEWSGO-DASHBOARD.md` documents v0.2.0 dashboard work and an older `NEWSGO-v0.2.20.md` states that generic shared headlines are hidden on initial load and shown when an explicit built-in publisher is selected, while personal-subscription headlines load as before. **These notes do not establish the current installed version.** The currently running NewsGo version/build, Git status, last pushed commit and the full list of new features are **not verified here**; check local source, app display and current relevant files before proposing a next build.

Do not assume a NewsGo semantic/build version from ScraperGo's v0.2.55 or from any filename ending in v0.2.20. Preserve FeedCheck CLI and the personal vs built-in feed behavior.

## Functional and security invariants

- Feed subscriptions are stored in `feeds.json`; don't overwrite a user's subscriptions in patches. Preserve existing feed configuration and migration behavior.
- NewsGo merges and sorts headlines from multiple RSS feeds; deduplication and cached stories prevent duplicate display and avoid losing earlier fetched stories when one feed fails.
- Dashboard headline links open ScraperGo by default in the documented integration; **Original** opens the publisher directly. Do not accidentally alter ScraperGo's separate code or URL scheme.
- Preserve separation between personal subscription feeds, built-in publisher/discovery sources, feed filtering and explicit selection. The historical v0.2.20 note identifies intentional hidden-on-initial-load behavior for generic shared headlines.
- Treat publisher feed responses as untrusted input; retain request limits, timeouts, filtering and HTML safety. Existing dashboard notes describe 4 concurrent feed requests, 4 MiB per feed, 12-second timeout, in-memory 600-headline limit and refresh minimum of one minute; **verify current values in source** before treating them as guaranteed.
- Local loopback binding is intentional; no built-in authentication is documented. Don't expose open dashboards or feeds without appropriate access controls.
- Never overwrite `.cloudflared/token`, tunnel paths, `feeds.json`, `newsgo.sh`, `devbuild.sh`, `go.mod`, or unrelated code in a small reader/dashboard patch unless the user explicitly asks and changes are reviewed.

## Build, patch and operational workflow

Use the complete Go package, not one file:

```sh
cd ~/projects/newsgo
go test ./...
./devbuild.sh
./newsgo.sh
```

A source-only development run is documented as `go run ./cmd/newsgo`, after ensuring there isn't a competing service on port 8084. `go run ./cmd/feedcheck` is the separate FeedCheck command. The existing `newsgo.sh` starts the NewsGo application and its own Cloudflare Tunnel; do not edit or restart another project's tunnel.

Patch only necessary changed files, with updated complete `AGENTS.md` and `README.md` when delivering a handoff, and do not replace personal `feeds.json`. Follow the established version/patch convention **after inspecting current source and prior documented packages**; do not infer a version or impose another project's patch rules. Run `go test ./...`, check dashboard load and filters, personal feeds vs default sources, headline links to ScraperGo, Original links, refresh/error handling, and FeedCheck behavior. Separate actual local tests from unverified publisher network behavior.

## Non-negotiable agent conduct

- Read this file and `README.md` in full before changing code, advising Git operations, or claiming completion.
- The user's verified **local source tree** is authoritative whenever it differs from GitHub. Obtain current affected files and reconcile local changes; never overwrite them from older remote sources.
- Trace **both definitions and actual call sites**, handlers, UI uses, tests, and configuration dependencies. A new function that is not called is **not** a finished fix.
- Protect unrelated working behavior. Keep changes minimal and test existing regression areas.
- Distinguish **drafted/generated**, **verified file exists**, **formatted**, **tested in full project**, **tested in isolated harness**, **installed**, **user-verified live**, **committed**, **pushed**, **tagged**, and **released**. Never imply one means another.
- If you cannot run tests or produce a real artifact, say so. Never invent a download link, test result, deployment, or Git action.
- Do not commit, push, tag, release, restart services, deploy, or modify tunnel configuration without the user's explicit authorization.

## Mandatory Git pre-commit audit — learned from ScraperGo failure

**Never compose a `git add` filename list from memory, an earlier ZIP manifest, or an assumption about changed files.** A previous agent omitted a crucial `readability.go` integration change in ScraperGo; the user found it before pushing. This is a cross-project safeguard.

For a commit request, **separately from normal patch installation**:

1. Obtain fresh `git status`, including staged, unstaged and untracked files, and ahead/behind state.
2. Examine `git diff --check`, `git diff`, `git diff --cached`, and relevant untracked files. If commits already exist locally, review `git log origin/main..HEAD` and commit file inventories.
3. Check all changed source files **and integration points**, tests, both documentation files, and build/configuration artifacts as applicable. Exclude secrets, local tokens, binaries, logs and unrelated temporary files.
4. Only after that evidence recommend exact staging commands. Before committing, verify `git diff --cached --name-status` and `git diff --cached` against the intended scope.
5. Do not infer that a clean working tree means the correct commits were made, or that ahead-of-origin commits were pushed. Verify publication independently.


## Mandatory synchronized version update (regression prevention)

**On EVERY installable NewsGo regeneration, update the version in ALL version-dependent locations together, before packaging.** The earlier v0.2.29 patch updated the UI to v0.2.29 but left five tests asserting v0.2.28, causing `go test ./...` to fail. This is a recurrent AI-agent error and a hard release blocker.

Before distributing a patch:

1. Inspect the current dashboard header version, script cache-busting URL (`/personal.js?v=...`), relevant CLI/build version indicators, `AGENTS.md`, `README.md`, and **every test containing an explicit version string**.
2. Update all necessary occurrences to the SAME new version. Do not simply change expectations to hide a real bug; verify the tests still assert meaningful behavior.
3. Search the entire affected project for stale versions from the previous build; decide explicitly whether historical references are intentionally retained.
4. Add the changed test files to the patch inventory and ZIP. Never omit version-dependent tests from a regeneration.
5. Run `go test ./...` on the full project when possible; if unavailable, state that clearly and require a passing test suite before build/restart. A patch that knowingly fails version assertions is incomplete.
6. Confirm the running page shows the new version and Safari receives the matching versioned JavaScript. An old Dock-launcher/server process can show an outdated app; inspect the actual process and URL before attributing a UI failure to the new build.

A **version change is an integration change**, and the changed-file inventory must include its matching test and documentation updates.

## Cross-chat handoff obligations

At the start of a new chat, use `AGENTS.md` for the **current engineering handoff** and `README.md` for the public-facing project context. Confirm current local source and Git state before applying edits. Maintain this file with a short, dated snapshot of the most recent verified baseline, pending work, confirmed tests, known regressions and next action; **replace obsolete current-state claims**, rather than accumulating contradictory updates. Keep the README focused on users and general architecture, not chat history or agent rules. Preserve the existing project's patch/install convention; do not transplant another project's conventions blindly.


## Development candidate v0.2.26 — publisher exclusion list (October 20, 2026)

**Baseline:** user-supplied `newsgo-source.zip`, containing current `internal/dashboard/` plus both documentation files. Its dashboard displayed **v0.2.25**. GitHub may differ; do not replace this development copy with older remote files. The user requested a NewsGo-only publisher exclusion feature. **Do not modify ScraperGo**.

**Change:** browser preference key `newsgo.excludedDomains.v1` stores normalized lower-case publisher domains separately from the existing `newsgo.subscriptions.v2` registry. Domain matching uses full hostname boundaries (exact or subdomain); URLs, `www.` names, and bare hosts are normalized. Exclusion is reversible and **never unsubscribes feeds**. NewsGo's shared and personal headline presentations hide excluded article destinations; discovery results and existing subscriptions show an exclusion marker with allow/exclude controls. Added `Excluded Sites` manager and one-click controls next to headlines. New exports contain subscriptions and exclusions in `newsgo-preferences-v3`; imports still accept legacy subscription arrays and leave existing exclusions unchanged when importing old files.

**Source call sites:** `internal/dashboard/dashboard.go` adds the manager UI, source-domain lookup on included cards, and increments visible/cache version to `v0.2.26`; `internal/dashboard/personal_script.go` handles localStorage, matching, UI rendering, filtering and export/import; `internal/dashboard/exclusions_test.go` checks wiring. Preserve the HTML template JS `<script>` loading, `renderReading`/`arrangeSharedLeads` shared cards and `renderPersonalArticles` personal cards; regression-test both paths and subscriptions.

**Verification:** JavaScript syntax validated with Node, plus a focused Node test of domain normalization, subdomain containment, and hostile lookalike domains. Go formatting performed. The upload does **not** include `go.mod` or the entire `internal/feeds` source, so full `go test ./...` could not run in the provided partial package; user must run it in `~/projects/newsgo` before building. Browser appearance, live import/export, and public deployment are not yet verified. This is an **unreleased candidate**, not a confirmed working feature. Next step: install patch, run full tests, verify page/UI, and report any failures.

**Installation:** the direct-unzip ZIP includes only modified `internal/dashboard/dashboard.go`, `internal/dashboard/personal_script.go`, new `internal/dashboard/exclusions_test.go`, complete `AGENTS.md`, and complete `README.md`. Run:

```sh
cd ~/projects/newsgo
unzip -o ~/Downloads/newsgo-v0.2.26-publisher-exclusions.zip
go test ./...
./devbuild.sh
./newsgo.sh
```

If tests fail, **stop** before rebuilding and report the exact failure. After browser testing, keep the source-only export/import, subscription persistence and exclusions independently validated. Before any commit, follow the mandatory actual-Git-status staging audit above.


## NewsGo v0.2.27 — Follow-up to failed v0.2.26 candidate (October 20, 2026)

The v0.2.26 direct-unzip candidate failed seven dashboard tests in the user's full checkout. Several legacy tests asserted the previous v0.2.25 HTML version string; the new exclusion test made an overly exact template-tree assertion; the combined-subscription test asserted an original shared-headline selection expression. This follow-up preserves the original selection expression and applies publisher exclusions using a separate `eligible` predicate. It also activates shared-headline exclusion controls in the render path and updates old version assertions to v0.2.27. No claim of end-to-end verification until `go test ./...` succeeds on the user's checkout and the browser flows are manually checked. The code belongs solely to NewsGo. Do not change ScraperGo.


## Current development candidate — v0.2.28 (October 20, 2026)

Publisher-domain exclusions were rebuilt from the user-uploaded v0.2.25 dashboard source, **not** layered on the failed v0.2.26/v0.2.27 candidates. This is a NewsGo-only feature; ScraperGo is unchanged. The subscription registry (`newsgo.subscriptions.v2`) is preserved and exclusions are stored separately (`newsgo.excludedDomains.v1`). The dashboard includes an Excluded Sites panel, exclusion indicators in Find Subscriptions and My Subscriptions, one-click exclude/allow actions, and article-domain filters for personal and shared headlines. Export now includes subscriptions and exclusions, and import accepts existing array-only exports. Existing browser storage is not cleared or reset. The visible page header and personal.js query string both identify v0.2.28.

Important acceptance tests: initial load shows no stray navigation or empty content when subscriptions exist; source discovery Subscribe operates; subscribed feeds remain subscribed when hidden by an exclusion; My Subscriptions and Find Subscriptions mark exclusions; headlines hide/restore; both old and new export formats import without wiping saved records; Safari reload retains data. Test against current local working tree and do not deploy on a failing `go test ./...`. GitHub source may lag local files. v0.2.28 is an unverified development candidate until full local Go tests and Safari behavior checks pass.

## NewsGo v0.2.29 — Excluded publisher navigation (October 20, 2026)

User verified v0.2.28 visually with 8 subscriptions and 1 excluded domain, but an excluded subscribed publisher retained an active Read Articles button. Clicking it opened an empty publisher headline page. This candidate hides the Read Articles action for excluded shared and personal subscriptions, automatically returns to the combined All Subscriptions view when the current publisher becomes excluded, and prevents direct/selector navigation to excluded publishers. Subscriptions and exclusions remain independent and persistent. `internal/dashboard/personal_script.go` and `internal/dashboard/dashboard.go` are the only changed implementation files; v0.2.29 updates both UI version and JS cache-busting URL. `exclusion_navigation_test.go` adds contract tests. No ScraperGo code changes. Source baseline: user-installed v0.2.28 candidate, per screenshots. Full `go test ./...` and live Safari testing remain required before commit/push; source/JS syntax verification alone is not end-to-end verification.

## NewsGo v0.2.30 — synchronized version-contract correction (October 20, 2026)

The user ran `go test ./...` on v0.2.29 and five tests failed solely because they checked v0.2.28 header or `/personal.js` strings. v0.2.30 keeps v0.2.29's navigation behavior, changes dashboard header and cache-busting script reference to v0.2.30, and updates the five version-dependent regression tests. The rule above makes this synchronization mandatory in future patches. The complete Go suite and live UI must still be checked on the user's Mac before running `./devbuild.sh` and `./newsgo.sh`. Do not commit/push until verified.

## Current handoff — v0.2.31 (October 20, 2026)

**User-observed v0.2.30 baseline:** My Subscriptions and Excluded Sites work; user demonstrated the excluded publisher is retained and its Read Articles button is hidden, but the personal-feed row showed only Unsubscribe. The requirement is for a reversible **Allow this site** button on every excluded subscription, shared or personal.

**v0.2.31 change:** Both subscription row types now add the existing `excludeButton(domain)` action, so allowed domains show **Exclude this site** and excluded domains show **Allow this site**. No changes to the subscription registry, storage keys, filtering, or server APIs. The displayed dashboard version and JavaScript URL, all five version-dependent tests, and both docs are updated together. `exclusion_actions_v31_test.go` protects the subscription action. This is a **development candidate**; full `go test ./...` and Safari testing on the user's actual checkout remain required.

**Mandatory version-sync gate:** Every new installable build **MUST** update (1) visible dashboard/app version, (2) `/personal.js?v=...` cache-busting reference, (3) any exact-string version assertions in the existing Go tests, (4) `AGENTS.md` and `README.md` handoff/version text, and (5) artifact filename. Search every affected test file for prior version strings before packaging. An inconsistent version or an unupdated test assertion blocks delivery; never make the test pass by reverting the app version. Preserve actual user preferences and exclude ScraperGo from this work.
