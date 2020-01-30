package apiserver

import (
	"net/http"
	"os"
)

func (s *Server) handleBreak() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		os.Exit(0)
	}
}
