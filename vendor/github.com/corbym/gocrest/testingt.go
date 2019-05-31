package gocrest

//TestingT supplies a convenience interface that matches the testing.T interface.
type TestingT interface {
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	FailNow()
	Helper()
}
