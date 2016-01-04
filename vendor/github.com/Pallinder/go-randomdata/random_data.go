// Package randomdata implements a bunch of simple ways to generate (pseudo) random data
package randomdata

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	Male         int = 0
	Female       int = 1
	RandomGender int = 2
)

const (
	Small int = 0
	Large int = 1
)

const (
	FullCountry      = 0
	TwoCharCountry   = 1
	ThreeCharCountry = 2
)

type jsonContent struct {
	Adjectives          []string `json:adjectives`
	Nouns               []string `json:nouns`
	FirstNamesFemale    []string `json:firstNamesFemale`
	FirstNamesMale      []string `json:firstNamesMale`
	LastNames           []string `json:lastNames`
	Domains             []string `json:domains`
	People              []string `json:people`
	StreetTypes         []string `json:streetTypes` // Taken from https://github.com/tomharris/random_data/blob/master/lib/random_data/locations.rb
	Paragraphs          []string `json:paragraphs`  // Taken from feedbooks.com
	Countries           []string `json:countries`   // Fetched from the world bank at http://siteresources.worldbank.org/DATASTATISTICS/Resources/CLASS.XLS
	CountriesThreeChars []string `json:countriesThreeChars`
	CountriesTwoChars   []string `json:countriesTwoChars`
	Currencies          []string `json:currencies` //https://github.com/OpenBookPrices/country-data
	Cities              []string `json:cities`
	States              []string `json:states`
	StatesSmall         []string `json:statesSmall`
}

var jsonData = jsonContent{}

func init() {
	jsonData = jsonContent{}

	err := json.Unmarshal(data, &jsonData)

	if err != nil {
		log.Fatal(err)
	}
}

func seedAndReturnRandom(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func seedAndReturnRandomFloat() float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()
}

// Returns a random part of a slice
func randomFrom(source []string) string {
	return source[seedAndReturnRandom(len(source))]
}

// Returns a random first name, gender decides the gender of the name
func FirstName(gender int) string {
	var name = ""
	switch gender {
	case Male:
		name = randomFrom(jsonData.FirstNamesMale)
		break
	case Female:
		name = randomFrom(jsonData.FirstNamesFemale)
		break
	default:
		rand.Seed(time.Now().UnixNano())
		name = FirstName(rand.Intn(2))
		break
	}
	return name
}

// Returns a random last name
func LastName() string {
	return randomFrom(jsonData.LastNames)
}

// Returns a combinaton of FirstName LastName randomized, gender decides the gender of the name
func FullName(gender int) string {
	return FirstName(gender) + " " + LastName()
}

// Returns a random email
func Email() string {
	return strings.ToLower(FirstName(RandomGender)+LastName()) + "@" + randomFrom(jsonData.Domains)
}

// Returns a random country, countryStyle decides what kind of format the returned country will have
func Country(countryStyle int64) string {
	country := ""
	switch countryStyle {

	default:

	case FullCountry:
		country = randomFrom(jsonData.Countries)
		break

	case TwoCharCountry:
		country = randomFrom(jsonData.CountriesTwoChars)
		break

	case ThreeCharCountry:
		country = randomFrom(jsonData.CountriesThreeChars)
		break
	}
	return country
}

// Returns a random currency under ISO 4217 format
func Currency() string {
	return randomFrom(jsonData.Currencies)
}

// Returns a random city
func City() string {
	return randomFrom(jsonData.Cities)
}

// Returns a random american state
func State(typeOfState int) string {
	if typeOfState == Small {
		return randomFrom(jsonData.StatesSmall)
	} else {
		return randomFrom(jsonData.States)
	}
	return ""
}

// Returns a random fake street name
func Street() string {
	return fmt.Sprintf("%s %s", randomFrom(jsonData.People), randomFrom(jsonData.StreetTypes))
}

// Returns an american style address
func Address() string {
	return fmt.Sprintf("%d %s,\n%s, %s, %s", Number(100), Street(), City(), State(Small), PostalCode("US"))
}

// Returns a random paragraph
func Paragraph() string {
	return randomFrom(jsonData.Paragraphs)
}

// Returns a random number, if only one integer is supplied it is treated as the max value to return
// if a second argument is supplied it returns a number between (and including) the two numbers
func Number(numberRange ...int) int {
	nr := 0
	rand.Seed(time.Now().UnixNano())
	if len(numberRange) > 1 {
		nr = 1
		nr = seedAndReturnRandom(numberRange[1]-numberRange[0]) + numberRange[0]
	} else {
		nr = seedAndReturnRandom(numberRange[0])
	}
	return nr
}

func Decimal(numberRange ...int) float64 {
	nr := 0.0
	rand.Seed(time.Now().UnixNano())
	if len(numberRange) > 1 {
		nr = 1.0
		nr = seedAndReturnRandomFloat()*(float64(numberRange[1])-float64(numberRange[0])) + float64(numberRange[0])
	} else {
		nr = seedAndReturnRandomFloat() * float64(numberRange[0])
	}

	if len(numberRange) > 2 {
		sf := strconv.FormatFloat(nr, 'f', numberRange[2], 64)
		nr, _ = strconv.ParseFloat(sf, 64)
	}
	return nr
}

func StringNumberExt(numberPairs int, separator string, numberOfDigits int) string {
	numberString := ""

	for i := 0; i < numberPairs; i++ {
		for d := 0; d < numberOfDigits; d++ {
			numberString += fmt.Sprintf("%d", Number(0, 9))
		}

		if i+1 != numberPairs {
			numberString += separator
		}
	}

	return numberString
}

// Returns a random number as a string
func StringNumber(numberPairs int, separator string) string {
	return StringNumberExt(numberPairs, separator, 2)
}

func Boolean() bool {
	nr := seedAndReturnRandom(2)
	return nr != 0
}

// Returns a random noun
func Noun() string {
	return randomFrom(jsonData.Nouns)
}

// Returns a random adjective
func Adjective() string {
	return randomFrom(jsonData.Adjectives)
}

func uppercaseFirstLetter(word string) string {
	a := []rune(word)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}

func lowercaseFirstLetter(word string) string {
	a := []rune(word)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

// Returns a silly name, useful for randomizing naming of things
func SillyName() string {
	return uppercaseFirstLetter(Noun()) + Adjective()
}

// Returns a valid IPv4 address as string
func IpV4Address() string {
	blocks := []string{}
	for i := 0; i < 4; i++ {
		number := seedAndReturnRandom(255)
		blocks = append(blocks, strconv.Itoa(number))
	}

	return strings.Join(blocks, ".")
}
