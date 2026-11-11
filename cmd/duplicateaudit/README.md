# NewsGo → ScraperGo duplicate-text audit

This standalone Go command loads NewsGo's administrator **Download Article URLs** JSON export (schema `newsgo.scrapergo-batch.v1`) and posts each URL to ScraperGo's existing `/api/read` endpoint. It does **not** fetch publishers directly, change ScraperGo, or save complete articles. It creates a local CSV of exact repeated text blocks and substantial prefix/suffix overlaps, including image-caption versus paragraph matches. Findings are **candidates**, not proof that editorial repetition is wrong.

## Quick start

From the NewsGo repository root (Go version as specified in `go.mod`):

For local ScraperGo on your Mac:

```sh
go run ./cmd/duplicateaudit -input ~/Downloads/newsgo-scrapergo-urls.json -limit 10 -output ~/Downloads/duplicate-audit-sample.csv
```

To use your already-public ScraperGo service remotely (only the approved HTTPS host is accepted):

```sh
go run ./cmd/duplicateaudit -input ~/Downloads/newsgo-scrapergo-urls.json -scrapergo https://scrapergo.themauris.org/api/read -public -limit 100 -output ~/Downloads/duplicate-audit-remote.csv
```

For a full export, specify the URL count, up to 500, e.g. `-limit 500`. The tool checks URLs sequentially, defaults to a 1-second delay between requests, and allows per-request `-timeout 45s`. It deliberately does not follow HTTP redirects. The public service may be unreachable or restrict traffic from some networks; unsuccessful requests appear as `read_failed` and yield a nonzero exit code. For a retry, choose a **new output filename**; existing reports are never overwritten.

## Report

Each CSV row shows the input article index, article URL, result status, ScraperGo HTTP status, block count, suspected match type, 1-based block positions, block kinds, and an excerpt up to 130 characters. A successful article without matches has one `ok` row; a failed read gets a `read_failed` row. Matching categories include `exact`, `caption_or_text_overlap`, and `image_`-prefixed variants. The audit also prints a progress line and final aggregate counts to the terminal.

The report includes article URLs (which may contain query strings) and text snippets. It is created with restrictive file permissions (`0600`) and should not be posted publicly without review. No API credentials are used or saved. The tool only accepts a numerical loopback HTTP address by default, or the explicitly opted-in `https://scrapergo.themauris.org/api/read` endpoint. It does not access localhost-only diagnostic endpoints.

## Accuracy limitations

This is a screening tool. Short duplicate bylines and credits under 24 characters are omitted to reduce false positives; they may still be recognized as part of longer repeated captions. Overlaps must cover a meaningful prefix or suffix, and noncontiguous partial matches may be missed. Legitimate repetitions (for instance a quote deliberately repeated in a summary) can be flagged. The CSV is most useful for comparing **before and after** a ScraperGo extraction change using the same URL batch, noting that publishers can change content between runs.

Run the synthetic tests from the NewsGo repository root with:

```sh
go test -count=1 ./cmd/duplicateaudit
```

## Relationship to the other diagnostics

`cmd/scrapertrace` measures classification/read outcomes and latency for the URL batch. This command inspects successful `/api/read` article blocks for candidate duplicate text. Both are development tools, not NewsGo server routes or privileged export APIs. The actual duplicate-text correctness work is tracked in [ScraperGo Issue #119](https://github.com/richard-mauri/scrapergo/issues/119) and [PR #120](https://github.com/richard-mauri/scrapergo/pull/120).
