// Package randomdata implements a bunch of simple ways to generate (pseudo) random data
package randomdata

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

// Supported formats obtained from:
// * http://www.geopostcodes.com/GeoPC_Postal_codes_formats

// PostalCode yields a random postal/zip code for the given 2-letter country code.
//
// These codes are not guaranteed to refer to actually locations.
// They merely follow the correct format as far as letters and digits goes.
// Where possible, the function enforces valid ranges of letters and digits.
func PostalCode(countrycode string) string {
	switch strings.ToUpper(countrycode) {
	case "LS", "MG", "IS", "OM", "PG":
		return Digits(3)

	case "AM", "GE", "NZ", "NE", "NO", "PY", "ZA", "MZ", "SJ", "LI", "AL",
		"BD", "CV", "GL":
		return Digits(4)

	case "DZ", "BA", "KH", "DO", "EG", "EE", "GP", "GT", "ID", "IL", "JO",
		"KW", "MQ", "MX", "LK", "SD", "TR", "UA", "US", "CR", "IQ", "KV", "MY",
		"MN", "ME", "PK", "SM", "MA", "UY", "EH", "ZM":
		return Digits(5)

	case "BY", "CN", "IN", "KZ", "KG", "NG", "RO", "RU", "SG", "TJ", "TM", "UZ", "VN":
		return Digits(6)

	case "CL":
		return Digits(7)

	case "IR":
		return Digits(10)

	case "FO":
		return "FO " + Digits(3)

	case "AF":
		return BoundedDigits(2, 10, 43) + BoundedDigits(2, 1, 99)

	case "AU", "AT", "BE", "BG", "CY", "DK", "ET", "GW", "HU", "LR", "MK", "PH",
		"CH", "TN", "VE":
		return BoundedDigits(4, 1000, 9999)

	case "SV":
		return "CP " + BoundedDigits(4, 1000, 9999)

	case "HT":
		return "HT" + Digits(4)

	case "LB":
		return Digits(4) + " " + Digits(4)

	case "LU":
		return BoundedDigits(4, 6600, 6999)

	case "MD":
		return "MD-" + BoundedDigits(4, 1000, 9999)

	case "HR":
		return "HR-" + Digits(5)

	case "CU":
		return "CP " + BoundedDigits(5, 10000, 99999)

	case "FI":
		// Last digit is usually 0 but can, in some cases, be 1 or 5.
		switch rand.Intn(2) {
		case 0:
			return Digits(4) + "0"
		case 1:
			return Digits(4) + "1"
		}

		return Digits(4) + "5"

	case "FR", "GF", "PF", "YT", "MC", "RE", "BL", "MF", "PM", "RS", "TH":
		return BoundedDigits(5, 10000, 99999)

	case "DE":
		return BoundedDigits(5, 1000, 99999)

	case "GR":
		return BoundedDigits(3, 100, 999) + " " + Digits(2)

	case "HN":
		return "CM" + Digits(4)

	case "IT", "VA":
		return BoundedDigits(5, 10, 99999)

	case "KE":
		return BoundedDigits(5, 100, 99999)

	case "LA":
		return BoundedDigits(5, 1000, 99999)

	case "MH":
		return BoundedDigits(5, 96960, 96970)

	case "FM":
		return "FM" + BoundedDigits(5, 96941, 96944)

	case "MM":
		return BoundedDigits(2, 1, 14) + Digits(3)

	case "NP":
		return BoundedDigits(5, 10700, 56311)

	case "NC":
		return "98" + Digits(3)

	case "PW":
		return "PW96940"

	case "PR":
		return "PR " + Digits(5)

	case "SA":
		return BoundedDigits(5, 10000, 99999) + "-" + BoundedDigits(4, 1000, 9999)

	case "ES":
		return BoundedDigits(2, 1, 52) + BoundedDigits(3, 100, 999)

	case "WF":
		return "986" + Digits(2)

	case "SZ":
		return Letters(1) + Digits(3)

	case "BM":
		return Letters(2) + Digits(2)

	case "AD":
		return Letters(2) + Digits(3)

	case "BN", "AZ", "VG", "PE":
		return Letters(2) + Digits(4)

	case "BB":
		return Letters(2) + Digits(5)

	case "EC":
		return Letters(2) + Digits(6)

	case "MT":
		return Letters(3) + Digits(4)

	case "JM":
		return "JM" + Letters(3) + Digits(2)

	case "AR":
		return Letters(1) + Digits(4) + Letters(3)

	case "CA":
		return Letters(1) + Digits(1) + Letters(1) + Digits(1) + Letters(1) + Digits(1)

	case "FK", "TC":
		return Letters(4) + Digits(1) + Letters(2)

	case "GG", "IM", "JE", "GB":
		return Letters(2) + Digits(2) + Letters(2)

	case "KY":
		return Letters(2) + Digits(1) + "-" + Digits(4)

	case "JP":
		return Digits(3) + "-" + Digits(4)

	case "LV", "SI":
		return Letters(2) + "-" + Digits(4)

	case "LT":
		return Letters(2) + "-" + Digits(5)

	case "SE", "TW":
		return Digits(5)

	case "MV":
		return Digits(2) + "-" + Digits(2)

	case "PL":
		return Digits(2) + "-" + Digits(3)

	case "NI":
		return Digits(3) + "-" + Digits(3) + "-" + Digits(1)

	case "KR":
		return Digits(3) + "-" + Digits(3)

	case "PT":
		return Digits(4) + "-" + Digits(3)

	case "NL":
		return Digits(4) + Letters(2)

	case "BR":
		return Digits(5) + "-" + Digits(3)
	}

	return ""
}

// Letters generates a string of N random leters (A-Z).
func Letters(letters int) string {
	list := make([]byte, letters)

	for i := range list {
		list[i] = byte(rand.Intn('Z'-'A') + 'A')
	}

	return string(list)
}

// Digits generates a string of N random digits, padded with zeros if necessary.
func Digits(digits int) string {
	max := int(math.Pow10(digits)) - 1
	num := rand.Intn(max)
	format := fmt.Sprintf("%%0%dd", digits)
	return fmt.Sprintf(format, num)
}

// BoundedDigits generates a string of N random digits, padded with zeros if necessary.
// The output is restricted to the given range.
func BoundedDigits(digits, low, high int) string {
	if low > high {
		low, high = high, low
	}

	max := (int(math.Pow10(digits)) - 1) & high
	num := rand.Intn(max-low) + low
	format := fmt.Sprintf("%%0%dd", digits)
	return fmt.Sprintf(format, num)
}
