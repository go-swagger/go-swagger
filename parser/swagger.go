package parser

import "github.com/casualjim/go-swagger/spec"

func newSwaggerParser(specDoc *spec.Swagger) *docCommentParser {
	// this one should be called first in a pipeline if it's called in any sort of pipeline
	// people can specify a Consumes and Produces key which has a new content type
	// on each line
	// Schemes is a tag that is required and allow for
	// Host and BasePath can be specified but those values will be defaults,
	// they should get substituded when serving the swagger spec
	// Default parameters and responses are not supported at this stage
	// Tags are a mapping of tag to package
	return nil
}
