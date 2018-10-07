package bson_test

import (
	"github.com/globalsign/mgo/bson"
	. "gopkg.in/check.v1"
)

type mixedTagging struct {
	First  string
	Second string `bson:"second_field"`
	Third  string `json:"third_field"`
	Fourth string `bson:"fourth_field" json:"alternate"`
}

// TestTaggingFallback checks that tagging fallback can be used/works as expected.
func (s *S) TestTaggingFallback(c *C) {
	initial := &mixedTagging{
		First:  "One",
		Second: "Two",
		Third:  "Three",
		Fourth: "Four",
	}

	// Take only testing.T, leave only footprints.
	initialState := bson.JSONTagFallbackState()
	defer bson.SetJSONTagFallback(initialState)

	// Marshal with the new mode applied.
	bson.SetJSONTagFallback(true)
	bsonState, errBSON := bson.Marshal(initial)
	c.Assert(errBSON, IsNil)

	// Unmarshal into a generic map so that we can pick up the actual field names
	// selected.
	target := make(map[string]string)
	errUnmarshal := bson.Unmarshal(bsonState, target)
	c.Assert(errUnmarshal, IsNil)

	// No tag, so standard naming
	_, firstExists := target["first"]
	c.Assert(firstExists, Equals, true)

	// Just a BSON tag
	_, secondExists := target["second_field"]
	c.Assert(secondExists, Equals, true)

	// Just a JSON tag
	_, thirdExists := target["third_field"]
	c.Assert(thirdExists, Equals, true)

	// Should marshal 4th as fourth_field (since we have both tags)
	_, fourthExists := target["fourth_field"]
	c.Assert(fourthExists, Equals, true)
}
