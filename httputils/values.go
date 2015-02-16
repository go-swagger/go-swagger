package httputils

import (
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/util"
	"github.com/casualjim/go-swagger/validate"
)

// ReadStringSliceParam reads a string slice param
func ReadStringSliceParam(name, in, format string, defaultValue []string, values Gettable, required bool) ([]string, *errors.Validation) {
	value := values.Get(name)

	if required && value == "" {
		if len(defaultValue) > 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return nil, err
		}
	}

	if value == "" && len(defaultValue) > 0 {
		return defaultValue, nil
	}

	return split(value, format), nil
}

// ReadBoolSliceParam reads a bool slice param
func ReadBoolSliceParam(name, in, format string, defaultValue []bool, values Gettable, required bool) ([]bool, *errors.Validation) {
	value := values.Get(name)
	if required && value == "" {
		if len(defaultValue) > 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return nil, err
		}
	}

	if value == "" && len(defaultValue) > 0 {
		return defaultValue, nil
	}

	var result []bool
	for _, r := range split(value, format) {
		v, err := util.ConvertBool(r)
		if err != nil {
			return nil, errors.InvalidType(name, in, "bool", r)
		}
		result = append(result, v)
	}
	return result, nil
}

// ReadFloat32SliceParam reads a float32 slice param
func ReadFloat32SliceParam(name, in, format string, defaultValue []float32, values Gettable, required bool) ([]float32, *errors.Validation) {
	value := values.Get(name)

	if required && value == "" {
		if len(defaultValue) > 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return nil, err
		}
	}

	if value == "" && len(defaultValue) > 0 {
		return defaultValue, nil
	}

	var result []float32
	for _, r := range split(value, format) {
		v, err := util.ConvertFloat32(r)
		if err != nil {
			return nil, errors.InvalidType(name, in, "float32", r)
		}
		result = append(result, v)
	}
	return result, nil
}

// ReadFloat64SliceParam reads a float64 slice param
func ReadFloat64SliceParam(name, in, format string, defaultValue []float64, values Gettable, required bool) ([]float64, *errors.Validation) {
	value := values.Get(name)

	if required && value == "" {
		if len(defaultValue) > 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return nil, err
		}
	}

	if value == "" && len(defaultValue) > 0 {
		return defaultValue, nil
	}

	var result []float64
	for _, r := range split(value, format) {
		v, err := util.ConvertFloat64(r)
		if err != nil {
			return nil, errors.InvalidType(name, in, "float64", r)
		}
		result = append(result, v)
	}
	return result, nil
}

// ReadInt8SliceParam reads a int8 slice param
func ReadInt8SliceParam(name, in, format string, defaultValue []int8, values Gettable, required bool) ([]int8, *errors.Validation) {
	value := values.Get(name)

	if required && value == "" {
		if len(defaultValue) > 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return nil, err
		}
	}

	if value == "" && len(defaultValue) > 0 {
		return defaultValue, nil
	}

	var result []int8
	for _, r := range split(value, format) {
		v, err := util.ConvertInt8(r)
		if err != nil {
			return nil, errors.InvalidType(name, in, "int8", r)
		}
		result = append(result, v)
	}
	return result, nil
}

// ReadInt16SliceParam reads a int16 slice param
func ReadInt16SliceParam(name, in, format string, defaultValue []int16, values Gettable, required bool) ([]int16, *errors.Validation) {
	value := values.Get(name)

	if required && value == "" {
		if len(defaultValue) > 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return nil, err
		}
	}

	if value == "" && len(defaultValue) > 0 {
		return defaultValue, nil
	}

	var result []int16
	for _, r := range split(value, format) {
		v, err := util.ConvertInt16(r)
		if err != nil {
			return nil, errors.InvalidType(name, in, "int16", r)
		}
		result = append(result, v)
	}
	return result, nil
}

// ReadInt32SliceParam reads a int32 slice param
func ReadInt32SliceParam(name, in, format string, defaultValue []int32, values Gettable, required bool) ([]int32, *errors.Validation) {
	value := values.Get(name)

	if required && value == "" {
		if len(defaultValue) > 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return nil, err
		}
	}

	if value == "" && len(defaultValue) > 0 {
		return defaultValue, nil
	}

	var result []int32
	for _, r := range split(value, format) {
		v, err := util.ConvertInt32(r)
		if err != nil {
			return nil, errors.InvalidType(name, in, "int32", r)
		}
		result = append(result, v)
	}
	return result, nil
}

// ReadInt64SliceParam reads a int64 slice param
func ReadInt64SliceParam(name, in, format string, defaultValue []int64, values Gettable, required bool) ([]int64, *errors.Validation) {
	value := values.Get(name)

	if required && value == "" {
		if len(defaultValue) > 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return nil, err
		}
	}

	if value == "" && len(defaultValue) > 0 {
		return defaultValue, nil
	}

	var result []int64
	for _, r := range split(value, format) {
		v, err := util.ConvertInt64(r)
		if err != nil {
			return nil, errors.InvalidType(name, in, "int64", r)
		}
		result = append(result, v)
	}
	return result, nil
}

// ReadUint8SliceParam reads a uint8 slice param
func ReadUint8SliceParam(name, in, format string, defaultValue []uint8, values Gettable, required bool) ([]uint8, *errors.Validation) {
	value := values.Get(name)

	if required && value == "" {
		if len(defaultValue) > 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return nil, err
		}
	}

	if value == "" && len(defaultValue) > 0 {
		return defaultValue, nil
	}

	var result []uint8
	for _, r := range split(value, format) {
		v, err := util.ConvertUint8(r)
		if err != nil {
			return nil, errors.InvalidType(name, in, "uint8", r)
		}
		result = append(result, v)
	}
	return result, nil
}

// ReadUint16SliceParam reads a uint16 slice param
func ReadUint16SliceParam(name, in, format string, defaultValue []uint16, values Gettable, required bool) ([]uint16, *errors.Validation) {
	value := values.Get(name)

	if required && value == "" {
		if len(defaultValue) > 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return nil, err
		}
	}

	if value == "" && len(defaultValue) > 0 {
		return defaultValue, nil
	}

	var result []uint16
	for _, r := range split(value, format) {
		v, err := util.ConvertUint16(r)
		if err != nil {
			return nil, errors.InvalidType(name, in, "uint16", r)
		}
		result = append(result, v)
	}
	return result, nil
}

// ReadUint32SliceParam reads a uint32 slice param
func ReadUint32SliceParam(name, in, format string, defaultValue []uint32, values Gettable, required bool) ([]uint32, *errors.Validation) {
	value := values.Get(name)

	if required && value == "" {
		if len(defaultValue) > 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return nil, err
		}
	}

	if value == "" && len(defaultValue) > 0 {
		return defaultValue, nil
	}

	var result []uint32
	for _, r := range split(value, format) {
		v, err := util.ConvertUint32(r)
		if err != nil {
			return nil, errors.InvalidType(name, in, "uint32", r)
		}
		result = append(result, v)
	}
	return result, nil
}

// ReadUint64SliceParam reads a uint64 slice param
func ReadUint64SliceParam(name, in, format string, defaultValue []uint64, values Gettable, required bool) ([]uint64, *errors.Validation) {
	value := values.Get(name)

	if required && value == "" {
		if len(defaultValue) > 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return nil, err
		}
	}

	if value == "" && len(defaultValue) > 0 {
		return defaultValue, nil
	}

	var result []uint64
	for _, r := range split(value, format) {
		v, err := util.ConvertUint64(r)
		if err != nil {
			return nil, errors.InvalidType(name, in, "uint64", r)
		}
		result = append(result, v)
	}
	return result, nil
}

// ReadStringParam reads a string param
func ReadStringParam(name, in, defaultValue string, values Gettable, required bool) (string, *errors.Validation) {
	value := values.Get(name)
	if required && value == "" {
		if defaultValue != "" {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return "", err
		}
	}
	if value == "" && defaultValue != "" {
		value = defaultValue
	}
	return value, nil
}

// ReadBoolParam reads a boolean param value
func ReadBoolParam(name, in string, defaultValue bool, values Gettable, required bool) (bool, *errors.Validation) {
	value := values.Get(name)
	if required && value == "" {
		if defaultValue {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return false, err
		}
	}

	v, err := util.ConvertBool(value)
	if err != nil {
		return false, errors.InvalidType(name, in, "bool", value)
	}
	return v, nil
}

// ReadFloat32Param reads a boolean param value
func ReadFloat32Param(name, in string, defaultValue float32, values Gettable, required bool) (float32, *errors.Validation) {
	value := values.Get(name)
	if required && value == "" {
		if defaultValue != 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return 0, err
		}
	}

	v, err := util.ConvertFloat32(value)
	if err != nil {
		return 0, errors.InvalidType(name, in, "float32", value)
	}
	return v, nil
}

// ReadFloat64Param reads a boolean param value
func ReadFloat64Param(name, in string, defaultValue float64, values Gettable, required bool) (float64, *errors.Validation) {
	value := values.Get(name)
	if required && value == "" {
		if defaultValue != 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return 0, err
		}
	}

	v, err := util.ConvertFloat64(value)
	if err != nil {
		return 0, errors.InvalidType(name, in, "float64", value)
	}
	return v, nil
}

// ReadInt8Param reads a int8 param value
func ReadInt8Param(name, in string, defaultValue int8, values Gettable, required bool) (int8, *errors.Validation) {
	value := values.Get(name)
	if required && value == "" {
		if defaultValue != 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return 0, err
		}
	}

	v, err := util.ConvertInt8(value)
	if err != nil {
		return 0, errors.InvalidType(name, in, "int8", value)
	}
	return v, nil
}

// ReadInt16Param reads a int16 param value
func ReadInt16Param(name, in string, defaultValue int16, values Gettable, required bool) (int16, *errors.Validation) {
	value := values.Get(name)
	if required && value == "" {
		if defaultValue != 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return 0, err
		}
	}
	v, err := util.ConvertInt16(value)
	if err != nil {
		return 0, errors.InvalidType(name, in, "int16", value)
	}
	return v, nil
}

// ReadInt32Param reads a int32 param value
func ReadInt32Param(name, in string, defaultValue int32, values Gettable, required bool) (int32, *errors.Validation) {
	value := values.Get(name)
	if required && value == "" {
		if defaultValue != 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return 0, err
		}
	}

	v, err := util.ConvertInt32(value)
	if err != nil {
		return 0, errors.InvalidType(name, in, "int32", value)
	}
	return v, nil
}

// ReadInt64Param reads a int64 param value
func ReadInt64Param(name, in string, defaultValue int64, values Gettable, required bool) (int64, *errors.Validation) {
	value := values.Get(name)
	if required && value == "" {
		if defaultValue != 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return 0, err
		}
	}

	v, err := util.ConvertInt64(value)
	if err != nil {
		return 0, errors.InvalidType(name, in, "int64", value)
	}
	return v, nil
}

// ReadIntParam reads a int param value
func ReadIntParam(name, in string, defaultValue int, values Gettable, required bool) (int, *errors.Validation) {
	value := values.Get(name)
	if required && value == "" {
		if defaultValue != 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return 0, err
		}
	}

	v, err := util.ConvertInt64(value)
	if err != nil {
		return 0, errors.InvalidType(name, in, "int64", value)
	}
	return int(v), nil
}

// ReadUint8Param reads a uint8 param value
func ReadUint8Param(name, in string, defaultValue uint8, values Gettable, required bool) (uint8, *errors.Validation) {
	value := values.Get(name)
	if required && value == "" {
		if defaultValue != 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return 0, err
		}
	}

	v, err := util.ConvertUint8(value)
	if err != nil {
		return 0, errors.InvalidType(name, in, "uint8", value)
	}
	return v, nil
}

// ReadUint16Param reads a uint16 param value
func ReadUint16Param(name, in string, defaultValue uint16, values Gettable, required bool) (uint16, *errors.Validation) {
	value := values.Get(name)
	if required && value == "" {
		if defaultValue != 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return 0, err
		}
	}

	v, err := util.ConvertUint16(value)
	if err != nil {
		return 0, errors.InvalidType(name, in, "uint16", value)
	}
	return v, nil
}

// ReadUint32Param reads a uint32 param value
func ReadUint32Param(name, in string, defaultValue uint32, values Gettable, required bool) (uint32, *errors.Validation) {
	value := values.Get(name)
	if required && value == "" {
		if defaultValue != 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return 0, err
		}
	}

	v, err := util.ConvertUint32(value)
	if err != nil {
		return 0, errors.InvalidType(name, in, "uint32", value)
	}
	return v, nil
}

// ReadUint64Param reads a uint64 param value
func ReadUint64Param(name, in string, defaultValue uint64, values Gettable, required bool) (uint64, *errors.Validation) {
	value := values.Get(name)
	if required && value == "" {
		if defaultValue != 0 {
			return defaultValue, nil
		}
		if err := validate.RequiredString(name, in, value); err != nil {
			return 0, err
		}
	}

	v, err := util.ConvertUint64(value)
	if err != nil {
		return 0, errors.InvalidType(name, in, "uint64", value)
	}
	return v, nil
}
