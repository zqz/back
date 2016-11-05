package files

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"github.com/vattle/sqlboiler/boil"
	"github.com/zqzca/back/lib"
	"github.com/zqzca/back/models"

	"github.com/vattle/sqlboiler/queries/qm"
)

type fileStatus struct {
	ID             string   `json:"id"`
	State          string   `json:"state"`
	ChunksReceived []string `json:"chunks_received,omitempty"`
	ChunksNeeded   int      `json:"chunks_needed,omitempty"`
	Slug           string   `json:"slug,omitempty"`
}

// FileStatus represents the files... status...
type fileState int

func (s fileState) String() string {
	switch s {
	case lib.FileIncomplete:
		return "incomplete"
	case lib.FileProcessing:
		return "processing"
	case lib.FileFinished:
		return "finished"
	default:
		return "unknown"
	}
}

func statusForFile(ex boil.Executor, f *models.File) (*fileStatus, error) {
	chunksNeeded := 0
	chunks, err := models.Chunks(ex, qm.Where("file_id=$1", f.ID)).All()
	if err != nil {
		return nil, err
	}

	numChunks := len(chunks)
	chunksReceived := []string{}

	if int(f.State) == lib.FileIncomplete {
		chunksNeeded = f.NumChunks - numChunks

		for _, c := range chunks {
			chunksReceived = append(chunksReceived, c.Hash)
		}
	}

	state := fileState(f.State)
	return &fileStatus{
		ID:             f.ID,
		State:          state.String(),
		ChunksReceived: chunksReceived,
		ChunksNeeded:   chunksNeeded,
		Slug:           f.Slug,
	}, nil
}

// Status returns JSON with the current state of the file.
func (f Controller) Status(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	if len(hash) == 0 {
		http.Error(w, "Hash not specified", 500)
		return
	}

	file, err := models.Files(f.DB, qm.Where("hash=$1", hash)).One()
	if err != nil {
		f.Debug("Failed to find file with hash", "hash", hash)
		http.Error(w, "", http.StatusNoContent)
		return
	}

	fs, err := statusForFile(f.DB, file)

	if err != nil {
		f.Info("Failed to fetch status for file", "err", err)
		http.Error(w, "Failed to fetch status for file", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, fs)
}
