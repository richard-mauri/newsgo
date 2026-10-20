# ScraperGo — Clean Web Article Reader

**Read web articles without the surrounding clutter.** ScraperGo is a self-hosted web reader that retrieves available article HTML and presents the main story in a simplified layout.

**Public reader:** <https://scrapergo.themauris.org/>  
**Source:** <https://github.com/richard-mauri/scrapergo>

> **For AI coding agents and future chat sessions:** Read [`AGENTS.md`](AGENTS.md) before changing code or proposing development/Git commands. `AGENTS.md` contains the development workflow, current chat handoff, regression history, patch conventions, and mandatory checks. This README is written primarily for readers and people learning what the project does.

## What it does

ScraperGo aims to retain the parts of a web article that matter while removing unrelated navigation, advertising, and promotional material. Depending on the publisher's HTML, the reading view can include the headline, author, paragraphs, article images, captions, section headings, and hyperlinks.

The application also offers a **TXT download** of the extracted article. This is useful for saving readable text and checking whether an article's paragraphs were actually captured.

ScraperGo is not a browser emulator or a subscription bypass. It cannot guarantee access to articles requiring logins, subscriptions, JavaScript-only rendering, or publishers that reject automated requests. Availability and completeness depend on the HTML actually returned by the publisher.

## Using the reader

Open <https://scrapergo.themauris.org/> and enter the web address of an article. ScraperGo retrieves the publicly available document and displays a simplified reading view.

If the result is incomplete or formatted incorrectly, compare it with the original publisher page. The **TXT download** helps narrow down the problem: content missing in both the TXT file and reading view may have been lost during retrieval or extraction; content present in TXT but missing on-screen points more directly to presentation.

Publisher images and links are preserved where possible, but their positions and styles may differ from the original website.

## Reading quality and known limitations

Websites structure articles in different ways. Some put editorial paragraphs inside containers shared with advertisements; some use complex photo galleries or custom heading markup. ScraperGo attempts to identify the complete article without including unrelated page elements, but extraction can still be imperfect.

As of the **locally verified v0.2.55 development build (October 20, 2026)**, ScraperGo successfully recovered the full text and previously missing showroom photographs in a long Ars Technica article that had been truncated in earlier builds. The user confirmed this improvement. **Formatting fidelity remains incomplete**: for example, a large bold section heading, “The shopping experience,” is currently displayed like an ordinary paragraph. Heading and emphasis preservation is an active improvement area.

That local development status does **not** imply v0.2.55 has been pushed, tagged, or publicly deployed; consult the repository's releases for published versions.

## Where ScraperGo runs

The application is written in **Go** and serves its reading interface over HTTP. In the established local setup, it listens at:

`http://localhost:8083/`

The public HTTPS endpoint is:

<https://scrapergo.themauris.org/>

The public address uses the existing Cloudflare Tunnel. The tunnel is a routing mechanism and should not, by itself, be treated as an authentication system.

This repository has existing scripts for building and running the service (`devbuild.sh` and `scrapergo.sh`). Coding agents and maintainers should use the development and operational instructions in [`AGENTS.md`](AGENTS.md) rather than relying on ad hoc setup or release commands.

## Relationship to NewsGo

[NewsGo](https://github.com/richard-mauri/newsgo) is a separate RSS/news aggregation project. ScraperGo can provide a cleaner view of articles discovered through NewsGo, but ScraperGo itself is not an RSS reader. The two projects have separate codebases and development processes.

## Security and publisher access

ScraperGo processes web URLs with request restrictions and publisher-access handling. It is designed to display content made available by the publisher, not to circumvent paywalls or account requirements. Publisher copyright and usage terms still apply to the source articles and their photographs.

## Reporting a problem

The most useful details are the **article URL**, the **ScraperGo version shown in the interface**, what is **missing or incorrectly formatted**, and whether it also appears in the downloaded **TXT** file. A comparison with the original article is helpful for issues involving images, captions, section headings, or missing ending paragraphs.

For engineering diagnostics, previously verified regressions, development patch procedures, Git checks, and cross-chat continuity, see [`AGENTS.md`](AGENTS.md).
