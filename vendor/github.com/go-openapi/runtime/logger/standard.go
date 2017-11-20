package logger

import (
	"fmt"
	"os"
)

type StandardLogger struct{}

func (StandardLogger) Printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format, args...)
}

func (StandardLogger) Debugf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}
