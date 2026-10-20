# NewsGo — Subscription-Based News Reader

NewsGo brings together headlines from news sources **you choose to subscribe to**. The dashboard presents one unified briefing: there is no separate built-in news stream, and no unsolicited headlines from the source catalog.

To begin, open **Find Subscriptions**, browse the available sources or discover a publication's RSS feed, and choose **Subscribe**. Your selections appear in **My Subscriptions** and all their non-excluded headlines appear together in **Today's briefing**. With no subscriptions, the briefing stays empty until you choose sources.

**Excluded Sites** lets you hide publishers by domain without deleting subscriptions. An excluded subscription can be restored with **Allow this site**. Preferences are kept in browser storage and may be exported and imported with the dashboard controls. They are not synchronized automatically between devices.

Headlines open in [ScraperGo](https://scrapergo.themauris.org/), a separate article reader. NewsGo does not alter ScraperGo's extraction behavior.

The service is written in Go, uses an RSS/Atom/JSON feed catalog (`feeds.json`) for discoverability, and presents article headlines from the feeds you have subscribed to. Provider access and publication frequencies vary. NewsGo is an aggregator, not a subscription/paywall bypass.

**Development candidate:** v0.2.32 (October 20, 2026). This version changes the article pipeline to one subscription-driven view and migrates recognized earlier publisher subscriptions to feed-backed subscriptions. Existing browser preferences are preserved. Legacy publisher entries without a matching catalog feed remain visible for manual review. Production deployment and live-browser validation have not been confirmed.

**Developers and future AI agents:** See [`AGENTS.md`](AGENTS.md) for the authoritative current engineering handoff, development procedure, version synchronization rules, and Git safeguards.
