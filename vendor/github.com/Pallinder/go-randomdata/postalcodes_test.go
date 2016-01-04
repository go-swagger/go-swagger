package randomdata

import (
	"testing"
)

var postalcodeTests = []struct {
	Country string
	Size    int
}{
	{"PE", 6},
	{"FO", 6},
	{"AF", 4},
	{"DZ", 5},
	{"BY", 6},
	{"CL", 7},
	{"SZ", 4},
	{"BM", 4},
	{"AD", 5},
	{"BN", 6},
	{"BB", 7},
	{"MT", 7},
	{"JM", 7},
	{"AR", 8},
	{"CA", 6},
	{"FK", 7},
	{"GG", 6},
	{"NL", 6},
	{"BR", 9},
	{"KY", 8},
	{"JP", 8},
	{"LV", 7},
	{"LT", 8},
	{"MV", 5},
	{"NI", 9},
	{"PL", 6},
	{"PT", 8},
	{"KR", 7},
	{"TW", 5},
}

func TestPostalCode(t *testing.T) {
	for _, pt := range postalcodeTests {
		code := PostalCode(pt.Country)

		if len(code) == pt.Size {
			continue
		}

		t.Fatalf("Invalid length for country %q: Expected %d, have %d.",
			pt.Country, pt.Size, len(code))
	}
}
