package apiserver

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/naduda/sector51-golang/internal/app/store"
)

const (
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectPhoneOrPassword = errors.New("incorrect phone or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type Server struct {
	router    *mux.Router
	logger    *logrus.Logger
	store     store.Store
	jwtSecret string
}

func newServer(store store.Store, jwtSecret string, logger *logrus.Logger) *Server {
	logger.SetOutput(os.Stdout)
	logger.Debug("Start newServer")
	s := &Server{
		router:    mux.NewRouter(),
		logger:    logger,
		store:     store,
		jwtSecret: jwtSecret,
	}

	logger.Debug("End newServer")
	s.configureRouter()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}
		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}
