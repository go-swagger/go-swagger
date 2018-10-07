// BSON library for Go
//
// Copyright (c) 2010-2012 - Gustavo Niemeyer <gustavo@niemeyer.net>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// gobson - BSON library for Go.

package bson_test

import (
	"encoding/binary"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	. "gopkg.in/check.v1"
)

func TestAll(t *testing.T) {
	TestingT(t)
}

type S struct{}

var _ = Suite(&S{})

// Wrap up the document elements contained in data, prepending the int32
// length of the data, and appending the '\x00' value closing the document.
func wrapInDoc(data string) string {
	result := make([]byte, len(data)+5)
	binary.LittleEndian.PutUint32(result, uint32(len(result)))
	copy(result[4:], []byte(data))
	return string(result)
}

func makeZeroDoc(value interface{}) (zero interface{}) {
	v := reflect.ValueOf(value)
	t := v.Type()
	switch t.Kind() {
	case reflect.Map:
		mv := reflect.MakeMap(t)
		zero = mv.Interface()
	case reflect.Ptr:
		pv := reflect.New(v.Type().Elem())
		zero = pv.Interface()
	case reflect.Slice, reflect.Int, reflect.Int64, reflect.Struct:
		zero = reflect.New(t).Interface()
	default:
		panic("unsupported doc type: " + t.Name())
	}
	return zero
}

func testUnmarshal(c *C, data string, obj interface{}) {
	zero := makeZeroDoc(obj)
	err := bson.Unmarshal([]byte(data), zero)
	c.Assert(err, IsNil)
	c.Assert(zero, DeepEquals, obj)

	testUnmarshalRawElements(c, []byte(data))
}

func testUnmarshalRawElements(c *C, data []byte) {
	elems := []bson.RawDocElem{}
	err := bson.Unmarshal(data, &elems)
	c.Assert(err, IsNil)
	for _, elem := range elems {
		if elem.Value.Kind == bson.ElementDocument || elem.Value.Kind == bson.ElementArray {
			testUnmarshalRawElements(c, elem.Value.Data)
		}
	}
}

type testItemType struct {
	obj  interface{}
	data string
}

// --------------------------------------------------------------------------
// Samples from bsonspec.org:

var sampleItems = []testItemType{
	{bson.M{"hello": "world"},
		"\x16\x00\x00\x00\x02hello\x00\x06\x00\x00\x00world\x00\x00"},
	{bson.M{"BSON": []interface{}{"awesome", float64(5.05), 1986}},
		"1\x00\x00\x00\x04BSON\x00&\x00\x00\x00\x020\x00\x08\x00\x00\x00" +
			"awesome\x00\x011\x00333333\x14@\x102\x00\xc2\x07\x00\x00\x00\x00"},
	{bson.M{"slice": []uint8{1, 2}},
		"\x13\x00\x00\x00\x05slice\x00\x02\x00\x00\x00\x00\x01\x02\x00"},
	{bson.M{"slice": []byte{1, 2}},
		"\x13\x00\x00\x00\x05slice\x00\x02\x00\x00\x00\x00\x01\x02\x00"},
}

func (s *S) TestMarshalSampleItems(c *C) {
	for i, item := range sampleItems {
		data, err := bson.Marshal(item.obj)
		c.Assert(err, IsNil)
		c.Assert(string(data), Equals, item.data, Commentf("Failed on item %d", i))
	}
}

func (s *S) TestUnmarshalSampleItems(c *C) {
	for i, item := range sampleItems {
		value := bson.M{}
		err := bson.Unmarshal([]byte(item.data), value)
		c.Assert(err, IsNil)
		c.Assert(value, DeepEquals, item.obj, Commentf("Failed on item %d", i))
	}
}

// --------------------------------------------------------------------------
// Every type, ordered by the type flag. These are not wrapped with the
// length and last \x00 from the document. wrapInDoc() computes them.
// Note that all of them should be supported as two-way conversions.

var allItems = []testItemType{
	{bson.M{},
		""},
	{bson.M{"_": float64(5.05)},
		"\x01_\x00333333\x14@"},
	{bson.M{"_": "yo"},
		"\x02_\x00\x03\x00\x00\x00yo\x00"},
	{bson.M{"_": bson.M{"a": true}},
		"\x03_\x00\x09\x00\x00\x00\x08a\x00\x01\x00"},
	{bson.M{"_": []interface{}{true, false}},
		"\x04_\x00\r\x00\x00\x00\x080\x00\x01\x081\x00\x00\x00"},
	{bson.M{"_": []byte("yo")},
		"\x05_\x00\x02\x00\x00\x00\x00yo"},
	{bson.M{"_": bson.Binary{Kind: 0x80, Data: []byte("udef")}},
		"\x05_\x00\x04\x00\x00\x00\x80udef"},
	{bson.M{"_": bson.Undefined}, // Obsolete, but still seen in the wild.
		"\x06_\x00"},
	{bson.M{"_": bson.ObjectId("0123456789ab")},
		"\x07_\x000123456789ab"},
	{bson.M{"_": bson.DBPointer{Namespace: "testnamespace", Id: bson.ObjectId("0123456789ab")}},
		"\x0C_\x00\x0e\x00\x00\x00testnamespace\x000123456789ab"},
	{bson.M{"_": false},
		"\x08_\x00\x00"},
	{bson.M{"_": true},
		"\x08_\x00\x01"},
	{bson.M{"_": time.Unix(0, 258e6).UTC()}, // Note the NS <=> MS conversion.
		"\x09_\x00\x02\x01\x00\x00\x00\x00\x00\x00"},
	{bson.M{"_": nil},
		"\x0A_\x00"},
	{bson.M{"_": bson.RegEx{Pattern: "ab", Options: "cd"}},
		"\x0B_\x00ab\x00cd\x00"},
	{bson.M{"_": bson.JavaScript{Code: "code", Scope: nil}},
		"\x0D_\x00\x05\x00\x00\x00code\x00"},
	{bson.M{"_": bson.Symbol("sym")},
		"\x0E_\x00\x04\x00\x00\x00sym\x00"},
	{bson.M{"_": bson.JavaScript{Code: "code", Scope: bson.M{"": nil}}},
		"\x0F_\x00\x14\x00\x00\x00\x05\x00\x00\x00code\x00" +
			"\x07\x00\x00\x00\x0A\x00\x00"},
	{bson.M{"_": 258},
		"\x10_\x00\x02\x01\x00\x00"},
	{bson.M{"_": bson.MongoTimestamp(258)},
		"\x11_\x00\x02\x01\x00\x00\x00\x00\x00\x00"},
	{bson.M{"_": int64(258)},
		"\x12_\x00\x02\x01\x00\x00\x00\x00\x00\x00"},
	{bson.M{"_": int64(258 << 32)},
		"\x12_\x00\x00\x00\x00\x00\x02\x01\x00\x00"},
	{bson.M{"_": bson.MaxKey},
		"\x7F_\x00"},
	{bson.M{"_": bson.MinKey},
		"\xFF_\x00"},
}

func (s *S) TestMarshalAllItems(c *C) {
	for i, item := range allItems {
		data, err := bson.Marshal(item.obj)
		c.Assert(err, IsNil)
		c.Assert(string(data), Equals, wrapInDoc(item.data), Commentf("Failed on item %d: %#v", i, item))
	}
}

func (s *S) TestUnmarshalAllItems(c *C) {
	for i, item := range allItems {
		value := bson.M{}
		err := bson.Unmarshal([]byte(wrapInDoc(item.data)), value)
		c.Assert(err, IsNil)
		c.Assert(value, DeepEquals, item.obj, Commentf("Failed on item %d: %#v", i, item))
	}
}

func (s *S) TestUnmarshalRawAllItems(c *C) {
	for i, item := range allItems {
		if len(item.data) == 0 {
			continue
		}
		value := item.obj.(bson.M)["_"]
		if value == nil {
			continue
		}
		pv := reflect.New(reflect.ValueOf(value).Type())
		raw := bson.Raw{Kind: item.data[0], Data: []byte(item.data[3:])}
		c.Logf("Unmarshal raw: %#v, %#v", raw, pv.Interface())
		err := raw.Unmarshal(pv.Interface())
		c.Assert(err, IsNil)
		c.Assert(pv.Elem().Interface(), DeepEquals, value, Commentf("Failed on item %d: %#v", i, item))
	}
}

func (s *S) TestUnmarshalRawIncompatible(c *C) {
	raw := bson.Raw{Kind: 0x08, Data: []byte{0x01}} // true
	err := raw.Unmarshal(&struct{}{})
	c.Assert(err, ErrorMatches, "BSON kind 0x08 isn't compatible with type struct \\{\\}")
}

func (s *S) TestUnmarshalZeroesStruct(c *C) {
	data, err := bson.Marshal(bson.M{"b": 2})
	c.Assert(err, IsNil)
	type T struct{ A, B int }
	v := T{A: 1}
	err = bson.Unmarshal(data, &v)
	c.Assert(err, IsNil)
	c.Assert(v.A, Equals, 0)
	c.Assert(v.B, Equals, 2)
}

func (s *S) TestUnmarshalZeroesMap(c *C) {
	data, err := bson.Marshal(bson.M{"b": 2})
	c.Assert(err, IsNil)
	m := bson.M{"a": 1}
	err = bson.Unmarshal(data, &m)
	c.Assert(err, IsNil)
	c.Assert(m, DeepEquals, bson.M{"b": 2})
}

func (s *S) TestUnmarshalNonNilInterface(c *C) {
	data, err := bson.Marshal(bson.M{"b": 2})
	c.Assert(err, IsNil)
	m := bson.M{"a": 1}
	var i interface{}
	i = m
	err = bson.Unmarshal(data, &i)
	c.Assert(err, IsNil)
	c.Assert(i, DeepEquals, bson.M{"b": 2})
	c.Assert(m, DeepEquals, bson.M{"a": 1})
}

func (s *S) TestMarshalBuffer(c *C) {
	buf := make([]byte, 0, 256)
	data, err := bson.MarshalBuffer(bson.M{"a": 1}, buf)
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, buf[:len(data)])
}

func (s *S) TestPtrInline(c *C) {
	cases := []struct {
		In  interface{}
		Out bson.M
	}{
		{
			In:  inlinePtrStruct{A: 1, MStruct: &MStruct{M: 3}},
			Out: bson.M{"a": 1, "m": 3},
		},
		{ // go deeper
			In:  inlinePtrPtrStruct{B: 10, inlinePtrStruct: &inlinePtrStruct{A: 20, MStruct: &MStruct{M: 30}}},
			Out: bson.M{"b": 10, "a": 20, "m": 30},
		},
		{
			// nil embed struct
			In:  &inlinePtrStruct{A: 3},
			Out: bson.M{"a": 3},
		},
		{
			// nil embed struct
			In:  &inlinePtrPtrStruct{B: 5},
			Out: bson.M{"b": 5},
		},
	}

	for _, cs := range cases {
		data, err := bson.Marshal(cs.In)
		c.Assert(err, IsNil)
		var dataBSON bson.M
		err = bson.Unmarshal(data, &dataBSON)
		c.Assert(err, IsNil)

		c.Assert(dataBSON, DeepEquals, cs.Out)
	}
}

// --------------------------------------------------------------------------
// Some one way marshaling operations which would unmarshal differently.

var oneWayMarshalItems = []testItemType{
	// These are being passed as pointers, and will unmarshal as values.
	{bson.M{"": &bson.Binary{Kind: 0x02, Data: []byte("old")}},
		"\x05\x00\x07\x00\x00\x00\x02\x03\x00\x00\x00old"},
	{bson.M{"": &bson.Binary{Kind: 0x80, Data: []byte("udef")}},
		"\x05\x00\x04\x00\x00\x00\x80udef"},
	{bson.M{"": &bson.RegEx{Pattern: "ab", Options: "cd"}},
		"\x0B\x00ab\x00cd\x00"},
	{bson.M{"": &bson.JavaScript{Code: "code", Scope: nil}},
		"\x0D\x00\x05\x00\x00\x00code\x00"},
	{bson.M{"": &bson.JavaScript{Code: "code", Scope: bson.M{"": nil}}},
		"\x0F\x00\x14\x00\x00\x00\x05\x00\x00\x00code\x00" +
			"\x07\x00\x00\x00\x0A\x00\x00"},

	// There's no float32 type in BSON.  Will encode as a float64.
	{bson.M{"": float32(5.05)},
		"\x01\x00\x00\x00\x00@33\x14@"},

	// The array will be unmarshaled as a slice instead.
	{bson.M{"": [2]bool{true, false}},
		"\x04\x00\r\x00\x00\x00\x080\x00\x01\x081\x00\x00\x00"},

	// The typed slice will be unmarshaled as []interface{}.
	{bson.M{"": []bool{true, false}},
		"\x04\x00\r\x00\x00\x00\x080\x00\x01\x081\x00\x00\x00"},

	// Will unmarshal as a []byte.
	{bson.M{"": bson.Binary{Kind: 0x00, Data: []byte("yo")}},
		"\x05\x00\x02\x00\x00\x00\x00yo"},
	{bson.M{"": bson.Binary{Kind: 0x02, Data: []byte("old")}},
		"\x05\x00\x07\x00\x00\x00\x02\x03\x00\x00\x00old"},

	// No way to preserve the type information here. We might encode as a zero
	// value, but this would mean that pointer values in structs wouldn't be
	// able to correctly distinguish between unset and set to the zero value.
	{bson.M{"": (*byte)(nil)},
		"\x0A\x00"},

	// No int types smaller than int32 in BSON. Could encode this as a char,
	// but it would still be ambiguous, take more, and be awkward in Go when
	// loaded without typing information.
	{bson.M{"": byte(8)},
		"\x10\x00\x08\x00\x00\x00"},

	// There are no unsigned types in BSON.  Will unmarshal as int32 or int64.
	{bson.M{"": uint32(258)},
		"\x10\x00\x02\x01\x00\x00"},
	{bson.M{"": uint64(258)},
		"\x12\x00\x02\x01\x00\x00\x00\x00\x00\x00"},
	{bson.M{"": uint64(258 << 32)},
		"\x12\x00\x00\x00\x00\x00\x02\x01\x00\x00"},

	// This will unmarshal as int.
	{bson.M{"": int32(258)},
		"\x10\x00\x02\x01\x00\x00"},

	// That's a special case. The unsigned value is too large for an int32,
	// so an int64 is used instead.
	{bson.M{"": uint32(1<<32 - 1)},
		"\x12\x00\xFF\xFF\xFF\xFF\x00\x00\x00\x00"},
	{bson.M{"": uint(1<<32 - 1)},
		"\x12\x00\xFF\xFF\xFF\xFF\x00\x00\x00\x00"},
}

func (s *S) TestOneWayMarshalItems(c *C) {
	for i, item := range oneWayMarshalItems {
		data, err := bson.Marshal(item.obj)
		c.Assert(err, IsNil)
		c.Assert(string(data), Equals, wrapInDoc(item.data),
			Commentf("Failed on item %d", i))
	}
}

// --------------------------------------------------------------------------
// Some ops marshaling operations which would encode []uint8 or []byte in array.

var arrayOpsMarshalItems = []testItemType{
	{bson.M{"_": bson.M{"$in": []uint8{1, 2}}},
		"\x03_\x00\x1d\x00\x00\x00\x04$in\x00\x13\x00\x00\x00\x100\x00\x01\x00\x00\x00\x101\x00\x02\x00\x00\x00\x00\x00"},
	{bson.M{"_": bson.M{"$nin": []uint8{1, 2}}},
		"\x03_\x00\x1e\x00\x00\x00\x04$nin\x00\x13\x00\x00\x00\x100\x00\x01\x00\x00\x00\x101\x00\x02\x00\x00\x00\x00\x00"},
	{bson.M{"_": bson.M{"$all": []uint8{1, 2}}},
		"\x03_\x00\x1e\x00\x00\x00\x04$all\x00\x13\x00\x00\x00\x100\x00\x01\x00\x00\x00\x101\x00\x02\x00\x00\x00\x00\x00"},
}

func (s *S) TestArrayOpsMarshalItems(c *C) {
	for i, item := range arrayOpsMarshalItems {
		data, err := bson.Marshal(item.obj)
		c.Assert(err, IsNil)
		c.Assert(string(data), Equals, wrapInDoc(item.data),
			Commentf("Failed on item %d", i))
	}
}

// --------------------------------------------------------------------------
// Two-way tests for user-defined structures using the samples
// from bsonspec.org.

type specSample1 struct {
	Hello string
}

type specSample2 struct {
	BSON []interface{} `bson:"BSON"`
}

var structSampleItems = []testItemType{
	{&specSample1{"world"},
		"\x16\x00\x00\x00\x02hello\x00\x06\x00\x00\x00world\x00\x00"},
	{&specSample2{[]interface{}{"awesome", float64(5.05), 1986}},
		"1\x00\x00\x00\x04BSON\x00&\x00\x00\x00\x020\x00\x08\x00\x00\x00" +
			"awesome\x00\x011\x00333333\x14@\x102\x00\xc2\x07\x00\x00\x00\x00"},
}

func (s *S) TestMarshalStructSampleItems(c *C) {
	for i, item := range structSampleItems {
		data, err := bson.Marshal(item.obj)
		c.Assert(err, IsNil)
		c.Assert(string(data), Equals, item.data,
			Commentf("Failed on item %d", i))
	}
}

func (s *S) TestUnmarshalStructSampleItems(c *C) {
	for _, item := range structSampleItems {
		testUnmarshal(c, item.data, item.obj)
	}
}

func (s *S) Test64bitInt(c *C) {
	var i int64 = (1 << 31)
	if int(i) > 0 {
		data, err := bson.Marshal(bson.M{"i": int(i)})
		c.Assert(err, IsNil)
		c.Assert(string(data), Equals, wrapInDoc("\x12i\x00\x00\x00\x00\x80\x00\x00\x00\x00"))

		var result struct{ I int }
		err = bson.Unmarshal(data, &result)
		c.Assert(err, IsNil)
		c.Assert(int64(result.I), Equals, i)
	}
}

// --------------------------------------------------------------------------
// Generic two-way struct marshaling tests.

type prefixPtr string
type prefixVal string

func (t *prefixPtr) GetBSON() (interface{}, error) {
	if t == nil {
		return nil, nil
	}
	return "foo-" + string(*t), nil
}

func (t *prefixPtr) SetBSON(raw bson.Raw) error {
	var s string
	if raw.Kind == 0x0A {
		return bson.ErrSetZero
	}
	if err := raw.Unmarshal(&s); err != nil {
		return err
	}
	if !strings.HasPrefix(s, "foo-") {
		return errors.New("Prefix not found: " + s)
	}
	*t = prefixPtr(s[4:])
	return nil
}

func (t prefixVal) GetBSON() (interface{}, error) {
	return "foo-" + string(t), nil
}

func (t *prefixVal) SetBSON(raw bson.Raw) error {
	var s string
	if raw.Kind == 0x0A {
		return bson.ErrSetZero
	}
	if err := raw.Unmarshal(&s); err != nil {
		return err
	}
	if !strings.HasPrefix(s, "foo-") {
		return errors.New("Prefix not found: " + s)
	}
	*t = prefixVal(s[4:])
	return nil
}

var bytevar = byte(8)
var byteptr = &bytevar
var prefixptr = prefixPtr("bar")
var prefixval = prefixVal("bar")

var structItems = []testItemType{
	{&struct{ Ptr *byte }{nil},
		"\x0Aptr\x00"},
	{&struct{ Ptr *byte }{&bytevar},
		"\x10ptr\x00\x08\x00\x00\x00"},
	{&struct{ Ptr **byte }{&byteptr},
		"\x10ptr\x00\x08\x00\x00\x00"},
	{&struct{ Byte byte }{8},
		"\x10byte\x00\x08\x00\x00\x00"},
	{&struct{ Byte byte }{0},
		"\x10byte\x00\x00\x00\x00\x00"},
	{&struct {
		V byte `bson:"Tag"`
	}{8},
		"\x10Tag\x00\x08\x00\x00\x00"},
	{&struct {
		V *struct {
			Byte byte
		}
	}{&struct{ Byte byte }{8}},
		"\x03v\x00" + "\x0f\x00\x00\x00\x10byte\x00\b\x00\x00\x00\x00"},
	{&struct{ priv byte }{}, ""},

	// The order of the dumped fields should be the same in the struct.
	{&struct{ A, C, B, D, F, E *byte }{},
		"\x0Aa\x00\x0Ac\x00\x0Ab\x00\x0Ad\x00\x0Af\x00\x0Ae\x00"},

	{&struct{ V bson.Raw }{bson.Raw{Kind: 0x03, Data: []byte("\x0f\x00\x00\x00\x10byte\x00\b\x00\x00\x00\x00")}},
		"\x03v\x00" + "\x0f\x00\x00\x00\x10byte\x00\b\x00\x00\x00\x00"},
	{&struct{ V bson.Raw }{bson.Raw{Kind: 0x10, Data: []byte("\x00\x00\x00\x00")}},
		"\x10v\x00" + "\x00\x00\x00\x00"},

	// Byte arrays.
	{&struct{ V [2]byte }{[2]byte{'y', 'o'}},
		"\x05v\x00\x02\x00\x00\x00\x00yo"},

	{&struct{ V prefixPtr }{prefixPtr("buzz")},
		"\x02v\x00\x09\x00\x00\x00foo-buzz\x00"},

	{&struct{ V *prefixPtr }{&prefixptr},
		"\x02v\x00\x08\x00\x00\x00foo-bar\x00"},

	{&struct{ V *prefixPtr }{nil},
		"\x0Av\x00"},

	{&struct{ V prefixVal }{prefixVal("buzz")},
		"\x02v\x00\x09\x00\x00\x00foo-buzz\x00"},

	{&struct{ V *prefixVal }{&prefixval},
		"\x02v\x00\x08\x00\x00\x00foo-bar\x00"},

	{&struct{ V *prefixVal }{nil},
		"\x0Av\x00"},
}

func (s *S) TestMarshalStructItems(c *C) {
	for i, item := range structItems {
		data, err := bson.Marshal(item.obj)
		c.Assert(err, IsNil)
		c.Assert(string(data), Equals, wrapInDoc(item.data),
			Commentf("Failed on item %d", i))
	}
}

func (s *S) TestUnmarshalStructItems(c *C) {
	for _, item := range structItems {
		testUnmarshal(c, wrapInDoc(item.data), item.obj)
	}
}

func (s *S) TestUnmarshalRawStructItems(c *C) {
	for i, item := range structItems {
		raw := bson.Raw{Kind: 0x03, Data: []byte(wrapInDoc(item.data))}
		zero := makeZeroDoc(item.obj)
		err := raw.Unmarshal(zero)
		c.Assert(err, IsNil)
		c.Assert(zero, DeepEquals, item.obj, Commentf("Failed on item %d: %#v", i, item))
	}
}

func (s *S) TestUnmarshalRawNil(c *C) {
	// Regression test: shouldn't try to nil out the pointer itself,
	// as it's not settable.
	raw := bson.Raw{Kind: 0x0A, Data: []byte{}}
	err := raw.Unmarshal(&struct{}{})
	c.Assert(err, IsNil)
}

// --------------------------------------------------------------------------
// One-way marshaling tests.

type dOnIface struct {
	D interface{}
}

type ignoreField struct {
	Before string
	Ignore string `bson:"-"`
	After  string
}

var marshalItems = []testItemType{
	// Ordered document dump.  Will unmarshal as a dictionary by default.
	{bson.D{{Name: "a", Value: nil}, {Name: "c", Value: nil}, {Name: "b", Value: nil}, {Name: "d", Value: nil}, {Name: "f", Value: nil}, {Name: "e", Value: true}},
		"\x0Aa\x00\x0Ac\x00\x0Ab\x00\x0Ad\x00\x0Af\x00\x08e\x00\x01"},
	{MyD{{Name: "a", Value: nil}, {Name: "c", Value: nil}, {Name: "b", Value: nil}, {Name: "d", Value: nil}, {Name: "f", Value: nil}, {Name: "e", Value: true}},
		"\x0Aa\x00\x0Ac\x00\x0Ab\x00\x0Ad\x00\x0Af\x00\x08e\x00\x01"},
	{&dOnIface{bson.D{{Name: "a", Value: nil}, {Name: "c", Value: nil}, {Name: "b", Value: nil}, {Name: "d", Value: true}}},
		"\x03d\x00" + wrapInDoc("\x0Aa\x00\x0Ac\x00\x0Ab\x00\x08d\x00\x01")},

	{bson.RawD{{Name: "a", Value: bson.Raw{Kind: 0x0A, Data: nil}}, {Name: "c", Value: bson.Raw{Kind: 0x0A, Data: nil}}, {Name: "b", Value: bson.Raw{Kind: 0x08, Data: []byte{0x01}}}},
		"\x0Aa\x00" + "\x0Ac\x00" + "\x08b\x00\x01"},
	{MyRawD{{Name: "a", Value: bson.Raw{Kind: 0x0A, Data: nil}}, {Name: "c", Value: bson.Raw{Kind: 0x0A, Data: nil}}, {Name: "b", Value: bson.Raw{Kind: 0x08, Data: []byte{0x01}}}},
		"\x0Aa\x00" + "\x0Ac\x00" + "\x08b\x00\x01"},
	{&dOnIface{bson.RawD{{Name: "a", Value: bson.Raw{Kind: 0x0A, Data: nil}}, {Name: "c", Value: bson.Raw{Kind: 0x0A, Data: nil}}, {Name: "b", Value: bson.Raw{Kind: 0x08, Data: []byte{0x01}}}}},
		"\x03d\x00" + wrapInDoc("\x0Aa\x00"+"\x0Ac\x00"+"\x08b\x00\x01")},

	{&ignoreField{"before", "ignore", "after"},
		"\x02before\x00\a\x00\x00\x00before\x00\x02after\x00\x06\x00\x00\x00after\x00"},

	// Marshalling a Raw document does nothing.
	{bson.Raw{Kind: 0x03, Data: []byte(wrapInDoc("anything"))},
		"anything"},
	{bson.Raw{Data: []byte(wrapInDoc("anything"))},
		"anything"},
}

func (s *S) TestMarshalOneWayItems(c *C) {
	for _, item := range marshalItems {
		data, err := bson.Marshal(item.obj)
		c.Assert(err, IsNil)
		c.Assert(string(data), Equals, wrapInDoc(item.data))
	}
}

// --------------------------------------------------------------------------
// One-way unmarshaling tests.

type intAlias int

var unmarshalItems = []testItemType{
	// Field is private.  Should not attempt to unmarshal it.
	{&struct{ priv byte }{},
		"\x10priv\x00\x08\x00\x00\x00"},

	// Wrong casing. Field names are lowercased.
	{&struct{ Byte byte }{},
		"\x10Byte\x00\x08\x00\x00\x00"},

	// Ignore non-existing field.
	{&struct{ Byte byte }{9},
		"\x10boot\x00\x08\x00\x00\x00" + "\x10byte\x00\x09\x00\x00\x00"},

	// Do not unmarshal on ignored field.
	{&ignoreField{"before", "", "after"},
		"\x02before\x00\a\x00\x00\x00before\x00" +
			"\x02-\x00\a\x00\x00\x00ignore\x00" +
			"\x02after\x00\x06\x00\x00\x00after\x00"},

	// Ignore unsuitable types silently.
	{map[string]string{"str": "s"},
		"\x02str\x00\x02\x00\x00\x00s\x00" + "\x10int\x00\x01\x00\x00\x00"},
	{map[string][]int{"array": {5, 9}},
		"\x04array\x00" + wrapInDoc("\x100\x00\x05\x00\x00\x00"+"\x021\x00\x02\x00\x00\x00s\x00"+"\x102\x00\x09\x00\x00\x00")},

	// Wrong type. Shouldn't init pointer.
	{&struct{ Str *byte }{},
		"\x02str\x00\x02\x00\x00\x00s\x00"},
	{&struct{ Str *struct{ Str string } }{},
		"\x02str\x00\x02\x00\x00\x00s\x00"},

	// Ordered document.
	{&struct{ bson.D }{bson.D{{Name: "a", Value: nil}, {Name: "c", Value: nil}, {Name: "b", Value: nil}, {Name: "d", Value: true}}},
		"\x03d\x00" + wrapInDoc("\x0Aa\x00\x0Ac\x00\x0Ab\x00\x08d\x00\x01")},

	// Raw document.
	{&bson.Raw{Kind: 0x03, Data: []byte(wrapInDoc("\x10byte\x00\x08\x00\x00\x00"))},
		"\x10byte\x00\x08\x00\x00\x00"},

	// RawD document.
	{&struct{ bson.RawD }{bson.RawD{{Name: "a", Value: bson.Raw{Kind: 0x0A, Data: []byte{}}}, {Name: "c", Value: bson.Raw{Kind: 0x0A, Data: []byte{}}}, {Name: "b", Value: bson.Raw{Kind: 0x08, Data: []byte{0x01}}}}},
		"\x03rawd\x00" + wrapInDoc("\x0Aa\x00\x0Ac\x00\x08b\x00\x01")},

	// Decode old binary.
	{bson.M{"_": []byte("old")},
		"\x05_\x00\x07\x00\x00\x00\x02\x03\x00\x00\x00old"},

	// Decode old binary without length. According to the spec, this shouldn't happen.
	{bson.M{"_": []byte("old")},
		"\x05_\x00\x03\x00\x00\x00\x02old"},

	// Decode a doc within a doc in to a slice within a doc; shouldn't error
	{&struct{ Foo []string }{},
		"\x03\x66\x6f\x6f\x00\x05\x00\x00\x00\x00"},

	// int key maps
	{map[int]string{10: "s"},
		"\x0210\x00\x02\x00\x00\x00s\x00"},

	//// event if type is alias to int
	{map[intAlias]string{10: "s"},
		"\x0210\x00\x02\x00\x00\x00s\x00"},
}

func (s *S) TestUnmarshalOneWayItems(c *C) {
	for _, item := range unmarshalItems {
		testUnmarshal(c, wrapInDoc(item.data), item.obj)
	}
}

func (s *S) TestUnmarshalNilInStruct(c *C) {
	// Nil is the default value, so we need to ensure it's indeed being set.
	b := byte(1)
	v := &struct{ Ptr *byte }{&b}
	err := bson.Unmarshal([]byte(wrapInDoc("\x0Aptr\x00")), v)
	c.Assert(err, IsNil)
	c.Assert(v, DeepEquals, &struct{ Ptr *byte }{nil})
}

// --------------------------------------------------------------------------
// Marshalling error cases.

type structWithDupKeys struct {
	Name  byte
	Other byte `bson:"name"` // Tag should precede.
}

var marshalErrorItems = []testItemType{
	{bson.M{"": uint64(1 << 63)},
		"BSON has no uint64 type, and value is too large to fit correctly in an int64"},
	{bson.M{"": bson.ObjectId("tooshort")},
		"ObjectIDs must be exactly 12 bytes long \\(got 8\\)"},
	{int64(123),
		"Can't marshal int64 as a BSON document"},
	{bson.M{"": 1i},
		"Can't marshal complex128 in a BSON document"},
	{&structWithDupKeys{},
		"Duplicated key 'name' in struct bson_test.structWithDupKeys"},
	{bson.Raw{Kind: 0xA, Data: []byte{}},
		"Attempted to marshal Raw kind 10 as a document"},
	{bson.Raw{Kind: 0x3, Data: []byte{}},
		"Attempted to marshal empty Raw document"},
	{bson.M{"w": bson.Raw{Kind: 0x3, Data: []byte{}}},
		"Attempted to marshal empty Raw document"},
	{&inlineDupName{1, struct{ A, B int }{2, 3}},
		"Duplicated key 'a' in struct bson_test.inlineDupName"},
	{&inlineDupMap{},
		"Multiple ,inline maps in struct bson_test.inlineDupMap"},
	{&inlineBadKeyMap{},
		"Option ,inline needs a map with string keys in struct bson_test.inlineBadKeyMap"},
	{&inlineMap{A: 1, M: map[string]interface{}{"a": 1}},
		`Can't have key "a" in inlined map; conflicts with struct field`},
}

func (s *S) TestMarshalErrorItems(c *C) {
	for _, item := range marshalErrorItems {
		data, err := bson.Marshal(item.obj)
		c.Assert(err, ErrorMatches, item.data)
		c.Assert(data, IsNil)
	}
}

// --------------------------------------------------------------------------
// Unmarshalling error cases.

type unmarshalErrorType struct {
	obj   interface{}
	data  string
	error string
}

var unmarshalErrorItems = []unmarshalErrorType{
	// Tag name conflicts with existing parameter.
	{&structWithDupKeys{},
		"\x10name\x00\x08\x00\x00\x00",
		"Duplicated key 'name' in struct bson_test.structWithDupKeys"},

	{nil,
		"\xEEname\x00",
		"Unknown element kind \\(0xEE\\)"},

	{struct{ Name bool }{},
		"\x10name\x00\x08\x00\x00\x00",
		"unmarshal can't deal with struct values. Use a pointer"},

	{123,
		"\x10name\x00\x08\x00\x00\x00",
		"unmarshal needs a map or a pointer to a struct"},

	{nil,
		"\x08\x62\x00\x02",
		"encoded boolean must be 1 or 0, found 2"},

	// Non-string and not numeric map key.
	{map[bool]interface{}{true: 1},
		"\x10true\x00\x01\x00\x00\x00",
		"BSON map must have string or decimal keys. Got: map\\[bool\\]interface \\{\\}"},
}

func (s *S) TestUnmarshalErrorItems(c *C) {
	for _, item := range unmarshalErrorItems {
		data := []byte(wrapInDoc(item.data))
		var value interface{}
		switch reflect.ValueOf(item.obj).Kind() {
		case reflect.Map, reflect.Ptr:
			value = makeZeroDoc(item.obj)
		case reflect.Invalid:
			value = bson.M{}
		default:
			value = item.obj
		}
		err := bson.Unmarshal(data, value)
		c.Assert(err, ErrorMatches, item.error)
	}
}

type unmarshalRawErrorType struct {
	obj   interface{}
	raw   bson.Raw
	error string
}

var unmarshalRawErrorItems = []unmarshalRawErrorType{
	// Tag name conflicts with existing parameter.
	{&structWithDupKeys{},
		bson.Raw{Kind: 0x03, Data: []byte("\x10byte\x00\x08\x00\x00\x00")},
		"Duplicated key 'name' in struct bson_test.structWithDupKeys"},

	{&struct{}{},
		bson.Raw{Kind: 0xEE, Data: []byte{}},
		"Unknown element kind \\(0xEE\\)"},

	{struct{ Name bool }{},
		bson.Raw{Kind: 0x10, Data: []byte("\x08\x00\x00\x00")},
		"raw Unmarshal can't deal with struct values. Use a pointer"},

	{123,
		bson.Raw{Kind: 0x10, Data: []byte("\x08\x00\x00\x00")},
		"raw Unmarshal needs a map or a valid pointer"},
}

func (s *S) TestUnmarshalRawErrorItems(c *C) {
	for i, item := range unmarshalRawErrorItems {
		err := item.raw.Unmarshal(item.obj)
		c.Assert(err, ErrorMatches, item.error, Commentf("Failed on item %d: %#v\n", i, item))
	}
}

var corruptedData = []string{
	"\x04\x00\x00\x00\x00",         // Document shorter than minimum
	"\x06\x00\x00\x00\x00",         // Not enough data
	"\x05\x00\x00",                 // Broken length
	"\x05\x00\x00\x00\xff",         // Corrupted termination
	"\x0A\x00\x00\x00\x0Aooop\x00", // Unfinished C string

	// Array end past end of string (s[2]=0x07 is correct)
	wrapInDoc("\x04\x00\x09\x00\x00\x00\x0A\x00\x00"),

	// Array end within string, but past acceptable.
	wrapInDoc("\x04\x00\x08\x00\x00\x00\x0A\x00\x00"),

	// Document end within string, but past acceptable.
	wrapInDoc("\x03\x00\x08\x00\x00\x00\x0A\x00\x00"),

	// String with corrupted end.
	wrapInDoc("\x02\x00\x03\x00\x00\x00yo\xFF"),

	// String with negative length (issue #116).
	"\x0c\x00\x00\x00\x02x\x00\xff\xff\xff\xff\x00",

	// String with zero length (must include trailing '\x00')
	"\x0c\x00\x00\x00\x02x\x00\x00\x00\x00\x00\x00",

	// Binary with negative length.
	"\r\x00\x00\x00\x05x\x00\xff\xff\xff\xff\x00\x00",
}

func (s *S) TestUnmarshalMapDocumentTooShort(c *C) {
	for _, data := range corruptedData {
		err := bson.Unmarshal([]byte(data), bson.M{})
		c.Assert(err, ErrorMatches, "Document is corrupted")

		err = bson.Unmarshal([]byte(data), &struct{}{})
		c.Assert(err, ErrorMatches, "Document is corrupted")
	}
}

// --------------------------------------------------------------------------
// Setter test cases.

var setterResult = map[string]error{}

type setterType struct {
	received interface{}
}

func (o *setterType) SetBSON(raw bson.Raw) error {
	err := raw.Unmarshal(&o.received)
	if err != nil {
		panic("The panic:" + err.Error())
	}
	if s, ok := o.received.(string); ok {
		if result, ok := setterResult[s]; ok {
			return result
		}
	}
	return nil
}

type ptrSetterDoc struct {
	Field *setterType `bson:"_"`
}

type valSetterDoc struct {
	Field setterType `bson:"_"`
}

func (s *S) TestUnmarshalAllItemsWithPtrSetter(c *C) {
	for _, item := range allItems {
		for i := 0; i != 2; i++ {
			var field *setterType
			if i == 0 {
				obj := &ptrSetterDoc{}
				err := bson.Unmarshal([]byte(wrapInDoc(item.data)), obj)
				c.Assert(err, IsNil)
				field = obj.Field
			} else {
				obj := &valSetterDoc{}
				err := bson.Unmarshal([]byte(wrapInDoc(item.data)), obj)
				c.Assert(err, IsNil)
				field = &obj.Field
			}
			if item.data == "" {
				// Nothing to unmarshal. Should be untouched.
				if i == 0 {
					c.Assert(field, IsNil)
				} else {
					c.Assert(field.received, IsNil)
				}
			} else {
				expected := item.obj.(bson.M)["_"]
				c.Assert(field, NotNil, Commentf("Pointer not initialized (%#v)", expected))
				c.Assert(field.received, DeepEquals, expected)
			}
		}
	}
}

func (s *S) TestUnmarshalWholeDocumentWithSetter(c *C) {
	obj := &setterType{}
	err := bson.Unmarshal([]byte(sampleItems[0].data), obj)
	c.Assert(err, IsNil)
	c.Assert(obj.received, DeepEquals, bson.M{"hello": "world"})
}

func (s *S) TestUnmarshalSetterOmits(c *C) {
	setterResult["2"] = &bson.TypeError{}
	setterResult["4"] = &bson.TypeError{}
	defer func() {
		delete(setterResult, "2")
		delete(setterResult, "4")
	}()

	m := map[string]*setterType{}
	data := wrapInDoc("\x02abc\x00\x02\x00\x00\x001\x00" +
		"\x02def\x00\x02\x00\x00\x002\x00" +
		"\x02ghi\x00\x02\x00\x00\x003\x00" +
		"\x02jkl\x00\x02\x00\x00\x004\x00")
	err := bson.Unmarshal([]byte(data), m)
	c.Assert(err, IsNil)
	c.Assert(m["abc"], NotNil)
	c.Assert(m["def"], IsNil)
	c.Assert(m["ghi"], NotNil)
	c.Assert(m["jkl"], IsNil)

	c.Assert(m["abc"].received, Equals, "1")
	c.Assert(m["ghi"].received, Equals, "3")
}

func (s *S) TestUnmarshalSetterErrors(c *C) {
	boom := errors.New("BOOM")
	setterResult["2"] = boom
	defer delete(setterResult, "2")

	m := map[string]*setterType{}
	data := wrapInDoc("\x02abc\x00\x02\x00\x00\x001\x00" +
		"\x02def\x00\x02\x00\x00\x002\x00" +
		"\x02ghi\x00\x02\x00\x00\x003\x00")
	err := bson.Unmarshal([]byte(data), m)
	c.Assert(err, Equals, boom)
	c.Assert(m["abc"], NotNil)
	c.Assert(m["def"], IsNil)
	c.Assert(m["ghi"], IsNil)

	c.Assert(m["abc"].received, Equals, "1")
}

func (s *S) TestDMap(c *C) {
	d := bson.D{{Name: "a", Value: 1}, {Name: "b", Value: 2}}
	c.Assert(d.Map(), DeepEquals, bson.M{"a": 1, "b": 2})
}

func (s *S) TestUnmarshalSetterErrSetZero(c *C) {
	setterResult["foo"] = bson.ErrSetZero
	defer delete(setterResult, "field")

	data, err := bson.Marshal(bson.M{"field": "foo"})
	c.Assert(err, IsNil)

	m := map[string]*setterType{}
	err = bson.Unmarshal([]byte(data), m)
	c.Assert(err, IsNil)

	value, ok := m["field"]
	c.Assert(ok, Equals, true)
	c.Assert(value, IsNil)
}

// --------------------------------------------------------------------------
// Getter test cases.

type typeWithGetter struct {
	result interface{}
	err    error
}

func (t *typeWithGetter) GetBSON() (interface{}, error) {
	if t == nil {
		return "<value is nil>", nil
	}
	return t.result, t.err
}

type docWithGetterField struct {
	Field *typeWithGetter `bson:"_"`
}

func (s *S) TestMarshalAllItemsWithGetter(c *C) {
	for i, item := range allItems {
		if item.data == "" {
			continue
		}
		obj := &docWithGetterField{}
		obj.Field = &typeWithGetter{result: item.obj.(bson.M)["_"]}
		data, err := bson.Marshal(obj)
		c.Assert(err, IsNil)
		c.Assert(string(data), Equals, wrapInDoc(item.data),
			Commentf("Failed on item #%d", i))
	}
}

func (s *S) TestMarshalWholeDocumentWithGetter(c *C) {
	obj := &typeWithGetter{result: sampleItems[0].obj}
	data, err := bson.Marshal(obj)
	c.Assert(err, IsNil)
	c.Assert(string(data), Equals, sampleItems[0].data)
}

func (s *S) TestGetterErrors(c *C) {
	e := errors.New("oops")

	obj1 := &docWithGetterField{}
	obj1.Field = &typeWithGetter{sampleItems[0].obj, e}
	data, err := bson.Marshal(obj1)
	c.Assert(err, ErrorMatches, "oops")
	c.Assert(data, IsNil)

	obj2 := &typeWithGetter{sampleItems[0].obj, e}
	data, err = bson.Marshal(obj2)
	c.Assert(err, ErrorMatches, "oops")
	c.Assert(data, IsNil)
}

type intGetter int64

func (t intGetter) GetBSON() (interface{}, error) {
	return int64(t), nil
}

type typeWithIntGetter struct {
	V intGetter `bson:",minsize"`
}

func (s *S) TestMarshalShortWithGetter(c *C) {
	obj := typeWithIntGetter{42}
	data, err := bson.Marshal(obj)
	c.Assert(err, IsNil)
	m := bson.M{}
	err = bson.Unmarshal(data, m)
	c.Assert(err, IsNil)
	c.Assert(m["v"], Equals, 42)
}

func (s *S) TestMarshalWithGetterNil(c *C) {
	obj := docWithGetterField{}
	data, err := bson.Marshal(obj)
	c.Assert(err, IsNil)
	m := bson.M{}
	err = bson.Unmarshal(data, m)
	c.Assert(err, IsNil)
	c.Assert(m, DeepEquals, bson.M{"_": "<value is nil>"})
}

// --------------------------------------------------------------------------
// Cross-type conversion tests.

type crossTypeItem struct {
	obj1 interface{}
	obj2 interface{}
}

type condStr struct {
	V string `bson:",omitempty"`
}
type condStrNS struct {
	V string `a:"A" bson:",omitempty" b:"B"`
}
type condBool struct {
	V bool `bson:",omitempty"`
}
type condInt struct {
	V int `bson:",omitempty"`
}
type condUInt struct {
	V uint `bson:",omitempty"`
}
type condFloat struct {
	V float64 `bson:",omitempty"`
}
type condIface struct {
	V interface{} `bson:",omitempty"`
}
type condPtr struct {
	V *bool `bson:",omitempty"`
}
type condSlice struct {
	V []string `bson:",omitempty"`
}
type condMap struct {
	V map[string]int `bson:",omitempty"`
}
type namedCondStr struct {
	V string `bson:"myv,omitempty"`
}
type condTime struct {
	V time.Time `bson:",omitempty"`
}
type condStruct struct {
	V struct{ A []int } `bson:",omitempty"`
}
type condRaw struct {
	V bson.Raw `bson:",omitempty"`
}

type shortInt struct {
	V int64 `bson:",minsize"`
}
type shortUint struct {
	V uint64 `bson:",minsize"`
}
type shortIface struct {
	V interface{} `bson:",minsize"`
}
type shortPtr struct {
	V *int64 `bson:",minsize"`
}
type shortNonEmptyInt struct {
	V int64 `bson:",minsize,omitempty"`
}

type inlineInt struct {
	V struct{ A, B int } `bson:",inline"`
}
type inlineCantPtr struct {
	V *struct{ A, B int } `bson:",inline"`
}
type inlineDupName struct {
	A int
	V struct{ A, B int } `bson:",inline"`
}
type inlineMap struct {
	A int
	M map[string]interface{} `bson:",inline"`
}
type inlineMapInt struct {
	A int
	M map[string]int `bson:",inline"`
}
type inlineMapMyM struct {
	A int
	M MyM `bson:",inline"`
}
type inlineDupMap struct {
	M1 map[string]interface{} `bson:",inline"`
	M2 map[string]interface{} `bson:",inline"`
}
type inlineBadKeyMap struct {
	M map[int]int `bson:",inline"`
}
type inlineUnexported struct {
	M          map[string]interface{} `bson:",inline"`
	unexported `bson:",inline"`
}
type MStruct struct {
	M int `bson:"m,omitempty"`
}
type inlinePtrStruct struct {
	A        int
	*MStruct `bson:",inline"`
}
type inlinePtrPtrStruct struct {
	B                int
	*inlinePtrStruct `bson:",inline"`
}
type unexported struct {
	A int
}

type getterSetterD bson.D

func (s getterSetterD) GetBSON() (interface{}, error) {
	if len(s) == 0 {
		return bson.D{}, nil
	}
	return bson.D(s[:len(s)-1]), nil
}

func (s *getterSetterD) SetBSON(raw bson.Raw) error {
	var doc bson.D
	err := raw.Unmarshal(&doc)
	doc = append(doc, bson.DocElem{Name: "suffix", Value: true})
	*s = getterSetterD(doc)
	return err
}

type getterSetterInt int

func (i getterSetterInt) GetBSON() (interface{}, error) {
	return bson.D{{Name: "a", Value: int(i)}}, nil
}

func (i *getterSetterInt) SetBSON(raw bson.Raw) error {
	var doc struct{ A int }
	err := raw.Unmarshal(&doc)
	*i = getterSetterInt(doc.A)
	return err
}

type ifaceType interface {
	Hello()
}

type ifaceSlice []ifaceType

func (s *ifaceSlice) SetBSON(raw bson.Raw) error {
	var ns []int
	if err := raw.Unmarshal(&ns); err != nil {
		return err
	}
	*s = make(ifaceSlice, ns[0])
	return nil
}

func (s ifaceSlice) GetBSON() (interface{}, error) {
	return []int{len(s)}, nil
}

type (
	MyString string
	MyBytes  []byte
	MyBool   bool
	MyD      []bson.DocElem
	MyRawD   []bson.RawDocElem
	MyM      map[string]interface{}
)

var (
	truevar  = true
	falsevar = false

	int64var = int64(42)
	int64ptr = &int64var
	intvar   = int(42)
	intptr   = &intvar

	gsintvar = getterSetterInt(42)
)

func parseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

// That's a pretty fun test.  It will dump the first item, generate a zero
// value equivalent to the second one, load the dumped data onto it, and then
// verify that the resulting value is deep-equal to the untouched second value.
// Then, it will do the same in the *opposite* direction!
var twoWayCrossItems = []crossTypeItem{
	// int<=>int
	{&struct{ I int }{42}, &struct{ I int8 }{42}},
	{&struct{ I int }{42}, &struct{ I int32 }{42}},
	{&struct{ I int }{42}, &struct{ I int64 }{42}},
	{&struct{ I int8 }{42}, &struct{ I int32 }{42}},
	{&struct{ I int8 }{42}, &struct{ I int64 }{42}},
	{&struct{ I int32 }{42}, &struct{ I int64 }{42}},

	// uint<=>uint
	{&struct{ I uint }{42}, &struct{ I uint8 }{42}},
	{&struct{ I uint }{42}, &struct{ I uint32 }{42}},
	{&struct{ I uint }{42}, &struct{ I uint64 }{42}},
	{&struct{ I uint8 }{42}, &struct{ I uint32 }{42}},
	{&struct{ I uint8 }{42}, &struct{ I uint64 }{42}},
	{&struct{ I uint32 }{42}, &struct{ I uint64 }{42}},

	// float32<=>float64
	{&struct{ I float32 }{42}, &struct{ I float64 }{42}},

	// int<=>uint
	{&struct{ I uint }{42}, &struct{ I int }{42}},
	{&struct{ I uint }{42}, &struct{ I int8 }{42}},
	{&struct{ I uint }{42}, &struct{ I int32 }{42}},
	{&struct{ I uint }{42}, &struct{ I int64 }{42}},
	{&struct{ I uint8 }{42}, &struct{ I int }{42}},
	{&struct{ I uint8 }{42}, &struct{ I int8 }{42}},
	{&struct{ I uint8 }{42}, &struct{ I int32 }{42}},
	{&struct{ I uint8 }{42}, &struct{ I int64 }{42}},
	{&struct{ I uint32 }{42}, &struct{ I int }{42}},
	{&struct{ I uint32 }{42}, &struct{ I int8 }{42}},
	{&struct{ I uint32 }{42}, &struct{ I int32 }{42}},
	{&struct{ I uint32 }{42}, &struct{ I int64 }{42}},
	{&struct{ I uint64 }{42}, &struct{ I int }{42}},
	{&struct{ I uint64 }{42}, &struct{ I int8 }{42}},
	{&struct{ I uint64 }{42}, &struct{ I int32 }{42}},
	{&struct{ I uint64 }{42}, &struct{ I int64 }{42}},

	// int <=> float
	{&struct{ I int }{42}, &struct{ I float64 }{42}},

	// int <=> bool
	{&struct{ I int }{1}, &struct{ I bool }{true}},
	{&struct{ I int }{0}, &struct{ I bool }{false}},

	// uint <=> float64
	{&struct{ I uint }{42}, &struct{ I float64 }{42}},

	// uint <=> bool
	{&struct{ I uint }{1}, &struct{ I bool }{true}},
	{&struct{ I uint }{0}, &struct{ I bool }{false}},

	// float64 <=> bool
	{&struct{ I float64 }{1}, &struct{ I bool }{true}},
	{&struct{ I float64 }{0}, &struct{ I bool }{false}},

	// string <=> string and string <=> []byte
	{&struct{ S []byte }{[]byte("abc")}, &struct{ S string }{"abc"}},
	{&struct{ S []byte }{[]byte("def")}, &struct{ S bson.Symbol }{"def"}},
	{&struct{ S string }{"ghi"}, &struct{ S bson.Symbol }{"ghi"}},

	// map <=> struct
	{&struct {
		A struct {
			B, C int
		}
	}{struct{ B, C int }{1, 2}},
		map[string]map[string]int{"a": {"b": 1, "c": 2}}},

	{&struct{ A bson.Symbol }{"abc"}, map[string]string{"a": "abc"}},
	{&struct{ A bson.Symbol }{"abc"}, map[string][]byte{"a": []byte("abc")}},
	{&struct{ A []byte }{[]byte("abc")}, map[string]string{"a": "abc"}},
	{&struct{ A uint }{42}, map[string]int{"a": 42}},
	{&struct{ A uint }{42}, map[string]float64{"a": 42}},
	{&struct{ A uint }{1}, map[string]bool{"a": true}},
	{&struct{ A int }{42}, map[string]uint{"a": 42}},
	{&struct{ A int }{42}, map[string]float64{"a": 42}},
	{&struct{ A int }{1}, map[string]bool{"a": true}},
	{&struct{ A float64 }{42}, map[string]float32{"a": 42}},
	{&struct{ A float64 }{42}, map[string]int{"a": 42}},
	{&struct{ A float64 }{42}, map[string]uint{"a": 42}},
	{&struct{ A float64 }{1}, map[string]bool{"a": true}},
	{&struct{ A bool }{true}, map[string]int{"a": 1}},
	{&struct{ A bool }{true}, map[string]uint{"a": 1}},
	{&struct{ A bool }{true}, map[string]float64{"a": 1}},
	{&struct{ A **byte }{&byteptr}, map[string]byte{"a": 8}},

	// url.URL <=> string
	{&struct{ URL *url.URL }{parseURL("h://e.c/p")}, map[string]string{"url": "h://e.c/p"}},
	{&struct{ URL url.URL }{*parseURL("h://e.c/p")}, map[string]string{"url": "h://e.c/p"}},

	// Slices
	{&struct{ S []int }{[]int{1, 2, 3}}, map[string][]int{"s": {1, 2, 3}}},
	{&struct{ S *[]int }{&[]int{1, 2, 3}}, map[string][]int{"s": {1, 2, 3}}},

	// Conditionals
	{&condBool{true}, map[string]bool{"v": true}},
	{&condBool{}, map[string]bool{}},
	{&condInt{1}, map[string]int{"v": 1}},
	{&condInt{}, map[string]int{}},
	{&condUInt{1}, map[string]uint{"v": 1}},
	{&condUInt{}, map[string]uint{}},
	{&condFloat{}, map[string]int{}},
	{&condStr{"yo"}, map[string]string{"v": "yo"}},
	{&condStr{}, map[string]string{}},
	{&condStrNS{"yo"}, map[string]string{"v": "yo"}},
	{&condStrNS{}, map[string]string{}},
	{&condSlice{[]string{"yo"}}, map[string][]string{"v": {"yo"}}},
	{&condSlice{}, map[string][]string{}},
	{&condMap{map[string]int{"k": 1}}, bson.M{"v": bson.M{"k": 1}}},
	{&condMap{}, map[string][]string{}},
	{&condIface{"yo"}, map[string]string{"v": "yo"}},
	{&condIface{""}, map[string]string{"v": ""}},
	{&condIface{}, map[string]string{}},
	{&condPtr{&truevar}, map[string]bool{"v": true}},
	{&condPtr{&falsevar}, map[string]bool{"v": false}},
	{&condPtr{}, map[string]string{}},

	{&condTime{time.Unix(123456789, 123e6).UTC()}, map[string]time.Time{"v": time.Unix(123456789, 123e6).UTC()}},
	{&condTime{}, map[string]string{}},

	{&condStruct{struct{ A []int }{[]int{1}}}, bson.M{"v": bson.M{"a": []interface{}{1}}}},
	{&condStruct{struct{ A []int }{}}, bson.M{}},

	{&condRaw{bson.Raw{Kind: 0x0A, Data: []byte{}}}, bson.M{"v": nil}},
	{&condRaw{bson.Raw{Kind: 0x00}}, bson.M{}},

	{&namedCondStr{"yo"}, map[string]string{"myv": "yo"}},
	{&namedCondStr{}, map[string]string{}},

	{&shortInt{1}, map[string]interface{}{"v": 1}},
	{&shortInt{1 << 30}, map[string]interface{}{"v": 1 << 30}},
	{&shortInt{1 << 31}, map[string]interface{}{"v": int64(1 << 31)}},
	{&shortUint{1 << 30}, map[string]interface{}{"v": 1 << 30}},
	{&shortUint{1 << 31}, map[string]interface{}{"v": int64(1 << 31)}},
	{&shortIface{int64(1) << 31}, map[string]interface{}{"v": int64(1 << 31)}},
	{&shortPtr{int64ptr}, map[string]interface{}{"v": intvar}},

	{&shortNonEmptyInt{1}, map[string]interface{}{"v": 1}},
	{&shortNonEmptyInt{1 << 31}, map[string]interface{}{"v": int64(1 << 31)}},
	{&shortNonEmptyInt{}, map[string]interface{}{}},

	{&inlineInt{struct{ A, B int }{1, 2}}, map[string]interface{}{"a": 1, "b": 2}},
	{&inlineMap{A: 1, M: map[string]interface{}{"b": 2}}, map[string]interface{}{"a": 1, "b": 2}},
	{&inlineMap{A: 1, M: nil}, map[string]interface{}{"a": 1}},
	{&inlineMapInt{A: 1, M: map[string]int{"b": 2}}, map[string]int{"a": 1, "b": 2}},
	{&inlineMapInt{A: 1, M: nil}, map[string]int{"a": 1}},
	{&inlineMapMyM{A: 1, M: MyM{"b": MyM{"c": 3}}}, map[string]interface{}{"a": 1, "b": map[string]interface{}{"c": 3}}},
	{&inlineUnexported{M: map[string]interface{}{"b": 1}, unexported: unexported{A: 2}}, map[string]interface{}{"b": 1, "a": 2}},

	// []byte <=> Binary
	{&struct{ B []byte }{[]byte("abc")}, map[string]bson.Binary{"b": {Data: []byte("abc")}}},

	// []byte <=> MyBytes
	{&struct{ B MyBytes }{[]byte("abc")}, map[string]string{"b": "abc"}},
	{&struct{ B MyBytes }{[]byte{}}, map[string]string{"b": ""}},
	{&struct{ B MyBytes }{}, map[string]bool{}},
	{&struct{ B []byte }{[]byte("abc")}, map[string]MyBytes{"b": []byte("abc")}},

	// bool <=> MyBool
	{&struct{ B MyBool }{true}, map[string]bool{"b": true}},
	{&struct{ B MyBool }{}, map[string]bool{"b": false}},
	{&struct{ B MyBool }{}, map[string]string{}},
	{&struct{ B bool }{}, map[string]MyBool{"b": false}},

	// arrays
	{&struct{ V [2]int }{[...]int{1, 2}}, map[string][2]int{"v": {1, 2}}},
	{&struct{ V [2]byte }{[...]byte{1, 2}}, map[string][2]byte{"v": {1, 2}}},

	// zero time
	{&struct{ V time.Time }{}, map[string]interface{}{"v": time.Time{}}},

	// zero time + 1 second + 1 millisecond; overflows int64 as nanoseconds
	{&struct{ V time.Time }{time.Unix(-62135596799, 1e6).UTC()},
		map[string]interface{}{"v": time.Unix(-62135596799, 1e6).UTC()}},

	// bson.D <=> []DocElem
	{&bson.D{{Name: "a", Value: bson.D{{Name: "b", Value: 1}, {Name: "c", Value: 2}}}}, &bson.D{{Name: "a", Value: bson.D{{Name: "b", Value: 1}, {Name: "c", Value: 2}}}}},
	{&bson.D{{Name: "a", Value: bson.D{{Name: "b", Value: 1}, {Name: "c", Value: 2}}}}, &MyD{{Name: "a", Value: MyD{{Name: "b", Value: 1}, {Name: "c", Value: 2}}}}},
	{&struct{ V MyD }{MyD{{Name: "a", Value: 1}}}, &bson.D{{Name: "v", Value: bson.D{{Name: "a", Value: 1}}}}},

	// bson.RawD <=> []RawDocElem
	{&bson.RawD{{Name: "a", Value: bson.Raw{Kind: 0x08, Data: []byte{0x01}}}}, &bson.RawD{{Name: "a", Value: bson.Raw{Kind: 0x08, Data: []byte{0x01}}}}},
	{&bson.RawD{{Name: "a", Value: bson.Raw{Kind: 0x08, Data: []byte{0x01}}}}, &MyRawD{{Name: "a", Value: bson.Raw{Kind: 0x08, Data: []byte{0x01}}}}},

	// bson.M <=> map
	{bson.M{"a": bson.M{"b": 1, "c": 2}}, MyM{"a": MyM{"b": 1, "c": 2}}},
	{bson.M{"a": bson.M{"b": 1, "c": 2}}, map[string]interface{}{"a": map[string]interface{}{"b": 1, "c": 2}}},

	// bson.M <=> map[MyString]
	{bson.M{"a": bson.M{"b": 1, "c": 2}}, map[MyString]interface{}{"a": map[MyString]interface{}{"b": 1, "c": 2}}},

	// json.Number <=> int64, float64
	{&struct{ N json.Number }{"5"}, map[string]interface{}{"n": int64(5)}},
	{&struct{ N json.Number }{"5.05"}, map[string]interface{}{"n": 5.05}},
	{&struct{ N json.Number }{"9223372036854776000"}, map[string]interface{}{"n": float64(1 << 63)}},

	// bson.D <=> non-struct getter/setter
	{&bson.D{{Name: "a", Value: 1}}, &getterSetterD{{Name: "a", Value: 1}, {Name: "suffix", Value: true}}},
	{&bson.D{{Name: "a", Value: 42}}, &gsintvar},

	// Interface slice setter.
	{&struct{ V ifaceSlice }{ifaceSlice{nil, nil, nil}}, bson.M{"v": []interface{}{3}}},
}

// Same thing, but only one way (obj1 => obj2).
var oneWayCrossItems = []crossTypeItem{
	// map <=> struct
	{map[string]interface{}{"a": 1, "b": "2", "c": 3}, map[string]int{"a": 1, "c": 3}},

	// inline map elides badly typed values
	{map[string]interface{}{"a": 1, "b": "2", "c": 3}, &inlineMapInt{A: 1, M: map[string]int{"c": 3}}},

	// Can't decode int into struct.
	{bson.M{"a": bson.M{"b": 2}}, &struct{ A bool }{}},

	// Would get decoded into a int32 too in the opposite direction.
	{&shortIface{int64(1) << 30}, map[string]interface{}{"v": 1 << 30}},

	// Ensure omitempty on struct with private fields works properly.
	{&struct {
		V struct{ v time.Time } `bson:",omitempty"`
	}{}, map[string]interface{}{}},

	// Attempt to marshal slice into RawD (issue #120).
	{bson.M{"x": []int{1, 2, 3}}, &struct{ X bson.RawD }{}},
}

func testCrossPair(c *C, dump interface{}, load interface{}) {
	c.Logf("Dump: %#v", dump)
	c.Logf("Load: %#v", load)
	zero := makeZeroDoc(load)
	data, err := bson.Marshal(dump)
	c.Assert(err, IsNil)
	c.Logf("Dumped: %#v", string(data))
	err = bson.Unmarshal(data, zero)
	c.Assert(err, IsNil)
	c.Logf("Loaded: %#v", zero)
	c.Assert(zero, DeepEquals, load)
}

func (s *S) TestTwoWayCrossPairs(c *C) {
	for _, item := range twoWayCrossItems {
		testCrossPair(c, item.obj1, item.obj2)
		testCrossPair(c, item.obj2, item.obj1)
	}
}

func (s *S) TestOneWayCrossPairs(c *C) {
	for _, item := range oneWayCrossItems {
		testCrossPair(c, item.obj1, item.obj2)
	}
}

// --------------------------------------------------------------------------
// ObjectId hex representation test.

func (s *S) TestObjectIdHex(c *C) {
	id := bson.ObjectIdHex("4d88e15b60f486e428412dc9")
	c.Assert(id.String(), Equals, `ObjectIdHex("4d88e15b60f486e428412dc9")`)
	c.Assert(id.Hex(), Equals, "4d88e15b60f486e428412dc9")
}

func (s *S) TestIsObjectIdHex(c *C) {
	test := []struct {
		id    string
		valid bool
	}{
		{"4d88e15b60f486e428412dc9", true},
		{"4d88e15b60f486e428412dc", false},
		{"4d88e15b60f486e428412dc9e", false},
		{"4d88e15b60f486e428412dcx", false},
	}
	for _, t := range test {
		c.Assert(bson.IsObjectIdHex(t.id), Equals, t.valid)
	}
}

// --------------------------------------------------------------------------
// ObjectId parts extraction tests.

type objectIdParts struct {
	id        bson.ObjectId
	timestamp int64
	machine   []byte
	pid       uint16
	counter   int32
}

var objectIds = []objectIdParts{
	{
		bson.ObjectIdHex("4d88e15b60f486e428412dc9"),
		1300816219,
		[]byte{0x60, 0xf4, 0x86},
		0xe428,
		4271561,
	},
	{
		bson.ObjectIdHex("000000000000000000000000"),
		0,
		[]byte{0x00, 0x00, 0x00},
		0x0000,
		0,
	},
	{
		bson.ObjectIdHex("00000000aabbccddee000001"),
		0,
		[]byte{0xaa, 0xbb, 0xcc},
		0xddee,
		1,
	},
}

func (s *S) TestObjectIdPartsExtraction(c *C) {
	for i, v := range objectIds {
		t := time.Unix(v.timestamp, 0)
		c.Assert(v.id.Time(), Equals, t, Commentf("#%d Wrong timestamp value", i))
		c.Assert(v.id.Machine(), DeepEquals, v.machine, Commentf("#%d Wrong machine id value", i))
		c.Assert(v.id.Pid(), Equals, v.pid, Commentf("#%d Wrong pid value", i))
		c.Assert(v.id.Counter(), Equals, v.counter, Commentf("#%d Wrong counter value", i))
	}
}

func (s *S) TestNow(c *C) {
	before := time.Now()
	time.Sleep(1e6)
	now := bson.Now()
	time.Sleep(1e6)
	after := time.Now()
	c.Assert(now.After(before) && now.Before(after), Equals, true, Commentf("now=%s, before=%s, after=%s", now, before, after))
}

// --------------------------------------------------------------------------
// ObjectId generation tests.

func (s *S) TestNewObjectId(c *C) {
	// Generate 10 ids
	ids := make([]bson.ObjectId, 10)
	for i := 0; i < 10; i++ {
		ids[i] = bson.NewObjectId()
	}
	for i := 1; i < 10; i++ {
		prevId := ids[i-1]
		id := ids[i]
		// Test for uniqueness among all other 9 generated ids
		for j, tid := range ids {
			if j != i {
				c.Assert(id, Not(Equals), tid, Commentf("Generated ObjectId is not unique"))
			}
		}
		// Check that timestamp was incremented and is within 30 seconds of the previous one
		secs := id.Time().Sub(prevId.Time()).Seconds()
		c.Assert((secs >= 0 && secs <= 30), Equals, true, Commentf("Wrong timestamp in generated ObjectId"))
		// Check that machine ids are the same
		c.Assert(id.Machine(), DeepEquals, prevId.Machine())
		// Check that pids are the same
		c.Assert(id.Pid(), Equals, prevId.Pid())
		// Test for proper increment
		delta := int(id.Counter() - prevId.Counter())
		c.Assert(delta, Equals, 1, Commentf("Wrong increment in generated ObjectId"))
	}
}

func (s *S) TestNewObjectIdWithTime(c *C) {
	t := time.Unix(12345678, 0)
	id := bson.NewObjectIdWithTime(t)
	c.Assert(id.Time(), Equals, t)
	c.Assert(id.Machine(), DeepEquals, []byte{0x00, 0x00, 0x00})
	c.Assert(int(id.Pid()), Equals, 0)
	c.Assert(int(id.Counter()), Equals, 0)
}

// --------------------------------------------------------------------------
// ObjectId JSON marshalling.

type jsonType struct {
	Id bson.ObjectId
}

var jsonIdTests = []struct {
	value     jsonType
	json      string
	marshal   bool
	unmarshal bool
	error     string
}{{
	value:     jsonType{Id: bson.ObjectIdHex("4d88e15b60f486e428412dc9")},
	json:      `{"Id":"4d88e15b60f486e428412dc9"}`,
	marshal:   true,
	unmarshal: true,
}, {
	value:     jsonType{},
	json:      `{"Id":""}`,
	marshal:   true,
	unmarshal: true,
}, {
	value:     jsonType{},
	json:      `{"Id":null}`,
	marshal:   false,
	unmarshal: true,
}, {
	json:      `{"Id":"4d88e15b60f486e428412dc9A"}`,
	error:     `invalid ObjectId in JSON: "4d88e15b60f486e428412dc9A"`,
	marshal:   false,
	unmarshal: true,
}, {
	json:      `{"Id":"4d88e15b60f486e428412dcZ"}`,
	error:     `invalid ObjectId in JSON: "4d88e15b60f486e428412dcZ" .*`,
	marshal:   false,
	unmarshal: true,
}}

func (s *S) TestObjectIdJSONMarshaling(c *C) {
	for _, test := range jsonIdTests {
		if test.marshal {
			data, err := json.Marshal(&test.value)
			if test.error == "" {
				c.Assert(err, IsNil)
				c.Assert(string(data), Equals, test.json)
			} else {
				c.Assert(err, ErrorMatches, test.error)
			}
		}

		if test.unmarshal {
			var value jsonType
			err := json.Unmarshal([]byte(test.json), &value)
			if test.error == "" {
				c.Assert(err, IsNil)
				c.Assert(value, DeepEquals, test.value)
			} else {
				c.Assert(err, ErrorMatches, test.error)
			}
		}
	}
}

// --------------------------------------------------------------------------
// ObjectId Text encoding.TextUnmarshaler.

var textIdTests = []struct {
	value     bson.ObjectId
	text      string
	marshal   bool
	unmarshal bool
	error     string
}{{
	value:     bson.ObjectIdHex("4d88e15b60f486e428412dc9"),
	text:      "4d88e15b60f486e428412dc9",
	marshal:   true,
	unmarshal: true,
}, {
	text:      "",
	marshal:   true,
	unmarshal: true,
}, {
	text:      "4d88e15b60f486e428412dc9A",
	marshal:   false,
	unmarshal: true,
	error:     `invalid ObjectId: 4d88e15b60f486e428412dc9A`,
}, {
	text:      "4d88e15b60f486e428412dcZ",
	marshal:   false,
	unmarshal: true,
	error:     `invalid ObjectId: 4d88e15b60f486e428412dcZ .*`,
}}

func (s *S) TestObjectIdTextMarshaling(c *C) {
	for _, test := range textIdTests {
		if test.marshal {
			data, err := test.value.MarshalText()
			if test.error == "" {
				c.Assert(err, IsNil)
				c.Assert(string(data), Equals, test.text)
			} else {
				c.Assert(err, ErrorMatches, test.error)
			}
		}

		if test.unmarshal {
			err := test.value.UnmarshalText([]byte(test.text))
			if test.error == "" {
				c.Assert(err, IsNil)
				if test.value != "" {
					value := bson.ObjectIdHex(test.text)
					c.Assert(value, DeepEquals, test.value)
				}
			} else {
				c.Assert(err, ErrorMatches, test.error)
			}
		}
	}
}

// --------------------------------------------------------------------------
// ObjectId XML marshalling.

type xmlType struct {
	Id bson.ObjectId
}

var xmlIdTests = []struct {
	value     xmlType
	xml       string
	marshal   bool
	unmarshal bool
	error     string
}{{
	value:     xmlType{Id: bson.ObjectIdHex("4d88e15b60f486e428412dc9")},
	xml:       "<xmlType><Id>4d88e15b60f486e428412dc9</Id></xmlType>",
	marshal:   true,
	unmarshal: true,
}, {
	value:     xmlType{},
	xml:       "<xmlType><Id></Id></xmlType>",
	marshal:   true,
	unmarshal: true,
}, {
	xml:       "<xmlType><Id>4d88e15b60f486e428412dc9A</Id></xmlType>",
	marshal:   false,
	unmarshal: true,
	error:     `invalid ObjectId: 4d88e15b60f486e428412dc9A`,
}, {
	xml:       "<xmlType><Id>4d88e15b60f486e428412dcZ</Id></xmlType>",
	marshal:   false,
	unmarshal: true,
	error:     `invalid ObjectId: 4d88e15b60f486e428412dcZ .*`,
}}

func (s *S) TestObjectIdXMLMarshaling(c *C) {
	for _, test := range xmlIdTests {
		if test.marshal {
			data, err := xml.Marshal(&test.value)
			if test.error == "" {
				c.Assert(err, IsNil)
				c.Assert(string(data), Equals, test.xml)
			} else {
				c.Assert(err, ErrorMatches, test.error)
			}
		}

		if test.unmarshal {
			var value xmlType
			err := xml.Unmarshal([]byte(test.xml), &value)
			if test.error == "" {
				c.Assert(err, IsNil)
				c.Assert(value, DeepEquals, test.value)
			} else {
				c.Assert(err, ErrorMatches, test.error)
			}
		}
	}
}

// --------------------------------------------------------------------------
// Some simple benchmarks.

type BenchT struct {
	A, B, C, D, E, F string
}

type BenchRawT struct {
	A string
	B int
	C bson.M
	D []float64
}

func (s *S) BenchmarkUnmarhsalStruct(c *C) {
	v := BenchT{A: "A", D: "D", E: "E"}
	data, err := bson.Marshal(&v)
	if err != nil {
		panic(err)
	}
	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		err = bson.Unmarshal(data, &v)
	}
	if err != nil {
		panic(err)
	}
}

func (s *S) BenchmarkUnmarhsalMap(c *C) {
	m := bson.M{"a": "a", "d": "d", "e": "e"}
	data, err := bson.Marshal(&m)
	if err != nil {
		panic(err)
	}
	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		err = bson.Unmarshal(data, &m)
	}
	if err != nil {
		panic(err)
	}
}

func (s *S) BenchmarkUnmarshalRaw(c *C) {
	var err error
	m := BenchRawT{
		A: "test_string",
		B: 123,
		C: bson.M{
			"subdoc_int": 12312,
			"subdoc_doc": bson.M{"1": 1},
		},
		D: []float64{0.0, 1.3333, -99.9997, 3.1415},
	}
	data, err := bson.Marshal(&m)
	if err != nil {
		panic(err)
	}
	raw := bson.Raw{}
	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		err = bson.Unmarshal(data, &raw)
	}
	if err != nil {
		panic(err)
	}
}

func (s *S) BenchmarkNewObjectId(c *C) {
	for i := 0; i < c.N; i++ {
		bson.NewObjectId()
	}
}

func (s *S) TestMarshalRespectNil(c *C) {
	type T struct {
		Slice    []int
		SlicePtr *[]int
		Ptr      *int
		Map      map[string]interface{}
		MapPtr   *map[string]interface{}
	}

	bson.SetRespectNilValues(true)
	defer bson.SetRespectNilValues(false)

	testStruct1 := T{}

	c.Assert(testStruct1.Slice, IsNil)
	c.Assert(testStruct1.SlicePtr, IsNil)
	c.Assert(testStruct1.Map, IsNil)
	c.Assert(testStruct1.MapPtr, IsNil)
	c.Assert(testStruct1.Ptr, IsNil)

	b, _ := bson.Marshal(testStruct1)

	testStruct2 := T{}

	bson.Unmarshal(b, &testStruct2)

	c.Assert(testStruct2.Slice, IsNil)
	c.Assert(testStruct2.SlicePtr, IsNil)
	c.Assert(testStruct2.Map, IsNil)
	c.Assert(testStruct2.MapPtr, IsNil)
	c.Assert(testStruct2.Ptr, IsNil)

	testStruct1 = T{
		Slice:    []int{},
		SlicePtr: &[]int{},
		Map:      map[string]interface{}{},
		MapPtr:   &map[string]interface{}{},
	}

	c.Assert(testStruct1.Slice, NotNil)
	c.Assert(testStruct1.SlicePtr, NotNil)
	c.Assert(testStruct1.Map, NotNil)
	c.Assert(testStruct1.MapPtr, NotNil)

	b, _ = bson.Marshal(testStruct1)

	testStruct2 = T{}

	bson.Unmarshal(b, &testStruct2)

	c.Assert(testStruct2.Slice, NotNil)
	c.Assert(testStruct2.SlicePtr, NotNil)
	c.Assert(testStruct2.Map, NotNil)
	c.Assert(testStruct2.MapPtr, NotNil)
}

func (s *S) TestMongoTimestampTime(c *C) {
	t := time.Now()
	ts, err := bson.NewMongoTimestamp(t, 123)
	c.Assert(err, IsNil)
	c.Assert(ts.Time().Unix(), Equals, t.Unix())
}

func (s *S) TestMongoTimestampCounter(c *C) {
	rnd := rand.Uint32()
	ts, err := bson.NewMongoTimestamp(time.Now(), rnd)
	c.Assert(err, IsNil)
	c.Assert(ts.Counter(), Equals, rnd)
}

func (s *S) TestMongoTimestampError(c *C) {
	t := time.Date(1969, time.December, 31, 23, 59, 59, 999, time.UTC)
	ts, err := bson.NewMongoTimestamp(t, 321)
	c.Assert(int64(ts), Equals, int64(-1))
	c.Assert(err, ErrorMatches, "invalid value for time")
}

func ExampleNewMongoTimestamp() {

	var counter uint32 = 1
	var t time.Time

	for i := 1; i <= 3; i++ {

		if c := time.Now(); t.Unix() == c.Unix() {
			counter++
		} else {
			t = c
			counter = 1
		}

		ts, err := bson.NewMongoTimestamp(t, counter)
		if err != nil {
			fmt.Printf("NewMongoTimestamp error: %v", err)
		} else {
			fmt.Printf("NewMongoTimestamp encoded timestamp: %d\n", ts)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
