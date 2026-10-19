# NewsGo discovery service evaluation (experimental)

This diagnostic tool **does not alter the NewsGo dashboard**, saved subscriptions, or ScraperGo.
It compares keyword-to-feed lookup APIs, not just searches over already followed feeds.

Run from `~/projects/newsgo` after extracting the ZIP:

```bash
go test ./cmd/discoverycheck
go run ./cmd/discoverycheck -output discovery-report.json
```

The default searches are **college football**, **sailing**, and **gardening**. You can run other terms:

```bash
go run ./cmd/discoverycheck -q 'local California news' -q 'woodworking' -limit 5
```

Services sampled:

- **Feedly** `https://cloud.feedly.com/v3/search/feeds?query=...` (a historically documented public feed-search endpoint; present authentication and access policies have **not** been confirmed).
- **RSSAmplifier** `https://rssamplifier.com/api/search?q=...` (an open searchable feed directory, with an evolving response schema).

The utility records HTTP status, elapsed time, candidate titles, feed URLs and errors. A non-200 result or unrecognized JSON format is reported rather than treated as successful discovery. Even a valid-looking URL must still pass NewsGo FeedCheck before subscription. This test does not claim that either API is production-ready, free for third-party services, or comprehensive. Review provider terms and rate limits before integrating.

The command does **not** make subscriptions, store credentials, scrape publisher pages, or open arbitrary discovered URLs. It only queries the two fixed service endpoints.

The optional `discovery-report.json` contains lookup results; don't commit it by accident. The ZIP deliberately leaves `.gitignore`, `go.mod`, current dashboard files, and startup scripts untouched.
