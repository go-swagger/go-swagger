package rest

import (
	"net/http"

	"github.com/casualjim/go-swagger/fixtures/goparsing/petstore/rest/handlers"
	"github.com/naoina/denco"
)

// ServeAPI serves this api
func ServeAPI() error {
	mux := denco.NewMux()

	routes := []denco.Handler{
		// +swagger:route GET /pets pets listPets
		//
		// Lists the pets known to the store.
		//
		// Responses:
		// default: genericError
		// 200: []pets
		mux.GET("/pets", handlers.GetPets),

		// +swagger:route POST /pets pets createPet
		//
		// Creates a new pet in the store.
		mux.POST("/pets", handlers.CreatePet),

		// +swagger:route GET /pets/{id} pets getPetById
		//
		// Gets the details for a pet.
		mux.GET("/pets/:id", handlers.GetPetByID),

		// +swagger:route GET /pets/{id} pets updatePet
		//
		// Updates the details for a pet.
		mux.PUT("/pets/:id", handlers.UpdatePet),

		// +swagger:route DELETE /pets/{id} pets deletePet
		//
		// Deletes the pet from the store.
		mux.Handler("DELETE", "/pets/:id", handlers.DeletePet),

		// +swagger:route GET /orders/{id}
		//
		// Gets the details for the specified order
		mux.GET("/orders/:id", handlers.GetOrderDetails),

		// +swagger:route DELETE /orders/{id}
		//
		// Cancels an order.
		mux.Handler("DELETE", "/orders/:id", handlers.CancelOrder),
	}
	handler, err := mux.Build(routes)
	if err != nil {
		return err
	}
	return http.ListenAndServe(":8000", handler)
}
