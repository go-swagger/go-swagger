// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package codescan

import (
	"flag"
	"io"
	"log"
	"os"
	"sort"
	"testing"

	"github.com/go-openapi/spec"

	"github.com/go-openapi/testify/v2/assert"

	"github.com/go-openapi/testify/v2/require"
)

var (
	petstoreCtx       *scanCtx
	classificationCtx *scanCtx
)

var (
	enableSpecOutput bool
	enableDebug      bool
)

func init() {
	flag.BoolVar(&enableSpecOutput, "enable-spec-output", false, "enable spec gen test to write output to a file")
	flag.BoolVar(&enableDebug, "enable-debug", false, "enable debug output in tests")
}

func TestMain(m *testing.M) {
	// initializations to run tests in this package
	flag.Parse()

	if !enableDebug {
		log.SetOutput(io.Discard)
	} else {
		// enable full debug when test is run with -enable-debug arg
		Debug = true
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.SetOutput(os.Stderr)
	}

	os.Exit(m.Run())
}

func TestApplication_LoadCode(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	require.NotNil(t, sctx)
	require.Len(t, sctx.app.Models, 39)
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
		assert.TrueT(t, ok, "Should include cross repo structs")
		_, ok = doc.Definitions["Customer"]
		assert.TrueT(t, ok, "Should include package structs with swagger:model")
		_, ok = doc.Definitions["DateRange"]
		assert.TrueT(t, ok, "Should include package structs that are used in responses")
		_, ok = doc.Definitions["BookingResponse"]
		assert.FalseT(t, ok, "Should not include responses")
		_, ok = doc.Definitions["IgnoreMe"]
		assert.FalseT(t, ok, "Should not include un-annotated/un-referenced structs")
	}
}

func loadPetstorePkgsCtx(t *testing.T) *scanCtx {
	t.Helper()

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

func loadClassificationPkgsCtx(t *testing.T, extra ...string) *scanCtx {
	t.Helper()

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

func verifyParsedPetStore(t *testing.T, doc *spec.Swagger) {
	t.Helper()

	verifyTop(t, doc)
	verifyInfo(t, doc.Info)
	verifyModels(t, doc.Definitions)
	verifyCommonResponses(t, doc.Responses)

	t.Run("with API paths", func(t *testing.T) {
		require.NotNil(t, doc.Paths)
		paths := doc.Paths.Paths
		require.Len(t, paths, 5)

		t.Run("with path /pets", func(t *testing.T) {
			op, ok := paths["/pets"]
			assert.TrueT(t, ok)
			assert.NotNil(t, op)

			t.Run("with GET: listPets", func(t *testing.T) {
				require.NotNil(t, op.Get)
				assert.EqualT(t, "Lists the pets known to the store.", op.Get.Summary)
				assert.EqualT(t, "By default it will only lists pets that are available for sale.\nThis can be changed with the status flag.", op.Get.Description)
				assert.EqualT(t, "listPets", op.Get.ID)
				assert.Equal(t, []string{"pets"}, op.Get.Tags)
				assert.TrueT(t, op.Get.Deprecated)
				names := make(namedParams, 0, len(op.Get.Parameters))
				for i, v := range op.Get.Parameters {
					names = append(names, namedParam{Index: i, Name: v.Name})
				}
				sort.Sort(names)
				sparam := op.Get.Parameters[names[1].Index]
				assert.EqualT(t, "Status\navailable STATUS_AVAILABLE\npending STATUS_PENDING\nsold STATUS_SOLD", sparam.Description)
				assert.EqualT(t, "query", sparam.In)
				assert.EqualT(t, "string", sparam.Type)
				assert.Empty(t, sparam.Format)
				assert.Equal(t, []any{"available", "pending", "sold"}, sparam.Enum)
				assert.FalseT(t, sparam.Required)
				assert.Equal(t, "Status", sparam.Extensions["x-go-name"])
				assert.EqualT(t, "#/responses/genericError", op.Get.Responses.Default.Ref.String())
				assert.Len(t, op.Get.Parameters, 2)
				sparam1 := op.Get.Parameters[names[0].Index]
				assert.EqualT(t, "Birthday", sparam1.Description)
				assert.EqualT(t, "query", sparam1.In)
				assert.EqualT(t, "string", sparam1.Type)
				assert.EqualT(t, "date", sparam1.Format)
				assert.FalseT(t, sparam1.Required)
				assert.Equal(t, "Birthday", sparam1.Extensions["x-go-name"])
				rs, ok := op.Get.Responses.StatusCodeResponses[200]
				require.TrueT(t, ok)
				require.NotNil(t, rs.Schema)
				aprop := rs.Schema
				assert.EqualT(t, "array", aprop.Type[0])
				assert.NotNil(t, aprop.Items)
				assert.NotNil(t, aprop.Items.Schema)
				assert.EqualT(t, "#/definitions/pet", aprop.Items.Schema.Ref.String())
			})

			t.Run("with POST: createPet", func(t *testing.T) {
				require.NotNil(t, op.Post)
				assert.EqualT(t, "Creates a new pet in the store.", op.Post.Summary)
				assert.Empty(t, op.Post.Description)
				assert.EqualT(t, "createPet", op.Post.ID)
				assert.Equal(t, []string{"pets"}, op.Post.Tags)
				verifyRefParam(t, op.Post.Parameters[0], "The pet to submit.", "pet")
				assert.EqualT(t, "#/responses/genericError", op.Post.Responses.Default.Ref.String())
				rs, ok := op.Post.Responses.StatusCodeResponses[200]
				assert.TrueT(t, ok)
				assert.NotNil(t, rs.Schema)
				aprop := rs.Schema
				assert.EqualT(t, "#/definitions/pet", aprop.Ref.String())
			})
		})

		t.Run("with path /pets/{id}", func(t *testing.T) {
			op, ok := paths["/pets/{id}"]
			require.TrueT(t, ok)
			require.NotNil(t, op)

			t.Run("with GET: getPetById", func(t *testing.T) {
				require.NotNil(t, op.Get)
				assert.EqualT(t, "Gets the details for a pet.", op.Get.Summary)
				assert.Empty(t, op.Get.Description)
				assert.EqualT(t, "getPetById", op.Get.ID)
				assert.Equal(t, []string{"pets"}, op.Get.Tags)
				verifyIDParam(t, op.Get.Parameters[0], "The ID of the pet")
				assert.EqualT(t, "#/responses/genericError", op.Get.Responses.Default.Ref.String())
				rs, ok := op.Get.Responses.StatusCodeResponses[200]
				require.TrueT(t, ok)
				require.NotNil(t, rs.Schema)
				aprop := rs.Schema
				assert.EqualT(t, "#/definitions/pet", aprop.Ref.String())
			})

			t.Run("with PUT: updatePet", func(t *testing.T) {
				require.NotNil(t, op.Put)
				assert.EqualT(t, "Updates the details for a pet.", op.Put.Summary)
				assert.Empty(t, op.Put.Description)
				assert.EqualT(t, "updatePet", op.Put.ID)
				assert.Equal(t, []string{"pets"}, op.Put.Tags)
				verifyIDParam(t, op.Put.Parameters[0], "The ID of the pet")
				verifyRefParam(t, op.Put.Parameters[1], "The pet to submit.", "pet")
				assert.EqualT(t, "#/responses/genericError", op.Put.Responses.Default.Ref.String())
				rs, ok := op.Put.Responses.StatusCodeResponses[200]
				require.TrueT(t, ok)
				require.NotNil(t, rs.Schema)
				aprop := rs.Schema
				assert.EqualT(t, "#/definitions/pet", aprop.Ref.String())
			})

			t.Run("with DELETE: deletePet", func(t *testing.T) {
				require.NotNil(t, op.Delete)
				assert.EqualT(t, "Deletes a pet from the store.", op.Delete.Summary)
				assert.Empty(t, op.Delete.Description)
				assert.EqualT(t, "deletePet", op.Delete.ID)
				assert.Equal(t, []string{"pets"}, op.Delete.Tags)
				verifyIDParam(t, op.Delete.Parameters[0], "The ID of the pet")
				assert.EqualT(t, "#/responses/genericError", op.Delete.Responses.Default.Ref.String())
				_, ok := op.Delete.Responses.StatusCodeResponses[204]
				assert.TrueT(t, ok)
			})
		})

		t.Run("with path /orders/{id}", func(t *testing.T) {
			op, ok := paths["/orders/{id}"]
			require.TrueT(t, ok)
			require.NotNil(t, op)

			t.Run("with GET: getOrderDetails", func(t *testing.T) {
				require.NotNil(t, op.Get)
				assert.EqualT(t, "Gets the details for an order.", op.Get.Summary)
				assert.Empty(t, op.Get.Description)
				assert.EqualT(t, "getOrderDetails", op.Get.ID)
				assert.Equal(t, []string{"orders"}, op.Get.Tags)
				verifyIDParam(t, op.Get.Parameters[0], "The ID of the order")
				assert.EqualT(t, "#/responses/genericError", op.Get.Responses.Default.Ref.String())
				rs, ok := op.Get.Responses.StatusCodeResponses[200]
				require.TrueT(t, ok)
				assert.EqualT(t, "#/responses/orderResponse", rs.Ref.String())
				rsm := doc.Responses["orderResponse"]
				assert.NotNil(t, rsm.Schema)
				assert.EqualT(t, "#/definitions/order", rsm.Schema.Ref.String())
			})

			t.Run("with DELETE: cancelOrder", func(t *testing.T) {
				require.NotNil(t, op.Delete)
				assert.EqualT(t, "Deletes an order.", op.Delete.Summary)
				assert.Empty(t, op.Delete.Description)
				assert.EqualT(t, "cancelOrder", op.Delete.ID)
				assert.Equal(t, []string{"orders"}, op.Delete.Tags)
				verifyIDParam(t, op.Delete.Parameters[0], "The ID of the order")
				assert.EqualT(t, "#/responses/genericError", op.Delete.Responses.Default.Ref.String())
				_, ok := op.Delete.Responses.StatusCodeResponses[204]
				assert.TrueT(t, ok)
			})

			t.Run("with PUT: updateOrder", func(t *testing.T) {
				require.NotNil(t, op.Put)
				assert.EqualT(t, "Updates an order.", op.Put.Summary)
				assert.Empty(t, op.Put.Description)
				assert.EqualT(t, "updateOrder", op.Put.ID)
				assert.Equal(t, []string{"orders"}, op.Put.Tags)
				verifyIDParam(t, op.Put.Parameters[0], "The ID of the order")
				verifyRefParam(t, op.Put.Parameters[1], "The order to submit", "order")
				assert.EqualT(t, "#/responses/genericError", op.Put.Responses.Default.Ref.String())
				rs, ok := op.Put.Responses.StatusCodeResponses[200]
				require.TrueT(t, ok)
				require.NotNil(t, rs.Schema)
				aprop := rs.Schema
				assert.EqualT(t, "#/definitions/order", aprop.Ref.String())
			})
		})

		t.Run("with path /orders", func(t *testing.T) {
			op, ok := paths["/orders"]
			require.TrueT(t, ok)
			require.NotNil(t, op)

			t.Run("with POST: createOrder", func(t *testing.T) {
				require.NotNil(t, op.Post)
				assert.EqualT(t, "Creates an order.", op.Post.Summary)
				assert.Empty(t, op.Post.Description)
				assert.EqualT(t, "createOrder", op.Post.ID)
				assert.Equal(t, []string{"orders"}, op.Post.Tags)
				verifyRefParam(t, op.Post.Parameters[0], "The order to submit", "order")
				assert.EqualT(t, "#/responses/genericError", op.Post.Responses.Default.Ref.String())
				rs, ok := op.Post.Responses.StatusCodeResponses[200]
				require.TrueT(t, ok)
				assert.EqualT(t, "#/responses/orderResponse", rs.Ref.String())
				rsm := doc.Responses["orderResponse"]
				require.NotNil(t, rsm.Schema)
				assert.EqualT(t, "#/definitions/order", rsm.Schema.Ref.String())
			})
		})
	})
}

func verifyTop(t *testing.T, doc *spec.Swagger) {
	t.Helper()

	t.Run("with top level specification", func(t *testing.T) {
		t.Run("should consume and produce JSON", func(t *testing.T) {
			assert.Equal(t, []string{"application/json"}, doc.Consumes)
			assert.Equal(t, []string{"application/json"}, doc.Produces)
		})
		t.Run("schemes should be http and https", func(t *testing.T) {
			assert.Equal(t, []string{"http", "https"}, doc.Schemes)
		})
		t.Run("API host should be localhost", func(t *testing.T) {
			assert.EqualT(t, "localhost", doc.Host)
		})
		t.Run("check API base path", func(t *testing.T) {
			assert.EqualT(t, "/v2", doc.BasePath)
		})
	})
}

func verifyCommonResponses(t *testing.T, responses map[string]spec.Response) {
	t.Helper()

	t.Run("with responses", func(t *testing.T) {
		require.Len(t, responses, 4)

		t.Run("should define genericError", func(t *testing.T) {
			resp, ok := responses["genericError"]
			require.TrueT(t, ok)
			require.NotNil(t, resp.Schema)
			assert.Len(t, resp.Schema.Properties, 2)
			assertProperty(t, resp.Schema, "integer", "code", "int32", "Code")
			assertProperty(t, resp.Schema, "string", "message", "", "Message")
		})

		t.Run("should define validationError", func(t *testing.T) {
			resp, ok := responses["validationError"]
			require.TrueT(t, ok)
			require.NotNil(t, resp.Schema)
			assert.Len(t, resp.Schema.Properties, 3)
			assertProperty(t, resp.Schema, "integer", "code", "int32", "Code")
			assertProperty(t, resp.Schema, "string", "message", "", "Message")
			assertProperty(t, resp.Schema, "string", "field", "", "Field")
		})

		t.Run("should define MarkdownRender", func(t *testing.T) {
			resp, ok := responses["MarkdownRender"]
			require.TrueT(t, ok)
			require.NotNil(t, resp.Schema)
			assert.TrueT(t, resp.Schema.Type.Contains("string"))
		})
	})
}

func verifyModels(t *testing.T, definitions spec.Definitions) {
	t.Helper()

	t.Run("with models definitions", func(t *testing.T) {
		keys := make([]string, 0, len(definitions))
		for k := range definitions {
			keys = append(keys, k)
		}
		require.Len(t, keys, 3)

		mod, ok := definitions["tag"]
		require.TrueT(t, ok)
		assert.Equal(t, spec.StringOrArray([]string{"object"}), mod.Type)
		assert.EqualT(t, "A Tag is an extra piece of data to provide more information about a pet.", mod.Title)
		assert.EqualT(t, "It is used to describe the animals available in the store.", mod.Description)
		assert.Len(t, mod.Required, 2)

		assertProperty(t, &mod, "integer", "id", "int64", "ID")
		prop, ok := mod.Properties["id"]
		require.TrueT(t, ok, "should have had an 'id' property")
		assert.EqualT(t, "The id of the tag.", prop.Description)

		assertProperty(t, &mod, "string", "value", "", "Value")
		prop, ok = mod.Properties["value"]
		require.TrueT(t, ok)
		assert.EqualT(t, "The value of the tag.", prop.Description)

		mod, ok = definitions["pet"]
		require.TrueT(t, ok)
		assert.Equal(t, spec.StringOrArray([]string{"object"}), mod.Type)
		assert.EqualT(t, "A Pet is the main product in the store.", mod.Title)
		assert.EqualT(t, "It is used to describe the animals available in the store.", mod.Description)
		assert.Len(t, mod.Required, 2)

		assertProperty(t, &mod, "integer", "id", "int64", "ID")
		prop, ok = mod.Properties["id"]
		require.TrueT(t, ok, "should have had an 'id' property")
		assert.EqualT(t, "The id of the pet.", prop.Description)

		assertProperty(t, &mod, "string", "name", "", "Name")
		prop, ok = mod.Properties["name"]
		require.TrueT(t, ok)
		assert.EqualT(t, "The name of the pet.", prop.Description)
		assert.EqualValues(t, 3, *prop.MinLength)
		assert.EqualValues(t, 50, *prop.MaxLength)
		assert.EqualT(t, "\\w[\\w-]+", prop.Pattern)

		assertArrayProperty(t, &mod, "string", "photoUrls", "", "PhotoURLs")
		prop, ok = mod.Properties["photoUrls"]
		require.TrueT(t, ok)
		assert.EqualT(t, "The photo urls for the pet.\nThis only accepts jpeg or png images.", prop.Description)
		if assert.NotNil(t, prop.Items) && assert.NotNil(t, prop.Items.Schema) {
			assert.EqualT(t, "\\.(jpe?g|png)$", prop.Items.Schema.Pattern)
		}

		assertProperty(t, &mod, "string", "status", "", "Status")
		prop, ok = mod.Properties["status"]
		assert.TrueT(t, ok)
		assert.EqualT(t, "The current status of the pet in the store.\navailable STATUS_AVAILABLE\npending STATUS_PENDING\nsold STATUS_SOLD", prop.Description)
		assert.Equal(t, []any{"available", "pending", "sold"}, prop.Enum)

		assertProperty(t, &mod, "string", "birthday", "date", "Birthday")
		prop, ok = mod.Properties["birthday"]
		assert.TrueT(t, ok)
		assert.EqualT(t, "The pet's birthday", prop.Description)

		assertArrayRef(t, &mod, "tags", "Tags", "#/definitions/tag")
		prop, ok = mod.Properties["tags"]
		assert.TrueT(t, ok)
		assert.EqualT(t, "Extra bits of information attached to this pet.", prop.Description)

		mod, ok = definitions["order"]
		assert.TrueT(t, ok)
		assert.Len(t, mod.Properties, 4)
		assert.Len(t, mod.Required, 3)

		assertProperty(t, &mod, "integer", "id", "int64", "ID")
		prop, ok = mod.Properties["id"]
		assert.TrueT(t, ok, "should have had an 'id' property")
		assert.EqualT(t, "the ID of the order", prop.Description)

		assertProperty(t, &mod, "integer", "userId", "int64", "UserID")
		prop, ok = mod.Properties["userId"]
		assert.TrueT(t, ok, "should have had an 'userId' property")
		assert.EqualT(t, "the id of the user who placed the order.", prop.Description)

		assertProperty(t, &mod, "string", "orderedAt", "date-time", "OrderedAt")
		prop, ok = mod.Properties["orderedAt"]
		assert.EqualT(t, "the time at which this order was made.", prop.Description)
		assert.TrueT(t, ok, "should have an 'orderedAt' property")

		assertArrayProperty(t, &mod, "object", "items", "", "Items")
		prop, ok = mod.Properties["items"]
		assert.TrueT(t, ok, "should have an 'items' slice")
		assert.NotNil(t, prop.Items, "items should have had an items property")
		assert.NotNil(t, prop.Items.Schema, "items.items should have had a schema property")

		itprop := prop.Items.Schema
		assert.Len(t, itprop.Properties, 2)
		assert.Len(t, itprop.Required, 2)

		assertProperty(t, itprop, "integer", "petId", "int64", "PetID")
		iprop, ok := itprop.Properties["petId"]
		assert.TrueT(t, ok, "should have had a 'petId' property")
		assert.EqualT(t, "the id of the pet to order", iprop.Description)

		assertProperty(t, itprop, "integer", "qty", "int32", "Quantity")
		iprop, ok = itprop.Properties["qty"]
		assert.TrueT(t, ok, "should have had a 'qty' property")
		assert.EqualT(t, "the quantity of this pet to order", iprop.Description)
		require.NotNil(t, iprop.Minimum)
		assert.InDeltaT(t, 1.00, *iprop.Minimum, epsilon)
	})
}

func verifyIDParam(t *testing.T, param spec.Parameter, description string) {
	t.Helper()

	assert.EqualT(t, description, param.Description)
	assert.EqualT(t, "path", param.In)
	assert.EqualT(t, "integer", param.Type)
	assert.EqualT(t, "int64", param.Format)
	assert.TrueT(t, param.Required)
	assert.Equal(t, "ID", param.Extensions["x-go-name"])
}

func verifyRefParam(t *testing.T, param spec.Parameter, description, refed string) {
	t.Helper()

	assert.EqualT(t, description, param.Description)
	assert.EqualT(t, "body", param.In)
	// TODO: this may fail sometimes (seen on go1.12 windows test): require pointer to be valid and avoid panicking
	require.NotNil(t, param)
	require.NotNil(t, param.Schema)
	assert.EqualT(t, "#/definitions/"+refed, param.Schema.Ref.String())
	assert.TrueT(t, param.Required)
}

type namedParam struct {
	Index int
	Name  string
}

type namedParams []namedParam

func (g namedParams) Len() int           { return len(g) }
func (g namedParams) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g namedParams) Less(i, j int) bool { return g[i].Name < g[j].Name }
