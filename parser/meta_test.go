package parser

import (
	goparser "go/parser"
	"go/token"
	"testing"

	"github.com/casualjim/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

func TestSetInfoVersion(t *testing.T) {
	info := new(spec.Swagger)
	err := setInfoVersion(info, []string{"0.0.1"})
	assert.NoError(t, err)
	assert.Equal(t, "0.0.1", info.Info.Version)
}

func TestSetInfoTitle(t *testing.T) {
	info := new(spec.Swagger)
	err := setInfoTitle(info, []string{"A title in", "2 parts"})
	assert.NoError(t, err)
	assert.Equal(t, "A title in\n2 parts", info.Info.Title)
}

func TestSetInfoTOS(t *testing.T) {
	info := new(spec.Swagger)
	err := setInfoTOS(info, []string{"A TOS in", "2 parts"})
	assert.NoError(t, err)
	assert.Equal(t, "A TOS in\n2 parts", info.Info.TermsOfService)
}

func TestSetInfoDescription(t *testing.T) {
	info := new(spec.Swagger)
	err := setInfoDescription(info, []string{"A description in", "2 parts"})
	assert.NoError(t, err)
	assert.Equal(t, "A description in\n2 parts", info.Info.Description)
}

func TestSetInfoLicense(t *testing.T) {
	info := new(spec.Swagger)
	err := setInfoLicense(info, []string{"MIT http://license.org/MIT"})
	assert.NoError(t, err)
	assert.Equal(t, "MIT", info.Info.License.Name)
	assert.Equal(t, "http://license.org/MIT", info.Info.License.URL)
}

func TestSetInfoContact(t *testing.T) {
	info := new(spec.Swagger)
	err := setInfoContact(info, []string{"Homer J. Simpson <homer@simpsons.com> http://simpsons.com"})
	assert.NoError(t, err)
	assert.Equal(t, "Homer J. Simpson", info.Info.Contact.Name)
	assert.Equal(t, "homer@simpsons.com", info.Info.Contact.Email)
	assert.Equal(t, "http://simpsons.com", info.Info.Contact.URL)
}

func TestParseInfo(t *testing.T) {
	parser := newMetaParser()
	docFile := "../fixtures/goparsing/classification/doc.go"
	fileSet := token.NewFileSet()
	fileTree, err := goparser.ParseFile(fileSet, docFile, nil, goparser.ParseComments)
	if err != nil {
		t.FailNow()
	}
	swspec := new(spec.Swagger)
	err = parser.Parse(fileTree, swspec)

	assert.NoError(t, err)
	assert.Equal(t, "0.0.1", swspec.Info.Version)
	assert.Equal(t, "there are no TOS at this moment, use at your own risk we take no responsibility", swspec.Info.TermsOfService)
	assert.Equal(t, "Petstore API", swspec.Info.Title)
	descr := `the purpose of this application is to provide an application
that is using plain go code to define an API

This should demonstrate all the possible comment annotations
that are available to turn go code into a fully compliant swagger 2.0 spec`
	assert.Equal(t, descr, swspec.Info.Description)

	assert.NotNil(t, swspec.Info.License)
	assert.Equal(t, "MIT", swspec.Info.License.Name)
	assert.Equal(t, "http://opensource.org/licenses/MIT", swspec.Info.License.URL)

	assert.NotNil(t, swspec.Info.Contact)
	assert.Equal(t, "John Doe", swspec.Info.Contact.Name)
	assert.Equal(t, "john.doe@example.com", swspec.Info.Contact.Email)
	assert.Equal(t, "http://john.doe.com", swspec.Info.Contact.URL)

}
