package processors

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/zqzca/back/controller"
	"github.com/zqzca/back/models"
)

func chunksExist(hashes []string) (bool, error) {
	for _, hash := range hashes {
		path := filepath.Join("files", "chunks", hash)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			return false, err
		}
	}

	return true, nil
}

// BuildFile builds a file from chunks.
func BuildFile(deps controller.Dependencies, f *models.File) (io.ReadSeeker, error) {
	chunks, err := models.Chunks(
		deps.DB,
		qm.Where("file_id=$1", f.ID),
		qm.OrderBy("position asc"),
	).All()

	var hashes []string
	for _, c := range chunks {
		hashes = append(hashes, c.Hash)
	}

	if ok, err := chunksExist(hashes); !ok {
		deps.Error("Missing chunks!", "id", f.ID, "name", f.Name)
		return nil, err
	}

	if err != nil {
		fmt.Println("Failed to find chunks for file:", f.ID)
		return nil, err
	}

	fs := deps.Fs
	fullFilePath := filepath.Join("files", f.Hash)
	fullFile, err := fs.Create(fullFilePath)
	if err != nil {
		fmt.Println("Failed because", err)
		return nil, err
	}
	defer fullFile.Close()
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

	// for _, c := range chunks {
	// 	path := filepath.Join("files", "chunks", c.Hash)
	// 	if err != nil {
	// 		fmt.Println("Failed to delete chunk entry:", c.ID)
	// 	}
	// 	fs.Remove(path)
	// }

	bs := bytes.NewReader(fullFileBuffer.Bytes())

	return bs, nil
}
