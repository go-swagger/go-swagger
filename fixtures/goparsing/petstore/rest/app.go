package rest

import (
	"net/http"

	"github.com/go-swagger/go-swagger/fixtures/goparsing/petstore/rest/handlers"
	"github.com/naoina/denco"
)

// ServeAPI serves this api
func ServeAPI() error {
	mux := denco.NewMux()

	routes := []denco.Handler{
		mux.GET("/pets", handlers.GetPets),
		mux.POST("/pets", handlers.CreatePet),
		mux.GET("/pets/:id", handlers.GetPetByID),
		mux.PUT("/pets/:id", handlers.UpdatePet),
		mux.Handler("DELETE", "/pets/:id", handlers.DeletePet),
		mux.GET("/orders/:id", handlers.GetOrderDetails),
		mux.POST("/orders", handlers.CreateOrder),
		mux.PUT("/orders/:id", handlers.UpdateOrder),
		mux.Handler("DELETE", "/orders/:id", handlers.CancelOrder),
	}
	handler, err := mux.Build(routes)
	if err != nil {
		return err
	}
	return http.ListenAndServe(":8000", handler)
}
