# Current work — NewsGo

**Active documentation review:** Issue #95, branch `docs/95-session-markdown-refresh` based on the merged `main` snapshot; documentation-only draft PR, awaiting owner approval. This branch does not change runtime behavior.

GitHub `main` at `540c766` includes Homepage-verified guest/named identity and activity, revoked-heartbeat handling, a non-dismissible sign-in-required dialog, and current-tab article navigation to ScraperGo. All prior session PRs/issues were merged or closed; no issues or PRs remain open and no implementation branch is active.

Homepage remains the session and renewal authority. NewsGo uses a validated, parent-bound child handoff, treats revoked supplied identity as 401 and temporary Homepage verification outage as 503, and permits no Close/Escape dismissal of the sign-in-required dialog. One-hour guests do not gain named-user account persistence. Feeds, reader, storage and publisher-link behavior are unchanged.

**Next:** bring Markdown into line with merged code in a fresh documentation-only draft PR. Verify paths/links and exact-head CI; no restart or Mac/Safari acceptance is needed for docs-only changes. For runtime work, protect unpushed local changes, use an issue-scoped branch, and require owner approval to merge or deploy.
