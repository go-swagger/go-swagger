// +build !appengine
// +build linux darwin freebsd netbsd openbsd solaris

/*
Copyright 2012 The Camlistore Authors.

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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

// if non-nil, osSelfPath is used from selfPath.
var osSelfPath func() (string, error)

// TODO(mpl): document the symlink behaviour in SelfPath for the BSDs when
// I know for sure.

// SelfPath returns the path of the executable for the currently running
// process. At least on linux, the returned path is a symlink to the actual
// executable.
func SelfPath() (string, error) {
	if f := osSelfPath; f != nil {
		return f()
	}
	switch runtime.GOOS {
	case "linux":
		return "/proc/self/exe", nil
	case "netbsd":
		return "/proc/curproc/exe", nil
	case "openbsd":
		return "/proc/curproc/file", nil
	case "darwin":
		// TODO(mpl): maybe do the right thing for darwin too, but that may require changes to runtime.
		// See https://codereview.appspot.com/6736069/
		return exec.LookPath(os.Args[0])
	}
	return "", errors.New("SelfPath not implemented for " + runtime.GOOS)
}

// RestartProcess restarts the process with the given arguments, if any,
// replacing the original process's arguments. It defaults to os.Args otherwise. It
// returns an error if things couldn't be restarted. On success, this function
// never returns because the process becomes the new process.
func RestartProcess(arg ...string) error {
	path, err := SelfPath()
	if err != nil {
		return fmt.Errorf("RestartProcess failed: %v", err)
	}

	var args []string
	if len(arg) > 0 {
		args = append(args, os.Args[0])
		for _, v := range arg {
			args = append(args, v)
		}
	} else {
		args = os.Args
	}

	return syscall.Exec(path, args, os.Environ())
}
