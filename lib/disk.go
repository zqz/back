package lib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/satori/go.uuid"
)

func LocalPath(hash string) string {
	return filepath.Join("files", hash)
}

func TempFilePath(prefix string) string {
	uid := uuid.NewV4()
	return LocalPath(fmt.Sprintf("%s-%s", prefix, uid.String()))
}

func ExistsOnDisk(hash string) bool {
	path := LocalPath(hash)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
