package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
	"github.com/naduda/sector51-golang/internal/app/model"
)

func (s *Server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Card     string `json:"card"`
		Roles    string `json:"roles"`
		IsMan    bool   `json:"isMan"`
	}
	//{"card":"1100000001102","isMan":true,"password":"secret","phone":"+380505555555","roles":"USER"}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Phone:    req.Phone,
			Password: req.Password,
			Card:     req.Card,
			Roles:    req.Roles,
			IsMan:    req.IsMan,
		}
		if err := s.store.User().Create(u); err != nil {
			httputils.SendError(w, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		httputils.Respond(w, http.StatusCreated, u)
	}
}
