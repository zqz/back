package file_test

import (
	"database/sql"
	"testing"

	"github.com/zqzca/back/models"
)

func TestCopy(t *testing.T) {
	t.Parallel()
	models.TxWrapper(func(tx *sql.Tx) {
		// a := assert.New(t)

		// Not sure yet :(
	})
}
