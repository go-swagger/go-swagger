package spec

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnknownSpecVersion(t *testing.T) {
	_, err := New([]byte{}, "0.9")
	assert.Error(t, err)
}

func TestDefaultsTo20(t *testing.T) {
	d, err := New(PetStoreJSONMessage, "")

	assert.NoError(t, err)
	assert.NotNil(t, d)
	assert.Equal(t, "2.0", d.Version())
	assert.Equal(t, "2.0", d.data["swagger"].(string))
}

func TestValidatesValidSchema(t *testing.T) {
	d, err := New(PetStoreJSONMessage, "")

	assert.NoError(t, err)
	assert.NotNil(t, d)
	res := d.Validate()
	assert.NotNil(t, res)
	assert.True(t, res.Valid())
	assert.Empty(t, res.Errors())
}

func TestFailsInvalidSchema(t *testing.T) {
	d, err := New(InvalidJSONMessage, "")

	assert.NoError(t, err)
	assert.NotNil(t, d)

	res := d.Validate()
	assert.NotNil(t, res)
	assert.False(t, res.Valid())
	assert.NotEmpty(t, res.Errors())
}

func TestFailsInvalidJSON(t *testing.T) {
	_, err := New(json.RawMessage([]byte("{]")), "")

	assert.Error(t, err)
}
