// Code generated by go-swagger; DO NOT EDIT.

package store

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// OrderCreateHandlerFunc turns a function with the right signature into a order create handler
type OrderCreateHandlerFunc func(OrderCreateParams, any) middleware.Responder

// Handle executing the request and returning a response
func (fn OrderCreateHandlerFunc) Handle(params OrderCreateParams, principal any) middleware.Responder {
	return fn(params, principal)
}

// OrderCreateHandler interface for that can handle valid order create params
type OrderCreateHandler interface {
	Handle(OrderCreateParams, any) middleware.Responder
}

// NewOrderCreate creates a new http.Handler for the order create operation
func NewOrderCreate(ctx *middleware.Context, handler OrderCreateHandler) *OrderCreate {
	return &OrderCreate{Context: ctx, Handler: handler}
}

/*
	OrderCreate swagger:route POST /store/order store orderCreate

Place an order for a pet
*/
type OrderCreate struct {
	Context *middleware.Context
	Handler OrderCreateHandler
}

func (o *OrderCreate) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewOrderCreateParams()
	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		*r = *aCtx
	}
	var principal any
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
