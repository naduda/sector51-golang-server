package model_test

import (
	"github.com/naduda/sector51-golang/internal/app/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBox_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		e       func() *model.Box
		isValid bool
	}{
		{
			name: "valid",
			e: func() *model.Box {
				return model.TestBox(t)
			},
			isValid: true,
		},
		{
			name: "invalid value",
			e: func() *model.Box {
				b := model.TestBox(t)
				b.Value = 0
				return b
			},
			isValid: false,
		},
		{
			name: "invalid type",
			e: func() *model.Box {
				b := model.TestBox(t)
				b.IdType = 8
				return b
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.e().Validate())
			} else {
				assert.Error(t, tc.e().Validate())
			}
		})
	}
}
