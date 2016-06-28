// +build appengine

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
	"log"
)

func DieOnParentDeath() {
	// TODO(mpl): maybe the way it's done in findproc_normal.go actually works
	// on appengine too? Verify that.
	log.Fatal("DieOnParentDeath not implemented on appengine.")
}
