package thumbnail

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"

	"github.com/disintegration/imaging"
	"github.com/zqzca/back/lib"
)

func Generate(r []byte) {
	raw, format, err := image.Decode(bytes.NewReader(r))

	fmt.Println("fmt", format)

	if err != nil {
		fmt.Println("failed to decode image", err)
		return
	}

	dst := imaging.Fill(raw, 200, 200, imaging.Center, imaging.Lanczos)

	var b bytes.Buffer
	err = imaging.Encode(&b, dst, imaging.PNG)

	if err != nil {
		fmt.Println("failed to encode data", err)
	}

	buf := bytes.NewReader(b.Bytes())
	hash, err := lib.Hash(buf)

	if err != nil {
		fmt.Println("thumbnail error:", err)
		return
	}

	path := lib.LocalPath(hash)
	ioutil.WriteFile(path, b.Bytes(), 0644)
}
