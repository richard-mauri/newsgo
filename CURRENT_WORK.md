# Current work

**Merged baseline:** GitHub main `83f773e911dd304c2b8ed6f48a0fc689c5789cb9` at feature start, verified November 7, 2026.

**Active feature:** [Issue #65](https://github.com/richard-mauri/newsgo/issues/65), branch `feat/65-verified-reading-cache` from main. Consolidates superseded Issues #62/#63 and draft PR #64 into **one issue, one branch, one draft PR**. Diagnostic Issue #60 / PR #61 is superseded and closed; the safe sequential classification-aware scrapertrace implementation and regression tests were incorporated into this one feature branch. GitHub main is authoritative; do not merge without explicit owner acceptance and authorization.

**Implementation:** Both normal and RSS metadata-only briefing modes require ScraperGo `/api/classify` to confirm `article` eligibility; homepages, unreadable, unavailable and unknown are filtered and counted. RSS-only still uses RSS headline metadata, dates and images; no longer bypasses classification. Existing server classification cache is shared; a default-on **Use classification cache** reader preference controls optional cache lookup/insertion, persists through localStorage, Homepage account and legacy Server Sync, and affects both modes. The default positive-result TTL is now 1440 minutes (24h), accepted administrator range 0–1440; saved explicit TTL values remain unchanged. No ScraperGo or unrelated suite changes.

**Diagnostic:** The integrated `cmd/scrapertrace` validates NewsGo accepted-URL exports with sequential ScraperGo classification and reading, reporting separate request HTTP status and wall-clock timing. Defaults: 10 URLs and one-second pacing; explicitly use `-limit 276` for the owner's 276-URL export. Do not use the older read-only runner or the superseded diagnostic branch.
\n**Tests:** Regression tests cover RSS-only classification gate and filtered accounting, cache reuse across modes and bypass without purging shared entries, compatibility of absent preference fields and UI/account persistence, classification TTL migration and bounds. CI for final PR head must run `go test -count=1 ./...` and `go build ./...`. Mac/Safari live acceptance not yet conducted.

**Mac commands:**

```sh
cd ~/projects/newsgo
git fetch origin
git checkout feat/65-verified-reading-cache
git pull --ff-only origin feat/65-verified-reading-cache
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

Safari acceptance: confirm Reading Preferences has cache enabled by default, disable/enable remains saved after reload and account restore, run Read All in both RSS-only and regular modes, and verify known homepages appear in Filtered / unprocessed rather than visible headlines. Check server admin classification cache minutes and ensure no excessive CPU/memory or new timeouts; stop testing if machine instability occurs. Note bypass only controls NewsGo's cache, not ScraperGo's internal caches. First verify a sample of 3–10 URLs; a full 276-URL run requires explicit `-limit 276`, may take several minutes due to pacing and publisher response times, and should be attempted only if machine load is acceptable.

**Gate:** One draft PR, no version bump, no tag or merge without owner approval. Superseded diagnostic PR #61 / Issue #60 are closed without merging.
