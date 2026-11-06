# Current work

**Merged baseline:** GitHub NewsGo `main` at `83f773e911dd304c2b8ed6f48a0fc689c5789cb9` (verified November 6, 2026). Article URL export PR #59 was merged and Issue #58 closed by the owner; those features are not being changed. Open issues #55–#57 remain outside the task.

**Active issue:** [#60](https://github.com/richard-mauri/newsgo/issues/60), local URL batch replay reporting ScraperGo classification and returned HTTP statuses. Working branch `diag/60-classify-read-status`, based on merged main; one new draft PR targeting main. No ScraperGo changes or news-reader classification changes authorized.

**Implementation:** For each validated URL the existing local replay CLI calls `/api/classify` and `/api/read` in sequence, prints the article index, classification, both ScraperGo HTTP status codes, bounded reader error kind, and upstream publisher HTTP status when supplied. Marks when classification says article but read fails, and totals such mismatches. Retains request pacing, numeric-loopback restriction, no redirects, and privacy: console output contains no full URLs, article text, response errors, or query strings. Exit status reflects read failures, preserving existing behavior. Note that `/api/classify` invokes ScraperGo's reader and can use a classification cache, so repeated upstream access and disagreement are possible.

**Validation:** Go regression tests include classification/read disagreement, read-kind/publisher status parsing, malformed classification payload and output privacy. Exact final-head CI and Mac acceptance must be recorded below before merge.

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

With ScraperGo running, replay the existing diagnostic file: `go run ./cmd/scrapertrace -input ~/Downloads/newsgo-scrapergo-urls.json`. Inspect classification and both HTTP status fields for each article. No second download is required. Optionally run `go run ./cmd/recoveryreport` in ScraperGo if its separate tracing PR is running. CI does not verify publisher responses or live Mac/Safari behavior.

**Merge gate:** Remain draft until exact-head CI and owner Mac acceptance; never merge, tag, release, bump version, or remove branches without approval. Dead-code removal and production NewsGo readability-policy changes remain separate and deferred.
