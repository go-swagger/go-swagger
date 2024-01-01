package go118

// An Interfaced struct contains objects with interface definitions
type Interfaced struct {
	CustomData any `json:"custom_data"`
}

type MarshalTextMap map[string]any

func (cm MarshalTextMap) MarshalText() ([]byte, error) {
	return []byte("hola desde CustomMap"), nil
}

// swagger:type object
type SomeObjectMap any

// swagger:model namedWithType
type NamedWithType struct {
	SomeMap SomeObjectMap `json:"some_map"`
}

// SomeObject is a type that refines an untyped map
type SomeObject map[string]any

// swagger:parameters putNumPlate
type NumPlates struct {
	// in: body
	NumPlates any `json:"num_plates"`
}

// swagger:response
type NumPlatesResp struct {
	// in: body
	NumPlates any `json:"num_plates"`
}

// swagger:model transportErr
type transportErr struct {
	// Message is a human-readable description of the error.
	// Required: true
	Message string `json:"message"`
	// Data is additional data about the error.
	Data any `json:"data,omitempty"` // <- for this use case. Unsupported type "invalid type"
}
