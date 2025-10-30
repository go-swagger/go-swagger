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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	// Debug when the env var DEBUG or SWAGGER_DEBUG is not empty
	// the generators will be very noisy about what they are doing.
	Debug = os.Getenv("DEBUG") != "" || os.Getenv("SWAGGER_DEBUG") != ""
	// generatorLogger is a debug logger for this package.
	generatorLogger *log.Logger
)

func debugOptions() {
	generatorLogger = log.New(os.Stdout, "generator:", log.LstdFlags)
}

// debugLog wraps log.Printf with a debug-specific logger.
func debugLog(frmt string, args ...any) {
	if Debug {
		_, file, pos, _ := runtime.Caller(1)
		safeArgs := sanitizeDebugLogArgs(args...)
		generatorLogger.Printf("%s:%d: %s", filepath.Base(file), pos,
			fmt.Sprintf(frmt, safeArgs...))
	}
}

// debugLogAsJSON unmarshals its last arg as pretty JSON.
func debugLogAsJSON(frmt string, args ...any) {
	if Debug {
		var dfrmt string
		_, file, pos, _ := runtime.Caller(1)
		dargs := make([]any, 0, len(args)+2)
		dargs = append(dargs, filepath.Base(file), pos)

		if len(args) > 0 {
			dfrmt = "%s:%d: " + frmt + "\n%s"
			bbb, _ := json.MarshalIndent(args[len(args)-1], "", " ") //nolint:errchkjson // it's okay for debug
			dargs = append(dargs, args[0:len(args)-1]...)
			dargs = append(dargs, string(bbb))
		} else {
			dfrmt = "%s:%d: " + frmt
		}

		generatorLogger.Printf(dfrmt, dargs...)
	}
}

// sanitizeDebugLogArgs traverses arguments to debugLog and redacts fields
// that may contain sensitive information, such as API keys or credentials.
func sanitizeDebugLogArgs(args ...any) []any {
	safeArgs := make([]any, len(args))
	for i, arg := range args {
		safeArgs[i] = sanitizeValue(arg)
	}
	return safeArgs
}

// sanitizeValue redacts sensitive information from known data structures.
// It can be expanded for more types over time as needed.
func sanitizeValue(val any) any {
	switch v := val.(type) {
	case map[string]any:
		// Recursively sanitize map values
		res := make(map[string]any, len(v))
		for k, subv := range v {
			if k == "IsAPIKeyAuth" || k == "TokenURL" { // false positive: this is a bool indicator, not a sensitive value
				continue
			}
			lower := strings.ToLower(k)
			if lower == "apikey" || lower == "token" ||
				lower == "secret" ||
				strings.Contains(lower, "password") ||
				strings.Contains(lower, "apikey") ||
				strings.Contains(lower, "token") {
				res[k] = "***REDACTED***"

				continue
			}

			res[k] = sanitizeValue(subv)
		}
		return res
	case []any:
		res := make([]any, len(v))
		for i, subv := range v {
			res[i] = sanitizeValue(subv)
		}
		return res
	case string:
		// heuristic: redact if looks like a key/secret
		lower := strings.ToLower(v)
		if strings.Contains(lower, "apikey") || strings.Contains(lower, "token") || strings.Contains(lower, "secret") ||
			strings.Contains(lower, "password") {
			return "***REDACTED***"
		}
		return v
	default:
		// Optionally, process struct types for known sensitive fields
		return v
	}
}
