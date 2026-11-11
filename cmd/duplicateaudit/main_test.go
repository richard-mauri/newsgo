package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestBatchAndEndpoint(t *testing.T) {
	b := []byte(`{"schema":"newsgo.scrapergo-batch.v1","urls":["https://example.org/a#top","https://example.org/a","https://example.org/b"]}`)
	urls, err := parseBatch(b)
	if err != nil || len(urls) != 2 || urls[0] != "https://example.org/a" {
		t.Fatalf("urls=%v err=%v", urls, err)
	}
	for _, s := range []string{"http://example.org/api/read", "http://localhost:8083/api/read", "http://127.0.0.1:8083/api/read?x=1", "https://127.0.0.1:8083/api/read"} {
		if _, err := parseEndpoint(s, false); err == nil {
			t.Fatalf("accepted %s", s)
		}
	}
	if _, err := parseEndpoint("http://127.0.0.1:8083/api/read", false); err != nil {
		t.Fatal(err)
	}
}
func TestApprovedPublicEndpoint(t *testing.T) {
	if _, err := parseEndpoint("https://scrapergo.themauris.org/api/read", true); err != nil {
		t.Fatal(err)
	}
	for _, v := range []string{"http://scrapergo.themauris.org/api/read", "https://other.example/api/read", "https://scrapergo.themauris.org:8443/api/read"} {
		if _, err := parseEndpoint(v, true); err == nil {
			t.Errorf("accepted public endpoint %s", v)
		}
	}
}
func TestFindDuplicateAndPreserveNonDuplicates(t *testing.T) {
	a := article{Blocks: []block{
		{Kind: "img", Text: "The northern harbor under cloudy skies. Photo: Observer"},
		{Kind: "p", Text: "The northern harbor under cloudy skies."},
		{Kind: "p", Text: "Photo: Observer"},
		{Kind: "p", Text: "The northern harbor under cloudy skies. Photo: Observer"},
		{Kind: "p", Text: "An unrelated description of this weekend's weather and the boats in port."},
		{Kind: "p", Text: "An unrelated description of this weekend's weather and the boats in port."},
		{Kind: "p", Text: "Brief heading"},
	}}
	f := audit(a)
	if len(f) != 4 {
		t.Fatalf("findings=%+v", f)
	}
	types := map[string]int{}
	for _, v := range f {
		types[v.Type]++
	}
	if types["image_exact"] != 1 || types["image_caption_or_text_overlap"] != 1 || types["exact"] != 1 || types["caption_or_text_overlap"] != 1 {
		t.Fatalf("types=%v", types)
	}
}
func TestBatchHTTPAndCSV(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/read" || r.Method != "POST" {
			t.Errorf("invalid request %v %v", r.Method, r.URL.Path)
		}
		if strings.Contains(r.Header.Get("Content-Type"), "application/json") == false {
			t.Error("missing JSON content type")
		}
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r.Body)
		if strings.Contains(buf.String(), "/bad") {
			w.WriteHeader(422)
			fmt.Fprint(w, `{"kind":"access_error","error":"secret internal information"}`)
			return
		}
		fmt.Fprint(w, `{"blocks":[{"kind":"img","text":"A long caption describing an image of the blue riverbank."},{"kind":"p","text":"A long caption describing an image of the blue riverbank."}]}`)
	}))
	defer server.Close()
	var out bytes.Buffer
	s, err := run(context.Background(), []string{"https://example.org/good?private=one", "https://example.org/bad"}, server.URL+"/api/read", server.Client(), 0, &out)
	if err != nil || s.ReadOK != 1 || s.ReadFailed != 1 || s.Flagged != 1 || s.Findings != 1 {
		t.Fatalf("summary=%+v err=%v", s, err)
	}
	rows, err := csv.NewReader(&out).ReadAll()
	if err != nil || len(rows) != 3 {
		t.Fatalf("csv=%v err=%v", rows, err)
	}
	if rows[1][5] != "image_exact" || rows[2][2] != "read_failed" {
		t.Fatalf("rows=%v", rows)
	}
	if strings.Contains(out.String(), "secret internal information") {
		t.Fatal("error body leaked")
	}
}
func TestPacing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"blocks":[{"kind":"p","text":"Some sufficiently lengthy and original article paragraph."}]}`)
	}))
	defer server.Close()
	start := time.Now()
	var out bytes.Buffer
	_, err := run(context.Background(), []string{"https://example.org/a", "https://example.org/b"}, server.URL+"/api/read", server.Client(), 20*time.Millisecond, &out)
	if err != nil || time.Since(start) < 20*time.Millisecond {
		t.Fatalf("missing pacing: %v", err)
	}
}
