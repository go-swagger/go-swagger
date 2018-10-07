package swag

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func assertSingleValue(t *testing.T, inElem, elem reflect.Value, expectPointer bool, idx int) {
	if !assert.Truef(t,
		(elem.Kind() == reflect.Ptr) == expectPointer,
		"Unexpected expectPointer=%t value type", expectPointer) {
		return
	}
	if inElem.Kind() == reflect.Ptr && !inElem.IsNil() {
		inElem = reflect.Indirect(inElem)
	}
	if elem.Kind() == reflect.Ptr && !elem.IsNil() {
		elem = reflect.Indirect(elem)
	}

	if !assert.Truef(t,
		(elem.Kind() == reflect.Ptr && elem.IsNil()) || IsZero(elem.Interface()) ==
			(inElem.Kind() == reflect.Ptr && inElem.IsNil()) || IsZero(inElem.Interface()),
		"Unexpected nil pointer at idx %d", idx) {
		return
	}

	if !((elem.Kind() == reflect.Ptr && elem.IsNil()) || IsZero(elem.Interface())) {
		if !assert.IsTypef(t, inElem.Interface(), elem.Interface(), "Expected in/out to match types") {
			return
		}
		assert.EqualValuesf(t, inElem.Interface(), elem.Interface(), "Unexpected value at idx %d: %v", idx, elem.Interface())
	}
}

// assertValues checks equivalent representation pointer vs values for single var, slices and maps
func assertValues(t *testing.T, in, out interface{}, expectPointer bool, idx int) {
	vin := reflect.ValueOf(in)
	vout := reflect.ValueOf(out)
	switch vin.Kind() {
	case reflect.Slice, reflect.Map:
		if !assert.Equalf(t, vin.Kind(), vout.Kind(), "Unexpected output type at idx %d", idx) ||
			!assert.Equalf(t, vin.Len(), vout.Len(), "Unexpected len at idx %d", idx) {
			break
		}
		var elem, inElem reflect.Value
		for i := 0; i < vin.Len(); i++ {
			if vin.Kind() == reflect.Slice {
				elem = vout.Index(i)
				inElem = vin.Index(i)
			} else if vin.Kind() == reflect.Map {
				keys := vin.MapKeys()
				elem = vout.MapIndex(keys[i])
				inElem = vout.MapIndex(keys[i])
			}
			assertSingleValue(t, inElem, elem, expectPointer, idx)
		}
	default:
		inElem := vin
		elem := vout
		assertSingleValue(t, inElem, elem, expectPointer, idx)
	}
}

var testCasesStringSlice = [][]string{
	{"a", "b", "c", "d", "e"},
	{"a", "b", "", "", "e"},
}

func TestStringSlice(t *testing.T) {
	for idx, in := range testCasesStringSlice {
		if in == nil {
			continue
		}
		out := StringSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := StringValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesStringValueSlice = [][]*string{
	{String("a"), String("b"), nil, String("c")},
}

func TestStringValueSlice(t *testing.T) {
	for idx, in := range testCasesStringValueSlice {
		if in == nil {
			continue
		}
		out := StringValueSlice(in)
		assertValues(t, in, out, false, idx)

		out2 := StringSlice(out)
		assertValues(t, in, out2, true, idx)
	}
}

var testCasesStringMap = []map[string]string{
	{"a": "1", "b": "2", "c": "3"},
}

func TestStringMap(t *testing.T) {
	for idx, in := range testCasesStringMap {
		if in == nil {
			continue
		}
		out := StringMap(in)
		assertValues(t, in, out, true, idx)

		out2 := StringValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesBoolSlice = [][]bool{
	{true, true, false, false},
}

func TestBoolSlice(t *testing.T) {
	for idx, in := range testCasesBoolSlice {
		if in == nil {
			continue
		}
		out := BoolSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := BoolValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesBoolValueSlice = [][]*bool{
	{Bool(true), Bool(true), Bool(false), Bool(false)},
}

func TestBoolValueSlice(t *testing.T) {
	for idx, in := range testCasesBoolValueSlice {
		if in == nil {
			continue
		}
		out := BoolValueSlice(in)
		assertValues(t, in, out, false, idx)

		out2 := BoolSlice(out)
		assertValues(t, in, out2, true, idx)
	}
}

var testCasesBoolMap = []map[string]bool{
	{"a": true, "b": false, "c": true},
}

func TestBoolMap(t *testing.T) {
	for idx, in := range testCasesBoolMap {
		if in == nil {
			continue
		}
		out := BoolMap(in)
		assertValues(t, in, out, true, idx)

		out2 := BoolValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesIntSlice = [][]int{
	{1, 2, 3, 4},
}

func TestIntSlice(t *testing.T) {
	for idx, in := range testCasesIntSlice {
		if in == nil {
			continue
		}
		out := IntSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := IntValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesIntValueSlice = [][]*int{
	{Int(1), Int(2), Int(3), Int(4)},
}

func TestIntValueSlice(t *testing.T) {
	for idx, in := range testCasesIntValueSlice {
		if in == nil {
			continue
		}
		out := IntValueSlice(in)
		assertValues(t, in, out, false, idx)

		out2 := IntSlice(out)
		assertValues(t, in, out2, true, idx)
	}
}

var testCasesIntMap = []map[string]int{
	{"a": 3, "b": 2, "c": 1},
}

func TestIntMap(t *testing.T) {
	for idx, in := range testCasesIntMap {
		if in == nil {
			continue
		}
		out := IntMap(in)
		assertValues(t, in, out, true, idx)

		out2 := IntValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesInt64Slice = [][]int64{
	{1, 2, 3, 4},
}

func TestInt64Slice(t *testing.T) {
	for idx, in := range testCasesInt64Slice {
		if in == nil {
			continue
		}
		out := Int64Slice(in)
		assertValues(t, in, out, true, idx)

		out2 := Int64ValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesInt64ValueSlice = [][]*int64{
	{Int64(1), Int64(2), Int64(3), Int64(4)},
}

func TestInt64ValueSlice(t *testing.T) {
	for idx, in := range testCasesInt64ValueSlice {
		if in == nil {
			continue
		}
		out := Int64ValueSlice(in)
		assertValues(t, in, out, false, idx)

		out2 := Int64Slice(out)
		assertValues(t, in, out2, true, idx)
	}
}

var testCasesInt64Map = []map[string]int64{
	{"a": 3, "b": 2, "c": 1},
}

func TestInt64Map(t *testing.T) {
	for idx, in := range testCasesInt64Map {
		if in == nil {
			continue
		}
		out := Int64Map(in)
		assertValues(t, in, out, true, idx)

		out2 := Int64ValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesFloat64Slice = [][]float64{
	{1, 2, 3, 4},
}

func TestFloat64Slice(t *testing.T) {
	for idx, in := range testCasesFloat64Slice {
		if in == nil {
			continue
		}
		out := Float64Slice(in)
		assertValues(t, in, out, true, idx)

		out2 := Float64ValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesUintSlice = [][]uint{
	{1, 2, 3, 4},
}

func TestUintSlice(t *testing.T) {
	for idx, in := range testCasesUintSlice {
		if in == nil {
			continue
		}
		out := UintSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := UintValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesUintValueSlice = [][]*uint{}

func TestUintValueSlice(t *testing.T) {
	for idx, in := range testCasesUintValueSlice {
		if in == nil {
			continue
		}
		out := UintValueSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := UintSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesUintMap = []map[string]uint{
	{"a": 3, "b": 2, "c": 1},
}

func TestUintMap(t *testing.T) {
	for idx, in := range testCasesUintMap {
		if in == nil {
			continue
		}
		out := UintMap(in)
		assertValues(t, in, out, true, idx)

		out2 := UintValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesUint64Slice = [][]uint64{
	{1, 2, 3, 4},
}

func TestUint64Slice(t *testing.T) {
	for idx, in := range testCasesUint64Slice {
		if in == nil {
			continue
		}
		out := Uint64Slice(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint64ValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesUint64ValueSlice = [][]*uint64{}

func TestUint64ValueSlice(t *testing.T) {
	for idx, in := range testCasesUint64ValueSlice {
		if in == nil {
			continue
		}
		out := Uint64ValueSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint64Slice(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesUint64Map = []map[string]uint64{
	{"a": 3, "b": 2, "c": 1},
}

func TestUint64Map(t *testing.T) {
	for idx, in := range testCasesUint64Map {
		if in == nil {
			continue
		}
		out := Uint64Map(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint64ValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesFloat64ValueSlice = [][]*float64{}

func TestFloat64ValueSlice(t *testing.T) {
	for idx, in := range testCasesFloat64ValueSlice {
		if in == nil {
			continue
		}
		out := Float64ValueSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := Float64Slice(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesFloat64Map = []map[string]float64{
	{"a": 3, "b": 2, "c": 1},
}

func TestFloat64Map(t *testing.T) {
	for idx, in := range testCasesFloat64Map {
		if in == nil {
			continue
		}
		out := Float64Map(in)
		assertValues(t, in, out, true, idx)

		out2 := Float64ValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesTimeSlice = [][]time.Time{
	{time.Now(), time.Now().AddDate(100, 0, 0)},
}

func TestTimeSlice(t *testing.T) {
	for idx, in := range testCasesTimeSlice {
		if in == nil {
			continue
		}
		out := TimeSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := TimeValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesTimeValueSlice = [][]*time.Time{
	{Time(time.Now()), Time(time.Now().AddDate(100, 0, 0))},
}

func TestTimeValueSlice(t *testing.T) {
	for idx, in := range testCasesTimeValueSlice {
		if in == nil {
			continue
		}
		out := TimeValueSlice(in)
		assertValues(t, in, out, false, idx)

		out2 := TimeSlice(out)
		assertValues(t, in, out2, true, idx)
	}
}

var testCasesTimeMap = []map[string]time.Time{
	{"a": time.Now().AddDate(-100, 0, 0), "b": time.Now()},
}

func TestTimeMap(t *testing.T) {
	for idx, in := range testCasesTimeMap {
		if in == nil {
			continue
		}
		out := TimeMap(in)
		assertValues(t, in, out, true, idx)

		out2 := TimeValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesInt32Slice = [][]int32{
	{1, 2, 3, 4},
}

func TestInt32Slice(t *testing.T) {
	for idx, in := range testCasesInt32Slice {
		if in == nil {
			continue
		}
		out := Int32Slice(in)
		assertValues(t, in, out, true, idx)

		out2 := Int32ValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesInt32ValueSlice = [][]*int32{
	{Int32(1), Int32(2), Int32(3), Int32(4)},
}

func TestInt32ValueSlice(t *testing.T) {
	for idx, in := range testCasesInt32ValueSlice {
		if in == nil {
			continue
		}
		out := Int32ValueSlice(in)
		assertValues(t, in, out, false, idx)

		out2 := Int32Slice(out)
		assertValues(t, in, out2, true, idx)
	}
}

var testCasesInt32Map = []map[string]int32{
	{"a": 3, "b": 2, "c": 1},
}

func TestInt32Map(t *testing.T) {
	for idx, in := range testCasesInt32Map {
		if in == nil {
			continue
		}
		out := Int32Map(in)
		assertValues(t, in, out, true, idx)

		out2 := Int32ValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesUint32Slice = [][]uint32{
	{1, 2, 3, 4},
}

func TestUint32Slice(t *testing.T) {
	for idx, in := range testCasesUint32Slice {
		if in == nil {
			continue
		}
		out := Uint32Slice(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint32ValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesUint32ValueSlice = [][]*uint32{
	{Uint32(1), Uint32(2), Uint32(3), Uint32(4)},
}

func TestUint32ValueSlice(t *testing.T) {
	for idx, in := range testCasesUint32ValueSlice {
		if in == nil {
			continue
		}
		out := Uint32ValueSlice(in)
		assertValues(t, in, out, false, idx)

		out2 := Uint32Slice(out)
		assertValues(t, in, out2, true, idx)
	}
}

var testCasesUint32Map = []map[string]uint32{
	{"a": 3, "b": 2, "c": 1},
}

func TestUint32Map(t *testing.T) {
	for idx, in := range testCasesUint32Map {
		if in == nil {
			continue
		}
		out := Uint32Map(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint32ValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

var testCasesString = []string{"a", "b", "c", "d", "e", ""}

func TestStringValue(t *testing.T) {
	for idx, in := range testCasesString {
		out := String(in)
		assertValues(t, in, out, true, idx)

		out2 := StringValue(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, StringValue(nil), "expected conversion from nil to return zero value")
}

var testCasesBool = []bool{true, false}

func TestBoolValue(t *testing.T) {
	for idx, in := range testCasesBool {
		out := Bool(in)
		assertValues(t, in, out, true, idx)

		out2 := BoolValue(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, BoolValue(nil), "expected conversion from nil to return zero value")
}

var testCasesInt = []int{1, 2, 3, 0}

func TestIntValue(t *testing.T) {
	for idx, in := range testCasesInt {
		out := Int(in)
		assertValues(t, in, out, true, idx)

		out2 := IntValue(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, IntValue(nil), "expected conversion from nil to return zero value")
}

var testCasesInt32 = []int32{1, 2, 3, 0}

func TestInt32Value(t *testing.T) {
	for idx, in := range testCasesInt32 {
		out := Int32(in)
		assertValues(t, in, out, true, idx)

		out2 := Int32Value(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, Int32Value(nil), "expected conversion from nil to return zero value")
}

var testCasesInt64 = []int64{1, 2, 3, 0}

func TestInt64Value(t *testing.T) {
	for idx, in := range testCasesInt64 {
		out := Int64(in)
		assertValues(t, in, out, true, idx)

		out2 := Int64Value(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, Int64Value(nil), "expected conversion from nil to return zero value")
}

var testCasesUint = []uint{1, 2, 3, 0}

func TestUintValue(t *testing.T) {
	for idx, in := range testCasesUint {
		out := Uint(in)
		assertValues(t, in, out, true, idx)

		out2 := UintValue(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, UintValue(nil), "expected conversion from nil to return zero value")
}

var testCasesUint32 = []uint32{1, 2, 3, 0}

func TestUint32Value(t *testing.T) {
	for idx, in := range testCasesUint32 {
		out := Uint32(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint32Value(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, Uint32Value(nil), "expected conversion from nil to return zero value")
}

var testCasesUint64 = []uint64{1, 2, 3, 0}

func TestUint64Value(t *testing.T) {
	for idx, in := range testCasesUint64 {
		out := Uint64(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint64Value(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, Uint64Value(nil), "expected conversion from nil to return zero value")
}

var testCasesFloat64 = []float64{1, 2, 3, 0}

func TestFloat64Value(t *testing.T) {
	for idx, in := range testCasesFloat64 {
		out := Float64(in)
		assertValues(t, in, out, true, idx)

		out2 := Float64Value(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, Float64Value(nil), "expected conversion from nil to return zero value")
}

var testCasesTime = []time.Time{
	time.Now().AddDate(-100, 0, 0), time.Now(),
}

func TestTimeValue(t *testing.T) {
	for idx, in := range testCasesTime {
		out := Time(in)
		assertValues(t, in, out, true, idx)

		out2 := TimeValue(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, TimeValue(nil), "expected conversion from nil to return zero value")
}
