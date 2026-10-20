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
