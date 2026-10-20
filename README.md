**v0.2.41:** Corrected the visible Server Sync token-entry placeholder to “Enter the 64-character approved token distributed by Richard.” No change to token handling, sync, subscriptions, or exclusions.

**v0.2.41:** Presentation-only changes: header tagline is “NEWS READING WITHOUT PAGE CLUTTER USING SCRAPERGO” and invalid/absent sync token guidance reads “Enter the 64-character approved token distributed by Richard.” No change to synchronization, subscription data, exclusions, or ScraperGo.

**v0.2.38:** Removes the redundant upper-right “Updated subscribed feeds” status. The detailed line showing updated feed counts and visible headlines remains the single source of refresh progress. No changes to subscriptions, exclusions, or server sync.

**v0.2.37:** Removes the redundant standalone ← All Subscriptions button from the selected-feed reading view. Use **My Subscriptions → Read All Articles** to return to the combined feed. Existing feed selection and exclusion fallback still use `chooseAllSubscriptions()` internally. No changes to browser/server persistence, sync, or ScraperGo.

**v0.2.36:** Excluded Sites now shows the saved publisher name(s) (for example, Hacker News) alongside the underlying excluded domain (`news.ycombinator.com`), and refreshes the panel when opened. Domain keys, counts, preferences, and sync APIs are unchanged.

# NewsGo — Subscription-Based News Reader

NewsGo brings together headlines from news sources **you choose to subscribe to**. The dashboard presents one unified briefing: there is no separate built-in news stream, and no unsolicited headlines from the source catalog.

To begin, open **Find Subscriptions**, browse the available sources or discover a publication's RSS feed, and choose **Subscribe**. Your selections appear in **My Subscriptions** and all their non-excluded headlines appear together in **Today's briefing**. With no subscriptions, the briefing stays empty until you choose sources.

**Excluded Sites** lets you hide publishers by domain without deleting subscriptions. An excluded subscription can be restored with **Allow this site**. Preferences are kept in browser storage and may be exported and imported with the dashboard controls. They are synchronized across devices only when optional Server Sync is connected.

Headlines open in [ScraperGo](https://scrapergo.themauris.org/), a separate article reader. NewsGo does not alter ScraperGo's extraction behavior.

The service is written in Go, uses an RSS/Atom/JSON feed catalog (`feeds.json`) for discoverability, and presents article headlines from the feeds you have subscribed to. Provider access and publication frequencies vary. NewsGo is an aggregator, not a subscription/paywall bypass.

**Development candidate:** v0.2.41 (October 20, 2026). This version adds optional per-user JSON-backed synchronization of subscriptions and excluded sites. The prior unified subscription reader and feed-discovery fixes remain intact. Existing browser preferences are preserved. Production deployment and live-browser validation have not been confirmed.

**Developers and future AI agents:** See [`AGENTS.md`](AGENTS.md) for the authoritative current engineering handoff, development procedure, version synchronization rules, and Git safeguards.

### v0.2.33 discovery and cleanup recap

Removed the obsolete server-side catalog headline fetching and image enrichment path. Improved Feedly publisher-name lookup with a prefix fallback when no results are returned, clearer external-feed HTTP error messages, and discovery result card layout. The feed catalog remains a source of discoverable subscriptions, not an automatic headline stream.


## Optional server synchronization (v0.2.34)

The NewsGo server can persist **both subscriptions and excluded sites** per approved user token.
The same token works across MacBook Safari, Docked Safari, and iPhone Safari.
Token entry is optional; browser-only operation remains available.

The server defaults to `-sync-dir data` relative to its working directory. Do not commit this directory; it holds user preferences and token digests. Back it up securely.

Issue the first token on the MacBook running NewsGo:

```sh
go run ./cmd/newsgo-token -dir data -user richard -label 'Primary user'
```

Copy the printed secret token into **Server Sync** on each trusted browser once. The first connection automatically uploads both local lists if no server profile exists; existing server profiles load on Connect. Both subscriptions and excluded sites save automatically after edits. Other devices check for updates every 30 seconds while visible and when returning to the tab. **Advanced** contains optional Load from Server, Save Now and Disconnect controls. Stale revisions cannot overwrite newer server data; after a conflict use Advanced > Load from Server (with confirmation when unsaved local changes exist).

Revoke a user token: `go run ./cmd/newsgo-token -dir data -user richard -revoke`. Reissue after revocation if necessary. Tokens are stored as SHA-256 digests; each token should be treated as a secret. The browser remembers the token in localStorage for automatic reconnection. Only enter tokens on trusted devices: browser storage is accessible to scripts running on this origin. Disconnect erases the locally remembered token, not the server account. The server must be reachable over HTTPS through your existing Cloudflare Tunnel. No additional hosted database is required.


### v0.2.41: safe reconnection

When a valid token already has preferences on the server and the local browser lists differ, NewsGo offers **Use Server Preferences**, **Merge Preferences**, and **Cancel**. The decision applies to both subscriptions and excluded sites. Merge uses subscription IDs and unique excluded domains, and saves against the server revision; stale writes still receive HTTP 409. A new server profile is initialized from browser preferences as before. This applies to remembered-token reconnects as well as manual Connect.
