package chunks

import (
	"sync"

	"github.com/zqzca/back/controllers"
)

// Controller carries dependencies
type Controller struct {
	controllers.Dependencies

	wsFileIDsLock *sync.RWMutex
	wsFileIDs     map[string]string
}

// NewController ..
func NewController(deps controllers.Dependencies) *Controller {
	c := &Controller{Dependencies: deps}
	c.wsFileIDs = make(map[string]string)
	c.wsFileIDsLock = &sync.RWMutex{}
	return c
}
