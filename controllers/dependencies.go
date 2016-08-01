package controllers

import (
	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
)

// Dependencies for each controller. This allows us to provide things like
// loggers and database handles. Also makes it easy to test.
type Dependencies struct {
	*logrus.Logger
	*sqlx.DB
}

// Info provides log15 api over logrus
func (d Dependencies) Info(msg string, args ...interface{}) {
	if len(args) == 0 {
		d.Logger.Info(msg)
	}
	d.Logger.WithFields(logConvert(args)).Info(msg)
}

// Warn provides log15 api over logrus
func (d Dependencies) Warn(msg string, args ...interface{}) {
	if len(args) == 0 {
		d.Logger.Warn(msg)
	}
	d.Logger.WithFields(logConvert(args)).Warn(msg)
}

// Debug provides log15 api over logrus
func (d Dependencies) Debug(msg string, args ...interface{}) {
	if len(args) == 0 {
		d.Logger.Debug(msg)
	}
	d.Logger.WithFields(logConvert(args)).Debug(msg)
}

// Error provides log15 api over logrus
func (d Dependencies) Error(msg string, args ...interface{}) {
	if len(args) == 0 {
		d.Logger.Error(msg)
	}
	d.Logger.WithFields(logConvert(args)).Error(msg)
}

func logConvert(intf []interface{}) map[string]interface{} {
	if len(intf)%2 != 0 {
		panic("structured logging requires an even amount of key to values")
	}

	m := make(map[string]interface{})

	for i := 0; i < len(intf); i += 2 {
		keyIntf := intf[i]
		valIntf := intf[i+1]

		key, ok := keyIntf.(string)
		if !ok {
			panic("structured logging requires string keys")
		}

		m[key] = valIntf
	}

	return m
}
