package processors

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zqzca/back/models"
)

// BuildFile builds a file from chunks.
func BuildFile(f *models.File) (io.Reader, error) {
	chunks, err := models.Chunks(
		Where("file_id=$1", f.ID),
		OrderBy("position asc"),
	).All()

	if err != nil {
		fmt.Println("Failed to find chunks for file:", f.ID)
		return nil, err
	}

	fullFilePath := filepath.Join("files", f.Hash)
	fullFile, err = os.Create(fullFilePath)
	fullFileBuffer := &bytes.Buffer{}
	defer fullFile.Close()

	if err != nil {
		fmt.Println("Failed because", err)
		return nil, err
	}

	mw := io.MultiWriter(fullFile, fullFileBuffer)

	for _, c := range chunks {
		path := filepath.Join("files", "chunks", c.Hash)
		chunkData, err := os.Open(path)

		n, err := io.Copy(mw, chunkData)
		f.Close()

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
		os.Remove(path)
	}

	return fileDataBuffer, nil
}
