# NewsGo v0.2.0 — dashboard

NewsGo adds a local web dashboard to the existing v0.1.1 FeedCheck project. No database, Docker, external Go modules, or ScraperGo modifications are necessary.

From `~/projects/newsgo`:

```sh
unzip -o ~/Downloads/newsgo-v0.2.0-dashboard.zip
go test ./...
go run ./cmd/newsgo
```

Open http://127.0.0.1:8084/ in Safari. Headlines open ScraperGo, which must be running at http://127.0.0.1:8083/; an **Original** link opens the publisher directly. FeedCheck remains available with `go run ./cmd/feedcheck`.

Options: `-listen 127.0.0.1:8084`, `-feeds feeds.json`, `-refresh 15m`, `-scrapergo http://127.0.0.1:8083`. Refresh interval minimum 1 minute. Fetches at most four feeds simultaneously; response size is limited to 4 MiB per feed, HTTP client timeout is 12 seconds. The dashboard caches up to 600 headlines in memory, orders them by publication date, deduplicates by article URL, and preserves previously retrieved stories for individual feeds that fail a refresh. A restarted service will fetch again. Filters are server-side and do not require JavaScript. Initial fetching runs asynchronously, so the page may briefly show no headlines. It has a `/health` endpoint. No authentication is provided; keep it bound to loopback until authentication/access controls are planned.

Sources are edited in the existing `feeds.json`. The ZIP does not overwrite `feeds.json`, `go.mod`, FeedCheck, or git configuration.
