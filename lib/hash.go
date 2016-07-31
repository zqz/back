package lib

import (
	"crypto/sha1"
	"fmt"
	"io"
)

func Hash(src io.Reader) (string, error) {
	h := sha1.New()

	if _, err := io.Copy(h, src); err != nil {
		fmt.Println(err)
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
