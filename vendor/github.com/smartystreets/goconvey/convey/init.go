
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

package convey

import (
	"flag"
	"os"

	"github.com/jtolds/gls"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey/reporting"
)

func init() {
	assertions.GoConveyMode(true)

	declareFlags()

	ctxMgr = gls.NewContextManager()
}

func declareFlags() {
	flag.BoolVar(&json, "json", false, "When true, emits results in JSON blocks. Default: 'false'")
	flag.BoolVar(&silent, "silent", false, "When true, all output from GoConvey is suppressed.")
	flag.BoolVar(&story, "story", false, "When true, emits story output, otherwise emits dot output. When not provided, this flag mirros the value of the '-test.v' flag")

	if noStoryFlagProvided() {
		story = verboseEnabled
	}

	// FYI: flag.Parse() is called from the testing package.
}

func noStoryFlagProvided() bool {
	return !story && !storyDisabled
}

func buildReporter() reporting.Reporter {
	switch {
	case testReporter != nil:
		return testReporter
	case json:
		return reporting.BuildJsonReporter()
	case silent:
		return reporting.BuildSilentReporter()
	case story:
		return reporting.BuildStoryReporter()
	default:
		return reporting.BuildDotReporter()
	}
}

var (
	ctxMgr *gls.ContextManager

	// only set by internal tests
	testReporter reporting.Reporter
)

var (
	json   bool
	silent bool
	story  bool

	verboseEnabled = flagFound("-test.v=true")
	storyDisabled  = flagFound("-story=false")
)

// flagFound parses the command line args manually for flags defined in other
// packages. Like the '-v' flag from the "testing" package, for instance.
func flagFound(flagValue string) bool {
	for _, arg := range os.Args {
		if arg == flagValue {
			return true
		}
	}
	return false
}
