# Current work — NewsGo

**Banner correction:** NewsGo mounts the banner outside its dashboard `<main>` and consumes pinned UI module commit `80deb3b13f6a`, whose viewport-based CSS controls all horizontal geometry without decorative borders. Owner Safari review is pending.

**Active:** [Issue #108](https://github.com/richard-mauri/newsgo/issues/108), branch `ui/108-shared-banner`, [draft PR #109](https://github.com/richard-mauri/newsgo/pull/109) targeting `main` at `81b71ada1b0fc7fca9fe4dc1627afc6b30df0181`. Homepage shared UI foundation PR #55 remains open; this consumer pins module commit `80deb3b13f6a`. Prior Homepage session recovery PR #107 is merged and preserved.

**Candidate:** Replace NewsGo's legacy topbar with the shared banner and exact original Homepage blurred teal ambient shapes. Retain same-origin shared buttons, verified account state through `suite-auth-status`, NewsGo's editorial headline typography, subscriptions, exclusions, settings and all processing behavior. No auth protocol or persistence changes.

**Verification:** Source/whitespace and embedded JavaScript syntax checked; Go CI [44778349234](https://github.com/richard-mauri/newsgo/actions/runs/44778349234) passed on implementation head `a23383b82781cf618e5033cf1599e23c2aa25868`; repeat exact-head check after this documentation commit. Local `go` is not a working compiler; owner Mac/Safari/live acceptance unverified. Preserve user subscriptions, account data, browser histories, tokens, tunnels and build counters. No merge, issue closure, branch deletion, deployment, restart, tag or release authorized.
