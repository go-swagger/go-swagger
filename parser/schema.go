package parser

import (
	"fmt"
	"go/ast"
	"regexp"

	"github.com/casualjim/go-swagger/spec"
)

type structCommentParser struct {
	taggers      []*sectionTagger
	otherTags    []string
	headerParser *docCommentParser
	classifier   *regexp.Regexp
}

func (scp *structCommentParser) Parse(gofile *ast.File, target interface{}) error {
	for _, decl := range gofile.Decls {
		switch gd := decl.(type) {
		case *ast.GenDecl:
			if len(gd.Specs) > 0 {
				if ts, ok := gd.Specs[0].(*ast.TypeSpec); ok {
					// check if this struct has a +swagger:model tag in its comments
					// only the structs with that in their comments will be considered for parsing
					if scp.classifier.MatchString(ts.Doc.Text()) {
						schema := new(spec.Schema)
						schema.Type = spec.StringOrArray([]string{"object"})
						// analyze doc comment for the model
						// first line of the doc comment is the title
						// all following lines are description
						// Type is the type name for the struct
						// Default is a tag that accepts a json structure
						// Example is a tag that accepts a json structure
						// all other things are ignored and by definition added to the last matched tag unless
						// preceded by 2 new lines

						//parseDocComments(ts.Doc, schema, nil, nil)

						// analyze struct body for fields etc
						// each exported struct field:
						// * gets a type mapped to a go primitive
						// * perhaps gets a format
						// * has to document the validations that apply for the type and the field
						// * when the struct field points to a model it becomes a ref: #/definitions/ModelName
						// * the first line of the comment is the title
						// * the following lines are the description
					}
				}
			}
		default:
			fmt.Println("unhandled decl:", gd)
		}
	}
	return nil
}

func schemaParser() *structCommentParser {
	return &structCommentParser{}
}
