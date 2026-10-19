package dashboard

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/richard-mauri/newsgo/internal/feeds"
)

//go:embed tanner.jpg
var tannerPhoto embed.FS

type State struct {
	mu        sync.RWMutex
	articles  []feeds.Article
	errors    map[string]string
	refreshed time.Time
	sources   []feeds.Source
	reader    string
	client    *http.Client
}

func New(sources []feeds.Source, reader string, client *http.Client) (*State, error) {
	if err := feeds.ValidateSources(sources); err != nil {
		return nil, err
	}
	u, err := url.Parse(reader)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil {
		return nil, fmt.Errorf("invalid ScraperGo address")
	}
	return &State{sources: append([]feeds.Source(nil), sources...), reader: strings.TrimRight(reader, "/"), client: client, errors: map[string]string{}}, nil
}
func (s *State) Refresh(ctx context.Context) {
	type result struct {
		src      feeds.Source
		articles []feeds.Article
		err      error
	}
	ch := make(chan result, len(s.sources))
	sem := make(chan struct{}, 4)
	var wg sync.WaitGroup
	for _, src := range s.sources {
		wg.Add(1)
		go func(src feeds.Source) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				ch <- result{src: src, err: ctx.Err()}
				return
			}
			defer func() { <-sem }()
			a, e := feeds.FetchArticles(ctx, s.client, src)
			ch <- result{src, a, e}
		}(src)
	}
	wg.Wait()
	close(ch)
	s.mu.RLock()
	old := append([]feeds.Article(nil), s.articles...)
	s.mu.RUnlock()
	errors := make(map[string]string)
	articles := make([]feeds.Article, 0)
	okSources := map[string]bool{}
	bySource := map[string]result{}
	for r := range ch {
		bySource[r.src.Name] = r
	}
	for _, src := range s.sources {
		r := bySource[src.Name]
		if r.err != nil {
			errors[src.Name] = r.err.Error()
		} else {
			okSources[src.Name] = true
			articles = append(articles, r.articles...)
		}
	}
	// Preserve previously retrieved headlines for sources that temporarily fail.
	for _, a := range old {
		if !okSources[a.Source] {
			articles = append(articles, a)
		}
	}
	seen := map[string]bool{}
	dedup := make([]feeds.Article, 0, len(articles))
	for _, a := range articles {
		if seen[a.URL] {
			continue
		}
		seen[a.URL] = true
		dedup = append(dedup, a)
	}
	sort.SliceStable(dedup, func(i, j int) bool { return dedup[i].Published.After(dedup[j].Published) })
	if len(dedup) > 600 {
		dedup = dedup[:600]
	}
	s.mu.Lock()
	s.articles = dedup
	s.errors = errors
	s.refreshed = time.Now()
	s.mu.Unlock()
}
func (s *State) Run(ctx context.Context, interval time.Duration) {
	s.Refresh(ctx)
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			s.Refresh(ctx)
		}
	}
}

type viewArticle struct{ Title, Source, Category, Date, Open, Original string }
type view struct {
	Articles                         []viewArticle
	Sources                          []string
	Categories                       []string
	SelectedSource, SelectedCategory string
	Updated                          string
	Total                            int
	Errors                           []string
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
	if r.URL.Path == "/api/discover" {
		s.serveDiscover(w, r)
		return
	}
	if r.URL.Path == "/api/personal" {
		s.servePersonal(w, r)
		return
	}
	if r.URL.Path == "/personal.js" {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")
		fmt.Fprint(w, personalJS)
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
	s.mu.RLock()
	articles := append([]feeds.Article(nil), s.articles...)
	updated := s.refreshed
	errs := make(map[string]string, len(s.errors))
	for k, v := range s.errors {
		errs[k] = v
	}
	s.mu.RUnlock()
	source := r.URL.Query().Get("source")
	category := r.URL.Query().Get("category")
	v := view{SelectedSource: source, SelectedCategory: category, Total: len(articles)}
	if !updated.IsZero() {
		v.Updated = updated.Local().Format("Jan 2, 2006 3:04 PM")
	}
	sm := map[string]bool{}
	cm := map[string]bool{}
	for _, src := range s.sources {
		sm[src.Name] = true
		cm[src.Category] = true
	}
	for n := range sm {
		v.Sources = append(v.Sources, n)
	}
	sort.Strings(v.Sources)
	for c := range cm {
		v.Categories = append(v.Categories, c)
	}
	sort.Strings(v.Categories)
	for k, e := range errs {
		v.Errors = append(v.Errors, k+": "+e)
	}
	sort.Strings(v.Errors)
	for _, a := range articles {
		if source != "" && source != a.Source || category != "" && category != a.Category {
			continue
		}
		date := "Date unavailable"
		if !a.Published.IsZero() {
			date = a.Published.Local().Format("Jan 2, 3:04 PM")
		}
		v.Articles = append(v.Articles, viewArticle{a.Title, a.Source, a.Category, date, s.reader + "/?url=" + url.QueryEscape(a.URL), a.URL})
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; base-uri 'none'; form-action 'self'; script-src 'self'; connect-src 'self'; img-src 'self'")
	if err := page.Execute(w, v); err != nil {
		return
	}
}

var page = template.Must(template.New("index").Funcs(template.FuncMap{"add": func(a, b int) int { return a + b }}).Parse(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><meta name="color-scheme" content="light dark"><title>NewsGo — Your daily briefing</title>
<style>
:root{color-scheme:light dark;--bg:light-dark(#f5f5f0,#0e141a);--panel:light-dark(#fff,#19232b);--text:light-dark(#182631,#eef2f4);--muted:light-dark(#60707b,#a5b3bb);--line:light-dark(#dce2e1,#34434a);--accent:light-dark(#126877,#69c6d3);--soft:light-dark(#e6f3f3,#203c43);font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif}
*{box-sizing:border-box}body{margin:0;background:var(--bg);color:var(--text)}a{color:inherit}main{max-width:1200px;margin:0 auto;padding:0 24px 72px}
.topbar{border-bottom:1px solid var(--line);background:var(--panel)}.topbar-inner{max-width:1200px;margin:auto;padding:19px 24px;display:flex;justify-content:space-between;align-items:center;gap:14px}.brand{font-weight:850;font-size:1.75rem;letter-spacing:-.075rem}.brand span{color:var(--accent)}.tagline{font-size:.82rem;color:var(--muted);letter-spacing:.03em}
.header-right{display:flex;align-items:center;gap:22px}.tanner{display:flex;align-items:center;gap:8px;color:var(--muted);font-size:.82rem;font-weight:700}.tanner img{width:64px;height:64px;object-fit:cover;object-position:center 35%;border-radius:12px;border:1px solid var(--line)}@media(max-width:540px){.header-right{gap:6px}.tanner img{width:49px;height:49px}.tanner span{display:none}}
.eyebrow{font-size:.73rem;font-weight:780;letter-spacing:.16em;text-transform:uppercase;color:var(--accent)}.intro{padding:38px 0 25px;display:flex;justify-content:space-between;align-items:end;gap:18px}.intro h1{font-family:Georgia,serif;font-size:clamp(2.6rem,5vw,4.6rem);line-height:1;margin:9px 0 8px;font-weight:700;letter-spacing:-.035em}.intro p{margin:0;color:var(--muted);font-size:.93rem}.time{color:var(--muted);font-size:.85rem;text-align:right}
.filters{padding:17px 19px;background:var(--panel);border:1px solid var(--line);border-radius:15px;display:flex;align-items:end;gap:15px;flex-wrap:wrap;margin-bottom:31px}.filters label{display:flex;flex-direction:column;gap:6px;font-size:.74rem;font-weight:720;color:var(--muted);text-transform:uppercase;letter-spacing:.065em}.filters select{min-width:175px;max-width:100%;font:inherit;font-size:.92rem;color:var(--text);background:var(--bg);border:1px solid var(--line);border-radius:9px;padding:10px 12px}.filter-btn{background:var(--accent);color:light-dark(#fff,#082028);padding:11px 20px;border-radius:9px;border:0;font-size:.89rem;font-weight:700;cursor:pointer}.clear{font-size:.88rem;color:var(--muted);padding-bottom:10px;text-decoration:underline;text-underline-offset:3px}
.section-label{display:flex;align-items:center;gap:14px;font-weight:770;font-size:.77rem;letter-spacing:.13em;text-transform:uppercase;margin:25px 0 17px}.section-label:after{content:"";height:1px;flex:1;background:var(--line)}
.feature-grid{display:grid;grid-template-columns:1.35fr 1fr 1fr;gap:16px}.feature{min-width:0;min-height:282px;display:flex;flex-direction:column;justify-content:space-between;gap:28px;background:var(--panel);border:1px solid var(--line);border-radius:17px;padding:23px;box-shadow:0 10px 24px #00000008}.feature:first-child{background:linear-gradient(130deg,light-dark(#d7eeeb,#16434a),var(--panel) 85%);border-color:light-dark(#c5dddd,#31585f)}.feature .num{font-family:Georgia,serif;font-size:3.25rem;color:var(--accent);opacity:.46;line-height:1}.story-meta{color:var(--muted);font-size:.77rem;line-height:1.6}.source{font-weight:800;text-transform:uppercase;letter-spacing:.075em;color:var(--accent);font-size:.72rem}.feature .story-title{font-family:Georgia,serif;font-size:clamp(1.15rem,1.65vw,1.57rem);line-height:1.22;font-weight:700;letter-spacing:-.015em;text-decoration:none;display:block;margin:10px 0 12px}.feature:first-child .story-title{font-size:clamp(1.5rem,2.3vw,2.2rem)}.story-title:hover{text-decoration:underline;text-underline-offset:3px}.original{color:var(--muted);font-size:.79rem;text-decoration:underline;text-underline-offset:3px}.feature .original{display:inline-block;margin-top:7px}
.latest{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:0 28px}.item{display:flex;gap:17px;padding:21px 0;border-bottom:1px solid var(--line);min-width:0}.item-num{font-family:Georgia,serif;font-size:1.6rem;color:var(--muted);opacity:.55;min-width:30px}.item-body{min-width:0}.item .story-title{display:block;font-family:Georgia,serif;font-size:1.23rem;font-weight:700;line-height:1.32;margin:7px 0 7px;text-decoration:none}.item .story-title:hover{text-decoration:underline}.item .original{margin-left:8px}.notice{border:1px solid var(--line);border-radius:9px;padding:12px;margin:12px 0;color:var(--muted)}footer{margin-top:35px;padding-top:17px;border-top:1px solid var(--line);font-size:.8rem;color:var(--muted)}
@media(max-width:820px){.feature-grid{grid-template-columns:1fr 1fr}.feature:first-child{grid-column:1/-1;min-height:260px}.feature{min-height:235px}.latest{grid-template-columns:1fr}.intro{padding-top:28px}}
@media(max-width:540px){main{padding:0 16px 50px}.topbar-inner{padding:15px 16px}.tagline{display:none}.intro{display:block}.time{text-align:left;margin-top:12px}.filters{gap:12px}.filters label{flex:1;min-width:130px}.filters select{width:100%;min-width:0}.feature-grid{grid-template-columns:1fr}.feature:first-child{grid-column:auto}.feature{min-height:200px}.item{padding:18px 0}}

.personal-tools{display:flex;align-items:center;justify-content:space-between;gap:12px;margin-bottom:20px;flex-wrap:wrap}.personal-tools button,.personal-form button,.personal-row button{cursor:pointer;background:var(--soft);border:1px solid var(--line);color:var(--text);border-radius:9px;padding:9px 13px;font:inherit}.personal-form{display:flex;gap:10px;flex-wrap:wrap;margin:12px 0}.personal-form input,.personal-form select{flex:1;min-width:170px;border:1px solid var(--line);background:var(--bg);color:var(--text);border-radius:9px;padding:10px;font:inherit}.personal-panel{padding:18px;border:1px solid var(--line);border-radius:15px;background:var(--panel);margin:14px 0 25px}.personal-row{display:flex;align-items:center;justify-content:space-between;border-top:1px solid var(--line);padding:10px 0;gap:10px}.personal-row small{color:var(--muted);overflow-wrap:anywhere}.personal-status{color:var(--muted);font-size:.86rem}.personal-list{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:0 24px}.personal-suggestion{margin:7px 8px 7px 0}.personal-hidden{display:none}.personal-panel h2{margin:0 0 9px}.personal-panel p{color:var(--muted)}@media(max-width:720px){.personal-list{grid-template-columns:1fr}}
</style><script defer src="/personal.js?v=0.2.7"></script></head><body><div class="topbar"><div class="topbar-inner"><div class="brand">News<span>Go</span></div><div class="header-right"><div class="tagline">NEWS DISCOVERY FOR SCRAPERGO · V0.2.7</div><div class="tanner"><img src="/tanner.jpg" width="64" height="64" alt="Tanner, the puppy"><span>Tanner</span></div></div></div></div><main>
<div class="intro"><div><div class="eyebrow">Your news, without the noise</div><h1>Today's briefing.</h1><p>{{.Total}} cached headlines from trusted feeds.</p></div><div class="time">{{if .Updated}}Updated {{.Updated}}{{else}}Fetching headlines…{{end}}</div></div>
<form class="filters" method="get" action="/"><label>Publisher<select name="source"><option value="">All sources</option>{{range .Sources}}<option value="{{.}}" {{if eq . $.SelectedSource}}selected{{end}}>{{.}}</option>{{end}}</select></label><label>Section<select name="category"><option value="">All categories</option>{{range .Categories}}<option value="{{.}}" {{if eq . $.SelectedCategory}}selected{{end}}>{{.}}</option>{{end}}</select></label><button class="filter-btn" type="submit">Apply filters</button><a class="clear" href="/">Reset</a></form>
<div class="personal-tools"><strong>Personal news feeds <span id="personal-count"></span></strong><button type="button" id="manage-personal">Manage feeds +</button></div>
<section id="personal-panel" class="personal-panel personal-hidden" aria-label="Manage personal feeds"><h2>My news sources</h2><p>Subscriptions are saved in this browser only. They aren't shared with other visitors.</p>
<form id="discover-form" class="personal-form"><input id="discover-url" type="url" maxlength="2048" placeholder="Paste a news website address, e.g. https://www.foxsports.com/college-football" aria-label="News website URL" required><button id="discover-button" type="button">Find feeds</button></form><div id="discover-results" aria-live="polite"></div><details><summary>Advanced: add an RSS feed manually</summary><form id="personal-form" class="personal-form"><input id="personal-name" maxlength="70" placeholder="Source name" required><input id="personal-url" type="url" maxlength="2048" placeholder="https://example.com/rss" required><select id="personal-category"><option>World</option><option>Business</option><option>Technology</option><option>Science</option><option>Local</option><option>Other</option></select><button>Add RSS feed</button></form></details>
<div><small>Suggested additions:</small><div id="personal-suggestions"></div></div><div id="personal-sources"></div>
<button id="personal-export" type="button">Export subscriptions</button> <button id="personal-import" type="button">Import subscriptions</button><input id="personal-file" type="file" accept="application/json,.json" hidden><p class="personal-status" id="personal-message" role="status">Feed discovery is loading…</p></section>
<section id="personal-feed" class="personal-hidden" aria-label="Personal headlines"><div class="section-label">My sources</div><p class="personal-status" id="personal-loading"></p><div class="personal-list" id="personal-articles"></div></section>
{{range .Errors}}<div class="notice">Feed unavailable: {{.}}</div>{{end}}
{{if not .Articles}}<div class="notice">No headlines match the selected filters, or feeds have not loaded yet.</div>{{else}}
<div class="section-label">The lead stories</div><section class="feature-grid" aria-label="Lead stories">{{range $i,$a := .Articles}}{{if lt $i 3}}<article class="feature"><div class="num">0{{add $i 1}}</div><div><div class="source">{{$a.Source}} · {{$a.Category}}</div><a class="story-title" href="{{$a.Open}}" target="newsgo-scrapergo">{{$a.Title}}</a><div class="story-meta">{{$a.Date}}</div><a class="original" href="{{$a.Original}}" target="_blank" rel="noopener noreferrer">Original publisher ↗</a></div></article>{{end}}{{end}}</section>
<div class="section-label">Latest headlines</div><section class="latest" aria-label="Latest headlines">{{range $i,$a := .Articles}}{{if ge $i 3}}<article class="item"><div class="item-num">{{add $i 1}}</div><div class="item-body"><div class="source">{{$a.Source}} · {{$a.Category}}</div><a class="story-title" href="{{$a.Open}}" target="newsgo-scrapergo">{{$a.Title}}</a><div class="story-meta">{{$a.Date}} <a class="original" href="{{$a.Original}}" target="_blank" rel="noopener noreferrer">Original ↗</a></div></div></article>{{end}}{{end}}</section>{{end}}
<footer>Choose a headline to read in ScraperGo. NewsGo reuses a single reader tab; original links open the publisher. Feed availability and article readability vary.</footer></main></body></html>`))
