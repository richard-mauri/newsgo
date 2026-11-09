# Current work — NewsGo

## Active draft — declutter the main briefing (Issue #76, November 9, 2026)

- **Issue:** [#76](https://github.com/richard-mauri/newsgo/issues/76); branch `fix/76-diagnostics-declutter`; draft [PR #77](https://github.com/richard-mauri/newsgo/pull/77) targets merged `main` baseline `776677d1c184ea7f2fd26626f69e7163a83c5c06`.
- **Scope:** Move the existing `Filtered / unprocessed articles` details audit from the reading page to Diagnostics, and replace the main-navigation `Excluded Domains` button with a collapsed `Excluded Domains (N)` disclosure inside Diagnostics, matching the `Filtered / unprocessed articles` disclosure. The manager expands inline and has no separate Close button. Raw Feed Responses is now a third collapsed native disclosure matching the other two panels; its per-feed sub-disclosures, original URLs, Copy/Download controls, and exact preserved feed bytes remain unchanged. The filtered tabs, per-feed filtering advice and unsubscribe confirmation, saved exclusions, confirmation dialogs, Homepage profile sync, and reading selection remain.
- **Tests:** Added `internal/dashboard/diagnostics_layout_test.go` to verify layout and navigation contracts. Regression tests cover the native disclosures, hidden main navigation, exclusion-management handlers, and filtered audit placement. The new raw-feed outer disclosure and existing per-feed controls have regression assertions. GitHub Actions at the updated draft head and owner Mac/Safari acceptance remain pending.
- **Owner Mac/Safari acceptance:** Run the commands below, then load an All Subscriptions briefing and a selected subscription in Safari. Confirm the main page shows no filtered audit panel and the top navigation has no Excluded Domains button. In Diagnostics confirm filtered groups/counts/selection-specific recommendation remain functional. Expand Excluded Domains, exclude/allow a domain (including affected-subscription confirmation), collapse it using the summary, and verify the saved excluded count remains correct. Expand Raw Feed Responses, verify the feed count and per-feed disclosures, original source links, Copy and Download Raw Feed actions, and collapse it again. Verify diagnostics counters, preferences and Homepage account persistence remain intact.
- **Merge boundary:** Do not merge, release, tag or alter `main` without owner approval.

```bash
cd ~/projects/newsgo
git fetch origin
git checkout fix/76-diagnostics-declutter
git pull --ff-only origin fix/76-diagnostics-declutter
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

## Other work (independent)

- Issue #74 / PR #75 was merged into `main` on November 9, 2026 (merge `776677d1c184ea7f2fd26626f69e7163a83c5c06`).
- Issue #57 / [PR #69](https://github.com/richard-mauri/newsgo/pull/69) remains a separate Homepage session renewal effort.
