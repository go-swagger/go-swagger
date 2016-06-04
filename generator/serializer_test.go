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
					fmt.Println(res)
					assertInCode(t, "type Category struct {", res)
					assertInCode(t, "ID int64 `json:\"id,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"id\\\":\")", res)
					assertInCode(t, "out.Int64(m.ID)", res)
					assertInCode(t, "m.ID = 0", res)
					assertInCode(t, "m.ID = in.Int64()", res)
					assertNotInCode(t, "m.ID = in.Int64()()", res)
					assertInCode(t, "Name string `json:\"name,omitempty\"`", res)
					assertInCode(t, "out.RawString(\"\\\"name\\\":\")", res)
					assertInCode(t, "out.String(m.Name)", res)
					assertInCode(t, "m.Name = \"\"", res)
					assertInCode(t, "m.Name = in.String()", res)
					assertNotInCode(t, "m.Name = in.String()()", res)
				} else {
					fmt.Println(string(ct))
				}
			} else {
				fmt.Println(buf.String())
			}
		}
	}
}
