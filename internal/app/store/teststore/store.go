package teststore

import (
	"github.com/jmoiron/sqlx"
	"github.com/naduda/sector51-golang/internal/app/model"
	"github.com/naduda/sector51-golang/internal/app/store"
)

// store ...
type Store struct {
	userRepository    *UserRepository
	serviceRepository *ServiceRepository
	boxRepository     *BoxRepository
}

// test store doesn't have db
func (s *Store) DB() *sqlx.DB {
	return nil
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

// Boxes ...
func (s *Store) Boxes() store.BoxRepository {
	if s.boxRepository != nil {
		return s.boxRepository
	}

	boxesForTest := make(map[int]model.Box)

	s.boxRepository = &BoxRepository{
		store: s,
		boxes: boxesForTest,
	}

	for i := 0; i < 50; i++ {
		for t := 0; t < 3; t++ {
			boxesForTest[t*50+i] = model.Box{
				Card:   "",
				IdType: t + 1,
				Value:  i + 1,
				Time:   "",
			}
		}
	}

	return s.boxRepository
}
