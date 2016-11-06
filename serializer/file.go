package serializer

import (
	"time"

	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models"
)

type File struct {
	Slug      string    `json:"slug"`
	Size      int       `json:"size"`
	Name      string    `json:"name"`
	Hash      string    `json:"hash"`
	Type      string    `json:"type"`
	Downloads int       `json:"downloads"`
	CreatedAt time.Time `json:"created_at"`
}

var FileDownloads func(db.Executor, *models.File) int

// ForFile
func ForFile(db db.Executor, f *models.File) File {
	return File{
		Slug:      f.Slug,
		Size:      f.Size,
		Name:      f.Name,
		Hash:      f.Hash,
		Type:      f.Type,
		Downloads: FileDownloads(db, f),
		CreatedAt: f.CreatedAt,
	}
}

func init() {
	FileDownloads = func(ex db.Executor, f *models.File) int {
		return int(f.Downloads(ex).CountP())
	}
}
