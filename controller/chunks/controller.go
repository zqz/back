package chunks

import (
	"sync"

	"github.com/zqzca/back/controller"
)

// Controller carries dependencies
type Controller struct {
	controller.Dependencies

	wsFileIDsLock *sync.RWMutex
	wsFileIDs     map[string]string
}

// NewController ..
func NewController(deps controller.Dependencies) *Controller {
	c := &Controller{Dependencies: deps}
	c.wsFileIDs = make(map[string]string)
	c.wsFileIDsLock = &sync.RWMutex{}
	return c
}
