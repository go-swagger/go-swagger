package bson_test

import (
	"bytes"

	"github.com/globalsign/mgo/bson"
	. "gopkg.in/check.v1"
)

var invalidSizeDocuments = [][]byte{
	// Empty document
	[]byte{},
	// Incomplete header
	[]byte{0x04},
	// Negative size
	[]byte{0xff, 0xff, 0xff, 0xff},
	// Full, valid size header but too small (less than 5 bytes)
	[]byte{0x04, 0x00, 0x00, 0x00},
	// Valid header, valid size but incomplete document
	[]byte{0xff, 0x00, 0x00, 0x00, 0x00},
	// Too big
	[]byte{0xff, 0xff, 0xff, 0x7f},
}

// Reusing sampleItems from bson_test

func (s *S) TestEncodeSampleItems(c *C) {
	for i, item := range sampleItems {
		buf := bytes.NewBuffer(nil)
		enc := bson.NewEncoder(buf)

		err := enc.Encode(item.obj)
		c.Assert(err, IsNil)
		c.Assert(string(buf.Bytes()), Equals, item.data, Commentf("Failed on item %d", i))
	}
}

func (s *S) TestDecodeSampleItems(c *C) {
	for i, item := range sampleItems {
		buf := bytes.NewBuffer([]byte(item.data))
		dec := bson.NewDecoder(buf)

		value := bson.M{}
		err := dec.Decode(&value)
		c.Assert(err, IsNil)
		c.Assert(value, DeepEquals, item.obj, Commentf("Failed on item %d", i))
	}
}

func (s *S) TestStreamRoundTrip(c *C) {
	buf := bytes.NewBuffer(nil)
	enc := bson.NewEncoder(buf)

	for _, item := range sampleItems {
		err := enc.Encode(item.obj)
		c.Assert(err, IsNil)
	}

	// Ensure that everything that was encoded is decodable in the same order.
	dec := bson.NewDecoder(buf)
	for i, item := range sampleItems {
		value := bson.M{}
		err := dec.Decode(&value)
		c.Assert(err, IsNil)
		c.Assert(value, DeepEquals, item.obj, Commentf("Failed on item %d", i))
	}
}

func (s *S) TestDecodeDocumentTooSmall(c *C) {
	for i, item := range invalidSizeDocuments {
		buf := bytes.NewBuffer(item)
		dec := bson.NewDecoder(buf)
		value := bson.M{}
		err := dec.Decode(&value)
		c.Assert(err, NotNil, Commentf("Failed on invalid size item %d", i))
	}
}
