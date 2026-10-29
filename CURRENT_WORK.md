- Added scoped timing instrumentation: total refresh, batch requests, summed per-feed fetch/wait and classification durations, render time and mean completed thumbnail event latency. Cache hit rates and image failure-stage timing are not implemented. Final-head CI and owner Safari tests pending.

# Issue #36 — admin Diagnostics and classification-grouped filtered articles (2026-10-29)

- Issue: https://github.com/richard-mauri/newsgo/issues/36. Branch `feat/36-admin-diagnostics-classification` from authoritative `main`; scoped draft PR planned, not merged or released.
- Repro supplied by owner: All Subscriptions repeatedly shows zero thumbnails after completing, while an individual subscription shows thumbnails.
- Follow-up 2026-10-29: filtered entries now use one accessible tab per classification/reason, with counts; unprocessed is its own tab. Restored original eager thumbnail loading (no unauthorized lazy behavior change). Safari acceptance pending.
- Candidate: Admin-only Diagnostics panel reports latest completed browser read totals, actual image loading outcomes, and per-feed errors. Non-featured thumbnail images are lazy-loaded instead of all eager; failed image status is recorded before the image is removed. Existing filtered view is grouped by classification/exclusion reason with unprocessed kept separate.
- Correction at `50af191`: initial diagnostics image instrumentation omitted `img.src=a.image`, producing empty image rectangles in Safari; restored. Original draft-head Go CI #202 passed but did not catch this browser behavior. New-head CI and Mac/Safari acceptance are pending.
- No ScraperGo changes, no server persistence or user data migrations. Go CI, Mac build, all-subscriptions Safari test and image outcome verification remain pending. Owner must approve merge, tag and release.

---

# Integrated Mac/Safari acceptance candidate — Issues #31 and #33 (2026-10-26)

- Branch `integration/31-33-cnnmoney-acceptance` combines charset parsing and Subscribe Failure dialog from PR #32 with the >90% filtered-article unsubscribe recommendation from PR #34.
- The sole numeric condition is filtered / processed articles > 90%, regardless of age or other filtering reason, with explicit confirmation and no automatic removal.
- Integration-only draft PR targets main; neither source PR nor main merged or changed. Final-head CI and owner Safari acceptance required before promotion.

---

# Active candidate — Issue #31 CnnMon subscription failure (2026-10-26)

- Baseline: GitHub `main`; branch `fix/31-charset-subscribe-feedback`; issue https://github.com/richard-mauri/newsgo/issues/31.
- Fix XML decoding for ISO-8859-1 declarations in RSS/Atom/RDF using a constrained charset reader; retain rejection for unrecognized charsets.
- Failed Subscribe becomes **Failure** with a click-to-view themed dialog and Retry control. Preserve the existing bottom-of-panel message for context; do not change stored subscriptions or Server Sync formats.
- Added targeted charset parser and subscription UI-contract tests. GitHub Actions exact-head outcome and Mac/Safari acceptance pending. Owner must approve merge or release.

---

# Active candidate — Issue #26 selected-subscription unsubscribe recommendation

- Baseline: GitHub main; branch `feat/26-filtered-subscription-unsubscribe`; issue https://github.com/richard-mauri/newsgo/issues/26.
- UI-only: selected subscription name continues to appear in the existing reading title, with a new recommendation after it only if >90% of at least 10 classified articles were filtered. All-subscriptions selection, incomplete feed fetches and unprocessed entries suppress the recommendation.
- Only classification-filtered entries count; age limits and user exclusions do not influence the percentage.
- Explicit confirm before unsubscribing; reuse saved subscriptions/Server Sync and return to All. No automatic unsubscribe; no ScraperGo, filtering, feed, cache or backend changes.
- Owner Safari feedback on Build 000033: restored saved feed required clicking Read Articles, the unsubscribe control looked unstyled, and native confirmation didn't match NewsGo. Updated to avoid final startup render overwriting sync restoration, retain sync's own one-time saved-selection restoration without a duplicate feed request, and use existing sync-choice dialog and confirmation button classes. Static coverage extended. Latest-head CI and owner Mac/Safari retest pending. Other draft PRs untouched; merge and release require owner authorization.

# Active draft work — Issue #18 (2026-10-25)

- Baseline: `main`; branch `fix/18-configurable-feed-search-timeout`; no merge/release authorization.
- Scope: expose feed-directory search HTTP timeout as a persistent admin setting, retaining the independently configurable classification timeout; backward-migrate existing server-settings.json records without resetting other values.
- Admin-observed symptom: classification diagnostic showed effective 7 seconds after a prior attempted 30-second admin save. Root cause of live reset not yet verified; startup/GET/PUT source already persists the classification setting. Added 30-second restart regression rather than blindly changing the default.
- Hardcoded operational values audited: feed search 7s (fixed in branch), global classification semaphore 16, feed fetch slots 4, feed transport dial 5s/header 8s/overall 12s, failed feed cache 20s, classification positive cache 30m, Feedly query cache 10m and 3s rate limit. Some represent safety boundaries and need separate design review before becoming admin controls.
- Owner accepted revised responsive Admin Settings layout, but requested an unsaved-changes warning before Close or Exit Admin Mode. Added three-way save/discard/keep-editing modal, plus protection when Reload is clicked and a native browser leave warning. Save occurs before leaving; failed saves leave the panel open. Mac/Safari acceptance of this new behavior is pending.
- Next: verify exact-head CI and owner Mac/Safari Save/Discard/Keep Editing scenarios, plus saved-settings restart; no merge until approval.

---

# NewsGo — current development state

## Unified exclusions candidate (Issues #14 and #16, 2026-10-24)

- Owner explicitly authorized widening draft PR #15 to resolve unified exclusions, rather than opening a second PR.
- Combined Excluded Sites and Excluded Article Domains into a single visible Excluded Domains manager. Both legacy storage fields and sync data remain readable; opening an old profile does not delete subscriptions.
- A new exclusion checks affected feed subscriptions and asks for explicit confirmation naming them; on confirmation, it stores the exclusion and removes the affected subscriptions. Cancel preserves both. Allowing does not restore deleted subscriptions.
- Article destinations and feed publishers both observe the effective exclusions. Migration behavior requires owner Mac/Safari acceptance, including local preferences, imports, server synchronization, and domain submatch cases.
- PR #15 remains draft; CI exact head and Safari acceptance pending. Do not merge or release without owner approval.

## Issue #14 draft candidate (2026-10-24)

- Branch `fix/14-search-exclusion-feedback` from `main` at `08ee9ef`; scoped to Issue #14, search candidate exclusion feedback.
- The search candidate buttons already invoke `toggleDomain`, but a successful exclusion previously refreshed the exclusions and subscriptions panels without updating visible search cards. Added `refreshSearchExclusionCards()` to update all currently visible candidates, labels, exclusion badges, and card state. Following owner Safari feedback, excluded candidates also disable Subscribe until allowed again, with guards against in-flight subscription validation. For My Subscriptions, excluded entries now show their Excluded status without a redundant exclusion button; re-allowing remains available in Excluded Sites. Owner screenshots revealed a second mismatch: `reg.bom.gov.au` is excluded while a subscribed feed displays `www.bom.gov.au`. The current candidate now groups sibling `*.bom.gov.au` hosts by their `bom.gov.au` agency domain, without broadly collapsing unrelated public suffixes; existing stored exclusion values are preserved and reversible.
- Added a source-contract regression test. GitHub CI and Mac/Safari acceptance must be verified at final PR head; no runtime data, subscriptions or exclusions are intentionally changed by this patch.
- This candidate is not merged; owner must review Safari behavior before approving promotion. The established semantic version remains unchanged.

## Current documentation-only branch (2026-10-24)

- Working branch: `docs/branch-lifecycle-20261024` → `main`. Scope: document safe post-merge branch deletion, prevent accumulated stale heads, preserve owner-controlled releases.
- Added `DEVELOPMENT.md`; updated `AGENTS.md` to delegate workflow details.
- Existing historical branches are **not deleted**. ScraperGo currently has 66 remote branches, 63 matching heads of merged PRs; bulk cleanup requires a separately approved safety audit.
- No runtime code or user state changes. Check final-head CI before handoff; owner approval required to merge this documentation PR.
- Previous PR-cycle notes below are historical context and should be reconciled with live GitHub state, not treated as current open PRs.


**As of 2026-10-24. Update this file when work or acceptance changes. This is not a release-history log.**

## Authority and candidate
- Repository: `richard-mauri/newsgo`.
- **Authoritative baseline for future work:** `main`. New work must use issue-specific branches from current `main`, with one draft PR per issue. Legacy `dev` remains only to reconcile existing PR #10; do not reuse it for new features.
- **Merged PR:** [#10 — NewsGo v0.2.58 fixes](https://github.com/richard-mauri/newsgo/pull/10), `dev` → `main` (merged 2026-10-24). Related [issue #9](https://github.com/richard-mauri/newsgo/issues/9).
- Dashboard semantic version remains **v0.2.58**. Mac-local **Build NNNNNN** increments with `./devbuild.sh`; numbers vary by Mac and are not committed.
- **PR #10 has merged into `main`.** Its earlier integration conflicts were resolved in legacy `dev`; all future work starts from current `main`. Preserve its tested behavior and four-document structure, check CI and obtain explicit final owner merge approval.

## Implemented on dev and observed in Mac/Safari
- Admin-only expandable classification diagnostics for filtered entries, backed by a server-authorized on-demand classification endpoint. Earlier diagnostic returned HTTP 400 for the Seahorse Magazine feed link ending `#new_tab`.
- Root cause identified: NewsGo's URL validation disallowed fragments, so normal feed classification returned `unknown` before reaching ScraperGo and the new diagnostic rejected the URL. Both paths now strip browser-only fragments before validating/classifying while preserving original navigation URLs and strict filtering.
- **Owner-confirmed live result:** Sailing Anarchy's Seahorse Magazine story **“Control from the cockpit”** appears in the headline grid with thumbnail; screenshot shows **10 visible headlines**, **0 filtered / 0 unprocessed** where previously five were filtered.
- Server Sync's up-to-date connection path now refreshes the subscription manager/count. **Owner-confirmed live result:** **My Subscriptions (11)** and feed selection restored, replacing an inconsistent display of zero subscriptions while discovery still marked a feed subscribed.
- **Exit Admin Mode owner-confirmed working** in Safari. It uses existing Server Sync disconnect and does not revoke the server token.
- Six-digit Mac-local build numbering is implemented in both architectures, CLI `-version`, startup log, launcher and dashboard.
- Seven obsolete root-level historical/evaluation Markdown files were deleted from `dev`; they remain recoverable through Git history.

## Outstanding work and acceptance
- **Exit Admin Mode styling correction implemented on dev:** the button now sits inside the existing admin form and inherits the same visual styling as Save Server Settings and Reload. This preserves its owner-accepted click handler and disconnect behavior. **Safari visual acceptance remains pending.**
- Legacy `dev` ↔ `main` divergence is reconciled at merge commit `25892f6` with the unchanged tested `dev` tree and both histories retained. Revalidate exact-head CI and Mac/Safari if needed. **Do not promote PR #10 until explicit owner approval.**
- After final changes, verify latest-head GitHub Actions, refresh the PR title/body to match work and evidence, and obtain any remaining Mac/Safari acceptance. Ensure current documentation and PR status agree.
- Four-file documentation consolidation is complete; `AGENTS.md` now mandates separate issue branches and focused draft PRs. Update exact reconciliation commit and CI outcome before handoff.
- Existing runtime data, tokens, feed subscriptions, exclusions and Cloudflare setup must remain untouched.

## Verified automation and live checks
- At previous code head `adb3982`, GitHub Actions passed for the fragment/classification regression tests.
- At documentation-cleanup head `51d81b2`, GitHub Actions passed ([run 40104223344](https://github.com/richard-mauri/newsgo/actions/runs/40104223344)).
- The owner independently verified the classification/subscription fixes and Exit Admin Mode behavior in Safari. **The Exit Admin Mode styling correction is committed; Safari visual acceptance has not yet been performed.** Passing CI does not establish visual acceptance.

## Mac reproduction and next acceptance
```bash
cd ~/projects/newsgo
git fetch origin
git checkout dev
git pull --ff-only origin dev
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

In Safari: confirm the current build label and 11 subscriptions; open Sailing Anarchy and verify “Control from the cockpit” is in headlines, not the filtered audit. Connect an admin token via Server Sync, inspect a filtered article if present, then use Exit Admin Mode and reconnect. Compare Exit Admin Mode appearance with Save Server Settings and Reload, and verify disconnect still works. **Do not reset/import preferences or overwrite user profiles merely to repair UI counts.**

For rules see [AGENTS.md](AGENTS.md), for operation see [README.md](README.md), and for design see [ARCHITECTURE.md](ARCHITECTURE.md).
