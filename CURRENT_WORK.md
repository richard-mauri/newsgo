# Current work — NewsGo

## Issue #89 — optional Homepage suite activity reporting (draft)

- Authoritative merged base: GitHub `main`, commit `071505af16a67fd0282ed291e0036ce9df5d451b` (merged duplicate-audit utility on November 11, 2026). Work is scoped to [Issue #89](https://github.com/richard-mauri/newsgo/issues/89), branch `feat/89-suite-activity`, draft [PR #90](https://github.com/richard-mauri/newsgo/pull/90). Prior Issue #87 was merged on `main`; this task does not touch that CLI or ScraperGo.
- Scope: optional same-origin per-tab anonymous/unverified presence, in-memory rolling visitor windows, in-flight `/api/personal` work counter, and server-to-server aggregate reports to Homepage's Issue #28 / draft PR #30 collector. Requires explicit `NEWSGO_ACTIVITY_URL` and private `HOMEPAGE_ACTIVITY_KEY_NEWSGO` (minimum 32 characters). Without configuration, NewsGo remains unchanged.
- Security: no browser access to service key; no URLs, publisher content, sessions, accounts or identities in telemetry. Signed-in browser counts are currently **unverified** and included in anonymous; no claim of reliable signed-in visitor classification. Only HTTPS and loopback HTTP collector URLs, no redirects. ScraperGo work detached from `/api/personal` requires an independent reporter.
- Validation: GitHub Actions Go CI #586 passed on the earlier reporter head `c8088d7` (Go tests and build). The new two-service integration harness and its syntax/safety unit test are awaiting exact-head CI. New unit tests cover disabled/unsafe configuration, same-origin ping, rolling windows, in-flight work and report shape. Mac/Safari, live Homepage collector and end-to-end telemetry **not tested**. No release, semantic version update or restart. Owner approval required for runtime acceptance.
- Next: verify exact-head CI, then owner runs the new isolated two-service Safari harness to see NewsGo change to Reporting in the Homepage dashboard. The script uses only disposable test data, ports 8788 and 8085, and a temporary server-only shared key. Later add server-verified signed-in classification and other application instrumentation. Do not merge until approved.

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

For pre-merge integration use `TEST_ROOT="$TEST_ROOT" bash scripts/activity-integration-test.sh` from a separate NewsGo worktree, after stopping ONLY the previous isolated Homepage 8788 process. The guarded harness builds NewsGo, confirms 8788/8085 are available, starts the existing test Homepage with preserved disposable admin/history, uses an isolated NewsGo settings directory and one temporary private matching key, and stops its own two processes on Ctrl-C. The owner then checks Safari at 8085 and the Homepage admin panel at 8788 after 75 seconds. Mac integration and Safari acceptance are still unverified.
