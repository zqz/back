package chunks

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/zqzca/back/lib"
)

const maxChunkSize = 5 * 1024 * 1024

type upload struct {
	chunkID    int
	fileID     string
	size       int
	wsID       string
	remoteHash string
	localHash  string
	data       []byte

	request *http.Request
}

func (u *upload) validData() (bool, error) {
	if u.remoteHash != u.localHash {
		return false, errors.New("Hash does not match")
	}

	if u.size != len(u.data) {
		return false, errors.New("Incorrect size given")
	}

	return true, nil
}

func (u *upload) validRequest() (bool, error) {
	if u.size == 0 {
		return false, errors.New("Chunk has no size")
	}

	if u.size > maxChunkSize {
		return false, errors.New("Chunk is too big")
	}

	if len(u.fileID) == 0 {
		return false, errors.New("No FileID specified")
	}

	if u.chunkID < 0 {
		return false, errors.New("No ChunkID specified")
	}

	if len(u.remoteHash) == 0 {
		return false, errors.New("No Hash specified")
	}

	return true, nil
}

func (u *upload) loadData() error {
	var err error
	u.data, err = ioutil.ReadAll(u.request.Body)
	return err
}

func (u *upload) hashData() error {
	br := bytes.NewReader(u.data)
	var err error
	u.localHash, err = lib.Hash(br)

	if err != nil {
		return err
	}

	_, err = br.Seek(0, os.SEEK_SET)
	return err
}

func (u *upload) storeData() error {
	dst, err := os.Create(u.localPath())
	if err != nil {
		return err
	}
	defer dst.Close()

	br := bytes.NewReader(u.data)
	if _, err := io.Copy(dst, br); err != nil {
		return errors.Wrap(err, "Failed to copy chunk data to dst")
	}

	return nil
}

func (u *upload) localPath() string {
	return filepath.Join("files", "chunks", u.localHash)
}
