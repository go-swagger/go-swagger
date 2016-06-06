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
					assertInCode(t, "func (m *Category) idIWriteJSON(value int64, out *jwriter.Writer) error", res)
					assertInCode(t, "out.Int64(value)", res)
					assertInCode(t, "if err := m.idIWriteJSON(m.ID, out); err != nil", res)
					assertInCode(t, "m.ID = 0", res)
					assertInCode(t, "func (m *Category) idIReadJSON(in *jlexer.Lexer) (int64, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if idValue, err := m.idIReadJSON(in); err != nil", res)
					assertInCode(t, "m.ID = idValue", res)

					assertInCode(t, "Name string `json:\"name,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"name\\\":\")", res)
					assertInCode(t, "func (m *Category) nameIWriteJSON(value string, out *jwriter.Writer) error", res)
					assertInCode(t, "out.String(value)", res)
					assertInCode(t, "if err := m.nameIWriteJSON(m.Name, out); err != nil", res)
					assertInCode(t, "m.Name = \"\"", res)
					assertInCode(t, "func (m *Category) nameIReadJSON(in *jlexer.Lexer) (string, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if nameValue, err := m.nameIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Name = nameValue", res)

					assertInCode(t, "URL strfmt.URI `json:\"url,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"url\\\":\")", res)
					assertInCode(t, "func (m *Category) urlIWriteJSON(value strfmt.URI, out *jwriter.Writer) error", res)
					assertInCode(t, "out.String(value)", res)
					assertInCode(t, "if err := m.urlIWriteJSON(m.URL, out); err != nil", res)
					assertInCode(t, "m.URL = strfmt.URI(\"\")", res)
					assertInCode(t, "func (m *Category) urlIReadJSON(in *jlexer.Lexer) (strfmt.URI, error) {", res)
					assertInCode(t, "data := in.Raw(); in.Ok()", res)
					assertInCode(t, "var result strfmt.URI", res)
					assertInCode(t, "return result, nil", res)
					assertInCode(t, "if urlValue, err := m.urlIReadJSON(in); err != nil", res)
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
					assertInCode(t, "out.Raw(swag.WriteJSON(value))", res)
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
					assertInCode(t, "func (m *Product) idIWriteJSON(value int64, out *jwriter.Writer) error", res)
					assertInCode(t, "out.Int64(value)", res)
					assertInCode(t, "if err := m.idIWriteJSON(m.ID, out); err != nil", res)
					assertInCode(t, "m.ID = 0", res)
					assertInCode(t, "func (m *Product) idIReadJSON(in *jlexer.Lexer) (int64, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if idValue, err := m.idIReadJSON(in); err != nil", res)
					assertInCode(t, "m.ID = idValue", res)

					assertInCode(t, "Name string `json:\"name,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"name\\\":\")", res)
					assertInCode(t, "func (m *Product) nameIWriteJSON(value string, out *jwriter.Writer) error", res)
					assertInCode(t, "out.String(value)", res)
					assertInCode(t, "if err := m.nameIWriteJSON(m.Name, out); err != nil", res)
					assertInCode(t, "m.Name = \"\"", res)
					assertInCode(t, "func (m *Product) nameIReadJSON(in *jlexer.Lexer) (string, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if nameValue, err := m.nameIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Name = nameValue", res)

					assertInCode(t, "Categories Categories `json:\"categories,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"categories\\\":\")", res)
					assertInCode(t, "func (m *Product) categoriesIWriteJSON(value Categories, out *jwriter.Writer) error", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(value))", res)
					assertInCode(t, "if err := m.categoriesIWriteJSON(m.Categories, out); err != nil", res)
					assertInCode(t, "m.Categories = nil", res)
					assertInCode(t, "func (m *Product) categoriesIReadJSON(in *jlexer.Lexer) (Categories, error) {", res)
					assertInCode(t, "var result Categories", res)
					assertInCode(t, "if data := in.Raw(); in.Ok()", res)
					assertInCode(t, "if err := swag.ReadJSON(data, &result)", res)
					assertInCode(t, "if categoriesValue, err := m.categoriesIReadJSON(in); err != nil", res)
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
					assertInCode(t, "func (m *ProductLine) idIWriteJSON(value int64, out *jwriter.Writer) error", res)
					assertInCode(t, "out.Int64(value)", res)
					assertInCode(t, "if err := m.idIWriteJSON(m.ID, out); err != nil", res)
					assertInCode(t, "m.ID = 0", res)
					assertInCode(t, "func (m *ProductLine) idIReadJSON(in *jlexer.Lexer) (int64, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if idValue, err := m.idIReadJSON(in); err != nil", res)
					assertInCode(t, "m.ID = idValue", res)

					assertInCode(t, "Name string `json:\"name,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"name\\\":\")", res)
					assertInCode(t, "func (m *ProductLine) nameIWriteJSON(value string, out *jwriter.Writer) error", res)
					assertInCode(t, "out.String(value)", res)
					assertInCode(t, "if err := m.nameIWriteJSON(m.Name, out); err != nil", res)
					assertInCode(t, "m.Name = \"\"", res)
					assertInCode(t, "func (m *ProductLine) nameIReadJSON(in *jlexer.Lexer) (string, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if nameValue, err := m.nameIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Name = nameValue", res)

					assertInCode(t, "Category *Category `json:\"category,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"category\\\":\")", res)
					assertInCode(t, "func (m *ProductLine) categoryIWriteJSON(value *Category, out *jwriter.Writer) error", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(value))", res)
					assertInCode(t, "if err := m.categoryIWriteJSON(m.Category, out); err != nil", res)
					assertInCode(t, "m.Category = nil", res)
					assertInCode(t, "func (m *ProductLine) categoryIReadJSON(in *jlexer.Lexer) (*Category, error) {", res)
					assertInCode(t, "var result Category", res)
					assertInCode(t, "if data := in.Raw(); in.Ok()", res)
					assertInCode(t, "if err := swag.ReadJSON(data, &result)", res)
					assertInCode(t, "if categoryValue, err := m.categoryIReadJSON(in); err != nil", res)
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

func TestSerializer_Scores(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["Scores"]
		genModel, err := makeGenDefinition("Scores", "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.Equal(t, "Scores", genModel.Name)
			assert.Equal(t, "[]float32", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, _ := formatGoFile("scores.go", buf.Bytes())
				// fmt.Println(string(ct))
				if assert.NoError(t, err) {
					res := string(ct)
					fmt.Println(res)
					assertInCode(t, "type Scores []float32", res)

					assertInCode(t, "scoresItemsFn := func(value float32, out *jwriter.Writer) error {", res)
					assertInCode(t, "out.Float32(value)", res)
					assertInCode(t, "err := scoresItemsFn(v, out); err != nil", res)
					assertInCode(t, "for i, v := range m {", res)
					assertInCode(t, "var result []float32", res)
					assertInCode(t, "iReadFn := func(in *jlexer.Lexer) (float32, error)", res)
					assertInCode(t, "return in.Float32(), nil", res)
					assertInCode(t, "for !in.IsDelim(']')", res)
					assertInCode(t, "result = make([]float32, 0, 64)", res)
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

func TestSerializer_JaggedScores(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["JaggedScores"]
		genModel, err := makeGenDefinition("JaggedScores", "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.Equal(t, "JaggedScores", genModel.Name)
			assert.Equal(t, "[][][]float32", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, _ := formatGoFile("JaggedScores.go", buf.Bytes())
				// fmt.Println(string(ct))
				if assert.NoError(t, err) {
					res := string(ct)
					fmt.Println(res)
					assertInCode(t, "type JaggedScores [][][]float32", res)

					assertInCode(t, "jaggedScoresItemsFn := func(value [][]float32, out *jwriter.Writer) error {", res)
					assertInCode(t, "iiFn := func(value []float32, out *jwriter.Writer) error {", res)
					assertInCode(t, "iiiFn := func(value float32, out *jwriter.Writer) error {", res)
					assertInCode(t, "out.Float32(value)", res)
					assertInCode(t, "err := jaggedScoresItemsFn(v, out); err != nil", res)
					assertInCode(t, "if err := iiFn(v, out); err != nil", res)
					assertInCode(t, "if err := iiiFn(v, out); err != nil", res)
					assertInCode(t, "for i, v := range m {", res)
					assertInCode(t, "var result []float32", res)
					assertInCode(t, "iiiReadFn := func(in *jlexer.Lexer) (float32, error)", res)
					assertInCode(t, "iiReadFn := func(in *jlexer.Lexer) ([]float32, error)", res)
					assertInCode(t, "iReadFn := func(in *jlexer.Lexer) ([][]float32, error)", res)

					assertInCode(t, "return in.Float32(), nil", res)
					assertInCode(t, "for !in.IsDelim(']')", res)
					assertInCode(t, "result = make([]float32, 0, 64)", res)
					assertInCode(t, "result = make([][]float32, 0, 64)", res)
					assertInCode(t, "result = make([][][]float32, 0, 64)", res)
					assertInCode(t, "wv, err := iReadFn(in)", res)
					assertInCode(t, "wv, err := iiReadFn(in)", res)
					assertInCode(t, "wv, err := iiiReadFn(in)", res)
				} else {
					fmt.Println(string(ct))
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_Notables(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["Notables"]
		genModel, err := makeGenDefinition("Notables", "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.False(t, genModel.IsAnonymous)
			assert.Equal(t, "Notables", genModel.Name)
			assert.Equal(t, "[]*Notable", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, _ := formatGoFile("product_line.go", buf.Bytes())
				// fmt.Println(string(ct))
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type Notables []*Notable", res)
					assertInCode(t, "notablesItemsFn := func(value *Notable, out *jwriter.Writer) error", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(value))", res)
					assertInCode(t, "if err := notablesItemsFn(v, out); err != nil", res)
					assertInCode(t, "iReadFn := func(in *jlexer.Lexer) (*Notable, error)", res)
					assertInCode(t, "var result Notable", res)
					assertInCode(t, "var result []*Notable", res)
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

func TestSerializer_Notablix(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["Notablix"]
		genModel, err := makeGenDefinition("Notablix", "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.False(t, genModel.IsAnonymous)
			assert.Equal(t, "Notablix", genModel.Name)
			assert.Equal(t, "[][][]*Notable", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, _ := formatGoFile("product_line.go", buf.Bytes())
				// fmt.Println(string(ct))
				if assert.NoError(t, err) {
					res := string(ct)
					fmt.Println(res)
					assertInCode(t, "type Notablix [][][]*Notable", res)
					assertInCode(t, "notablixItemsFn := func(value [][]*Notable, out *jwriter.Writer) error", res)
					assertInCode(t, "iiFn := func(value []*Notable, out *jwriter.Writer) error", res)
					assertInCode(t, "iiiFn := func(value *Notable, out *jwriter.Writer) error", res)
					assertInCode(t, "if err := iiiFn(v, out); err != nil", res)
					assertInCode(t, "if err := iiFn(v, out); err != nil", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(value))", res)
					assertInCode(t, "if err := notablixItemsFn(v, out); err != nil", res)
					assertInCode(t, "iReadFn := func(in *jlexer.Lexer) ([][]*Notable, error)", res)
					assertInCode(t, "iiReadFn := func(in *jlexer.Lexer) ([]*Notable, error)", res)
					assertInCode(t, "iiiReadFn := func(in *jlexer.Lexer) (*Notable, error)", res)
					assertInCode(t, "var result Notable", res)
					assertInCode(t, "var result []*Notable", res)
					assertInCode(t, "var result [][]*Notable", res)
					assertInCode(t, "var result [][][]*Notable", res)
					assertInCode(t, "result = make([]*Notable, 0, 64)", res)
					assertInCode(t, "result = make([][]*Notable, 0, 64)", res)
					assertInCode(t, "result = make([][][]*Notable, 0, 64)", res)
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

func TestSerializer_WithComplexItems(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["WithComplexItems"]
		genModel, err := makeGenDefinition("WithComplexItems", "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "WithComplexItems", genModel.Name)
			assert.Equal(t, "WithComplexItems", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, _ := formatGoFile("product_line.go", buf.Bytes())
				// fmt.Println(string(ct))
				if assert.NoError(t, err) {
					res := string(ct)
					fmt.Println(res)
					assertInCode(t, "type WithComplexItems struct {", res)
					assertInCode(t, "Tags []*WithComplexItemsTagsItems0 `json:\"tags,omitempty\"`", res)
					assertInCode(t, "if err := m.tagsIWriteJSON(m.Tags, out); err != nil", res)
					assertInCode(t, "if tagsValue, err := m.tagsIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Tags = nil", res)
					assertInCode(t, "func (m *WithComplexItems) tagsIWriteJSON(value []*WithComplexItemsTagsItems0, out *jwriter.Writer) error", res)
					assertInCode(t, "iFn := func(value *WithComplexItemsTagsItems0, out *jwriter.Writer) error", res)
					assertInCode(t, "if err := iFn(v, out); err != nil", res)
					assertInCode(t, "func (m *WithComplexItems) tagsIReadJSON(in *jlexer.Lexer) ([]*WithComplexItemsTagsItems0, error)", res)
					assertInCode(t, "iReadFn := func(in *jlexer.Lexer) (*WithComplexItemsTagsItems0, error)", res)
					assertInCode(t, "var result WithComplexItemsTagsItems0", res)
					assertInCode(t, "if data := in.Raw(); in.Ok()", res)
					assertInCode(t, "if err := swag.ReadJSON(data, &result); err != nil", res)
					assertInCode(t, "result = make([]*WithComplexItemsTagsItems0, 0, 64)", res)
					assertInCode(t, "wv, err := iReadFn(in)", res)
					assertInCode(t, "type WithComplexItemsTagsItems0 struct", res)
					assertInCode(t, "Points []int64 `json:\"points,omitempty\"`", res)
					assertInCode(t, "if err := m.pointsIWriteJSON(m.Points, out); err != nil", res)
					assertInCode(t, "m.Points = nil", res)
					assertInCode(t, "if pointsValue, err := m.pointsIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Points = pointsValue", res)
					assertInCode(t, "func (m *WithComplexItemsTagsItems0) pointsIWriteJSON(value []int64, out *jwriter.Writer) error", res)
					assertInCode(t, "func (m *WithComplexItemsTagsItems0) pointsIReadJSON(in *jlexer.Lexer) ([]int64, error)", res)
				} else {
					fmt.Println(string(ct))
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_WithItemsAndAdditional2(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["WithItemsAndAdditional2"]
		genModel, err := makeGenDefinition("WithItemsAndAdditional2", "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "WithItemsAndAdditional2", genModel.Name)
			assert.Equal(t, "WithItemsAndAdditional2", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, _ := formatGoFile("product_line.go", buf.Bytes())
				// fmt.Println(string(ct))
				if assert.NoError(t, err) {
					res := string(ct)
					fmt.Println(res)
					assertInCode(t, "type WithItemsAndAdditional2 struct {", res)
					assertInCode(t, "Tags *WithItemsAndAdditional2TagsTuple0 `json:\"tags,omitempty\"`", res)
					assertInCode(t, "if err := m.tagsIWriteJSON(m.Tags, out); err != nil", res)
					assertInCode(t, "if tagsValue, err := m.tagsIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Tags = nil", res)
					assertInCode(t, "m.Tags = tagsValue", res)
					assertInCode(t, "func (m *WithItemsAndAdditional2) tagsIWriteJSON(value *WithItemsAndAdditional2TagsTuple0, out *jwriter.Writer) error", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(value))", res)
					assertInCode(t, "func (m *WithItemsAndAdditional2) tagsIReadJSON(in *jlexer.Lexer) (*WithItemsAndAdditional2TagsTuple0, error)", res)
					assertInCode(t, "var result WithItemsAndAdditional2TagsTuple0", res)
					assertInCode(t, "if data := in.Raw(); in.Ok()", res)
					assertInCode(t, "if err := swag.ReadJSON(data, &result); err != nil", res)
					assertInCode(t, "type WithItemsAndAdditional2TagsTuple0 struct", res)
					assertInCode(t, "P0 *string `json:\"-\"`", res)
					assertInCode(t, "WithItemsAndAdditional2TagsTuple0Items []int32 `json:\"-\"`", res)
					assertInCode(t, "if err := m.p0IWriteJSON(m.P0, out); err != nil", res)
					assertInCode(t, "case \"P0\":", res)
					assertInCode(t, "m.P0 = nil", res)
					assertInCode(t, "case \"WithItemsAndAdditional2TagsTuple0Items\":", res)
					assertInCode(t, "m.P0 = p0Value", res)
					// assertInCode(t, "", res)
				} else {
					fmt.Println(string(ct))
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}
