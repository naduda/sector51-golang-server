package model_test

import (
	"testing"
	"time"

	"github.com/naduda/sector51-golang/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestService_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		e       func() *model.Service
		isValid bool
	}{
		{
			name: "valid",
			e: func() *model.Service {
				return model.TestService(t)
			},
			isValid: true,
		},
		{
			name: "invalid",
			e: func() *model.Service {
				s := model.TestService(t)
				s.Name = "sm"
				return s
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.e().Validate())
			} else {
				assert.Error(t, tc.e().Validate())
			}
		})
	}
}

func TestUserService_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		e       func() *model.UserService
		isValid bool
	}{
		{
			name: "valid",
			e: func() *model.UserService {
				return model.TestUserService(t)
			},
			isValid: true,
		},
		{
			name: "invalid",
			e: func() *model.UserService {
				us := model.TestUserService(t)
				us.DtEnd = time.Date(2000, time.January, 01, 0, 0, 0, 0, time.UTC)
				return us
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.e().Validate())
			} else {
				assert.Error(t, tc.e().Validate())
			}
		})
	}
}
