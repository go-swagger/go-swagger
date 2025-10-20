package transparentalias

// ResponseEnvelope is the canonical struct referenced by aliases in responses.
type ResponseEnvelope struct {
	// Payload uses an alias that should resolve transparently.
	Payload TransparentPayloadAlias `json:"payload"`
}

// ResponseEnvelopeAlias is an exported alias to ResponseEnvelope.
type ResponseEnvelopeAlias = ResponseEnvelope

type transparentAliasResponse struct {
	// Body exercises alias handling for response bodies.
	//
	// in: body
	Body ResponseEnvelopeAlias `json:"body"`
}

// TransparentAliasResponse is an exported alias annotated as swagger response.
//
// swagger:response transparentAliasResponse
type TransparentAliasResponse = transparentAliasResponse
