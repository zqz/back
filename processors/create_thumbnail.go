package processors

import (
	"crypto/sha1"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/disintegration/imaging"
	"github.com/spf13/afero"
	"github.com/zqzca/back/controllers"
	"github.com/zqzca/back/lib"
)

func CreateThumbnail(deps controllers.Dependencies, r io.Reader) (string, int, error) {
	raw, format, err := image.Decode(r)

	// var jpgData bytes.Buffer

	// err = jpeg.Encode(jpgData, raw, &jpeg.Options{90})

	// if err != nil {
	// 	fmt.Println("failed to create a jpg")
	// 	return "", 0, err
	// }

	fmt.Println("format:", format)

	if err != nil {
		fmt.Println("failed to decode image", err)
		return "", 0, err
	}

	fs := deps.Fs
	tmpFile, err := afero.TempFile(fs, ".", "thumbnail")

	// Make sure we close.
	closeTmpFile := func() {
		if tmpFile != nil {
			tmpFile.Close()
			tmpFile = nil
		}
	}

	defer closeTmpFile()

	if err != nil {
		fmt.Println("Failed to create temp file", err)
		return "", 0, err
	}

	h := sha1.New()
	var wc writeCounter
	mw := io.MultiWriter(tmpFile, h, wc)

	dst := imaging.Fill(raw, 200, 200, imaging.Center, imaging.Lanczos)
	err = imaging.Encode(mw, dst, imaging.JPEG)
	if err != nil {
		fmt.Println("failed to encode data", err)
		return "", 0, err
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))

	deps.Debug("Thumbnail hash", "hash:", hash)
	newPath := lib.LocalPath(hash)

	err = fs.Rename(tmpFile.Name(), newPath)
	if err != nil {
		fmt.Println("Failed to rename file!?", err)

		// Todo delete file
		return hash, int(wc), err
	}

	err = fs.Chmod(newPath, 0644)
	if err != nil {
		fmt.Println("Failed to set permissions on", newPath, err)
		// Todo delete file
		return hash, int(wc), err
	}

	return hash, int(wc), nil
}

type writeCounter int64

func (w writeCounter) Write(b []byte) (int, error) {
	w += writeCounter(len(b))

	return len(b), nil
}
