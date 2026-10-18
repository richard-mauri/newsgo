# NewsGo v0.2.1 — Dashboard refinement

Replaces the v0.2.0 dashboard with a responsive editorial layout: lead-story cards, two-column recent stories, refined typography, automatic light/dark color schemes, and more legible filter controls.

All article-title links use the same `target="newsgo-scrapergo"` so Safari can reuse a single dedicated ScraperGo reader tab rather than opening a new one for each headline. Original publisher links still open separately. A named target must not include `rel="noopener"` or `rel="noreferrer"`, because those attributes can force a new browsing context and defeat reuse.

No external images have been added: the existing RSS normalization records titles, URLs, sources, categories and dates, not image metadata. The lead stories use typographic cards rather than unrelated stock images. Future feed-thumbnail support should use publisher-provided media when trustworthy.

Install from ~/projects/newsgo by extracting ZIP, then `go test ./...` and `go run ./cmd/newsgo`. Reload http://127.0.0.1:8084/. Stop the old NewsGo server first. No changes to FeedCheck, feed config, go.mod, git settings, or ScraperGo.
