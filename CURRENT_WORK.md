# Current work — NewsGo

## Issue #81 / draft PR — seamless account preferences (November 10, 2026)

- **Merged baseline:** GitHub `main` at `73223995943bd2c4c71745c77d588f74e7c1a366`; PR #80 (Homepage-only session management and direct-visit session recovery) is merged, Issue #79 closed. This candidate does not alter authentication authority, session lifecycle, admin authorization or other apps.
- **New issue/branch:** [Issue #81](https://github.com/richard-mauri/newsgo/issues/81), branch `fix/81-seamless-account-preferences`, draft PR pending creation.
- **Observed:** Account Sync tab still requires an explicit first Save or Load before NewsGo's automatic account writes begin. User requested removal and fully automatic behavior like SailingGo.
- **Implementation:** Remove obsolete Account Sync nav/panel/manual handlers; after validated Homepage profile, GET the user-scoped existing record, validate and apply it, set revision and automatically save changes. If 404, initialize a user-scoped profile, but do not copy an earlier verified user's local preferences into a different user. On failed GET, do not perform a potentially destructive initial PUT. Retain conflict checking, browser-only fallback, admin permissions, subscription/exclusion/reading preference schema and per-user settings isolation.
- **Validation:** New Go regression test checks absent manual UI and automatic hydration invariants. JavaScript syntax and latest-head CI pending. No local Mac, Safari, Cloudflare or subscription roundtrip testing performed remotely.
- **Authorization:** Owner authorized removal/automation but not merge, release or restart. Retain Issue → dedicated branch → draft PR workflow. GitHub main remains unmodified pending approval.

## Owner-side Mac validation

```sh
cd ~/projects/newsgo
git status --short --branch
git fetch origin
git checkout fix/81-seamless-account-preferences
git pull --ff-only origin fix/81-seamless-account-preferences
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

Do not restart until separately authorized. Test both authenticated and browser-only sessions, subscription add/remove, exclusion toggles, max article age/lead settings, account persistence after Safari restart and Homepage roundtrip, and second user isolation. Verify an expired/revoked Homepage bearer, unsupported record and 409 conflict never overwrite an account. Keep the previous Mac binary and settings data intact until accepted.
