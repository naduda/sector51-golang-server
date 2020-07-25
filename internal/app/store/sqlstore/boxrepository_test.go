package sqlstore_test

import (
	"github.com/naduda/sector51-golang/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBoxRepository_List(t *testing.T) {
	db, _ := sqlstore.TestDB(t, databaseURL)
	r := sqlstore.New(db)

	all, err := r.Boxes().List()
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, 150, len(all))
}
