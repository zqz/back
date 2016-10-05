package chunks

import (
	"sync"

	"github.com/zqzca/back/controllers"
)

// Controller carries dependencies
type Controller struct {
	controllers.Dependencies

	wsFileIDs     map[string]string
	wsFileIDsLock *sync.RWMutex
}

// NewController ..
func NewController(deps controllers.Dependencies) *Controller {
	c := &Controller{Dependencies: deps}
	c.wsFileIDs = make(map[string]string)
	c.wsFileIDsLock = &sync.RWMutex{}
	return c
}
