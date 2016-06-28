// +build !appengine

/*
Copyright 2013 The Camlistore Authors.

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
	"os"
	"time"
)

// DieOnParentDeath starts a goroutine that regularly checks that
// the current process can find its parent, and calls os.Exit(0)
// as soon as it cannot.
func DieOnParentDeath() {
	// TODO: on Linux, use PR_SET_PDEATHSIG later. For now, the portable way:
	go func() {
		pollParent(30 * time.Second)
		os.Exit(0)
	}()
}

// pollParent checks every t that the ppid of the current
// process has not changed (i.e that the process has not
// been orphaned). It returns as soon as that ppid changes.
func pollParent(t time.Duration) {
	for initial := os.Getppid(); initial == os.Getppid(); time.Sleep(t) {
	}
}
