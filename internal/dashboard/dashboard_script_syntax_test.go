package dashboard

import (
 "os"
 "os/exec"
 "path/filepath"
 "strings"
 "testing"
)

func TestDashboardJavaScriptSyntax(t *testing.T) {
 if strings.Contains(personalJS, "async async function") {
  t.Fatal("duplicate async keyword prevents dashboard JavaScript from loading")
 }
 node, err := exec.LookPath("node")
 if err != nil {
  t.Log("node unavailable: browser JavaScript syntax check skipped")
  return
 }
 dir := t.TempDir()
 path := filepath.Join(dir, "personal.js")
 if err := os.WriteFile(path, []byte(personalJS), 0600); err != nil { t.Fatal(err) }
 output, err := exec.Command(node, "--check", path).CombinedOutput()
 if err != nil { t.Fatalf("dashboard JavaScript syntax invalid: %v\n%s", err, output) }
}
