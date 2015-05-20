package operations

// ListPetParams the params for the list pets query
type ListPetParams struct {
	// OutOfStock when set to true only the pets that are out of stock will be returned
	OutOfStock bool
}

// ServeAPI serves the API for this record store
func ServeAPI(host, basePath string, schemes []string) error {

	// +swagger:route GET /pets pets:listPets ListPetParams
	//
	// lists pets filtered by some parameters
	mountItem("GET", basePath+"/pets", nil)

	// +swagger:route POST /pets pets:createPet CreatePetParams
	//
	// Create a pet based on the parameters
	mountItem("POST", basePath+"/pets", nil)

	// +swagger:route GET /orders orders:listOrders
	//
	// lists pets filtered by some parameters
	mountItem("GET", basePath+"/orders", nil)

	// +swagger:route POST /orders orders:createOrder
	//
	// create an order based on the parameters
	mountItem("POST", basePath+"/orders", nil)

	return nil
}

// not really used but I need a method to decorate the calls to
func mountItem(method, path string, handler interface{}) {}
