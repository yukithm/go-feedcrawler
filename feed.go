package feedcrawler

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/naoina/toml"
)

// Feed is a feed configuration to be subscribed.
type Feed struct {
	ID                string `toml:"id"`
	URI               string `toml:"uri"`
	TitleFilter       string `toml:"title_filter,omitempty"`
	DescriptionFilter string `toml:"description_filter,omitempty"`
	ContentFilter     string `toml:"content_filter,omitempty"`
	AuthorFilter      string `toml:"author_filter,omitempty"`
	CategoryFilter    string `toml:"category_filter,omitempty"`
}

// Subscription returns a Subscription.
func (f *Feed) Subscription() (Subscription, error) {
	s := &SimpleSubscription{
		FeedID:  FeedID(f.ID),
		FeedURI: f.URI,
	}

	if re, err := newFilter(f.ID, "title", f.TitleFilter); err != nil {
		return s, err
	} else {
		s.TitleFilter = re
	}

	if re, err := newFilter(f.ID, "description", f.DescriptionFilter); err != nil {
		return s, err
	} else {
		s.DescriptionFilter = re
	}

	if re, err := newFilter(f.ID, "content", f.ContentFilter); err != nil {
		return s, err
	} else {
		s.ContentFilter = re
	}

	if re, err := newFilter(f.ID, "author", f.AuthorFilter); err != nil {
		return s, err
	} else {
		s.AuthorFilter = re
	}

	if re, err := newFilter(f.ID, "category", f.CategoryFilter); err != nil {
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

// LoadFeeds loads feeds from io.Reader.
func LoadFeeds(r io.Reader) ([]Feed, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if len(buf) == 0 {
		return make([]Feed, 0), nil
	}

	cfg := struct {
		Feeds []Feed `toml:"feeds"`
	}{}

	if err := toml.Unmarshal(buf, &cfg); err != nil {
		return nil, err
	}

	return cfg.Feeds, nil
}

// LoadFeedsFile loads feeds from file.
func LoadFeedsFile(file string) ([]Feed, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return LoadFeeds(f)
}
