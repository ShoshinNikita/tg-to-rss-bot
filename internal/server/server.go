package server

import (
	"context"
	"net/http"
	"time"

	"github.com/ShoshinNikita/tg-to-rss-bot/cmd"
	"github.com/ShoshinNikita/tg-to-rss-bot/internal/params"
)

const address = ":80"

type Server struct {
	feed cmd.FeedInterface

	server *http.Server
}

func NewServer(feed cmd.FeedInterface) *Server {
	return &Server{feed: feed}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.index)
	mux.HandleFunc("/feed", s.serveFeed)
	// serve files
	mux.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("data/"))))

	s.server = &http.Server{Addr: address, Handler: mux}

	listenAndServe := s.server.ListenAndServe
	if params.TLS {
		listenAndServe = func() error {
			return s.server.ListenAndServeTLS("ssl/cert.cert", "ssl/key.key")
		}
	}

	if err := listenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	s.server.SetKeepAlivesEnabled(false)
	return s.server.Shutdown(ctx)
}
