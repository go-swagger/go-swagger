// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package strfmt

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFormatURI(t *testing.T) {
	uri := URI("http://somewhere.com")
	str := "http://somewhereelse.com"
	testStringFormat(t, &uri, "uri", str, []string{}, []string{"somewhere.com"})
}

func TestFormatEmail(t *testing.T) {
	email := Email("somebody@somewhere.com")
	str := string("somebodyelse@somewhere.com")
	validEmails := []string{
		"blah@gmail.com",
		"test@d.verylongtoplevel",
		"email+tag@gmail.com",
		`" "@example.com`,
		`"Abc\@def"@example.com`,
		`"Fred Bloggs"@example.com`,
		`"Joe\\Blow"@example.com`,
		`"Abc@def"@example.com`,
		"customer/department=shipping@example.com",
		"$A12345@example.com",
		"!def!xyz%abc@example.com",
		"_somename@example.com",
		"!#$%&'*+-/=?^_`{}|~@example.com",
		"Miles.O'Brian@example.com",
		"postmaster@☁→❄→☃→☀→☺→☂→☹→✝.ws",
		"root@localhost",
		"john@com",
		"api@piston.ninja",
	}

	testStringFormat(t, &email, "email", str, validEmails, []string{"somebody@somewhere@com"})
}

func TestFormatHostname(t *testing.T) {
	hostname := Hostname("somewhere.com")
	str := string("somewhere.com")
	veryLongStr := strings.Repeat("a", 256)
	longStr := strings.Repeat("a", 64)
	longAddrSegment := strings.Join([]string{"x", "y", longStr}, ".")
	invalidHostnames := []string{
		"somewhere.com!",
		"user@email.domain",
		"1.1.1.1",
		veryLongStr,
		longAddrSegment,
	}
	validHostnames := []string{
		"somewhere.com",
		"888.com",
		"a.com",
		"a.b.com",
		"a.b.c.com",
		"a.b.c.d.com",
		"a.b.c.d.e.com",
		"1.com",
		"1.2.com",
		"1.2.3.com",
		"1.2.3.4.com",
		"99.domain.com",
		"99.99.domain.com",
		"xn-80ak6aa92e.co",
		"xn-80ak6aa92e.com",
		"xn--ls8h.la",
		"☁→❄→☃→☀→☺→☂→☹→✝.ws",
		"www.example.onion",
		"www.example.ôlà",
		"ôlà.ôlà",
		"ôlà.ôlà.ôlà",
	}

	testStringFormat(t, &hostname, "hostname", str, []string{}, invalidHostnames)
	testStringFormat(t, &hostname, "hostname", str, validHostnames, []string{})
}

func TestFormatIPv4(t *testing.T) {
	ipv4 := IPv4("192.168.254.1")
	str := string("192.168.254.2")
	testStringFormat(t, &ipv4, "ipv4", str, []string{}, []string{"198.168.254.2.2"})
}

func TestFormatIPv6(t *testing.T) {
	ipv6 := IPv6("::1")
	str := string("::2")
	// TODO: test ipv6 zones
	testStringFormat(t, &ipv6, "ipv6", str, []string{}, []string{"127.0.0.1"})
}

func TestFormatMAC(t *testing.T) {
	mac := MAC("01:02:03:04:05:06")
	str := string("06:05:04:03:02:01")
	testStringFormat(t, &mac, "mac", str, []string{}, []string{"01:02:03:04:05"})
}

func TestFormatUUID3(t *testing.T) {
	first3 := uuid.NewMD5(uuid.NameSpace_URL, []byte("somewhere.com"))
	other3 := uuid.NewMD5(uuid.NameSpace_URL, []byte("somewhereelse.com"))
	uuid3 := UUID3(first3.String())
	str := other3.String()
	testStringFormat(t, &uuid3, "uuid3", str, []string{}, []string{"not-a-uuid"})

	// special case for zero UUID
	var uuidZero UUID3
	err := uuidZero.UnmarshalJSON([]byte(jsonNull))
	assert.NoError(t, err)
	assert.EqualValues(t, UUID3(""), uuidZero)
}

func TestFormatUUID4(t *testing.T) {
	first4 := uuid.NewRandom()
	other4 := uuid.NewRandom()
	uuid4 := UUID4(first4.String())
	str := other4.String()
	testStringFormat(t, &uuid4, "uuid4", str, []string{}, []string{"not-a-uuid"})

	// special case for zero UUID
	var uuidZero UUID4
	err := uuidZero.UnmarshalJSON([]byte(jsonNull))
	assert.NoError(t, err)
	assert.EqualValues(t, UUID4(""), uuidZero)
}

func TestFormatUUID5(t *testing.T) {
	first5 := uuid.NewSHA1(uuid.NameSpace_URL, []byte("somewhere.com"))
	other5 := uuid.NewSHA1(uuid.NameSpace_URL, []byte("somewhereelse.com"))
	uuid5 := UUID5(first5.String())
	str := other5.String()
	testStringFormat(t, &uuid5, "uuid5", str, []string{}, []string{"not-a-uuid"})

	// special case for zero UUID
	var uuidZero UUID5
	err := uuidZero.UnmarshalJSON([]byte(jsonNull))
	assert.NoError(t, err)
	assert.EqualValues(t, UUID5(""), uuidZero)
}

func TestFormatUUID(t *testing.T) {
	first5 := uuid.NewSHA1(uuid.NameSpace_URL, []byte("somewhere.com"))
	other5 := uuid.NewSHA1(uuid.NameSpace_URL, []byte("somewhereelse.com"))
	uuid := UUID(first5.String())
	str := other5.String()
	testStringFormat(t, &uuid, "uuid", str, []string{}, []string{"not-a-uuid"})

	// special case for zero UUID
	var uuidZero UUID
	err := uuidZero.UnmarshalJSON([]byte(jsonNull))
	assert.NoError(t, err)
	assert.EqualValues(t, UUID(""), uuidZero)
}

func TestFormatISBN(t *testing.T) {
	isbn := ISBN("0321751043")
	str := string("0321751043")
	testStringFormat(t, &isbn, "isbn", str, []string{}, []string{"836217463"}) // bad checksum
}

func TestFormatISBN10(t *testing.T) {
	isbn10 := ISBN10("0321751043")
	str := string("0321751043")
	testStringFormat(t, &isbn10, "isbn10", str, []string{}, []string{"836217463"}) // bad checksum
}

func TestFormatISBN13(t *testing.T) {
	isbn13 := ISBN13("978-0321751041")
	str := string("978-0321751041")
	testStringFormat(t, &isbn13, "isbn13", str, []string{}, []string{"978-0321751042"}) // bad checksum
}

func TestFormatHexColor(t *testing.T) {
	hexColor := HexColor("#FFFFFF")
	str := string("#000000")
	testStringFormat(t, &hexColor, "hexcolor", str, []string{}, []string{"#fffffffz"})
}

func TestFormatRGBColor(t *testing.T) {
	rgbColor := RGBColor("rgb(255,255,255)")
	str := string("rgb(0,0,0)")
	testStringFormat(t, &rgbColor, "rgbcolor", str, []string{}, []string{"rgb(300,0,0)"})
}

func TestFormatSSN(t *testing.T) {
	ssn := SSN("111-11-1111")
	str := string("999 99 9999")
	testStringFormat(t, &ssn, "ssn", str, []string{}, []string{"999 99 999"})
}

func TestFormatCreditCard(t *testing.T) {
	creditCard := CreditCard("4111-1111-1111-1111")
	str := string("4012-8888-8888-1881")
	testStringFormat(t, &creditCard, "creditcard", str, []string{}, []string{"9999-9999-9999-999"})
}

func TestFormatPassword(t *testing.T) {
	password := Password("super secret stuff here")
	testStringFormat(t, &password, "password", "super secret!!!", []string{"even more secret"}, []string{})
}

func TestFormatBase64(t *testing.T) {
	b64 := Base64("ZWxpemFiZXRocG9zZXk=")
	str := string("ZWxpemFiZXRocG9zZXk=")
	b := []byte(str)
	bj := []byte("\"" + str + "\"")

	err := b64.UnmarshalText(b)
	assert.NoError(t, err)
	assert.EqualValues(t, Base64("ZWxpemFiZXRocG9zZXk="), string(b))

	b, err = b64.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("ZWxpemFiZXRocG9zZXk="), b)

	err = b64.UnmarshalJSON(bj)
	assert.NoError(t, err)
	assert.EqualValues(t, Base64(str), string(b))

	b, err = b64.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, bj, b)

	bsonData, err := bson.Marshal(&b64)
	assert.NoError(t, err)

	var b64Copy Base64
	err = bson.Unmarshal(bsonData, &b64Copy)
	assert.NoError(t, err)
	assert.Equal(t, b64, b64Copy)

	testValid(t, "byte", str)
	testInvalid(t, "byte", "ZWxpemFiZXRocG9zZXk") // missing pad char

	// Valuer interface
	b64 = Base64("ZWxpemFiZXRocG9zZXk=")
	sqlvalue, err := b64.Value()
	assert.NoError(t, err)
	sqlvalueAsString, ok := sqlvalue.(string)
	if assert.Truef(t, ok, "[%s]Value: expected driver value to be a string", "byte") {
		assert.EqualValuesf(t, str, sqlvalueAsString, "[%s]Value: expected %v and %v to be equal", "byte", sqlvalue, str)
	}
	// Scanner interface
	b64 = Base64("")
	err = b64.Scan(str)
	assert.NoError(t, err)
	b64AsStr := b64.String()
	assert.EqualValues(t, str, b64AsStr)

	err = b64.Scan([]byte(str))
	assert.NoError(t, err)
	b64AsStr = b64.String()
	assert.EqualValues(t, str, b64AsStr)

	err = b64.Scan(123)
	assert.Error(t, err)
}

type testableFormat interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	json.Marshaler
	json.Unmarshaler
	bson.Getter
	bson.Setter
	fmt.Stringer
	sql.Scanner
	driver.Valuer
}

func testStringFormat(t *testing.T, what testableFormat, format, with string, validSamples, invalidSamples []string) {
	// text encoding interface
	b := []byte(with)
	err := what.UnmarshalText(b)
	assert.NoError(t, err)

	val := reflect.Indirect(reflect.ValueOf(what))
	strVal := val.String()
	assert.Equalf(t, with, strVal, "[%s]UnmarshalText: expected %v and %v to be value equal", format, strVal, with)

	b, err = what.MarshalText()
	assert.NoError(t, err)
	assert.Equalf(t, []byte(with), b, "[%s]MarshalText: expected %v and %v to be value equal as []byte", format, string(b), with)

	// Stringer
	strVal = what.String()
	assert.Equalf(t, []byte(with), b, "[%s]String: expected %v and %v to be equal", strVal, with)

	// JSON encoding interface
	bj := []byte("\"" + with + "\"")
	err = what.UnmarshalJSON(bj)
	assert.NoError(t, err)
	val = reflect.Indirect(reflect.ValueOf(what))
	strVal = val.String()
	assert.EqualValuesf(t, with, strVal, "[%s]UnmarshalJSON: expected %v and %v to be value equal", format, strVal, with)

	b, err = what.MarshalJSON()
	assert.NoError(t, err)
	assert.Equalf(t, bj, b, "[%s]MarshalJSON: expected %v and %v to be value equal as []byte", format, string(b), with)

	// bson encoding interface
	bsonData, err := bson.Marshal(&what)
	assert.NoError(t, err)

	resetValue(t, format, what)

	err = bson.Unmarshal(bsonData, what)
	assert.NoError(t, err)
	val = reflect.Indirect(reflect.ValueOf(what))
	strVal = val.String()
	assert.EqualValuesf(t, with, strVal, "[%s]bson.Unmarshal: expected %v and %v to be equal (reset value) ", format, what, with)

	// Scanner interface
	resetValue(t, format, what)
	err = what.Scan(with)
	assert.NoError(t, err)
	val = reflect.Indirect(reflect.ValueOf(what))
	strVal = val.String()
	assert.EqualValuesf(t, with, strVal, "[%s]Scan: expected %v and %v to be value equal", format, strVal, with)

	err = what.Scan([]byte(with))
	assert.NoError(t, err)
	val = reflect.Indirect(reflect.ValueOf(what))
	strVal = val.String()
	assert.EqualValuesf(t, with, strVal, "[%s]Scan: expected %v and %v to be value equal", format, strVal, with)

	err = what.Scan(123)
	assert.Error(t, err)

	// Valuer interface
	sqlvalue, err := what.Value()
	assert.NoError(t, err)
	sqlvalueAsString, ok := sqlvalue.(string)
	if assert.Truef(t, ok, "[%s]Value: expected driver value to be a string", format) {
		assert.EqualValuesf(t, with, sqlvalueAsString, "[%s]Value: expected %v and %v to be equal", format, sqlvalue, with)
	}

	// validation with Registry
	for _, valid := range append(validSamples, with) {
		testValid(t, format, valid)
	}

	for _, invalid := range invalidSamples {
		testInvalid(t, format, invalid)
	}
}

func resetValue(t *testing.T, format string, what encoding.TextUnmarshaler) {
	err := what.UnmarshalText([]byte("reset value"))
	assert.NoError(t, err)
	val := reflect.Indirect(reflect.ValueOf(what))
	strVal := val.String()
	assert.Equalf(t, "reset value", strVal, "[%s]UnmarshalText: expected %v and %v to be equal (reset value) ", format, strVal, "reset value")
}

func testValid(t *testing.T, name, value string) {
	ok := Default.Validates(name, value)
	if !ok {
		t.Errorf("expected %q of type %s to be valid", value, name)
	}
}

func testInvalid(t *testing.T, name, value string) {
	ok := Default.Validates(name, value)
	if ok {
		t.Errorf("expected %q of type %s to be invalid", value, name)
	}
}
