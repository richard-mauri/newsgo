# NewsGo — app-specific development and validation

Read [AGENTS.md](AGENTS.md) and the complete [generated suite policy](docs/suite-development-policy.md). The suite policy, not this file, controls Issues, branches, draft PRs, exact-head CI, visible six-item handoff, owner FINISH and merged branch cleanup. Keep this document for NewsGo-specific validation and commands.

## NewsGo regression and data boundaries

Protect subscriptions, exclusions, feeds catalog, browser search histories, Homepage account storage and tokens, and per-user settings. Preserve existing publisher URL validation, classification, SSRF limits, import/export and regression coverage. Normal code changes require `gofmt`, `go test -count=1 ./...`, `go build ./...`, and exact-head GitHub Go CI. A Go CI pass cannot verify live publisher access, Mac Safari, or the tunnel.

## Owner Mac acceptance for runtime PRs

First run `git status --short --branch` and stop if dirty. Use the exact branch from the PR in place of the example:

```bash
cd ~/projects/newsgo
git status --short --branch
# Stop if uncommitted changes are shown.
git fetch origin
git checkout ACTUAL-PR-BRANCH
git pull --ff-only origin ACTUAL-PR-BRANCH
go test -count=1 ./...
./devbuild.sh
./newsgo.sh
```

`newsgo.sh` restarts NewsGo and its own tunnel, and requires distinct owner authorization. Check feed discovery and subscriptions, exclusions, article reader and export/import, Homepage user/guest handoff, token expiry/revocation, browser history and light/dark UI in Safari. Do not restart for documentation-only PRs. For safe post-FINISH checkout and local branch deletion, use the shared policy's exact clean-tree-guarded template with `~/projects/newsgo`.
