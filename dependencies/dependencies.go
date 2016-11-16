package dependencies

import (
	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/afero"
)

// WebsocketClientWriter can send a message to a client
type WebsocketClientWriter interface {
	WriteClient(string, string, interface{})
	Broadcast(string, interface{})
}

// Dependencies are used throughout the app.
type Dependencies struct {
	*logrus.Logger
	*sqlx.DB
	afero.Fs
	WS WebsocketClientWriter
}

// New dependencies for non test
func New() Dependencies {
	return Dependencies{}
}

// Test dependencies
func Test() Dependencies {
	return Dependencies{
		Fs:     afero.NewMemMapFs(),
		Logger: nil,
		DB:     nil,
		WS:     nil,
	}
}
