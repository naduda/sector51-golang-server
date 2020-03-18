package store

import (
	"github.com/naduda/sector51-golang/internal/app/model"
)

// UserRepository ...
type UserRepository interface {
	Create(*model.User) error
	Find(string) (*model.User, error)
	FindAll() ([]*model.User, error)
	FindByPhone(string) (*model.User, error)
}
