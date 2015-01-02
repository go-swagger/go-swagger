package httputils

import "strings"

// CanHaveBody returns true if this method can have a body
func CanHaveBody(method string) bool {
	mn := strings.ToUpper(method)
	return mn == "POST" || mn == "PUT" || mn == "PATCH"
}
