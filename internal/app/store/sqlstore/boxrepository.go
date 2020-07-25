package sqlstore

import (
	"database/sql"
	"github.com/naduda/sector51-golang/internal/app/model"
	"github.com/naduda/sector51-golang/internal/app/store"
)

// ServiceRepository ...
type BoxRepository struct {
	store *Store
}

// Boxes ...
func (r *BoxRepository) List() ([]model.Box, error) {
	boxes := []model.Box{}
	err := r.store.db.Select(&boxes, "SELECT COALESCE(card, '') as card, idtype, number as value, time "+
		"FROM box ORDER BY value")
	if err != nil {
		return nil, err
	}
	return boxes, nil
}

// GetByCard ...
func (r *BoxRepository) GetByCard(card string) (*model.Box, error) {
	query := "SELECT COALESCE(card, '') as card, idtype, number as value, time FROM box WHERE card = $1"
	result := model.Box{}
	err := r.store.db.Get(&result, query, card)
	if err == sql.ErrNoRows {
		return nil, store.ErrRecordNotFound
	}

	return &result, err
}

// Update ...
func (r *BoxRepository) Update(box model.Box) error {
	if err := box.Validate(); err != nil {
		return err
	}

	tx := r.store.db.MustBegin()
	query := "UPDATE box set card = :card, time = now() " +
		"WHERE idtype = :idtype AND number = :value"
	res, err := tx.NamedExec(query, box)
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
