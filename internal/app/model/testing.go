package model

import (
	"testing"
	"time"
)

// TestUser ...
func TestUser(t *testing.T) *User {
	t.Helper()

	return &User{
		Phone:    "+380661122334",
		Password: "password",
		Name:     "Name",
		Surname:  "Surname",
	}
}

// TestService ...
func TestService(t *testing.T) *Service {
	t.Helper()

	return &Service{
		ID:    6,
		Name:  "ABONEMENT (Morning) 3",
		Desc:  "-",
		Price: 1600,
	}
}

// TestUserService ...
func TestUserService(t *testing.T) *UserService {
	t.Helper()

	id := time.Now().UTC().Format("2006-01-02T15:04:05.999999Z")
	beg := time.Now()
	end := time.Now()
	return &UserService{
		IdService: 6,
		IdUser:    id,
		DtBeg:     beg,
		DtEnd:     end,
		Value:     "",
	}
}
