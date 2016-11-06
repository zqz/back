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
	Downloads int       `json:"download"`
	CreatedAt time.Time `json:"created_at"`
}

func NewFile(db db.Executor, file *models.File) File {
	return File{
		Slug:      file.Slug,
		Size:      file.Size,
		Name:      file.Name,
		Hash:      file.Hash,
		Type:      file.Type,
		Downloads: int(file.Downloads(db).CountP()),
		CreatedAt: file.CreatedAt,
	}
}
