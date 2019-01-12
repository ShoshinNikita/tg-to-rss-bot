package cmd

import (
	"io"
	"time"
)

type ServerInterface interface {
	// Start starts server. It doesn't finish until server is shutdowned
	Start() error

	// Shutdown shutdowns server. It unlocks Start()
	Shutdown() error
}

type BotInterface interface {
	// Start starts bot
	Start() error

	// Shutdown gracefully shutdowns bot
	Shutdown() error
}

type FeedInterface interface {
	// Init inits feed (create files and directories)
	Init() error

	// Add adds new item into feed
	Add(author, title, description, link string, created time.Time) error

	// Write writes RSS feed into w
	Write(w io.Writer) error
}
