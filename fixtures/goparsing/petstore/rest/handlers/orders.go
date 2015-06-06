package handlers

import (
	"net/http"

	"github.com/naoina/denco"
)

// GetOrderDetails gets the details for an order
//
// +swagger:route GET /orders/{id} orders getOrderDetails
func GetOrderDetails(rw http.ResponseWriter, req *http.Request, params denco.Params) {}

// CancelOrder deletes an order
//
// +swagger:route DELETE /orders/{id} orders cancelOrder
func CancelOrder(rw http.ResponseWriter, req *http.Request, params denco.Params) {}

// UpdateOrder updates an order
//
// +swagger:route PUT /orders/{id} orders updateOrder
func UpdateOrder(rw http.ResponseWriter, req *http.Request, params denco.Params) {}

// CreateOrder creates an order
//
// +swagger:route POST /orders orders createOrder
func CreateOrder(rw http.ResponseWriter, req *http.Request, params denco.Params) {}
