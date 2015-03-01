package commands

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/casualjim/go-swagger/swagger-ui"
	"github.com/jessevdk/go-flags"
)

type ui struct {
}

// NewUI creates a new doc command.
func NewUI() flags.Commander {
	return &ui{}
}

func (d *ui) Execute(args []string) error {
	fileName := "./swagger.json"
	if len(args) > 0 {
		fileName = args[0]
	}

	fi, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return fmt.Errorf("%s is a directory, expected a file", fileName)
	}

	orig := fileName
	if !filepath.IsAbs(fileName) {
		if abs, err := filepath.Abs(fileName); err == nil {
			fileName = abs
		}
	}

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return err
	}

	fmt.Printf("serving swagger ui for %s at http://%s/swagger-ui\n", orig, listener.Addr())

	return http.Serve(listener, swaggerui.Middleware(fileName, nil))
}
