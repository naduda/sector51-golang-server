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
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	// return r.store.db.QueryRow(
	// 	"INSERT INTO users (phone, encrypted_password) VALUES ($1, $2) RETURNING id",
	// 	u.Phone,
	// 	u.EncryptedPassword,
	// ).Scan(&u.ID)
	if err := r.store.db.QueryRow(
		"INSERT INTO usersecurity (created, password) VALUES ($1, $2) RETURNING created",
		time.Now().UTC(),
		u.EncryptedPassword,
	).Scan(&u.ID); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO userinfo (created, phone, name, surname, email) VALUES ($1, $2, $3, $4, 'fake@gmail.com') RETURNING created",
		u.ID,
		u.Phone,
		u.Name,
		u.Surname,
	).Scan(&u.ID)
}

// Find ...
func (r *UserRepository) Find(id string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT ui.created, ui.phone, ui.name, ui.surname, us.password FROM usersecurity us LEFT JOIN userinfo ui ON us.created = ui.created WHERE us.created = $1",
		id,
	).Scan(
		&u.ID,
		&u.Phone,
		&u.Name,
		&u.Surname,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

// FindByPhone ...
func (r *UserRepository) FindByPhone(phone string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT ui.created, ui.phone, ui.name, ui.surname, us.password FROM usersecurity us LEFT JOIN userinfo ui ON us.created = ui.created WHERE ui.phone = $1",
		phone,
	).Scan(
		&u.ID,
		&u.Phone,
		&u.Name,
		&u.Surname,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}
