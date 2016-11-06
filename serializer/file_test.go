package serializer_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"
	"github.com/zqzca/back/serializer"
)

func TestForFile(t *testing.T) {
	now := time.Now()

	f := &models.File{
		Name:      "foo",
		Hash:      "123",
		Type:      "image",
		CreatedAt: now,
		Slug:      "abc",
		Size:      100,
		State:     lib.FileProcessing,
	}

	s := serializer.ForFile(nil, f)
	js := renderJSON(s)

	assert.Equal(t, "foo", js["name"])
	assert.Equal(t, "123", js["hash"])
	assert.Equal(t, "image", js["type"])
	assert.Equal(t, now.Format(time.RFC3339Nano), js["created_at"])
	assert.Equal(t, 100.0, js["size"])
}
