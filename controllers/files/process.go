package files

import (
	"fmt"
	"net/http"

	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models/file"
	"github.com/zqzca/back/models/thumbnail"
	"github.com/zqzca/echo"
)

// Process builds thumbnails
func Process(c echo.Context) error {
	fmt.Println("processing")
	go func() {
		fmt.Println("inprocessing")

		fileID := c.Param("id")
		f, err := file.FindByID(db.Connection, fileID)

		if f.State == file.Processing {
			fmt.Println("file already being processed")
			return
		}

		if err != nil {
			fmt.Println("failed to find file", err)
			return
		}

		tx := db.StartTransaction()
		thumbnail.DeleteByFileID(tx, f.ID)
		f.Process(tx)
		tx.Commit()
	}()

	return c.NoContent(http.StatusOK)
}
