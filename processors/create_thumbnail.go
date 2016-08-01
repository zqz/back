package processors

import (
	"crypto/sha1"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"

	"github.com/disintegration/imaging"
	"github.com/zqzca/back/lib"
)

func CreateThumbnail(r io.Reader) (*Thumbnail, int, error) {
	raw, format, err := image.Decode(r)

	fmt.Println("format:", format)

	if err != nil {
		fmt.Println("failed to decode image", err)
		return nil, 0, err
	}

	tmpFile, err := ioutil.TempFile("", "thumbnail")

	// Make sure we close.
	closeTmpFile := func() {
		if tmpFile {
			tmpFile.Close()
			tmpFile = nil
		}
	}

	defer closeTmpFile()

	if err != nil {
		fmt.Println("Failed to create temp file", err)
		return nil, 0, err
	}

	h := sha1.New()
	wc := new(writeCounter)
	mw := io.MultiWriter(tmpFile, h, wc)

	dst := imaging.Fill(raw, 200, 200, imaging.Center, imaging.Lanczos)
	err = imaging.Encode(mw, dst, imaging.JPEG)
	if err != nil {
		fmt.Println("failed to encode data", err)
		return nil, 0, err
	}

	// hash, err := lib.Hash(b)
	// if hash != f.Hash {
	// 	// return errors
	// }

	closeTmpFile()
	hash := h.Sum(nil)
	newPath := lib.LocalPath(hash)

	err = os.Rename(tmpFile.Name(), newPath)

	if err != nil {
		fmt.Println("Failed to rename file!?", err)
		return nil, 0, err
	}

	err = os.Chmod(newPath, 0644)
	if err != nil {
		fmt.Println("Failed to set permissions on", newPath, err)
	}

	return hash, wc, nil
}

type writeCounter int64

func (w *writeCounter) Write(b []byte) (int, error) {
	*w += int64(len(b))

	return len(b), nil
}
