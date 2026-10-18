package main

import (
	"os"
	"strings"
	"testing"
)

func TestPublicScraperGoDefault(t *testing.T) {
	b, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"https://scrapergo.themauris.org", "ScraperGo reader base URL"`) {
		t.Fatal("default must use public HTTPS reader")
	}
}
