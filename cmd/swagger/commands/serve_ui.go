package commands

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/rs/cors"
	"github.com/toqueteos/webbrowser"
)

//ServeUI Show given swagger spec in swagger-ui
type ServeUI struct {
}

// Execute this command
func (s *ServeUI) Execute(args []string) error {

	if len(args) == 0 {
		return errors.New("The serve-ui command requires the swagger document url to be specified")
	}
	swaggerDoc := args[0]
	specDoc, err := spec.Load(swaggerDoc)
	if err != nil {
		log.Fatalln(err)
	}
	serveUI(specDoc)
	return nil
}

func serveUI(doc *spec.Document) error {

	// parse the url and open in default browser
	u, err := url.Parse("http://petstore.swagger.io/")
	if err != nil {
		return err
	}
	q := u.Query()
	q.Add("url", "http://localhost:8080/swagger.json")
	u.RawQuery = q.Encode()
	webbrowser.Open(u.String())

	// serve the swagger.json
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(doc.Raw()))
	})

	handler := cors.Default().Handler(mux)

	http.ListenAndServe(":8080", handler)

	return nil
}
