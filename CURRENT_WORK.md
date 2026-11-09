# Current work — NewsGo

## Active draft — declutter the main briefing (Issue #76, November 9, 2026)

- **Issue:** [#76](https://github.com/richard-mauri/newsgo/issues/76); branch `fix/76-diagnostics-declutter`; draft [PR #77](https://github.com/richard-mauri/newsgo/pull/77) targets merged `main` baseline `776677d1c184ea7f2fd26626f69e7163a83c5c06`.
- **Scope:** Move the existing `Filtered / unprocessed articles` details audit from the reading page to Diagnostics, and move the `Excluded Domains` button out of the main navigation into Diagnostics. The existing exclusion manager is nested within Diagnostics; Close returns to Diagnostics. The filtered tabs, per-feed filtering advice and unsubscribe confirmation, saved exclusions, confirmation dialogs, Homepage profile sync, reading selection, and raw feed diagnostics remain.
- **Tests:** Added `internal/dashboard/diagnostics_layout_test.go` to verify layout and navigation contracts. GitHub Actions Go test/build on final PR head and owner Mac/Safari UI acceptance are pending.
- **Owner Mac/Safari acceptance:** Run the commands below, then load an All Subscriptions briefing and a selected subscription in Safari. Confirm the main page shows no filtered audit panel and the top navigation has no Excluded Domains button. In Diagnostics confirm filtered groups/counts/selection-specific recommendation remain functional. Open Excluded Domains, exclude/allow a domain (including affected-subscription confirmation), and verify Close returns to Diagnostics. Verify raw feed viewer, diagnostics counters, preferences and Homepage account persistence remain intact.
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
