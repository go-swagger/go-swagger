package handlers

import (
	"net/http"

	"github.com/naoina/denco"
)

// GetOrderDetails gets the details for an order
func GetOrderDetails(rw http.ResponseWriter, req *http.Request, params denco.Params) {}

// CancelOrder deletes an order
func CancelOrder(rw http.ResponseWriter, req *http.Request, params denco.Params) {}
