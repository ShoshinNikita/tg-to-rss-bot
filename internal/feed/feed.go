package feed

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"

	"github.com/gorilla/feeds"
	"github.com/pkg/errors"

	"github.com/ShoshinNikita/tg-to-rss-bot/internal/params"
)

const (
	rssFolder = "rss"
	rssFile   = "rss/rss.json"
)

type Feed struct {
	feed  *feeds.Feed
	mutex *sync.RWMutex
}

func NewFeed() *Feed {
	return &Feed{
		feed:  new(feeds.Feed),
		mutex: new(sync.RWMutex),
	}
}

func (feed *Feed) Init() error {
	err := os.MkdirAll(rssFolder, 0666)
	if err != nil {
		return errors.Wrapf(err, "can't create folder %s", rssFolder)
	}

	f, err := os.Open(rssFile)
	if err != nil {
		// Need to create a file
		f, err = os.Create(rssFile)
		if err != nil {
			return errors.Wrapf(err, "can't create a new file %s", rssFile)
		}
		defer f.Close()

		feed.feed = &feeds.Feed{
			Title:  "tg-to-rss-bot",
			Link:   &feeds.Link{Href: "github.com/ShoshinNikita/tg-to-rss-bot"},
			Author: &feeds.Author{Name: "ShoshinNikita"},
		}

		err = feed.feed.WriteJSON(f)
		if err != nil {
			return errors.Wrap(err, "can't write RSS feed into file")
		}

		return nil
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(feed.feed)
	if err != nil {
		return errors.Wrap(err, "can't decode RSS feed from a file")
	}

	return nil
}

func (feed *Feed) Add(author, title, description, link string, created time.Time) error {
	feed.mutex.Lock()
	defer feed.mutex.Unlock()

	feed.feed.Add(&feeds.Item{
		Author:      &feeds.Author{Name: author},
		Created:     created,
		Description: description,
		Title:       title,
		Link:        &feeds.Link{Href: params.Host + "/" + link},
	})

	// Write into disk
	f, err := os.OpenFile(rssFile, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return errors.Wrapf(err, "can't open file %s", rssFile)
	}
	defer f.Close()

	err = feed.feed.WriteJSON(f)
	if err != nil {
		return errors.Wrap(err, "can't write RSS feed into file")
	}

	return nil
}

func (feed *Feed) Write(w io.Writer) error {
	feed.mutex.RLock()
	defer feed.mutex.RUnlock()

	err := feed.feed.WriteRss(w)
	if err != nil {
		errors.Wrap(err, "can't write RSS feed into io.Writer")
	}

	return nil
}
