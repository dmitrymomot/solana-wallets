package kitlog

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

type (
	// Logger is a logrus wrapper for go-kit/log
	Logger struct {
		*logrus.Entry
	}
)

var errMissingValue = errors.New("(MISSING)")

// NewLogger returns a new Logger instance.
// It uses logrus as a backend.
func NewLogger(l *logrus.Entry) Logger {
	return Logger{l}
}

func (l Logger) Log(keyvals ...interface{}) error {
	fields := logrus.Fields{}
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			fields[fmt.Sprint(keyvals[i])] = keyvals[i+1]
		} else {
			fields[fmt.Sprint(keyvals[i])] = errMissingValue
		}
	}
	l.WithFields(fields).Print()
	return nil
}
