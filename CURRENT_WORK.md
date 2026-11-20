# Current work — NewsGo

**Active:** [Issue #108](https://github.com/richard-mauri/newsgo/issues/108), branch `ui/108-shared-banner` targeting `main` at `81b71ada1b0fc7fca9fe4dc1627afc6b30df0181`. Homepage shared UI foundation PR #55 remains open; this consumer pins its immutable module commit `ee3fd0da8e9f`. Prior Homepage session recovery PR #107 is merged and preserved.

**Candidate:** Replace NewsGo's legacy topbar with the shared banner and exact original Homepage blurred teal ambient shapes. Retain same-origin shared buttons, verified account state through `suite-auth-status`, NewsGo's editorial headline typography, subscriptions, exclusions, settings and all processing behavior. No auth protocol or persistence changes.

**Verification:** Source/whitespace and regression tests reviewed; exact-head GitHub CI required. Local `go` is not a working compiler; owner Mac/Safari/live acceptance unverified. Preserve user subscriptions, account data, browser histories, tokens, tunnels and build counters. No merge, issue closure, branch deletion, deployment, restart, tag or release authorized.
