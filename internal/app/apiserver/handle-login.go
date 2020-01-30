package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
)

func (s *Server) handleLogin() http.HandlerFunc {
	type request struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByPhone(req.Phone)
		if err != nil || !u.ComparePassword(req.Password) {
			httputils.SendError(w, http.StatusUnauthorized, errIncorrectPhoneOrPassword)
			return
		}

		token, err := s.generateJWT(u)
		if err != nil {
			httputils.SendError(w, http.StatusInternalServerError, err)
			return
		}

		httputils.Respond(w, http.StatusOK, map[string]string{"token": token})
	}
}
