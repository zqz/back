package lib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/lib"
)

func TestLocalPath(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	a.Equal("files/foo", lib.LocalPath("foo"))
}
