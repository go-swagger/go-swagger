package is

import "github.com/corbym/gocrest"

//False returns true if the actual matches false. Confusing but true.
func False() *gocrest.Matcher {
	return &gocrest.Matcher{
		Describe: "is false",
		Matches: func(actual interface{}) bool {
			return actual == false
		},
	}
}
