// Package noop provides a No-Operation logger
package noop

import (
	"fmt"

	"github.com/fclairamb/go-log"
)

// NewNoOpLogger creates a no-op logger
func NewNoOpLogger() log.Logger {
	return &noLogger{}
}

type noLogger struct{}

func (nl *noLogger) Debug(string, ...interface{}) {
}

func (nl *noLogger) Info(string, ...interface{}) {
}

func (nl *noLogger) Warn(string, ...interface{}) {
}

func (nl *noLogger) Error(string, ...interface{}) {
}

// We don't _log_ anything but we need to keep a consistent behavior

func (nl *noLogger) Panic(msg string, args ...interface{}) {
	panic(fmt.Errorf("%s: %s", msg, args)) //nolint:goerr113
}

func (nl *noLogger) With(...interface{}) log.Logger {
	return nl
}
