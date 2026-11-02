# Active work — Homepage authentication and reader migration

**2026-11-02.** Merged baseline: `main`. [Issue #41](https://github.com/richard-mauri/newsgo/issues/41); branch `feat/41-homepage-auth-sync`; [draft PR #42](https://github.com/richard-mauri/newsgo/pull/42). No merge or deployment authorized.

## Implemented in current draft
- Added validated legacy NewsGo reader-profile ↔ Homepage record-data mapping and regression tests. Outer Homepage revisions remain separate from legacy NewsGo revisions.
- Added read-only Homepage session discovery in the NewsGo browser on the trusted HTTPS `newsgo.themauris.org` origin, using credentialed `GET /api/session/profile`. Shows the signed-in Homepage username in the header while explicitly stating NewsGo Server Sync remains in use. Unauthenticated or unavailable requests hide the indicator. No bearer tokens are copied.
- Added a read-only Homepage `GET /api/storage/newsgo` helper for subsequent migration work; it is not called automatically and cannot overwrite either side. CSP permits connecting only to the exact Homepage HTTPS origin.
- Existing local NewsGo token registry, `/api/sync`, admin settings, user data and all legacy server-sync flows are unchanged.

## Remaining / blockers
Implement an explicit migration flow with backup and profile-owner mapping; validate old and central profile revisions before any write. Then connect subscriptions/preferences to Homepage storage with correct 409/offline behavior, switch admin authorization to central role checks, and retire legacy token code only after verified Safari acceptance and rollback. Maintain compatibility with legacy data until owner approves cutover.

## Validation
Initial profile mapping CI passed at prior head; final-head CI pending. Owner Mac: `cd ~/projects/newsgo; git fetch origin; git checkout feat/41-homepage-auth-sync; git pull --ff-only origin feat/41-homepage-auth-sync; go test -count=1 ./...; ./devbuild.sh; ./newsgo.sh`. Test Safari signed-out/sign-in indicator, existing token-based sync, subscriptions and admin controls. Do not deploy or restart without owner action.
