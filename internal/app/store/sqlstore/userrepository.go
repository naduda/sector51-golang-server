package sqlstore

import (
	"database/sql"
	"time"

	"github.com/naduda/sector51-golang/internal/app/model"
	"github.com/naduda/sector51-golang/internal/app/store"
)

// UserRepository ...
type UserRepository struct {
	store *Store
}

// Create ...
func (r *UserRepository) Create(u *model.User) error {
	u.ID = time.Now().UTC().Format("2006-01-02T15:04:05.999999Z")

	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	tx := r.store.db.MustBegin()
	insertUserSecurityQuery := "INSERT INTO usersecurity (created, password) VALUES ($1, $2)"
	tx.MustExec(insertUserSecurityQuery, u.ID, u.EncryptedPassword)

	insertUserInfoQuery := "INSERT INTO userinfo (created, phone, name, surname, email) VALUES ($1, $2, $3, $4, '')"
	tx.MustExec(insertUserInfoQuery, u.ID, u.Phone, u.Name, u.Surname)
	return tx.Commit()
}

// Find ...
func (r *UserRepository) Find(id string) (*model.User, error) {
	query := "SELECT ui.created as id, ui.phone, ui.name, ui.surname, us.password as EncryptedPassword " +
		"FROM usersecurity us LEFT JOIN userinfo ui ON us.created = ui.created WHERE us.created = $1"

	u := model.User{}
	err := r.store.db.Get(&u, query, id)
	if err == sql.ErrNoRows {
		return nil, store.ErrRecordNotFound
	}

	return &u, err
}

// FindAll ...
func (r *UserRepository) FindAll() ([]model.User, error) {
	query := "SELECT ui.created as id, ui.phone, ui.name, ui.surname, us.password as EncryptedPassword " +
		"FROM usersecurity us LEFT JOIN userinfo ui ON us.created = ui.created " +
		"ORDER BY ui.surname, ui.name"

	users := []model.User{}
	if err := r.store.db.Select(&users, query); err != nil {
		return nil, err
	}

	return users, nil
}

// FindByPhone ...
func (r *UserRepository) FindByPhone(phone string) (*model.User, error) {
	query := "SELECT ui.created as id, ui.phone, ui.name, ui.surname, us.password as EncryptedPassword " +
		"FROM usersecurity us LEFT JOIN userinfo ui ON us.created = ui.created WHERE ui.phone = $1"

	u := model.User{}
	err := r.store.db.Get(&u, query, phone)
	if err == sql.ErrNoRows {
		return nil, store.ErrRecordNotFound
	}

	return &u, err
}
