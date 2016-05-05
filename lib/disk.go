package lib

import (
	"os"
	"path"
)

func LocalPath(hash string) string {
	return path.Join("files", hash)
}

func ExistsOnDisk(hash string) bool {
	path := LocalPath(hash)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
