package cmd

import (
	"io"
	"time"
)

type ServerInterface interface {
	Start() error

	Shutdown() error
}

type BotInterface interface {
	Start() error

	Shutdown() error
}

type FeedInterface interface {
	Init() error

	Add(author, title, description, link string, created time.Time) error

	Write(w io.Writer) error
}
