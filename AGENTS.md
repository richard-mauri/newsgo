# NewsGo — development agent rules

This file is the **authoritative instruction set for agents**. Read it before making changes. Do not put feature history, architecture explanations, or active task notes here: see [README.md](README.md), [ARCHITECTURE.md](ARCHITECTURE.md), [DEVELOPMENT.md](DEVELOPMENT.md), and [CURRENT_WORK.md](CURRENT_WORK.md).

## Mandatory development entry point

Read [CURRENT_WORK.md](CURRENT_WORK.md), [ARCHITECTURE.md](ARCHITECTURE.md) and [DEVELOPMENT.md](DEVELOPMENT.md) before any development action. **DEVELOPMENT.md is authoritative** for Issues, branch creation, draft PRs, tests, Mac/Safari acceptance, promotion gates and post-merge branch cleanup. GitHub `main` is the authoritative merged baseline; use a fresh issue-scoped branch for every new task, **never the legacy `dev` branch**.

Before making changes, tell the owner the verified baseline, Issue and branch, actual PR state, tests and acceptance path, and merge authorization boundary. Work only within the authorized scope. Preserve user data, tokens, subscriptions, exclusions, preferences, local build counter, feeds catalog and Cloudflare settings.

Never merge a PR, modify `main`, tag or release without the owner's explicit approval. **Normal cleanup of a merged task's short-lived branch follows the safety checks in DEVELOPMENT.md; bulk deletion of historical branches requires separate approval.** No forced deletion or force-pushing.

Before handoff, update the actual PR title and full description and CURRENT_WORK.md, reporting exact final-head CI and what was or wasn't tested on the Mac. For code work, give the complete branch-specific Mac command block prescribed in DEVELOPMENT.md. Documentation-only PRs need documentation checks and CI, not runtime acceptance.

## Security and quality invariants
- Never print, log, commit, or disclose tokens, cookies, tunnel secrets, personally identifying user data, or private profile contents.
- Keep role authorization on the server, not merely in browser presentation. Treat publisher URLs and fetched content as untrusted; preserve SSRF and request bounds. Do not loosen strict article classification to mask upstream errors.
- Do not modify ScraperGo or other applications during NewsGo work unless separately authorized.
- If a step or integration cannot be verified, report that fact instead of asserting success.

## Documentation ownership
- **README.md:** user-facing setup, operation, features and commands.
- **ARCHITECTURE.md:** implemented technical design, data flow, endpoints, security, persistence.
- **CURRENT_WORK.md:** live development state, evidence, PR/branch/CI, outstanding checks and next tasks.
- **DEVELOPMENT.md:** Git workflow, validation and post-merge cleanup.
- **AGENTS.md:** these mandatory agent rules alone.

Git history is the record of older work. Do not accumulate dated handoff sections or duplicate content across these files.
