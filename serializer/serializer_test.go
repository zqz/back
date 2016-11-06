package serializer_test

import (
	"encoding/json"

	"github.com/zqzca/back/db"
	"github.com/zqzca/back/models"
	"github.com/zqzca/back/serializer"
)

func init() {
	serializer.FileDownloads = func(_ db.Executor, _ *models.File) int {
		return 1
	}
}

func renderJSON(d interface{}) map[string]interface{} {
	var output map[string]interface{}
	b, _ := json.Marshal(d)
	_ = json.Unmarshal(b, &output)

	return output
}
