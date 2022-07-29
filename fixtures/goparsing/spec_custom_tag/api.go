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
package spec_custom_tag

import (
	"github.com/go-swagger/go-swagger/fixtures/goparsing/spec_custom_tag/makeplans"
	"net/http"
)

// Customer of the site.
//
// swagger:model Customer
type Customer struct {
	Name string `swagger:"name"`
}

// IgnoreMe should not be added to definitions since it is not annotated.
type IgnoreMe struct {
	Name string `swagger:"name"`
}

// DateRange represents a scheduled appointments time
// DateRange should be in definitions since it's being used in a response
type DateRange struct {
	Start string `swagger:"start"`
	End   string `swagger:"end"`
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
		Booking  makeplans.Booking `swagger:"booking"`
		Customer Customer          `swagger:"customer"`
		Dates    DateRange         `swagger:"dates"`
		// example: {"key": "value"}
		Map map[string]string `swagger:"map"`
		// example: [1, 2]
		Slice []int `swagger:"slice"`
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
// Deprecated: true
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
