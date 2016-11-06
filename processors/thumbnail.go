package processors

import (
	"crypto/sha1"
	"fmt"
	"image"
	_ "image/gif"  // GIF Support
	_ "image/jpeg" // JPG Support
	_ "image/png"  // PNG Support
	"io"
	"os"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/zqzca/back/dependencies"
	"github.com/zqzca/back/lib"
)

func rotate(img image.Image, orientation int) image.Image {
	fmt.Println("orientation:", orientation)
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

	var out image.Image
	switch orientation {
	case 1:
		out = img
		// nothing;
	case 2:
		out = imaging.FlipH(img)
		// flip Horiz L to R
	case 3:
		out = imaging.Rotate180(img)
		// rotate 180 ccw
	case 4:
		out = imaging.FlipV(img)
		// flip Vert T to B
	case 5:
		out = imaging.Transpose(img)
		// transpose
	case 6:
		out = imaging.Rotate90(img)
		// rotate 90
	case 7:
		out = imaging.Transverse(img)
		// transverse
	case 8:
		out = imaging.Rotate270(img)
		// rotate 270
	default:
		out = img
		// nothing;
	}

	return out
}

func readOrientation(r io.ReadSeeker) (int, error) {
	_, err := r.Seek(0, os.SEEK_SET)

	if err != nil {
		fmt.Println("Failed to seek to beginning of stream")
		return 0, err
	}

	x, err := exif.Decode(r)

	if err != nil {
		fmt.Println("Failed to decode EXIF", err)
		return 0, err
	}

	orientationData, err := x.Get(exif.Orientation)

	if err != nil {
		fmt.Println("Failed to read orientation property")
		return 0, err
	}

	orientation, err := orientationData.Int(0)

	if err != nil {
		fmt.Println("Failed to decode orientation")
		return 0, err
	}

	return orientation, nil
}

// CreateThumnail builds a JPG thumbnail and can rotate if an exif bit is set.
func CreateThumbnail(deps dependencies.Dependencies, r io.ReadSeeker) (string, int, error) {
	raw, format, err := image.Decode(r)

	if format == "" {
		return "", 0, nil
	}

	if format == "jpeg" || format == "jpg" {
		deps.Debug("Received JPG")
		orientation, err := readOrientation(r)

		if err == nil {
			deps.Debug("Rotating JPG", "orientation", orientation)
			raw = rotate(raw, orientation)
		}
	}

	deps.Debug("Thumbnail format", "fmt", format)

	if err != nil {
		deps.Error("Failed to decode image")
		return "", 0, err
	}

	fs := deps.Fs
	tmpFilePath := lib.TempFilePath("thumbnail")
	tmpFile, err := fs.Create(tmpFilePath)
	if err != nil {
		deps.Error("Failed to create temp file", "path", tmpFilePath)
		return "", 0, err
	}

	// Make sure we close.
	closeTmpFile := func() {
		if tmpFile != nil {
			tmpFile.Close()
			tmpFile = nil
		}
	}

	defer closeTmpFile()

	h := sha1.New()
	var wc writeCounter
	mw := io.MultiWriter(tmpFile, h, wc)

	// Generate Thumbnail image data
	dst := imaging.Fill(raw, 200, 200, imaging.Center, imaging.Lanczos)
	// Write it
	err = imaging.Encode(mw, dst, imaging.JPEG)
	if err != nil {
		deps.Error("Failed to encode data")
		return "", 0, err
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))
	deps.Debug("Thumbnail hash", "hash:", hash)
	newPath := lib.LocalPath(hash)

	// Move temp thumbnail to final destination.
	err = os.Rename(tmpFilePath, newPath)
	if err != nil {
		deps.Error("Failed to rename file")

		// Todo delete file
		return hash, int(wc), err
	}

	// Set permissons
	err = fs.Chmod(newPath, 0644)
	if err != nil {
		deps.Error("Failed to set permissions", "path", newPath)
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
