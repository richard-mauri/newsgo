package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/richard-mauri/newsgo/internal/feeds"
)

func main() {
	path := flag.String("feeds", "feeds.json", "JSON list of feeds")
	out := flag.String("out", "report.json", "report output filename")
	timeout := flag.Duration("timeout", 12*time.Second, "per-feed timeout")
	parallel := flag.Int("parallel", 4, "maximum simultaneous requests")
	scraperBase := flag.String("scrapergo", "", "optional ScraperGo base URL, e.g. http://127.0.0.1:8083")
	samples := flag.Int("samples", 2, "maximum article samples per feed (with -scrapergo)")
	flag.Parse()
	if *samples < 1 || *samples > 5 {
		log.Fatal("samples must be between 1 and 5")
	}
	endpoint := ""
	if *scraperBase != "" {
		var e error
		endpoint, e = feeds.ScraperGoEndpoint(*scraperBase)
		if e != nil {
			log.Fatal(e)
		}
	}
	if *parallel < 1 || *parallel > 16 {
		log.Fatal("parallel must be between 1 and 16")
	}
	b, err := os.ReadFile(*path)
	if err != nil {
		log.Fatal(err)
	}
	var sources []feeds.Source
	if err = json.Unmarshal(b, &sources); err != nil {
		log.Fatal(err)
	}
	if err = feeds.ValidateSources(sources); err != nil {
		log.Fatal(err)
	}
	client := &http.Client{Timeout: *timeout, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return fmt.Errorf("too many redirects")
		}
		if !feeds.ValidURL(req.URL.String()) {
			return fmt.Errorf("invalid redirect URL")
		}
		return nil
	}}
	results := make([]feeds.Report, len(sources))
	sem := make(chan struct{}, *parallel)
	var wg sync.WaitGroup
	for i, s := range sources {
		wg.Add(1)
		go func(i int, s feeds.Source) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			results[i] = feeds.Check(context.Background(), client, s, time.Now())
			if endpoint != "" && results[i].Items > 0 {
				urls, e := feeds.ArticleURLs(context.Background(), client, s.URL)
				if e != nil {
					results[i].Compatibility = &feeds.Compatibility{Results: []feeds.Sample{{Status: "other", Detail: e.Error()}}}
				} else {
					urls = urls[:min(len(urls), *samples)]
					results[i].Compatibility = feeds.SampleScraperGo(context.Background(), client, endpoint, urls, *samples)
				}
			}
		}(i, s)
	}
	wg.Wait()
	feeds.SortReports(results)
	data, err := json.MarshalIndent(struct {
		GeneratedAt string         `json:"generated_at"`
		Feeds       []feeds.Report `json:"feeds"`
	}{time.Now().UTC().Format(time.RFC3339), results}, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err = os.WriteFile(*out, append(data, '\n'), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-25s %-15s %5s %5s %5s %5s  %s\n", "SOURCE", "STATUS", "ITEMS", "TITLE", "LINK", "DATE", "LATEST (UTC)")
	for _, r := range results {
		fmt.Printf("%-25.25s %-15s %5d %5d %5d %5d  %s\n", r.Name, r.Status, r.Items, r.WithTitle, r.WithLink, r.WithDate, r.Latest)
		if r.Error != "" {
			fmt.Printf("  %s\n", r.Error)
		}
	}
	if endpoint != "" {
		fmt.Println("\nSCRAPERGO SAMPLES (readable / restricted / other)")
		for _, r := range results {
			if r.Compatibility != nil {
				c := r.Compatibility
				fmt.Printf("%-25.25s %d / %d / %d\n", r.Name, c.Readable, c.Restricted, c.Other)
			}
		}
	}
	fmt.Printf("\nFull report: %s\n", *out)
}
