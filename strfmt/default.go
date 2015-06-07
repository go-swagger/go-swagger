package strfmt

import (
	"encoding/base64"
	"net/url"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
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
)

var (
	rxHostname = regexp.MustCompile(HostnamePattern)
)

// IsStrictURI returns true when the string is an absolute URI
func IsStrictURI(str string) bool {
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

func init() {
	u := URI("")
	Default.Add("uri", &u, IsStrictURI)

	eml := Email("")
	Default.Add("email", &eml, govalidator.IsEmail)

	hn := Hostname("")
	Default.Add("hostname", &hn, IsHostname)

	ip4 := IPv4("")
	Default.Add("ipv4", &ip4, govalidator.IsIPv4)

	ip6 := IPv6("")
	Default.Add("ipv6", &ip6, govalidator.IsIPv6)

	uid := UUID("")
	Default.Add("uuid", &uid, govalidator.IsUUID)

	uid3 := UUID3("")
	Default.Add("uuid3", &uid3, govalidator.IsUUIDv3)

	uid4 := UUID4("")
	Default.Add("uuid4", &uid4, govalidator.IsUUIDv4)

	uid5 := UUID5("")
	Default.Add("uuid5", &uid5, govalidator.IsUUIDv5)

	isbn := ISBN("")
	Default.Add("isbn", &isbn, func(str string) bool { return govalidator.IsISBN10(str) || govalidator.IsISBN13(str) })

	isbn10 := ISBN10("")
	Default.Add("isbn10", &isbn10, govalidator.IsISBN10)

	isbn13 := ISBN13("")
	Default.Add("isbn13", &isbn13, govalidator.IsISBN13)

	cc := CreditCard("")
	Default.Add("creditcard", &cc, govalidator.IsCreditCard)

	ssn := SSN("")
	Default.Add("ssn", &ssn, govalidator.IsSSN)

	hc := HexColor("")
	Default.Add("hexcolor", &hc, govalidator.IsHexcolor)

	rc := RGBColor("")
	Default.Add("rgbcolor", &rc, govalidator.IsRGBcolor)

	b64 := Base64([]byte(nil))
	Default.Add("byte", &b64, govalidator.IsBase64)
}

var formatCheckers = map[string]Validator{
	"byte": govalidator.IsBase64,
}

// Base64 represents a base64 encoded string
//
// swagger:strfmt byte
type Base64 []byte

// MarshalText turns this instance into text
func (b Base64) MarshalText() ([]byte, error) {
	enc := base64.URLEncoding
	src := []byte(b)
	buf := make([]byte, enc.EncodedLen(len(src)))
	enc.Encode(buf, src)
	return buf, nil
}

// UnmarshalText hydrates this instance from text
func (b *Base64) UnmarshalText(data []byte) error { // validation is performed later on
	enc := base64.URLEncoding
	dbuf := make([]byte, enc.DecodedLen(len(data)))

	n, err := enc.Decode(dbuf, data)
	if err != nil {
		return err
	}

	*b = dbuf[:n]
	return nil
}

// URI represents the uri string format as specified by the json schema spec
//
// swagger:strfmt uri
type URI string

// MarshalText turns this instance into text
func (u URI) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

// UnmarshalText hydrates this instance from text
func (u *URI) UnmarshalText(data []byte) error { // validation is performed later on
	*u = URI(string(data))
	return nil
}

// Email represents the email string format as specified by the json schema spec
//
// swagger:strfmt email
type Email string

// MarshalText turns this instance into text
func (e Email) MarshalText() ([]byte, error) {
	return []byte(string(e)), nil
}

// UnmarshalText hydrates this instance from text
func (e *Email) UnmarshalText(data []byte) error { // validation is performed later on
	*e = Email(string(data))
	return nil
}

// Hostname represents the hostname string format as specified by the json schema spec
//
// swagger:strfmt hostname
type Hostname string

// MarshalText turns this instance into text
func (h Hostname) MarshalText() ([]byte, error) {
	return []byte(string(h)), nil
}

// UnmarshalText hydrates this instance from text
func (h *Hostname) UnmarshalText(data []byte) error { // validation is performed later on
	*h = Hostname(string(data))
	return nil
}

// IPv4 represents an IP v4 address
//
// swagger:strfmt ipv4
type IPv4 string

// MarshalText turns this instance into text
func (u IPv4) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

// UnmarshalText hydrates this instance from text
func (u *IPv4) UnmarshalText(data []byte) error { // validation is performed later on
	*u = IPv4(string(data))
	return nil
}

// IPv6 represents an IP v6 address
//
// swagger:strfmt ipv6
type IPv6 string

// MarshalText turns this instance into text
func (u IPv6) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

// UnmarshalText hydrates this instance from text
func (u *IPv6) UnmarshalText(data []byte) error { // validation is performed later on
	*u = IPv6(string(data))
	return nil
}

// UUID represents a uuid string format
//
// swagger:strfmt uuid
type UUID string

// MarshalText turns this instance into text
func (u UUID) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

// UnmarshalText hydrates this instance from text
func (u *UUID) UnmarshalText(data []byte) error { // validation is performed later on
	*u = UUID(string(data))
	return nil
}

// UUID3 represents a uuid3 string format
//
// swagger:strfmt uuid3
type UUID3 string

// MarshalText turns this instance into text
func (u UUID3) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

// UnmarshalText hydrates this instance from text
func (u *UUID3) UnmarshalText(data []byte) error { // validation is performed later on
	*u = UUID3(string(data))
	return nil
}

// UUID4 represents a uuid4 string format
//
// swagger:strfmt uuid4
type UUID4 string

// MarshalText turns this instance into text
func (u UUID4) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

// UnmarshalText hydrates this instance from text
func (u *UUID4) UnmarshalText(data []byte) error { // validation is performed later on
	*u = UUID4(string(data))
	return nil
}

// UUID5 represents a uuid5 string format
//
// swagger:strfmt uuid5
type UUID5 string

// MarshalText turns this instance into text
func (u UUID5) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

// UnmarshalText hydrates this instance from text
func (u *UUID5) UnmarshalText(data []byte) error { // validation is performed later on
	*u = UUID5(string(data))
	return nil
}

// ISBN represents an isbn string format
//
// swagger:strfmt isbn
type ISBN string

// MarshalText turns this instance into text
func (u ISBN) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

// UnmarshalText hydrates this instance from text
func (u *ISBN) UnmarshalText(data []byte) error { // validation is performed later on
	*u = ISBN(string(data))
	return nil
}

// ISBN10 represents an isbn 10 string format
//
// swagger:strfmt isbn10
type ISBN10 string

// MarshalText turns this instance into text
func (u ISBN10) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

// UnmarshalText hydrates this instance from text
func (u *ISBN10) UnmarshalText(data []byte) error { // validation is performed later on
	*u = ISBN10(string(data))
	return nil
}

// ISBN13 represents an isbn 13 string format
//
// swagger:strfmt isbn13
type ISBN13 string

// MarshalText turns this instance into text
func (u ISBN13) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

// UnmarshalText hydrates this instance from text
func (u *ISBN13) UnmarshalText(data []byte) error { // validation is performed later on
	*u = ISBN13(string(data))
	return nil
}

// CreditCard represents a credit card string format
//
// swagger:strfmt creditcard
type CreditCard string

// MarshalText turns this instance into text
func (u CreditCard) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

// UnmarshalText hydrates this instance from text
func (u *CreditCard) UnmarshalText(data []byte) error { // validation is performed later on
	*u = CreditCard(string(data))
	return nil
}

// SSN represents a social security string format
//
// swagger:strfmt ssn
type SSN string

// MarshalText turns this instance into text
func (u SSN) MarshalText() ([]byte, error) {
	return []byte(string(u)), nil
}

// UnmarshalText hydrates this instance from text
func (u *SSN) UnmarshalText(data []byte) error { // validation is performed later on
	*u = SSN(string(data))
	return nil
}

// HexColor represents a hex color string format
//
// swagger:strfmt hexcolor
type HexColor string

// MarshalText turns this instance into text
func (h HexColor) MarshalText() ([]byte, error) {
	return []byte(string(h)), nil
}

// UnmarshalText hydrates this instance from text
func (h *HexColor) UnmarshalText(data []byte) error { // validation is performed later on
	*h = HexColor(string(data))
	return nil
}

// RGBColor represents a RGB color string format
//
// swagger:strfmt rgbcolor
type RGBColor string

// MarshalText turns this instance into text
func (r RGBColor) MarshalText() ([]byte, error) {
	return []byte(string(r)), nil
}

// UnmarshalText hydrates this instance from text
func (r *RGBColor) UnmarshalText(data []byte) error { // validation is performed later on
	*r = RGBColor(string(data))
	return nil
}
