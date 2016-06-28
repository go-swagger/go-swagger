/*
Copyright 2014 Google Inc.

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

// Package osutil provides operating system-specific path information,
// and other utility functions.
package osutil // import "camlistore.org/pkg/osutil"

import (
	"errors"
	"os"
)

// ErrNotSupported is returned by functions (like Mkfifo and Mksocket)
// when the underlying operating system or environment doesn't support
// the operation.
var ErrNotSupported = errors.New("operation not supported")

// DirExists reports whether dir exists. Errors are ignored and are
// reported as false.
func DirExists(dir string) bool {
	fi, err := os.Stat(dir)
	return err == nil && fi.IsDir()
}
