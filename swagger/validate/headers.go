package validate

import (
	"mime"

	"github.com/casualjim/go-swagger/swagger/util"
)

// ContentType validates the content type of a request
func ContentType(allowed []string, actual string) *Error {
	mt, _, err := mime.ParseMediaType(actual)
	if err != nil {
		return invalidContentType(actual, allowed)
	}
	if util.ContainsStringsCI(allowed, mt) {
		return nil
	}
	return invalidContentType(actual, allowed)
}
