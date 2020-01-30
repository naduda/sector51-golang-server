package apiserver

import (
	"net/http"

	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
	"github.com/naduda/sector51-golang/internal/app/model"
)

func (s *Server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httputils.Respond(w, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}
