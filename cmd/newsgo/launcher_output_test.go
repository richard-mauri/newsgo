package main

import (
 "os"
 "strings"
 "testing"
)

func TestLauncherPrintsExecutableVersionOnlyOnce(t *testing.T) {
 source, err := os.ReadFile("../../newsgo.sh")
 if err != nil { t.Fatal(err) }
 script := string(source)
 if !strings.Contains(script, `if ! version_output="$("$APP" -version 2>&1)"; then`) {
  t.Fatal("launcher must verify the actual executable before stopping the previous process")
 }
 if !strings.Contains(script, `printf 'NewsGo executable version: %s\n' "$version_output"`) {
  t.Fatal("launcher must display the verified version")
 }
 if strings.Count(script, `printf 'NewsGo executable version: %s\n' "$version_output"`) != 1 {
  t.Fatalf("launcher version output uses changed unexpectedly: %d",strings.Count(script,`printf 'NewsGo executable version: %s\n' "$version_output"`))
 }
 if !strings.Contains(script,`echo "Started NewsGo on 127.0.0.1:8084 and its dedicated Cloudflare Tunnel."`) {
  t.Fatal("launcher must still report the start and endpoint")
 }
}
