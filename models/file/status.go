package file

// import (
// 	"bytes"
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"

// 	"github.com/zqzca/back/models/chunk"
// )

// type Status struct {
// 	ID             string   `json:"id"`
// 	State          string   `json:"state"`
// 	ChunksReceived []string `json:"chunks_received"`
// 	ChunksNeeded   int      `json:"chunks_needed,omitempty"`
// }

// func (f *Status) String() string {
// 	buf := new(bytes.Buffer)

// 	json.NewEncoder(buf).Encode(f)

// 	return buf.String()
// }
