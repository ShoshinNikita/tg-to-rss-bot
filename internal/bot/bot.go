package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ShoshinNikita/tg-to-rss-bot/internal/params"
)

var (
	bot *tgbotapi.BotAPI

	commandHandlers = []struct {
		command string
		handler func(*tgbotapi.Message)
	}{
		{"start", start},
		{"help", help},
		{"", video},
	}
)

func Init() (err error) {
	bot, err = tgbotapi.NewBotAPI(params.BotToken)
	if err != nil {
		return err
	}

	update := tgbotapi.NewUpdate(0)
	update.Timeout = 60

	updatesChan, err := bot.GetUpdatesChan(update)
	if err != nil {
		return err
	}

	go listenAndServe(updatesChan)

	return nil
}

func listenAndServe(updatesChan tgbotapi.UpdatesChannel) {
	for update := range updatesChan {
		if update.Message != nil {
			go serve(update.Message)
		}
	}
}

func serve(msg *tgbotapi.Message) {
	cmd := msg.Command()
	for _, hand := range commandHandlers {
		if hand.command == cmd {
			hand.handler(msg)
			return
		}
	}

	wrongCommand(msg)
}
