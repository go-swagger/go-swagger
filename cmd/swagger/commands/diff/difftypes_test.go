package diff

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpecChangeCode(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	c := NoChangeDetected
	assert.Equal(t, toLongStringSpecChangeCode[NoChangeDetected], c.Description())
	assert.Equal(t, "UNDEFINED", SpecChangeCode(9999999999999).Description())

	res, err := json.Marshal(c)
	require.NoError(t, err)
	assert.JSONEq(t, `"NoChangeDetected"`, string(res))

	var d SpecChangeCode
	in := []byte(`"NoChangeDetected"`)
	err = json.Unmarshal(in, &d)
	require.NoError(t, err)
	assert.Equal(t, NoChangeDetected, d)

	in = []byte(`"dummy"`) // invalid enum
	err = json.Unmarshal(in, &d)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown enum value")

	in = []byte(`{"dummy"`) // invalid json
	err = json.Unmarshal(in, &d)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JSON")
}

func TestCompatibiliyCode(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	c := Breaking
	assert.Equal(t, toStringCompatibility[Breaking], c.String())

	res, err := json.Marshal(c)
	require.NoError(t, err)
	assert.JSONEq(t, `"Breaking"`, string(res))

	var d Compatibility
	in := []byte(`"Breaking"`)
	err = json.Unmarshal(in, &d)
	require.NoError(t, err)
	assert.Equal(t, Breaking, d)

	in = []byte(`"dummy"`) // invalid enum
	err = json.Unmarshal(in, &d)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown enum value")

	in = []byte(`{"dummy"`) // invalid json
	err = json.Unmarshal(in, &d)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JSON")
}
