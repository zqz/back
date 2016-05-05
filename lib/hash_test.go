package lib_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/lib"
)

func TestHash(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	b := bytes.NewReader([]byte("boo"))

	h, err := lib.Hash(b)
	a.Nil(err)
	a.Equal("78b371f0ea1410abc62ccb9b7f40c34288a72e1a", h)
}
