package transparentalias

// TransparentPayload is the canonical struct referenced by aliases in tests.
//
// swagger:model TransparentPayload
type TransparentPayload struct {
	// ID of the payload.
	//
	// required: true
	ID int64 `json:"id"`

	// Name of the payload.
	Name string `json:"name"`
}

// TransparentPayloadAlias is an exported alias to TransparentPayload.
type TransparentPayloadAlias = TransparentPayload

// QueryValue is the base type used for aliasing query parameters.
type QueryValue string

// QueryValueAlias is an exported alias to QueryValue.
type QueryValueAlias = QueryValue

type transparentAliasParams struct {
	// AliasBody exercises alias handling for body parameters.
	//
	// in: body
	// required: true
	AliasBody TransparentPayloadAlias `json:"aliasBody"`

	// AliasQuery exercises alias handling for non-body parameters.
	//
	// in: query
	AliasQuery QueryValueAlias `json:"aliasQuery"`
}

// TransparentAliasParams is an exported alias annotated as swagger parameters.
//
// swagger:parameters transparentAlias
type TransparentAliasParams = transparentAliasParams
