package codescan

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestRouteExpression(t *testing.T) {
	assert.Regexp(t, rxRoute, "swagger:route DELETE /orders/{id} deleteOrder")
	assert.Regexp(t, rxRoute, "swagger:route GET /v1.2/something deleteOrder")
}

func TestRoutesParser(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	var ops spec.Paths
	for _, apiPath := range sctx.app.Routes {
		prs := &routesBuilder{
			ctx:        sctx,
			route:      apiPath,
			operations: make(map[string]*spec.Operation),
		}
		require.NoError(t, prs.Build(&ops))
	}

	assert.Len(t, ops.Paths, 3)

	po, ok := ops.Paths["/pets"]
	assert.True(t, ok)
	assert.NotNil(t, po.Get)
	assertOperation(t,
		po.Get,
		"listPets",
		"Lists pets filtered by some parameters.",
		"This will show all available pets by default.\nYou can get the pets that are out of stock",
		[]string{"pets", "users"},
		[]string{"read", "write"},
	)
	assertOperation(t,
		po.Post,
		"createPet",
		"Create a pet based on the parameters.",
		"",
		[]string{"pets", "users"},
		[]string{"read", "write"},
	)

	po, ok = ops.Paths["/orders"]
	assert.True(t, ok)
	assert.NotNil(t, po.Get)
	assertOperation(t,
		po.Get,
		"listOrders",
		"lists orders filtered by some parameters.",
		"",
		[]string{"orders"},
		[]string{"orders:read", "https://www.googleapis.com/auth/userinfo.email"},
	)
	assertOperation(t,
		po.Post,
		"createOrder",
		"create an order based on the parameters.",
		"",
		[]string{"orders"},
		[]string{"read", "write"},
	)

	po, ok = ops.Paths["/orders/{id}"]
	assert.True(t, ok)
	assert.NotNil(t, po.Get)
	assertOperation(t,
		po.Get,
		"orderDetails",
		"gets the details for an order.",
		"",
		[]string{"orders"},
		[]string{"read", "write"},
	)

	assertOperation(t,
		po.Put,
		"updateOrder",
		"Update the details for an order.",
		"When the order doesn't exist this will return an error.",
		[]string{"orders"},
		[]string{"read", "write"},
	)

	assertOperation(t,
		po.Delete,
		"deleteOrder",
		"delete a particular order.",
		"",
		nil,
		[]string{"read", "write"},
	)

	// additional check description tag at Responses
	rsp, ok := po.Delete.Responses.StatusCodeResponses[202]
	assert.True(t, ok)
	assert.Equal(t, "Some description", rsp.Description)
	assert.Equal(t, "", rsp.Ref.String())
}

func TestRoutesParserBody(t *testing.T) {
	sctx, err := newScanCtx(&Options{
		Packages: []string{
			"github.com/go-swagger/go-swagger/fixtures/goparsing/classification",
			"github.com/go-swagger/go-swagger/fixtures/goparsing/classification/models",
			"github.com/go-swagger/go-swagger/fixtures/goparsing/classification/operations",
			"github.com/go-swagger/go-swagger/fixtures/goparsing/classification/operations_body",
		},
	})
	require.NoError(t, err)
	var ops spec.Paths
	for _, apiPath := range sctx.app.Routes {
		prs := &routesBuilder{
			ctx:        sctx,
			route:      apiPath,
			operations: make(map[string]*spec.Operation),
		}
		require.NoError(t, prs.Build(&ops))
	}

	assert.Len(t, ops.Paths, 4)

	po, ok := ops.Paths["/pets"]
	assert.True(t, ok)
	assert.NotNil(t, po.Get)
	assertOperationBody(t,
		po.Get,
		"listPets",
		"Lists pets filtered by some parameters.",
		"This will show all available pets by default.\nYou can get the pets that are out of stock",
		[]string{"pets", "users"},
		[]string{"read", "write"},
	)
	assert.NotNil(t, po.Post)

	assertOperationBody(t,
		po.Post,
		"createPet",
		"Create a pet based on the parameters.",
		"",
		[]string{"pets", "users"},
		[]string{"read", "write"},
	)

	po, ok = ops.Paths["/orders"]
	assert.True(t, ok)
	assert.NotNil(t, po.Get)
	assertOperationBody(t,
		po.Get,
		"listOrders",
		"lists orders filtered by some parameters.",
		"",
		[]string{"orders"},
		[]string{"orders:read", "https://www.googleapis.com/auth/userinfo.email"},
	)
	assert.NotNil(t, po.Post)

	assertOperationBody(t,
		po.Post,
		"createOrder",
		"create an order based on the parameters.",
		"",
		[]string{"orders"},
		[]string{"read", "write"},
	)

	po, ok = ops.Paths["/orders/{id}"]
	assert.True(t, ok)
	assert.NotNil(t, po.Get)
	assertOperationBody(t,
		po.Get,
		"orderDetails",
		"gets the details for an order.",
		"",
		[]string{"orders"},
		[]string{"read", "write"},
	)

	assertOperationBody(t,
		po.Put,
		"updateOrder",
		"Update the details for an order.",
		"When the order doesn't exist this will return an error.",
		[]string{"orders"},
		[]string{"read", "write"},
	)

	assertOperationBody(t,
		po.Delete,
		"deleteOrder",
		"delete a particular order.",
		"",
		nil,
		[]string{"read", "write"},
	)

	validateRoutesParameters(t, ops)
}

func validateRoutesParameters(t *testing.T, ops spec.Paths) {
	po := ops.Paths["/pets"]
	assert.Equal(t, 2, len(po.Post.Parameters))

	// Testing standard param properties
	p := po.Post.Parameters[0]
	assert.Equal(t, "request", p.Name)
	assert.Equal(t, "body", p.In)
	assert.Equal(t, "The request model.", p.Description)

	// Testing "required" and "allowEmpty"
	p = po.Post.Parameters[1]
	assert.Equal(t, "id", p.Name)
	assert.Equal(t, "The pet id", p.Description)
	assert.Equal(t, "path", p.In)
	assert.Equal(t, true, p.Required)
	assert.Equal(t, false, p.AllowEmptyValue)

	po = ops.Paths["/orders"]
	assert.Equal(t, 2, len(po.Post.Parameters))

	// Testing invalid value for "in"
	p = po.Post.Parameters[0]
	assert.Equal(t, "id", p.Name)
	assert.Equal(t, "The order id", p.Description)
	assert.Equal(t, "", p.In) // Invalid value should not be set
	assert.Equal(t, false, p.Required)
	assert.Equal(t, true, p.AllowEmptyValue)

	p = po.Post.Parameters[1]
	assert.Equal(t, "request", p.Name)
	assert.Equal(t, "body", p.In)
	assert.Equal(t, "The request model.", p.Description)

	po = ops.Paths["/param-test"]
	assert.Equal(t, 6, len(po.Post.Parameters))

	// Testing number param with "max" and "min" constraints
	p = po.Post.Parameters[0]
	assert.Equal(t, "someNumber", p.Name)
	assert.Equal(t, "some number", p.Description)
	assert.Equal(t, "path", p.In)
	assert.Equal(t, true, p.Required)
	assert.Equal(t, "number", p.Type)
	assert.Equal(t, "number", p.Schema.Type[0])
	min, max, def := float64(10), float64(20), float64(15)
	assert.Equal(t, &max, p.Schema.Maximum)
	assert.Equal(t, &min, p.Schema.Minimum)
	assert.Equal(t, def, p.Schema.Default)

	// Testing array param provided as query string. Testing "minLength" and "maxLength" constraints for "array" types
	p = po.Post.Parameters[1]
	assert.Equal(t, "someQuery", p.Name)
	assert.Equal(t, "some query values", p.Description)
	assert.Equal(t, "query", p.In)
	assert.Equal(t, false, p.Required)
	assert.Equal(t, "array", p.Type)
	assert.Equal(t, "array", p.Schema.Type[0])
	minLen, maxLen := int64(5), int64(20)
	assert.Equal(t, &maxLen, p.Schema.MaxLength)
	assert.Equal(t, &minLen, p.Schema.MinLength)

	// Testing boolean param with default value
	p = po.Post.Parameters[2]
	assert.Equal(t, "someBoolean", p.Name)
	assert.Equal(t, "some boolean", p.Description)
	assert.Equal(t, "path", p.In)
	assert.Equal(t, false, p.Required)
	assert.Equal(t, "boolean", p.Type)
	assert.Equal(t, "boolean", p.Schema.Type[0])
	assert.Equal(t, true, p.Schema.Default)

	// Testing that "min", "max", "minLength" and "maxLength" constraints will only be considered if the right type is provided
	p = po.Post.Parameters[3]
	assert.Equal(t, "constraintsOnInvalidType", p.Name)
	assert.Equal(t, "test constraints on invalid types", p.Description)
	assert.Equal(t, "query", p.In)
	assert.Equal(t, "boolean", p.Type)
	assert.Equal(t, "boolean", p.Schema.Type[0])
	assert.Nil(t, p.Schema.Maximum)
	assert.Nil(t, p.Schema.Minimum)
	assert.Nil(t, p.Schema.MaxLength)
	assert.Nil(t, p.Schema.MinLength)
	assert.Equal(t, "abcde", p.Schema.Format)
	assert.Equal(t, false, p.Schema.Default)

	// Testing that when "type" is not provided, a schema will not be created
	p = po.Post.Parameters[4]
	assert.Equal(t, "noType", p.Name)
	assert.Equal(t, "test no type", p.Description)
	assert.Equal(t, "", p.Type)
	assert.Nil(t, p.Schema)

	// Testing a request body that takes a string value defined by a list of possible values in "enum"
	p = po.Post.Parameters[5]
	assert.Equal(t, "request", p.Name)
	assert.Equal(t, "The request model.", p.Description)
	assert.Equal(t, "body", p.In)
	assert.Equal(t, "string", p.Type)
	assert.Equal(t, "orange", p.Schema.Default)
	assert.Equal(t, []interface{}{"apple", "orange", "pineapple", "peach", "plum"}, p.Schema.Enum)
}

func assertOperation(t *testing.T, op *spec.Operation, id, summary, description string, tags, scopes []string) {
	assert.NotNil(t, op)
	assert.Equal(t, summary, op.Summary)
	assert.Equal(t, description, op.Description)
	assert.Equal(t, id, op.ID)
	assert.EqualValues(t, tags, op.Tags)
	assert.EqualValues(t, []string{"application/json", "application/x-protobuf"}, op.Consumes)
	assert.EqualValues(t, []string{"application/json", "application/x-protobuf"}, op.Produces)
	assert.EqualValues(t, []string{"http", "https", "ws", "wss"}, op.Schemes)
	assert.Len(t, op.Security, 2)
	akv, ok := op.Security[0]["api_key"]
	assert.True(t, ok)
	// akv must be defined & not empty
	assert.NotNil(t, akv)
	assert.Empty(t, akv)

	vv, ok := op.Security[1]["oauth"]
	assert.True(t, ok)
	assert.EqualValues(t, scopes, vv)

	assert.NotNil(t, op.Responses.Default)
	assert.Equal(t, "#/responses/genericError", op.Responses.Default.Ref.String())

	rsp, ok := op.Responses.StatusCodeResponses[200]
	assert.True(t, ok)
	assert.Equal(t, "#/responses/someResponse", rsp.Ref.String())
	rsp, ok = op.Responses.StatusCodeResponses[422]
	assert.True(t, ok)
	assert.Equal(t, "#/responses/validationError", rsp.Ref.String())
}

func assertOperationBody(t *testing.T, op *spec.Operation, id, summary, description string, tags, scopes []string) {
	assert.NotNil(t, op)
	assert.Equal(t, summary, op.Summary)
	assert.Equal(t, description, op.Description)
	assert.Equal(t, id, op.ID)
	assert.EqualValues(t, tags, op.Tags)
	assert.EqualValues(t, []string{"application/json", "application/x-protobuf"}, op.Consumes)
	assert.EqualValues(t, []string{"application/json", "application/x-protobuf"}, op.Produces)
	assert.EqualValues(t, []string{"http", "https", "ws", "wss"}, op.Schemes)
	assert.Len(t, op.Security, 2)
	akv, ok := op.Security[0]["api_key"]
	assert.True(t, ok)
	// akv must be defined & not empty
	assert.NotNil(t, akv)
	assert.Empty(t, akv)

	vv, ok := op.Security[1]["oauth"]
	assert.True(t, ok)
	assert.EqualValues(t, scopes, vv)

	assert.NotNil(t, op.Responses.Default)
	assert.Equal(t, "", op.Responses.Default.Ref.String())
	assert.Equal(t, "#/definitions/genericError", op.Responses.Default.Schema.Ref.String())

	rsp, ok := op.Responses.StatusCodeResponses[200]
	assert.True(t, ok)
	assert.Equal(t, "", rsp.Ref.String())
	assert.Equal(t, "#/definitions/someResponse", rsp.Schema.Ref.String())
	rsp, ok = op.Responses.StatusCodeResponses[422]
	assert.True(t, ok)
	assert.Equal(t, "", rsp.Ref.String())
	assert.Equal(t, "#/definitions/validationError", rsp.Schema.Ref.String())
}
