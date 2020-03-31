package apiserver

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/naduda/sector51-golang/internal/app/model"
	"github.com/naduda/sector51-golang/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

var jwtSecret = "jwt-secret-test"

func TestServer_AuthUser(t *testing.T) {
	store := teststore.New()
	u := model.TestUser(t)
	if err := store.User().Create(u); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name         string
		hasToken     bool
		expectedCode int
	}{
		{
			name:         "authenticated",
			hasToken:     true,
			expectedCode: http.StatusOK,
		},
		{
			name:         "not authenticated",
			hasToken:     false,
			expectedCode: http.StatusForbidden,
		},
	}

	logger := logrus.New()
	s := newServer(store, jwtSecret, logger)
	mw := s.authUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	token, _ := s.generateJWT(u)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			if tc.hasToken {
				req.Header.Set("Authorization", token)
			}
			mw.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleUsersCreate(t *testing.T) {
	logger := logrus.New()
	s := newServer(teststore.New(), jwtSecret, logger)
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"phone":    "+380505555555",
				"password": "secret",
				"card":     "1100000001102",
				"roles":    "USER",
				"isMan":    true,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]interface{}{
				"phone":    "invalid",
				"password": "short",
				"card":     "1100000001102",
				"roles":    "USER",
				"isMan":    true,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			if err := json.NewEncoder(b).Encode(tc.payload); err != nil {
				t.Fatal(err)
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/users", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleLogin(t *testing.T) {
	store := teststore.New()
	u := model.TestUser(t)
	if err := store.User().Create(u); err != nil {
		t.Fatal(err)
	}
	logger := logrus.New()
	s := newServer(store, jwtSecret, logger)
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"phone":    u.Phone,
				"password": u.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid phone",
			payload: map[string]interface{}{
				"phone":    "invalid",
				"password": u.Password,
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]interface{}{
				"phone":    u.Phone,
				"password": "invalid",
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			if err := json.NewEncoder(b).Encode(tc.payload); err != nil {
				t.Fatal(err)
			}
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/login", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
