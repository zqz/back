package thumbnail

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/daddye/vips"
	"github.com/zqzca/back/lib"
)

func Generate(r []byte) {
	options := vips.Options{
		Width:        200,
		Height:       200,
		Crop:         true,
		Extend:       vips.EXTEND_BLACK,
		Interpolator: vips.NOHALO,
		Gravity:      vips.CENTRE,
		Quality:      95,
	}

	buf, err := vips.Resize(r, options)

	if err != nil {
		fmt.Println("thumbnail error:", err)
		return
	}

	b := bytes.NewReader(buf)
	hash, err := lib.Hash(b)

	if err != nil {
		fmt.Println("thumbnail error:", err)
		return
	}

	path := lib.LocalPath(hash)
	ioutil.WriteFile(path, buf, 0644)
}
