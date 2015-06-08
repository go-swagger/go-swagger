package swag

import (
	"math"
	"strconv"
	"strings"
)

// same as ECMA Number.MAX_SAFE_INTEGER and Number.MIN_SAFE_INTEGER
const (
	maxJSONFloat = float64(1<<53 - 1)  // 9007199254740991.0 	 	 2^53 - 1
	minJSONFloat = -float64(1<<53 - 1) //-9007199254740991.0	-2^53 - 1
)

// IsFloat64AJSONInteger allow for integers [-2^53, 2^53-1] inclusive
func IsFloat64AJSONInteger(f float64) bool {
	if math.IsNaN(f) || math.IsInf(f, 0) || f < minJSONFloat || f > maxJSONFloat {
		return false
	}

	return f == float64(int64(f)) || f == float64(uint64(f))
}

var evaluatesAsTrue = map[string]struct{}{
	"true":     struct{}{},
	"1":        struct{}{},
	"yes":      struct{}{},
	"ok":       struct{}{},
	"y":        struct{}{},
	"on":       struct{}{},
	"selected": struct{}{},
	"checked":  struct{}{},
	"t":        struct{}{},
	"enabled":  struct{}{},
}

// ConvertBool turn a string into a boolean
func ConvertBool(str string) (bool, error) {
	_, ok := evaluatesAsTrue[strings.ToLower(str)]
	return ok, nil
}

// ConvertFloat32 turn a string into a float32
func ConvertFloat32(str string) (float32, error) {
	f, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

// ConvertFloat64 turn a string into a float64
func ConvertFloat64(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

// ConvertInt8 turn a string into int8 boolean
func ConvertInt8(str string) (int8, error) {
	i, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(i), nil
}

// ConvertInt16 turn a string into a int16
func ConvertInt16(str string) (int16, error) {
	i, err := strconv.ParseInt(str, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(i), nil
}

// ConvertInt32 turn a string into a int32
func ConvertInt32(str string) (int32, error) {
	i, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

// ConvertInt64 turn a string into a int64
func ConvertInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

// ConvertUint8 turn a string into a uint8
func ConvertUint8(str string) (uint8, error) {
	i, err := strconv.ParseUint(str, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(i), nil
}

// ConvertUint16 turn a string into a uint16
func ConvertUint16(str string) (uint16, error) {
	i, err := strconv.ParseUint(str, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(i), nil
}

// ConvertUint32 turn a string into a uint32
func ConvertUint32(str string) (uint32, error) {
	i, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}

// ConvertUint64 turn a string into a uint64
func ConvertUint64(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}
