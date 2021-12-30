package apiserver

import (
	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
	"net/http"
	"os"
)

func (s *Server) handleBreak() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		os.Exit(0)
	}
}

func (s Server) handleFixPhones() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.store.User().FixPhones(); err != nil {
			httputils.SendError(w, http.StatusInternalServerError, err)
			return
		}
		httputils.Respond(w, http.StatusOK, nil)
	}
}
