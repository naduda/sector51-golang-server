package apiserver

import (
	"context"
	"net/http"

	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
)

func (s *Server) authUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			httputils.SendError(w, http.StatusForbidden, errNotAuthenticated)
			return
		}

		claims, err := s.parseJWT(token)
		if err != nil {
			httputils.SendError(w, http.StatusForbidden, err)
			return
		}

		userID := claims["id"].(string)

		u, err := s.store.User().Find(userID)
		if err != nil {
			httputils.SendError(w, http.StatusForbidden, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}
