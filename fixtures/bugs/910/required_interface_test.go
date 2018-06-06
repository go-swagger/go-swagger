//+build ignore

package main

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-swagger/go-swagger/fixtures/bugs/910/gen-fixture-910/models"
	"github.com/stretchr/testify/assert"
)

func Test_Required(t *testing.T) {
	x := models.GetMytestOKBody{}

	res := x.Validate(strfmt.Default)
	assert.Error(t, res)
	t.Logf("Empty err=%v", res)

	x.Bar = swag.Int64(10)
	res = x.Validate(strfmt.Default)
	assert.Error(t, res)
	t.Logf("Empty err=%v", res)

	x.Foo = map[string]string{"a": "val1", "b": "val2"}
	res = x.Validate(strfmt.Default)
	assert.NoError(t, res)

}
