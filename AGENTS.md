# ScraperGo — AI Development Instructions and Chat Handoff

**This file is the authoritative development process and cross-chat handoff for AI coding agents working on ScraperGo. Read it in full before changing code, generating patches, recommending Git commands, or asserting project status.** `README.md` is the end-user/project overview, not the engineering handoff.

## Identity and working environment

- Repository: <https://github.com/richard-mauri/scrapergo>
- Public reader: <https://scrapergo.themauris.org/>
- User's local project: `~/projects/scrapergo`
- Local service: `http://localhost:8083/`
- Language: Go; extraction uses `codeberg.org/readeck/go-readability/v2`.
- Existing build and startup scripts: `./devbuild.sh` and `./scrapergo.sh`. Preserve these and the Cloudflare Tunnel configuration.
- NewsGo is a **separate** RSS application. Never modify NewsGo merely to fix ScraperGo.

## Current handoff (October 20, 2026)

**Last user-confirmed working local version: v0.2.55.** The user verified that previously missing introductory and concluding text, including the gallery images of the Slate Truck showroom, is now recovered. The runtime logs supplied by the user showed:

```text
slate_diagnostic stage=download html_bytes=182095 raw_slate_garage=true visible_slate_garage=true raw_repairpal=true visible_repairpal=true raw_closing=true visible_closing=true raw_pegboard=true raw_interior_models=true
slate_diagnostic stage=extraction blocks=32 images=6 slate_garage=true repairpal=true closing=true pegboard=true interior_models=true
```

The user then said, **“yes it works.”** This confirms content recovery on the actual article, **not** a guarantee that all original formatting or caption deduplication is correct.

**Outstanding work:** Restore semantic article typography. In particular, the original large, bold section heading **“The shopping experience”** appears as normal body text. Readability's `RenderText` output is initially converted to `Block{Kind: "p"}`, which can lose heading identity. Preserve h2/h3 hierarchy, bold/italic runs, inline links, and spacing while retaining the successful v0.2.55 completeness behavior. The proposed next development version is **v0.2.56**, not yet generated or verified in this handoff.

**Git status at last user report:** Local `main` was clean and ahead of `origin/main` by **two commits**:

```text
30ff390 ScraperGo v0.2.55: integrate article content recovery
2ac2f86 ScraperGo v0.2.55: restore complete article content
```

The user was advised to run `git push origin main` but has **not confirmed the result here**. Do not assume GitHub contains the latest local code. Check the current Git state if needed; do not treat this dated handoff as a live status report.

**Temporary diagnostics:** `slate_diagnostic.go` and associated calls were added in v0.2.54 to distinguish retrieval from extraction loss. They should not be removed without deliberate review and user consent when preparing a clean release.

## Source authority and evidence

The user's verified **local working tree outranks GitHub main** when the two diverge. Before generating a patch, establish the exact local baseline and obtain the *current* affected files if necessary. Never overwrite newer working local files using older GitHub versions. Do not infer complete source contents from summaries, a previous patch name, or README entries.

Use precise status language:

- **Generated** — files were created and their actual locations verified.
- **Statically checked** — syntax, formatting, or ZIP integrity checked.
- **Tested** — specified tests actually executed, with their outputs observed.
- **Installed** — the user installed the build.
- **Live-verified** — the user verified actual browser/TXT behavior.
- **Committed / pushed / tagged / released** — distinct events, each requiring evidence.

Never represent synthetic fixture tests, isolated harness tests, source inspection, or successful compilation as full `go test ./...` or real-site verification. Do not assert any downloadable ZIP exists until the exact file has been created and verified in the current environment.

## Critical source integration points

- `main.go`: displayed `appVersion`, HTTP fetch, URL/access controls, legacy parser, and post-extraction stages.
- `readability.go`: Readability pipeline and **integration calls** to leading and continuity recovery. Missing edits here can leave recovery code unused.
- `leading_recovery.go`: recovering missing introductory paragraphs.
- `continuity_recovery.go`: recovering omitted downstream paragraphs and figures; v0.2.55 uses anchored tail continuation.
- `structure.go`: structured article recovery.
- `ui.go`: rendering text, headings, inline links, images, and captions.
- `slate_diagnostic.go`: targeted temporary logs from download and extraction stages.

These filenames describe the known local development lineage. Inspect actual current contents. A feature is **not integrated** merely because its function exists: trace the call path from the live request handler to the function and the resulting rendered block.

## Non-negotiable development workflow

1. **Establish baseline:** identify the current local version, current source, desired fix, evidence, and regression risks. Ask for current affected files only when needed.
2. **Trace the whole path:** review the new/changed functions **and their call sites**. Inspect all interfaces or downstream stages affected. A patch missing a required calling file is incomplete.
3. **Prefer minimal fixes:** preserve working behavior, cross-publisher handling, access restrictions, safe URLs, images, captions, links, browser UI, and startup operations. Do not add publisher-specific hardcoded article text or bypass access controls.
4. **Build meaningful regression tests:** test the reported failure and previously successful behavior; distinguish synthetic fixtures from real fetched HTML.
5. **Run checks where possible:** `go test ./...` in a complete project checkout, plus formatting/build checks as appropriate. Report exact commands actually run, failures, and unavailable checks. If a required test fails, **do not declare success** or suggest proceeding to restart/release.
6. **Verify the complete patch inventory:** compare changed-file list to call sites, tests, version change, and documentation. Include every file necessary to make the fix active.
7. **Advance version:** increment the displayed `appVersion` in `main.go` for each installable candidate, and record the new version in this handoff. Do not conflate development version with public release/tag state.
8. **Package correctly:** a **direct-unzip ZIP**, with only complete changed Go source/test files plus **the complete updated `AGENTS.md` and `README.md` at ZIP root**, without a parent directory or installer scripts. Preserve existing operational scripts unless explicitly requested.
9. **Report precisely:** enumerate files, critical call-site changes, tests actually run, live checks outstanding, and risks. Do not invent successful checks or paths.
10. **After user verification:** update this handoff's current state; replace stale present-tense claims rather than collecting contradictory 'current' versions.

### Standard patch installation

The user deliberately uses a short direct-unzip workflow. Do **not** add Python, shell patch installers, ZIP inspection requirements, or preliminary `git status` to routine **installation**:

```sh
cd ~/projects/scrapergo
unzip -o ~/Downloads/<actual-patch-name>.zip
go test ./...
./devbuild.sh
./scrapergo.sh
```

If `go test ./...` fails, stop before build/restart, diagnose, and regenerate. The existing startup script may not terminate an old server instance; follow the user's established operational restart procedure. Confirm the displayed version and actual browser/TXT extraction before concluding that an article fix works.

### Mandatory article-completeness regression

Test article: <https://arstechnica.com/cars/2026/10/visiting-slates-first-showroom-as-its-25000-electric-truck-nears-production/>

- Opening paragraph begins “As I sat in the cab of the pickup truck ...”; introductory paragraphs formerly omitted must remain.
- Article continues **past** the Slate Garage customization paragraph.
- RepairPal section and final conclusion ending “Now it’s time to deliver.” are included.
- Pegboard and miniature-interior-model photos are retained in a sensible reading position.
- No repeated captions or photographic credits; inline links remain functional.
- “The shopping experience” displays as a real prominent heading, rather than ordinary paragraph text (**known unresolved issue as of v0.2.55**).
- Test other publishers and ensure promotional/ad containers do not cause editorial content loss.

## Git is a separate review and publication phase

**Critical historical mistake:** A previous AI agent gave an explicit `git add` file list that **omitted `readability.go`**, although it contained the live integration calls `recoverLeadingParagraphs(...)` and `recoverEditorialContinuity(...)`. The user detected the omission and committed it separately. This undermined trust and must not recur.

**NEVER compose a `git add` file list based on memory, the ZIP manifest, or an assumed set of changed files.** Git state must be inspected immediately before giving staging instructions. Installation and Git review are intentionally separate workflows.

When asked to commit:

1. Obtain **fresh** `git status`, including staged, unstaged, and untracked entries; check ahead/behind state.
2. Review `git diff --check`, `git diff`, `git diff --cached`, and relevant untracked source files. If there are already local commits, inspect `git log origin/main..HEAD` and the changed-file inventories (`git show --stat` or equivalent).
3. Reconcile changes against actual application integration points: `readability.go`, `main.go`, recovery sources, UI files, tests, and both documentation files, **as applicable**.
4. Identify unrelated work and exclude temporary files (`foo`), logs, generated binaries, secrets, and tunnel configuration.
5. Recommend precise staging **only after** checking the real inventory. Before committing, review `git diff --cached --name-status` and the staged diff for missing or unintended files.
6. Do not assume a clean tree means correct commits or a successful push. Explicitly verify the intended commits and remote status.
7. **Do not commit, push, tag, or release without explicit user direction.** A push and a release tag are separate operations; do not tag development builds automatically.

## Cross-chat handoff procedure

When a new agent receives `AGENTS.md` and `README.md`:

1. Read both, using **this file** for development instructions and current handoff.
2. Confirm the last **user-verified local version**, outstanding issue, last test evidence, and actual local-vs-GitHub state.
3. Obtain current source for affected files before producing a patch if the GitHub state might lag.
4. Preserve v0.2.55's established content-completeness improvement while working on semantic formatting.
5. Do not request screenshots, previously supplied article text, or repeated diagnostics when existing evidence is sufficient.
6. Keep these files coherent for the next agent. Update the **current handoff** above, and keep history below short and factual.

### Compact history and verified lessons

- **v0.2.49:** prior documented GitHub baseline; may no longer reflect remote state.
- **v0.2.50–v0.2.53:** experimental leading and continuity recovery. Initial article paragraphs improved, but real article still truncated after Slate Garage; duplicated captions remained. One v0.2.52 regression test failed in the user's checkout. Passing synthetic tests did not establish live correctness.
- **v0.2.54:** diagnostics showed HTTP 200, complete visible article HTML (`visible_repairpal=true`, `visible_closing=true`), but only 24 extracted blocks / 3 images without the ending/gallery. A separate shell `curl` got HTTP 403; that was **not** evidence of what ScraperGo's HTTP client received.
- **v0.2.55:** anchored-tail correction yielded 32 extracted blocks / 6 images and markers for RepairPal, closing paragraph, pegboard, and interior models. User confirmed the content recovery works. Semantic heading formatting remains defective.
- **Git lesson:** `readability.go` was left out of an AI-generated staging list. Mandatory review of actual local changes and integration call sites was established in response.

**For every future patch and chat, leave clear, accurate evidence—not assurances.**
