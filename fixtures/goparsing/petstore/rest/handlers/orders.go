package handlers

import (
	"net/http"

	"github.com/go-swagger/go-swagger/fixtures/goparsing/petstore/models"
	"github.com/naoina/denco"
)

// An OrderID parameter model.
//
// This is used for operations that want the ID of an order in the path
// swagger:parameters getOrderDetails cancelOrder updateOrder
type OrderID struct {
	// The ID of the order
	//
	// in: path
	// required: true
	ID int64 `json:"id"`
}

// An OrderBodyParams model.
//
// This is used for operations that want an Order as body of the request
// swagger:parameters updateOrder createOrder
type OrderBodyParams struct {
	// The order to submit
	//
	// required: true
	// in: body
	Order *models.Order `json:"order"`
}

// GetOrderDetails swagger:route GET /orders/{id} orders getOrderDetails
//
// Gets the details for an order.
//
// Responses:
//    default: genericError
//        200: order
func GetOrderDetails(rw http.ResponseWriter, req *http.Request, params denco.Params) {
	// some actual stuff should happen in here
}

// CancelOrder swagger:route DELETE /orders/{id} orders cancelOrder
//
// Deletes an order.
//
// Responses:
//    default: genericError
//        204:
func CancelOrder(rw http.ResponseWriter, req *http.Request, params denco.Params) {
	// some actual stuff should happen in here
}

// UpdateOrder swagger:route PUT /orders/{id} orders updateOrder
//
// Updates an order.
//
// Responses:
//    default: genericError
//        200: order
//        422: validationError
func UpdateOrder(rw http.ResponseWriter, req *http.Request, params denco.Params) {
	// some actual stuff should happen in here
}

// CreateOrder swagger:route POST /orders orders createOrder
//
// Creates an order.
//
// Responses:
//    default: genericError
//        200: order
//        422: validationError
func CreateOrder(rw http.ResponseWriter, req *http.Request, params denco.Params) {
	// some actual stuff should happen in here
}
