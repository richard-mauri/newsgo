# Current work — NewsGo

**Main baseline:** GitHub merged `main` at `3aeeb3426c1eb896dbda6c719c3986a2456a8ac6`. Active [Issue #110](https://github.com/richard-mauri/newsgo/issues/110), branch `auth/110-verified-suite-access`, draft PR targeting `main` (verify number and head on GitHub). Older Issue #108/PR #109 was completed in merged `main`; shared banner and button contracts remain unchanged.

**Scoped implementation:** Adopt immutable-pinned Homepage `shared/auth` module, with an outer HTTP gate on every NewsGo `/api/*` endpoint, trusting only live Homepage profile verification for an app-specific bearer. An anonymous/revoked caller receives 401; upstream verification failures fail closed (503), and verified personal/admin/Guest access remains available while existing admin routes retain additional authorization. The initial HTML and same-origin assets are only a recovery bootstrap. Interactive reading is hidden until verified; browser API requests attach the handoff bearer; the old local-mode fragment cannot authorize NewsGo. No consumer-specific session issuer, UI styling change, data migration or new local dependency.

**Homepage integration:** Homepage [Issue #60](https://github.com/richard-mauri/homepage/issues/60) / draft [PR #61](https://github.com/richard-mauri/homepage/pull/61) supplies `shared/auth` and fail-closed recovery. Shared verifier injection is adequate; no framework change required this phase. Homepage PR #59 (Safari CORS) is a separate concern. NewsGo deployment must be coordinated with matching Homepage behavior, and owner real personal/Guest/private browsing and account-storage acceptance is outstanding.

**Data safety:** Owner Mac-local uncommitted and unpushed changes inaccessible; no live subscriptions, exclusions, publisher catalog, search histories, account settings, tokens, user identities, browser storage, activity configuration, data files, build counters, processes, tunnels or DNS modified. No new local tooling.

**Documentation:** Read current Homepage merged policy v1.1.0 and NewsGo AGENTS, README, ARCHITECTURE, DEVELOPMENT, CURRENT_WORK; update the latter four to describe the active security contract. AGENTS and policy unchanged.

**Verification:** GitHub CI must pass on exact latest draft PR head, including all existing feed/classification tests, new auth role/route/failure tests and executable build. Local developer Mac build and Safari/live integration remain unverified. No separate `go test` or `go build` beyond `./devbuild.sh` for owner.

**Mac handoff:** `cd ~/projects/newsgo`; `git status --short --branch` and stop if dirty; `git fetch origin`; `git checkout auth/110-verified-suite-access`; `git pull --ff-only origin auth/110-verified-suite-access`; `./devbuild.sh`. `./newsgo.sh` is a separate owner-authorized restart of NewsGo and its tunnel.

**Promotion:** Issue open, PR draft. No FINISH, merge, Issue closure, branch deletion, restart, deployment, tag, release or version/builder-counter change authorized.
