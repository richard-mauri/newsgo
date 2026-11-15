# Current work — NewsGo

**Active runtime fix:** [Issue #97](https://github.com/richard-mauri/newsgo/issues/97), branch `fix/97-discovery-feed-labels`, [draft PR #98](https://github.com/richard-mauri/newsgo/pull/98) targeting `main`. Exact-head CI and Mac/Safari acceptance must be confirmed; owner merge approval is required.

GitHub `main` at `791d4d9` includes merged PR #96 (documentation refresh) and the earlier Homepage session improvements. No other issues or PRs were open when Issue #97 was created.

Homepage remains the session and renewal authority. NewsGo uses a validated, parent-bound child handoff, treats revoked supplied identity as 401 and temporary Homepage verification outage as 503, and permits no Close/Escape dismissal of the sign-in-required dialog. One-hour guests do not gain named-user account persistence. Feeds, reader, storage and publisher-link behavior are unchanged.

**Current change:** use each discovered feed's RSS/Atom/JSON title or meaningful URL-based fallback, infer finance and additional topical categories, display full discovered feed addresses without collapsible disclosures, and transition successful Subscribe buttons from Checking to Subscribed. Align action controls throughout discovery with explicit grid tracks, including precise grid placement of the dynamically inserted exclusion badge and actions in the included-source cards; ensure the discovered-feed grid wins over the generic row display rule. Apply the same accent-colored hover label as the Find Articles navigation, plus border and keyboard-focus highlighting, while preserving contrast on selected/filled buttons and disabled states. A failed subscription save is surfaced as a failure. No changes to accounts, subscription schema, server settings or deployment.

**Next:** run exact-head CI (`go test -count=1 ./...`, `go build ./...`), then owner verifies actual publisher results and Subscribe in Mac/Safari. Do not merge or deploy until explicitly approved.
