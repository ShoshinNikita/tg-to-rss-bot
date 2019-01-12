package bot

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/ShoshinNikita/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ShoshinNikita/tg-to-rss-bot/internal/download"
	"github.com/ShoshinNikita/tg-to-rss-bot/internal/params"
)

func (b *Bot) start(msg *tgbotapi.Message) {
	b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "For start send link to a YouTube video"))
}

func (b *Bot) help(msg *tgbotapi.Message) {
	b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "It is a bot, that adds received videos from YouTube into RSS feed"))
}

func (b *Bot) wrongCommand(msg *tgbotapi.Message) {
	b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Error: wrong command"))
}

func (b *Bot) video(msg *tgbotapi.Message) {
	stringURL := msg.Text
	// Add http:// or https://
	if !(strings.HasPrefix(stringURL, "http://") || strings.HasPrefix(stringURL, "https://")) {
		stringURL = "https://" + stringURL
	}

	u, err := url.ParseRequestURI(stringURL)
	if err != nil {
		log.Errorf("can't parse url %s: %s\n", stringURL, err)
		b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Error: invalid link"))
		return
	}

	v, err := download.NewVideo(u)
	if err != nil {
		log.Errorf("can't get video info %s: %s\n", stringURL, err)
		b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Error: invalid link to a video"))
		return
	}

	msgText := fmt.Sprintf("Video: %s\n", v.Title)
	botMsg, err := b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, msgText))
	if err != nil {
		log.Errorf("can't send a message: %s", err)
		return
	}
	msgID := botMsg.MessageID

	var lastStatus interface{}

	for st := range v.Download() {
		lastStatus = st
		var status string
		switch st.(type) {
		case error:
			status = "❌ " + st.(error).Error()
		case string:
			status = st.(string)
		}

		msgText += "\n " + status
		b.bot.Send(tgbotapi.NewEditMessageText(msg.Chat.ID, msgID, msgText))
	}

	audioLink := params.DataFolder + v.Filename
	// We can define, was downloading success with lastStatus.(type): if type == string, downloading was success
	switch lastStatus.(type) {
	case string:
		// Add into feed
		err := b.feed.Add(v.Author, v.Title, v.Description, audioLink, time.Now())
		if err != nil {
			log.Errorf("can't add item into RSS feed: %s", err)
		} else {
			log.Infof("Add new item. Title: %s\n", v.Title)
		}
	}
}
