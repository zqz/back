package models

import (
	"errors"
	"log"

	"upper.io/db"
)

type File struct {
	ID     string `json:"id" db:"id,omitempty"`
	Size   uint32 `json:"size" valid:"required,numeric" db:"size"`
	Hash   string `json:"hash" valid:"required,alphanum" db:"hash"`
	Done   bool   `json:"done" db:"done"`
	Chunks uint32 `json:"numchunks" valid:"required,numeric" db:"chunks"`
	Type   string `json:"type" valid:"required" db:"type"`
}

func (f *File) Save() bool {
	if len(f.ID) == 0 {
		return f.Create()
	}

	return false
}

func (f *File) Create() bool {
	fc := fileCollection()

	if fc == nil {
		return false
	}

	var err error
	if _, err = fc.Append(f); err != nil {
		log.Println("failed to create user", err.Error())
		return false
	}

	return true
}

func FileFindByHash(hash string) (*File, error) {
	fc := fileCollection()
	res := fc.Find(db.Cond{"hash": hash})
	var f File

	if count, _ := res.Count(); count > 0 {
		res.One(&f)
	}

	if len(f.ID) == 0 {
		return nil, errors.New("No File Found with that Hash")
	}

	return &f, nil
}

func fileCollection() db.Collection {
	col, err := database.Collection("files")

	if err != nil {
		log.Fatalln("Failed to find files collection", err.Error())
	}

	return col
}

// SetID allows us to update the struct after the DB sets the ID
func (f *File) SetID(values map[string]interface{}) error {
	if valueInterface, ok := values["id"]; ok {
		f.ID = valueInterface.(string)
	}
	return nil
}
