// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

func initFixture2494() {
	f := newModelFixture("../fixtures/bugs/2494/fixture-2494.yaml", "map of nullable array")
	flattenRun := f.AddRun(false).WithMinimalFlatten(true)

	flattenRun.AddExpectations("port_map.go", []string{
		`type PortMap map[string][]PortBinding`,
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("port_list.go", []string{
		`type PortList [][]*PortBinding`,
	}, todo, noLines, noLines)
}

func initFixture2444() {
	f := newModelFixture("../fixtures/enhancements/2444/fixture-2244.yaml", "min/maxProperties")
	flattenRun := f.AddRun(false).WithMinimalFlatten(true)

	flattenRun.AddExpectations("all_of_with_min_max_properties.go", []string{
		`func (m *AllOfWithMinMaxProperties) Validate(formats strfmt.Registry) error {`,
		`	if err := m.AllOfWithMinMaxPropertiesAO0P0.Validate(formats); err != nil {`,
		`func (m *AllOfWithMinMaxPropertiesAO0P0) Validate(formats strfmt.Registry) error {`,
		`	if m == nil {`,
		`		return errors.TooFewProperties("", "body", 3)`,
		`	props := make(map[string]json.RawMessage, 1+10)`,
		`	j, err := swag.WriteJSON(m)`,
		`	if err = swag.ReadJSON(j, &props); err != nil {`,
		`	nprops := len(props)`,
		`	if nprops < 3 {`,
		`		return errors.TooFewProperties("", "body", 3)`,
		`	if nprops > 5 {`,
		`		return errors.TooManyProperties("", "body", 5)`,
		`	if err := m.validateUID(formats); err != nil {`,
		`	for k := range m.AllOfWithMinMaxPropertiesAO0P0 {`,
		`		if err := validate.MaximumUint(k, "body", uint64(m.AllOfWithMinMaxPropertiesAO0P0[k]), 100, false); err != nil {`,
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("array_items_with_min_max_properties.go", []string{
		`type ArrayItemsWithMinMaxProperties []map[string]interface{}`,
		`func (m ArrayItemsWithMinMaxProperties) Validate(formats strfmt.Registry) error {`,
		`	for i := 0; i < len(m); i++ {`,
		`		nprops := len(m[i])`,
		`		if nprops < 3 {`,
		`			return errors.TooFewProperties(strconv.Itoa(i), "body", 3)`,
		`		if nprops > 5 {`,
		`			return errors.TooManyProperties(strconv.Itoa(i), "body", 5)`,
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("has_max_properties.go", []string{
		`	HasMaxPropertiesAdditionalProperties map[string]interface{}`,
		`	props := make(map[string]json.RawMessage, 1+10)`,
		`	j, err := swag.WriteJSON(m)`,
		`	if err = swag.ReadJSON(j, &props); err != nil {`,
		`	nprops := len(props)`,
		`	if nprops > 2 {`,
		`		return errors.TooManyProperties("", "body", 2)`,
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("has_min_max_properties.go", []string{
		`	HasMinMaxPropertiesAdditionalProperties map[string]interface{}`,
		`	if m == nil {`,
		`		return errors.TooFewProperties("", "body", 3)`,
		`	props := make(map[string]json.RawMessage, 1+10)`,
		`	j, err := swag.WriteJSON(m)`,
		`	if err = swag.ReadJSON(j, &props); err != nil {`,
		`	nprops := len(props)`,
		`	if nprops < 3 {`,
		`		return errors.TooFewProperties("", "body", 3)`,
		`	if nprops > 5 {`,
		`		return errors.TooManyProperties("", "body", 5)`,
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("has_min_properties.go", []string{
		`HasMinPropertiesAdditionalProperties map[string]interface{}`,
		`	if m == nil {`,
		`		return errors.TooFewProperties("", "body", 2)`,
		`	props := make(map[string]json.RawMessage, 1+10)`,
		`	j, err := swag.WriteJSON(m)`,
		`	if err = swag.ReadJSON(j, &props); err != nil {`,
		`	nprops := len(props)`,
		`	if nprops < 2 {`,
		`		return errors.TooFewProperties("", "body", 2)`,
		`	if err := m.validateA(formats); err != nil {`,
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("map_of_arrays_with_min_max_properties.go", []string{
		`type MapOfArraysWithMinMaxProperties map[string][]HasMaxProperties`,
		`	nprops := len(m)`,
		`	if nprops < 3 {`,
		`		return errors.TooFewProperties("", "body", 3)`,
		`	if nprops > 5 {`,
		`		return errors.TooManyProperties("", "body", 5)`,
		`	for k := range m {`,
		`		if err := validate.Required(k, "body", m[k]); err != nil {`,
		`		for i := 0; i < len(m[k]); i++ {`,
		`			if err := m[k][i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName(k + "." + strconv.Itoa(i))`,
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("map_of_integers_with_min_max_properties.go", []string{
		`type MapOfIntegersWithMinMaxProperties map[string]int64`,
		`	nprops := len(m)`,
		`	if nprops < 3 {`,
		`		return errors.TooFewProperties("", "body", 3)`,
		`	if nprops > 5 {`,
		`		return errors.TooManyProperties("", "body", 5)`,
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("map_of_objects_with_min_max_properties.go", []string{
		`type MapOfObjectsWithMinMaxProperties map[string]HasMaxProperties`,
		`	nprops := len(m)`,
		`	if nprops < 3 {`,
		`		return errors.TooFewProperties("", "body", 3)`,
		`	if nprops > 5 {`,
		`		return errors.TooManyProperties("", "body", 5)`,
		`	for k := range m {`,
		`		if err := validate.Required(k, "body", m[k]); err != nil {`,
		`		if val, ok := m[k]; ok {`,
		`			if err := val.Validate(formats); err != nil {`,
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("map_with_min_max_properties.go", []string{
		`type MapWithMinMaxProperties map[string]interface{}`,
		`	nprops := len(m)`,
		`	if nprops < 3 {`,
		`		return errors.TooFewProperties("", "body", 3)`,
		`	if nprops > 5 {`,
		`		return errors.TooManyProperties("", "body", 5)`,
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("object_with_min_max_properties.go", []string{
		`	ObjectWithMinMaxProperties map[string]*HasMaxProperties`,
		`	if m == nil {`,
		`		return errors.TooFewProperties("", "body", 3)`,
		`	props := make(map[string]json.RawMessage, 2+10)`,
		`	j, err := swag.WriteJSON(m)`,
		`	if err = swag.ReadJSON(j, &props); err != nil {`,
		`	nprops := len(props)`,
		`	if nprops < 3 {`,
		`		return errors.TooFewProperties("", "body", 3)`,
		`	if nprops > 5 {`,
		`		return errors.TooManyProperties("", "body", 5)`,
		`	if err := m.validateB(formats); err != nil {`,
		`	if err := m.validateID(formats); err != nil {`,
		`	for k := range m.ObjectWithMinMaxProperties {`,
		`		if err := validate.Required(k, "body", m.ObjectWithMinMaxProperties[k]); err != nil {`,
		`		if val, ok := m.ObjectWithMinMaxProperties[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("untyped_with_min_max_properties.go", []string{
		`type UntypedWithMinMaxProperties map[string]interface{}`,
		`	nprops := len(m)`,
		`	if nprops < 3 {`,
		`		return errors.TooFewProperties("", "body", 3)`,
		`	if nprops > 5 {`,
		`		return errors.TooManyProperties("", "body", 5)`,
	}, todo, noLines, noLines)
}

func initFixtureGuardFormats() {
	f := newModelFixture("../fixtures/enhancements/guard-formats/fixture-guard-formats.yaml", "guard format validations")
	flattenRun := f.AddRun(false).WithMinimalFlatten(true)

	flattenRun.AddExpectations("aliased_date.go", []string{
		`type AliasedDate strfmt.Date`,
		`func (m AliasedDate) Validate(formats strfmt.Registry) error {`,
		`if err := validate.MinLength("", "body", strfmt.Date(m).String(), 16); err != nil {`,
		`if err := validate.MaxLength("", "body", strfmt.Date(m).String(), 20); err != nil {`,
		"if err := validate.Pattern(\"\", \"body\", strfmt.Date(m).String(), `(\\d+)/(\\d+)`); err != nil {",
		`if err := validate.FormatOf("", "body", "date", strfmt.Date(m).String(), formats); err != nil {`,
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("aliased_reader.go",
		[]string{
			`type AliasedReader io.ReadCloser`,
		},
		// no validations for binary
		[]string{
			`validate.MinLength(`,
			`validate.MaxLength(`,
			`validate.Pattern(`,
			`validate.FormatOf(`,
		},
		[]string{
			// disable log assertions (dodgy with parallel tests)
			// `warning: validation pattern`,
			// `warning: validation minLength (value: 16) not compatible with type for format "binary". Skipped`,
			// `warning: validation maxLength (value: 20) not compatible with type for format "binary". Skipped`,
		}, noLines,
	)

	flattenRun.AddExpectations("aliased_string.go", []string{
		`type AliasedString string`,
		`func (m AliasedString) Validate(formats strfmt.Registry) error {`,
		`if err := validate.MinLength("", "body", string(m), 16); err != nil {`,
		`if err := validate.MaxLength("", "body", string(m), 20); err != nil {`,
		"if err := validate.Pattern(\"\", \"body\", string(m), `(\\d+)/(\\d+)`); err != nil {",
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("object.go", []string{
		`func (m *Object) Validate(formats strfmt.Registry) error {`,
		`if err := m.validateP0(formats); err != nil {`,
		`if err := m.validateP1(formats); err != nil {`,
		`if err := m.validateP10(formats); err != nil {`,
		`if err := m.validateP11(formats); err != nil {`,
		`if err := m.validateP1Nullable(formats); err != nil {`,
		`if err := m.validateP2(formats); err != nil {`,
		`if err := m.validateP3(formats); err != nil {`,
		`if err := m.validateP4(formats); err != nil {`,
		`if err := m.validateP5(formats); err != nil {`,
		`if err := m.validateP6(formats); err != nil {`,
		`if err := m.validateP7(formats); err != nil {`,
		`if err := m.validateP8(formats); err != nil {`,
		`if err := m.validateP9(formats); err != nil {`,
		`func (m *Object) validateP0(formats strfmt.Registry) error {`,
		`if err := validate.Required("p0", "body", m.P0); err != nil {`,
		`if err := validate.MinLength("p0", "body", *m.P0, 12); err != nil {`,
		`if err := validate.MaxLength("p0", "body", *m.P0, 16); err != nil {`,
		"if err := validate.Pattern(\"p0\", \"body\", *m.P0, `(\\d+)-(\\d+)`); err != nil {",
		`func (m *Object) validateP1(formats strfmt.Registry) error {`,
		`if err := validate.Required("p1", "body", m.P1); err != nil {`,
		`if err := validate.MinLength("p1", "body", m.P1.String(), 12); err != nil {`,
		`if err := validate.MaxLength("p1", "body", m.P1.String(), 16); err != nil {`,
		"if err := validate.Pattern(\"p1\", \"body\", m.P1.String(), `(\\d+)/(\\d+)`); err != nil {",
		`if err := validate.FormatOf("p1", "body", "date", m.P1.String(), formats); err != nil {`,
		`func (m *Object) validateP10(formats strfmt.Registry) error {`,
		`if err := validate.RequiredString("p10", "body", m.P10); err != nil {`,
		`if err := validate.MinLength("p10", "body", m.P10, 12); err != nil {`,
		`if err := validate.MaxLength("p10", "body", m.P10, 16); err != nil {`,
		"if err := validate.Pattern(\"p10\", \"body\", m.P10, `(\\d+)-(\\d+)`); err != nil {",
		`func (m *Object) validateP11(formats strfmt.Registry) error {`,
		`	if err := validate.Required("p11", "body", m.P11); err != nil {`,
		`	if err := validate.Required("p11", "body", m.P11); err != nil {`,
		`	if m.P11 != nil {`,
		`		if err := m.P11.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("p11")`,
		`func (m *Object) validateP1Nullable(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.P1Nullable) {`,
		`	if err := validate.MinLength("p1Nullable", "body", m.P1Nullable.String(), 12); err != nil {`,
		`	if err := validate.MaxLength("p1Nullable", "body", m.P1Nullable.String(), 16); err != nil {`,
		"	if err := validate.Pattern(\"p1Nullable\", \"body\", m.P1Nullable.String(), `(\\d+)/(\\d+)`); err != nil {",
		`	if err := validate.FormatOf("p1Nullable", "body", "date", m.P1Nullable.String(), formats); err != nil {`,
		`func (m *Object) validateP2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.P2) {`,
		`	if err := validate.MinLength("p2", "body", m.P2.String(), 12); err != nil {`,
		`	if err := validate.MaxLength("p2", "body", m.P2.String(), 16); err != nil {`,
		"	if err := validate.Pattern(\"p2\", \"body\", m.P2.String(), `(\\d+)-(\\d+)`); err != nil {",
		`	if err := validate.FormatOf("p2", "body", "uuid", m.P2.String(), formats); err != nil {`,
		`func (m *Object) validateP3(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.P3) {`,
		`	if err := validate.MinLength("p3", "body", m.P3.String(), 12); err != nil {`,
		`	if err := validate.MaxLength("p3", "body", m.P3.String(), 16); err != nil {`,
		"	if err := validate.Pattern(\"p3\", \"body\", m.P3.String(), `(\\d+)-(\\d+)`); err != nil {",
		`	if err := validate.FormatOf("p3", "body", "datetime", m.P3.String(), formats); err != nil {`,
		`func (m *Object) validateP4(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.P4) {`,
		`	if err := validate.MinLength("p4", "body", m.P4.String(), 12); err != nil {`,
		`	if err := validate.MaxLength("p4", "body", m.P4.String(), 16); err != nil {`,
		"	if err := validate.Pattern(\"p4\", \"body\", m.P4.String(), `(\\d+)-(\\d+)`); err != nil {",
		`	if err := validate.FormatOf("p4", "body", "bsonobjectid", m.P4.String(), formats); err != nil {`,
		`func (m *Object) validateP5(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.P5) {`,
		`	if err := validate.MinLength("p5", "body", m.P5.String(), 12); err != nil {`,
		`	if err := validate.MaxLength("p5", "body", m.P5.String(), 16); err != nil {`,
		"	if err := validate.Pattern(\"p5\", \"body\", m.P5.String(), `(\\d+)-(\\d+)`); err != nil {",
		`	if err := validate.FormatOf("p5", "body", "duration", m.P5.String(), formats); err != nil {`,
		`func (m *Object) validateP6(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.P6) {`,
		`	if err := validate.MinLength("p6", "body", m.P6.String(), 12); err != nil {`,
		`	if err := validate.MaxLength("p6", "body", m.P6.String(), 16); err != nil {`,
		"	if err := validate.Pattern(\"p6\", \"body\", m.P6.String(), `(\\d+)-(\\d+)`); err != nil {",
		`func (m *Object) validateP7(formats strfmt.Registry) error {`,
		`	if err := validate.Required("p7", "body", io.ReadCloser(m.P7)); err != nil {`,
		`func (m *Object) validateP8(formats strfmt.Registry) error {`,
		`	if err := validate.Required("p8", "body", m.P8); err != nil {`,
		`	if err := validate.Required("p8", "body", m.P8); err != nil {`,
		`	if m.P8 != nil {`,
		`		if err := m.P8.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("p8")`,
		`func (m *Object) validateP9(formats strfmt.Registry) error {`,
		`	if err := validate.Required("p9", "body", AliasedReader(m.P9)); err != nil {`,
	}, todo,
		[]string{
			// disable log assertions (dodgy with parallel tests)
			// warnings about no validations for binary
			// `warning: validation minLength (value: 12) not compatible with type for format "binary". Skipped`,
			// `warning: validation maxLength (value: 16) not compatible with type for format "binary". Skipped`,
			// `warning: validation pattern`,
		}, noLines)
}

func initFixture2448() {
	f := newModelFixture("../fixtures/bugs/2448/fixture-2448.yaml", "numerical validations")
	flattenRun := f.AddRun(false).WithMinimalFlatten(true)

	flattenRun.AddExpectations("integers.go", []string{
		`if err := validate.MinimumInt("i0", "body", m.I0, 10, true); err != nil {`,
		`if err := validate.MaximumInt("i0", "body", m.I0, 100, true); err != nil {`,
		`if err := validate.MultipleOfInt("i0", "body", m.I0, 10); err != nil {`,
		`if err := validate.MinimumInt("i1", "body", int64(m.I1), 10, true); err != nil {`,
		`if err := validate.MaximumInt("i1", "body", int64(m.I1), 100, true); err != nil {`,
		`if err := validate.MultipleOfInt("i1", "body", int64(m.I1), 10); err != nil {`,
		`if err := validate.MinimumInt("i2", "body", m.I2, 10, true); err != nil {`,
		`if err := validate.MaximumInt("i2", "body", m.I2, 100, true); err != nil {`,
		`if err := validate.MultipleOfInt("i2", "body", m.I2, 10); err != nil {`,
		`if err := validate.MinimumInt("i3", "body", int64(m.I3), 10, true); err != nil {`,
		`if err := validate.MaximumInt("i3", "body", int64(m.I3), 100, true); err != nil {`,
		`if err := validate.MultipleOfInt("i3", "body", int64(m.I3), 10); err != nil {`,
		`if err := validate.MultipleOf("i4", "body", float64(m.I4), 10.5); err != nil {`,
		`if err := validate.MinimumUint("ui1", "body", uint64(m.Ui1), 10, true); err != nil {`,
		`if err := validate.MaximumUint("ui1", "body", uint64(m.Ui1), 100, true); err != nil {`,
		`if err := validate.MultipleOfUint("ui1", "body", uint64(m.Ui1), 10); err != nil {`,
		`if err := validate.MinimumUint("ui2", "body", m.Ui2, 10, true); err != nil {`,
		`if err := validate.MaximumUint("ui2", "body", m.Ui2, 100, true); err != nil {`,
		`if err := validate.MultipleOfUint("ui2", "body", m.Ui2, 10); err != nil {`,
		`if err := validate.MinimumUint("ui3", "body", uint64(m.Ui3), 10, true); err != nil {`,
		`if err := validate.MaximumUint("ui3", "body", uint64(m.Ui3), 100, true); err != nil {`,
		`if err := validate.MultipleOfUint("ui3", "body", uint64(m.Ui3), 10); err != nil {`,
		`if err := validate.MultipleOf("ui4", "body", float64(m.Ui4), 10.5); err != nil {`,
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("numbers.go", []string{
		`if err := validate.Minimum("f0", "body", m.F0, 10, true); err != nil {`,
		`if err := validate.Maximum("f0", "body", m.F0, 100, true); err != nil {`,
		`if err := validate.MultipleOf("f0", "body", m.F0, 10); err != nil {`,
		`if err := validate.Minimum("f1", "body", float64(m.F1), 10, true); err != nil {`,
		`if err := validate.Maximum("f1", "body", float64(m.F1), 100, true); err != nil {`,
		`if err := validate.MultipleOf("f1", "body", float64(m.F1), 10); err != nil {`,
		`if err := validate.Minimum("f2", "body", m.F2, 10, true); err != nil {`,
		`if err := validate.Maximum("f2", "body", m.F2, 100, true); err != nil {`,
		`if err := validate.MultipleOf("f2", "body", m.F2, 10); err != nil {`,
	}, todo, noLines, noLines)
}
func initFixture2400() {
	f := newModelFixture("../fixtures/bugs/2400/fixture-2400.yaml", "required aliased primitive")
	flattenRun := f.AddRun(false).WithMinimalFlatten(true)

	flattenRun.AddExpectations("signup_request.go", []string{
		`type SignupRequest struct {`,
		`Email *string`,
	}, todo, noLines, noLines)
	flattenRun.AddExpectations("signup_request2.go", []string{
		`type SignupRequest2 struct {`,
		`Email *Email`,
	}, todo, noLines, noLines)
}

func initFixture2381() {
	f := newModelFixture("../fixtures/bugs/2381/fixture-2381.yaml", "required $ref primitive")
	flattenRun := f.AddRun(false).WithMinimalFlatten(true)
	expandRun := f.AddRun(true)

	flattenRun.AddExpectations("my_object.go", []string{
		`RequiredReferencedObject MyObjectRef`,  // this is an interface{}
		`RequiredReferencedString *MyStringRef`, // alias to primitive
		`RequiredString *string`,                // unaliased primitive
		`RequiredReferencedArray MyArrayRef`,    // no need to use a pointer
		`RequiredReferencedMap MyMapRef`,        // no need to use a pointer
		`RequiredReferencedStruct *MyStructRef`, // pointer to struct
		`func (m *MyObject) validateRequiredReferencedObject(formats strfmt.Registry) error {`,
		`	if m.RequiredReferencedObject == nil {`,
		`		return errors.Required("required_referenced_object", "body", nil)`,
		`func (m *MyObject) validateRequiredReferencedString(formats strfmt.Registry) error {`,
		`	if err := validate.Required("required_referenced_string", "body", m.RequiredReferencedString); err != nil {`,
		`	if err := m.RequiredReferencedString.Validate(formats); err != nil {`,
		`		if ve, ok := err.(*errors.Validation); ok {`,
		`			return ve.ValidateName("required_referenced_string")`,
		`func (m *MyObject) validateRequiredString(formats strfmt.Registry) error {`,
		`	if err := validate.Required("required_string", "body", m.RequiredString); err != nil {`,
		`func (m *MyObject) validateRequiredReferencedStruct(formats strfmt.Registry) error {`,
		`if err := validate.Required("required_referenced_struct", "body", m.RequiredReferencedStruct); err != nil {`,
		`func (m *MyObject) validateRequiredReferencedArray(formats strfmt.Registry) error {`,
		`if err := validate.Required("required_referenced_array", "body", m.RequiredReferencedArray); err != nil {`,
		`func (m *MyObject) validateRequiredReferencedMap(formats strfmt.Registry) error {`,
		`if err := validate.Required("required_referenced_map", "body", m.RequiredReferencedMap); err != nil {`,
	}, todo, noLines, noLines)

	expandRun.AddExpectations("my_object.go", []string{
		`func (m *MyObject) validateRequiredReferencedObject(formats strfmt.Registry) error {`,
		`	if m.RequiredReferencedObject == nil {`,
		`		return errors.Required("required_referenced_object", "body", nil)`,
		`func (m *MyObject) validateRequiredReferencedString(formats strfmt.Registry) error {`,
		`	if err := validate.Required("required_referenced_string", "body", m.RequiredReferencedString); err != nil {`,
		`func (m *MyObject) validateRequiredString(formats strfmt.Registry) error {`,
		`	if err := validate.Required("required_string", "body", m.RequiredString); err != nil {`,
	}, todo, noLines, noLines)
}

func initFixture2300() {
	f := newModelFixture("../fixtures/bugs/2300/fixture-2300.yaml", "required interface{} is validated with against nil")
	flattenRun := f.AddRun(false).WithMinimalFlatten(true)

	// test behaviour with all structs made anonymous (inlined)
	expandRun := f.AddRun(true)

	flattenRun.AddExpectations("obj.go", []string{
		`func (m *Obj) validateReq1(formats strfmt.Registry) error {`,
		`	if m.Req1 == nil {`,
		`		return errors.Required("req1", "body", nil)`,
		`func (m *Obj) validateReq2(formats strfmt.Registry) error {`,
		`	if m.Req2 == nil {`,
		`		return errors.Required("req2", "body", nil)`,
		`func (m *Obj) validateReq3(formats strfmt.Registry) error {`,
		`	if m.Req3 == nil {`,
		`		return errors.Required("req3", "body", nil)`,
		`func (m *Obj) validateReq4(formats strfmt.Registry) error {`,
		`	if err := validate.Required("req4", "body", m.Req4); err != nil {`,
		`	if m.Req4 != nil {`,
		`		if err := m.Req4.Validate(formats); err != nil {`,
		`func (m *Obj) validateReq5(formats strfmt.Registry) error {`,
		`	if err := validate.Required("req5", "body", m.Req5); err != nil {`,
		`	if m.Req5 != nil {`,
		`		if err := m.Req5.Validate(formats); err != nil {`,
	}, todo, noLines, noLines)

	// on anonymous
	expandRun.AddExpectations("obj.go", []string{
		`Req4 map[string]interface{}`,
		`func (m *Obj) validateReq1(formats strfmt.Registry) error {`,
		`	if m.Req1 == nil {`,
		`		return errors.Required("req1", "body", nil)`,
		`func (m *Obj) validateReq2(formats strfmt.Registry) error {`,
		`	if m.Req2 == nil {`,
		`		return errors.Required("req2", "body", nil)`,
		`func (m *Obj) validateReq3(formats strfmt.Registry) error {`,
		`	if m.Req3 == nil {`,
		`		return errors.Required("req3", "body", nil)`,
		`func (m *Obj) validateReq4(formats strfmt.Registry) error {`,
		`if err := validate.Required("req4", "body", m.Req4); err != nil {`,
		`for k := range m.Req4 {`,
		`if err := validate.Required("req4"+"."+k, "body", m.Req4[k]); err != nil {`,
		`func (m *Obj) validateReq5(formats strfmt.Registry) error {`,
		`	if err := validate.Required("req5", "body", m.Req5); err != nil {`,
		`	if m.Req5 != nil {`,
		`		if err := m.Req5.Validate(formats); err != nil {`,
	}, todo, noLines, noLines)
}

func initFixture2081() {
	f := newModelFixture("../fixtures/bugs/2081/fixture-2081.yaml", "required interface{} is validated with against nil")
	flattenRun := f.AddRun(false).WithMinimalFlatten(true)

	// interface{}
	flattenRun.AddExpectations("event.go", []string{
		`func (m *Event) validateValue(formats strfmt.Registry) error {`,
		`if m.Value == nil {`,
		`return errors.Required("value", "body", nil)`,
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("events.go", []string{
		`func (m *EventsItems0) validateA(formats strfmt.Registry) error {`,
		`if m.A == nil {`,
		`return errors.Required("a", "body", nil)`,
	}, todo, noLines, noLines)
}

func initFixture936ReadOnly() {
	f := newModelFixture("../fixtures/enhancements/936/fixture-936.yml", "check ReadOnly ContextValidate is generated properly")
	flattenRun := f.AddRun(false).WithMinimalFlatten(true)

	// object simple has 2 read only feilds
	flattenRun.AddExpectations("object1.go", []string{
		// object1
		`func (m *Object1) ContextValidate(ctx context.Context, formats strfmt.Registry)`,
		`m.contextValidateID(ctx, formats)`,
		`m.contextValidateName(ctx, formats)`,
		`) contextValidateID(ctx context.Context, formats strfmt.Registry)`,
		`) contextValidateName(ctx context.Context, formats strfmt.Registry)`,
		`validate.ReadOnly(ctx`,
	}, todo, noLines, noLines)

	// object2 composed of object1
	flattenRun.AddExpectations("object2.go", []string{
		`) ContextValidate(ctx context.Context, formats strfmt.Registry)`,
		`contextValidateObjectRef(ctx, formats)`,
		`) contextValidateObjectRef(ctx context.Context, formats strfmt.Registry) `,
		`m.ObjectRef.ContextValidate(ctx, formats)`,
	}, todo, noLines, noLines)

	// object3 has inline object
	flattenRun.AddExpectations("object3.go", []string{
		`) ContextValidate(ctx context.Context, formats strfmt.Registry)`,
		`contextValidateName(ctx, formats)`,
		`) contextValidateName(ctx context.Context, formats strfmt.Registry)`,
		`validate.ReadOnly(ctx`,
	}, todo, noLines, noLines)

	// object4 is array with readonly string element
	flattenRun.AddExpectations("object4.go", []string{
		`) ContextValidate(ctx context.Context, formats strfmt.Registry)`,
		`for i := 0; i < len(m); i++`,
		`validate.ReadOnly(ctx`,
	}, todo, noLines, noLines)

	// object5 is string alias
	flattenRun.AddExpectations("object5.go", []string{
		`) ContextValidate(ctx context.Context, formats strfmt.Registry)`,
		`validate.ReadOnly(ctx`,
	}, todo, noLines, noLines)

	// object6 is array of object5
	flattenRun.AddExpectations("object6.go", []string{
		`) ContextValidate(ctx context.Context, formats strfmt.Registry)`,
		`for i := 0; i < len(m); i++`,
		`m[i].ContextValidate(ctx, formats)`,
	}, todo, noLines, noLines)

	// object7 is all of object5 and object4 and fields
	flattenRun.AddExpectations("object7.go", []string{
		`) ContextValidate(ctx context.Context, formats strfmt.Registry)`,
		`m.Object5.ContextValidate(ctx, formats)`,
		`m.Object4.ContextValidate(ctx, formats)`,
		// field one
		`m.contextValidateObject7Field1(ctx, formats)`,
		`contextValidateObject7Field1(ctx context.Context, formats strfmt.Registry)`,
		// field two should missing since not readOnly
		// field three
		`m.contextValidateObject7Field3(ctx, formats)`,
		`contextValidateObject7Field3(ctx context.Context, formats strfmt.Registry)`,
		`m.Object7Field3.ContextValidate(ctx, formats)`,
	}, todo, noLines,
		[]string{
			`m.contextValidateObject7Field2(ctx, formats)`,
		})
	// x go type
	flattenRun.AddExpectations("time.go", []string{
		`) ContextValidate(ctx context.Context, formats strfmt.Registry) error`,
		`var f interface{} = m.Time`,
	}, todo, noLines, noLines)
	// additional properties
	flattenRun.AddExpectations("object8.go", []string{
		`) ContextValidate(ctx context.Context, formats strfmt.Registry) error`,
		`for k := range m`,
		`validate.ReadOnly(ctx`,
	}, todo, noLines, noLines)
	flattenRun.AddExpectations("object9.go", []string{
		`) ContextValidate(ctx context.Context, formats strfmt.Registry) error`,
		`m.contextValidateObject9Field1(ctx, formats)`,
		`validate.ReadOnly(ctx`,
	}, todo, noLines, noLines)
}

func initFixture2220() {
	// NOTE(fred): this test merely asserts that template refactoring (essentially dealing with hite space gobbling etc)
	// properly runs against the case of base type with additionalProperties.
	//
	// TODO(fred): should actually fix the problem in base type model rendering
	f := newModelFixture("../fixtures/bugs/2220/fixture-2220.yaml", "check base type with additional properties")
	flattenRun := f.AddRun(false).WithMinimalFlatten(true)

	flattenRun.AddExpectations("object.go", []string{
		// This asserts our template announcement about forthcoming fix (used to  be a func commented out of luck)
		`// AdditionalProperties in base type shoud be handled just like regular properties`,
		`// At this moment, the base type property is pushed down to the subtype`,
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("component.go", []string{
		// This asserts the current schema layout, which works but does not honour inheritance from the base type
		"ObjectAdditionalProperties map[string]interface{} `json:\"-\"`",
	}, todo, noLines, noLines)
}

func initFixture2116() {
	f := newModelFixture("../fixtures/bugs/2116/fixture-2116.yaml", "check x-omitempty and x-nullable with $ref")
	flattenRun := f.AddRun(false).WithMinimalFlatten(true)

	flattenRun.AddExpectations("case1_fail_omitempty_false_not_hoisted_by_ref.go", []string{
		"Body *ObjectWithOmitemptyFalse `json:\"Body\"`",
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("case2_fail_omitempty_false_not_overridden_by_ref_sibling.go", []string{
		"Body *ObjectWithOmitemptyTrue `json:\"Body,omitempty\"`",
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("case3_pass_object_nullable_false_hoisted_by_ref.go", []string{
		"Body ObjectWithNullableFalse `json:\"Body,omitempty\"`",
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("case4_pass_object_nullable_false_overriden_by_ref_sibling.go", []string{
		"Body *ObjectWithNullableTrue `json:\"Body,omitempty\"`",
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("array_with_default.go", []string{
		"type ArrayWithDefault []string",
	}, append(todo, "omitempty"), noLines, noLines)

	flattenRun.AddExpectations("array_with_no_omit_empty.go", []string{
		"type ArrayWithNoOmitEmpty []string",
	}, append(todo, "omitempty"), noLines, noLines)

	flattenRun.AddExpectations("array_with_nullable.go", []string{
		"type ArrayWithNullable []string",
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("array_with_nullable_items.go", []string{
		"type ArrayWithNullableItems []*string",
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("array_with_omit_empty.go", []string{
		"type ArrayWithOmitEmpty []string",
	}, append(todo, "omitempty"), noLines, noLines)

	flattenRun.AddExpectations("object_with_arrays.go", []string{
		"Array0 ArrayWithDefault `json:\"array0,omitempty\"`",
		"Array1 []string `json:\"array1\"`",
		"Array11 []string `json:\"array11,omitempty\"`",
		"Array12 []string `json:\"array12\"`",
		"Array2 ArrayWithOmitEmpty `json:\"array2,omitempty\"`",
		"Array3 ArrayWithNoOmitEmpty `json:\"array3\"`",
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("object_with_nullable_false.go", []string{
		"Data interface{} `json:\"Data,omitempty\"`",
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("object_with_nullable_true.go", []string{
		"Data interface{} `json:\"Data,omitempty\"`",
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("object_with_omitempty_false.go", []string{
		"Data interface{} `json:\"Data,omitempty\"`",
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("object_with_omitempty_true.go", []string{
		"Data interface{} `json:\"Data,omitempty\"`",
	}, todo, noLines, noLines)

	flattenRun.AddExpectations("array_with_omit_empty_items.go", []string{
		"type ArrayWithOmitEmptyItems []string",
	}, append(todo, "omitempty"), noLines, noLines)
}

func initFixture2071() {
	f := newModelFixture("../fixtures/bugs/2071/fixture-2071.yaml", "check allOf serializer when x-go-name is present")
	flattenRun := f.AddRun(false).WithMinimalFlatten(true)
	flattenRun.AddExpectations("cat.go", []string{
		"var dataAO1 struct {\n\t\tSomeAbility *string `json:\"ability,omitempty\"`",
	},
		// not expected
		append(todo, "SomeAbility *string `json:\"SomeAbility,omitempty\"`"),
		// output in log
		noLines,
		noLines)
}

func initFixture1479Part() {
	// testing ../fixtures/bugs/1479/fixture-1479-part.yaml with flatten and expand (--skip-flatten)

	/*
		The breakage with allOf occurs when a schema with an allOf has itself a
		property which is an allOf construct
	*/

	f := newModelFixture("../fixtures/bugs/1479/fixture-1479-part.yaml", "check nested AllOf validations (from Pouch Engine API)")
	flattenRun := f.AddRun(false)
	expandRun := f.AddRun(true)

	// load expectations for model: container_create_config_all_of1.go
	flattenRun.AddExpectations("container_create_config_all_of1.go", []string{
		`type ContainerCreateConfigAllOf1 struct {`,
		"	HostConfig *HostConfig `json:\"HostConfig,omitempty\"`",
		"	NetworkingConfig *NetworkingConfig `json:\"NetworkingConfig,omitempty\"`",
		`func (m *ContainerCreateConfigAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateHostConfig(formats); err != nil {`,
		`	if err := m.validateNetworkingConfig(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ContainerCreateConfigAllOf1) validateHostConfig(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.HostConfig) {`,
		`	if m.HostConfig != nil {`,
		`		if err := m.HostConfig.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("HostConfig"`,
		`func (m *ContainerCreateConfigAllOf1) validateNetworkingConfig(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.NetworkingConfig) {`,
		`	if m.NetworkingConfig != nil {`,
		`		if err := m.NetworkingConfig.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("NetworkingConfig"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: container_config.go
	flattenRun.AddExpectations("container_config.go", []string{
		`type ContainerConfig struct {`,
		"	ArgsEscaped bool `json:\"ArgsEscaped,omitempty\"`",
		"	AttachStderr bool `json:\"AttachStderr,omitempty\"`",
		"	AttachStdin bool `json:\"AttachStdin,omitempty\"`",
		"	AttachStdout bool `json:\"AttachStdout,omitempty\"`",
		"	Cmd []string `json:\"Cmd\"`",
		"	DiskQuota map[string]string `json:\"DiskQuota,omitempty\"`",
		"	Domainname string `json:\"Domainname,omitempty\"`",
		"	Entrypoint []string `json:\"Entrypoint\"`",
		"	Env []string `json:\"Env\"`",
		"	ExposedPorts map[string]interface{} `json:\"ExposedPorts,omitempty\"`",
		"	Hostname strfmt.Hostname `json:\"Hostname,omitempty\"`",
		"	Image string `json:\"Image\"`",
		"	InitScript string `json:\"InitScript,omitempty\"`",
		"	Labels map[string]string `json:\"Labels,omitempty\"`",
		"	MacAddress string `json:\"MacAddress,omitempty\"`",
		"	NetworkDisabled bool `json:\"NetworkDisabled,omitempty\"`",
		"	OnBuild []string `json:\"OnBuild\"`",
		"	OpenStdin bool `json:\"OpenStdin,omitempty\"`",
		"	QuotaID string `json:\"QuotaID,omitempty\"`",
		"	Rich bool `json:\"Rich,omitempty\"`",
		"	RichMode string `json:\"RichMode,omitempty\"`",
		"	Shell []string `json:\"Shell\"`",
		"	SpecAnnotation map[string]string `json:\"SpecAnnotation,omitempty\"`",
		"	StdinOnce bool `json:\"StdinOnce,omitempty\"`",
		"	StopSignal string `json:\"StopSignal,omitempty\"`",
		"	StopTimeout *int64 `json:\"StopTimeout,omitempty\"`",
		"	Tty bool `json:\"Tty,omitempty\"`",
		"	User string `json:\"User,omitempty\"`",
		"	Volumes map[string]interface{} `json:\"Volumes,omitempty\"`",
		"	WorkingDir string `json:\"WorkingDir,omitempty\"`",
		`func (m *ContainerConfig) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateExposedPorts(formats); err != nil {`,
		`	if err := m.validateHostname(formats); err != nil {`,
		`	if err := m.validateImage(formats); err != nil {`,
		`	if err := m.validateRichMode(formats); err != nil {`,
		`	if err := m.validateVolumes(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`var containerConfigExposedPortsValueEnum []interface{`,
		`	var res []interface{`,
		"	if err := json.Unmarshal([]byte(`[{}]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		containerConfigExposedPortsValueEnum = append(containerConfigExposedPortsValueEnum, v`,
		`func (m *ContainerConfig) validateExposedPortsValueEnum(path, location string, value interface{}) error {`,
		`	if err := validate.EnumCase(path, location, value, containerConfigExposedPortsValueEnum, true); err != nil {`,
		`func (m *ContainerConfig) validateExposedPorts(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ExposedPorts) {`,
		`	for k := range m.ExposedPorts {`,
		`		if err := m.validateExposedPortsValueEnum("ExposedPorts"+"."+k, "body", m.ExposedPorts[k]); err != nil {`,
		`func (m *ContainerConfig) validateHostname(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Hostname) {`,
		`	if err := validate.MinLength("Hostname", "body", m.Hostname.String(), 1); err != nil {`,
		`	if err := validate.FormatOf("Hostname", "body", "hostname", m.Hostname.String(), formats); err != nil {`,
		`func (m *ContainerConfig) validateImage(formats strfmt.Registry) error {`,
		`	if err := validate.RequiredString("Image", "body", m.Image); err != nil {`,
		`var containerConfigTypeRichModePropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"dumb-init\",\"sbin-init\",\"systemd\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		containerConfigTypeRichModePropEnum = append(containerConfigTypeRichModePropEnum, v`,
		`	ContainerConfigRichModeDumbDashInit string = "dumb-init"`,
		`	ContainerConfigRichModeSbinDashInit string = "sbin-init"`,
		`	ContainerConfigRichModeSystemd string = "systemd"`,
		`func (m *ContainerConfig) validateRichModeEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, containerConfigTypeRichModePropEnum, true); err != nil {`,
		`func (m *ContainerConfig) validateRichMode(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.RichMode) {`,
		`	if err := m.validateRichModeEnum("RichMode", "body", m.RichMode); err != nil {`,
		`var containerConfigVolumesValueEnum []interface{`,
		`	var res []interface{`,
		"	if err := json.Unmarshal([]byte(`[{}]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		containerConfigVolumesValueEnum = append(containerConfigVolumesValueEnum, v`,
		`func (m *ContainerConfig) validateVolumesValueEnum(path, location string, value interface{}) error {`,
		`	if err := validate.EnumCase(path, location, value, containerConfigVolumesValueEnum, true); err != nil {`,
		`func (m *ContainerConfig) validateVolumes(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Volumes) {`,
		`	for k := range m.Volumes {`,
		`		if err := m.validateVolumesValueEnum("Volumes"+"."+k, "body", m.Volumes[k]); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("container_config.go", flattenRun.ExpectedFor("ContainerConfig").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: host_config_all_of0_log_config.go
	flattenRun.AddExpectations("host_config_all_of0_log_config.go", []string{
		`type HostConfigAllOf0LogConfig struct {`,
		"	Config map[string]string `json:\"Config,omitempty\"`",
		"	Type string `json:\"Type,omitempty\"`",
		`func (m *HostConfigAllOf0LogConfig) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateType(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`var hostConfigAllOf0LogConfigTypeTypePropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"json-file\",\"syslog\",\"journald\",\"gelf\",\"fluentd\",\"awslogs\",\"splunk\",\"etwlogs\",\"none\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		hostConfigAllOf0LogConfigTypeTypePropEnum = append(hostConfigAllOf0LogConfigTypeTypePropEnum, v`,
		`	HostConfigAllOf0LogConfigTypeJSONDashFile string = "json-file"`,
		`	HostConfigAllOf0LogConfigTypeSyslog string = "syslog"`,
		`	HostConfigAllOf0LogConfigTypeJournald string = "journald"`,
		`	HostConfigAllOf0LogConfigTypeGelf string = "gelf"`,
		`	HostConfigAllOf0LogConfigTypeFluentd string = "fluentd"`,
		`	HostConfigAllOf0LogConfigTypeAwslogs string = "awslogs"`,
		`	HostConfigAllOf0LogConfigTypeSplunk string = "splunk"`,
		`	HostConfigAllOf0LogConfigTypeEtwlogs string = "etwlogs"`,
		`	HostConfigAllOf0LogConfigTypeNone string = "none"`,
		`func (m *HostConfigAllOf0LogConfig) validateTypeEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, hostConfigAllOf0LogConfigTypeTypePropEnum, true); err != nil {`,
		`func (m *HostConfigAllOf0LogConfig) validateType(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Type) {`,
		`	if err := m.validateTypeEnum("Type", "body", m.Type); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: restart_policy.go
	flattenRun.AddExpectations("restart_policy.go", []string{
		`type RestartPolicy struct {`,
		"	MaximumRetryCount int64 `json:\"MaximumRetryCount,omitempty\"`",
		"	Name string `json:\"Name,omitempty\"`",
		// empty validation
		"func (m *RestartPolicy) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("restart_policy.go", flattenRun.ExpectedFor("RestartPolicy").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: endpoint_ip_a_m_config.go
	flattenRun.AddExpectations("endpoint_ip_a_m_config.go", []string{
		`type EndpointIPAMConfig struct {`,
		"	IPV4Address string `json:\"IPv4Address,omitempty\"`",
		"	IPV6Address string `json:\"IPv6Address,omitempty\"`",
		"	LinkLocalIPs []string `json:\"LinkLocalIPs\"`",
		// empty validation
		"func (m *EndpointIPAMConfig) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("endpoint_ip_a_m_config.go", flattenRun.ExpectedFor("EndpointIPAMConfig").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: host_config_all_of0.go
	flattenRun.AddExpectations("host_config_all_of0.go", []string{
		`type HostConfigAllOf0 struct {`,
		"	AutoRemove bool `json:\"AutoRemove,omitempty\"`",
		"	Binds []string `json:\"Binds\"`",
		"	CapAdd []string `json:\"CapAdd\"`",
		"	CapDrop []string `json:\"CapDrop\"`",
		"	Cgroup string `json:\"Cgroup,omitempty\"`",
		"	ConsoleSize []*int64 `json:\"ConsoleSize\"`",
		"	ContainerIDFile string `json:\"ContainerIDFile,omitempty\"`",
		"	DNS []string `json:\"Dns\"`",
		"	DNSOptions []string `json:\"DnsOptions\"`",
		"	DNSSearch []string `json:\"DnsSearch\"`",
		"	EnableLxcfs bool `json:\"EnableLxcfs,omitempty\"`",
		"	ExtraHosts []string `json:\"ExtraHosts\"`",
		"	GroupAdd []string `json:\"GroupAdd\"`",
		"	InitScript string `json:\"InitScript,omitempty\"`",
		"	IpcMode string `json:\"IpcMode,omitempty\"`",
		"	Isolation string `json:\"Isolation,omitempty\"`",
		"	Links []string `json:\"Links\"`",
		"	LogConfig *HostConfigAllOf0LogConfig `json:\"LogConfig,omitempty\"`",
		"	NetworkMode string `json:\"NetworkMode,omitempty\"`",
		"	OomScoreAdj int64 `json:\"OomScoreAdj,omitempty\"`",
		"	PidMode string `json:\"PidMode,omitempty\"`",
		"	Privileged bool `json:\"Privileged,omitempty\"`",
		"	PublishAllPorts bool `json:\"PublishAllPorts,omitempty\"`",
		"	ReadonlyRootfs bool `json:\"ReadonlyRootfs,omitempty\"`",
		"	RestartPolicy *RestartPolicy `json:\"RestartPolicy,omitempty\"`",
		"	Rich bool `json:\"Rich,omitempty\"`",
		"	RichMode string `json:\"RichMode,omitempty\"`",
		"	Runtime string `json:\"Runtime,omitempty\"`",
		"	SecurityOpt []string `json:\"SecurityOpt\"`",
		"	ShmSize *int64 `json:\"ShmSize,omitempty\"`",
		"	StorageOpt map[string]string `json:\"StorageOpt,omitempty\"`",
		"	Sysctls map[string]string `json:\"Sysctls,omitempty\"`",
		"	Tmpfs map[string]string `json:\"Tmpfs,omitempty\"`",
		"	UTSMode string `json:\"UTSMode,omitempty\"`",
		"	UsernsMode string `json:\"UsernsMode,omitempty\"`",
		"	VolumeDriver string `json:\"VolumeDriver,omitempty\"`",
		"	VolumesFrom []string `json:\"VolumesFrom\"`",
		`func (m *HostConfigAllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateConsoleSize(formats); err != nil {`,
		`	if err := m.validateIsolation(formats); err != nil {`,
		`	if err := m.validateLogConfig(formats); err != nil {`,
		`	if err := m.validateOomScoreAdj(formats); err != nil {`,
		`	if err := m.validateRestartPolicy(formats); err != nil {`,
		`	if err := m.validateRichMode(formats); err != nil {`,
		`	if err := m.validateShmSize(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *HostConfigAllOf0) validateConsoleSize(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ConsoleSize) {`,
		`	iConsoleSizeSize := int64(len(m.ConsoleSize)`,
		`	if err := validate.MinItems("ConsoleSize", "body", iConsoleSizeSize, 2); err != nil {`,
		`	if err := validate.MaxItems("ConsoleSize", "body", iConsoleSizeSize, 2); err != nil {`,
		`	for i := 0; i < len(m.ConsoleSize); i++ {`,
		// do we need...?
		`		if swag.IsZero(m.ConsoleSize[i]) {`,
		`		if err := validate.MinimumInt("ConsoleSize"+"."+strconv.Itoa(i), "body", *m.ConsoleSize[i], 0, false); err != nil {`,
		`var hostConfigAllOf0TypeIsolationPropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"default\",\"process\",\"hyperv\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		hostConfigAllOf0TypeIsolationPropEnum = append(hostConfigAllOf0TypeIsolationPropEnum, v`,
		`	HostConfigAllOf0IsolationDefault string = "default"`,
		`	HostConfigAllOf0IsolationProcess string = "process"`,
		`	HostConfigAllOf0IsolationHyperv string = "hyperv"`,
		`func (m *HostConfigAllOf0) validateIsolationEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, hostConfigAllOf0TypeIsolationPropEnum, true); err != nil {`,
		`func (m *HostConfigAllOf0) validateIsolation(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Isolation) {`,
		`	if err := m.validateIsolationEnum("Isolation", "body", m.Isolation); err != nil {`,
		`func (m *HostConfigAllOf0) validateLogConfig(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.LogConfig) {`,
		`	if m.LogConfig != nil {`,
		`		if err := m.LogConfig.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("LogConfig"`,
		`func (m *HostConfigAllOf0) validateOomScoreAdj(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.OomScoreAdj) {`,
		`	if err := validate.MinimumInt("OomScoreAdj", "body", m.OomScoreAdj, -1000, false); err != nil {`,
		`	if err := validate.MaximumInt("OomScoreAdj", "body", m.OomScoreAdj, 1000, false); err != nil {`,
		`func (m *HostConfigAllOf0) validateRestartPolicy(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.RestartPolicy) {`,
		`	if m.RestartPolicy != nil {`,
		`		if err := m.RestartPolicy.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("RestartPolicy"`,
		`var hostConfigAllOf0TypeRichModePropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"dumb-init\",\"sbin-init\",\"systemd\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		hostConfigAllOf0TypeRichModePropEnum = append(hostConfigAllOf0TypeRichModePropEnum, v`,
		`	HostConfigAllOf0RichModeDumbDashInit string = "dumb-init"`,
		`	HostConfigAllOf0RichModeSbinDashInit string = "sbin-init"`,
		`	HostConfigAllOf0RichModeSystemd string = "systemd"`,
		`func (m *HostConfigAllOf0) validateRichModeEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, hostConfigAllOf0TypeRichModePropEnum, true); err != nil {`,
		`func (m *HostConfigAllOf0) validateRichMode(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.RichMode) {`,
		`	if err := m.validateRichModeEnum("RichMode", "body", m.RichMode); err != nil {`,
		`func (m *HostConfigAllOf0) validateShmSize(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ShmSize) {`,
		`	if err := validate.MinimumInt("ShmSize", "body", *m.ShmSize, 0, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: host_config.go
	flattenRun.AddExpectations("host_config.go", []string{
		`type HostConfig struct {`,
		`	HostConfigAllOf0`,
		`	Resources`,
		`func (m *HostConfig) Validate(formats strfmt.Registry) error {`,
		`	if err := m.HostConfigAllOf0.Validate(formats); err != nil {`,
		`	if err := m.Resources.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("host_config.go", []string{
		`type HostConfig struct {`,
		"	AutoRemove bool `json:\"AutoRemove,omitempty\"`",
		"	Binds []string `json:\"Binds\"`",
		"	CapAdd []string `json:\"CapAdd\"`",
		"	CapDrop []string `json:\"CapDrop\"`",
		"	Cgroup string `json:\"Cgroup,omitempty\"`",
		"	ConsoleSize []*int64 `json:\"ConsoleSize\"`",
		"	ContainerIDFile string `json:\"ContainerIDFile,omitempty\"`",
		"	DNS []string `json:\"Dns\"`",
		"	DNSOptions []string `json:\"DnsOptions\"`",
		"	DNSSearch []string `json:\"DnsSearch\"`",
		"	EnableLxcfs bool `json:\"EnableLxcfs,omitempty\"`",
		"	ExtraHosts []string `json:\"ExtraHosts\"`",
		"	GroupAdd []string `json:\"GroupAdd\"`",
		"	InitScript string `json:\"InitScript,omitempty\"`",
		"	IpcMode string `json:\"IpcMode,omitempty\"`",
		"	Isolation string `json:\"Isolation,omitempty\"`",
		"	Links []string `json:\"Links\"`",
		"	LogConfig *HostConfigAO0LogConfig `json:\"LogConfig,omitempty\"`",
		"	NetworkMode string `json:\"NetworkMode,omitempty\"`",
		"	OomScoreAdj int64 `json:\"OomScoreAdj,omitempty\"`",
		"	PidMode string `json:\"PidMode,omitempty\"`",
		"	Privileged bool `json:\"Privileged,omitempty\"`",
		"	PublishAllPorts bool `json:\"PublishAllPorts,omitempty\"`",
		"	ReadonlyRootfs bool `json:\"ReadonlyRootfs,omitempty\"`",
		"	RestartPolicy *HostConfigAO0RestartPolicy `json:\"RestartPolicy,omitempty\"`",
		"	Rich bool `json:\"Rich,omitempty\"`",
		"	RichMode string `json:\"RichMode,omitempty\"`",
		"	Runtime string `json:\"Runtime,omitempty\"`",
		"	SecurityOpt []string `json:\"SecurityOpt\"`",
		"	ShmSize *int64 `json:\"ShmSize,omitempty\"`",
		"	StorageOpt map[string]string `json:\"StorageOpt,omitempty\"`",
		"	Sysctls map[string]string `json:\"Sysctls,omitempty\"`",
		"	Tmpfs map[string]string `json:\"Tmpfs,omitempty\"`",
		"	UTSMode string `json:\"UTSMode,omitempty\"`",
		"	UsernsMode string `json:\"UsernsMode,omitempty\"`",
		"	VolumeDriver string `json:\"VolumeDriver,omitempty\"`",
		"	VolumesFrom []string `json:\"VolumesFrom\"`",
		"	BlkioWeight uint16 `json:\"BlkioWeight,omitempty\"`",
		"	CgroupParent string `json:\"CgroupParent,omitempty\"`",
		"	CPUShares int64 `json:\"CpuShares,omitempty\"`",
		"	Memory int64 `json:\"Memory,omitempty\"`",
		`func (m *HostConfig) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateConsoleSize(formats); err != nil {`,
		`	if err := m.validateIsolation(formats); err != nil {`,
		`	if err := m.validateLogConfig(formats); err != nil {`,
		`	if err := m.validateOomScoreAdj(formats); err != nil {`,
		`	if err := m.validateRestartPolicy(formats); err != nil {`,
		`	if err := m.validateRichMode(formats); err != nil {`,
		`	if err := m.validateShmSize(formats); err != nil {`,
		`	if err := m.validateBlkioWeight(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *HostConfig) validateConsoleSize(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ConsoleSize) {`,
		`	iConsoleSizeSize := int64(len(m.ConsoleSize)`,
		`	if err := validate.MinItems("ConsoleSize", "body", iConsoleSizeSize, 2); err != nil {`,
		`	if err := validate.MaxItems("ConsoleSize", "body", iConsoleSizeSize, 2); err != nil {`,
		`	for i := 0; i < len(m.ConsoleSize); i++ {`,
		// do we need...
		`		if swag.IsZero(m.ConsoleSize[i]) {`,
		`		if err := validate.MinimumInt("ConsoleSize"+"."+strconv.Itoa(i), "body", *m.ConsoleSize[i], 0, false); err != nil {`,
		`var hostConfigTypeIsolationPropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"default\",\"process\",\"hyperv\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		hostConfigTypeIsolationPropEnum = append(hostConfigTypeIsolationPropEnum, v`,
		`func (m *HostConfig) validateIsolationEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, hostConfigTypeIsolationPropEnum, true); err != nil {`,
		`func (m *HostConfig) validateIsolation(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Isolation) {`,
		`	if err := m.validateIsolationEnum("Isolation", "body", m.Isolation); err != nil {`,
		`func (m *HostConfig) validateLogConfig(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.LogConfig) {`,
		`	if m.LogConfig != nil {`,
		`		if err := m.LogConfig.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("LogConfig"`,
		`func (m *HostConfig) validateOomScoreAdj(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.OomScoreAdj) {`,
		`	if err := validate.MinimumInt("OomScoreAdj", "body", m.OomScoreAdj, -1000, false); err != nil {`,
		`	if err := validate.MaximumInt("OomScoreAdj", "body", m.OomScoreAdj, 1000, false); err != nil {`,
		`func (m *HostConfig) validateRestartPolicy(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.RestartPolicy) {`,
		`	if m.RestartPolicy != nil {`,
		`		if err := m.RestartPolicy.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("RestartPolicy"`,
		`var hostConfigTypeRichModePropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"dumb-init\",\"sbin-init\",\"systemd\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		hostConfigTypeRichModePropEnum = append(hostConfigTypeRichModePropEnum, v`,
		`func (m *HostConfig) validateRichModeEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, hostConfigTypeRichModePropEnum, true); err != nil {`,
		`func (m *HostConfig) validateRichMode(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.RichMode) {`,
		`	if err := m.validateRichModeEnum("RichMode", "body", m.RichMode); err != nil {`,
		`func (m *HostConfig) validateShmSize(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ShmSize) {`,
		`	if err := validate.MinimumInt("ShmSize", "body", *m.ShmSize, 0, false); err != nil {`,
		`func (m *HostConfig) validateBlkioWeight(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.BlkioWeight) {`,
		`	if err := validate.MinimumUint("BlkioWeight", "body", uint64(m.BlkioWeight), 0, false); err != nil {`,
		`	if err := validate.MaximumUint("BlkioWeight", "body", uint64(m.BlkioWeight), 1000, false); err != nil {`,
		`type HostConfigAO0LogConfig struct {`,
		"	Config map[string]string `json:\"Config,omitempty\"`",
		"	Type string `json:\"Type,omitempty\"`",
		`func (m *HostConfigAO0LogConfig) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateType(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`var hostConfigAO0LogConfigTypeTypePropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"json-file\",\"syslog\",\"journald\",\"gelf\",\"fluentd\",\"awslogs\",\"splunk\",\"etwlogs\",\"none\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		hostConfigAO0LogConfigTypeTypePropEnum = append(hostConfigAO0LogConfigTypeTypePropEnum, v`,
		`	HostConfigAO0LogConfigTypeJSONDashFile string = "json-file"`,
		`	HostConfigAO0LogConfigTypeSyslog string = "syslog"`,
		`	HostConfigAO0LogConfigTypeJournald string = "journald"`,
		`	HostConfigAO0LogConfigTypeGelf string = "gelf"`,
		`	HostConfigAO0LogConfigTypeFluentd string = "fluentd"`,
		`	HostConfigAO0LogConfigTypeAwslogs string = "awslogs"`,
		`	HostConfigAO0LogConfigTypeSplunk string = "splunk"`,
		`	HostConfigAO0LogConfigTypeEtwlogs string = "etwlogs"`,
		`	HostConfigAO0LogConfigTypeNone string = "none"`,
		`func (m *HostConfigAO0LogConfig) validateTypeEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, hostConfigAO0LogConfigTypeTypePropEnum, true); err != nil {`,
		`func (m *HostConfigAO0LogConfig) validateType(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Type) {`,
		`	if err := m.validateTypeEnum("LogConfig"+"."+"Type", "body", m.Type); err != nil {`,
		`type HostConfigAO0RestartPolicy struct {`,
		"	MaximumRetryCount int64 `json:\"MaximumRetryCount,omitempty\"`",
		"	Name string `json:\"Name,omitempty\"`",
		`func (m *HostConfigAO0RestartPolicy) Validate(formats strfmt.Registry) error {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: container_create_config.go
	flattenRun.AddExpectations("container_create_config.go", []string{
		`type ContainerCreateConfig struct {`,
		`	ContainerConfig`,
		`	ContainerCreateConfigAllOf1`,
		`func (m *ContainerCreateConfig) Validate(formats strfmt.Registry) error {`,
		`	if err := m.ContainerConfig.Validate(formats); err != nil {`,
		`	if err := m.ContainerCreateConfigAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("container_create_config.go", []string{
		`type ContainerCreateConfig struct {`,
		"	ArgsEscaped bool `json:\"ArgsEscaped,omitempty\"`",
		"	AttachStderr bool `json:\"AttachStderr,omitempty\"`",
		"	AttachStdin bool `json:\"AttachStdin,omitempty\"`",
		"	AttachStdout bool `json:\"AttachStdout,omitempty\"`",
		"	Cmd []string `json:\"Cmd\"`",
		"	DiskQuota map[string]string `json:\"DiskQuota,omitempty\"`",
		"	Domainname string `json:\"Domainname,omitempty\"`",
		"	Entrypoint []string `json:\"Entrypoint\"`",
		"	Env []string `json:\"Env\"`",
		"	ExposedPorts map[string]interface{} `json:\"ExposedPorts,omitempty\"`",
		"	Hostname strfmt.Hostname `json:\"Hostname,omitempty\"`",
		"	Image string `json:\"Image\"`",
		"	InitScript string `json:\"InitScript,omitempty\"`",
		"	Labels map[string]string `json:\"Labels,omitempty\"`",
		"	MacAddress string `json:\"MacAddress,omitempty\"`",
		"	NetworkDisabled bool `json:\"NetworkDisabled,omitempty\"`",
		"	OnBuild []string `json:\"OnBuild\"`",
		"	OpenStdin bool `json:\"OpenStdin,omitempty\"`",
		"	QuotaID string `json:\"QuotaID,omitempty\"`",
		"	Rich bool `json:\"Rich,omitempty\"`",
		"	RichMode string `json:\"RichMode,omitempty\"`",
		"	Shell []string `json:\"Shell\"`",
		"	SpecAnnotation map[string]string `json:\"SpecAnnotation,omitempty\"`",
		"	StdinOnce bool `json:\"StdinOnce,omitempty\"`",
		"	StopSignal string `json:\"StopSignal,omitempty\"`",
		"	StopTimeout *int64 `json:\"StopTimeout,omitempty\"`",
		"	Tty bool `json:\"Tty,omitempty\"`",
		"	User string `json:\"User,omitempty\"`",
		"	Volumes map[string]interface{} `json:\"Volumes,omitempty\"`",
		"	WorkingDir string `json:\"WorkingDir,omitempty\"`",
		`	HostConfig struct {`,
		"		AutoRemove bool `json:\"AutoRemove,omitempty\"`",
		"		Binds []string `json:\"Binds\"`",
		"		CapAdd []string `json:\"CapAdd\"`",
		"		CapDrop []string `json:\"CapDrop\"`",
		"		Cgroup string `json:\"Cgroup,omitempty\"`",
		"		ConsoleSize []*int64 `json:\"ConsoleSize\"`",
		"		ContainerIDFile string `json:\"ContainerIDFile,omitempty\"`",
		"		DNS []string `json:\"Dns\"`",
		"		DNSOptions []string `json:\"DnsOptions\"`",
		"		DNSSearch []string `json:\"DnsSearch\"`",
		"		EnableLxcfs bool `json:\"EnableLxcfs,omitempty\"`",
		"		ExtraHosts []string `json:\"ExtraHosts\"`",
		"		GroupAdd []string `json:\"GroupAdd\"`",
		"		InitScript string `json:\"InitScript,omitempty\"`",
		"		IpcMode string `json:\"IpcMode,omitempty\"`",
		"		Isolation string `json:\"Isolation,omitempty\"`",
		"		Links []string `json:\"Links\"`",
		"		LogConfig *ContainerCreateConfigHostConfigAO0LogConfig `json:\"LogConfig,omitempty\"`",
		"		NetworkMode string `json:\"NetworkMode,omitempty\"`",
		"		OomScoreAdj int64 `json:\"OomScoreAdj,omitempty\"`",
		"		PidMode string `json:\"PidMode,omitempty\"`",
		"		Privileged bool `json:\"Privileged,omitempty\"`",
		"		PublishAllPorts bool `json:\"PublishAllPorts,omitempty\"`",
		"		ReadonlyRootfs bool `json:\"ReadonlyRootfs,omitempty\"`",
		"		RestartPolicy *ContainerCreateConfigHostConfigAO0RestartPolicy `json:\"RestartPolicy,omitempty\"`",
		"		Rich bool `json:\"Rich,omitempty\"`",
		"		RichMode string `json:\"RichMode,omitempty\"`",
		"		Runtime string `json:\"Runtime,omitempty\"`",
		"		SecurityOpt []string `json:\"SecurityOpt\"`",
		"		ShmSize *int64 `json:\"ShmSize,omitempty\"`",
		"		StorageOpt map[string]string `json:\"StorageOpt,omitempty\"`",
		"		Sysctls map[string]string `json:\"Sysctls,omitempty\"`",
		"		Tmpfs map[string]string `json:\"Tmpfs,omitempty\"`",
		"		UTSMode string `json:\"UTSMode,omitempty\"`",
		"		UsernsMode string `json:\"UsernsMode,omitempty\"`",
		"		VolumeDriver string `json:\"VolumeDriver,omitempty\"`",
		"		VolumesFrom []string `json:\"VolumesFrom\"`",
		"		BlkioWeight uint16 `json:\"BlkioWeight,omitempty\"`",
		"		CgroupParent string `json:\"CgroupParent,omitempty\"`",
		"		CPUShares int64 `json:\"CpuShares,omitempty\"`",
		"		Memory int64 `json:\"Memory,omitempty\"`",
		"	} `json:\"HostConfig,omitempty\"`",
		"	NetworkingConfig *ContainerCreateConfigAO1NetworkingConfig `json:\"NetworkingConfig,omitempty\"`",
		`func (m *ContainerCreateConfig) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateExposedPorts(formats); err != nil {`,
		`	if err := m.validateHostname(formats); err != nil {`,
		`	if err := m.validateImage(formats); err != nil {`,
		`	if err := m.validateRichMode(formats); err != nil {`,
		`	if err := m.validateVolumes(formats); err != nil {`,
		`	if err := m.validateHostConfig(formats); err != nil {`,
		`	if err := m.validateNetworkingConfig(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`var containerCreateConfigExposedPortsValueEnum []interface{`,
		`	var res []interface{`,
		"	if err := json.Unmarshal([]byte(`[{}]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		containerCreateConfigExposedPortsValueEnum = append(containerCreateConfigExposedPortsValueEnum, v`,
		`func (m *ContainerCreateConfig) validateExposedPortsValueEnum(path, location string, value interface{}) error {`,
		`	if err := validate.EnumCase(path, location, value, containerCreateConfigExposedPortsValueEnum, true); err != nil {`,
		`func (m *ContainerCreateConfig) validateExposedPorts(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ExposedPorts) {`,
		`	for k := range m.ExposedPorts {`,
		`		if err := m.validateExposedPortsValueEnum("ExposedPorts"+"."+k, "body", m.ExposedPorts[k]); err != nil {`,
		`func (m *ContainerCreateConfig) validateHostname(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Hostname) {`,
		`	if err := validate.MinLength("Hostname", "body", m.Hostname.String(), 1); err != nil {`,
		`	if err := validate.FormatOf("Hostname", "body", "hostname", m.Hostname.String(), formats); err != nil {`,
		`func (m *ContainerCreateConfig) validateImage(formats strfmt.Registry) error {`,
		`	if err := validate.RequiredString("Image", "body", m.Image); err != nil {`,
		`var containerCreateConfigTypeRichModePropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"dumb-init\",\"sbin-init\",\"systemd\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		containerCreateConfigTypeRichModePropEnum = append(containerCreateConfigTypeRichModePropEnum, v`,
		`func (m *ContainerCreateConfig) validateRichModeEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, containerCreateConfigTypeRichModePropEnum, true); err != nil {`,
		`func (m *ContainerCreateConfig) validateRichMode(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.RichMode) {`,
		`	if err := m.validateRichModeEnum("RichMode", "body", m.RichMode); err != nil {`,
		`var containerCreateConfigVolumesValueEnum []interface{`,
		`	var res []interface{`,
		"	if err := json.Unmarshal([]byte(`[{}]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		containerCreateConfigVolumesValueEnum = append(containerCreateConfigVolumesValueEnum, v`,
		`func (m *ContainerCreateConfig) validateVolumesValueEnum(path, location string, value interface{}) error {`,
		`	if err := validate.EnumCase(path, location, value, containerCreateConfigVolumesValueEnum, true); err != nil {`,
		`func (m *ContainerCreateConfig) validateVolumes(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Volumes) {`,
		`	for k := range m.Volumes {`,
		`		if err := m.validateVolumesValueEnum("Volumes"+"."+k, "body", m.Volumes[k]); err != nil {`,
		`func (m *ContainerCreateConfig) validateHostConfig(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.HostConfig) {`,
		`	iConsoleSizeSize := int64(len(m.HostConfig.ConsoleSize)`,
		`	if err := validate.MinItems("HostConfig"+"."+"ConsoleSize", "body", iConsoleSizeSize, 2); err != nil {`,
		`	if err := validate.MaxItems("HostConfig"+"."+"ConsoleSize", "body", iConsoleSizeSize, 2); err != nil {`,
		`	for i := 0; i < len(m.HostConfig.ConsoleSize); i++ {`,
		// do we need... ?
		`		if swag.IsZero(m.HostConfig.ConsoleSize[i]) {`,
		`		if err := validate.MinimumInt("HostConfig"+"."+"ConsoleSize"+"."+strconv.Itoa(i), "body", *m.HostConfig.ConsoleSize[i], 0, false); err != nil {`,
		// TODO: enum if anonymous allOf is not honored (missing func)
		// => will do that with Enum refactoring
		`	if err := m.validateIsolationEnum("HostConfig"+"."+"Isolation", "body", m.HostConfig.Isolation); err != nil {`,
		`	if m.HostConfig.LogConfig != nil {`,
		`		if err := m.HostConfig.LogConfig.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("HostConfig" + "." + "LogConfig"`,
		`	if err := validate.MinimumInt("HostConfig"+"."+"OomScoreAdj", "body", m.HostConfig.OomScoreAdj, -1000, false); err != nil {`,
		`	if err := validate.MaximumInt("HostConfig"+"."+"OomScoreAdj", "body", m.HostConfig.OomScoreAdj, 1000, false); err != nil {`,
		`	if m.HostConfig.RestartPolicy != nil {`,
		`		if err := m.HostConfig.RestartPolicy.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("HostConfig" + "." + "RestartPolicy"`,
		`	if err := m.validateRichModeEnum("HostConfig"+"."+"RichMode", "body", m.HostConfig.RichMode); err != nil {`,
		`	if err := validate.MinimumInt("HostConfig"+"."+"ShmSize", "body", *m.HostConfig.ShmSize, 0, false); err != nil {`,
		`	if err := validate.MinimumUint("HostConfig"+"."+"BlkioWeight", "body", uint64(m.HostConfig.BlkioWeight), 0, false); err != nil {`,
		`	if err := validate.MaximumUint("HostConfig"+"."+"BlkioWeight", "body", uint64(m.HostConfig.BlkioWeight), 1000, false); err != nil {`,
		`func (m *ContainerCreateConfig) validateNetworkingConfig(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.NetworkingConfig) {`,
		`	if m.NetworkingConfig != nil {`,
		`		if err := m.NetworkingConfig.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("NetworkingConfig"`,
		`type ContainerCreateConfigAO1NetworkingConfig struct {`,
		"	EndpointsConfig map[string]ContainerCreateConfigAO1NetworkingConfigEndpointsConfigAnon `json:\"EndpointsConfig,omitempty\"`",
		`func (m *ContainerCreateConfigAO1NetworkingConfig) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateEndpointsConfig(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ContainerCreateConfigAO1NetworkingConfig) validateEndpointsConfig(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.EndpointsConfig) {`,
		`	for k := range m.EndpointsConfig {`,
		`		if swag.IsZero(m.EndpointsConfig[k]) {`,
		`		if val, ok := m.EndpointsConfig[k]; ok {`,
		// NOTE: fixed incorrect IsNullable status in map element
		`				if err := val.Validate(formats); err != nil {`,
		`type ContainerCreateConfigAO1NetworkingConfigEndpointsConfigAnon struct {`,
		"	Aliases []string `json:\"Aliases\"`",
		"	DriverOpts map[string]string `json:\"DriverOpts,omitempty\"`",
		"	EndpointID string `json:\"EndpointID,omitempty\"`",
		"	Gateway string `json:\"Gateway,omitempty\"`",
		"	GlobalIPV6Address string `json:\"GlobalIPv6Address,omitempty\"`",
		"	GlobalIPV6PrefixLen int64 `json:\"GlobalIPv6PrefixLen,omitempty\"`",
		"	IPAMConfig *ContainerCreateConfigAO1NetworkingConfigEndpointsConfigAnonIPAMConfig `json:\"IPAMConfig,omitempty\"`",
		"	IPAddress string `json:\"IPAddress,omitempty\"`",
		"	IPPrefixLen int64 `json:\"IPPrefixLen,omitempty\"`",
		"	IPV6Gateway string `json:\"IPv6Gateway,omitempty\"`",
		"	Links []string `json:\"Links\"`",
		"	MacAddress string `json:\"MacAddress,omitempty\"`",
		"	NetworkID string `json:\"NetworkID,omitempty\"`",
		`func (m *ContainerCreateConfigAO1NetworkingConfigEndpointsConfigAnon) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateIPAMConfig(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ContainerCreateConfigAO1NetworkingConfigEndpointsConfigAnon) validateIPAMConfig(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.IPAMConfig) {`,
		`	if m.IPAMConfig != nil {`,
		`		if err := m.IPAMConfig.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("IPAMConfig"`,
		`type ContainerCreateConfigAO1NetworkingConfigEndpointsConfigAnonIPAMConfig struct {`,
		"	IPV4Address string `json:\"IPv4Address,omitempty\"`",
		"	IPV6Address string `json:\"IPv6Address,omitempty\"`",
		"	LinkLocalIPs []string `json:\"LinkLocalIPs\"`",
		`func (m *ContainerCreateConfigAO1NetworkingConfigEndpointsConfigAnonIPAMConfig) Validate(formats strfmt.Registry) error {`,
		`		return errors.CompositeValidationError(res...`,
		`type ContainerCreateConfigHostConfigAO0LogConfig struct {`,
		"	Config map[string]string `json:\"Config,omitempty\"`",
		"	Type string `json:\"Type,omitempty\"`",
		`func (m *ContainerCreateConfigHostConfigAO0LogConfig) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateType(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`var containerCreateConfigHostConfigAO0LogConfigTypeTypePropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"json-file\",\"syslog\",\"journald\",\"gelf\",\"fluentd\",\"awslogs\",\"splunk\",\"etwlogs\",\"none\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		containerCreateConfigHostConfigAO0LogConfigTypeTypePropEnum = append(containerCreateConfigHostConfigAO0LogConfigTypeTypePropEnum, v`,
		`	ContainerCreateConfigHostConfigAO0LogConfigTypeJSONDashFile string = "json-file"`,
		`	ContainerCreateConfigHostConfigAO0LogConfigTypeSyslog string = "syslog"`,
		`	ContainerCreateConfigHostConfigAO0LogConfigTypeJournald string = "journald"`,
		`	ContainerCreateConfigHostConfigAO0LogConfigTypeGelf string = "gelf"`,
		`	ContainerCreateConfigHostConfigAO0LogConfigTypeFluentd string = "fluentd"`,
		`	ContainerCreateConfigHostConfigAO0LogConfigTypeAwslogs string = "awslogs"`,
		`	ContainerCreateConfigHostConfigAO0LogConfigTypeSplunk string = "splunk"`,
		`	ContainerCreateConfigHostConfigAO0LogConfigTypeEtwlogs string = "etwlogs"`,
		`	ContainerCreateConfigHostConfigAO0LogConfigTypeNone string = "none"`,
		`func (m *ContainerCreateConfigHostConfigAO0LogConfig) validateTypeEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, containerCreateConfigHostConfigAO0LogConfigTypeTypePropEnum, true); err != nil {`,
		`func (m *ContainerCreateConfigHostConfigAO0LogConfig) validateType(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Type) {`,
		`	if err := m.validateTypeEnum("HostConfig"+"."+"LogConfig"+"."+"Type", "body", m.Type); err != nil {`,
		`type ContainerCreateConfigHostConfigAO0RestartPolicy struct {`,
		"	MaximumRetryCount int64 `json:\"MaximumRetryCount,omitempty\"`",
		"	Name string `json:\"Name,omitempty\"`",
		`func (m *ContainerCreateConfigHostConfigAO0RestartPolicy) Validate(formats strfmt.Registry) error {`,
		`		return errors.CompositeValidationError(res...`,
	}, []string{
		// not expected
		`			if val != nil {`,
	},
		// output in log
		noLines,
		noLines)

	// load expectations for model: resources.go
	flattenRun.AddExpectations("resources.go", []string{
		`type Resources struct {`,
		"	BlkioWeight uint16 `json:\"BlkioWeight,omitempty\"`",
		"	CgroupParent string `json:\"CgroupParent,omitempty\"`",
		"	CPUShares int64 `json:\"CpuShares,omitempty\"`",
		"	Memory int64 `json:\"Memory,omitempty\"`",
		`func (m *Resources) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateBlkioWeight(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *Resources) validateBlkioWeight(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.BlkioWeight) {`,
		`	if err := validate.MinimumUint("BlkioWeight", "body", uint64(m.BlkioWeight), 0, false); err != nil {`,
		`	if err := validate.MaximumUint("BlkioWeight", "body", uint64(m.BlkioWeight), 1000, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("resources.go", flattenRun.ExpectedFor("Resources").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: networking_config.go
	flattenRun.AddExpectations("networking_config.go", []string{
		`type NetworkingConfig struct {`,
		// maps are now simple types
		"	EndpointsConfig map[string]*EndpointSettings `json:\"EndpointsConfig,omitempty\"`",
		`func (m *NetworkingConfig) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateEndpointsConfig(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NetworkingConfig) validateEndpointsConfig(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.EndpointsConfig) {`,
		`       for k := range m.EndpointsConfig {`,
		`	if err := validate.Required("EndpointsConfig"+"."+k, "body", m.EndpointsConfig[k]); err != nil {`,
		`       	if val, ok := m.EndpointsConfig[k]; ok {`,
		`          		if err := val.Validate(formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// NOTE(fredbi): maps are now simple types: this definition disappears

	// load expectations for model: endpoint_settings.go
	flattenRun.AddExpectations("endpoint_settings.go", []string{
		`type EndpointSettings struct {`,
		"	Aliases []string `json:\"Aliases\"`",
		"	DriverOpts map[string]string `json:\"DriverOpts,omitempty\"`",
		"	EndpointID string `json:\"EndpointID,omitempty\"`",
		"	Gateway string `json:\"Gateway,omitempty\"`",
		"	GlobalIPV6Address string `json:\"GlobalIPv6Address,omitempty\"`",
		"	GlobalIPV6PrefixLen int64 `json:\"GlobalIPv6PrefixLen,omitempty\"`",
		"	IPAMConfig *EndpointIPAMConfig `json:\"IPAMConfig,omitempty\"`",
		"	IPAddress string `json:\"IPAddress,omitempty\"`",
		"	IPPrefixLen int64 `json:\"IPPrefixLen,omitempty\"`",
		"	IPV6Gateway string `json:\"IPv6Gateway,omitempty\"`",
		"	Links []string `json:\"Links\"`",
		"	MacAddress string `json:\"MacAddress,omitempty\"`",
		"	NetworkID string `json:\"NetworkID,omitempty\"`",
		`func (m *EndpointSettings) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateIPAMConfig(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *EndpointSettings) validateIPAMConfig(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.IPAMConfig) {`,
		`	if m.IPAMConfig != nil {`,
		`		if err := m.IPAMConfig.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("IPAMConfig"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("endpoint_settings.go", []string{
		`type EndpointSettings struct {`,
		"	Aliases []string `json:\"Aliases\"`",
		"	DriverOpts map[string]string `json:\"DriverOpts,omitempty\"`",
		"	EndpointID string `json:\"EndpointID,omitempty\"`",
		"	Gateway string `json:\"Gateway,omitempty\"`",
		"	GlobalIPV6Address string `json:\"GlobalIPv6Address,omitempty\"`",
		"	GlobalIPV6PrefixLen int64 `json:\"GlobalIPv6PrefixLen,omitempty\"`",
		"	IPAMConfig *EndpointSettingsIPAMConfig `json:\"IPAMConfig,omitempty\"`",
		"	IPAddress string `json:\"IPAddress,omitempty\"`",
		"	IPPrefixLen int64 `json:\"IPPrefixLen,omitempty\"`",
		"	IPV6Gateway string `json:\"IPv6Gateway,omitempty\"`",
		"	Links []string `json:\"Links\"`",
		"	MacAddress string `json:\"MacAddress,omitempty\"`",
		"	NetworkID string `json:\"NetworkID,omitempty\"`",
		`func (m *EndpointSettings) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateIPAMConfig(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *EndpointSettings) validateIPAMConfig(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.IPAMConfig) {`,
		`	if m.IPAMConfig != nil {`,
		`		if err := m.IPAMConfig.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("IPAMConfig"`,
		`type EndpointSettingsIPAMConfig struct {`,
		"	IPV4Address string `json:\"IPv4Address,omitempty\"`",
		"	IPV6Address string `json:\"IPv6Address,omitempty\"`",
		"	LinkLocalIPs []string `json:\"LinkLocalIPs\"`",
		`func (m *EndpointSettingsIPAMConfig) Validate(formats strfmt.Registry) error {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: networking_config_endpoints_config.go
	// NOTE(fredbi): maps are now simple types - this definition disappears
}
func initFixtureSimpleAllOf() {
	// testing ../fixtures/bugs/1487/fixture-simple-allOf.yaml with flatten and expand (--skip-flatten)

	/* we test various composition combinations, including nested, and nested isolated with a properties (e.g. issue #1479) */

	f := newModelFixture("../fixtures/bugs/1487/fixture-simple-allOf.yaml", "fixture for nested allOf with ref")
	flattenRun := f.AddRun(false)
	expandRun := f.AddRun(true)

	// load expectations for model: not_really_composed_thing_all_of0.go
	flattenRun.AddExpectations("not_really_composed_thing_all_of0.go", []string{
		`type NotReallyComposedThingAllOf0 struct {`,
		"	Prop0 strfmt.UUID `json:\"prop0,omitempty\"`",
		`func (m *NotReallyComposedThingAllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp0(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NotReallyComposedThingAllOf0) validateProp0(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop0) {`,
		`	if err := validate.FormatOf("prop0", "body", "uuid", m.Prop0.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: not_really_composed_thing.go
	expandRun.AddExpectations("not_really_composed_thing.go", []string{
		`type NotReallyComposedThing struct {`,
		"	Prop0 strfmt.UUID `json:\"prop0,omitempty\"`",
		`func (m *NotReallyComposedThing) UnmarshalJSON(raw []byte) error {`,
		`	var dataAO0 struct {`,
		"		Prop0 strfmt.UUID `json:\"prop0,omitempty\"`",
		`	if err := swag.ReadJSON(raw, &dataAO0); err != nil {`,
		`	m.Prop0 = dataAO0.Prop0`,
		`func (m NotReallyComposedThing) MarshalJSON() ([]byte, error) {`,
		`	_parts := make([][]byte, 0, 1`,
		`	var dataAO0 struct {`,
		"		Prop0 strfmt.UUID `json:\"prop0,omitempty\"`",
		`	dataAO0.Prop0 = m.Prop0`,
		`	jsonDataAO0, errAO0 := swag.WriteJSON(dataAO0`,
		`	if errAO0 != nil {`,
		`		return nil, errAO0`,
		`	_parts = append(_parts, jsonDataAO0`,
		`	return swag.ConcatJSON(_parts...), nil`,
		`func (m *NotReallyComposedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp0(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NotReallyComposedThing) validateProp0(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop0) {`,
		`	if err := validate.FormatOf("prop0", "body", "uuid", m.Prop0.String(), formats); err != nil {`,
		`func (m *NotReallyComposedThing) MarshalBinary() ([]byte, error) {`,
		`	if m == nil {`,
		`		return nil, nil`,
		`	return swag.WriteJSON(m`,
		`func (m *NotReallyComposedThing) UnmarshalBinary(b []byte) error {`,
		`	var res NotReallyComposedThing`,
		`	if err := swag.ReadJSON(b, &res); err != nil {`,
		`	*m = res`,
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: simple_nested_object_all_of1.go
	flattenRun.AddExpectations("simple_nested_object_all_of1.go", []string{
		`type SimpleNestedObjectAllOf1 struct {`,
		"	Prop3 strfmt.UUID `json:\"prop3,omitempty\"`",
		`func (m *SimpleNestedObjectAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp3(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *SimpleNestedObjectAllOf1) validateProp3(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop3) {`,
		`	if err := validate.FormatOf("prop3", "body", "uuid", m.Prop3.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: break_nested_object_all_of1_prop7.go
	flattenRun.AddExpectations("break_nested_object_all_of1_prop7.go", []string{
		`type BreakNestedObjectAllOf1Prop7 struct {`,
		`	BreakNestedObjectAllOf1Prop7AllOf0`,
		`	BreakNestedObjectAllOf1Prop7AllOf1`,
		`func (m *BreakNestedObjectAllOf1Prop7) Validate(formats strfmt.Registry) error {`,
		`	if err := m.BreakNestedObjectAllOf1Prop7AllOf0.Validate(formats); err != nil {`,
		`	if err := m.BreakNestedObjectAllOf1Prop7AllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: composed_thing.go
	flattenRun.AddExpectations("composed_thing.go", []string{
		`type ComposedThing struct {`,
		`	ComposedThingAllOf0`,
		`	ComposedThingAllOf1`,
		`func (m *ComposedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.ComposedThingAllOf0.Validate(formats); err != nil {`,
		`	if err := m.ComposedThingAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: composed_thing.go
	expandRun.AddExpectations("composed_thing.go", []string{
		`type ComposedThing struct {`,
		"	Prop1 strfmt.UUID `json:\"prop1,omitempty\"`",
		"	Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		`func (m *ComposedThing) UnmarshalJSON(raw []byte) error {`,
		`	var dataAO0 struct {`,
		"		Prop1 strfmt.UUID `json:\"prop1,omitempty\"`",
		`	if err := swag.ReadJSON(raw, &dataAO0); err != nil {`,
		`	m.Prop1 = dataAO0.Prop1`,
		`	var dataAO1 struct {`,
		"		Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		`	if err := swag.ReadJSON(raw, &dataAO1); err != nil {`,
		`	m.Prop2 = dataAO1.Prop2`,
		`func (m ComposedThing) MarshalJSON() ([]byte, error) {`,
		`	_parts := make([][]byte, 0, 2`,
		`	var dataAO0 struct {`,
		"		Prop1 strfmt.UUID `json:\"prop1,omitempty\"`",
		`	dataAO0.Prop1 = m.Prop1`,
		`	jsonDataAO0, errAO0 := swag.WriteJSON(dataAO0`,
		`	if errAO0 != nil {`,
		`		return nil, errAO0`,
		`	_parts = append(_parts, jsonDataAO0`,
		`	var dataAO1 struct {`,
		"		Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		`	dataAO1.Prop2 = m.Prop2`,
		`	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1`,
		`	if errAO1 != nil {`,
		`		return nil, errAO1`,
		`	_parts = append(_parts, jsonDataAO1`,
		`	return swag.ConcatJSON(_parts...), nil`,
		`func (m *ComposedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`	if err := m.validateProp2(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ComposedThing) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	if err := validate.FormatOf("prop1", "body", "uuid", m.Prop1.String(), formats); err != nil {`,
		`func (m *ComposedThing) validateProp2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop2) {`,
		`	if err := validate.FormatOf("prop2", "body", "uuid", m.Prop2.String(), formats); err != nil {`,
		`func (m *ComposedThing) MarshalBinary() ([]byte, error) {`,
		`	if m == nil {`,
		`		return nil, nil`,
		`	return swag.WriteJSON(m`,
		`func (m *ComposedThing) UnmarshalBinary(b []byte) error {`,
		`	var res ComposedThing`,
		`	if err := swag.ReadJSON(b, &res); err != nil {`,
		`	*m = res`,
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: break_nested_object.go
	flattenRun.AddExpectations("break_nested_object.go", []string{
		`type BreakNestedObject struct {`,
		`	SimpleNestedObject`,
		`	BreakNestedObjectAllOf1`,
		`func (m *BreakNestedObject) Validate(formats strfmt.Registry) error {`,
		`	if err := m.SimpleNestedObject.Validate(formats); err != nil {`,
		`	if err := m.BreakNestedObjectAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: break_nested_object.go
	expandRun.AddExpectations("break_nested_object.go", []string{
		`type BreakNestedObject struct {`,
		`	BreakNestedObjectAllOf0`,
		"	Prop6 strfmt.UUID `json:\"prop6,omitempty\"`",
		`	Prop7 struct {`,
		"		Prop8 int64 `json:\"prop8,omitempty\"`",
		"		Prop9 int64 `json:\"prop9,omitempty\"`",
		"	} `json:\"prop7,omitempty\"`",
		`func (m *BreakNestedObject) UnmarshalJSON(raw []byte) error {`,
		`	var aO0 BreakNestedObjectAllOf0`,
		`	if err := swag.ReadJSON(raw, &aO0); err != nil {`,
		`	m.BreakNestedObjectAllOf0 = aO0`,
		`	var dataAO1 struct {`,
		"		Prop6 strfmt.UUID `json:\"prop6,omitempty\"`",
		`		Prop7 struct {`,
		"			Prop8 int64 `json:\"prop8,omitempty\"`",
		"			Prop9 int64 `json:\"prop9,omitempty\"`",
		"		} `json:\"prop7,omitempty\"`",
		`	if err := swag.ReadJSON(raw, &dataAO1); err != nil {`,
		`	m.Prop6 = dataAO1.Prop6`,
		`	m.Prop7 = dataAO1.Prop7`,
		`func (m BreakNestedObject) MarshalJSON() ([]byte, error) {`,
		`	_parts := make([][]byte, 0, 2`,
		`	aO0, err := swag.WriteJSON(m.BreakNestedObjectAllOf0`,
		`	if err != nil {`,
		`		return nil, err`,
		`	_parts = append(_parts, aO0`,
		`	var dataAO1 struct {`,
		"		Prop6 strfmt.UUID `json:\"prop6,omitempty\"`",
		`		Prop7 struct {`,
		"			Prop8 int64 `json:\"prop8,omitempty\"`",
		"			Prop9 int64 `json:\"prop9,omitempty\"`",
		"		} `json:\"prop7,omitempty\"`",
		`	dataAO1.Prop6 = m.Prop6`,
		`	dataAO1.Prop7 = m.Prop7`,
		`	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1`,
		`	if errAO1 != nil {`,
		`		return nil, errAO1`,
		`	_parts = append(_parts, jsonDataAO1`,
		`	return swag.ConcatJSON(_parts...), nil`,
		`func (m *BreakNestedObject) Validate(formats strfmt.Registry) error {`,
		`	if err := m.BreakNestedObjectAllOf0.Validate(formats); err != nil {`,
		`	if err := m.validateProp6(formats); err != nil {`,
		`	if err := m.validateProp7(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *BreakNestedObject) validateProp6(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop6) {`,
		`	if err := validate.FormatOf("prop6", "body", "uuid", m.Prop6.String(), formats); err != nil {`,
		`func (m *BreakNestedObject) validateProp7(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop7) {`,
		`	if err := validate.MinimumInt("prop7"+"."+"prop8", "body", m.Prop7.Prop8, 12, false); err != nil {`,
		`	if err := validate.MaximumInt("prop7"+"."+"prop9", "body", m.Prop7.Prop9, 12, false); err != nil {`,
		`func (m *BreakNestedObject) MarshalBinary() ([]byte, error) {`,
		`	if m == nil {`,
		`		return nil, nil`,
		`	return swag.WriteJSON(m`,
		`func (m *BreakNestedObject) UnmarshalBinary(b []byte) error {`,
		`	var res BreakNestedObject`,
		`	if err := swag.ReadJSON(b, &res); err != nil {`,
		`	*m = res`,
		`type BreakNestedObjectAllOf0 struct {`,
		`	BreakNestedObjectAllOf0AllOf0`,
		"	Prop3 strfmt.UUID `json:\"prop3,omitempty\"`",
		`func (m *BreakNestedObjectAllOf0) UnmarshalJSON(raw []byte) error {`,
		`	var aO0 BreakNestedObjectAllOf0AllOf0`,
		`	if err := swag.ReadJSON(raw, &aO0); err != nil {`,
		`	m.BreakNestedObjectAllOf0AllOf0 = aO0`,
		`	var dataAO1 struct {`,
		"		Prop3 strfmt.UUID `json:\"prop3,omitempty\"`",
		`	if err := swag.ReadJSON(raw, &dataAO1); err != nil {`,
		`	m.Prop3 = dataAO1.Prop3`,
		`func (m BreakNestedObjectAllOf0) MarshalJSON() ([]byte, error) {`,
		`	_parts := make([][]byte, 0, 2`,
		`	aO0, err := swag.WriteJSON(m.BreakNestedObjectAllOf0AllOf0`,
		`	if err != nil {`,
		`		return nil, err`,
		`	_parts = append(_parts, aO0`,
		`	var dataAO1 struct {`,
		"		Prop3 strfmt.UUID `json:\"prop3,omitempty\"`",
		`	dataAO1.Prop3 = m.Prop3`,
		`	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1`,
		`	if errAO1 != nil {`,
		`		return nil, errAO1`,
		`	_parts = append(_parts, jsonDataAO1`,
		`	return swag.ConcatJSON(_parts...), nil`,
		`func (m *BreakNestedObjectAllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.BreakNestedObjectAllOf0AllOf0.Validate(formats); err != nil {`,
		`	if err := m.validateProp3(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *BreakNestedObjectAllOf0) validateProp3(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop3) {`,
		`	if err := validate.FormatOf("prop3", "body", "uuid", m.Prop3.String(), formats); err != nil {`,
		`func (m *BreakNestedObjectAllOf0) MarshalBinary() ([]byte, error) {`,
		`	if m == nil {`,
		`		return nil, nil`,
		`	return swag.WriteJSON(m`,
		`func (m *BreakNestedObjectAllOf0) UnmarshalBinary(b []byte) error {`,
		`	var res BreakNestedObjectAllOf0`,
		`	if err := swag.ReadJSON(b, &res); err != nil {`,
		`	*m = res`,
		`type BreakNestedObjectAllOf0AllOf0 struct {`,
		"	Prop1 strfmt.UUID `json:\"prop1,omitempty\"`",
		"	Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		`func (m *BreakNestedObjectAllOf0AllOf0) UnmarshalJSON(raw []byte) error {`,
		`	var dataAO0 struct {`,
		"		Prop1 strfmt.UUID `json:\"prop1,omitempty\"`",
		`	if err := swag.ReadJSON(raw, &dataAO0); err != nil {`,
		`	m.Prop1 = dataAO0.Prop1`,
		`	var dataAO1 struct {`,
		"		Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		`	if err := swag.ReadJSON(raw, &dataAO1); err != nil {`,
		`	m.Prop2 = dataAO1.Prop2`,
		`func (m BreakNestedObjectAllOf0AllOf0) MarshalJSON() ([]byte, error) {`,
		`	_parts := make([][]byte, 0, 2`,
		`	var dataAO0 struct {`,
		"		Prop1 strfmt.UUID `json:\"prop1,omitempty\"`",
		`	dataAO0.Prop1 = m.Prop1`,
		`	jsonDataAO0, errAO0 := swag.WriteJSON(dataAO0`,
		`	if errAO0 != nil {`,
		`		return nil, errAO0`,
		`	_parts = append(_parts, jsonDataAO0`,
		`	var dataAO1 struct {`,
		"		Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		`	dataAO1.Prop2 = m.Prop2`,
		`	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1`,
		`	if errAO1 != nil {`,
		`		return nil, errAO1`,
		`	_parts = append(_parts, jsonDataAO1`,
		`	return swag.ConcatJSON(_parts...), nil`,
		`func (m *BreakNestedObjectAllOf0AllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`	if err := m.validateProp2(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *BreakNestedObjectAllOf0AllOf0) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	if err := validate.FormatOf("prop1", "body", "uuid", m.Prop1.String(), formats); err != nil {`,
		`func (m *BreakNestedObjectAllOf0AllOf0) validateProp2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop2) {`,
		`	if err := validate.FormatOf("prop2", "body", "uuid", m.Prop2.String(), formats); err != nil {`,
		`func (m *BreakNestedObjectAllOf0AllOf0) MarshalBinary() ([]byte, error) {`,
		`	if m == nil {`,
		`		return nil, nil`,
		`	return swag.WriteJSON(m`,
		`func (m *BreakNestedObjectAllOf0AllOf0) UnmarshalBinary(b []byte) error {`,
		`	var res BreakNestedObjectAllOf0AllOf0`,
		`	if err := swag.ReadJSON(b, &res); err != nil {`,
		`	*m = res`,
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: deep_nested_object_all_of1_all_of1.go
	flattenRun.AddExpectations("deep_nested_object_all_of1_all_of1.go", []string{
		`type DeepNestedObjectAllOf1AllOf1 struct {`,
		"	Prop5 strfmt.Date `json:\"prop5,omitempty\"`",
		`func (m *DeepNestedObjectAllOf1AllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp5(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *DeepNestedObjectAllOf1AllOf1) validateProp5(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop5) {`,
		`	if err := validate.FormatOf("prop5", "body", "date", m.Prop5.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: deep_nested_object.go
	flattenRun.AddExpectations("deep_nested_object.go", []string{
		`type DeepNestedObject struct {`,
		`	SimpleNestedObject`,
		`	DeepNestedObjectAllOf1`,
		`func (m *DeepNestedObject) Validate(formats strfmt.Registry) error {`,
		`	if err := m.SimpleNestedObject.Validate(formats); err != nil {`,
		`	if err := m.DeepNestedObjectAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: deep_nested_object.go
	expandRun.AddExpectations("deep_nested_object.go", []string{
		`type DeepNestedObject struct {`,
		`	DeepNestedObjectAllOf0`,
		`	DeepNestedObjectAllOf1`,
		`func (m *DeepNestedObject) UnmarshalJSON(raw []byte) error {`,
		`	var aO0 DeepNestedObjectAllOf0`,
		`	if err := swag.ReadJSON(raw, &aO0); err != nil {`,
		`	m.DeepNestedObjectAllOf0 = aO0`,
		`	var aO1 DeepNestedObjectAllOf1`,
		`	if err := swag.ReadJSON(raw, &aO1); err != nil {`,
		`	m.DeepNestedObjectAllOf1 = aO1`,
		`func (m DeepNestedObject) MarshalJSON() ([]byte, error) {`,
		`	_parts := make([][]byte, 0, 2`,
		`	aO0, err := swag.WriteJSON(m.DeepNestedObjectAllOf0`,
		`	if err != nil {`,
		`		return nil, err`,
		`	_parts = append(_parts, aO0`,
		`	aO1, err := swag.WriteJSON(m.DeepNestedObjectAllOf1`,
		`	if err != nil {`,
		`		return nil, err`,
		`	_parts = append(_parts, aO1`,
		`	return swag.ConcatJSON(_parts...), nil`,
		`func (m *DeepNestedObject) Validate(formats strfmt.Registry) error {`,
		`	if err := m.DeepNestedObjectAllOf0.Validate(formats); err != nil {`,
		`	if err := m.DeepNestedObjectAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *DeepNestedObject) MarshalBinary() ([]byte, error) {`,
		`	if m == nil {`,
		`		return nil, nil`,
		`	return swag.WriteJSON(m`,
		`func (m *DeepNestedObject) UnmarshalBinary(b []byte) error {`,
		`	var res DeepNestedObject`,
		`	if err := swag.ReadJSON(b, &res); err != nil {`,
		`	*m = res`,
		`type DeepNestedObjectAllOf0 struct {`,
		`	DeepNestedObjectAllOf0AllOf0`,
		"	Prop3 strfmt.UUID `json:\"prop3,omitempty\"`",
		`func (m *DeepNestedObjectAllOf0) UnmarshalJSON(raw []byte) error {`,
		`	var aO0 DeepNestedObjectAllOf0AllOf0`,
		`	if err := swag.ReadJSON(raw, &aO0); err != nil {`,
		`	m.DeepNestedObjectAllOf0AllOf0 = aO0`,
		`	var dataAO1 struct {`,
		"		Prop3 strfmt.UUID `json:\"prop3,omitempty\"`",
		`	if err := swag.ReadJSON(raw, &dataAO1); err != nil {`,
		`	m.Prop3 = dataAO1.Prop3`,
		`func (m DeepNestedObjectAllOf0) MarshalJSON() ([]byte, error) {`,
		`	_parts := make([][]byte, 0, 2`,
		`	aO0, err := swag.WriteJSON(m.DeepNestedObjectAllOf0AllOf0`,
		`	if err != nil {`,
		`		return nil, err`,
		`	_parts = append(_parts, aO0`,
		`	var dataAO1 struct {`,
		"		Prop3 strfmt.UUID `json:\"prop3,omitempty\"`",
		`	dataAO1.Prop3 = m.Prop3`,
		`	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1`,
		`	if errAO1 != nil {`,
		`		return nil, errAO1`,
		`	_parts = append(_parts, jsonDataAO1`,
		`	return swag.ConcatJSON(_parts...), nil`,
		`func (m *DeepNestedObjectAllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.DeepNestedObjectAllOf0AllOf0.Validate(formats); err != nil {`,
		`	if err := m.validateProp3(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *DeepNestedObjectAllOf0) validateProp3(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop3) {`,
		`	if err := validate.FormatOf("prop3", "body", "uuid", m.Prop3.String(), formats); err != nil {`,
		`func (m *DeepNestedObjectAllOf0) MarshalBinary() ([]byte, error) {`,
		`	if m == nil {`,
		`		return nil, nil`,
		`	return swag.WriteJSON(m`,
		`func (m *DeepNestedObjectAllOf0) UnmarshalBinary(b []byte) error {`,
		`	var res DeepNestedObjectAllOf0`,
		`	if err := swag.ReadJSON(b, &res); err != nil {`,
		`	*m = res`,
		`type DeepNestedObjectAllOf0AllOf0 struct {`,
		"	Prop1 strfmt.UUID `json:\"prop1,omitempty\"`",
		"	Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		`func (m *DeepNestedObjectAllOf0AllOf0) UnmarshalJSON(raw []byte) error {`,
		`	var dataAO0 struct {`,
		"		Prop1 strfmt.UUID `json:\"prop1,omitempty\"`",
		`	if err := swag.ReadJSON(raw, &dataAO0); err != nil {`,
		`	m.Prop1 = dataAO0.Prop1`,
		`	var dataAO1 struct {`,
		"		Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		`	if err := swag.ReadJSON(raw, &dataAO1); err != nil {`,
		`	m.Prop2 = dataAO1.Prop2`,
		`func (m DeepNestedObjectAllOf0AllOf0) MarshalJSON() ([]byte, error) {`,
		`	_parts := make([][]byte, 0, 2`,
		`	var dataAO0 struct {`,
		"		Prop1 strfmt.UUID `json:\"prop1,omitempty\"`",
		`	dataAO0.Prop1 = m.Prop1`,
		`	jsonDataAO0, errAO0 := swag.WriteJSON(dataAO0`,
		`	if errAO0 != nil {`,
		`		return nil, errAO0`,
		`	_parts = append(_parts, jsonDataAO0`,
		`	var dataAO1 struct {`,
		"		Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		`	dataAO1.Prop2 = m.Prop2`,
		`	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1`,
		`	if errAO1 != nil {`,
		`		return nil, errAO1`,
		`	_parts = append(_parts, jsonDataAO1`,
		`	return swag.ConcatJSON(_parts...), nil`,
		`func (m *DeepNestedObjectAllOf0AllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`	if err := m.validateProp2(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *DeepNestedObjectAllOf0AllOf0) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	if err := validate.FormatOf("prop1", "body", "uuid", m.Prop1.String(), formats); err != nil {`,
		`func (m *DeepNestedObjectAllOf0AllOf0) validateProp2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop2) {`,
		`	if err := validate.FormatOf("prop2", "body", "uuid", m.Prop2.String(), formats); err != nil {`,
		`func (m *DeepNestedObjectAllOf0AllOf0) MarshalBinary() ([]byte, error) {`,
		`	if m == nil {`,
		`		return nil, nil`,
		`	return swag.WriteJSON(m`,
		`func (m *DeepNestedObjectAllOf0AllOf0) UnmarshalBinary(b []byte) error {`,
		`	var res DeepNestedObjectAllOf0AllOf0`,
		`	if err := swag.ReadJSON(b, &res); err != nil {`,
		`	*m = res`,
		`type DeepNestedObjectAllOf1 struct {`,
		"	Prop4 strfmt.UUID `json:\"prop4,omitempty\"`",
		"	Prop5 strfmt.Date `json:\"prop5,omitempty\"`",
		`func (m *DeepNestedObjectAllOf1) UnmarshalJSON(raw []byte) error {`,
		`	var dataAO0 struct {`,
		"		Prop4 strfmt.UUID `json:\"prop4,omitempty\"`",
		`	if err := swag.ReadJSON(raw, &dataAO0); err != nil {`,
		`	m.Prop4 = dataAO0.Prop4`,
		`	var dataAO1 struct {`,
		"		Prop5 strfmt.Date `json:\"prop5,omitempty\"`",
		`	if err := swag.ReadJSON(raw, &dataAO1); err != nil {`,
		`	m.Prop5 = dataAO1.Prop5`,
		`func (m DeepNestedObjectAllOf1) MarshalJSON() ([]byte, error) {`,
		`	_parts := make([][]byte, 0, 2`,
		`	var dataAO0 struct {`,
		"		Prop4 strfmt.UUID `json:\"prop4,omitempty\"`",
		`	dataAO0.Prop4 = m.Prop4`,
		`	jsonDataAO0, errAO0 := swag.WriteJSON(dataAO0`,
		`	if errAO0 != nil {`,
		`		return nil, errAO0`,
		`	_parts = append(_parts, jsonDataAO0`,
		`	var dataAO1 struct {`,
		"		Prop5 strfmt.Date `json:\"prop5,omitempty\"`",
		`	dataAO1.Prop5 = m.Prop5`,
		`	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1`,
		`	if errAO1 != nil {`,
		`		return nil, errAO1`,
		`	_parts = append(_parts, jsonDataAO1`,
		`	return swag.ConcatJSON(_parts...), nil`,
		`func (m *DeepNestedObjectAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp4(formats); err != nil {`,
		`	if err := m.validateProp5(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *DeepNestedObjectAllOf1) validateProp4(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop4) {`,
		`	if err := validate.FormatOf("prop4", "body", "uuid", m.Prop4.String(), formats); err != nil {`,
		`func (m *DeepNestedObjectAllOf1) validateProp5(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop5) {`,
		`	if err := validate.FormatOf("prop5", "body", "date", m.Prop5.String(), formats); err != nil {`,
		`func (m *DeepNestedObjectAllOf1) MarshalBinary() ([]byte, error) {`,
		`	if m == nil {`,
		`		return nil, nil`,
		`	return swag.WriteJSON(m`,
		`func (m *DeepNestedObjectAllOf1) UnmarshalBinary(b []byte) error {`,
		`	var res DeepNestedObjectAllOf1`,
		`	if err := swag.ReadJSON(b, &res); err != nil {`,
		`	*m = res`,
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: break_nested_object_all_of1.go
	flattenRun.AddExpectations("break_nested_object_all_of1.go", []string{
		`type BreakNestedObjectAllOf1 struct {`,
		"	Prop6 strfmt.UUID `json:\"prop6,omitempty\"`",
		"	Prop7 *BreakNestedObjectAllOf1Prop7 `json:\"prop7,omitempty\"`",
		`func (m *BreakNestedObjectAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp6(formats); err != nil {`,
		`	if err := m.validateProp7(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *BreakNestedObjectAllOf1) validateProp6(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop6) {`,
		`	if err := validate.FormatOf("prop6", "body", "uuid", m.Prop6.String(), formats); err != nil {`,
		`func (m *BreakNestedObjectAllOf1) validateProp7(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop7) {`,
		`	if m.Prop7 != nil {`,
		`		if err := m.Prop7.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("prop7"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: not_really_composed_thing.go
	flattenRun.AddExpectations("not_really_composed_thing.go", []string{
		`type NotReallyComposedThing struct {`,
		`	NotReallyComposedThingAllOf0`,
		`func (m *NotReallyComposedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.NotReallyComposedThingAllOf0.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("not_really_composed_thing.go", []string{
		`type NotReallyComposedThing struct {`,
		"	Prop0 strfmt.UUID `json:\"prop0,omitempty\"`",
		`func (m *NotReallyComposedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp0(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NotReallyComposedThing) validateProp0(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop0) {`,
		`	if err := validate.FormatOf("prop0", "body", "uuid", m.Prop0.String(), formats); err != nil {`,
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: simple_nested_object.go
	flattenRun.AddExpectations("simple_nested_object.go", []string{
		`type SimpleNestedObject struct {`,
		`	ComposedThing`,
		`	SimpleNestedObjectAllOf1`,
		`func (m *SimpleNestedObject) Validate(formats strfmt.Registry) error {`,
		`	if err := m.ComposedThing.Validate(formats); err != nil {`,
		`	if err := m.SimpleNestedObjectAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: simple_nested_object.go
	expandRun.AddExpectations("simple_nested_object.go", []string{
		`type SimpleNestedObject struct {`,
		`	SimpleNestedObjectAllOf0`,
		"	Prop3 strfmt.UUID `json:\"prop3,omitempty\"`",
		`func (m *SimpleNestedObject) UnmarshalJSON(raw []byte) error {`,
		`	var aO0 SimpleNestedObjectAllOf0`,
		`	if err := swag.ReadJSON(raw, &aO0); err != nil {`,
		`	m.SimpleNestedObjectAllOf0 = aO0`,
		`	var dataAO1 struct {`,
		"		Prop3 strfmt.UUID `json:\"prop3,omitempty\"`",
		`	if err := swag.ReadJSON(raw, &dataAO1); err != nil {`,
		`	m.Prop3 = dataAO1.Prop3`,
		`func (m SimpleNestedObject) MarshalJSON() ([]byte, error) {`,
		`	_parts := make([][]byte, 0, 2`,
		`	aO0, err := swag.WriteJSON(m.SimpleNestedObjectAllOf0`,
		`	if err != nil {`,
		`		return nil, err`,
		`	_parts = append(_parts, aO0`,
		`	var dataAO1 struct {`,
		"		Prop3 strfmt.UUID `json:\"prop3,omitempty\"`",
		`	dataAO1.Prop3 = m.Prop3`,
		`	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1`,
		`	if errAO1 != nil {`,
		`		return nil, errAO1`,
		`	_parts = append(_parts, jsonDataAO1`,
		`	return swag.ConcatJSON(_parts...), nil`,
		`func (m *SimpleNestedObject) Validate(formats strfmt.Registry) error {`,
		`	if err := m.SimpleNestedObjectAllOf0.Validate(formats); err != nil {`,
		`	if err := m.validateProp3(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *SimpleNestedObject) validateProp3(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop3) {`,
		`	if err := validate.FormatOf("prop3", "body", "uuid", m.Prop3.String(), formats); err != nil {`,
		`func (m *SimpleNestedObject) MarshalBinary() ([]byte, error) {`,
		`	if m == nil {`,
		`		return nil, nil`,
		`	return swag.WriteJSON(m`,
		`func (m *SimpleNestedObject) UnmarshalBinary(b []byte) error {`,
		`	var res SimpleNestedObject`,
		`	if err := swag.ReadJSON(b, &res); err != nil {`,
		`	*m = res`,
		`type SimpleNestedObjectAllOf0 struct {`,
		"	Prop1 strfmt.UUID `json:\"prop1,omitempty\"`",
		"	Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		`func (m *SimpleNestedObjectAllOf0) UnmarshalJSON(raw []byte) error {`,
		`	var dataAO0 struct {`,
		"		Prop1 strfmt.UUID `json:\"prop1,omitempty\"`",
		`	if err := swag.ReadJSON(raw, &dataAO0); err != nil {`,
		`	m.Prop1 = dataAO0.Prop1`,
		`	var dataAO1 struct {`,
		"		Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		`	if err := swag.ReadJSON(raw, &dataAO1); err != nil {`,
		`	m.Prop2 = dataAO1.Prop2`,
		`func (m SimpleNestedObjectAllOf0) MarshalJSON() ([]byte, error) {`,
		`	_parts := make([][]byte, 0, 2`,
		`	var dataAO0 struct {`,
		"		Prop1 strfmt.UUID `json:\"prop1,omitempty\"`",
		`	dataAO0.Prop1 = m.Prop1`,
		`	jsonDataAO0, errAO0 := swag.WriteJSON(dataAO0`,
		`	if errAO0 != nil {`,
		`		return nil, errAO0`,
		`	_parts = append(_parts, jsonDataAO0`,
		`	var dataAO1 struct {`,
		"		Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		`	dataAO1.Prop2 = m.Prop2`,
		`	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1`,
		`	if errAO1 != nil {`,
		`		return nil, errAO1`,
		`	_parts = append(_parts, jsonDataAO1`,
		`	return swag.ConcatJSON(_parts...), nil`,
		`func (m *SimpleNestedObjectAllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`	if err := m.validateProp2(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *SimpleNestedObjectAllOf0) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	if err := validate.FormatOf("prop1", "body", "uuid", m.Prop1.String(), formats); err != nil {`,
		`func (m *SimpleNestedObjectAllOf0) validateProp2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop2) {`,
		`	if err := validate.FormatOf("prop2", "body", "uuid", m.Prop2.String(), formats); err != nil {`,
		`func (m *SimpleNestedObjectAllOf0) MarshalBinary() ([]byte, error) {`,
		`	if m == nil {`,
		`		return nil, nil`,
		`	return swag.WriteJSON(m`,
		`func (m *SimpleNestedObjectAllOf0) UnmarshalBinary(b []byte) error {`,
		`	var res SimpleNestedObjectAllOf0`,
		`	if err := swag.ReadJSON(b, &res); err != nil {`,
		`	*m = res`,
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: break_nested_object_all_of1_prop7_all_of0.go
	flattenRun.AddExpectations("break_nested_object_all_of1_prop7_all_of0.go", []string{
		`type BreakNestedObjectAllOf1Prop7AllOf0 struct {`,
		"	Prop8 int64 `json:\"prop8,omitempty\"`",
		`func (m *BreakNestedObjectAllOf1Prop7AllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp8(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *BreakNestedObjectAllOf1Prop7AllOf0) validateProp8(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop8) {`,
		`	if err := validate.MinimumInt("prop8", "body", m.Prop8, 12, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: deep_nested_object_all_of1_all_of0.go
	flattenRun.AddExpectations("deep_nested_object_all_of1_all_of0.go", []string{
		`type DeepNestedObjectAllOf1AllOf0 struct {`,
		"	Prop4 strfmt.UUID `json:\"prop4,omitempty\"`",
		`func (m *DeepNestedObjectAllOf1AllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp4(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *DeepNestedObjectAllOf1AllOf0) validateProp4(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop4) {`,
		`	if err := validate.FormatOf("prop4", "body", "uuid", m.Prop4.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: break_nested_object_all_of1_prop7_all_of1.go
	flattenRun.AddExpectations("break_nested_object_all_of1_prop7_all_of1.go", []string{
		`type BreakNestedObjectAllOf1Prop7AllOf1 struct {`,
		"	Prop9 int64 `json:\"prop9,omitempty\"`",
		`func (m *BreakNestedObjectAllOf1Prop7AllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp9(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *BreakNestedObjectAllOf1Prop7AllOf1) validateProp9(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop9) {`,
		`	if err := validate.MaximumInt("prop9", "body", m.Prop9, 12, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: composed_thing_all_of0.go
	flattenRun.AddExpectations("composed_thing_all_of0.go", []string{
		`type ComposedThingAllOf0 struct {`,
		"	Prop1 strfmt.UUID `json:\"prop1,omitempty\"`",
		`func (m *ComposedThingAllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ComposedThingAllOf0) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	if err := validate.FormatOf("prop1", "body", "uuid", m.Prop1.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: composed_thing_all_of1.go
	flattenRun.AddExpectations("composed_thing_all_of1.go", []string{
		`type ComposedThingAllOf1 struct {`,
		"	Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		`func (m *ComposedThingAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp2(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ComposedThingAllOf1) validateProp2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop2) {`,
		`	if err := validate.FormatOf("prop2", "body", "uuid", m.Prop2.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: deep_nested_object_all_of1.go
	flattenRun.AddExpectations("deep_nested_object_all_of1.go", []string{
		`type DeepNestedObjectAllOf1 struct {`,
		`	DeepNestedObjectAllOf1AllOf0`,
		`	DeepNestedObjectAllOf1AllOf1`,
		`func (m *DeepNestedObjectAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.DeepNestedObjectAllOf1AllOf0.Validate(formats); err != nil {`,
		`	if err := m.DeepNestedObjectAllOf1AllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

}

func initFixtureComplexAllOf() {
	// testing ../fixtures/bugs/1487/fixture-complex-allOf.yaml with flatten and expand (--skip-flatten)

	/*
	 */
	f := newModelFixture("../fixtures/bugs/1487/fixture-complex-allOf.yaml", "fixture for nested allOf with ref")
	flattenRun := f.AddRun(false)
	expandRun := f.AddRun(true)

	// load expectations for model: aliased_date.go
	flattenRun.AddExpectations("aliased_date.go", []string{
		`type AliasedDate strfmt.Date`,
		`func (m AliasedDate) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.FormatOf("", "body", "date", strfmt.Date(m).String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("aliased_date.go", flattenRun.ExpectedFor("AliasedDate").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: object_mix_all_of2.go
	flattenRun.AddExpectations("object_mix_all_of2.go", []string{
		`type ObjectMixAllOf2 struct {`,
		"	Prop2 *ObjectMixAllOf2Prop2 `json:\"prop2,omitempty\"`",
		`func (m *ObjectMixAllOf2) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp2(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ObjectMixAllOf2) validateProp2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop2) {`,
		`	if m.Prop2 != nil {`,
		`		if err := m.Prop2.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("prop2"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: object_mix.go
	flattenRun.AddExpectations("object_mix.go", []string{
		`type ObjectMix struct {`,
		`	ObjectMixAllOf1`,
		`	ObjectMixAllOf2`,
		`func (m *ObjectMix) Validate(formats strfmt.Registry) error {`,
		`	if err := m.ObjectMixAllOf1.Validate(formats); err != nil {`,
		`	if err := m.ObjectMixAllOf2.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("object_mix.go", []string{
		`type ObjectMix struct {`,
		`	Prop1 struct {`,
		`		ObjectMixProp1AllOf0`,
		`		ObjectMixProp1AllOf1`,
		"	} `json:\"prop1,omitempty\"`",
		`	Prop2 struct {`,
		`		ObjectMixProp2AllOf0`,
		`		ObjectMixProp2AllOf1`,
		"	} `json:\"prop2,omitempty\"`",
		`func (m *ObjectMix) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`	if err := m.validateProp2(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ObjectMix) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`func (m *ObjectMix) validateProp2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop2) {`,
		`type ObjectMixProp1AllOf0 strfmt.Date`,
		`func (m ObjectMixProp1AllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.FormatOf("", "body", "date", strfmt.Date(m).String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`type ObjectMixProp1AllOf1 strfmt.Date`,
		`func (m ObjectMixProp1AllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.FormatOf("", "body", "date", strfmt.Date(m).String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`type ObjectMixProp2AllOf0 strfmt.Date`,
		`func (m ObjectMixProp2AllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.FormatOf("", "body", "date", strfmt.Date(m).String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`type ObjectMixProp2AllOf1 strfmt.Date`,
		`func (m ObjectMixProp2AllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.FormatOf("", "body", "date", strfmt.Date(m).String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: all_of_slices_of_aliases.go
	flattenRun.AddExpectations("all_of_slices_of_aliases.go", []string{
		`type AllOfSlicesOfAliases struct {`,
		`	AllOfSlicesOfAliasesAllOf0`,
		`	AllOfSlicesOfAliasesAllOf1`,
		`func (m *AllOfSlicesOfAliases) Validate(formats strfmt.Registry) error {`,
		`	if err := m.AllOfSlicesOfAliasesAllOf0.Validate(formats); err != nil {`,
		`	if err := m.AllOfSlicesOfAliasesAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("all_of_slices_of_aliases.go", []string{
		`type AllOfSlicesOfAliases struct {`,
		"	Prop1 []strfmt.Date `json:\"prop1\"`",
		"	Prop2 []*strfmt.Date `json:\"prop2\"`",
		`func (m *AllOfSlicesOfAliases) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`	if err := m.validateProp2(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AllOfSlicesOfAliases) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	iProp1Size := int64(len(m.Prop1)`,
		`	if err := validate.MaxItems("prop1", "body", iProp1Size, 10); err != nil {`,
		`	for i := 0; i < len(m.Prop1); i++ {`,
		`		if err := validate.FormatOf("prop1"+"."+strconv.Itoa(i), "body", "date", m.Prop1[i].String(), formats); err != nil {`,
		`func (m *AllOfSlicesOfAliases) validateProp2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop2) {`,
		`	iProp2Size := int64(len(m.Prop2)`,
		`	if err := validate.MaxItems("prop2", "body", iProp2Size, 20); err != nil {`,
		`	for i := 0; i < len(m.Prop2); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`		if swag.IsZero(m.Prop2[i]) {`,
		// nullable required:
		`		if err := validate.FormatOf("prop2"+"."+strconv.Itoa(i), "body", "date", m.Prop2[i].String(), formats); err != nil {`,
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: all_of_aliases.go
	flattenRun.AddExpectations("all_of_aliases.go", []string{
		`type AllOfAliases struct {`,
		`	AliasedDate`,
		`	AliasedNullableDate`,
		`func (m *AllOfAliases) Validate(formats strfmt.Registry) error {`,
		`	if err := m.AliasedDate.Validate(formats); err != nil {`,
		`	if err := m.AliasedNullableDate.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("all_of_aliases.go", []string{
		`type AllOfAliases struct {`,
		`	AllOfAliasesAllOf0`,
		`	AllOfAliasesAllOf1`,
		`func (m *AllOfAliases) Validate(formats strfmt.Registry) error {`,
		`	if err := m.AllOfAliasesAllOf0.Validate(formats); err != nil {`,
		`	if err := m.AllOfAliasesAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`type AllOfAliasesAllOf0 strfmt.Date`,
		`func (m AllOfAliasesAllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.FormatOf("", "body", "date", strfmt.Date(m).String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		// NOTE: the x-nullable has not been honored here
		// so we don't have: `type AllOfAliasesAllOf1 *strfmt.Date`,
		// this is by design, since nullability is honored by the container of the alias, not the
		// alias itself. An allOf branch container is composing types, not pointers.
		`type AllOfAliasesAllOf1 strfmt.Date`,
		`func (m AllOfAliasesAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.FormatOf("", "body", "date", strfmt.Date(m).String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: object_mix_all_of1.go
	flattenRun.AddExpectations("object_mix_all_of1.go", []string{
		`type ObjectMixAllOf1 struct {`,
		"	Prop1 *ObjectMixAllOf1Prop1 `json:\"prop1,omitempty\"`",
		`func (m *ObjectMixAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ObjectMixAllOf1) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	if m.Prop1 != nil {`,
		`		if err := m.Prop1.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("prop1"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: all_of_slices_of_aliases_all_of0.go
	flattenRun.AddExpectations("all_of_slices_of_aliases_all_of0.go", []string{
		`type AllOfSlicesOfAliasesAllOf0 struct {`,
		"	Prop1 []AliasedDate `json:\"prop1\"`",
		`func (m *AllOfSlicesOfAliasesAllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AllOfSlicesOfAliasesAllOf0) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	iProp1Size := int64(len(m.Prop1)`,
		`	if err := validate.MaxItems("prop1", "body", iProp1Size, 10); err != nil {`,
		`	for i := 0; i < len(m.Prop1); i++ {`,
		`		if err := m.Prop1[i].Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("prop1" + "." + strconv.Itoa(i)`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: slice_of_all_of.go
	flattenRun.AddExpectations("slice_of_all_of.go", []string{
		`type SliceOfAllOf []*SliceOfAllOfItems`,
		`func (m SliceOfAllOf) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.UniqueItems("", "body", m); err != nil {`,
		`	for i := 0; i < len(m); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`		if swag.IsZero(m[i]) {`,
		// nullable required:
		`		if m[i] != nil {`,
		`			if err := m[i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName(strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("slice_of_all_of.go", []string{
		`type SliceOfAllOf []*SliceOfAllOfItems0`,
		`func (m SliceOfAllOf) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.UniqueItems("", "body", m); err != nil {`,
		`	for i := 0; i < len(m); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`		if swag.IsZero(m[i]) {`,
		// nullable required:
		`		if m[i] != nil {`,
		`			if err := m[i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName(strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
		`type SliceOfAllOfItems0 struct {`,
		"	Prop0 strfmt.UUID `json:\"prop0,omitempty\"`",
		`	SliceOfAllOfItems0AllOf1`,
		`func (m *SliceOfAllOfItems0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp0(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *SliceOfAllOfItems0) validateProp0(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop0) {`,
		`	if err := validate.FormatOf("prop0", "body", "uuid", m.Prop0.String(), formats); err != nil {`,
		`type SliceOfAllOfItems0AllOf1 []interface{`,
		// empty validation
		"func (m SliceOfAllOfItems0AllOf1) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: object_mix_all_of2_prop2.go
	flattenRun.AddExpectations("object_mix_all_of2_prop2.go", []string{
		`type ObjectMixAllOf2Prop2 struct {`,
		`	AliasedDate`,
		`	AliasedNullableDate`,
		`func (m *ObjectMixAllOf2Prop2) Validate(formats strfmt.Registry) error {`,
		`	if err := m.AliasedDate.Validate(formats); err != nil {`,
		`	if err := m.AliasedNullableDate.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: slice_of_all_of_items_all_of0.go
	flattenRun.AddExpectations("slice_of_all_of_items_all_of0.go", []string{
		`type SliceOfAllOfItemsAllOf0 struct {`,
		"	Prop0 strfmt.UUID `json:\"prop0,omitempty\"`",
		`func (m *SliceOfAllOfItemsAllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp0(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *SliceOfAllOfItemsAllOf0) validateProp0(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop0) {`,
		`	if err := validate.FormatOf("prop0", "body", "uuid", m.Prop0.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: slice_of_interfaces.go
	flattenRun.AddExpectations("slice_of_interfaces.go", []string{
		`type SliceOfInterfaces []interface{`,
		// empty validation
		"func (m SliceOfInterfaces) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("slice_of_interfaces.go", flattenRun.ExpectedFor("SliceOfInterfaces").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: slice_of_interfaces_with_validation.go
	flattenRun.AddExpectations("slice_of_interfaces_with_validation.go", []string{
		`type SliceOfInterfacesWithValidation []interface{`,
		`func (m SliceOfInterfacesWithValidation) Validate(formats strfmt.Registry) error {`,
		`	iSliceOfInterfacesWithValidationSize := int64(len(m)`,
		`	if err := validate.MaxItems("", "body", iSliceOfInterfacesWithValidationSize, 10); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("slice_of_interfaces_with_validation.go", flattenRun.ExpectedFor("SliceOfInterfacesWithValidation").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: aliased_nullable_date.go
	flattenRun.AddExpectations("aliased_nullable_date.go", []string{
		`type AliasedNullableDate strfmt.Date`,
		`func (m AliasedNullableDate) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.FormatOf("", "body", "date", strfmt.Date(m).String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("aliased_nullable_date.go", flattenRun.ExpectedFor("AliasedNullableDate").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: slice_mix.go
	flattenRun.AddExpectations("slice_mix.go", []string{
		`type SliceMix struct {`,
		`	SliceOfAllOf`,
		`	SliceOfInterfaces`,
		`func (m *SliceMix) Validate(formats strfmt.Registry) error {`,
		`	if err := m.SliceOfAllOf.Validate(formats); err != nil {`,
		`	if err := m.SliceOfInterfaces.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// disable log assertions (dodgy with parallel tests)
		// output in log
		// warning,
		noLines,
		noLines)

	expandRun.AddExpectations("slice_mix.go", []string{
		`type SliceMix struct {`,
		`	SliceMixAllOf0`,
		`	SliceMixAllOf1`,
		`func (m *SliceMix) Validate(formats strfmt.Registry) error {`,
		`	if err := m.SliceMixAllOf0.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`type SliceMixAllOf0 []*SliceMixAllOf0Items0`,
		`func (m SliceMixAllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.UniqueItems("", "body", m); err != nil {`,
		`	for i := 0; i < len(m); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`		if swag.IsZero(m[i]) {`,
		// nullable required:
		`		if m[i] != nil {`,
		`			if err := m[i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName(strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
		`type SliceMixAllOf0Items0 struct {`,
		"	Prop0 strfmt.UUID `json:\"prop0,omitempty\"`",
		`	SliceMixAllOf0Items0AllOf1`,
		`func (m *SliceMixAllOf0Items0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp0(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *SliceMixAllOf0Items0) validateProp0(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop0) {`,
		`	if err := validate.FormatOf("prop0", "body", "uuid", m.Prop0.String(), formats); err != nil {`,
		`type SliceMixAllOf0Items0AllOf1 []interface{`,
		`type SliceMixAllOf1 []interface{`,
		// empty validation
		"func (m SliceMixAllOf0Items0AllOf1) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: object_mix_all_of1_prop1.go
	flattenRun.AddExpectations("object_mix_all_of1_prop1.go", []string{
		`type ObjectMixAllOf1Prop1 struct {`,
		`	AliasedDate`,
		`	AliasedNullableDate`,
		`func (m *ObjectMixAllOf1Prop1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.AliasedDate.Validate(formats); err != nil {`,
		`	if err := m.AliasedNullableDate.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: slice_of_all_of_items.go
	flattenRun.AddExpectations("slice_of_all_of_items.go", []string{
		`type SliceOfAllOfItems struct {`,
		`	SliceOfAllOfItemsAllOf0`,
		`	SliceOfInterfaces`,
		`func (m *SliceOfAllOfItems) Validate(formats strfmt.Registry) error {`,
		`	if err := m.SliceOfAllOfItemsAllOf0.Validate(formats); err != nil {`,
		`	if err := m.SliceOfInterfaces.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: all_of_slices_of_aliases_all_of1.go
	flattenRun.AddExpectations("all_of_slices_of_aliases_all_of1.go", []string{
		`type AllOfSlicesOfAliasesAllOf1 struct {`,
		"	Prop2 []*AliasedNullableDate `json:\"prop2\"`",
		`func (m *AllOfSlicesOfAliasesAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp2(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AllOfSlicesOfAliasesAllOf1) validateProp2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop2) {`,
		`	iProp2Size := int64(len(m.Prop2)`,
		`	if err := validate.MaxItems("prop2", "body", iProp2Size, 20); err != nil {`,
		`	for i := 0; i < len(m.Prop2); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`		if swag.IsZero(m.Prop2[i]) {`,
		// nullable required:
		`		if m.Prop2[i] != nil {`,
		`			if err := m.Prop2[i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName("prop2" + "." + strconv.Itoa(i)`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

}

func initFixtureIsNullable() {
	// testing ../fixtures/bugs/1487/fixture-is-nullable.yaml with flatten and expand (--skip-flatten)

	/* just an elementary check with the x-nullable tag
	 */
	f := newModelFixture("../fixtures/bugs/1487/fixture-is-nullable.yaml", "fixture for x-nullable flag")
	flattenRun := f.AddRun(false)
	expandRun := f.AddRun(true)

	// load expectations for model: thing_with_nullable_dates.go
	flattenRun.AddExpectations("thing_with_nullable_dates.go", []string{
		`type ThingWithNullableDates struct {`,
		"	Prop1 strfmt.Date `json:\"prop1,omitempty\"`",
		"	Prop2 *strfmt.Date `json:\"prop2,omitempty\"`",
		`func (m *ThingWithNullableDates) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`	if err := m.validateProp2(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ThingWithNullableDates) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	if err := validate.FormatOf("prop1", "body", "date", m.Prop1.String(), formats); err != nil {`,
		`func (m *ThingWithNullableDates) validateProp2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop2) {`,
		`	if err := validate.FormatOf("prop2", "body", "date", m.Prop2.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("thing_with_nullable_dates.go", flattenRun.ExpectedFor("ThingWithNullableDates").ExpectedLines, todo, noLines, noLines)

}

func initFixtureItching() {
	// testing ../fixtures/bugs/1487/fixture-itching.yaml with flatten and expand (--skip-flatten)

	/*
		This one regroups a number of itching cases, essentially around additionalProperties.
		In particular, we test some things with empty objects (no properties) which have additionalProperties of diverse sorts.
		We also added here some funny models using the special types Files, string format: binary and interface{}
		These special cases do not correspond to actual API specs: we use them to verify the internal behavior of the general.
	*/
	f := newModelFixture("../fixtures/bugs/1487/fixture-itching.yaml", "fixture for additionalProperties")
	flattenRun := f.AddRun(false)
	expandRun := f.AddRun(true)

	// load expectations for model: top_level_format_issue_my_alternate_file.go
	flattenRun.AddExpectations("top_level_format_issue_my_alternate_file.go", []string{
		`import "io"`,
		`type TopLevelFormatIssueMyAlternateFile io.ReadCloser`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: not_validated_additional_props.go
	flattenRun.AddExpectations("not_validated_additional_props.go", []string{
		`type NotValidatedAdditionalProps struct {`,
		"	Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		"	NotValidatedAdditionalProps map[string]map[string]map[string]string `json:\"-\"`",
		`func (m *NotValidatedAdditionalProps) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp2(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NotValidatedAdditionalProps) validateProp2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop2) {`,
		`	if err := validate.FormatOf("prop2", "body", "uuid", m.Prop2.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("not_validated_additional_props.go", flattenRun.ExpectedFor("NotValidatedAdditionalProps").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: aliased_nullable_file.go
	flattenRun.AddExpectations("aliased_nullable_file.go", []string{
		`import "io"`,
		`type AliasedNullableFile io.ReadCloser`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("aliased_nullable_file.go", flattenRun.ExpectedFor("AliasedNullableFile").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: empty_object_with_additional_nullable_primitive.go
	flattenRun.AddExpectations("empty_object_with_additional_nullable_primitive.go", []string{
		`type EmptyObjectWithAdditionalNullablePrimitive map[string]*strfmt.Date`,
		`func (m EmptyObjectWithAdditionalNullablePrimitive) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if swag.IsZero(m[k]) {`,
		`		if err := validate.FormatOf(k, "body", "date", m[k].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("empty_object_with_additional_nullable_primitive.go", flattenRun.ExpectedFor("EmptyObjectWithAdditionalNullablePrimitive").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: empty_object_with_additional_alias.go
	flattenRun.AddExpectations("empty_object_with_additional_alias.go", []string{
		`type EmptyObjectWithAdditionalAlias map[string]AliasedThing`,
		`func (m EmptyObjectWithAdditionalAlias) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if err := validate.Required(k, "body", m[k]); err != nil {`,
		`		if val, ok := m[k]; ok {`,
		`			if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("empty_object_with_additional_alias.go", []string{
		`type EmptyObjectWithAdditionalAlias map[string]EmptyObjectWithAdditionalAliasAnon`,
		`func (m EmptyObjectWithAdditionalAlias) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if swag.IsZero(m[k]) {`,
		`		if val, ok := m[k]; ok {`,
		`			if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`type EmptyObjectWithAdditionalAliasAnon struct {`,
		"	Prop1 strfmt.Date `json:\"prop1,omitempty\"`",
		`func (m *EmptyObjectWithAdditionalAliasAnon) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *EmptyObjectWithAdditionalAliasAnon) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	if err := validate.FormatOf("prop1", "body", "date", m.Prop1.String(), formats); err != nil {`,
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: nullable_thing.go
	flattenRun.AddExpectations("nullable_thing.go", []string{
		`type NullableThing strfmt.Date`,
		`func (m NullableThing) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.FormatOf("", "body", "date", strfmt.Date(m).String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("nullable_thing.go", flattenRun.ExpectedFor("NullableThing").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: slice_of_aliased_files.go
	flattenRun.AddExpectations("slice_of_aliased_files.go", []string{
		`type SliceOfAliasedFiles []AliasedFile`,
		`func (m SliceOfAliasedFiles) Validate(formats strfmt.Registry) error {`,
		`	iSliceOfAliasedFilesSize := int64(len(m)`,
		`	if err := validate.MinItems("", "body", iSliceOfAliasedFilesSize, 4); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("slice_of_aliased_files.go", []string{
		`type SliceOfAliasedFiles []io.ReadCloser`,
		`func (m SliceOfAliasedFiles) Validate(formats strfmt.Registry) error {`,
		`	iSliceOfAliasedFilesSize := int64(len(m)`,
		`	if err := validate.MinItems("", "body", iSliceOfAliasedFilesSize, 4); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: empty_object_with_additional_non_nullable_primitive.go
	flattenRun.AddExpectations("empty_object_with_additional_non_nullable_primitive.go", []string{
		`type EmptyObjectWithAdditionalNonNullablePrimitive map[string]strfmt.Date`,
		`func (m EmptyObjectWithAdditionalNonNullablePrimitive) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		// fix undue IsZero call
		`		if err := validate.FormatOf(k, "body", "date", m[k].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("empty_object_with_additional_non_nullable_primitive.go", flattenRun.ExpectedFor("EmptyObjectWithAdditionalNonNullablePrimitive").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: good_old_format_issue.go
	flattenRun.AddExpectations("good_old_format_issue.go", []string{
		`type GoodOldFormatIssue struct {`,
		"	AlternateFile GoodOldFormatIssueAlternateFile `json:\"alternateFile,omitempty\"`",
		"	AnotherFile io.ReadCloser `json:\"anotherFile,omitempty\"`",
		"	MyBytes strfmt.Base64 `json:\"myBytes,omitempty\"`",
		"	MyFile io.ReadCloser `json:\"myFile\"`",
		"	ThisAliasedFile AliasedFile `json:\"thisAliasedFile,omitempty\"`",
		"	ThisAlternateAliasedFile AliasedTypeFile `json:\"thisAlternateAliasedFile,omitempty\"`",
		"	ThisNullableAliasedFile *AliasedNullableFile `json:\"thisNullableAliasedFile,omitempty\"`",
		"	ThisNullableAlternateAliasedFile *AliasedTypeNullableFile `json:\"thisNullableAlternateAliasedFile,omitempty\"`",
		`func (m *GoodOldFormatIssue) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMyFile(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *GoodOldFormatIssue) validateMyFile(formats strfmt.Registry) error {`,
		`	if err := validate.Required("myFile", "body", io.ReadCloser(m.MyFile)); err != nil {`,
	},
		// not expected
		[]string{
			`	if err := m.validateMyBytes(formats); err != nil {`,
			`func (m *GoodOldFormatIssue) validateMyBytes(formats strfmt.Registry) error {`,
			`	if err := validate.FormatOf("myBytes", "body", "byte", m.MyBytes.String(), formats); err != nil {`,
		},
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("good_old_format_issue.go", []string{
		`type GoodOldFormatIssue struct {`,
		"	AlternateFile io.ReadCloser `json:\"alternateFile,omitempty\"`",
		"	AnotherFile io.ReadCloser `json:\"anotherFile,omitempty\"`",
		"	MyBytes strfmt.Base64 `json:\"myBytes,omitempty\"`",
		"	MyFile io.ReadCloser `json:\"myFile\"`",
		"	ThisAliasedFile io.ReadCloser `json:\"thisAliasedFile,omitempty\"`",
		"	ThisAlternateAliasedFile io.ReadCloser `json:\"thisAlternateAliasedFile,omitempty\"`",
		"	ThisNullableAliasedFile io.ReadCloser `json:\"thisNullableAliasedFile,omitempty\"`",
		"	ThisNullableAlternateAliasedFile io.ReadCloser `json:\"thisNullableAlternateAliasedFile,omitempty\"`",
		`func (m *GoodOldFormatIssue) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMyFile(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *GoodOldFormatIssue) validateMyFile(formats strfmt.Registry) error {`,
		`	if err := validate.Required("myFile", "body", io.ReadCloser(m.MyFile)); err != nil {`,
	},
		// not expected
		[]string{
			`	if err := m.validateMyBytes(formats); err != nil {`,
			`func (m *GoodOldFormatIssue) validateMyBytes(formats strfmt.Registry) error {`,
			`	if err := validate.FormatOf("myBytes", "body", "byte", m.MyBytes.String(), formats); err != nil {`,
		},
		// output in log
		noLines,
		noLines)

	// load expectations for model: empty_object_with_additional_slice_additional_properties_items.go
	flattenRun.AddExpectations("empty_object_with_additional_slice_additional_properties_items.go", []string{
		`type EmptyObjectWithAdditionalSliceAdditionalPropertiesItems struct {`,
		"	DummyProp1 strfmt.Date `json:\"dummyProp1,omitempty\"`",
		`func (m *EmptyObjectWithAdditionalSliceAdditionalPropertiesItems) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateDummyProp1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *EmptyObjectWithAdditionalSliceAdditionalPropertiesItems) validateDummyProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.DummyProp1) {`,
		`	if err := validate.FormatOf("dummyProp1", "body", "date", m.DummyProp1.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: not_validated_additional_props_slice.go
	flattenRun.AddExpectations("not_validated_additional_props_slice.go", []string{
		`type NotValidatedAdditionalPropsSlice struct {`,
		"	Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		"	NotValidatedAdditionalPropsSlice map[string][]map[string]map[string]string `json:\"-\"`",
		`func (m *NotValidatedAdditionalPropsSlice) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp2(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NotValidatedAdditionalPropsSlice) validateProp2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop2) {`,
		`	if err := validate.FormatOf("prop2", "body", "uuid", m.Prop2.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: aliased_type_file.go
	flattenRun.AddExpectations("aliased_type_file.go", []string{
		`import "io"`,
		`type AliasedTypeFile io.ReadCloser`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("aliased_type_file.go", flattenRun.ExpectedFor("AliasedTypeFile").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: object_with_empty_object.go
	flattenRun.AddExpectations("object_with_empty_object.go", []string{
		`type ObjectWithEmptyObject struct {`,
		"	EmptyObj EmptyObjectWithAdditionalAlias `json:\"emptyObj,omitempty\"`",
		"	NonEmptyObj *NullableThing `json:\"nonEmptyObj,omitempty\"`",
		`func (m *ObjectWithEmptyObject) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateEmptyObj(formats); err != nil {`,
		`	if err := m.validateNonEmptyObj(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ObjectWithEmptyObject) validateEmptyObj(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.EmptyObj) {`,
		`	if err := m.EmptyObj.Validate(formats); err != nil {`,
		`		if ve, ok := err.(*errors.Validation); ok {`,
		`			return ve.ValidateName("emptyObj"`,
		`func (m *ObjectWithEmptyObject) validateNonEmptyObj(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.NonEmptyObj) {`,
		`	if m.NonEmptyObj != nil {`,
		`		if err := m.NonEmptyObj.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("nonEmptyObj"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("object_with_empty_object.go", []string{
		`type ObjectWithEmptyObject struct {`,
		"	EmptyObj map[string]ObjectWithEmptyObjectEmptyObjAnon `json:\"emptyObj,omitempty\"`",
		"	NonEmptyObj *strfmt.Date `json:\"nonEmptyObj,omitempty\"`",
		`func (m *ObjectWithEmptyObject) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateEmptyObj(formats); err != nil {`,
		`	if err := m.validateNonEmptyObj(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ObjectWithEmptyObject) validateEmptyObj(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.EmptyObj) {`,
		`	for k := range m.EmptyObj {`,
		`		if swag.IsZero(m.EmptyObj[k]) {`,
		`		if val, ok := m.EmptyObj[k]; ok {`,
		`			if err := val.Validate(formats); err != nil {`,
		`func (m *ObjectWithEmptyObject) validateNonEmptyObj(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.NonEmptyObj) {`,
		`	if err := validate.FormatOf("nonEmptyObj", "body", "date", m.NonEmptyObj.String(), formats); err != nil {`,
		`type ObjectWithEmptyObjectEmptyObjAnon struct {`,
		"	Prop1 strfmt.Date `json:\"prop1,omitempty\"`",
		`func (m *ObjectWithEmptyObjectEmptyObjAnon) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ObjectWithEmptyObjectEmptyObjAnon) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	if err := validate.FormatOf("prop1", "body", "date", m.Prop1.String(), formats); err != nil {`,
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: aliased_file.go
	flattenRun.AddExpectations("aliased_file.go", []string{
		`import "io"`,
		`type AliasedFile io.ReadCloser`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("aliased_file.go", flattenRun.ExpectedFor("AliasedFile").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: empty_object_with_additional_slice.go
	flattenRun.AddExpectations("empty_object_with_additional_slice.go", []string{
		`type EmptyObjectWithAdditionalSlice map[string][]EmptyObjectWithAdditionalSliceAdditionalPropertiesItems`,
		`func (m EmptyObjectWithAdditionalSlice) Validate(formats strfmt.Registry) error {`,
		// fixed undue Required on this aliased type
		`	for k := range m {`,
		`		if err := validate.Required(k, "body", m[k]); err != nil {`,
		`		for i := 0; i < len(m[k]); i++ {`,
		`			if err := m[k][i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName(k + "." + strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("empty_object_with_additional_slice.go", []string{
		`type EmptyObjectWithAdditionalSlice map[string][]EmptyObjectWithAdditionalSliceItems0`,
		`func (m EmptyObjectWithAdditionalSlice) Validate(formats strfmt.Registry) error {`,
		// fixed undue Required on this aliased type
		`	for k := range m {`,
		`		if err := validate.Required(k, "body", m[k]); err != nil {`,
		`		for i := 0; i < len(m[k]); i++ {`,
		`			if err := m[k][i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName(k + "." + strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
		`type EmptyObjectWithAdditionalSliceItems0 struct {`,
		"	DummyProp1 strfmt.Date `json:\"dummyProp1,omitempty\"`",
		`func (m *EmptyObjectWithAdditionalSliceItems0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateDummyProp1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *EmptyObjectWithAdditionalSliceItems0) validateDummyProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.DummyProp1) {`,
		`	if err := validate.FormatOf("dummyProp1", "body", "date", m.DummyProp1.String(), formats); err != nil {`,
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_aliased_file.go
	flattenRun.AddExpectations("additional_aliased_file.go", []string{
		`type AdditionalAliasedFile interface{`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_aliased_file.go", flattenRun.ExpectedFor("AdditionalAliasedFile").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: good_old_format_issue_alternate_file.go
	flattenRun.AddExpectations("good_old_format_issue_alternate_file.go", []string{
		`import "io"`,
		`type GoodOldFormatIssueAlternateFile io.ReadCloser`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: empty_object_with_additional_nested_slice_additional_properties_items_items_items.go
	flattenRun.AddExpectations("empty_object_with_additional_nested_slice_additional_properties_items_items_items.go", []string{
		`type EmptyObjectWithAdditionalNestedSliceAdditionalPropertiesItemsItemsItems struct {`,
		"	DummyProp1 strfmt.Date `json:\"dummyProp1,omitempty\"`",
		`func (m *EmptyObjectWithAdditionalNestedSliceAdditionalPropertiesItemsItemsItems) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateDummyProp1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *EmptyObjectWithAdditionalNestedSliceAdditionalPropertiesItemsItemsItems) validateDummyProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.DummyProp1) {`,
		`	if err := validate.FormatOf("dummyProp1", "body", "date", m.DummyProp1.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: aliased_thing.go
	flattenRun.AddExpectations("aliased_thing.go", []string{
		`type AliasedThing struct {`,
		"	Prop1 strfmt.Date `json:\"prop1,omitempty\"`",
		`func (m *AliasedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AliasedThing) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	if err := validate.FormatOf("prop1", "body", "date", m.Prop1.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("aliased_thing.go", flattenRun.ExpectedFor("AliasedThing").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: additional_file.go
	flattenRun.AddExpectations("additional_file.go", []string{
		`type AdditionalFile struct {`,
		"	DirName string `json:\"dirName,omitempty\"`",
		"	AdditionalFile map[string]io.ReadCloser `json:\"-\"`",
		// empty validation
		"func (m *AdditionalFile) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_file.go", flattenRun.ExpectedFor("AdditionalFile").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: aliased_type_nullable_file.go
	flattenRun.AddExpectations("aliased_type_nullable_file.go", []string{
		`import "io"`,
		`type AliasedTypeNullableFile io.ReadCloser`,
	},
		// not expected
		validatable,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("aliased_type_nullable_file.go", flattenRun.ExpectedFor("AliasedTypeNullableFile").ExpectedLines, validatable, noLines, noLines)

	// load expectations for model: top_level_format_issue.go
	flattenRun.AddExpectations("top_level_format_issue.go", []string{
		`type TopLevelFormatIssue struct {`,
		"	MyAlternateFile TopLevelFormatIssueMyAlternateFile `json:\"myAlternateFile,omitempty\"`",
		"	MyFile io.ReadCloser `json:\"myFile,omitempty\"`",
		// empty validation
		"func (m *TopLevelFormatIssue) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("top_level_format_issue.go", []string{
		`type TopLevelFormatIssue struct {`,
		"	MyAlternateFile io.ReadCloser `json:\"myAlternateFile,omitempty\"`",
		"	MyFile io.ReadCloser `json:\"myFile,omitempty\"`",
		// empty validation
		"func (m *TopLevelFormatIssue) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: enums_with_additional_props.go
	flattenRun.AddExpectations("enums_with_additional_props.go", []string{
		`type EnumsWithAdditionalProps map[string]interface{`,
		`var enumsWithAdditionalPropsEnum []interface{`,
		`	var res []EnumsWithAdditionalProps`,
		"	if err := json.Unmarshal([]byte(`[\"{ \\\"a\\\": 1, \\\"b\\\": 2 }\",\"{ \\\"a\\\": 3, \\\"b\\\": 4 }\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		enumsWithAdditionalPropsEnum = append(enumsWithAdditionalPropsEnum, v`,
		`func (m *EnumsWithAdditionalProps) validateEnumsWithAdditionalPropsEnum(path, location string, value EnumsWithAdditionalProps) error {`,
		`	if err := validate.EnumCase(path, location, value, enumsWithAdditionalPropsEnum, true); err != nil {`,
		`var enumsWithAdditionalPropsValueEnum []interface{`,
		`	var res []interface{`,
		"	if err := json.Unmarshal([]byte(`[\"{ \\\"b\\\": 2 }\",\"{ \\\"b\\\": 4 }\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		enumsWithAdditionalPropsValueEnum = append(enumsWithAdditionalPropsValueEnum, v`,
		`func (m *EnumsWithAdditionalProps) validateEnumsWithAdditionalPropsValueEnum(path, location string, value interface{}) error {`,
		`	if err := validate.EnumCase(path, location, value, enumsWithAdditionalPropsValueEnum, true); err != nil {`,
		`func (m EnumsWithAdditionalProps) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if err := m.validateEnumsWithAdditionalPropsValueEnum(k, "body", m[k]); err != nil {`,
		`	if err := m.validateEnumsWithAdditionalPropsEnum("", "body", m); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("enums_with_additional_props.go", flattenRun.ExpectedFor("EnumsWithAdditionalProps").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: empty_object_with_additional_nested_slice.go
	flattenRun.AddExpectations("empty_object_with_additional_nested_slice.go", []string{
		`type EmptyObjectWithAdditionalNestedSlice map[string][][][]EmptyObjectWithAdditionalNestedSliceAdditionalPropertiesItemsItemsItems`,
		`func (m EmptyObjectWithAdditionalNestedSlice) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if err := validate.Required(k, "body", m[k]); err != nil {`,
		`		for i := 0; i < len(m[k]); i++ {`,
		`			for ii := 0; ii < len(m[k][i]); ii++ {`,
		`				for iii := 0; iii < len(m[k][i][ii]); iii++ {`,
		`					if err := m[k][i][ii][iii].Validate(formats); err != nil {`,
		`						if ve, ok := err.(*errors.Validation); ok {`,
		`							return ve.ValidateName(k + "." + strconv.Itoa(i) + "." + strconv.Itoa(ii) + "." + strconv.Itoa(iii)`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("empty_object_with_additional_nested_slice.go", []string{
		`type EmptyObjectWithAdditionalNestedSlice map[string][][][]EmptyObjectWithAdditionalNestedSliceItems0`,
		`func (m EmptyObjectWithAdditionalNestedSlice) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if err := validate.Required(k, "body", m[k]); err != nil {`,
		`		for i := 0; i < len(m[k]); i++ {`,
		`			for ii := 0; ii < len(m[k][i]); ii++ {`,
		`				for iii := 0; iii < len(m[k][i][ii]); iii++ {`,
		`					if err := m[k][i][ii][iii].Validate(formats); err != nil {`,
		`						if ve, ok := err.(*errors.Validation); ok {`,
		`							return ve.ValidateName(k + "." + strconv.Itoa(i) + "." + strconv.Itoa(ii) + "." + strconv.Itoa(iii)`,
		`		return errors.CompositeValidationError(res...`,
		`type EmptyObjectWithAdditionalNestedSliceItems0 struct {`,
		"	DummyProp1 strfmt.Date `json:\"dummyProp1,omitempty\"`",
		`func (m *EmptyObjectWithAdditionalNestedSliceItems0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateDummyProp1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *EmptyObjectWithAdditionalNestedSliceItems0) validateDummyProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.DummyProp1) {`,
		`	if err := validate.FormatOf("dummyProp1", "body", "date", m.DummyProp1.String(), formats); err != nil {`,
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: empty_object_with_additional_nullable.go
	// fixed nullability of aliased type
	flattenRun.AddExpectations("empty_object_with_additional_nullable.go", []string{
		`type EmptyObjectWithAdditionalNullable map[string]*NullableThing`,
		`func (m EmptyObjectWithAdditionalNullable) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if swag.IsZero(m[k]) {`,
		`		if val, ok := m[k]; ok {`,
		`			if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("empty_object_with_additional_nullable.go", []string{
		`type EmptyObjectWithAdditionalNullable map[string]*strfmt.Date`,
		`func (m EmptyObjectWithAdditionalNullable) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if swag.IsZero(m[k]) {`,
		`		if err := validate.FormatOf(k, "body", "date", m[k].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		noLines,
		// output in log
		noLines,
		noLines)

	// load expectations for model: not_validated_at_all.go
	flattenRun.AddExpectations("not_validated_at_all.go", []string{
		`type NotValidatedAtAll struct {`,
		"	Prop2 string `json:\"prop2,omitempty\"`",
		"	NotValidatedAtAll map[string][]map[string]map[string]string `json:\"-\"`",
		// empty validation
		"func (m *NotValidatedAtAll) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("not_validated_at_all.go", flattenRun.ExpectedFor("NotValidatedAtAll").ExpectedLines, todo, noLines, noLines)
}

func initFixtureAdditionalProps() {
	// testing ../fixtures/bugs/1487/fixture-additionalProps.yaml with flatten and expand (--skip-flatten)

	/*
		various patterns of additionalProperties
	*/
	f := newModelFixture("../fixtures/bugs/1487/fixture-additionalProps.yaml", "fixture for additionalProperties")
	flattenRun := f.AddRun(false)
	expandRun := f.AddRun(true)

	// load expectations for model: additional_object_with_formated_thing.go
	flattenRun.AddExpectations("additional_object_with_formated_thing.go", []string{
		`type AdditionalObjectWithFormatedThing struct {`,
		"	Blob *int64 `json:\"blob\"`",
		"	AdditionalObjectWithFormatedThing map[string]strfmt.Date `json:\"-\"`",
		`func (m *AdditionalObjectWithFormatedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateBlob(formats); err != nil {`,
		`	for k := range m.AdditionalObjectWithFormatedThing {`,
		// removed undue IZero call
		`		if err := validate.FormatOf(k, "body", "date", m.AdditionalObjectWithFormatedThing[k].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalObjectWithFormatedThing) validateBlob(formats strfmt.Registry) error {`,
		`	if err := validate.Required("blob", "body", m.Blob); err != nil {`,
		`	if err := validate.MinimumInt("blob", "body", *m.Blob, 1, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_object_with_formated_thing.go", flattenRun.ExpectedFor("AdditionalObjectWithFormatedThing").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: aliased_date.go
	flattenRun.AddExpectations("aliased_date.go", []string{
		`type AliasedDate strfmt.Date`,
		`func (m AliasedDate) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.FormatOf("", "body", "date", strfmt.Date(m).String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("aliased_date.go", flattenRun.ExpectedFor("AliasedDate").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: additional_array_of_refed_thing.go
	flattenRun.AddExpectations("additional_array_of_refed_thing.go", []string{
		`type AdditionalArrayOfRefedThing struct {`,
		"	ThisOneNotRequired int64 `json:\"thisOneNotRequired,omitempty\"`",
		"	AdditionalArrayOfRefedThing map[string][]AliasedDate `json:\"-\"`",
		`func (m *AdditionalArrayOfRefedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequired(formats); err != nil {`,
		`	for k := range m.AdditionalArrayOfRefedThing {`,
		// removed undue IsZero call
		`		if err := validate.UniqueItems(k, "body", m.AdditionalArrayOfRefedThing[k]); err != nil {`,
		`		for i := 0; i < len(m.AdditionalArrayOfRefedThing[k]); i++ {`,
		`			if err := m.AdditionalArrayOfRefedThing[k][i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName(k + "." + strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalArrayOfRefedThing) validateThisOneNotRequired(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequired) {`,
		`	if err := validate.MaximumInt("thisOneNotRequired", "body", m.ThisOneNotRequired, 10, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_array_of_refed_thing.go", []string{
		`type AdditionalArrayOfRefedThing struct {`,
		"	ThisOneNotRequired int64 `json:\"thisOneNotRequired,omitempty\"`",
		"	AdditionalArrayOfRefedThing map[string][]strfmt.Date `json:\"-\"`",
		`func (m *AdditionalArrayOfRefedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequired(formats); err != nil {`,
		`	for k := range m.AdditionalArrayOfRefedThing {`,
		// removed undue IsZero() call
		`		if err := validate.UniqueItems(k, "body", m.AdditionalArrayOfRefedThing[k]); err != nil {`,
		`		for i := 0; i < len(m.AdditionalArrayOfRefedThing[k]); i++ {`,
		`			if err := validate.FormatOf(k+"."+strconv.Itoa(i), "body", "date", m.AdditionalArrayOfRefedThing[k][i].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalArrayOfRefedThing) validateThisOneNotRequired(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequired) {`,
		`	if err := validate.MaximumInt("thisOneNotRequired", "body", m.ThisOneNotRequired, 10, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_object_with_nullable_thing.go
	flattenRun.AddExpectations("additional_object_with_nullable_thing.go", []string{
		`type AdditionalObjectWithNullableThing struct {`,
		"	Blob int64 `json:\"blob,omitempty\"`",
		"	AdditionalObjectWithNullableThing map[string]*AliasedNullableDate `json:\"-\"`",
		`func (m *AdditionalObjectWithNullableThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateBlob(formats); err != nil {`,
		`	for k := range m.AdditionalObjectWithNullableThing {`,
		`		if swag.IsZero(m.AdditionalObjectWithNullableThing[k]) {`,
		`		if val, ok := m.AdditionalObjectWithNullableThing[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalObjectWithNullableThing) validateBlob(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Blob) {`,
		`	if err := validate.MinimumInt("blob", "body", m.Blob, 1, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_object_with_nullable_thing.go", []string{
		`type AdditionalObjectWithNullableThing struct {`,
		"	Blob int64 `json:\"blob,omitempty\"`",
		"	AdditionalObjectWithNullableThing map[string]*strfmt.Date `json:\"-\"`",
		`func (m *AdditionalObjectWithNullableThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateBlob(formats); err != nil {`,
		`	for k := range m.AdditionalObjectWithNullableThing {`,
		`		if swag.IsZero(m.AdditionalObjectWithNullableThing[k]) {`,
		`		if err := validate.FormatOf(k, "body", "date", m.AdditionalObjectWithNullableThing[k].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalObjectWithNullableThing) validateBlob(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Blob) {`,
		`	if err := validate.MinimumInt("blob", "body", m.Blob, 1, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_things.go
	flattenRun.AddExpectations("additional_things.go", []string{
		`type AdditionalThings struct {`,
		"	Origin *string `json:\"origin\"`",
		"	Status string `json:\"status,omitempty\"`",
		"	AdditionalThings map[string]string `json:\"-\"`",
		`var additionalThingsValueEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"bookshop\",\"amazon\",\"library\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		additionalThingsValueEnum = append(additionalThingsValueEnum, v`,
		`func (m *AdditionalThings) validateAdditionalThingsValueEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, additionalThingsValueEnum, true); err != nil {`,
		`func (m *AdditionalThings) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateOrigin(formats); err != nil {`,
		`	if err := m.validateStatus(formats); err != nil {`,
		`	for k := range m.AdditionalThings {`,
		// removed undue IsZero call
		`		if err := m.validateAdditionalThingsValueEnum(k, "body", m.AdditionalThings[k]); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`var additionalThingsTypeOriginPropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"print\",\"e-book\",\"collection\",\"museum\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		additionalThingsTypeOriginPropEnum = append(additionalThingsTypeOriginPropEnum, v`,
		`	AdditionalThingsOriginPrint string = "print"`,
		`	AdditionalThingsOriginEDashBook string = "e-book"`,
		`	AdditionalThingsOriginCollection string = "collection"`,
		`	AdditionalThingsOriginMuseum string = "museum"`,
		`func (m *AdditionalThings) validateOriginEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, additionalThingsTypeOriginPropEnum, true); err != nil {`,
		`func (m *AdditionalThings) validateOrigin(formats strfmt.Registry) error {`,
		`	if err := validate.Required("origin", "body", m.Origin); err != nil {`,
		`	if err := m.validateOriginEnum("origin", "body", *m.Origin); err != nil {`,
		`var additionalThingsTypeStatusPropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"OK\",\"KO\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		additionalThingsTypeStatusPropEnum = append(additionalThingsTypeStatusPropEnum, v`,
		`	AdditionalThingsStatusOK string = "OK"`,
		`	AdditionalThingsStatusKO string = "KO"`,
		`func (m *AdditionalThings) validateStatusEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, additionalThingsTypeStatusPropEnum, true); err != nil {`,
		`func (m *AdditionalThings) validateStatus(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Status) {`,
		`	if err := m.validateStatusEnum("status", "body", m.Status); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_things.go", flattenRun.ExpectedFor("AdditionalThings").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: transitive_refed_thing_additional_properties.go
	flattenRun.AddExpectations("transitive_refed_thing_additional_properties.go", []string{
		`type TransitiveRefedThingAdditionalProperties struct {`,
		"	A1 strfmt.DateTime `json:\"a1,omitempty\"`",
		"	TransitiveRefedThingAdditionalProperties map[string]*NoValidationThing `json:\"-\"`",
		`func (m *TransitiveRefedThingAdditionalProperties) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateA1(formats); err != nil {`,
		`	for k := range m.TransitiveRefedThingAdditionalProperties {`,
		`		if val, ok := m.TransitiveRefedThingAdditionalProperties[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *TransitiveRefedThingAdditionalProperties) validateA1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.A1) {`,
		`	if err := validate.FormatOf("a1", "body", "date-time", m.A1.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_object.go
	flattenRun.AddExpectations("additional_object.go", []string{
		`type AdditionalObject struct {`,
		"	MockID float64 `json:\"mockId,omitempty\"`",
		"	AdditionalObject map[string]*AdditionalObjectAdditionalProperties `json:\"-\"`",
		`func (m *AdditionalObject) Validate(formats strfmt.Registry) error {`,
		`	for k := range m.AdditionalObject {`,
		`		if val, ok := m.AdditionalObject[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_object.go", []string{
		`type AdditionalObject struct {`,
		"	MockID float64 `json:\"mockId,omitempty\"`",
		"	AdditionalObject map[string]*AdditionalObjectAnon `json:\"-\"`",
		`func (m *AdditionalObject) Validate(formats strfmt.Registry) error {`,
		`	for k := range m.AdditionalObject {`,
		`		if val, ok := m.AdditionalObject[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`type AdditionalObjectAnon struct {`,
		"	MockA string `json:\"mockA,omitempty\"`",
		"	MockB *string `json:\"mockB\"`",
		"	MockC float64 `json:\"mockC,omitempty\"`",
		`func (m *AdditionalObjectAnon) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMockA(formats); err != nil {`,
		`	if err := m.validateMockB(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalObjectAnon) validateMockA(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.MockA) {`,
		"	if err := validate.Pattern(\"mockA\", \"body\", m.MockA, `^[A-Z]$`); err != nil {",
		`func (m *AdditionalObjectAnon) validateMockB(formats strfmt.Registry) error {`,
		`	if err := validate.Required("mockB", "body", m.MockB); err != nil {`,
		`	if err := validate.MinLength("mockB", "body", *m.MockB, 1); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_slice_of_objects_additional_properties_items.go
	flattenRun.AddExpectations("additional_slice_of_objects_additional_properties_items.go", []string{
		`type AdditionalSliceOfObjectsAdditionalPropertiesItems struct {`,
		"	Prop2 int64 `json:\"prop2,omitempty\"`",
		// empty validation
		"func (m *AdditionalSliceOfObjectsAdditionalPropertiesItems) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_slice_of_aliased_nullable_primitives.go
	flattenRun.AddExpectations("additional_slice_of_aliased_nullable_primitives.go", []string{
		`type AdditionalSliceOfAliasedNullablePrimitives struct {`,
		"	Prop3 strfmt.UUID `json:\"prop3,omitempty\"`",
		"	AdditionalSliceOfAliasedNullablePrimitives map[string][]*AliasedNullableDate `json:\"-\"`",
		`func (m *AdditionalSliceOfAliasedNullablePrimitives) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp3(formats); err != nil {`,
		`	for k := range m.AdditionalSliceOfAliasedNullablePrimitives {`,
		// removed undue IsSzero call
		`		iAdditionalSliceOfAliasedNullablePrimitivesSize := int64(len(m.AdditionalSliceOfAliasedNullablePrimitives[k])`,
		`		if err := validate.MinItems(k, "body", iAdditionalSliceOfAliasedNullablePrimitivesSize, 10); err != nil {`,
		`		for i := 0; i < len(m.AdditionalSliceOfAliasedNullablePrimitives[k]); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`			if swag.IsZero(m.AdditionalSliceOfAliasedNullablePrimitives[k][i]) {`,
		// nullable required:
		`			if m.AdditionalSliceOfAliasedNullablePrimitives[k][i] != nil {`,
		`				if err := m.AdditionalSliceOfAliasedNullablePrimitives[k][i].Validate(formats); err != nil {`,
		`					if ve, ok := err.(*errors.Validation); ok {`,
		`						return ve.ValidateName(k + "." + strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalSliceOfAliasedNullablePrimitives) validateProp3(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop3) {`,
		`	if err := validate.FormatOf("prop3", "body", "uuid", m.Prop3.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_slice_of_aliased_nullable_primitives.go", []string{
		`type AdditionalSliceOfAliasedNullablePrimitives struct {`,
		"	Prop3 strfmt.UUID `json:\"prop3,omitempty\"`",
		"	AdditionalSliceOfAliasedNullablePrimitives map[string][]*strfmt.Date `json:\"-\"`",
		`func (m *AdditionalSliceOfAliasedNullablePrimitives) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp3(formats); err != nil {`,
		`	for k := range m.AdditionalSliceOfAliasedNullablePrimitives {`,
		`		iAdditionalSliceOfAliasedNullablePrimitivesSize := int64(len(m.AdditionalSliceOfAliasedNullablePrimitives[k])`,
		`		if err := validate.MinItems(k, "body", iAdditionalSliceOfAliasedNullablePrimitivesSize, 10); err != nil {`,
		`		for i := 0; i < len(m.AdditionalSliceOfAliasedNullablePrimitives[k]); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`			if swag.IsZero(m.AdditionalSliceOfAliasedNullablePrimitives[k][i]) {`,
		// nullable required:
		`			if err := validate.FormatOf(k+"."+strconv.Itoa(i), "body", "date", m.AdditionalSliceOfAliasedNullablePrimitives[k][i].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalSliceOfAliasedNullablePrimitives) validateProp3(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop3) {`,
		`	if err := validate.FormatOf("prop3", "body", "uuid", m.Prop3.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_slice_of_slice.go
	flattenRun.AddExpectations("additional_slice_of_slice.go", []string{
		`type AdditionalSliceOfSlice struct {`,
		"	Prop4 strfmt.UUID `json:\"prop4,omitempty\"`",
		"	AdditionalSliceOfSlice map[string][][]*AdditionalSliceOfSliceAdditionalPropertiesItemsItems `json:\"-\"`",
		`func (m *AdditionalSliceOfSlice) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp4(formats); err != nil {`,
		`	for k := range m.AdditionalSliceOfSlice {`,
		`		if err := validate.Required(k, "body", m.AdditionalSliceOfSlice[k]); err != nil {`,
		`		for i := 0; i < len(m.AdditionalSliceOfSlice[k]); i++ {`,
		`			iiAdditionalSliceOfSliceSize := int64(len(m.AdditionalSliceOfSlice[k][i])`,
		`			if err := validate.MaxItems(k+"."+strconv.Itoa(i), "body", iiAdditionalSliceOfSliceSize, 10); err != nil {`,
		`			for ii := 0; ii < len(m.AdditionalSliceOfSlice[k][i]); ii++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`				if swag.IsZero(m.AdditionalSliceOfSlice[k][i][ii]) {`,
		// nullable not required:
		`				if m.AdditionalSliceOfSlice[k][i][ii] != nil {`,
		`					if err := m.AdditionalSliceOfSlice[k][i][ii].Validate(formats); err != nil {`,
		`						if ve, ok := err.(*errors.Validation); ok {`,
		`							return ve.ValidateName(k + "." + strconv.Itoa(i) + "." + strconv.Itoa(ii)`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalSliceOfSlice) validateProp4(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop4) {`,
		`	if err := validate.FormatOf("prop4", "body", "uuid", m.Prop4.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_slice_of_slice.go", []string{
		`type AdditionalSliceOfSlice struct {`,
		"	Prop4 strfmt.UUID `json:\"prop4,omitempty\"`",
		"	AdditionalSliceOfSlice map[string][][]*AdditionalSliceOfSliceItems0 `json:\"-\"`",
		`func (m *AdditionalSliceOfSlice) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp4(formats); err != nil {`,
		`	for k := range m.AdditionalSliceOfSlice {`,
		`		if err := validate.Required(k, "body", m.AdditionalSliceOfSlice[k]); err != nil {`,
		`		for i := 0; i < len(m.AdditionalSliceOfSlice[k]); i++ {`,
		`			iiAdditionalSliceOfSliceSize := int64(len(m.AdditionalSliceOfSlice[k][i])`,
		`			if err := validate.MaxItems(k+"."+strconv.Itoa(i), "body", iiAdditionalSliceOfSliceSize, 10); err != nil {`,
		`			for ii := 0; ii < len(m.AdditionalSliceOfSlice[k][i]); ii++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`				if swag.IsZero(m.AdditionalSliceOfSlice[k][i][ii]) {`,
		// nullable required:
		`				if m.AdditionalSliceOfSlice[k][i][ii] != nil {`,
		`					if err := m.AdditionalSliceOfSlice[k][i][ii].Validate(formats); err != nil {`,
		`						if ve, ok := err.(*errors.Validation); ok {`,
		`							return ve.ValidateName(k + "." + strconv.Itoa(i) + "." + strconv.Itoa(ii)`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalSliceOfSlice) validateProp4(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop4) {`,
		`	if err := validate.FormatOf("prop4", "body", "uuid", m.Prop4.String(), formats); err != nil {`,
		`type AdditionalSliceOfSliceItems0 struct {`,
		"	Prop5 int64 `json:\"prop5,omitempty\"`",
		`func (m *AdditionalSliceOfSliceItems0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp5(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalSliceOfSliceItems0) validateProp5(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop5) {`,
		`	if err := validate.MaximumInt("prop5", "body", m.Prop5, 10, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_object_with_aliased_thing.go
	flattenRun.AddExpectations("additional_object_with_aliased_thing.go", []string{
		`type AdditionalObjectWithAliasedThing struct {`,
		"	Blob int64 `json:\"blob,omitempty\"`",
		"	AdditionalObjectWithAliasedThing map[string]AliasedDate `json:\"-\"`",
		`func (m *AdditionalObjectWithAliasedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateBlob(formats); err != nil {`,
		`	for k := range m.AdditionalObjectWithAliasedThing {`,
		// removed undue IsZero call
		`		if val, ok := m.AdditionalObjectWithAliasedThing[k]; ok {`,
		`			if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalObjectWithAliasedThing) validateBlob(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Blob) {`,
		`	if err := validate.MinimumInt("blob", "body", m.Blob, 1, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_object_with_aliased_thing.go", []string{
		`type AdditionalObjectWithAliasedThing struct {`,
		"	Blob int64 `json:\"blob,omitempty\"`",
		"	AdditionalObjectWithAliasedThing map[string]strfmt.Date `json:\"-\"`",
		`func (m *AdditionalObjectWithAliasedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateBlob(formats); err != nil {`,
		`	for k := range m.AdditionalObjectWithAliasedThing {`,
		`		if err := validate.FormatOf(k, "body", "date", m.AdditionalObjectWithAliasedThing[k].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalObjectWithAliasedThing) validateBlob(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Blob) {`,
		`	if err := validate.MinimumInt("blob", "body", m.Blob, 1, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_things_nested_additional_properties.go
	flattenRun.AddExpectations("additional_things_nested_additional_properties.go", []string{
		`type AdditionalThingsNestedAdditionalProperties struct {`,
		"	PrinterAddress string `json:\"printerAddress,omitempty\"`",
		"	PrinterCountry string `json:\"printerCountry,omitempty\"`",
		"	PrinterDate strfmt.Date `json:\"printerDate,omitempty\"`",
		"	AdditionalThingsNestedAdditionalProperties map[string]*AdditionalThingsNestedAdditionalPropertiesAdditionalProperties `json:\"-\"`",
		`func (m *AdditionalThingsNestedAdditionalProperties) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validatePrinterCountry(formats); err != nil {`,
		`	if err := m.validatePrinterDate(formats); err != nil {`,
		`	for k := range m.AdditionalThingsNestedAdditionalProperties {`,
		`		if val, ok := m.AdditionalThingsNestedAdditionalProperties[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`var additionalThingsNestedAdditionalPropertiesTypePrinterCountryPropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"US\",\"FR\",\"UK\",\"BE\",\"CA\",\"DE\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		additionalThingsNestedAdditionalPropertiesTypePrinterCountryPropEnum = append(additionalThingsNestedAdditionalPropertiesTypePrinterCountryPropEnum, v`,
		`	AdditionalThingsNestedAdditionalPropertiesPrinterCountryUS string = "US"`,
		`	AdditionalThingsNestedAdditionalPropertiesPrinterCountryFR string = "FR"`,
		`	AdditionalThingsNestedAdditionalPropertiesPrinterCountryUK string = "UK"`,
		`	AdditionalThingsNestedAdditionalPropertiesPrinterCountryBE string = "BE"`,
		`	AdditionalThingsNestedAdditionalPropertiesPrinterCountryCA string = "CA"`,
		`	AdditionalThingsNestedAdditionalPropertiesPrinterCountryDE string = "DE"`,
		`func (m *AdditionalThingsNestedAdditionalProperties) validatePrinterCountryEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, additionalThingsNestedAdditionalPropertiesTypePrinterCountryPropEnum, true); err != nil {`,
		`func (m *AdditionalThingsNestedAdditionalProperties) validatePrinterCountry(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.PrinterCountry) {`,
		`	if err := m.validatePrinterCountryEnum("printerCountry", "body", m.PrinterCountry); err != nil {`,
		`func (m *AdditionalThingsNestedAdditionalProperties) validatePrinterDate(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.PrinterDate) {`,
		`	if err := validate.FormatOf("printerDate", "body", "date", m.PrinterDate.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: empty_object_with_additional_slice_additional_properties_items.go
	flattenRun.AddExpectations("empty_object_with_additional_slice_additional_properties_items.go", []string{
		`type EmptyObjectWithAdditionalSliceAdditionalPropertiesItems struct {`,
		"	DummyProp1 strfmt.Date `json:\"dummyProp1,omitempty\"`",
		`func (m *EmptyObjectWithAdditionalSliceAdditionalPropertiesItems) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateDummyProp1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *EmptyObjectWithAdditionalSliceAdditionalPropertiesItems) validateDummyProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.DummyProp1) {`,
		`	if err := validate.FormatOf("dummyProp1", "body", "date", m.DummyProp1.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_things_nested_additional_properties_additional_properties.go
	flattenRun.AddExpectations("additional_things_nested_additional_properties_additional_properties.go", []string{
		`type AdditionalThingsNestedAdditionalPropertiesAdditionalProperties struct {`,
		"	AverageDelay strfmt.Duration `json:\"averageDelay,omitempty\"`",
		`func (m *AdditionalThingsNestedAdditionalPropertiesAdditionalProperties) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAverageDelay(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalThingsNestedAdditionalPropertiesAdditionalProperties) validateAverageDelay(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.AverageDelay) {`,
		`	if err := validate.FormatOf("averageDelay", "body", "duration", m.AverageDelay.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_slice_of_slice_additional_properties_items_items.go
	flattenRun.AddExpectations("additional_slice_of_slice_additional_properties_items_items.go", []string{
		`type AdditionalSliceOfSliceAdditionalPropertiesItemsItems struct {`,
		"	Prop5 int64 `json:\"prop5,omitempty\"`",
		`func (m *AdditionalSliceOfSliceAdditionalPropertiesItemsItems) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp5(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalSliceOfSliceAdditionalPropertiesItemsItems) validateProp5(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop5) {`,
		`	if err := validate.MaximumInt("prop5", "body", m.Prop5, 10, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_object_additional_properties.go
	flattenRun.AddExpectations("additional_object_additional_properties.go", []string{
		`type AdditionalObjectAdditionalProperties struct {`,
		"	MockA string `json:\"mockA,omitempty\"`",
		"	MockB *string `json:\"mockB\"`",
		"	MockC float64 `json:\"mockC,omitempty\"`",
		`func (m *AdditionalObjectAdditionalProperties) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMockA(formats); err != nil {`,
		`	if err := m.validateMockB(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalObjectAdditionalProperties) validateMockA(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.MockA) {`,
		"	if err := validate.Pattern(\"mockA\", \"body\", m.MockA, `^[A-Z]$`); err != nil {",
		`func (m *AdditionalObjectAdditionalProperties) validateMockB(formats strfmt.Registry) error {`,
		`	if err := validate.Required("mockB", "body", m.MockB); err != nil {`,
		`	if err := validate.MinLength("mockB", "body", *m.MockB, 1); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_transitive_refed_thing.go
	flattenRun.AddExpectations("additional_transitive_refed_thing.go", []string{
		`type AdditionalTransitiveRefedThing struct {`,
		"	ThisOneNotRequired int64 `json:\"thisOneNotRequired,omitempty\"`",
		"	AdditionalTransitiveRefedThing map[string][]*TransitiveRefedThing `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequired(formats); err != nil {`,
		`	for k := range m.AdditionalTransitiveRefedThing {`,
		`		if err := validate.Required(k, "body", m.AdditionalTransitiveRefedThing[k]); err != nil {`,
		`		if err := validate.UniqueItems(k, "body", m.AdditionalTransitiveRefedThing[k]); err != nil {`,
		`		for i := 0; i < len(m.AdditionalTransitiveRefedThing[k]); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`			if swag.IsZero(m.AdditionalTransitiveRefedThing[k][i]) {`,
		// nullable required:
		`			if m.AdditionalTransitiveRefedThing[k][i] != nil {`,
		`				if err := m.AdditionalTransitiveRefedThing[k][i].Validate(formats); err != nil {`,
		`					if ve, ok := err.(*errors.Validation); ok {`,
		`						return ve.ValidateName(k + "." + strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedThing) validateThisOneNotRequired(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequired) {`,
		`	if err := validate.MaximumInt("thisOneNotRequired", "body", m.ThisOneNotRequired, 10, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_transitive_refed_thing.go", []string{
		`type AdditionalTransitiveRefedThing struct {`,
		"	ThisOneNotRequired int64 `json:\"thisOneNotRequired,omitempty\"`",
		"	AdditionalTransitiveRefedThing map[string][]*AdditionalTransitiveRefedThingItems0 `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequired(formats); err != nil {`,
		`	for k := range m.AdditionalTransitiveRefedThing {`,
		`		if err := validate.Required(k, "body", m.AdditionalTransitiveRefedThing[k]); err != nil {`,
		`		if err := validate.UniqueItems(k, "body", m.AdditionalTransitiveRefedThing[k]); err != nil {`,
		`		for i := 0; i < len(m.AdditionalTransitiveRefedThing[k]); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`			if swag.IsZero(m.AdditionalTransitiveRefedThing[k][i]) {`,
		// nullable required:
		`			if m.AdditionalTransitiveRefedThing[k][i] != nil {`,
		`				if err := m.AdditionalTransitiveRefedThing[k][i].Validate(formats); err != nil {`,
		`					if ve, ok := err.(*errors.Validation); ok {`,
		`						return ve.ValidateName(k + "." + strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedThing) validateThisOneNotRequired(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequired) {`,
		`	if err := validate.MaximumInt("thisOneNotRequired", "body", m.ThisOneNotRequired, 10, false); err != nil {`,
		`type AdditionalTransitiveRefedThingItems0 struct {`,
		"	ThisOneNotRequiredEither int64 `json:\"thisOneNotRequiredEither,omitempty\"`",
		"	AdditionalTransitiveRefedThingItems0 map[string]*AdditionalTransitiveRefedThingItems0Anon `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedThingItems0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequiredEither(formats); err != nil {`,
		`	for k := range m.AdditionalTransitiveRefedThingItems0 {`,
		`		if val, ok := m.AdditionalTransitiveRefedThingItems0[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedThingItems0) validateThisOneNotRequiredEither(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequiredEither) {`,
		`	if err := validate.MaximumInt("thisOneNotRequiredEither", "body", m.ThisOneNotRequiredEither, 20, false); err != nil {`,
		`type AdditionalTransitiveRefedThingItems0Anon struct {`,
		"	A1 strfmt.DateTime `json:\"a1,omitempty\"`",
		"	AdditionalTransitiveRefedThingItems0Anon map[string]*AdditionalTransitiveRefedThingItems0AnonAnon `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedThingItems0Anon) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateA1(formats); err != nil {`,
		`	for k := range m.AdditionalTransitiveRefedThingItems0Anon {`,
		`		if val, ok := m.AdditionalTransitiveRefedThingItems0Anon[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedThingItems0Anon) validateA1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.A1) {`,
		`	if err := validate.FormatOf("a1", "body", "date-time", m.A1.String(), formats); err != nil {`,
		`type AdditionalTransitiveRefedThingItems0AnonAnon struct {`,
		"	Discourse string `json:\"discourse,omitempty\"`",
		"	HoursSpent float64 `json:\"hoursSpent,omitempty\"`",
		"	AdditionalTransitiveRefedThingItems0AnonAnonAdditionalProperties map[string]interface{} `json:\"-\"`",
		// empty validation
		"func (m *AdditionalTransitiveRefedThingItems0AnonAnon) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_nullable_array_thing.go
	flattenRun.AddExpectations("additional_nullable_array_thing.go", []string{
		`type AdditionalNullableArrayThing struct {`,
		"	ThisOneNotRequired int64 `json:\"thisOneNotRequired,omitempty\"`",
		"	AdditionalNullableArrayThing map[string][]strfmt.ISBN `json:\"-\"`",
		`func (m *AdditionalNullableArrayThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequired(formats); err != nil {`,
		`	for k := range m.AdditionalNullableArrayThing {`,
		`		if err := validate.UniqueItems(k, "body", m.AdditionalNullableArrayThing[k]); err != nil {`,
		`		for i := 0; i < len(m.AdditionalNullableArrayThing[k]); i++ {`,
		`			if err := validate.FormatOf(k+"."+strconv.Itoa(i), "body", "isbn", m.AdditionalNullableArrayThing[k][i].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalNullableArrayThing) validateThisOneNotRequired(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequired) {`,
		`	if err := validate.MaximumInt("thisOneNotRequired", "body", m.ThisOneNotRequired, 10, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_nullable_array_thing.go", flattenRun.ExpectedFor("AdditionalNullableArrayThing").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: additional_slice_of_primitives.go
	flattenRun.AddExpectations("additional_slice_of_primitives.go", []string{
		`type AdditionalSliceOfPrimitives struct {`,
		"	Prop1 string `json:\"prop1,omitempty\"`",
		"	AdditionalSliceOfPrimitives map[string][]strfmt.Date `json:\"-\"`",
		`func (m *AdditionalSliceOfPrimitives) Validate(formats strfmt.Registry) error {`,
		`	for k := range m.AdditionalSliceOfPrimitives {`,
		`		if err := validate.UniqueItems(k, "body", m.AdditionalSliceOfPrimitives[k]); err != nil {`,
		`		for i := 0; i < len(m.AdditionalSliceOfPrimitives[k]); i++ {`,
		`			if err := validate.FormatOf(k+"."+strconv.Itoa(i), "body", "date", m.AdditionalSliceOfPrimitives[k][i].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_slice_of_primitives.go", flattenRun.ExpectedFor("AdditionalSliceOfPrimitives").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: additional_array_thing.go
	flattenRun.AddExpectations("additional_array_thing.go", []string{
		`type AdditionalArrayThing struct {`,
		"	ThisOneNotRequired int64 `json:\"thisOneNotRequired,omitempty\"`",
		"	AdditionalArrayThing map[string][]strfmt.UUID `json:\"-\"`",
		`func (m *AdditionalArrayThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequired(formats); err != nil {`,
		`	for k := range m.AdditionalArrayThing {`,
		`		if err := validate.UniqueItems(k, "body", m.AdditionalArrayThing[k]); err != nil {`,
		`		for i := 0; i < len(m.AdditionalArrayThing[k]); i++ {`,
		`			if err := validate.FormatOf(k+"."+strconv.Itoa(i), "body", "uuid", m.AdditionalArrayThing[k][i].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalArrayThing) validateThisOneNotRequired(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequired) {`,
		`	if err := validate.MaximumInt("thisOneNotRequired", "body", m.ThisOneNotRequired, 10, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_array_thing.go", flattenRun.ExpectedFor("AdditionalArrayThing").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: interface_thing.go
	flattenRun.AddExpectations("interface_thing.go", []string{
		`type InterfaceThing interface{}`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("interface_thing.go", flattenRun.ExpectedFor("InterfaceThing").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: empty_object_with_additional_slice.go
	flattenRun.AddExpectations("empty_object_with_additional_slice.go", []string{
		`type EmptyObjectWithAdditionalSlice map[string][]EmptyObjectWithAdditionalSliceAdditionalPropertiesItems`,
		`func (m EmptyObjectWithAdditionalSlice) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if err := validate.Required(k, "body", m[k]); err != nil {`,
		`		for i := 0; i < len(m[k]); i++ {`,
		`			if err := m[k][i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName(k + "." + strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("empty_object_with_additional_slice.go", []string{
		`type EmptyObjectWithAdditionalSlice map[string][]EmptyObjectWithAdditionalSliceItems0`,
		`func (m EmptyObjectWithAdditionalSlice) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if err := validate.Required(k, "body", m[k]); err != nil {`,
		`		for i := 0; i < len(m[k]); i++ {`,
		`			if err := m[k][i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName(k + "." + strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
		`type EmptyObjectWithAdditionalSliceItems0 struct {`,
		"	DummyProp1 strfmt.Date `json:\"dummyProp1,omitempty\"`",
		`func (m *EmptyObjectWithAdditionalSliceItems0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateDummyProp1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *EmptyObjectWithAdditionalSliceItems0) validateDummyProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.DummyProp1) {`,
		`	if err := validate.FormatOf("dummyProp1", "body", "date", m.DummyProp1.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_slice_of_objects.go
	flattenRun.AddExpectations("additional_slice_of_objects.go", []string{
		`type AdditionalSliceOfObjects struct {`,
		"	Prop1 string `json:\"prop1,omitempty\"`",
		"	AdditionalSliceOfObjects map[string][]*AdditionalSliceOfObjectsAdditionalPropertiesItems `json:\"-\"`",
		`func (m *AdditionalSliceOfObjects) Validate(formats strfmt.Registry) error {`,
		`	for k := range m.AdditionalSliceOfObjects {`,
		`		if err := validate.Required(k, "body", m.AdditionalSliceOfObjects[k]); err != nil {`,
		`		if err := validate.UniqueItems(k, "body", m.AdditionalSliceOfObjects[k]); err != nil {`,
		`		for i := 0; i < len(m.AdditionalSliceOfObjects[k]); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`			if swag.IsZero(m.AdditionalSliceOfObjects[k][i]) {`,
		// nullable required:
		`			if m.AdditionalSliceOfObjects[k][i] != nil {`,
		`				if err := m.AdditionalSliceOfObjects[k][i].Validate(formats); err != nil {`,
		`					if ve, ok := err.(*errors.Validation); ok {`,
		`						return ve.ValidateName(k + "." + strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_slice_of_objects.go", []string{
		`type AdditionalSliceOfObjects struct {`,
		"	Prop1 string `json:\"prop1,omitempty\"`",
		"	AdditionalSliceOfObjects map[string][]*AdditionalSliceOfObjectsItems0 `json:\"-\"`",
		`func (m *AdditionalSliceOfObjects) Validate(formats strfmt.Registry) error {`,
		`	for k := range m.AdditionalSliceOfObjects {`,
		`		if err := validate.Required(k, "body", m.AdditionalSliceOfObjects[k]); err != nil {`,
		`		if err := validate.UniqueItems(k, "body", m.AdditionalSliceOfObjects[k]); err != nil {`,
		`		for i := 0; i < len(m.AdditionalSliceOfObjects[k]); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`			if swag.IsZero(m.AdditionalSliceOfObjects[k][i]) {`,
		// nullable required:
		`			if m.AdditionalSliceOfObjects[k][i] != nil {`,
		`				if err := m.AdditionalSliceOfObjects[k][i].Validate(formats); err != nil {`,
		`					if ve, ok := err.(*errors.Validation); ok {`,
		`						return ve.ValidateName(k + "." + strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
		`type AdditionalSliceOfObjectsItems0 struct {`,
		"	Prop2 int64 `json:\"prop2,omitempty\"`",
		// empty validation
		"func (m *AdditionalSliceOfObjectsItems0) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_things_nested.go
	flattenRun.AddExpectations("additional_things_nested.go", []string{
		`type AdditionalThingsNested struct {`,
		"	Origin string `json:\"origin,omitempty\"`",
		"	AdditionalThingsNested map[string]*AdditionalThingsNestedAdditionalProperties `json:\"-\"`",
		`func (m *AdditionalThingsNested) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateOrigin(formats); err != nil {`,
		`	for k := range m.AdditionalThingsNested {`,
		`		if val, ok := m.AdditionalThingsNested[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`var additionalThingsNestedTypeOriginPropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"goPrint\",\"goE-book\",\"goCollection\",\"goMuseum\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		additionalThingsNestedTypeOriginPropEnum = append(additionalThingsNestedTypeOriginPropEnum, v`,
		`	AdditionalThingsNestedOriginGoPrint string = "goPrint"`,
		`	AdditionalThingsNestedOriginGoEDashBook string = "goE-book"`,
		`	AdditionalThingsNestedOriginGoCollection string = "goCollection"`,
		`	AdditionalThingsNestedOriginGoMuseum string = "goMuseum"`,
		`func (m *AdditionalThingsNested) validateOriginEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, additionalThingsNestedTypeOriginPropEnum, true); err != nil {`,
		`func (m *AdditionalThingsNested) validateOrigin(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Origin) {`,
		`	if err := m.validateOriginEnum("origin", "body", m.Origin); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_things_nested.go", []string{
		`type AdditionalThingsNested struct {`,
		"	Origin string `json:\"origin,omitempty\"`",
		"	AdditionalThingsNested map[string]*AdditionalThingsNestedAnon `json:\"-\"`",
		`func (m *AdditionalThingsNested) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateOrigin(formats); err != nil {`,
		`	for k := range m.AdditionalThingsNested {`,
		`		if val, ok := m.AdditionalThingsNested[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`var additionalThingsNestedTypeOriginPropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"goPrint\",\"goE-book\",\"goCollection\",\"goMuseum\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		additionalThingsNestedTypeOriginPropEnum = append(additionalThingsNestedTypeOriginPropEnum, v`,
		`	AdditionalThingsNestedOriginGoPrint string = "goPrint"`,
		`	AdditionalThingsNestedOriginGoEDashBook string = "goE-book"`,
		`	AdditionalThingsNestedOriginGoCollection string = "goCollection"`,
		`	AdditionalThingsNestedOriginGoMuseum string = "goMuseum"`,
		`func (m *AdditionalThingsNested) validateOriginEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, additionalThingsNestedTypeOriginPropEnum, true); err != nil {`,
		`func (m *AdditionalThingsNested) validateOrigin(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Origin) {`,
		`	if err := m.validateOriginEnum("origin", "body", m.Origin); err != nil {`,
		`type AdditionalThingsNestedAnon struct {`,
		"	PrinterAddress string `json:\"printerAddress,omitempty\"`",
		"	PrinterCountry string `json:\"printerCountry,omitempty\"`",
		"	PrinterDate strfmt.Date `json:\"printerDate,omitempty\"`",
		"	AdditionalThingsNestedAnon map[string]*AdditionalThingsNestedAnonAnon `json:\"-\"`",
		`func (m *AdditionalThingsNestedAnon) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validatePrinterCountry(formats); err != nil {`,
		`	if err := m.validatePrinterDate(formats); err != nil {`,
		`	for k := range m.AdditionalThingsNestedAnon {`,
		`		if val, ok := m.AdditionalThingsNestedAnon[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`var additionalThingsNestedAnonTypePrinterCountryPropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"US\",\"FR\",\"UK\",\"BE\",\"CA\",\"DE\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		additionalThingsNestedAnonTypePrinterCountryPropEnum = append(additionalThingsNestedAnonTypePrinterCountryPropEnum, v`,
		`	AdditionalThingsNestedAnonPrinterCountryUS string = "US"`,
		`	AdditionalThingsNestedAnonPrinterCountryFR string = "FR"`,
		`	AdditionalThingsNestedAnonPrinterCountryUK string = "UK"`,
		`	AdditionalThingsNestedAnonPrinterCountryBE string = "BE"`,
		`	AdditionalThingsNestedAnonPrinterCountryCA string = "CA"`,
		`	AdditionalThingsNestedAnonPrinterCountryDE string = "DE"`,
		`func (m *AdditionalThingsNestedAnon) validatePrinterCountryEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, additionalThingsNestedAnonTypePrinterCountryPropEnum, true); err != nil {`,
		`func (m *AdditionalThingsNestedAnon) validatePrinterCountry(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.PrinterCountry) {`,
		`	if err := m.validatePrinterCountryEnum("printerCountry", "body", m.PrinterCountry); err != nil {`,
		`func (m *AdditionalThingsNestedAnon) validatePrinterDate(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.PrinterDate) {`,
		`	if err := validate.FormatOf("printerDate", "body", "date", m.PrinterDate.String(), formats); err != nil {`,
		`type AdditionalThingsNestedAnonAnon struct {`,
		"	AverageDelay strfmt.Duration `json:\"averageDelay,omitempty\"`",
		`func (m *AdditionalThingsNestedAnonAnon) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAverageDelay(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalThingsNestedAnonAnon) validateAverageDelay(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.AverageDelay) {`,
		`	if err := validate.FormatOf("averageDelay", "body", "duration", m.AverageDelay.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: no_validation_thing.go
	flattenRun.AddExpectations("no_validation_thing.go", []string{
		`type NoValidationThing struct {`,
		"	Discourse string `json:\"discourse,omitempty\"`",
		"	HoursSpent float64 `json:\"hoursSpent,omitempty\"`",
		"	NoValidationThingAdditionalProperties map[string]interface{} `json:\"-\"`",
		// empty validation
		"func (m *NoValidationThing) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("no_validation_thing.go", flattenRun.ExpectedFor("NoValidationThing").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: additional_array_of_interface.go
	flattenRun.AddExpectations("additional_array_of_interface.go", []string{
		`type AdditionalArrayOfInterface struct {`,
		"	ThisOneNotRequired int64 `json:\"thisOneNotRequired,omitempty\"`",
		"	AdditionalArrayOfInterface map[string][]interface{} `json:\"-\"`",
		`func (m *AdditionalArrayOfInterface) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequired(formats); err != nil {`,
		`	for k := range m.AdditionalArrayOfInterface {`,
		// remove undue IsZero call
		`		if err := validate.UniqueItems(k, "body", m.AdditionalArrayOfInterface[k]); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalArrayOfInterface) validateThisOneNotRequired(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequired) {`,
		`	if err := validate.MaximumInt("thisOneNotRequired", "body", m.ThisOneNotRequired, 10, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_array_of_interface.go", flattenRun.ExpectedFor("AdditionalArrayOfInterface").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: additional_formated_thing.go
	flattenRun.AddExpectations("additional_formated_thing.go", []string{
		`type AdditionalFormatedThing map[string]strfmt.Date`,
		`func (m AdditionalFormatedThing) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if err := validate.FormatOf(k, "body", "date", m[k].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_formated_thing.go", flattenRun.ExpectedFor("AdditionalFormatedThing").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: aliased_nullable_date.go
	flattenRun.AddExpectations("aliased_nullable_date.go", []string{
		`type AliasedNullableDate strfmt.Date`,
		`func (m AliasedNullableDate) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.FormatOf("", "body", "date", strfmt.Date(m).String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("aliased_nullable_date.go", flattenRun.ExpectedFor("AliasedNullableDate").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: additional_array_of_refed_object.go
	flattenRun.AddExpectations("additional_array_of_refed_object.go", []string{
		`type AdditionalArrayOfRefedObject struct {`,
		"	ThisOneNotRequired int64 `json:\"thisOneNotRequired,omitempty\"`",
		"	AdditionalArrayOfRefedObject map[string][]*NoValidationThing `json:\"-\"`",
		`func (m *AdditionalArrayOfRefedObject) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequired(formats); err != nil {`,
		`	for k := range m.AdditionalArrayOfRefedObject {`,
		`		if err := validate.Required(k, "body", m.AdditionalArrayOfRefedObject[k]); err != nil {`,
		`		if err := validate.UniqueItems(k, "body", m.AdditionalArrayOfRefedObject[k]); err != nil {`,
		`		for i := 0; i < len(m.AdditionalArrayOfRefedObject[k]); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`			if swag.IsZero(m.AdditionalArrayOfRefedObject[k][i]) {`,
		// nullable required:
		`			if m.AdditionalArrayOfRefedObject[k][i] != nil {`,
		`				if err := m.AdditionalArrayOfRefedObject[k][i].Validate(formats); err != nil {`,
		`					if ve, ok := err.(*errors.Validation); ok {`,
		`						return ve.ValidateName(k + "." + strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalArrayOfRefedObject) validateThisOneNotRequired(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequired) {`,
		`	if err := validate.MaximumInt("thisOneNotRequired", "body", m.ThisOneNotRequired, 10, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_array_of_refed_object.go", []string{
		`type AdditionalArrayOfRefedObject struct {`,
		"	ThisOneNotRequired int64 `json:\"thisOneNotRequired,omitempty\"`",
		"	AdditionalArrayOfRefedObject map[string][]*AdditionalArrayOfRefedObjectItems0 `json:\"-\"`",
		`func (m *AdditionalArrayOfRefedObject) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequired(formats); err != nil {`,
		`	for k := range m.AdditionalArrayOfRefedObject {`,
		`		if err := validate.Required(k, "body", m.AdditionalArrayOfRefedObject[k]); err != nil {`,
		`		if err := validate.UniqueItems(k, "body", m.AdditionalArrayOfRefedObject[k]); err != nil {`,
		`		for i := 0; i < len(m.AdditionalArrayOfRefedObject[k]); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`			if swag.IsZero(m.AdditionalArrayOfRefedObject[k][i]) {`,
		// nullable required:
		`			if m.AdditionalArrayOfRefedObject[k][i] != nil {`,
		`				if err := m.AdditionalArrayOfRefedObject[k][i].Validate(formats); err != nil {`,
		`					if ve, ok := err.(*errors.Validation); ok {`,
		`						return ve.ValidateName(k + "." + strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalArrayOfRefedObject) validateThisOneNotRequired(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequired) {`,
		`	if err := validate.MaximumInt("thisOneNotRequired", "body", m.ThisOneNotRequired, 10, false); err != nil {`,
		`type AdditionalArrayOfRefedObjectItems0 struct {`,
		"	Discourse string `json:\"discourse,omitempty\"`",
		"	HoursSpent float64 `json:\"hoursSpent,omitempty\"`",
		"	AdditionalArrayOfRefedObjectItems0AdditionalProperties map[string]interface{} `json:\"-\"`",
		// empty validation
		"func (m *AdditionalArrayOfRefedObjectItems0) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_slice_of_aliased_primitives.go
	flattenRun.AddExpectations("additional_slice_of_aliased_primitives.go", []string{
		`type AdditionalSliceOfAliasedPrimitives struct {`,
		"	Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		"	AdditionalSliceOfAliasedPrimitives map[string][]AliasedDate `json:\"-\"`",
		`func (m *AdditionalSliceOfAliasedPrimitives) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp2(formats); err != nil {`,
		`	for k := range m.AdditionalSliceOfAliasedPrimitives {`,
		// removed undue IsZero call
		`		iAdditionalSliceOfAliasedPrimitivesSize := int64(len(m.AdditionalSliceOfAliasedPrimitives[k])`,
		`		if err := validate.MaxItems(k, "body", iAdditionalSliceOfAliasedPrimitivesSize, 10); err != nil {`,
		`		for i := 0; i < len(m.AdditionalSliceOfAliasedPrimitives[k]); i++ {`,
		`			if err := m.AdditionalSliceOfAliasedPrimitives[k][i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName(k + "." + strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalSliceOfAliasedPrimitives) validateProp2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop2) {`,
		`	if err := validate.FormatOf("prop2", "body", "uuid", m.Prop2.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_slice_of_aliased_primitives.go", []string{
		`type AdditionalSliceOfAliasedPrimitives struct {`,
		"	Prop2 strfmt.UUID `json:\"prop2,omitempty\"`",
		"	AdditionalSliceOfAliasedPrimitives map[string][]strfmt.Date `json:\"-\"`",
		`func (m *AdditionalSliceOfAliasedPrimitives) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp2(formats); err != nil {`,
		`	for k := range m.AdditionalSliceOfAliasedPrimitives {`,
		`		iAdditionalSliceOfAliasedPrimitivesSize := int64(len(m.AdditionalSliceOfAliasedPrimitives[k])`,
		`		if err := validate.MaxItems(k, "body", iAdditionalSliceOfAliasedPrimitivesSize, 10); err != nil {`,
		`		for i := 0; i < len(m.AdditionalSliceOfAliasedPrimitives[k]); i++ {`,
		`			if err := validate.FormatOf(k+"."+strconv.Itoa(i), "body", "date", m.AdditionalSliceOfAliasedPrimitives[k][i].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalSliceOfAliasedPrimitives) validateProp2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop2) {`,
		`	if err := validate.FormatOf("prop2", "body", "uuid", m.Prop2.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: transitive_refed_thing.go
	flattenRun.AddExpectations("transitive_refed_thing.go", []string{
		`type TransitiveRefedThing struct {`,
		"	ThisOneNotRequiredEither int64 `json:\"thisOneNotRequiredEither,omitempty\"`",
		"	TransitiveRefedThing map[string]*TransitiveRefedThingAdditionalProperties `json:\"-\"`",
		`func (m *TransitiveRefedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequiredEither(formats); err != nil {`,
		`	for k := range m.TransitiveRefedThing {`,
		`		if val, ok := m.TransitiveRefedThing[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *TransitiveRefedThing) validateThisOneNotRequiredEither(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequiredEither) {`,
		`	if err := validate.MaximumInt("thisOneNotRequiredEither", "body", m.ThisOneNotRequiredEither, 20, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("transitive_refed_thing.go", []string{
		`type TransitiveRefedThing struct {`,
		"	ThisOneNotRequiredEither int64 `json:\"thisOneNotRequiredEither,omitempty\"`",
		"	TransitiveRefedThing map[string]*TransitiveRefedThingAnon `json:\"-\"`",
		`func (m *TransitiveRefedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequiredEither(formats); err != nil {`,
		`	for k := range m.TransitiveRefedThing {`,
		`		if val, ok := m.TransitiveRefedThing[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *TransitiveRefedThing) validateThisOneNotRequiredEither(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequiredEither) {`,
		`	if err := validate.MaximumInt("thisOneNotRequiredEither", "body", m.ThisOneNotRequiredEither, 20, false); err != nil {`,
		`type TransitiveRefedThingAnon struct {`,
		"	A1 strfmt.DateTime `json:\"a1,omitempty\"`",
		"	TransitiveRefedThingAnon map[string]*TransitiveRefedThingAnonAnon `json:\"-\"`",
		`func (m *TransitiveRefedThingAnon) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateA1(formats); err != nil {`,
		`	for k := range m.TransitiveRefedThingAnon {`,
		`		if val, ok := m.TransitiveRefedThingAnon[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *TransitiveRefedThingAnon) validateA1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.A1) {`,
		`	if err := validate.FormatOf("a1", "body", "date-time", m.A1.String(), formats); err != nil {`,
		`type TransitiveRefedThingAnonAnon struct {`,
		"	Discourse string `json:\"discourse,omitempty\"`",
		"	HoursSpent float64 `json:\"hoursSpent,omitempty\"`",
		"	TransitiveRefedThingAnonAnonAdditionalProperties map[string]interface{} `json:\"-\"`",
		// empty validation
		"func (m *TransitiveRefedThingAnonAnon) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_empty_object.go
	flattenRun.AddExpectations("additional_empty_object.go", []string{
		`type AdditionalEmptyObject struct {`,
		"	PropA interface{} `json:\"propA,omitempty\"`",
		"	AdditionalEmptyObject map[string]interface{} `json:\"-\"`",
		// empty validation
		"func (m *AdditionalEmptyObject) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_empty_object.go", flattenRun.ExpectedFor("AdditionalEmptyObject").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: additional_date_with_nullable_thing.go
	flattenRun.AddExpectations("additional_date_with_nullable_thing.go", []string{
		`type AdditionalDateWithNullableThing struct {`,
		"	Blob int64 `json:\"blob,omitempty\"`",
		"	NullableDate *AliasedNullableDate `json:\"nullableDate,omitempty\"`",
		"	AdditionalDateWithNullableThing map[string]*AliasedNullableDate `json:\"-\"`",
		`func (m *AdditionalDateWithNullableThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateBlob(formats); err != nil {`,
		`	if err := m.validateNullableDate(formats); err != nil {`,
		`	for k := range m.AdditionalDateWithNullableThing {`,
		`		if swag.IsZero(m.AdditionalDateWithNullableThing[k]) {`,
		`		if val, ok := m.AdditionalDateWithNullableThing[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalDateWithNullableThing) validateBlob(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Blob) {`,
		`	if err := validate.MinimumInt("blob", "body", m.Blob, 1, false); err != nil {`,
		`func (m *AdditionalDateWithNullableThing) validateNullableDate(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.NullableDate) {`,
		`	if m.NullableDate != nil {`,
		`		if err := m.NullableDate.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("nullableDate"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("additional_date_with_nullable_thing.go", []string{
		`type AdditionalDateWithNullableThing struct {`,
		"	Blob int64 `json:\"blob,omitempty\"`",
		"	NullableDate *strfmt.Date `json:\"nullableDate,omitempty\"`",
		"	AdditionalDateWithNullableThing map[string]*strfmt.Date `json:\"-\"`",
		`func (m *AdditionalDateWithNullableThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateBlob(formats); err != nil {`,
		`	if err := m.validateNullableDate(formats); err != nil {`,
		`	for k := range m.AdditionalDateWithNullableThing {`,
		`		if swag.IsZero(m.AdditionalDateWithNullableThing[k]) {`,
		`		if err := validate.FormatOf(k, "body", "date", m.AdditionalDateWithNullableThing[k].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalDateWithNullableThing) validateBlob(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Blob) {`,
		`	if err := validate.MinimumInt("blob", "body", m.Blob, 1, false); err != nil {`,
		`func (m *AdditionalDateWithNullableThing) validateNullableDate(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.NullableDate) {`,
		`	if err := validate.FormatOf("nullableDate", "body", "date", m.NullableDate.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

}

func initFixtureTuple() {
	// testing ../fixtures/bugs/1487/fixture-tuple.yaml with expand (--skip-flatten)

	/* check different patterns of additionalItems validations or absence thereof
	 */
	f := newModelFixture("../fixtures/bugs/1487/fixture-tuple.yaml", "fixture for tuples and additionalItems")
	flattenRun := f.AddRun(false)
	expandRun := f.AddRun(true)

	// load expectations for model: classics.go
	flattenRun.AddExpectations("classics.go", []string{
		`type Classics struct {`,
		"	P0 *int64 `json:\"-\"`",
		"	P1 *strfmt.ISBN `json:\"-\"`",
		"	P2 Comics `json:\"-\"`",
		"	ClassicsItems []ClassicsTupleAdditionalItems `json:\"-\"`",
		`func (m *Classics) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateP2(formats); err != nil {`,
		`	if err := m.validateClassicsItems(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *Classics) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`func (m *Classics) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("1", "body", m.P1); err != nil {`,
		`	if err := validate.FormatOf("1", "body", "isbn", m.P1.String(), formats); err != nil {`,
		`func (m *Classics) validateP2(formats strfmt.Registry) error {`,
		`	if err := m.P2.Validate(formats); err != nil {`,
		`		if ve, ok := err.(*errors.Validation); ok {`,
		`			return ve.ValidateName("2"`,
		`func (m *Classics) validateClassicsItems(formats strfmt.Registry) error {`,
		`	for i := range m.ClassicsItems {`,
		`		if err := m.ClassicsItems[i].Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName(strconv.Itoa(i + 3)`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("classics.go", []string{
		`type Classics struct {`,
		"	P0 *int64 `json:\"-\"`",
		"	P1 *strfmt.ISBN `json:\"-\"`",
		"	P2 *ClassicsTuple0 `json:\"-\"`",
		// TODO: items should not be pointer
		"	ClassicsItems []*ClassicsClassicsItemsTuple0 `json:\"-\"`",
		`func (m *Classics) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateP2(formats); err != nil {`,
		`	if err := m.validateClassicsItems(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *Classics) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`func (m *Classics) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("1", "body", m.P1); err != nil {`,
		`	if err := validate.FormatOf("1", "body", "isbn", m.P1.String(), formats); err != nil {`,
		`func (m *Classics) validateP2(formats strfmt.Registry) error {`,
		`	if err := validate.Required("2", "body", m.P2); err != nil {`,
		`	if m.P2 != nil {`,
		`		if err := m.P2.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("2"`,
		`func (m *Classics) validateClassicsItems(formats strfmt.Registry) error {`,
		`	for i := range m.ClassicsItems {`,
		`		if m.ClassicsItems[i] != nil {`,
		`			if err := m.ClassicsItems[i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName(strconv.Itoa(i + 3)`,
		`type ClassicsClassicsItemsTuple0 struct {`,
		"	P0 *ClassicsClassicsItemsTuple0P0 `json:\"-\"`",
		"	P1 []strfmt.Date `json:\"-\"`",
		"	P2 *ClassicsClassicsItemsTuple0P2 `json:\"-\"`",
		"	P3 *ClassicsClassicsItemsTuple0P3Tuple0 `json:\"-\"`",
		`func (m *ClassicsClassicsItemsTuple0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateP2(formats); err != nil {`,
		`	if err := m.validateP3(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ClassicsClassicsItemsTuple0) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P0", "body", m.P0); err != nil {`,
		`	if m.P0 != nil {`,
		`		if err := m.P0.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("P0"`,
		`func (m *ClassicsClassicsItemsTuple0) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P1", "body", m.P1); err != nil {`,
		`	for i := 0; i < len(m.P1); i++ {`,
		`		if err := validate.FormatOf("P1"+"."+strconv.Itoa(i), "body", "date", m.P1[i].String(), formats); err != nil {`,
		`func (m *ClassicsClassicsItemsTuple0) validateP2(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P2", "body", m.P2); err != nil {`,
		`	if m.P2 != nil {`,
		`		if err := m.P2.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("P2"`,
		`func (m *ClassicsClassicsItemsTuple0) validateP3(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P3", "body", m.P3); err != nil {`,
		`	if m.P3 != nil {`,
		`		if err := m.P3.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("P3"`,
		`type ClassicsClassicsItemsTuple0P0 struct {`,
		"	Period *string `json:\"period,omitempty\"`",
		"	Title *string `json:\"title,omitempty\"`",
		`func (m *ClassicsClassicsItemsTuple0P0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateTitle(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`var classicsClassicsItemsTuple0P0TypeTitlePropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"Les Misérables\",\"Bleak House\",\"Sherlock Holmes\",\"Siddhartha\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		classicsClassicsItemsTuple0P0TypeTitlePropEnum = append(classicsClassicsItemsTuple0P0TypeTitlePropEnum, v`,
		`	ClassicsClassicsItemsTuple0P0TitleLesMisérables string = "Les Misérables"`,
		`	ClassicsClassicsItemsTuple0P0TitleBleakHouse string = "Bleak House"`,
		`	ClassicsClassicsItemsTuple0P0TitleSherlockHolmes string = "Sherlock Holmes"`,
		`	ClassicsClassicsItemsTuple0P0TitleSiddhartha string = "Siddhartha"`,
		`func (m *ClassicsClassicsItemsTuple0P0) validateTitleEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, classicsClassicsItemsTuple0P0TypeTitlePropEnum, true); err != nil {`,
		`func (m *ClassicsClassicsItemsTuple0P0) validateTitle(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Title) {`,
		`	if err := m.validateTitleEnum("P0"+"."+"title", "body", *m.Title); err != nil {`,
		`type ClassicsClassicsItemsTuple0P2 struct {`,
		"	Origin *string `json:\"origin,omitempty\"`",
		"	ClassicsClassicsItemsTuple0P2 map[string]string `json:\"-\"`",
		`var classicsClassicsItemsTuple0P2ValueEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"bookshop\",\"amazon\",\"library\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		classicsClassicsItemsTuple0P2ValueEnum = append(classicsClassicsItemsTuple0P2ValueEnum, v`,
		`func (m *ClassicsClassicsItemsTuple0P2) validateClassicsClassicsItemsTuple0P2ValueEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, classicsClassicsItemsTuple0P2ValueEnum, true); err != nil {`,
		`func (m *ClassicsClassicsItemsTuple0P2) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateOrigin(formats); err != nil {`,
		`	for k := range m.ClassicsClassicsItemsTuple0P2 {`,
		// removed undue IsZero() call
		`		if err := m.validateClassicsClassicsItemsTuple0P2ValueEnum("P2"+"."+k, "body", m.ClassicsClassicsItemsTuple0P2[k]); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`var classicsClassicsItemsTuple0P2TypeOriginPropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"print\",\"e-book\",\"collection\",\"museum\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		classicsClassicsItemsTuple0P2TypeOriginPropEnum = append(classicsClassicsItemsTuple0P2TypeOriginPropEnum, v`,
		`	ClassicsClassicsItemsTuple0P2OriginPrint string = "print"`,
		`	ClassicsClassicsItemsTuple0P2OriginEDashBook string = "e-book"`,
		`	ClassicsClassicsItemsTuple0P2OriginCollection string = "collection"`,
		`	ClassicsClassicsItemsTuple0P2OriginMuseum string = "museum"`,
		`func (m *ClassicsClassicsItemsTuple0P2) validateOriginEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, classicsClassicsItemsTuple0P2TypeOriginPropEnum, true); err != nil {`,
		`func (m *ClassicsClassicsItemsTuple0P2) validateOrigin(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Origin) {`,
		`	if err := m.validateOriginEnum("P2"+"."+"origin", "body", *m.Origin); err != nil {`,
		`type ClassicsClassicsItemsTuple0P3Tuple0 struct {`,
		"	P0 *string `json:\"-\"`",
		"	P1 *ClassicsClassicsItemsTuple0P3Tuple0P1 `json:\"-\"`",
		"	P2 *ClassicsClassicsItemsTuple0P3Tuple0P2 `json:\"-\"`",
		"	P3 *ClassicsClassicsItemsTuple0P3Tuple0P3 `json:\"-\"`",
		"	P4 []strfmt.ISBN `json:\"-\"`",
		"	P5 *int64 `json:\"-\"`",
		`func (m *ClassicsClassicsItemsTuple0P3Tuple0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateP2(formats); err != nil {`,
		`	if err := m.validateP3(formats); err != nil {`,
		`	if err := m.validateP4(formats); err != nil {`,
		`	if err := m.validateP5(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ClassicsClassicsItemsTuple0P3Tuple0) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P0", "body", m.P0); err != nil {`,
		`func (m *ClassicsClassicsItemsTuple0P3Tuple0) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P1", "body", m.P1); err != nil {`,
		`	if m.P1 != nil {`,
		`		if err := m.P1.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("P1"`,
		`func (m *ClassicsClassicsItemsTuple0P3Tuple0) validateP2(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P2", "body", m.P2); err != nil {`,
		`	if m.P2 != nil {`,
		`		if err := m.P2.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("P2"`,
		`func (m *ClassicsClassicsItemsTuple0P3Tuple0) validateP3(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P3", "body", m.P3); err != nil {`,
		`	if m.P3 != nil {`,
		`		if err := m.P3.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("P3"`,
		`func (m *ClassicsClassicsItemsTuple0P3Tuple0) validateP4(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P4", "body", m.P4); err != nil {`,
		`	for i := 0; i < len(m.P4); i++ {`,
		`		if err := validate.FormatOf("P4"+"."+strconv.Itoa(i), "body", "isbn", m.P4[i].String(), formats); err != nil {`,
		`func (m *ClassicsClassicsItemsTuple0P3Tuple0) validateP5(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P5", "body", m.P5); err != nil {`,
		`type ClassicsClassicsItemsTuple0P3Tuple0P1 struct {`,
		"	Narrative *string `json:\"narrative\"`",
		`func (m *ClassicsClassicsItemsTuple0P3Tuple0P1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateNarrative(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ClassicsClassicsItemsTuple0P3Tuple0P1) validateNarrative(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P1"+"."+"narrative", "body", m.Narrative); err != nil {`,
		`type ClassicsClassicsItemsTuple0P3Tuple0P2 struct {`,
		"	MarketingBS *string `json:\"marketingBS,omitempty\"`",
		`func (m *ClassicsClassicsItemsTuple0P3Tuple0P2) Validate(formats strfmt.Registry) error {`,
		`		return errors.CompositeValidationError(res...`,
		`type ClassicsClassicsItemsTuple0P3Tuple0P3 struct {`,
		"	Author *string `json:\"author,omitempty\"`",
		"	Character *string `json:\"character,omitempty\"`",
		`func (m *ClassicsClassicsItemsTuple0P3Tuple0P3) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAuthor(formats); err != nil {`,
		`	if err := m.validateCharacter(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ClassicsClassicsItemsTuple0P3Tuple0P3) validateAuthor(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Author) {`,
		`	if err := validate.MinLength("P3"+"."+"author", "body", *m.Author, 1); err != nil {`,
		`func (m *ClassicsClassicsItemsTuple0P3Tuple0P3) validateCharacter(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Character) {`,
		"	if err := validate.Pattern(\"P3\"+\".\"+\"character\", \"body\", *m.Character, `^[A-Z]+$`); err != nil {",
		`type ClassicsTuple0 struct {`,
		"	P0 *string `json:\"-\"`",
		"	P1 *ClassicsTuple0P1 `json:\"-\"`",
		"	P2 *ClassicsTuple0P2 `json:\"-\"`",
		"	P3 *ClassicsTuple0P3 `json:\"-\"`",
		"	P4 []strfmt.ISBN `json:\"-\"`",
		"	P5 *int64 `json:\"-\"`",
		`func (m *ClassicsTuple0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateP2(formats); err != nil {`,
		`	if err := m.validateP3(formats); err != nil {`,
		`	if err := m.validateP4(formats); err != nil {`,
		`	if err := m.validateP5(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ClassicsTuple0) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P0", "body", m.P0); err != nil {`,
		`func (m *ClassicsTuple0) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P1", "body", m.P1); err != nil {`,
		`	if m.P1 != nil {`,
		`		if err := m.P1.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("P1"`,
		`func (m *ClassicsTuple0) validateP2(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P2", "body", m.P2); err != nil {`,
		`	if m.P2 != nil {`,
		`		if err := m.P2.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("P2"`,
		`func (m *ClassicsTuple0) validateP3(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P3", "body", m.P3); err != nil {`,
		`	if m.P3 != nil {`,
		`		if err := m.P3.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("P3"`,
		`func (m *ClassicsTuple0) validateP4(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P4", "body", m.P4); err != nil {`,
		`	for i := 0; i < len(m.P4); i++ {`,
		`		if err := validate.FormatOf("P4"+"."+strconv.Itoa(i), "body", "isbn", m.P4[i].String(), formats); err != nil {`,
		`func (m *ClassicsTuple0) validateP5(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P5", "body", m.P5); err != nil {`,
		`type ClassicsTuple0P1 struct {`,
		"	Narrative *string `json:\"narrative\"`",
		`func (m *ClassicsTuple0P1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateNarrative(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ClassicsTuple0P1) validateNarrative(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P1"+"."+"narrative", "body", m.Narrative); err != nil {`,
		`type ClassicsTuple0P2 struct {`,
		"	MarketingBS *string `json:\"marketingBS,omitempty\"`",
		`func (m *ClassicsTuple0P2) Validate(formats strfmt.Registry) error {`,
		`		return errors.CompositeValidationError(res...`,
		`type ClassicsTuple0P3 struct {`,
		"	Author *string `json:\"author,omitempty\"`",
		"	Character *string `json:\"character,omitempty\"`",
		`func (m *ClassicsTuple0P3) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAuthor(formats); err != nil {`,
		`	if err := m.validateCharacter(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ClassicsTuple0P3) validateAuthor(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Author) {`,
		`	if err := validate.MinLength("P3"+"."+"author", "body", *m.Author, 1); err != nil {`,
		`func (m *ClassicsTuple0P3) validateCharacter(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Character) {`,
		"	if err := validate.Pattern(\"P3\"+\".\"+\"character\", \"body\", *m.Character, `^[A-Z]+$`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: comics_items2.go
	flattenRun.AddExpectations("comics_items2.go", []string{
		`type ComicsItems2 struct {`,
		"	MarketingBS string `json:\"marketingBS,omitempty\"`",
		// empty validation
		"func (m *ComicsItems2) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: classics_items_additional_items_items2.go
	flattenRun.AddExpectations("classics_items_additional_items_items2.go", []string{
		`type ClassicsItemsAdditionalItemsItems2 struct {`,
		"	Origin string `json:\"origin,omitempty\"`",
		"	ClassicsItemsAdditionalItemsItems2 map[string]string `json:\"-\"`",
		`var classicsItemsAdditionalItemsItems2ValueEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"bookshop\",\"amazon\",\"library\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		classicsItemsAdditionalItemsItems2ValueEnum = append(classicsItemsAdditionalItemsItems2ValueEnum, v`,
		`func (m *ClassicsItemsAdditionalItemsItems2) validateClassicsItemsAdditionalItemsItems2ValueEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, classicsItemsAdditionalItemsItems2ValueEnum, true); err != nil {`,
		`func (m *ClassicsItemsAdditionalItemsItems2) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateOrigin(formats); err != nil {`,
		`	for k := range m.ClassicsItemsAdditionalItemsItems2 {`,
		// removed undue IsZero()
		`		if err := m.validateClassicsItemsAdditionalItemsItems2ValueEnum(k, "body", m.ClassicsItemsAdditionalItemsItems2[k]); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`var classicsItemsAdditionalItemsItems2TypeOriginPropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"print\",\"e-book\",\"collection\",\"museum\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		classicsItemsAdditionalItemsItems2TypeOriginPropEnum = append(classicsItemsAdditionalItemsItems2TypeOriginPropEnum, v`,
		`	ClassicsItemsAdditionalItemsItems2OriginPrint string = "print"`,
		`	ClassicsItemsAdditionalItemsItems2OriginEDashBook string = "e-book"`,
		`	ClassicsItemsAdditionalItemsItems2OriginCollection string = "collection"`,
		`	ClassicsItemsAdditionalItemsItems2OriginMuseum string = "museum"`,
		`func (m *ClassicsItemsAdditionalItemsItems2) validateOriginEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, classicsItemsAdditionalItemsItems2TypeOriginPropEnum, true); err != nil {`,
		`func (m *ClassicsItemsAdditionalItemsItems2) validateOrigin(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Origin) {`,
		`	if err := m.validateOriginEnum("origin", "body", m.Origin); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: comics.go
	flattenRun.AddExpectations("comics.go", []string{
		`type Comics struct {`,
		"	P0 *string `json:\"-\"`",
		"	P1 *ComicsItems1 `json:\"-\"`",
		"	P2 *ComicsItems2 `json:\"-\"`",
		"	P3 *ComicsItems3 `json:\"-\"`",
		"	P4 []strfmt.ISBN `json:\"-\"`",
		"	P5 *int64 `json:\"-\"`",
		`func (m *Comics) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateP2(formats); err != nil {`,
		`	if err := m.validateP3(formats); err != nil {`,
		`	if err := m.validateP4(formats); err != nil {`,
		`	if err := m.validateP5(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *Comics) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`func (m *Comics) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("1", "body", m.P1); err != nil {`,
		`	if m.P1 != nil {`,
		`		if err := m.P1.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("1"`,
		`func (m *Comics) validateP2(formats strfmt.Registry) error {`,
		`	if err := validate.Required("2", "body", m.P2); err != nil {`,
		`	if m.P2 != nil {`,
		`		if err := m.P2.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("2"`,
		`func (m *Comics) validateP3(formats strfmt.Registry) error {`,
		`	if err := validate.Required("3", "body", m.P3); err != nil {`,
		`	if m.P3 != nil {`,
		`		if err := m.P3.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("3"`,
		`func (m *Comics) validateP4(formats strfmt.Registry) error {`,
		`	if err := validate.Required("4", "body", m.P4); err != nil {`,
		`	for i := 0; i < len(m.P4); i++ {`,
		`		if err := validate.FormatOf("4"+"."+strconv.Itoa(i), "body", "isbn", m.P4[i].String(), formats); err != nil {`,
		`func (m *Comics) validateP5(formats strfmt.Registry) error {`,
		`	if err := validate.Required("5", "body", m.P5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("comics.go", []string{
		`type Comics struct {`,
		"	P0 *string `json:\"-\"`",
		"   P1 *ComicsItems1 `json:\"-\"`",
		"   P2 *ComicsItems2 `json:\"-\"`",
		"   P3 *ComicsItems3 `json:\"-\"`",
		"	P4 []strfmt.ISBN `json:\"-\"`",
		"	P5 *int64 `json:\"-\"`",
		`func (m *Comics) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateP2(formats); err != nil {`,
		`	if err := m.validateP3(formats); err != nil {`,
		`	if err := m.validateP4(formats); err != nil {`,
		`	if err := m.validateP5(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *Comics) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`func (m *Comics) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("1", "body", m.P1); err != nil {`,
		`func (m *Comics) validateP2(formats strfmt.Registry) error {`,
		`	if err := validate.Required("2", "body", m.P2); err != nil {`,
		`func (m *Comics) validateP3(formats strfmt.Registry) error {`,
		`	if err := validate.Required("3", "body", m.P3); err != nil {`,
		`func (m *Comics) validateP4(formats strfmt.Registry) error {`,
		`	if err := validate.Required("4", "body", m.P4); err != nil {`,
		`	for i := 0; i < len(m.P4); i++ {`,
		`		if err := validate.FormatOf("4"+"."+strconv.Itoa(i), "body", "isbn", m.P4[i].String(), formats); err != nil {`,
		`func (m *Comics) validateP5(formats strfmt.Registry) error {`,
		`	if err := validate.Required("5", "body", m.P5); err != nil {`,
		`type ComicsItems1 struct {`,
		`type ComicsItems2 struct {`,
		`type ComicsItems3 struct {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: classics_items_additional_items_items0.go
	flattenRun.AddExpectations("classics_items_additional_items_items0.go", []string{
		`type ClassicsItemsAdditionalItemsItems0 struct {`,
		"	Period string `json:\"period,omitempty\"`",
		"	Title string `json:\"title,omitempty\"`",
		`func (m *ClassicsItemsAdditionalItemsItems0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateTitle(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`var classicsItemsAdditionalItemsItems0TypeTitlePropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"Les Misérables\",\"Bleak House\",\"Sherlock Holmes\",\"Siddhartha\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		classicsItemsAdditionalItemsItems0TypeTitlePropEnum = append(classicsItemsAdditionalItemsItems0TypeTitlePropEnum, v`,
		`	ClassicsItemsAdditionalItemsItems0TitleLesMisérables string = "Les Misérables"`,
		`	ClassicsItemsAdditionalItemsItems0TitleBleakHouse string = "Bleak House"`,
		`	ClassicsItemsAdditionalItemsItems0TitleSherlockHolmes string = "Sherlock Holmes"`,
		`	ClassicsItemsAdditionalItemsItems0TitleSiddhartha string = "Siddhartha"`,
		`func (m *ClassicsItemsAdditionalItemsItems0) validateTitleEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, classicsItemsAdditionalItemsItems0TypeTitlePropEnum, true); err != nil {`,
		`func (m *ClassicsItemsAdditionalItemsItems0) validateTitle(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Title) {`,
		`	if err := m.validateTitleEnum("title", "body", m.Title); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: comics_items1.go
	flattenRun.AddExpectations("comics_items1.go", []string{
		`type ComicsItems1 struct {`,
		"	Narrative *string `json:\"narrative\"`",
		`func (m *ComicsItems1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateNarrative(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ComicsItems1) validateNarrative(formats strfmt.Registry) error {`,
		`	if err := validate.Required("narrative", "body", m.Narrative); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: comics_items3.go
	flattenRun.AddExpectations("comics_items3.go", []string{
		`type ComicsItems3 struct {`,
		"	Author string `json:\"author,omitempty\"`",
		"	Character string `json:\"character,omitempty\"`",
		`func (m *ComicsItems3) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAuthor(formats); err != nil {`,
		`	if err := m.validateCharacter(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ComicsItems3) validateAuthor(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Author) {`,
		`	if err := validate.MinLength("author", "body", m.Author, 1); err != nil {`,
		`func (m *ComicsItems3) validateCharacter(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Character) {`,
		"	if err := validate.Pattern(\"character\", \"body\", m.Character, `^[A-Z]+$`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: classics_tuple_additional_items.go
	flattenRun.AddExpectations("classics_tuple_additional_items.go", []string{
		`type ClassicsTupleAdditionalItems struct {`,
		"	P0 *ClassicsItemsAdditionalItemsItems0 `json:\"-\"`",
		"	P1 []strfmt.Date `json:\"-\"`",
		"	P2 *ClassicsItemsAdditionalItemsItems2 `json:\"-\"`",
		"	P3 Comics `json:\"-\"`",
		`func (m *ClassicsTupleAdditionalItems) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateP2(formats); err != nil {`,
		`	if err := m.validateP3(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ClassicsTupleAdditionalItems) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	if m.P0 != nil {`,
		`		if err := m.P0.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("0"`,
		`func (m *ClassicsTupleAdditionalItems) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("1", "body", m.P1); err != nil {`,
		`	for i := 0; i < len(m.P1); i++ {`,
		`		if err := validate.FormatOf("1"+"."+strconv.Itoa(i), "body", "date", m.P1[i].String(), formats); err != nil {`,
		`func (m *ClassicsTupleAdditionalItems) validateP2(formats strfmt.Registry) error {`,
		`	if err := validate.Required("2", "body", m.P2); err != nil {`,
		`	if m.P2 != nil {`,
		`		if err := m.P2.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("2"`,
		`func (m *ClassicsTupleAdditionalItems) validateP3(formats strfmt.Registry) error {`,
		`	if err := m.P3.Validate(formats); err != nil {`,
		`		if ve, ok := err.(*errors.Validation); ok {`,
		`			return ve.ValidateName("3"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

}

func initFixture1198() {
	// testing ../fixtures/bugs/1487/fixture-1198.yaml with expand (--skip-flatten)

	f := newModelFixture("../fixtures/bugs/1198/fixture-1198.yaml", "string-body-api")
	flattenRun := f.AddRun(false)

	// load expectations for model: pet.go
	flattenRun.AddExpectations("pet.go", []string{
		`type Pet struct {`,
		"	Date interface{} `json:\"date\"`",
		`func (m *Pet) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateDate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *Pet) validateDate(formats strfmt.Registry) error {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

}

func initFixture1042() {
	// testing ../fixtures/bugs/1487/fixture-1042.yaml with expand (--skip-flatten)

	/* when the specification incorrectly defines the allOf,
	generated unmarshalling is wrong.
	This fixture asserts that with correct spec, the generated models are correct.

	*/

	f := newModelFixture("../fixtures/bugs/1042/fixture-1042.yaml", "allOf marshalling")
	flattenRun := f.AddRun(false)

	// load expectations for model: b.go
	flattenRun.AddExpectations("b.go", []string{
		`type B struct {`,
		`	A`,
		`	BAllOf1`,
		`func (m *B) UnmarshalJSON(raw []byte) error {`,
		`	var aO0 A`,
		`	if err := swag.ReadJSON(raw, &aO0); err != nil {`,
		`	m.A = aO0`,
		`	var aO1 BAllOf1`,
		`	if err := swag.ReadJSON(raw, &aO1); err != nil {`,
		`	m.BAllOf1 = aO1`,
		`func (m B) MarshalJSON() ([]byte, error) {`,
		// slight optimization of allocations
		`	_parts := make([][]byte, 0, 2)`,
		`	aO0, err := swag.WriteJSON(m.A`,
		`	if err != nil {`,
		`		return nil, err`,
		`	_parts = append(_parts, aO0`,
		`	aO1, err := swag.WriteJSON(m.BAllOf1`,
		`	if err != nil {`,
		`		return nil, err`,
		`	_parts = append(_parts, aO1`,
		`	return swag.ConcatJSON(_parts...), nil`,
		`func (m *B) Validate(formats strfmt.Registry) error {`,
		`	if err := m.A.Validate(formats); err != nil {`,
		`	if err := m.BAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: b_all_of1.go
	flattenRun.AddExpectations("b_all_of1.go", []string{
		`type BAllOf1 struct {`,
		"	F3 *string `json:\"f3\"`",
		"	F4 []string `json:\"f4\"`",
		`func (m *BAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateF3(formats); err != nil {`,
		`	if err := m.validateF4(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *BAllOf1) validateF3(formats strfmt.Registry) error {`,
		`	if err := validate.Required("f3", "body", m.F3); err != nil {`,
		`func (m *BAllOf1) validateF4(formats strfmt.Registry) error {`,
		`	if err := validate.Required("f4", "body", m.F4); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: a.go
	flattenRun.AddExpectations("a.go", []string{
		`type A struct {`,
		"	F1 string `json:\"f1,omitempty\"`",
		"	F2 string `json:\"f2,omitempty\"`",
		// empty validation
		"func (m *A) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

}

func initFixture1042V2() {
	// testing ../fixtures/bugs/1487/fixture-1042-2.yaml with expand (--skip-flatten)

	/* when the specification incorrectly defines the allOf,
	generated unmarshalling is wrong.
	This fixture asserts that with correct spec, the generated models are correct.

	*/

	f := newModelFixture("../fixtures/bugs/1042/fixture-1042-2.yaml", "allOf marshalling")
	flattenRun := f.AddRun(false)

	// load expectations for model: error_model.go
	flattenRun.AddExpectations("error_model.go", []string{
		`type ErrorModel struct {`,
		"	Code *int64 `json:\"code\"`",
		"	Message *string `json:\"message\"`",
		`func (m *ErrorModel) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateCode(formats); err != nil {`,
		`	if err := m.validateMessage(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ErrorModel) validateCode(formats strfmt.Registry) error {`,
		`	if err := validate.Required("code", "body", m.Code); err != nil {`,
		`	if err := validate.MinimumInt("code", "body", *m.Code, 100, false); err != nil {`,
		`	if err := validate.MaximumInt("code", "body", *m.Code, 600, false); err != nil {`,
		`func (m *ErrorModel) validateMessage(formats strfmt.Registry) error {`,
		`	if err := validate.Required("message", "body", m.Message); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: extended_error_model.go
	flattenRun.AddExpectations("extended_error_model.go", []string{
		`type ExtendedErrorModel struct {`,
		`	ErrorModel`,
		`	ExtendedErrorModelAllOf1`,
		`func (m *ExtendedErrorModel) UnmarshalJSON(raw []byte) error {`,
		`	var aO0 ErrorModel`,
		`	if err := swag.ReadJSON(raw, &aO0); err != nil {`,
		`	m.ErrorModel = aO0`,
		`	var aO1 ExtendedErrorModelAllOf1`,
		`	if err := swag.ReadJSON(raw, &aO1); err != nil {`,
		`	m.ExtendedErrorModelAllOf1 = aO1`,
		`func (m ExtendedErrorModel) MarshalJSON() ([]byte, error) {`,
		// slight optimization of allocations
		`	_parts := make([][]byte, 0, 2)`,
		`	aO0, err := swag.WriteJSON(m.ErrorModel`,
		`	if err != nil {`,
		`		return nil, err`,
		`	_parts = append(_parts, aO0`,
		`	aO1, err := swag.WriteJSON(m.ExtendedErrorModelAllOf1`,
		`	if err != nil {`,
		`		return nil, err`,
		`	_parts = append(_parts, aO1`,
		`	return swag.ConcatJSON(_parts...), nil`,
		`func (m *ExtendedErrorModel) Validate(formats strfmt.Registry) error {`,
		`	if err := m.ErrorModel.Validate(formats); err != nil {`,
		`	if err := m.ExtendedErrorModelAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: extended_error_model_all_of1.go
	flattenRun.AddExpectations("extended_error_model_all_of1.go", []string{
		`type ExtendedErrorModelAllOf1 struct {`,
		"	RootCause *string `json:\"rootCause\"`",
		`func (m *ExtendedErrorModelAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateRootCause(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ExtendedErrorModelAllOf1) validateRootCause(formats strfmt.Registry) error {`,
		`	if err := validate.Required("rootCause", "body", m.RootCause); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

}

func initFixture979() {
	// testing ../fixtures/bugs/1487/fixture-979.yaml with expand (--skip-flatten)

	/* checking that properties is enough to figure out an object schema
	 */

	f := newModelFixture("../fixtures/bugs/979/fixture-979.yaml", "allOf without the explicit type object")
	flattenRun := f.AddRun(false)

	// load expectations for model: cluster.go
	flattenRun.AddExpectations("cluster.go", []string{
		`type Cluster struct {`,
		`	NewCluster`,
		`	ClusterAllOf1`,
		`func (m *Cluster) UnmarshalJSON(raw []byte) error {`,
		`	var aO0 NewCluster`,
		`	if err := swag.ReadJSON(raw, &aO0); err != nil {`,
		`	m.NewCluster = aO0`,
		`	var aO1 ClusterAllOf1`,
		`	if err := swag.ReadJSON(raw, &aO1); err != nil {`,
		`	m.ClusterAllOf1 = aO1`,
		`func (m Cluster) MarshalJSON() ([]byte, error) {`,
		// slight optimization of allocations
		`	_parts := make([][]byte, 0, 2)`,
		`	aO0, err := swag.WriteJSON(m.NewCluster`,
		`	if err != nil {`,
		`		return nil, err`,
		`	_parts = append(_parts, aO0`,
		`	aO1, err := swag.WriteJSON(m.ClusterAllOf1`,
		`	if err != nil {`,
		`		return nil, err`,
		`	_parts = append(_parts, aO1`,
		`	return swag.ConcatJSON(_parts...), nil`,
		`func (m *Cluster) Validate(formats strfmt.Registry) error {`,
		`	if err := m.NewCluster.Validate(formats); err != nil {`,
		`	if err := m.ClusterAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: new_cluster.go
	flattenRun.AddExpectations("new_cluster.go", []string{
		`type NewCluster struct {`,
		"	DummyProp1 int64 `json:\"dummyProp1,omitempty\"`",
		// empty validation
		"func (m *NewCluster) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: cluster_all_of1.go
	flattenRun.AddExpectations("cluster_all_of1.go", []string{
		`type ClusterAllOf1 struct {`,
		"	Result string `json:\"result,omitempty\"`",
		"	Status string `json:\"status,omitempty\"`",
		// empty validation
		"func (m *ClusterAllOf1) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

}

func initFixture842() {
	// testing ../fixtures/bugs/1487/fixture-842.yaml with expand (--skip-flatten)

	/* codegen fails to produce code that builds
	 */

	f := newModelFixture("../fixtures/bugs/842/fixture-842.yaml", "polymorphic type containing an array of the base type")
	flattenRun := f.AddRun(false)

	// load expectations for model: value_array_all_of1.go
	flattenRun.AddExpectations("value_array_all_of1.go", []string{
		`type ValueArrayAllOf1 struct {`,
		`	valuesField []Value`,
		`func (m *ValueArrayAllOf1) Values() []Value {`,
		`	return m.valuesField`,
		`func (m *ValueArrayAllOf1) SetValues(val []Value) {`,
		`	m.valuesField = val`,
		`func (m *ValueArrayAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateValues(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ValueArrayAllOf1) validateValues(formats strfmt.Registry) error {`,
		`	if err := validate.Required("Values", "body", m.Values()); err != nil {`,
		`	for i := 0; i < len(m.Values()); i++ {`,
		`		if err := m.valuesField[i].Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("Values" + "." + strconv.Itoa(i)`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: value_array.go
	flattenRun.AddExpectations("value_array.go", []string{
		`type ValueArray struct {`,
		`	ValueArrayAllOf1`,
		`func (m *ValueArray) ValueType() string {`,
		`	return "ValueArray"`,
		`func (m *ValueArray) SetValueType(val string) {`,
		`func (m *ValueArray) Validate(formats strfmt.Registry) error {`,
		`	if err := m.ValueArrayAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: value.go
	flattenRun.AddExpectations("value.go", []string{
		`type Value interface {`,
		`	runtime.Validatable`,
		`	ValueType() string`,
		`	SetValueType(string`,
		`type value struct {`,
		`	valueTypeField string`,
		`func (m *value) ValueType() string {`,
		`	return "Value"`,
		`func (m *value) SetValueType(val string) {`,
		`func UnmarshalValueSlice(reader io.Reader, consumer runtime.Consumer) ([]Value, error) {`,
		`	var elements []json.RawMessage`,
		`	if err := consumer.Consume(reader, &elements); err != nil {`,
		`		return nil, err`,
		`	var result []Value`,
		`	for _, element := range elements {`,
		`		obj, err := unmarshalValue(element, consumer`,
		`		if err != nil {`,
		`			return nil, err`,
		`		result = append(result, obj`,
		`	return result, nil`,
		`func UnmarshalValue(reader io.Reader, consumer runtime.Consumer) (Value, error) {`,
		`	data, err := ioutil.ReadAll(reader`,
		`	if err != nil {`,
		`		return nil, err`,
		`	return unmarshalValue(data, consumer`,
		`func unmarshalValue(data []byte, consumer runtime.Consumer) (Value, error) {`,
		`	buf := bytes.NewBuffer(data`,
		`	buf2 := bytes.NewBuffer(data`,
		`	var getType struct {`,
		"		ValueType string `json:\"ValueType\"`",
		`	if err := consumer.Consume(buf, &getType); err != nil {`,
		`		return nil, err`,
		`	if err := validate.RequiredString("ValueType", "body", getType.ValueType); err != nil {`,
		`		return nil, err`,
		`	switch getType.ValueType {`,
		`	case "Value":`,
		`		var result value`,
		`		if err := consumer.Consume(buf2, &result); err != nil {`,
		`			return nil, err`,
		`		return &result, nil`,
		`	case "ValueArray":`,
		`		var result ValueArray`,
		`		if err := consumer.Consume(buf2, &result); err != nil {`,
		`			return nil, err`,
		`		return &result, nil`,
		`	return nil, errors.New(422, "invalid ValueType value: %q", getType.ValueType`,
		// empty validation
		"func (m *value) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

}

func initFixture607() {
	// testing ../fixtures/bugs/1487/fixture-607.yaml with expand (--skip-flatten)

	/* broken code produced on polymorphic type
	 */

	f := newModelFixture("../fixtures/bugs/607/fixture-607.yaml", "broken code when using array of polymorphic type")
	flattenRun := f.AddRun(false)

	// load expectations for model: range_filter_all_of1.go
	flattenRun.AddExpectations("range_filter_all_of1.go", []string{
		`type RangeFilterAllOf1 struct {`,
		"	Config *RangeFilterAllOf1Config `json:\"config\"`",
		`func (m *RangeFilterAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateConfig(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *RangeFilterAllOf1) validateConfig(formats strfmt.Registry) error {`,
		`	if err := validate.Required("config", "body", m.Config); err != nil {`,
		`	if m.Config != nil {`,
		`		if err := m.Config.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("config"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: filter.go
	flattenRun.AddExpectations("filter.go", []string{
		`type Filter interface {`,
		`	runtime.Validatable`,
		`	Type() string`,
		`	SetType(string`,
		`type filter struct {`,
		`	typeField string`,
		`func (m *filter) Type() string {`,
		`	return "Filter"`,
		`func (m *filter) SetType(val string) {`,
		`func UnmarshalFilterSlice(reader io.Reader, consumer runtime.Consumer) ([]Filter, error) {`,
		`	var elements []json.RawMessage`,
		`	if err := consumer.Consume(reader, &elements); err != nil {`,
		`		return nil, err`,
		`	var result []Filter`,
		`	for _, element := range elements {`,
		`		obj, err := unmarshalFilter(element, consumer`,
		`		if err != nil {`,
		`			return nil, err`,
		`		result = append(result, obj`,
		`	return result, nil`,
		`func UnmarshalFilter(reader io.Reader, consumer runtime.Consumer) (Filter, error) {`,
		`	data, err := ioutil.ReadAll(reader`,
		`	if err != nil {`,
		`		return nil, err`,
		`	return unmarshalFilter(data, consumer`,
		`func unmarshalFilter(data []byte, consumer runtime.Consumer) (Filter, error) {`,
		`	buf := bytes.NewBuffer(data`,
		`	buf2 := bytes.NewBuffer(data`,
		`	var getType struct {`,
		"		Type string `json:\"type\"`",
		`	if err := consumer.Consume(buf, &getType); err != nil {`,
		`		return nil, err`,
		`	if err := validate.RequiredString("type", "body", getType.Type); err != nil {`,
		`		return nil, err`,
		`	switch getType.Type {`,
		`	case "AndFilter":`,
		`		var result AndFilter`,
		`		if err := consumer.Consume(buf2, &result); err != nil {`,
		`			return nil, err`,
		`		return &result, nil`,
		`	case "Filter":`,
		`		var result filter`,
		`		if err := consumer.Consume(buf2, &result); err != nil {`,
		`			return nil, err`,
		`		return &result, nil`,
		`	case "RangeFilter":`,
		`		var result RangeFilter`,
		`		if err := consumer.Consume(buf2, &result); err != nil {`,
		`			return nil, err`,
		`		return &result, nil`,
		`	return nil, errors.New(422, "invalid type value: %q", getType.Type`,
		// empty validation
		"func (m *filter) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: and_filter_all_of1.go
	flattenRun.AddExpectations("and_filter_all_of1.go", []string{
		`type AndFilterAllOf1 struct {`,
		`	configField []Filter`,
		`func (m *AndFilterAllOf1) Config() []Filter {`,
		`	return m.configField`,
		`func (m *AndFilterAllOf1) SetConfig(val []Filter) {`,
		`	m.configField = val`,
		`func (m *AndFilterAllOf1) UnmarshalJSON(raw []byte) error {`,
		`	var data struct {`,
		"		Config json.RawMessage `json:\"config\"`",
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&data); err != nil {`,
		`	propConfig, err := UnmarshalFilterSlice(bytes.NewBuffer(data.Config), runtime.JSONConsumer()`,
		`	if err != nil && err != io.EOF {`,
		`	var result AndFilterAllOf1`,
		`	result.configField = propConfig`,
		`	*m = result`,
		`func (m AndFilterAllOf1) MarshalJSON() ([]byte, error) {`,
		`	var b1, b2, b3 []byte`,
		`	var err error`,
		`	b1, err = json.Marshal(struct {`,
		`	}{})`,
		`	if err != nil {`,
		`		return nil, err`,
		`	b2, err = json.Marshal(struct {`,
		"		Config []Filter `json:\"config\"`",
		`	}{`,
		`		Config: m.configField,`,
		`	})`,
		`	if err != nil {`,
		`		return nil, err`,
		`	return swag.ConcatJSON(b1, b2, b3), nil`,
		`func (m *AndFilterAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateConfig(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AndFilterAllOf1) validateConfig(formats strfmt.Registry) error {`,
		`	if err := validate.Required("config", "body", m.Config()); err != nil {`,
		`	for i := 0; i < len(m.Config()); i++ {`,
		`		if err := m.configField[i].Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("config" + "." + strconv.Itoa(i)`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: and_filter.go
	flattenRun.AddExpectations("and_filter.go", []string{
		`type AndFilter struct {`,
		`	AndFilterAllOf1`,
		`func (m *AndFilter) Type() string {`,
		`	return "AndFilter"`,
		`func (m *AndFilter) SetType(val string) {`,
		`func (m *AndFilter) UnmarshalJSON(raw []byte) error {`,
		`	var data struct {`,
		`		AndFilterAllOf1`,
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&data); err != nil {`,
		`	var base struct {`,
		"		Type string `json:\"type\"`",
		`	buf = bytes.NewBuffer(raw`,
		`	dec = json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&base); err != nil {`,
		`	var result AndFilter`,
		`	if base.Type != result.Type() {`,
		`		return errors.New(422, "invalid type value: %q", base.Type`,
		`	result.AndFilterAllOf1 = data.AndFilterAllOf1`,
		`	*m = result`,
		`func (m AndFilter) MarshalJSON() ([]byte, error) {`,
		`	var b1, b2, b3 []byte`,
		`	var err error`,
		`	b1, err = json.Marshal(struct {`,
		`		AndFilterAllOf1`,
		`	}{`,
		`		AndFilterAllOf1: m.AndFilterAllOf1,`,
		`	})`,
		`	if err != nil {`,
		`		return nil, err`,
		`	b2, err = json.Marshal(struct {`,
		"		Type string `json:\"type\"`",
		`	}{`,
		`		Type: m.Type(),`,
		`	})`,
		`	if err != nil {`,
		`		return nil, err`,
		`	return swag.ConcatJSON(b1, b2, b3), nil`,
		`func (m *AndFilter) Validate(formats strfmt.Registry) error {`,
		`	if err := m.AndFilterAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: range_filter.go
	flattenRun.AddExpectations("range_filter.go", []string{
		`type RangeFilter struct {`,
		`	RangeFilterAllOf1`,
		`func (m *RangeFilter) Type() string {`,
		`	return "RangeFilter"`,
		`func (m *RangeFilter) SetType(val string) {`,
		`func (m *RangeFilter) UnmarshalJSON(raw []byte) error {`,
		`	var data struct {`,
		`		RangeFilterAllOf1`,
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&data); err != nil {`,
		`	var base struct {`,
		"		Type string `json:\"type\"`",
		`	buf = bytes.NewBuffer(raw`,
		`	dec = json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&base); err != nil {`,
		`	var result RangeFilter`,
		`	if base.Type != result.Type() {`,
		`		return errors.New(422, "invalid type value: %q", base.Type`,
		`	result.RangeFilterAllOf1 = data.RangeFilterAllOf1`,
		`	*m = result`,
		`func (m RangeFilter) MarshalJSON() ([]byte, error) {`,
		`	var b1, b2, b3 []byte`,
		`	var err error`,
		`	b1, err = json.Marshal(struct {`,
		`		RangeFilterAllOf1`,
		`	}{`,
		`		RangeFilterAllOf1: m.RangeFilterAllOf1,`,
		`	})`,
		`	if err != nil {`,
		`		return nil, err`,
		`	b2, err = json.Marshal(struct {`,
		"		Type string `json:\"type\"`",
		`	}{`,
		`		Type: m.Type(),`,
		`	})`,
		`	if err != nil {`,
		`		return nil, err`,
		`	return swag.ConcatJSON(b1, b2, b3), nil`,
		`func (m *RangeFilter) Validate(formats strfmt.Registry) error {`,
		`	if err := m.RangeFilterAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: range_filter_all_of1_config.go
	flattenRun.AddExpectations("range_filter_all_of1_config.go", []string{
		`type RangeFilterAllOf1Config struct {`,
		"	Gt float64 `json:\"gt,omitempty\"`",
		"	Gte float64 `json:\"gte,omitempty\"`",
		"	Lt float64 `json:\"lt,omitempty\"`",
		"	Lte float64 `json:\"lte,omitempty\"`",
		// empty validation
		"func (m *RangeFilterAllOf1Config) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

}

func initFixture1336() {
	// testing ../fixtures/bugs/1487/fixture-1336.yaml with expand (--skip-flatten)

	/* broken code produced on polymorphic type
	 */

	f := newModelFixture("../fixtures/bugs/1336/fixture-1336.yaml", "broken code when using array of polymorphic type")
	flattenRun := f.AddRun(false)

	// load expectations for model: node.go
	flattenRun.AddExpectations("node.go", []string{
		`type Node interface {`,
		`	runtime.Validatable`,
		`	NodeType() string`,
		`	SetNodeType(string`,
		`type node struct {`,
		`	nodeTypeField string`,
		`func (m *node) NodeType() string {`,
		`	return "Node"`,
		`func (m *node) SetNodeType(val string) {`,
		`func UnmarshalNodeSlice(reader io.Reader, consumer runtime.Consumer) ([]Node, error) {`,
		`	var elements []json.RawMessage`,
		`	if err := consumer.Consume(reader, &elements); err != nil {`,
		`		return nil, err`,
		`	var result []Node`,
		`	for _, element := range elements {`,
		`		obj, err := unmarshalNode(element, consumer`,
		`		if err != nil {`,
		`			return nil, err`,
		`		result = append(result, obj`,
		`	return result, nil`,
		`func UnmarshalNode(reader io.Reader, consumer runtime.Consumer) (Node, error) {`,
		`	data, err := ioutil.ReadAll(reader`,
		`	if err != nil {`,
		`		return nil, err`,
		`	return unmarshalNode(data, consumer`,
		`func unmarshalNode(data []byte, consumer runtime.Consumer) (Node, error) {`,
		`	buf := bytes.NewBuffer(data`,
		`	buf2 := bytes.NewBuffer(data`,
		`	var getType struct {`,
		"		NodeType string `json:\"NodeType\"`",
		`	if err := consumer.Consume(buf, &getType); err != nil {`,
		`		return nil, err`,
		`	if err := validate.RequiredString("NodeType", "body", getType.NodeType); err != nil {`,
		`		return nil, err`,
		`	switch getType.NodeType {`,
		`	case "CodeBlockNode":`,
		`		var result CodeBlockNode`,
		`		if err := consumer.Consume(buf2, &result); err != nil {`,
		`			return nil, err`,
		`		return &result, nil`,
		`	case "DocBlockNode":`,
		`		var result DocBlockNode`,
		`		if err := consumer.Consume(buf2, &result); err != nil {`,
		`			return nil, err`,
		`		return &result, nil`,
		`	case "Node":`,
		`		var result node`,
		`		if err := consumer.Consume(buf2, &result); err != nil {`,
		`			return nil, err`,
		`		return &result, nil`,
		`	return nil, errors.New(422, "invalid NodeType value: %q", getType.NodeType`,
		// empty validation
		"func (m *node) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: code_block_node_all_of1.go
	flattenRun.AddExpectations("code_block_node_all_of1.go", []string{
		`type CodeBlockNodeAllOf1 struct {`,
		"	Code string `json:\"Code,omitempty\"`",
		// empty validation
		"func (m *CodeBlockNodeAllOf1) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: graph.go
	flattenRun.AddExpectations("graph.go", []string{
		`type Graph struct {`,
		`	nodesField []Node`,
		`func (m *Graph) Nodes() []Node {`,
		`	return m.nodesField`,
		`func (m *Graph) SetNodes(val []Node) {`,
		`	m.nodesField = val`,
		`func (m *Graph) UnmarshalJSON(raw []byte) error {`,
		`	var data struct {`,
		"		Nodes json.RawMessage `json:\"Nodes\"`",
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&data); err != nil {`,
		`		nodes, err := UnmarshalNodeSlice(bytes.NewBuffer(data.Nodes), runtime.JSONConsumer()`,
		`	if err != nil && err != io.EOF {`,
		`	var result Graph`,
		`	result.nodesField = propNodes`,
		`	*m = result`,
		`func (m Graph) MarshalJSON() ([]byte, error) {`,
		`	var b1, b2, b3 []byte`,
		`	var err error`,
		`	b1, err = json.Marshal(struct {`,
		`	}{})`,
		`	if err != nil {`,
		`		return nil, err`,
		`	b2, err = json.Marshal(struct {`,
		"		Nodes []Node `json:\"Nodes\"`",
		`	}{`,
		`		Nodes: m.nodesField,`,
		`	})`,
		`	if err != nil {`,
		`		return nil, err`,
		`	return swag.ConcatJSON(b1, b2, b3), nil`,
		`func (m *Graph) Validate(formats strfmt.Registry) error {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: doc_block_node_all_of1.go
	flattenRun.AddExpectations("doc_block_node_all_of1.go", []string{
		`type DocBlockNodeAllOf1 struct {`,
		"	Doc string `json:\"Doc,omitempty\"`",
		// empty validation
		"func (m *DocBlockNodeAllOf1) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: doc_block_node.go
	flattenRun.AddExpectations("doc_block_node.go", []string{
		`type DocBlockNode struct {`,
		`	DocBlockNodeAllOf1`,
		`func (m *DocBlockNode) NodeType() string {`,
		`	return "DocBlockNode"`,
		`func (m *DocBlockNode) SetNodeType(val string) {`,
		`func (m *DocBlockNode) UnmarshalJSON(raw []byte) error {`,
		`	var data struct {`,
		`		DocBlockNodeAllOf1`,
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&data); err != nil {`,
		`	var base struct {`,
		"		NodeType string `json:\"NodeType\"`",
		`	buf = bytes.NewBuffer(raw`,
		`	dec = json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&base); err != nil {`,
		`	var result DocBlockNode`,
		`	if base.NodeType != result.NodeType() {`,
		`		return errors.New(422, "invalid NodeType value: %q", base.NodeType`,
		`	result.DocBlockNodeAllOf1 = data.DocBlockNodeAllOf1`,
		`	*m = result`,
		`func (m DocBlockNode) MarshalJSON() ([]byte, error) {`,
		`	var b1, b2, b3 []byte`,
		`	var err error`,
		`	b1, err = json.Marshal(struct {`,
		`		DocBlockNodeAllOf1`,
		`	}{`,
		`		DocBlockNodeAllOf1: m.DocBlockNodeAllOf1,`,
		`	})`,
		`	if err != nil {`,
		`		return nil, err`,
		`	b2, err = json.Marshal(struct {`,
		"		NodeType string `json:\"NodeType\"`",
		`	}{`,
		`		NodeType: m.NodeType(),`,
		`	})`,
		`	if err != nil {`,
		`		return nil, err`,
		`	return swag.ConcatJSON(b1, b2, b3), nil`,
		`func (m *DocBlockNode) Validate(formats strfmt.Registry) error {`,
		`	if err := m.DocBlockNodeAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: code_block_node.go
	flattenRun.AddExpectations("code_block_node.go", []string{
		`type CodeBlockNode struct {`,
		`	CodeBlockNodeAllOf1`,
		`func (m *CodeBlockNode) NodeType() string {`,
		`	return "CodeBlockNode"`,
		`func (m *CodeBlockNode) SetNodeType(val string) {`,
		`func (m *CodeBlockNode) UnmarshalJSON(raw []byte) error {`,
		`	var data struct {`,
		`		CodeBlockNodeAllOf1`,
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&data); err != nil {`,
		`	var base struct {`,
		"		NodeType string `json:\"NodeType\"`",
		`	buf = bytes.NewBuffer(raw`,
		`	dec = json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&base); err != nil {`,
		`	var result CodeBlockNode`,
		`	if base.NodeType != result.NodeType() {`,
		`		return errors.New(422, "invalid NodeType value: %q", base.NodeType`,
		`	result.CodeBlockNodeAllOf1 = data.CodeBlockNodeAllOf1`,
		`	*m = result`,
		`func (m CodeBlockNode) MarshalJSON() ([]byte, error) {`,
		`	var b1, b2, b3 []byte`,
		`	var err error`,
		`	b1, err = json.Marshal(struct {`,
		`		CodeBlockNodeAllOf1`,
		`	}{`,
		`		CodeBlockNodeAllOf1: m.CodeBlockNodeAllOf1,`,
		`	})`,
		`	if err != nil {`,
		`		return nil, err`,
		`	b2, err = json.Marshal(struct {`,
		"		NodeType string `json:\"NodeType\"`",
		`	}{`,
		`		NodeType: m.NodeType(),`,
		`	})`,
		`	if err != nil {`,
		`		return nil, err`,
		`	return swag.ConcatJSON(b1, b2, b3), nil`,
		`func (m *CodeBlockNode) Validate(formats strfmt.Registry) error {`,
		`	if err := m.CodeBlockNodeAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

}

func initFixtureErrors() {
	// testing ../fixtures/bugs/1487/fixture-errors.yaml with expand (--skip-flatten)

	/*
		invalid specs supported by go-swagger
	*/

	f := newModelFixture("../fixtures/bugs/1487/fixture-errors.yaml", "broken spec to exercise error handling")
	flattenRun := f.AddRun(false)
	expandRun := f.AddRun(true)

	// load expectations for model: node.go
	flattenRun.AddExpectations("array_without_items.go", []string{
		`type ArrayWithoutItems []interface{}`,
		// empty validation
		"func (m ArrayWithoutItems) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		// NOTE would expect warning for a non-swagger compliant, but nonetheless supposed construct (not implemented)
		noLines,
		noLines)

	expandRun.AddExpectations("array_without_items.go", flattenRun.ExpectedFor("ArrayWithoutItems").ExpectedLines, todo, noLines, noLines)

	flattenRun.AddExpectations("multiple_types.go", []string{
		`type MultipleTypes interface{`,
	},
		// not expected
		validatable,
		// output in log
		// expect warning
		// warning,
		noLines,
		noLines)

	expandRun.AddExpectations("multiple_types.go", flattenRun.ExpectedFor("MultipleTypes").ExpectedLines, validatable, noLines, noLines)
}

func initTodolistSchemavalidation() {
	// testing todolist.schemavalidation.yaml with flatten and expand (--skip-flatten)

	/*
	   A very simple api description that makes a json only API to submit to do's.

	*/

	f := newModelFixture("../fixtures/codegen/todolist.schemavalidation.yml", "Private to-do list")
	flattenRun := f.AddRun(false)
	expandRun := f.AddRun(true)

	// load expectations for model: all_of_validations_meta_all_of6.go
	flattenRun.AddExpectations("all_of_validations_meta_all_of6.go", []string{
		`type AllOfValidationsMetaAllOf6 struct {`,
		"	Coords *AllOfValidationsMetaAllOf6Coords `json:\"coords,omitempty\"`",
		`func (m *AllOfValidationsMetaAllOf6) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateCoords(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AllOfValidationsMetaAllOf6) validateCoords(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Coords) {`,
		`	if m.Coords != nil {`,
		`		if err := m.Coords.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("coords"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: nested_array_validations.go
	flattenRun.AddExpectations("nested_array_validations.go", []string{
		`type NestedArrayValidations struct {`,
		"	Tags [][][]string `json:\"tags\"`",
		`func (m *NestedArrayValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateTags(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NestedArrayValidations) validateTags(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Tags) {`,
		`	iTagsSize := int64(len(m.Tags)`,
		`	if err := validate.MinItems("tags", "body", iTagsSize, 3); err != nil {`,
		`	if err := validate.MaxItems("tags", "body", iTagsSize, 10); err != nil {`,
		`	for i := 0; i < len(m.Tags); i++ {`,
		`		iiTagsSize := int64(len(m.Tags[i])`,
		`		if err := validate.MinItems("tags"+"."+strconv.Itoa(i), "body", iiTagsSize, 3); err != nil {`,
		`		if err := validate.MaxItems("tags"+"."+strconv.Itoa(i), "body", iiTagsSize, 10); err != nil {`,
		`		for ii := 0; ii < len(m.Tags[i]); ii++ {`,
		`			iiiTagsSize := int64(len(m.Tags[i][ii])`,
		`			if err := validate.MinItems("tags"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiTagsSize, 3); err != nil {`,
		`			if err := validate.MaxItems("tags"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiTagsSize, 10); err != nil {`,
		`			for iii := 0; iii < len(m.Tags[i][ii]); iii++ {`,
		`				if err := validate.MinLength("tags"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", m.Tags[i][ii][iii], 3); err != nil {`,
		`				if err := validate.MaxLength("tags"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", m.Tags[i][ii][iii], 10); err != nil {`,
		"				if err := validate.Pattern(\"tags\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii), \"body\", m.Tags[i][ii][iii], `\\w+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("nested_array_validations.go", flattenRun.ExpectedFor("NestedArrayValidations").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: all_of_validations_meta_all_of4.go
	flattenRun.AddExpectations("all_of_validations_meta_all_of4.go", []string{
		`type AllOfValidationsMetaAllOf4 struct {`,
		"	Opts map[string]int32 `json:\"opts,omitempty\"`",
		`func (m *AllOfValidationsMetaAllOf4) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateOpts(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AllOfValidationsMetaAllOf4) validateOpts(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Opts) {`,
		`	for k := range m.Opts {`,
		`		if err := validate.MinimumInt("opts"+"."+k, "body", int64(m.Opts[k]), 2, false); err != nil {`,
		`		if err := validate.MaximumInt("opts"+"."+k, "body", int64(m.Opts[k]), 50, false); err != nil {`,
		`		if err := validate.MultipleOf("opts"+"."+k, "body", float64(m.Opts[k]), 1.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: simple_zero_allowed.go
	flattenRun.AddExpectations("simple_zero_allowed.go", []string{
		`type SimpleZeroAllowed struct {`,
		"	ID string `json:\"id,omitempty\"`",
		"	Name *string `json:\"name\"`",
		"	Urls []string `json:\"urls\"`",
		`func (m *SimpleZeroAllowed) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateID(formats); err != nil {`,
		`	if err := m.validateName(formats); err != nil {`,
		`	if err := m.validateUrls(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *SimpleZeroAllowed) validateID(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ID) {`,
		`	if err := validate.MinLength("id", "body", m.ID, 2); err != nil {`,
		`	if err := validate.MaxLength("id", "body", m.ID, 50); err != nil {`,
		"	if err := validate.Pattern(\"id\", \"body\", m.ID, `[A-Za-z0-9][\\w- ]+`); err != nil {",
		`func (m *SimpleZeroAllowed) validateName(formats strfmt.Registry) error {`,
		`	if err := validate.Required("name", "body", m.Name); err != nil {`,
		`	if err := validate.MinLength("name", "body", *m.Name, 2); err != nil {`,
		`	if err := validate.MaxLength("name", "body", *m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", *m.Name, `[A-Za-z0-9][\\w- ]+`); err != nil {",
		`func (m *SimpleZeroAllowed) validateUrls(formats strfmt.Registry) error {`,
		`	if err := validate.Required("urls", "body", m.Urls); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("simple_zero_allowed.go", flattenRun.ExpectedFor("SimpleZeroAllowed").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: named_all_of_all_of6_coords_all_of0.go
	flattenRun.AddExpectations("named_all_of_all_of6_coords_all_of0.go", []string{
		`type NamedAllOfAllOf6CoordsAllOf0 struct {`,
		"	Name string `json:\"name,omitempty\"`",
		`func (m *NamedAllOfAllOf6CoordsAllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateName(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedAllOfAllOf6CoordsAllOf0) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 2); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `[A-Za-z0-9][\\w- ]+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_all_of_all_of6.go
	flattenRun.AddExpectations("named_all_of_all_of6.go", []string{
		`type NamedAllOfAllOf6 struct {`,
		"	Coords *NamedAllOfAllOf6Coords `json:\"coords,omitempty\"`",
		`func (m *NamedAllOfAllOf6) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateCoords(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedAllOfAllOf6) validateCoords(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Coords) {`,
		`	if m.Coords != nil {`,
		`		if err := m.Coords.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("coords"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_array_multi.go
	flattenRun.AddExpectations("named_array_multi.go", []string{
		`type NamedArrayMulti struct {`,
		"	P0 *string `json:\"-\"`",
		"	P1 *float64 `json:\"-\"`",
		`func (m *NamedArrayMulti) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedArrayMulti) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	if err := validate.MinLength("0", "body", *m.P0, 3); err != nil {`,
		`	if err := validate.MaxLength("0", "body", *m.P0, 10); err != nil {`,
		"	if err := validate.Pattern(\"0\", \"body\", *m.P0, `\\w+`); err != nil {",
		`func (m *NamedArrayMulti) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("1", "body", m.P1); err != nil {`,
		`	if err := validate.Minimum("1", "body", *m.P1, 3, false); err != nil {`,
		`	if err := validate.Maximum("1", "body", *m.P1, 12, false); err != nil {`,
		`	if err := validate.MultipleOf("1", "body", *m.P1, 1.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("named_array_multi.go", flattenRun.ExpectedFor("NamedArrayMulti").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: named_array.go
	flattenRun.AddExpectations("named_array.go", []string{
		`type NamedArray []string`,
		`func (m NamedArray) Validate(formats strfmt.Registry) error {`,
		`	iNamedArraySize := int64(len(m)`,
		`	if err := validate.MinItems("", "body", iNamedArraySize, 3); err != nil {`,
		`	if err := validate.MaxItems("", "body", iNamedArraySize, 10); err != nil {`,
		`	for i := 0; i < len(m); i++ {`,
		`		if err := validate.MinLength(strconv.Itoa(i), "body", m[i], 3); err != nil {`,
		`		if err := validate.MaxLength(strconv.Itoa(i), "body", m[i], 10); err != nil {`,
		"		if err := validate.Pattern(strconv.Itoa(i), \"body\", m[i], `\\w+`); err != nil {",
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("named_array.go", flattenRun.ExpectedFor("NamedArray").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: named_number.go
	flattenRun.AddExpectations("named_number.go", []string{
		`type NamedNumber int32`,
		`func (m NamedNumber) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.MinimumInt("", "body", int64(m), 0, true); err != nil {`,
		`	if err := validate.MaximumInt("", "body", int64(m), 500, false); err != nil {`,
		`	if err := validate.MultipleOf("", "body", float64(m), 1.5); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("named_number.go", flattenRun.ExpectedFor("NamedNumber").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: nested_map_validations.go
	flattenRun.AddExpectations("nested_map_validations.go", []string{
		`type NestedMapValidations struct {`,
		"	Meta map[string]map[string]map[string]int64 `json:\"meta,omitempty\"`",
		`func (m *NestedMapValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMeta(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NestedMapValidations) validateMeta(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Meta) {`,
		`	for k := range m.Meta {`,
		`		for kk := range m.Meta[k] {`,
		`			for kkk := range m.Meta[k][kk] {`,
		`				if err := validate.MinimumInt("meta"+"."+k+"."+kk+"."+kkk, "body", m.Meta[k][kk][kkk], 3, false); err != nil {`,
		`				if err := validate.MaximumInt("meta"+"."+k+"."+kk+"."+kkk, "body", m.Meta[k][kk][kkk], 6, false); err != nil {`,
		`				if err := validate.MultipleOfInt("meta"+"."+k+"."+kk+"."+kkk, "body", m.Meta[k][kk][kkk], 1); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("nested_map_validations.go", flattenRun.ExpectedFor("NestedMapValidations").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: array_multi_validations_args.go
	flattenRun.AddExpectations("array_multi_validations_args.go", []string{
		`type ArrayMultiValidationsArgs struct {`,
		"	P0 *string `json:\"-\"`",
		"	P1 *float64 `json:\"-\"`",
		`func (m *ArrayMultiValidationsArgs) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ArrayMultiValidationsArgs) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	if err := validate.MinLength("0", "body", *m.P0, 3); err != nil {`,
		`	if err := validate.MaxLength("0", "body", *m.P0, 10); err != nil {`,
		"	if err := validate.Pattern(\"0\", \"body\", *m.P0, `\\w+`); err != nil {",
		`func (m *ArrayMultiValidationsArgs) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("1", "body", m.P1); err != nil {`,
		`	if err := validate.Minimum("1", "body", *m.P1, 3, false); err != nil {`,
		`	if err := validate.Maximum("1", "body", *m.P1, 12, false); err != nil {`,
		`	if err := validate.MultipleOf("1", "body", *m.P1, 1.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_map_complex_additional_properties.go
	flattenRun.AddExpectations("named_map_complex_additional_properties.go", []string{
		`type NamedMapComplexAdditionalProperties struct {`,
		"	Age int32 `json:\"age,omitempty\"`",
		"	Name string `json:\"name,omitempty\"`",
		`func (m *NamedMapComplexAdditionalProperties) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`	if err := m.validateName(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedMapComplexAdditionalProperties) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 1, true); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 200, true); err != nil {`,
		`	if err := validate.MultipleOfInt("age", "body", int64(m.Age), 1); err != nil {`,
		`func (m *NamedMapComplexAdditionalProperties) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 10); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `\\w+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_nested_map_complex.go
	flattenRun.AddExpectations("named_nested_map_complex.go", []string{
		// maps are now simple types
		`type NamedNestedMapComplex map[string]map[string]map[string]NamedNestedMapComplexAdditionalPropertiesAdditionalPropertiesAdditionalProperties`,
		`func (m NamedNestedMapComplex) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		for kk := range m[k] {`,
		`			for kkk := range m[k][kk] {`,
		`				if err := validate.Required(k+"."+kk+"."+kkk, "body", m[k][kk][kkk]); err != nil {`,
		`				if val, ok := m[k][kk][kkk]; ok {`,
		`					if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("named_nested_map_complex.go", []string{
		`type NamedNestedMapComplex map[string]map[string]map[string]NamedNestedMapComplexAnon`,
		`func (m NamedNestedMapComplex) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		for kk := range m[k] {`,
		`			for kkk := range m[k][kk] {`,
		`				if val, ok := m[k][kk][kkk]; ok {`,
		`					if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`type NamedNestedMapComplexAnon struct {`,
		"	Age int32 `json:\"age,omitempty\"`",
		"	Name string `json:\"name,omitempty\"`",
		`func (m *NamedNestedMapComplexAnon) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`	if err := m.validateName(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedNestedMapComplexAnon) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 1, true); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 200, true); err != nil {`,
		`	if err := validate.MultipleOfInt("age", "body", int64(m.Age), 1); err != nil {`,
		`func (m *NamedNestedMapComplexAnon) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 10); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `\\w+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: all_of_validations_meta_all_of1.go
	flattenRun.AddExpectations("all_of_validations_meta_all_of1.go", []string{
		`type AllOfValidationsMetaAllOf1 struct {`,
		"	Age int32 `json:\"age,omitempty\"`",
		`func (m *AllOfValidationsMetaAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AllOfValidationsMetaAllOf1) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 2, false); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 50, false); err != nil {`,
		`	if err := validate.MultipleOf("age", "body", float64(m.Age), 1.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: nested_map_complex_validations_meta_additional_properties_additional_properties.go
	// NOTE(fredbi): maps are now simple types - this definition is no more generated

	// load expectations for model: tag.go
	flattenRun.AddExpectations("tag.go", []string{
		`type Tag struct {`,
		"	ID int64 `json:\"id,omitempty\"`",
		"	Name string `json:\"name,omitempty\"`",
		// empty validation
		"func (m *Tag) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("tag.go", flattenRun.ExpectedFor("Tag").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: nested_object_validations_args.go
	flattenRun.AddExpectations("nested_object_validations_args.go", []string{
		`type NestedObjectValidationsArgs struct {`,
		"	Meta *NestedObjectValidationsArgsMeta `json:\"meta,omitempty\"`",
		`func (m *NestedObjectValidationsArgs) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMeta(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NestedObjectValidationsArgs) validateMeta(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Meta) {`,
		`	if m.Meta != nil {`,
		`		if err := m.Meta.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("meta"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_all_of_all_of6_coords_all_of1.go
	flattenRun.AddExpectations("named_all_of_all_of6_coords_all_of1.go", []string{
		`type NamedAllOfAllOf6CoordsAllOf1 struct {`,
		"	Age int32 `json:\"age,omitempty\"`",
		`func (m *NamedAllOfAllOf6CoordsAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedAllOfAllOf6CoordsAllOf1) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 2, false); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 50, false); err != nil {`,
		`	if err := validate.MultipleOf("age", "body", float64(m.Age), 1.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_all_of_all_of6_coords.go
	flattenRun.AddExpectations("named_all_of_all_of6_coords.go", []string{
		`type NamedAllOfAllOf6Coords struct {`,
		`	NamedAllOfAllOf6CoordsAllOf0`,
		`	NamedAllOfAllOf6CoordsAllOf1`,
		`func (m *NamedAllOfAllOf6Coords) Validate(formats strfmt.Registry) error {`,
		`	if err := m.NamedAllOfAllOf6CoordsAllOf0.Validate(formats); err != nil {`,
		`	if err := m.NamedAllOfAllOf6CoordsAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: array_multi_validations.go
	flattenRun.AddExpectations("array_multi_validations.go", []string{
		`type ArrayMultiValidations struct {`,
		"	Args ArrayMultiValidationsArgs `json:\"args,omitempty\"`",
		`func (m *ArrayMultiValidations) Validate(formats strfmt.Registry) error {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("array_multi_validations.go", []string{
		`type ArrayMultiValidations struct {`,
		"	Args *ArrayMultiValidationsArgsTuple0 `json:\"args,omitempty\"`",
		`func (m *ArrayMultiValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateArgs(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ArrayMultiValidations) validateArgs(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Args) {`,
		`	if m.Args != nil {`,
		`		if err := m.Args.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("args"`,
		`type ArrayMultiValidationsArgsTuple0 struct {`,
		"	P0 *string `json:\"-\"`",
		"	P1 *float64 `json:\"-\"`",
		`func (m *ArrayMultiValidationsArgsTuple0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ArrayMultiValidationsArgsTuple0) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P0", "body", m.P0); err != nil {`,
		`	if err := validate.MinLength("P0", "body", *m.P0, 3); err != nil {`,
		`	if err := validate.MaxLength("P0", "body", *m.P0, 10); err != nil {`,
		"	if err := validate.Pattern(\"P0\", \"body\", *m.P0, `\\w+`); err != nil {",
		`func (m *ArrayMultiValidationsArgsTuple0) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P1", "body", m.P1); err != nil {`,
		`	if err := validate.Minimum("P1", "body", *m.P1, 3, false); err != nil {`,
		`	if err := validate.Maximum("P1", "body", *m.P1, 12, false); err != nil {`,
		`	if err := validate.MultipleOf("P1", "body", *m.P1, 1.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: string_validations.go
	flattenRun.AddExpectations("string_validations.go", []string{
		`type StringValidations struct {`,
		"	Name string `json:\"name,omitempty\"`",
		`func (m *StringValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateName(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *StringValidations) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 2); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `[A-Za-z0-9][\\w- ]+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("string_validations.go", flattenRun.ExpectedFor("StringValidations").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: required_props.go
	flattenRun.AddExpectations("required_props.go", []string{
		`type RequiredProps struct {`,
		"	Age *int32 `json:\"age\"`",
		"	CreatedAt *strfmt.DateTime `json:\"createdAt\"`",
		"	ID *int64 `json:\"id\"`",
		"	Name *string `json:\"name\"`",
		"	Score *float32 `json:\"score\"`",
		"	Tags []string `json:\"tags\"`",
		`func (m *RequiredProps) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`	if err := m.validateCreatedAt(formats); err != nil {`,
		`	if err := m.validateID(formats); err != nil {`,
		`	if err := m.validateName(formats); err != nil {`,
		`	if err := m.validateScore(formats); err != nil {`,
		`	if err := m.validateTags(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *RequiredProps) validateAge(formats strfmt.Registry) error {`,
		`	if err := validate.Required("age", "body", m.Age); err != nil {`,
		`func (m *RequiredProps) validateCreatedAt(formats strfmt.Registry) error {`,
		`	if err := validate.Required("createdAt", "body", m.CreatedAt); err != nil {`,
		`	if err := validate.FormatOf("createdAt", "body", "date-time", m.CreatedAt.String(), formats); err != nil {`,
		`func (m *RequiredProps) validateID(formats strfmt.Registry) error {`,
		`	if err := validate.Required("id", "body", m.ID); err != nil {`,
		`func (m *RequiredProps) validateName(formats strfmt.Registry) error {`,
		`	if err := validate.Required("name", "body", m.Name); err != nil {`,
		`func (m *RequiredProps) validateScore(formats strfmt.Registry) error {`,
		`	if err := validate.Required("score", "body", m.Score); err != nil {`,
		`func (m *RequiredProps) validateTags(formats strfmt.Registry) error {`,
		`	if err := validate.Required("tags", "body", m.Tags); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("required_props.go", flattenRun.ExpectedFor("RequiredProps").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: named_all_of_all_of5.go
	flattenRun.AddExpectations("named_all_of_all_of5.go", []string{
		`type NamedAllOfAllOf5 struct {`,
		"	ExtOpts map[string]map[string]map[string]int32 `json:\"extOpts,omitempty\"`",
		`func (m *NamedAllOfAllOf5) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateExtOpts(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedAllOfAllOf5) validateExtOpts(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ExtOpts) {`,
		`	for k := range m.ExtOpts {`,
		`		for kk := range m.ExtOpts[k] {`,
		`			for kkk := range m.ExtOpts[k][kk] {`,
		`				if err := validate.MinimumInt("extOpts"+"."+k+"."+kk+"."+kkk, "body", int64(m.ExtOpts[k][kk][kkk]), 2, false); err != nil {`,
		`				if err := validate.MaximumInt("extOpts"+"."+k+"."+kk+"."+kkk, "body", int64(m.ExtOpts[k][kk][kkk]), 50, false); err != nil {`,
		`				if err := validate.MultipleOf("extOpts"+"."+k+"."+kk+"."+kkk, "body", float64(m.ExtOpts[k][kk][kkk]), 1.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_map.go
	flattenRun.AddExpectations("named_map.go", []string{
		`type NamedMap map[string]int64`,
		`func (m NamedMap) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if err := validate.MinimumInt(k, "body", m[k], 3, false); err != nil {`,
		`		if err := validate.MaximumInt(k, "body", m[k], 6, false); err != nil {`,
		`		if err := validate.MultipleOfInt(k, "body", m[k], 1); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("named_map.go", flattenRun.ExpectedFor("NamedMap").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: named_string.go
	flattenRun.AddExpectations("named_string.go", []string{
		`type NamedString string`,
		`func (m NamedString) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.MinLength("", "body", string(m), 2); err != nil {`,
		`	if err := validate.MaxLength("", "body", string(m), 50); err != nil {`,
		"	if err := validate.Pattern(\"\", \"body\", string(m), `[A-Za-z0-9][\\w- ]+`); err != nil {",
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("named_string.go", flattenRun.ExpectedFor("NamedString").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: named_all_of_all_of3.go
	flattenRun.AddExpectations("named_all_of_all_of3.go", []string{
		`type NamedAllOfAllOf3 struct {`,
		"	Assoc [][][]string `json:\"assoc\"`",
		`func (m *NamedAllOfAllOf3) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAssoc(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedAllOfAllOf3) validateAssoc(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Assoc) {`,
		`	iAssocSize := int64(len(m.Assoc)`,
		`	if err := validate.MinItems("assoc", "body", iAssocSize, 5); err != nil {`,
		`	if err := validate.MaxItems("assoc", "body", iAssocSize, 20); err != nil {`,
		`	for i := 0; i < len(m.Assoc); i++ {`,
		`		iiAssocSize := int64(len(m.Assoc[i])`,
		`		if err := validate.MinItems("assoc"+"."+strconv.Itoa(i), "body", iiAssocSize, 5); err != nil {`,
		`		if err := validate.MaxItems("assoc"+"."+strconv.Itoa(i), "body", iiAssocSize, 20); err != nil {`,
		`		for ii := 0; ii < len(m.Assoc[i]); ii++ {`,
		`			iiiAssocSize := int64(len(m.Assoc[i][ii])`,
		`			if err := validate.MinItems("assoc"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiAssocSize, 5); err != nil {`,
		`			if err := validate.MaxItems("assoc"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiAssocSize, 20); err != nil {`,
		`			for iii := 0; iii < len(m.Assoc[i][ii]); iii++ {`,
		`				if err := validate.MinLength("assoc"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", m.Assoc[i][ii][iii], 2); err != nil {`,
		`				if err := validate.MaxLength("assoc"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", m.Assoc[i][ii][iii], 50); err != nil {`,
		"				if err := validate.Pattern(\"assoc\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii), \"body\", m.Assoc[i][ii][iii], `[A-Za-z0-9][\\w- ]+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: map_complex_validations.go
	flattenRun.AddExpectations("map_complex_validations.go", []string{
		`type MapComplexValidations struct {`,
		// maps are now simple types
		"Meta map[string]MapComplexValidationsMetaAdditionalProperties `json:\"meta,omitempty\"`",
		`func (m *MapComplexValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMeta(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *MapComplexValidations) validateMeta(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Meta) {`,
		`            		for k := range m.Meta {`,
		`            			if val, ok := m.Meta[k]; ok {`,
		`            				if err := val.Validate(formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("map_complex_validations.go", []string{
		`type MapComplexValidations struct {`,
		"	Meta map[string]MapComplexValidationsMetaAnon `json:\"meta,omitempty\"`",
		`func (m *MapComplexValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMeta(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *MapComplexValidations) validateMeta(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Meta) {`,
		`	for k := range m.Meta {`,
		`		if val, ok := m.Meta[k]; ok {`,
		`			if err := val.Validate(formats); err != nil {`,
		`type MapComplexValidationsMetaAnon struct {`,
		"	Age int32 `json:\"age,omitempty\"`",
		"	Name string `json:\"name,omitempty\"`",
		`func (m *MapComplexValidationsMetaAnon) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`	if err := m.validateName(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *MapComplexValidationsMetaAnon) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 1, true); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 200, true); err != nil {`,
		`	if err := validate.MultipleOfInt("age", "body", int64(m.Age), 1); err != nil {`,
		`func (m *MapComplexValidationsMetaAnon) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 10); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `\\w+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_nested_map_complex_additional_properties_additional_properties_additional_properties.go
	flattenRun.AddExpectations("named_nested_map_complex_additional_properties_additional_properties_additional_properties.go", []string{
		`type NamedNestedMapComplexAdditionalPropertiesAdditionalPropertiesAdditionalProperties struct {`,
		"	Age int32 `json:\"age,omitempty\"`",
		"	Name string `json:\"name,omitempty\"`",
		`func (m *NamedNestedMapComplexAdditionalPropertiesAdditionalPropertiesAdditionalProperties) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`	if err := m.validateName(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedNestedMapComplexAdditionalPropertiesAdditionalPropertiesAdditionalProperties) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 1, true); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 200, true); err != nil {`,
		`	if err := validate.MultipleOfInt("age", "body", int64(m.Age), 1); err != nil {`,
		`func (m *NamedNestedMapComplexAdditionalPropertiesAdditionalPropertiesAdditionalProperties) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 10); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `\\w+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: all_of_validations_meta_all_of6_coords.go
	flattenRun.AddExpectations("all_of_validations_meta_all_of6_coords.go", []string{
		`type AllOfValidationsMetaAllOf6Coords struct {`,
		`	AllOfValidationsMetaAllOf6CoordsAllOf0`,
		`	AllOfValidationsMetaAllOf6CoordsAllOf1`,
		`func (m *AllOfValidationsMetaAllOf6Coords) Validate(formats strfmt.Registry) error {`,
		`	if err := m.AllOfValidationsMetaAllOf6CoordsAllOf0.Validate(formats); err != nil {`,
		`	if err := m.AllOfValidationsMetaAllOf6CoordsAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: array_validations.go
	flattenRun.AddExpectations("array_validations.go", []string{
		`type ArrayValidations struct {`,
		"	Tags []string `json:\"tags\"`",
		`func (m *ArrayValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateTags(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ArrayValidations) validateTags(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Tags) {`,
		`	iTagsSize := int64(len(m.Tags)`,
		`	if err := validate.MinItems("tags", "body", iTagsSize, 3); err != nil {`,
		`	if err := validate.MaxItems("tags", "body", iTagsSize, 10); err != nil {`,
		`	for i := 0; i < len(m.Tags); i++ {`,
		`		if err := validate.MinLength("tags"+"."+strconv.Itoa(i), "body", m.Tags[i], 3); err != nil {`,
		`		if err := validate.MaxLength("tags"+"."+strconv.Itoa(i), "body", m.Tags[i], 10); err != nil {`,
		"		if err := validate.Pattern(\"tags\"+\".\"+strconv.Itoa(i), \"body\", m.Tags[i], `\\w+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("array_validations.go", flattenRun.ExpectedFor("ArrayValidations").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: all_of_validations_meta.go
	flattenRun.AddExpectations("all_of_validations_meta.go", []string{
		`type AllOfValidationsMeta struct {`,
		`	AllOfValidationsMetaAllOf0`,
		`	AllOfValidationsMetaAllOf1`,
		`	AllOfValidationsMetaAllOf2`,
		`	AllOfValidationsMetaAllOf3`,
		`	AllOfValidationsMetaAllOf4`,
		`	AllOfValidationsMetaAllOf5`,
		`	AllOfValidationsMetaAllOf6`,
		`func (m *AllOfValidationsMeta) Validate(formats strfmt.Registry) error {`,
		`	if err := m.AllOfValidationsMetaAllOf0.Validate(formats); err != nil {`,
		`	if err := m.AllOfValidationsMetaAllOf1.Validate(formats); err != nil {`,
		`	if err := m.AllOfValidationsMetaAllOf2.Validate(formats); err != nil {`,
		`	if err := m.AllOfValidationsMetaAllOf3.Validate(formats); err != nil {`,
		`	if err := m.AllOfValidationsMetaAllOf4.Validate(formats); err != nil {`,
		`	if err := m.AllOfValidationsMetaAllOf5.Validate(formats); err != nil {`,
		`	if err := m.AllOfValidationsMetaAllOf6.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: map_complex_validations_meta_additional_properties.go
	flattenRun.AddExpectations("map_complex_validations_meta_additional_properties.go", []string{
		`type MapComplexValidationsMetaAdditionalProperties struct {`,
		"	Age int32 `json:\"age,omitempty\"`",
		"	Name string `json:\"name,omitempty\"`",
		`func (m *MapComplexValidationsMetaAdditionalProperties) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`	if err := m.validateName(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *MapComplexValidationsMetaAdditionalProperties) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 1, true); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 200, true); err != nil {`,
		`	if err := validate.MultipleOfInt("age", "body", int64(m.Age), 1); err != nil {`,
		`func (m *MapComplexValidationsMetaAdditionalProperties) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 10); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `\\w+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: nested_map_complex_validations_meta_additional_properties.go
	// NOTE(fredbi): maps are simple types and this definition is no more generated

	// load expectations for model: all_of_validations_meta_all_of6_coords_all_of1.go
	flattenRun.AddExpectations("all_of_validations_meta_all_of6_coords_all_of1.go", []string{
		`type AllOfValidationsMetaAllOf6CoordsAllOf1 struct {`,
		"	Age int32 `json:\"age,omitempty\"`",
		`func (m *AllOfValidationsMetaAllOf6CoordsAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AllOfValidationsMetaAllOf6CoordsAllOf1) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 2, false); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 50, false); err != nil {`,
		`	if err := validate.MultipleOf("age", "body", float64(m.Age), 1.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: all_of_validations_meta_all_of3.go
	flattenRun.AddExpectations("all_of_validations_meta_all_of3.go", []string{
		`type AllOfValidationsMetaAllOf3 struct {`,
		"	Assoc [][][]string `json:\"assoc\"`",
		`func (m *AllOfValidationsMetaAllOf3) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAssoc(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AllOfValidationsMetaAllOf3) validateAssoc(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Assoc) {`,
		`	iAssocSize := int64(len(m.Assoc)`,
		`	if err := validate.MinItems("assoc", "body", iAssocSize, 5); err != nil {`,
		`	if err := validate.MaxItems("assoc", "body", iAssocSize, 20); err != nil {`,
		`	for i := 0; i < len(m.Assoc); i++ {`,
		`		iiAssocSize := int64(len(m.Assoc[i])`,
		`		if err := validate.MinItems("assoc"+"."+strconv.Itoa(i), "body", iiAssocSize, 5); err != nil {`,
		`		if err := validate.MaxItems("assoc"+"."+strconv.Itoa(i), "body", iiAssocSize, 20); err != nil {`,
		`		for ii := 0; ii < len(m.Assoc[i]); ii++ {`,
		`			iiiAssocSize := int64(len(m.Assoc[i][ii])`,
		`			if err := validate.MinItems("assoc"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiAssocSize, 5); err != nil {`,
		`			if err := validate.MaxItems("assoc"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiAssocSize, 20); err != nil {`,
		`			for iii := 0; iii < len(m.Assoc[i][ii]); iii++ {`,
		`				if err := validate.MinLength("assoc"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", m.Assoc[i][ii][iii], 2); err != nil {`,
		`				if err := validate.MaxLength("assoc"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", m.Assoc[i][ii][iii], 50); err != nil {`,
		"				if err := validate.Pattern(\"assoc\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii), \"body\", m.Assoc[i][ii][iii], `[A-Za-z0-9][\\w- ]+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: nested_object_validations.go
	flattenRun.AddExpectations("nested_object_validations.go", []string{
		`type NestedObjectValidations struct {`,
		"	Args *NestedObjectValidationsArgs `json:\"args,omitempty\"`",
		`func (m *NestedObjectValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateArgs(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NestedObjectValidations) validateArgs(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Args) {`,
		`	if m.Args != nil {`,
		`		if err := m.Args.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("args"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("nested_object_validations.go", []string{
		`type NestedObjectValidations struct {`,
		"	Args *NestedObjectValidationsArgs `json:\"args,omitempty\"`",
		`func (m *NestedObjectValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateArgs(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NestedObjectValidations) validateArgs(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Args) {`,
		`	if m.Args != nil {`,
		`		if err := m.Args.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("args"`,
		`type NestedObjectValidationsArgs struct {`,
		"	Meta *NestedObjectValidationsArgsMeta `json:\"meta,omitempty\"`",
		`func (m *NestedObjectValidationsArgs) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMeta(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NestedObjectValidationsArgs) validateMeta(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Meta) {`,
		`	if m.Meta != nil {`,
		`		if err := m.Meta.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("args" + "." + "meta"`,
		`type NestedObjectValidationsArgsMeta struct {`,
		"	First string `json:\"first,omitempty\"`",
		"	Fourth [][][]float32 `json:\"fourth\"`",
		"	Second float64 `json:\"second,omitempty\"`",
		"	Third []float32 `json:\"third\"`",
		`func (m *NestedObjectValidationsArgsMeta) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateFirst(formats); err != nil {`,
		`	if err := m.validateFourth(formats); err != nil {`,
		`	if err := m.validateSecond(formats); err != nil {`,
		`	if err := m.validateThird(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NestedObjectValidationsArgsMeta) validateFirst(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.First) {`,
		`	if err := validate.MinLength("args"+"."+"meta"+"."+"first", "body", m.First, 2); err != nil {`,
		`	if err := validate.MaxLength("args"+"."+"meta"+"."+"first", "body", m.First, 50); err != nil {`,
		"	if err := validate.Pattern(\"args\"+\".\"+\"meta\"+\".\"+\"first\", \"body\", m.First, `^\\w+`); err != nil {",
		`func (m *NestedObjectValidationsArgsMeta) validateFourth(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Fourth) {`,
		`	iFourthSize := int64(len(m.Fourth)`,
		`	if err := validate.MinItems("args"+"."+"meta"+"."+"fourth", "body", iFourthSize, 5); err != nil {`,
		`	if err := validate.MaxItems("args"+"."+"meta"+"."+"fourth", "body", iFourthSize, 93); err != nil {`,
		`	for i := 0; i < len(m.Fourth); i++ {`,
		`		iiFourthSize := int64(len(m.Fourth[i])`,
		`		if err := validate.MinItems("args"+"."+"meta"+"."+"fourth"+"."+strconv.Itoa(i), "body", iiFourthSize, 5); err != nil {`,
		`		if err := validate.MaxItems("args"+"."+"meta"+"."+"fourth"+"."+strconv.Itoa(i), "body", iiFourthSize, 93); err != nil {`,
		`		for ii := 0; ii < len(m.Fourth[i]); ii++ {`,
		`			iiiFourthSize := int64(len(m.Fourth[i][ii])`,
		`			if err := validate.MinItems("args"+"."+"meta"+"."+"fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiFourthSize, 5); err != nil {`,
		`			if err := validate.MaxItems("args"+"."+"meta"+"."+"fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiFourthSize, 93); err != nil {`,
		`			for iii := 0; iii < len(m.Fourth[i][ii]); iii++ {`,
		`				if err := validate.Minimum("args"+"."+"meta"+"."+"fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", float64(m.Fourth[i][ii][iii]), 3, false); err != nil {`,
		`				if err := validate.Maximum("args"+"."+"meta"+"."+"fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", float64(m.Fourth[i][ii][iii]), 6, false); err != nil {`,
		`				if err := validate.MultipleOf("args"+"."+"meta"+"."+"fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", float64(m.Fourth[i][ii][iii]), 0.5); err != nil {`,
		`func (m *NestedObjectValidationsArgsMeta) validateSecond(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Second) {`,
		`	if err := validate.Minimum("args"+"."+"meta"+"."+"second", "body", m.Second, 3, false); err != nil {`,
		`	if err := validate.Maximum("args"+"."+"meta"+"."+"second", "body", m.Second, 51, false); err != nil {`,
		`	if err := validate.MultipleOf("args"+"."+"meta"+"."+"second", "body", m.Second, 1.5); err != nil {`,
		`func (m *NestedObjectValidationsArgsMeta) validateThird(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Third) {`,
		`	iThirdSize := int64(len(m.Third)`,
		`	if err := validate.MinItems("args"+"."+"meta"+"."+"third", "body", iThirdSize, 5); err != nil {`,
		`	if err := validate.MaxItems("args"+"."+"meta"+"."+"third", "body", iThirdSize, 93); err != nil {`,
		`	for i := 0; i < len(m.Third); i++ {`,
		`		if err := validate.Minimum("args"+"."+"meta"+"."+"third"+"."+strconv.Itoa(i), "body", float64(m.Third[i]), 3, false); err != nil {`,
		`		if err := validate.Maximum("args"+"."+"meta"+"."+"third"+"."+strconv.Itoa(i), "body", float64(m.Third[i]), 6, false); err != nil {`,
		`		if err := validate.MultipleOf("args"+"."+"meta"+"."+"third"+"."+strconv.Itoa(i), "body", float64(m.Third[i]), 0.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_array_additional.go
	flattenRun.AddExpectations("named_array_additional.go", []string{
		`type NamedArrayAdditional struct {`,
		"	P0 *string `json:\"-\"`",
		"	P1 *float64 `json:\"-\"`",
		"	NamedArrayAdditionalItems []int64 `json:\"-\"`",
		`func (m *NamedArrayAdditional) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateNamedArrayAdditionalItems(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedArrayAdditional) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	if err := validate.MinLength("0", "body", *m.P0, 3); err != nil {`,
		`	if err := validate.MaxLength("0", "body", *m.P0, 10); err != nil {`,
		"	if err := validate.Pattern(\"0\", \"body\", *m.P0, `\\w+`); err != nil {",
		`func (m *NamedArrayAdditional) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("1", "body", m.P1); err != nil {`,
		`	if err := validate.Minimum("1", "body", *m.P1, 3, false); err != nil {`,
		`	if err := validate.Maximum("1", "body", *m.P1, 12, false); err != nil {`,
		`	if err := validate.MultipleOf("1", "body", *m.P1, 1.5); err != nil {`,
		`func (m *NamedArrayAdditional) validateNamedArrayAdditionalItems(formats strfmt.Registry) error {`,
		`	for i := range m.NamedArrayAdditionalItems {`,
		`		if err := validate.MinimumInt(strconv.Itoa(i+2), "body", m.NamedArrayAdditionalItems[i], 3, false); err != nil {`,
		`		if err := validate.MaximumInt(strconv.Itoa(i+2), "body", m.NamedArrayAdditionalItems[i], 6, false); err != nil {`,
		`		if err := validate.MultipleOfInt(strconv.Itoa(i+2), "body", m.NamedArrayAdditionalItems[i], 1); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("named_array_additional.go", flattenRun.ExpectedFor("NamedArrayAdditional").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: pet.go
	flattenRun.AddExpectations("pet.go", []string{
		`type Pet struct {`,
		"	Category *Category `json:\"category,omitempty\"`",
		"	ID int64 `json:\"id,omitempty\"`",
		"	Name *string `json:\"name\"`",
		"	PhotoUrls []string `json:\"photoUrls\" xml:\"photoUrl\"`",
		"	Status string `json:\"status,omitempty\"`",
		"	Tags []*Tag `json:\"tags\" xml:\"tag\"`",
		`func (m *Pet) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateCategory(formats); err != nil {`,
		`	if err := m.validateName(formats); err != nil {`,
		`	if err := m.validatePhotoUrls(formats); err != nil {`,
		`	if err := m.validateStatus(formats); err != nil {`,
		`	if err := m.validateTags(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *Pet) validateCategory(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Category) {`,
		`	if m.Category != nil {`,
		`		if err := m.Category.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("category"`,
		`func (m *Pet) validateName(formats strfmt.Registry) error {`,
		`	if err := validate.Required("name", "body", m.Name); err != nil {`,
		`func (m *Pet) validatePhotoUrls(formats strfmt.Registry) error {`,
		`	if err := validate.Required("photoUrls", "body", m.PhotoUrls); err != nil {`,
		`var petTypeStatusPropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"available\",\"pending\",\"sold\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		petTypeStatusPropEnum = append(petTypeStatusPropEnum, v`,
		`	PetStatusAvailable string = "available"`,
		`	PetStatusPending string = "pending"`,
		`	PetStatusSold string = "sold"`,
		`func (m *Pet) validateStatusEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, petTypeStatusPropEnum, true); err != nil {`,
		`func (m *Pet) validateStatus(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Status) {`,
		`	if err := m.validateStatusEnum("status", "body", m.Status); err != nil {`,
		`func (m *Pet) validateTags(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Tags) {`,
		`	for i := 0; i < len(m.Tags); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`		if swag.IsZero(m.Tags[i]) {`,
		// nullable required:
		`		if m.Tags[i] != nil {`,
		`			if err := m.Tags[i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName("tags" + "." + strconv.Itoa(i)`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("pet.go", []string{
		`type Pet struct {`,
		"	Category *PetCategory `json:\"category,omitempty\"`",
		"	ID int64 `json:\"id,omitempty\"`",
		"	Name *string `json:\"name\"`",
		"	PhotoUrls []string `json:\"photoUrls\" xml:\"photoUrl\"`",
		"	Status string `json:\"status,omitempty\"`",
		"	Tags []*PetTagsItems0 `json:\"tags\" xml:\"tag\"`",
		`func (m *Pet) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateCategory(formats); err != nil {`,
		`	if err := m.validateName(formats); err != nil {`,
		`	if err := m.validatePhotoUrls(formats); err != nil {`,
		`	if err := m.validateStatus(formats); err != nil {`,
		`	if err := m.validateTags(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *Pet) validateCategory(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Category) {`,
		`	if m.Category != nil {`,
		`		if err := m.Category.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("category"`,
		`func (m *Pet) validateName(formats strfmt.Registry) error {`,
		`	if err := validate.Required("name", "body", m.Name); err != nil {`,
		`func (m *Pet) validatePhotoUrls(formats strfmt.Registry) error {`,
		`	if err := validate.Required("photoUrls", "body", m.PhotoUrls); err != nil {`,
		`var petTypeStatusPropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"available\",\"pending\",\"sold\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		petTypeStatusPropEnum = append(petTypeStatusPropEnum, v`,
		`	PetStatusAvailable string = "available"`,
		`	PetStatusPending string = "pending"`,
		`	PetStatusSold string = "sold"`,
		`func (m *Pet) validateStatusEnum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, petTypeStatusPropEnum, true); err != nil {`,
		`func (m *Pet) validateStatus(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Status) {`,
		`	if err := m.validateStatusEnum("status", "body", m.Status); err != nil {`,
		`func (m *Pet) validateTags(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Tags) {`,
		`	for i := 0; i < len(m.Tags); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`		if swag.IsZero(m.Tags[i]) {`,
		// nullable required:
		`		if m.Tags[i] != nil {`,
		`			if err := m.Tags[i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName("tags" + "." + strconv.Itoa(i)`,
		`type PetCategory struct {`,
		"	ID int64 `json:\"id,omitempty\"`",
		"	Name string `json:\"name,omitempty\"`",
		// empty validation
		"func (m *PetCategory) Validate(formats strfmt.Registry) error {\n	return nil\n}",
		`type PetTagsItems0 struct {`,
		"	ID int64 `json:\"id,omitempty\"`",
		"	Name string `json:\"name,omitempty\"`",
		// empty validation
		"func (m *PetTagsItems0) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: map_complex_validations_meta.go
	// NOTE(fredbi): maps are simple types and this definition is no more generated

	// load expectations for model: array_additional_validations_args.go
	flattenRun.AddExpectations("array_additional_validations_args.go", []string{
		`type ArrayAdditionalValidationsArgs struct {`,
		"	P0 *string `json:\"-\"`",
		"	P1 *float64 `json:\"-\"`",
		"	ArrayAdditionalValidationsArgsItems []int64 `json:\"-\"`",
		`func (m *ArrayAdditionalValidationsArgs) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateArrayAdditionalValidationsArgsItems(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ArrayAdditionalValidationsArgs) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	if err := validate.MinLength("0", "body", *m.P0, 3); err != nil {`,
		`	if err := validate.MaxLength("0", "body", *m.P0, 10); err != nil {`,
		"	if err := validate.Pattern(\"0\", \"body\", *m.P0, `\\w+`); err != nil {",
		`func (m *ArrayAdditionalValidationsArgs) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("1", "body", m.P1); err != nil {`,
		`	if err := validate.Minimum("1", "body", *m.P1, 3, false); err != nil {`,
		`	if err := validate.Maximum("1", "body", *m.P1, 12, false); err != nil {`,
		`	if err := validate.MultipleOf("1", "body", *m.P1, 1.5); err != nil {`,
		`func (m *ArrayAdditionalValidationsArgs) validateArrayAdditionalValidationsArgsItems(formats strfmt.Registry) error {`,
		`	for i := range m.ArrayAdditionalValidationsArgsItems {`,
		`		if err := validate.MinimumInt(strconv.Itoa(i+2), "body", m.ArrayAdditionalValidationsArgsItems[i], 3, false); err != nil {`,
		`		if err := validate.MaximumInt(strconv.Itoa(i+2), "body", m.ArrayAdditionalValidationsArgsItems[i], 6, false); err != nil {`,
		`		if err := validate.MultipleOfInt(strconv.Itoa(i+2), "body", m.ArrayAdditionalValidationsArgsItems[i], 1); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: all_of_validations_meta_all_of2.go
	flattenRun.AddExpectations("all_of_validations_meta_all_of2.go", []string{
		`type AllOfValidationsMetaAllOf2 struct {`,
		"	Args []string `json:\"args\"`",
		`func (m *AllOfValidationsMetaAllOf2) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateArgs(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AllOfValidationsMetaAllOf2) validateArgs(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Args) {`,
		`	iArgsSize := int64(len(m.Args)`,
		`	if err := validate.MinItems("args", "body", iArgsSize, 5); err != nil {`,
		`	if err := validate.MaxItems("args", "body", iArgsSize, 20); err != nil {`,
		`	for i := 0; i < len(m.Args); i++ {`,
		`		if err := validate.MinLength("args"+"."+strconv.Itoa(i), "body", m.Args[i], 2); err != nil {`,
		`		if err := validate.MaxLength("args"+"."+strconv.Itoa(i), "body", m.Args[i], 50); err != nil {`,
		"		if err := validate.Pattern(\"args\"+\".\"+strconv.Itoa(i), \"body\", m.Args[i], `[A-Za-z0-9][\\w- ]+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: all_of_validations_meta_all_of0.go
	flattenRun.AddExpectations("all_of_validations_meta_all_of0.go", []string{
		`type AllOfValidationsMetaAllOf0 struct {`,
		"	Name string `json:\"name,omitempty\"`",
		`func (m *AllOfValidationsMetaAllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateName(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AllOfValidationsMetaAllOf0) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 2); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `[A-Za-z0-9][\\w- ]+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_all_of_all_of4.go
	flattenRun.AddExpectations("named_all_of_all_of4.go", []string{
		`type NamedAllOfAllOf4 struct {`,
		"	Opts map[string]float64 `json:\"opts,omitempty\"`",
		`func (m *NamedAllOfAllOf4) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateOpts(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedAllOfAllOf4) validateOpts(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Opts) {`,
		`	for k := range m.Opts {`,
		`		if err := validate.Minimum("opts"+"."+k, "body", m.Opts[k], 2, false); err != nil {`,
		`		if err := validate.Maximum("opts"+"."+k, "body", m.Opts[k], 50, false); err != nil {`,
		`		if err := validate.MultipleOf("opts"+"."+k, "body", m.Opts[k], 1.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_all_of_all_of0.go
	flattenRun.AddExpectations("named_all_of_all_of0.go", []string{
		`type NamedAllOfAllOf0 struct {`,
		"	Name string `json:\"name,omitempty\"`",
		`func (m *NamedAllOfAllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateName(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedAllOfAllOf0) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 2); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `[A-Za-z0-9][\\w- ]+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: nested_map_complex_validations_meta_additional_properties_additional_properties_additional_properties.go
	flattenRun.AddExpectations("nested_map_complex_validations_meta_additional_properties_additional_properties_additional_properties.go", []string{
		`type NestedMapComplexValidationsMetaAdditionalPropertiesAdditionalPropertiesAdditionalProperties struct {`,
		"	Age int32 `json:\"age,omitempty\"`",
		"	Name string `json:\"name,omitempty\"`",
		`func (m *NestedMapComplexValidationsMetaAdditionalPropertiesAdditionalPropertiesAdditionalProperties) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`	if err := m.validateName(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NestedMapComplexValidationsMetaAdditionalPropertiesAdditionalPropertiesAdditionalProperties) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 1, true); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 200, true); err != nil {`,
		`	if err := validate.MultipleOfInt("age", "body", int64(m.Age), 1); err != nil {`,
		`func (m *NestedMapComplexValidationsMetaAdditionalPropertiesAdditionalPropertiesAdditionalProperties) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 10); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `\\w+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: refed_all_of_validations.go
	flattenRun.AddExpectations("refed_all_of_validations.go", []string{
		`type RefedAllOfValidations struct {`,
		`	NamedString`,
		`	NamedNumber`,
		`func (m *RefedAllOfValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.NamedString.Validate(formats); err != nil {`,
		`	if err := m.NamedNumber.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("refed_all_of_validations.go", []string{
		`type RefedAllOfValidations struct {`,
		`	RefedAllOfValidationsAllOf0`,
		`	RefedAllOfValidationsAllOf1`,
		`func (m *RefedAllOfValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.RefedAllOfValidationsAllOf0.Validate(formats); err != nil {`,
		`	if err := m.RefedAllOfValidationsAllOf1.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`type RefedAllOfValidationsAllOf0 string`,
		`func (m RefedAllOfValidationsAllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.MinLength("", "body", string(m), 2); err != nil {`,
		`	if err := validate.MaxLength("", "body", string(m), 50); err != nil {`,
		"	if err := validate.Pattern(\"\", \"body\", string(m), `[A-Za-z0-9][\\w- ]+`); err != nil {",
		`		return errors.CompositeValidationError(res...`,
		`type RefedAllOfValidationsAllOf1 int32`,
		`func (m RefedAllOfValidationsAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.MinimumInt("", "body", int64(m), 0, true); err != nil {`,
		`	if err := validate.MaximumInt("", "body", int64(m), 500, false); err != nil {`,
		`	if err := validate.MultipleOf("", "body", float64(m), 1.5); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: update_org.go
	flattenRun.AddExpectations("update_org.go", []string{
		`type UpdateOrg struct {`,
		"	Email string `json:\"email,omitempty\"`",
		"	InvoiceEmail bool `json:\"invoice_email,omitempty\"`",
		"	TagExpiration *int64 `json:\"tag_expiration,omitempty\"`",
		`func (m *UpdateOrg) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateTagExpiration(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *UpdateOrg) validateTagExpiration(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.TagExpiration) {`,
		`	if err := validate.MinimumInt("tag_expiration", "body", *m.TagExpiration, 0, false); err != nil {`,
		`	if err := validate.MaximumInt("tag_expiration", "body", *m.TagExpiration, 2.592e+06, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("update_org.go", flattenRun.ExpectedFor("UpdateOrg").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: number_validations.go
	flattenRun.AddExpectations("number_validations.go", []string{
		`type NumberValidations struct {`,
		"	Age int32 `json:\"age,omitempty\"`",
		`func (m *NumberValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NumberValidations) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 0, true); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 500, false); err != nil {`,
		`	if err := validate.MultipleOf("age", "body", float64(m.Age), 1.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("number_validations.go", flattenRun.ExpectedFor("NumberValidations").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: nested_map_complex_validations.go
	flattenRun.AddExpectations("nested_map_complex_validations.go", []string{
		`type NestedMapComplexValidations struct {`,
		// maps are now simple types
		"	Meta map[string]map[string]map[string]NestedMapComplexValidationsMetaAdditionalPropertiesAdditionalPropertiesAdditionalProperties `json:\"meta,omitempty\"`",
		`func (m *NestedMapComplexValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMeta(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NestedMapComplexValidations) validateMeta(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Meta) {`,
		`            		for k := range m.Meta {`,
		`            			for kk := range m.Meta[k] {`,
		`            				for kkk := range m.Meta[k][kk] {`,
		`	            				if err := validate.Required("meta"+"."+k+"."+kk+"."+kkk, "body", m.Meta[k][kk][kkk]); err != nil {`,
		`            					if val, ok := m.Meta[k][kk][kkk]; ok {`,
		`            						if err := val.Validate(formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("nested_map_complex_validations.go", []string{
		`type NestedMapComplexValidations struct {`,
		"	Meta map[string]map[string]map[string]NestedMapComplexValidationsMetaAnon `json:\"meta,omitempty\"`",
		`func (m *NestedMapComplexValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMeta(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NestedMapComplexValidations) validateMeta(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Meta) {`,
		`	for k := range m.Meta {`,
		`		for kk := range m.Meta[k] {`,
		`			for kkk := range m.Meta[k][kk] {`,
		`				if val, ok := m.Meta[k][kk][kkk]; ok {`,
		`					if err := val.Validate(formats); err != nil {`,
		`type NestedMapComplexValidationsMetaAnon struct {`,
		"	Age int32 `json:\"age,omitempty\"`",
		"	Name string `json:\"name,omitempty\"`",
		`func (m *NestedMapComplexValidationsMetaAnon) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`	if err := m.validateName(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NestedMapComplexValidationsMetaAnon) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 1, true); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 200, true); err != nil {`,
		`	if err := validate.MultipleOfInt("age", "body", int64(m.Age), 1); err != nil {`,
		`func (m *NestedMapComplexValidationsMetaAnon) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 10); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `\\w+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: category.go
	flattenRun.AddExpectations("category.go", []string{
		`type Category struct {`,
		"	ID int64 `json:\"id,omitempty\"`",
		"	Name string `json:\"name,omitempty\"`",
		// empty validation
		"func (m *Category) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("category.go", flattenRun.ExpectedFor("Category").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: named_all_of_all_of2.go
	flattenRun.AddExpectations("named_all_of_all_of2.go", []string{
		`type NamedAllOfAllOf2 struct {`,
		"	Args []string `json:\"args\"`",
		`func (m *NamedAllOfAllOf2) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateArgs(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedAllOfAllOf2) validateArgs(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Args) {`,
		`	iArgsSize := int64(len(m.Args)`,
		`	if err := validate.MinItems("args", "body", iArgsSize, 5); err != nil {`,
		`	if err := validate.MaxItems("args", "body", iArgsSize, 20); err != nil {`,
		`	for i := 0; i < len(m.Args); i++ {`,
		`		if err := validate.MinLength("args"+"."+strconv.Itoa(i), "body", m.Args[i], 2); err != nil {`,
		`		if err := validate.MaxLength("args"+"."+strconv.Itoa(i), "body", m.Args[i], 50); err != nil {`,
		"		if err := validate.Pattern(\"args\"+\".\"+strconv.Itoa(i), \"body\", m.Args[i], `[A-Za-z0-9][\\w- ]+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_nested_map_complex_additional_properties_additional_properties.go
	// NOTE(fredbi): maps are now simple types - this definition is no more generated

	// load expectations for model: named_nested_map_complex_additional_properties.go
	// NOTE(fredbi): maps are now simple types - this definition is no more generated

	// load expectations for model: named_nested_array.go
	flattenRun.AddExpectations("named_nested_array.go", []string{
		`type NamedNestedArray [][][]string`,
		`func (m NamedNestedArray) Validate(formats strfmt.Registry) error {`,
		`	iNamedNestedArraySize := int64(len(m)`,
		`	if err := validate.MinItems("", "body", iNamedNestedArraySize, 3); err != nil {`,
		`	if err := validate.MaxItems("", "body", iNamedNestedArraySize, 10); err != nil {`,
		`	for i := 0; i < len(m); i++ {`,
		`		iiNamedNestedArraySize := int64(len(m[i])`,
		`		if err := validate.MinItems(strconv.Itoa(i), "body", iiNamedNestedArraySize, 3); err != nil {`,
		`		if err := validate.MaxItems(strconv.Itoa(i), "body", iiNamedNestedArraySize, 10); err != nil {`,
		`		for ii := 0; ii < len(m[i]); ii++ {`,
		`			iiiNamedNestedArraySize := int64(len(m[i][ii])`,
		`			if err := validate.MinItems(strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiNamedNestedArraySize, 3); err != nil {`,
		`			if err := validate.MaxItems(strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiNamedNestedArraySize, 10); err != nil {`,
		`			for iii := 0; iii < len(m[i][ii]); iii++ {`,
		`				if err := validate.MinLength(strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", m[i][ii][iii], 3); err != nil {`,
		`				if err := validate.MaxLength(strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", m[i][ii][iii], 10); err != nil {`,
		"				if err := validate.Pattern(strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii), \"body\", m[i][ii][iii], `\\w+`); err != nil {",
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("named_nested_array.go", flattenRun.ExpectedFor("NamedNestedArray").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: named_all_of.go
	flattenRun.AddExpectations("named_all_of.go", []string{
		`type NamedAllOf struct {`,
		`	NamedAllOfAllOf0`,
		`	NamedAllOfAllOf1`,
		`	NamedAllOfAllOf2`,
		`	NamedAllOfAllOf3`,
		`	NamedAllOfAllOf4`,
		`	NamedAllOfAllOf5`,
		`	NamedAllOfAllOf6`,
		`func (m *NamedAllOf) Validate(formats strfmt.Registry) error {`,
		`	if err := m.NamedAllOfAllOf0.Validate(formats); err != nil {`,
		`	if err := m.NamedAllOfAllOf1.Validate(formats); err != nil {`,
		`	if err := m.NamedAllOfAllOf2.Validate(formats); err != nil {`,
		`	if err := m.NamedAllOfAllOf3.Validate(formats); err != nil {`,
		`	if err := m.NamedAllOfAllOf4.Validate(formats); err != nil {`,
		`	if err := m.NamedAllOfAllOf5.Validate(formats); err != nil {`,
		`	if err := m.NamedAllOfAllOf6.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("named_all_of.go", []string{
		`type NamedAllOf struct {`,
		"	Name string `json:\"name,omitempty\"`",
		"	Age int32 `json:\"age,omitempty\"`",
		"	Args []string `json:\"args\"`",
		"	Assoc [][][]string `json:\"assoc\"`",
		"	Opts map[string]float64 `json:\"opts,omitempty\"`",
		"	ExtOpts map[string]map[string]map[string]int32 `json:\"extOpts,omitempty\"`",
		`	Coords struct {`,
		"		Name string `json:\"name,omitempty\"`",
		"		Age int32 `json:\"age,omitempty\"`",
		"	} `json:\"coords,omitempty\"`",
		`func (m *NamedAllOf) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateName(formats); err != nil {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`	if err := m.validateArgs(formats); err != nil {`,
		`	if err := m.validateAssoc(formats); err != nil {`,
		`	if err := m.validateOpts(formats); err != nil {`,
		`	if err := m.validateExtOpts(formats); err != nil {`,
		`	if err := m.validateCoords(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedAllOf) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 2); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `[A-Za-z0-9][\\w- ]+`); err != nil {",
		`func (m *NamedAllOf) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 2, false); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 50, false); err != nil {`,
		`	if err := validate.MultipleOf("age", "body", float64(m.Age), 1.5); err != nil {`,
		`func (m *NamedAllOf) validateArgs(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Args) {`,
		`	iArgsSize := int64(len(m.Args)`,
		`	if err := validate.MinItems("args", "body", iArgsSize, 5); err != nil {`,
		`	if err := validate.MaxItems("args", "body", iArgsSize, 20); err != nil {`,
		`	for i := 0; i < len(m.Args); i++ {`,
		`		if err := validate.MinLength("args"+"."+strconv.Itoa(i), "body", m.Args[i], 2); err != nil {`,
		`		if err := validate.MaxLength("args"+"."+strconv.Itoa(i), "body", m.Args[i], 50); err != nil {`,
		"		if err := validate.Pattern(\"args\"+\".\"+strconv.Itoa(i), \"body\", m.Args[i], `[A-Za-z0-9][\\w- ]+`); err != nil {",
		`func (m *NamedAllOf) validateAssoc(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Assoc) {`,
		`	iAssocSize := int64(len(m.Assoc)`,
		`	if err := validate.MinItems("assoc", "body", iAssocSize, 5); err != nil {`,
		`	if err := validate.MaxItems("assoc", "body", iAssocSize, 20); err != nil {`,
		`	for i := 0; i < len(m.Assoc); i++ {`,
		`		iiAssocSize := int64(len(m.Assoc[i])`,
		`		if err := validate.MinItems("assoc"+"."+strconv.Itoa(i), "body", iiAssocSize, 5); err != nil {`,
		`		if err := validate.MaxItems("assoc"+"."+strconv.Itoa(i), "body", iiAssocSize, 20); err != nil {`,
		`		for ii := 0; ii < len(m.Assoc[i]); ii++ {`,
		`			iiiAssocSize := int64(len(m.Assoc[i][ii])`,
		`			if err := validate.MinItems("assoc"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiAssocSize, 5); err != nil {`,
		`			if err := validate.MaxItems("assoc"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiAssocSize, 20); err != nil {`,
		`			for iii := 0; iii < len(m.Assoc[i][ii]); iii++ {`,
		`				if err := validate.MinLength("assoc"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", m.Assoc[i][ii][iii], 2); err != nil {`,
		`				if err := validate.MaxLength("assoc"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", m.Assoc[i][ii][iii], 50); err != nil {`,
		"				if err := validate.Pattern(\"assoc\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii), \"body\", m.Assoc[i][ii][iii], `[A-Za-z0-9][\\w- ]+`); err != nil {",
		`func (m *NamedAllOf) validateOpts(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Opts) {`,
		`	for k := range m.Opts {`,
		`		if err := validate.Minimum("opts"+"."+k, "body", m.Opts[k], 2, false); err != nil {`,
		`		if err := validate.Maximum("opts"+"."+k, "body", m.Opts[k], 50, false); err != nil {`,
		`		if err := validate.MultipleOf("opts"+"."+k, "body", m.Opts[k], 1.5); err != nil {`,
		`func (m *NamedAllOf) validateExtOpts(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ExtOpts) {`,
		`	for k := range m.ExtOpts {`,
		`		for kk := range m.ExtOpts[k] {`,
		`			for kkk := range m.ExtOpts[k][kk] {`,
		`				if err := validate.MinimumInt("extOpts"+"."+k+"."+kk+"."+kkk, "body", int64(m.ExtOpts[k][kk][kkk]), 2, false); err != nil {`,
		`				if err := validate.MaximumInt("extOpts"+"."+k+"."+kk+"."+kkk, "body", int64(m.ExtOpts[k][kk][kkk]), 50, false); err != nil {`,
		`				if err := validate.MultipleOf("extOpts"+"."+k+"."+kk+"."+kkk, "body", float64(m.ExtOpts[k][kk][kkk]), 1.5); err != nil {`,
		`func (m *NamedAllOf) validateCoords(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Coords) {`,
		`	if err := validate.MinLength("coords"+"."+"name", "body", m.Coords.Name, 2); err != nil {`,
		`	if err := validate.MaxLength("coords"+"."+"name", "body", m.Coords.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"coords\"+\".\"+\"name\", \"body\", m.Coords.Name, `[A-Za-z0-9][\\w- ]+`); err != nil {",
		`	if err := validate.MinimumInt("coords"+"."+"age", "body", int64(m.Coords.Age), 2, false); err != nil {`,
		`	if err := validate.MaximumInt("coords"+"."+"age", "body", int64(m.Coords.Age), 50, false); err != nil {`,
		`	if err := validate.MultipleOf("coords"+"."+"age", "body", float64(m.Coords.Age), 1.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_map_complex.go
	flattenRun.AddExpectations("named_map_complex.go", []string{
		`type NamedMapComplex map[string]NamedMapComplexAdditionalProperties`,
		`func (m NamedMapComplex) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if err := validate.Required(k, "body", m[k]); err != nil {`,
		`		if val, ok := m[k]; ok {`,
		`			if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("named_map_complex.go", []string{
		`type NamedMapComplex map[string]NamedMapComplexAnon`,
		`func (m NamedMapComplex) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if val, ok := m[k]; ok {`,
		`			if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`type NamedMapComplexAnon struct {`,
		"	Age int32 `json:\"age,omitempty\"`",
		"	Name string `json:\"name,omitempty\"`",
		`func (m *NamedMapComplexAnon) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`	if err := m.validateName(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedMapComplexAnon) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 1, true); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 200, true); err != nil {`,
		`	if err := validate.MultipleOfInt("age", "body", int64(m.Age), 1); err != nil {`,
		`func (m *NamedMapComplexAnon) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 10); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `\\w+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: nested_map_complex_validations_meta.go
	// NOTE(fredbi): maps are now simple types - this definition is no more generated

	// load expectations for model: array_additional_validations.go
	flattenRun.AddExpectations("array_additional_validations.go", []string{
		`type ArrayAdditionalValidations struct {`,
		"	Args ArrayAdditionalValidationsArgs `json:\"args,omitempty\"`",
		`func (m *ArrayAdditionalValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateArgs(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ArrayAdditionalValidations) validateArgs(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Args) {`,
		`	if err := m.Args.Validate(formats); err != nil {`,
		`		if ve, ok := err.(*errors.Validation); ok {`,
		`			return ve.ValidateName("args"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("array_additional_validations.go", []string{
		`type ArrayAdditionalValidations struct {`,
		"	Args *ArrayAdditionalValidationsArgsTuple0 `json:\"args,omitempty\"`",
		`func (m *ArrayAdditionalValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateArgs(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ArrayAdditionalValidations) validateArgs(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Args) {`,
		`	if m.Args != nil {`,
		`		if err := m.Args.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("args"`,
		`type ArrayAdditionalValidationsArgsTuple0 struct {`,
		"	P0 *string `json:\"-\"`",
		"	P1 *float64 `json:\"-\"`",
		"	ArrayAdditionalValidationsArgsTuple0Items []int64 `json:\"-\"`",
		`func (m *ArrayAdditionalValidationsArgsTuple0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateArrayAdditionalValidationsArgsTuple0Items(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ArrayAdditionalValidationsArgsTuple0) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P0", "body", m.P0); err != nil {`,
		`	if err := validate.MinLength("P0", "body", *m.P0, 3); err != nil {`,
		`	if err := validate.MaxLength("P0", "body", *m.P0, 10); err != nil {`,
		"	if err := validate.Pattern(\"P0\", \"body\", *m.P0, `\\w+`); err != nil {",
		`func (m *ArrayAdditionalValidationsArgsTuple0) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P1", "body", m.P1); err != nil {`,
		`	if err := validate.Minimum("P1", "body", *m.P1, 3, false); err != nil {`,
		`	if err := validate.Maximum("P1", "body", *m.P1, 12, false); err != nil {`,
		`	if err := validate.MultipleOf("P1", "body", *m.P1, 1.5); err != nil {`,
		`func (m *ArrayAdditionalValidationsArgsTuple0) validateArrayAdditionalValidationsArgsTuple0Items(formats strfmt.Registry) error {`,
		`	for i := range m.ArrayAdditionalValidationsArgsTuple0Items {`,
		`		if err := validate.MinimumInt(strconv.Itoa(i), "body", m.ArrayAdditionalValidationsArgsTuple0Items[i], 3, false); err != nil {`,
		`		if err := validate.MaximumInt(strconv.Itoa(i), "body", m.ArrayAdditionalValidationsArgsTuple0Items[i], 6, false); err != nil {`,
		`		if err := validate.MultipleOfInt(strconv.Itoa(i), "body", m.ArrayAdditionalValidationsArgsTuple0Items[i], 1); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: nested_object_validations_args_meta.go
	flattenRun.AddExpectations("nested_object_validations_args_meta.go", []string{
		`type NestedObjectValidationsArgsMeta struct {`,
		"	First string `json:\"first,omitempty\"`",
		"	Fourth [][][]float32 `json:\"fourth\"`",
		"	Second float64 `json:\"second,omitempty\"`",
		"	Third []float32 `json:\"third\"`",
		`func (m *NestedObjectValidationsArgsMeta) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateFirst(formats); err != nil {`,
		`	if err := m.validateFourth(formats); err != nil {`,
		`	if err := m.validateSecond(formats); err != nil {`,
		`	if err := m.validateThird(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NestedObjectValidationsArgsMeta) validateFirst(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.First) {`,
		`	if err := validate.MinLength("first", "body", m.First, 2); err != nil {`,
		`	if err := validate.MaxLength("first", "body", m.First, 50); err != nil {`,
		"	if err := validate.Pattern(\"first\", \"body\", m.First, `^\\w+`); err != nil {",
		`func (m *NestedObjectValidationsArgsMeta) validateFourth(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Fourth) {`,
		`	iFourthSize := int64(len(m.Fourth)`,
		`	if err := validate.MinItems("fourth", "body", iFourthSize, 5); err != nil {`,
		`	if err := validate.MaxItems("fourth", "body", iFourthSize, 93); err != nil {`,
		`	for i := 0; i < len(m.Fourth); i++ {`,
		`		iiFourthSize := int64(len(m.Fourth[i])`,
		`		if err := validate.MinItems("fourth"+"."+strconv.Itoa(i), "body", iiFourthSize, 5); err != nil {`,
		`		if err := validate.MaxItems("fourth"+"."+strconv.Itoa(i), "body", iiFourthSize, 93); err != nil {`,
		`		for ii := 0; ii < len(m.Fourth[i]); ii++ {`,
		`			iiiFourthSize := int64(len(m.Fourth[i][ii])`,
		`			if err := validate.MinItems("fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiFourthSize, 5); err != nil {`,
		`			if err := validate.MaxItems("fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiFourthSize, 93); err != nil {`,
		`			for iii := 0; iii < len(m.Fourth[i][ii]); iii++ {`,
		`				if err := validate.Minimum("fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", float64(m.Fourth[i][ii][iii]), 3, false); err != nil {`,
		`				if err := validate.Maximum("fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", float64(m.Fourth[i][ii][iii]), 6, false); err != nil {`,
		`				if err := validate.MultipleOf("fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", float64(m.Fourth[i][ii][iii]), 0.5); err != nil {`,
		`func (m *NestedObjectValidationsArgsMeta) validateSecond(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Second) {`,
		`	if err := validate.Minimum("second", "body", m.Second, 3, false); err != nil {`,
		`	if err := validate.Maximum("second", "body", m.Second, 51, false); err != nil {`,
		`	if err := validate.MultipleOf("second", "body", m.Second, 1.5); err != nil {`,
		`func (m *NestedObjectValidationsArgsMeta) validateThird(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Third) {`,
		`	iThirdSize := int64(len(m.Third)`,
		`	if err := validate.MinItems("third", "body", iThirdSize, 5); err != nil {`,
		`	if err := validate.MaxItems("third", "body", iThirdSize, 93); err != nil {`,
		`	for i := 0; i < len(m.Third); i++ {`,
		`		if err := validate.Minimum("third"+"."+strconv.Itoa(i), "body", float64(m.Third[i]), 3, false); err != nil {`,
		`		if err := validate.Maximum("third"+"."+strconv.Itoa(i), "body", float64(m.Third[i]), 6, false); err != nil {`,
		`		if err := validate.MultipleOf("third"+"."+strconv.Itoa(i), "body", float64(m.Third[i]), 0.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: map_validations.go
	flattenRun.AddExpectations("map_validations.go", []string{
		`type MapValidations struct {`,
		"	Meta map[string]int64 `json:\"meta,omitempty\"`",
		`func (m *MapValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMeta(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *MapValidations) validateMeta(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Meta) {`,
		`	for k := range m.Meta {`,
		`		if err := validate.MinimumInt("meta"+"."+k, "body", m.Meta[k], 3, false); err != nil {`,
		`		if err := validate.MaximumInt("meta"+"."+k, "body", m.Meta[k], 6, false); err != nil {`,
		`		if err := validate.MultipleOfInt("meta"+"."+k, "body", m.Meta[k], 1); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("map_validations.go", flattenRun.ExpectedFor("MapValidations").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: named_all_of_all_of1.go
	flattenRun.AddExpectations("named_all_of_all_of1.go", []string{
		`type NamedAllOfAllOf1 struct {`,
		"	Age int32 `json:\"age,omitempty\"`",
		`func (m *NamedAllOfAllOf1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedAllOfAllOf1) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 2, false); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 50, false); err != nil {`,
		`	if err := validate.MultipleOf("age", "body", float64(m.Age), 1.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: all_of_validations_meta_all_of5.go
	flattenRun.AddExpectations("all_of_validations_meta_all_of5.go", []string{
		`type AllOfValidationsMetaAllOf5 struct {`,
		"	ExtOpts map[string]map[string]map[string]int32 `json:\"extOpts,omitempty\"`",
		`func (m *AllOfValidationsMetaAllOf5) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateExtOpts(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AllOfValidationsMetaAllOf5) validateExtOpts(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ExtOpts) {`,
		`	for k := range m.ExtOpts {`,
		`		for kk := range m.ExtOpts[k] {`,
		`			for kkk := range m.ExtOpts[k][kk] {`,
		`				if err := validate.MinimumInt("extOpts"+"."+k+"."+kk+"."+kkk, "body", int64(m.ExtOpts[k][kk][kkk]), 2, false); err != nil {`,
		`				if err := validate.MaximumInt("extOpts"+"."+k+"."+kk+"."+kkk, "body", int64(m.ExtOpts[k][kk][kkk]), 50, false); err != nil {`,
		`				if err := validate.MultipleOf("extOpts"+"."+k+"."+kk+"."+kkk, "body", float64(m.ExtOpts[k][kk][kkk]), 1.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: all_of_validations.go
	flattenRun.AddExpectations("all_of_validations.go", []string{
		`type AllOfValidations struct {`,
		"	Meta *AllOfValidationsMeta `json:\"meta,omitempty\"`",
		`func (m *AllOfValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMeta(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AllOfValidations) validateMeta(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Meta) {`,
		`	if m.Meta != nil {`,
		`		if err := m.Meta.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("meta"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("all_of_validations.go", []string{
		`type AllOfValidations struct {`,
		`	Meta struct {`,
		"		Name string `json:\"name,omitempty\"`",
		"		Age int32 `json:\"age,omitempty\"`",
		"		Args []string `json:\"args\"`",
		"		Assoc [][][]string `json:\"assoc\"`",
		"		Opts map[string]int32 `json:\"opts,omitempty\"`",
		"		ExtOpts map[string]map[string]map[string]int32 `json:\"extOpts,omitempty\"`",
		`		Coords struct {`,
		"			Name string `json:\"name,omitempty\"`",
		"			Age int32 `json:\"age,omitempty\"`",
		"		} `json:\"coords,omitempty\"`",
		"	} `json:\"meta,omitempty\"`",
		`func (m *AllOfValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMeta(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AllOfValidations) validateMeta(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Meta) {`,
		`	if err := validate.MinLength("meta"+"."+"name", "body", m.Meta.Name, 2); err != nil {`,
		`	if err := validate.MaxLength("meta"+"."+"name", "body", m.Meta.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"meta\"+\".\"+\"name\", \"body\", m.Meta.Name, `[A-Za-z0-9][\\w- ]+`); err != nil {",
		`	if err := validate.MinimumInt("meta"+"."+"age", "body", int64(m.Meta.Age), 2, false); err != nil {`,
		`	if err := validate.MaximumInt("meta"+"."+"age", "body", int64(m.Meta.Age), 50, false); err != nil {`,
		`	if err := validate.MultipleOf("meta"+"."+"age", "body", float64(m.Meta.Age), 1.5); err != nil {`,
		`	iArgsSize := int64(len(m.Meta.Args)`,
		`	if err := validate.MinItems("meta"+"."+"args", "body", iArgsSize, 5); err != nil {`,
		`	if err := validate.MaxItems("meta"+"."+"args", "body", iArgsSize, 20); err != nil {`,
		`	for i := 0; i < len(m.Meta.Args); i++ {`,
		`		if err := validate.MinLength("meta"+"."+"args"+"."+strconv.Itoa(i), "body", m.Meta.Args[i], 2); err != nil {`,
		`		if err := validate.MaxLength("meta"+"."+"args"+"."+strconv.Itoa(i), "body", m.Meta.Args[i], 50); err != nil {`,
		"		if err := validate.Pattern(\"meta\"+\".\"+\"args\"+\".\"+strconv.Itoa(i), \"body\", m.Meta.Args[i], `[A-Za-z0-9][\\w- ]+`); err != nil {",
		`	iAssocSize := int64(len(m.Meta.Assoc)`,
		`	if err := validate.MinItems("meta"+"."+"assoc", "body", iAssocSize, 5); err != nil {`,
		`	if err := validate.MaxItems("meta"+"."+"assoc", "body", iAssocSize, 20); err != nil {`,
		`	for i := 0; i < len(m.Meta.Assoc); i++ {`,
		`		iiAssocSize := int64(len(m.Meta.Assoc[i])`,
		`		if err := validate.MinItems("meta"+"."+"assoc"+"."+strconv.Itoa(i), "body", iiAssocSize, 5); err != nil {`,
		`		if err := validate.MaxItems("meta"+"."+"assoc"+"."+strconv.Itoa(i), "body", iiAssocSize, 20); err != nil {`,
		`		for ii := 0; ii < len(m.Meta.Assoc[i]); ii++ {`,
		`			iiiAssocSize := int64(len(m.Meta.Assoc[i][ii])`,
		`			if err := validate.MinItems("meta"+"."+"assoc"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiAssocSize, 5); err != nil {`,
		`			if err := validate.MaxItems("meta"+"."+"assoc"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiAssocSize, 20); err != nil {`,
		`			for iii := 0; iii < len(m.Meta.Assoc[i][ii]); iii++ {`,
		`				if err := validate.MinLength("meta"+"."+"assoc"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", m.Meta.Assoc[i][ii][iii], 2); err != nil {`,
		`				if err := validate.MaxLength("meta"+"."+"assoc"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", m.Meta.Assoc[i][ii][iii], 50); err != nil {`,
		"				if err := validate.Pattern(\"meta\"+\".\"+\"assoc\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii), \"body\", m.Meta.Assoc[i][ii][iii], `[A-Za-z0-9][\\w- ]+`); err != nil {",
		`	for k := range m.Meta.Opts {`,
		`		if err := validate.MinimumInt("meta"+"."+"opts"+"."+k, "body", int64(m.Meta.Opts[k]), 2, false); err != nil {`,
		`		if err := validate.MaximumInt("meta"+"."+"opts"+"."+k, "body", int64(m.Meta.Opts[k]), 50, false); err != nil {`,
		`		if err := validate.MultipleOf("meta"+"."+"opts"+"."+k, "body", float64(m.Meta.Opts[k]), 1.5); err != nil {`,
		`	for k := range m.Meta.ExtOpts {`,
		`		for kk := range m.Meta.ExtOpts[k] {`,
		`			for kkk := range m.Meta.ExtOpts[k][kk] {`,
		`				if err := validate.MinimumInt("meta"+"."+"extOpts"+"."+k+"."+kk+"."+kkk, "body", int64(m.Meta.ExtOpts[k][kk][kkk]), 2, false); err != nil {`,
		`				if err := validate.MaximumInt("meta"+"."+"extOpts"+"."+k+"."+kk+"."+kkk, "body", int64(m.Meta.ExtOpts[k][kk][kkk]), 50, false); err != nil {`,
		`				if err := validate.MultipleOf("meta"+"."+"extOpts"+"."+k+"."+kk+"."+kkk, "body", float64(m.Meta.ExtOpts[k][kk][kkk]), 1.5); err != nil {`,
		`	if err := validate.MinLength("meta"+"."+"coords"+"."+"name", "body", m.Meta.Coords.Name, 2); err != nil {`,
		`	if err := validate.MaxLength("meta"+"."+"coords"+"."+"name", "body", m.Meta.Coords.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"meta\"+\".\"+\"coords\"+\".\"+\"name\", \"body\", m.Meta.Coords.Name, `[A-Za-z0-9][\\w- ]+`); err != nil {",
		`	if err := validate.MinimumInt("meta"+"."+"coords"+"."+"age", "body", int64(m.Meta.Coords.Age), 2, false); err != nil {`,
		`	if err := validate.MaximumInt("meta"+"."+"coords"+"."+"age", "body", int64(m.Meta.Coords.Age), 50, false); err != nil {`,
		`	if err := validate.MultipleOf("meta"+"."+"coords"+"."+"age", "body", float64(m.Meta.Coords.Age), 1.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_nested_object_meta.go
	flattenRun.AddExpectations("named_nested_object_meta.go", []string{
		`type NamedNestedObjectMeta struct {`,
		"	First string `json:\"first,omitempty\"`",
		"	Fourth [][][]float32 `json:\"fourth\"`",
		"	Second float64 `json:\"second,omitempty\"`",
		"	Third []float32 `json:\"third\"`",
		`func (m *NamedNestedObjectMeta) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateFirst(formats); err != nil {`,
		`	if err := m.validateFourth(formats); err != nil {`,
		`	if err := m.validateSecond(formats); err != nil {`,
		`	if err := m.validateThird(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedNestedObjectMeta) validateFirst(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.First) {`,
		`	if err := validate.MinLength("first", "body", m.First, 2); err != nil {`,
		`	if err := validate.MaxLength("first", "body", m.First, 50); err != nil {`,
		"	if err := validate.Pattern(\"first\", \"body\", m.First, `^\\w+`); err != nil {",
		`func (m *NamedNestedObjectMeta) validateFourth(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Fourth) {`,
		`	iFourthSize := int64(len(m.Fourth)`,
		`	if err := validate.MinItems("fourth", "body", iFourthSize, 5); err != nil {`,
		`	if err := validate.MaxItems("fourth", "body", iFourthSize, 93); err != nil {`,
		`	for i := 0; i < len(m.Fourth); i++ {`,
		`		iiFourthSize := int64(len(m.Fourth[i])`,
		`		if err := validate.MinItems("fourth"+"."+strconv.Itoa(i), "body", iiFourthSize, 5); err != nil {`,
		`		if err := validate.MaxItems("fourth"+"."+strconv.Itoa(i), "body", iiFourthSize, 93); err != nil {`,
		`		for ii := 0; ii < len(m.Fourth[i]); ii++ {`,
		`			iiiFourthSize := int64(len(m.Fourth[i][ii])`,
		`			if err := validate.MinItems("fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiFourthSize, 5); err != nil {`,
		`			if err := validate.MaxItems("fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiFourthSize, 93); err != nil {`,
		`			for iii := 0; iii < len(m.Fourth[i][ii]); iii++ {`,
		`				if err := validate.Minimum("fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", float64(m.Fourth[i][ii][iii]), 3, false); err != nil {`,
		`				if err := validate.Maximum("fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", float64(m.Fourth[i][ii][iii]), 6, false); err != nil {`,
		`				if err := validate.MultipleOf("fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", float64(m.Fourth[i][ii][iii]), 0.5); err != nil {`,
		`func (m *NamedNestedObjectMeta) validateSecond(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Second) {`,
		`	if err := validate.Minimum("second", "body", m.Second, 3, false); err != nil {`,
		`	if err := validate.Maximum("second", "body", m.Second, 51, false); err != nil {`,
		`	if err := validate.MultipleOf("second", "body", m.Second, 1.5); err != nil {`,
		`func (m *NamedNestedObjectMeta) validateThird(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Third) {`,
		`	iThirdSize := int64(len(m.Third)`,
		`	if err := validate.MinItems("third", "body", iThirdSize, 5); err != nil {`,
		`	if err := validate.MaxItems("third", "body", iThirdSize, 93); err != nil {`,
		`	for i := 0; i < len(m.Third); i++ {`,
		`		if err := validate.Minimum("third"+"."+strconv.Itoa(i), "body", float64(m.Third[i]), 3, false); err != nil {`,
		`		if err := validate.Maximum("third"+"."+strconv.Itoa(i), "body", float64(m.Third[i]), 6, false); err != nil {`,
		`		if err := validate.MultipleOf("third"+"."+strconv.Itoa(i), "body", float64(m.Third[i]), 0.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_nested_object.go
	flattenRun.AddExpectations("named_nested_object.go", []string{
		`type NamedNestedObject struct {`,
		"	Meta *NamedNestedObjectMeta `json:\"meta,omitempty\"`",
		`func (m *NamedNestedObject) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMeta(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedNestedObject) validateMeta(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Meta) {`,
		`	if m.Meta != nil {`,
		`		if err := m.Meta.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("meta"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("named_nested_object.go", []string{
		`type NamedNestedObject struct {`,
		"	Meta *NamedNestedObjectMeta `json:\"meta,omitempty\"`",
		`func (m *NamedNestedObject) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMeta(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedNestedObject) validateMeta(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Meta) {`,
		`	if m.Meta != nil {`,
		`		if err := m.Meta.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("meta"`,
		`type NamedNestedObjectMeta struct {`,
		"	First string `json:\"first,omitempty\"`",
		"	Fourth [][][]float32 `json:\"fourth\"`",
		"	Second float64 `json:\"second,omitempty\"`",
		"	Third []float32 `json:\"third\"`",
		`func (m *NamedNestedObjectMeta) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateFirst(formats); err != nil {`,
		`	if err := m.validateFourth(formats); err != nil {`,
		`	if err := m.validateSecond(formats); err != nil {`,
		`	if err := m.validateThird(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedNestedObjectMeta) validateFirst(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.First) {`,
		`	if err := validate.MinLength("meta"+"."+"first", "body", m.First, 2); err != nil {`,
		`	if err := validate.MaxLength("meta"+"."+"first", "body", m.First, 50); err != nil {`,
		"	if err := validate.Pattern(\"meta\"+\".\"+\"first\", \"body\", m.First, `^\\w+`); err != nil {",
		`func (m *NamedNestedObjectMeta) validateFourth(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Fourth) {`,
		`	iFourthSize := int64(len(m.Fourth)`,
		`	if err := validate.MinItems("meta"+"."+"fourth", "body", iFourthSize, 5); err != nil {`,
		`	if err := validate.MaxItems("meta"+"."+"fourth", "body", iFourthSize, 93); err != nil {`,
		`	for i := 0; i < len(m.Fourth); i++ {`,
		`		iiFourthSize := int64(len(m.Fourth[i])`,
		`		if err := validate.MinItems("meta"+"."+"fourth"+"."+strconv.Itoa(i), "body", iiFourthSize, 5); err != nil {`,
		`		if err := validate.MaxItems("meta"+"."+"fourth"+"."+strconv.Itoa(i), "body", iiFourthSize, 93); err != nil {`,
		`		for ii := 0; ii < len(m.Fourth[i]); ii++ {`,
		`			iiiFourthSize := int64(len(m.Fourth[i][ii])`,
		`			if err := validate.MinItems("meta"+"."+"fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiFourthSize, 5); err != nil {`,
		`			if err := validate.MaxItems("meta"+"."+"fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii), "body", iiiFourthSize, 93); err != nil {`,
		`			for iii := 0; iii < len(m.Fourth[i][ii]); iii++ {`,
		`				if err := validate.Minimum("meta"+"."+"fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", float64(m.Fourth[i][ii][iii]), 3, false); err != nil {`,
		`				if err := validate.Maximum("meta"+"."+"fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", float64(m.Fourth[i][ii][iii]), 6, false); err != nil {`,
		`				if err := validate.MultipleOf("meta"+"."+"fourth"+"."+strconv.Itoa(i)+"."+strconv.Itoa(ii)+"."+strconv.Itoa(iii), "body", float64(m.Fourth[i][ii][iii]), 0.5); err != nil {`,
		`func (m *NamedNestedObjectMeta) validateSecond(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Second) {`,
		`	if err := validate.Minimum("meta"+"."+"second", "body", m.Second, 3, false); err != nil {`,
		`	if err := validate.Maximum("meta"+"."+"second", "body", m.Second, 51, false); err != nil {`,
		`	if err := validate.MultipleOf("meta"+"."+"second", "body", m.Second, 1.5); err != nil {`,
		`func (m *NamedNestedObjectMeta) validateThird(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Third) {`,
		`	iThirdSize := int64(len(m.Third)`,
		`	if err := validate.MinItems("meta"+"."+"third", "body", iThirdSize, 5); err != nil {`,
		`	if err := validate.MaxItems("meta"+"."+"third", "body", iThirdSize, 93); err != nil {`,
		`	for i := 0; i < len(m.Third); i++ {`,
		`		if err := validate.Minimum("meta"+"."+"third"+"."+strconv.Itoa(i), "body", float64(m.Third[i]), 3, false); err != nil {`,
		`		if err := validate.Maximum("meta"+"."+"third"+"."+strconv.Itoa(i), "body", float64(m.Third[i]), 6, false); err != nil {`,
		`		if err := validate.MultipleOf("meta"+"."+"third"+"."+strconv.Itoa(i), "body", float64(m.Third[i]), 0.5); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: all_of_validations_meta_all_of6_coords_all_of0.go
	flattenRun.AddExpectations("all_of_validations_meta_all_of6_coords_all_of0.go", []string{
		`type AllOfValidationsMetaAllOf6CoordsAllOf0 struct {`,
		"	Name string `json:\"name,omitempty\"`",
		`func (m *AllOfValidationsMetaAllOf6CoordsAllOf0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateName(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AllOfValidationsMetaAllOf6CoordsAllOf0) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 2); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `[A-Za-z0-9][\\w- ]+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_nested_map.go
	flattenRun.AddExpectations("named_nested_map.go", []string{
		`type NamedNestedMap map[string]map[string]map[string]int64`,
		`func (m NamedNestedMap) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		for kk := range m[k] {`,
		`			for kkk := range m[k][kk] {`,
		`				if err := validate.MinimumInt(k+"."+kk+"."+kkk, "body", m[k][kk][kkk], 3, false); err != nil {`,
		`				if err := validate.MaximumInt(k+"."+kk+"."+kkk, "body", m[k][kk][kkk], 6, false); err != nil {`,
		`				if err := validate.MultipleOfInt(k+"."+kk+"."+kkk, "body", m[k][kk][kkk], 1); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("named_nested_map.go", flattenRun.ExpectedFor("NamedNestedMap").ExpectedLines, todo, noLines, noLines)
}

func initFixtureNestedMaps() {
	// testing fixture-nested-maps.yaml with flatten and expand (--skip-flatten)

	/*
	   Test specifically focused on nested maps (e.g.nested additionalProperties)

	*/

	f := newModelFixture("../fixtures/bugs/1487/fixture-nested-maps.yaml", "Nested maps")
	flattenRun := f.AddRun(false)
	expandRun := f.AddRun(true)

	// load expectations for model: alias_interface.go
	flattenRun.AddExpectations("alias_interface.go", []string{
		`type AliasInterface interface{`,
	},
		// not expected
		validatable,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("alias_interface.go", flattenRun.ExpectedFor("AliasInterface").ExpectedLines, validatable, noLines, noLines)

	// load expectations for model: test_nested_aliased_interface.go
	flattenRun.AddExpectations("test_nested_aliased_interface.go", []string{
		`type TestNestedAliasedInterface struct {`,
		"	Meta map[string]map[string]map[string]AliasInterface `json:\"meta,omitempty\"`",
		// empty validation
		"func (m *TestNestedAliasedInterface) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("test_nested_aliased_interface.go", []string{
		`type TestNestedAliasedInterface struct {`,
		"	Meta map[string]map[string]map[string]interface{} `json:\"meta,omitempty\"`",
		// empty validation
		"func (m *TestNestedAliasedInterface) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: nested_map_validations.go
	flattenRun.AddExpectations("nested_map_validations.go", []string{
		`type NestedMapValidations struct {`,
		"	Meta map[string]map[string]map[string]int64 `json:\"meta,omitempty\"`",
		`func (m *NestedMapValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMeta(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NestedMapValidations) validateMeta(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Meta) {`,
		`	for k := range m.Meta {`,
		`		for kk := range m.Meta[k] {`,
		`			for kkk := range m.Meta[k][kk] {`,
		`				if err := validate.MinimumInt("meta"+"."+k+"."+kk+"."+kkk, "body", m.Meta[k][kk][kkk], 3, false); err != nil {`,
		`				if err := validate.MaximumInt("meta"+"."+k+"."+kk+"."+kkk, "body", m.Meta[k][kk][kkk], 6, false); err != nil {`,
		`				if err := validate.MultipleOfInt("meta"+"."+k+"."+kk+"."+kkk, "body", m.Meta[k][kk][kkk], 1); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("nested_map_validations.go", flattenRun.ExpectedFor("NestedMapValidations").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: named_nested_map_complex.go
	flattenRun.AddExpectations("named_nested_map_complex.go", []string{
		`type NamedNestedMapComplex map[string]map[string]map[string]NamedNestedMapComplexAdditionalPropertiesAdditionalPropertiesAdditionalProperties`,
		`func (m NamedNestedMapComplex) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		for kk := range m[k] {`,
		`			for kkk := range m[k][kk] {`,
		`				if err := validate.Required(k+"."+kk+"."+kkk, "body", m[k][kk][kkk]); err != nil {`,
		`				if val, ok := m[k][kk][kkk]; ok {`,
		`					if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("named_nested_map_complex.go", []string{
		`type NamedNestedMapComplex map[string]map[string]map[string]NamedNestedMapComplexAnon`,
		`func (m NamedNestedMapComplex) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		for kk := range m[k] {`,
		`			for kkk := range m[k][kk] {`,
		`				if val, ok := m[k][kk][kkk]; ok {`,
		`					if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`type NamedNestedMapComplexAnon struct {`,
		"	Age int32 `json:\"age,omitempty\"`",
		"	Name string `json:\"name,omitempty\"`",
		`func (m *NamedNestedMapComplexAnon) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`	if err := m.validateName(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedNestedMapComplexAnon) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 1, true); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 200, true); err != nil {`,
		`	if err := validate.MultipleOfInt("age", "body", int64(m.Age), 1); err != nil {`,
		`func (m *NamedNestedMapComplexAnon) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 10); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `\\w+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: nested_map_complex_validations_meta_additional_properties_additional_properties.go
	// NOTE(fredbi): maps are now simple types - this definition is no more generated

	// load expectations for model: nested_map_no_validations_additional_properties_additional_properties.go
	// NOTE(fredbi): maps are now simple types - this definition is no more generated

	// load expectations for model: test_nested_interface.go
	flattenRun.AddExpectations("test_nested_interface.go", []string{
		`type TestNestedInterface struct {`,
		"	Meta map[string]map[string]map[string]interface{} `json:\"meta,omitempty\"`",
		// empty validation
		"func (m *TestNestedInterface) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("test_nested_interface.go", flattenRun.ExpectedFor("TestNestedInterface").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: named_nested_map_complex_additional_properties_additional_properties_additional_properties.go
	flattenRun.AddExpectations("named_nested_map_complex_additional_properties_additional_properties_additional_properties.go", []string{
		`type NamedNestedMapComplexAdditionalPropertiesAdditionalPropertiesAdditionalProperties struct {`,
		"	Age int32 `json:\"age,omitempty\"`",
		"	Name string `json:\"name,omitempty\"`",
		`func (m *NamedNestedMapComplexAdditionalPropertiesAdditionalPropertiesAdditionalProperties) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`	if err := m.validateName(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NamedNestedMapComplexAdditionalPropertiesAdditionalPropertiesAdditionalProperties) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 1, true); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 200, true); err != nil {`,
		`	if err := validate.MultipleOfInt("age", "body", int64(m.Age), 1); err != nil {`,
		`func (m *NamedNestedMapComplexAdditionalPropertiesAdditionalPropertiesAdditionalProperties) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 10); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `\\w+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: nested_map_no_validations_additional_properties.go
	// NOTE(fredbi): maps are simple types and this definition is no more generated

	expandRun.AddExpectations("nested_map_no_validations.go", []string{
		`type NestedMapNoValidations map[string]map[string]map[string]NestedMapNoValidationsAnon`,
		`func (m NestedMapNoValidations) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		for kk := range m[k] {`,
		`			for kkk := range m[k][kk] {`,
		`				if val, ok := m[k][kk][kkk]; ok {`,
		`					if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`type NestedMapNoValidationsAnon struct {`,
		"	Age int64 `json:\"age,omitempty\"`",
		"	Name string `json:\"name,omitempty\"`",
		// empty validation
		"func (m *NestedMapNoValidationsAnon) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: nested_map_complex_validations_meta_additional_properties.go
	// NOTE(fredbi): maps are simple types and this definition is no more generated

	// load expectations for model: nested_map_no_validations_additional_properties_additional_properties_additional_properties.go
	// NOTE: maps are now simple types - this definition is no more generated

	// load expectations for model: nested_map_no_validations.go
	flattenRun.AddExpectations("nested_map_no_validations.go", []string{
		`type NestedMapNoValidations map[string]map[string]map[string]NestedMapNoValidationsAdditionalPropertiesAdditionalPropertiesAdditionalProperties`,
		`func (m NestedMapNoValidations) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`       	for kk := range m[k] {`,
		`            		for kkk := range m[k][kk] {`,
		`            			if val, ok := m[k][kk][kkk]; ok {`,
		`            			if err := validate.Required(k+"."+kk+"."+kkk, "body", m[k][kk][kkk]); err != nil {`,
		`            				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: nested_map_complex_validations_meta_additional_properties_additional_properties_additional_properties.go
	flattenRun.AddExpectations("nested_map_complex_validations_meta_additional_properties_additional_properties_additional_properties.go", []string{
		`type NestedMapComplexValidationsMetaAdditionalPropertiesAdditionalPropertiesAdditionalProperties struct {`,
		"	Age int32 `json:\"age,omitempty\"`",
		"	Name string `json:\"name,omitempty\"`",
		`func (m *NestedMapComplexValidationsMetaAdditionalPropertiesAdditionalPropertiesAdditionalProperties) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateAge(formats); err != nil {`,
		`	if err := m.validateName(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NestedMapComplexValidationsMetaAdditionalPropertiesAdditionalPropertiesAdditionalProperties) validateAge(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Age) {`,
		`	if err := validate.MinimumInt("age", "body", int64(m.Age), 1, true); err != nil {`,
		`	if err := validate.MaximumInt("age", "body", int64(m.Age), 200, true); err != nil {`,
		`	if err := validate.MultipleOfInt("age", "body", int64(m.Age), 1); err != nil {`,
		`func (m *NestedMapComplexValidationsMetaAdditionalPropertiesAdditionalPropertiesAdditionalProperties) validateName(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Name) {`,
		`	if err := validate.MinLength("name", "body", m.Name, 10); err != nil {`,
		`	if err := validate.MaxLength("name", "body", m.Name, 50); err != nil {`,
		"	if err := validate.Pattern(\"name\", \"body\", m.Name, `\\w+`); err != nil {",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: nested_map_complex_validations.go
	flattenRun.AddExpectations("nested_map_complex_validations.go", []string{
		`type NestedMapComplexValidations struct {`,
		// maps are now simple types
		"	Meta map[string]map[string]map[string]NestedMapComplexValidationsMetaAdditionalPropertiesAdditionalPropertiesAdditionalProperties `json:\"meta,omitempty\"`",
		`func (m *NestedMapComplexValidations) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateMeta(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *NestedMapComplexValidations) validateMeta(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Meta) {`,
		`          		for k := range m.Meta {`,
		`          			for kk := range m.Meta[k] {`,
		`          				for kkk := range m.Meta[k][kk] {`,
		`          				if err := validate.Required("meta"+"."+k+"."+kk+"."+kkk, "body", m.Meta[k][kk][kkk]); err != nil {`,
		`          					if val, ok := m.Meta[k][kk][kkk]; ok {`,
		`          						if err := val.Validate(formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: named_nested_map_complex_additional_properties_additional_properties.go
	// NOTE(fredbi): maps are now simple types - this definition is no more generated

	// load expectations for model: named_nested_map_complex_additional_properties.go
	// NOTE(fredbi): maps are now simple types - this definition is no more generated

	// load expectations for model: nested_map_complex_validations_meta.go
	// NOTE(fredbi): maps are now simple types - this definition is no more generated
}

func initFixture844Variations() {
	// testing fixture-844-variations.yaml with flatten and expand (--skip-flatten)

	/*
	   repro
	*/

	f := newModelFixture("../fixtures/bugs/1487/fixture-844-variations.yaml", "allOf bugs with empty objects")
	flattenRun := f.AddRun(false)
	expandRun := f.AddRun(true)

	// load expectations for model: foo.go
	flattenRun.AddExpectations("foo.go", []string{
		`type Foo interface{`,
	},
		// not expected
		validatable,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("foo.go", flattenRun.ExpectedFor("Foo").ExpectedLines, validatable, noLines, noLines)

	// load expectations for model: variation2.go
	flattenRun.AddExpectations("variation2.go", []string{
		`type Variation2 struct {`,
		"	Prop1 EmptyEnum `json:\"prop1,omitempty\"`",
		// empty validation
		"func (m *Variation2) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("variation2.go", []string{
		`type Variation2 struct {`,
		"	Prop1 interface{} `json:\"prop1,omitempty\"`",
		// empty validation
		"func (m *Variation2) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: bar.go
	flattenRun.AddExpectations("bar.go", []string{
		`type Bar interface{`,
	},
		// not expected
		validatable,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("bar.go", flattenRun.ExpectedFor("Bar").ExpectedLines, validatable, noLines, noLines)

	// load expectations for model: variation3.go
	flattenRun.AddExpectations("variation3.go", []string{
		`type Variation3 struct {`,
		"	Prop1 []EmptyEnum `json:\"prop1\"`",
		`func (m *Variation3) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *Variation3) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	iProp1Size := int64(len(m.Prop1)`,
		`	if err := validate.MinItems("prop1", "body", iProp1Size, 10); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("variation3.go", []string{
		`type Variation3 struct {`,
		"	Prop1 []interface{} `json:\"prop1\"`",
		`func (m *Variation3) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`var variation3Prop1ItemsEnum []interface{`,
		`	var res []interface{`,
		"	if err := json.Unmarshal([]byte(`[\"abc\",\"def\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		variation3Prop1ItemsEnum = append(variation3Prop1ItemsEnum, v`,
		`func (m *Variation3) validateProp1ItemsEnum(path, location string, value interface{}) error {`,
		`	if err := validate.EnumCase(path, location, value, variation3Prop1ItemsEnum, true); err != nil {`,
		`func (m *Variation3) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	iProp1Size := int64(len(m.Prop1)`,
		`	if err := validate.MinItems("prop1", "body", iProp1Size, 10); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: tuple_variation.go
	flattenRun.AddExpectations("tuple_variation.go", []string{
		`type TupleVariation struct {`,
		"	P0 *int64 `json:\"-\"`",
		"	P1 Bar `json:\"-\"`",
		"	P2 NonInterface `json:\"-\"`",
		"	P3 []Bar `json:\"-\"`",
		"	TupleVariationItems []interface{} `json:\"-\"`",
		`func (m *TupleVariation) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateP2(formats); err != nil {`,
		`	if err := m.validateP3(formats); err != nil {`,
		`	if err := m.validateTupleVariationItems(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *TupleVariation) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	if err := validate.MaximumInt("0", "body", *m.P0, 10, false); err != nil {`,
		`func (m *TupleVariation) validateP1(formats strfmt.Registry) error {`,
		`if m.P1 == nil {`, // now required interface{} checked against nil
		`return errors.Required("1", "body", m.P1)`,
		`func (m *TupleVariation) validateP2(formats strfmt.Registry) error {`,
		`	if err := m.P2.Validate(formats); err != nil {`,
		`		if ve, ok := err.(*errors.Validation); ok {`,
		`			return ve.ValidateName("2"`,
		`func (m *TupleVariation) validateP3(formats strfmt.Registry) error {`,
		`	if err := validate.Required("3", "body", m.P3); err != nil {`,
		`	iP3Size := int64(len(m.P3)`,
		`	if err := validate.MaxItems("3", "body", iP3Size, 10); err != nil {`,
		// empty validation
		"func (m *TupleVariation) validateTupleVariationItems(formats strfmt.Registry) error {\n\n	return nil\n}",
	},
		// not expected
		append(todo,
			`	if err := validate.Required("1", "body", m.P1)`, // check we don't have redundant validations
			`	if err := validate.Required("1", "body", Bar(m.P1))`,
		),
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("tuple_variation.go", []string{
		`type TupleVariation struct {`,
		"	P0 *int64 `json:\"-\"`",
		"	P1 interface{} `json:\"-\"`",
		"	P2 map[string]strfmt.Date `json:\"-\"`",
		"	P3 []interface{} `json:\"-\"`",
		"	TupleVariationItems []interface{} `json:\"-\"`",
		`func (m *TupleVariation) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateP2(formats); err != nil {`,
		`	if err := m.validateP3(formats); err != nil {`,
		`	if err := m.validateTupleVariationItems(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *TupleVariation) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	if err := validate.MaximumInt("0", "body", *m.P0, 10, false); err != nil {`,
		`func (m *TupleVariation) validateP1(formats strfmt.Registry) error {`,
		`if m.P1 == nil {`, // now required interface{} checked against nil
		`return errors.Required("1", "body", m.P1)`,
		`func (m *TupleVariation) validateP2(formats strfmt.Registry) error {`,
		`	for k := range m.P2 {`,
		`		if err := validate.FormatOf("2"+"."+k, "body", "date", m.P2[k].String(), formats); err != nil {`,
		`func (m *TupleVariation) validateP3(formats strfmt.Registry) error {`,
		`	if err := validate.Required("3", "body", m.P3); err != nil {`,
		`	iP3Size := int64(len(m.P3)`,
		`	if err := validate.MaxItems("3", "body", iP3Size, 10); err != nil {`,
		// empty validation
		"func (m *TupleVariation) validateTupleVariationItems(formats strfmt.Registry) error {\n\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: add_items_variation.go
	flattenRun.AddExpectations("add_items_variation.go", []string{
		`type AddItemsVariation struct {`,
		"	P0 *int64 `json:\"-\"`",
		"	P1 Bar `json:\"-\"`",
		"	AddItemsVariationItems [][]Foo `json:\"-\"`",
		`func (m *AddItemsVariation) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateAddItemsVariationItems(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AddItemsVariation) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	if err := validate.MaximumInt("0", "body", *m.P0, 10, false); err != nil {`,
		`func (m *AddItemsVariation) validateP1(formats strfmt.Registry) error {`,
		`if m.P1 == nil {`, // now required interface{} checked against nil
		`return errors.Required("1", "body", m.P1)`,
		`func (m *AddItemsVariation) validateAddItemsVariationItems(formats strfmt.Registry) error {`,
		`	for i := range m.AddItemsVariationItems {`,
		`		if err := validate.UniqueItems(strconv.Itoa(i+2), "body", m.AddItemsVariationItems[i]); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("add_items_variation.go", []string{
		`type AddItemsVariation struct {`,
		"	P0 *int64 `json:\"-\"`",
		"	P1 interface{} `json:\"-\"`",
		"	AddItemsVariationItems [][]interface{} `json:\"-\"`",
		`func (m *AddItemsVariation) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateAddItemsVariationItems(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AddItemsVariation) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	if err := validate.MaximumInt("0", "body", *m.P0, 10, false); err != nil {`,
		`func (m *AddItemsVariation) validateP1(formats strfmt.Registry) error {`,
		`if m.P1 == nil {`, // now required interface{} checked against nil
		`return errors.Required("1", "body", m.P1)`,
		`func (m *AddItemsVariation) validateAddItemsVariationItems(formats strfmt.Registry) error {`,
		`	for i := range m.AddItemsVariationItems {`,
		`		if err := validate.UniqueItems(strconv.Itoa(i+2), "body", m.AddItemsVariationItems[i]); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: non_interface.go
	flattenRun.AddExpectations("non_interface.go", []string{
		`type NonInterface map[string]strfmt.Date`,
		`func (m NonInterface) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if err := validate.FormatOf(k, "body", "date", m[k].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("non_interface.go", flattenRun.ExpectedFor("NonInterface").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: variation0.go
	flattenRun.AddExpectations("variation0.go", []string{
		`type Variation0 struct {`,
		`	Foo`,
		`	Bar`,
		// empty validation
		"func (m *Variation0) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("variation0.go", []string{
		`type Variation0 struct {`,
		`	Variation0AllOf0`,
		`	Variation0AllOf1`,
		`type Variation0AllOf0 interface{}`,
		`type Variation0AllOf1 interface{}`,
		// empty validation
		"func (m *Variation0) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: variation1.go
	flattenRun.AddExpectations("variation1.go", []string{
		`type Variation1 struct {`,
		`	Foo`,
		`	NonInterface`,
		`func (m *Variation1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.NonInterface.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("variation1.go", []string{
		`type Variation1 struct {`,
		`	Variation1AllOf0`,
		"	AO1 map[string]strfmt.Date `json:\"-\"`",
		`func (m *Variation1) Validate(formats strfmt.Registry) error {`,
		`	for k := range m.AO1 {`,
		`		if err := validate.FormatOf(k, "body", "date", m.AO1[k].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`type Variation1AllOf0 interface{}`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: empty_enum.go
	flattenRun.AddExpectations("empty_enum.go", []string{
		`type EmptyEnum interface{}`,
	},
		// not expected
		validatable,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("empty_enum.go", flattenRun.ExpectedFor("EmptyEnum").ExpectedLines, validatable, noLines, noLines)

	// load expectations for model: get_o_k_body.go
	flattenRun.AddExpectations("get_o_k_body.go", []string{
		`type GetOKBody struct {`,
		`	Foo`,
		`	Bar`,
		// empty validation
		"func (m *GetOKBody) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

}

func initFixtureMoreAddProps() {
	// testing fixture-moreAddProps.yaml with flatten and expand (--skip-flatten)

	/*
	   various patterns of additionalProperties
	*/

	f := newModelFixture("../fixtures/bugs/1487/fixture-moreAddProps.yaml", "fixture for additionalProperties")
	flattenRun := f.AddRun(false)
	expandRun := f.AddRun(true)

	// load expectations for model: trial.go
	flattenRun.AddExpectations("trial.go", []string{
		`type Trial struct {`,
		"	A1 strfmt.DateTime `json:\"a1,omitempty\"`",
		"	AdditionalProperties *TrialAdditionalProperties `json:\"additionalProperties,omitempty\"`",
		`func (m *Trial) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateA1(formats); err != nil {`,
		`	if err := m.validateAdditionalProperties(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *Trial) validateA1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.A1) {`,
		`	if err := validate.FormatOf("a1", "body", "date-time", m.A1.String(), formats); err != nil {`,
		`func (m *Trial) validateAdditionalProperties(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.AdditionalProperties) {`,
		`	if m.AdditionalProperties != nil {`,
		`		if err := m.AdditionalProperties.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("additionalProperties"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: trial_additional_properties.go
	flattenRun.AddExpectations("trial_additional_properties.go", []string{
		`type TrialAdditionalProperties struct {`,
		"	Discourse string `json:\"discourse,omitempty\"`",
		"	HoursSpent float64 `json:\"hoursSpent,omitempty\"`",
		"	TrialAdditionalPropertiesAdditionalProperties map[string]interface{} `json:\"-\"`",
		// empty validation
		"func (m *TrialAdditionalProperties) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_transitive_refed_object_thing_additional_properties.go
	flattenRun.AddExpectations("additional_transitive_refed_object_thing_additional_properties.go", []string{
		`type AdditionalTransitiveRefedObjectThingAdditionalProperties struct {`,
		"	Prop1 *AdditionalTransitiveRefedObjectThingAdditionalPropertiesProp1 `json:\"prop1,omitempty\"`",
		"	AdditionalTransitiveRefedObjectThingAdditionalProperties map[string]*AdditionalTransitiveRefedObjectThingAdditionalPropertiesAdditionalProperties `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedObjectThingAdditionalProperties) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`	for k := range m.AdditionalTransitiveRefedObjectThingAdditionalProperties {`,
		`		if val, ok := m.AdditionalTransitiveRefedObjectThingAdditionalProperties[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedObjectThingAdditionalProperties) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	if m.Prop1 != nil {`,
		`		if err := m.Prop1.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("prop1"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_transitive_refed_object_thing_additional_properties_prop1.go
	flattenRun.AddExpectations("additional_transitive_refed_object_thing_additional_properties_prop1.go", []string{
		`type AdditionalTransitiveRefedObjectThingAdditionalPropertiesProp1 struct {`,
		"	ThisOneNotRequiredEither int64 `json:\"thisOneNotRequiredEither,omitempty\"`",
		"	AdditionalTransitiveRefedObjectThingAdditionalPropertiesProp1 map[string]*AdditionalTransitiveRefedObjectThingAdditionalPropertiesProp1AdditionalProperties `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedObjectThingAdditionalPropertiesProp1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequiredEither(formats); err != nil {`,
		`	for k := range m.AdditionalTransitiveRefedObjectThingAdditionalPropertiesProp1 {`,
		`		if val, ok := m.AdditionalTransitiveRefedObjectThingAdditionalPropertiesProp1[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedObjectThingAdditionalPropertiesProp1) validateThisOneNotRequiredEither(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequiredEither) {`,
		`	if err := validate.MaximumInt("thisOneNotRequiredEither", "body", m.ThisOneNotRequiredEither, 20, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_transitive_refed_thing.go
	flattenRun.AddExpectations("additional_transitive_refed_thing.go", []string{
		`type AdditionalTransitiveRefedThing struct {`,
		"	ThisOneNotRequired int64 `json:\"thisOneNotRequired,omitempty\"`",
		"	AdditionalTransitiveRefedThing map[string][]*AdditionalTransitiveRefedThingAdditionalPropertiesItems `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequired(formats); err != nil {`,
		`	for k := range m.AdditionalTransitiveRefedThing {`,
		`		if err := validate.Required(k, "body", m.AdditionalTransitiveRefedThing[k]); err != nil {`,
		`		if err := validate.UniqueItems(k, "body", m.AdditionalTransitiveRefedThing[k]); err != nil {`,
		`		for i := 0; i < len(m.AdditionalTransitiveRefedThing[k]); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`			if swag.IsZero(m.AdditionalTransitiveRefedThing[k][i]) {`,
		// nullable not required:
		`			if m.AdditionalTransitiveRefedThing[k][i] != nil {`,
		`				if err := m.AdditionalTransitiveRefedThing[k][i].Validate(formats); err != nil {`,
		`					if ve, ok := err.(*errors.Validation); ok {`,
		`						return ve.ValidateName(k + "." + strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedThing) validateThisOneNotRequired(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequired) {`,
		`	if err := validate.MaximumInt("thisOneNotRequired", "body", m.ThisOneNotRequired, 10, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_transitive_refed_thing_additional_properties_items_additional_properties_additional_properties.go
	flattenRun.AddExpectations("additional_transitive_refed_thing_additional_properties_items_additional_properties_additional_properties.go", []string{
		`type AdditionalTransitiveRefedThingAdditionalPropertiesItemsAdditionalPropertiesAdditionalProperties struct {`,
		"	Discourse string `json:\"discourse,omitempty\"`",
		"	HoursSpent float64 `json:\"hoursSpent,omitempty\"`",
		"	AdditionalTransitiveRefedThingAdditionalPropertiesItemsAdditionalPropertiesAdditionalPropertiesAdditionalProperties map[string]interface{} `json:\"-\"`",
		// empty validation
		"func (m *AdditionalTransitiveRefedThingAdditionalPropertiesItemsAdditionalPropertiesAdditionalProperties) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_transitive_refed_object_thing.go
	flattenRun.AddExpectations("additional_transitive_refed_object_thing.go", []string{
		`type AdditionalTransitiveRefedObjectThing struct {`,
		"	ThisOneNotRequired int64 `json:\"thisOneNotRequired,omitempty\"`",
		"	AdditionalTransitiveRefedObjectThing map[string]*AdditionalTransitiveRefedObjectThingAdditionalProperties `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedObjectThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequired(formats); err != nil {`,
		`	for k := range m.AdditionalTransitiveRefedObjectThing {`,
		`		if val, ok := m.AdditionalTransitiveRefedObjectThing[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedObjectThing) validateThisOneNotRequired(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequired) {`,
		`	if err := validate.MaximumInt("thisOneNotRequired", "body", m.ThisOneNotRequired, 10, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_transitive_refed_object_thing_additional_properties_additional_properties.go
	flattenRun.AddExpectations("additional_transitive_refed_object_thing_additional_properties_additional_properties.go", []string{
		`type AdditionalTransitiveRefedObjectThingAdditionalPropertiesAdditionalProperties struct {`,
		"	Discourse string `json:\"discourse,omitempty\"`",
		"	HoursSpent float64 `json:\"hoursSpent,omitempty\"`",
		"	AdditionalTransitiveRefedObjectThingAdditionalPropertiesAdditionalPropertiesAdditionalProperties map[string]interface{} `json:\"-\"`",
		// empty validation
		"func (m *AdditionalTransitiveRefedObjectThingAdditionalPropertiesAdditionalProperties) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_transitive_refed_object_thing_additional_properties_prop1_additional_properties.go
	flattenRun.AddExpectations("additional_transitive_refed_object_thing_additional_properties_prop1_additional_properties.go", []string{
		`type AdditionalTransitiveRefedObjectThingAdditionalPropertiesProp1AdditionalProperties struct {`,
		"	A1 strfmt.DateTime `json:\"a1,omitempty\"`",
		"	B1 strfmt.Date `json:\"b1,omitempty\"`",
		"	AdditionalTransitiveRefedObjectThingAdditionalPropertiesProp1AdditionalPropertiesAdditionalProperties map[string]interface{} `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedObjectThingAdditionalPropertiesProp1AdditionalProperties) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateA1(formats); err != nil {`,
		`	if err := m.validateB1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedObjectThingAdditionalPropertiesProp1AdditionalProperties) validateA1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.A1) {`,
		`	if err := validate.FormatOf("a1", "body", "date-time", m.A1.String(), formats); err != nil {`,
		`func (m *AdditionalTransitiveRefedObjectThingAdditionalPropertiesProp1AdditionalProperties) validateB1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.B1) {`,
		`	if err := validate.FormatOf("b1", "body", "date", m.B1.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_transitive_refed_thing_additional_properties_items.go
	flattenRun.AddExpectations("additional_transitive_refed_thing_additional_properties_items.go", []string{
		`type AdditionalTransitiveRefedThingAdditionalPropertiesItems struct {`,
		"	ThisOneNotRequiredEither int64 `json:\"thisOneNotRequiredEither,omitempty\"`",
		"	AdditionalTransitiveRefedThingAdditionalPropertiesItems map[string]*AdditionalTransitiveRefedThingAdditionalPropertiesItemsAdditionalProperties `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedThingAdditionalPropertiesItems) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequiredEither(formats); err != nil {`,
		`	for k := range m.AdditionalTransitiveRefedThingAdditionalPropertiesItems {`,
		`		if val, ok := m.AdditionalTransitiveRefedThingAdditionalPropertiesItems[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedThingAdditionalPropertiesItems) validateThisOneNotRequiredEither(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequiredEither) {`,
		`	if err := validate.MaximumInt("thisOneNotRequiredEither", "body", m.ThisOneNotRequiredEither, 20, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_transitive_refed_thing_additional_properties_items_additional_properties.go
	flattenRun.AddExpectations("additional_transitive_refed_thing_additional_properties_items_additional_properties.go", []string{
		`type AdditionalTransitiveRefedThingAdditionalPropertiesItemsAdditionalProperties struct {`,
		"	A1 strfmt.DateTime `json:\"a1,omitempty\"`",
		"	B1 strfmt.DateTime `json:\"b1,omitempty\"`",
		"	AdditionalTransitiveRefedThingAdditionalPropertiesItemsAdditionalProperties map[string]*AdditionalTransitiveRefedThingAdditionalPropertiesItemsAdditionalPropertiesAdditionalProperties `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedThingAdditionalPropertiesItemsAdditionalProperties) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateA1(formats); err != nil {`,
		`	if err := m.validateB1(formats); err != nil {`,
		`	for k := range m.AdditionalTransitiveRefedThingAdditionalPropertiesItemsAdditionalProperties {`,
		`		if val, ok := m.AdditionalTransitiveRefedThingAdditionalPropertiesItemsAdditionalProperties[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedThingAdditionalPropertiesItemsAdditionalProperties) validateA1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.A1) {`,
		`	if err := validate.FormatOf("a1", "body", "date-time", m.A1.String(), formats); err != nil {`,
		`func (m *AdditionalTransitiveRefedThingAdditionalPropertiesItemsAdditionalProperties) validateB1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.B1) {`,
		`	if err := validate.FormatOf("b1", "body", "date-time", m.B1.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: trial.go
	expandRun.AddExpectations("trial.go", []string{
		`type Trial struct {`,
		"	A1 strfmt.DateTime `json:\"a1,omitempty\"`",
		"	AdditionalProperties *TrialAdditionalProperties `json:\"additionalProperties,omitempty\"`",
		`func (m *Trial) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateA1(formats); err != nil {`,
		`	if err := m.validateAdditionalProperties(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *Trial) validateA1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.A1) {`,
		`	if err := validate.FormatOf("a1", "body", "date-time", m.A1.String(), formats); err != nil {`,
		`func (m *Trial) validateAdditionalProperties(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.AdditionalProperties) {`,
		`	if m.AdditionalProperties != nil {`,
		`		if err := m.AdditionalProperties.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("additionalProperties"`,
		`type TrialAdditionalProperties struct {`,
		"	Discourse string `json:\"discourse,omitempty\"`",
		"	HoursSpent float64 `json:\"hoursSpent,omitempty\"`",
		"	TrialAdditionalPropertiesAdditionalProperties map[string]interface{} `json:\"-\"`",
		// empty validation
		"func (m *TrialAdditionalProperties) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_transitive_refed_thing.go
	expandRun.AddExpectations("additional_transitive_refed_thing.go", []string{
		`type AdditionalTransitiveRefedThing struct {`,
		"	ThisOneNotRequired int64 `json:\"thisOneNotRequired,omitempty\"`",
		"	AdditionalTransitiveRefedThing map[string][]*AdditionalTransitiveRefedThingItems0 `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequired(formats); err != nil {`,
		`	for k := range m.AdditionalTransitiveRefedThing {`,
		`		if err := validate.Required(k, "body", m.AdditionalTransitiveRefedThing[k]); err != nil {`,
		`		if err := validate.UniqueItems(k, "body", m.AdditionalTransitiveRefedThing[k]); err != nil {`,
		`		for i := 0; i < len(m.AdditionalTransitiveRefedThing[k]); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`			if swag.IsZero(m.AdditionalTransitiveRefedThing[k][i]) {`,
		// nullable required:
		`			if m.AdditionalTransitiveRefedThing[k][i] != nil {`,
		`				if err := m.AdditionalTransitiveRefedThing[k][i].Validate(formats); err != nil {`,
		`					if ve, ok := err.(*errors.Validation); ok {`,
		`						return ve.ValidateName(k + "." + strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedThing) validateThisOneNotRequired(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequired) {`,
		`	if err := validate.MaximumInt("thisOneNotRequired", "body", m.ThisOneNotRequired, 10, false); err != nil {`,
		`type AdditionalTransitiveRefedThingItems0 struct {`,
		"	ThisOneNotRequiredEither int64 `json:\"thisOneNotRequiredEither,omitempty\"`",
		"	AdditionalTransitiveRefedThingItems0 map[string]*AdditionalTransitiveRefedThingItems0Anon `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedThingItems0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequiredEither(formats); err != nil {`,
		`	for k := range m.AdditionalTransitiveRefedThingItems0 {`,
		`		if val, ok := m.AdditionalTransitiveRefedThingItems0[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedThingItems0) validateThisOneNotRequiredEither(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequiredEither) {`,
		`	if err := validate.MaximumInt("thisOneNotRequiredEither", "body", m.ThisOneNotRequiredEither, 20, false); err != nil {`,
		`type AdditionalTransitiveRefedThingItems0Anon struct {`,
		"	A1 strfmt.DateTime `json:\"a1,omitempty\"`",
		"	B1 strfmt.DateTime `json:\"b1,omitempty\"`",
		"	AdditionalTransitiveRefedThingItems0Anon map[string]*AdditionalTransitiveRefedThingItems0AnonAnon `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedThingItems0Anon) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateA1(formats); err != nil {`,
		`	if err := m.validateB1(formats); err != nil {`,
		`	for k := range m.AdditionalTransitiveRefedThingItems0Anon {`,
		`		if val, ok := m.AdditionalTransitiveRefedThingItems0Anon[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedThingItems0Anon) validateA1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.A1) {`,
		`	if err := validate.FormatOf("a1", "body", "date-time", m.A1.String(), formats); err != nil {`,
		`func (m *AdditionalTransitiveRefedThingItems0Anon) validateB1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.B1) {`,
		`	if err := validate.FormatOf("b1", "body", "date-time", m.B1.String(), formats); err != nil {`,
		`type AdditionalTransitiveRefedThingItems0AnonAnon struct {`,
		"	Discourse string `json:\"discourse,omitempty\"`",
		"	HoursSpent float64 `json:\"hoursSpent,omitempty\"`",
		"	AdditionalTransitiveRefedThingItems0AnonAnonAdditionalProperties map[string]interface{} `json:\"-\"`",
		// empty validation
		"func (m *AdditionalTransitiveRefedThingItems0AnonAnon) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: additional_transitive_refed_object_thing.go
	expandRun.AddExpectations("additional_transitive_refed_object_thing.go", []string{
		`type AdditionalTransitiveRefedObjectThing struct {`,
		"	ThisOneNotRequired int64 `json:\"thisOneNotRequired,omitempty\"`",
		"	AdditionalTransitiveRefedObjectThing map[string]*AdditionalTransitiveRefedObjectThingAnon `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedObjectThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequired(formats); err != nil {`,
		`	for k := range m.AdditionalTransitiveRefedObjectThing {`,
		`		if val, ok := m.AdditionalTransitiveRefedObjectThing[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedObjectThing) validateThisOneNotRequired(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequired) {`,
		`	if err := validate.MaximumInt("thisOneNotRequired", "body", m.ThisOneNotRequired, 10, false); err != nil {`,
		`type AdditionalTransitiveRefedObjectThingAnon struct {`,
		"	Prop1 *AdditionalTransitiveRefedObjectThingAnonProp1 `json:\"prop1,omitempty\"`",
		"	AdditionalTransitiveRefedObjectThingAnon map[string]*AdditionalTransitiveRefedObjectThingAnonAnon `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedObjectThingAnon) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`	for k := range m.AdditionalTransitiveRefedObjectThingAnon {`,
		`		if val, ok := m.AdditionalTransitiveRefedObjectThingAnon[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedObjectThingAnon) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	if m.Prop1 != nil {`,
		`		if err := m.Prop1.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("prop1"`,
		`type AdditionalTransitiveRefedObjectThingAnonAnon struct {`,
		"	Discourse string `json:\"discourse,omitempty\"`",
		"	HoursSpent float64 `json:\"hoursSpent,omitempty\"`",
		"	AdditionalTransitiveRefedObjectThingAnonAnonAdditionalProperties map[string]interface{} `json:\"-\"`",
		`type AdditionalTransitiveRefedObjectThingAnonProp1 struct {`,
		"	ThisOneNotRequiredEither int64 `json:\"thisOneNotRequiredEither,omitempty\"`",
		"	AdditionalTransitiveRefedObjectThingAnonProp1 map[string]*AdditionalTransitiveRefedObjectThingAnonProp1Anon `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedObjectThingAnonProp1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateThisOneNotRequiredEither(formats); err != nil {`,
		`	for k := range m.AdditionalTransitiveRefedObjectThingAnonProp1 {`,
		`		if val, ok := m.AdditionalTransitiveRefedObjectThingAnonProp1[k]; ok {`,
		`			if val != nil {`,
		`				if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedObjectThingAnonProp1) validateThisOneNotRequiredEither(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ThisOneNotRequiredEither) {`,
		`	if err := validate.MaximumInt("prop1"+"."+"thisOneNotRequiredEither", "body", m.ThisOneNotRequiredEither, 20, false); err != nil {`,
		`type AdditionalTransitiveRefedObjectThingAnonProp1Anon struct {`,
		"	A1 strfmt.DateTime `json:\"a1,omitempty\"`",
		"	B1 strfmt.Date `json:\"b1,omitempty\"`",
		"	AdditionalTransitiveRefedObjectThingAnonProp1AnonAdditionalProperties map[string]interface{} `json:\"-\"`",
		`func (m *AdditionalTransitiveRefedObjectThingAnonProp1Anon) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateA1(formats); err != nil {`,
		`	if err := m.validateB1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *AdditionalTransitiveRefedObjectThingAnonProp1Anon) validateA1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.A1) {`,
		`	if err := validate.FormatOf("a1", "body", "date-time", m.A1.String(), formats); err != nil {`,
		`func (m *AdditionalTransitiveRefedObjectThingAnonProp1Anon) validateB1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.B1) {`,
		`	if err := validate.FormatOf("b1", "body", "date", m.B1.String(), formats); err != nil {`,
		// empty validation
		"func (m *AdditionalTransitiveRefedObjectThingAnonAnon) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

}

func initFixture1537() {
	// testing fixture-1537.yaml with flatten and expand (--skip-flatten)
	// TODO: expand

	/*
	   repro issue 1537
	*/

	f := newModelFixture("../fixtures/bugs/1537/fixture-1537.yaml", "param body required with array of objects")
	thisRun := f.AddRun(false)

	// load expectations for model: profile_array.go
	thisRun.AddExpectations("profile_array.go", []string{
		`type ProfileArray struct {`,
		"	ProfileCfg []*ProfileCfg `json:\"profileCfg\"`",
		`func (m *ProfileArray) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProfileCfg(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ProfileArray) validateProfileCfg(formats strfmt.Registry) error {`,
		`	if err := validate.Required("profileCfg", "body", m.ProfileCfg); err != nil {`,
		`	for i := 0; i < len(m.ProfileCfg); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`		if swag.IsZero(m.ProfileCfg[i]) {`,
		// nullable required:
		`		if m.ProfileCfg[i] != nil {`,
		`			if err := m.ProfileCfg[i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName("profileCfg" + "." + strconv.Itoa(i)`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: profile.go
	thisRun.AddExpectations("profile.go", []string{
		`type Profile struct {`,
		"	ProfileCfg ProfileCfgs `json:\"profileCfg,omitempty\"`",
		`func (m *Profile) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProfileCfg(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *Profile) validateProfileCfg(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.ProfileCfg) {`,
		`	if err := m.ProfileCfg.Validate(formats); err != nil {`,
		`		if ve, ok := err.(*errors.Validation); ok {`,
		`			return ve.ValidateName("profileCfg"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: profile_cfgs.go
	thisRun.AddExpectations("profile_cfgs.go", []string{
		`type ProfileCfgs []*ProfileCfg`,
		`func (m ProfileCfgs) Validate(formats strfmt.Registry) error {`,
		`	for i := 0; i < len(m); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`		if swag.IsZero(m[i]) {`,
		// nullable required:
		`		if m[i] != nil {`,
		`			if err := m[i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName(strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: profile_cfg.go
	thisRun.AddExpectations("profile_cfg.go", []string{
		`type ProfileCfg struct {`,
		"	Value1 int32 `json:\"value1,omitempty\"`",
		"	Value2 int32 `json:\"value2,omitempty\"`",
		// empty validation
		"func (m *ProfileCfg) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: profile_required.go
	thisRun.AddExpectations("profile_required.go", []string{
		`type ProfileRequired struct {`,
		"	ProfileCfg ProfileCfgs `json:\"profileCfg\"`",
		`func (m *ProfileRequired) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProfileCfg(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ProfileRequired) validateProfileCfg(formats strfmt.Registry) error {`,
		`	if err := validate.Required("profileCfg", "body", m.ProfileCfg); err != nil {`,
		`	if err := m.ProfileCfg.Validate(formats); err != nil {`,
		`		if ve, ok := err.(*errors.Validation); ok {`,
		`			return ve.ValidateName("profileCfg"`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

}

func initFixture1537v2() {
	// testing fixture-1537-2.yaml with flatten and expand (--skip-flatten)
	// TODO: expand

	/*
	   repro issue 1537, with aliased items
	*/

	f := newModelFixture("../fixtures/bugs/1537/fixture-1537-2.yaml", "param body required with array of aliased items")
	thisRun := f.AddRun(false)

	// load expectations for model: profiles.go
	thisRun.AddExpectations("profiles.go", []string{
		`type Profiles []ProfileCfgs`,
		`func (m Profiles) Validate(formats strfmt.Registry) error {`,
		`	for i := 0; i < len(m); i++ {`,
		`		if err := m[i].Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName(strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: profile_cfgs_no_validation.go
	thisRun.AddExpectations("profile_cfgs_no_validation.go", []string{
		`type ProfileCfgsNoValidation []*ProfileCfg`,
		`func (m ProfileCfgsNoValidation) Validate(formats strfmt.Registry) error {`,
		`	for i := 0; i < len(m); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`		if swag.IsZero(m[i]) {`,
		// nullable required:
		`		if m[i] != nil {`,
		`			if err := m[i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName(strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: profile_cfgs.go
	thisRun.AddExpectations("profile_cfgs.go", []string{
		`type ProfileCfgs []*ProfileCfg`,
		`func (m ProfileCfgs) Validate(formats strfmt.Registry) error {`,
		`	iProfileCfgsSize := int64(len(m)`,
		`	if err := validate.MaxItems("", "body", iProfileCfgsSize, 10); err != nil {`,
		`	for i := 0; i < len(m); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`		if swag.IsZero(m[i]) {`,
		// nullable required:
		`		if m[i] != nil {`,
		`			if err := m[i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName(strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: profile_cfg.go
	thisRun.AddExpectations("profile_cfg.go", []string{
		`type ProfileCfg struct {`,
		"	Value1 int32 `json:\"value1,omitempty\"`",
		"	Value2 int32 `json:\"value2,omitempty\"`",
		// empty validation
		"func (m *ProfileCfg) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: profiles_no_validation.go
	thisRun.AddExpectations("profiles_no_validation.go", []string{
		`type ProfilesNoValidation []ProfileCfgsNoValidation`,
		`func (m ProfilesNoValidation) Validate(formats strfmt.Registry) error {`,
		`	for i := 0; i < len(m); i++ {`,
		`		if err := m[i].Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName(strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

}

func initFixture15365() {
	// testing fixture-1536-5.yaml with flatten but NOT expand (--skip-flatten)

	f := newModelFixture("../fixtures/bugs/1536/fixture-1536-5.yaml", "param body with maps")
	thisRun := f.AddRun(false)

	// load expectations for model: model_array_of_nullable.go
	thisRun.AddExpectations("model_array_of_nullable.go", []string{
		`type ModelArrayOfNullable []*int64`,
		`func (m ModelArrayOfNullable) Validate(formats strfmt.Registry) error {`,
		`	for i := 0; i < len(m); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`		if swag.IsZero(m[i]) {`,
		// nullable required:
		`		if err := validate.MinimumInt(strconv.Itoa(i), "body", *m[i], 0, false); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: model_map_of_nullable_primitive.go
	thisRun.AddExpectations("model_map_of_nullable_primitive.go", []string{
		`type ModelMapOfNullablePrimitive map[string]*int64`,
		`func (m ModelMapOfNullablePrimitive) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		// do we need Required when element is nullable?
		// nullable not required:
		`		if swag.IsZero(m[k]) {`,
		`		if err := validate.MinimumInt(k, "body", *m[k], 0, false); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: model_array_with_max.go
	thisRun.AddExpectations("model_array_with_max.go", []string{
		`type ModelArrayWithMax []interface{`,
		`func (m ModelArrayWithMax) Validate(formats strfmt.Registry) error {`,
		`	iModelArrayWithMaxSize := int64(len(m)`,
		`	if err := validate.MaxItems("", "body", iModelArrayWithMaxSize, 10); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: model_array_of_x_nullable.go
	thisRun.AddExpectations("model_array_of_x_nullable.go", []string{
		`type ModelArrayOfXNullable []*int64`,
		// do we need Required when item is nullable?
		// nullable not required:
		"func (m ModelArrayOfXNullable) Validate(formats strfmt.Registry) error {\n	return nil\n}",
		// nullable required:
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: model_array_of_nullable_format.go
	thisRun.AddExpectations("model_array_of_nullable_format.go", []string{
		`type ModelArrayOfNullableFormat []*strfmt.UUID`,
		`func (m ModelArrayOfNullableFormat) Validate(formats strfmt.Registry) error {`,
		`	for i := 0; i < len(m); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`		if swag.IsZero(m[i]) {`,
		// nullable required:
		`		if err := validate.FormatOf(strconv.Itoa(i), "body", "uuid", m[i].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: model_map_of_of_slice_of_nullable_primitive.go
	thisRun.AddExpectations("model_map_of_of_slice_of_nullable_primitive.go", []string{
		`type ModelMapOfOfSliceOfNullablePrimitive map[string][]*int64`,
		`func (m ModelMapOfOfSliceOfNullablePrimitive) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		// do we need Required when element is nullable?
		// nullable not required:
		`		for i := 0; i < len(m[k]); i++ {`,
		// do we need Required when item is nullable?
		// nullable not required:
		`			if err := validate.MinimumInt(k+"."+strconv.Itoa(i), "body", *m[k][i], 0, false); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: model_map_of_ref.go
	thisRun.AddExpectations("model_map_of_ref.go", []string{
		`type ModelMapOfRef map[string]ModelArrayWithMax`,
		`func (m ModelMapOfRef) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if err := m[k].Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName(k`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: model_array_of_format.go
	thisRun.AddExpectations("model_array_of_format.go", []string{
		`type ModelArrayOfFormat []strfmt.UUID`,
		`func (m ModelArrayOfFormat) Validate(formats strfmt.Registry) error {`,
		`	for i := 0; i < len(m); i++ {`,
		`		if err := validate.FormatOf(strconv.Itoa(i), "body", "uuid", m[i].String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: model_array_of_ref_no_validations.go
	thisRun.AddExpectations("model_array_of_ref_no_validations.go", []string{
		`type ModelArrayOfRefNoValidations []ModelInterface`,
		// empty validation
		"func (m ModelArrayOfRefNoValidations) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: model_array_of_not_nullable.go
	thisRun.AddExpectations("model_array_of_not_nullable.go", []string{
		`type ModelArrayOfNotNullable []int64`,
		`func (m ModelArrayOfNotNullable) Validate(formats strfmt.Registry) error {`,
		`	for i := 0; i < len(m); i++ {`,
		`		if err := validate.MinimumInt(strconv.Itoa(i), "body", m[i], 10, false); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: model_array_of_nullable_string.go
	thisRun.AddExpectations("model_array_of_nullable_string.go", []string{
		`type ModelArrayOfNullableString []*string`,
		// do we need Required when item is nullable?
		// nullable not required:
		// empty validation
		"func (m ModelArrayOfNullableString) Validate(formats strfmt.Registry) error {\n	return nil\n}",
		// nullable required:
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: model_map_of_x_nullable_primitive.go
	thisRun.AddExpectations("model_map_of_x_nullable_primitive.go", []string{
		`type ModelMapOfXNullablePrimitive map[string]*int64`,
		`func (m ModelMapOfXNullablePrimitive) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		// do we need...?
		`		if swag.IsZero(m[k]) {`,
		`		if err := validate.MinimumInt(k, "body", *m[k], 100, false); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: model_array_of_not_nullable_string.go
	thisRun.AddExpectations("model_array_of_not_nullable_string.go", []string{
		`type ModelArrayOfNotNullableString []string`,
		// empty validation
		"func (m ModelArrayOfNotNullableString) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: model_array_of_ref_slice_validations.go
	thisRun.AddExpectations("model_array_of_ref_slice_validations.go", []string{
		`type ModelArrayOfRefSliceValidations []ModelInterface`,
		`func (m ModelArrayOfRefSliceValidations) Validate(formats strfmt.Registry) error {`,
		`	iModelArrayOfRefSliceValidationsSize := int64(len(m)`,
		`	if err := validate.MaxItems("", "body", iModelArrayOfRefSliceValidationsSize, 10); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: model_interface.go
	thisRun.AddExpectations("model_interface.go", []string{
		`type ModelInterface interface{`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: model_array_of_ref.go
	thisRun.AddExpectations("model_array_of_ref.go", []string{
		`type ModelArrayOfRef []ModelArrayOfXNullable`,
		`func (m ModelArrayOfRef) Validate(formats strfmt.Registry) error {`,
		`	for i := 0; i < len(m); i++ {`,
		`		if err := m[i].Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName(strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)
}

func initFixture1548() {
	// testing fixture-1548.yaml with flatten

	/*
		My App API: check that there is no format validation on Base64 types
	*/

	f := newModelFixture("../fixtures/bugs/1548/fixture-1548.yaml", "My App API")
	thisRun := f.AddRun(false)

	// load expectations for model: base64_alias.go
	thisRun.AddExpectations("base64_alias.go", []string{
		`type Base64Alias strfmt.Base64`,
		`func (m *Base64Alias) UnmarshalJSON(b []byte) error {`,
		`	return ((*strfmt.Base64)(m)).UnmarshalJSON(b`,
		`func (m Base64Alias) MarshalJSON() ([]byte, error) {`,
		`	return (strfmt.Base64(m)).MarshalJSON(`,
		`func (m Base64Alias) Validate(formats strfmt.Registry) error {`,
	},
		// not expected
		[]string{"TODO",
			"validate.FormatOf(",
			`return errors.CompositeValidationError(res...`,
		},
		// output in log
		noLines,
		noLines)

	// load expectations for model: base64_map.go
	thisRun.AddExpectations("base64_map.go", []string{
		`type Base64Map map[string]strfmt.Base64`,
		`func (m Base64Map) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if err := validate.MaxLength(k, "body", m[k].String(), 100); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		[]string{"TODO", "validate.FormatOf("},
		// output in log
		noLines,
		noLines)

	// load expectations for model: base64_array.go
	thisRun.AddExpectations("base64_array.go", []string{
		`type Base64Array []strfmt.Base64`,
		`func (m Base64Array) Validate(formats strfmt.Registry) error {`,
	},
		// not expected
		[]string{
			"TODO",
			"validate.FormatOf(",
			`	for i := 0; i < len(m); i++ {`,
			`		return errors.CompositeValidationError(res...`,
		},
		// output in log
		noLines,
		noLines)

	// load expectations for model: base64_model.go
	thisRun.AddExpectations("base64_model.go", []string{
		`type Base64Model struct {`,
		"	Prop1 strfmt.Base64 `json:\"prop1,omitempty\"`",
		`func (m *Base64Model) Validate(formats strfmt.Registry) error {`,
	},
		// not expected
		[]string{
			"TODO",
			"validate.FormatOf(",
			`	if err := m.validateProp1(formats); err != nil {`,
			`func (m *Base64Model) validateProp1(formats strfmt.Registry) error {`,
		},
		// output in log
		noLines,
		noLines)
}

func initFixtureSimpleTuple() {
	// testing fixture-simple-tuple.yaml with flatten

	/*
	   A basic test of for serialization generation for tuples and additionalItems.

	*/
	f := newModelFixture("../fixtures/bugs/1571/fixture-simple-tuple.yaml", "fixture for serializing tuples")
	flattenRun := f.AddRun(false).WithMinimalFlatten(true)
	expandRun := f.AddRun(true)

	// load expectations for model: tuple_thing_with_map_element.go
	flattenRun.AddExpectations("tuple_thing_with_map_element.go", []string{
		`type TupleThingWithMapElement struct {`,
		"	P0 map[string]string `json:\"-\"`",
		"	P1 map[string]int64 `json:\"-\"`",
		"	TupleThingWithMapElementItems []map[string]strfmt.Date `json:\"-\"`",
		`func (m *TupleThingWithMapElement) UnmarshalJSON(raw []byte) error {`,
		`	var stage1 []json.RawMessage`,
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&stage1); err != nil {`,
		`	var lastIndex int`,
		`	if len(stage1) > 0 {`,
		`		var dataP0 map[string]string`,
		`		buf = bytes.NewBuffer(stage1[0]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP0); err != nil {`,
		`		m.P0 = dataP0`,
		`		lastIndex = 0`,
		`	if len(stage1) > 1 {`,
		`		var dataP1 map[string]int64`,
		`		buf = bytes.NewBuffer(stage1[1]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP1); err != nil {`,
		`		m.P1 = dataP1`,
		`		lastIndex = 1`,
		`	if len(stage1) > lastIndex+1 {`,
		`		for _, val := range stage1[lastIndex+1:] {`,
		`			var toadd map[string]strfmt.Date`,
		`			buf = bytes.NewBuffer(val`,
		`			dec := json.NewDecoder(buf`,
		`			dec.UseNumber(`,
		`			if err := dec.Decode(&toadd); err != nil {`,
		`			m.TupleThingWithMapElementItems = append(m.TupleThingWithMapElementItems, toadd`,
		`func (m TupleThingWithMapElement) MarshalJSON() ([]byte, error) {`,
		`	data := []interface{}{`,
		`		m.P0,`,
		`		m.P1,`,
		`	for _, v := range m.TupleThingWithMapElementItems {`,
		`		data = append(data, v`,
		`	return json.Marshal(data`,
		`func (m *TupleThingWithMapElement) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateTupleThingWithMapElementItems(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *TupleThingWithMapElement) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	for k := range m.P0 {`,
		`		if err := validate.MaxLength("0"+"."+k, "body", m.P0[k], 10); err != nil {`,
		`func (m *TupleThingWithMapElement) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("1", "body", m.P1); err != nil {`,
		`	for k := range m.P1 {`,
		`		if err := validate.MinimumInt("1"+"."+k, "body", m.P1[k], 10, false); err != nil {`,
		`func (m *TupleThingWithMapElement) validateTupleThingWithMapElementItems(formats strfmt.Registry) error {`,
		`	for i := range m.TupleThingWithMapElementItems {`,
		`		for k := range m.TupleThingWithMapElementItems[i] {`,
		`			if err := validate.FormatOf(strconv.Itoa(i+2)+"."+k, "body", "date", ` +
			`m.TupleThingWithMapElementItems[i][k].String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("tuple_thing_with_map_element.go",
		flattenRun.ExpectedFor("TupleThingWithMapElement").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: true_tuple_thing.go
	flattenRun.AddExpectations("true_tuple_thing.go", []string{
		`type TrueTupleThing struct {`,
		"	P0 *float64 `json:\"-\"`",
		"	P1 *string `json:\"-\"`",
		"	TrueTupleThingItems []interface{} `json:\"-\"`",
		`func (m *TrueTupleThing) UnmarshalJSON(raw []byte) error {`,
		`	var stage1 []json.RawMessage`,
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&stage1); err != nil {`,
		`	var lastIndex int`,
		`	if len(stage1) > 0 {`,
		`		var dataP0 float64`,
		`		buf = bytes.NewBuffer(stage1[0]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP0); err != nil {`,
		`		m.P0 = &dataP0`,
		`		lastIndex = 0`,
		`	if len(stage1) > 1 {`,
		`		var dataP1 string`,
		`		buf = bytes.NewBuffer(stage1[1]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP1); err != nil {`,
		`		m.P1 = &dataP1`,
		`		lastIndex = 1`,
		`	if len(stage1) > lastIndex+1 {`,
		`		for _, val := range stage1[lastIndex+1:] {`,
		`			var toadd interface{`,
		`			buf = bytes.NewBuffer(val`,
		`			dec := json.NewDecoder(buf`,
		`			dec.UseNumber(`,
		`			if err := dec.Decode(&toadd); err != nil {`,
		`			m.TrueTupleThingItems = append(m.TrueTupleThingItems, toadd`,
		`func (m TrueTupleThing) MarshalJSON() ([]byte, error) {`,
		`	data := []interface{}{`,
		`		m.P0,`,
		`		m.P1,`,
		`	for _, v := range m.TrueTupleThingItems {`,
		`		data = append(data, v`,
		`	return json.Marshal(data`,
		`func (m *TrueTupleThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateTrueTupleThingItems(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *TrueTupleThing) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`func (m *TrueTupleThing) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("1", "body", m.P1); err != nil {`,
		`func (m *TrueTupleThing) validateTrueTupleThingItems(formats strfmt.Registry) error {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("true_tuple_thing.go",
		flattenRun.ExpectedFor("TrueTupleThing").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: false_tuple_thing.go
	flattenRun.AddExpectations("false_tuple_thing.go", []string{
		`type FalseTupleThing struct {`,
		"	P0 *float64 `json:\"-\"`",
		"	P1 *string `json:\"-\"`",
		`func (m *FalseTupleThing) UnmarshalJSON(raw []byte) error {`,
		`	var stage1 []json.RawMessage`,
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&stage1); err != nil {`,
		`	if len(stage1) > 0 {`,
		`		var dataP0 float64`,
		`		buf = bytes.NewBuffer(stage1[0]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP0); err != nil {`,
		`		m.P0 = &dataP0`,
		`	if len(stage1) > 1 {`,
		`		var dataP1 string`,
		`		buf = bytes.NewBuffer(stage1[1]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP1); err != nil {`,
		`		m.P1 = &dataP1`,
		`func (m FalseTupleThing) MarshalJSON() ([]byte, error) {`,
		`	data := []interface{}{`,
		`		m.P0,`,
		`		m.P1,`,
		`	return json.Marshal(data`,
		`func (m *FalseTupleThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *FalseTupleThing) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`func (m *FalseTupleThing) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("1", "body", m.P1); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("false_tuple_thing.go",
		flattenRun.ExpectedFor("FalseTupleThing").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: tuple_thing_with_not_nullable.go
	flattenRun.AddExpectations("tuple_thing_with_not_nullable.go", []string{
		`type TupleThingWithNotNullable struct {`,
		"	P0 string `json:\"-\"`",
		"	P1 *int64 `json:\"-\"`",
		"	TupleThingWithNotNullableItems []interface{} `json:\"-\"`",
		`func (m *TupleThingWithNotNullable) UnmarshalJSON(raw []byte) error {`,
		`	var stage1 []json.RawMessage`,
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&stage1); err != nil {`,
		`	var lastIndex int`,
		`	if len(stage1) > 0 {`,
		`		var dataP0 string`,
		`		buf = bytes.NewBuffer(stage1[0]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP0); err != nil {`,
		`		m.P0 = dataP0`,
		`		lastIndex = 0`,
		`	if len(stage1) > 1 {`,
		`		var dataP1 int64`,
		`		buf = bytes.NewBuffer(stage1[1]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP1); err != nil {`,
		`		m.P1 = &dataP1`,
		`		lastIndex = 1`,
		`	if len(stage1) > lastIndex+1 {`,
		`		for _, val := range stage1[lastIndex+1:] {`,
		`			var toadd interface{`,
		`			buf = bytes.NewBuffer(val`,
		`			dec := json.NewDecoder(buf`,
		`			dec.UseNumber(`,
		`			if err := dec.Decode(&toadd); err != nil {`,
		`			m.TupleThingWithNotNullableItems = append(m.TupleThingWithNotNullableItems, toadd`,
		`func (m TupleThingWithNotNullable) MarshalJSON() ([]byte, error) {`,
		`	data := []interface{}{`,
		`		m.P0,`,
		`		m.P1,`,
		`	for _, v := range m.TupleThingWithNotNullableItems {`,
		`		data = append(data, v`,
		`	return json.Marshal(data`,
		`func (m *TupleThingWithNotNullable) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateTupleThingWithNotNullableItems(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *TupleThingWithNotNullable) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.RequiredString("0", "body", m.P0); err != nil {`,
		`	if err := validate.MaxLength("0", "body", m.P0, 10); err != nil {`,
		`func (m *TupleThingWithNotNullable) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("1", "body", m.P1); err != nil {`,
		`	if err := validate.MaximumInt("1", "body", *m.P1, 10, false); err != nil {`,
		`func (m *TupleThingWithNotNullable) validateTupleThingWithNotNullableItems(formats strfmt.Registry) error {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("tuple_thing_with_not_nullable.go",
		flattenRun.ExpectedFor("TupleThingWithNotNullable").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: tuple_thing.go
	flattenRun.AddExpectations("tuple_thing.go", []string{
		`type TupleThing struct {`,
		"	P0 *string `json:\"-\"`",
		"	P1 *string `json:\"-\"`",
		`func (m *TupleThing) UnmarshalJSON(raw []byte) error {`,
		`	var stage1 []json.RawMessage`,
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&stage1); err != nil {`,
		`	if len(stage1) > 0 {`,
		`		var dataP0 string`,
		`		buf = bytes.NewBuffer(stage1[0]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP0); err != nil {`,
		`		m.P0 = &dataP0`,
		`	if len(stage1) > 1 {`,
		`		var dataP1 string`,
		`		buf = bytes.NewBuffer(stage1[1]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP1); err != nil {`,
		`		m.P1 = &dataP1`,
		`func (m TupleThing) MarshalJSON() ([]byte, error) {`,
		`	data := []interface{}{`,
		`		m.P0,`,
		`		m.P1,`,
		`	return json.Marshal(data`,
		`func (m *TupleThing) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`var tupleThingTypeP0PropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"CONST1\",\"CONST2\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		tupleThingTypeP0PropEnum = append(tupleThingTypeP0PropEnum, v`,
		`	TupleThingP0CONST1 string = "CONST1"`,
		`	TupleThingP0CONST2 string = "CONST2"`,
		`func (m *TupleThing) validateP0Enum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, tupleThingTypeP0PropEnum, true); err != nil {`,
		`func (m *TupleThing) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	if err := m.validateP0Enum("0", "body", *m.P0); err != nil {`,
		`var tupleThingTypeP1PropEnum []interface{`,
		`	var res []string`,
		"	if err := json.Unmarshal([]byte(`[\"CONST3\",\"CONST4\"]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		tupleThingTypeP1PropEnum = append(tupleThingTypeP1PropEnum, v`,
		`	TupleThingP1CONST3 string = "CONST3"`,
		`	TupleThingP1CONST4 string = "CONST4"`,
		`func (m *TupleThing) validateP1Enum(path, location string, value string) error {`,
		`	if err := validate.EnumCase(path, location, value, tupleThingTypeP1PropEnum, true); err != nil {`,
		`func (m *TupleThing) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("1", "body", m.P1); err != nil {`,
		`	if err := m.validateP1Enum("1", "body", *m.P1); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("tuple_thing.go",
		flattenRun.ExpectedFor("TupleThing").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: tuple_thing_with_additional_items.go
	flattenRun.AddExpectations("tuple_thing_with_additional_items.go", []string{
		`type TupleThingWithAdditionalItems struct {`,
		"	P0 *string `json:\"-\"`",
		"	P1 *int64 `json:\"-\"`",
		"	TupleThingWithAdditionalItemsItems []int64 `json:\"-\"`",
		`func (m *TupleThingWithAdditionalItems) UnmarshalJSON(raw []byte) error {`,
		`	var stage1 []json.RawMessage`,
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&stage1); err != nil {`,
		`	var lastIndex int`,
		`	if len(stage1) > 0 {`,
		`		var dataP0 string`,
		`		buf = bytes.NewBuffer(stage1[0]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP0); err != nil {`,
		`		m.P0 = &dataP0`,
		`		lastIndex = 0`,
		`	if len(stage1) > 1 {`,
		`		var dataP1 int64`,
		`		buf = bytes.NewBuffer(stage1[1]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP1); err != nil {`,
		`		m.P1 = &dataP1`,
		`		lastIndex = 1`,
		`	if len(stage1) > lastIndex+1 {`,
		`		for _, val := range stage1[lastIndex+1:] {`,
		`			var toadd int64`,
		`			buf = bytes.NewBuffer(val`,
		`			dec := json.NewDecoder(buf`,
		`			dec.UseNumber(`,
		`			if err := dec.Decode(&toadd); err != nil {`,
		`			m.TupleThingWithAdditionalItemsItems = append(m.TupleThingWithAdditionalItemsItems, toadd`,
		`func (m TupleThingWithAdditionalItems) MarshalJSON() ([]byte, error) {`,
		`	data := []interface{}{`,
		`		m.P0,`,
		`		m.P1,`,
		`	for _, v := range m.TupleThingWithAdditionalItemsItems {`,
		`		data = append(data, v`,
		`	return json.Marshal(data`,
		`func (m *TupleThingWithAdditionalItems) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateTupleThingWithAdditionalItemsItems(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *TupleThingWithAdditionalItems) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`func (m *TupleThingWithAdditionalItems) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("1", "body", m.P1); err != nil {`,
		`	if err := validate.MaximumInt("1", "body", *m.P1, 10, false); err != nil {`,
		`var tupleThingWithAdditionalItemsItemsEnum []interface{`,
		`	var res []int64`,
		"	if err := json.Unmarshal([]byte(`[1,2]`), &res); err != nil {",
		`	for _, v := range res {`,
		`		tupleThingWithAdditionalItemsItemsEnum = append(tupleThingWithAdditionalItemsItemsEnum, v`,
		`func (m *TupleThingWithAdditionalItems) validateTupleThingWithAdditionalItemsItemsEnum(path,` +
			` location string, value int64) error {`,
		`	if err := validate.EnumCase(path, location, value, tupleThingWithAdditionalItemsItemsEnum, true); err != nil {`,
		`func (m *TupleThingWithAdditionalItems) validateTupleThingWithAdditionalItemsItems(formats strfmt.Registry)` +
			` error {`,
		`	for i := range m.TupleThingWithAdditionalItemsItems {`,
		`		if err := m.validateTupleThingWithAdditionalItemsItemsEnum(strconv.Itoa(i+2), "body", ` +
			`m.TupleThingWithAdditionalItemsItems[i]); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("tuple_thing_with_additional_items.go",
		flattenRun.ExpectedFor("TupleThingWithAdditionalItems").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: tuple_thing_with_array_element.go
	flattenRun.AddExpectations("tuple_thing_with_array_element.go", []string{
		`type TupleThingWithArrayElement struct {`,
		"	P0 []string `json:\"-\"`",
		"	P1 []int64 `json:\"-\"`",
		"	TupleThingWithArrayElementItems [][]strfmt.Date `json:\"-\"`",
		`func (m *TupleThingWithArrayElement) UnmarshalJSON(raw []byte) error {`,
		`	var stage1 []json.RawMessage`,
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&stage1); err != nil {`,
		`	var lastIndex int`,
		`	if len(stage1) > 0 {`,
		`		var dataP0 []string`,
		`		buf = bytes.NewBuffer(stage1[0]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP0); err != nil {`,
		`		m.P0 = dataP0`,
		`		lastIndex = 0`,
		`	if len(stage1) > 1 {`,
		`		var dataP1 []int64`,
		`		buf = bytes.NewBuffer(stage1[1]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP1); err != nil {`,
		`		m.P1 = dataP1`,
		`		lastIndex = 1`,
		`	if len(stage1) > lastIndex+1 {`,
		`		for _, val := range stage1[lastIndex+1:] {`,
		`			var toadd []strfmt.Date`,
		`			buf = bytes.NewBuffer(val`,
		`			dec := json.NewDecoder(buf`,
		`			dec.UseNumber(`,
		`			if err := dec.Decode(&toadd); err != nil {`,
		`			m.TupleThingWithArrayElementItems = append(m.TupleThingWithArrayElementItems, toadd`,
		`func (m TupleThingWithArrayElement) MarshalJSON() ([]byte, error) {`,
		`	data := []interface{}{`,
		`		m.P0,`,
		`		m.P1,`,
		`	for _, v := range m.TupleThingWithArrayElementItems {`,
		`		data = append(data, v`,
		`	return json.Marshal(data`,
		`func (m *TupleThingWithArrayElement) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateTupleThingWithArrayElementItems(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *TupleThingWithArrayElement) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	iP0Size := int64(len(m.P0)`,
		`	if err := validate.MaxItems("0", "body", iP0Size, 10); err != nil {`,
		`	for i := 0; i < len(m.P0); i++ {`,
		`		if err := validate.MaxLength("0"+"."+strconv.Itoa(i), "body", m.P0[i], 10); err != nil {`,
		`func (m *TupleThingWithArrayElement) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	iP1Size := int64(len(m.P1)`,
		`	if err := validate.MinItems("1", "body", iP1Size, 20); err != nil {`,
		`	for i := 0; i < len(m.P1); i++ {`,
		`		if err := validate.MinimumInt("1"+"."+strconv.Itoa(i), "body", m.P1[i], 10, false); err != nil {`,
		`func (m *TupleThingWithArrayElement) validateTupleThingWithArrayElementItems(formats strfmt.Registry) error {`,
		`	for i := range m.TupleThingWithArrayElementItems {`,
		`		for ii := 0; ii < len(m.TupleThingWithArrayElementItems[i]); ii++ {`,
		`			if err := validate.FormatOf(strconv.Itoa(i+2)+"."+strconv.Itoa(ii), "body", ` +
			`"date", m.TupleThingWithArrayElementItems[i][ii].String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("tuple_thing_with_array_element.go",
		flattenRun.ExpectedFor("TupleThingWithArrayElement").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: array_of_tuples.go
	flattenRun.AddExpectations("array_of_tuples.go", []string{
		`type ArrayOfTuples []ArrayOfTuplesTuple0`,
		`func (m ArrayOfTuples) Validate(formats strfmt.Registry) error {`,
		`	iArrayOfTuplesSize := int64(len(m)`,
		`	if err := validate.MinItems("", "body", iArrayOfTuplesSize, 1); err != nil {`,
		`	for i := 0; i < len(m); i++ {`,
		`		if err := m[i].Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName(strconv.Itoa(i)`,
		`		return errors.CompositeValidationError(res...`,
		`type ArrayOfTuplesTuple0 struct {`,
		"	P0 []string `json:\"-\"`",
		"	P1 []int64 `json:\"-\"`",
		"	ArrayOfTuplesTuple0Items [][]strfmt.Date `json:\"-\"`",
		`func (m *ArrayOfTuplesTuple0) UnmarshalJSON(raw []byte) error {`,
		`	var stage1 []json.RawMessage`,
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&stage1); err != nil {`,
		`	var lastIndex int`,
		`	if len(stage1) > 0 {`,
		`		var dataP0 []string`,
		`		buf = bytes.NewBuffer(stage1[0]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP0); err != nil {`,
		`		m.P0 = dataP0`,
		`		lastIndex = 0`,
		`	if len(stage1) > 1 {`,
		`		var dataP1 []int64`,
		`		buf = bytes.NewBuffer(stage1[1]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP1); err != nil {`,
		`		m.P1 = dataP1`,
		`		lastIndex = 1`,
		`	if len(stage1) > lastIndex+1 {`,
		`		for _, val := range stage1[lastIndex+1:] {`,
		`			var toadd []strfmt.Date`,
		`			buf = bytes.NewBuffer(val`,
		`			dec := json.NewDecoder(buf`,
		`			dec.UseNumber(`,
		`			if err := dec.Decode(&toadd); err != nil {`,
		`			m.ArrayOfTuplesTuple0Items = append(m.ArrayOfTuplesTuple0Items, toadd`,
		`func (m ArrayOfTuplesTuple0) MarshalJSON() ([]byte, error) {`,
		`	data := []interface{}{`,
		`		m.P0,`,
		`		m.P1,`,
		`	for _, v := range m.ArrayOfTuplesTuple0Items {`,
		`		data = append(data, v`,
		`	return json.Marshal(data`,
		`func (m *ArrayOfTuplesTuple0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateArrayOfTuplesTuple0Items(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ArrayOfTuplesTuple0) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P0", "body", m.P0); err != nil {`,
		`	iP0Size := int64(len(m.P0)`,
		`	if err := validate.MaxItems("P0", "body", iP0Size, 10); err != nil {`,
		`	for i := 0; i < len(m.P0); i++ {`,
		`		if err := validate.MaxLength("P0"+"."+strconv.Itoa(i), "body", m.P0[i], 10); err != nil {`,
		`func (m *ArrayOfTuplesTuple0) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("P0", "body", m.P0); err != nil {`,
		`	iP1Size := int64(len(m.P1)`,
		`	if err := validate.MinItems("P1", "body", iP1Size, 20); err != nil {`,
		`	for i := 0; i < len(m.P1); i++ {`,
		`		if err := validate.MinimumInt("P1"+"."+strconv.Itoa(i), ` +
			`"body", m.P1[i], 10, false); err != nil {`,
		`func (m *ArrayOfTuplesTuple0) validateArrayOfTuplesTuple0Items(formats strfmt.Registry) error {`,
		`	for i := range m.ArrayOfTuplesTuple0Items {`,
		`		for ii := 0; ii < len(m.ArrayOfTuplesTuple0Items[i]); ii++ {`,
		`			if err := validate.FormatOf(strconv.Itoa(i)+"."+strconv.Itoa(ii), ` +
			`"body", "date", m.ArrayOfTuplesTuple0Items[i][ii].String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("array_of_tuples.go",
		flattenRun.ExpectedFor("ArrayOfTuples").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: tuple_thing_with_object_element.go
	flattenRun.AddExpectations("tuple_thing_with_object_element.go", []string{
		`type TupleThingWithObjectElement struct {`,
		"	P0 *TupleThingWithObjectElementItems0 `json:\"-\"`",
		"	P1 map[string]int64 `json:\"-\"`",
		"	TupleThingWithObjectElementItems []map[string]strfmt.Date `json:\"-\"`",
		`func (m *TupleThingWithObjectElement) UnmarshalJSON(raw []byte) error {`,
		`	var stage1 []json.RawMessage`,
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&stage1); err != nil {`,
		`	var lastIndex int`,
		`	if len(stage1) > 0 {`,
		`		var dataP0 TupleThingWithObjectElementItems0`,
		`		buf = bytes.NewBuffer(stage1[0]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP0); err != nil {`,
		`		m.P0 = &dataP0`,
		`		lastIndex = 0`,
		`	if len(stage1) > 1 {`,
		`		var dataP1 map[string]int64`,
		`		buf = bytes.NewBuffer(stage1[1]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP1); err != nil {`,
		`		m.P1 = dataP1`,
		`		lastIndex = 1`,
		`	if len(stage1) > lastIndex+1 {`,
		`		for _, val := range stage1[lastIndex+1:] {`,
		`			var toadd map[string]strfmt.Date`,
		`			buf = bytes.NewBuffer(val`,
		`			dec := json.NewDecoder(buf`,
		`			dec.UseNumber(`,
		`			if err := dec.Decode(&toadd); err != nil {`,
		`			m.TupleThingWithObjectElementItems = append(m.TupleThingWithObjectElementItems, toadd`,
		`func (m TupleThingWithObjectElement) MarshalJSON() ([]byte, error) {`,
		`	data := []interface{}{`,
		`		m.P0,`,
		`		m.P1,`,
		`	for _, v := range m.TupleThingWithObjectElementItems {`,
		`		data = append(data, v`,
		`	return json.Marshal(data`,
		`func (m *TupleThingWithObjectElement) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateTupleThingWithObjectElementItems(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *TupleThingWithObjectElement) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	if m.P0 != nil {`,
		`		if err := m.P0.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("0"`,
		`func (m *TupleThingWithObjectElement) validateP1(formats strfmt.Registry) error {`,
		`if err := validate.Required("1", "body", m.P1); err != nil {`,
		`	for k := range m.P1 {`,
		`		if err := validate.MinimumInt("1"+"."+k, "body", m.P1[k], 10, false); err != nil {`,
		`func (m *TupleThingWithObjectElement) validateTupleThingWithObjectElementItems(formats strfmt.Registry)` +
			` error {`,
		`	for i := range m.TupleThingWithObjectElementItems {`,
		`		for k := range m.TupleThingWithObjectElementItems[i] {`,
		`			if err := validate.FormatOf(strconv.Itoa(i+2)+"."+k, "body", "date", ` +
			`m.TupleThingWithObjectElementItems[i][k].String(), formats); err != nil {`,
		`type TupleThingWithObjectElementItems0 struct {`,
		"	Prop1 string `json:\"prop1,omitempty\"`",
		`func (m *TupleThingWithObjectElementItems0) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *TupleThingWithObjectElementItems0) validateProp1(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1) {`,
		`	if err := validate.MaxLength("prop1", "body", m.Prop1, 10); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("tuple_thing_with_object_element.go",
		flattenRun.ExpectedFor("TupleThingWithObjectElement").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: tuple_thing_with_no_additional_items.go
	flattenRun.AddExpectations("tuple_thing_with_no_additional_items.go", []string{
		`type TupleThingWithNoAdditionalItems struct {`,
		"	P0 *string `json:\"-\"`",
		"	P1 *int64 `json:\"-\"`",
		`func (m *TupleThingWithNoAdditionalItems) UnmarshalJSON(raw []byte) error {`,
		`	var stage1 []json.RawMessage`,
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&stage1); err != nil {`,
		`	if len(stage1) > 0 {`,
		`		var dataP0 string`,
		`		buf = bytes.NewBuffer(stage1[0]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP0); err != nil {`,
		`		m.P0 = &dataP0`,
		`	if len(stage1) > 1 {`,
		`		var dataP1 int64`,
		`		buf = bytes.NewBuffer(stage1[1]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP1); err != nil {`,
		`		m.P1 = &dataP1`,
		`func (m TupleThingWithNoAdditionalItems) MarshalJSON() ([]byte, error) {`,
		`	data := []interface{}{`,
		`		m.P0,`,
		`		m.P1,`,
		`	return json.Marshal(data`,
		`func (m *TupleThingWithNoAdditionalItems) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *TupleThingWithNoAdditionalItems) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`func (m *TupleThingWithNoAdditionalItems) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	if err := validate.MaximumInt("1", "body", *m.P1, 10, false); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("tuple_thing_with_no_additional_items.go",
		flattenRun.ExpectedFor("TupleThingWithNoAdditionalItems").ExpectedLines, todo, noLines, noLines)

	// load expectations for model: tuple_thing_with_any_additional_items.go
	flattenRun.AddExpectations("tuple_thing_with_any_additional_items.go", []string{
		`type TupleThingWithAnyAdditionalItems struct {`,
		"	P0 *string `json:\"-\"`",
		"	P1 *int64 `json:\"-\"`",
		"	TupleThingWithAnyAdditionalItemsItems []interface{} `json:\"-\"`",
		`func (m *TupleThingWithAnyAdditionalItems) UnmarshalJSON(raw []byte) error {`,
		`	var stage1 []json.RawMessage`,
		`	buf := bytes.NewBuffer(raw`,
		`	dec := json.NewDecoder(buf`,
		`	dec.UseNumber(`,
		`	if err := dec.Decode(&stage1); err != nil {`,
		`	var lastIndex int`,
		`	if len(stage1) > 0 {`,
		`		var dataP0 string`,
		`		buf = bytes.NewBuffer(stage1[0]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP0); err != nil {`,
		`		m.P0 = &dataP0`,
		`		lastIndex = 0`,
		`	if len(stage1) > 1 {`,
		`		var dataP1 int64`,
		`		buf = bytes.NewBuffer(stage1[1]`,
		`		dec := json.NewDecoder(buf`,
		`		dec.UseNumber(`,
		`		if err := dec.Decode(&dataP1); err != nil {`,
		`		m.P1 = &dataP1`,
		`		lastIndex = 1`,
		`	if len(stage1) > lastIndex+1 {`,
		`		for _, val := range stage1[lastIndex+1:] {`,
		`			var toadd interface{`,
		`			buf = bytes.NewBuffer(val`,
		`			dec := json.NewDecoder(buf`,
		`			dec.UseNumber(`,
		`			if err := dec.Decode(&toadd); err != nil {`,
		`			m.TupleThingWithAnyAdditionalItemsItems = append(m.TupleThingWithAnyAdditionalItemsItems, toadd`,
		`func (m TupleThingWithAnyAdditionalItems) MarshalJSON() ([]byte, error) {`,
		`	data := []interface{}{`,
		`		m.P0,`,
		`		m.P1,`,
		`	for _, v := range m.TupleThingWithAnyAdditionalItemsItems {`,
		`		data = append(data, v`,
		`	return json.Marshal(data`,
		`func (m *TupleThingWithAnyAdditionalItems) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateP0(formats); err != nil {`,
		`	if err := m.validateP1(formats); err != nil {`,
		`	if err := m.validateTupleThingWithAnyAdditionalItemsItems(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *TupleThingWithAnyAdditionalItems) validateP0(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	if err := validate.MaxLength("0", "body", *m.P0, 10); err != nil {`,
		`func (m *TupleThingWithAnyAdditionalItems) validateP1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("0", "body", m.P0); err != nil {`,
		`	if err := validate.MaximumInt("1", "body", *m.P1, 10, false); err != nil {`,
		`func (m *TupleThingWithAnyAdditionalItems) ` +
			`validateTupleThingWithAnyAdditionalItemsItems(formats strfmt.Registry) error {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	expandRun.AddExpectations("tuple_thing_with_any_additional_items.go",
		flattenRun.ExpectedFor("TupleThingWithAnyAdditionalItems").ExpectedLines, todo, noLines, noLines)
}

func initFixtureDeepMaps() {
	// testing fixture-deepMaps.yaml with minimal flatten

	f := newModelFixture("../fixtures/enhancements/1572/fixture-deepMaps.yaml", "issue 1572 - deep maps")
	thisRun := f.AddRun(false).WithMinimalFlatten(true)

	// load expectations for model: model_object_vanilla.go
	thisRun.AddExpectations("model_object_vanilla.go", []string{
		`type ModelObjectVanilla struct {`,
		"	Prop0 *ModelSanity `json:\"prop0,omitempty\"`",
		"	Prop1 *ModelSanity `json:\"prop1\"`",
		"	Prop2 []*ModelSanity `json:\"prop2\"`",
		"	Prop3 *ModelSanity `json:\"prop3,omitempty\"`",
		"	Prop4 map[string]ModelSanity `json:\"prop4,omitempty\"`",
		"	Prop5 int64 `json:\"prop5,omitempty\"`",
		"	ModelObjectVanilla map[string]map[string]map[string]ModelSanity `json:\"-\"`",
		`func (m *ModelObjectVanilla) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp0(formats); err != nil {`,
		`	if err := m.validateProp1(formats); err != nil {`,
		`	if err := m.validateProp2(formats); err != nil {`,
		`	if err := m.validateProp3(formats); err != nil {`,
		`	if err := m.validateProp4(formats); err != nil {`,
		`	for k := range m.ModelObjectVanilla {`,
		`		for kk := range m.ModelObjectVanilla[k] {`,
		`			for kkk := range m.ModelObjectVanilla[k][kk] {`,
		`				if err := validate.Required(k+"."+kk+"."+kkk, "body",` +
			` m.ModelObjectVanilla[k][kk][kkk]); err != nil {`,
		`				if val, ok := m.ModelObjectVanilla[k][kk][kkk]; ok {`,
		`					if err := val.Validate(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ModelObjectVanilla) validateProp0(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop0) {`,
		`	if m.Prop0 != nil {`,
		`		if err := m.Prop0.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("prop0"`,
		`func (m *ModelObjectVanilla) validateProp1(formats strfmt.Registry) error {`,
		`	if err := validate.Required("prop1", "body", m.Prop1); err != nil {`,
		`	if m.Prop1 != nil {`,
		`		if err := m.Prop1.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("prop1"`,
		`func (m *ModelObjectVanilla) validateProp2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop2) {`,
		`	for i := 0; i < len(m.Prop2); i++ {`,
		`		if swag.IsZero(m.Prop2[i]) {`,
		`		if m.Prop2[i] != nil {`,
		`			if err := m.Prop2[i].Validate(formats); err != nil {`,
		`				if ve, ok := err.(*errors.Validation); ok {`,
		`					return ve.ValidateName("prop2" + "." + strconv.Itoa(i)`,
		`func (m *ModelObjectVanilla) validateProp3(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop3) {`,
		`	if m.Prop3 != nil {`,
		`		if err := m.Prop3.Validate(formats); err != nil {`,
		`			if ve, ok := err.(*errors.Validation); ok {`,
		`				return ve.ValidateName("prop3"`,
		`func (m *ModelObjectVanilla) validateProp4(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop4) {`,
		`	for k := range m.Prop4 {`,
		`		if err := validate.Required("prop4"+"."+k, "body", m.Prop4[k]); err != nil {`,
		`		if val, ok := m.Prop4[k]; ok {`,
		`			if err := val.Validate(formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: model_sanity.go
	thisRun.AddExpectations("model_sanity.go", []string{
		`type ModelSanity struct {`,
		"	PropA string `json:\"propA,omitempty\"`",
		"	PropB *string `json:\"propB\"`",
		`func (m *ModelSanity) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validatePropB(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ModelSanity) validatePropB(formats strfmt.Registry) error {`,
		`	if err := validate.Required("propB", "body", m.PropB); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)
}

func initFixture1617() {
	// testing fixture-1617.yaml with flatten and expand (--skip-flatten)

	f := newModelFixture("../fixtures/bugs/1617/fixture-1617.yaml", "aaa")
	thisRun := f.AddRun(false).WithMinimalFlatten(true)

	// load expectations for model: artifact_info.go
	thisRun.AddExpectations("artifact_info.go", []string{
		`type ArtifactInfo struct {`,
		`	ArtifactDescription`,
		"	Path ArtifactPath `json:\"Path,omitempty\"`",
		"	Status ArtifactStatus `json:\"Status,omitempty\"`",
		"	Timestamp strfmt.DateTime `json:\"Timestamp,omitempty\"`",
		`func (m *ArtifactInfo) UnmarshalJSON(raw []byte) error {`,
		`	var aO0 ArtifactDescription`,
		`	if err := swag.ReadJSON(raw, &aO0); err != nil {`,
		`	m.ArtifactDescription = aO0`,
		`	var propsArtifactInfo struct {`,
		"		Path ArtifactPath `json:\"Path,omitempty\"`",
		"		Status ArtifactStatus `json:\"Status,omitempty\"`",
		"		Timestamp strfmt.DateTime `json:\"Timestamp,omitempty\"`",
		`	if err := swag.ReadJSON(raw, &propsArtifactInfo); err != nil {`,
		`	m.Path = propsArtifactInfo.Path`,
		`	m.Status = propsArtifactInfo.Status`,
		`	m.Timestamp = propsArtifactInfo.Timestamp`,
		`func (m ArtifactInfo) MarshalJSON() ([]byte, error) {`,
		`	_parts := make([][]byte, 0, 1`,
		`	aO0, err := swag.WriteJSON(m.ArtifactDescription`,
		`	if err != nil {`,
		`		return nil, err`,
		`	_parts = append(_parts, aO0`,
		`	var propsArtifactInfo struct {`,
		"		Path ArtifactPath `json:\"Path,omitempty\"`",
		"		Status ArtifactStatus `json:\"Status,omitempty\"`",
		"		Timestamp strfmt.DateTime `json:\"Timestamp,omitempty\"`",
		`	propsArtifactInfo.Path = m.Path`,
		`	propsArtifactInfo.Status = m.Status`,
		`	propsArtifactInfo.Timestamp = m.Timestamp`,
		`	jsonDataPropsArtifactInfo, errArtifactInfo := swag.WriteJSON(propsArtifactInfo`,
		`	if errArtifactInfo != nil {`,
		`		return nil, errArtifactInfo`,
		`	_parts = append(_parts, jsonDataPropsArtifactInfo`,
		`	return swag.ConcatJSON(_parts...), nil`,
		`func (m *ArtifactInfo) Validate(formats strfmt.Registry) error {`,
		`	if err := m.ArtifactDescription.Validate(formats); err != nil {`,
		`	if err := m.validatePath(formats); err != nil {`,
		`	if err := m.validateStatus(formats); err != nil {`,
		`	if err := m.validateTimestamp(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *ArtifactInfo) validatePath(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Path) {`,
		`	if err := m.Path.Validate(formats); err != nil {`,
		`		if ve, ok := err.(*errors.Validation); ok {`,
		`			return ve.ValidateName("Path"`,
		`func (m *ArtifactInfo) validateStatus(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Status) {`,
		`	if err := m.Status.Validate(formats); err != nil {`,
		`		if ve, ok := err.(*errors.Validation); ok {`,
		`			return ve.ValidateName("Status"`,
		`func (m *ArtifactInfo) validateTimestamp(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Timestamp) {`,
		`	if err := validate.FormatOf("Timestamp", "body", "date-time", m.Timestamp.String(), formats); err != nil {`,
		`func (m *ArtifactInfo) MarshalBinary() ([]byte, error) {`,
		`	if m == nil {`,
		`		return nil, nil`,
		`	return swag.WriteJSON(m`,
		`func (m *ArtifactInfo) UnmarshalBinary(b []byte) error {`,
		`	var res ArtifactInfo`,
		`	if err := swag.ReadJSON(b, &res); err != nil {`,
		`	*m = res`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)
}

func initFixtureRealiasedTypes() {
	/*
	   realiased types
	*/

	f := newModelFixture("../fixtures/bugs/1260/fixture-realiased-types.yaml", "test type realiasing")
	thisRun := f.AddRun(false).WithMinimalFlatten(true)

	// load expectations for model: g1.go
	thisRun.AddExpectations("g1.go", []string{
		`type G1 struct {`,
		"	Prop1 int64 `json:\"prop1,omitempty\"`",
		// empty validation
		"func (m *G1) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: e2v.go
	thisRun.AddExpectations("e2v.go", []string{
		`type E2v = E0v`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: a1v.go
	thisRun.AddExpectations("a1v.go", []string{
		`type A1v []int64`,
		`func (m A1v) Validate(formats strfmt.Registry) error {`,
		`	iA1vSize := int64(len(m)`,
		`	if err := validate.MaxItems("", "body", iA1vSize, 100); err != nil {`,
		`	for i := 0; i < len(m); i++ {`,
		`		if err := validate.MaximumInt(strconv.Itoa(i), "body", m[i], 100, false); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: f2v.go
	thisRun.AddExpectations("f2v.go", []string{
		`type F2v = F0v`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: hsubtype1.go
	thisRun.AddExpectations("hsubtype1.go", []string{
		`type Hsubtype1 struct {`,
		`	h1p1Field string`,
		`	h1p2Field strfmt.Date`,
		"	Hsp1 uint32 `json:\"hsp1,omitempty\"`",
		`func (m *Hsubtype1) H1p1() string {`,
		`	return m.h1p1Field`,
		`func (m *Hsubtype1) SetH1p1(val string) {`,
		`	m.h1p1Field = val`,
		`func (m *Hsubtype1) H1p2() strfmt.Date {`,
		`	return m.h1p2Field`,
		`func (m *Hsubtype1) SetH1p2(val strfmt.Date) {`,
		`	m.h1p2Field = val`,
		`func (m *Hsubtype1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateH1p2(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *Hsubtype1) validateH1p2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.H1p2()) {`,
		`	if err := validate.FormatOf("h1p2", "body", "date", m.H1p2().String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: f2.go
	thisRun.AddExpectations("f2.go", []string{
		`type F2 = F0`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: h0.go
	thisRun.AddExpectations("h0.go", []string{
		`type H0 = H1`,
		`func UnmarshalH0(reader io.Reader, consumer runtime.Consumer) (H0, error) {`,
		`	return UnmarshalH1(reader, consumer`,
		`func UnmarshalH0Slice(reader io.Reader, consumer runtime.Consumer) ([]H0, error) {`,
		`	return UnmarshalH1Slice(reader, consumer`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: c1v.go
	thisRun.AddExpectations("c1v.go", []string{
		`type C1v interface{`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: d0v.go
	thisRun.AddExpectations("d0v.go", []string{
		`type D0v = D1v`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: e2.go
	thisRun.AddExpectations("e2.go", []string{
		`type E2 = E0`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: d2v.go
	thisRun.AddExpectations("d2v.go", []string{
		`type D2v = D0v`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: b2v.go
	thisRun.AddExpectations("b2v.go", []string{
		`type B2v = B0v`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: a1.go
	thisRun.AddExpectations("a1.go", []string{
		`type A1 []int64`,
		// empty validation
		"func (m A1) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: hsubtype2.go
	thisRun.AddExpectations("hsubtype2.go", []string{
		`type Hsubtype2 struct {`,
		`	h1p1Field string`,
		`	h1p2Field strfmt.Date`,
		"	Hsp2 strfmt.DateTime `json:\"hsp2,omitempty\"`",
		`func (m *Hsubtype2) H1p1() string {`,
		`	return m.h1p1Field`,
		`func (m *Hsubtype2) SetH1p1(val string) {`,
		`	m.h1p1Field = val`,
		`func (m *Hsubtype2) H1p2() strfmt.Date {`,
		`	return m.h1p2Field`,
		`func (m *Hsubtype2) SetH1p2(val strfmt.Date) {`,
		`	m.h1p2Field = val`,
		`func (m *Hsubtype2) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateH1p2(formats); err != nil {`,
		`	if err := m.validateHsp2(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *Hsubtype2) validateH1p2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.H1p2()) {`,
		`	if err := validate.FormatOf("h1p2", "body", "date", m.H1p2().String(), formats); err != nil {`,
		`func (m *Hsubtype2) validateHsp2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Hsp2) {`,
		`	if err := validate.FormatOf("hsp2", "body", "date-time", m.Hsp2.String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: b2.go
	thisRun.AddExpectations("b2.go", []string{
		`type B2 = B0`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: b1.go
	thisRun.AddExpectations("b1.go", []string{
		`type B1 map[string]int64`,
		// empty validation
		"func (m B1) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: d0.go
	thisRun.AddExpectations("d0.go", []string{
		`type D0 = D1`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: g1v.go
	thisRun.AddExpectations("g1v.go", []string{
		`type G1v struct {`,
		"	Prop1v int64 `json:\"prop1v,omitempty\"`",
		"	Prop2v *int64 `json:\"prop2v\"`",
		`func (m *G1v) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateProp1v(formats); err != nil {`,
		`	if err := m.validateProp2v(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *G1v) validateProp1v(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.Prop1v) {`,
		`	if err := validate.MaximumInt("prop1v", "body", m.Prop1v, 100, false); err != nil {`,
		`func (m *G1v) validateProp2v(formats strfmt.Registry) error {`,
		`	if err := validate.Required("prop2v", "body", m.Prop2v); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: e0.go
	thisRun.AddExpectations("e0.go", []string{
		`type E0 = E1`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: f0v.go
	thisRun.AddExpectations("f0v.go", []string{
		`type F0v = F1v`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: h2.go
	thisRun.AddExpectations("h2.go", []string{
		`type H2 = H0`,
		`func UnmarshalH2(reader io.Reader, consumer runtime.Consumer) (H2, error) {`,
		`	return UnmarshalH0(reader, consumer`,
		`func UnmarshalH2Slice(reader io.Reader, consumer runtime.Consumer) ([]H2, error) {`,
		`	return UnmarshalH0Slice(reader, consumer`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: b1v.go
	thisRun.AddExpectations("b1v.go", []string{
		`type B1v map[string]int64`,
		`func (m B1v) Validate(formats strfmt.Registry) error {`,
		`	for k := range m {`,
		`		if err := validate.MaximumInt(k, "body", m[k], 100, false); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: f0.go
	thisRun.AddExpectations("f0.go", []string{
		`type F0 = F1`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: a2v.go
	thisRun.AddExpectations("a2v.go", []string{
		`type A2v = A0v`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: hs2.go
	thisRun.AddExpectations("hs2.go", []string{
		`type Hs2 struct {`,
		`	Hs0`,
		// empty validation
		"func (m *Hs2) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: d1v.go
	thisRun.AddExpectations("d1v.go", []string{
		`type D1v int64`,
		`func (m D1v) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.MaximumInt("", "body", int64(m), 100, false); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: c1.go
	thisRun.AddExpectations("c1.go", []string{
		`type C1 interface{`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: c2.go
	thisRun.AddExpectations("c2.go", []string{
		`type C2 = C0`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: e0v.go
	thisRun.AddExpectations("e0v.go", []string{
		`type E0v = E1v`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: d1.go
	thisRun.AddExpectations("d1.go", []string{
		`type D1 int64`,
		// empty validation
		"func (m D1) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: g2.go
	thisRun.AddExpectations("g2.go", []string{
		`type G2 struct {`,
		`	G0`,
		// empty validation
		"func (m *G2) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: d2.go
	thisRun.AddExpectations("d2.go", []string{
		`type D2 = D0`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: b0.go
	thisRun.AddExpectations("b0.go", []string{
		`type B0 = B1`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: c2v.go
	thisRun.AddExpectations("c2v.go", []string{
		`type C2v = C0v`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: g0.go
	thisRun.AddExpectations("g0.go", []string{
		`type G0 struct {`,
		`	G1`,
		// empty validation
		"func (m *G0) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: a0.go
	thisRun.AddExpectations("a0.go", []string{
		`type A0 = A1`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: a2.go
	thisRun.AddExpectations("a2.go", []string{
		`type A2 = A0`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: h1.go
	thisRun.AddExpectations("h1.go", []string{
		`type H1 interface {`,
		`	runtime.Validatable`,
		`	H1p1() string`,
		`	SetH1p1(string`,
		`	H1p2() strfmt.Date`,
		`	SetH1p2(strfmt.Date`,
		`type h1 struct {`,
		`	h1p1Field string`,
		`	h1p2Field strfmt.Date`,
		`func (m *h1) H1p1() string {`,
		`	return m.h1p1Field`,
		`func (m *h1) SetH1p1(val string) {`,
		`	m.h1p1Field = val`,
		`func (m *h1) H1p2() strfmt.Date {`,
		`	return m.h1p2Field`,
		`func (m *h1) SetH1p2(val strfmt.Date) {`,
		`	m.h1p2Field = val`,
		`func UnmarshalH1Slice(reader io.Reader, consumer runtime.Consumer) ([]H1, error) {`,
		`	var elements []json.RawMessage`,
		`	if err := consumer.Consume(reader, &elements); err != nil {`,
		`		return nil, err`,
		`	var result []H1`,
		`	for _, element := range elements {`,
		`		obj, err := unmarshalH1(element, consumer`,
		`		if err != nil {`,
		`			return nil, err`,
		`		result = append(result, obj`,
		`	return result, nil`,
		`func UnmarshalH1(reader io.Reader, consumer runtime.Consumer) (H1, error) {`,
		`	data, err := ioutil.ReadAll(reader`,
		`	if err != nil {`,
		`		return nil, err`,
		`	return unmarshalH1(data, consumer`,
		`func unmarshalH1(data []byte, consumer runtime.Consumer) (H1, error) {`,
		`	buf := bytes.NewBuffer(data`,
		`	buf2 := bytes.NewBuffer(data`,
		`	var getType struct {`,
		"		Htype string `json:\"htype\"`",
		`	if err := consumer.Consume(buf, &getType); err != nil {`,
		`		return nil, err`,
		`	if err := validate.RequiredString("htype", "body", getType.Htype); err != nil {`,
		`		return nil, err`,
		`	switch getType.Htype {`,
		`	case "h1":`,
		`		var result h1`,
		`		if err := consumer.Consume(buf2, &result); err != nil {`,
		`			return nil, err`,
		`		return &result, nil`,
		`	case "hsubtype1":`,
		`		var result Hsubtype1`,
		`		if err := consumer.Consume(buf2, &result); err != nil {`,
		`			return nil, err`,
		`		return &result, nil`,
		`	case "hsubtype2":`,
		`		var result Hsubtype2`,
		`		if err := consumer.Consume(buf2, &result); err != nil {`,
		`			return nil, err`,
		`		return &result, nil`,
		`	return nil, errors.New(422, "invalid htype value: %q", getType.Htype`,
		`func (m *h1) Validate(formats strfmt.Registry) error {`,
		`	if err := m.validateH1p2(formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
		`func (m *h1) validateH1p2(formats strfmt.Registry) error {`,
		`	if swag.IsZero(m.H1p2()) {`,
		`	if err := validate.FormatOf("h1p2", "body", "date", m.H1p2().String(), formats); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: c0v.go
	thisRun.AddExpectations("c0v.go", []string{
		`type C0v = C1v`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: g2v.go
	thisRun.AddExpectations("g2v.go", []string{
		`type G2v struct {`,
		`	G0v`,
		// empty validation
		"func (m *G2v) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: f1.go
	thisRun.AddExpectations("f1.go", []string{
		`type F1 strfmt.UUID`,
		`func (m F1) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.FormatOf("", "body", "uuid", strfmt.UUID(m).String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: e1.go
	thisRun.AddExpectations("e1.go", []string{
		`type E1 strfmt.Date`,
		`func (m E1) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.FormatOf("", "body", "date", strfmt.Date(m).String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: a0v.go
	thisRun.AddExpectations("a0v.go", []string{
		`type A0v = A1v`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: c0.go
	thisRun.AddExpectations("c0.go", []string{
		`type C0 = C1`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: f1v.go
	thisRun.AddExpectations("f1v.go", []string{
		`type F1v strfmt.UUID`,
		`func (m F1v) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.FormatOf("", "body", "uuid", strfmt.UUID(m).String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: e1v.go
	thisRun.AddExpectations("e1v.go", []string{
		`type E1v strfmt.Date`,
		`func (m E1v) Validate(formats strfmt.Registry) error {`,
		`	if err := validate.FormatOf("", "body", "date", strfmt.Date(m).String(), formats); err != nil {`,
		`		return errors.CompositeValidationError(res...`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: b0v.go
	thisRun.AddExpectations("b0v.go", []string{
		`type B0v = B1v`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: g0v.go
	thisRun.AddExpectations("g0v.go", []string{
		`type G0v struct {`,
		`	G1v`,
		// empty validation
		"func (m *G0v) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: hs0.go
	thisRun.AddExpectations("hs0.go", []string{
		`type Hs0 struct {`,
		`	Hsubtype1`,
		// empty validation
		"func (m *Hs0) Validate(formats strfmt.Registry) error {\n	return nil\n}",
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)
}

func initFixture1993() {
	/*
	   required / non required base type
	*/

	f := newModelFixture("../fixtures/bugs/1993/fixture-1993.yaml", "test required/non required base type")
	thisRun := f.AddRun(false).WithMinimalFlatten(true)

	// load expectations for model: house.go
	thisRun.AddExpectations("house.go", []string{
		`if err := validate.Required("pet", "body", m.Pet()); err != nil {`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)

	// load expectations for model: empty_house.go
	thisRun.AddExpectations("empty_house.go", []string{
		`if swag.IsZero(m.Pet())`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)
}

func initFixture2364() {
	f := newModelFixture("../fixtures/bugs/2364/fixture-2364.yaml", "test non-nullable allOf")
	thisRun := f.AddRun(false).WithMinimalFlatten(true)

	thisRun.AddExpectations("bundle_attributes_response.go", []string{
		`type BundleAttributesResponse struct {`,
		`Items []BundleItemResponse`,
		`Sections []ItemBundleSectionResponse`,
		`NullableSections []*NullableItemBundleSectionResponse`,
		`OtherSections []*OtherItemBundleSectionResponse`,
		`Type BundleType`,
	},
		// not expected
		todo,
		// output in log
		noLines,
		noLines)
}
