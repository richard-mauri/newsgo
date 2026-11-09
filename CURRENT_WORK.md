# Current work — NewsGo

## Active draft — remove headline exclusion buttons (Issue #74, November 9, 2026)

- **Issue:** [#74](https://github.com/richard-mauri/newsgo/issues/74); branch `fix/74-remove-headline-exclusion-buttons`; draft [PR #75](https://github.com/richard-mauri/newsgo/pull/75) against authoritative merged `main` (base `6fff6daa8b555118a74db61a0c85db8eb66943e5`).
- **Change:** Remove both `Exclude This Domain` and `Exclude articles from <domain>` action buttons from all article headline cards (Lead Stories and Latest Headlines). Article links, original publisher links, dates and thumbnails remain. Existing exclusion controls in My Subscriptions, Find Subscriptions and Excluded Domains remain available, with the same confirmation, saved exclusions and profile sync.
- **Regression:** `internal/dashboard/headline_exclusion_controls_test.go` protects against reintroduction of the two headline controls and checks the subscription exclusion handlers remain. Go CI and owner Mac/Safari runtime acceptance are pending.
- **Owner Mac acceptance:** `cd ~/projects/newsgo`, `git fetch origin`, `git checkout fix/74-remove-headline-exclusion-buttons`, `git pull --ff-only origin fix/74-remove-headline-exclusion-buttons`, `go test -count=1 ./...`, `./devbuild.sh`, `./newsgo.sh`. In Safari open All Subscriptions, read Lead Stories and Latest Headlines, verify neither contains an Exclude button, then use My Subscriptions to exclude a domain and verify confirmation and that the Excluded Domains panel still permits allowing it. Verify subscription selections and preferences survive a refresh.
- **Merge boundary:** Draft PR #75. Do not merge, tag, release or modify `main` without explicit owner approval.

## Other work (independent)

- Issue #55 / PR #73 for subscription failure details was merged into `main` on November 9, 2026 (merge `6fff6daa8b555118a74db61a0c85db8eb66943e5`). It is no longer the active draft.
- Issue #57 / [PR #69](https://github.com/richard-mauri/newsgo/pull/69) addresses Homepage session renewal and requires separately coordinated Homepage changes and acceptance.
