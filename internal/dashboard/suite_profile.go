package dashboard

import (
 "encoding/json"
 "errors"
)

// suiteProfileData is the schema stored in Homepage's /api/storage/newsgo
// record.data. Its revision belongs to the outer Homepage record, not this payload.
type suiteProfileData struct {
 Subscriptions []syncSubscription `json:"subscriptions"`
 ExcludedDomains []string `json:"excludedDomains"`
 ExcludedArticleDomains []string `json:"excludedArticleDomains"`
 MaxArticleAgeDays int `json:"maxArticleAgeDays,omitempty"`
 MaximumLeadStories int `json:"maximumLeadStories,omitempty"`
 LeadOnlyAllSubscriptions bool `json:"leadOnlyAllSubscriptions,omitempty"`
 LastReadingSelection string `json:"lastReadingSelection,omitempty"`
}

func suiteDataFromLegacy(p storedPrefs) (suiteProfileData,error) {
 if !validPrefs(p) {return suiteProfileData{},errors.New("invalid NewsGo reader profile")}
 return suiteProfileData{
  Subscriptions:append([]syncSubscription(nil),p.Subscriptions...),
  ExcludedDomains:append([]string(nil),p.ExcludedDomains...),
  ExcludedArticleDomains:append([]string(nil),p.ExcludedArticleDomains...),
  MaxArticleAgeDays:p.MaxArticleAgeDays,
  MaximumLeadStories:p.MaximumLeadStories,
  LeadOnlyAllSubscriptions:p.LeadOnlyAllSubscriptions,
  LastReadingSelection:p.LastReadingSelection,
 },nil
}

func legacyPrefsFromSuite(data json.RawMessage,revision uint64) (storedPrefs,error){
 var p suiteProfileData
 if err:=json.Unmarshal(data,&p);err!=nil{return storedPrefs{},err}
 legacy:=storedPrefs{
  Revision:revision,
  Subscriptions:p.Subscriptions,
  ExcludedDomains:p.ExcludedDomains,
  ExcludedArticleDomains:p.ExcludedArticleDomains,
  MaxArticleAgeDays:p.MaxArticleAgeDays,
  MaximumLeadStories:p.MaximumLeadStories,
  LeadOnlyAllSubscriptions:p.LeadOnlyAllSubscriptions,
  LastReadingSelection:p.LastReadingSelection,
 }
 if !validPrefs(legacy){return storedPrefs{},errors.New("invalid Homepage NewsGo profile")}
 return legacy,nil
}
