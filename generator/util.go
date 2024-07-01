package generator

import "strings"

func getCustomTagKeyMap(tagStr string) map[string]bool {
	res := map[string]bool{}
	for {
		output := strings.SplitN(tagStr, ":", 2)
		if len(output) != 2 {
			break
		}
		res[output[0]] = true
		output = strings.SplitN(output[1], " ", 2)
		if len(output) != 2 {
			break
		}
		tagStr = output[1]
	}

	return res
}
