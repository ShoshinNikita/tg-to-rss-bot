package params

import (
	"os"

	"github.com/ShoshinNikita/log"
)

const (
	DataFolder = "data/"
	RssFolder  = "rss/"
	RssFile    = "rss/rss.xml"
)

var (
	// BotToken is a token for Telegram bot
	BotToken string
	// Host serves for creating link to audio (for example, "http://1.1.1.1")
	Host string
)

func init() {
	BotToken = os.Getenv("TOKEN")
	Host = os.Getenv("HOST")
	if Host == "" {
		log.Fatal("HOST can't be empty")
	}
}
