# Current work — NewsGo

## Verified repository state (November 7, 2026)

- **Merged baseline:** GitHub `main`, commit `046783f1c020938a4d431711eb5bb905ea179727`, after the owner-approved PR #66 merge and Issue #65 closure.
- **Active workstream:** [Issue #67](https://github.com/richard-mauri/newsgo/issues/67), branch `feat/67-homepage-auth-diagnostics`, [draft PR #68](https://github.com/richard-mauri/newsgo/pull/68) targeting `main`. The code at the last verified development head before this documentation change was `780038161f21110dff4aea44ae991adb1eaef487`; retrieve the **actual latest PR head** at session start.
- **Owner approval boundary:** Owner approved removing all obsolete Server Sync code, backward compatibility and migration logic. **No approval to merge PR #68, modify `main`, tag or release.**
- At the start of a new session, read `AGENTS.md`, `DEVELOPMENT.md`, `ARCHITECTURE.md`, and this file **from this feature branch**, then check the actual GitHub PR and CI status. The branch, not an older assistant summary, is authoritative for unmerged work.

## Implemented and previously tested

- Homepage is NewsGo's sole live account and administrator identity and the sole per-user account store. The browser uses Homepage's revisioned `/api/storage/newsgo` for account save/load. Privileged NewsGo endpoints independently verify the active Homepage `admin` bearer against Homepage's fixed session profile endpoint; no legacy token fallback.
- Legacy Server Sync implementation, token-issuing CLI, `/api/sync` endpoint, browser login/auto-sync workflow, legacy migration helpers and obsolete hidden article-exclusion controls/listeners have been removed from PR #68. Current operational server settings use independent `SettingsStore` and `data/server-settings.json`. Existing historical token/profile files, if still on disk, are not used by runtime. Owner is sole current user and does not require migration or preservation of compatibility behavior; do not add it back.
- Both reader modes require ScraperGo-confirmed `article` classification. Reading Diagnostics shows latest-read cache hits, misses, bypasses, disabled-cache counts, configured TTL, hit rate and timing.
- Owner Safari tests established that an old saved TTL of zero disabled caching. Setting the TTL to 1440 minutes produced 329 hits and 6 misses (98%) and a warm read of 245 ms versus a cold 11.71 s. **These tests prove in-process cache reuse, not persistence across a NewsGo restart.**
- On November 7 around 10:07 AM Pacific, the owner showed a subsequent read with **13 cache hits, 321 misses, 4% hit rate, TTL 1440 minutes, 41.68 s total refresh**. This can result from a cold process-local cache; the precise restart history for that read was not verified. Do not conclude it was a valid 60-minute persistent-cache test.

## Highest-priority open requirement — NOT IMPLEMENTED

The owner explicitly clarified: **a persistent NewsGo classification cache with a one-hour (60-minute) TTL**. “Persistent” means cached classification decisions **survive stopping/restarting/rebuilding the NewsGo process**, not just closing a Safari tab. This is **not** satisfied by the current process-local `newsClassificationCache` in `internal/dashboard/personal.go`, even when its TTL is positive. The currently configured 1440-minute (24-hour) TTL is **not** the requested one-hour policy.

**Next engineering task within draft PR #68:** implement durable on-disk storage of validated classification decisions and their absolute expiration times; restore valid entries on startup; expire/reclassify entries after 60 minutes, including across restarts; preserve bounded storage, concurrency safety, strict eligibility and error behavior, URL normalization, and the per-reader cache-bypass semantics. Use the existing NewsGo settings/data directory or a clearly documented separate file; do not recreate Server Sync or a new identity system. Keep disk I/O bounded and protect the cache from malformed or corrupted state. Determine whether the administrator cache TTL control should be fixed at 60 or retain a configurable value with a 60-minute default, but **do not describe a 1440-minute default as fulfilling the owner's requirement**. Favor the explicit one-hour retention request and seek owner input only if a material unresolved choice blocks correct implementation.

**Required regression tests:** cold first read populates durable cache; restart/new NewsGo instance yields cache hits for the same normalized URLs; after simulated expiration the same URL triggers fresh ScraperGo classification; failures/unavailable/unknown do not become durable positive cache entries; reader opt-out bypasses lookup and insertion without erasing cached entries; malformed durable data fails safely; capacity bounds hold; server settings and user/account data remain separate and unaffected. For durability tests, avoid long sleeps by injecting/testing timestamps or clock boundaries.

## Remaining validation and release gate

- Verify exact **current PR head** and rerun/inspect GitHub Actions (`go test -count=1 ./...` and build) after any changes. The earlier head `7800381` passed CI; this is **not** evidence for later commits or the pending persistent-cache feature.
- Owner must still validate Homepage-only account save/load, administrator read/save, ordinary user denied admin, and account/session expiry in Safari. Existing classification/diagnostic Safari screenshots do not establish all those behaviors.
- Once durable caching is implemented, have owner read the same subscription batch, stop and restart **NewsGo**, read again within 60 minutes, and confirm substantial cache hits and faster classification, then verify expiration behavior with tests. Distinguish a Safari tab restart from a server process restart.
- Keep PR #68 **draft and unmerged** until explicit owner approval. Do not modify ScraperGo, Cloudflare configuration, private data, or the semantic application version without separate authorization.

## Owner Mac commands for the active branch

```bash
cd ~/projects/newsgo
git status --short
git fetch origin
git checkout feat/67-homepage-auth-diagnostics
git pull --ff-only origin feat/67-homepage-auth-diagnostics
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

## Handoff instruction for the next agent

**Resume existing Issue #67 / draft PR #68, not a new PR.** First inspect repository state and this file on the feature branch. The **first unfinished functional requirement** is restart-surviving, 60-minute classification caching; the present TTL setting and in-memory cache are not sufficient. Implement and test it before treating PR #68 as ready for owner acceptance. Update this file, `ARCHITECTURE.md` and the PR description to reflect what is actually implemented; cite exact-head CI. Do not merge without owner authorization.
