package feeds

import "testing"

func TestEmbeddedFeedImages(t *testing.T) {
	tests := []struct{ name, data, want string }{
		{"RSS description", `<rss><channel><item><title>Dogs</title><link>https://example.com/dogs</link><description><![CDATA[<p>Story</p><img src="https://images.example.com/dog.jpg">]]></description></item></channel></rss>`, "https://images.example.com/dog.jpg"},
		{"RSS content encoded", `<rss xmlns:content="http://purl.org/rss/1.0/modules/content/"><channel><item><title>Dogs</title><link>https://example.com/dogs</link><content:encoded><![CDATA[<img data-src="/photo.jpg">]]></content:encoded></item></channel></rss>`, "https://example.com/photo.jpg"},
		{"Atom summary", `<feed xmlns="http://www.w3.org/2005/Atom"><entry><title>Dogs</title><link href="https://example.com/dogs"/><summary type="html">&lt;img src="https://images.example.com/dog.jpg"/&gt;</summary></entry></feed>`, "https://images.example.com/dog.jpg"},
		{"No invented image", `<rss><channel><item><title>Dogs</title><link>https://example.com/dogs</link><description>No picture</description></item></channel></rss>`, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, items, e := Parse([]byte(tc.data))
			if e != nil {
				t.Fatal(e)
			}
			if len(items) != 1 || items[0].image != tc.want {
				t.Fatalf("got %+v want %s", items, tc.want)
			}
		})
	}
}
