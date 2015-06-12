package swag

import (
	"math"
	"regexp"
	"strings"
)

// Taken from https://github.com/golang/lint/blob/1fab560e16097e5b69afb66eb93aab843ef77845/lint.go#L663-L698
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
}

// JoinByFormat joins a string array by a known format:
//		ssv: space separated value
//		tsv: tab separated value
//		pipes: pipe (|) separated value
//		csv: comma separated value (default)
func JoinByFormat(data []string, format string) []string {
	if len(data) == 0 {
		return data
	}
	var sep string
	switch format {
	case "ssv":
		sep = " "
	case "tsv":
		sep = "\t"
	case "pipes":
		sep = "|"
	case "multi":
		return data
	default:
		sep = ","
	}
	return []string{strings.Join(data, sep)}
}

// SplitByFormat splits a string by a known format:
//		ssv: space separated value
//		tsv: tab separated value
//		pipes: pipe (|) separated value
//		csv: comma separated value (default)
func SplitByFormat(data, format string) []string {
	if data == "" {
		return nil
	}
	var sep string
	switch format {
	case "ssv":
		sep = " "
	case "tsv":
		sep = "\t"
	case "pipes":
		sep = "|"
	case "multi":
		return nil
	default:
		sep = ","
	}
	var result []string
	for _, s := range strings.Split(data, sep) {
		if ts := strings.TrimSpace(s); ts != "" {
			result = append(result, ts)
		}
	}
	return result
}

// Prepares strings by splitting by caps, spaces, dashes, and underscore
func split(str string) (words []string) {
	repl := strings.NewReplacer("-", " ", "_", " ")

	rex1 := regexp.MustCompile(`(\p{Lu})`)
	rex2 := regexp.MustCompile(`(\pL|\pM|\pN|\p{Pc})+`)

	str = trim(str)

	// Convert dash and underscore to spaces
	str = repl.Replace(str)

	// Split when uppercase is found (needed for Snake)
	str = rex1.ReplaceAllString(str, " $1")

	// Get the final list of words
	words = rex2.FindAllString(str, -1)

	return
}

// Removes leading whitespaces
func trim(str string) string {
	return strings.Trim(str, " ")
}

// Shortcut to strings.ToUpper()
func upper(str string) string {
	return strings.ToUpper(trim(str))
}

// Shortcut to strings.ToLower()
func lower(str string) string {
	return strings.ToLower(trim(str))
}

// ToFileName lowercases and underscores a go type name
func ToFileName(name string) string {
	var out []string
	for _, w := range split(name) {
		out = append(out, lower(w))
	}
	return strings.Join(out, "_")
}

// ToCommandName lowercases and underscores a go type name
func ToCommandName(name string) string {
	var out []string
	for _, w := range split(name) {
		out = append(out, lower(w))
	}
	return strings.Join(out, "-")
}

// ToHumanNameLower represents a code name as a human series of words
func ToHumanNameLower(name string) string {
	var out []string
	for _, w := range split(name) {
		out = append(out, lower(w))
	}
	return strings.Join(out, " ")
}

// ToJSONName camelcases a name which can be underscored or pascal cased
func ToJSONName(name string) string {
	var out []string
	for i, w := range split(name) {
		if i == 0 {
			out = append(out, lower(w))
			continue
		}
		out = append(out, upper(w[:1])+lower(w[1:]))
	}
	return strings.Join(out, "")
}

// ToGoName translates a swagger name which can be underscored or camel cased to a name that golint likes
func ToGoName(name string) string {
	var out []string
	for _, w := range split(name) {
		uw := upper(w)
		mod := int(math.Min(float64(len(uw)), 2))
		if !commonInitialisms[uw] && !commonInitialisms[uw[:len(uw)-mod]] {
			uw = upper(w[:1]) + lower(w[1:])
		}
		out = append(out, uw)
	}
	return strings.Join(out, "")
}

// ContainsStringsCI searches a slice of strings for a case-insensitive match
func ContainsStringsCI(coll []string, item string) bool {
	for _, a := range coll {
		if strings.EqualFold(a, item) {
			return true
		}
	}
	return false
}
