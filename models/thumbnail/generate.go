package thumbnail

import (
	"fmt"
	"io/ioutil"

	"github.com/daddye/vips"
	"github.com/zqzca/back/lib"
)

const thumbPath = "files/thumbs/"

func Generate(r []byte, path string) {
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

	hash, err := lib.HashFile(buf)

	if err != nil {
		fmt.Println("thumbnail error:", err)
		return
	}

	path := fmt.Sprintf("%s%s", thumbPath, hash)
	ioutil.WriteFile(path, buf, 0644)
}
