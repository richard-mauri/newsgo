package dashboard

import "testing"

func TestArticleSettingsDefaultsAndSafetyBounds(t *testing.T) {
 cfg:=defaultArticleSettings()
 if cfg.MaxArticlesPerFeed != 60 || cfg.ClassificationTimeoutSeconds != 7 || cfg.ClassificationConcurrency != 8 || cfg.FeedCacheMinutes != 10 { t.Fatalf("unexpected defaults: %+v",cfg) }
 if !validArticleSettings(cfg) { t.Fatal("default settings rejected") }
 cases:=[]ArticleSettings{
 {MaxArticlesPerFeed:0,ClassificationTimeoutSeconds:7,ClassificationConcurrency:8,FeedCacheMinutes:10},
 {MaxArticlesPerFeed:251,ClassificationTimeoutSeconds:7,ClassificationConcurrency:8,FeedCacheMinutes:10},
 {MaxArticlesPerFeed:60,ClassificationTimeoutSeconds:31,ClassificationConcurrency:8,FeedCacheMinutes:10},
 {MaxArticlesPerFeed:60,ClassificationTimeoutSeconds:7,ClassificationConcurrency:0,FeedCacheMinutes:10},
 {MaxArticlesPerFeed:60,ClassificationTimeoutSeconds:7,ClassificationConcurrency:8,FeedCacheMinutes:61},
 }
 for _,c:=range cases { if validArticleSettings(c) { t.Errorf("unsafe settings accepted: %+v",c) } }
}

func TestArticleSettingsFeedCacheUpdatesAtRuntime(t *testing.T) {
 previous:=currentArticleSettings()
 defer activeArticleSettings.Store(previous)
 cfg:=defaultArticleSettings()
 cfg.FeedCacheMinutes=3
 activeArticleSettings.Store(cfg)
 if got:=personalCacheTTL(cacheEntry{});got.Minutes()!=3 { t.Fatalf("cache TTL = %s",got) }
}
