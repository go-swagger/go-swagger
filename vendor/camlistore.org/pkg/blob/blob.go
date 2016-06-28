/*
Copyright 2014 The Camlistore Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package blob

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"unicode/utf8"

	"camlistore.org/pkg/constants"
	"go4.org/readerutil"
)

// Blob represents a blob. Use the methods Size, SizedRef and
// Open to query and get data from Blob.
type Blob struct {
	ref       Ref
	size      uint32
	newReader func() readerutil.ReadSeekCloser
	mem       []byte // if in memory
}

// NewBlob constructs a Blob from its Ref, size and a function that
// returns an io.ReadCloser from which the blob can be read. Any error
// in the function newReader when constructing the io.ReadCloser should
// be returned upon the first call to Read or Close.
func NewBlob(ref Ref, size uint32, newReader func() readerutil.ReadSeekCloser) *Blob {
	return &Blob{
		ref:       ref,
		size:      size,
		newReader: newReader,
	}
}

// Size returns the size of the blob (in bytes).
func (b *Blob) Size() uint32 {
	return b.size
}

// SizedRef returns the SizedRef corresponding to the blob.
func (b *Blob) SizedRef() SizedRef {
	return SizedRef{b.ref, b.size}
}

// Ref returns the blob's reference.
func (b *Blob) Ref() Ref { return b.ref }

// Open returns an io.ReadCloser that can be used to read the blob
// data. The caller must close the io.ReadCloser when finished.
func (b *Blob) Open() readerutil.ReadSeekCloser {
	return b.newReader()
}

// ValidContents reports whether the hash of blob's content matches
// its reference.
func (b *Blob) ValidContents() bool {
	h := b.ref.Hash()
	if b.mem != nil {
		h.Write(b.mem)
	} else {
		rc := b.Open()
		defer rc.Close()
		_, err := io.Copy(h, rc)
		if err != nil {
			return false
		}
	}
	return b.ref.HashMatches(h)
}

// IsUTF8 reports whether the blob is entirely UTF-8.
func (b *Blob) IsUTF8() bool {
	if b.mem != nil {
		return utf8.Valid(b.mem)
	}
	rc := b.Open()
	defer rc.Close()
	slurp, err := ioutil.ReadAll(rc)
	if err != nil {
		return false
	}
	return utf8.Valid(slurp)
}

// A reader reads a blob's contents.
// It adds a no-op Close method to a *bytes.Reader.
type reader struct {
	*bytes.Reader
}

func (reader) Close() error { return nil }

// FromFetcher fetches br from fetcher and slurps its contents to
// memory. It does not validate the blob's digest.  Use the
// Blob.ValidContents method for that.
func FromFetcher(fetcher Fetcher, br Ref) (*Blob, error) {
	rc, size, err := fetcher.Fetch(br)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return FromReader(br, rc, size)
}

// FromReader slurps the given blob from r to memory.
// It does not validate the blob's digest.  Use the
// Blob.ValidContents method for that.
func FromReader(br Ref, r io.Reader, size uint32) (*Blob, error) {
	if size > constants.MaxBlobSize {
		return nil, fmt.Errorf("blob: %v with reported size %d is over max size of %d", br, size, constants.MaxBlobSize)
	}
	buf := make([]byte, size)
	if n, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("blob: after reading %d bytes of %v: %v", n, br, err)
	}
	n, _ := io.CopyN(ioutil.Discard, r, 1)
	if n > 0 {
		return nil, fmt.Errorf("blob: %v had more than reported %d bytes", br, size)
	}
	opener := func() readerutil.ReadSeekCloser {
		return reader{bytes.NewReader(buf)}
	}
	b := NewBlob(br, uint32(size), opener)
	b.mem = buf
	return b, nil
}
