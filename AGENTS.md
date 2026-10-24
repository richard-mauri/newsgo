# NewsGo — development agent rules

This file is the **authoritative instruction set for agents**. Read it before making changes. Do not put feature history, architecture explanations, or active task notes here: see [README.md](README.md), [ARCHITECTURE.md](ARCHITECTURE.md), and [CURRENT_WORK.md](CURRENT_WORK.md).

## Source of truth and owner authority
- GitHub `richard-mauri/newsgo` branch `dev` is the development integration branch. Read the actual current files on `dev` and check the current PR and CI state before editing. `main` is the owner-controlled release branch.
- During an authorized development request, commit focused, tested changes to `dev` and maintain the draft `dev` → `main` PR. Never merge, tag, publish a release, rewrite branch history, force-push, or change `main` without explicit owner approval.
- Respect the owner's stated scope; perform work only in the active interaction. No implied asynchronous follow-up.
- Preserve the local runtime environment and data: `data/`, tokens, subscriptions, exclusions, reading preferences, `feeds.json`, Cloudflare credentials/tunnels, ports, launcher behavior, and existing built artifacts. Do not silently reset or migrate profiles.

## Required workflow
1. Read this file, the other three documentation files, and the relevant source/tests. Confirm the current `dev` head, open PR, and constraints. State the intended workflow to the owner before editing.
2. Implement the smallest coherent change with regression tests for the failure. Format touched Go files and update documentation in its **single owning file**. Do not put chronological release notes in README or AGENTS.
3. Push to `dev` as authorized. Inspect the GitHub Actions workflow at the **exact final head**. It executes `go test -count=1 ./...` and `go build ./...`; resolve any failure and rerun. Green CI does not substitute for live Safari/integration acceptance.
4. Update the draft PR **title and full description** with latest code, actual CI outcome, Mac acceptance status, open problems, and explicit no-merge status; verify the update took effect. Keep [CURRENT_WORK.md](CURRENT_WORK.md) synchronized with current state, including pending owner acceptance.
5. Give the owner complete copy/paste Mac acceptance commands with the **actual branch**, including tests, build, and launcher:
   ```bash
   cd ~/projects/newsgo
   git fetch origin
   git checkout dev
   git pull --ff-only origin dev
   go test -count=1 ./...
   ./devbuild.sh
   ./newsgo.sh
   ```
   The owner's Mac uses older Git: use `git checkout`, not `git switch`. Build with prebuilt Go tooling; avoid unrelated heavyweight installs.
6. Ask the owner to verify observable changes in Safari, live feeds, sync, classification, images, and tunnel behavior as applicable. Never claim an unperformed live test succeeded.
7. Only after explicit owner promotion: check out and pull `main`, run tests and rebuild with existing scripts, and verify the promoted app.

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
