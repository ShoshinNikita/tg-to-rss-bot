package server

import "net/http"

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	s.serveFeed(w, r)
}

func (s *Server) serveFeed(w http.ResponseWriter, r *http.Request) {
	s.feed.Write(w)
}
