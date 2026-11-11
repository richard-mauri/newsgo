# Current work — NewsGo

## Issue #87 — preserve standalone duplicate-audit CLI (November 11, 2026)

- **Authoritative source:** GitHub `main`; current work is scoped to branch `feat/87-preserve-duplicate-audit`, [Issue #87](https://github.com/richard-mauri/newsgo/issues/87). This issue was reopened at the owner's explicit request to preserve and merge the useful standalone duplicate-analysis tool. It does **not** reopen the ScraperGo extraction-fix issue.
- **Change:** Add `cmd/duplicateaudit/main.go`, `main_test.go`, and `README.md` from the previously delivered and actively used ZIP. Update the root README alongside the existing `cmd/scrapertrace` export instructions. This tool reads `newsgo.scrapergo-batch.v1` JSON, calls the existing ScraperGo `POST /api/read`, reports exact/overlap duplicate *candidates* in a CSV, and records failed reads separately.
- **Security/integration:** No NewsGo server, UI, authentication, feed processing, cache, launcher, or settings changes. Defaults to sequential requests with pacing and a local loopback ScraperGo endpoint; the approved public HTTPS hostname requires explicit `-public`. Reports are private local files that can contain URLs and short excerpts. No publisher content or token is committed.
- **Validation:** Standalone command's Go tests passed locally on Linux Go 1.23.2 before import. GitHub Actions on exact PR head and owner's Mac test/build have not yet been confirmed at the time of this handoff; do not claim a green CI until checked. No NewsGo restart or Safari acceptance is needed for this CLI-only change. The owner explicitly authorized preserving and merging the utility, subject to verification.
- **Boundary:** [ScraperGo Issue #119](https://github.com/richard-mauri/scrapergo/issues/119) / [PR #120](https://github.com/richard-mauri/scrapergo/pull/120) remain solely responsible for correcting duplicate article text. This NewsGo PR contains no extraction logic or policy changes.

## Previous work

NewsGo Issue #85 / PR #86 concerned avoiding repeated classification of excluded domains and better classification-cache miss diagnostics. Its previous handoff is retained in Git history; it is unrelated to the duplicate-audit CLI. Check live GitHub issue/PR state rather than assuming a historical handoff reflects current status.

## Owner Mac verification (CLI code, no server restart)

```sh
cd ~/projects/newsgo
git fetch origin
git checkout feat/87-preserve-duplicate-audit
git pull --ff-only origin feat/87-preserve-duplicate-audit
go test -count=1 ./...
./devbuild.sh
go run ./cmd/duplicateaudit -input ~/Downloads/newsgo-scrapergo-urls.json -limit 10 -output ~/Downloads/duplicate-audit-integrated.csv
```

Review the CSV for duplicate candidates, legitimate repetitions, and explicit `read_failed` records; choose a new output name for each run. `./newsgo.sh` is **not needed** for this tool-only change, and no release or deployment is authorized.
