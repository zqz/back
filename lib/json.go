package lib

import (
	"bytes"
	"encoding/json"
)

// ToJSON Converts a struct to a JSON string.
func ToJSON(o interface{}) string {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(o)
	return buf.String()
}
