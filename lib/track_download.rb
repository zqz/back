package lib

import (
	"fmt"

	null "gopkg.in/nullbio/null.v5"

	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models"
	"github.com/zqzca/echo"
)

// TrackDownload stores a record for the download.
func TrackDownload(db db.Executor, fileID string, e echo.Context, hit bool) {
	d := models.Download{
		FileID:   null.StringFrom(fileID),
		Ip:       null.StringFrom(e.Request().RemoteAddress()),
		CacheHit: hit,
	}

	if err := d.Insert(db); err != nil {
		fmt.Println("failed to insert download thing", err.Error())
	}
}
