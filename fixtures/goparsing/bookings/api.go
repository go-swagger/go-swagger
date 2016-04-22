// Package booking API.
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
//
//     Schemes: https
//     Host: localhost
//     Version: 0.0.1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//
// swagger:meta
package booking

import (
	"net/http"

	"github.com/go-swagger/scan-repo-boundary/makeplans"
)

// Customer of the site.
//
// swagger:model Customer
type Customer struct {
	Name string `json:"name"`
}

// IgnoreMe should not be added to definitions since it is not annotated.
type IgnoreMe struct {
	Name string `json:"name"`
}

// DateRange represents a scheduled appointments time
// DateRange should be in definitions since it's being used in a response
type DateRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// BookingResponse represents a scheduled appointment
//
// swagger:response BookingResponse
type BookingResponse struct {
	// Booking struct
	//
	// in: body
	// required: true
	Body struct {
		Booking  makeplans.Booking `json:"booking"`
		Customer Customer          `json:"customer"`
		Dates    DateRange         `json:"dates"`
	}
}

// Bookings swagger:route GET /admin/bookings/ booking Bookings
//
// Bookings lists all the appointments that have been made on the site.
//
//
// Consumes:
// application/json
//
// Schemes: http, https
//
// Produces:
// application/json
//
// Responses:
// 200: BookingResponse
func bookings(w http.ResponseWriter, r *http.Request) {

}
