/*
Copyright 2011 Google Inc.

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
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sync"

	"go4.org/readerutil"
)

var (
	ErrNegativeSubFetch         = errors.New("invalid negative subfetch parameters")
	ErrOutOfRangeOffsetSubFetch = errors.New("subfetch offset greater than blob size")
)

// Fetcher is the minimal interface for retrieving a blob from storage.
// The full storage interface is blobserver.Storage.
type Fetcher interface {
	// Fetch returns a blob.  If the blob is not found then
	// os.ErrNotExist should be returned for the error (not a wrapped
	// error with a ErrNotExist inside)
	//
	// The contents are not guaranteed to match the digest of the
	// provided Ref (e.g. when streamed over HTTP). Paranoid
	// callers should verify them.
	//
	// The caller must close blob.
	Fetch(Ref) (blob io.ReadCloser, size uint32, err error)
}

// A SubFetcher is a Fetcher that can retrieve part of a blob.
type SubFetcher interface {
	// SubFetch returns part of a blob.
	// The caller must close the returned io.ReadCloser.
	// The Reader may return fewer than 'length' bytes. Callers should
	// check. The returned error should be: ErrNegativeSubFetch if any of
	// offset or length is negative, or os.ErrNotExist if the blob
	// doesn't exist, or ErrOutOfRangeOffsetSubFetch if offset goes over
	// the size of the blob.
	SubFetch(ref Ref, offset, length int64) (io.ReadCloser, error)
}

func NewSerialFetcher(fetchers ...Fetcher) Fetcher {
	return &serialFetcher{fetchers}
}

func NewSimpleDirectoryFetcher(dir string) *DirFetcher {
	return &DirFetcher{dir, "camli"}
}

type serialFetcher struct {
	fetchers []Fetcher
}

func (sf *serialFetcher) Fetch(r Ref) (file io.ReadCloser, size uint32, err error) {
	for _, fetcher := range sf.fetchers {
		file, size, err = fetcher.Fetch(r)
		if err == nil {
			return
		}
	}
	return
}

type DirFetcher struct {
	directory, extension string
}

func (df *DirFetcher) Fetch(r Ref) (file io.ReadCloser, size uint32, err error) {
	fileName := fmt.Sprintf("%s/%s.%s", df.directory, r.String(), df.extension)
	var stat os.FileInfo
	stat, err = os.Stat(fileName)
	if err != nil {
		return
	}
	if stat.Size() > math.MaxUint32 {
		err = errors.New("file size too big")
		return
	}
	file, err = os.Open(fileName)
	if err != nil {
		return
	}
	size = uint32(stat.Size())
	return
}

// NewLazyReadSeekCloser returns a ReadSeekCloser that does no work
// until one of its Read, Seek, or Close methods is called, but then
// fetches the ref from src. Any fetch error is returned in the Read,
// Seek, or Close call.
func NewLazyReadSeekCloser(src Fetcher, br Ref) readerutil.ReadSeekCloser {
	return &lazyReadSeekCloser{src: src, br: br}
}

type lazyReadSeekCloser struct {
	once sync.Once // guards init
	src  Fetcher
	br   Ref

	// after init, exactly one is set:
	err error
	rsc readerutil.ReadSeekCloser
}

func (r *lazyReadSeekCloser) init() {
	b, err := FromFetcher(r.src, r.br)
	if err != nil {
		r.err = err
		return
	}
	r.rsc = b.Open()
}

func (r *lazyReadSeekCloser) Read(p []byte) (n int, err error) {
	r.once.Do(r.init)
	if r.err != nil {
		return 0, r.err
	}
	return r.rsc.Read(p)
}

func (r *lazyReadSeekCloser) Seek(offset int64, whence int) (int64, error) {
	r.once.Do(r.init)
	if r.err != nil {
		return 0, r.err
	}
	return r.rsc.Seek(offset, whence)
}

func (r *lazyReadSeekCloser) Close() error {
	r.once.Do(r.init)
	if r.err != nil {
		return r.err
	}
	return r.rsc.Close()
}

// ReaderAt returns an io.ReaderAt of br, fetching against sf.
func ReaderAt(sf SubFetcher, br Ref) io.ReaderAt {
	return readerAt{sf, br}
}

type readerAt struct {
	sf SubFetcher
	br Ref
}

func (ra readerAt) ReadAt(p []byte, off int64) (n int, err error) {
	rc, err := ra.sf.SubFetch(ra.br, off, int64(len(p)))
	if err != nil {
		return 0, err
	}
	defer rc.Close()
	return io.ReadFull(rc, p)
}
