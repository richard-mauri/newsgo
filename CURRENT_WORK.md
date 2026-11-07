# Current work

**Merged baseline:** GitHub NewsGo `main` at `83f773e911dd304c2b8ed6f48a0fc689c5789cb9`, checked November 7, 2026. Draft PR #61 (ScraperGo diagnostic timings, safety review) and Issues #60/#62 are independent; nothing in their branch is merged here.

**Active task:** [Issue #63](https://github.com/richard-mauri/newsgo/issues/63) — add Reading Preferences **Use classification cache**, checked by default and applicable to verified NewsGo article classifications. Issue branch `feat/63-reader-cache-preference`, from merged main, one draft PR against main. No changes to ScraperGo, classification rejection rules or global cache TTL.

**Implementation:** A browser-local default-on checkbox controls a boolean sent with every `/api/personal` call, including preflight subscription validation. The server treats an omitted value as true for backward compatibility; false skips NewsGo cache lookup and insertion for that request but never erases shared cache entries. Preference is preserved through browser localStorage, Homepage account save/load, and legacy Server Sync storedPrefs, using optional JSON to default old profiles to on. RSS Metadata Only does not invoke classifier regardless of cache preference. Server administrator cache TTL remains independently configurable.

**Validation:** Go regression tests check bypass isolation, legacy request/profile migration and UI request/persistence wiring. GitHub Actions must pass `go test -count=1 ./...` and `go build ./...` on the exact final PR head; Mac Safari acceptance is required and not yet performed.

**Mac acceptance commands:**

```sh
cd ~/projects/newsgo
git fetch origin
git checkout feat/63-reader-cache-preference
git pull --ff-only origin feat/63-reader-cache-preference
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

Check Reading Preferences shows the checkbox checked on a fresh/default profile; switching it off persists after reload and bypasses NewsGo cache on the next verified refresh. Check Homepage account save/load and Server Sync if applicable. RSS Metadata Only remains unverified. Do not merge, tag, release, bump app version or modify main without explicit owner approval.
