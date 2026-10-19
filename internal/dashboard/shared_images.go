package dashboard

import (
	"context"
	"sync"
	"time"

	"github.com/richard-mauri/newsgo/internal/feeds"
)

// sharedThumbnailCache avoids repeatedly fetching publisher article metadata on
// every 15-minute shared-feed refresh. Cache misses as well as successes.
type sharedThumbnailEntry struct {
	url string
	at  time.Time
}

var sharedThumbnailCache = struct {
	sync.Mutex
	entries map[string]sharedThumbnailEntry
}{entries: make(map[string]sharedThumbnailEntry)}

// enrichSharedImages supplements missing RSS image metadata for built-in feeds.
// Work is bounded independently of the number of feed articles, and only the
// most recent image-less articles per publisher are probed during a refresh.
func enrichSharedImages(ctx context.Context, articles []feeds.Article, lookup func(context.Context, string) string) {
	const perPublisher = 8
	const maxConcurrent = 8
	counts := map[string]int{}
	type job struct {
		index   int
		address string
	}
	jobs := []job{}
	now := time.Now()
	sharedThumbnailCache.Lock()
	for i := range articles {
		if articles[i].Image != "" || !safeURL(articles[i].URL) {
			continue
		}
		if c, ok := sharedThumbnailCache.entries[articles[i].URL]; ok && now.Sub(c.at) < 4*time.Hour {
			articles[i].Image = c.url
			continue
		}
		if counts[articles[i].Source] >= perPublisher {
			continue
		}
		counts[articles[i].Source]++
		jobs = append(jobs, job{i, articles[i].URL})
	}
	sharedThumbnailCache.Unlock()
	if len(jobs) == 0 {
		return
	}
	limited, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	for _, j := range jobs {
		wg.Add(1)
		go func(j job) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-limited.Done():
				return
			}
			image := lookup(limited, j.address)
			if limited.Err() != nil {
				return
			} // Do not cache canceled or timed-out attempts.
			articles[j.index].Image = image
			sharedThumbnailCache.Lock()
			if len(sharedThumbnailCache.entries) > 1200 {
				sharedThumbnailCache.entries = make(map[string]sharedThumbnailEntry)
			}
			sharedThumbnailCache.entries[j.address] = sharedThumbnailEntry{image, time.Now()}
			sharedThumbnailCache.Unlock()
		}(j)
	}
	wg.Wait()
}
