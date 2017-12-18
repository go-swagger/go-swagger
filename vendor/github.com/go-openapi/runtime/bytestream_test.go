package runtime

import (
	"bytes"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteStreamConsumer(t *testing.T) {
	cons := ByteStreamConsumer()

	expected := "the data for the stream to be sent over the wire"

	// can consume as a Writer
	var b bytes.Buffer
	if assert.NoError(t, cons.Consume(bytes.NewBufferString(expected), &b)) {
		assert.Equal(t, expected, b.String())
	}

	// can consume as an UnmarshalBinary
	var bu binaryUnmarshalDummy
	if assert.NoError(t, cons.Consume(bytes.NewBufferString(expected), &bu)) {
		assert.Equal(t, expected, bu.str)
	}

	// can consume as a binary slice
	var bs []byte
	if assert.NoError(t, cons.Consume(bytes.NewBufferString(expected), &bs)) {
		assert.Equal(t, expected, string(bs))
	}
	type binarySlice []byte
	var bs2 binarySlice
	if assert.NoError(t, cons.Consume(bytes.NewBufferString(expected), &bs2)) {
		assert.Equal(t, expected, string(bs2))
	}

	// passing in a nilslice wil result in an error
	var ns *[]byte
	assert.Error(t, cons.Consume(bytes.NewBufferString(expected), &ns))

	// passing in nil wil result in an error as well
	assert.Error(t, cons.Consume(bytes.NewBufferString(expected), nil))

	// a reader who results in an error, will make it fail
	assert.Error(t, cons.Consume(new(nopReader), &bu))
	assert.Error(t, cons.Consume(new(nopReader), &bs))

	// the readers can also not be nil
	assert.Error(t, cons.Consume(nil, &bs))
}

type binaryUnmarshalDummy struct {
	str string
}

func (b *binaryUnmarshalDummy) UnmarshalBinary(bytes []byte) error {
	if len(bytes) == 0 {
		return errors.New("no text given")
	}

	b.str = string(bytes)
	return nil
}

func TestByteStreamProducer(t *testing.T) {
	cons := ByteStreamProducer()
	expected := "the data for the stream to be sent over the wire"

	var rdr bytes.Buffer

	// can produce using a reader
	if assert.NoError(t, cons.Produce(&rdr, bytes.NewBufferString(expected))) {
		assert.Equal(t, expected, rdr.String())
		rdr.Reset()
	}

	// can produce using a binary marshaller
	if assert.NoError(t, cons.Produce(&rdr, &binaryMarshalDummy{expected})) {
		assert.Equal(t, expected, rdr.String())
		rdr.Reset()
	}

	// binary slices can also be used to produce
	if assert.NoError(t, cons.Produce(&rdr, []byte(expected))) {
		assert.Equal(t, expected, rdr.String())
		rdr.Reset()
	}

	// errors can also be used to produce
	if assert.NoError(t, cons.Produce(&rdr, errors.New(expected))) {
		assert.Equal(t, expected, rdr.String())
		rdr.Reset()
	}

	// structs can also be used to produce
	if assert.NoError(t, cons.Produce(&rdr, Error{Message: expected})) {
		assert.Equal(t, fmt.Sprintf(`{"message":%q}`, expected), rdr.String())
		rdr.Reset()
	}

	// struct pointers can also be used to produce
	if assert.NoError(t, cons.Produce(&rdr, &Error{Message: expected})) {
		assert.Equal(t, fmt.Sprintf(`{"message":%q}`, expected), rdr.String())
		rdr.Reset()
	}

	// slices can also be used to produce
	if assert.NoError(t, cons.Produce(&rdr, []string{expected})) {
		assert.Equal(t, fmt.Sprintf(`[%q]`, expected), rdr.String())
		rdr.Reset()
	}

	type binarySlice []byte
	if assert.NoError(t, cons.Produce(&rdr, binarySlice(expected))) {
		assert.Equal(t, expected, rdr.String())
		rdr.Reset()
	}

	// when binaryMarshal data is used, its potential error gets propagated
	assert.Error(t, cons.Produce(&rdr, new(binaryMarshalDummy)))
	// nil data should never be accepted either
	assert.Error(t, cons.Produce(&rdr, nil))
	// nil readers should also never be acccepted
	assert.Error(t, cons.Produce(nil, bytes.NewBufferString(expected)))
}

type binaryMarshalDummy struct {
	str string
}

func (b *binaryMarshalDummy) MarshalBinary() ([]byte, error) {
	if len(b.str) == 0 {
		return nil, errors.New("no text set")
	}

	return []byte(b.str), nil
}

type closingWriter struct {
	calledClose int64
	calledWrite int64
	b           bytes.Buffer
}

func (c *closingWriter) Close() error {
	atomic.AddInt64(&c.calledClose, 1)
	return nil
}

func (c *closingWriter) Write(p []byte) (n int, err error) {
	atomic.AddInt64(&c.calledWrite, 1)
	return c.b.Write(p)
}

func (c *closingWriter) String() string {
	return c.b.String()
}

type closingReader struct {
	calledClose int64
	calledRead  int64
	b           *bytes.Buffer
}

func (c *closingReader) Close() error {
	atomic.AddInt64(&c.calledClose, 1)
	return nil
}

func (c *closingReader) Read(p []byte) (n int, err error) {
	atomic.AddInt64(&c.calledRead, 1)
	return c.b.Read(p)
}

func TestBytestreamConsumer_Close(t *testing.T) {
	cons := ByteStreamConsumer(ClosesStream)
	expected := "the data for the stream to be sent over the wire"

	// can consume as a Writer
	var b bytes.Buffer
	r := &closingReader{b: bytes.NewBufferString(expected)}
	if assert.NoError(t, cons.Consume(r, &b)) {
		assert.Equal(t, expected, b.String())
		assert.EqualValues(t, 1, r.calledClose)
	}

	// can consume as a Writer
	cons = ByteStreamConsumer()
	b.Reset()
	r = &closingReader{b: bytes.NewBufferString(expected)}
	if assert.NoError(t, cons.Consume(r, &b)) {
		assert.Equal(t, expected, b.String())
		assert.EqualValues(t, 0, r.calledClose)
	}
}

func TestBytestreamProducer_Close(t *testing.T) {
	cons := ByteStreamProducer(ClosesStream)
	expected := "the data for the stream to be sent over the wire"

	// can consume as a Writer
	r := &closingWriter{}
	// can produce using a reader
	if assert.NoError(t, cons.Produce(r, bytes.NewBufferString(expected))) {
		assert.Equal(t, expected, r.String())
		assert.EqualValues(t, 1, r.calledClose)
	}

	cons = ByteStreamProducer()
	r = &closingWriter{}
	// can produce using a reader
	if assert.NoError(t, cons.Produce(r, bytes.NewBufferString(expected))) {
		assert.Equal(t, expected, r.String())
		assert.EqualValues(t, 0, r.calledClose)
	}
}
