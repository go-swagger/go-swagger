//go:build testintegration

package main

import (
	"log"
	"net"
	"net/http"

	"swagger/api"
)

func main() {
	// Route => handler
	http.HandleFunc("POST /foobar", api.FooBarHandler)

	// Start server
	listener, err := net.Listen("tcp", ":1323")
	if err != nil {
		log.Fatal(err)
	}

	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal(err)
	}
}
