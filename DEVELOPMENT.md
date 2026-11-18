# NewsGo — app-specific development and validation

Read [AGENTS.md](AGENTS.md) and [Homepage's authoritative suite development policy](https://github.com/richard-mauri/homepage/blob/main/SUITE_DEVELOPMENT_POLICY.md). The suite policy, not this file, controls Issues, branches, draft PRs, exact-head CI, visible six-item handoff, owner FINISH and merged branch cleanup. Keep this document for NewsGo-specific validation and commands.

## NewsGo regression and data boundaries

### Mandatory shared navigation button styling

**Canonical implementation:** Homepage now owns `shared/ui/buttons.css` in the nested Go module `github.com/richard-mauri/homepage/shared/ui`. NewsGo pins that module in `go.mod` and serves its embedded CSS from its own origin. The local `dashboard.go` stylesheet may define tokens and button *placement*, but must not reintroduce shared component colors, geometry or hover/focus/disabled styling. Update the dependency through a scoped PR; review the Homepage source and immutable pin, then rerun the UI regression checks. No cross-origin stylesheet or CDN is allowed.

The `.subscription-read` control intentionally behaves visually like an inline underlined link, rather than a standard action button; that explicit exception is not a competing generic-button theme. All ordinary buttons, including compact article exclusion actions, must use the canonical appearance.

All top-level NewsGo navigation buttons, including **Settings**, must use the existing shared `ui-button` and `ui-button-large` styling. Their resting backgrounds, text and borders must match whether selected or not: **no persistent accent fill or accent border** from `aria-pressed="true"`. Apply this rule consistently to every navigation destination; do not introduce a one-off visual treatment for Settings or another section. Preserve accessible pressed-state semantics and visible keyboard focus. Transient hover and `:focus-visible` may highlight the border. UI changes must review the shared selectors in `internal/dashboard/dashboard.go` and the button initialization in `internal/dashboard/personal_script.go`, and retain a regression assertion for the selected-state contract.

The same shared resting appearance applies to **all action buttons**, including Save Server Settings, Reload, Download Article URLs, confirmation-dialog actions, discovery/search, subscriptions, and Diagnostics tabs. No always-on accent border or solid accent/cyan fill may be added for a submit, emphasized, pressed or selected action. All server-rendered buttons must receive `ui-button` at startup; buttons created later must use the shared `button()` helper or explicitly add `ui-button`. Preserve meaningful transient hover/focus, disabled-state behavior and accessible selected-state semantics. Remove obsolete panel-specific color overrides rather than adding higher-priority CSS or `!important` exceptions. Whenever adding a button or editing a button rule, audit panel-specific selectors and update the common-style regression checks.

Protect subscriptions, exclusions, feeds catalog, browser search histories, Homepage account storage and tokens, and per-user settings. Preserve existing publisher URL validation, classification, SSRF limits, import/export and regression coverage. For both agents and the owner, **`./devbuild.sh` is the sole ordinary development test-and-build command**: it already runs `go test ./...` before compiling both Mac binaries. Do not run separate `go test` or `go build` commands before or after it. Normal code changes also require `gofmt` and exact-head GitHub Go CI. If the agent's environment cannot run `./devbuild.sh`, report local test/build as unverified rather than substituting another procedure. A Go CI pass cannot verify live publisher access, Mac Safari, or the tunnel.

## Owner Mac acceptance for runtime PRs

First run `git status --short --branch` and stop if dirty. Use the exact branch from the PR in place of the example:

```bash
cd ~/projects/newsgo
git status --short --branch
# Stop if uncommitted changes are shown.
git fetch origin
git checkout ACTUAL-PR-BRANCH
git pull --ff-only origin ACTUAL-PR-BRANCH
./devbuild.sh
./newsgo.sh
```

`newsgo.sh` restarts NewsGo and its own tunnel, and requires distinct owner authorization. Check feed discovery and subscriptions, exclusions, article reader and export/import, Homepage user/guest handoff, token expiry/revocation, browser history and light/dark UI in Safari. Do not restart for documentation-only PRs. For safe post-FINISH checkout and local branch deletion, use the shared policy's exact clean-tree-guarded template with `~/projects/newsgo`.
