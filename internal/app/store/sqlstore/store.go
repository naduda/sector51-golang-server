package sqlstore

import (
	"github.com/jmoiron/sqlx"

	"github.com/naduda/sector51-golang/internal/app/store"
)

// store ...
type Store struct {
	db                *sqlx.DB
	userRepository    *UserRepository
	serviceRepository *ServiceRepository
	boxRepository     *BoxRepository
}

// DB ...
func (s *Store) DB() *sqlx.DB {
	return s.db
}

// New ...
func New(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

// User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

// Service ...
func (s *Store) Service() store.ServiceRepository {
	if s.serviceRepository != nil {
		return s.serviceRepository
	}

	s.serviceRepository = &ServiceRepository{
		store: s,
	}

	return s.serviceRepository
}

// Boxes ...
func (s *Store) Boxes() store.BoxRepository {
	if s.boxRepository != nil {
		return s.boxRepository
	}

	s.boxRepository = &BoxRepository{
		store: s,
	}

	return s.boxRepository
}
