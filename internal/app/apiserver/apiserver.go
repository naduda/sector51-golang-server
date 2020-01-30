package apiserver

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/lib/pq" // ...
	"github.com/naduda/sector51-golang/internal/app/store/sqlstore"
)

// Start ...
func Start() error {
	connStr := os.Getenv("SECTOR_DB")
	bindAddr := os.Getenv("SECTOR_PORT")
	jwtSecret := os.Getenv("SECTOR_JWT")
	db, err := newDB(connStr)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	srv := newServer(store, jwtSecret)
	srv.logger.Infof("Server was started on port: %s", bindAddr)

	return http.ListenAndServe(":"+bindAddr, srv)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
