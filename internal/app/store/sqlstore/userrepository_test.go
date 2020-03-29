package sqlstore_test

import (
	"testing"

	"github.com/naduda/sector51-golang/internal/app/model"
	"github.com/naduda/sector51-golang/internal/app/store"
	"github.com/naduda/sector51-golang/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("usersecurity", "userinfo")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u.ID)
}

func TestUserRepository_Find(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("usersecurity", "userinfo")

	s := sqlstore.New(db)
	u1 := model.TestUser(t)
	if err := s.User().Create(u1); err != nil {
		t.Fatal(err)
	}
	u2, err := s.User().Find(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)

	_, err = s.User().Find("2000-01-01T00:00:00.000000Z")
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
}

func TestUserRepository_FindByPhone(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("usersecurity", "userinfo")

	s := sqlstore.New(db)
	u1 := model.TestUser(t)
	_, err := s.User().FindByPhone(u1.Phone)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	if err := s.User().Create(u1); err != nil {
		t.Fatal(err)
	}
	u2, err := s.User().FindByPhone(u1.Phone)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
