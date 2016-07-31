package files

import (
	"fmt"
	"net/http"

	"github.com/zqzca/echo"
	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models/file"
)

func fileExistsWithHash(hash string) bool {
	// Todo write an exists? for this
	existing, _ := file.FindByHash(db.Connection, hash)

	return len(existing.ID) > 0
}

// Create creates a file container in the database.
func Create(c echo.Context) error {
	f := &file.File{}

	if err := c.Bind(f); err != nil {
		return err
	}

	if fileExistsWithHash(f.Hash) {
		fmt.Println("file exists with hash:", f.Hash)
		return c.NoContent(http.StatusConflict)
	}

	fmt.Println("file doesnt exists with hash:", f.Hash)

	tx := db.StartTransaction()
	f.State = file.Incomplete
	f.Create(tx)
	tx.Commit()

	return c.JSON(http.StatusOK, f)
}
