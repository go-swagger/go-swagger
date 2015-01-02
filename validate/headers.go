package validate

import (
	"mime"

	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/util"
)

// ContentType validates the content type of a request
func ContentType(allowed []string, actual string) *errors.Validation {
	mt, _, err := mime.ParseMediaType(actual)
	if err != nil {
		return errors.InvalidContentType(actual, allowed)
	}
	if util.ContainsStringsCI(allowed, mt) {
		return nil
	}
	return errors.InvalidContentType(actual, allowed)
}
