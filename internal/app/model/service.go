package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

// Service ...
type Service struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Price int    `json:"price"`
}

// Validate Service ...
func (s *Service) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Name, validation.Required, validation.Length(3, 25)),
		validation.Field(&s.Desc, validation.Required, validation.Length(1, 50)),
		validation.Field(&s.Price, validation.Required, validation.Min(0)),
	)
}

// UserService ...
type UserService struct {
	IdService int       `json:"id_service"`
	IdUser    string    `json:"id_user"`
	DtBeg     time.Time `json:"dt_beg"`
	DtEnd     time.Time `json:"dt_end"`
	Value     string    `json:"value"`
}

// Validate UserService ...
func (s *UserService) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.DtBeg, validation.Required),
		validation.Field(&s.DtEnd, validation.Required, validation.Min(s.DtBeg)),
		validation.Field(&s.IdUser, validation.Required),
		validation.Field(&s.IdService, validation.Required),
	)
}
