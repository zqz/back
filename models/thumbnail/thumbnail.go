package thumbnail

import (
	"database/sql"
	"time"
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

// Create a thumbnail inside of a transaction.
func (t *Thumbnail) Create(tx *sql.Tx) error {
	err := tx.
		QueryRow(insertSQL, t.Size, t.Hash, t.FileID).
		Scan(&t.ID)

	return err
}

func FindByFileID(tx *sql.Tx, id string) (*Thumbnail, error) {
	var t Thumbnail
	t.FileID = id
	err := tx.QueryRow(findByFileIDSQL, id).Scan(
		&t.ID, &t.Size, &t.Hash, &t.CreatedAt, &t.UpdatedAt,
	)

	return &t, err
}
