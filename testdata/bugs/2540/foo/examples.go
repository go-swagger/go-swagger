// Package docs SwagTest
//
// Test Swagger
//
//	Schemes: https
//	Version: 1.0.0
//	BasePath: /test/
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package foo

// swagger:route GET /book  book
// Get book
// responses:
//   200: getBook

type Author struct {
	Name string
}

// Book holds all relevant information about a book.
//
// At this moment, a book is only described by its publishing date
// and author.
//
// example: { "Published": 2026, "Author": "Fred" }
//
// default: { "Published": 1900, "Author": "Unknown" }
//
// swagger:model
type Book struct {
	// min: 0
	//
	// example: 2021
	Published int
	// example: { "Name": "Tolkien" }
	Author Author
}

// OK.
// swagger:response getBook
type response struct {
	// in:body
	Body Book
}
