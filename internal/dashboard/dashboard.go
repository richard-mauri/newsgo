package dashboard

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/richard-mauri/newsgo/internal/feeds"
)

//go:embed tanner.jpg
var tannerPhoto embed.FS

type State struct {
	sources   []feeds.Source // Discovery catalog only, never fetched automatically.
	reader    string
	syncStore *SyncStore
}

func New(sources []feeds.Source, reader string, _ *http.Client) (*State, error) {
	if err := feeds.ValidateSources(sources); err != nil {
		return nil, err
	}
	u, err := url.Parse(reader)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil {
		return nil, fmt.Errorf("invalid ScraperGo address")
	}
	return &State{sources: append([]feeds.Source(nil), sources...), reader: strings.TrimRight(reader, "/")}, nil
}

func (s *State) SetSyncStore(store *SyncStore) { s.syncStore = store }

type view struct {
	Catalog []feeds.Source
 BuildLabel string
}

func (s *State) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/tanner.jpg" {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		data, err := tannerPhoto.ReadFile("tanner.jpg")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = w.Write(data)
		return
	}
	if r.URL.Path == "/api/admin/classification-diagnostic" { s.serveAdminClassificationDiagnostic(w,r); return }
	if r.URL.Path == "/api/admin/settings" {
		if s.syncStore == nil {
			http.Error(w, "Sync unavailable", 503)
		} else {
			s.syncStore.ServeAdminSettings(w, r)
		}
		return
	}
	if r.URL.Path == "/api/sync" {
		if s.syncStore == nil {
			http.Error(w, "Sync unavailable", 503)
		} else {
			s.syncStore.ServeSync(w, r)
		}
		return
	}
	if r.URL.Path == "/api/discover" {
		s.serveDiscover(w, r)
		return
	}
	if r.URL.Path == "/api/search-feeds" {
		s.serveFeedSearch(w, r)
		return
	}
	if r.URL.Path == "/api/diagnostic/feed-trace" { s.serveFeedTrace(w,r); return }
	if r.URL.Path == "/api/personal" {
		s.servePersonal(w, r)
		return
	}
	if r.URL.Path == "/personal.js" {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")
		io.WriteString(w, personalJS)
		return
	}
	if r.URL.Path == "/health" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "ok")
		return
	}
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	v := view{Catalog: append([]feeds.Source(nil), s.sources...), BuildLabel: BuildLabel()}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; base-uri 'none'; form-action 'self'; script-src 'self'; connect-src 'self'; img-src 'self' https:")
	if err := page.Execute(w, v); err != nil {
		return
	}
}

var page = template.Must(template.New("index").Parse(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><meta name="color-scheme" content="light dark"><title>NewsGo — Your daily briefing</title>
<style>
:root{color-scheme:light dark;--bg:light-dark(#f5f5f0,#0e141a);--panel:light-dark(#fff,#19232b);--text:light-dark(#182631,#eef2f4);--muted:light-dark(#60707b,#a5b3bb);--line:light-dark(#dce2e1,#34434a);--accent:light-dark(#126877,#69c6d3);--soft:light-dark(#e6f3f3,#203c43);font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif}
*{box-sizing:border-box}body{margin:0;background:var(--bg);color:var(--text)}a{color:inherit}main{max-width:1200px;margin:0 auto;padding:0 24px 72px}
.topbar{border-bottom:1px solid var(--line);background:var(--panel)}.topbar-inner{max-width:1200px;margin:auto;padding:19px 24px;display:flex;justify-content:space-between;align-items:center;gap:14px}.brand{font-weight:850;font-size:1.75rem;letter-spacing:-.075rem}.brand span{color:var(--accent)}.tagline{font-size:.82rem;color:var(--muted);letter-spacing:.03em}
.header-right{display:flex;align-items:center;gap:22px}.tanner{display:flex;align-items:center;gap:8px;color:var(--muted);font-size:.82rem;font-weight:700}.tanner img{width:64px;height:64px;object-fit:cover;object-position:center 35%;border-radius:12px;border:1px solid var(--line)}@media(max-width:540px){.header-right{gap:6px}.tanner img{width:49px;height:49px}.tanner span{display:none}}
.eyebrow{font-size:.73rem;font-weight:780;letter-spacing:.16em;text-transform:uppercase;color:var(--accent)}.intro{padding:38px 0 25px;display:flex;justify-content:space-between;align-items:end;gap:18px}.intro h1{font-family:Georgia,serif;font-size:clamp(2.6rem,5vw,4.6rem);line-height:1;margin:9px 0 8px;font-weight:700;letter-spacing:-.035em}.intro p{margin:0;color:var(--muted);font-size:.93rem}.time{color:var(--muted);font-size:.85rem;text-align:right}
.filters{padding:17px 19px;background:var(--panel);border:1px solid var(--line);border-radius:15px;display:flex;align-items:end;gap:15px;flex-wrap:wrap;margin-bottom:31px}.filters label{display:flex;flex-direction:column;gap:6px;font-size:.74rem;font-weight:720;color:var(--muted);text-transform:uppercase;letter-spacing:.065em}.filters select{min-width:175px;max-width:100%;font:inherit;font-size:.92rem;color:var(--text);background:var(--bg);border:1px solid var(--line);border-radius:9px;padding:10px 12px}.filter-btn{background:var(--accent);color:light-dark(#fff,#082028);padding:11px 20px;border-radius:9px;border:0;font-size:.89rem;font-weight:700;cursor:pointer}.clear{font-size:.88rem;color:var(--muted);padding-bottom:10px;text-decoration:underline;text-underline-offset:3px}
.section-label{display:flex;align-items:center;gap:14px;font-weight:770;font-size:.77rem;letter-spacing:.13em;text-transform:uppercase;margin:25px 0 17px}.section-label:after{content:"";height:1px;flex:1;background:var(--line)}
.feature-grid{display:grid;grid-template-columns:1.35fr 1fr 1fr;gap:16px}.feature{min-width:0;min-height:282px;display:flex;flex-direction:column;justify-content:space-between;gap:28px;background:var(--panel);border:1px solid var(--line);border-radius:17px;padding:23px;box-shadow:0 10px 24px #00000008}.feature:first-child{background:linear-gradient(130deg,light-dark(#d7eeeb,#16434a),var(--panel) 85%);border-color:light-dark(#c5dddd,#31585f)}.story-meta{color:var(--muted);font-size:.77rem;line-height:1.6}.source{font-weight:800;text-transform:uppercase;letter-spacing:.075em;color:var(--accent);font-size:.72rem}.feature .story-title{font-family:Georgia,serif;font-size:clamp(1.15rem,1.65vw,1.57rem);line-height:1.22;font-weight:700;letter-spacing:-.015em;text-decoration:none;display:block;margin:10px 0 12px}.feature:first-child .story-title{font-size:clamp(1.5rem,2.3vw,2.2rem)}.story-title:hover{text-decoration:underline;text-underline-offset:3px}.original{color:var(--muted);font-size:.79rem;text-decoration:underline;text-underline-offset:3px}.feature .original{display:inline-block;margin-top:7px}
.latest{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:0 28px}.item{display:flex;gap:17px;padding:21px 0;border-bottom:1px solid var(--line);min-width:0}.item-num{font-family:Georgia,serif;font-size:1.6rem;color:var(--muted);opacity:.55;min-width:30px}.item-body{min-width:0}.item .story-title{display:block;font-family:Georgia,serif;font-size:1.23rem;font-weight:700;line-height:1.32;margin:7px 0 7px;text-decoration:none}.item .story-title:hover{text-decoration:underline}.item .original{margin-left:8px}.notice{border:1px solid var(--line);border-radius:9px;padding:12px;margin:12px 0;color:var(--muted)}footer{margin-top:35px;padding-top:17px;border-top:1px solid var(--line);font-size:.8rem;color:var(--muted)}
@media(max-width:820px){.feature-grid{grid-template-columns:1fr 1fr}.feature:first-child{grid-column:1/-1;min-height:260px}.feature{min-height:235px}.latest{grid-template-columns:1fr}.intro{padding-top:28px}}
@media(max-width:540px){main{padding:0 16px 50px}.topbar-inner{padding:15px 16px}.tagline{display:none}.intro{display:block}.time{text-align:left;margin-top:12px}.filters{gap:12px}.filters label{flex:1;min-width:130px}.filters select{width:100%;min-width:0}.feature-grid{grid-template-columns:1fr}.feature:first-child{grid-column:auto}.feature{min-height:200px}.item{padding:18px 0}}

.catalog-search-form{display:flex;align-items:stretch;gap:9px}.catalog-search-form .catalog-search{flex:1;min-width:0}.catalog-search-form button{cursor:pointer;background:var(--soft);border:1px solid var(--line);color:var(--text);border-radius:10px;padding:10px 20px;font:inherit;font-weight:650}.catalog-search-form button:disabled{opacity:.55;cursor:wait}.catalog-search{width:100%;padding:12px;border:1px solid var(--line);background:var(--bg);color:var(--text);border-radius:10px;font:inherit}.catalog-list{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:8px 15px;margin:14px 0}.catalog-card{padding:13px;border:1px solid var(--line);border-radius:10px;display:flex;align-items:center;justify-content:space-between;gap:12px}/* Keep long publisher names and controls in their own grid tracks. */
.catalog-card{display:grid;grid-template-columns:minmax(0,1fr) auto;align-items:start}
.catalog-card>div:first-child{min-width:0;overflow-wrap:anywhere;grid-row:1 / span 2}
.catalog-card>button{grid-column:2;min-width:0;max-width:100%;white-space:normal}
@media(max-width:540px){.catalog-card{grid-template-columns:minmax(0,1fr)}.catalog-card>div:first-child{grid-row:auto}.catalog-card>button{grid-column:1}}
.personal-row>div:first-child{min-width:0;overflow-wrap:anywhere}.personal-row>button{flex-shrink:0}
.catalog-card strong{display:block}.catalog-card small{color:var(--muted)}.catalog-card button{cursor:pointer;background:var(--soft);border:1px solid var(--line);color:var(--text);border-radius:8px;padding:8px 13px}.catalog-card button:disabled{opacity:.6;cursor:default}@media(max-width:720px){.catalog-list{grid-template-columns:1fr}}.personal-tools{display:flex;align-items:center;justify-content:space-between;gap:12px;margin-bottom:20px;flex-wrap:wrap}.personal-tools button,.personal-form button,.personal-row button{cursor:pointer;background:var(--soft);border:1px solid var(--line);color:var(--text);border-radius:9px;padding:9px 13px;font:inherit}.personal-form{display:flex;gap:10px;flex-wrap:wrap;margin:12px 0}.personal-form input,.personal-form select{flex:1;min-width:170px;border:1px solid var(--line);background:var(--bg);color:var(--text);border-radius:9px;padding:10px;font:inherit}.personal-panel{padding:18px;border:1px solid var(--line);border-radius:15px;background:var(--panel);margin:14px 0 25px}.personal-row{display:flex;align-items:center;justify-content:space-between;border-top:1px solid var(--line);padding:10px 0;gap:10px}.personal-row small{color:var(--muted);overflow-wrap:anywhere}.personal-status{color:var(--muted);font-size:.86rem}.personal-status.is-refreshing::before{content:"";display:inline-block;width:13px;height:13px;margin-right:8px;vertical-align:-2px;border:2px solid currentColor;border-right-color:transparent;border-radius:50%;animation:newsgo-spin .75s linear infinite}@keyframes newsgo-spin{to{transform:rotate(360deg)}}@media(prefers-reduced-motion:reduce){.personal-status.is-refreshing::before{animation:none;border-right-color:currentColor}}.personal-list{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:0 24px}.personal-suggestion{margin:7px 8px 7px 0}.personal-hidden{display:none}.personal-panel h2{margin:0 0 9px}.personal-panel p{color:var(--muted)}@media(max-width:720px){.personal-list{grid-template-columns:1fr}}
.included-list{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:10px}.included-card{display:flex;align-items:center;justify-content:space-between;gap:12px;padding:12px;border:1px solid var(--line);border-radius:10px}.included-card strong{overflow-wrap:anywhere}.included-card button{flex-shrink:0;cursor:pointer;background:var(--panel);color:var(--text);border:1px solid var(--line);border-radius:8px;padding:8px 12px;font:inherit}.browse-included h3{margin:20px 0 4px}.discovery-divider{font-weight:700;margin:22px 0 10px;color:var(--muted)}@media(max-width:720px){.included-list{grid-template-columns:1fr}}.filters{margin:0 0 22px}[hidden]{display:none!important}.filters label{text-transform:none}.reading-age-row{display:flex;align-items:center;gap:14px;flex-wrap:wrap;margin:20px 0 8px}.reading-age-row label{font-weight:700}.reading-age-row select{font:inherit;color:var(--text);background:var(--bg);border:1px solid var(--line);border-radius:9px;min-width:215px;max-width:100%;padding:11px 13px}.reading-age-row select:focus-visible{outline:2px solid var(--accent);outline-offset:2px}
.news-nav{display:flex;gap:10px;flex-wrap:wrap;margin:0 0 22px}.news-nav button,.panel-close{cursor:pointer;background:var(--panel);border:1px solid var(--line);color:var(--text);border-radius:10px;padding:12px 17px;font:inherit;font-weight:700}.news-nav button[aria-pressed="true"]{background:var(--accent);color:light-dark(#fff,#082028)}.panel-heading{display:flex;justify-content:space-between;align-items:center;gap:12px;margin-bottom:10px}.panel-heading h2{margin:0}.panel-heading .panel-close{padding:8px 13px}.subscription-panel{margin:0 0 24px}.subscription-panel .personal-row small{display:block;margin-top:4px}.subscription-panel .personal-row button{flex-shrink:0}.subscription-actions{display:flex;gap:10px;align-items:center;flex-wrap:wrap}.subscription-actions button:first-child{background:var(--accent);color:#082028;font-weight:700}.included-card button:disabled{opacity:.6;cursor:default}.personal-source-filter{display:flex;align-items:center;gap:12px;flex-wrap:wrap;margin:15px 0 8px}.personal-source-filter label{font-size:.86rem;font-weight:700}.personal-source-filter select{font:inherit;background:var(--panel);border:1px solid var(--line);border-radius:8px;padding:9px 12px;color:var(--text);max-width:100%}.selected-publication{margin:0 0 18px;font-size:1.5rem}.reading-navigation{display:flex;align-items:center;justify-content:space-between;gap:12px;flex-wrap:wrap;margin:0 0 18px}.reading-navigation .selected-publication{margin:0}.reading-navigation button{cursor:pointer;background:var(--panel);border:1px solid var(--line);color:var(--text);border-radius:10px;padding:9px 13px;font:inherit;font-weight:700}.subscription-read{display:block;border:0!important;background:transparent!important;color:var(--accent)!important;padding:0!important;text-align:left;font:inherit!important;font-weight:700!important;cursor:pointer;text-decoration:underline;text-underline-offset:3px}
.article-thumb{width:128px;height:88px;object-fit:cover;border-radius:10px;flex-shrink:0;background:transparent}.item .article-thumb{margin-left:auto}.feature .article-thumb{width:128px;height:88px;margin:0 0 14px 0}@media(max-width:540px){.article-thumb{width:96px;height:72px}.feature .article-thumb{width:96px;height:72px}}
.item-body>.exclude-site,.item-body>.exclude-article-domain{display:inline-block;vertical-align:top;max-width:100%;width:auto;white-space:normal;overflow-wrap:anywhere;text-align:left;font-family:inherit;font-size:12px;font-weight:500;line-height:1.25;padding:6px 9px;margin:4px 6px 2px 0;border:1px solid var(--line);background:var(--soft);color:var(--text);border-radius:6px;cursor:pointer}.item-body>.exclude-site{white-space:nowrap}.item-body>.exclude-article-domain{overflow-wrap:anywhere}.exclusion-label{color:#e9b978;margin-left:8px}.catalog-card[data-excluded=true],.included-card[data-excluded=true]{border-color:#b5864c}
[hidden]{display:none!important}
.exclusion-confirm-list{margin:10px 0 16px;padding-left:23px;max-height:220px;overflow:auto;color:var(--text)}.exclusion-confirm-list li{margin:5px 0}.exclusion-confirm-action{background:var(--accent)!important;color:#0a2026!important;font-weight:700!important;border-color:var(--accent)!important}
.sync-choice-backdrop{position:fixed;inset:0;z-index:9999;background:rgba(0,0,0,.58);display:flex;align-items:center;justify-content:center;padding:16px}.sync-choice-dialog{width:min(540px,100%);box-sizing:border-box;background:var(--panel);border:1px solid var(--line);color:var(--text);border-radius:14px;padding:20px;box-shadow:0 12px 36px rgba(0,0,0,.3)}.sync-choice-dialog h3{margin:0 0 10px}.sync-choice-dialog p{line-height:1.5}.sync-choice-actions{display:flex;flex-wrap:wrap;gap:9px;margin-top:17px}.sync-choice-actions button{cursor:pointer;font:inherit;color:var(--text);background:var(--soft);border:1px solid var(--line);border-radius:9px;padding:9px 12px}

.admin-settings-form{display:block;margin:12px 0}.admin-settings-fields{min-width:0;border:1px solid var(--line);border-radius:12px;padding:18px 20px 12px}.admin-settings-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:18px 24px;margin:8px 0 18px}.admin-settings-grid label{display:flex;flex-direction:column;align-items:stretch;gap:7px;min-width:0;font-weight:650;line-height:1.4}.admin-settings-grid label small{font-size:.8rem;color:var(--muted);font-weight:400}.admin-settings-grid input{width:100%;max-width:360px;min-width:0;box-sizing:border-box;font:inherit;color:var(--text);background:var(--bg);border:1px solid var(--line);border-radius:9px;padding:10px 12px}.admin-settings-actions{display:flex;gap:12px;align-items:center;flex-wrap:wrap;margin-top:16px}.admin-settings-actions button{cursor:pointer;background:var(--soft);border:1px solid var(--line);color:var(--text);border-radius:9px;padding:11px 16px;font:inherit;font-weight:650}.admin-settings-actions button[type=submit]{background:var(--accent);color:#082028}@media(max-width:700px){.admin-settings-grid{grid-template-columns:1fr}.admin-settings-fields{padding:14px}.admin-settings-grid input{max-width:100%}}
.admin-unsaved-dialog{width:min(520px,calc(100% - 32px));border:1px solid var(--line);border-radius:14px;padding:22px;color:var(--text);background:var(--panel);box-shadow:0 20px 55px #0008}.admin-unsaved-dialog::backdrop{background:#000a}.admin-unsaved-dialog h3{margin:0 0 12px}.admin-unsaved-dialog p{line-height:1.5;color:var(--muted)}
</style><script defer src="/personal.js?v=0.2.58"></script></head><body><div class="topbar"><div class="topbar-inner"><div class="brand">News<span>Go</span></div><div class="header-right"><div class="tagline">NEWS READING WITHOUT PAGE CLUTTER USING SCRAPERGO · V0.2.58 · {{.BuildLabel}}</div><div class="tanner"><img src="/tanner.jpg" width="64" height="64" alt="Tanner, the puppy"><span>Tanner</span></div></div></div></div><main>
<div class="intro"><div><div class="eyebrow">Your news, without the noise</div><h1>Today's briefing.</h1><p id="briefing-summary">Headlines from your subscriptions only.</p></div></div>

<nav class="news-nav" aria-label="NewsGo sections"><button type="button" id="manage-personal" aria-pressed="false">Find Subscriptions</button><button type="button" id="view-subscriptions" aria-pressed="false">My Subscriptions <span id="personal-count">(0)</span></button><button type="button" id="manage-preferences">Reading Preferences</button><button type="button" id="manage-exclusions" aria-pressed="false">Excluded Domains <span id="excluded-count">(0)</span></button><button type="button" id="manage-sync">Server Sync</button><button type="button" id="manage-admin" hidden>Admin Settings</button></nav>
<section id="preferences-panel" class="personal-panel personal-hidden" aria-label="Reading preferences"><div class="panel-heading"><h2>Reading Preferences</h2><button type="button" id="close-preferences" class="panel-close">Close</button></div><p>Choose how old an article can be before it is hidden from your briefing. Undated articles remain visible. This preference belongs to your reader profile, not the administrator.</p><div class="reading-age-row"><label for="maximum-article-age">Maximum article age</label><select id="maximum-article-age"><option value="0">All articles (no age limit)</option><option value="1">1 day</option><option value="3">3 days</option><option value="7">7 days</option><option value="14">14 days</option><option value="30">30 days</option><option value="60">60 days</option></select></div><p class="personal-status">Filtered older stories appear in the Filtered / unprocessed articles card.</p></section>
<section id="admin-panel" class="personal-panel personal-hidden" aria-label="Server administrator settings"><div class="panel-heading"><h2>Admin Settings</h2><button type="button" id="close-admin" class="panel-close">Close</button></div><p>Server-wide operational controls. Article processing settings apply to new refreshes. ScraperGo manages thumbnail retrieval separately.</p><form id="admin-form" class="admin-settings-form"><fieldset class="admin-settings-fields"><legend>Article processing (active)</legend><div class="admin-settings-grid"><label>Maximum articles per feed <small>(1–250)</small><input type="number" id="admin-articles" min="1" max="250" required></label><label>Classification timeout <small>seconds (2–30)</small><input type="number" id="admin-classification-timeout" min="2" max="30" required></label><label>Feed directory search timeout <small>seconds (2–30)</small><input type="number" id="admin-feed-search-timeout" min="2" max="30" required></label><label>Parallel classifications <small>(1–16)</small><input type="number" id="admin-classification-parallel" min="1" max="16" required></label><label>Successful feed cache <small>minutes (1–60)</small><input type="number" id="admin-feed-cache" min="1" max="60" required></label></div><p class="personal-status">Feed requests are grouped into batches of eight; unprocessed entries and rejected classifications appear below headlines.</p></fieldset><div class="admin-settings-actions"><button type="submit">Save Server Settings</button><button id="admin-reload" type="button">Reload</button><button id="admin-exit" type="button">Exit Admin Mode</button></div></form><p class="personal-status">Disconnects this browser from Server Sync and removes its saved token. Your server token remains valid; reconnect through Server Sync to administer again.</p><p id="admin-status" class="personal-status" role="status" aria-live="polite"></p><dialog id="admin-unsaved-dialog" class="admin-unsaved-dialog" aria-labelledby="admin-unsaved-title"><form method="dialog"><h3 id="admin-unsaved-title">Unsaved admin settings</h3><p>You have changes that have not been saved to the server. What would you like to do?</p><div class="admin-settings-actions"><button value="save" type="submit">Save changes</button><button value="discard" type="submit">Discard changes</button><button value="cancel" type="submit" autofocus>Keep editing</button></div></form></dialog><dialog id="admin-exit-dialog" class="admin-unsaved-dialog" aria-labelledby="admin-exit-title"><form method="dialog"><h3 id="admin-exit-title">Exit Admin Mode?</h3><p>This will discard your saved admin token from this browser and disconnect Server Sync. You will need to enter your token again to access Admin Settings.</p><div class="admin-settings-actions"><button value="cancel" type="submit" autofocus>Cancel</button><button value="confirm" type="submit">Exit and remove token</button></div></form></dialog></section>
<section id="sync-panel" class="personal-panel personal-hidden" aria-label="Server synchronization"><div class="panel-heading"><h2>Server Sync</h2><button type="button" class="panel-close" id="close-sync">Close</button></div><p>Use one approved token for all your devices. Subscriptions and excluded sites save automatically and refresh when newer server preferences become available.</p><p><strong id="sync-connection">Not connected</strong></p><div class="personal-form"><input type="password" id="sync-token" autocomplete="off" placeholder="Enter the 64-character approved token distributed by Richard." aria-label="Approved sync token"><button id="sync-connect" type="button">Connect</button></div><p class="personal-status" id="sync-status" role="status" aria-live="polite">Not connected. Browser preferences remain local.</p><details id="sync-advanced"><summary>Advanced</summary><div class="personal-form"><button id="sync-load" type="button">Load from Server</button><button id="sync-save" type="button">Save Now</button><button id="sync-disconnect" type="button">Disconnect</button></div><p>Manual controls for troubleshooting. Loading replaces both local lists; a confirmation is required if local changes are pending.</p></details></section>
<section id="personal-panel" class="personal-panel personal-hidden" aria-label="Manage personal feeds"><div class="panel-heading"><h2>Find Subscriptions</h2><button type="button" class="panel-close" id="close-discovery">Close</button></div><p>Search the web for news feeds by topic or publisher, then choose Subscribe. Your selections stay in this browser unless Server Sync is enabled.</p><div class="browse-included"><h3>Browse available sources</h3><p>Browse this feed catalog or search for more sources. Nothing appears in your briefing until you subscribe.</p><div class="included-list">{{range .Catalog}}<div class="included-card"><strong>{{.Name}}</strong><button type="button" class="read-included" data-included-source="{{.Name}}" data-feed-url="{{.URL}}" data-category="{{.Category}}">Subscribe</button></div>{{end}}</div></div><div class="discovery-divider">Or search for more subscriptions</div><form id="catalog-form" class="catalog-search-form"><input class="catalog-search" id="catalog-search" type="search" placeholder="Try sailing, gardening, college football…" aria-label="Search news sources"><button type="submit" id="catalog-submit">Search</button></form><div class="catalog-list" id="catalog-list" aria-live="polite"></div><p class="personal-status" id="catalog-message"></p><details><summary>Add a website not listed here</summary>
<form id="discover-form" class="personal-form"><input id="discover-url" type="url" maxlength="2048" placeholder="Paste a news website address, e.g. https://www.foxsports.com/college-football" aria-label="News website URL" required><button id="discover-button" type="button">Find feeds</button></form><div id="discover-results" aria-live="polite"></div><details><summary>Advanced: add an RSS feed manually</summary><form id="personal-form" class="personal-form"><input id="personal-name" maxlength="70" placeholder="Source name" required><input id="personal-url" type="url" maxlength="2048" placeholder="https://example.com/rss" required><select id="personal-category"><option>World</option><option>Business</option><option>Technology</option><option>Science</option><option>Sports</option><option>Local</option><option>Other</option></select><button>Add RSS feed</button></form></details>
<div><small>Suggested additions:</small><div id="personal-suggestions"></div></div></details>
<button id="personal-export" type="button">Export subscriptions</button> <button id="personal-import" type="button">Import subscriptions</button><input id="personal-file" type="file" accept="application/json,.json" hidden><p class="personal-status" id="personal-message" role="status">Feed discovery is loading…</p></section>
<section id="subscriptions-panel" class="personal-panel subscription-panel personal-hidden" aria-label="Manage subscriptions"><div class="panel-heading"><h2>My subscriptions <span id="subscription-count"></span></h2><button type="button" class="panel-close" id="close-subscriptions">Close</button></div><p>Select Read Articles to view headlines, or Unsubscribe to remove a subscription.</p><div id="personal-sources"></div><p class="personal-status" id="subscriptions-empty"></p></section>
<section id="exclusions-panel" class="personal-panel personal-hidden" aria-label="Manage excluded domains"><div class="panel-heading"><h2>Excluded Domains</h2><button type="button" class="panel-close" id="close-exclusions">Close</button></div><p>Exclude a domain and its subdomains from subscriptions and article links. If any current subscriptions are affected, NewsGo will ask before removing them. Allowing a domain does not restore removed subscriptions. Existing saved exclusions are retained. Preferences remain in this browser unless synchronized.</p><div class="personal-form"><input id="exclude-domain" placeholder="example.com or https://example.com" aria-label="Domain to exclude"><button type="button" id="exclude-add">Exclude domain</button></div><div id="excluded-list"></div><div class="personal-form" hidden><input id="article-exclude-domain" aria-label="Legacy article domain"><button type="button" id="article-exclude-add">Legacy article exclusion</button></div></section><div id="news-content">
<div id="subscriptions-empty-notice" class="notice" hidden>No subscriptions yet. Use Find Subscriptions to add a publication.</div>
<div class="reading-navigation"><h2 id="reading-title" class="selected-publication" hidden></h2></div><select id="read-source" hidden><option value="">All sources</option><optgroup id="followed-options" label="Followed sources"></optgroup></select>
<section id="personal-feed" class="personal-hidden" aria-label="Personal headlines"><div class="section-label" id="personal-section-label">My subscriptions</div><p class="personal-status" id="personal-loading" role="status" aria-live="polite" aria-atomic="true"></p><details id="filtered-articles-card" class="personal-panel" hidden><summary><strong>Filtered / unprocessed articles</strong> <span id="filtered-articles-count"></span></summary><p class="personal-status">Entries excluded by classification, publisher exclusion rules, or the per-feed processing cap. Not-processed entries were never classified.</p><div id="filtered-articles-list"></div></details><div id="personal-lead-section" hidden><div class="section-label">The lead stories</div><section class="feature-grid" id="personal-lead-stories" aria-label="Lead stories"></section></div><div id="personal-latest-section" hidden><div class="section-label">Latest headlines</div><section class="latest" id="personal-articles" aria-label="Latest headlines"></section></div></section>
</div><footer>Choose a headline to read in ScraperGo. NewsGo reuses a single reader tab; original links open the publisher. Feed availability and article readability vary.</footer></main></body></html>`))
