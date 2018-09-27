package feedcrawler

import (
	"regexp"

	"github.com/mmcdole/gofeed"
)

// FeedID is an identifier of a feed.
type FeedID string

// Subscription is a subscription information of a feed.
type Subscription interface {
	ID() FeedID
	URI() string
	Filter(*gofeed.Item) bool
}

// SimpleSubscription is a feed configuration to be subscribed.
type SimpleSubscription struct {
	FeedID            FeedID
	FeedURI           string
	TitleFilter       *regexp.Regexp
	DescriptionFilter *regexp.Regexp
	ContentFilter     *regexp.Regexp
	AuthorFilter      *regexp.Regexp
	CategoryFilter    *regexp.Regexp
	FilterFunc        func(*gofeed.Item) bool
}

// ID returns the feed ID.
func (s *SimpleSubscription) ID() FeedID { return s.FeedID }

// URI returns the feed URI.
func (s *SimpleSubscription) URI() string { return s.FeedURI }

// Filter returns true if the item is acceptable.
func (s *SimpleSubscription) Filter(item *gofeed.Item) bool {
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
	if s.FilterFunc != nil && !s.FilterFunc(item) {
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
