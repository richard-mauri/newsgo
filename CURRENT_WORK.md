# Current work — NewsGo

**Active runtime fix:** [Issue #97](https://github.com/richard-mauri/newsgo/issues/97), branch `fix/97-discovery-feed-labels`, [draft PR #98](https://github.com/richard-mauri/newsgo/pull/98) targeting `main`. Exact-head CI and Mac/Safari acceptance must be confirmed; owner merge approval is required.

GitHub `main` at `791d4d9` includes merged PR #96 (documentation refresh) and the earlier Homepage session improvements. No other issues or PRs were open when Issue #97 was created.

Homepage remains the session and renewal authority. NewsGo uses a validated, parent-bound child handoff, treats revoked supplied identity as 401 and temporary Homepage verification outage as 503, and permits no Close/Escape dismissal of the sign-in-required dialog. One-hour guests do not gain named-user account persistence. Feeds, reader, storage and publisher-link behavior are unchanged.

**Current change:** use each discovered feed's RSS/Atom/JSON title or meaningful URL-based fallback, infer finance and additional topical categories, and display full feed addresses. Align discovery controls, compact action widths and use actionable Unsubscribe with existing safeguards. Group discovered feeds by publisher with one exclusion action, remove Suggested additions and standardize shared button CSS. Use Exclude Domain/Allow Domain consistently. Improve discovery presentation by making feed metadata block-spaced (no explicit line-break elements), adding vertical separation above Export/Import, and labeling publisher counts as feeds found. The search's existing 24-probe / 8-result bounds remain unchanged; the interface now states the eight-feed display limit. No changes to accounts, subscription schema, server settings or deployment.

**Next:** run exact-head CI (`go test -count=1 ./...`, `go build ./...`), then owner verifies actual publisher results and Subscribe in Mac/Safari. Do not merge or deploy until explicitly approved.
