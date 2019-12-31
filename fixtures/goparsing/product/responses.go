package product

// swagger:response GetProductsResponse
type GetProductsResponse struct {
	// in:body
	Body map[string]Product `json:"body"`
}
