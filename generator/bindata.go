// Code generated by go-bindata.
// sources:
// templates/model.gotmpl
// templates/modelvalidator.gotmpl
// templates/server/builder.gotmpl
// templates/server/main.gotmpl
// templates/server/operation.gotmpl
// templates/server/parameter.gotmpl
// DO NOT EDIT!

package generator

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesModelGotmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x84\x92\x3f\x4f\xc3\x30\x10\xc5\xf7\x7c\x8a\x53\x54\x24\x90\x20\xdd\x2b\x31\x51\x06\x24\x40\x48\x74\x60\xac\x49\x2e\x89\xa9\xff\x04\xfb\xd2\x12\x59\xfe\xee\xd8\x4e\xd2\x3f\x2c\x6c\xb1\xfd\xbb\x77\xf7\xde\xc5\xb9\x0a\x6b\xae\x10\x72\xa9\x2b\x14\x9d\xd1\x1d\x1a\x1a\x72\xef\x33\xe7\x78\x0d\xc5\x5a\x97\xef\x64\xb8\x6a\xbc\x77\xee\xf2\x84\xaa\x4a\x58\xf1\x36\x55\xbd\x32\x89\xde\x43\xe4\x18\xb1\xcd\xd0\xc5\xd3\xf6\xcb\x6a\xb5\xca\x23\xc6\x0c\x93\x23\x93\x8f\xe2\x1f\x2f\xcf\x53\xcd\x8f\x14\x89\x39\xde\x4c\xf2\xf9\x36\x9b\x1b\x65\x1d\x2b\x77\xac\x41\x48\x52\xe9\x33\xde\x2e\x97\xb0\x69\xb9\x85\x9a\x0b\x84\x03\xb3\xd0\xa0\x42\xc3\x08\x2b\xf8\x1c\x80\x5a\x04\x7b\x60\x4d\x83\x06\x48\x6b\x51\x44\xfe\xb1\xe2\x14\x3c\x84\xc7\xb9\x4e\xf2\xa6\x25\x08\xee\xf7\x08\x75\x4f\x49\xaa\x45\x05\x83\xee\xc1\xe0\x9d\xe9\xd5\x85\xd2\xdc\x02\x4a\x2d\x25\x53\x55\x36\x85\xf5\x24\x3b\x6d\xc8\x7a\xcf\xc7\x0f\xb8\xce\x20\x8c\x6b\x98\x0a\x63\x17\x6b\xac\x59\x2f\xe8\x08\x39\xd7\x85\x28\xa9\x86\xfc\xea\x3b\x87\x22\x98\x89\xf0\x68\xf6\x54\xb6\xd8\xe1\x70\x0b\x8b\x3d\x13\x3d\xc2\xea\xfe\xac\x89\x73\xf1\x2d\x05\x0e\xe7\x4a\x23\x7b\x21\x77\x73\x4a\xf1\x9f\xb5\xce\x20\x85\xed\xc5\xa4\x1f\x04\xb3\x76\x5a\x92\x25\xd3\x97\x04\x2e\x3b\x5a\x9a\x36\xcf\xd1\xa6\x4a\x42\xd9\x89\x98\xcb\x9f\xbf\x29\x99\x9b\x27\xf0\xd9\x6f\x00\x00\x00\xff\xff\x35\x7f\x74\xf7\x75\x02\x00\x00")

func templatesModelGotmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesModelGotmpl,
		"templates/model.gotmpl",
	)
}

func templatesModelGotmpl() (*asset, error) {
	bytes, err := templatesModelGotmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/model.gotmpl", size: 629, mode: os.FileMode(420), modTime: time.Unix(1435797486, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _templatesModelvalidatorGotmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xcc\x56\xdd\x6e\xdb\x36\x14\xbe\xf7\x53\x9c\x09\x19\x10\xaf\x89\xd2\x01\xc3\x2e\xd6\x65\x40\xd1\x7a\x68\x80\x2e\x0d\x9a\x6e\x37\xc3\x80\x32\x32\x25\x73\x91\x48\x85\xa4\x9c\x78\x82\xde\x7d\x87\xa4\x28\x4b\xb2\xe4\x58\x0e\x06\xec\x4a\xd4\xf9\xe3\xc7\xf3\xf3\x91\x65\xb9\xa4\x31\xe3\x14\x82\x5c\xb2\x8c\x69\xb6\xa6\x6b\x92\xb2\x25\xd1\x42\x06\x55\x35\x2b\x4b\x16\x43\xf8\x99\x3e\x14\x4c\xd2\x25\x0a\xf0\x97\x4a\x09\x3f\x5d\x42\x6d\x47\x1b\xed\x69\x59\x86\x37\x44\xaf\xaa\xea\x0c\x02\x5c\x7f\x14\x11\xd1\x4c\xf0\xaa\x0a\xce\x00\xff\xff\x20\x69\x41\x17\x4f\xb9\xa4\x4a\x59\xf1\xfc\x8d\x8d\xf5\xcd\x25\x70\x96\x42\x39\x03\x90\x54\x17\x92\x1b\xe9\xcc\xec\x4d\xf9\xb2\xc1\xf0\x1b\xe3\x1f\x29\x4f\x4c\xf8\x21\x10\x8d\x7a\x32\x0a\x2b\x6d\x45\x9f\x86\x8a\x3c\xed\x45\xe5\xd5\x47\xa2\xda\x46\x9f\x84\x0a\x77\xd2\x54\xf2\x61\x4c\xb5\xf2\x08\x44\x5f\x9d\x8b\x0b\xfd\x75\x6a\xf5\x58\x56\x64\xa3\xb5\x33\xca\xbd\x88\xe2\x54\x10\xfd\xe3\x0f\xa7\x83\x7d\xe4\x4b\xe8\xb6\xb0\x7f\x8b\xa7\x28\x2d\x14\xb6\x73\x23\x9e\x5a\xd7\x3d\x78\x9d\xf2\xa5\x78\xfd\x16\x3d\xbc\x5e\x3c\x0d\x6f\x91\x6a\x96\xa7\xf4\x53\x3c\x02\xb9\xd1\xbf\x14\x75\x6b\xa3\x49\x08\x17\x7c\x2c\x9d\x46\x73\xdc\x7c\xb8\x98\x07\xc3\xf0\x5f\x4f\x79\x51\xa1\xb4\xc8\x62\x21\x33\xa2\x3b\xac\x37\x00\xf2\x57\x6b\xf5\x4c\xfa\x8c\xc0\x19\xda\x5f\xa5\x25\xe3\xc9\x58\x32\xdd\xbe\xea\x60\xf4\x1e\xb5\x4a\x59\x34\x44\xd2\xd7\x94\x2e\xd5\x2d\xfb\x87\x5a\x09\x82\x94\x24\xbb\x26\x19\xfe\x1a\xa1\x39\x0c\xe3\xa6\xb6\x29\xe5\xc3\x90\xe6\xbb\x33\x7b\xa5\x69\xa6\x46\x87\xd6\x6a\x9f\xab\x5c\x0f\x87\x1f\xd5\x3a\xf2\xd4\xa1\xdc\x07\xa8\xd6\x1e\x05\xa8\x89\x3c\x09\xd0\xef\x9c\x3d\x14\x74\x0f\xa6\x96\xc1\x7f\x7b\x3b\xfe\x0f\xa6\xcb\xc0\xb8\xc5\x7e\x4f\xe9\x6d\xb4\xa2\x19\xb9\x35\x7d\x0a\xa8\xba\xb8\x00\x65\xe5\xa0\xac\x62\x70\xc7\x19\x8e\x03\x30\x83\xfc\xf5\x1b\xfc\xfe\x0c\xa3\x6d\x8a\xea\x57\xaf\x10\x48\x59\x4a\xc2\x13\x0a\xa1\xcf\x3f\x60\x60\x5c\xe6\x29\x1e\xdb\xbc\x67\x44\x4e\xa5\xde\x6c\x27\x05\xc2\x16\x0b\xd8\x55\xaa\xa8\xc3\xc7\x85\xde\xc5\x78\x53\x47\x70\xbd\xf2\xc2\xfd\x5c\x7e\xde\x2e\x97\xcc\x24\x9e\xa4\xdb\x20\xcd\xc1\x71\x4b\x2b\xc5\x2b\xbf\xaa\x4c\x12\x30\x0b\x76\x5a\xe7\x70\xde\x55\x1a\xc1\xf7\xc6\xc2\x26\x02\xe0\x20\x24\x00\xad\x33\x23\x98\xd1\x04\xc3\x2f\xdd\xdd\x7a\x45\x17\x52\xf5\xcf\x71\x2d\xf4\xdb\x34\x15\x8f\xf8\x06\x0c\x86\x42\x06\x3b\x6d\x37\x1f\x24\xe6\x3e\xd5\x89\xbb\xbf\x69\xd4\xa5\x66\x2c\x96\xa3\x6d\x70\x4a\x03\xf5\x3d\xd1\xe4\xcb\x26\xa7\x9d\x01\x18\xc2\x61\x24\x76\x2a\x4e\x8f\x25\xdf\xdd\xdc\x36\xb5\xbd\x52\x37\xfe\x09\x5d\x55\xdd\x7a\xec\xbc\xac\x7d\x6b\x60\x31\xc0\xf9\xbe\xb3\x87\x72\x97\x07\x3e\xb0\xba\x11\x86\x2f\xaa\x81\x20\x82\x6b\x82\x38\x7b\xee\xbd\x1b\x63\xc8\x0f\x4d\xe9\xd3\x27\x9b\xd1\xae\x6f\xbf\x04\xc6\xb9\x5f\xb0\x9c\x44\xf7\x04\x27\xc3\xb2\x8c\x5d\xa2\xd0\x54\xea\xcb\x8a\x29\x88\x19\x4e\xd5\x23\x51\x90\x50\x44\x86\x41\x97\x70\xb7\x01\xbd\xc2\x51\x7b\x24\x49\x42\x25\x68\x21\xd2\xd0\xd8\x2f\x4c\x57\xf1\x04\x95\xde\x2f\x63\xc9\x4a\x03\x66\x7d\x4d\x21\x2e\xb4\x0d\xb5\xa2\x1c\x36\xa2\xc0\x52\x9d\xcb\x82\x77\x22\xf9\x2d\x20\x12\x59\x46\xf8\x72\x36\x63\x59\x2e\xa4\x86\x53\x2c\x6d\x90\x30\xbd\x2a\xee\x42\xd4\x5d\x24\xe2\xbc\xf6\x69\x2f\x5d\x77\x07\x07\xd9\xe2\xe5\x1e\x67\xfa\x30\xdb\x95\xd6\xf9\x3d\xd3\x17\x9e\x94\x83\x99\x25\x8f\x9a\x4f\xde\xd3\x98\xe0\xab\xea\xca\x22\x55\x26\xbd\xd8\x31\x5c\xc7\x10\x7c\xfb\xe0\x47\xd7\xa7\x7a\xeb\x76\x72\x4f\x37\x67\x70\xb2\x36\x3d\x6e\xfa\x3d\x6c\xf9\x1b\x9d\x99\xdb\x12\xda\x91\x9c\x6d\x27\xdc\xdc\x96\xc9\x4f\x45\x73\x69\x28\x57\x01\xac\xe7\x87\x02\xd3\xf8\x2e\x25\x4a\xd5\x2c\x18\x17\x3c\x02\xc3\x1b\x9f\x69\x44\xb1\xa7\xa5\x93\xc3\x77\x28\x6a\xd9\xcd\xa1\x3f\x6a\xe0\x12\x86\x7e\x09\xc3\xe5\x66\xee\xb8\xc4\x0e\x9d\x1b\xa0\x0f\x44\xd5\x4e\x38\xac\x8e\x1b\xd7\x44\x62\x99\x15\xfc\xf9\x97\x35\xee\xa4\xad\xa6\x67\x46\x3d\x17\x8f\xc4\xe8\x50\x42\x17\x75\xe8\x0f\xbc\xc3\xf6\xe3\x04\x01\x16\xd0\x25\x90\x3c\xc7\x1c\x9e\xe2\xcf\x99\x31\x99\x5b\x7e\xed\x16\xca\xad\x1c\x04\xc3\xb6\x68\x6b\xc8\xf5\x75\x13\xa7\x4d\xa9\x66\x06\x85\x62\x9a\x6e\xf1\x2f\x8c\xc6\x78\x85\x61\xb8\x1b\xbf\x76\x47\x5c\x96\xa9\xe0\x24\xf2\xe9\xb7\xed\xd0\x14\x03\xda\x77\x57\x27\x69\x23\x29\xdb\x57\xe3\xed\x26\xa6\xc6\xcf\xa6\x6f\x7f\xd1\x9f\xbd\xb2\x76\x4f\xd9\x25\x9e\x7f\x03\x00\x00\xff\xff\xbd\x9d\x11\x8f\xbe\x10\x00\x00")

func templatesModelvalidatorGotmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesModelvalidatorGotmpl,
		"templates/modelvalidator.gotmpl",
	)
}

func templatesModelvalidatorGotmpl() (*asset, error) {
	bytes, err := templatesModelvalidatorGotmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/modelvalidator.gotmpl", size: 4286, mode: os.FileMode(420), modTime: time.Unix(1435797486, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _templatesServerBuilderGotmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xc4\x59\x5b\x6f\x1b\xbb\x11\x7e\xae\x7e\xc5\x40\x38\xa7\xd8\x0d\x74\xd6\x7e\x2c\x0c\xb8\x80\x1b\x9f\x22\xee\x25\x35\xec\xa0\x7d\x30\x82\x82\xde\xa5\x24\xc2\x7b\x3b\x24\xd7\xaa\x2a\xe8\xbf\x77\x86\x97\x5d\xee\x45\x8a\x14\x3b\xa8\x11\x04\x12\x39\x9c\xf9\xe6\x3e\xa4\x6a\x96\xbe\xb0\x15\x87\xdd\x2e\xb9\xb7\x1f\xf7\xfb\xd9\xec\xe2\x02\xbe\xac\x85\x82\xa5\xc8\x39\x6c\x98\x82\x15\x2f\xb9\x64\x9a\x67\xf0\xbc\x05\xbd\xe6\xa0\x36\x6c\xb5\xe2\x12\x74\x55\xe5\x09\xd1\xff\x9a\x09\x2d\xca\x15\x6e\xfa\x73\x85\x58\xad\x35\xd4\xb2\x7a\xe5\xb0\x6c\xb4\x61\xb5\xe6\x25\x6c\xab\x06\x24\xff\x45\x36\x65\x8f\x93\x17\x01\x69\x55\x14\xac\xcc\x66\x33\x51\xd4\x95\xd4\x10\xcd\x00\xe6\x4a\x4b\xe4\xae\xe6\xf4\xb9\xe4\xfa\x62\xad\x75\x6d\xbe\xac\x84\x5e\x37\xcf\x09\x1e\xba\x58\x55\xbf\x38\x66\xe1\x47\xa2\x7c\x11\xfa\x34\x62\x55\xf3\xf4\x44\x4a\x2d\x97\xc5\x89\x5c\x1d\x84\x8b\x42\x64\x59\xce\x37\x4c\xf2\xf3\xce\x29\x9e\x36\x52\xe8\xed\x7c\x86\xc7\x76\x3b\xc9\x4a\xf4\x58\x72\xcb\x97\xac\xc9\xf5\x9d\xb1\x92\xda\xef\x77\xbb\x1a\x6d\xa4\x97\x30\xff\xf9\xb7\x39\x24\xe8\x47\x22\xe6\x65\xe6\x3e\xd9\x63\x3f\xbd\xf0\xed\x02\x7e\x7a\x65\x79\xc3\xe1\xea\x1a\x92\xe0\x3c\xed\xed\xf7\x48\x0a\x21\x27\x4b\xdb\x63\x17\x9b\x10\xf9\xcc\x37\x18\x36\x37\x75\xfd\x99\x15\xb8\x7f\x73\x7f\x07\xa9\xe4\xe8\x42\x05\x0c\x4a\xbe\x81\x70\x17\x44\xa9\x34\x2b\x53\x3e\x5b\x36\x65\x3a\x71\x36\x22\xdb\xc3\x07\xfa\x3f\xb9\xad\xd2\xa6\xe0\xa5\x8e\xe1\xc3\x50\xc2\xce\xc0\x48\x1e\x78\xca\xc5\x2b\x97\x8e\x39\x2a\xf2\xfb\x01\x25\x11\x02\x10\xbb\x2b\xf0\x9f\x16\x66\x6d\x8d\xd1\x95\x73\xa9\xae\xa0\x60\x2f\x3c\x2a\x58\xfd\x64\xc3\xeb\x2b\x19\x3c\xf9\x64\xb7\x63\x4b\xbc\xac\x64\xc1\x34\xd2\x82\xf5\xb8\x37\xbb\xdd\xcd\xec\x97\x8f\x55\xa9\x10\x30\x52\xcd\x11\xc5\x6d\x7f\x71\xbf\x9f\xf7\x88\xef\x65\x95\x35\xe9\x80\xd8\x2f\x3a\xe2\x3d\x79\x5a\x72\xdd\xc8\x72\xac\xed\xcc\x66\xe8\xc8\x32\xbb\xe4\xae\x5c\x56\xc8\x51\xa5\x52\xd4\x5a\x54\x25\xd2\xea\x6d\xcd\x47\xa4\xa8\x4a\x93\x6a\x63\x4b\x63\xf5\xe0\xaf\xef\x00\x24\x48\xab\x52\xf3\xff\xe8\x8e\xa0\x8b\xe2\xe4\xa3\xdd\x9b\x75\x36\xf5\x54\x07\x8c\x3a\x6b\x0d\xda\xf2\x73\x66\x7d\xe0\x2b\x81\x1f\xb7\xb3\x91\x51\xc1\xf2\x99\x8d\x0c\xd8\x6d\xb4\x39\xd1\xd9\xdc\x1a\xe8\x63\xce\x94\xb2\x7a\xbb\x2d\x89\x66\x25\x49\x84\x95\x91\x72\x76\x11\x51\xe1\x57\x72\xc8\xdf\x79\x26\xd8\x17\xb4\x1a\xba\x02\x6b\x58\xc1\x81\x4c\x68\xa3\x6e\x8a\x9d\x4b\x52\x2f\x5a\x4e\x66\x5d\xd2\xf9\x77\x04\xcc\x6d\xf5\x81\xd5\x7e\xf1\x6c\x60\x2d\x3b\x0f\xcc\x2f\x4c\x03\x7b\x74\xb5\x05\xe3\x50\x94\x82\x82\x46\x39\x02\xb1\xc4\xe2\xa0\xfe\xc4\x94\x48\x6f\x1a\xbd\x9e\x40\x4e\xcb\x3d\xd4\x94\xda\xc4\x02\x0b\x3b\xd3\xa0\x31\xbb\x14\x34\x8a\xcb\x12\xc9\x01\x23\x00\x6a\x3c\xbb\xa9\x64\x66\xbe\xd8\xf8\xb6\xda\x8a\x32\x15\x35\xcb\x51\x30\x4a\x11\xd8\x36\xb8\xa4\x40\xc1\x4d\x94\x81\x81\x28\x52\x66\x18\x6f\xb0\x66\xc2\x33\x61\x32\x3b\x23\xed\x0d\x24\x82\x11\xd9\xe0\x58\xb8\x20\x89\x21\xa2\x52\x72\xef\x05\xed\xf7\x0b\xe0\x52\x56\x32\xee\xcc\xe2\x55\xc6\x0c\xf9\x2b\xdf\xbe\x45\x67\x86\x7d\xf1\x05\x5b\xdd\xf7\x6a\x89\x0a\x62\xab\xad\x88\x01\xb0\x5a\x00\xd6\x65\x82\xe1\x8a\x1d\xb5\x54\x91\x21\x81\xb0\x1d\x14\x77\x1e\xab\x46\xa6\xbe\x46\x1f\xb3\xc7\x29\x76\x98\x0e\x94\x7f\xd4\xd4\x9e\x6d\x7c\x8c\xac\xe2\xd2\x1b\x14\xc7\xcc\x26\x4c\x95\xa7\xf6\x95\xc1\x04\xb2\x43\xfb\xa9\xc1\x06\x1f\x9c\xee\xa8\xdb\xc0\x6b\x87\x91\x69\x39\xe1\xb8\x92\x4c\x92\x58\x25\x72\x75\x8c\xc5\xa1\x53\x23\x23\xa0\xbe\x8f\x5c\xbe\xf2\x5f\xc9\x52\x80\x03\x4e\xca\xf2\x1c\x1d\x60\xe6\x19\xf4\x11\xf7\xeb\xd2\x16\xea\x6c\x41\xaa\x4a\x4e\x4b\xcc\x97\x2d\x6f\x09\xcb\xef\xb9\xd1\x66\x12\x4a\xf1\x38\x5a\x8d\x3e\x4b\xa8\x36\x2e\xc2\x69\x8a\x42\xba\x40\xa8\xe9\x45\xe4\x47\x53\x4e\x1f\xb8\xaa\xd1\x13\xfc\x5f\x98\xba\x5c\x2e\xe0\x83\x5b\xfd\xad\xe1\x4a\xb7\x1e\xb5\x7d\xe2\x91\xeb\xdb\x61\xe1\xf4\x6e\xf2\xd0\x6a\xbf\x53\x50\x91\xb1\x85\xc5\xf4\xe9\x68\xdc\x6b\x87\x2d\x39\x9e\x90\x10\x15\xbe\x58\xb5\xf9\xb7\x9b\xfd\x6e\xc4\x2b\x19\x56\xf4\x6b\x68\x0f\x8e\xd0\xb7\xfd\xc0\x67\x54\xa8\x40\xea\x37\xdf\xa8\x80\x17\x72\xa6\x02\x2d\xb6\xb1\x02\x43\xdb\x4f\xa1\x7f\x9b\xf9\x87\xb6\x8f\x1d\x64\x42\x7c\x68\x86\x18\x5a\xbe\x0f\xf6\x07\x9a\x7a\x68\xe7\x73\xc0\xfa\x43\x0e\xec\x9f\xdd\x24\x11\x82\xf4\x95\x19\x93\xd3\xf1\x75\xf3\xc6\x19\x10\x1d\x5f\x0b\x2d\x9c\x4d\x8e\x62\xf4\x72\x2c\xb6\x07\x87\xc3\xf2\xea\xcf\x1c\x8d\xd2\x55\xe1\x70\x01\x0e\xd7\x22\x63\xba\x92\x67\x00\xec\x33\x8f\x4c\x77\xf5\xed\xce\xb1\x75\xc8\x2d\xc5\xa2\x93\xe2\x37\xfe\xe9\x17\xe2\xe9\x89\xda\xab\x93\xdc\x64\x99\x11\xe0\x39\x07\xbc\x7c\x81\x71\xbc\xb8\xdf\xe1\xa1\x2b\x5c\xcf\x08\x7a\x55\xa8\xcb\x19\x4a\x7b\x29\xe8\x16\x5b\x6e\x09\xf7\x2b\x93\xd0\x94\x81\xd3\x9f\xbe\x1e\x1b\x0a\x71\x15\x9b\xcb\x58\xd9\x03\xa3\xdd\xf5\x35\x94\x22\x07\x7b\x93\xe8\x89\xb9\xc6\xb6\x5c\x63\x77\x88\xc2\xd5\x85\x19\xd3\x26\x18\xcd\x63\x33\xd2\x7f\x63\x30\x3c\x0d\x5c\x3b\xde\xbd\x15\x9c\x67\x74\x0c\xdc\xa1\xe1\xf0\x04\x9c\x66\xf0\x78\x2b\x46\x62\x72\x0c\x5f\x38\x93\x9c\x06\xcb\x77\xff\xb7\x22\x73\x7c\x46\xe0\x2c\x8a\x9c\x97\xbd\xe3\x31\xfc\x11\x2e\x9d\x30\x57\x40\x28\x09\x4d\x67\x5f\x46\xf3\x42\x28\x45\xa5\x2a\xcc\x98\x2b\xf8\x59\xcd\xfd\xf4\xaa\x92\xbf\x54\xa2\x1c\x22\xc2\x7f\x71\x3c\xb8\x2c\xa2\x52\x98\x95\xbd\x79\x05\x6b\x00\xac\xa8\xe1\x33\x97\x38\xe1\x44\xc6\x60\x85\xb6\x2a\x83\x79\x4d\x64\x67\x35\xce\x40\x4a\xd4\x32\xb9\xbb\x6d\xbb\xe6\x99\x33\x8b\x31\xd2\xc1\x1a\xdb\x89\xb3\x4a\xde\x74\x63\x73\x25\x55\xab\x28\x15\x1a\xd6\xdb\x6a\xa7\x4f\xba\xdd\x8a\xa5\xa0\xf6\xe0\x62\x1b\x54\xba\xe6\xd4\x54\x4e\xd7\x7a\x24\x36\x72\x3c\xc2\x6b\xaf\xb9\x47\xfb\x04\x7a\x34\xfb\xf1\xf0\x5a\x4c\xd7\xb3\x1e\x33\x57\x8c\x69\x02\x3e\x94\x7b\x92\x2b\xea\xc2\x57\xd7\x93\xaf\x17\x23\x8e\xb1\xbd\x72\x83\xad\xe1\x16\x27\x1d\xb6\x19\xe4\x71\xbb\xc7\x12\x9c\x3d\xd3\xb5\x21\x75\x2b\x27\xd4\x02\xfa\x4b\xf1\xbe\x62\x32\xc4\x1a\x69\x7e\x35\xf3\xb7\xfb\x89\x6b\xa4\x55\xe0\x89\xa4\x7c\xc5\x6c\xf3\x7e\x48\x5a\x92\xc8\x7a\xa2\x59\x40\xdd\xdd\xde\x44\x89\x41\xb3\x64\x29\xdf\xed\xbb\x58\x39\x1c\x29\xe3\x3a\x62\xf8\xc5\xfb\xb8\x2b\x23\x7d\x84\xe1\xad\xef\x10\xc4\x8e\xc6\x79\xdc\x28\xec\xcd\x9a\xdc\x95\x0b\x1b\xef\x78\xfb\x7b\x4f\xe4\xc8\x2e\x86\x21\xf2\xf0\xdb\xde\x55\x21\xc7\xd4\xc2\x0f\x6e\x40\xfd\xfa\xd0\x9d\xb5\xfd\xdb\xb7\xa9\x7e\x02\xf9\xa7\x91\xa9\xdc\xe9\x26\xbf\x73\xd2\x26\x94\xd3\xcd\xd7\xaa\xed\xd9\x93\xe9\xd1\x36\xe3\x2e\x33\x7a\xfd\xfc\xdb\xe9\xe0\x39\xf8\x4c\xf8\xf7\x02\x0a\xdd\xa5\x40\x00\xa4\x97\x05\x85\x1e\xe7\x40\x4f\x72\x6f\xe7\x26\xcf\xb1\x38\x09\x9c\x51\xfe\x8b\x0a\x8e\x13\x23\x7c\xbc\xe9\xb2\xc3\xc5\xd9\x90\x80\x62\xee\xd4\x21\x65\x22\x1a\xde\x33\x36\xfc\x94\xd0\x8f\x0d\xff\x3a\xf5\x7e\xb1\x11\xca\x39\x39\x36\xda\x59\xc8\xc7\x46\x7f\x9a\xfa\x76\x68\x78\x06\xef\x10\x1a\x3d\xc9\xff\xdf\xd0\x08\x1e\xfc\x7e\x64\x68\xb8\x11\x28\x18\x2f\xc2\x97\xde\x36\x32\xda\xd7\xaa\xef\x1c\x31\x3a\x31\x93\xf3\x45\x14\x0a\x5d\xc0\x73\x55\xe5\x76\x88\x98\x1c\x06\xdb\x67\xea\xde\xfc\xd7\x29\x89\xf5\x9b\xa1\xea\xce\x2e\xeb\x05\x60\x21\xbf\x9a\xb2\xb8\x67\xf4\x14\x60\xfa\xda\xd9\xcb\x9c\x24\x3b\x9d\xae\x27\x35\x55\xa7\xc6\x47\x86\xfd\x24\x3a\xa2\x86\x7f\x93\xef\x69\x71\x84\x0c\x82\x37\xfb\xcf\x7c\xf3\x50\x35\x9a\x3d\xe7\xdc\x3d\xdf\x8f\xe1\x25\xe6\xc7\x92\x31\xc7\x05\x89\xeb\x46\x5e\x2a\xc6\x83\x11\xfc\x98\xc9\x8f\xff\xdc\x72\x64\xae\x1f\x3c\x09\x1e\x15\xf3\x14\x8c\x21\x2e\x59\xba\x97\x42\xfb\xab\x53\x90\x2a\x13\xaa\x3b\xa3\x4d\x68\x3f\x7d\x95\x88\xfb\x09\x73\x3a\xb2\x1f\x08\x66\xe2\x15\x37\xcc\x5c\x33\x46\x07\xbf\xd5\x91\x1f\xda\x5b\x81\xae\x70\xde\xa1\x7d\x4a\x5e\xfa\xb1\xa8\x42\x99\xf0\xe9\xcb\x97\x7b\x3a\x4a\xcf\x95\xcf\x9c\x1e\xf5\x33\xc8\x84\xe4\xa9\xce\xb7\x74\xb7\x37\xae\xfc\x1b\xdd\x4d\xca\x9b\x32\x33\x02\xa2\xf9\xd5\x1f\x2e\x2f\x2f\xf1\x9a\xc2\x6a\x61\x47\xf7\x08\xef\x2b\x67\x5e\x2e\x30\x0d\x7a\x65\x65\xd7\xdd\xb0\x0e\x9b\x3a\xa6\xcc\xb8\x3c\x98\x17\xe3\x54\xfb\xd6\x4f\x6e\xde\x11\x34\x01\xba\x93\x11\x3d\x7f\xfc\x2f\x00\x00\xff\xff\x11\x72\x87\xb2\x3f\x1f\x00\x00")

func templatesServerBuilderGotmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesServerBuilderGotmpl,
		"templates/server/builder.gotmpl",
	)
}

func templatesServerBuilderGotmpl() (*asset, error) {
	bytes, err := templatesServerBuilderGotmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/server/builder.gotmpl", size: 7999, mode: os.FileMode(420), modTime: time.Unix(1435797486, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _templatesServerMainGotmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x56\x4d\x73\xdb\x36\x13\x3e\x5b\xbf\x62\xc3\x49\x66\x28\xbf\x12\xf5\x9e\xdd\xd1\xc1\x4d\x9b\x5a\x6d\xfd\x31\x56\x32\x99\xe9\xf4\x02\x93\x4b\x0a\x31\x05\xb0\x00\x68\x45\xd5\xf0\xbf\x77\x17\x24\x28\x4a\x76\x13\x27\x69\x6f\xbd\x88\xe2\x62\xf7\xd9\xdd\x67\x3f\xc0\x4a\xa4\xf7\xa2\x40\x58\x0b\xa9\x46\x23\xb9\xae\xb4\x71\x10\x8f\x00\xa2\x52\x17\x11\x3f\xb5\xf5\x0f\x85\x6e\xb6\x72\xae\x8a\x46\xfc\x56\x48\xb7\xaa\xef\x92\x54\xaf\x67\x85\x9e\xda\x8d\x28\x0a\x34\xc3\xbf\xb6\xc2\xd4\xab\xee\x76\x46\x28\x72\x90\xfc\x80\xb9\xa8\x4b\xb7\xf0\x2e\x6c\xd3\xec\x76\x95\x91\xca\xe5\x10\xbd\xfa\x23\x82\xa4\x69\xbc\x32\xaa\xac\xfb\xd7\x9a\xbd\xbc\xc7\xed\x04\x5e\x3e\x88\xb2\x46\x38\x9b\x43\x32\xb0\xe7\xb3\xa6\x21\x55\x18\x22\xb5\xba\x07\x70\xe3\xd1\x68\x36\x83\xb7\x2b\x69\x21\x97\x25\xc2\x46\x58\x28\x50\xa1\x11\x0e\x33\xb8\xdb\x82\x5b\x21\x74\x91\x83\xd3\xba\x4c\x58\xff\x52\xdc\x93\xb4\x36\x08\x4a\x3b\x12\x83\x7e\x40\xb3\x31\xd2\x21\xe9\x07\x28\x91\x3b\xb2\xd9\xea\x7a\x00\x28\x1d\xdc\x61\x2a\x6a\x4b\xc7\x65\xc9\x87\x06\x30\x93\xce\xc2\x46\xd7\x25\x39\x44\x28\xb5\x75\x2f\xd8\xc9\xc2\x75\x42\xad\xca\x2d\x9f\x04\x27\x0e\x15\xc8\xdc\x23\xe3\xc7\xaa\x94\xa9\x74\xa4\xc0\xb4\xca\x7c\x0b\xd3\xa9\x54\x69\x59\x67\x38\xe5\xc2\x41\xae\x8d\xcf\x21\xc4\xe0\xfd\x92\xcc\xd6\x95\x2f\x28\xd5\x69\x2d\x54\x66\xc9\x63\xa1\xcf\x7a\xad\x90\xf2\x5e\x80\x86\xdc\xc3\xd4\x41\x92\xcc\x92\x04\xa6\xe7\xc4\x61\x72\x5e\x55\x57\x62\x8d\x4c\x39\x45\x94\xdc\x10\xd9\xa9\xac\x44\x49\xdc\x4f\xa7\x55\x78\x63\xcd\xc1\x51\xe0\x7e\xf4\x20\x4c\x70\xf4\xf3\xf2\xfa\x0a\xe6\xf0\xc1\x6a\x95\xdc\x8a\xcd\x25\x5a\x4b\xbd\x17\x93\xe1\x72\xaf\xd0\x34\x54\xae\xbc\x56\xa9\xef\xc9\x78\x0c\x3b\xaa\x64\x07\xb0\xa4\xf4\x27\x80\xc6\x70\x2b\x30\x17\xc9\x15\x6e\xe2\x01\xfa\x04\xa2\x68\x4c\xfa\x14\x27\x6b\xbd\x98\x83\x92\xa5\x47\x00\xe2\xbc\x48\xde\x08\x47\xd4\xa8\x98\x0e\x59\xad\xe1\x0e\xf5\x0c\x11\x9e\xb6\xc9\x4f\x48\xb4\x3f\xc4\xd1\xcd\xf5\xed\xdb\x80\xe3\x8f\xe7\x73\x02\xee\x70\x5a\x01\x44\xff\x8f\x02\xc2\x8a\xaa\x79\x84\x70\x71\xbd\xec\x11\xfc\xf1\x10\xa1\x15\xf0\x88\xa5\xa2\xe4\x97\x1e\x49\x54\x92\x81\x98\xca\x76\x30\x9b\x86\x53\x1c\x16\xe1\xfc\x66\x11\x0f\xe8\x60\x1f\xa9\x56\xb9\x2c\xa8\x53\xf9\x8c\x20\xc6\x0c\x55\x4a\xeb\xb8\xae\x3d\x5f\x34\xc1\xc9\xaf\x5e\x18\x47\x2e\xad\xa2\x49\x1b\xc7\xff\x20\x3a\x8b\xe8\x97\xd3\x1a\x8f\x4e\x8e\x99\x3b\x39\x79\xc4\xdb\x09\x85\x7a\x92\xaf\x9d\x2f\xb7\xcb\xe3\x88\xdb\x46\xaa\x82\xc3\xbe\xa8\xa9\xd1\xfa\x58\x41\x38\xe0\xa5\x71\x36\x9b\xbd\xb2\xbf\x2b\x72\x19\xa2\x4a\xce\xb3\xcc\xc4\x63\x1f\x69\xe7\x92\x42\x64\xdd\x64\xc9\x4d\x18\xef\xc3\xa7\x84\x5a\x59\xdb\x7e\x8b\xb6\xf1\xdf\x2d\x9a\xe6\x3d\x6d\xa1\x77\x8b\xae\xd3\x08\xec\xbb\x67\x16\x9d\xe2\x6f\x3b\xec\x98\x38\x38\x3d\x60\xfe\x88\xf6\xb6\x15\x69\x68\x7b\x33\x3f\x74\x6c\xb7\x42\x83\x6d\xf9\xda\x50\x7f\x34\x86\xc6\x6f\xce\xf1\x68\x63\x07\xb2\x83\x95\xf8\x5a\x2b\x5b\xaf\xd1\x86\xc9\xa2\xdd\x56\xe2\x1a\x95\x13\x4e\x6a\xd5\x34\x0c\x47\x31\xbc\x2e\x85\xb5\x6d\x14\x9d\x05\x43\xd3\xc1\xb1\x7e\x3c\x6e\x97\x5e\x69\xf1\x33\xc6\x4c\xf4\xbd\x74\x21\x02\xf3\x86\xd8\x88\x99\x92\xd8\x80\xd4\xc9\x2d\x8a\x8c\xa9\x77\xc2\x14\xe8\x80\xaa\x8c\x26\x17\x29\xee\x9a\x71\x9b\x52\xc7\xae\x41\x57\x1b\x15\xb2\xbc\xd2\xae\x8f\x08\xb3\x38\x22\xef\x5d\x1b\xa4\xc1\xf3\x8a\x36\x2f\x6f\xd3\x2d\xf2\x8e\xe4\x15\xb7\x37\xf0\xe3\xd2\x8c\x87\x77\xc0\xf1\x6d\x40\x1d\xa7\xb3\x3a\xfd\x12\xc6\x3a\x8b\xaf\x63\x6c\x60\x1c\x18\x0b\xa2\x3d\x63\x1b\x66\xec\x3d\xdf\x0b\xc4\x58\x46\xbd\xf6\xed\x7c\x55\xc1\xef\xb7\xf2\xb5\xc4\xb4\xa6\xc8\xb6\x74\xf9\x4a\x25\x39\x67\xdb\x29\x78\xf6\xec\xf7\xc2\xca\xf4\xbc\x76\x2b\x2f\x7d\x4c\x00\x1f\x51\xf2\x3e\x4f\xba\xcf\x68\x93\x3b\x9a\xf9\x62\x02\x15\xa9\x74\x2f\x63\x88\x4f\x0f\x77\xff\xa4\xcd\x70\x7c\x98\x35\x0d\xe5\xe4\xef\x52\xbf\xe3\x38\x40\xb0\xb7\xcf\xa7\xbc\x4f\x35\xa4\x41\xc3\xf9\x0b\x6e\x9f\x99\x87\xd3\xf7\x84\xfa\xcf\xc5\xce\xf3\x4f\x5f\x22\x6d\xf4\xfb\x1a\xe6\x46\xaf\xf9\x75\x49\xd7\x7f\xca\x82\x2f\x49\xec\xe9\x6a\x5e\x57\x7c\x4d\xb7\x45\xec\xee\xe2\xb0\xab\x1e\xa7\x7c\x41\x17\x7e\x19\xfa\xfe\x60\xa7\x3d\x56\xda\x37\x73\x80\x35\x62\x4d\x4e\x2a\xff\xfc\x14\x40\xab\x39\x2c\x07\x21\x42\xc2\x7c\x6b\x23\xff\xa4\xaf\xa1\x1e\x6c\x72\x58\xb5\xbd\x0a\xf9\xe9\x3f\x22\x4e\x9f\xfc\x8a\xa0\x32\xb5\x46\xcb\x3a\xa5\xe9\xb7\x97\x3a\xc3\x32\x00\xdd\xfa\x12\xd9\xd7\x9a\xe9\xfc\x78\x7d\xf7\x01\x53\xd7\x34\xa7\xbd\xb3\x23\xa3\x3e\x8c\xa7\xea\xfc\xb4\x97\x20\xf8\x0d\x8d\x3e\x06\x78\xdc\x0f\x3a\x54\x69\xd0\x0c\xcf\x98\xe3\x4f\x2e\xa2\x83\x62\x7e\x4d\xfd\xfe\x2b\xd9\xbf\x56\xb2\xa3\x79\xa5\x8f\x8b\xbf\x02\x00\x00\xff\xff\x43\xcc\x61\x58\x55\x0d\x00\x00")

func templatesServerMainGotmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesServerMainGotmpl,
		"templates/server/main.gotmpl",
	)
}

func templatesServerMainGotmpl() (*asset, error) {
	bytes, err := templatesServerMainGotmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/server/main.gotmpl", size: 3413, mode: os.FileMode(420), modTime: time.Unix(1436298580, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _templatesServerOperationGotmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xdc\x56\xdf\x6f\xe4\x34\x10\x7e\xcf\x5f\x31\xac\x8e\x63\x53\x85\xf4\xbd\xa8\x0f\xd0\x03\xf5\x1e\x38\x4e\x6d\x05\x8f\xc8\x4d\x26\x1b\xd3\xc4\x4e\x1d\x7b\xb7\x4b\x94\xff\x9d\xf1\x8f\x64\x93\x6d\xba\x42\x40\x2b\x74\x52\xa5\x6e\xe2\xf1\x37\xfe\xbe\x99\xf9\x9c\x86\x65\x0f\x6c\x83\xd0\x75\xe9\x67\xff\xb3\xef\xa3\xe8\xfc\x1c\xee\x4a\xde\x42\xc1\x2b\x84\x1d\x6b\x61\x83\x02\x15\xd3\x98\xc3\xfd\x1e\x74\x89\xd0\xee\xd8\x66\x83\x0a\xb4\x94\x55\x6a\xe3\x7f\xcc\xb9\xe6\x62\x43\x8b\xc3\xbe\x9a\x6f\x4a\x0d\x8d\x92\x5b\x84\xc2\x68\x07\x55\xa2\x80\xbd\x34\xa0\xf0\x5b\x65\x84\x43\x1a\xa0\x21\x93\x75\xcd\x44\x1e\x45\xbc\x6e\xa4\xd2\xb0\x8e\x00\x56\x02\xf5\x79\xa9\x75\xb3\x8a\xe8\xa9\xeb\x14\x13\x74\xd8\xf4\x03\x16\xcc\x54\xfa\xa3\x0b\x6c\xfb\xbe\xeb\x1a\xc5\x85\x2e\x60\xf5\xf5\xe3\x0a\x52\xa2\x60\x83\x51\xe4\xe1\x97\xdf\xf6\xee\x01\xf7\x09\xbc\xdb\xb2\xca\x20\x5c\x5c\x42\x3a\xd9\x6f\xd7\xfa\x9e\x42\x61\x8a\xe4\x63\x67\x70\xb1\x53\x87\xe4\xba\xaa\x58\xdb\x7e\x62\x35\x2d\x5f\xd3\xb1\x2b\x54\x3f\x19\x91\x81\x36\x4a\xb4\xc0\x88\xb1\xc8\x34\x97\x02\x76\x5c\x97\x8e\xa8\x72\x7a\xb4\x7c\x23\x18\x05\x21\x50\x1a\x49\x81\x04\x75\x6d\x88\xf8\x04\x0f\x4a\x0f\x18\xe9\x7d\x83\x27\x72\xd9\x1c\xeb\xae\xe3\x05\x50\xf1\x14\xab\x1d\x93\x69\xb0\x7f\x1b\x8e\xee\x02\x69\x37\xa4\xdf\x1b\x5d\x4a\xc5\xff\xa4\x72\x8e\x1b\x13\x98\x86\x4d\x42\xfa\xfe\xcc\x36\x07\xa9\x92\xf1\x86\x55\x36\xc0\xc5\xc5\x10\x52\xdf\x9a\x2c\xc3\xb6\xfd\x59\xe6\x58\x0d\xdb\x6f\xd0\xe9\x70\x25\xeb\xa6\xc2\xa7\x5f\xee\xff\xc0\x4c\x3b\xa0\x90\xe2\x68\xd3\x98\x1c\x95\x92\x8a\x24\xb6\xcc\x60\x5d\x88\x97\xc9\xc7\xe0\x1f\x8e\xf8\x37\xee\x3f\xbc\x86\x0c\xcd\x20\x01\xbc\xb9\x20\xd0\x51\x07\x2a\x07\x01\x85\x58\xa4\xfc\x9f\xb1\x1b\xd8\x44\xfd\xcb\x9d\x6e\x7b\x17\x55\xc1\x32\x9a\x6c\x49\x26\x50\x32\x0d\x19\x13\xa1\x6f\x81\xa6\x86\xe7\x8b\x8d\xed\xcf\x7a\xa2\xaf\x27\xc8\x96\xf3\x62\x8d\xbf\x8c\x1e\xf7\xf2\x7e\xc2\xdd\x9c\x0e\x64\x0a\xc9\x0d\xad\x85\x08\xdc\x81\xf5\xbe\x74\xd0\xc6\x6b\x8d\x8b\xca\xca\xc6\xba\x28\x19\x8e\x9f\x9d\x67\xb8\xeb\x4c\x3f\xc1\x59\xcd\x73\x42\xda\x31\x85\xe9\x95\x24\xa5\x9f\x74\x32\x98\xcd\x72\x3d\x62\xd7\xec\xd3\x44\x93\x56\x7c\x3f\x5f\xea\x02\xe4\x05\x50\xae\x24\xd4\x4e\x5d\x0c\x09\x7a\x4b\xd9\x4b\xf7\x41\x66\xb7\x9a\xd4\xde\x38\x9d\x66\x4f\xde\x65\x17\x1a\x04\x5a\xad\x4c\xa6\x5d\xfe\x90\x68\x89\x8f\xb3\xea\x69\xb7\xf8\xff\xb0\x68\x09\x07\x5f\xbf\x3e\x25\x82\x3d\xb8\x77\x24\x5a\xbe\xc1\x0c\xf9\x16\x55\x38\xd5\x91\x3c\x31\xdc\xa2\xda\xe2\xf5\xdd\xdd\xe7\xb5\x0a\xe5\xbb\xc1\xb6\x91\xa2\xc5\xdf\x14\xa7\xde\x4e\x40\xc1\x59\x78\xff\x68\xb0\xd5\x61\xba\xa5\xd1\x98\xc0\xef\xf6\x62\x7a\x96\x65\x20\x97\xde\xd8\xa8\x8f\xa2\x90\x6b\xeb\x92\x03\xd5\x69\x23\x1b\x37\xc8\x09\x50\x97\x9d\x86\x1a\x37\xad\xed\x91\x2c\x6e\x4c\x80\x04\x67\x77\x7e\x75\x09\x82\x57\xee\x60\x70\xea\x38\x8e\x59\x4e\x4c\x09\x22\xa0\xd0\x18\xc9\xdc\x50\xe3\x27\x03\x27\x02\x8c\x1d\x90\x6f\x1b\xfa\x69\x2f\xd4\x2d\x53\x70\xb0\x54\x47\x44\x48\xba\xf5\xf1\x11\x0e\x93\x08\xab\xd1\x0f\xba\x7e\x15\xcf\xc6\x6b\x32\xae\xfe\xe0\x9e\xfa\xfc\xec\x87\x0c\x97\x3e\xc7\x09\xf8\x41\x3c\x4a\x51\xb5\x38\x3c\xa5\xeb\x23\x6f\x88\x81\xe6\x96\xeb\x6f\x5a\x90\x0f\xfe\x73\x87\xfe\x68\x68\xab\x6a\xef\xaf\xf3\xe7\x3e\xe2\x28\xcf\xbe\x49\x82\xce\x27\x2b\xf4\x03\x17\xf9\xaf\xd6\x4a\x43\xa3\x8c\x85\x4a\x8e\x5a\xfc\xfd\x73\x8c\xd1\x15\x1d\x13\xd2\x63\x70\xb4\xef\x66\xf5\xb5\x54\xee\x29\xcd\x60\xcc\xaf\x55\xee\xc5\x56\x1d\x5f\xce\x8d\x52\x59\xac\x83\x57\x2e\xc6\x5c\x84\xf5\x25\xf9\xc2\xd0\xa6\x2f\xdd\x1d\x8b\x4a\x8d\x19\xc7\x86\x71\x65\x66\x99\x36\xae\xb0\xe1\x66\x73\xdf\x71\xbe\x1a\x6f\x38\x2c\xd1\xbf\x87\x7d\x41\xe8\x85\x06\xf1\x6d\xea\xde\xfe\xcd\x02\x2d\x36\xf1\x3f\xaa\xc2\x78\xed\xfe\x6f\xa4\x7f\x7b\xe5\x9d\x41\xf4\xd1\x5f\x01\x00\x00\xff\xff\xf6\x60\xe1\x76\x99\x0d\x00\x00")

func templatesServerOperationGotmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesServerOperationGotmpl,
		"templates/server/operation.gotmpl",
	)
}

func templatesServerOperationGotmpl() (*asset, error) {
	bytes, err := templatesServerOperationGotmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/server/operation.gotmpl", size: 3481, mode: os.FileMode(420), modTime: time.Unix(1435797486, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _templatesServerParameterGotmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xd4\x59\x4b\x73\xdb\x38\x12\xbe\xeb\x57\xf4\xaa\x32\x29\xc9\xa3\xa1\x7d\x98\xda\x43\x66\xbd\x87\x38\x9e\x8d\xab\xb2\xd9\xac\x3d\xe3\x4b\x92\xda\xc0\x24\x24\xa1\xcc\x97\x01\xd0\xb6\xa2\xd2\x7f\xdf\x6e\x00\x24\x41\x8a\xa4\x25\xc5\xae\xad\xcd\x21\x26\x5e\x8d\x0f\xdd\x5f\x3f\x00\xad\xd7\x11\x9f\x8b\x94\xc3\x38\x97\x22\x11\x5a\xdc\xf3\x7b\x16\x8b\x88\xe9\x4c\x8e\x37\x9b\xd1\x7a\x2d\xe6\x10\xfc\x53\xa4\x1f\x78\xba\xd0\x4b\xec\xc1\x36\x97\x12\xde\x9c\x82\x9b\xc8\xeb\xe1\xc9\x7a\x1d\x7c\x62\x34\x6d\x06\x63\xfc\xfe\x90\x85\x4c\x8b\x2c\xdd\x6c\xc6\x33\xc0\xf6\x35\x8b\x0b\x7e\xfe\x98\x4b\xae\x94\xe9\x36\xbd\x9e\xf4\xe9\x6f\x46\xf8\x5f\x4e\x21\x15\x31\xac\x47\x00\x92\xeb\x42\xa6\xd4\x3b\x22\x34\x3c\x8d\x6a\x54\xec\x71\x10\x55\x39\x7c\x20\xaa\x5a\xfa\x5e\xa8\x70\x27\xcd\x65\xda\x8d\xc9\x0d\x1e\x80\xe8\x9b\x5d\x62\x45\x7f\xdb\x4f\x4f\x22\x15\x49\x91\xf4\xda\x8e\x06\x07\x11\xcd\xe3\x8c\xe9\xbf\xfe\x3a\xe9\x42\x36\x2d\x4d\x68\xb7\x30\xad\xf3\xc7\x30\x2e\x14\x52\xa9\xea\xde\xd7\xae\x03\x78\xed\xe0\x8f\xe2\x2d\xb7\x68\xe1\x2d\xbb\xf7\xc3\x5b\xc4\x5a\xe4\x31\xff\xd7\xbc\x07\x72\x35\xfe\xa3\xa8\xbd\x8d\xf6\x42\x78\x9e\xf6\xa9\x93\x46\x0e\xf3\x0f\x2b\x73\x67\x18\xe6\x6f\x1d\x6d\xc2\x42\xe9\x2c\x99\x67\x32\x61\xba\x11\x70\x3a\x30\xfe\x6e\x66\x3d\xa1\x3d\xea\xb0\x13\x4d\x53\x69\x29\xd2\x45\x9f\x2e\xed\xbe\x6a\x37\xf0\x35\x68\x15\x8b\xb0\x2b\x3c\x7e\xe4\x3c\x52\x57\xe2\x3b\x37\x3d\x88\x51\xb2\xe4\x23\x4b\xb0\x49\x9d\x74\x16\x91\x92\x65\x63\x9e\x76\x23\x9a\x6e\x7b\xec\x85\xe6\x89\xea\x75\x59\x33\xfa\x94\xdd\x5a\x38\x4a\x47\x75\x92\xf7\x75\xc9\x21\x40\x6e\xf4\x20\x40\x95\xe4\xbd\x00\xfd\x99\x8a\xbb\x82\x0f\x60\xf2\x26\xec\xcd\xef\xff\x73\xdf\xca\x65\x96\x73\xa9\x57\x1d\x4c\xbd\x50\x9f\xca\x34\x4f\x2b\x50\x3b\x79\x8c\x50\x3b\xb3\x3f\x04\x34\xc5\x3f\xea\x85\x3a\x33\x6e\x6b\xfd\x0c\x53\x51\x53\x46\xb7\x4f\x77\x8a\xc9\x52\xcd\x10\x6b\x4b\x40\xcb\xbf\x9a\x2b\x5b\x87\xbc\x11\x69\x54\x81\x1e\x77\xcd\x30\xd2\x68\x1a\xf7\x14\x80\x14\xe4\xa9\xa6\x69\xc1\x05\x8e\x3c\x5e\x33\xc4\x10\x92\xd9\xd4\x03\x5b\x04\x57\x79\x2c\xf4\xdb\x95\x3d\xa0\xb5\x1d\xcd\xf7\xe7\x7e\xee\xea\xfd\x6a\xad\x7b\x96\xc5\x31\x0f\xc9\xbe\x55\x28\x32\xae\x1d\x2b\xde\xb5\xa5\x64\x0f\xf5\xf9\xbc\x41\xf5\xdd\x00\x42\x17\x19\xdd\x33\x09\x8d\x31\xd3\xfc\x63\x95\xf3\xf6\xa2\x6b\x47\xbb\xf3\x98\x27\x08\x8e\x24\xcc\x8b\x34\x9c\x34\x26\x51\x20\x32\x0c\x3b\x5b\x8a\x38\xda\x66\x5f\x3d\x64\xb7\x98\xc2\x11\x92\x2d\x93\x2a\x70\xe2\x71\x96\x61\x62\x93\x3b\x6d\xbe\x81\x15\x82\x10\x2b\xce\x22\x85\x91\xb3\x23\x64\x47\xf3\x3c\x84\xf3\xe4\xb7\x56\xdf\xdf\xa0\xa5\x8f\xd6\x84\x9f\x7f\x76\x20\xd0\xa4\x28\xd0\x41\xde\xa2\x67\x3d\xd0\x60\x3d\xf1\xc0\x0e\x20\x0f\xef\x11\x39\xf1\xf0\x9e\x54\x31\x2b\x7d\xb8\x52\x83\x37\xa3\xa9\x49\xc3\x03\x8f\x00\x53\xc4\xe3\x62\x80\xe7\xb0\xbe\xcb\x92\x16\x2f\x52\xa3\x24\x52\xee\xa4\xda\x63\x30\xa7\xf9\xd6\xb0\x21\x63\x18\x03\xea\xb8\x02\x62\x0f\xd2\x4b\x91\xe6\x81\x86\x68\xb1\x1d\x89\x1a\xb1\x88\x76\x85\x36\x4d\x4f\x81\xe5\x39\x92\xbb\xb9\x8b\x9c\x81\xd1\xf4\xd4\x2c\xb0\x8e\x61\xc4\x1d\x0c\x79\x40\x1d\x1d\xa8\x5b\xb8\xf7\x43\x3e\xbc\x5b\x15\x80\xe8\x54\x50\x93\xac\x11\xee\x5a\xae\xe3\xc7\x28\xdf\x69\x7e\xd8\x84\x1e\xee\x97\x50\xc3\xf6\x26\x65\x20\xab\x02\x71\xce\xc2\x5b\xb6\xe0\x36\xef\x9b\x4f\x1c\x1d\x1d\x1f\xc3\x1f\x4b\xa1\x60\x2e\x62\x0e\x0f\x4c\xc1\x82\xa3\x5e\xf0\x40\x11\xdc\xac\x40\x2f\xb9\x89\xc3\x0b\xf4\x5d\x9d\x65\x71\x40\xf3\xcf\x23\xf4\xdc\x74\x81\x83\xe5\xba\x44\x2c\x96\x1a\x30\xec\xdc\x73\x8c\x71\xda\x88\x5a\xf2\x14\x56\x59\x81\xe7\xfa\x45\x16\x69\x43\x52\xb9\x05\x84\x59\x92\xb0\x34\x1a\x8d\x44\x92\x67\x52\xc3\x04\x0f\x3d\x4e\xb9\x3e\x5e\x6a\x9d\x8f\xa9\xb1\x10\x7a\x59\xdc\x04\x38\xf1\x78\x91\xfd\xe2\x04\xf8\x9f\x34\xf3\x56\xe8\xdd\x26\xd3\xdf\xdd\x66\xda\xf0\xb0\x17\x84\xe3\xb2\xce\xd8\x13\xb8\x31\xb4\x64\x29\x9a\x26\x78\xc7\xe7\x0c\xef\x14\x17\x46\x1d\x8a\xb8\x8b\x29\x35\xd5\x73\x18\xff\x74\x67\xd2\xaf\xf5\xd2\x34\x72\x5f\x76\xd9\xab\x5b\xbe\x9a\xc1\x2b\xe3\xc7\x44\xd0\xc0\x5b\x4f\x63\x26\x8d\x80\x2f\xc9\xce\x6d\x88\x9b\x1a\x2e\x10\x95\x62\xa6\x94\x2d\x0a\x4d\x7d\xa8\xd0\x4c\xc6\x5d\x14\xb0\x38\x36\x86\xbc\xc9\x8a\x34\x82\xdc\x8e\x52\x06\xa1\x4e\x5c\xfa\xbe\x40\x73\x7a\xeb\x81\xf2\x90\x09\x9f\x24\x5b\xaf\x72\x11\xa2\x08\x43\x2b\xf4\x48\xcc\xd9\x90\xdd\x18\x47\x8c\x60\x2e\xb3\x04\x18\x90\x56\x82\x4b\x8e\x95\xa2\xd2\x23\x5c\xc0\xbb\x11\xe1\x6d\xa2\x08\xb5\xcb\x39\x4e\x77\x76\xa8\xcc\x27\xef\xb8\x0a\xa5\xc8\x6d\xe8\xb6\x07\x6b\x74\xf9\x5a\x0c\x3e\xb9\x84\xe9\x50\xd7\x09\xbd\x56\x8f\x75\x95\xb7\x18\x19\x1c\x3a\x54\x82\x5e\x02\x85\x0a\xd4\x0b\x6a\xa3\xb4\x3e\xb6\x90\xf3\x66\xca\x0c\x84\x06\x84\x5e\x24\xd8\xab\x97\x4c\x13\xe1\xf1\xba\xf8\x48\xae\x93\x2e\x14\x08\x6a\x99\xe2\x80\x81\x0b\x24\xec\x26\xe6\x13\x3c\xde\x3c\xd1\xa8\x87\x85\xc0\xcf\xd5\xd4\x66\x2b\xaa\x15\xb8\x9c\xb3\x90\x13\x14\x52\xbb\x32\x02\x6c\x00\x57\xb4\xd9\x83\x40\x0b\x15\xa8\x5b\x5c\xc6\x8c\x53\x26\x5c\x2f\xb3\x08\x48\xef\x6a\x44\xf5\x07\x50\xf8\xb8\xe4\x21\xc7\xe4\x2b\xdd\x81\x8f\xba\x94\x3c\xf5\x4f\x3b\x91\x70\xe4\xdb\x66\x06\x32\x2b\xd0\x83\x8f\x12\x11\x45\x31\x7f\x40\x5b\xe2\xcd\x41\x87\x4b\x1e\x5d\xd2\x40\x09\x99\x2c\x44\x25\x13\x66\x2e\xf8\xfc\xd5\xf4\x95\x75\x42\xf0\x9e\xa9\x7f\x17\x5c\xae\x4a\xc3\xdd\x29\x53\x83\x05\x7f\x5e\x7e\x08\xcc\xc0\xa4\x4e\x4a\xe0\x16\x50\x29\x51\xce\xf7\xac\xd3\xc5\x83\x72\x9f\x34\xd3\x5b\x25\xae\xad\x7a\xeb\xdd\x37\x9b\x46\x78\x6f\xaa\x27\x20\x23\x6f\xb1\x64\x72\xa7\x82\x7f\x70\x5d\xdf\x27\xa6\x4e\x27\xee\xd6\xdb\x71\x99\x05\xa3\x86\x2a\x8c\x63\xc3\xd4\x37\xd3\x2a\x5f\x57\x27\xc5\x02\x09\x65\x1e\x0c\xcd\xe2\xb0\x8a\x78\x49\x90\xef\x39\xc3\x44\x79\x38\xcc\xc0\x0a\x78\x49\x88\x15\x61\x6a\xb3\xff\x8e\xf9\xa9\xea\xf2\xef\xc0\xed\x3b\xb1\x45\x57\xd5\xa0\xd2\x20\xa2\xd5\x1e\xd8\xde\x2a\xb3\x03\x20\x15\x9c\x1f\xf9\xc3\xe4\xd7\x93\x13\xac\x25\x25\x4a\xa7\x34\x6a\x32\xe8\x97\x71\x73\xeb\x2f\x63\x98\x33\x1c\x88\xde\xc0\x4f\xf7\x63\x7b\x3c\x73\x3e\x30\x67\xb3\x9b\x6c\xeb\x79\x3b\x96\x9d\x82\x4b\x34\x01\x01\x5f\xbf\xc3\x08\xf3\x06\xda\xc7\xb6\x07\x6d\xf7\xdb\xde\x4d\x43\xab\x87\x99\x99\xf4\x66\x6a\xd8\xe7\xb5\xb2\xab\xee\xaa\x38\xee\x99\xfd\xd9\xbd\xbd\xe3\x36\xda\x11\x00\xfa\xee\x9c\xcf\x47\x69\x4a\x35\x4d\x5a\x3f\xcb\x59\xfa\x6c\xf4\x82\x07\xf2\xad\x57\xe5\x84\x0b\xf5\x36\x8b\x4a\x2b\x35\x2e\x4e\x76\x3f\xb4\x2b\xa5\x53\x59\x7e\x20\x70\x5a\x30\x83\xd7\x3b\x38\xc3\xce\x28\x9d\xa3\x22\x0c\xc5\xcf\xa9\x39\x69\xb9\xe7\xb8\xeb\x72\xd8\xef\xa6\x5b\x9c\xa4\xe4\xfd\x9f\xd6\xfd\x65\x3b\x2d\x9b\x17\x89\xd4\x96\xec\x4f\x3a\x7a\x79\x8f\xe8\xb9\xa7\x6c\x8b\x28\x6f\x2e\x93\x27\x2d\x39\x68\x4d\xfb\xef\x06\xa3\xd9\xad\x6b\x6d\x46\xf5\xff\xbb\x44\x8d\xad\xb3\xec\x05\x6d\x00\x58\x05\xc1\x45\x08\xc7\x33\x2f\x60\xb8\x2f\x44\x47\x4f\xc1\xb8\x7a\x0a\x7f\x87\x93\xce\x97\x82\x33\x2c\xdd\x32\x25\x34\xaf\x1f\x5e\x2c\x35\x70\x55\x10\x04\x25\xb1\x9b\xaf\x2b\x58\x76\xbf\x0a\xcb\xc2\xca\x94\xe6\x55\x99\x05\xe6\xbd\xa8\x5d\xb3\xf8\x15\x8b\xef\x09\xd5\xcb\x8a\xf7\x74\xd2\xf9\xfe\x37\x54\xe3\xd5\x50\xea\x1a\xaf\x27\x64\xb3\x07\xf7\x62\x5f\xbd\xcd\x43\x4f\x51\x5a\xbd\xfb\x50\x64\x9a\x38\xe8\x55\xf5\x32\x05\x53\x31\x0a\xc9\x23\x9f\x04\xd5\x93\x6c\x39\x78\x55\xfd\x3c\xd0\xfb\xf8\x82\x98\x76\x7c\xf6\xa8\x0d\x6c\x2e\xfd\x43\x6f\x4a\xde\x6b\x12\xc9\x3f\xe4\xcd\x68\xf0\xb5\xa8\x7a\x27\x72\xd2\xdd\x85\x63\xfb\x9d\xef\xd4\x56\xf1\x5e\xa2\xed\x99\x86\x82\x3a\x0e\x89\xa5\x71\xcd\x4b\x35\xec\x6f\xa5\xf2\xb7\xcd\xde\xfb\x33\xcc\xa0\x9a\x1b\x94\xef\xcb\xc2\xcf\xc7\xcb\xcf\x5f\xf7\x60\xa6\x72\xbf\xf7\x18\xf7\x26\x13\x54\x1a\x6b\xd0\xd2\x4c\x3b\x3d\xed\x71\xfd\x72\xea\x80\xb5\x5b\x79\xad\xde\xc6\xdd\xe8\xaf\xed\x7d\x3b\xe2\xf3\xeb\xf2\x96\xde\xfd\xbe\xdd\x9c\x3f\xf0\x8a\x6d\x98\x5a\xe3\x7e\xfd\xda\x9c\xb1\xdc\xc0\x8f\x63\x3d\x44\x2a\xa7\x36\x6b\xbb\x1e\x4d\xa0\x71\xfd\x38\x3a\xf4\x68\xb6\x19\xe0\x78\xf3\xe5\xca\x67\xef\x15\xc9\xf8\x1f\x51\xb8\x83\xc3\xf5\xaf\x23\x14\x74\x9b\xde\xd5\x83\x77\x5f\x86\x3f\x79\x86\xe1\x88\x3b\xfc\xdc\xdf\xe9\x98\xfe\x0f\x37\xd5\xdf\xff\x06\x00\x00\xff\xff\xfe\x65\x2a\xb5\x71\x22\x00\x00")

func templatesServerParameterGotmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesServerParameterGotmpl,
		"templates/server/parameter.gotmpl",
	)
}

func templatesServerParameterGotmpl() (*asset, error) {
	bytes, err := templatesServerParameterGotmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/server/parameter.gotmpl", size: 8817, mode: os.FileMode(420), modTime: time.Unix(1435797486, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if (err != nil) {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"templates/model.gotmpl": templatesModelGotmpl,
	"templates/modelvalidator.gotmpl": templatesModelvalidatorGotmpl,
	"templates/server/builder.gotmpl": templatesServerBuilderGotmpl,
	"templates/server/main.gotmpl": templatesServerMainGotmpl,
	"templates/server/operation.gotmpl": templatesServerOperationGotmpl,
	"templates/server/parameter.gotmpl": templatesServerParameterGotmpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"templates": &bintree{nil, map[string]*bintree{
		"model.gotmpl": &bintree{templatesModelGotmpl, map[string]*bintree{
		}},
		"modelvalidator.gotmpl": &bintree{templatesModelvalidatorGotmpl, map[string]*bintree{
		}},
		"server": &bintree{nil, map[string]*bintree{
			"builder.gotmpl": &bintree{templatesServerBuilderGotmpl, map[string]*bintree{
			}},
			"main.gotmpl": &bintree{templatesServerMainGotmpl, map[string]*bintree{
			}},
			"operation.gotmpl": &bintree{templatesServerOperationGotmpl, map[string]*bintree{
			}},
			"parameter.gotmpl": &bintree{templatesServerParameterGotmpl, map[string]*bintree{
			}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        // File
        if err != nil {
                return RestoreAsset(dir, name)
        }
        // Dir
        for _, child := range children {
                err = RestoreAssets(dir, path.Join(name, child))
                if err != nil {
                        return err
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

