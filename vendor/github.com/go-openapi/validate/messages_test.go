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

package validate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var (
	// This debug environment variable allows to report and capture actual validation messages
	// during testing. It should be disabled (undefined) during CI tests.
	DebugTest = os.Getenv("SWAGGER_DEBUG_TEST") != ""
)

func init() {
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
}

type ExpectedMessage struct {
	Message              string `yaml:"message"`
	WithContinueOnErrors bool   `yaml:"withContinueOnErrors"` // should be expected only when SetContinueOnErrors(true)
	IsRegexp             bool   `yaml:"isRegexp"`             // expected message is interpreted as regexp (with regexp.MatchString())
}

type ExpectedFixture struct {
	Comment           string            `yaml:"comment,omitempty"`
	Todo              string            `yaml:"todo,omitempty"`
	ExpectedLoadError bool              `yaml:"expectedLoadError"` // expect error on load: skip validate step
	ExpectedValid     bool              `yaml:"expectedValid"`     // expect valid spec
	ExpectedMessages  []ExpectedMessage `yaml:"expectedMessages"`
	ExpectedWarnings  []ExpectedMessage `yaml:"expectedWarnings"`
}

type ExpectedMap map[string]ExpectedFixture

// Test message improvements, issue #44 and some more
// ContinueOnErrors mode on
// WARNING: this test is very demanding and constructed with varied scenarios,
// which are not necessarily "unitary". Expect multiple changes in messages whenever
// altering the validator.
func Test_MessageQualityContinueOnErrors_Issue44(t *testing.T) {
	if !enableLongTests {
		skipNotify(t)
		t.SkipNow()
	}
	state := continueOnErrors
	SetContinueOnErrors(true)
	defer func() {
		SetContinueOnErrors(state)
	}()
	errs := testMessageQuality(t, true) /* set haltOnErrors=true to iterate spec by spec */
	assert.Zero(t, errs, "Message testing didn't match expectations")
}

// ContinueOnErrors mode off
func Test_MessageQualityStopOnErrors_Issue44(t *testing.T) {
	if !enableLongTests {
		skipNotify(t)
		t.SkipNow()
	}
	state := continueOnErrors
	SetContinueOnErrors(false)
	defer func() {
		SetContinueOnErrors(state)
	}()
	errs := testMessageQuality(t, true) /* set haltOnErrors=true to iterate spec by spec */
	assert.Zero(t, errs, "Message testing didn't match expectations")
}

// Verifies the production of validation error messages in multiple
// spec scenarios.
//
// The objective is to demonstrate:
// - messages are stable
// - validation continues as much as possible, even in presence of many errors
//
// haltOnErrors is used in dev mode to study and fix testcases step by step (output is pretty verbose)
//
// set SWAGGER_DEBUG_TEST=1 env to get a report of messages at the end of each test.
// expectedMessage{"", false, false},
//
// expected messages and warnings are configured in ./fixtures/validation/expected_messages.yaml
//
func testMessageQuality(t *testing.T, haltOnErrors bool) (errs int) {
	expectedConfig, ferr := ioutil.ReadFile("./fixtures/validation/expected_messages.yaml")
	if ferr != nil {
		t.Logf("Cannot read expected messages. Skip this test: %v", ferr)
		return
	}

	tested := ExpectedMap{}
	yerr := yaml.Unmarshal(expectedConfig, &tested)
	if yerr != nil {
		t.Logf("Cannot unmarshall expected messages. Skip this test: %v", yerr)
		return
	}
	err := filepath.Walk(filepath.Join("fixtures", "validation"),
		func(path string, info os.FileInfo, err error) error {
			_, found := tested[info.Name()]
			errs := 0
			if !info.IsDir() && found && tested[info.Name()].ExpectedValid == false {
				// Checking invalid specs
				t.Logf("Testing messages for invalid spec: %s", path)
				if DebugTest {
					if tested[info.Name()].Comment != "" {
						t.Logf("\tDEVMODE: Comment: %s", tested[info.Name()].Comment)
					}
					if tested[info.Name()].Todo != "" {
						t.Logf("\tDEVMODE: Todo: %s", tested[info.Name()].Todo)
					}
				}
				doc, err := loads.Spec(path)

				// Check specs with load errors (error is located in pkg loads or spec)
				if tested[info.Name()].ExpectedLoadError == true {
					// Expect a load error: no further validation may possibly be conducted.
					if assert.Error(t, err, "Expected this spec to return a load error") {
						errs += verifyLoadErrors(t, err, tested[info.Name()].ExpectedMessages)
						if errs == 0 {
							// spec does not load as expected
							return nil
						}
					} else {
						errs++
					}
				}
				if errs > 0 {
					if haltOnErrors {
						return fmt.Errorf("Test halted: stop on error mode")
					}
					return nil
				}

				if assert.NoError(t, err, "Expected this spec to load properly") {
					// Validate the spec document
					validator := NewSpecValidator(doc.Schema(), strfmt.Default)
					res, warn := validator.Validate(doc)

					// Check specs with load errors (error is located in pkg loads or spec)
					if !assert.False(t, res.IsValid(), "Expected this spec to be invalid") {
						errs++
					}

					errs += verifyErrorsVsWarnings(t, res, warn)
					errs += verifyErrors(t, res, tested[info.Name()].ExpectedMessages, "error")
					errs += verifyErrors(t, warn, tested[info.Name()].ExpectedWarnings, "warning")

					// DEVMODE allows developers to experiment and tune expected results
					if DebugTest && errs > 0 {
						reportTest(t, path, res, tested[info.Name()].ExpectedMessages, "error")
						reportTest(t, path, warn, tested[info.Name()].ExpectedWarnings, "warning")
					}
				} else {
					errs++
				}

				if errs > 0 {
					t.Logf("Message qualification on Spec validation failed for %s", path)
				}
			} else {
				// Expecting no message (e.g.valid spec): 0 message expected
				if !info.IsDir() && found && tested[info.Name()].ExpectedValid {
					t.Logf("Testing valid spec: %s", path)
					if DebugTest {
						if tested[info.Name()].Comment != "" {
							t.Logf("\tDEVMODE: Comment: %s", tested[info.Name()].Comment)
						}
						if tested[info.Name()].Todo != "" {
							t.Logf("\tDEVMODE: Todo: %s", tested[info.Name()].Todo)
						}
					}
					doc, err := loads.Spec(path)
					if assert.NoError(t, err, "Expected this spec to load without error") {
						validator := NewSpecValidator(doc.Schema(), strfmt.Default)
						res, warn := validator.Validate(doc)
						if !assert.True(t, res.IsValid(), "Expected this spec to be valid") {
							errs++
						}
						errs += verifyErrors(t, warn, tested[info.Name()].ExpectedWarnings, "warning")
						if DebugTest && errs > 0 {
							reportTest(t, path, res, tested[info.Name()].ExpectedMessages, "error")
							reportTest(t, path, warn, tested[info.Name()].ExpectedWarnings, "warning")
						}
					} else {
						errs++
					}
				}
			}
			if haltOnErrors && errs > 0 {
				return fmt.Errorf("Test halted: stop on error mode")
			}
			return nil
		})
	if err != nil {
		t.Logf("%v", err)
		errs++
	}
	return
}

// Prints out a recap of error messages. To be enabled during development / test iterations
func reportTest(t *testing.T, path string, res *Result, expectedMessages []ExpectedMessage, msgtype string) {
	var verifiedErrors, lines []string
	for _, e := range res.Errors {
		verifiedErrors = append(verifiedErrors, e.Error())
	}
	t.Logf("DEVMODE:Recap of returned %s messages while validating %s ", msgtype, path)
	for _, v := range verifiedErrors {
		status := fmt.Sprintf("Unexpected %s", msgtype)
		for _, s := range expectedMessages {
			if (s.WithContinueOnErrors == true && continueOnErrors == true) || s.WithContinueOnErrors == false {
				if s.IsRegexp {
					if matched, _ := regexp.MatchString(s.Message, v); matched {
						status = fmt.Sprintf("Expected %s", msgtype)
						break
					}
				} else {
					if strings.Contains(v, s.Message) {
						status = fmt.Sprintf("Expected %s", msgtype)
						break
					}
				}
			}
		}
		lines = append(lines, fmt.Sprintf("[%s]%s", status, v))
	}

	for _, s := range expectedMessages {
		if (s.WithContinueOnErrors == true && continueOnErrors == true) || s.WithContinueOnErrors == false {
			status := fmt.Sprintf("Missing %s", msgtype)
			for _, v := range verifiedErrors {
				if s.IsRegexp {
					if matched, _ := regexp.MatchString(s.Message, v); matched {
						status = fmt.Sprintf("Expected %s", msgtype)
						break
					}
				} else {
					if strings.Contains(v, s.Message) {
						status = fmt.Sprintf("Expected %s", msgtype)
						break
					}
				}
			}
			if status != fmt.Sprintf("Expected %s", msgtype) {
				lines = append(lines, fmt.Sprintf("[%s]%s", status, s.Message))
			}
		}
	}
	if len(lines) > 0 {
		sort.Strings(lines)
		for _, line := range lines {
			t.Logf(line)
		}
	}
}

func verifyErrorsVsWarnings(t *testing.T, res, warn *Result) (errs int) {
	// First verification of result conventions: results are redundant, just a matter of presentation
	w := len(warn.Errors)
	if !assert.Len(t, res.Warnings, w) {
		errs++
	}
	if !assert.Len(t, warn.Warnings, 0) {
		errs++
	}
	if !assert.Subset(t, res.Warnings, warn.Errors) {
		errs++
	}
	if !assert.Subset(t, warn.Errors, res.Warnings) {
		errs++
	}
	return
}

func verifyErrors(t *testing.T, res *Result, expectedMessages []ExpectedMessage, msgtype string) (errs int) {
	var verifiedErrors []string
	var numExpected int

	for _, e := range res.Errors {
		verifiedErrors = append(verifiedErrors, e.Error())
	}
	for _, s := range expectedMessages {
		if (s.WithContinueOnErrors == true && continueOnErrors == true) || s.WithContinueOnErrors == false {
			numExpected++
		}
	}

	// We got the expected number of messages (e.g. no duplicates, no uncontrolled side-effect, ...)
	if !assert.Len(t, verifiedErrors, numExpected, "Unexpected number of %s messages returned. Wanted %d, got %d", msgtype, numExpected, len(verifiedErrors)) {
		errs++
	}

	// Check that all expected messages are here
	for _, s := range expectedMessages {
		found := false
		if (s.WithContinueOnErrors == true && continueOnErrors == true) || s.WithContinueOnErrors == false {
			for _, v := range verifiedErrors {
				if s.IsRegexp {
					if matched, _ := regexp.MatchString(s.Message, v); matched {
						found = true
						break
					}
				} else {
					if strings.Contains(v, s.Message) {
						found = true
						break
					}
				}
			}
			if !assert.True(t, found, "Missing expected %s message: %s", msgtype, s.Message) {
				errs++
			}
		}
	}

	// Check for no unexpected message
	for _, v := range verifiedErrors {
		found := false
		for _, s := range expectedMessages {
			if (s.WithContinueOnErrors == true && continueOnErrors == true) || s.WithContinueOnErrors == false {
				if s.IsRegexp {
					if matched, _ := regexp.MatchString(s.Message, v); matched {
						found = true
						break
					}
				} else {
					if strings.Contains(v, s.Message) {
						found = true
						break
					}
				}
			}
		}
		if !assert.True(t, found, "Unexpected %s message: %s", msgtype, v) {
			errs++
		}
	}
	return
}

// Perform several matchedes on single error message
// Process here error messages from loads (normally unit tested in the load package:
// we just want to figure out how all this is captured at the validate package level.
func verifyLoadErrors(t *testing.T, err error, expectedMessages []ExpectedMessage) (errs int) {
	v := err.Error()
	for _, s := range expectedMessages {
		found := false
		if s.IsRegexp {
			if matched, _ := regexp.MatchString(s.Message, v); matched {
				found = true
				break
			}
		} else {
			if strings.Contains(v, s.Message) {
				found = true
				break
			}
		}
		if !assert.True(t, found, "Unexpected load error: %s", v) {
			// DEBUG CI
			t.Logf("Expecting one of the following:")
			for _, s := range expectedMessages {
				smode := "Contains"
				if s.IsRegexp {
					smode = "MatchString"
				}
				t.Logf("[%s]:%s", smode, s.Message)
			}
			errs++
		}
	}
	return
}
