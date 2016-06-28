// +build windows

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
	"log"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

// SelfPath returns the path of the executable for the currently running
// process.
func SelfPath() (string, error) {
	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return "", err
	}
	sysproc, err := kernel32.FindProc("GetModuleFileNameW")
	if err != nil {
		return "", err
	}
	b := make([]uint16, syscall.MAX_PATH)
	r, _, err := sysproc.Call(0, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
	n := uint32(r)
	if n == 0 {
		return "", err
	}
	return string(utf16.Decode(b[0:n])), nil
}

// RestartProcess returns an error if things couldn't be
// restarted.  On success, this function never returns
// because the process becomes the new process.
func RestartProcess(args ...string) error {
	log.Print("RestartProcess not implemented on this platform.")
	return nil
}
