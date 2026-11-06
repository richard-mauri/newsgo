# Current work

**Merged baseline:** GitHub NewsGo `main` at `83f773e911dd304c2b8ed6f48a0fc689c5789cb9` (verified November 6, 2026). Article URL export PR #59 was merged and Issue #58 closed by the owner; those features are not being changed. Open issues #55–#57 remain outside the task.

**Active issue:** [#60](https://github.com/richard-mauri/newsgo/issues/60), local URL batch replay reporting ScraperGo classification and returned HTTP statuses. Working branch `diag/60-classify-read-status`, based on merged main; one new draft PR targeting main. No ScraperGo changes or news-reader classification changes authorized.

**Implementation:** After the owner reported a hard Mac lockup during a four-worker, zero-delay replay, concurrent execution was removed. The CLI now runs one URL pair at a time with a 1-second default delay and processes only the first 10 URLs unless the operator explicitly requests a higher `-limit` (up to 500). Optional `-delay` affects pacing but not measured requests; `classify_wall` and `read_wall` measure each request independently and the summary provides min/mean/p50/p95/max plus total batch elapsed, prints the article index, classification, both ScraperGo HTTP status codes, bounded reader error kind, and upstream publisher HTTP status when supplied. Marks when classification says article but read fails, and totals such mismatches. Retains optional pacing, numeric-loopback restriction, no redirects, and privacy: console output contains no full URLs, article text, response errors, or query strings. Exit status reflects read failures, preserving existing behavior. Note that `/api/classify` invokes ScraperGo's reader and can use a classification cache, so repeated upstream access and disagreement are possible.

**Validation:** Go regression tests include classification/read disagreement, read-kind/publisher status parsing, malformed classification payload, output privacy, sequential execution, measured delay exclusion and latency statistics. GitHub Actions [run 42299091038](https://github.com/richard-mauri/newsgo/actions/runs/42299091038) passed tests and build at code head `54df3c29d850b5b0dfaa0a0e82167d8fdc3f0b44`. Earlier four-worker code passed CI but **failed owner Mac acceptance** with a hard system lockup. That code has been removed. Updated sequential-only safety changes need fresh exact-head CI; live owner validation remains blocked pending system-health review. Owner Mac acceptance is still pending.

**Owner Mac review:**

```sh
cd ~/projects/newsgo
git fetch origin
git checkout diag/60-classify-read-status
git pull --ff-only origin diag/60-classify-read-status
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

Only after checking the Mac is stable and ScraperGo healthy, a cautious test of the first 10 URLs is available: `go run ./cmd/scrapertrace -input ~/Downloads/newsgo-scrapergo-urls.json -limit 10`. Do not retry the full 300-URL batch as part of this acceptance. Inspect classification, both HTTP status fields and `classify_wall` for each article, together with the latency distribution. Requests and results are sequential; no concurrent workers remain. No second download is required. Optionally run `go run ./cmd/recoveryreport` in ScraperGo if its separate tracing PR is running. CI does not verify publisher responses or live Mac/Safari behavior.

**Merge gate:** Remain draft until exact-head CI and owner Mac acceptance; never merge, tag, release, bump version, or remove branches without approval. Dead-code removal is deferred. Future opt-in Verified RSS design is separately tracked in [Issue #62](https://github.com/richard-mauri/newsgo/issues/62), pending measured classification latency before any reader behavior changes.
