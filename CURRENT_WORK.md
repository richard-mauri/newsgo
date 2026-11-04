# NewsGo — current work (2026-11-04)

Authoritative merged baseline: GitHub `main` (Homepage return navigation PR #50 merged). Active Issue [#48](https://github.com/richard-mauri/newsgo/issues/48), branch `fix/48-signin-feedback`, [draft PR #51](https://github.com/richard-mauri/newsgo/pull/51) targeting `main`.

Scope: Improve the current token-based Server Sync Connect control in NewsGo. The source currently calls it `Connect` in the Server Sync panel, rather than `Sign in` in a separate Personal settings popup. Changes add clearly differentiated hover, keyboard focus and pressed styling, disable repeated submissions and show a `Connecting…` busy state, a 15-second request timeout, and accessible status for success, invalid/expired token and network failure. No token or local preference data is printed or changed as part of the feedback.

The browser-local preferences, existing sync conflict-resolution flow, role checks, Homepage account handoff, feed catalog, ScraperGo interaction, version and launch configuration are unchanged.

Go regression tests cover the UI contract and sign-in status/state changes. Exact-head GitHub Actions and owner Mac/Safari testing are pending and must be recorded before acceptance. Owner should inspect normal/invalid/expired token, disconnected network, slow request, repeated click, light/dark theme and keyboard focus using non-sensitive test credentials.

Do not merge, tag, release, restart or alter production services without separate owner approval. Owner's local working tree, build counter, private data and tunnels have not been inspected.
