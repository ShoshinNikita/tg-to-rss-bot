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

const helpMsg = `It is a bot, that adds received videos from YouTube into RSS feed
Commands:
/help – get help
/link – get link to RSS feed

Bot repo: https://github.com/ShoshinNikita/tg-to-rss-bot`

func (b *Bot) start(msg *tgbotapi.Message) {
	b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "For start send link to a YouTube video"))
}

func (b *Bot) help(msg *tgbotapi.Message) {
	b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, helpMsg))
}

func (b *Bot) sendLink(msg *tgbotapi.Message) {
	b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, params.Host+"/feed"))
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
		b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ invalid link"))
		return
	}

	v, err := download.NewVideo(u)
	if err != nil {
		b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ invalid link to a video"))
		return
	}

	msgText := fmt.Sprintf("Video: %s\n", v.Title)
	botMsg, err := b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, msgText))
	// Check error because we will use botMsg.MessageID
	if err != nil {
		log.Errorf("can't send a message: %s\n", err)
		return
	}

	msgID := botMsg.MessageID
	status := ""
	ok := false
	for m := range v.Download() {
		switch {
		case m.IsFatalError:
			ok = false
			status = "❌ " + m.Msg
		case m.IsFinished:
			ok = true
			status = "✅ " + m.Msg
		default:
			status = "* " + m.Msg
		}

		msgText += "\n " + status
		b.bot.Send(tgbotapi.NewEditMessageText(msg.Chat.ID, msgID, msgText))

		if m.IsFatalError || m.IsFinished {
			break
		}
	}

	if ok {
		// Add into feed
		audioLink := params.DataFolder + v.Filename
		err := b.feed.Add(v.Author, v.Title, v.Description, audioLink, time.Now())
		if err != nil {
			log.Errorf("can't add item into RSS feed: %s\n", err)
		} else {
			log.Infof("add new item. Title: \"%s\"\n", v.Title)
		}
	}
}
