package commands

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	// initializations to run tests in this package
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

func fixtureBase() string {
	return filepath.FromSlash("../../../fixtures")
}
