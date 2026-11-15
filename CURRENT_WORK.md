# Current work — NewsGo

**Active history feature:** [Issue #99](https://github.com/richard-mauri/newsgo/issues/99), branch `feat/99-recent-discovery-history`, draft PR pending, based on authoritative GitHub `main` at `791d4d9`. This change is independent of ongoing draft [PR #98](https://github.com/richard-mauri/newsgo/pull/98) for feed discovery and duplicate handling. Both require separate review and explicit owner merge approval.

Homepage remains responsible for authentication and per-user account settings. This feature keeps two separate, bounded browser-local lists: recent topic/publisher searches and recent website URLs submitted for feed discovery. Up to eight unique entries per list appear as one-click repeat actions below the corresponding form. The latest entry appears first; each list has its own Clear history control. Failed network requests still retain the entered query for retries. Storage failures do not block search, and histories are never included in Homepage profile saves.

**Next:** run GitHub Actions Go CI (`go test -count=1 ./...`, `go build ./...`) at the final PR head; owner tests Safari persistence, repeat actions, clearing, and mobile layout using the branch-specific Mac commands from DEVELOPMENT.md. No main modification, merge, deploy, release, or subscription reset without explicit owner approval.
