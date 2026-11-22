# Current work — NewsGo

**Verified merged main at session start:** `cc86480caafbf3299b1064c380e0d7e5bd17ff3c`. **Current scoped issue:** [#114](https://github.com/richard-mauri/newsgo/issues/114), **branch:** `fix/114-authenticated-thumbnails`. This is a new draft candidate based on current `main`; historical auth branch `auth/110-verified-suite-access` was not present and must not be used. ScraperGo merged main inspected at `e22b2e7427a6884074d10e72b24b47fe544aded8`. The documented earlier auth PRs #111 and #146 are historical references, not current candidate heads.

## Scope and evidence

Owner Safari showed 196 accepted NewsGo articles, **190 thumbnail failures and 0 successes** on both cold and cached reads. A specific ScraperGo `GET /api/thumbnail` yielded HTTP 401 with no Authorization header. RSS metadata-only thumbnails rendered successfully. NewsGo used direct `img.src` pointing to protected ScraperGo; ScraperGo's server gate requires a Homepage bearer for all media endpoints and grants the existing service token solely to `POST /api/classify`.

This NewsGo draft makes normal reading thumbnail URLs same-origin `/api/thumbnail` and loads them through `authorizedNewsgoFetch` with the verified NewsGo child bearer; browser blob URLs are revoked on rerender, requests aborted and diagnostics retained. NewsGo's server proxy validates requests, uses only the configured ScraperGo destination (HTTPS or explicit loopback HTTP), refuses redirects and non-images, limits body size and concurrency, and requires **separate** private `SUITE_NEWSGO_SCRAPERGO_THUMBNAIL_TOKEN`. This credential is **not** the existing classification-only token. RSS mode is unchanged.

## Cross-repository blocker

ScraperGo does not yet accept a distinct machine-thumbnail credential on `GET /api/thumbnail`. A **separately authorized, scoped ScraperGo Issue/PR** is required before integration acceptance or deployment. Do not relax ScraperGo's API gate, reuse the classification token, place credentials in browser state, or restart services as part of this NewsGo draft.

## Verification and safety

GitHub-only code and documentation changes. Added `internal/dashboard/thumbnail_proxy_test.go` for method/credential/request validation, image passthrough, redirect and MIME failures, and protected URL construction. **GitHub Go CI run #735 passed on code/test head `5036a39411ca0a2b165f9d66dc908522b01abd28`**, including Go tests/build, embedded JS syntax and CI-only browser checks. A subsequent documentation head requires its own exact-head CI confirmation. **Local `./devbuild.sh`, actual Mac/Safari image acceptance, ScraperGo credential provisioning and cross-service integration remain unverified.** Owner Mac uncommitted/unpushed files, private tokens, browser storage, subscriptions, accounts, tunnel, process and build-counter state cannot be inspected. No local files, tokens, processes, tunnels, data stores or production deployment changed.

Read current Homepage merged `SUITE_DEVELOPMENT_POLICY.md` v1.1.0 and NewsGo `AGENTS.md`, `README.md`, `ARCHITECTURE.md`, `DEVELOPMENT.md`, `CURRENT_WORK.md`. Updated README, ARCHITECTURE, DEVELOPMENT, CURRENT_WORK for #114; AGENTS and Homepage policy unchanged. CI outcome and draft PR link should be updated here when verified.

## Owner Mac handoff (not yet authorized for production launch)

```bash
cd ~/projects/newsgo
git status --short --branch
# Stop if uncommitted work is present.
git fetch origin
git checkout fix/114-authenticated-thumbnails
git pull --ff-only origin fix/114-authenticated-thumbnails
./devbuild.sh
./dist/newsgo-darwin-arm64 -version
# ./newsgo.sh restarts NewsGo and Cloudflare; do not execute without separate authorization.
```

No merge, Issue closure, remote branch deletion, deployment, restart, tag or release authorized. Owner `FINISH` is separately required for promotion.

## Safari regression found during local review (2026-11-22)

Owner screenshots show a valid twelve-hour Homepage personal login, while NewsGo briefly signs in and shows articles then immediately hides them and shows `Not signed in`; tab focus reproduces the flash. Source trace establishes that an unconfigured private thumbnail machine token causes the new NewsGo thumbnail proxy to return **HTTP 503**, which `authorizedNewsgoFetch` mistakes for a Homepage verification outage. Fixed within the same draft PR #115 by returning **502 Bad Gateway** for the thumbnail proxy's own configuration failures (without altering the outer Homepage auth 503) and added Go plus CI-only executable browser regression proving the reader remains signed in on 502 while a true 503 still denies access. Matching ScraperGo work is draft [PR #148](https://github.com/richard-mauri/scrapergo/pull/148), Issue #147. Existing auth and user data remain unchanged. Owner Safari still must verify no flash on initial load or focus, and image rendering after private credential provisioning. Exact-head CI status must be checked after this update.
