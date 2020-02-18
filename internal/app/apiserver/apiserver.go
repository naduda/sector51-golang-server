package apiserver

import (
	"database/sql"
	"github.com/naduda/sector51-golang/internal/app/backup"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq" // ...
	"github.com/naduda/sector51-golang/internal/app/store/sqlstore"
)

// Start ...
func Start() error {
	connStr := os.Getenv("CONNECTION_STR")
	bindAddr := os.Getenv("WEB_PORT")
	jwtSecret := os.Getenv("JWT_SECRET")
	logger := logrus.New()
	db := newDB(connStr, logger)
	defer db.Close()

	b := backup.New(connStr, logger)
	go b.Start()

	store := sqlstore.New(db)
	srv := newServer(store, jwtSecret, logger)
	srv.logger.Infof("Server was started on port: %s", bindAddr)

	return http.ListenAndServe(":"+bindAddr, srv)
}

func newDB(connStr string, log *logrus.Logger) *sql.DB {
	for {
		db, err := sql.Open("postgres", connStr)
		if err == nil {
			if err := db.Ping(); err == nil {
				return db
			}
		}
		log.Info("Waiting for DB connection...")
		time.Sleep(5 * time.Second)
	}
}
