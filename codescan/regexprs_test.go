package codescan

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaValueExtractors(t *testing.T) {
	strfmts := []string{
		"// swagger:strfmt ",
		"* swagger:strfmt ",
		"* swagger:strfmt ",
		" swagger:strfmt ",
		"swagger:strfmt ",
		"// swagger:strfmt    ",
		"* swagger:strfmt     ",
		"* swagger:strfmt    ",
		" swagger:strfmt     ",
		"swagger:strfmt      ",
	}
	models := []string{
		"// swagger:model ",
		"* swagger:model ",
		"* swagger:model ",
		" swagger:model ",
		"swagger:model ",
		"// swagger:model    ",
		"* swagger:model     ",
		"* swagger:model    ",
		" swagger:model     ",
		"swagger:model      ",
	}

	allOf := []string{
		"// swagger:allOf ",
		"* swagger:allOf ",
		"* swagger:allOf ",
		" swagger:allOf ",
		"swagger:allOf ",
		"// swagger:allOf    ",
		"* swagger:allOf     ",
		"* swagger:allOf    ",
		" swagger:allOf     ",
		"swagger:allOf      ",
	}

	parameters := []string{
		"// swagger:parameters ",
		"* swagger:parameters ",
		"* swagger:parameters ",
		" swagger:parameters ",
		"swagger:parameters ",
		"// swagger:parameters    ",
		"* swagger:parameters     ",
		"* swagger:parameters    ",
		" swagger:parameters     ",
		"swagger:parameters      ",
	}

	validParams := []string{
		"yada123",
		"date",
		"date-time",
		"long-combo-1-with-combo-2-and-a-3rd-one-too",
	}
	invalidParams := []string{
		"1-yada-3",
		"1-2-3",
		"-yada-3",
		"-2-3",
		"*blah",
		"blah*",
	}

	verifySwaggerOneArgSwaggerTag(t, rxStrFmt, strfmts, validParams, append(invalidParams, "", "  ", " "))
	verifySwaggerOneArgSwaggerTag(t, rxModelOverride, models, append(validParams, "", "  ", " "), invalidParams)

	verifySwaggerOneArgSwaggerTag(t, rxAllOf, allOf, append(validParams, "", "  ", " "), invalidParams)

	verifySwaggerMultiArgSwaggerTag(t, rxParametersOverride, parameters, validParams, invalidParams)

	verifyMinMax(t, rxf(rxMinimumFmt, ""), "min", []string{"", ">", "="})
	verifyMinMax(t, rxf(rxMinimumFmt, fmt.Sprintf(rxItemsPrefixFmt, 1)), "items.min", []string{"", ">", "="})
	verifyMinMax(t, rxf(rxMaximumFmt, ""), "max", []string{"", "<", "="})
	verifyMinMax(t, rxf(rxMaximumFmt, fmt.Sprintf(rxItemsPrefixFmt, 1)), "items.max", []string{"", "<", "="})
	verifyNumeric2Words(t, rxf(rxMultipleOfFmt, ""), "multiple", "of")
	verifyNumeric2Words(t, rxf(rxMultipleOfFmt, fmt.Sprintf(rxItemsPrefixFmt, 1)), "items.multiple", "of")

	verifyIntegerMinMaxManyWords(t, rxf(rxMinLengthFmt, ""), "min", []string{"len", "length"})
	// pattern
	extraSpaces := []string{"", " ", "  ", "     "}
	prefixes := []string{"//", "*", ""}
	patArgs := []string{"^\\w+$", "[A-Za-z0-9-.]*"}
	patNames := []string{"pattern", "Pattern"}
	for _, pref := range prefixes {
		for _, es1 := range extraSpaces {
			for _, nm := range patNames {
				for _, es2 := range extraSpaces {
					for _, es3 := range extraSpaces {
						for _, arg := range patArgs {
							line := strings.Join([]string{pref, es1, nm, es2, ":", es3, arg}, "")
							matches := rxf(rxPatternFmt, "").FindStringSubmatch(line)
							assert.Len(t, matches, 2)
							assert.Equal(t, arg, matches[1])
						}
					}
				}
			}
		}
	}

	verifyIntegerMinMaxManyWords(t, rxf(rxMinItemsFmt, ""), "min", []string{"items"})
	verifyBoolean(t, rxf(rxUniqueFmt, ""), []string{"unique"}, nil)

	verifyBoolean(t, rxReadOnly, []string{"read"}, []string{"only"})
	verifyBoolean(t, rxRequired, []string{"required"}, nil)
}

func makeMinMax(lower string) (res []string) {
	for _, a := range []string{"", "imum"} {
		res = append(res, lower+a, strings.Title(lower)+a)
	}
	return
}

func verifyBoolean(t *testing.T, matcher *regexp.Regexp, names, names2 []string) {
	extraSpaces := []string{"", " ", "  ", "     "}
	prefixes := []string{"//", "*", ""}
	validArgs := []string{"true", "false"}
	invalidArgs := []string{"TRUE", "FALSE", "t", "f", "1", "0", "True", "False", "true*", "false*"}
	var nms []string
	for _, nm := range names {
		nms = append(nms, nm, strings.Title(nm))
	}

	var nms2 []string
	for _, nm := range names2 {
		nms2 = append(nms2, nm, strings.Title(nm))
	}

	var rnms []string
	if len(nms2) > 0 {
		for _, nm := range nms {
			for _, es := range append(extraSpaces, "-") {
				for _, nm2 := range nms2 {
					rnms = append(rnms, strings.Join([]string{nm, es, nm2}, ""))
				}
			}
		}
	} else {
		rnms = nms
	}

	var cnt int
	for _, pref := range prefixes {
		for _, es1 := range extraSpaces {
			for _, nm := range rnms {
				for _, es2 := range extraSpaces {
					for _, es3 := range extraSpaces {
						for _, vv := range validArgs {
							line := strings.Join([]string{pref, es1, nm, es2, ":", es3, vv}, "")
							matches := matcher.FindStringSubmatch(line)
							assert.Len(t, matches, 2)
							assert.Equal(t, vv, matches[1])
							cnt++
						}
						for _, iv := range invalidArgs {
							line := strings.Join([]string{pref, es1, nm, es2, ":", es3, iv}, "")
							matches := matcher.FindStringSubmatch(line)
							assert.Empty(t, matches)
							cnt++
						}
					}
				}
			}
		}
	}
	var nm2 string
	if len(names2) > 0 {
		nm2 = " " + names2[0]
	}
	var Debug = os.Getenv("DEBUG") != "" || os.Getenv("SWAGGER_DEBUG") != ""
	if Debug {
		fmt.Printf("tested %d %s%s combinations\n", cnt, names[0], nm2)
	}
}

func verifyIntegerMinMaxManyWords(t *testing.T, matcher *regexp.Regexp, name1 string, words []string) {
	extraSpaces := []string{"", " ", "  ", "     "}
	prefixes := []string{"//", "*", ""}
	validNumericArgs := []string{"0", "1234"}
	invalidNumericArgs := []string{"1A3F", "2e10", "*12", "12*", "-1235", "0.0", "1234.0394", "-2948.484"}

	var names []string
	for _, w := range words {
		names = append(names, w, strings.Title(w))
	}

	var cnt int
	for _, pref := range prefixes {
		for _, es1 := range extraSpaces {
			for _, nm1 := range makeMinMax(name1) {
				for _, es2 := range append(extraSpaces, "-") {
					for _, nm2 := range names {
						for _, es3 := range extraSpaces {
							for _, es4 := range extraSpaces {
								for _, vv := range validNumericArgs {
									line := strings.Join([]string{pref, es1, nm1, es2, nm2, es3, ":", es4, vv}, "")
									matches := matcher.FindStringSubmatch(line)
									assert.Len(t, matches, 2)
									assert.Equal(t, vv, matches[1])
									cnt++
								}
								for _, iv := range invalidNumericArgs {
									line := strings.Join([]string{pref, es1, nm1, es2, nm2, es3, ":", es4, iv}, "")
									matches := matcher.FindStringSubmatch(line)
									assert.Empty(t, matches)
									cnt++
								}
							}
						}
					}
				}
			}
		}
	}
	var nm2 string
	if len(words) > 0 {
		nm2 = " " + words[0]
	}
	var Debug = os.Getenv("DEBUG") != "" || os.Getenv("SWAGGER_DEBUG") != ""
	if Debug {
		fmt.Printf("tested %d %s%s combinations\n", cnt, name1, nm2)

	}
}

func verifyNumeric2Words(t *testing.T, matcher *regexp.Regexp, name1, name2 string) {
	extraSpaces := []string{"", " ", "  ", "     "}
	prefixes := []string{"//", "*", ""}
	validNumericArgs := []string{"0", "1234", "-1235", "0.0", "1234.0394", "-2948.484"}
	invalidNumericArgs := []string{"1A3F", "2e10", "*12", "12*"}

	var cnt int
	for _, pref := range prefixes {
		for _, es1 := range extraSpaces {
			for _, es2 := range extraSpaces {
				for _, es3 := range extraSpaces {
					for _, es4 := range extraSpaces {
						for _, vv := range validNumericArgs {
							lines := []string{
								strings.Join([]string{pref, es1, name1, es2, name2, es3, ":", es4, vv}, ""),
								strings.Join([]string{pref, es1, strings.Title(name1), es2, strings.Title(name2), es3, ":", es4, vv}, ""),
								strings.Join([]string{pref, es1, strings.Title(name1), es2, name2, es3, ":", es4, vv}, ""),
								strings.Join([]string{pref, es1, name1, es2, strings.Title(name2), es3, ":", es4, vv}, ""),
							}
							for _, line := range lines {
								matches := matcher.FindStringSubmatch(line)
								assert.Len(t, matches, 2)
								assert.Equal(t, vv, matches[1])
								cnt++
							}
						}
						for _, iv := range invalidNumericArgs {
							lines := []string{
								strings.Join([]string{pref, es1, name1, es2, name2, es3, ":", es4, iv}, ""),
								strings.Join([]string{pref, es1, strings.Title(name1), es2, strings.Title(name2), es3, ":", es4, iv}, ""),
								strings.Join([]string{pref, es1, strings.Title(name1), es2, name2, es3, ":", es4, iv}, ""),
								strings.Join([]string{pref, es1, name1, es2, strings.Title(name2), es3, ":", es4, iv}, ""),
							}
							for _, line := range lines {
								matches := matcher.FindStringSubmatch(line)
								assert.Empty(t, matches)
								cnt++
							}
						}
					}
				}
			}
		}
	}
	var Debug = os.Getenv("DEBUG") != "" || os.Getenv("SWAGGER_DEBUG") != ""
	if Debug {
		fmt.Printf("tested %d %s %s combinations\n", cnt, name1, name2)
	}
}

func verifyMinMax(t *testing.T, matcher *regexp.Regexp, name string, operators []string) {
	extraSpaces := []string{"", " ", "  ", "     "}
	prefixes := []string{"//", "*", ""}
	validNumericArgs := []string{"0", "1234", "-1235", "0.0", "1234.0394", "-2948.484"}
	invalidNumericArgs := []string{"1A3F", "2e10", "*12", "12*"}

	var cnt int
	for _, pref := range prefixes {
		for _, es1 := range extraSpaces {
			for _, wrd := range makeMinMax(name) {
				for _, es2 := range extraSpaces {
					for _, es3 := range extraSpaces {
						for _, op := range operators {
							for _, es4 := range extraSpaces {
								for _, vv := range validNumericArgs {
									line := strings.Join([]string{pref, es1, wrd, es2, ":", es3, op, es4, vv}, "")
									matches := matcher.FindStringSubmatch(line)
									// fmt.Printf("matching %q with %q, matches (%d): %v\n", line, matcher, len(matches), matches)
									assert.Len(t, matches, 3)
									assert.Equal(t, vv, matches[2])
									cnt++
								}
								for _, iv := range invalidNumericArgs {
									line := strings.Join([]string{pref, es1, wrd, es2, ":", es3, op, es4, iv}, "")
									matches := matcher.FindStringSubmatch(line)
									assert.Empty(t, matches)
									cnt++
								}
							}
						}
					}
				}
			}
		}
	}
	var Debug = os.Getenv("DEBUG") != "" || os.Getenv("SWAGGER_DEBUG") != ""
	if Debug {
		fmt.Printf("tested %d %s combinations\n", cnt, name)
	}
}

func verifySwaggerOneArgSwaggerTag(t *testing.T, matcher *regexp.Regexp, prefixes, validParams, invalidParams []string) {
	for _, pref := range prefixes {
		for _, param := range validParams {
			line := pref + param
			matches := matcher.FindStringSubmatch(line)
			if assert.Len(t, matches, 2) {
				assert.Equal(t, strings.TrimSpace(param), matches[1])
			}
		}
	}

	for _, pref := range prefixes {
		for _, param := range invalidParams {
			line := pref + param
			matches := matcher.FindStringSubmatch(line)
			assert.Empty(t, matches)
		}
	}
}

func verifySwaggerMultiArgSwaggerTag(t *testing.T, matcher *regexp.Regexp, prefixes, validParams, invalidParams []string) {
	var actualParams []string
	for i := 0; i < len(validParams); i++ {
		var vp []string
		for j := 0; j < (i + 1); j++ {
			vp = append(vp, validParams[j])
		}
		actualParams = append(actualParams, strings.Join(vp, " "))
	}
	for _, pref := range prefixes {
		for _, param := range actualParams {
			line := pref + param
			matches := matcher.FindStringSubmatch(line)
			// fmt.Printf("matching %q with %q, matches (%d): %v\n", line, matcher, len(matches), matches)
			assert.Len(t, matches, 2)
			assert.Equal(t, strings.TrimSpace(param), matches[1])
		}
	}

	for _, pref := range prefixes {
		for _, param := range invalidParams {
			line := pref + param
			matches := matcher.FindStringSubmatch(line)
			assert.Empty(t, matches)
		}
	}
}
