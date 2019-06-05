package main

import (
	"log"
	"os"
	"text/template"

	"github.com/yukithm/go-feedcrawler"
)

const configFileName = "feedcrawler.toml"

var outputTemplate = `[{{.Feed.Title}} ({{.Subscription.URI}})]
{{range .NewItems}}{{.Title}} ({{.Link}})
{{end}}
`

func main() {
	cfg, err := loadConfig(configFileName)
	if err != nil {
		log.Fatal(err)
	}

	crawler, err := newFeedCrawler(cfg)
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
	feedcrawler.SaveStatesFile(crawler.States, cfg.StatesFile)
}

func newFeedCrawler(cfg *Config) (*feedcrawler.Crawler, error) {
	crawler := &feedcrawler.Crawler{
		NumWorkers: cfg.NumWorkers,
	}

	feeds, err := feedcrawler.LoadFeedsFile(cfg.FeedsFile)
	if err != nil {
		return nil, err
	}

	subscriptions, err := feedsToSubscriptions(feeds)
	if err != nil {
		return nil, err
	}
	crawler.Subscriptions = subscriptions

	if cfg.StatesFile != "" {
		states, err := feedcrawler.LoadStatesFile(cfg.StatesFile)
		if err != nil {
			return nil, err
		}
		crawler.States = states
	}

	return crawler, nil
}

func feedsToSubscriptions(feeds []feedcrawler.Feed) ([]feedcrawler.Subscription, error) {
	subscriptions := make([]feedcrawler.Subscription, 0, len(feeds))

	for _, feed := range feeds {
		s, err := feed.Subscription()
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, s)
	}

	return subscriptions, nil
}
