# NewsGo — RSS Feed Dashboard

NewsGo is a lightweight, self-hosted news aggregator built in Go. It collects headlines from configured RSS feeds, presents them in a local browser dashboard and offers an easy path to open an article in **ScraperGo**, a separate clean-article reader.

**Source:** https://github.com/richard-mauri/newsgo  
**Local dashboard (documented default):** http://127.0.0.1:8084/

> **For AI coding agents and future development chats:** Read [AGENTS.md](AGENTS.md) for current project status, source-baseline verification, test/patch workflow, Git safeguards and the development handoff. This README is for NewsGo users and general project context.

## What NewsGo does

NewsGo gathers headlines from RSS news sources, orders available stories by publication date and presents a browsing dashboard. It deduplicates the same article URL and may retain previously fetched headlines when an individual feed has a temporary retrieval failure.

It supports source filtering and personal feed configuration through the project's `feeds.json` file. Some built-in/shared publisher headlines are intentionally not shown on initial load and appear only after selecting the corresponding built-in publisher; personal-subscription headlines follow their own configured behavior. Exact options depend on the running build.

Click an article headline to open it in [ScraperGo](https://scrapergo.themauris.org/), which attempts to show its publicly available article text in a simplified reader. The **Original** link opens the publisher's own page. NewsGo and ScraperGo are separate applications: an RSS feed is not the same as a full publisher article.

## Running locally

NewsGo is a Go application. In the documented Mac setup, from the local repository:

```sh
cd ~/projects/newsgo
go run ./cmd/newsgo
```

Then open http://127.0.0.1:8084/ in Safari. The separate FeedCheck command can be started with `go run ./cmd/feedcheck` when checking source feeds. The project also contains a `discoverycheck` command and established build/start scripts for its normal hosted setup.

The dashboard may briefly be empty while initial feeds are fetched asynchronously. Upstream RSS hosts can be slow, unavailable or rate-limit requests; a missing headline does not necessarily mean NewsGo is malfunctioning.

## Feed configuration and operations

Configured feeds are kept in `feeds.json`. Preserve existing subscriptions when upgrading. The dashboard caches a limited number of stories in memory, so a service restart triggers fresh fetching. The established local server uses loopback by default and is not intended as an unauthenticated public service.

NewsGo's operational Mac setup uses an application binary and its **own Cloudflare Tunnel**, independent of ScraperGo's tunnel. Tunnel availability and URL configuration depend on the server setup; use the existing project scripts rather than copying procedures from another application.

## Relationship to ScraperGo

[ScraperGo](https://github.com/richard-mauri/scrapergo) retrieves available web-article HTML and attempts to preserve the story's paragraphs, photographs, headings and links. It is separate from NewsGo. Publisher access controls, JavaScript-only content and feed-summary limitations can prevent either application from showing all article text.

## Data limitations and security

RSS data and publication timestamps come from external providers; content, availability and accuracy vary. Neither NewsGo nor ScraperGo bypasses paywalls. If you expose the dashboard beyond localhost, protect it with appropriate access controls.

For engineering history, development changes and cross-chat continuity, consult [AGENTS.md](AGENTS.md). Historical `NEWSGO-*.md` files may describe older milestones and should not be mistaken for the current software version.

## Excluding publishers from your news view (v0.2.26 development candidate)

NewsGo can hide headlines from publishers you choose not to see, including sites whose articles require subscriptions. Use **Excluded Sites** to add, review or remove website domains. You can also select **Exclude this site** next to a headline, a discovered feed, or a subscribed source. Excluded domains and their subdomains are hidden from NewsGo's shared and personal headline views. Discovery and My Subscriptions continue to display the source with an **Excluded** indicator, so you can identify and reverse the setting.

Excluding a site **does not unsubscribe** from its feeds. Allowing it again restores visibility when the next feed data is available. The setting applies only in the current browser and is saved in its local storage; it does not change ScraperGo or the original publisher. If you clear the browser's site data, these preferences may be lost.

**Export subscriptions** now downloads a preferences file containing subscriptions and the exclusion list; **Import subscriptions** accepts both the new format and older subscription-only JSON arrays. Importing an older file retains the browser's existing exclusions. A source supplied through a third-party feed can list a different feed domain from the actual article publisher: headlines use the article's destination URL for filtering.

This is documented as a **development candidate** until the current build has been installed and the browser behavior has been verified. The developer workflow and current handoff are in [AGENTS.md](AGENTS.md).


## NewsGo v0.2.27 — Follow-up to failed v0.2.26 candidate (October 20, 2026)

The v0.2.26 direct-unzip candidate failed seven dashboard tests in the user's full checkout. Several legacy tests asserted the previous v0.2.25 HTML version string; the new exclusion test made an overly exact template-tree assertion; the combined-subscription test asserted an original shared-headline selection expression. This follow-up preserves the original selection expression and applies publisher exclusions using a separate `eligible` predicate. It also activates shared-headline exclusion controls in the render path and updates old version assertions to v0.2.27. No claim of end-to-end verification until `go test ./...` succeeds on the user's checkout and the browser flows are manually checked. The code belongs solely to NewsGo. Do not change ScraperGo.


## Publisher exclusions (v0.2.28 development candidate)

NewsGo includes an **Excluded Sites** list for hiding articles from specified publisher domains, such as subscription-only sites. This is a NewsGo preference, not a ScraperGo restriction. Find Subscriptions and My Subscriptions mark excluded publishers without canceling existing subscriptions. A user can remove an exclusion to restore articles; exclusions persist in this browser and can travel with exported preferences. Older subscription-only exports are still accepted. The v0.2.28 build has not yet been confirmed in the local Safari dashboard.


### Excluded publisher navigation (v0.2.29 development candidate)

An excluded publisher remains subscribed and labeled as excluded, but its **Read Articles** action is hidden. When a publisher is excluded while its articles are selected, NewsGo returns to **All Subscriptions** so the reader isn't left on an empty page. Removing the exclusion restores the publisher's reading action and articles. This is a local development candidate; browser verification is still required.

### Version v0.2.30 (development candidate)

This maintenance update aligns the displayed NewsGo version, browser JavaScript cache-busting reference, and version-dependent tests while retaining the v0.2.29 excluded-publisher navigation behavior. It requires a successful local test run and browser verification before release.

## v0.2.31 — Publisher exclusion control refinement

In My Subscriptions, both shared and personal followed sources present **Exclude this site** when allowed and **Allow this site** when excluded. An excluded subscription remains subscribed, cannot be opened through Read Articles while excluded, and can be restored directly from its row or the Excluded Sites panel. The dashboard displays v0.2.31. Development candidate pending full user-side Go tests and Safari verification.
