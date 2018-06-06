package client

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newCountingReader(rdr io.Reader, readOnce bool) *countingReadCloser {
	return &countingReadCloser{
		rdr:      rdr,
		readOnce: readOnce,
	}
}

type countingReadCloser struct {
	rdr         io.Reader
	readOnce    bool
	readCalled  int
	closeCalled int
}

func (c *countingReadCloser) Read(b []byte) (int, error) {
	c.readCalled++
	if c.readCalled > 1 && c.readOnce {
		return 0, io.EOF
	}
	return c.rdr.Read(b)
}

func (c *countingReadCloser) Close() error {
	c.closeCalled++
	return nil
}

func TestDrainingReadCloser(t *testing.T) {
	rdr := newCountingReader(bytes.NewBufferString("There are many things to do"), false)
	prevDisc := ioutil.Discard
	disc := bytes.NewBuffer(nil)
	ioutil.Discard = disc
	defer func() { ioutil.Discard = prevDisc }()

	buf := make([]byte, 5)
	ts := &drainingReadCloser{rdr: rdr}
	ts.Read(buf)
	ts.Close()
	assert.Equal(t, "There", string(buf))
	assert.Equal(t, " are many things to do", disc.String())
	assert.Equal(t, 3, rdr.readCalled)
	assert.Equal(t, 1, rdr.closeCalled)
}

func TestDrainingReadCloser_SeenEOF(t *testing.T) {
	rdr := newCountingReader(bytes.NewBufferString("There are many things to do"), true)
	prevDisc := ioutil.Discard
	disc := bytes.NewBuffer(nil)
	ioutil.Discard = disc
	defer func() { ioutil.Discard = prevDisc }()

	buf := make([]byte, 5)
	ts := &drainingReadCloser{rdr: rdr}
	ts.Read(buf)
	_, err := ts.Read(nil)
	assert.Equal(t, io.EOF, err)
	ts.Close()
	assert.Equal(t, string(buf), "There")
	assert.Equal(t, disc.String(), "")
	assert.Equal(t, 2, rdr.readCalled)
	assert.Equal(t, 1, rdr.closeCalled)
}
