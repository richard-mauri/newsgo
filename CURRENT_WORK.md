# Current work — NewsGo

## Active draft — main navigation hover highlight (Issue #56, November 9, 2026)

- **Issue:** [#56](https://github.com/richard-mauri/newsgo/issues/56); branch `fix/56-navigation-hover-highlight` based on authoritative merged `main` at `e12609f65193335264c268649c37cd707f9d91c3`. Draft [PR #78](https://github.com/richard-mauri/newsgo/pull/78) against `main`.
- **Scope:** Apply Homepage Manage Tokens-style accent hover highlighting to the border and text of unselected buttons in the main NewsGo navigation row. Provide matching `:focus-visible` color and accessible outline, while retaining the current cyan-filled active-button styling, dark/light theme colors, and all navigation behavior.
- **Regression:** `internal/dashboard/navigation_hover_test.go` checks the navigation hover, keyboard focus, and active styles. GitHub Actions [run #42644378651](https://github.com/richard-mauri/newsgo/actions/runs/42644378651) passed Go tests and build on code head `8d1c2fe6adc73372286841ed758f5cd7ca3b5770`. Final documentation-head CI and owner Mac/Safari visual acceptance are pending.
- **Owner Mac acceptance:** Run the commands below, then in Safari hover each visible button in the main NewsGo navigation row, including Find Subscriptions, My Subscriptions, Reading Preferences and Diagnostics. Confirm text and border turn cyan without changing background, active navigation remains legible, and keyboard Tab focus is clearly visible in both light and dark themes. Check click/navigation still works.
- **Merge boundary:** No merge, release, tag, or modification to `main` without explicit owner approval.

```bash
cd ~/projects/newsgo
git fetch origin
git checkout fix/56-navigation-hover-highlight
git pull --ff-only origin fix/56-navigation-hover-highlight
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

## Other work (independent)

- Issue #76 / PR #77 was merged into `main` on November 9, 2026 (merge `e12609f65193335264c268649c37cd707f9d91c3`).
- Issue #57 / [PR #69](https://github.com/richard-mauri/newsgo/pull/69) is a separate Homepage session-renewal effort.
