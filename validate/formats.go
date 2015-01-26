package validate

import (
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
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

// IsStrictURI returns true when the string is an absolute URI
func IsStrictURI(str string) bool {
	_, err := url.ParseRequestURI(str)
	return err == nil
}

// IsURI returns true when the string resembles an URI
func IsURI(str string) bool {
	// this makes little sense really but it needs
	// https://github.com/swagger-api/swagger-spec/issues/249
	// to be resolved before this can be changed
	return IsStrictURI(str) || IsHostname(str)
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

// FormatValidator validates if a string matches a format
type FormatValidator func(string) bool

var formatCheckers = map[string]FormatValidator{
	"datetime":   IsDateTime,
	"date":       IsDate,
	"byte":       govalidator.IsBase64,
	"uri":        IsStrictURI,
	"email":      govalidator.IsEmail,
	"hostname":   IsHostname,
	"ipv4":       govalidator.IsIPv4,
	"ipv6":       govalidator.IsIPv6,
	"uuid":       govalidator.IsUUID,
	"uuid3":      govalidator.IsUUIDv3,
	"uuid4":      govalidator.IsUUIDv4,
	"uuid5":      govalidator.IsUUIDv5,
	"isbn":       func(str string) bool { return govalidator.IsISBN10(str) || govalidator.IsISBN13(str) },
	"isbn10":     govalidator.IsISBN10,
	"isbn13":     govalidator.IsISBN13,
	"creditcard": govalidator.IsCreditCard,
	"ssn":        govalidator.IsSSN,
	"hexcolor":   govalidator.IsHexcolor,
	"rgbcolor":   govalidator.IsRGBcolor,
}

type formatValidator struct {
	Default      interface{}
	Format       string
	Path         string
	In           string
	KnownFormats map[string]FormatValidator
}

func (f *formatValidator) SetPath(path string) {
	f.Path = path
}

func (f *formatValidator) Applies(source interface{}, kind reflect.Kind) bool {
	doit := func() bool {
		if source == nil {
			return false
		}
		switch source.(type) {
		case *spec.Items:
			it := source.(*spec.Items)
			_, known := f.KnownFormats[strings.Replace(it.Format, "-", "", -1)]
			return kind == reflect.String && known
		case *spec.Parameter:
			par := source.(*spec.Parameter)
			_, known := f.KnownFormats[strings.Replace(par.Format, "-", "", -1)]
			return kind == reflect.String && known
		case *spec.Schema:
			sch := source.(*spec.Schema)
			_, known := f.KnownFormats[strings.Replace(sch.Format, "-", "", -1)]
			return kind == reflect.String && known
		}
		return false
	}
	r := doit()
	// fmt.Printf("schema props validator for %q applies %t for %T (kind: %v)\n", f.Path, r, source, kind)
	return r
}

func (f *formatValidator) Validate(val interface{}) *Result {
	result := new(Result)

	var valid bool
	if validate, ok := f.KnownFormats[strings.Replace(f.Format, "-", "", -1)]; ok {
		valid = validate(val.(string))
	}

	if !valid {
		result.AddErrors(errors.InvalidType(f.Path, f.In, f.Format, val))
	}
	result.Inc()
	return result
}
