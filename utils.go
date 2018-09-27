package feedcrawler

import (
	"time"

	"github.com/mmcdole/gofeed"
)

func latestFeedTime(feed *gofeed.Feed) *time.Time {
	if feed == nil {
		return nil
	}

	t := latestTime(feed.PublishedParsed, feed.UpdatedParsed)
	for _, item := range feed.Items {
		it := latestTime(item.PublishedParsed, item.UpdatedParsed)
		if it != nil && (t == nil || it.After(*t)) {
			t = it
		}
	}

	return t
}

func latestTime(a, b *time.Time) *time.Time {
	if a == nil {
		return b
	} else if b == nil {
		return a
	}

	if a.After(*b) {
		return a
	}
	return b
}
