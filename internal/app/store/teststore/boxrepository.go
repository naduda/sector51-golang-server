package teststore

import "github.com/naduda/sector51-golang/internal/app/model"

// ServiceRepository ...
type BoxRepository struct {
	store *Store
	boxes map[int]model.Box
}

// Boxes ...
func (r *BoxRepository) List() ([]model.Box, error) {
	res := make([]model.Box, 0)
	for _, u := range r.boxes {
		res = append(res, u)
	}
	return res, nil
}

// Update ...
func (r *BoxRepository) Update(box model.Box) error {
	return nil
}

// GetByCard ...
func (r *BoxRepository) GetByCard(card string) (*model.Box, error) {
	return nil, nil
}
