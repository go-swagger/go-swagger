package diff

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func Test_getRef(t *testing.T) {
	type args struct {
		item interface{}
	}
	aRef, _ := spec.NewRef("hello")
	tests := []struct {
		name string
		args args
		want spec.Ref
	}{
		{name: "rando object",
			args: args{item: "bob"},
			want: spec.Ref{},
		},
		{name: "refable",
			args: args{&spec.Refable{Ref: aRef}},
			want: aRef,
		},
		{name: "schema",
			args: args{&spec.Schema{SchemaProps: spec.SchemaProps{Ref: aRef}}},
			want: aRef,
		},
		{name: "schemaProps",
			args: args{&spec.SchemaProps{Ref: aRef}},
			want: aRef,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRef(tt.args.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRef() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestCheckToFromArrayType(t *testing.T) {
// 	type args struct {
// 		diffs []TypeDiff
// 		type1 interface{}
// 		type2 interface{}
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want []TypeDiff
// 	}{
// 		{
// 			name: "to",
// 			args: args{
// 				type1: spec.Int32Property(),
// 				type2: arraySchemaOf("string"),
// 			},
// 			want: []TypeDiff{{Change: ChangedType, FromType: "<integer>", ToType: "<array[string]>"}},
// 		},
// 		{
// 			name: "from",
// 			args: args{
// 				type1: arraySchemaOf("string"),
// 				type2: spec.Int32Property(),
// 			},
// 			want: []TypeDiff{{Change: ChangedType, ToType: "<integer>", FromType: "<array[string]>"}},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := CheckToFromArrayType(tt.args.diffs, tt.args.type1, tt.args.type2); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("CheckToFromArrayType() = %s, want %s", jsonStr(got), jsonStr(tt.want))
// 			}
// 		})
// 	}
// }

/*
func arraySchemaOf(typename string) *spec.Schema {
	return &spec.Schema{SchemaProps: spec.SchemaProps{
		Type: spec.StringOrArray{"array"},
		Items: &spec.SchemaOrArray{
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{typename}}}},
	},
	}
}
*/

func TestCheckToFromPrimitiveType(t *testing.T) {
	type args struct {
		diffs []TypeDiff
		type1 interface{}
		type2 interface{}
	}
	tests := []struct {
		name string
		args args
		want []TypeDiff
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckToFromPrimitiveType(tt.args.diffs, tt.args.type1, tt.args.type2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckToFromPrimitiveType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckRefChange(t *testing.T) {
	type args struct {
		diffs []TypeDiff
		type1 interface{}
		type2 interface{}
	}
	tests := []struct {
		name           string
		args           args
		wantDiffReturn []TypeDiff
	}{
		{
			name: "reftarget",
			args: args{
				type1: spec.RefProperty("#/definitions/FirstObject"),
				type2: spec.RefProperty("#/definitions/SecondObject"),
			},
			wantDiffReturn: []TypeDiff{{Change: RefTargetChanged, FromType: "<FirstObject>", ToType: "<SecondObject>"}},
		},
		{
			name: "toref",
			args: args{
				type1: spec.Int32Property(),
				type2: spec.RefProperty("#/definitions/SecondObject"),
			},
			wantDiffReturn: []TypeDiff{{Change: ChangedType, FromType: "<integer>", ToType: "<SecondObject>"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDiffReturn := CheckRefChange(tt.args.diffs, tt.args.type1, tt.args.type2); !reflect.DeepEqual(gotDiffReturn, tt.wantDiffReturn) {
				t.Errorf("CheckRefChange() = %s, want %s", jsonStr(gotDiffReturn), jsonStr(tt.wantDiffReturn))
			}
		})
	}
}

func Test_isRef(t *testing.T) {
	r := spec.RefSchema("#/definitions/Bob")
	p := spec.Int16Property()

	assert.True(t, isRefType(r))
	assert.False(t, isRefType(p))

	refb := spec.Refable{Ref: r.Ref}
	assert.True(t, isRefType(refb))

	ss := spec.SimpleSchema{}
	assert.False(t, isRefType(&ss))

	ro := time.Timer{}
	assert.False(t, isRefType(ro))

}

func Test_compareEnums(t *testing.T) {
	type args struct {
		left  []interface{}
		right []interface{}
	}
	tests := []struct {
		name string
		args args
		want []TypeDiff
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareEnums(tt.args.left, tt.args.right); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compareEnums() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkNumericTypeChanges(t *testing.T) {
	type args struct {
		diffs []TypeDiff
		type1 *spec.SchemaProps
		type2 *spec.SchemaProps
	}
	tests := []struct {
		name string
		args args
		want []TypeDiff
	}{
		{
			name: "ExclusiveMin",
			args: args{
				type1: &spec.Int32Property().WithMinimum(100, true).SchemaProps,
				type2: &spec.Int32Property().SchemaProps,
			},
			want: []TypeDiff{{Change: WidenedType, Description: "Exclusive Minimum Removed:false->false"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkNumericTypeChanges(tt.args.diffs, tt.args.type1, tt.args.type2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkNumericTypeChanges() = %s, want %s", jsonStr(got), jsonStr(tt.want))
			}
		})
	}
}

func TestCompareFloatValues(t *testing.T) {
	type args struct {
		name   string
		field1 *float64
		field2 *float64
	}
	tests := []struct {
		name string
		args args
		want []TypeDiff
	}{
		{
			name: "both null",
			args: args{name: "bob", field1: nil, field2: nil},
			want: []TypeDiff{},
		},
		{
			name: "greater",
			args: args{name: "bob", field1: floatPointerOf(1.0), field2: floatPointerOf(2.0)},
			want: []TypeDiff{{Change: WidenedType, Description: "bob 1.000000->2.000000"}},
		},
		{
			name: "less",
			args: args{name: "bob", field1: floatPointerOf(2.0), field2: floatPointerOf(1.0)},
			want: []TypeDiff{{Change: NarrowedType, Description: "bob 2.000000->1.000000"}},
		},
		{
			name: "firstNil",
			args: args{name: "bob", field1: nil, field2: floatPointerOf(1.0)},
			want: []TypeDiff{{Change: AddedConstraint, Description: "bob(1.000000)"}},
		},
		{
			name: "secondNil",
			args: args{name: "bob", field1: floatPointerOf(2.0), field2: nil},
			want: []TypeDiff{{Change: DeletedConstraint, Description: "bob(2.000000)"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareFloatValues(tt.args.name, tt.args.field1, tt.args.field2, WidenedType, NarrowedType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckStringTypeChanges() = %s, want %s", jsonStr(got), jsonStr(tt.want))
			}
		})
	}
}

func TestCompareIntValues(t *testing.T) {
	type args struct {
		name   string
		field1 *int64
		field2 *int64
	}
	tests := []struct {
		name string
		args args
		want []TypeDiff
	}{
		{
			name: "both null",
			args: args{name: "bob", field1: nil, field2: nil},
			want: []TypeDiff{},
		},
		{
			name: "greater",
			args: args{name: "bob", field1: intPointerOf(1), field2: intPointerOf(2)},
			want: []TypeDiff{{Change: WidenedType, Description: "bob 1->2"}},
		},
		{
			name: "less",
			args: args{name: "bob", field1: intPointerOf(2), field2: intPointerOf(1)},
			want: []TypeDiff{{Change: NarrowedType, Description: "bob 2->1"}},
		},
		{
			name: "firstNil",
			args: args{name: "bob", field1: nil, field2: intPointerOf(1)},
			want: []TypeDiff{{Change: AddedConstraint, Description: "bob(1)"}},
		},
		{
			name: "secondNil",
			args: args{name: "bob", field1: intPointerOf(2), field2: nil},
			want: []TypeDiff{{Change: DeletedConstraint, Description: "bob(2)"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareIntValues(tt.args.name, tt.args.field1, tt.args.field2, WidenedType, NarrowedType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckStringTypeChanges() = %s, want %s", jsonStr(got), jsonStr(tt.want))
			}
		})
	}
}

func floatPointerOf(f float64) *float64 {
	return &f
}

func intPointerOf(f int64) *int64 {
	return &f
}

func TestCheckToFromRequired(t *testing.T) {
	type args struct {
		required1 bool
		required2 bool
	}
	tests := []struct {
		name      string
		args      args
		wantDiffs []TypeDiff
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDiffs := CheckToFromRequired(tt.args.required1, tt.args.required2); !reflect.DeepEqual(gotDiffs, tt.wantDiffs) {
				t.Errorf("CheckToFromRequired() = %v, want %v", gotDiffs, tt.wantDiffs)
			}
		})
	}
}

func Test_compareProperties(t *testing.T) {
	type args struct {
		location  DifferenceLocation
		schema1   *spec.Schema
		schema2   *spec.Schema
		getRefFn1 SchemaFromRefFn
		getRefFn2 SchemaFromRefFn
		cmp       CompareSchemaFn
	}
	tests := []struct {
		name string
		args args
		want []SpecDifference
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareProperties(tt.args.location, tt.args.schema1, tt.args.schema2, tt.args.getRefFn1, tt.args.getRefFn2, tt.args.cmp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compareProperties() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_propertiesFor(t *testing.T) {
	type args struct {
		schema   *spec.Schema
		getRefFn SchemaFromRefFn
	}
	tests := []struct {
		name string
		args args
		want PropertyMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := propertiesFor(tt.args.schema, tt.args.getRefFn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("propertiesFor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func jsonStr(thing interface{}) string {
	bstr, _ := JSONMarshal(thing)
	return string(bstr)
}
