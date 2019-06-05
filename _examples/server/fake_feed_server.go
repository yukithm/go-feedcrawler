package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/feeds"
)

type FakeFeedServer struct {
	http.Server
}

func (s *FakeFeedServer) ListenAndServe() error {
	if s.Server.Addr == "" {
		s.Server.Addr = ":8080"
	}
	s.Server.ReadTimeout = 30 * time.Second

	mux := http.NewServeMux()
	mux.HandleFunc("/foo.rss", fooFeedFunc)
	mux.HandleFunc("/foo.atom", fooFeedFunc)
	s.Handler = mux

	return s.Server.ListenAndServe()
}

var fooFeed = feeds.Feed{
	Title:       "Foo feed",
	Link:        &feeds.Link{Href: "http://localhost/foo"},
	Description: "Fake feed for /foo.rss",
	Author:      &feeds.Author{Name: "yukithm", Email: "yukithm@example.com"},
	Created:     time.Date(2016, time.April, 12, 0, 0, 0, 0, time.UTC),
	Items: []*feeds.Item{
		{
			Title:       "Version 1.3.0 released!",
			Link:        &feeds.Link{Href: "http://localhost/foo/v1.3.0"},
			Description: "Version 1.3.0 released!",
			Author:      &feeds.Author{Name: "yukithm", Email: "yukithm@example.com"},
			Created:     time.Date(2016, time.April, 12, 0, 0, 0, 0, time.UTC),
		},
		{
			Title:       "Version 1.2.6 released!",
			Link:        &feeds.Link{Href: "http://localhost/foo/v1.2.6"},
			Description: "Version 1.2.6 released!",
			Author:      &feeds.Author{Name: "yukithm", Email: "yukithm@example.com"},
			Created:     time.Date(2016, time.April, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			Title:       "Version 1.2.5 released!",
			Link:        &feeds.Link{Href: "http://localhost/foo/v1.2.5"},
			Description: "Version 1.2.5 released!",
			Author:      &feeds.Author{Name: "yukithm", Email: "yukithm@example.com"},
			Created:     time.Date(2016, time.April, 2, 0, 0, 0, 0, time.UTC),
		},
	},
}

var fooCount = 0
var lastUpdate = time.Now()

func fooFeedFunc(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	d := now.Sub(lastUpdate)
	if d.Seconds() > 5.0 {
		fooCount++
		updateFooFeed(now, fooCount)
		lastUpdate = now
	}

	if path.Ext(r.URL.Path) == ".atom" {
		w.Header().Set("Content-Type", "application/atom+xml")
		fooFeed.WriteAtom(w)
	} else {
		w.Header().Set("Content-Type", "application/rss+xml")
		fooFeed.WriteRss(w)
	}
}

func updateFooFeed(now time.Time, n int) {
	version := fmt.Sprintf("1.3.%d", n)
	newItems := []*feeds.Item{
		{
			Title:       fmt.Sprintf("About v%s", version),
			Link:        &feeds.Link{Href: fmt.Sprintf("http://localhost/foo/v%s", version)},
			Description: fmt.Sprintf("About v%s", version),
			Author:      &feeds.Author{Name: "yukithm", Email: "yukithm@example.com"},
			Created:     now,
		},
		{
			Title:       fmt.Sprintf("Version %s released!", version),
			Link:        &feeds.Link{Href: fmt.Sprintf("http://localhost/foo/v%s", version)},
			Description: fmt.Sprintf("Version %s released!", version),
			Author:      &feeds.Author{Name: "yukithm", Email: "yukithm@example.com"},
			Created:     now,
		},
	}
	fooFeed.Items = append(newItems, fooFeed.Items...)
	fooFeed.Created = now
}

func main() {
	server := FakeFeedServer{}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
