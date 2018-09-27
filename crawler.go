package feedcrawler

import (
	"time"

	"github.com/mmcdole/gofeed"
)

var defaultNumWorkers = 3

// Result is a result of a feed crawling.
type Result struct {
	Subscription Subscription
	Feed         *gofeed.Feed
	NewItems     []*gofeed.Item
	Err          error
}

// Crawler is a crawler for RSS and Atom feeds.
type Crawler struct {
	Subscriptions []Subscription
	States        States
	NumWorkers    int
	Parser        *gofeed.Parser
}

// Crawl crawls subscribed feeds.
func (fc *Crawler) Crawl() ([]Result, error) {
	var results []Result
	err := fc.CrawlFunc(func(r Result) {
		results = append(results, r)
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// CrawlFunc crawls subscribed feeds and call the func with each result.
func (fc *Crawler) CrawlFunc(f func(Result)) error {
	if fc.States == nil {
		fc.States = make(States, 0)
	}
	subscriptions := make(chan Subscription, len(fc.Subscriptions))
	results := make(chan Result, len(fc.Subscriptions))
	defer close(results)

	nw := fc.NumWorkers
	if nw <= 0 {
		nw = defaultNumWorkers
	}

	for w := 1; w <= nw; w++ {
		go fc.worker(w, subscriptions, results)
	}

	for _, s := range fc.Subscriptions {
		subscriptions <- s
	}
	close(subscriptions)

	for i := 1; i <= len(fc.Subscriptions); i++ {
		result := <-results
		if fc.States != nil {
			fc.States.UpdateState(result)
		}
		f(result)
	}

	return nil
}

func (fc *Crawler) worker(id int, subscriptions <-chan Subscription, results chan<- Result) {
	var fp *gofeed.Parser
	if fc.Parser != nil {
		fp = fc.Parser
	} else {
		fp = gofeed.NewParser()
	}

	for s := range subscriptions {
		feed, err := fp.ParseURL(s.URI())
		if err != nil {
			results <- Result{
				Subscription: s,
				Feed:         nil,
				NewItems:     nil,
				Err:          err,
			}
		} else {
			results <- Result{
				Subscription: s,
				Feed:         feed,
				NewItems:     fc.selectNewItems(s, feed),
				Err:          nil,
			}
		}
	}
}

func (fc *Crawler) selectNewItems(s Subscription, feed *gofeed.Feed) []*gofeed.Item {
	if fc.States == nil || len(fc.States) == 0 {
		return feed.Items
	}

	var updatedAt time.Time
	if state, ok := fc.States[s.ID()]; ok {
		updatedAt = state.UpdatedAt
	}

	var newItems []*gofeed.Item
	for _, item := range feed.Items {
		published := latestTime(item.PublishedParsed, item.UpdatedParsed)
		if published.Before(updatedAt) || published.Equal(updatedAt) || !s.Filter(item) {
			continue
		}
		newItems = append(newItems, item)
	}
	return newItems
}
