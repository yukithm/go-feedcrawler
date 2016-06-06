package feedcrawler

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/mmcdole/gofeed"
)

// FeedID is an identifier of a feed.
type FeedID string

// Subscription is a feed configuration to be subscribed.
type Subscription struct {
	ID                FeedID
	URI               string
	TitleFilter       *regexp.Regexp
	DescriptionFilter *regexp.Regexp
	ContentFilter     *regexp.Regexp
	AuthorFilter      *regexp.Regexp
	CategoryFilter    *regexp.Regexp
	Filter            func(*gofeed.Item) bool
	Meta              interface{}
}

// Result is a result of a feed crawling.
type Result struct {
	Subscription Subscription
	Feed         *gofeed.Feed
	NewItems     []*gofeed.Item
	Err          error
}

type state struct {
	CrawledAt time.Time `json:"crawled_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

var defaultNumWorkers = 3

// Crawler is a crawler for RSS and Atom feeds.
type Crawler struct {
	NumWorkers    int
	Parser        *gofeed.Parser
	Subscriptions []Subscription
	StateFile     string
	states        map[FeedID]*state
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
	if err := fc.loadState(); err != nil {
		return err
	}

	subscriptions := make(chan Subscription, len(fc.Subscriptions))
	results := make(chan Result, len(fc.Subscriptions))

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
		fc.updateState(result)
		f(result)
	}

	return fc.saveState()
}

func (fc *Crawler) worker(id int, subscriptions <-chan Subscription, result chan<- Result) {
	var fp *gofeed.Parser
	if fc.Parser != nil {
		fp = fc.Parser
	} else {
		fp = gofeed.NewParser()
	}

	for s := range subscriptions {
		feed, err := fp.ParseURL(s.URI)
		if err != nil {
			result <- Result{
				Subscription: s,
				Feed:         nil,
				NewItems:     nil,
				Err:          err,
			}
		} else {
			result <- Result{
				Subscription: s,
				Feed:         feed,
				NewItems:     fc.selectNewItems(s, feed),
				Err:          nil,
			}
		}
	}
}

func (fc *Crawler) selectNewItems(s Subscription, feed *gofeed.Feed) []*gofeed.Item {
	var updatedAt time.Time
	if state, ok := fc.states[s.ID]; ok {
		updatedAt = state.UpdatedAt
	}

	var newItems []*gofeed.Item
	for _, item := range feed.Items {
		published := latestTime(item.PublishedParsed, item.UpdatedParsed)
		if published.Before(updatedAt) || published.Equal(updatedAt) || fc.filtered(s, item) {
			continue
		}
		newItems = append(newItems, item)
	}
	return newItems
}

func (fc *Crawler) filtered(s Subscription, item *gofeed.Item) bool {
	return !fc.matched(s, item)
}

func (fc *Crawler) matched(s Subscription, item *gofeed.Item) bool {
	if s.AuthorFilter != nil && !matchAuthor(s.AuthorFilter, item.Author) {
		return false
	}
	if s.CategoryFilter != nil && !matchCategories(s.CategoryFilter, item.Categories) {
		return false
	}
	if s.TitleFilter != nil && !matchString(s.TitleFilter, item.Title) {
		return false
	}
	if s.DescriptionFilter != nil && !matchString(s.DescriptionFilter, item.Description) {
		return false
	}
	if s.ContentFilter != nil && !matchString(s.ContentFilter, item.Content) {
		return false
	}
	if s.Filter != nil && !s.Filter(item) {
		return false
	}
	return true
}

func matchAuthor(r *regexp.Regexp, author *gofeed.Person) bool {
	if author == nil {
		return false
	}
	return matchString(r, author.Name) || matchString(r, author.Email)
}

func matchCategories(r *regexp.Regexp, categories []string) bool {
	for _, cat := range categories {
		if matchString(r, cat) {
			return true
		}
	}

	return false
}

func matchString(r *regexp.Regexp, value string) bool {
	if value != "" && r.MatchString(value) {
		return true
	}
	return false
}

func (fc *Crawler) updateState(result Result) {
	id := result.Subscription.ID

	st, ok := fc.states[id]
	if !ok {
		st = &state{}
		fc.states[id] = st
	}

	st.CrawledAt = time.Now()
	if result.Err == nil && result.Feed != nil {
		published := latestFeedTime(result.Feed)
		if published != nil {
			st.UpdatedAt = published.Local()
		}
	}
}

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

func (fc *Crawler) loadState() error {
	if fc.states == nil {
		fc.states = make(map[FeedID]*state, 0)
	}

	if fc.StateFile == "" {
		return nil
	}

	f, err := os.Open(fc.StateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)

	if err := decoder.Decode(&fc.states); err != nil {
		return err
	}

	return nil
}

func (fc *Crawler) saveState() error {
	if fc.StateFile == "" {
		return nil
	}

	buf, err := json.MarshalIndent(fc.states, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fc.StateFile, buf, 0666); err != nil {
		return err
	}

	return nil
}
