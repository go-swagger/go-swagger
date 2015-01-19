package validate

import (
	"net/url"
	"regexp"
	"strings"
)

const (
	// HostnamePattern http://json-schema.org/latest/json-schema-validation.html#anchor114
	//  A string instance is valid against this attribute if it is a valid
	//  representation for an Internet host name, as defined by RFC 1034, section 3.1 [RFC1034].
	//  http://tools.ietf.org/html/rfc1034#section-3.5
	//  <digit> ::= any one of the ten digits 0 through 9
	//  var digit = /[0-9]/;
	//  <letter> ::= any one of the 52 alphabetic characters A through Z in upper case and a through z in lower case
	//  var letter = /[a-zA-Z]/;
	//  <let-dig> ::= <letter> | <digit>
	//  var letDig = /[0-9a-zA-Z]/;
	//  <let-dig-hyp> ::= <let-dig> | "-"
	//  var letDigHyp = /[-0-9a-zA-Z]/;
	//  <ldh-str> ::= <let-dig-hyp> | <let-dig-hyp> <ldh-str>
	//  var ldhStr = /[-0-9a-zA-Z]+/;
	//  <label> ::= <letter> [ [ <ldh-str> ] <let-dig> ]
	//  var label = /[a-zA-Z](([-0-9a-zA-Z]+)?[0-9a-zA-Z])?/;
	//  <subdomain> ::= <label> | <subdomain> "." <label>
	//  var subdomain = /^[a-zA-Z](([-0-9a-zA-Z]+)?[0-9a-zA-Z])?(\.[a-zA-Z](([-0-9a-zA-Z]+)?[0-9a-zA-Z])?)*$/;
	//  <domain> ::= <subdomain> | " "
	HostnamePattern = `^[a-zA-Z](([-0-9a-zA-Z]+)?[0-9a-zA-Z])?(\.[a-zA-Z](([-0-9a-zA-Z]+)?[0-9a-zA-Z])?)*$`

	// DatePattern pattern to match for the date format from http://tools.ietf.org/html/rfc3339#section-5.6
	DatePattern = `^([0-9]{4})-([0-9]{2})-([0-9]{2})`

	// DateTimePattern pattern to match for the date-time format from http://tools.ietf.org/html/rfc3339#section-5.6
	DateTimePattern = `^([0-9]{2}):([0-9]{2}):([0-9]{2})(.[0-9]+)?(z|([+-][0-9]{2}:[0-9]{2}))$`
)

var (
	rxHostname = regexp.MustCompile(HostnamePattern)
	rxDate     = regexp.MustCompile(DatePattern)
	rxDateTime = regexp.MustCompile(DateTimePattern)
)

func IsURI(str string) bool {
	_, err := url.ParseRequestURI(str)
	return err == nil
}

// IsHostname returns true when the string is a valid hostname
func IsHostname(str string) bool {
	if !rxHostname.MatchString(str) {
		return false
	}

	// the sum of all label octets and label lengths is limited to 255.
	if len(str) > 255 {
		return false
	}

	// Each node has a label, which is zero to 63 octets in length
	parts := strings.Split(str, ".")
	valid := true
	for _, p := range parts {
		if len(p) > 63 {
			valid = false
		}
	}
	return valid
}

// IsDate returns true when the string is a valid date
func IsDate(str string) bool {
	matches := rxDate.FindAllStringSubmatch(str, -1)
	if len(matches) == 0 || len(matches[0]) == 0 {
		return false
	}
	m := matches[0]
	return !(m[2] < "01" || m[2] > "12" || m[3] < "01" || m[3] > "31")
}

// IsDateTime returns true when the string is a valid date-time
func IsDateTime(str string) bool {
	s := strings.Split(strings.ToLower(str), "t")
	if !IsDate(s[0]) {
		return false
	}

	matches := rxDateTime.FindAllStringSubmatch(s[1], -1)
	if len(matches) == 0 || len(matches[0]) == 0 {
		return false
	}
	m := matches[0]
	res := m[1] <= "23" && m[2] <= "59" && m[3] <= "59"
	return res
}
