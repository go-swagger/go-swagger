// +build appengine

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
	"log"
)

// SelfPath returns the path of the executable for the currently running
// process.
func SelfPath() (string, error) {
	return "", errors.New("SelfPath not implemented on App Engine.")
}

// RestartProcess returns an error if things couldn't be
// restarted.  On success, this function never returns
// because the process becomes the new process.
func RestartProcess(arg ...string) error {
	log.Print("RestartProcess not implemented on this platform.")
	return nil
}
