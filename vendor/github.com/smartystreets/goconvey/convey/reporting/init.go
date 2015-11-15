
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

package reporting

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func init() {
	if !isXterm() {
		monochrome()
	}

	if runtime.GOOS == "windows" {
		success, failure, error_ = dotSuccess, dotFailure, dotError
	}
}

func BuildJsonReporter() Reporter {
	out := NewPrinter(NewConsole())
	return NewReporters(
		NewGoTestReporter(),
		NewJsonReporter(out))
}
func BuildDotReporter() Reporter {
	out := NewPrinter(NewConsole())
	return NewReporters(
		NewGoTestReporter(),
		NewDotReporter(out),
		NewProblemReporter(out),
		consoleStatistics)
}
func BuildStoryReporter() Reporter {
	out := NewPrinter(NewConsole())
	return NewReporters(
		NewGoTestReporter(),
		NewStoryReporter(out),
		NewProblemReporter(out),
		consoleStatistics)
}
func BuildSilentReporter() Reporter {
	out := NewPrinter(NewConsole())
	return NewReporters(
		NewGoTestReporter(),
		NewProblemReporter(out))
}

var (
	newline         = "\n"
	success         = "✔"
	failure         = "✘"
	error_          = "🔥"
	skip            = "⚠"
	dotSuccess      = "."
	dotFailure      = "x"
	dotError        = "E"
	dotSkip         = "S"
	errorTemplate   = "* %s \nLine %d: - %v \n%s\n"
	failureTemplate = "* %s \nLine %d:\n%s\n"
)

var (
	greenColor  = "\033[32m"
	yellowColor = "\033[33m"
	redColor    = "\033[31m"
	resetColor  = "\033[0m"
)

var consoleStatistics = NewStatisticsReporter(NewPrinter(NewConsole()))

func SuppressConsoleStatistics() { consoleStatistics.Suppress() }
func PrintConsoleStatistics()    { consoleStatistics.PrintSummary() }

// QuiteMode disables all console output symbols. This is only meant to be used
// for tests that are internal to goconvey where the output is distracting or
// otherwise not needed in the test output.
func QuietMode() {
	success, failure, error_, skip, dotSuccess, dotFailure, dotError, dotSkip = "", "", "", "", "", "", "", ""
}

func monochrome() {
	greenColor, yellowColor, redColor, resetColor = "", "", "", ""
}

func isXterm() bool {
	env := fmt.Sprintf("%v", os.Environ())
	return strings.Contains(env, " TERM=isXterm") ||
		strings.Contains(env, " TERM=xterm")
}

// This interface allows us to pass the *testing.T struct
// throughout the internals of this tool without ever
// having to import the "testing" package.
type T interface {
	Fail()
}
