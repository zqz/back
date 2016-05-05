package file

import (
	"database/sql"
	"time"
)

type File struct {
	ID        string    `json:"id"`
	Size      int       `json:"size" valid:"required,numeric"`
	Hash      string    `json:"hash" valid:"required,alphanum"`
	Chunks    int       `json:"chunks" valid:"required,numeric"`
	Name      string    `json:"name" valid:"required"`
	Type      string    `json:"type" valid:"required"`
	State     int       `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	Incomplete = iota
	Processing
	Finished
)

const findByIDSQL = `
	SELECT
	size, hash, chunks, name, type, state, created_at, updated_at
	FROM files
	WHERE id = $1`

const findByHashSQL = `
	SELECT
	id, size, chunks, name, type, state, created_at, updated_at
	FROM files
	WHERE hash = $1`

const insertSQL = `
	INSERT INTO files
	(size, hash, chunks, name, type, state)
	VALUES
	($1, $2, $3, $4, $5, $6)
	RETURNING id
`

const setStateSQL = `
	UPDATE files
	SET state = $2
	WHERE id = $1
	RETURNING state
`

// Create a File inside of a transaction.
func (f *File) Create(tx *sql.Tx) error {
	err := tx.
		QueryRow(insertSQL, f.Size, f.Hash, f.Chunks, f.Name, f.Type, f.State).
		Scan(&f.ID)

	return err
}

// FindByHash returns a File with the specified hash.
func FindByHash(tx *sql.Tx, hash string) (*File, error) {
	var f File
	f.Hash = hash
	err := tx.QueryRow(findByHashSQL, hash).Scan(
		&f.ID, &f.Size, &f.Chunks, &f.Name, &f.Type, &f.State,
		&f.CreatedAt, &f.UpdatedAt,
	)
	return &f, err
}

// FindByID returns a File with the specified id.
func FindByID(tx *sql.Tx, id string) (*File, error) {
	var f File
	f.ID = id
	err := tx.QueryRow(findByIDSQL, id).Scan(
		&f.Size, &f.Hash, &f.Chunks, &f.Name, &f.Type, &f.State,
		&f.CreatedAt, &f.UpdatedAt,
	)
	return &f, err
}

// SetState sets the state of the File.
func (f *File) SetState(tx *sql.Tx, state int) error {
	err := tx.QueryRow(setStateSQL, f.ID, state).Scan(&f.State)
	return err
}
