# Current work — NewsGo

## Issue #83 — grouped filtered audit domain exclusion (November 10, 2026)

- **Merged baseline:** GitHub `main` at `a005d76644ad107fd51317b839363271c297d2de` (Account Sync cleanup PR #82 merged). New issue [#83](https://github.com/richard-mauri/newsgo/issues/83); branch `feat/83-filtered-domain-exclusion`.
- **Request:** Diagnostics screenshot shows Unknown (75) filtered articles apparently from one Investor's Business Daily feed. Determine normalized publisher domain from article URLs and offer a single domain-level Exclude control instead of one per article.
- **Implementation:** Each classification tab groups entries by `publisherScope(article.url)`, displays domain and count, and offers a single Exclude domain button per non-excluded, valid domain. It reuses `setUnifiedDomain` with the existing confirmation, local/account exclusion persistence and refresh. No automatic exclusion, altered classifier decisions, shortened per-item audit, or admin diagnostics changes. Invalid domains get no exclusion action.
- **Validation:** Added grouped audit Go contract tests, JavaScript syntax checked. GitHub CI exact head pending; Mac/Safari/live classification and account persistence acceptance unverified.
- **Authorization:** Owner requested feature development, not merge, restart, release, or Cloudflare change. A focused draft PR will target main.

## Owner Mac acceptance
```sh
cd ~/projects/newsgo
git status --short --branch
git fetch origin
git checkout feat/83-filtered-domain-exclusion
git pull --ff-only origin feat/83-filtered-domain-exclusion
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

Restart only with owner approval. In Safari, visit Diagnostics → Filtered / unprocessed articles, select Unknown and verify that many entries sharing one domain produce one Exclude control; confirm cancellation leaves data intact and confirmation removes matching feed subscriptions only with the existing warning. Validate unknown URL handling, per-item admin diagnostics, local and authenticated account exclusion persistence, and unaffected classification counts.
