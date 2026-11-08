package dashboard

import (
 "bytes"
 "encoding/json"
 "testing"
)

func TestRawFeedDiagnosticPreservesOriginalBytes(t *testing.T) {
 raw := []byte("<?xml version=\"1.0\"?>\r\n<rss><channel><item><pubDate>Fri, Nov 06 2026 03:05:45 PM</pubDate></item></channel></rss>\r\n")
 diagnostic := rawFeedDiagnostic{URL:"https://example.com/rss",Body:raw,ContentType:"application/rss+xml",HTTPStatus:200}
 data,err:=json.Marshal(diagnostic)
 if err!=nil {t.Fatal(err)}
 var roundTrip rawFeedDiagnostic
 if err:=json.Unmarshal(data,&roundTrip);err!=nil {t.Fatal(err)}
 if !bytes.Equal(roundTrip.Body,raw)||roundTrip.URL!=diagnostic.URL||roundTrip.ContentType!=diagnostic.ContentType||roundTrip.HTTPStatus!=200 {t.Fatalf("feed response changed: %+v",roundTrip)}
}
