# NewsGo — app-specific development and validation

Read [AGENTS.md](AGENTS.md) and [Homepage's authoritative suite development policy](https://github.com/richard-mauri/homepage/blob/main/SUITE_DEVELOPMENT_POLICY.md). The suite policy, not this file, controls Issues, branches, draft PRs, exact-head CI, visible six-item handoff, owner FINISH and merged branch cleanup. Keep this document for NewsGo-specific validation and commands.

## NewsGo regression and data boundaries

When validating banner spacing after a module update, verify that `/suite-ui/banner.css?v=` changes with the embedded shared stylesheet content. The banner CSS response is cacheable, and a URL revision derived only from `personal.js` can leave Safari displaying a stale layout for up to an hour despite a new executable. The `TestBannerStylesheetURLTracksEmbeddedCSS` regression protects this contract.

For Issue #108, review the shared masthead's NewsGo version/build, small teal **The Mauris** portrait caption, return link and verified Homepage status in desktop/narrow Safari, both themes. Confirm the **original Homepage** ambient shapes (660×530 and 530×500, 70px blur) are visible without changing NewsGo's editorial hero, feed controls, reader, or print. Keep `suite-auth-status` synchronized with the existing profile-verification and recovery flow; do not create a second status line. The stylesheet is same-origin and versioned in the dashboard HTML.

### Mandatory shared navigation button styling

**Canonical implementation:** Homepage now owns `shared/ui/buttons.css` in the nested Go module `github.com/richard-mauri/homepage/shared/ui`. NewsGo pins that module in `go.mod` and serves its embedded CSS from its own origin. The local `dashboard.go` stylesheet may define tokens and button *placement*, but must not reintroduce shared component colors, geometry or hover/focus/disabled styling. Update the dependency through a scoped PR; review the Homepage source and immutable pin, then rerun the UI regression checks. No cross-origin stylesheet or CDN is allowed.

The `.subscription-read` control intentionally behaves visually like an inline underlined link, rather than a standard action button; that explicit exception is not a competing generic-button theme. All ordinary buttons, including compact article exclusion actions, must use the canonical appearance.

All top-level NewsGo navigation buttons, including **Settings**, must use the existing shared `ui-button` and `ui-button-large` styling. Their resting backgrounds, text and borders must match whether selected or not: **no persistent accent fill or accent border** from `aria-pressed="true"`. Apply this rule consistently to every navigation destination; do not introduce a one-off visual treatment for Settings or another section. Preserve accessible pressed-state semantics and visible keyboard focus. Transient hover and `:focus-visible` may highlight the border. UI changes must review the shared selectors in `internal/dashboard/dashboard.go` and the button initialization in `internal/dashboard/personal_script.go`, and retain a regression assertion for the selected-state contract.

The same shared resting appearance applies to **all action buttons**, including Save Server Settings, Reload, Download Article URLs, confirmation-dialog actions, discovery/search, subscriptions, and Diagnostics tabs. No always-on accent border or solid accent/cyan fill may be added for a submit, emphasized, pressed or selected action. All server-rendered buttons must receive `ui-button` at startup; buttons created later must use the shared `button()` helper or explicitly add `ui-button`. Preserve meaningful transient hover/focus, disabled-state behavior and accessible selected-state semantics. Remove obsolete panel-specific color overrides rather than adding higher-priority CSS or `!important` exceptions. Whenever adding a button or editing a button rule, audit panel-specific selectors and update the common-style regression checks.

Protect subscriptions, exclusions, feeds catalog, browser search histories, Homepage account storage and tokens, and per-user settings. Preserve existing publisher URL validation, classification, SSRF limits, import/export and regression coverage. For both agents and the owner, **`./devbuild.sh` is the sole ordinary development test-and-build command**: it already runs `go test ./...` before compiling both Mac binaries. Do not run separate `go test` or `go build` commands before or after it. Normal code changes also require `gofmt` and exact-head GitHub Go CI. If the agent's environment cannot run `./devbuild.sh`, report local test/build as unverified rather than substituting another procedure. A Go CI pass cannot verify live publisher access, Mac Safari, or the tunnel.

## Owner Mac acceptance for runtime PRs

First run `git status --short --branch` and stop if dirty. Use the exact branch from the PR in place of the example:

```bash
cd ~/projects/newsgo
git status --short --branch
# Stop if uncommitted changes are shown.
git fetch origin
git checkout ACTUAL-PR-BRANCH
git pull --ff-only origin ACTUAL-PR-BRANCH
./devbuild.sh
./newsgo.sh
```

`newsgo.sh` restarts NewsGo and its own tunnel, and requires distinct owner authorization. Check feed discovery and subscriptions, exclusions, article reader and export/import, Homepage user/guest handoff, token expiry/revocation, browser history and light/dark UI in Safari. Do not restart for documentation-only PRs. For safe post-FINISH checkout and local branch deletion, use the shared policy's exact clean-tree-guarded template with `~/projects/newsgo`.

## Mandatory Homepage/Guest authorization acceptance — Issue #110

Use `./devbuild.sh` once on a clean `auth/110-verified-suite-access` checkout; it runs Go tests and builds both macOS binaries. No Node.js, additional package manager or new local tooling is needed. This candidate uses the unmerged Homepage `shared/auth` module from Issue #60/PR #61; coordinate Homepage deployment separately and **do not activate a consumer before its matching Homepage recovery route is available**.

In a fresh Safari Private window visit `https://newsgo.themauris.org/` directly. The public bootstrap may load but must not display usable news/reader controls or permit API data until Homepage personal or Guest verification. If not signed in, the Homepage sign-in/Guest choice must appear; the old `#suite_mode=local` cannot bypass it. Confirm legitimate personal and Guest sessions can search feeds, discover publishers, read subscribed feeds and inspect exclusions, but Guest cannot read/write stored personal account settings and only authorized administrators can use admin settings. Direct requests to `/api/personal`, `/api/discover`, `/api/search-feeds`, `/api/activity/ping`, admin routes and unknown `/api/*` must return 401 without a bearer, 200 or normal handler responses with a verified bearer, 401 for expired/revoked bearer and 503 when Homepage profile verification is unavailable. Validate account settings loading and explicit user-triggered Save Changes (no automatic synchronization), including private-browser isolation; confirm no subscriptions, exclusions, histories, tokens or catalog entries are lost. `/health` must remain accessible. No live tokens in screenshots or logs. Run/restart `./newsgo.sh` only with separate owner authorization; it restarts NewsGo and its dedicated tunnel.

## Explicit-save consistency with MauriStocks — Issue #112

**Navigation layout acceptance:** On desktop Safari, Save Changes and its status must share the Find Subscriptions / My Subscriptions / Settings / Diagnostics navigation row. On narrow Safari widths the same flex row may wrap naturally without overlapping controls. The status stays beside Save Changes. Saving, authorization and button styling behavior must not change. The Go UI regression checks the markup and flex rules.


On a clean `auth/110-verified-suite-access` checkout run `./devbuild.sh` **once**; the script performs Go tests and builds the standard macOS binaries. **No Node.js or additional developer tool is required on the Mac.** GitHub CI runs the separate JavaScript behavioral Save Changes test using its existing Node.js runner. The app's internal data schema and Homepage settings endpoint are unchanged. The latest Mac build is Go-only; Node.js behavioral checks are GitHub-CI-only.

In Safari, sign in to Homepage as a personal user and test Subscribe, Unsubscribe, Exclude/Allow Domain, age limit, lead count, classification cache, RSS metadata-only and reading selection. Each change must affect the current view, enable **Save Changes**, and leave browser/account settings uncommitted until clicked. Navigating back to Homepage or refreshing with staged changes must trigger Safari's native unsaved-changes warning; staying must preserve edits. Clicking Save Changes must make one Homepage account PUT and display **Settings saved** only after verified success. Simulate a failed browser storage write (no Homepage PUT, rollback, retry enabled), failed/409 Homepage PUT (no automatic retry; original browser preferences preserved, manual retry enabled), and an account GET that finishes after edits (no overwrite of staged edits). Verify a failed PUT survives reload **for the same account only**. Changing to another personal account must not silently transfer the first account's pending preferences. A newly created personal account must not initialize its storage until an explicit save. Guest settings may be saved locally using the same button but cannot write a personal Homepage record. Confirm search and website-discovery history continue to save and clear automatically, and administrator **Save Server Settings** remains unchanged with its own unsaved form confirmation. Authentication (Issue #110) and saving (Issue #112) are integrated into the same NewsGo PR #111 branch and must be accepted together. Coordinate its live handoff with Homepage PR #61; do not test the superseded PR #113 as the combined candidate.

Do not run `./newsgo.sh` without separate owner authorization: it restarts NewsGo and its tunnel. No new local dependencies; authentication and explicit saving are tested together on the single branch `auth/110-verified-suite-access`.


## Cross-app classification service authentication acceptance (draft)

On separately authorized isolated test processes, set the **same privately generated** `SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN` (64 lowercase hexadecimal digits) in each process's environment. Never check it into Git or paste its value into test reports. Confirm NewsGo article selection and admin classification diagnostics can POST to ScraperGo while preserving existing classification semantics, and that a missing/incorrect service token receives 401. Verify the service token does **not** authorize `/api/read`, `/api/pdf`, `/api/suite/identity`, image, diagnostics or activity endpoints; Guest/personal Homepage bearers retain their normal access. The machine token is not a browser bearer. Redirect targets and non-loopback HTTP must never receive it. Do not alter launchers, private environment files or restart either production process without separate owner approval. When testing, continue using the existing single `./devbuild.sh` for each application and its established version check; no Node.js dependency is added to the Mac.

## Read Articles navigation acceptance

In Safari with a verified personal or Guest Homepage session, choose **My Subscriptions → Read Articles** for one publisher and then choose a headline. Selecting a reading source must not enable **Save Changes** or cause a Safari leave-page warning unless an actual, separate preference edit remains unsaved. The last reading selection remains browser-local operational history. Verify switching to All Subscriptions, unsubscribing the current source, and storage-unavailable fallback. A real preference change must still enable Save Changes and warn before leaving; saving then clears the warning. `scripts/test-explicit-save.cjs` covers the staged-save and non-staging navigation behavior in CI only; no Node.js setup is needed for owner Mac.

## Initial Homepage sign-in / restored-feed acceptance

In Safari, visit NewsGo through Homepage as Guest or personal user while a valid previously selected subscription exists. On first load, NewsGo must not show the transient `Could not refresh feeds: Homepage sign-in or Guest access required` error: it should defer that request until Homepage verifies the session, then populate the selected headlines without another click. Verify an unselected landing page does not fetch all subscriptions, a rejected/revoked bearer does not fetch feeds, and normal manual source changes still work. CI-only `scripts/test-explicit-save.cjs` executes the pre-auth and post-verification refresh paths, with a static assertion that the verification-success branch resumes the pending refresh. Mac builds do not require Node.js.


## Protected thumbnails — Issue #114 draft acceptance

The Safari regression was 190/190 unauthenticated direct ScraperGo thumbnails returning 401 while RSS metadata-only images worked. NewsGo normal reading now uses an authenticated same-origin proxy which forwards the **existing** privately configured `SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN`; no second thumbnail token should be generated or provisioned. ScraperGo PR #148 must be built/deployed in coordination to allow `GET /api/thumbnail` as well as `POST /api/classify` with this single shared credential. Validate that machine credentials cannot reach unrelated ScraperGo APIs, and that NewsGo proxy 502 errors leave an otherwise valid Homepage personal/Guest session signed in. Test real images in Mac Safari across initial view/tab switching, and verify RSS-only mode, failed images, bounded responses and session expiry. Run the existing `./devbuild.sh` **once** on a clean `fix/114-authenticated-thumbnails` checkout; it includes Go tests and both macOS binaries. Do not change private credentials or restart NewsGo/ScraperGo without explicit owner authorization. GitHub CI and live acceptance are distinct.


## Issue #116 candidate — optional thumbnail loopback backend

This phase is intentionally narrower than the original issue: no browser request-scheduler, semaphore/timeout redesign, transport tuning or classification rerouting. The new optional `-thumbnail-backend http://127.0.0.1:8083` affects only NewsGo's server-to-server thumbnail GET requests; `-scrapergo https://scrapergo.themauris.org` remains the public article-reader URL and classification target. Production `newsgo.sh` is unchanged, so the feature is not active by default. To consider activating it, first obtain separate owner permission for a runtime/launcher change and restart. Never manually change the shared token. The single ordinary owner Mac test/build command is `./devbuild.sh`; CI-only browser checks and Go tests should pass on the exact PR head. Owner Safari acceptance should target one normal All Subscriptions load, a refresh and a subscription switch, recording loaded/failed/pending thumbnail counts and any 502 response origin. The brief owner test cannot replace production tunnel/log correlation.
