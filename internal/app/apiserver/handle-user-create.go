package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
	"github.com/naduda/sector51-golang/internal/app/model"
)

type request struct {
	ID       string `json:"id"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Card     string `json:"card"`
	Roles    string `json:"roles"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	IsMan    bool   `json:"isMan"`
}

func requestToUser(req *request) *model.User {
	return &model.User{
		ID:       req.ID,
		Phone:    req.Phone,
		Password: req.Password,
		Card:     req.Card,
		Roles:    req.Roles,
		Name:     req.Name,
		Surname:  req.Surname,
		IsMan:    req.IsMan,
	}
}

func (s *Server) handleUsersCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		u := requestToUser(req)
		if err := s.store.User().Create(u); err != nil {
			httputils.SendError(w, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		httputils.Respond(w, http.StatusCreated, u)
	}
}

func (s *Server) handleUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		u := requestToUser(req)
		if err := s.store.User().Update(*u); err != nil {
			httputils.SendError(w, http.StatusUnprocessableEntity, err)
			return
		}

		httputils.Respond(w, http.StatusCreated, u)
	}
}
