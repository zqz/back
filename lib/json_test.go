package lib_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/lib"
)

type Foobar struct {
	A int `json:"a"`
	B int `json:"c"`
}

func TestToJSON(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	x := &Foobar{1, 3}
	a.Equal(`{"a":1,"c":3}`, strings.TrimSpace(lib.ToJSON(x)))
}
