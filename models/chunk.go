package models

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"upper.io/db"
)

type Chunk struct {
	ID       string `json:"id" db:"id,omitempty"`
	Size     int    `json:"size" db:"size"`
	Hash     string `json:"hash" db:"hash"`
	FileID   string `json:"file_id" db:"file_id"`
	Position int    `json:"position" db:"position"`
}

func (c *Chunk) Path() string {
	return fmt.Sprintf("files/chunks/%s", c.Hash)
}

func ChunkFindByFileIDAndPosition(fileID string, position int) (*Chunk, error) {
	cc := chunkCollection()

	if cc == nil {
		return nil, errors.New("Can't lookup chunk collection")
	}

	res := cc.Find(db.Cond{"file_id": fileID, "position": position})
	var c Chunk

	if count, _ := res.Count(); count > 0 {
		res.One(&c)
	}

	if len(c.ID) == 0 {
		return nil, errors.New("No Chunk Found with that position and id")
	}

	return &c, nil
}

func ChunkFindByHash(hash string) (*Chunk, error) {
	cc := chunkCollection()

	if cc == nil {
		return nil, errors.New("Can't lookup chunk collection")
	}

	res := cc.Find(db.Cond{"hash": hash})
	var c Chunk

	if count, _ := res.Count(); count > 0 {
		res.One(&c)
	}

	if len(c.ID) == 0 {
		return nil, errors.New("No Chunk Found with that Hash")
	}

	return &c, nil
}

func ChunksByFileID(fileID string) ([]Chunk, error) {
	cc := chunkCollection()

	if cc == nil {
		return nil, errors.New("Can't lookup chunk collection")
	}

	res := cc.Find(db.Cond{"file_id": fileID})

	var chunks []Chunk
	err := res.Sort("position").All(&chunks)

	if err != nil {
		return nil, err
	}

	return chunks, nil
}

func (c *Chunk) Save() bool {
	if len(c.ID) == 0 {
		return c.Create()
	}

	return false
}

func (c *Chunk) Create() bool {
	cc := chunkCollection()

	if cc == nil {
		return false
	}

	var err error
	if _, err = cc.Append(c); err != nil {
		log.Println("Failed to create Chunk", err.Error())
		return false
	}

	return true
}

// SetID allows us to update the struct after the DB sets the ID
func (c *Chunk) SetID(values map[string]interface{}) error {
	if valueInterface, ok := values["id"]; ok {
		c.ID = valueInterface.(string)
	}

	return nil
}

func chunkCollection() db.Collection {
	col, err := Database.Collection("chunks")

	if err != nil {
		log.Fatalln("Failed to find chunks collection", err.Error())
	}

	return col
}

func (c *Chunk) Fd() (io.ReadCloser, error) {
	return os.Open(c.Path())
}
