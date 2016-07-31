package chunk

import (
	"log"
	"time"

	"database/sql"

	"github.com/zqzca/back/db"
)

type Chunk struct {
	ID        string    `json:"id"`
	FileID    string    `json:"file_id"`
	Size      int       `json:"size"`
	Hash      string    `json:"hash"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const findByIDSQL = `
	SELECT
	file_id, size, hash, position, created_at, updated_at
	FROM chunks
	WHERE id = $1
`

const findByFileIDSQL = `
	SELECT
	id, size, hash, position
	FROM chunks
	WHERE file_id = $1
	ORDER BY position asc
`

const haveChunkForFileSQL = `
	SELECT EXISTS (
		SELECT 1
		FROM chunks
		WHERE file_id = $1 AND position = $2
	)
`

const haveChunkForFileWithHashSQL = `
	SELECT EXISTS (
		SELECT 1
		FROM chunks
		WHERE file_id = $1 AND hash = $2
	)
`

const insertSQL = `
	INSERT INTO chunks
	(file_id, size, hash, position)
	VALUES
	($1, $2, $3, $4)
	RETURNING id
`

const updateChunkSQL = `
	UPDATE chunks
	SET file_id = $2, size = $3, hash = $4, position = $5
	WHERE id = $1;
`

// FindByID returns a chunk with the specified id.
func FindByID(ex db.Executor, id string) (*Chunk, error) {
	var c Chunk
	c.ID = id
	err := ex.QueryRow(findByIDSQL, id).Scan(
		&c.FileID, &c.Size, &c.Hash, &c.Position, &c.CreatedAt, &c.UpdatedAt,
	)
	return &c, err
}

// FindByFileID return all chunks with the specified FileID.
// TODO: cleanup
func FindByFileID(ex db.Executor, id string) ([]*Chunk, error) {
	var chunks []*Chunk
	var err error
	var rows *sql.Rows

	if rows, err = ex.Query(findByFileIDSQL, id); err != nil {
		return chunks, err
	}
	defer rows.Close()

	for rows.Next() {
		c := &Chunk{}

		if err = rows.Scan(&c.ID, &c.Size, &c.Hash, &c.Position); err != nil {
			log.Fatal(err)
		}

		chunks = append(chunks, c)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return chunks, err
}

func HaveChunkForFile(ex db.Executor, fileID string, position int) bool {
	var exists bool
	err := ex.QueryRow(haveChunkForFileSQL, fileID, position).Scan(&exists)

	if err != nil {
		return false
	}
	return exists
}

func HaveChunkForFileWithHash(ex db.Executor, fid string, hash string) bool {
	var exists bool

	err := ex.QueryRow(haveChunkForFileWithHashSQL, fid, hash).Scan(&exists)

	if err != nil {
		return false
	}
	return exists
}

// Create a chunk inside of a transaction.
func (c *Chunk) Create(ex db.Executor) error {
	err := ex.
		QueryRow(insertSQL, c.FileID, c.Size, c.Hash, c.Position).
		Scan(&c.ID)

	return err
}
