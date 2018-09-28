package feedcrawler

import (
	"errors"
	"net/url"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/mmcdole/gofeed/atom"
)

type EnhancedAtomTranslator struct {
	defaultTranslator *gofeed.DefaultAtomTranslator
}

func NewEnhancedAtomTranslator() *EnhancedAtomTranslator {
	return &EnhancedAtomTranslator{
		defaultTranslator: &gofeed.DefaultAtomTranslator{},
	}
}

func (ct *EnhancedAtomTranslator) Translate(feed interface{}) (*gofeed.Feed, error) {
	af, ok := feed.(*atom.Feed)
	if !ok {
		return nil, errors.New("Feed did not match expected type of *atom.Feed")
	}

	f, err := ct.defaultTranslator.Translate(feed)
	if err != nil {
		return nil, err
	}

	ct.fillUpLink(f, af)

	return f, nil
}

func (ct *EnhancedAtomTranslator) fillUpLink(f *gofeed.Feed, af *atom.Feed) {
	if f.Link == "" {
		f.Link = extractFeedLink(af)
	}

	for _, item := range f.Items {
		if item.Link == "" {
			entry := findEntry(af.Entries, item.GUID)
			if entry != nil {
				item.Link = extractEntryLink(entry)
			}
		}
	}
}

func findEntry(entries []*atom.Entry, id string) *atom.Entry {
	for _, entry := range entries {
		if entry.ID == id {
			return entry
		}
	}
	return nil
}

func extractFeedLink(af *atom.Feed) string {
	var link string
	if af.Links != nil {
		link = extractLink(af.Links, "alternate")
	}

	if link == "" && isURL(af.ID) {
		link = af.ID
	}

	return link
}

func extractEntryLink(entry *atom.Entry) string {
	var link string
	if entry.Links != nil {
		link = extractLink(entry.Links, "alternate")
	}

	if link == "" && isURL(entry.ID) {
		link = entry.ID
	}

	return link
}

func extractLink(links []*atom.Link, rel string) string {
	if links == nil || len(links) == 0 {
		return ""
	}

	for _, link := range links {
		if rel == "" || link.Rel == rel {
			return link.Href
		}
	}

	return ""
}

func isURL(value string) bool {
	u, err := url.Parse(value)
	if err != nil {
		return false
	}
	return u.IsAbs() || strings.HasPrefix(value, "//")
}
