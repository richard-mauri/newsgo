# Active work — Header accent consistency

Baseline: GitHub `main` after merged PR #42. Issue [#44](https://github.com/richard-mauri/newsgo/issues/44), branch `fix/44-header-accent`. Candidate is a UI-only change; no release or merge authorized.

## Scope
- Match MauriStocks' accent-colored authenticated username: NewsGo's Homepage session line uses `var(--accent)`, 0.9rem and font-weight 650 **only when authenticated**. Signed-out or unavailable connection status remains muted.
- Match MauriStocks' accent-colored Tanner name using NewsGo's existing `var(--accent)` and 700 font weight. Preserve Tanner image sizing, layout and responsive behavior.
- No changes to auth verification, account persistence, legacy settings, tokens, subscriptions, or private storage.

## Acceptance
GitHub CI must run `go test -count=1 ./...` and `go build ./...` on final head. Owner to review Safari appearance with signed-in and signed-out state; confirm Tanner label and authenticated identity are teal. No local Mac acceptance performed by assistant.

## Mac validation
```bash
cd ~/projects/newsgo
git fetch origin
git checkout fix/44-header-accent
git pull --ff-only origin fix/44-header-accent
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

PR remains unmerged until owner approval.
