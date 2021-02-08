// +build ignore

package main

import (
	"log"
	"os"

	"github.com/go-openapi/runtime"
	"github.com/go-swagger/go-swagger/examples/file-server/client"
	"github.com/go-swagger/go-swagger/examples/file-server/client/uploads"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("must provide a file name as argument")
	}

	filename := os.Args[1]

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = upload(f)
	if err != nil {
		log.Fatal(err)
	}
}

func upload(reader runtime.NamedReadCloser) error {

	config := client.DefaultTransportConfig().WithHost("localhost:8000")

	uploader := client.NewHTTPClientWithConfig(nil, config)

	params := uploads.NewUploadFileParams().WithFile(reader)

	_, err := uploader.Uploads.UploadFile(params)

	return err
}
