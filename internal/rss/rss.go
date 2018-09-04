package rss

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ShoshinNikita/log"
	"github.com/gorilla/feeds"
	"github.com/pkg/errors"
)

const rssFile = "rss.json"

var (
	feed  = new(feeds.Feed)
	mutex = new(sync.RWMutex)
)

func Init() error {
	f, err := os.Open(rssFile)
	if err != nil {
		// Need to create a file
		f, err := os.Create(rssFile)
		if err != nil {
			return errors.Wrapf(err, "can't create a new file %s", rssFile)
		}
		defer f.Close()
		feed = &feeds.Feed{
			Title:  "tg-to-rss-bot",
			Link:   &feeds.Link{Href: "github.com/ShoshinNikita/tg-to-rss-bot"},
			Author: &feeds.Author{Name: "ShoshinNikita"},
		}

		err = feed.WriteJSON(f)
		if err != nil {
			return errors.Wrap(err, "can't write RSS feed into file")
		}

		go runServer()

		return nil
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(feed)
	if err != nil {
		return errors.Wrap(err, "can't decode RSS feed from a file")
	}

	go runServer()

	return nil
}

func Add(author, title, description, link string, created time.Time) error {
	mutex.Lock()
	defer mutex.Unlock()

	feed.Add(&feeds.Item{
		Author:      &feeds.Author{Name: author},
		Created:     created,
		Description: description,
		Title:       title,
		Link:        &feeds.Link{Href: link},
	})

	// Backup
	f, err := os.OpenFile(rssFile, os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrapf(err, "can't open file %s", rssFile)
	}
	defer f.Close()
	err = feed.WriteJSON(f)
	if err != nil {
		return errors.Wrap(err, "can't write RSS feed into file")
	}

	return nil
}

func runServer() {
	http.HandleFunc("/rss", serveRSS)
	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("data/"))))
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalf("Server crashed: %s\n", err)
	}
}

func serveRSS(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	err := feed.WriteRss(w)
	if err != nil {
		log.Errorf("can't write RSS feed into http.ResponseWriter: %s", err)
	}
}
