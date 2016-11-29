package files

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zqzca/back/app"
	"github.com/zqzca/back/dependencies"
)

func TestFileStatus(t *testing.T) {
	deps := dependencies.Test()
	s := httptest.NewServer(app.Routes(deps))
	defer s.Close()

	r, err := http.Get(s.URL + "/api/v1/files")
	if err != nil {
		t.Error(err)
	}

	greeting, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", greeting)
}
