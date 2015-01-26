package errors

// Unauthenticated returns an unauthenticated error
func Unauthenticated(scheme string) Error {
	return New(401, "unauthenticated for %s", scheme)
}
