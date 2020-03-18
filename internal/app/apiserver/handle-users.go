package apiserver

import (
	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
	"net/http"
)

func (s *Server) handleUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := s.store.User().FindAll()
		if err != nil {
			httputils.SendError(w, http.StatusInternalServerError, err)
			return
		}
		httputils.Respond(w, http.StatusOK, res)
	}
}
