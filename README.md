# NewsGo v0.1.0 — FeedCheck

A standalone Go feed preflight utility. No third-party modules, Docker, or database are required. Supports RSS 2.0, Atom, RSS 1.0 (RDF), and JSON Feed. The initial ten feed URLs are candidates, **not yet verified**. Run the checker from your Mac to measure actual availability.

## Install into your existing repository

You have already created `~/projects/newsgo` and initialized `go.mod` with `github.com/richard-mauri/newsgo`. Extract the patch contents directly into that folder. The ZIP deliberately excludes `go.mod` and `.git`.

```sh
cd ~/projects/newsgo
unzip -o ~/Downloads/newsgo-v0.1.0-feedcheck.zip
go test ./...
go run ./cmd/feedcheck
```

The program prints feed status and metadata coverage, and writes `report.json` with per-feed HTTP status, format, item count, timestamps, and errors. `report.json` is local output, not part of the distributed ZIP. Use `-feeds custom.json`, `-out custom-report.json`, `-timeout 20s`, and `-parallel 2` as needed.

Statuses: `healthy` (at least one headline and link, newest known story <=72 hours old), `stale`, `undated`, `empty`, `poor_metadata`, `invalid_feed`, `http_error`, `error`. A successful feed fetch says nothing about ScraperGo article extraction: publisher access restrictions must be separately tested later.

## Scope

FeedCheck is the first milestone, not a running web server. It does not crawl arbitrary URLs, bypass paywalls, or modify ScraperGo. Next milestones are cached aggregation and an HTML dashboard opening articles via the existing ScraperGo URL interface.

## v0.1.1 — RSS link recovery and optional ScraperGo sampling

Some RSS publishers supply article links as a permalink `<guid>` instead of a normal `<link>`, or use an `href` attribute. FeedCheck now recognizes both formats generically, without publisher-specific logic. A GUID explicitly marked `isPermaLink="false"` is not interpreted as an article URL. FeedCheck continues to report missing links when no legitimate HTTP(S) URL exists.

Optionally ask your **running local** ScraperGo service to extract at most two articles per feed:

```sh
go run ./cmd/feedcheck -scrapergo http://127.0.0.1:8083 -samples 2
```

The tool makes POST requests to ScraperGo's existing `/api/read` endpoint and reports `readable`, `restricted`, or `other` for each sample in `report.json`, with a terminal summary. Publisher denial and subscription/registration notices are *restricted* rather than failed extraction. `other` includes network and format errors. This is diagnostic information, not a source eligibility decision. Sample tests make additional network requests to publishers via ScraperGo; they are **off by default**. Use `-samples 1` for minimal traffic. No settings, database, tunnel, or service code changes are required.

**Limitations:** a successful response must contain a nonempty article title to count as readable; this does not guarantee a complete clean extraction. Sampling is not a ranking score or comprehensive site certification. FeedCheck refetches a feed once to obtain sample URLs when `-scrapergo` is enabled. Avoid pointing the option to an untrusted service because it receives article URLs.
