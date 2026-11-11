// duplicateaudit checks NewsGo's exported article URLs against ScraperGo's
// canonical read API. It never contacts publishers directly.
package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const schema = "newsgo.scrapergo-batch.v1"
const maxURLs = 500
const maxResponseBytes = 12 << 20

type batch struct {
	Schema     string   `json:"schema"`
	CapturedAt string   `json:"capturedAt"`
	URLs       []string `json:"urls"`
}
type block struct {
	Kind string `json:"kind"`
	Text string `json:"text"`
}
type article struct {
	Title  string  `json:"title"`
	Method string  `json:"method"`
	Blocks []block `json:"blocks"`
}
type finding struct {
	First      int
	Second     int
	Type       string
	FirstKind  string
	SecondKind string
	Sample     string
}
type summary struct {
	Total      int
	ReadOK     int
	ReadFailed int
	Flagged    int
	Findings   int
}

func parseBatch(data []byte) ([]string, error) {
	if len(data) > 1<<20 {
		return nil, errors.New("input exceeds 1 MiB")
	}
	var b batch
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&b); err != nil {
		return nil, fmt.Errorf("invalid batch JSON: %w", err)
	}
	var extra any
	if dec.Decode(&extra) != io.EOF {
		return nil, errors.New("batch must contain one JSON object")
	}
	if b.Schema != schema || len(b.URLs) == 0 || len(b.URLs) > maxURLs {
		return nil, errors.New("wrong schema or URL count outside 1..500")
	}
	out := []string{}
	seen := map[string]bool{}
	for _, raw := range b.URLs {
		if len(raw) > 2048 || strings.ContainsAny(raw, "\r\n\t") {
			return nil, errors.New("invalid URL in batch")
		}
		u, err := url.Parse(raw)
		if err != nil || u == nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || u.User != nil {
			return nil, errors.New("invalid HTTP(S) article URL")
		}
		u.Fragment = ""
		key := u.String()
		if !seen[key] {
			seen[key] = true
			out = append(out, key)
		}
	}
	return out, nil
}
func parseEndpoint(raw string, allowPublic bool) (string, error) {
	u, err := url.Parse(raw)
	if err != nil || u == nil || u.Path != "/api/read" || u.RawQuery != "" || u.Fragment != "" || u.User != nil {
		return "", errors.New("ScraperGo endpoint must be a loopback HTTP /api/read URL")
	}
	if allowPublic && u.Scheme == "https" && u.Hostname() == "scrapergo.themauris.org" && u.Port() == "" {
		return u.String(), nil
	}
	if u.Scheme != "http" {
		return "", errors.New("non-loopback ScraperGo requires -public with its approved HTTPS hostname")
	}
	ip := net.ParseIP(u.Hostname())
	if ip == nil || !ip.IsLoopback() {
		return "", errors.New("ScraperGo endpoint must have a numeric loopback IP")
	}
	return u.String(), nil
}
func normalize(s string) string {
	return strings.Join(strings.Fields(strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return ' '
		}
		return unicode.ToLower(r)
	}, s)), " ")
}
func short(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	r := []rune(s)
	if len(r) > 130 {
		return string(r[:130]) + "…"
	}
	return s
}
func audit(a article) []finding {
	findings := []finding{}
	for i, b := range a.Blocks {
		if b.Kind == "table" || b.Kind == "topic-link" {
			continue
		}
		x := normalize(b.Text)
		if len([]rune(x)) < 24 {
			continue
		} // Short labels/credits alone produce many false positives.
		for j := i + 1; j < len(a.Blocks); j++ {
			c := a.Blocks[j]
			if c.Kind == "table" || c.Kind == "topic-link" {
				continue
			}
			y := normalize(c.Text)
			if len([]rune(y)) < 24 {
				continue
			}
			kind := ""
			if x == y {
				kind = "exact"
			} else if len([]rune(x)) >= 32 && len([]rune(y)) >= 32 {
				shorter, longer := x, y
				if len([]rune(shorter)) > len([]rune(longer)) {
					shorter, longer = longer, shorter
				}
				// Requires a long, substantive overlap rather than a common word or phrase.
				if len([]rune(shorter)) >= 30 && (strings.HasPrefix(longer, shorter+" ") || strings.HasSuffix(longer, " "+shorter)) {
					kind = "caption_or_text_overlap"
				}
			}
			if kind == "" {
				continue
			}
			if (b.Kind == "img") != (c.Kind == "img") {
				kind = "image_" + kind
			}
			findings = append(findings, finding{First: i + 1, Second: j + 1, Type: kind, FirstKind: b.Kind, SecondKind: c.Kind, Sample: short(b.Text)})
		}
	}
	return findings
}
func fetch(ctx context.Context, client *http.Client, endpoint, address string) (article, int, error) {
	payload, _ := json.Marshal(map[string]string{"url": address})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return article{}, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := client.Do(req)
	if err != nil {
		return article{}, 0, err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		io.Copy(io.Discard, io.LimitReader(r.Body, 1024))
		return article{}, r.StatusCode, fmt.Errorf("HTTP %d", r.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(r.Body, maxResponseBytes+1))
	if err != nil {
		return article{}, r.StatusCode, err
	}
	if len(data) > maxResponseBytes {
		return article{}, r.StatusCode, errors.New("response exceeds 12 MiB")
	}
	var a article
	if err := json.Unmarshal(data, &a); err != nil {
		return article{}, r.StatusCode, errors.New("invalid article JSON")
	}
	if len(a.Blocks) == 0 {
		return article{}, r.StatusCode, errors.New("article contains no blocks")
	}
	return a, r.StatusCode, nil
}
func run(ctx context.Context, urls []string, endpoint string, client *http.Client, delay time.Duration, w io.Writer) (summary, error) {
	csvw := csv.NewWriter(w)
	if err := csvw.Write([]string{"article_index", "article_url", "status", "http_status", "block_count", "finding_type", "block_1", "block_2", "kind_1", "kind_2", "sample"}); err != nil {
		return summary{}, err
	}
	s := summary{Total: len(urls)}
	for i, address := range urls {
		if i > 0 && delay > 0 {
			select {
			case <-ctx.Done():
				return s, ctx.Err()
			case <-time.After(delay):
			}
		}
		a, status, err := fetch(ctx, client, endpoint, address)
		if err != nil {
			s.ReadFailed++
			_ = csvw.Write([]string{strconv.Itoa(i + 1), address, "read_failed", strconv.Itoa(status), "0", "", "", "", "", "", ""})
		} else {
			s.ReadOK++
			matches := audit(a)
			if len(matches) > 0 {
				s.Flagged++
				s.Findings += len(matches)
			}
			if len(matches) == 0 {
				_ = csvw.Write([]string{strconv.Itoa(i + 1), address, "ok", strconv.Itoa(status), strconv.Itoa(len(a.Blocks)), "", "", "", "", "", ""})
			}
			for _, f := range matches {
				_ = csvw.Write([]string{strconv.Itoa(i + 1), address, "duplicate_candidate", strconv.Itoa(status), strconv.Itoa(len(a.Blocks)), f.Type, strconv.Itoa(f.First), strconv.Itoa(f.Second), f.FirstKind, f.SecondKind, f.Sample})
			}
		}
		csvw.Flush()
		if err := csvw.Error(); err != nil {
			return s, err
		}
		fmt.Fprintf(os.Stderr, "%d/%d read_ok=%d failed=%d flagged=%d\n", i+1, len(urls), s.ReadOK, s.ReadFailed, s.Flagged)
	}
	return s, nil
}
func main() {
	input := flag.String("input", "", "NewsGo JSON article URL export")
	output := flag.String("output", "duplicate-audit.csv", "local CSV report (may contain article URLs and excerpts)")
	endpoint := flag.String("scrapergo", "http://127.0.0.1:8083/api/read", "ScraperGo read endpoint")
	public := flag.Bool("public", false, "permit approved public HTTPS ScraperGo hostname")
	limit := flag.Int("limit", 10, "number of URLs to audit (max 500)")
	delay := flag.Duration("delay", time.Second, "pause between article requests")
	timeout := flag.Duration("timeout", 45*time.Second, "per-request timeout")
	flag.Parse()
	if *input == "" || *limit < 1 || *limit > maxURLs || *delay < 0 || *delay > time.Minute || *timeout < time.Second || *timeout > 3*time.Minute || flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "Usage: go run ./cmd/duplicateaudit -input ~/Downloads/newsgo-scrapergo-urls.json [-limit 10] [-output duplicate-audit.csv]")
		os.Exit(2)
	}
	ep, err := parseEndpoint(*endpoint, *public)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	data, err := os.ReadFile(*input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	urls, err := parseBatch(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if len(urls) > *limit {
		urls = urls[:*limit]
	}
	f, err := os.OpenFile(*output, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot create report (choose a new -output path):", err)
		os.Exit(1)
	}
	defer f.Close()
	client := &http.Client{Timeout: *timeout, CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }}
	s, err := run(context.Background(), urls, ep, client, *delay, f)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Audit interrupted:", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Finished: %d checked, %d read successfully, %d failed, %d articles flagged, %d matches. Report: %s\n", s.Total, s.ReadOK, s.ReadFailed, s.Flagged, s.Findings, *output)
	if s.ReadFailed > 0 {
		os.Exit(1)
	}
}
