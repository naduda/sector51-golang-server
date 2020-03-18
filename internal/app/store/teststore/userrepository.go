package teststore

import (
	"github.com/naduda/sector51-golang/internal/app/model"
	"github.com/naduda/sector51-golang/internal/app/store"
	"time"
)

// UserRepository ...
type UserRepository struct {
	store *Store
	users map[string]*model.User
}

// Create ...
func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	u.ID = time.Now().String()
	r.users[u.ID] = u

	return nil
}

// Find ...
func (r *UserRepository) Find(id string) (*model.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}

// FindAll ...
func (r *UserRepository) FindAll() ([]*model.User, error) {
	res := make([]*model.User, 0)
	for _, u := range r.users {
		res = append(res, u)
	}
	return res, nil
}

// FindByPhone ...
func (r *UserRepository) FindByPhone(phone string) (*model.User, error) {
	for _, u := range r.users {
		if u.Phone == phone {
			return u, nil
		}
	}

	return nil, store.ErrRecordNotFound
}
