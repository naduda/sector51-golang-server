package sqlstore_test

import (
	"github.com/naduda/sector51-golang/internal/app/model"
	"github.com/naduda/sector51-golang/internal/app/store"
	"github.com/naduda/sector51-golang/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServiceRepository_List(t *testing.T) {
	db, _ := sqlstore.TestDB(t, databaseURL)
	r := sqlstore.New(db)

	all, err := r.Service().List()
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, 15, len(all))
}

func TestServiceRepository_UpdateService(t *testing.T) {
	db, _ := sqlstore.TestDB(t, databaseURL)
	r := sqlstore.New(db)

	all, err := r.Service().List()
	if err != nil {
		t.Fatal(err)
	}

	first := all[0]
	original := first.Price
	first.Price = 8888

	err = r.Service().UpdateService(first)
	assert.NoError(t, err)

	f, _ := r.Service().Find(first.ID)
	assert.Equal(t, 8888, f.Price)
	first.Price = original
	_ = r.Service().UpdateService(first)

	first.ID = 100
	err = r.Service().UpdateService(first)
	assert.Error(t, store.ErrRecordNotFound)
}

func TestServiceRepository_Find(t *testing.T) {
	db, _ := sqlstore.TestDB(t, databaseURL)
	r := sqlstore.New(db)

	s, err := r.Service().Find(13)
	assert.NoError(t, err)
	assert.Equal(t, 4400, s.Price)

	_, err = r.Service().Find(100)
	assert.Error(t, store.ErrRecordNotFound)
}

func TestServiceRepository_CreateUserService(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("user_service")
	r := sqlstore.New(db)
	us := model.TestUserService(t)
	us.IdService = 13
	err := r.Service().CreateUserService(us)
	assert.NoError(t, err)
}

func TestServiceRepository_UpdateUserService(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("user_service")
	r := sqlstore.New(db)
	us := model.TestUserService(t)
	us.IdService = 13
	err := r.Service().CreateUserService(us)
	assert.NoError(t, err)

	res, err := r.Service().GetUserServices(us.IdUser)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))

	res[0].Value = "50"
	err = r.Service().UpdateUserService(res[0])
	assert.NoError(t, err)

	_, err = r.Service().GetUserServices("")
	assert.Error(t, store.ErrRecordNotFound)

	exists, err := r.Service().GetUserServices(us.IdUser)
	assert.NoError(t, err)
	assert.Equal(t, "50", exists[0].Value)
}

func TestServiceRepository_GetUserServices(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("user_service")
	r := sqlstore.New(db)
	us := model.TestUserService(t)
	us.IdService = 13

	err := r.Service().CreateUserService(us)
	assert.NoError(t, err)

	res, err := r.Service().GetUserServices(us.IdUser)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))

	_, err = r.Service().GetUserServices("")
	assert.Error(t, store.ErrRecordNotFound)
}

func TestServiceRepository_DeleteUserService(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("user_service")
	r := sqlstore.New(db)
	us := model.TestUserService(t)
	us.IdService = 13

	err := r.Service().CreateUserService(us)
	_, err = r.Service().GetUserServices(us.IdUser)
	assert.NoError(t, err)

	err = r.Service().DeleteUserService(us.IdUser, 100)
	assert.Error(t, store.ErrRecordNotFound)

	err = r.Service().DeleteUserService(us.IdUser, us.IdService)
	assert.NoError(t, err)

	res, _ := r.Service().GetUserServices(us.IdUser)
	assert.Equal(t, 0, len(res))
}
