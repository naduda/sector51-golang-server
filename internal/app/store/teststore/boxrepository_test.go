package teststore_test

import (
	"github.com/naduda/sector51-golang/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBoxRepository_List(t *testing.T) {
	r := teststore.New()
	all, err := r.Boxes().List()
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, 150, len(all))
}
