# NewsGo — development and validation

This document owns **Git workflow, Issue/PR procedures, validation and branch cleanup**. [AGENTS.md](AGENTS.md) tells agents what rules to follow; [ARCHITECTURE.md](ARCHITECTURE.md) describes implementation boundaries; [CURRENT_WORK.md](CURRENT_WORK.md) records actual active state. [README.md](README.md) covers operation.

## Baseline and scope
GitHub `main` is authoritative. For each new work item, check the current `main`, existing Issues and PRs, create a scoped Issue and a **new dedicated branch from `main`** (`fix/123-summary`, `feat/124-summary`, or `docs/125-summary`). Work only within that Issue; open one focused **draft PR targeting `main`**. Do **not** reuse a permanent `dev` branch or an unrelated PR. Keep semantic version unchanged during normal development unless owner authorizes a release.

## Validation and review
1. At session start read AGENTS, DEVELOPMENT, ARCHITECTURE and CURRENT_WORK on the working branch. Identify actual branch, merged baseline, PR and pending constraints; explicitly summarize intended workflow before development.
2. Preserve active user subscriptions, exclusions, reading preferences, Homepage credentials, private settings data, Cloudflare tunnel configuration, launchers and local build counter unless the owner explicitly authorizes a clean reset. Add regression tests for changes, format touched Go code and honor security boundaries.
3. After commit/push, inspect GitHub Actions at the **exact PR head** (`go test -count=1 ./...` and `go build ./...`). Green CI does not validate live feeds, Safari, tunnel behavior or publisher access.
4. Update and verify the real PR title/description, current commit, Issue reference, change scope, CI result, owner acceptance and limitations. Keep **README.md, ARCHITECTURE.md, DEVELOPMENT.md, AGENTS.md and CURRENT_WORK.md** accurate for their respective documentation responsibilities; do not leave stale workflow, implementation or open-PR claims. Keep CURRENT_WORK.md a single factual current snapshot.
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

## Visible development checklist

For **every NewsGo development status or handoff**, include a short, factual checklist in the user-facing reply covering all six items:

1. **Baseline and scope:** authoritative `main`, active Issue(s), branch and PR; state whether related changes must be tested together.
2. **Data safety:** whether local modifications, subscriptions, exclusions, preferences, account records, settings, tokens and build counter were preserved.
3. **Documentation current:** verify README.md, ARCHITECTURE.md, DEVELOPMENT.md, AGENTS.md and CURRENT_WORK.md are up to date for the actual implementation and workflow; report any exception.
4. **Verification:** exact-head CI result, regression checks and what Mac/Safari or live behavior remains unverified.
5. **Local commands:** show complete branch-specific Mac commands for code acceptance; after an authorized `FINISH`, show the complete safe local transition back to `main`.
6. **Promotion authority:** clearly state whether merge, Issue closure, branch deletion, deployment, tag or release was authorized and completed. Never claim a check passed before it is observed.

A feature requested as part of the same owner acceptance session stays on the active draft PR when the owner directs integrated testing, even if separately tracked by Issue.

## Owner `FINISH` command

The owner's explicit `FINISH` in the active NewsGo development context authorizes **merge the active PR → close its associated Issue(s) → delete its merged short-lived remote branch → provide safe Mac commands to return the local checkout to `main`**. This instruction does not authorize deployment, release, tagging, unrelated branch deletion, or discarding local work. Execute the operations in order, checking each result:

1. Identify the active PR and its head/base, verify the latest head has passing required CI and no unresolved blocking merge condition. Confirm the owner's `FINISH` applies to this PR. If checks fail, stop and report them.
2. Merge the verified PR using a normal allowed merge method without bypassing protections; verify GitHub reports it merged to `main`.
3. Close only Issue(s) attached to this completed work; verify closure.
4. Verify the short-lived PR branch is merged and not used by any open PR, then delete that **exact remote branch** without forcing. Never delete `main`, unrelated branches or historical branches.
5. Give the owner the following safe local recovery sequence, adapted to the actual branch. It must stop if the working tree has uncommitted changes; never use `reset --hard` or forced checkout:

   ```bash
   cd ~/projects/newsgo
   git status --short --branch
   # Continue only if no uncommitted files are listed.
   git fetch --prune origin
   git checkout main
   git pull --ff-only origin main
   git branch -d ACTUAL-MERGED-BRANCH
   git status --short --branch
   ```

6. Report the exact merge commit, Issue states, branch cleanup result, CI and Mac/Safari acceptance status, and confirm that no deploy/release occurred. If GitHub already automatically deleted the branch, verify that state and skip a second deletion.

## Post-merge branch lifecycle
- **Every merged task finishes with short-lived branch cleanup.** Verify its PR was merged into `main` and that the branch is not `main`, not head of an open PR, and contains no unmerged work. Do not delete a closed-but-unmerged branch just because its PR is closed.
- Once verified, delete the merged feature branch on GitHub (or use the owner's repository **Automatically delete head branches** setting), then prune remote-tracking refs and delete the local branch safely with `git branch -d` after checking out `main`. Never use force-deletion to conceal an unmerged state.
- `archive/`, `release/`, or other intentionally retained branches require separate owner authorization for deletion. Existing historical branches require **individual review and separate owner approval for bulk cleanup**; this policy does not authorize blanket deletion.
- `main` retains merged commit history after head branch deletion. The next Issue starts from updated `main`, not from the previous head branch.

