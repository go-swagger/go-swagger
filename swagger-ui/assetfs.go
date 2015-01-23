package swaggerui

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var (
	fileTimestamp = time.Now()
)

// FakeFile implements os.FileInfo interface for a given path and size
type FakeFile struct {
	// Path is the path of this file
	Path string
	// Dir marks of the path is a directory
	Dir bool
	// Len is the length of the fake file, zero if it is a directory
	Len int64
}

func (f *FakeFile) Name() string {
	_, name := filepath.Split(f.Path)
	return name
}

func (f *FakeFile) Mode() os.FileMode {
	mode := os.FileMode(0644)
	if f.Dir {
		return mode | os.ModeDir
	}
	return mode
}

func (f *FakeFile) ModTime() time.Time {
	return fileTimestamp
}

func (f *FakeFile) Size() int64 {
	return f.Len
}

func (f *FakeFile) IsDir() bool {
	return f.Mode().IsDir()
}

func (f *FakeFile) Sys() interface{} {
	return nil
}

// AssetFile implements http.File interface for a no-directory file with content
type AssetFile struct {
	*bytes.Reader
	io.Closer
	FakeFile
}

func NewAssetFile(name string, content []byte) *AssetFile {
	return &AssetFile{
		bytes.NewReader(content),
		ioutil.NopCloser(nil),
		FakeFile{name, false, int64(len(content))}}
}

func (f *AssetFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, errors.New("not a directory")
}

func (f *AssetFile) Stat() (os.FileInfo, error) {
	return f, nil
}

// AssetDirectory implements http.File interface for a directory
type AssetDirectory struct {
	AssetFile
	ChildrenRead int
	Children     []os.FileInfo
}

func NewAssetDirectory(name string, children []string, fs *AssetFS) *AssetDirectory {
	fileinfos := make([]os.FileInfo, 0, len(children))
	for _, child := range children {
		_, err := fs.AssetDir(filepath.Join(name, child))
		fileinfos = append(fileinfos, &FakeFile{child, err == nil, 0})
	}
	return &AssetDirectory{
		AssetFile{
			bytes.NewReader(nil),
			ioutil.NopCloser(nil),
			FakeFile{name, true, 0},
		},
		0,
		fileinfos}
}

func (f *AssetDirectory) Readdir(count int) ([]os.FileInfo, error) {
	if count <= 0 {
		return f.Children, nil
	}
	if f.ChildrenRead+count > len(f.Children) {
		count = len(f.Children) - f.ChildrenRead
	}
	rv := f.Children[f.ChildrenRead : f.ChildrenRead+count]
	f.ChildrenRead += count
	return rv, nil
}

func (f *AssetDirectory) Stat() (os.FileInfo, error) {
	return f, nil
}

// AssetFS implements http.FileSystem, allowing
// embedded files to be served from net/http package.
type AssetFS struct {
	// Asset should return content of file in path if exists
	Asset func(path string) ([]byte, error)
	// AssetDir should return list of files in the path
	AssetDir func(path string) ([]string, error)
	// Prefix would be prepended to http requests
	Prefix string

	Strip string
}

func (fs *AssetFS) Open(name string) (http.File, error) {
	if name == fs.Strip {
		name = fs.Strip + "/index.html"
	}

	// name = path.Join(fs.Prefix, name)
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	if name == "" || name == fs.Strip {
		name = fs.Strip + "/index.html"
	}

	if strings.Count(name, "/") > 1 {
		if children, err := fs.AssetDir(name); err == nil {
			return NewAssetDirectory(name, children, fs), nil
		}

	}

	b, err := fs.Asset(name)
	if err != nil {
		return nil, err
	}
	return NewAssetFile(path.Join(fs.Prefix, name), b), nil
}

// Middleware creates a middleware to serve swagger-ui at /swagger-ui
func Middleware(next http.Handler) http.Handler {
	return middlewareAt("/swagger-ui", next)
}

// MiddlewareAt creates a middleware to serve swagger ui at the specified basePath
func middlewareAt(basePath string, next http.Handler) http.Handler {
	assetFS := func() *AssetFS {
		for k := range _bintree.Children {
			return &AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: basePath + "/" + k, Strip: basePath}
		}
		panic("unreachable")
	}

	fileServer := http.FileServer(assetFS())

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/favico") {
			rw.WriteHeader(http.StatusNotFound)
		}

		if strings.HasPrefix(r.URL.Path, basePath) {
			fileServer.ServeHTTP(rw, r)
			return
		}

		if next == nil {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			next.ServeHTTP(rw, r)
		}
	})
}
