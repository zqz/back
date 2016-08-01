package processors

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zqzca/back/models"
)

// BuildFile builds a file from chunks.
func BuildFile(file *models.File) error {
	chunks, err := models.Chunks(
		Where("file_id=$1", f.ID),
		OrderBy("position asc"),
	).All()

	if err != nil {
		return fmt.Println("Failed to find chunks for file:", f.ID)
	}

	fullFilePath := filepath.Join("files", f.Hash)
	fullFile, err = os.Create(fullFilePath)
	defer fullFile.Close()

	if err != nil {
		fmt.Println("Failed because", err)
		return err
	}

	for _, c := range chunks {
		path := filepath.Join("files", "chunks", c.Hash)
		chunkData, err := os.Open(path)
		defer f.Close()

		n, err := io.Copy(fullFile, chunkData)
	}

	fmt.Println("Finished building file")

	// Remove chunks
	for _, c := range chunks {
		path := filepath.Join("files", "chunks", c.Hash)
		err = c.Delete()
		if err != nil {
			fmt.Println("Failed to delete chunk entry:", c.ID)
		}

		err = os.Remove(path)
		if err != nil {
			fmt.Println("Failed to delete chunk at", path)
		}
	}
}
