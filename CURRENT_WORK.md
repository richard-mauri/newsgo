# Current work — NewsGo

## Issue #91 — Homepage session identity and temporary guests

- **Authoritative baseline:** GitHub `main` already contains merged optional suite activity reporting from former Issue #89 / PR #90. The old `feat/89-suite-activity` branch is gone; do not reopen it.
- **Scoped branch:** `feat/91-suite-session-identity` from current `main`. One draft PR for this phase, coordinated through [Homepage Issue #32](https://github.com/richard-mauri/homepage/issues/32).
- **Goal:** Use Homepage's existing named and one-hour `guest-xxx` handoffs as authenticated identity for activity reporting and reader access. Verify bearer server-to-server and report origin-scoped session/user aliases; no user credentials or raw browser identifiers leave NewsGo. Multiple tabs sharing a login count once.
- **UX:** Without validated Homepage identity, show a same-page explanation/modal with user-initiated Homepage link. No automatic redirect loop. Temporary guests may read public articles, but must not autosave/autorestore or use Homepage server account settings. Named user preferences remain isolated.
- **Known integration gap:** Homepage's merged profile endpoint currently emits `activitySession` and `activityUser` only for ScraperGo. NewsGo cannot safely fabricate a common parent-session identity from distinct child bearer tokens. A narrowly scoped Homepage follow-up will be needed, with explicit owner coordination before changing Homepage.
- **Boundaries:** Preserve feed/subscription pipeline, classification cache, SSRF limits, settings data, server tokens, Cloudflare configuration and Mac launcher. No change to other suite apps, no merge or deployment without owner approval.
- **Verification:** Add regression coverage for identity, guest, multiple tabs, expiry, validation outage, persistence denial, popup and heartbeat; run exact-head GitHub Go CI. Owner acceptance uses Apple Silicon Mac and Safari.

## Owner Mac acceptance (only after development and owner approval)

```bash
cd ~/projects/newsgo
git status -sb
git fetch origin
git checkout feat/91-suite-session-identity
git pull --ff-only origin feat/91-suite-session-identity
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

The launcher restarts NewsGo and its tunnel; do not run it before an explicit acceptance checkpoint.
