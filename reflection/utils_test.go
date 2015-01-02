package reflection

// import (
// 	"reflect"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// )

// func int64Ptr(i int64) *int64 { return &i }

// func TestIsNumber(t *testing.T) {
// 	values := []interface{}{
// 		1,
// 		int(1),
// 		int8(1),
// 		int16(1),
// 		int32(1),
// 		int64(1),
// 		uint(1),
// 		uint8(1),
// 		uint16(1),
// 		uint32(1),
// 		uint64(1),
// 		float32(1),
// 		float64(1),
// 		int64Ptr(1),
// 	}
// 	for _, v := range values {
// 		assert.True(t, IsNumeric(reflect.ValueOf(v)), "expected %#v to be numeric", v)
// 	}

// 	notNumbers := []interface{}{
// 		"a",
// 		&struct{ A string }{"A"},
// 		struct{ A string }{"A"},
// 		map[string]interface{}{},
// 		[]interface{}{},
// 	}

// 	for _, v := range notNumbers {
// 		assert.False(t, IsNumeric(reflect.ValueOf(v)), "expected %#v to NOT be numeric", v)
// 	}
// }

// func TestIsInteger(t *testing.T) {
// 	values := []interface{}{
// 		1,
// 		int(1),
// 		int8(1),
// 		int16(1),
// 		int32(1),
// 		int64(1),
// 		uint(1),
// 		uint8(1),
// 		uint16(1),
// 		uint32(1),
// 		uint64(1),
// 		1.0,
// 		float32(1.0),
// 		float64(1.0),
// 		int64Ptr(1),
// 	}
// 	for _, v := range values {
// 		assert.True(t, IsInteger(reflect.ValueOf(v)), "expected %#v to be an integer", v)
// 	}

// 	notNumbers := []interface{}{
// 		float32(1.1),
// 		float64(1.2),
// 		"a",
// 		&struct{ A string }{"A"},
// 		struct{ A string }{"A"},
// 		map[string]interface{}{},
// 		[]interface{}{},
// 	}

// 	for _, v := range notNumbers {
// 		assert.False(t, IsInteger(reflect.ValueOf(v)), "expected %#v to NOT be integer", v)
// 	}
// }

// func TestIsBasicType(t *testing.T) {
// 	values := []interface{}{
// 		true,
// 		1,
// 		int(1),
// 		int8(1),
// 		int16(1),
// 		int32(1),
// 		int64(1),
// 		uint(1),
// 		uint8(1),
// 		uint16(1),
// 		uint32(1),
// 		uint64(1),
// 		1.0,
// 		float32(1.0),
// 		float64(1.0),
// 		int64Ptr(1),
// 		complex64(1),
// 		complex128(1),
// 		byte(1),
// 		"blah",
// 		'b',
// 		uintptr(1),
// 		time.Now(),
// 	}

// 	for _, v := range values {
// 		assert.True(t, IsSimpleType(reflect.ValueOf(v)), "expected %#v to be a simple type", v)
// 	}

// 	notNumbers := []interface{}{
// 		&struct{ A string }{"A"},
// 		struct{ A string }{"A"},
// 		map[string]interface{}{},
// 		[]interface{}{},
// 	}
// 	for _, v := range notNumbers {
// 		assert.False(t, IsSimpleType(reflect.ValueOf(v)), "expected %#v to NOT be a simple type", v)
// 	}
// }
