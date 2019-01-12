package cmd

import "io"

type ServerInterface interface {
	Start() error

	Shutdown() error
}

type BotInterface interface {
	Start() error

	Shutdown() error
}

type FeedInterface interface {
	Add() error

	Write(w io.Writer)
}
