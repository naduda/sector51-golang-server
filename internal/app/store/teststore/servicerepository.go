package teststore

import (
	"github.com/naduda/sector51-golang/internal/app/model"
	"github.com/naduda/sector51-golang/internal/app/store"
)

// ServiceRepository ...
type ServiceRepository struct {
	store        *Store
	services     map[int]model.Service
	userServices map[string][]model.UserService
}

// Services ...
func (r *ServiceRepository) List() ([]model.Service, error) {
	res := make([]model.Service, 0)
	for _, u := range r.services {
		res = append(res, u)
	}
	return res, nil
}

// UpdateService ...
func (r *ServiceRepository) UpdateService(s model.Service) error {
	for i, v := range r.services {
		if s.ID == v.ID {
			r.services[i] = s
			return nil
		}
	}
	return store.ErrRecordNotFound
}

// Find ...
func (r *ServiceRepository) Find(id int) (*model.Service, error) {
	for _, v := range r.services {
		if v.ID == id {
			return &v, nil
		}
	}
	return nil, store.ErrRecordNotFound
}

// GetUserServices ...
func (r *ServiceRepository) GetUserServices(idUser string) ([]model.UserService, error) {
	res, ok := r.userServices[idUser]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return res, nil
}

// CreateUserService ...
func (r *ServiceRepository) CreateUserService(us *model.UserService) error {
	if err := us.Validate(); err != nil {
		return err
	}

	if r.userServices == nil {
		r.userServices = make(map[string][]model.UserService)
	}
	exists, ok := r.userServices[us.IdUser]
	if !ok {
		exists = []model.UserService{}
	}

	exists = append(exists, *us)
	r.userServices[us.IdUser] = exists

	return nil
}

// UpdateUserService ...
func (r *ServiceRepository) UpdateUserService(us model.UserService) error {
	if err := us.Validate(); err != nil {
		return err
	}

	exists, ok := r.userServices[us.IdUser]
	if !ok {
		return store.ErrRecordNotFound
	}

	exists = append(exists, us)
	r.userServices[us.IdUser] = exists

	return nil
}

// DeleteUserService ...
func (r *ServiceRepository) DeleteUserService(idUser string, idService int) error {
	if _, ok := r.services[idService]; !ok {
		return store.ErrRecordNotFound
	}

	if res, ok := r.userServices[idUser]; !ok {
		return store.ErrRecordNotFound
	} else {
		for i, v := range res {
			if v.IdService == idService {
				res = append(res[:i], res[i+1:]...)
				break
			}
		}
		r.userServices[idUser] = res
	}
	return nil
}
