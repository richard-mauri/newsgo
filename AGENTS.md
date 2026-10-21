# CURRENT CANDIDATE — NewsGo v0.2.47 admin server settings (2026-10-21)

Based on authoritative user-provided source matching GitHub commit `598ec5b`. This direct-unzip candidate adds local CLI `-set-role admin` / `-set-role user` to designate an existing hashed token without rotating it, admin-guarded `/api/admin/settings` GET/PUT, separately revisioned atomic `data/server-settings.json`, and a browser Admin Settings panel revealed to the connected admin token. Image metadata timeout, maximum HTML response, concurrent lookups, enrichment deadline, and missing-image cache TTL now use bounded live server settings. No changes to user preferences or secrets. The default values remain 2500 ms, 512 KiB, 12, 25 seconds, 90 seconds. Tests passed in an isolated Go 1.23.2 compatibility run (temporary go.mod only, original go.mod retained); the required Go 1.27.1 project tests, ./devbuild.sh and live Safari acceptance remain unverified. No Git commit, push, restart or release authorized. ScraperGo and NewsValidate will be handled separately after the CNBC investigation. The existing historical handoff below is retained for context; old "next step" statements are superseded by this paragraph.

## NewsGo administrator settings candidate — October 21, 2026

Based on the user's clean commit `598ec5b` and source archive. Server tokens support optional `role: "admin"`, assigned locally with `go run ./cmd/newsgo-token -dir data -user richard -set-role admin` (use the existing user ID; no token rotation). `GET/PUT /api/admin/settings` require the token's admin role, use revision checks, JSON atomic writes, input bounds, and server-wide image lookup parameters. The Admin Settings browser panel is available only after admin Server Sync connection; authorization is enforced server-side. Existing subscriptions, exclusions, tokens, data files and tunnels are not replaced. Candidate requires full local tests and browser acceptance prior to commit or deployment. ScraperGo and NewsValidate remain separate projects in the planned shared extraction architecture.

---

# NEW DEVELOPMENT PATCH — NewsGo executable-version introspection (2026-10-21)

Prepared against verified clean GitHub/local HEAD `9766238` (both local and `origin/main` confirmed by user). This source-only, direct-unzip v0.2.46 follow-up adds `newsgo -version`, updates the startup log to derive its version from the existing dashboard header, and updates `newsgo.sh` to print/check the actual `dist/newsgo-darwin-amd64` path and version *before* stopping the existing process. Files: `cmd/newsgo/main.go`, `cmd/newsgo/version_test.go`, `internal/dashboard/version.go`, `newsgo.sh`, `README.md`, `AGENTS.md`. This patch does NOT change the dashboard header version, `personal.js` cache buster, feeds, secrets, or the other projects. It does not add release tooling. Shell syntax and Go formatting have been checked in the packaging environment; full `go test ./...`, `./devbuild.sh`, running process and Safari behavior remain UNVERIFIED and must be checked on the user's Mac. Do not commit/push/restart without explicit authorization. `./devbuild.sh` must be used to refresh the launched `dist/` executable. Older v0.2.46 handoff statements below are historical until updated after local acceptance.

**Multi-project architecture remains on the roadmap:** ScraperGo v0.2.56 will supply reusable article/HTML/image extraction utilities from its existing repository; NewsGo will consume appropriate image metadata; NewsValidate v0.3.2 will migrate only low-level parsing helpers and preserve independent reference selection. No changes to ScraperGo or NewsValidate in this patch.

---

# ACTIVE HANDOFF — NewsGo (2026-10-21; read before historical notes)

**Scope:** This is the authoritative *conversation handoff*, not a claim about the current remote Git branch or running process. The user's local checkout is authoritative. Historical sections below describe prior candidates and are not instructions to revert the current work. Do not begin a new feature until checking local status.

## Verified current state vs outstanding checks

- NewsGo **v0.2.46 manual subscription search** is installed in the user's local source tree from `newsgo-v0.2.46-manual-search-corrected.zip`. The user ran **`go test ./...` successfully for every package** and ran **`go build -o newsgo ./cmd/newsgo` successfully**. `git diff --check` also passed. Do not represent this as a proven deployed build or pushed commit.
- The intended UI change: explicit **Search** button / Enter to submit a subscription-directory search; keystrokes alone do not submit; retain prior results during a request, show meaningful progress, and handle HTTP 429. Existing feed-search backend, subscription records, exclusions and sync should remain intact. **Live manual-search acceptance after rebuilding the correct executable has not been demonstrated here.**
- `git status --short` before commit showed modified `AGENTS.md`, `README.md`, dashboard implementation and eight existing tests; new `manual_search_v46_test.go` and `version_consistency_test.go`; and **untracked root executable `newsgo`**. `git diff --check` was clean. The user said “we should commit first” and was given staging/push commands, **but no subsequent confirmation that commit or push happened**. New agent MUST check `git status -sb`, `git log -3 --oneline`, upstream/ahead-behind, and staged diff before describing anything as committed/pushed.
- **Critical deployment trap:** `go build -o newsgo ./cmd/newsgo` writes a root executable **NOT** used by the established launcher. `newsgo.sh` launches `~/projects/newsgo/dist/newsgo-darwin-amd64`, which `./devbuild.sh` builds. User observed docked/running NewsGo **v0.2.44** despite local v0.2.46 code/tests. Treat the running version as **unverified pending explicit check**. Do NOT recommend generic `go build -o newsgo`; use the documented `./devbuild.sh` → `./newsgo.sh` path only after tests and with user approval. The untracked root binary should NOT be committed.
- User explicitly requests (still **NOT implemented/verified**) **`newsgo -version` CLI support** and `newsgo.sh` printing/reporting the *actual launched binary version*. Do not claim present. Implement as a tightly scoped subsequent change with tests, and be careful that any Dock launcher may invoke a different script/process.
- ScraperGo baseline discussed: **v0.2.56** with localhost diagnostic HTML extraction; NewsValidate baseline **v0.3.2**, with direct publisher HTML comparison, diagnostic recovery, and article-footer audit. No commits/pushes for these projects have been confirmed in this conversation. Resume NewsValidate later, not as part of NewsGo patch.

## Required next actions, in order

1. **Close the current NewsGo v0.2.46 commit checkpoint before architectural work.** Inspect the *actual local checkout* and its `AGENTS.md`/`README.md`, check Git status/diffs, confirm tests, run `./devbuild.sh`, run/restart via `./newsgo.sh` (only with user approval), and inspect the browser's visible version and explicit search. Stage only intended source/docs/tests; NEVER include the `newsgo` binary, `dist/`, `.cloudflared/`, `data/`, secrets, logs, `feeds.json`, or unrelated files. Commit/push **only on explicit user instruction**; user's preceding “we should commit first” establishes the desired sequencing, not proof of completion. Verify actual commit and remote status after any Git action.
2. **Next narrow NewsGo improvement:** Add one authoritative executable version constant/flag **`-version`** (should report v0.2.46 if implemented without a version bump, otherwise synchronize all visible versions), and adapt `newsgo.sh` to report the exact binary path and version it launches. Keep `./devbuild.sh` as the build command. Ensure Dock/startup consistency and test CLI behavior. Do not accidentally restart or reconfigure ScraperGo's tunnel.
3. **Then address missing CNBC thumbnail**: NewsGo shows no thumbnail for `https://www.cnbc.com/2026/10/13/inflation-hedge-tips-commodities-infrastructure.html`. User supplied NewsGo/ScraperGo screenshots and exported `article-3.md`. Markdown begins with the photographer credit “J Studios | Digitalvision | Getty Images” but provides **no usable lead-image URL**; the only Markdown image is an irrelevant “hide content” SVG. Do NOT mistake that credit for a resolved image URL or claim the peach-colored screenshot area proves a successfully rendered image. Inspect publisher HTML/feed metadata, NewsGo `articlePageImage`/`articleBodyImage`, and ScraperGo's structured `imageUrl` first; distinguish discovery from rendering and caching.
4. **Architecture decided in principle, not implemented**: ScraperGo (lowest service layer) should own reusable HTML/article parsing and lead-image selection, likely a public Go package inside the **existing** `github.com/richard-mauri/scrapergo` repository. No new Git repo. NewsGo may use a lightweight ScraperGo image-metadata endpoint or existing `/api/read` `imageUrl` field; assess payload/cost and reliability before adding an endpoint. NewsValidate should use only **low-level parsing utilities** from that shared code, retaining its independent article-reference selection so ScraperGo defects cannot validate themselves. Migrate incrementally with regression tests, avoid a cross-module dependency cycle, and do not rewrite all three applications at once.
5. Keep existing regression evidence: Seahorse Magazine Finn missing thumbnail (v0.2.44 fallback), CNBC lead image, Ars Technica where Go HTML reference found 16 paragraphs and Trafilatura only six, Guardian footer-resource false positive fixed in NewsValidate v0.3.2, Sailing Anarchy HTTP 202 normal retrieval vs diagnostic HTML success. Preserve tests for those behaviors where relevant.

## Mandatory engineering conduct (user has reiterated emphatically)

- **Read `AGENTS.md` and `README.md` fully *before* coding, packaging, offering CLI steps, or committing.** Inspect `devbuild.sh` and `newsgo.sh` before any build guidance. Follow existing instructions, do not replace them with generic suggestions. The user should not have to repeat them.
- **No false “ready/verified” claims.** Distinguish source edited, ZIP created, syntax checked, complete `go test ./...` passed, `./devbuild.sh` passed, deployed version verified, Git committed, and pushed. If full-project tests cannot run, identify the block and obtain the user's local output; do not deliver an unverified patch as a verified release.
- **Do not encode fixed release numbers into ordinary functional tests.** Version-consistency checks must derive expected cache-busting/header version from the single current version rather than test constants. On every change verify dashboard version, JS cache buster, CLI version (once added), launcher reported version, tests, README, AGENTS, and ZIP filenames agree.
- **Use direct-unzip ZIP patches with project-relative files at ZIP root**, containing updated `AGENTS.md` and `README.md`; never make the user manually edit code. Do not bundle unrelated files or remove existing work. Review old and new test expectations and all integration call sites before changing behavior.
- **Do not create additional build/release tooling just to compensate for agent inattention.** Existing scripts and instructions are sufficient; the agent must follow them. Respect the user's Intel Mac and prefer prebuilt tools over costly local compilation.
- Never edit NewsGo's user `data/`, Cloudflare tunnel tokens, `feeds.json`, other project services, or unrelated code. Inspect actual local sources when they differ from GitHub. Request missing local files rather than silently using stale GitHub versions.

---

## v0.2.44 — Missing-thumbnail retry and lookup budget

The Seahorse Magazine Finn article had no thumbnail even after the v0.2.44 body-image fallback. Preserve feed/metadata/body extraction; allow 12 parallel metadata lookups within a 25-second shared deadline so late article lookups are less likely to be abandoned. Cache successful feeds containing image-less articles for 90 seconds rather than ten minutes to recover from transient metadata failures. Regression covers the Finn lead-image HTML and retry TTL. Live publisher image retrieval remains a browser acceptance test; no ScraperGo changes or data migration.

## v0.2.44 — article body thumbnail fallback (October 20, 2026)

The existing feed image and OG/Twitter metadata fallbacks remain primary. If those are missing, the same bounded article HTML response is examined for a plausible `<img>` in the article/main content; cross-publisher HTTPS image URLs are permitted. No new ScraperGo dependency. Regression: `article_metadata_test.go` Seahorse Magazine case. Keep per-feed caches, SSRF controls, and separate publisher/article-domain exclusions.

## v0.2.42 — article destination exclusions (October 20, 2026)

From released v0.2.41, add a second exclusion registry for *article destination domains*, separate from publisher/whole-feed exclusions. For example, youtube.com blocks youtube.com and www.youtube.com article links and youtu.be links across all subscriptions, without unsubscribing from KCRA 3. New article-domain exclusions are held in newsgo.excludedArticleDomains.v1 browser storage and synchronized using optional excludedArticleDomains in the existing per-user JSON profile. Do not change the interpretation of excludedDomains (publisher exclusions). Import/export both lists; refresh visible stories after an article exclusion change. Version all visible version strings and tests. Avoid touching ScraperGo. Tests run with the installed project Go toolchain; Go 1.23 fallback test is not the required production toolchain. Do not commit private data or tunnel files.

## v0.2.42 — visible token placeholder (October 20, 2026)

Server Sync token input placeholder now explicitly says `Enter the 64-character approved token distributed by Richard.` (not only validation error text). Presentation-only change; no sync behavior or JSON data changes. Tests assert actual input markup; preserve the version and script-cache synchronization.

## v0.2.42 — presentation copy (October 20, 2026)

Update header tagline to `NEWS READING WITHOUT PAGE CLUTTER USING SCRAPERGO`, retaining displayed version `V0.2.42`. Change the missing/malformed token instruction to `Enter the 64-character approved token distributed by Richard.` Keep sync and storage behavior unchanged; update JavaScript query version and locked Go tests. Validate syntax, full Go tests, and staged file inventory before commit.

## v0.2.38 — duplicate refresh status removal (October 20, 2026)

Remove only the redundant upper-right `#subscription-updated` element and its personal_script.go assignment. Keep the detailed `#personal-loading` progress/status (updated feed count, errors, visible headlines) and spinner. The old v33 status test now asserts absence of the duplicate and presence of the detailed status. Sync, feeds, preferences, and exclusions unchanged. All exact version test assertions synchronized to 0.2.38. Validate full Go suite on the user Mac before commit.

## v0.2.37 reading navigation cleanup — October 20, 2026

Remove the redundant standalone `← All Subscriptions` button (`#all-subscriptions`) and its JS references; do **not** remove `chooseAllSubscriptions()`, the My Subscriptions `Read All Articles` action, automatic reset when selecting an excluded publisher, or the selected publication heading. Preserve all server sync and exclusion persistence behavior. `navigation_all_subscriptions_test.go` tests retained navigation and absence of the old button. Dashboard header/JS version and all version-locked tests are v0.2.37. Validate `go test ./...` on the user's Go toolchain before committing.

## v0.2.36 excluded-site display — October 20, 2026

The Excluded Sites panel previously showed raw domains only; a Hacker News exclusion was saved as `news.ycombinator.com`, making the publisher appear missing. Show the matching subscription names and the unchanged domain, using a comparison scoped to each excluded-domain entry. Refresh the panel on open. Preserve `excludedDomains` JSON, token sync, and subscriptions. The count must remain the number of excluded domain entries. Regression: `exclusion_manager_v36_test.go`.

## v0.2.34 cleanup and discovery fixes — October 20, 2026

The catalog exists for discovery only. `State.Refresh`, `State.Run`, shared-image enrichment, cached built-in articles, and the dead periodic-refresh CLI option were removed. Never reintroduce an automatic catalog headline pipeline. Preserve the v0.2.32 unified subscription registry, migration, exclusions, and reading behavior.

Discovery: Feedly exact-name searches may be empty while a prefix returns results, so an empty complete query is retried with one final rune removed. Results prefer exact and matching publisher names. Publisher feed failures must be shown as upstream errors, not NewsGo server failures. Subscription cards must not overlap on narrow layouts.

Tests required: `go test ./...`, JavaScript syntax, version sync, and regression tests. The direct-unzip patch should be applied to a clean v0.2.32 baseline; never stage or push changes without user approval.

## v0.2.34 unified subscription architecture — October 20, 2026

Implementation candidate from the complete user-provided `newsgo-v32-source.zip`, based on local v0.2.31. **Every displayed article now comes from a subscribed feed fetched through `/api/personal`.** The built-in feed catalog is for discovery only; the independent automatic `app.Run` feed refresh is no longer invoked at server startup. Browser `readingSelection` has one feed path. Lead stories, latest headlines, sorting, and exclusion all operate on the same subscribed article array. An empty registry yields an empty briefing with a discovery action.

Migration reads `newsgo.subscriptions.v2` and the older legacy keys without deleting them. Recognized `shared:` publisher records are converted to `feed:` records with the URLs configured in `feeds.json`, then stored in `newsgo.subscriptions.unified.v1`. Unknown legacy shared records are retained and labeled for review rather than discarded. Existing excluded domains remain in their original browser storage. Export/import retains the existing format. No automatic migration is possible for unknown catalog names or feeds that are missing from `feeds.json`.

**Testing:** Full Go suite was executed successfully in the supplied archive with an *isolated, temporary* `go.mod` declaration changed from `go 1.27.1` to `go 1.23.2` because the container lacks the required toolchain and cannot download it. The actual `go.mod` in the distributed patch is unchanged. JavaScript parser check and ZIP audit are separately required before delivery. This is not a verified result on the user's installed Go 1.27.1 or Safari; follow normal `go test ./...`, build, restart, and browser acceptance before Git commit.

**Acceptance:** Existing subscribed catalog publishers remain subscribed after conversion; no unsubscribed publisher's articles appear; manual Find Subscriptions and Subscribe work; no automatic 200+ headlines; excluded sources have no headlines; the summary counts only visible subscribed articles; and new header, script cache-buster, tests, and documentation all display v0.2.34. Do not run tests with any toolchain downgrade in the user's actual project.

**Next step:** Install this development patch, execute `go test ./...`, verify Safari subscription persistence, subscribe/unsubscribe, exclusions, lead stories, and empty-view behavior. Stop if anything fails. Do not push, tag, or release without user approval.

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


## v0.2.34 server sync (October 20, 2026)

Adds `internal/dashboard/sync.go` (per-user JSON preferences + hashed approved tokens), `cmd/newsgo-token` (manual issue/revoke), optional UI in dashboard/personal_script.go, and `-sync-dir` in cmd/newsgo. Server profile includes **both** `subscriptions` and `excludedDomains`; browser fallback preserved. Bearer tokens are not stored in localStorage; only held in memory. GET loads a profile, PUT uses optimistic revision concurrency control and atomic file replacement. Data directory must stay outside version control; patch installs must not overwrite user data. Test with three distinct Safari contexts via same tunnel origin.

## v0.2.36 development handoff

Server Sync remembers the per-user token in browser localStorage, loads the server profile on launch, saves both lists on edits, and polls every 30 seconds while visible plus on tab reactivation. Advanced holds manual Load, Save, Disconnect. Stale-write 409 errors preserve local changes and require explicit conflict resolution; polling must never replace unsaved edits. Keep token logging disabled. The original JSON persistence backend and unified subscriber-only headline model remain unchanged. Do not commit data/. Update all visible/template/cache/test versions together.


### v0.2.42: safe reconnection

When a valid token already has preferences on the server and the local browser lists differ, NewsGo offers **Use Server Preferences**, **Merge Preferences**, and **Cancel**. The decision applies to both subscriptions and excluded sites. Merge uses subscription IDs and unique excluded domains, and saves against the server revision; stale writes still receive HTTP 409. A new server profile is initialized from browser preferences as before. This applies to remembered-token reconnects as well as manual Connect.


## v0.2.45 development handoff

Change limited to subscription search UX in `internal/dashboard/personal_script.go` and HTML display version (`dashboard.go`). Existing `/api/search-feeds` per-IP cooldown is three seconds and remains unchanged. Client waits at least 3.2 seconds between requests, retries up to twice on HTTP 429, labels pending/waiting searches, retains earlier results while waiting, and ignores stale async results. Versioned UI tests updated and `search_feedback_v45_test.go` added. Next work: verify live typing `inves`→`invest` and browser accessibility feedback. NewsValidate remains v0.3.2 and is not touched.

### v0.2.46 corrected release gate
Version-sensitive tests must derive expectations from the displayed application version rather than encode a release number. Tests for changed behaviors must be updated with the implementation, not suppressed. The explicit subscription search requires a button or Enter; keystrokes and initial load must not trigger discovery requests. Run the complete `go test ./...` and `go build -o newsgo ./cmd/newsgo` in the actual project toolchain. Do not describe a ZIP as verified unless tests have passed.
