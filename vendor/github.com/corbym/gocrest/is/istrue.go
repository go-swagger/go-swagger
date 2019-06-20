package is

import "github.com/corbym/gocrest"

//True returns true if the actual matches true
func True() *gocrest.Matcher {
	return &gocrest.Matcher{
		Describe: "is true",
		Matches: func(actual interface{}) bool {
			return actual == true
		},
	}
}
