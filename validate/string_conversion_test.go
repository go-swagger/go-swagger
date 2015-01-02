package validate

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/casualjim/go-swagger"
	"github.com/stretchr/testify/assert"
)

type unmarshallerSlice []string

func (u *unmarshallerSlice) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return errors.New("an error")
	}
	*u = strings.Split(string(data), ",")
	return nil
}

type SomeOperationParams struct {
	Name        string
	ID          int64
	Confirmed   bool
	Age         int
	Visits      int32
	Count       int16
	Seq         int8
	UID         uint64
	UAge        uint
	UVisits     uint32
	UCount      uint16
	USeq        uint8
	Score       float32
	Rate        float64
	Timestamp   swagger.DateTime
	Birthdate   swagger.Date
	LastFailure *swagger.DateTime
	Unsupported struct{}
	Tags        []string
	Prefs       []int32
	Categories  unmarshallerSlice
}

func FloatParamTest(t *testing.T, fName, pName, format string, val reflect.Value, defVal, expectedDef interface{}, actual func() interface{}) {
	fld := val.FieldByName(pName)

	err := setFieldValue(fld, "5", defVal)
	assert.NoError(t, err)
	assert.Equal(t, 5, actual())

	err = setFieldValue(fld, "", defVal)
	assert.NoError(t, err)
	assert.Equal(t, expectedDef, actual())

	err = setFieldValue(fld, "yada", defVal)
	assert.Error(t, err)
}

func IntParamTest(t *testing.T, pName string, val reflect.Value, defVal, expectedDef interface{}, actual func() interface{}) {
	fld := val.FieldByName(pName)

	err := setFieldValue(fld, "5", defVal)
	assert.NoError(t, err)
	assert.Equal(t, 5, actual())

	err = setFieldValue(fld, "", defVal)
	assert.NoError(t, err)
	assert.Equal(t, expectedDef, actual())

	err = setFieldValue(fld, "yada", defVal)
	assert.Error(t, err)
}

func TestParamBinding(t *testing.T) {

	actual := new(SomeOperationParams)
	val := reflect.ValueOf(actual).Elem()
	fld := val.FieldByName("Name")

	err := setFieldValue(fld, "the name value", "some-name")
	assert.NoError(t, err)
	assert.Equal(t, "the name value", actual.Name)

	err = setFieldValue(fld, "", "some-name")
	assert.NoError(t, err)
	assert.Equal(t, "some-name", actual.Name)

	IntParamTest(t, "ID", val, 1, 1, func() interface{} { return actual.ID })
	IntParamTest(t, "ID", val, nil, 0, func() interface{} { return actual.ID })
	IntParamTest(t, "Age", val, 1, 1, func() interface{} { return actual.Age })
	IntParamTest(t, "Age", val, nil, 0, func() interface{} { return actual.Age })
	IntParamTest(t, "Visits", val, 1, 1, func() interface{} { return actual.Visits })
	IntParamTest(t, "Visits", val, nil, 0, func() interface{} { return actual.Visits })
	IntParamTest(t, "Count", val, 1, 1, func() interface{} { return actual.Count })
	IntParamTest(t, "Count", val, nil, 0, func() interface{} { return actual.Count })
	IntParamTest(t, "Seq", val, 1, 1, func() interface{} { return actual.Seq })
	IntParamTest(t, "Seq", val, nil, 0, func() interface{} { return actual.Seq })
	IntParamTest(t, "UID", val, uint64(1), 1, func() interface{} { return actual.UID })
	IntParamTest(t, "UID", val, uint64(0), 0, func() interface{} { return actual.UID })
	IntParamTest(t, "UAge", val, uint(1), 1, func() interface{} { return actual.UAge })
	IntParamTest(t, "UAge", val, nil, 0, func() interface{} { return actual.UAge })
	IntParamTest(t, "UVisits", val, uint32(1), 1, func() interface{} { return actual.UVisits })
	IntParamTest(t, "UVisits", val, nil, 0, func() interface{} { return actual.UVisits })
	IntParamTest(t, "UCount", val, uint16(1), 1, func() interface{} { return actual.UCount })
	IntParamTest(t, "UCount", val, nil, 0, func() interface{} { return actual.UCount })
	IntParamTest(t, "USeq", val, uint8(1), 1, func() interface{} { return actual.USeq })
	IntParamTest(t, "USeq", val, nil, 0, func() interface{} { return actual.USeq })

	FloatParamTest(t, "score", "Score", "float", val, 1.0, 1, func() interface{} { return actual.Score })
	FloatParamTest(t, "score", "Score", "float", val, nil, 0, func() interface{} { return actual.Score })
	FloatParamTest(t, "rate", "Rate", "double", val, 1.0, 1, func() interface{} { return actual.Rate })
	FloatParamTest(t, "rate", "Rate", "double", val, nil, 0, func() interface{} { return actual.Rate })

	confirmedField := val.FieldByName("Confirmed")

	for _, tv := range evaluatesAsTrue {
		err = setFieldValue(confirmedField, tv, true)
		assert.NoError(t, err)
		assert.True(t, actual.Confirmed)
	}

	err = setFieldValue(confirmedField, "", true)
	assert.NoError(t, err)
	assert.True(t, actual.Confirmed)

	err = setFieldValue(confirmedField, "0", nil)
	assert.NoError(t, err)
	assert.False(t, actual.Confirmed)

	dt := swagger.DateTime{Time: time.Date(2014, 3, 19, 2, 9, 0, 0, time.UTC)}
	exp := swagger.DateTime{Time: time.Date(2014, 5, 14, 2, 9, 0, 0, time.UTC)}
	timeField := val.FieldByName("Timestamp")

	err = setFieldValue(timeField, exp.String(), dt)
	assert.NoError(t, err)
	assert.Equal(t, exp, actual.Timestamp)

	err = setFieldValue(timeField, "", dt)
	assert.NoError(t, err)
	assert.Equal(t, dt, actual.Timestamp)

	err = setFieldValue(timeField, "yada", dt)
	assert.Error(t, err)

	ddt := swagger.Date{Time: time.Date(2014, 3, 19, 0, 0, 0, 0, time.UTC)}
	expd := swagger.Date{Time: time.Date(2014, 5, 14, 0, 0, 0, 0, time.UTC)}
	dateField := val.FieldByName("Birthdate")

	err = setFieldValue(dateField, expd.String(), ddt)
	assert.NoError(t, err)
	assert.Equal(t, expd, actual.Birthdate)

	err = setFieldValue(dateField, "", ddt)
	assert.NoError(t, err)
	assert.Equal(t, ddt, actual.Birthdate)

	err = setFieldValue(dateField, "yada", ddt)
	assert.Error(t, err)

	fdt := &swagger.DateTime{Time: time.Date(2014, 3, 19, 2, 9, 0, 0, time.UTC)}
	fexp := &swagger.DateTime{Time: time.Date(2014, 5, 14, 2, 9, 0, 0, time.UTC)}
	ftimeField := val.FieldByName("LastFailure")

	err = setFieldValue(ftimeField, fexp.String(), fdt)
	assert.NoError(t, err)
	assert.Equal(t, fexp, actual.LastFailure)

	err = setFieldValue(ftimeField, "", fdt)
	assert.NoError(t, err)
	assert.Equal(t, fdt, actual.LastFailure)

	err = setFieldValue(ftimeField, "", dt)
	assert.NoError(t, err)
	assert.Equal(t, &dt, actual.LastFailure)

	actual.LastFailure = nil
	err = setFieldValue(ftimeField, "yada", fdt)
	assert.Error(t, err)
	assert.Nil(t, actual.LastFailure)

	unsupportedField := val.FieldByName("Unsupported")
	err = setFieldValue(unsupportedField, "", nil)
	assert.Error(t, err)
}

func TestSliceConversion(t *testing.T) {

	actual := new(SomeOperationParams)
	val := reflect.ValueOf(actual).Elem()

	prefsField := val.FieldByName("Prefs")
	cData := "yada,2,3"
	err := setFormattedSliceFieldValue(prefsField, cData, "csv", nil)
	assert.Error(t, err)

	sliced := []string{"some", "string", "values"}
	seps := map[string]string{"ssv": " ", "tsv": "\t", "pipes": "|", "csv": ",", "": ","}

	tagsField := val.FieldByName("Tags")
	for k, sep := range seps {
		actual.Tags = nil
		cData := strings.Join(sliced, sep)
		err := setFormattedSliceFieldValue(tagsField, cData, k, nil)
		assert.NoError(t, err)
		assert.Equal(t, sliced, actual.Tags)
		cData = strings.Join(sliced, " "+sep+" ")
		err = setFormattedSliceFieldValue(tagsField, cData, k, nil)
		assert.NoError(t, err)
		assert.Equal(t, sliced, actual.Tags)
		err = setFormattedSliceFieldValue(tagsField, "", k, sliced)
		assert.NoError(t, err)
		assert.Equal(t, sliced, actual.Tags)
	}

	assert.Nil(t, split("yada", "multi"))
	assert.Nil(t, split("", ""))

	categoriesField := val.FieldByName("Categories")
	cData = strings.Join(sliced, ",")
	err = setFormattedSliceFieldValue(categoriesField, cData, "csv", nil)
	assert.NoError(t, err)
	assert.Equal(t, sliced, actual.Categories)
	err = setFormattedSliceFieldValue(categoriesField, "", "csv", nil)
	assert.Error(t, err)

}
