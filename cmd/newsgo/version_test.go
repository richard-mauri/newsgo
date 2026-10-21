package main

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/richard-mauri/newsgo/internal/dashboard"
)

// Exercise the actual CLI entry point, not just the dashboard's helper.
func TestVersionFlag(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-version", "-feeds", "/nonexistent-newsgo-feed-file")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("-version failed without a feed file: %v: %s", err, out)
	}
	expected := "NewsGo v" + dashboard.Version()
	if got := strings.TrimSpace(string(out)); got != expected {
		t.Fatalf("-version output = %q, want %q", got, expected)
	}
}
