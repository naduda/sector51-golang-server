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
	insertUserSecurityQuery := "INSERT INTO usersecurity (created, password, roles) VALUES ($1, $2, $3)"
	tx.MustExec(insertUserSecurityQuery, u.ID, u.EncryptedPassword, u.Roles)

	insertUserInfoQuery := "INSERT INTO userinfo (created, phone, name, surname, card, sex, email) " +
		"VALUES ($1, $2, $3, $4, $5, $6, '')"
	tx.MustExec(insertUserInfoQuery, u.ID, u.Phone, u.Name, u.Surname, u.Card, u.IsMan)
	return tx.Commit()
}

// Update ...
func (r *UserRepository) Update(u model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	tx := r.store.db.MustBegin()
	query := "UPDATE userinfo SET phone = $1, name = $2, surname = $3, card = $4, sex = $5 WHERE created = $6"
	res := tx.MustExec(query, u.Phone, u.Name, u.Surname, u.Card, u.IsMan, u.ID)
	err := tx.Commit()
	rows, _ := res.RowsAffected()
	if rows < 1 {
		return store.ErrRecordNotFound
	}
	return err
}

// Delete ...
func (r *UserRepository) Delete(id string) error {
	queryUI := "DELETE FROM userinfo WHERE created = $1"
	queryUS := "DELETE FROM usersecurity WHERE created = $1"
	tx := r.store.db.MustBegin()
	tx.MustExec(queryUI, id)
	tx.MustExec(queryUS, id)
	return tx.Commit()
}

// Find ...
func (r *UserRepository) Find(id string) (*model.User, error) {
	query := "SELECT ui.created as id, ui.phone, ui.name, ui.surname, ui.card, ui.sex as isMan, us.password as EncryptedPassword " +
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
	query := "SELECT s.dtend is not null and s.dtend > now() as active, ui.created as id, ui.phone, ui.name, " +
		"ui.surname, us.password as EncryptedPassword, ui.card, us.roles, ui.sex as isMan " +
		"FROM usersecurity us LEFT JOIN userinfo ui ON us.created = ui.created " +
		"LEFT JOIN user_service s ON us.created = s.iduser " +
		"ORDER BY active DESC, ui.surname, ui.name;"

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
