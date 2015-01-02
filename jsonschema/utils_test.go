package jsonschema

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestIsFloat64IntegerA(t *testing.T) {

	assert.False(t, isFloat64AnInteger(math.Inf(+1)))
	assert.False(t, isFloat64AnInteger(math.Inf(-1)))
	assert.False(t, isFloat64AnInteger(math.NaN()))
	assert.False(t, isFloat64AnInteger(math.Float64frombits((1<<11-1)<<52|1)))
	assert.True(t, isFloat64AnInteger(float64(0.0)))
	assert.True(t, isFloat64AnInteger(-float64(0.0)))
	assert.False(t, isFloat64AnInteger(float64(0.5)))
	assert.True(t, isFloat64AnInteger(float64(1.0)))
	assert.True(t, isFloat64AnInteger(-float64(1.0)))
	assert.False(t, isFloat64AnInteger(float64(1.1)))
	assert.True(t, isFloat64AnInteger(float64(131.0)))
	assert.True(t, isFloat64AnInteger(-float64(131.0)))
	assert.True(t, isFloat64AnInteger(float64(1<<52-1)))
	assert.True(t, isFloat64AnInteger(-float64(1<<52-1)))
	assert.True(t, isFloat64AnInteger(float64(1<<52)))
	assert.True(t, isFloat64AnInteger(-float64(1<<52)))
	assert.True(t, isFloat64AnInteger(float64(1<<53-1)))
	assert.True(t, isFloat64AnInteger(-float64(1<<53-1)))
	assert.True(t, isFloat64AnInteger(float64(1<<53-1)))
	assert.False(t, isFloat64AnInteger(float64(1<<53)))
	assert.False(t, isFloat64AnInteger(-float64(1<<53)))
	assert.False(t, isFloat64AnInteger(float64(1<<63)))
	assert.False(t, isFloat64AnInteger(-float64(1<<63)))
	assert.False(t, isFloat64AnInteger(math.Nextafter(float64(1<<63), math.MaxFloat64)))
	assert.False(t, isFloat64AnInteger(-math.Nextafter(float64(1<<63), math.MaxFloat64)))
	assert.False(t, isFloat64AnInteger(float64(1<<70+3<<21)))
	assert.False(t, isFloat64AnInteger(-float64(1<<70+3<<21)))

	assert.False(t, isFloat64AnInteger(math.Nextafter(float64(9007199254740991.0), math.MaxFloat64)))
	assert.True(t, isFloat64AnInteger(float64(9007199254740991.0)))

	assert.True(t, isFloat64AnInteger(float64(-9007199254740991.0)))
	assert.False(t, isFloat64AnInteger(math.Nextafter(float64(-9007199254740991.0), -math.MaxFloat64)))
}

func TestIsFloat64Integer(t *testing.T) {
	// fails. MaxUint64 is to large for JS anyway, so we can ignore it here.
	//assert.True(t, isFloat64AnInteger(float64(math.MaxUint64)))

	assert.False(t, isFloat64AnInteger(math.MaxInt64))
	assert.False(t, isFloat64AnInteger(1<<62))
	assert.False(t, isFloat64AnInteger(math.MinInt64))
	assert.True(t, isFloat64AnInteger(100100100100))
	assert.True(t, isFloat64AnInteger(-100100100100))
	assert.True(t, isFloat64AnInteger(100100100))
	assert.True(t, isFloat64AnInteger(-100100100))
	assert.True(t, isFloat64AnInteger(100100))
	assert.True(t, isFloat64AnInteger(-100100))
	assert.True(t, isFloat64AnInteger(100))
	assert.True(t, isFloat64AnInteger(-100))
	assert.True(t, isFloat64AnInteger(-0))
	assert.True(t, isFloat64AnInteger(-1))
	assert.True(t, isFloat64AnInteger(1))
	assert.True(t, isFloat64AnInteger(0))
	assert.True(t, isFloat64AnInteger(77))
	assert.True(t, isFloat64AnInteger(-77))
	assert.True(t, isFloat64AnInteger(1e10))
	assert.True(t, isFloat64AnInteger(-1e10))
	assert.True(t, isFloat64AnInteger(100100100.0))
	assert.True(t, isFloat64AnInteger(-100100100.0))

	assert.False(t, isFloat64AnInteger(100100100100.1))
	assert.False(t, isFloat64AnInteger(-100100100100.1))
	assert.False(t, isFloat64AnInteger(math.MaxFloat64))
	assert.False(t, isFloat64AnInteger(-math.MaxFloat64))
	assert.False(t, isFloat64AnInteger(1.1))
	assert.False(t, isFloat64AnInteger(-1.1))
	assert.False(t, isFloat64AnInteger(1.000000000001))
	assert.False(t, isFloat64AnInteger(-1.000000000001))
	assert.False(t, isFloat64AnInteger(1e-10))
	assert.False(t, isFloat64AnInteger(-1e-10))

	assert.False(t, isFloat64AnInteger(1.001))
	assert.False(t, isFloat64AnInteger(-1.001))

	assert.False(t, isFloat64AnInteger(0.0001))

	assert.False(t, isFloat64AnInteger(1<<62))
	assert.False(t, isFloat64AnInteger(math.MinInt64))
	assert.False(t, isFloat64AnInteger(math.MaxInt64))
	assert.False(t, isFloat64AnInteger(-1<<62))

	assert.False(t, isFloat64AnInteger(1e-10))
	assert.False(t, isFloat64AnInteger(-1e-10))

	assert.False(t, isFloat64AnInteger(1.000000000001))
	assert.False(t, isFloat64AnInteger(-1.000000000001))

	assert.False(t, isFloat64AnInteger(4.611686018427387904e7))
	assert.False(t, isFloat64AnInteger(-4.611686018427387904e7))
}

func TestValidationErrorFormatNumber(t *testing.T) {
	assert.Equal(t, "1", validationErrorFormatNumber(1))
	assert.Equal(t, "-1", validationErrorFormatNumber(-1))
	assert.Equal(t, "0", validationErrorFormatNumber(0))
	// unfortunately, can not be recognized as float
	assert.Equal(t, "0", validationErrorFormatNumber(0.0))

	assert.Equal(t, "1.001", validationErrorFormatNumber(1.001))
	assert.Equal(t, "-1.001", validationErrorFormatNumber(-1.001))
	assert.Equal(t, "0.0001", validationErrorFormatNumber(0.0001))

	// casting math.MaxInt64 (1<<63 -1) to float back to int64
	// becomes negative. obviousely because of bit missinterpretation.
	// so simply test a slightly smaller "large" integer here
	assert.Equal(t, "4.611686018427388e+18", validationErrorFormatNumber(1<<62))
	// with negative int64 max works
	assert.Equal(t, "-9.223372036854776e+18", validationErrorFormatNumber(math.MinInt64))
	assert.Equal(t, "-4.611686018427388e+18", validationErrorFormatNumber(-1<<62))

	assert.Equal(t, "10000000000", validationErrorFormatNumber(1e10))
	assert.Equal(t, "-10000000000", validationErrorFormatNumber(-1e10))

	assert.Equal(t, "1.000000000001", validationErrorFormatNumber(1.000000000001))
	assert.Equal(t, "-1.000000000001", validationErrorFormatNumber(-1.000000000001))
	assert.Equal(t, "1e-10", validationErrorFormatNumber(1e-10))
	assert.Equal(t, "-1e-10", validationErrorFormatNumber(-1e-10))
	assert.Equal(t, "4.6116860184273876e+07", validationErrorFormatNumber(4.611686018427387904e7))
	assert.Equal(t, "-4.6116860184273876e+07", validationErrorFormatNumber(-4.611686018427387904e7))
}
