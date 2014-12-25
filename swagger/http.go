package swagger

import "net/http"

// ServeAPI takes the untyped API and a router
// with those it will validate the registrations in the API.
// If there are missing consumers for registered media types it will return an error
// If there are missing producers for registered media types it will return an error
// If there are missing auth handlers for registered security schemes it will return an error
// If there are missing operation handlers for operationIds it will return an error
func ServeAPI(api *API, router Router) (http.Handler, error) {
	if router == nil {
		router = DefaultRouter()
	}
	// validate the api registrations against the swagger spec
	if err := validateAPIRegistrations(api); err != nil {
		return nil, err
	}

	// iterate over all the operations and build a handler func for each
	// and add the operation as a route to the router
	if err := convertOperations(api, router); err != nil {
		return nil, err
	}

	// invoke the build method on the router
	return router.Build()
}

func validateAPIRegistrations(api *API) error {
	// check all the consumers
	// check all the producers
	// check all the security schemes
	// check all the operations
	return nil
}

func convertOperations(api *API, router Router) error {
	// for each operation build a handler func that binds, validates and executes requests
	// register the handler func with a method and a swagger path pattern in the router
	return nil
}
