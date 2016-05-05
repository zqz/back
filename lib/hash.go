package lib

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func Hash(src io.ReadSeeker) (string, error) {
	h := sha1.New()

	if _, err := io.Copy(h, src); err != nil {
		fmt.Println(err)
		return "", err
	}

	if _, err := src.Seek(0, os.SEEK_SET); err != nil {
		fmt.Println(err)
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
