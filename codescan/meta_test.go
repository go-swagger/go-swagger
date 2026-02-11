// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package codescan

import (
	goparser "go/parser"
	"go/token"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"

	"github.com/go-openapi/spec"
)

func TestSetInfoVersion(t *testing.T) {
	info := new(spec.Swagger)
	err := setInfoVersion(info, []string{"0.0.1"})
	require.NoError(t, err)
	assert.EqualT(t, "0.0.1", info.Info.Version)
}

func TestSetInfoLicense(t *testing.T) {
	info := new(spec.Swagger)
	err := setInfoLicense(info, []string{"MIT http://license.org/MIT"})
	require.NoError(t, err)
	assert.EqualT(t, "MIT", info.Info.License.Name)
	assert.EqualT(t, "http://license.org/MIT", info.Info.License.URL)
}

func TestSetInfoContact(t *testing.T) {
	info := new(spec.Swagger)
	err := setInfoContact(info, []string{"Homer J. Simpson <homer@simpsons.com> http://simpsons.com"})
	require.NoError(t, err)
	assert.EqualT(t, "Homer J. Simpson", info.Info.Contact.Name)
	assert.EqualT(t, "homer@simpsons.com", info.Info.Contact.Email)
	assert.EqualT(t, "http://simpsons.com", info.Info.Contact.URL)
}

func TestParseInfo(t *testing.T) {
	swspec := new(spec.Swagger)
	parser := newMetaParser(swspec)
	docFile := "../fixtures/goparsing/classification/doc.go"
	fileSet := token.NewFileSet()
	fileTree, err := goparser.ParseFile(fileSet, docFile, nil, goparser.ParseComments)
	if err != nil {
		t.FailNow()
	}

	err = parser.Parse(fileTree.Doc)

	require.NoError(t, err)
	verifyInfo(t, swspec.Info)
}

func TestParseSwagger(t *testing.T) {
	swspec := new(spec.Swagger)
	parser := newMetaParser(swspec)
	docFile := "../fixtures/goparsing/classification/doc.go"
	fileSet := token.NewFileSet()
	fileTree, err := goparser.ParseFile(fileSet, docFile, nil, goparser.ParseComments)
	if err != nil {
		t.FailNow()
	}

	err = parser.Parse(fileTree.Doc)
	verifyMeta(t, swspec)

	require.NoError(t, err)
}

func verifyMeta(t *testing.T, doc *spec.Swagger) {
	assert.NotNil(t, doc)
	verifyInfo(t, doc.Info)
	assert.Equal(t, []string{"application/json", "application/xml"}, doc.Consumes)
	assert.Equal(t, []string{"application/json", "application/xml"}, doc.Produces)
	assert.Equal(t, []string{"http", "https"}, doc.Schemes)
	assert.Equal(t, []map[string][]string{{"api_key": {}}}, doc.Security)
	expectedSecuritySchemaKey := spec.SecurityScheme{
		SecuritySchemeProps: spec.SecuritySchemeProps{
			Type: "apiKey",
			In:   "header",
			Name: "KEY",
		},
	}
	expectedSecuritySchemaOAuth := spec.SecurityScheme{
		SecuritySchemeProps: spec.SecuritySchemeProps{
			Type:             "oauth2",
			In:               "header",
			AuthorizationURL: "/oauth2/auth",
			TokenURL:         "/oauth2/token",
			Flow:             "accessCode",
			Scopes: map[string]string{
				"bla1": "foo1",
				"bla2": "foo2",
			},
		},
	}
	expectedExtensions := spec.Extensions{
		"x-meta-array": []any{
			"value1",
			"value2",
		},
		"x-meta-array-obj": []any{
			map[string]any{
				"name":  "obj",
				"value": "field",
			},
		},
		"x-meta-value": "value",
	}
	expectedInfoExtensions := spec.Extensions{
		"x-info-array": []any{
			"value1",
			"value2",
		},
		"x-info-array-obj": []any{
			map[string]any{
				"name":  "obj",
				"value": "field",
			},
		},
		"x-info-value": "value",
	}
	assert.NotNil(t, doc.SecurityDefinitions["api_key"])
	assert.NotNil(t, doc.SecurityDefinitions["oauth2"])
	assert.Equal(t, spec.SecurityDefinitions{"api_key": &expectedSecuritySchemaKey, "oauth2": &expectedSecuritySchemaOAuth}, doc.SecurityDefinitions)
	assert.Equal(t, expectedExtensions, doc.Extensions)
	assert.Equal(t, expectedInfoExtensions, doc.Info.Extensions)
	assert.EqualT(t, "localhost", doc.Host)
	assert.EqualT(t, "/v2", doc.BasePath)
}

func verifyInfo(t *testing.T, info *spec.Info) {
	t.Helper()

	assert.NotNil(t, info)
	assert.EqualT(t, "0.0.1", info.Version)
	assert.EqualT(t, "there are no TOS at this moment, use at your own risk we take no responsibility", info.TermsOfService)
	assert.EqualT(t, "Petstore API.", info.Title)
	descr := `the purpose of this application is to provide an application
that is using plain go code to define an API

This should demonstrate all the possible comment annotations
that are available to turn go code into a fully compliant swagger 2.0 spec`
	assert.EqualT(t, descr, info.Description)

	require.NotNil(t, info.License)
	assert.EqualT(t, "MIT", info.License.Name)
	assert.EqualT(t, "http://opensource.org/licenses/MIT", info.License.URL)

	require.NotNil(t, info.Contact)
	assert.EqualT(t, "John Doe", info.Contact.Name)
	assert.EqualT(t, "john.doe@example.com", info.Contact.Email)
	assert.EqualT(t, "http://john.doe.com", info.Contact.URL)
}

func TestMoreParseMeta(t *testing.T) {
	for _, docFile := range []string{
		"../fixtures/goparsing/meta/v1/doc.go",
		"../fixtures/goparsing/meta/v2/doc.go",
		"../fixtures/goparsing/meta/v3/doc.go",
		"../fixtures/goparsing/meta/v4/doc.go",
	} {
		swspec := new(spec.Swagger)
		parser := newMetaParser(swspec)
		fileSet := token.NewFileSet()
		fileTree, err := goparser.ParseFile(fileSet, docFile, nil, goparser.ParseComments)
		if err != nil {
			t.FailNow()
		}

		err = parser.Parse(fileTree.Doc)
		require.NoError(t, err)
		assert.EqualT(t, "there are no TOS at this moment, use at your own risk we take no responsibility", swspec.Info.TermsOfService)
		/*
			jazon, err := json.MarshalIndent(swspec.Info, "", " ")
			require.NoError(t, err)
			t.Logf("%v", string(jazon))
		*/
	}
}
