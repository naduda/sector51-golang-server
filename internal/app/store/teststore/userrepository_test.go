package teststore_test

import (
	"testing"

	"github.com/naduda/sector51-golang/internal/app/model"
	"github.com/naduda/sector51-golang/internal/app/store"
	"github.com/naduda/sector51-golang/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u.ID)
}

func TestUserRepository_Find(t *testing.T) {
	s := teststore.New()
	u1 := model.TestUser(t)
	if err := s.User().Create(u1); err != nil {
		t.Fatal(err)
	}
	u2, err := s.User().Find(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestUserRepository_FindAll(t *testing.T) {
	s := teststore.New()
	u1 := model.TestUser(t)
	if err := s.User().Create(u1); err != nil {
		t.Fatal(err)
	}
	r, err := s.User().FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r))
}

func TestUserRepository_FindByPhone(t *testing.T) {
	s := teststore.New()
	u1 := model.TestUser(t)
	_, err := s.User().FindByPhone(u1.Phone)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	if err = s.User().Create(u1); err != nil {
		t.Fatal(err)
	}
	u2, err := s.User().FindByPhone(u1.Phone)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
