package model

import validation "github.com/go-ozzo/ozzo-validation"

// Box ...
type Box struct {
	Card   string `json:"card"`
	IdType int    `json:"idType"`
	Value  int    `json:"number"`
	Time   string `json:"time"`
}

// Validate Box ...
func (s *Box) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Card, validation.Length(13, 14)),
		validation.Field(&s.IdType, validation.Required, validation.Min(1), validation.Max(3)),
		validation.Field(&s.Value, validation.Required, validation.Min(1), validation.Max(50)),
	)
}
