package strfmt

import (
	"testing"

	"github.com/nu7hatch/gouuid"
	"github.com/stretchr/testify/assert"
)

func TestFormats(t *testing.T) {
	uri := URI("http://somewhere.com")
	str := string("http://somewhereelse.com")
	b := []byte(str)
	err := uri.UnmarshalText(b)
	assert.NoError(t, err)
	assert.Equal(t, URI("http://somewhereelse.com"), string(b))

	b, err = uri.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("http://somewhereelse.com"), b)

	email := Email("somebody@somewhere.com")
	str = string("somebodyelse@somewhere.com")
	b = []byte(str)
	err = email.UnmarshalText(b)
	assert.NoError(t, err)
	assert.Equal(t, Email("somebodyelse@somewhere.com"), string(b))

	b, err = email.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("somebodyelse@somewhere.com"), b)

	hostname := Hostname("somewhere.com")
	str = string("somewhere.com")
	b = []byte(str)
	err = hostname.UnmarshalText(b)
	assert.NoError(t, err)
	assert.Equal(t, Hostname("somewhere.com"), string(b))

	b, err = hostname.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("somewhere.com"), b)

	ipv4 := IPv4("192.168.254.1")
	str = string("192.168.254.2")
	b = []byte(str)
	err = ipv4.UnmarshalText(b)
	assert.NoError(t, err)
	assert.Equal(t, IPv4("192.168.254.2"), string(b))

	b, err = ipv4.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("192.168.254.2"), b)

	ipv6 := IPv6("::1")
	str = string("::2")
	b = []byte(str)
	err = ipv6.UnmarshalText(b)
	assert.NoError(t, err)
	assert.Equal(t, IPv6("::2"), string(b))

	b, err = ipv6.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("::2"), b)

	first3, _ := uuid.NewV3(uuid.NamespaceURL, []byte("somewhere.com"))
	other3, _ := uuid.NewV3(uuid.NamespaceURL, []byte("somewhereelse.com"))
	uuid3 := UUID3(first3.String())
	str = string(other3.String())
	b = []byte(str)
	err = uuid3.UnmarshalText(b)
	assert.NoError(t, err)
	assert.Equal(t, UUID3(other3.String()), string(b))

	b, err = uuid3.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte(other3.String()), b)

	first4, _ := uuid.NewV4()
	other4, _ := uuid.NewV4()
	uuid4 := UUID4(first4.String())
	str = string(other4.String())
	b = []byte(str)
	err = uuid4.UnmarshalText(b)
	assert.NoError(t, err)
	assert.Equal(t, UUID4(other4.String()), string(b))

	b, err = uuid4.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte(other4.String()), b)

	first5, _ := uuid.NewV5(uuid.NamespaceURL, []byte("somewhere.com"))
	other5, _ := uuid.NewV5(uuid.NamespaceURL, []byte("somewhereelse.com"))
	uuid5 := UUID5(first5.String())
	str = string(other5.String())
	b = []byte(str)
	err = uuid5.UnmarshalText(b)
	assert.NoError(t, err)
	assert.Equal(t, UUID5(other5.String()), string(b))

	b, err = uuid5.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte(other5.String()), b)

	uuid := UUID(first5.String())
	str = string(other5.String())
	b = []byte(str)
	err = uuid.UnmarshalText(b)
	assert.NoError(t, err)
	assert.Equal(t, UUID(other5.String()), string(b))

	b, err = uuid.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte(other5.String()), b)

	isbn := ISBN("0836217462")
	str = string("0836217463")
	b = []byte(str)
	err = isbn.UnmarshalText(b)
	assert.NoError(t, err)
	assert.Equal(t, ISBN("0836217463"), string(b))

	b, err = isbn.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("0836217463"), b)

	isbn10 := ISBN10("0836217462")
	str = string("0836217463")
	b = []byte(str)
	err = isbn10.UnmarshalText(b)
	assert.NoError(t, err)
	assert.Equal(t, ISBN10("0836217463"), string(b))

	b, err = isbn10.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("0836217463"), b)

	isbn13 := ISBN13("0836217462384")
	str = string("0836217462385")
	b = []byte(str)
	err = isbn13.UnmarshalText(b)
	assert.NoError(t, err)
	assert.Equal(t, ISBN13("0836217462385"), string(b))

	b, err = isbn13.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("0836217462385"), b)

	hexColor := HexColor("#FFFFFF")
	str = string("#000000")
	b = []byte(str)
	err = hexColor.UnmarshalText(b)
	assert.NoError(t, err)
	assert.Equal(t, HexColor("#000000"), string(b))

	b, err = hexColor.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("#000000"), b)

	rgbColor := RGBColor("rgb(255,255,255)")
	str = string("rgb(0,0,0)")
	b = []byte(str)
	err = rgbColor.UnmarshalText(b)
	assert.NoError(t, err)
	assert.Equal(t, RGBColor("rgb(0,0,0)"), string(b))

	b, err = rgbColor.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("rgb(0,0,0)"), b)

	ssn := SSN("111-11-1111")
	str = string("999 99 9999")
	b = []byte(str)
	err = ssn.UnmarshalText(b)
	assert.NoError(t, err)
	assert.Equal(t, SSN("999 99 9999"), string(b))

	b, err = ssn.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("999 99 9999"), b)

	creditCard := CreditCard("1111-1111-1111-1111")
	str = string("9999-9999-9999-9999")
	b = []byte(str)
	err = creditCard.UnmarshalText(b)
	assert.NoError(t, err)
	assert.Equal(t, CreditCard("9999-9999-9999-9999"), string(b))

	b, err = creditCard.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("9999-9999-9999-9999"), b)
}
