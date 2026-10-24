# NewsGo — current development state

**As of 2026-10-24. Update this file when work or acceptance changes. This is not a release-history log.**

## Authority and candidate
- Repository: `richard-mauri/newsgo`.
- **Development source:** GitHub `dev`; **release source:** `main` (owner controls promotion).
- **Open draft PR:** [#10 — Admin diagnostics, Exit Admin Mode and build identification](https://github.com/richard-mauri/newsgo/pull/10), `dev` → `main`. Related [issue #9](https://github.com/richard-mauri/newsgo/issues/9).
- Dashboard semantic version remains **v0.2.58**. Mac-local **Build NNNNNN** increments with `./devbuild.sh`; numbers vary by Mac and are not committed.
- **No release or merge authorized.** The last checked GitHub PR reported **not mergeable**; inspect whether `dev`/`main` require conflict resolution before seeking promotion. PR includes previous `dev` commits, not solely the latest feature changes.

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
- Check whether the draft PR remains non-mergeable and inspect/resolve branch conflicts **without merging, force-pushing or rewriting history without owner authorization**.
- After final changes, verify latest-head GitHub Actions, refresh the PR title/body to match work and evidence, and obtain any remaining Mac/Safari acceptance. Ensure current documentation and PR status agree.
- Documentation consolidation into README / AGENTS / ARCHITECTURE / CURRENT_WORK is being performed in this task; note its final commit/CI result before handoff.
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
