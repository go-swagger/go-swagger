// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// Debug when the env var DEBUG or SWAGGER_DEBUG is not empty
// the generators will be very noisy about what they are doing
var Debug = os.Getenv("DEBUG") != "" || os.Getenv("SWAGGER_DEBUG") != ""

func logDebug(frmt string, args ...interface{}) {
	if Debug {
		_, file, pos, _ := runtime.Caller(2)
		log.Printf("%s:%d: %s", filepath.Base(file), pos, fmt.Sprintf(frmt, args...))
	}
}

// debuglog is used to debug the typeResolver (types.go)
func debugLog(format string, args ...interface{}) {
	if Debug {
		_, file, pos, _ := runtime.Caller(2)
		log.Printf("%s:%d: "+format, append([]interface{}{filepath.Base(file), pos}, args...)...)
	}
}
