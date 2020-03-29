package teststore

import (
	"github.com/naduda/sector51-golang/internal/app/model"
	"github.com/naduda/sector51-golang/internal/app/store"
)

// store ...
type Store struct {
	userRepository    *UserRepository
	serviceRepository *ServiceRepository
}

// New ...
func New() *Store {
	return &Store{}
}

// User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[string]model.User),
	}

	return s.userRepository
}

// Service ...
func (s *Store) Service() store.ServiceRepository {
	if s.serviceRepository != nil {
		return s.serviceRepository
	}

	servicesForTest := make(map[int]model.Service)
	servicesForTest[0] = model.Service{
		ID:    0,
		Name:  "ABONEMENT",
		Desc:  "-",
		Price: 750,
	}
	servicesForTest[13] = model.Service{
		ID:    13,
		Name:  "ABONEMENT (Evening) 12",
		Desc:  "-",
		Price: 4600,
	}
	s.serviceRepository = &ServiceRepository{
		store:    s,
		services: servicesForTest,
	}

	return s.serviceRepository
}
