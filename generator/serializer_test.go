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
		genModel, err := makeGenDefinition("Category", "models", schema, specDoc, true, true)
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
		genModel, err := makeGenDefinition("Category", "models", schema, specDoc, true, true)
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
					assertInCode(t, "out.String(\"id\")", res)
					assertInCode(t, "func (m *Category) idIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.Int64(m.ID)", res)
					assertInCode(t, "if err := m.idIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.ID = 0", res)
					assertInCode(t, "func (m *Category) idIReadJSON(in *jlexer.Lexer) (int64, error) {", res)
					assertInCode(t, "return in.Int64(), nil", res)
					assertInCode(t, "if idValue, err := m.idIReadJSON(in); err != nil", res)
					assertInCode(t, "m.ID = idValue", res)

					assertInCode(t, "Name string `json:\"name,omitempty\"`", res)
					assertInCode(t, "out.String(\"name\")", res)
					assertInCode(t, "func (m *Category) nameIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.String(m.Name)", res)
					assertInCode(t, "if err := m.nameIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.Name = \"\"", res)
					assertInCode(t, "func (m *Category) nameIReadJSON(in *jlexer.Lexer) (string, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if nameValue, err := m.nameIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Name = nameValue", res)

					assertInCode(t, "URL strfmt.URI `json:\"url,omitempty\"`", res)
					assertInCode(t, "out.String(\"url\")", res)
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
		genModel, err := makeGenDefinition("Categories", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.Equal(t, "Categories", genModel.Name)
			assert.Equal(t, "Categories", genModel.GoType)
			assert.Equal(t, "[]*Category", genModel.AliasedType)
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
		genModel, err := makeGenDefinition("Product", "models", schema, specDoc, true, true)
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
					assertInCode(t, "out.String(\"id\")", res)
					assertInCode(t, "func (m *Product) idIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.Int64(m.ID)", res)
					assertInCode(t, "if err := m.idIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.ID = 0", res)
					assertInCode(t, "func (m *Product) idIReadJSON(in *jlexer.Lexer) (int64, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if idValue, err := m.idIReadJSON(in); err != nil", res)
					assertInCode(t, "m.ID = idValue", res)

					assertInCode(t, "Name string `json:\"name,omitempty\"`", res)
					assertInCode(t, "out.String(\"name\")", res)
					assertInCode(t, "func (m *Product) nameIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.String(m.Name)", res)
					assertInCode(t, "if err := m.nameIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.Name = \"\"", res)
					assertInCode(t, "func (m *Product) nameIReadJSON(in *jlexer.Lexer) (string, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if nameValue, err := m.nameIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Name = nameValue", res)

					assertInCode(t, "Categories Categories `json:\"categories,omitempty\"`", res)
					assertInCode(t, "out.String(\"categories\")", res)
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
		genModel, err := makeGenDefinition("ProductLine", "models", schema, specDoc, true, true)
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
					assertInCode(t, "out.String(\"id\")", res)
					assertInCode(t, "func (m *ProductLine) idIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.Int64(m.ID)", res)
					assertInCode(t, "if err := m.idIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.ID = 0", res)
					assertInCode(t, "func (m *ProductLine) idIReadJSON(in *jlexer.Lexer) (int64, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if idValue, err := m.idIReadJSON(in); err != nil", res)
					assertInCode(t, "m.ID = idValue", res)

					assertInCode(t, "Name string `json:\"name,omitempty\"`", res)
					assertInCode(t, "out.String(\"name\")", res)
					assertInCode(t, "func (m *ProductLine) nameIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.String(m.Name)", res)
					assertInCode(t, "if err := m.nameIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.Name = \"\"", res)
					assertInCode(t, "func (m *ProductLine) nameIReadJSON(in *jlexer.Lexer) (string, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if nameValue, err := m.nameIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Name = nameValue", res)

					assertInCode(t, "Category *Category `json:\"category,omitempty\"`", res)
					assertInCode(t, "out.String(\"category\")", res)
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
		genModel, err := makeGenDefinition("Scores", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.Equal(t, "Scores", genModel.Name)
			assert.Equal(t, "Scores", genModel.GoType)
			assert.Equal(t, "[]float32", genModel.AliasedType)
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
		genModel, err := makeGenDefinition("JaggedScores", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.Equal(t, "JaggedScores", genModel.Name)
			assert.Equal(t, "JaggedScores", genModel.GoType)
			assert.Equal(t, "[][][]float32", genModel.AliasedType)
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
		genModel, err := makeGenDefinition("Notables", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.False(t, genModel.IsAnonymous)
			assert.Equal(t, "Notables", genModel.Name)
			assert.Equal(t, "Notables", genModel.GoType)
			assert.Equal(t, "[]*Notable", genModel.AliasedType)
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
		genModel, err := makeGenDefinition("Notablix", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.False(t, genModel.IsAnonymous)
			assert.Equal(t, "Notablix", genModel.Name)
			assert.Equal(t, "Notablix", genModel.GoType)
			assert.Equal(t, "[][][]*Notable", genModel.AliasedType)
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
		genModel, err := makeGenDefinition("WithComplexItems", "models", schema, specDoc, true, true)
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
		genModel, err := makeGenDefinition("WithItemsAndAdditional", "models", schema, specDoc, true, true)
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
		genModel, err := makeGenDefinition("WithItemsAndAdditional2", "models", schema, specDoc, true, true)
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
		genModel, err := makeGenDefinition("WithComplexAdditional", "models", schema, specDoc, true, true)
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
		genModel, err := makeGenDefinition("Age", "models", schema, specDoc, true, true)
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

func TestSerializer_Flag(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["flag"]
		genModel, err := makeGenDefinition("flag", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsPrimitive)
			assert.True(t, genModel.IsAliased)
			assert.Equal(t, "flag", genModel.Name)
			assert.Equal(t, "Flag", genModel.GoType)
			assert.Equal(t, "string", genModel.AliasedType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("flag.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type Flag string", res)
					assertInCode(t, "func (m Flag) MarshalJSON() ([]byte, error) {", res)
					assertInCode(t, "func (m *Flag) MarshalEasyJSON(out *jwriter.Writer) {", res)
					assertInCode(t, "out.String(m)", res)
					assertInCode(t, "*m = in.String()", res)
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

func TestSerializer_FlagsList(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["flags_list"]
		genModel, err := makeGenDefinition("flags_list", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.True(t, genModel.IsAliased)
			assert.Equal(t, "flags_list", genModel.Name)
			assert.Equal(t, "FlagsList", genModel.GoType)
			assert.Equal(t, "[]Flag", genModel.AliasedType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("flags_list.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type FlagsList []Flag", res)
					assertInCode(t, "func (m FlagsList) MarshalJSON() ([]byte, error) {", res)
					assertInCode(t, "func (m *FlagsList) MarshalEasyJSON(out *jwriter.Writer) {", res)
					assertInCode(t, "out.RawByte('[')", res)
					assertInCode(t, "out.RawByte(',')", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(m[i]))", res)
					assertInCode(t, "out.RawByte(']')", res)
					assertInCode(t, "iReadFn := func(in *jlexer.Lexer) (Flag, error)", res)
					assertInCode(t, "var flagsListValue Flag", res)
					assertInCode(t, "data := in.Raw(); in.Ok()", res)
					assertInCode(t, "err := swag.ReadJSON(data, &flagsListValue)", res)
					assertInCode(t, "var result []Flag", res)
					assertInCode(t, "in.Delim('[')", res)
					assertInCode(t, "result = make([]Flag, 0, 64)", res)
					assertInCode(t, "flagsListValue, err := iReadFn(in)", res)
					assertInCode(t, "in.Delim(']')", res)
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

func TestSerializer_ImageTar(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["ImageTar"]
		genModel, err := makeGenDefinition("ImageTar", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsStream)
			assert.True(t, genModel.IsAliased)
			assert.Equal(t, "ImageTar", genModel.Name)
			assert.Equal(t, "ImageTar", genModel.GoType)
			assert.Equal(t, "io.ReadCloser", genModel.AliasedType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("image_tar.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type ImageTar io.ReadCloser", res)
					assertNotInCode(t, "func (m ImageTar) MarshalJSON() ([]byte, error)", res)
					assertNotInCode(t, "func (m *ImageTar) UnmarshalJSON(data []byte) error", res)
					assertNotInCode(t, "func (m *ImageTar) MarshalEasyJSON(out *jwriter.Writer)", res)
					assertNotInCode(t, "func (m *ImageTar) UnmarshalEasyJSON(in *jlexer.Lexer)", res)
					assertNotInCode(t, "func (m *ImageTar) Validate(formats strfmt.Registry)", res)
				} else {
					fmt.Println(buf.String())
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_Tag(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["Tag"]
		genModel, err := makeGenDefinition("Tag", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.False(t, genModel.IsAliased)
			assert.Equal(t, "Tag", genModel.Name)
			assert.Equal(t, "Tag", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("image_tar.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type Tag struct", res)
					assertInCode(t, "func (m Tag) MarshalJSON() ([]byte, error)", res)
					assertInCode(t, "func (m *Tag) UnmarshalJSON(data []byte) error", res)
					assertInCode(t, "func (m *Tag) MarshalEasyJSON(out *jwriter.Writer)", res)
					assertInCode(t, "func (m *Tag) UnmarshalEasyJSON(in *jlexer.Lexer)", res)
					assertInCode(t, "func (m *Tag) FlagNameSet()", res)
					assertInCode(t, "func (m *Tag) FlagNameUnset()", res)
					assertInCode(t, "func (m *Tag) FlagNameNil()", res)
					assertInCode(t, "func (m *Tag) FlagNameZero()", res)
					assertInCode(t, "func (m *Tag) IsNameNil() bool", res)
					assertInCode(t, "func (m *Tag) IsNameSet() bool", res)
					assertInCode(t, "func (m *Tag) HasNameValue() bool", res)
					assertInCode(t, "func (m *Tag) SetName(value *string)", res)
					assertInCode(t, "func (m *Tag) ClearName()", res)
					assertInCode(t, "func (m *Tag) GetName() (value *string, null bool, haskey bool)", res)
					assertInCode(t, "func (m *Tag) GetNamePtr() *string", res)
					assertInCode(t, "out.String(\"name\")", res)
					assertInCode(t, "err := m.nameIWriteJSON(out)", res)
					assertInCode(t, "nameValue, err := m.nameIReadJSON(in)", res)
					assertInCode(t, "out.String(m.Name)", res)
					assertInCode(t, "return &in.String(), nil", res)
				} else {
					fmt.Println(buf.String())
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_Stats(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["Stats"]
		genModel, err := makeGenDefinition("Stats", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.True(t, genModel.IsAliased)
			assert.Equal(t, "Stats", genModel.Name)
			assert.Equal(t, "Stats", genModel.GoType)
			assert.Equal(t, "[]*StatsItems0", genModel.AliasedType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("stats.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type Stats []*StatsItems0", res)
					assertInCode(t, "func (m Stats) MarshalJSON() ([]byte, error)", res)
					assertInCode(t, "func (m *Stats) UnmarshalJSON(data []byte) error", res)
					assertInCode(t, "func (m *Stats) MarshalEasyJSON(out *jwriter.Writer)", res)
					assertInCode(t, "func (m *Stats) UnmarshalEasyJSON(in *jlexer.Lexer)", res)
					assertInCode(t, "statsValue, err := iReadFn(in)", res)
					assertInCode(t, "iReadFn := func(in *jlexer.Lexer) (*StatsItems0, error)", res)
					assertInCode(t, "err := swag.ReadJSON(data, &statsValue)", res)
					assertInCode(t, "var result []*StatsItems0", res)
					assertInCode(t, "result = make([]*StatsItems0, 0, 64)", res)

					assertInCode(t, "type StatsItems0 struct", res)
					assertInCode(t, "func (m StatsItems0) MarshalJSON() ([]byte, error)", res)
					assertInCode(t, "func (m *StatsItems0) UnmarshalJSON(data []byte) error", res)
					assertInCode(t, "func (m *StatsItems0) MarshalEasyJSON(out *jwriter.Writer)", res)
					assertInCode(t, "func (m *StatsItems0) UnmarshalEasyJSON(in *jlexer.Lexer)", res)
					assertInCode(t, "Points []int64 `json:\"points,omitempty\"`", res)
					assertInCode(t, "isPointsFieldNil bool `json:\"-\"`", res)
					assertInCode(t, "isPointsFieldSet bool `json:\"-\"`", res)
					assertInCode(t, "func (m *StatsItems0) pointsIWriteJSON(out *jwriter.Writer)", res)
					assertInCode(t, "func (m *StatsItems0) pointsIReadJSON(in *jlexer.Lexer) ([]int64, error)", res)
					assertInCode(t, "out.Int64(m.Points[i])", res)
					assertInCode(t, "wv, err := iReadFn(in)", res)
					assertInCode(t, "pointsValue, err := m.pointsIReadJSON(in)", res)
					assertInCode(t, "err := m.pointsIWriteJSON(out)", res)
				} else {
					fmt.Println(buf.String())
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_Statix(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["Statix"]
		genModel, err := makeGenDefinition("Statix", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsArray)
			assert.True(t, genModel.IsAliased)
			assert.Equal(t, "Statix", genModel.Name)
			assert.Equal(t, "Statix", genModel.GoType)
			assert.Equal(t, "[][][]*StatixItems0", genModel.AliasedType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("Statix.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type Statix [][][]*StatixItems0", res)
					assertInCode(t, "func (m Statix) MarshalJSON() ([]byte, error)", res)
					assertInCode(t, "func (m *Statix) UnmarshalJSON(data []byte) error", res)
					assertInCode(t, "func (m *Statix) MarshalEasyJSON(out *jwriter.Writer)", res)
					assertInCode(t, "func (m *Statix) UnmarshalEasyJSON(in *jlexer.Lexer)", res)
					assertInCode(t, "statixValue, err := iReadFn(in)", res)
					assertInCode(t, "iReadFn := func(in *jlexer.Lexer) (*StatixItems0, error)", res)
					assertInCode(t, "err := swag.ReadJSON(data, &statixValue)", res)
					assertInCode(t, "var result []*StatixItems0", res)
					assertInCode(t, "result = make([]*StatixItems0, 0, 64)", res)

					assertInCode(t, "type StatixItems0 struct", res)
					assertInCode(t, "func (m StatixItems0) MarshalJSON() ([]byte, error)", res)
					assertInCode(t, "func (m *StatixItems0) UnmarshalJSON(data []byte) error", res)
					assertInCode(t, "func (m *StatixItems0) MarshalEasyJSON(out *jwriter.Writer)", res)
					assertInCode(t, "func (m *StatixItems0) UnmarshalEasyJSON(in *jlexer.Lexer)", res)
					assertInCode(t, "Points []int64 `json:\"points,omitempty\"`", res)
					assertInCode(t, "isPointsFieldNil bool `json:\"-\"`", res)
					assertInCode(t, "isPointsFieldSet bool `json:\"-\"`", res)
					assertInCode(t, "func (m *StatixItems0) pointsIWriteJSON(out *jwriter.Writer)", res)
					assertInCode(t, "func (m *StatixItems0) pointsIReadJSON(in *jlexer.Lexer) ([]int64, error)", res)
					assertInCode(t, "out.Int64(m.Points[i])", res)
					assertInCode(t, "wv, err := iReadFn(in)", res)
					assertInCode(t, "pointsValue, err := m.pointsIReadJSON(in)", res)
					assertInCode(t, "err := m.pointsIWriteJSON(out)", res)
				} else {
					fmt.Println(buf.String())
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_WithItems(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["WithItems"]
		genModel, err := makeGenDefinition("WithItems", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.False(t, genModel.IsAliased)
			assert.Equal(t, "WithItems", genModel.Name)
			assert.Equal(t, "WithItems", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("with_items.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type WithItems struct", res)
					assertInCode(t, "func (m WithItems) MarshalJSON() ([]byte, error)", res)
					assertInCode(t, "func (m *WithItems) UnmarshalJSON(data []byte) error", res)
					assertInCode(t, "func (m *WithItems) MarshalEasyJSON(out *jwriter.Writer)", res)
					assertInCode(t, "func (m *WithItems) UnmarshalEasyJSON(in *jlexer.Lexer)", res)
					assertInCode(t, "func (m *WithItems) Validate(formats strfmt.Registry)", res)

					assertInCode(t, "Tags []string `json:\"tags,omitempty\"`", res)
					assertInCode(t, "out.String(\"tags\")", res)
					assertInCode(t, "func (m *WithItems) tagsIWriteJSON(out *jwriter.Writer) error", res)
					assertInCode(t, "out.String(m.Tags[i])", res)
					assertInCode(t, "if err := m.tagsIWriteJSON(out); err != nil", res)
					assertInCode(t, "m.Tags = nil", res)
					assertInCode(t, "func (m *WithItems) tagsIReadJSON(in *jlexer.Lexer) ([]string, error) {", res)
					assertInCode(t, "return in.String(), nil", res)
					assertInCode(t, "if tagsValue, err := m.tagsIReadJSON(in); err != nil", res)
					assertInCode(t, "m.Tags = tagsValue", res)

				} else {
					fmt.Println(buf.String())
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_Nota(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["Nota"]
		genModel, err := makeGenDefinition("Nota", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.False(t, genModel.IsComplexObject)
			assert.True(t, genModel.IsAliased)
			assert.True(t, genModel.IsMap)
			assert.Equal(t, "Nota", genModel.Name)
			assert.Equal(t, "Nota", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("nota.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type Nota map[string]int32", res)
					assertInCode(t, "func (m Nota) MarshalJSON() ([]byte, error)", res)
					assertInCode(t, "func (m Nota) UnmarshalJSON(data []byte) error", res)
					assertInCode(t, "func (m Nota) MarshalEasyJSON(out *jwriter.Writer)", res)
					assertInCode(t, "func (m Nota) UnmarshalEasyJSON(in *jlexer.Lexer)", res)
					assertInCode(t, "func (m Nota) Validate(formats strfmt.Registry)", res)

				} else {
					fmt.Println(buf.String())
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_NotaWithRef(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["NotaWithRef"]
		genModel, err := makeGenDefinition("NotaWithRef", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.False(t, genModel.IsComplexObject)
			assert.True(t, genModel.IsAliased)
			assert.True(t, genModel.IsMap)
			assert.Equal(t, "NotaWithRef", genModel.Name)
			assert.Equal(t, "NotaWithRef", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("nota_with_ref.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type NotaWithRef map[string]Notable", res)
					assertInCode(t, "func (m NotaWithRef) MarshalJSON() ([]byte, error)", res)
					assertInCode(t, "func (m NotaWithRef) UnmarshalJSON(data []byte) error", res)
					assertInCode(t, "func (m NotaWithRef) MarshalEasyJSON(out *jwriter.Writer)", res)
					assertInCode(t, "func (m NotaWithRef) UnmarshalEasyJSON(in *jlexer.Lexer)", res)
					assertInCode(t, "func (m NotaWithRef) Validate(formats strfmt.Registry)", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(map[string]Notable(m)))", res)
					assertInCode(t, "var notaWithRefValue map[string]Notable", res)
					assertInCode(t, "err := swag.ReadJSON(data, &notaWithRefValue)", res)
					assertInCode(t, "*m = notaWithRefValue", res)
				} else {
					fmt.Println(buf.String())
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_NotaWithMeta(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["NotaWithMeta"]
		genModel, err := makeGenDefinition("NotaWithMeta", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.False(t, genModel.IsComplexObject)
			assert.True(t, genModel.IsAliased)
			assert.True(t, genModel.IsMap)
			assert.Equal(t, "NotaWithMeta", genModel.Name)
			assert.Equal(t, "NotaWithMeta", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("nota_with_meta.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type NotaWithMeta map[string]NotaWithMetaAnon", res)
					assertInCode(t, "func (m NotaWithMeta) MarshalJSON() ([]byte, error)", res)
					assertInCode(t, "func (m NotaWithMeta) UnmarshalJSON(data []byte) error", res)
					assertInCode(t, "func (m NotaWithMeta) MarshalEasyJSON(out *jwriter.Writer)", res)
					assertInCode(t, "func (m NotaWithMeta) UnmarshalEasyJSON(in *jlexer.Lexer)", res)
					assertInCode(t, "func (m NotaWithMeta) Validate(formats strfmt.Registry)", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(map[string]NotaWithMetaAnon(m)))", res)
					assertInCode(t, "var notaWithMetaValue map[string]NotaWithMetaAnon", res)
					assertInCode(t, "err := swag.ReadJSON(data, &notaWithMetaValue)", res)
					assertInCode(t, "*m = notaWithMetaValue", res)

					assertInCode(t, "type NotaWithMetaAnon struct", res)
					assertInCode(t, "Comment *string `json:\"comment\"`", res)
					assertInCode(t, "out.String(\"comment\")", res)
					assertInCode(t, "m.Comment = nil", res)
					assertInCode(t, "commentValue, err := m.commentIReadJSON(in)", res)
					assertInCode(t, "err := m.commentIWriteJSON(out)", res)
					assertInCode(t, "m.IsCommentNil()", res)
					assertInCode(t, "out.RawString(\"null\")", res)
					assertInCode(t, "out.String(m.Comment)", res)
					assertInCode(t, "out.String(\"count\")", res)
					assertInCode(t, "err := m.countIWriteJSON(out)", res)
					assertInCode(t, "m.Count = 0", res)
					assertInCode(t, "commentValue, err := m.commentIReadJSON(in)", res)
					assertInCode(t, "out.Int32(m.Count)", res)
				} else {
					fmt.Println(buf.String())
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_NotaWithName(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["NotaWithName"]
		genModel, err := makeGenDefinition("NotaWithName", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.False(t, genModel.IsComplexObject)
			assert.True(t, genModel.IsAdditionalProperties)
			assert.True(t, genModel.HasAdditionalProperties)
			assert.True(t, genModel.IsAliased)
			assert.False(t, genModel.IsMap)
			assert.Equal(t, "NotaWithName", genModel.Name)
			assert.Equal(t, "NotaWithName", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("nota_with_name.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type NotaWithName struct", res)
					assertInCode(t, "func (m NotaWithName) MarshalJSON() ([]byte, error)", res)
					assertInCode(t, "func (m *NotaWithName) UnmarshalJSON(data []byte) error", res)
					assertInCode(t, "func (m *NotaWithName) MarshalEasyJSON(out *jwriter.Writer)", res)
					assertInCode(t, "func (m *NotaWithName) UnmarshalEasyJSON(in *jlexer.Lexer)", res)
					assertInCode(t, "func (m *NotaWithName) Validate(formats strfmt.Registry)", res)

					assertInCode(t, "Name *string `json:\"name\"`", res)
					assertInCode(t, "NotaWithName map[string]*int32 `json:\"-\"`", res)
					assertInCode(t, "out.String(\"name\")", res)
					assertInCode(t, "out.String(m.Name)", res)
					assertInCode(t, "return &in.String(), nil", res)
					assertInCode(t, "m.NotaWithName != nil", res)
					assertInCode(t, "out.Raw(swag.WriteJSON(m.NotaWithName))", res)
					assertInCode(t, "var notaWithNameValue map[string]*int32", res)
					assertInCode(t, "notaWithNameValue = make(map[string]*int32)", res)
					assertInCode(t, "notaWithNameValue[key] = notaWithNameEntry", res)
					assertInCode(t, "notaWithNameValue[key] = nil", res)
				} else {
					fmt.Println(buf.String())
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_NotaWithRefRegistry(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["NotaWithRefRegistry"]
		genModel, err := makeGenDefinition("NotaWithRefRegistry", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.False(t, genModel.IsComplexObject)
			assert.False(t, genModel.IsAdditionalProperties)
			assert.True(t, genModel.HasAdditionalProperties)
			assert.True(t, genModel.IsAliased)
			assert.True(t, genModel.IsMap)
			assert.Equal(t, "NotaWithRefRegistry", genModel.Name)
			assert.Equal(t, "NotaWithRefRegistry", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("nota_with_ref_registry.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type NotaWithRefRegistry map[string]map[string]map[string]Notable", res)
					assertInCode(t, "func (m NotaWithRefRegistry) MarshalJSON() ([]byte, error)", res)
					assertInCode(t, "func (m NotaWithRefRegistry) UnmarshalJSON(data []byte) error", res)
					assertInCode(t, "func (m NotaWithRefRegistry) MarshalEasyJSON(out *jwriter.Writer)", res)
					assertInCode(t, "func (m NotaWithRefRegistry) UnmarshalEasyJSON(in *jlexer.Lexer)", res)
					assertInCode(t, "func (m NotaWithRefRegistry) Validate(formats strfmt.Registry)", res)

					assertInCode(t, "out.Raw(swag.WriteJSON(map[string]map[string]map[string]Notable(m)))", res)
					assertInCode(t, "var notaWithRefRegistryValue map[string]map[string]map[string]Notable", res)
					assertInCode(t, "err := swag.ReadJSON(data, &notaWithRefRegistryValue); err != nil", res)
					assertInCode(t, "*m = notaWithRefRegistryValue", res)
				} else {
					fmt.Println(buf.String())
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}

func TestSerializer_NotaWithMetaRegistry(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["NotaWithMetaRegistry"]
		genModel, err := makeGenDefinition("NotaWithMetaRegistry", "models", schema, specDoc, true, true)
		if assert.NoError(t, err) {
			assert.False(t, genModel.IsComplexObject)
			assert.False(t, genModel.IsAdditionalProperties)
			assert.True(t, genModel.HasAdditionalProperties)
			assert.True(t, genModel.IsAliased)
			assert.True(t, genModel.IsMap)
			assert.Equal(t, "NotaWithMetaRegistry", genModel.Name)
			assert.Equal(t, "NotaWithMetaRegistry", genModel.GoType)
			// pretty.Println(genModel)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				ct, err := formatGoFile("nota_with_ref_registry.go", buf.Bytes())
				// fmt.Println(buf.String())
				if assert.NoError(t, err) {
					res := string(ct)
					// fmt.Println(res)
					assertInCode(t, "type NotaWithMetaRegistry map[string]map[string]map[string]NotaWithMetaRegistryAnon", res)
					assertInCode(t, "func (m NotaWithMetaRegistry) MarshalJSON() ([]byte, error)", res)
					assertInCode(t, "func (m NotaWithMetaRegistry) UnmarshalJSON(data []byte) error", res)
					assertInCode(t, "func (m NotaWithMetaRegistry) MarshalEasyJSON(out *jwriter.Writer)", res)
					assertInCode(t, "func (m NotaWithMetaRegistry) UnmarshalEasyJSON(in *jlexer.Lexer)", res)
					assertInCode(t, "func (m NotaWithMetaRegistry) Validate(formats strfmt.Registry)", res)

					assertInCode(t, "type NotaWithMetaRegistryAnon struct", res)
					assertInCode(t, "func (m NotaWithMetaRegistryAnon) MarshalJSON() ([]byte, error)", res)
					assertInCode(t, "func (m *NotaWithMetaRegistryAnon) UnmarshalJSON(data []byte) error", res)
					assertInCode(t, "func (m *NotaWithMetaRegistryAnon) MarshalEasyJSON(out *jwriter.Writer)", res)
					assertInCode(t, "func (m *NotaWithMetaRegistryAnon) UnmarshalEasyJSON(in *jlexer.Lexer)", res)
					assertInCode(t, "func (m *NotaWithMetaRegistryAnon) Validate(formats strfmt.Registry)", res)

					assertInCode(t, "Comment *string `json:\"comment\"`", res)
					assertInCode(t, "out.String(\"comment\")", res)
					assertInCode(t, "m.Comment = nil", res)
					assertInCode(t, "commentValue, err := m.commentIReadJSON(in)", res)
					assertInCode(t, "err := m.commentIWriteJSON(out)", res)
					assertInCode(t, "m.IsCommentNil()", res)
					assertInCode(t, "out.RawString(\"null\")", res)
					assertInCode(t, "out.String(m.Comment)", res)
					assertInCode(t, "out.String(\"count\")", res)
					assertInCode(t, "err := m.countIWriteJSON(out)", res)
					assertInCode(t, "m.Count = 0", res)
					assertInCode(t, "commentValue, err := m.commentIReadJSON(in)", res)
					assertInCode(t, "out.Int32(m.Count)", res)
				} else {
					fmt.Println(buf.String())
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}
