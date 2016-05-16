package feedcrawler

import (
	"fmt"
	"regexp"
)

// Feeds is a map of Feed.
type Feeds map[string]Feed

// Subscriptions returns an array of Subscription.
func (fs Feeds) Subscriptions() ([]Subscription, error) {
	var subscriptions []Subscription

	for id, feed := range fs {
		s, err := feed.Subscription(id)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, s)
	}

	return subscriptions, nil
}

// Feed is a feed configuration to be subscribed.
type Feed struct {
	URI               string `toml:"uri"`
	TitleFilter       string `toml:"title_filter,omitempty"`
	DescriptionFilter string `toml:"description_filter,omitempty"`
	ContentFilter     string `toml:"content_filter,omitempty"`
	AuthorFilter      string `toml:"author_filter,omitempty"`
	CategoryFilter    string `toml:"category_filter,omitempty"`
}

// Subscription returns a Subscription.
func (f *Feed) Subscription(id string) (Subscription, error) {
	s := Subscription{
		ID:  FeedID(id),
		URI: f.URI,
	}

	if re, err := newFilter(id, "title", f.TitleFilter); err != nil {
		return s, err
	} else {
		s.TitleFilter = re
	}

	if re, err := newFilter(id, "description", f.DescriptionFilter); err != nil {
		return s, err
	} else {
		s.DescriptionFilter = re
	}

	if re, err := newFilter(id, "content", f.ContentFilter); err != nil {
		return s, err
	} else {
		s.ContentFilter = re
	}

	if re, err := newFilter(id, "author", f.AuthorFilter); err != nil {
		return s, err
	} else {
		s.AuthorFilter = re
	}

	if re, err := newFilter(id, "category", f.CategoryFilter); err != nil {
		return s, err
	} else {
		s.CategoryFilter = re
	}

	return s, nil
}

func newFilter(id, name, filter string) (*regexp.Regexp, error) {
	if filter == "" {
		return nil, nil
	}

	re, err := regexp.Compile(filter)
	if err != nil {
		return nil, fmt.Errorf("feed.%s: Invalid %s_filter: %s", id, name, err.Error())
	}

	return re, nil
}
