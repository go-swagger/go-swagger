package mgo

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	. "gopkg.in/check.v1"
)

type S struct{}

var _ = Suite(&S{})

// This file is for testing functions that are not exported outside the mgo
// package - avoid doing so if at all possible.

// Ensures indexed int64 fields do not cause mgo to panic.
//
// See https://github.com/globalsign/mgo/pull/23
func TestIndexedInt64FieldsBug(t *testing.T) {
	input := bson.D{
		{Name: "testkey", Value: int(1)},
		{Name: "testkey", Value: int64(1)},
		{Name: "testkey", Value: "test"},
		{Name: "testkey", Value: float64(1)},
	}

	_ = simpleIndexKey(input)
}

func (s *S) TestGetRFC2253NameStringSingleValued(c *C) {
	var RDNElements = pkix.RDNSequence{
		{{Type: asn1.ObjectIdentifier{2, 5, 4, 6}, Value: "GO"}},
		{{Type: asn1.ObjectIdentifier{2, 5, 4, 8}, Value: "MGO"}},
		{{Type: asn1.ObjectIdentifier{2, 5, 4, 7}, Value: "MGO"}},
		{{Type: asn1.ObjectIdentifier{2, 5, 4, 10}, Value: "MGO"}},
		{{Type: asn1.ObjectIdentifier{2, 5, 4, 11}, Value: "Client"}},
		{{Type: asn1.ObjectIdentifier{2, 5, 4, 3}, Value: "localhost"}},
	}

	c.Assert(getRFC2253NameString(&RDNElements), Equals, "CN=localhost,OU=Client,O=MGO,L=MGO,ST=MGO,C=GO")
}

func (s *S) TestGetRFC2253NameStringEscapeChars(c *C) {
	var RDNElements = pkix.RDNSequence{
		{{Type: asn1.ObjectIdentifier{2, 5, 4, 6}, Value: "GB"}},
		{{Type: asn1.ObjectIdentifier{2, 5, 4, 8}, Value: "MGO "}},
		{{Type: asn1.ObjectIdentifier{2, 5, 4, 10}, Value: "Sue, Grabbit and Runn < > ;"}},
		{{Type: asn1.ObjectIdentifier{2, 5, 4, 3}, Value: "L. Eagle"}},
	}

	c.Assert(getRFC2253NameString(&RDNElements), Equals, "CN=L. Eagle,O=Sue\\, Grabbit and Runn \\< \\> \\;,ST=MGO\\ ,C=GB")
}

func (s *S) TestGetRFC2253NameStringMultiValued(c *C) {
	var RDNElements = pkix.RDNSequence{
		{{Type: asn1.ObjectIdentifier{2, 5, 4, 6}, Value: "US"}},
		{{Type: asn1.ObjectIdentifier{2, 5, 4, 10}, Value: "Widget Inc."}},
		{{Type: asn1.ObjectIdentifier{2, 5, 4, 11}, Value: "Sales"}, {Type: asn1.ObjectIdentifier{2, 5, 4, 3}, Value: "J. Smith"}},
	}

	c.Assert(getRFC2253NameString(&RDNElements), Equals, "OU=Sales+CN=J. Smith,O=Widget Inc.,C=US")
}

func (s *S) TestDialTimeouts(c *C) {
	info := &DialInfo{}

	c.Assert(info.readTimeout(), Equals, time.Duration(0))
	c.Assert(info.writeTimeout(), Equals, time.Duration(0))
	c.Assert(info.roundTripTimeout(), Equals, time.Duration(0))

	info.Timeout = 60 * time.Second
	c.Assert(info.readTimeout(), Equals, 60*time.Second)
	c.Assert(info.writeTimeout(), Equals, 60*time.Second)
	c.Assert(info.roundTripTimeout(), Equals, 120*time.Second)

	info.ReadTimeout = time.Second
	c.Assert(info.readTimeout(), Equals, time.Second)

	info.WriteTimeout = time.Second
	c.Assert(info.writeTimeout(), Equals, time.Second)
}
