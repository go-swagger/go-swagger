package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYAMLToJSON(t *testing.T) {

	data := make(map[interface{}]interface{})
	data[1] = "the int key value"
	data["name"] = "a string value"

	d, err := YAMLToJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"1":"the int key value","name":"a string value"}`), d)

	data[true] = "the bool value"
	d, err = YAMLToJSON(data)
	assert.Error(t, err)
	assert.Nil(t, d)

	delete(data, true)

	tag := make(map[interface{}]interface{})
	tag["name"] = "tag name"
	data["tag"] = tag

	d, err = YAMLToJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"1":"the int key value","name":"a string value","tag":{"name":"tag name"}}`), d)

	tag = make(map[interface{}]interface{})
	tag[true] = "bool tag name"
	data["tag"] = tag

	d, err = YAMLToJSON(data)
	assert.Error(t, err)
	assert.Nil(t, d)

	var lst []interface{}
	lst = append(lst, "hello")

	d, err = YAMLToJSON(lst)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`["hello"]`), d)

	lst = append(lst, data)

	d, err = YAMLToJSON(lst)
	assert.Error(t, err)
	assert.Nil(t, d)
}
