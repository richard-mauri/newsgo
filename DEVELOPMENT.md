# NewsGo — development and validation

This document owns **Git workflow, Issue/PR procedures, validation and branch cleanup**. [AGENTS.md](AGENTS.md) tells agents what rules to follow; [ARCHITECTURE.md](ARCHITECTURE.md) describes implementation boundaries; [CURRENT_WORK.md](CURRENT_WORK.md) records actual active state. [README.md](README.md) covers operation.

## Baseline and scope
GitHub `main` is authoritative. For each new work item, check the current `main`, existing Issues and PRs, create a scoped Issue and a **new dedicated branch from `main`** (`fix/123-summary`, `feat/124-summary`, or `docs/125-summary`). Work only within that Issue; open one focused **draft PR targeting `main`**. Do **not** reuse a permanent `dev` branch or an unrelated PR. Keep semantic version unchanged during normal development unless owner authorizes a release.

## Validation and review
1. At session start read AGENTS, DEVELOPMENT, ARCHITECTURE and CURRENT_WORK on the working branch. Identify actual branch, merged baseline, PR and pending constraints; explicitly summarize intended workflow before development.
2. Preserve user subscriptions, exclusions, reading preferences, tokens, private `data/`, Cloudflare tunnel configuration, existing launchers and local build counter. Add regression tests for changes, format touched Go code and honor security boundaries.
3. After commit/push, inspect GitHub Actions at the **exact PR head** (`go test -count=1 ./...` and `go build ./...`). Green CI does not validate live feeds, Safari, tunnel behavior or publisher access.
4. Update and verify the real PR title/description, current commit, Issue reference, change scope, CI result, owner acceptance and limitations. Keep CURRENT_WORK.md a single factual current snapshot.
5. For a code/runtime PR give the owner complete Mac commands, replacing the branch placeholder with its exact PR head branch:
   ```bash
   cd ~/projects/newsgo
   git fetch origin
   git checkout ACTUAL-ISSUE-BRANCH
   git pull --ff-only origin ACTUAL-ISSUE-BRANCH
   go test -count=1 ./...
   ./devbuild.sh
   ./newsgo.sh
   ```
   Owner's Mac Git supports `git checkout` rather than `git switch`. Document-only PRs require Markdown/link/path checking and normal CI; no runtime rebuild or Safari acceptance is necessary.
6. Owner performs applicable Safari/live acceptance. **Do not merge, tag, release, enable auto-merge or modify `main` without explicit owner approval.** After authorized merge, sync local `main` and verify promoted build when relevant.

## Post-merge branch lifecycle
- **Every merged task finishes with short-lived branch cleanup.** Verify its PR was merged into `main` and that the branch is not `main`, not head of an open PR, and contains no unmerged work. Do not delete a closed-but-unmerged branch just because its PR is closed.
- Once verified, delete the merged feature branch on GitHub (or use the owner's repository **Automatically delete head branches** setting), then prune remote-tracking refs and delete the local branch safely with `git branch -d` after checking out `main`. Never use force-deletion to conceal an unmerged state.
- `archive/`, `release/`, or other intentionally retained branches require separate owner authorization for deletion. Existing historical branches require **individual review and separate owner approval for bulk cleanup**; this policy does not authorize blanket deletion.
- `main` retains merged commit history after head branch deletion. The next Issue starts from updated `main`, not from the previous head branch.

## Current legacy branch
NewsGo's former shared `dev` has already served merged PR #10. It is not a permitted head for new tasks; handle its eventual deletion separately after verifying local and remote branch safety.
