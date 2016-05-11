package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"text/template"

	"github.com/yukithm/go-feedcrawler"
)

const configFile = "feedcrawler.toml"

var outputTemplate = `[{{.Feed.Title}} ({{.Subscription.URI}})]
{{range .NewItems}}{{.Title}} ({{.Link}})
{{end}}
`

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	crawler, err := newFeedCrawler(config)
	if err != nil {
		log.Fatal(err)
	}
	err = crawler.CrawlFunc(func(result feedcrawler.Result) {
		if result.Err != nil {
			log.Print(result.Err)
		} else {
			tmpl := template.Must(template.New("output").Parse(outputTemplate))
			tmpl.Execute(os.Stdout, result)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}

func newFeedCrawler(config Config) (feedcrawler.Crawler, error) {
	crawler := feedcrawler.Crawler{
		StateFile:  config.FeedCrawler.StateFile,
		NumWorkers: config.FeedCrawler.NumWorkers,
	}

	subscriptions, err := newSubscriptions(config)
	if err != nil {
		return crawler, err
	}
	crawler.Subscriptions = subscriptions

	return crawler, nil
}

func newSubscriptions(config Config) ([]feedcrawler.Subscription, error) {
	var subscriptions []feedcrawler.Subscription

	for id, feed := range config.Feed {
		s, err := newSubscription(id, feed)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, s)
	}

	return subscriptions, nil
}

func newSubscription(id string, feed Feed) (feedcrawler.Subscription, error) {
	s := feedcrawler.Subscription{
		ID:  feedcrawler.FeedID(id),
		URI: feed.URI,
	}

	if re, err := newFilter(id, "title", feed.TitleFilter); err != nil {
		return s, err
	} else {
		s.TitleFilter = re
	}

	if re, err := newFilter(id, "description", feed.DescriptionFilter); err != nil {
		return s, err
	} else {
		s.DescriptionFilter = re
	}

	if re, err := newFilter(id, "content", feed.ContentFilter); err != nil {
		return s, err
	} else {
		s.ContentFilter = re
	}

	if re, err := newFilter(id, "author", feed.AuthorFilter); err != nil {
		return s, err
	} else {
		s.AuthorFilter = re
	}

	if re, err := newFilter(id, "category", feed.CategoryFilter); err != nil {
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
