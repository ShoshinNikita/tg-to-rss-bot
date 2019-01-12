package bot

import (
	"github.com/ShoshinNikita/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"

	"github.com/ShoshinNikita/tg-to-rss-bot/internal/params"
)

type Bot struct {
	bot *tgbotapi.BotAPI

	shutdownReq  chan struct{}
	shutdownResp chan struct{}
}

func NewBot() *Bot {
	return &Bot{}
}

func (b *Bot) Start() (err error) {
	b.bot, err = tgbotapi.NewBotAPI(params.BotToken)
	if err != nil {
		return errors.Wrap(err, "can't init bot")
	}

	update := tgbotapi.NewUpdate(0)
	update.Timeout = 60

	updatesChan, err := b.bot.GetUpdatesChan(update)
	if err != nil {
		return err
	}

	go b.listenAndServe(updatesChan)

	return nil
}

func (b *Bot) listenAndServe(updatesChan tgbotapi.UpdatesChannel) {
	for {
		select {
		case update := <-updatesChan:
			if update.Message != nil {
				go b.serve(update.Message)
			}
		case <-b.shutdownReq:
			close(b.shutdownResp)
			return
		}
	}
}

func (b *Bot) serve(msg *tgbotapi.Message) {
	log.Infof("User: %s ID: %d Text: %s\n", msg.Chat.UserName, msg.Chat.ID, msg.Text)

	cmd := msg.Command()
	switch cmd {
	case "start":
		b.start(msg)
	case "help":
		b.help(msg)
	case "":
		b.video(msg)
	default:
		b.wrongCommand(msg)
	}
}

func (b *Bot) Shutdown() error {
	close(b.shutdownReq)
	<-b.shutdownResp

	return nil
}
