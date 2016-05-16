package main

import (
	"log"
	"os"
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

	subscriptions, err := config.Feed.Subscriptions()
	if err != nil {
		return crawler, err
	}
	crawler.Subscriptions = subscriptions

	return crawler, nil
}
