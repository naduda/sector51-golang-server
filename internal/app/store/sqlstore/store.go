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
