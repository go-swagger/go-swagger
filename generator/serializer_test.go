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
		genModel, err := makeGenDefinition("Category", "models", schema, specDoc, true)
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("category.go", buf.Bytes())
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
					fmt.Println(buf.String())
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
		genModel, err := makeGenDefinition("Category", "models", schema, specDoc, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "Category", genModel.Name)
			assert.Equal(t, "Category", genModel.GoType)
			assert.True(t, genModel.IncludeValidator)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("category.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(ct)
					// pretty.Println(res)
					assertInCode(t, "type Category struct {", res)

					assertInCode(t, "ID int64 `json:\"id,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"id\\\":\")", res)
					assertInCode(t, "func (m *Category) idIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.Int64(m.ID)", res)
					assertInCode(t, "if err := m.idIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.ID = 0", res)
					assertInCode(t, "func (m *Category) idIReadJSON(in *jlexer.Lexer) (int64, error) {", res)
					assertInCode(t, "return in.Int64(), nil", res)
					assertInCode(t, "if idValue, err := m.idIReadJSON(in); err != nil", res)
					assertInCode(t, "m.ID = idValue", res)

					assertInCode(t, "Name string `json:\"name,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"name\\\":\")", res)
					assertInCode(t, "func (m *Category) nameIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.String(m.Name)", res)
					assertInCode(t, "if err := m.nameIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.Name = \"\"", res)
					assertInCode(t, "func (m *Category) nameIReadJSON(in *jlexer.Lexer) (string, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if nameValue, err := m.nameIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Name = nameValue", res)

					assertInCode(t, "URL strfmt.URI `json:\"url,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"url\\\":\")", res)
					assertInCode(t, "func (m *Category) urlIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(m.URL))", res)
					assertInCode(t, "if err := m.urlIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.URL = strfmt.URI(\"\")", res)
					assertInCode(t, "func (m *Category) urlIReadJSON(in *jlexer.Lexer) (strfmt.URI, error) {", res)
					assertInCode(t, "data := in.Raw(); in.Ok()", res)
					assertInCode(t, "var urlValue strfmt.URI", res)
					assertInCode(t, "return urlValue, nil", res)
					assertInCode(t, "return strfmt.URI(\"\"), err", res)
					assertInCode(t, "if urlValue, err := m.urlIReadJSON(in); err != nil", res)
					assertInCode(t, "swag.ReadJSON(data, &urlValue)", res)
					assertInCode(t, "m.URL = urlValue", res)
					assertNotInCode(t, "m.URL = in.String()()", res)
				} else {
					fmt.Println(buf.String())
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
		genModel, err := makeGenDefinition("Categories", "models", schema, specDoc, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.Equal(t, "Categories", genModel.Name)
			assert.Equal(t, "[]*Category", genModel.GoType)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("categories.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(ct)
					assertInCode(t, "type Categories []*Category", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(m[i]))", res)
					assertInCode(t, "var result []*Category", res)
					assertInCode(t, "iReadFn := func(in *jlexer.Lexer) (*Category, error)", res)
					assertInCode(t, "var categoriesValue Category", res)
					assertInCode(t, "if data := in.Raw(); in.Ok() {", res)
					assertInCode(t, "if err := swag.ReadJSON(data, &categoriesValue); err != nil", res)
					assertInCode(t, "categoriesValue, err := iReadFn(in)", res)
				} else {
					fmt.Println(buf.String())
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
		genModel, err := makeGenDefinition("Product", "models", schema, specDoc, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "Product", genModel.Name)
			assert.Equal(t, "Product", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("product.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type Product struct {", res)

					assertInCode(t, "ID int64 `json:\"id,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"id\\\":\")", res)
					assertInCode(t, "func (m *Product) idIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.Int64(m.ID)", res)
					assertInCode(t, "if err := m.idIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.ID = 0", res)
					assertInCode(t, "func (m *Product) idIReadJSON(in *jlexer.Lexer) (int64, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if idValue, err := m.idIReadJSON(in); err != nil", res)
					assertInCode(t, "m.ID = idValue", res)

					assertInCode(t, "Name string `json:\"name,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"name\\\":\")", res)
					assertInCode(t, "func (m *Product) nameIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.String(m.Name)", res)
					assertInCode(t, "if err := m.nameIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.Name = \"\"", res)
					assertInCode(t, "func (m *Product) nameIReadJSON(in *jlexer.Lexer) (string, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if nameValue, err := m.nameIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Name = nameValue", res)

					assertInCode(t, "Categories Categories `json:\"categories,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"categories\\\":\")", res)
					assertInCode(t, "func (m *Product) categoriesIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(m.Categories))", res)
					assertInCode(t, "if err := m.categoriesIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.Categories = nil", res)
					assertInCode(t, "func (m *Product) categoriesIReadJSON(in *jlexer.Lexer) (Categories, error) {", res)
					assertInCode(t, "var categoriesValue Categories", res)
					assertInCode(t, "if data := in.Raw(); in.Ok()", res)
					assertInCode(t, "if err := swag.ReadJSON(data, &categoriesValue)", res)
					assertInCode(t, "if categoriesValue, err := m.categoriesIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Categories = categoriesValue", res)
				} else {
					fmt.Println(buf.String())
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
		genModel, err := makeGenDefinition("ProductLine", "models", schema, specDoc, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "ProductLine", genModel.Name)
			assert.Equal(t, "ProductLine", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("product_line.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type ProductLine struct {", res)

					assertInCode(t, "ID int64 `json:\"id,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"id\\\":\")", res)
					assertInCode(t, "func (m *ProductLine) idIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.Int64(m.ID)", res)
					assertInCode(t, "if err := m.idIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.ID = 0", res)
					assertInCode(t, "func (m *ProductLine) idIReadJSON(in *jlexer.Lexer) (int64, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if idValue, err := m.idIReadJSON(in); err != nil", res)
					assertInCode(t, "m.ID = idValue", res)

					assertInCode(t, "Name string `json:\"name,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"name\\\":\")", res)
					assertInCode(t, "func (m *ProductLine) nameIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.String(m.Name)", res)
					assertInCode(t, "if err := m.nameIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.Name = \"\"", res)
					assertInCode(t, "func (m *ProductLine) nameIReadJSON(in *jlexer.Lexer) (string, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if nameValue, err := m.nameIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Name = nameValue", res)

					assertInCode(t, "Category *Category `json:\"category,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"category\\\":\")", res)
					assertInCode(t, "func (m *ProductLine) categoryIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(m.Category))", res)
					assertInCode(t, "if err := m.categoryIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.Category = nil", res)
					assertInCode(t, "func (m *ProductLine) categoryIReadJSON(in *jlexer.Lexer) (*Category, error) {", res)
					assertInCode(t, "var categoryValue Category", res)
					assertInCode(t, "if data := in.Raw(); in.Ok()", res)
					assertInCode(t, "if err := swag.ReadJSON(data, &categoryValue)", res)
					assertInCode(t, "if categoryValue, err := m.categoryIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Category = categoryValue", res)
				} else {
					fmt.Println(buf.String())
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
		genModel, err := makeGenDefinition("Scores", "models", schema, specDoc, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.Equal(t, "Scores", genModel.Name)
			assert.Equal(t, "[]float32", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("scores.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type Scores []float32", res)
					assertInCode(t, "out.Float32(m[i])", res)
					assertInCode(t, "for i := range m {", res)
					assertInCode(t, "var result []float32", res)
					assertInCode(t, "iReadFn := func(in *jlexer.Lexer) (float32, error)", res)
					assertInCode(t, "return in.Float32(), nil", res)
					assertInCode(t, "for !in.IsDelim(']')", res)
					assertInCode(t, "result = make([]float32, 0, 64)", res)
					assertInCode(t, "scoresValue, err := iReadFn(in)", res)
				} else {
					fmt.Println(buf.String())
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
		genModel, err := makeGenDefinition("JaggedScores", "models", schema, specDoc, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.Equal(t, "JaggedScores", genModel.Name)
			assert.Equal(t, "[][][]float32", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("JaggedScores.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type JaggedScores [][][]float32", res)
					assertInCode(t, "out.Float32(m[i][ii][iii])", res)
					assertInCode(t, "for iii := range m[i][ii]", res)
					assertInCode(t, "for ii := range m[i]", res)
					assertInCode(t, "for i := range m {", res)
					assertInCode(t, "var result []float32", res)
					assertInCode(t, "iiiReadFn := func(in *jlexer.Lexer) (float32, error)", res)
					assertInCode(t, "iiReadFn := func(in *jlexer.Lexer) ([]float32, error)", res)
					assertInCode(t, "iReadFn := func(in *jlexer.Lexer) ([][]float32, error)", res)
					assertInCode(t, "return in.Float32(), nil", res)
					assertInCode(t, "for !in.IsDelim(']')", res)
					assertInCode(t, "result = make([]float32, 0, 64)", res)
					assertInCode(t, "result = make([][]float32, 0, 64)", res)
					assertInCode(t, "result = make([][][]float32, 0, 64)", res)
					assertInCode(t, "jaggedScoresValue, err := iReadFn(in)", res)
					assertInCode(t, "wv, err := iiReadFn(in)", res)
					assertInCode(t, "wv, err := iiiReadFn(in)", res)
				} else {
					fmt.Println(buf.String())
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
		genModel, err := makeGenDefinition("Notables", "models", schema, specDoc, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.False(t, genModel.IsAnonymous)
			assert.Equal(t, "Notables", genModel.Name)
			assert.Equal(t, "[]*Notable", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("product_line.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type Notables []*Notable", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(m[i]))", res)
					assertInCode(t, "iReadFn := func(in *jlexer.Lexer) (*Notable, error)", res)
					assertInCode(t, "var notablesValue Notable", res)
					assertInCode(t, "result = make([]*Notable, 0, 64)", res)
					assertInCode(t, "notablesValue, err := iReadFn(in)", res)
				} else {
					fmt.Println(buf.String())
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
		genModel, err := makeGenDefinition("Notablix", "models", schema, specDoc, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.False(t, genModel.IsAnonymous)
			assert.Equal(t, "Notablix", genModel.Name)
			assert.Equal(t, "[][][]*Notable", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("product_line.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type Notablix [][][]*Notable", res)
					assertInCode(t, "for ii := range m[i]", res)
					assertInCode(t, "for iii := range m[i][ii]", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(m[i][ii][iii]))", res)
					assertInCode(t, "iReadFn := func(in *jlexer.Lexer) ([][]*Notable, error)", res)
					assertInCode(t, "iiReadFn := func(in *jlexer.Lexer) ([]*Notable, error)", res)
					assertInCode(t, "iiiReadFn := func(in *jlexer.Lexer) (*Notable, error)", res)
					assertInCode(t, "var notablixValue Notable", res)
					assertInCode(t, "var result []*Notable", res)
					assertInCode(t, "var result [][]*Notable", res)
					assertInCode(t, "var result [][][]*Notable", res)
					assertInCode(t, "result = make([]*Notable, 0, 64)", res)
					assertInCode(t, "result = make([][]*Notable, 0, 64)", res)
					assertInCode(t, "result = make([][][]*Notable, 0, 64)", res)
					assertInCode(t, "notablixValue, err := iReadFn(in)", res)
				} else {
					fmt.Println(buf.String())
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
		genModel, err := makeGenDefinition("WithComplexItems", "models", schema, specDoc, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "WithComplexItems", genModel.Name)
			assert.Equal(t, "WithComplexItems", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("product_line.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type WithComplexItems struct {", res)
					assertInCode(t, "Tags []*WithComplexItemsTagsItems0 `json:\"tags,omitempty\"`", res)
					assertInCode(t, "if err := m.tagsIWriteJSON(out); err != nil", res)
					assertInCode(t, "if tagsValue, err := m.tagsIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Tags = nil", res)
					assertInCode(t, "func (m *WithComplexItems) tagsIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "for i := range m.Tags", res)
					assertInCode(t, "func (m *WithComplexItems) tagsIReadJSON(in *jlexer.Lexer) ([]*WithComplexItemsTagsItems0, error)", res)
					assertInCode(t, "iReadFn := func(in *jlexer.Lexer) (*WithComplexItemsTagsItems0, error)", res)
					assertInCode(t, "var tagsValue WithComplexItemsTagsItems0", res)
					assertInCode(t, "if data := in.Raw(); in.Ok()", res)
					assertInCode(t, "if err := swag.ReadJSON(data, &tagsValue); err != nil", res)
					assertInCode(t, "result = make([]*WithComplexItemsTagsItems0, 0, 64)", res)
					assertInCode(t, "wv, err := iReadFn(in)", res)
					assertInCode(t, "type WithComplexItemsTagsItems0 struct", res)
					assertInCode(t, "Points []int64 `json:\"points,omitempty\"`", res)
					assertInCode(t, "if err := m.pointsIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.Points = nil", res)
					assertInCode(t, "if pointsValue, err := m.pointsIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Points = pointsValue", res)
					assertInCode(t, "func (m *WithComplexItemsTagsItems0) pointsIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "func (m *WithComplexItemsTagsItems0) pointsIReadJSON(in *jlexer.Lexer) ([]int64, error)", res)
				} else {
					fmt.Println(buf.String())
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_WithItemsAndAdditional(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["WithItemsAndAdditional"]
		genModel, err := makeGenDefinition("WithItemsAndAdditional", "models", schema, specDoc, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "WithItemsAndAdditional", genModel.Name)
			assert.Equal(t, "WithItemsAndAdditional", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("product_line.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type WithItemsAndAdditional struct {", res)
					assertInCode(t, "Tags *WithItemsAndAdditionalTagsTuple0 `json:\"tags,omitempty\"`", res)
					assertInCode(t, "if err := m.tagsIWriteJSON(out); err != nil", res)
					assertInCode(t, "if tagsValue, err := m.tagsIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Tags = nil", res)
					assertInCode(t, "m.Tags = tagsValue", res)
					assertInCode(t, "func (m *WithItemsAndAdditional) tagsIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(m.Tags))", res)
					assertInCode(t, "func (m *WithItemsAndAdditional) tagsIReadJSON(in *jlexer.Lexer) (*WithItemsAndAdditionalTagsTuple0, error)", res)
					assertInCode(t, "var withItemsAndAdditionalTagsTuple0Value WithItemsAndAdditionalTagsTuple0", res)
					assertInCode(t, "if data := in.Raw(); in.Ok()", res)
					assertInCode(t, "if err := swag.ReadJSON(data, &withItemsAndAdditionalTagsTuple0Value); err != nil", res)
					assertInCode(t, "type WithItemsAndAdditionalTagsTuple0 struct", res)
					assertInCode(t, "P0 *string `json:\"-\"`", res)
					assertInCode(t, "WithItemsAndAdditionalTagsTuple0Items []interface{} `json:\"-\"`", res)
					assertInCode(t, "if err := m.p0IWriteJSON(out); err != nil", res)
					assertInCode(t, "case \"P0\":", res)
					assertInCode(t, "m.P0 = nil", res)
					assertInCode(t, "case \"WithItemsAndAdditionalTagsTuple0Items\":", res)
					assertInCode(t, "m.P0 = p0Value", res)
					assertInCode(t, "m.withItemsAndAdditionalTagsTuple0ItemsIWriteJSON", res)
					assertInCode(t, "withItemsAndAdditionalTagsTuple0ItemsValue, err := m.withItemsAndAdditionalTagsTuple0ItemsIReadJSON(in)", res)
					assertInCode(t, "func (m *WithItemsAndAdditionalTagsTuple0) withItemsAndAdditionalTagsTuple0ItemsIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "func (m *WithItemsAndAdditionalTagsTuple0) withItemsAndAdditionalTagsTuple0IReadJSON(in *jlexer.Lexer) (*WithItemsAndAdditionalTagsTuple0, error)", res)
					// assertInCode(t, "", res)
				} else {
					fmt.Println(buf.String())
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
		genModel, err := makeGenDefinition("WithItemsAndAdditional2", "models", schema, specDoc, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "WithItemsAndAdditional2", genModel.Name)
			assert.Equal(t, "WithItemsAndAdditional2", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("product_line.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type WithItemsAndAdditional2 struct {", res)
					assertInCode(t, "Tags *WithItemsAndAdditional2TagsTuple0 `json:\"tags,omitempty\"`", res)
					assertInCode(t, "if err := m.tagsIWriteJSON(out); err != nil", res)
					assertInCode(t, "if tagsValue, err := m.tagsIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Tags = nil", res)
					assertInCode(t, "m.Tags = tagsValue", res)
					assertInCode(t, "func (m *WithItemsAndAdditional2) tagsIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(m.Tags))", res)
					assertInCode(t, "func (m *WithItemsAndAdditional2) tagsIReadJSON(in *jlexer.Lexer) (*WithItemsAndAdditional2TagsTuple0, error)", res)
					assertInCode(t, "var withItemsAndAdditional2TagsTuple0Value WithItemsAndAdditional2TagsTuple0", res)
					assertInCode(t, "if data := in.Raw(); in.Ok()", res)
					assertInCode(t, "if err := swag.ReadJSON(data, &withItemsAndAdditional2TagsTuple0Value); err != nil", res)
					assertInCode(t, "type WithItemsAndAdditional2TagsTuple0 struct", res)
					assertInCode(t, "P0 *string `json:\"-\"`", res)
					assertInCode(t, "WithItemsAndAdditional2TagsTuple0Items []int32 `json:\"-\"`", res)
					assertInCode(t, "if err := m.p0IWriteJSON(out); err != nil", res)
					assertInCode(t, "case \"P0\":", res)
					assertInCode(t, "m.P0 = nil", res)
					assertInCode(t, "case \"WithItemsAndAdditional2TagsTuple0Items\":", res)
					assertInCode(t, "m.P0 = p0Value", res)
					// assertInCode(t, "", res)
				} else {
					fmt.Println(buf.String())
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_WithComplexAdditional(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["WithComplexAdditional"]
		genModel, err := makeGenDefinition("WithComplexAdditional", "models", schema, specDoc, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "WithComplexAdditional", genModel.Name)
			assert.Equal(t, "WithComplexAdditional", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("product_line.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type WithComplexAdditional struct {", res)
					assertInCode(t, "Tags *WithComplexAdditionalTagsTuple0 `json:\"tags,omitempty\"`", res)
					assertInCode(t, "if err := m.tagsIWriteJSON(out); err != nil", res)
					assertInCode(t, "if tagsValue, err := m.tagsIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Tags = nil", res)
					assertInCode(t, "m.Tags = tagsValue", res)
					assertInCode(t, "func (m *WithComplexAdditional) tagsIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(m.Tags))", res)
					assertInCode(t, "func (m *WithComplexAdditional) tagsIReadJSON(in *jlexer.Lexer) (*WithComplexAdditionalTagsTuple0, error)", res)
					assertInCode(t, "var tagsValue WithComplexAdditionalTagsTuple0", res)
					assertInCode(t, "if data := in.Raw(); in.Ok()", res)
					assertInCode(t, "if err := swag.ReadJSON(data, &tagsValue); err != nil", res)
					assertInCode(t, "type WithComplexAdditionalTagsTuple0 struct", res)
					assertInCode(t, "Points []int64 `json:\"points,omitempty\"`", res)
					assertInCode(t, "m.pointsIWriteJSON(out)", res)
					assertInCode(t, "pointsValue, err := m.pointsIReadJSON(in)", res)
					assertInCode(t, "for i := range m.Points", res)
					assertInCode(t, "out.Int64(m.Points[i])", res)
					assertInCode(t, "result = make([]int64, 0, 64)", res)
					assertInCode(t, "P0 *string `json:\"-\"`", res)
					assertInCode(t, "WithComplexAdditionalTagsTuple0Items []*WithComplexAdditionalTagsItems `json:\"-\"`", res)
					assertInCode(t, "if err := m.p0IWriteJSON(out); err != nil", res)
					assertInCode(t, "case \"P0\":", res)
					assertInCode(t, "m.P0 = nil", res)
					assertInCode(t, "case \"WithComplexAdditionalTagsTuple0Items\":", res)
					assertInCode(t, "m.P0 = p0Value", res)
					// assertInCode(t, "", res)
				} else {
					fmt.Println(buf.String())
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_Age(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["Age"]
		genModel, err := makeGenDefinition("Age", "models", schema, specDoc, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsPrimitive)
			assert.True(t, genModel.IsAliased)
			assert.Equal(t, "Age", genModel.Name)
			assert.Equal(t, "Age", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("product_line.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type Age int32", res)
					assertInCode(t, "func (m Age) MarshalJSON() ([]byte, error) {", res)
					assertInCode(t, "func (m *Age) MarshalEasyJSON(out *jwriter.Writer) {", res)
					assertInCode(t, "out.Int32(m)", res)
					assertInCode(t, "*m = in.Int32()", res)
					// assertInCode(t, "", res)
				} else {
					fmt.Println(buf.String())
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}
