package parser

import (
	goparser "go/parser"
	"go/token"
	"testing"

	"github.com/casualjim/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

func TestSetInfoVersion(t *testing.T) {
	info := new(spec.Info)
	err := setInfoVersion(info, []string{"0.0.1"})
	assert.NoError(t, err)
	assert.Equal(t, "0.0.1", info.Version)
}

func TestSetInfoTitle(t *testing.T) {
	info := new(spec.Info)
	err := setInfoTitle(info, []string{"A title in", "2 parts"})
	assert.NoError(t, err)
	assert.Equal(t, "A title in\n2 parts", info.Title)
}

func TestSetInfoTOS(t *testing.T) {
	info := new(spec.Info)
	err := setInfoTOS(info, []string{"A TOS in", "2 parts"})
	assert.NoError(t, err)
	assert.Equal(t, "A TOS in\n2 parts", info.TermsOfService)
}

func TestSetInfoDescription(t *testing.T) {
	info := new(spec.Info)
	err := setInfoDescription(info, []string{"A description in", "2 parts"})
	assert.NoError(t, err)
	assert.Equal(t, "A description in\n2 parts", info.Description)
}

func TestSetInfoLicense(t *testing.T) {
	info := new(spec.Info)
	err := setInfoLicense(info, []string{"MIT http://license.org/MIT"})
	assert.NoError(t, err)
	assert.Equal(t, "MIT", info.License.Name)
	assert.Equal(t, "http://license.org/MIT", info.License.URL)
}

func TestSetInfoContact(t *testing.T) {
	info := new(spec.Info)
	err := setInfoContact(info, []string{"Homer J. Simpson <homer@simpsons.com> http://simpsons.com"})
	assert.NoError(t, err)
	assert.Equal(t, "Homer J. Simpson", info.Contact.Name)
	assert.Equal(t, "homer@simpsons.com", info.Contact.Email)
	assert.Equal(t, "http://simpsons.com", info.Contact.URL)
}

func TestParseInfo(t *testing.T) {
	parser := newAPIInfoParser()
	docFile := "../fixtures/goparsing/petstoreapp/doc.go"
	fileSet := token.NewFileSet()
	fileTree, err := goparser.ParseFile(fileSet, docFile, nil, goparser.ParseComments)
	if err != nil {
		t.FailNow()
	}
	info, err := parser.Parse(fileTree)
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, "0.0.1", info.Version)
	assert.Equal(t, "Petstore API", info.Title)
	assert.Equal(t, "there are no TOS at this moment, use at your own risk we take no responsibility", info.TermsOfService)

	descr := `the purpose of this application is to provide an application
that is using plain go code to define an API

This should demonstrate all the possible comment annotations
that are available to turn go code into a fully compliant swagger 2.0 spec`
	assert.Equal(t, descr, info.Description)

	assert.NotNil(t, info.License)
	assert.Equal(t, "MIT", info.License.Name)
	assert.Equal(t, "http://opensource.org/licenses/MIT", info.License.URL)

	assert.NotNil(t, info.Contact)
	assert.Equal(t, "John Doe", info.Contact.Name)
	assert.Equal(t, "john.doe@example.com", info.Contact.Email)
	assert.Equal(t, "http://john.doe.com", info.Contact.URL)

}
