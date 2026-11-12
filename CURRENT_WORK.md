# Current work — NewsGo

## Issue #89 — optional Homepage suite activity reporting (draft)

- Authoritative merged base: GitHub `main`, commit `071505af16a67fd0282ed291e0036ce9df5d451b` (merged duplicate-audit utility on November 11, 2026). Work is scoped to [Issue #89](https://github.com/richard-mauri/newsgo/issues/89), branch `feat/89-suite-activity`, draft [PR #90](https://github.com/richard-mauri/newsgo/pull/90). Prior Issue #87 was merged on `main`; this task does not touch that CLI or ScraperGo.
- Scope: optional same-origin per-tab anonymous/unverified presence, in-memory rolling visitor windows, in-flight `/api/personal` work counter, and server-to-server aggregate reports to Homepage's Issue #28 / draft PR #30 collector. Requires explicit `NEWSGO_ACTIVITY_URL` and private `HOMEPAGE_ACTIVITY_KEY_NEWSGO` (minimum 32 characters). Without configuration, NewsGo remains unchanged.
- Security: no browser access to service key; no URLs, publisher content, sessions, accounts or identities in telemetry. Signed-in browser counts are currently **unverified** and included in anonymous; no claim of reliable signed-in visitor classification. Only HTTPS and loopback HTTP collector URLs, no redirects. ScraperGo work detached from `/api/personal` requires an independent reporter.
- Validation: GitHub Actions Go CI #585 passed on code and test head `9a937a5675d045c5618e848977cf8d2cf076236f` (Go tests and build). This documentation update requires a new exact-head check. New unit tests cover disabled/unsafe configuration, same-origin ping, rolling windows, in-flight work and report shape. Mac/Safari, live Homepage collector and end-to-end telemetry **not tested**. No release, semantic version update or restart. Owner approval required for runtime acceptance.
- Next: run exact-head CI; verify NewsGo telemetry against an isolated Homepage collector, validate Safari browser-only behavior, then separately instrument trusted signed-in classification and remaining service applications. Do not merge until approved.

## Owner Mac acceptance (production restart ONLY after separate owner authorization)

```sh
cd ~/projects/newsgo
git fetch origin
git checkout feat/89-suite-activity
git pull --ff-only origin feat/89-suite-activity
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

For pre-merge testing, prefer a temporary test build with `-listen 127.0.0.1:8085`, an isolated `-settings-dir`, a disposable secret and `NEWSGO_ACTIVITY_URL=http://127.0.0.1:8788/api/admin/activity/ingest`; do not run the launcher, touch production data or restart another app. The test Homepage instance must be restarted with its own matching server-only `HOMEPAGE_ACTIVITY_KEY_NEWSGO` environment variable. Validate the NewsGo row becomes Reporting in the administrator panel only once it receives accepted reports; instrumented visitor counts may be anonymous/unverified.
