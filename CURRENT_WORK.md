# Current work

**Authoritative merged baseline:** GitHub `main` at `d62dbbfbdd0e676e1bb8e296e760772c80c8971a` (verified November 6, 2026). No open PRs at start of this issue; unrelated open Issues #55, #56 and #57 remain untouched. Older note about draft PR #54 is stale and not current.

**Active task:** [Issue #58](https://github.com/richard-mauri/newsgo/issues/58): export the result of **Read All Articles** as an explicit URL batch and provide a local replay command that sends those URLs to ScraperGo for persistent recovery-trace observations. Scope limited to NewsGo. Branch `feat/58-scrapergo-url-batch`; one draft PR against `main`. ScraperGo is not modified and its recovery tracer must be running independently for observations to be recorded.

**Current implementation:** The browser builds a deduplicated, exclusion-filtered URL list from the successful current Read All refresh and enables a **Download ScraperGo URLs** button in the completed briefing. Download happens only on click. It writes `newsgo-scrapergo-urls.json` with a `newsgo.scrapergo-batch.v1` marker and URL array. `cmd/scrapertrace` loads this file, validates bounds and HTTP(S) URLs, enforces numeric loopback-only ScraperGo target, submits sequential `POST /api/read` requests with pacing, and reports response totals without printing full URLs. Unit tests validate batch deduplication, safety, loopback and request behavior; dashboard tests validate export wiring. No NewsGo server routes, subscription data or auth behavior have changed.

**Pending verification:** Exact latest-head GitHub Actions `go test -count=1 ./...` and `go build ./...`, then owner Mac `./devbuild.sh` and `./newsgo.sh`, Safari Read All selection/export flow and manual command invoking a **running** ScraperGo at localhost:8083. Verify ScraperGo's own `go run ./cmd/recoveryreport` shows added observations. CI cannot establish real Safari or publisher access. Do not merge, tag, release, change version or delete branches without owner authorization.

**Mac review commands:**

```sh
cd ~/projects/newsgo
git fetch origin
git checkout feat/58-scrapergo-url-batch
git pull --ff-only origin feat/58-scrapergo-url-batch
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

Then select **Read All Articles** in Safari, wait for the refresh to finish, choose **Download ScraperGo URLs**, and run `go run ./cmd/scrapertrace -input ~/Downloads/newsgo-scrapergo-urls.json` in the NewsGo repository. After the batch, inspect `go run ./cmd/recoveryreport` in the ScraperGo repository.

**Deferred:** ScraperGo behavioral dead-code elimination is the higher-level goal; NewsGo only produces inputs to its independent diagnostic collector. Leave unrelated subscriber UI/admin/security work in its own issues.
