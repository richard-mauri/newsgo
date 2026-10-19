package feeds

import "testing"

func TestArticleImageMetadata(t *testing.T) {
	tests := []struct{ name, body, want string }{
		{"RSS media thumbnail", `<rss xmlns:media="http://search.yahoo.com/mrss/"><channel><item><title>Dog</title><link>https://example.org/dog</link><media:thumbnail url="https://cdn.example.org/dog.jpg"/></item></channel></rss>`, "https://cdn.example.org/dog.jpg"},
		{"RSS image enclosure", `<rss><channel><item><title>Dog</title><link>https://example.org/dog</link><enclosure type="image/jpeg" url="https://cdn.example.org/dog.jpg"/></item></channel></rss>`, "https://cdn.example.org/dog.jpg"},
		{"RSS no image", `<rss><channel><item><title>Dog</title><link>https://example.org/dog</link></item></channel></rss>`, ""},
		{"Atom media image", `<feed xmlns="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/"><entry><title>Dog</title><link href="https://example.org/dog"/><media:content url="https://cdn.example.org/dog.png" type="image/png" /></entry></feed>`, "https://cdn.example.org/dog.png"},
		{"JSON feed image", `{"version":"https://jsonfeed.org/version/1.1","items":[{"title":"Dog","url":"https://example.org/dog","image":"https://cdn.example.org/dog.jpg"}]}`, "https://cdn.example.org/dog.jpg"},
		{"Reject insecure image", `<rss><channel><item><title>Dog</title><link>https://example.org/dog</link><enclosure type="image/jpeg" url="http://cdn.example.org/dog.jpg"/></item></channel></rss>`, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, items, err := Parse([]byte(tc.body))
			if err != nil {
				t.Fatal(err)
			}
			if len(items) != 1 || items[0].image != tc.want {
				t.Fatalf("items=%+v, want %q", items, tc.want)
			}
		})
	}
}
