package dashboard

import ("strings";"testing")

func TestReadingSelectionHomepageAccountContract(t *testing.T){
 for _,needle:=range []string{
  "lastReadingSelection:validReadingSelection(activeReadingSelection)||''",
  "restoreReadingSelection(restored)",
  "rememberReadingSelection('all')",
  "rememberReadingSelection(url)",
  "restoreReadingSelection(localStorage.getItem(readingSelectionKey))",
  "const restored=validReadingSelection(data.lastReadingSelection)",
  "async function initializeHomepageAccount(owner,token)",
  "applyAccountPreferences(record.data)",
 }{if !strings.Contains(personalJS,needle){t.Errorf("missing Homepage reading-selection contract %q",needle)}}
}
