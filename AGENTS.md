# NewsGo — development agent rules

This file is the **authoritative instruction set for agents**. Read it before making changes. Do not put feature history, architecture explanations, or active task notes here: see [README.md](README.md), [ARCHITECTURE.md](ARCHITECTURE.md), and [CURRENT_WORK.md](CURRENT_WORK.md).

## Branching and owner authority
- **`main` is the authoritative integrated baseline.** For every new task, start from current `main`, create a GitHub Issue and **dedicated issue-scoped branch** (`fix/123-description` or `feat/124-description`), and open a focused **draft PR to `main`**. Never reuse a long-lived `dev` branch for new tasks, or add unrelated changes to an old PR.
- The existing `dev` and draft PR #10 are legacy work under reconciliation. Preserve them until reviewed; do not treat them as the pattern for new work. See CURRENT_WORK.md.
- No PR merge, release, tag, semantic version bump, branch deletion, force-push, history rewrite, auto-merge, or changes to `main` without express owner authorization.
- Preserve local data, tokens, subscriptions/exclusions, reading preferences, feeds.json, tunnel configuration, ports, launchers, and build counters.

## Mandatory workflow
1. Before edits read all four project documents on the working branch, check latest `main` and actual Issues/PRs, and **acknowledge** baseline, scoped issue/branch/PR, code constraints, validation plan and owner release gate.
2. Create a GitHub Issue and feature/fix branch **from current `main`**. Make only scoped changes with meaningful regression tests, format touched Go code, and update only the owning document.
3. Open a draft PR into `main`. Run and inspect GitHub Actions at the **exact final head** (`go test -count=1 ./...`, `go build ./...`); fix failures. CI is not Mac/Safari acceptance.
4. Before handoff update and verify the **saved PR title and description** with exact head, Issue, changes, CI, owner browser acceptance, limitations and unmerged status. Keep CURRENT_WORK.md current; no chronology dump.
5. For code changes repeat **complete Mac instructions with the ACTUAL feature branch**, never an example branch, including:
   ```bash
   cd ~/projects/newsgo
   git fetch origin
   git checkout ACTUAL-ISSUE-BRANCH
   git pull --ff-only origin ACTUAL-ISSUE-BRANCH
   go test -count=1 ./...
   ./devbuild.sh
   ./newsgo.sh
   ```
   Use `git checkout` rather than `git switch` on the owner's older Mac. Documentation-only work needs documentary review and CI, not a runtime rebuild.
6. Owner performs required local/browser/live integration acceptance. **Only the owner authorizes promotion**; after approval update `main` and rebuild/verify the merged candidate. Start the next issue from updated `main`, never from the prior feature branch.
7. Recheck mergeability before handoff; resolve integration conflicts without discarding functionality or overwriting private state. Work only in authorized interactions.
 
## Security and quality invariants
- Never print, log, commit, or disclose tokens, cookies, tunnel secrets, personally identifying user data, or private profile contents.
- Keep role authorization on the server, not merely in browser presentation. Treat publisher URLs and fetched content as untrusted; preserve SSRF and request bounds. Do not loosen strict article classification to mask upstream errors.
- Do not modify ScraperGo or other applications during NewsGo work unless separately authorized.
- If a step or integration cannot be verified, report that fact instead of asserting success.

## Documentation ownership
- **README.md:** user-facing setup, operation, features and commands.
- **ARCHITECTURE.md:** implemented technical design, data flow, endpoints, security, persistence.
- **CURRENT_WORK.md:** live development state, evidence, PR/branch/CI, outstanding checks and next tasks.
- **AGENTS.md:** these mandatory development rules alone.

Git history is the record of older work. Do not accumulate dated handoff sections or duplicate content across these files.
