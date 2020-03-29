package store

import (
	"github.com/naduda/sector51-golang/internal/app/model"
)

// UserRepository ...
type UserRepository interface {
	Create(*model.User) error
	Find(string) (*model.User, error)
	FindAll() ([]model.User, error)
	FindByPhone(string) (*model.User, error)
}

// ServiceRepository ...
type ServiceRepository interface {
	List() ([]model.Service, error)
	Find(int) (*model.Service, error)
	UpdateService(model.Service) error
	CreateUserService(*model.UserService) error
	GetUserServices(string) ([]model.UserService, error)
	DeleteUserService(string, int) error
}
