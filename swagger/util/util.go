package util

import (
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

// Prepares strings by splitting by caps, spaces, dashes, and underscore
func split(str string) (words []string) {
	repl := strings.NewReplacer("-", " ", "_", " ")

	rex1 := regexp.MustCompile("([A-Z])")
	rex2 := regexp.MustCompile("(\\w+)")

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

// ToGoName translates a swagger name which can be underscored or camel cased to a name that golint likes
func ToGoName(name string) string {
	var out []string
	for _, w := range split(name) {
		uw := upper(w)
		if !commonInitialisms[uw] {
			uw = upper(w[:1]) + w[1:]
		}
		out = append(out, uw)
	}
	return strings.Join(out, "")
}
