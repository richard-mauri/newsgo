# Current work — NewsGo

**Issue:** [#116](https://github.com/richard-mauri/newsgo/issues/116) — conservative thumbnail-only loopback candidate. **Branch:** `fix/116-thumbnail-loopback-backend`; **base:** NewsGo merged main `e9078c70bfa36cdec20762ac17ae205b70e92b43` (PR #115 merged). ScraperGo merged main `eab6019d68282dc36003fb5d447cb981c4146787` (PR #148 merged). Do not revive either merged authentication PR.

## Scope and owner decision

Owner initiated BEGIN then NEXT on November 22, 2026 and explicitly requested a low-risk incremental approach with minimal hands-on testing. This candidate adds an **optional, server-only** `-thumbnail-backend` flag restricted to an explicit HTTP loopback IP and port. It affects only the authenticated NewsGo thumbnail proxy, not public ScraperGo reader links or the existing classification route. The Issue #116 branch updates `newsgo.sh` to pass `-thumbnail-backend http://127.0.0.1:8083` while retaining the public `-scrapergo` URL and private token sourcing; the binary default remains the existing public upstream. New regression tests cover fixed backend isolation and rejected destinations. No browser concurrency scheduler, rerender rewrite, semaphore/timeout fix, HTTP transport tuning, or classified-article behavior change.

## Preservation and verification

The original shared `SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN` is retained without rotation or exposure. Homepage per-user/Guest authorization, account isolation, subscriptions, exclusions, Save Changes, RSS mode, image size/MIME checks, redirects, SSRF and classification semantics are unchanged. No owner-local state, tunnel, private env, process or data was accessed or modified; GitHub branch launcher text was modified but not executed. The agent's GitHub-only environment does not establish the owner's Mac working-tree status. `./devbuild.sh`, owner M6 Mac Safari acceptance, and tunnel comparison are **unverified** until performed. CI on exact draft PR head must be inspected before any successful owner handoff; do not assume a pass. Follow Homepage `SUITE_DEVELOPMENT_POLICY.md` v1.1.0 and NewsGo's five documents.

## Owner operations and promotion

Before fetching or building on the M6 Mac, inspect `cd ~/projects/newsgo && git status --short --branch` and stop if dirty. After passing exact-head CI and permission to test locally, fetch and check out `fix/116-thumbnail-loopback-backend` without overwriting local changes, run `./devbuild.sh` **once**, then validate the executable version. **Do not run `./newsgo.sh`** (it restarts NewsGo and its dedicated Cloudflare tunnel) without separate authorization. Activation requires separate owner authorization to run the updated launcher, which restarts NewsGo and its dedicated tunnel. No merge, FINISH, Issue closure, branch deletion, deployment, tag, release or version bump is authorized.
