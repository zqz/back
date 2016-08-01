package lib

import (
	"os"
	"path/filepath"
)

func LocalPath(hash string) string {
	return filepath.Join("files", hash)
}

func ExistsOnDisk(hash string) bool {
	path := LocalPath(hash)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
