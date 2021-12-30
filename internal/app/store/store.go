package store

import (
	"github.com/jmoiron/sqlx"
)

// store ...
type Store interface {
	DB() *sqlx.DB
	User() UserRepository
	Service() ServiceRepository
	Boxes() BoxRepository
}
