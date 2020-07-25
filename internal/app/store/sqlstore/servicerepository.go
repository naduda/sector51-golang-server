package sqlstore

import (
	"database/sql"
	"github.com/naduda/sector51-golang/internal/app/model"
	"github.com/naduda/sector51-golang/internal/app/store"
)

// ServiceRepository ...
type ServiceRepository struct {
	store *Store
}

// Services ...
func (r *ServiceRepository) List() ([]model.Service, error) {
	users := []model.Service{}
	err := r.store.db.Select(&users, "SELECT * FROM service order by price")
	if err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateService ...
func (r *ServiceRepository) UpdateService(s model.Service) error {
	if err := s.Validate(); err != nil {
		return err
	}

	tx := r.store.db.MustBegin()
	query := "UPDATE service SET name = :name, \"desc\" = :desc, price = :price WHERE id = :id"
	res, err := tx.NamedExec(query, s)
	if err != nil {
		return err
	}
	err = tx.Commit()
	rows, _ := res.RowsAffected()
	if rows < 1 {
		return store.ErrRecordNotFound
	}
	return err
}

// Find ...
func (r *ServiceRepository) Find(id int) (*model.Service, error) {
	s := model.Service{}
	if err := r.store.db.Get(&s, "SELECT * FROM service WHERE id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return &s, nil
}

// CreateUserService ...
func (r *ServiceRepository) CreateUserService(us *model.UserService) error {
	if err := us.Validate(); err != nil {
		return err
	}

	tx := r.store.db.MustBegin()
	query := "INSERT INTO user_service (iduser, idservice, dtbeg, dtend, value) " +
		"VALUES (:iduser, :idservice, :dtbeg, :dtend, :value)"
	_, err := tx.NamedExec(query, us)
	if err != nil {
		return err
	}

	if us.IdService == 2 {
		query = "UPDATE box SET card = (SELECT card FROM userinfo WHERE created = $1) WHERE number = $2::int AND idtype = 3;"
		tx.MustExec(query, us.IdUser, us.Value)
	}

	return tx.Commit()
}

// UpdateUserService ...
func (r *ServiceRepository) UpdateUserService(us model.UserService) error {
	if err := us.Validate(); err != nil {
		return err
	}

	tx := r.store.db.MustBegin()
	query := "UPDATE user_service set dtbeg = :dtbeg, dtend = :dtend, value = :value " +
		"WHERE iduser = :iduser AND idservice = :idservice"
	res, err := tx.NamedExec(query, us)
	if err != nil {
		return err
	}
	err = tx.Commit()
	rows, _ := res.RowsAffected()
	if rows < 1 {
		return store.ErrRecordNotFound
	}
	return err
}

// GetUserServices ...
func (r *ServiceRepository) GetUserServices(idUser string) ([]model.UserService, error) {
	query := "SELECT iduser, idservice, dtbeg, dtend, COALESCE(value, '') as value FROM user_service WHERE iduser = $1"
	services := []model.UserService{}
	err := r.store.db.Select(&services, query, idUser)
	return services, err
}

// DeleteUserService ...
func (r *ServiceRepository) DeleteUserService(idUser string, idService int) error {
	tx := r.store.db.MustBegin()
	query := "DELETE FROM user_service WHERE iduser = $1 AND idservice = $2;"
	sqlRes := tx.MustExec(query, idUser, idService)

	if idService == 2 {
		query = "UPDATE box SET card = null WHERE card = (SELECT card FROM userinfo WHERE created = $1);"
		tx.MustExec(query, idUser)
	}

	rows, _ := sqlRes.RowsAffected()
	err := tx.Commit()
	if rows < 1 {
		return store.ErrRecordNotFound
	}
	return err
}
