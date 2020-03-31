package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

const BCRYP_MIN_COST = 12

// User ...
type User struct {
	ID                string `json:"id"`
	Phone             string `json:"phone"`
	Name              string `json:"name"`
	Surname           string `json:"surname"`
	Password          string `json:"password,omitempty"`
	Card              string `json:"card"`
	Roles             string `json:"roles"`
	IsMan             bool   `json:"isMan"`
	EncryptedPassword string `json:"-"`
}

// Validate ...
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Phone, validation.Required, is.E164),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(4, 100)),
		validation.Field(&u.Card, validation.Required, validation.Length(13, 14)),
		validation.Field(&u.Roles, validation.Required, validation.Length(3, 50)),
	)
}

// BeforeCreate ...
func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}

	return nil
}

// Sanitize ...
func (u *User) Sanitize() {
	u.Password = ""
}

// ComparePassword ...
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), BCRYP_MIN_COST)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
