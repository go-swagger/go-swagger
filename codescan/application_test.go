package codescan

import (
	"sort"
	"testing"

	"github.com/go-openapi/spec"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

var (
	petstoreCtx       *scanCtx
	classificationCtx *scanCtx
)

func loadPetstorePkgsCtx(t testing.TB) *scanCtx {
	if petstoreCtx != nil {
		return petstoreCtx
	}
	sctx, err := newScanCtx(&Options{
		Packages: []string{"github.com/go-swagger/go-swagger/fixtures/goparsing/petstore/..."},
	})
	require.NoError(t, err)
	petstoreCtx = sctx
	return petstoreCtx
}

func loadClassificationPkgsCtx(t testing.TB, extra ...string) *scanCtx {
	if classificationCtx != nil {
		return classificationCtx
	}
	sctx, err := newScanCtx(&Options{
		Packages: append([]string{
			"github.com/go-swagger/go-swagger/fixtures/goparsing/classification",
			"github.com/go-swagger/go-swagger/fixtures/goparsing/classification/models",
			"github.com/go-swagger/go-swagger/fixtures/goparsing/classification/operations",
		}, extra...),
	})
	require.NoError(t, err)
	classificationCtx = sctx
	return classificationCtx
}

func TestApplication_LoadCode(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	require.NotNil(t, sctx)
	require.Len(t, sctx.app.Models, 30)
	require.Len(t, sctx.app.Meta, 1)
	require.Len(t, sctx.app.Routes, 7)
	require.Empty(t, sctx.app.Operations)
	require.Len(t, sctx.app.Parameters, 10)
	require.Len(t, sctx.app.Responses, 11)
}

func TestAppScanner_NewSpec(t *testing.T) {
	doc, err := Run(&Options{
		Packages: []string{"github.com/go-swagger/go-swagger/fixtures/goparsing/petstore/..."},
	})
	require.NoError(t, err)
	if assert.NotNil(t, doc) {
		// b, _ := json.MarshalIndent(doc.Responses, "", "  ")
		// log.Println(string(b))
		verifyParsedPetStore(t, doc)
	}
}

func TestAppScanner_Definitions(t *testing.T) {
	doc, err := Run(&Options{
		Packages:   []string{"github.com/go-swagger/go-swagger/fixtures/goparsing/bookings/..."},
		ScanModels: true,
	})
	require.NoError(t, err)
	if assert.NotNil(t, doc) {
		_, ok := doc.Definitions["Booking"]
		assert.True(t, ok, "Should include cross repo structs")
		_, ok = doc.Definitions["Customer"]
		assert.True(t, ok, "Should include package structs with swagger:model")
		_, ok = doc.Definitions["DateRange"]
		assert.True(t, ok, "Should include package structs that are used in responses")
		_, ok = doc.Definitions["BookingResponse"]
		assert.False(t, ok, "Should not include responses")
		_, ok = doc.Definitions["IgnoreMe"]
		assert.False(t, ok, "Should not include un-annotated/un-referenced structs")
	}
}

func verifyParsedPetStore(t testing.TB, doc *spec.Swagger) {
	assert.EqualValues(t, []string{"application/json"}, doc.Consumes)
	assert.EqualValues(t, []string{"application/json"}, doc.Produces)
	assert.EqualValues(t, []string{"http", "https"}, doc.Schemes)
	assert.Equal(t, "localhost", doc.Host)
	assert.Equal(t, "/v2", doc.BasePath)

	verifyInfo(t, doc.Info)

	if assert.NotNil(t, doc.Paths) {
		assert.Len(t, doc.Paths.Paths, 5)
	}
	var keys []string
	for k := range doc.Definitions {
		keys = append(keys, k)
	}
	assert.Len(t, keys, 3)
	assert.Len(t, doc.Responses, 4)

	definitions := doc.Definitions
	mod, ok := definitions["tag"]
	assert.True(t, ok)
	assert.Equal(t, spec.StringOrArray([]string{"object"}), mod.Type)
	assert.Equal(t, "A Tag is an extra piece of data to provide more information about a pet.", mod.Title)
	assert.Equal(t, "It is used to describe the animals available in the store.", mod.Description)
	assert.Len(t, mod.Required, 2)

	assertProperty(t, &mod, "integer", "id", "int64", "ID")
	prop, ok := mod.Properties["id"]
	assert.True(t, ok, "should have had an 'id' property")
	assert.Equal(t, "The id of the tag.", prop.Description)

	assertProperty(t, &mod, "string", "value", "", "Value")
	prop, ok = mod.Properties["value"]
	assert.True(t, ok)
	assert.Equal(t, "The value of the tag.", prop.Description)

	mod, ok = definitions["pet"]
	assert.True(t, ok)
	assert.Equal(t, spec.StringOrArray([]string{"object"}), mod.Type)
	assert.Equal(t, "A Pet is the main product in the store.", mod.Title)
	assert.Equal(t, "It is used to describe the animals available in the store.", mod.Description)
	assert.Len(t, mod.Required, 2)

	assertProperty(t, &mod, "integer", "id", "int64", "ID")
	prop, ok = mod.Properties["id"]
	assert.True(t, ok, "should have had an 'id' property")
	assert.Equal(t, "The id of the pet.", prop.Description)

	assertProperty(t, &mod, "string", "name", "", "Name")
	prop, ok = mod.Properties["name"]
	assert.True(t, ok)
	assert.Equal(t, "The name of the pet.", prop.Description)
	assert.EqualValues(t, 3, *prop.MinLength)
	assert.EqualValues(t, 50, *prop.MaxLength)
	assert.Equal(t, "\\w[\\w-]+", prop.Pattern)

	assertArrayProperty(t, &mod, "string", "photoUrls", "", "PhotoURLs")
	prop, ok = mod.Properties["photoUrls"]
	assert.True(t, ok)
	assert.Equal(t, "The photo urls for the pet.\nThis only accepts jpeg or png images.", prop.Description)
	if assert.NotNil(t, prop.Items) && assert.NotNil(t, prop.Items.Schema) {
		assert.Equal(t, "\\.(jpe?g|png)$", prop.Items.Schema.Pattern)
	}

	assertProperty(t, &mod, "string", "status", "", "Status")
	prop, ok = mod.Properties["status"]
	assert.True(t, ok)
	assert.Equal(t, "The current status of the pet in the store.", prop.Description)
	assert.Equal(t, []interface{}{"available", "pending", "sold"}, prop.Enum)

	assertProperty(t, &mod, "string", "birthday", "date", "Birthday")
	prop, ok = mod.Properties["birthday"]
	assert.True(t, ok)
	assert.Equal(t, "The pet's birthday", prop.Description)

	assertArrayRef(t, &mod, "tags", "Tags", "#/definitions/tag")
	prop, ok = mod.Properties["tags"]
	assert.True(t, ok)
	assert.Equal(t, "Extra bits of information attached to this pet.", prop.Description)

	mod, ok = definitions["order"]
	assert.True(t, ok)
	assert.Len(t, mod.Properties, 4)
	assert.Len(t, mod.Required, 3)

	assertProperty(t, &mod, "integer", "id", "int64", "ID")
	prop, ok = mod.Properties["id"]
	assert.True(t, ok, "should have had an 'id' property")
	assert.Equal(t, "the ID of the order", prop.Description)

	assertProperty(t, &mod, "integer", "userId", "int64", "UserID")
	prop, ok = mod.Properties["userId"]
	assert.True(t, ok, "should have had an 'userId' property")
	assert.Equal(t, "the id of the user who placed the order.", prop.Description)

	assertProperty(t, &mod, "string", "orderedAt", "date-time", "OrderedAt")
	prop, ok = mod.Properties["orderedAt"]
	assert.Equal(t, "the time at which this order was made.", prop.Description)
	assert.True(t, ok, "should have an 'orderedAt' property")

	assertArrayProperty(t, &mod, "object", "items", "", "Items")
	prop, ok = mod.Properties["items"]
	assert.True(t, ok, "should have an 'items' slice")
	assert.NotNil(t, prop.Items, "items should have had an items property")
	assert.NotNil(t, prop.Items.Schema, "items.items should have had a schema property")

	itprop := prop.Items.Schema
	assert.Len(t, itprop.Properties, 2)
	assert.Len(t, itprop.Required, 2)

	assertProperty(t, itprop, "integer", "petId", "int64", "PetID")
	iprop, ok := itprop.Properties["petId"]
	assert.True(t, ok, "should have had a 'petId' property")
	assert.Equal(t, "the id of the pet to order", iprop.Description)

	assertProperty(t, itprop, "integer", "qty", "int32", "Quantity")
	iprop, ok = itprop.Properties["qty"]
	assert.True(t, ok, "should have had a 'qty' property")
	assert.Equal(t, "the quantity of this pet to order", iprop.Description)
	assert.EqualValues(t, 1, *iprop.Minimum)

	// responses
	resp, ok := doc.Responses["genericError"]
	assert.True(t, ok)
	assert.NotNil(t, resp.Schema)
	assert.Len(t, resp.Schema.Properties, 2)
	assertProperty(t, resp.Schema, "integer", "code", "int32", "Code")
	assertProperty(t, resp.Schema, "string", "message", "", "Message")

	resp, ok = doc.Responses["validationError"]
	assert.True(t, ok)
	assert.NotNil(t, resp.Schema)
	assert.Len(t, resp.Schema.Properties, 3)
	assertProperty(t, resp.Schema, "integer", "code", "int32", "Code")
	assertProperty(t, resp.Schema, "string", "message", "", "Message")
	assertProperty(t, resp.Schema, "string", "field", "", "Field")

	resp, ok = doc.Responses["MarkdownRender"]
	assert.True(t, ok)
	assert.NotNil(t, resp.Schema)
	assert.True(t, resp.Schema.Type.Contains("string"))

	paths := doc.Paths.Paths

	// path /pets
	op, ok := paths["/pets"]
	assert.True(t, ok)
	assert.NotNil(t, op)

	// listPets
	assert.NotNil(t, op.Get)
	assert.Equal(t, "Lists the pets known to the store.", op.Get.Summary)
	assert.Equal(t, "By default it will only lists pets that are available for sale.\nThis can be changed with the status flag.", op.Get.Description)
	assert.Equal(t, "listPets", op.Get.ID)
	assert.EqualValues(t, []string{"pets"}, op.Get.Tags)
	assert.True(t, op.Get.Deprecated)
	var names namedParams
	for i, v := range op.Get.Parameters {
		names = append(names, namedParam{Index: i, Name: v.Name})
	}
	sort.Sort(names)
	sparam := op.Get.Parameters[names[1].Index]
	assert.Equal(t, "Status", sparam.Description)
	assert.Equal(t, "query", sparam.In)
	assert.Equal(t, "string", sparam.Type)
	assert.Equal(t, "", sparam.Format)
	assert.Equal(t, []interface{}{"available", "pending", "sold"}, sparam.Enum)
	assert.False(t, sparam.Required)
	assert.Equal(t, "Status", sparam.Extensions["x-go-name"])
	assert.Equal(t, "#/responses/genericError", op.Get.Responses.Default.Ref.String())
	assert.Len(t, op.Get.Parameters, 2)
	sparam1 := op.Get.Parameters[names[0].Index]
	assert.Equal(t, "Birthday", sparam1.Description)
	assert.Equal(t, "query", sparam1.In)
	assert.Equal(t, "string", sparam1.Type)
	assert.Equal(t, "date", sparam1.Format)
	assert.False(t, sparam1.Required)
	assert.Equal(t, "Birthday", sparam1.Extensions["x-go-name"])
	rs, ok := op.Get.Responses.StatusCodeResponses[200]
	assert.True(t, ok)
	assert.NotNil(t, rs.Schema)
	aprop := rs.Schema
	assert.Equal(t, "array", aprop.Type[0])
	assert.NotNil(t, aprop.Items)
	assert.NotNil(t, aprop.Items.Schema)
	assert.Equal(t, "#/definitions/pet", aprop.Items.Schema.Ref.String())

	// createPet
	assert.NotNil(t, op.Post)
	assert.Equal(t, "Creates a new pet in the store.", op.Post.Summary)
	assert.Equal(t, "", op.Post.Description)
	assert.Equal(t, "createPet", op.Post.ID)
	assert.EqualValues(t, []string{"pets"}, op.Post.Tags)
	verifyRefParam(t, op.Post.Parameters[0], "The pet to submit.", "pet")
	assert.Equal(t, "#/responses/genericError", op.Post.Responses.Default.Ref.String())
	rs, ok = op.Post.Responses.StatusCodeResponses[200]
	assert.True(t, ok)
	assert.NotNil(t, rs.Schema)
	aprop = rs.Schema
	assert.Equal(t, "#/definitions/pet", aprop.Ref.String())

	// path /pets/{id}
	op, ok = paths["/pets/{id}"]
	assert.True(t, ok)
	assert.NotNil(t, op)

	// getPetById
	assert.NotNil(t, op.Get)
	assert.Equal(t, "Gets the details for a pet.", op.Get.Summary)
	assert.Equal(t, "", op.Get.Description)
	assert.Equal(t, "getPetById", op.Get.ID)
	assert.EqualValues(t, []string{"pets"}, op.Get.Tags)
	verifyIDParam(t, op.Get.Parameters[0], "The ID of the pet")
	assert.Equal(t, "#/responses/genericError", op.Get.Responses.Default.Ref.String())
	rs, ok = op.Get.Responses.StatusCodeResponses[200]
	assert.True(t, ok)
	assert.NotNil(t, rs.Schema)
	aprop = rs.Schema
	assert.Equal(t, "#/definitions/pet", aprop.Ref.String())

	// updatePet
	assert.NotNil(t, op.Put)
	assert.Equal(t, "Updates the details for a pet.", op.Put.Summary)
	assert.Equal(t, "", op.Put.Description)
	assert.Equal(t, "updatePet", op.Put.ID)
	assert.EqualValues(t, []string{"pets"}, op.Put.Tags)
	verifyIDParam(t, op.Put.Parameters[0], "The ID of the pet")
	verifyRefParam(t, op.Put.Parameters[1], "The pet to submit.", "pet")
	assert.Equal(t, "#/responses/genericError", op.Put.Responses.Default.Ref.String())
	rs, ok = op.Put.Responses.StatusCodeResponses[200]
	assert.True(t, ok)
	assert.NotNil(t, rs.Schema)
	aprop = rs.Schema
	assert.Equal(t, "#/definitions/pet", aprop.Ref.String())

	// deletePet
	assert.NotNil(t, op.Delete)
	assert.Equal(t, "Deletes a pet from the store.", op.Delete.Summary)
	assert.Equal(t, "", op.Delete.Description)
	assert.Equal(t, "deletePet", op.Delete.ID)
	assert.EqualValues(t, []string{"pets"}, op.Delete.Tags)
	verifyIDParam(t, op.Delete.Parameters[0], "The ID of the pet")
	assert.Equal(t, "#/responses/genericError", op.Delete.Responses.Default.Ref.String())
	_, ok = op.Delete.Responses.StatusCodeResponses[204]
	assert.True(t, ok)

	// path /orders/{id}
	op, ok = paths["/orders/{id}"]
	assert.True(t, ok)
	assert.NotNil(t, op)

	// getOrderDetails
	assert.NotNil(t, op.Get)
	assert.Equal(t, "Gets the details for an order.", op.Get.Summary)
	assert.Equal(t, "", op.Get.Description)
	assert.Equal(t, "getOrderDetails", op.Get.ID)
	assert.EqualValues(t, []string{"orders"}, op.Get.Tags)
	verifyIDParam(t, op.Get.Parameters[0], "The ID of the order")
	assert.Equal(t, "#/responses/genericError", op.Get.Responses.Default.Ref.String())
	rs, ok = op.Get.Responses.StatusCodeResponses[200]
	assert.True(t, ok)
	assert.Equal(t, "#/responses/orderResponse", rs.Ref.String())
	rsm := doc.Responses["orderResponse"]
	assert.NotNil(t, rsm.Schema)
	assert.Equal(t, "#/definitions/order", rsm.Schema.Ref.String())

	// cancelOrder
	assert.NotNil(t, op.Delete)
	assert.Equal(t, "Deletes an order.", op.Delete.Summary)
	assert.Equal(t, "", op.Delete.Description)
	assert.Equal(t, "cancelOrder", op.Delete.ID)
	assert.EqualValues(t, []string{"orders"}, op.Delete.Tags)
	verifyIDParam(t, op.Delete.Parameters[0], "The ID of the order")
	assert.Equal(t, "#/responses/genericError", op.Delete.Responses.Default.Ref.String())
	_, ok = op.Delete.Responses.StatusCodeResponses[204]
	assert.True(t, ok)

	// updateOrder
	assert.NotNil(t, op.Put)
	assert.Equal(t, "Updates an order.", op.Put.Summary)
	assert.Equal(t, "", op.Put.Description)
	assert.Equal(t, "updateOrder", op.Put.ID)
	assert.EqualValues(t, []string{"orders"}, op.Put.Tags)
	verifyIDParam(t, op.Put.Parameters[0], "The ID of the order")
	verifyRefParam(t, op.Put.Parameters[1], "The order to submit", "order")
	assert.Equal(t, "#/responses/genericError", op.Put.Responses.Default.Ref.String())
	rs, ok = op.Put.Responses.StatusCodeResponses[200]
	assert.True(t, ok)
	assert.NotNil(t, rs.Schema)
	aprop = rs.Schema
	assert.Equal(t, "#/definitions/order", aprop.Ref.String())

	// path /orders
	op, ok = paths["/orders"]
	assert.True(t, ok)
	assert.NotNil(t, op)

	// createOrder
	assert.NotNil(t, op.Post)
	assert.Equal(t, "Creates an order.", op.Post.Summary)
	assert.Equal(t, "", op.Post.Description)
	assert.Equal(t, "createOrder", op.Post.ID)
	assert.EqualValues(t, []string{"orders"}, op.Post.Tags)
	verifyRefParam(t, op.Post.Parameters[0], "The order to submit", "order")
	assert.Equal(t, "#/responses/genericError", op.Post.Responses.Default.Ref.String())
	rs, ok = op.Post.Responses.StatusCodeResponses[200]
	assert.True(t, ok)
	assert.Equal(t, "#/responses/orderResponse", rs.Ref.String())
	rsm = doc.Responses["orderResponse"]
	assert.NotNil(t, rsm.Schema)
	assert.Equal(t, "#/definitions/order", rsm.Schema.Ref.String())
}

func verifyIDParam(t testing.TB, param spec.Parameter, description string) {
	assert.Equal(t, description, param.Description)
	assert.Equal(t, "path", param.In)
	assert.Equal(t, "integer", param.Type)
	assert.Equal(t, "int64", param.Format)
	assert.True(t, param.Required)
	assert.Equal(t, "ID", param.Extensions["x-go-name"])
}

func verifyRefParam(t testing.TB, param spec.Parameter, description, refed string) {
	assert.Equal(t, description, param.Description)
	assert.Equal(t, "body", param.In)
	// TODO: this may fail sometimes (seen on go1.12 windows test): require pointer to be valid and avoid panicking
	require.NotNil(t, param)
	require.NotNil(t, param.Schema)
	assert.Equal(t, "#/definitions/"+refed, param.Schema.Ref.String())
	assert.True(t, param.Required)
}

type namedParam struct {
	Index int
	Name  string
}

type namedParams []namedParam

func (g namedParams) Len() int           { return len(g) }
func (g namedParams) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g namedParams) Less(i, j int) bool { return g[i].Name < g[j].Name }
