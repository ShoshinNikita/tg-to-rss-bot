package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ShoshinNikita/log"

	"github.com/ShoshinNikita/tg-to-rss-bot/cmd"
	"github.com/ShoshinNikita/tg-to-rss-bot/internal/bot"
	"github.com/ShoshinNikita/tg-to-rss-bot/internal/feed"
	"github.com/ShoshinNikita/tg-to-rss-bot/internal/server"
)

type App struct {
	Server cmd.ServerInterface
	Bot    cmd.BotInterface
	Feed   cmd.FeedInterface
}

func main() {
	log.Infoln("Start")

	app := &App{}

	// Feed
	app.Feed = feed.NewFeed()
	err := app.Feed.Init()
	if err != nil {
		log.Fatalf("can't init feed: %s\n", err)
	}

	// Bot
	app.Bot = bot.NewBot(app.Feed)
	err = app.Bot.Start()
	if err != nil {
		log.Fatalf("can't init feed: %s\n", err)
	}

	// Server
	app.Server = server.NewServer(app.Feed)

	done := make(chan struct{})
	serverError := make(chan struct{})
	go func() {
		term := make(chan os.Signal)
		signal.Notify(term, syscall.SIGTERM, syscall.SIGINT)

		select {
		case <-term:
			log.Warnln("catch interrupt signal")
		case <-serverError:
			// Nothing
		}

		log.Infoln("shutdown bot")
		err := app.Bot.Shutdown()
		if err != nil {
			log.Errorf("can't shutdown bot gracefully: %s\n", err)
		}
		log.Infoln("shutdown server")
		err = app.Server.Shutdown()
		if err != nil {
			log.Errorf("can't shutdown server gracefully: %s\n", err)
		}

		close(done)
	}()

	err = app.Server.Start()
	if err != nil {
		log.Errorf("server error: %s\n", err)
		close(serverError)
	}

	<-done

	log.Infoln("Stop")
}
