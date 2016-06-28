// +build solaris

/*
Copyright 2014 The Camlistore Authors.

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

package osutil

import (
	"net"
	"os"
	"path/filepath"
	"syscall"
)

func Mkfifo(path string, mode uint32) error {
	// Mkfifo is missing from syscall, thus call Mknod as it does on Linux.
	return syscall.Mknod(path, mode|syscall.S_IFIFO, 0)
}

// Mksocket creates a socket file (a Unix Domain Socket) named path.
func Mksocket(path string) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	tmp := filepath.Join(dir, "."+base)
	l, err := net.ListenUnix("unix", &net.UnixAddr{tmp, "unix"})
	if err != nil {
		return err
	}

	err = os.Rename(tmp, path)
	if err != nil {
		l.Close()
		os.Remove(tmp) // Ignore error
		return err
	}

	l.Close()

	return nil
}
