# Current work

**Authoritative source:** GitHub `main` is the merged baseline. The previously outstanding feature/fix issues and PRs were closed as of 2026-11-04; do not interpret old handoff notes as active work.

**Active change now:** documentation-only reconciliation on `docs/current-suite-baseline-20261104`, [draft PR #54](https://github.com/richard-mauri/newsgo/pull/54) targeting `main`. The goal is to make the five standard Markdown files accurately describe the current implemented suite. There is no active runtime feature/fix PR in this repository.

**Current implementation:** NewsGo is a separate article-feed client using ScraperGo for reading/classification. Homepage owns the 12-hour account session; NewsGo validates its session-bound handoff and displays the user name without a countdown or Renew control. Browser-only preferences persist on sign-out.

**Next action:** review the documentation diff, check Markdown links and the PR's exact-head CI status, correct any factual errors, then request explicit owner approval to merge this documentation-only PR. Mac/Safari runtime retesting or app restart is not required for text-only changes. No merge, tag, release, deployment, branch deletion or tunnel change has been approved for this PR.
