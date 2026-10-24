package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/richard-mauri/newsgo/internal/dashboard"
	"github.com/richard-mauri/newsgo/internal/feeds"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	feedPath := flag.String("feeds", "feeds.json", "feed source JSON file")
	syncDir := flag.String("sync-dir", "data", "local directory for approved tokens and synced preferences")
	addr := flag.String("listen", "127.0.0.1:8084", "HTTP listen address")
	reader := flag.String("scrapergo", "https://scrapergo.themauris.org", "ScraperGo reader base URL")
	showVersion := flag.Bool("version", false, "print the executable version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("NewsGo v%s · %s\n", dashboard.Version(), dashboard.BuildLabel())
		return
	}
	b, e := os.ReadFile(*feedPath)
	if e != nil {
		log.Fatal(e)
	}
	var sources []feeds.Source
	if e = json.Unmarshal(b, &sources); e != nil {
		log.Fatal(e)
	}
	client := &http.Client{Timeout: 12 * time.Second, CheckRedirect: func(r *http.Request, via []*http.Request) error {
		if len(via) >= 5 || !feeds.ValidURL(r.URL.String()) {
			return http.ErrUseLastResponse
		}
		return nil
	}}
	app, e := dashboard.New(sources, *reader, client)
	if e != nil {
		log.Fatal(e)
	}
	syncStore, e := dashboard.NewSyncStore(*syncDir)
	if e != nil {
		log.Fatal(e)
	}
	if e = syncStore.ActivateServerSettings(); e != nil {
		log.Fatal(e)
	}
	app.SetSyncStore(syncStore)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// Feed catalog is for discovery only. Browser requests subscribed feeds via /api/personal.
	srv := &http.Server{Addr: *addr, Handler: app, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdown)
	}()
	log.Printf("NewsGo v%s · %s: http://%s", dashboard.Version(), dashboard.BuildLabel(), *addr)
	if e = srv.ListenAndServe(); e != nil && e != http.ErrServerClosed {
		log.Fatal(e)
	}
}
