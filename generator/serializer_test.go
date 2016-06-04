package generator

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
)

func TestSerializer_Methods(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.serializers.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["Category"]
		genModel, err := makeGenDefinition("Category", "models", schema, specDoc)
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, _ := formatGoFile("category.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(ct)
					assertInCode(t, "type Category struct {", res)
					assertInCode(t, "func (m Category) MarshalJSON() ([]byte, error) {", res)
					assertInCode(t, "func (m *Category) MarshalEasyJSON(out *jwriter.Writer) {", res)
					assertInCode(t, "out := jwriter.Writer{}", res)
					assertInCode(t, "m.MarshalEasyJSON(&out)", res)
					assertInCode(t, "return out.BuildBytes()", res)
					assertInCode(t, "func (m *Category) UnmarshalJSON(data []byte) error {", res)
					assertInCode(t, "func (m *Category) UnmarshalEasyJSON(in *jlexer.Lexer) {", res)
					assertInCode(t, "in := jlexer.Lexer{Data: data}", res)
					assertInCode(t, "m.UnmarshalEasyJSON(&in)", res)
					assertInCode(t, "return in.Error()", res)

				} else {
					fmt.Println(string(ct))
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_Category(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.serializers.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["Category"]
		genModel, err := makeGenDefinition("Category", "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "Category", genModel.Name)
			assert.Equal(t, "Category", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, _ := formatGoFile("category.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type Category struct {", res)
					assertInCode(t, "ID int64 `json:\"id,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"id\\\":\")", res)
					assertInCode(t, "idWriteFn := func(value int64, out *jwriter.Writer) error", res)
					assertInCode(t, "out.Int64(value)", res)
					assertInCode(t, "if err := idWriteFn(m.ID, out); err != nil", res)
					assertInCode(t, "m.ID = 0", res)
					assertInCode(t, "idValueFn := func(in *jlexer.Lexer) (int64, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if idValue, err := idValueFn(in); err != nil", res)
					assertInCode(t, "m.ID = idValue", res)
					assertNotInCode(t, "m.ID = in.Int64()()", res)

					assertInCode(t, "Name string `json:\"name,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"name\\\":\")", res)
					assertInCode(t, "nameWriteFn := func(value string, out *jwriter.Writer) error", res)
					assertInCode(t, "out.String(value)", res)
					assertInCode(t, "if err := nameWriteFn(m.Name, out); err != nil", res)
					assertInCode(t, "m.Name = \"\"", res)
					assertInCode(t, "nameValueFn := func(in *jlexer.Lexer) (string, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if nameValue, err := nameValueFn(in); err != nil", res)
					assertInCode(t, "m.Name = nameValue", res)
					assertNotInCode(t, "m.Name = in.String()()", res)

					assertInCode(t, "URL strfmt.URI `json:\"url,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"url\\\":\")", res)
					assertInCode(t, "urlWriteFn := func(value strfmt.URI, out *jwriter.Writer) error", res)
					assertInCode(t, "out.String(value)", res)
					assertInCode(t, "if err := urlWriteFn(m.URL, out); err != nil", res)
					assertInCode(t, "m.URL = strfmt.URI(\"\")", res)
					assertInCode(t, "urlValueFn := func(in *jlexer.Lexer) (strfmt.URI, error) {", res)
					assertInCode(t, "data := in.Raw(); in.Ok()", res)
					assertInCode(t, "var result strfmt.URI", res)
					assertInCode(t, "return result, nil", res)
					assertInCode(t, "if urlValue, err := urlValueFn(in); err != nil", res)
					assertInCode(t, "m.URL = urlValue", res)
					assertNotInCode(t, "m.URL = in.String()()", res)
				} else {
					fmt.Println(string(ct))
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_Categories(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.serializers.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["Categories"]
		genModel, err := makeGenDefinition("Categories", "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.Equal(t, "Categories", genModel.Name)
			assert.Equal(t, "[]*Category", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, _ := formatGoFile("categories.go", buf.Bytes())
				// fmt.Println(string(ct))
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type Categories []*Category", res)

					assertInCode(t, "categoriesItemsFn := func(value *Category, out *jwriter.Writer) error {", res)
					assertInCode(t, "b, err := swag.WriteJSON(value)", res)
					assertInCode(t, "out.Raw(b, nil)", res)
					assertInCode(t, "err := categoriesItemsFn(v, out); err != nil", res)

					assertInCode(t, "var result []*Category", res)
					assertInCode(t, "iReadFn := func(in *jlexer.Lexer) (*Category, error)", res)
					assertInCode(t, "var result Category", res)
					assertInCode(t, "if data := in.Raw(); in.Ok() {", res)
					assertInCode(t, "if err := swag.ReadJSON(data, &result); err != nil", res)
					assertInCode(t, "wv, err := iReadFn(in)", res)
				} else {
					fmt.Println(string(ct))
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_Product(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.serializers.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["Product"]
		genModel, err := makeGenDefinition("Product", "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "Product", genModel.Name)
			assert.Equal(t, "Product", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, _ := formatGoFile("product.go", buf.Bytes())
				// fmt.Println(string(ct))
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type Product struct {", res)

					assertInCode(t, "ID int64 `json:\"id,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"id\\\":\")", res)
					assertInCode(t, "idWriteFn := func(value int64, out *jwriter.Writer) error", res)
					assertInCode(t, "out.Int64(value)", res)
					assertInCode(t, "if err := idWriteFn(m.ID, out); err != nil", res)
					assertInCode(t, "m.ID = 0", res)
					assertInCode(t, "idValueFn := func(in *jlexer.Lexer) (int64, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if idValue, err := idValueFn(in); err != nil", res)
					assertInCode(t, "m.ID = idValue", res)

					assertInCode(t, "Name string `json:\"name,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"name\\\":\")", res)
					assertInCode(t, "nameWriteFn := func(value string, out *jwriter.Writer) error", res)
					assertInCode(t, "out.String(value)", res)
					assertInCode(t, "if err := nameWriteFn(m.Name, out); err != nil", res)
					assertInCode(t, "m.Name = \"\"", res)
					assertInCode(t, "nameValueFn := func(in *jlexer.Lexer) (string, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if nameValue, err := nameValueFn(in); err != nil", res)
					assertInCode(t, "m.Name = nameValue", res)

					assertInCode(t, "Categories Categories `json:\"categories,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"categories\\\":\")", res)
					assertInCode(t, "categoriesWriteFn := func(value Categories, out *jwriter.Writer) error", res)
					assertInCode(t, "b, err := swag.WriteJSON(value)", res)
					assertInCode(t, "out.Raw(b, nil)", res)
					assertInCode(t, "if err := categoriesWriteFn(m.Categories, out); err != nil", res)
					assertInCode(t, "m.Categories = nil", res)
					assertInCode(t, "categoriesValueFn := func(in *jlexer.Lexer) (Categories, error) {", res)
					assertInCode(t, "var result Categories", res)
					assertInCode(t, "if data := in.Raw(); in.Ok()", res)
					assertInCode(t, "if err := swag.ReadJSON(data, &result)", res)
					assertInCode(t, "if categoriesValue, err := categoriesValueFn(in); err != nil", res)
					assertInCode(t, "m.Categories = categoriesValue", res)
				} else {
					fmt.Println(string(ct))
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_ProductLine(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.serializers.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["ProductLine"]
		genModel, err := makeGenDefinition("ProductLine", "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "ProductLine", genModel.Name)
			assert.Equal(t, "ProductLine", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, _ := formatGoFile("product_line.go", buf.Bytes())
				// fmt.Println(string(ct))
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type ProductLine struct {", res)

					assertInCode(t, "ID int64 `json:\"id,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"id\\\":\")", res)
					assertInCode(t, "idWriteFn := func(value int64, out *jwriter.Writer) error", res)
					assertInCode(t, "out.Int64(value)", res)
					assertInCode(t, "if err := idWriteFn(m.ID, out); err != nil", res)
					assertInCode(t, "m.ID = 0", res)
					assertInCode(t, "idValueFn := func(in *jlexer.Lexer) (int64, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if idValue, err := idValueFn(in); err != nil", res)
					assertInCode(t, "m.ID = idValue", res)

					assertInCode(t, "Name string `json:\"name,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"name\\\":\")", res)
					assertInCode(t, "nameWriteFn := func(value string, out *jwriter.Writer) error", res)
					assertInCode(t, "out.String(value)", res)
					assertInCode(t, "if err := nameWriteFn(m.Name, out); err != nil", res)
					assertInCode(t, "m.Name = \"\"", res)
					assertInCode(t, "nameValueFn := func(in *jlexer.Lexer) (string, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if nameValue, err := nameValueFn(in); err != nil", res)
					assertInCode(t, "m.Name = nameValue", res)

					assertInCode(t, "Category *Category `json:\"category,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"category\\\":\")", res)
					assertInCode(t, "categoryWriteFn := func(value *Category, out *jwriter.Writer) error", res)
					assertInCode(t, "b, err := swag.WriteJSON(value)", res)
					assertInCode(t, "out.Raw(b, nil)", res)
					assertInCode(t, "if err := categoryWriteFn(m.Category, out); err != nil", res)
					assertInCode(t, "m.Category = nil", res)
					assertInCode(t, "categoryValueFn := func(in *jlexer.Lexer) (*Category, error) {", res)
					assertInCode(t, "var result Category", res)
					assertInCode(t, "if data := in.Raw(); in.Ok()", res)
					assertInCode(t, "if err := swag.ReadJSON(data, &result)", res)
					assertInCode(t, "if categoryValue, err := categoryValueFn(in); err != nil", res)
					assertInCode(t, "m.Category = categoryValue", res)
				} else {
					fmt.Println(string(ct))
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}
