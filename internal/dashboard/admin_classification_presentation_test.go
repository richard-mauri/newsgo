package dashboard

import (
 "strings"
 "testing"
)

func TestAdminClassificationDiagnosticsPresentationContract(t *testing.T) {
 required:=[]string{
  "const classification=result.classification||'unknown'",
  "Classification: ",
  "Request failure: ",
  "ScraperGo could not confirm a readable article.",
  "Classification request duration: ",
  "lines.join('\\n')",
 }
 for _,s:=range required {if !strings.Contains(personalJS,s){t.Errorf("missing diagnostic contract %q",s)}}
 forbidden:=[]string{"'\\\\nFailure: '","'\\\\nHTTP status: '","'\\\\nTimeout: '","'\\nTimeout: '","'\\nHTTP status: '"}
 for _,s:=range forbidden{if strings.Contains(personalJS,s){t.Errorf("misleading diagnostic remains: %q",s)}}
}
