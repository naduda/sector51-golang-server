package apiserver

import (
	"encoding/json"
	"github.com/naduda/sector51-golang/internal/app/store"
	"net/http"

	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
	"github.com/naduda/sector51-golang/internal/app/model"
)

func (s *Server) handleUsersCreate(isFirst bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isFirst {
			all, err := s.store.User().FindAll()
			if err != nil {
				httputils.SendError(w, http.StatusUnprocessableEntity, err)
				return
			}
			if len(all) > 0 {
				httputils.SendError(w, http.StatusUnprocessableEntity, store.ErrReject)
				return
			}
		}

		u := &model.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		//if err := s.store.User().Create(u); err != nil {
		//	httputils.SendError(w, http.StatusUnprocessableEntity, err)
		//	return
		//}

		u.Sanitize()
		httputils.Respond(w, http.StatusCreated, u)
	}
}

func (s *Server) handleUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &model.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		if err := s.store.User().Update(*u); err != nil {
			httputils.SendError(w, http.StatusUnprocessableEntity, err)
			return
		}

		httputils.Respond(w, http.StatusCreated, u)
	}
}
