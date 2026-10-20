# NewsGo — AI Development Process and Chat Handoff

**Read this file first when coding or resuming NewsGo in a new chat.** `README.md` is the NewsGo user guide, not a ScraperGo document.

## Repository identity and critical documentation correction

- Repository: https://github.com/richard-mauri/newsgo
- Local directory: `~/projects/newsgo` on the established Mac setup.
- Local dashboard: `http://127.0.0.1:8084/`.
- Go project with `cmd/newsgo`, `cmd/feedcheck`, `cmd/discoverycheck`, `internal/dashboard`, `internal/feeds`, and `feeds.json`.
- Existing scripts: `./devbuild.sh` and `./newsgo.sh`. The startup script launches NewsGo's own Intel executable and **dedicated Cloudflare Tunnel**, keeping ScraperGo and other tunnels untouched.
- ScraperGo is **a separate article reader** (historically `http://127.0.0.1:8083/` or `https://scrapergo.themauris.org/`). NewsGo sends articles there but does not share its extraction code.

**Documentation corruption on GitHub, observed October 20, 2026:** Both `README.md` **and existing `AGENTS.md`** in the NewsGo repository currently contain **ScraperGo** documents. These replacement files deliberately restore NewsGo identity. Existing `NEWSGO-DASHBOARD.md`, `NEWSGO-v0.2.20.md`, `NEWSGO-DISCOVERY-EVALUATION.md`, and other NewsGo-specific historical files remain useful background; do not delete them just because these two documents are corrected.

## Current handoff / uncertainty — October 20, 2026

The GitHub `NEWSGO-DASHBOARD.md` documents v0.2.0 dashboard work and an older `NEWSGO-v0.2.20.md` states that generic shared headlines are hidden on initial load and shown when an explicit built-in publisher is selected, while personal-subscription headlines load as before. **These notes do not establish the current installed version.** The currently running NewsGo version/build, Git status, last pushed commit and the full list of new features are **not verified here**; check local source, app display and current relevant files before proposing a next build.

Do not assume a NewsGo semantic/build version from ScraperGo's v0.2.55 or from any filename ending in v0.2.20. Preserve FeedCheck CLI and the personal vs built-in feed behavior.

## Functional and security invariants

- Feed subscriptions are stored in `feeds.json`; don't overwrite a user's subscriptions in patches. Preserve existing feed configuration and migration behavior.
- NewsGo merges and sorts headlines from multiple RSS feeds; deduplication and cached stories prevent duplicate display and avoid losing earlier fetched stories when one feed fails.
- Dashboard headline links open ScraperGo by default in the documented integration; **Original** opens the publisher directly. Do not accidentally alter ScraperGo's separate code or URL scheme.
- Preserve separation between personal subscription feeds, built-in publisher/discovery sources, feed filtering and explicit selection. The historical v0.2.20 note identifies intentional hidden-on-initial-load behavior for generic shared headlines.
- Treat publisher feed responses as untrusted input; retain request limits, timeouts, filtering and HTML safety. Existing dashboard notes describe 4 concurrent feed requests, 4 MiB per feed, 12-second timeout, in-memory 600-headline limit and refresh minimum of one minute; **verify current values in source** before treating them as guaranteed.
- Local loopback binding is intentional; no built-in authentication is documented. Don't expose open dashboards or feeds without appropriate access controls.
- Never overwrite `.cloudflared/token`, tunnel paths, `feeds.json`, `newsgo.sh`, `devbuild.sh`, `go.mod`, or unrelated code in a small reader/dashboard patch unless the user explicitly asks and changes are reviewed.

## Build, patch and operational workflow

Use the complete Go package, not one file:

```sh
cd ~/projects/newsgo
go test ./...
./devbuild.sh
./newsgo.sh
```

A source-only development run is documented as `go run ./cmd/newsgo`, after ensuring there isn't a competing service on port 8084. `go run ./cmd/feedcheck` is the separate FeedCheck command. The existing `newsgo.sh` starts the NewsGo application and its own Cloudflare Tunnel; do not edit or restart another project's tunnel.

Patch only necessary changed files, with updated complete `AGENTS.md` and `README.md` when delivering a handoff, and do not replace personal `feeds.json`. Follow the established version/patch convention **after inspecting current source and prior documented packages**; do not infer a version or impose another project's patch rules. Run `go test ./...`, check dashboard load and filters, personal feeds vs default sources, headline links to ScraperGo, Original links, refresh/error handling, and FeedCheck behavior. Separate actual local tests from unverified publisher network behavior.

## Non-negotiable agent conduct

- Read this file and `README.md` in full before changing code, advising Git operations, or claiming completion.
- The user's verified **local source tree** is authoritative whenever it differs from GitHub. Obtain current affected files and reconcile local changes; never overwrite them from older remote sources.
- Trace **both definitions and actual call sites**, handlers, UI uses, tests, and configuration dependencies. A new function that is not called is **not** a finished fix.
- Protect unrelated working behavior. Keep changes minimal and test existing regression areas.
- Distinguish **drafted/generated**, **verified file exists**, **formatted**, **tested in full project**, **tested in isolated harness**, **installed**, **user-verified live**, **committed**, **pushed**, **tagged**, and **released**. Never imply one means another.
- If you cannot run tests or produce a real artifact, say so. Never invent a download link, test result, deployment, or Git action.
- Do not commit, push, tag, release, restart services, deploy, or modify tunnel configuration without the user's explicit authorization.

## Mandatory Git pre-commit audit — learned from ScraperGo failure

**Never compose a `git add` filename list from memory, an earlier ZIP manifest, or an assumption about changed files.** A previous agent omitted a crucial `readability.go` integration change in ScraperGo; the user found it before pushing. This is a cross-project safeguard.

For a commit request, **separately from normal patch installation**:

1. Obtain fresh `git status`, including staged, unstaged and untracked files, and ahead/behind state.
2. Examine `git diff --check`, `git diff`, `git diff --cached`, and relevant untracked files. If commits already exist locally, review `git log origin/main..HEAD` and commit file inventories.
3. Check all changed source files **and integration points**, tests, both documentation files, and build/configuration artifacts as applicable. Exclude secrets, local tokens, binaries, logs and unrelated temporary files.
4. Only after that evidence recommend exact staging commands. Before committing, verify `git diff --cached --name-status` and `git diff --cached` against the intended scope.
5. Do not infer that a clean working tree means the correct commits were made, or that ahead-of-origin commits were pushed. Verify publication independently.

## Cross-chat handoff obligations

At the start of a new chat, use `AGENTS.md` for the **current engineering handoff** and `README.md` for the public-facing project context. Confirm current local source and Git state before applying edits. Maintain this file with a short, dated snapshot of the most recent verified baseline, pending work, confirmed tests, known regressions and next action; **replace obsolete current-state claims**, rather than accumulating contradictory updates. Keep the README focused on users and general architecture, not chat history or agent rules. Preserve the existing project's patch/install convention; do not transplant another project's conventions blindly.

