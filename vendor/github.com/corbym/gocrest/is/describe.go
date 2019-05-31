package is

import (
	"fmt"
	"github.com/corbym/gocrest"
)

func describe(matchers []*gocrest.Matcher, conjunction string) string {
	var description string
	for x := 0; x < len(matchers); x++ {
		description += matchers[x].Describe
		if x+1 < len(matchers) {
			description += fmt.Sprintf(" %s ", conjunction)
		}
	}
	return description
}
