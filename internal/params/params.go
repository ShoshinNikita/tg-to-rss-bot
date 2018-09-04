package params

import (
	"os"
)

var (
	// BotToken is a token for Telegram bot
	BotToken string
)

func init() {
	BotToken = os.Getenv("TOKEN")
}
