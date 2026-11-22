package dashboard

import (
	"strings"
	"testing"
)

func TestDiagnosticsOwnsFilteredAuditAndExcludedDomains(t *testing.T) {
	html := page.Tree.Root.String()
	nav := strings.Index(html, `<nav class="news-nav"`)
	diagnostics := strings.Index(html, `id="diagnostics-panel"`)
	diagnosticsEnd := strings.Index(html, `id="personal-panel"`)
	news := strings.Index(html, `id="news-content"`)
	if nav < 0 || diagnostics < 0 || diagnosticsEnd <= diagnostics || news <= diagnosticsEnd {
		t.Fatal("unexpected dashboard section ordering")
	}
	if strings.Contains(html[nav:diagnostics], `id="exclusions-panel"`) {
		t.Fatal("Excluded Domains disclosure must not be on the main page")
	}
	panel := html[diagnostics:diagnosticsEnd]
	for _, id := range []string{
		`<details id="latest-read-panel" class="personal-panel"><summary><strong>Latest read</strong></summary>`,
		`id="diagnostics-summary"`,
		`id="diagnostics-timing"`,
		`id="diagnostics-note"`,
		`id="thumbnail-performance-panel"`,
		`id="thumbnail-performance"`,
		`<details id="exclusions-panel"`,
		`<summary><strong>Excluded Domains</strong> <span id="excluded-count">(0)</span></summary>`,
		`id="filtered-articles-card"`,
		`id="filtered-articles-list"`,
		`id="filtered-articles-count"`,
		`<details id="raw-feed-inspector" class="personal-panel" aria-label="Raw feed responses"><summary><strong>Raw feed responses</strong> <span id="raw-feed-count"></span></summary>`,
		`id="raw-feed-list"`,
		`id="raw-feed-status"`,
	} {
		if !strings.Contains(panel, id) {
			t.Errorf("Diagnostics missing %s", id)
		}
	}
	latestStart := strings.Index(panel, `<details id="latest-read-panel"`)
	exclusionsStart := strings.Index(panel, `<details id="exclusions-panel"`)
	if latestStart < 0 || exclusionsStart <= latestStart || !strings.Contains(panel[latestStart:exclusionsStart], `id="diagnostics-note"`) || !strings.Contains(panel[latestStart:exclusionsStart], `</details>`) {
		t.Fatal("latest-read metrics and timing must be inside their own disclosure before other panels")
	}
	if strings.Contains(panel[:latestStart], `id="diagnostics-summary"`) {
		t.Fatal("latest-read metrics rendered outside the disclosure")
	}
	if strings.Contains(panel, `id="manage-exclusions"`) || strings.Contains(panel, `id="close-exclusions"`) {
		t.Fatal("obsolete Excluded Domains navigation controls remain")
	}
	if strings.Contains(panel, `<section id="raw-feed-inspector"`) || strings.Contains(panel, `<h3>Raw feed responses`) {
		t.Fatal("raw feed inspector still uses non-disclosure heading")
	}
	if strings.Contains(html[news:], `id="filtered-articles-card"`) {
		t.Fatal("filtered audit still appears in main article content")
	}
}

func TestDiagnosticsExclusionNavigationAndFilteringRetained(t *testing.T) {
	for _, marker := range []string{
		"el('diagnostics-panel').classList.toggle('personal-hidden',view!=='diagnostics')",
		"el('manage-diagnostics').addEventListener('click',()=>{showView('diagnostics');updateDiagnostics()})",
		"function renderFilteredSubscriptionOffer()",
		"function renderFilteredArticles()",
		"renderFilteredArticles();",
		"function setUnifiedDomain(domain,allow)",
		"const details=document.createElement('details');details.className='raw-feed-entry'",
		"const download=button('Download Raw Feed',()=>{",
		"const copy=button('Copy Raw Feed',async()=>{",
		"list.replaceChildren();count.textContent='('+latestRawFeeds.length+' feeds)'",
		"el('excluded-count').textContent='('+domains.length+')",
	} {
		if !strings.Contains(personalJS, marker) {
			t.Errorf("missing diagnostics/exclusions behavior: %s", marker)
		}
	}
}

func TestThumbnailAttributionShowsMissingMarkerAsUnattributed(t *testing.T) {
 for _,needle:=range []string{
  "response.headers.get('X-NewsGo-Thumbnail-Proxy')==='v1'",
  "'proxy_'+reason",
  "'unattributed_http_'+response.status",
  "Unattributed HTTP responses (no NewsGo marker)",
  "browser_transport_or_auth_failure",
  "browser_decode_failure",
 } {
  if !strings.Contains(personalJS,needle) {t.Errorf("missing thumbnail attribution guard %q",needle)}
 }
}

func TestDiagnosticsTextExportAvailableWithoutExposingSensitiveData(t *testing.T) {
 html:=page.Tree.Root.String()
 if !strings.Contains(html,"id=\"download-diagnostics\"") {t.Fatal("missing diagnostic download button")}
 start:=strings.Index(personalJS,"function diagnosticsTextReport(){")
 end:=strings.Index(personalJS,"function downloadDiagnostics(){")
 if start<0 || end<=start {t.Fatal("missing standalone report construction")}
 report:=personalJS[start:end]
 for _,want:=range []string{
  "snapshot('diagnostics-summary')",
  "snapshot('diagnostics-timing')",
  "snapshot('thumbnail-performance')",
  "'LATEST READ'","'THUMBNAIL PERFORMANCE'","'INTERPRETATION'",
  "return lines.join('\\n')+'\\n'",
 }{
  if !strings.Contains(report,want){t.Errorf("report missing %q",want)}
 }
 // Only explicitly allowed aggregate metric values may enter the text file.
 for _,forbidden:=range []string{
  "document.body.innerText","document.documentElement","outerHTML",
  "latestRawFeeds","allPersonalArticles","suiteHandoff",
  "sessionStorage","localStorage","document.cookie",
 }{
  if strings.Contains(report,forbidden){t.Errorf("report includes sensitive source %q",forbidden)}
 }
 downloadStart:=strings.Index(personalJS,"function downloadDiagnostics(){")
 downloadEnd:=strings.Index(personalJS[downloadStart:],"el('download-diagnostics').addEventListener")+downloadStart
 if downloadStart<0||downloadEnd<=downloadStart {t.Fatal("missing direct-download wiring")}
 download:=personalJS[downloadStart:downloadEnd]
 for _,want:=range []string{
  "diagnosticsTextReport()","text/plain;charset=utf-8",
  "URL.createObjectURL(blob)","link.download='newsgo-diagnostics-'",
  "URL.revokeObjectURL(objectURL)",
 }{
  if !strings.Contains(download,want){t.Errorf("missing local TXT download behavior %q",want)}
 }
 if strings.Contains(download,"fetch(") || strings.Contains(download,"adminRequest("){t.Fatal("download must not make a network request")}
}

func TestProgressiveThumbnailBatchBrowserContract(t *testing.T){
 for _,needle:=range []string{
  "async function loadThumbnailBatches(items,generation)",
  "items.slice(offset,offset+20)",
  "authorizedNewsgoFetch('/api/thumbnail/batch'",
  "response.headers.get('X-Thumbnail-Batch')!=='v1'",
  "response.body.getReader()",
  "new TextDecoder()",
  "deliver(JSON.parse(line))",
  "thumbnailBatchGeneration++",
  "controller.abort()",
  "URL.revokeObjectURL(src)",
  "thumbnailBatchItems.push({img,article:a.url,imageURL:a.image,started:imageStart})",
  "void loadThumbnailBatches(thumbnailBatchItems,batchGeneration)",
  "Browser thumbnail batch requests",
  "if(!rssMetadataOnly)thumbnailBatchItems.push",
  "record.status!==200",
  "outcome:'browser_batch_transport_failure'",
 } {
  if !strings.Contains(personalJS,needle) {t.Errorf("missing progressive batching browser contract: %q",needle)}
 }
 if strings.Contains(personalJS,"if(!rssMetadataOnly)void loadProtectedThumbnail(img") {
  t.Fatal("normal thumbnail path still requests all images individually")
 }
}
