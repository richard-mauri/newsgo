# Active work — Homepage Account Sync

**2026-11-02.** Baseline: `main`. Issue [#41](https://github.com/richard-mauri/newsgo/issues/41), branch `feat/41-homepage-auth-sync`, draft [PR #42](https://github.com/richard-mauri/newsgo/pull/42). No merge or deployment.

## Current candidate
- Aligned NewsGo with MauriStocks: settings remain browser-local, `GET/PUT https://home.themauris.org/api/storage/newsgo` uses the signed-in Homepage user's cookie and server revision. Profile is established with **Save to account** or explicitly restored with **Load from account**; later edits autosave after the first account save/load.
- Replaced four migration actions with the initial Save to account and Load from account; once account sync is established the initial save button disappears. Removed the displayed legacy recovery section and renamed the main navigation to Account Sync. Legacy NewsGo Server Sync records remain untouched on disk; old token controls are hidden from the normal Account Sync panel but retained in the DOM/backend while authentication migration is incomplete.
- Initial account save checks for existing Homepage profile and requests confirmation before replacement. Explicit load requests confirmation and disconnects legacy browser sync first. Account saves use revision checks, pause on 409, retain browser state on failure, and retry transient errors.
- Homepage session is verified through a same-origin NewsGo proxy; Safari ordinary tabs verified username `richard`. Dock web-app cookie isolation is separate from regular Safari.
- Go tests cover reader profile mapping and session endpoint. No automatic cross-device load at sign-in; account load is intentional, like MauriStocks. Admin settings still use legacy role checks pending centralized authorization.

## Remaining
Owner must test the latest account workflow on Mac/Safari, including saved subscriptions and column-equivalent reader preferences, 409 conflicts, offline behavior and profile survival. Refactor NewsGo admin authorization to Homepage roles before retiring token-based admin login. Preserve NewsGo `data/` and Homepage `private/` backups.

## Mac validation
```bash
cd ~/projects/newsgo
git fetch origin
git checkout feat/41-homepage-auth-sync
git pull --ff-only origin feat/41-homepage-auth-sync
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```
Final-head CI and Safari acceptance pending.
