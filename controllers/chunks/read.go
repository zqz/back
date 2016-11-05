package chunks

import "net/http"

// Read does nothing
func Read(w http.ResponseWriter, r *http.Request) {
	return
}
