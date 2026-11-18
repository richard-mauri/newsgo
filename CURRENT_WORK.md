# Current work — NewsGo

**Active scope:** [Issue #106](https://github.com/richard-mauri/newsgo/issues/106), branch `fix/106-homepage-session-recovery`, based on merged `main` at `dc28094e109c3594ab1fba828dccf6f68e5c683f`. Draft PR and exact-head CI details follow publication. Prior PR #105 is merged; its shared-button changes remain intact.

**Reproduction:** In Safari, direct navigation to NewsGo in a fresh tab showed `Not signed in` and a blocking sign-in dialog despite Homepage remaining authenticated as Richard. Navigating to Homepage's `/launch/newsgo` bridge returned an authenticated NewsGo tab. Commit `855c38f` removed the fallback added by `7322399`.

**Change:** Retain sibling-tab handoff recovery and then attempt Homepage bridge recovery once if needed. A rejected cached child bearer also receives one controlled recovery attempt. Tab-scoped storage prevents redirect loops, resets on verified identity or explicit local-only return, and fails closed when unavailable. Guest access, per-user preferences, server-side authorization, Homepage's 12-hour parent session, buttons and other application behavior remain unchanged.

**Verification and next:** Run `./devbuild.sh` if the toolchain permits, then confirm exact-head GitHub CI after a scoped draft PR. Owner Mac/Safari acceptance must test direct visits, sibling tabs, genuinely signed-out/guest returns, rejected tokens, persisted subscriptions and normal Homepage-card launches. Owner local checkout and browser data are inaccessible; no merge, issue closure, deployment, restart, tag, semantic version or release is authorized.
