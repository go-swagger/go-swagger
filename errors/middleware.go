package errors

import (
	"bytes"
	"fmt"
	"strings"
)

// APIVerificationFailed is an error that contains all the missing info for a mismatched section
// between the api registrations and the api spec
type APIVerificationFailed struct {
	Section              string
	MissingSpecification []string
	MissingRegistration  []string
}

//
func (v *APIVerificationFailed) Error() string {
	buf := bytes.NewBuffer(nil)

	hasRegMissing := len(v.MissingRegistration) > 0
	hasSpecMissing := len(v.MissingSpecification) > 0

	if hasRegMissing {
		buf.WriteString(fmt.Sprintf("missing [%s] %s registrations", strings.Join(v.MissingRegistration, ", "), v.Section))
	}

	if hasRegMissing && hasSpecMissing {
		buf.WriteString("\n")
	}

	if hasSpecMissing {
		buf.WriteString(fmt.Sprintf("missing from spec file [%s] %s", strings.Join(v.MissingSpecification, ", "), v.Section))
	}

	return buf.String()
}
