package swagger

import (
	"errors"
	"net/http"
)

// InitializeRouter takes the untyped API and a router
// with those it will validate the registrations in the API.
// If there are missing consumers for registered media types it will return an error
// If there are missing producers for registered media types it will return an error
// If there are missing auth handlers for registered security schemes it will return an error
// If there are missing operation handlers for operationIds it will return an error
func InitializeRouter(api *API, router Router) (http.Handler, error) {
	if api == nil {
		return nil, errors.New("the api to serve can't be nil, but it was")
	}

	if router == nil {
		router = DefaultRouter()
	}

	analyzer := NewAnalyzer(api.Spec())

	if err := api.ValidateWith(analyzer); err != nil {
		return nil, err
	}

	return initializeRoutes(api, router)
}

func initializeRoutes(api *API, router Router) (http.Handler, error) {
	for path, pathItem := range api.Spec().Paths.Paths {
		if op := pathItem.Get; op != nil {
			if h := api.OperationHandlerFor(op.ID); h != nil {
				router.AddRoute("GET", path, createHandler(h))
			}
		}
		if op := pathItem.Head; op != nil {
			if h := api.OperationHandlerFor(op.ID); h != nil {
				router.AddRoute("HEAD", path, createHandler(h))
			}
		}
		if op := pathItem.Options; op != nil {
			if h := api.OperationHandlerFor(op.ID); h != nil {
				router.AddRoute("OPTIONS", path, createHandler(h))
			}
		}
		if op := pathItem.Put; op != nil {
			if h := api.OperationHandlerFor(op.ID); h != nil {
				router.AddRoute("PUT", path, createHandler(h))
			}
		}
		if op := pathItem.Post; op != nil {
			if h := api.OperationHandlerFor(op.ID); h != nil {
				router.AddRoute("POST", path, createHandler(h))
			}
		}
		if op := pathItem.Patch; op != nil {
			if h := api.OperationHandlerFor(op.ID); h != nil {
				router.AddRoute("PATCH", path, createHandler(h))
			}
		}
		if op := pathItem.Delete; op != nil {
			if h := api.OperationHandlerFor(op.ID); h != nil {
				router.AddRoute("DELETE", path, createHandler(h))
			}
		}
	}
	return router.Build()
}

func createHandler(h *OperationHandler) HandlerFunc {
	// make request authorizer
	// make request binder
	// make request validator
	// get result
	// render response
	return func(rw http.ResponseWriter, r *http.Request, routeParams RouteParams) {
		rw.WriteHeader(http.StatusNotImplemented)
	}
}
