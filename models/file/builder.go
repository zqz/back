package file

import (
	"fmt"
	"io"
	"os"

	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models/chunk"
)

type Builder struct {
	file   *File
	chunks []*chunk.Chunk
	ex     db.Executor
}

func NewBuilder(ex db.Executor, f *File) *Builder {
	fb := Builder{
		ex:   ex,
		file: f,
	}

	if fb.fetchChunks() != nil {
		return nil
	}

	return &fb
}

func (f *Builder) Copy(w io.Writer, after func()) (int, error) {
	var bytesRead int

	for _, c := range f.chunks {
		f, err := chunkReadCloser(c.Hash)
		defer f.Close()
		n, err := io.Copy(w, f)
		after() // Useful for flushing file out to internets

		bytesRead += int(n)

		if err != nil {
			return bytesRead, err
		}
	}

	return bytesRead, nil
}

func chunkReadCloser(hash string) (io.ReadCloser, error) {
	path := fmt.Sprintf("files/chunks/%s", hash)
	return os.Open(path)
}

func (fb *Builder) fetchChunks() error {
	chunks, err := chunk.FindByFileID(fb.ex, fb.file.ID)

	if err != nil {
		return err
	}

	fb.chunks = chunks

	return nil
}
