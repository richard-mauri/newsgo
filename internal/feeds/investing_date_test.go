package feeds

import (
 "strings"
 "testing"
 "time"
)

func TestInvestingRSSPublicationDate(t *testing.T) {
 raw := `<rss version="2.0"><channel><item>
 <enclosure url="https://content-media.investing.com/trkd-images/LYNXMPEM9K0IR_L.jpg" length="72616" type="image/jpeg" />
 <title>Government websites go dark after AI-discovered flaw sparks security scramble</title>
 <pubDate>2026-10-21 20:48:33</pubDate>
 <author>Reuters</author>
 <link>https://www.investing.com/news/stock-market-news/cyber-agency-warns-websites-may-go-dark-as-php-prepares-critical-fix-4961963</link>
 </item></channel></rss>`
 format,items,err:=Parse([]byte(raw))
 if err!=nil {t.Fatal(err)}
 if !strings.Contains(strings.ToLower(format),"rss"){t.Fatalf("unexpected format %q",format)}
 if len(items)!=1 {t.Fatalf("got %d items",len(items))}
 expected:=time.Date(2026,10,21,20,48,33,0,time.UTC)
 if !items[0].date.Equal(expected){t.Fatalf("published=%v, want %v",items[0].date,expected)}
 if items[0].date.IsZero(){t.Fatal("Investing.com date was discarded")}
}
