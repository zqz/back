package processors

import (
	"fmt"
	"path/filepath"

	"github.com/zqzca/back/controllers"
	"github.com/zqzca/back/models"
)

// Cleanup removes old chunks
func Cleanup(deps controllers.Dependencies, f *models.File) error {
	chunks, err := f.Chunks(deps.DB).All()

	if err != nil {
		fmt.Println("Failed to lookup chunks for file", f.ID)
		return nil
	}

	fs := deps.Fs
	for _, c := range chunks {
		path := filepath.Join("files", "chunks", c.Hash)
		fmt.Println("Removing chunk at", path)
		if err != nil {
			fmt.Println("Failed to delete chunk entry:", c.ID)
		}
		fs.Remove(path)
	}

	return nil
}
