package processors

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"

	. "github.com/vattle/sqlboiler/boil/qm"
	"github.com/zqzca/back/controllers"
	"github.com/zqzca/back/models"
)

// BuildFile builds a file from chunks.
func BuildFile(deps controllers.Dependencies, f *models.File) (io.ReadSeeker, error) {
	chunks, err := models.Chunks(
		deps.DB,
		Where("file_id=$1", f.ID),
		OrderBy("position asc"),
	).All()

	if err != nil {
		fmt.Println("Failed to find chunks for file:", f.ID)
		return nil, err
	}

	fs := deps.Fs
	fullFilePath := filepath.Join("files", f.Hash)
	fullFile, err := fs.Create(fullFilePath)
	defer fullFile.Close()
	if err != nil {
		fmt.Println("Failed because", err)
		return nil, err
	}
	fullFileBuffer := &bytes.Buffer{}

	mw := io.MultiWriter(fullFile, fullFileBuffer)

	for _, c := range chunks {
		fmt.Println("pos:", c.Position)
		path := filepath.Join("files", "chunks", c.Hash)
		chunkData, err := fs.Open(path)

		if err != nil {
			fmt.Println("chunk Failed because", err)
			return nil, err
		}
		_, err = io.Copy(mw, chunkData)
		chunkData.Close()

		if err != nil {
			fmt.Println("Failed to copy chunk to full file")
		}
	}

	fmt.Println("Finished building file")

	for _, c := range chunks {
		path := filepath.Join("files", "chunks", c.Hash)
		if err != nil {
			fmt.Println("Failed to delete chunk entry:", c.ID)
		}
		fs.Remove(path)
	}

	bs := bytes.NewReader(fullFileBuffer.Bytes())

	return bs, nil
}
