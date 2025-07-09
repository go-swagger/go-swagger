package _3125

// swagger:model Item
type Item struct {
	// Nullable value
	// required: true
	// Extensions:
	// ---
	// x-nullable: true
	Value1 *ValueStruct `json:"value1"`

	// Non-nullable value
	// required: true
	// example: {"value": 42}
	Value2 ValueStruct `json:"value2"`
}

type ValueStruct struct {
	Value int `json:"value"`
}
