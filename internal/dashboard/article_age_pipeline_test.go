package dashboard

import (
 "context"
 "testing"
 "time"
 "github.com/richard-mauri/newsgo/internal/feeds"
)

func TestAgeFilteringSkipsOldBeforeClassification(t *testing.T) {
 now := time.Now()
 articles := []feeds.Article{
  {Title:"Old daily news",URL:"https://example.org/old",Published:now.Add(-20*24*time.Hour)},
  {Title:"Undated story",URL:"https://example.org/undated"},
 }
 decisions := classifyFeedArticlesWithAge(context.Background(),"http://127.0.0.1:1",articles,3)
 if decisions[0]!="too-old" {t.Fatalf("old story classified as %q, want too-old",decisions[0])}
 if decisions[1]=="too-old" {t.Fatal("undated story filtered by age")}
 if len(decisions)!=len(articles) {t.Fatal("classification results not aligned with feed")}
}
