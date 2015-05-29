package operations

// ListPetParams the params for the list pets query
type ListPetParams struct {
	// OutOfStock when set to true only the pets that are out of stock will be returned
	OutOfStock bool
}

// ServeAPI serves the API for this record store
func ServeAPI(host, basePath string, schemes []string) error {

	// +swagger:route GET /pets ListPetParams pets
	//
	// lists pets filtered by some parameters.
	//
	// This will show all available pets by default.
	// You can get the pets that are out of stock
	mountItem("GET", basePath+"/pets", nil)

	// +swagger:route POST /pets createPet pets
	//
	// Create a pet based on the parameters
	mountItem("POST", basePath+"/pets", nil)

	// +swagger:route GET /orders listOrders orders
	//
	// lists pets filtered by some parameters
	mountItem("GET", basePath+"/orders", nil)

	// +swagger:route POST /orders createOrder orders
	//
	// create an order based on the parameters
	mountItem("POST", basePath+"/orders", nil)

	// +swagger:route GET /orders/{id} orderDetails orders
	//
	// gets the details for an order
	//
	// +swagger:param route:id number:int32
	// required: true
	//
	// +swagger:param query:purchaseDates array
	// items: string
	// items.min length: 10
	// min items: 2
	//
	// +swagger:response 200:order
	mountItem("GET", basePath+"/orders/:id", nil)

	return nil
}

// not really used but I need a method to decorate the calls to
func mountItem(method, path string, handler interface{}) {}
