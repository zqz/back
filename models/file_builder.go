package models

import "io"

type FileBuilder struct {
	file   *File
	chunks []Chunk
}

func NewFileBuilder(f *File) *FileBuilder {
	fb := FileBuilder{
		file: f,
	}

	if fb.fetchChunks() != nil {
		return nil
	}

	return &fb
}

func (fb *FileBuilder) fetchChunks() error {
	chunks, err := ChunksByFileID(fb.file.ID)

	if err != nil {
		return err
	}

	fb.chunks = chunks

	return nil
}

func (f *FileBuilder) Copy(w io.Writer, after func()) (int, error) {
	var bytesRead int

	for _, chunk := range f.chunks {
		f, err := chunk.Fd()
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
