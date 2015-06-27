package main

import (
	"log"
	"net/http"

	"github.com/go-swagger/go-swagger/examples/2.0/petstore/server/api"
)

func main() {
	petstoreAPI, err := api.NewPetstore()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Serving petstore api on http://127.0.0.1:8344/swagger-ui/")
	http.ListenAndServe(":8344", petstoreAPI)
}
