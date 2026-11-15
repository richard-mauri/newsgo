# NewsGo — agent entry point

**Mandatory:** Read the entire [generated local suite policy](docs/suite-development-policy.md) before any change, including when offline. This is a generated, pinned copy of Homepage's [authoritative policy](https://github.com/richard-mauri/homepage/blob/main/SUITE_DEVELOPMENT_POLICY.md). Run `python3 scripts/suite_policy.py --offline` for local integrity; when online run `python3 scripts/suite_policy.py` for upstream freshness. An absent, edited or stale policy is a blocker for shared-process edits. Do not hand-edit the copy; use `--sync` on a reviewed policy synchronization PR.

Also read [README.md](README.md), [ARCHITECTURE.md](ARCHITECTURE.md), [DEVELOPMENT.md](DEVELOPMENT.md) and [CURRENT_WORK.md](CURRENT_WORK.md). Verify actual GitHub and owner-local states; do not trust an older snapshot. Homepage policy owns six-part visible handoffs, FINISH, approval and branch lifecycle; NewsGo-specific checks and Mac commands are in DEVELOPMENT.

NewsGo receives Homepage bearer handoffs; preserve user subscriptions, exclusions, feed catalogs, reading preferences, browser search histories, per-user storage, tokens, local counters and private Cloudflare configuration. Keep authorization server-side. Treat publisher URLs and fetched content as untrusted, preserve SSRF/request bounds and article classification. NewsGo sends articles to ScraperGo but never shares ScraperGo extraction code. Do not change another app without a separately approved scoped PR.

No source/runtime changes, forced resets, merges, deployments, tags or restarts are authorized by merely reading this entry point. Report passed/failed/unverified checks factually.
