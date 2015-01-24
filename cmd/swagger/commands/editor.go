package commands

// import (
// 	"fmt"
// 	"net"
// 	"net/http"
// 	"os"
// 	"path/filepath"

// 	"github.com/jessevdk/go-flags"
// 	"github.com/rs/cors"
// )

// type editor struct {
// }

// // NewEditor creates a new doc command.
// func NewEditor() flags.Commander {
// 	return &editor{}
// }

// func (d *editor) Execute(args []string) error {
// 	fileName := "./swagger.json"
// 	if len(args) > 0 {
// 		fileName = args[0]
// 	}

// 	fi, err := os.Stat(fileName)
// 	if err != nil {
// 		return err
// 	}
// 	if fi.IsDir() {
// 		return fmt.Errorf("expected")
// 	}

// 	orig := fileName
// 	if !filepath.IsAbs(fileName) {
// 		if abs, err := filepath.Abs(fileName); err == nil {
// 			fileName = abs
// 		}
// 	}

// 	listener, err := net.Listen("tcp", "localhost:0")
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Printf("serving swagger editor for %s at http://%s/swagger-editor/#/edit\n", orig, listener.Addr())
// 	opts := cors.Options{
// 		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"},
// 		AllowedOrigins: []string{"*"},
// 	}
// 	corsMW := cors.New(opts)
// 	return http.Serve(listener, corsMW.Handler(http.FileServer(http.Dir("."))))
// }
