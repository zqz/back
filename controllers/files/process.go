package files

import (
	"net/http"

	"github.com/zqzca/echo"
)

// Process builds thumbnails
func (f FileController) Process(e echo.Context) error {
	f.Debug("processing")

	//go func() {
	//f.Debug("inprocessing")

	//fileID := c.Param("id")
	//file, err := file.FindByID(db.Connection, fileID)

	//if file.State == Processing {
	//	f.Debug("file already being processed")
	//	return
	//}

	//if err != nil {
	//	f.Debug("failed to find file", "err", err)
	//	return
	//}

	//tx, err := boil.Begin()
	//if err != nil {
	//	f.Debug("couldn't open transaction", "err", err)
	//	return
	//}
	//err := models.Thumbnails(Where("file_id=$1", f.ID)).DeleteAll()
	//// f.Process(tx) TODO(dylanj): Models derp derp derp
	//tx.Commit()
	//}()

	return e.NoContent(http.StatusOK)
}
