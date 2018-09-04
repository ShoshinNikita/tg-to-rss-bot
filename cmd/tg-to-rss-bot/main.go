package main

import (
	"github.com/ShoshinNikita/log"

	"github.com/ShoshinNikita/tg-to-rss-bot/internal/bot"
	"github.com/ShoshinNikita/tg-to-rss-bot/internal/rss"
)

func main() {
	log.Infoln("Start")

	err := bot.Init()
	if err != nil {
		log.Fatalln(err)
	}
	err = rss.Init()
	if err != nil {
		log.Fatalln(err)
	}

	c := make(chan struct{})
	<-c
}
