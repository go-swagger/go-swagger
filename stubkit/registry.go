package stubkit

import (
	"log"
	"math/rand"
	"strings"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/go-swagger/go-swagger/spec"
	"github.com/manveru/faker"
)

// StringGen a type for things that can generate strings
type StringGen func() string

var fak *faker.Faker

func init() {
	fak, _ = faker.New("en")
}

var stringGen = map[string]StringGen{
	"firstname": fak.FirstName,
	"lastname":  fak.LastName,
	"login":     fak.UserName,
	"username":  fak.UserName,
	"email":     fak.Email,
	"emailfree": fak.FreeEmail,
	"emailsafe": fak.SafeEmail,
	"jobtitle":  fak.JobTitle,
	"name":      fak.Name,
	"city":      fak.City,
	"street":    fak.StreetAddress,
	"postcode":  fak.PostCode,
	"stateabbr": fak.StateAbbr,
	"state":     fak.State,
	"country":   fak.Country,
	"ipv4":      func() string { return fak.IPv4Address().String() },
	"ip":        func() string { return fak.IPv4Address().String() },
	"ipv6":      func() string { return fak.IPv6Address().String() },
	"url":       fak.URL,
	"domain":    fak.DomainName,
	"phone":     fak.PhoneNumber,
	"mobile":    fak.CellPhoneNumber,
	"word":      randomdata.Noun,
	"sentence":  func() string { return fak.Sentence(rand.Intn(10), false) },
	"paragraph": func() string { return fak.Paragraph(rand.Intn(10), false) },
	"text":      func() string { return strings.Join(fak.Paragraphs(rand.Intn(6), false), "\n") },
}

// StringGenerator creates a string gen for a field name or generator hint
func StringGenerator(name string) StringGen {
	if gen, ok := stringGen[name]; ok {
		return gen
	}
	return func() string { return fak.Characters(12) }
}

// SchemaGen a type for things that can generate data for a schema definition
// composed out of the more primitive data generators
type SchemaGen func() interface{}

// SchemaGenerator creates a generator for a schema
func SchemaGenerator(schema *spec.Schema, aggregator SchemaGen) SchemaGen {
	if len(schema.Enum) > 0 {
		return func() interface{} { return schema.Enum[rand.Intn(len(schema.Enum))] }
	}
	for k, v := range schema.Properties {
		log.Printf("generating for %q: %v", k, v)
	}
	for k, v := range schema.PatternProperties {
		log.Printf("generating for pattern %q: %v", k, v)
	}
	if schema.AdditionalProperties != nil && schema.AdditionalProperties.Schema != nil {
		SchemaGenerator(schema.AdditionalProperties.Schema, nil)
	}
	for _, v := range schema.AllOf {
		SchemaGenerator(&v, nil)
	}
	// for _, v := range schema.
	return func() interface{} { return nil }
}
