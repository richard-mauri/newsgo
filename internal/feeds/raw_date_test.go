package feeds

import (
 "testing"
 "time"
)

func TestUnrecognizedRSSDatePreserved(t *testing.T) {
 raw:=`<rss version="2.0"><channel><item><title>Unknown date format</title><link>https://example.org/item</link><pubDate>October 21st at breakfast</pubDate></item><item><title>No date</title><link>https://example.org/missing</link></item></channel></rss>`
 _,entries,err:=Parse([]byte(raw))
 if err!=nil {t.Fatal(err)}
 if len(entries)!=2 {t.Fatalf("items=%d",len(entries))}
 if !entries[0].date.IsZero() || entries[0].rawDate!="October 21st at breakfast" {t.Fatalf("unrecognized date not preserved: %+v",entries[0])}
 if !entries[1].date.Equal(time.Time{}) || entries[1].rawDate!="" {t.Fatalf("missing date should remain missing: %+v",entries[1])}
}
