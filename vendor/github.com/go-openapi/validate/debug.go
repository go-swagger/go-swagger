package validate

import (
	"log"
	"os"
)

var (
	// Debug is true when the SWAGGER_DEBUG env var is not empty
	Debug = os.Getenv("SWAGGER_DEBUG") != ""
)

func debugLog(msg string, args ...interface{}) {
	// a private, trivial trace logger, based on go-openapi/spec/expander.go:debugLog()
	if Debug {
		log.Printf(msg, args...)
	}
}
