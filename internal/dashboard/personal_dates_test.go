package dashboard

import (
 "strings"
 "testing"
 "time"

 "github.com/richard-mauri/newsgo/internal/feeds"
)

func TestDateDiagnosticPreservesArticleURLAndRawDate(t *testing.T) {
 raw := "Fri, Nov 06 2026 03:05:45 PM"
 instant := time.Date(2026,11,6,15,5,45,0,time.UTC)
 a := feeds.Article{Title:"Example news",Source:"NBC Bay Area",URL:"https://www.nbcbayarea.com/news/example?src=rss",RawDate:raw,Published:instant}
 got := dateDiagnostic(a)
 if got.URL!=a.URL || got.RawDate!=raw || got.Title!=a.Title || got.Source!=a.Source {t.Fatalf("article provenance not preserved: %+v",got)}
 if got.ParseStatus!="recognized" || got.ParsedTimestamp!=instant.Format(time.RFC3339) {t.Fatalf("wrong date parse diagnostic: %+v",got)}
 if !strings.Contains(got.Timezone,"Not specified") {t.Fatalf("expected timezone ambiguity annotation: %+v",got)}
}
func TestDateDiagnosticUnknownDatePreserved(t *testing.T) {
 raw := "unknown-format-date"
 a := feeds.Article{Title:"Other",Source:"Feed",URL:"https://example.com/a",RawDate:raw}
 got := dateDiagnostic(a)
 if got.ParseStatus!="unrecognized" || got.DisplayedDate!=raw || got.RawDate!=raw || got.ParsedTimestamp!="" || got.URL!=a.URL {t.Fatalf("bad fallback provenance: %+v",got)}
}
