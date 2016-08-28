package processors

import (
	"crypto/sha1"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/spf13/afero"
	"github.com/zqzca/back/controllers"
	"github.com/zqzca/back/lib"
)

func CreateThumbnail(deps controllers.Dependencies, r io.ReadSeeker) (string, int, error) {
	raw, format, err := image.Decode(r)

	if format == "" {
		return "", 0, nil
	}

	if format == "jpeg" || format == "jpg" {

		r.Seek(0, os.SEEK_SET)

		// 1        2       3      4         5            6           7          8

		// 888888  888888      88  88      8888888888  88                  88  8888888888
		// 88          88      88  88      88  88      88  88          88  88      88  88
		// 8888      8888    8888  8888    88          8888888888  8888888888          88
		// 88          88      88  88
		// 88          88  888888  888888

		// func Rotate180(img image.Image) *image.NRGBA
		// func Rotate270(img image.Image) *image.NRGBA
		// func Rotate90(img image.Image) *image.NRGBA
		// func FlipH(img image.Image) *image.NRGBA
		// func FlipV(img image.Image) *image.NRGBA

		fmt.Println("Rotateing")
		x, err := exif.Decode(r)
		if err != nil {
			fmt.Println("failed to decode", err)
			return "", 0, err
		}

		rawo, _ := x.Get(exif.Orientation)

		orientation, err := rawo.Int(0)

		switch orientation {
		case 1:
			// nothing;
		case 2:
			// flip Horiz L to R
		case 3:
			raw = imaging.Rotate180(raw)
			// rotate 180 ccw
		case 4:
			// flip Vert T to B
		case 5:
			// transpose
		case 6:
			// rotate 90
		case 7:
			// transverse
		case 8:
			// rotate 270
		default:
			// nothing;
		}

		fmt.Println("orientation:", orientation)
	}

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
