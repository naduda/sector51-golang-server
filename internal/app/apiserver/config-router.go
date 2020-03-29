package apiserver

import (
	"github.com/naduda/sector51-golang/internal/app/apiserver/clients"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/naduda/sector51-golang/internal/app/backup/backuphandlers"
)

func (s *Server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/login", s.handleLogin()).Methods("POST")
	s.router.HandleFunc("/break", s.handleBreak()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authUser)
	private.HandleFunc("/has-google-credentials", backuphandlers.HandleHasGoogleCredentials())
	private.HandleFunc("/upload", s.handleUploadFile()).Methods("POST")
	private.HandleFunc("/create-google-token", backuphandlers.HandleCreateGoogleTokenFile()).Methods("POST")
	private.HandleFunc("/backup", backuphandlers.HandleBackup(s.logger)).Methods("POST")
	private.HandleFunc("/restore", backuphandlers.HandleRestore(s.logger)).Methods("POST")
	// clients
	private.HandleFunc("/whoami", s.handleWhoami())
	private.HandleFunc("/clients-list", s.handleUsers())
	private.HandleFunc("/services", clients.HandleServices(s.store.Service()))

	fs := http.Dir("static")
	fileHandler := http.FileServer(fs)
	s.router.PathPrefix("/").Handler(fileHandler)
	s.router.PathPrefix("/{_:.*}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
}
