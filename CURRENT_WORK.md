# NewsGo — current development state

## Issue #14 draft candidate (2026-10-24)

- Branch `fix/14-search-exclusion-feedback` from `main` at `08ee9ef`; scoped to Issue #14, search candidate exclusion feedback.
- The search candidate buttons already invoke `toggleDomain`, but a successful exclusion previously refreshed the exclusions and subscriptions panels without updating visible search cards. Added `refreshSearchExclusionCards()` to update all currently visible candidates, labels, exclusion badges, and card state. Following owner Safari feedback, excluded candidates also disable Subscribe until allowed again, with guards against in-flight subscription validation.
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
