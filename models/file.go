package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"upper.io/db"
)

type File struct {
	ID        string    `json:"id" db:"id,omitempty"`
	Size      uint32    `json:"size" valid:"required,numeric" db:"size"`
	Hash      string    `json:"hash" valid:"required,alphanum" db:"hash"`
	Chunks    int       `json:"chunks" valid:"required,numeric" db:"chunks"`
	Name      string    `json:"name" valid:"required" db:"name"`
	Type      string    `json:"type" valid:"required" db:"type"`
	State     uint8     `json:"state" db:"state"`
	CreatedAt time.Time `json:"created_at" db:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at,omitempty"`
}

const (
	Incomplete = iota
	Assembling
	Processing
	Finished
)

func (f *File) Assemble() {
	fmt.Println("Now we assemble")
	f.State = Assembling
	f.Save()
}

func (f *File) Save() bool {
	if len(f.ID) == 0 {
		return f.Create()
	}

	return f.Update()
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

// Update a file.
func (f *File) Update() bool {
	fc := fileCollection()

	if fc == nil {
		return false
	}

	res := fc.Find(db.Cond{"id": f.ID})

	if err := res.Update(f); err != nil {
		log.Println("failed to update file", err.Error())
		return false
	}

	return true
}

func FilePagination(page uint, per_page uint) ([]File, error) {
	fc := fileCollection()

	offset := page * per_page

	var files []File
	res := fc.Find().Skip(offset).Limit(per_page)

	err := res.All(&files)

	if err != nil {
		fmt.Println("Failed to fetch files", err)
		return nil, err
	}

	return files, nil
}

func FileFindByID(id string) (*File, error) {
	fc := fileCollection()
	res := fc.Find(db.Cond{"id": id})
	var f File

	if count, _ := res.Count(); count > 0 {
		res.One(&f)
	}

	if len(f.ID) == 0 {
		return nil, errors.New("No File Found with that ID")
	}

	return &f, nil
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
	col, err := Database.Collection("files")

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

func (f *File) Process() {

}
