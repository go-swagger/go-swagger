package strfmt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testFormat string

func (t testFormat) MarshalText() ([]byte, error) {
	return []byte(string(t)), nil
}

func (t *testFormat) UnmarshalText(b []byte) error {
	*t = testFormat(string(b))
	return nil
}

func isTestFormat(s string) bool {
	return strings.HasPrefix(s, "tf")
}

type tf2 string

func (t tf2) MarshalText() ([]byte, error) {
	return []byte(string(t)), nil
}

func (t *tf2) UnmarshalText(b []byte) error {
	*t = tf2(string(b))
	return nil
}

func istf2(s string) bool {
	return strings.HasPrefix(s, "af")
}

type bf string

func (t bf) MarshalText() ([]byte, error) {
	return []byte(string(t)), nil
}

func (t *bf) UnmarshalText(b []byte) error {
	*t = bf(string(b))
	return nil
}

func isbf(s string) bool {
	return strings.HasPrefix(s, "bf")
}

func istf3(s string) bool {
	return strings.HasPrefix(s, "ff")
}

func init() {
	tf := testFormat("")
	Default.Add("test-format", &tf, isTestFormat)
}

func TestFormatRegistry(t *testing.T) {
	f2 := tf2("")
	f3 := bf("")
	registry := NewFormats()

	assert.True(t, registry.ContainsName("test-format"))
	assert.True(t, registry.ContainsName("testformat"))
	assert.False(t, registry.ContainsName("ttt"))

	assert.True(t, registry.Validates("testformat", "tfa"))
	assert.False(t, registry.Validates("testformat", "ffa"))

	assert.True(t, registry.Add("tf2", &f2, istf2))
	assert.True(t, registry.ContainsName("tf2"))
	assert.False(t, registry.ContainsName("tfw"))
	assert.True(t, registry.Validates("tf2", "afa"))

	assert.False(t, registry.Add("tf2", &f3, isbf))
	assert.True(t, registry.ContainsName("tf2"))
	assert.False(t, registry.ContainsName("tfw"))
	assert.True(t, registry.Validates("tf2", "bfa"))
	assert.False(t, registry.Validates("tf2", "afa"))

	assert.False(t, registry.Add("tf2", &f2, istf2))
	assert.True(t, registry.Add("tf3", &f2, istf3))
	assert.True(t, registry.ContainsName("tf3"))
	assert.True(t, registry.ContainsName("tf2"))
	assert.False(t, registry.ContainsName("tfw"))
	assert.True(t, registry.Validates("tf3", "ffa"))

	assert.True(t, registry.DelByName("tf3"))
	assert.True(t, registry.Add("tf3", &f2, istf3))

	assert.True(t, registry.DelByName("tf3"))
	assert.False(t, registry.DelByName("unknown"))
	assert.False(t, registry.Validates("unknown", ""))
}
