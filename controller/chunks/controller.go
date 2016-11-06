package chunks

import (
	"sync"

	"github.com/zqzca/back/dependencies"
)

// Controller carries dependencies
type Controller struct {
	dependencies.Dependencies

	wsFileIDsLock *sync.RWMutex
	wsFileIDs     map[string]string
}

// NewController ..
func NewController(deps dependencies.Dependencies) *Controller {
	c := &Controller{Dependencies: deps}
	c.wsFileIDs = make(map[string]string)
	c.wsFileIDsLock = &sync.RWMutex{}
	return c
}
