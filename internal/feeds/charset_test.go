package feeds

import (
 "strings"
 "testing"
)

func TestParseISO88591RSS(t *testing.T) {
 prefix := "<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?><rss><channel><item><title>Caf"
 data := append([]byte(prefix), 0xe9)
 data = append(data, []byte("</title><link>https://example.org/story</link></item></channel></rss>")...)
 format, items, err := Parse(data)
 if err != nil { t.Fatal(err) }
 if format != "RSS" || len(items) != 1 || items[0].title != "Café" {
  t.Fatalf("format=%q items=%+v",format,items)
 }
}

func TestParseRejectsUnsupportedEncoding(t *testing.T) {
 _, _, err := Parse([]byte(`<?xml version="1.0" encoding="X-UNKNOWN"?><rss/>`))
 if err == nil || !strings.Contains(err.Error(),"unsupported XML encoding") { t.Fatalf("unexpected error: %v",err) }
}
