# NewsGo — agent entry point

**Mandatory before any development action:** Read the **current full text** of [Homepage's authoritative suite development policy](https://github.com/richard-mauri/homepage/blob/main/SUITE_DEVELOPMENT_POLICY.md) from Homepage `main`, then read this repository's five standard documents. Homepage alone owns shared workflow, six-item factual handoff, FINISH, security/data preservation, approvals and branch cleanup. This repository does not keep a policy copy. If the live Homepage policy cannot be retrieved or verified, **stop shared development and promotion** until it becomes available; do not rely on an outdated cached or remembered version. Policy exceptions require owner approval.

Read [README.md](README.md), [ARCHITECTURE.md](ARCHITECTURE.md), [DEVELOPMENT.md](DEVELOPMENT.md) and [CURRENT_WORK.md](CURRENT_WORK.md). Verify actual GitHub and owner-local states; do not trust an older snapshot. Homepage policy owns six-part visible handoffs, FINISH, approval and branch lifecycle; NewsGo-specific checks and Mac commands are in DEVELOPMENT.

NewsGo receives Homepage bearer handoffs; preserve user subscriptions, exclusions, feed catalogs, reading preferences, browser search histories, per-user storage, tokens, local counters and private Cloudflare configuration. Keep authorization server-side. Treat publisher URLs and fetched content as untrusted, preserve SSRF/request bounds and article classification. NewsGo sends articles to ScraperGo but never shares ScraperGo extraction code. Do not change another app without a separately approved scoped PR.

No source/runtime changes, forced resets, merges, deployments, tags or restarts are authorized by merely reading this entry point. Report passed/failed/unverified checks factually.
