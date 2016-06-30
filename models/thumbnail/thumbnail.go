package thumbnail

import (
	"time"

	"github.com/zqzca/back/db"
)

type Thumbnail struct {
	ID        string    `json:"id"`
	Size      int       `json:"size"`
	Hash      string    `json:"hash"`
	FileID    string    `json:"file_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const insertSQL = `
	INSERT INTO thumbnails
	(size, hash, file_id)
	VALUES
	($1, $2, $3)
	RETURNING id`

const findByFileIDSQL = `
	SELECT
	id, size, hash, created_at, updated_at
	FROM thumbnails
	WHERE id = $1
`

const deleteByFileIDSQL = `
	DELETE FROM thumbnails
	WHERE file_id = $1
`

// Create a thumbnail inside of a transaction.
func (t *Thumbnail) Create(ex db.Executor) error {
	err := ex.
		QueryRow(insertSQL, t.Size, t.Hash, t.FileID).
		Scan(&t.ID)

	return err
}

func FindByFileID(ex db.Executor, id string) (*Thumbnail, error) {
	var t Thumbnail
	t.FileID = id
	err := ex.QueryRow(findByFileIDSQL, id).Scan(
		&t.ID, &t.Size, &t.Hash, &t.CreatedAt, &t.UpdatedAt,
	)

	return &t, err
}

func DeleteByFileID(ex db.Executor, id string) error {
	_, err := ex.Exec(deleteByFileIDSQL, id)

	return err
}
