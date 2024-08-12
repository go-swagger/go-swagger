package pointers_nullable_by_default

// swagger:model Item
type Item struct {
	Value1 *int

	Value2 int

	// Extensions:
	// ---
	// x-nullable: false
	Value3 *int

	// Extensions:
	// ---
	// x-isnullable: false
	Value4 *int

	Value5 *int `json:"Value5,omitempty"`
}

// swagger:model ItemInterface
type ItemInterface interface {
	Value1() *int
	Value2() int

	// Value3 is a nullable value
	// Extensions:
	// ---
	// x-nullable: false
	Value3() *int

	// Value4 is a non-nullable value
	// Extensions:
	// ---
	// x-isnullable: false
	Value4() *int

	Value5() int
}
