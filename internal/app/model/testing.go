package model

import "testing"

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
