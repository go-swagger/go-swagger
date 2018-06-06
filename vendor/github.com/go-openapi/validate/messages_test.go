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
	Tested            bool              `yaml:"-"`
	Failed            bool              `yaml:"-"`
}

type ExpectedMap map[string]*ExpectedFixture

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
	errs := testMessageQuality(t, true, true) /* set haltOnErrors=true to iterate spec by spec */
	assert.Zero(t, errs, "Message testing didn't match expectations")
}

// ContinueOnErrors mode off
func Test_MessageQualityStopOnErrors_Issue44(t *testing.T) {
	if !enableLongTests {
		skipNotify(t)
		t.SkipNow()
	}
	errs := testMessageQuality(t, true, false) /* set haltOnErrors=true to iterate spec by spec */
	assert.Zero(t, errs, "Message testing didn't match expectations")
}

func testMessageQuality(t *testing.T, haltOnErrors bool, continueOnErrors bool) (errs int) {
	// Verifies the production of validation error messages in multiple
	// spec scenarios.
	//
	// The objective is to demonstrate that:
	//   - messages are stable
	//   - validation continues as much as possible, even in presence of many errors
	//
	// haltOnErrors is used in dev mode to study and fix testcases step by step (output is pretty verbose)
	//
	// set SWAGGER_DEBUG_TEST=1 env to get a report of messages at the end of each test.
	// expectedMessage{"", false, false},
	//
	// expected messages and warnings are configured in ./fixtures/validation/expected_messages.yaml
	//
	expectedConfig, ferr := ioutil.ReadFile("./fixtures/validation/expected_messages.yaml")
	if ferr != nil {
		t.Logf("Cannot read expected messages config file: %v", ferr)
		errs++
		return
	}

	tested := ExpectedMap{}
	yerr := yaml.Unmarshal(expectedConfig, &tested)
	if yerr != nil {
		t.Logf("Cannot unmarshall expected messages from config file : %v", yerr)
		errs++
		return
	}

	// Check config
	for fixture, expected := range tested {
		if err := UniqueItems("", "", expected.ExpectedMessages); err != nil {
			t.Logf("Duplicate messages configured for %s", fixture)
			errs++
		}
		if err := UniqueItems("", "", expected.ExpectedWarnings); err != nil {
			t.Logf("Duplicate messages configured for %s", fixture)
			errs++
		}
	}
	if errs > 0 {
		return
	}
	err := filepath.Walk(filepath.Join("fixtures", "validation"),
		func(path string, info os.FileInfo, err error) error {
			t.Run(path, func(t *testing.T) {
				t.Parallel()
				basename := info.Name()
				_, found := tested[basename]
				errs := 0

				defer func() {
					if found {
						tested[basename].Tested = true
						tested[basename].Failed = (errs != 0)
					}

				}()

				if !info.IsDir() && found && tested[basename].ExpectedValid == false {
					// Checking invalid specs
					t.Logf("Testing messages for invalid spec: %s", path)
					if DebugTest {
						if tested[basename].Comment != "" {
							t.Logf("\tDEVMODE: Comment: %s", tested[basename].Comment)
						}
						if tested[basename].Todo != "" {
							t.Logf("\tDEVMODE: Todo: %s", tested[basename].Todo)
						}
					}
					doc, err := loads.Spec(path)

					// Check specs with load errors (error is located in pkg loads or spec)
					if tested[basename].ExpectedLoadError == true {
						// Expect a load error: no further validation may possibly be conducted.
						if assert.Error(t, err, "Expected this spec to return a load error") {
							errs += verifyLoadErrors(t, err, tested[basename].ExpectedMessages)
							if errs == 0 {
								// spec does not load as expected
								return
							}
						} else {
							errs++
						}
					}
					if errs > 0 {
						if haltOnErrors {
							assert.FailNow(t, "Test halted: stop testing on message checking error mode")
							return
						}
						return
					}

					if assert.NoError(t, err, "Expected this spec to load properly") {
						// Validate the spec document
						validator := NewSpecValidator(doc.Schema(), strfmt.Default)
						validator.SetContinueOnErrors(continueOnErrors)
						res, warn := validator.Validate(doc)

						// Check specs with load errors (error is located in pkg loads or spec)
						if !assert.False(t, res.IsValid(), "Expected this spec to be invalid") {
							errs++
						}

						errs += verifyErrorsVsWarnings(t, res, warn)
						errs += verifyErrors(t, res, tested[basename].ExpectedMessages, "error", continueOnErrors)
						errs += verifyErrors(t, warn, tested[basename].ExpectedWarnings, "warning", continueOnErrors)

						// DEVMODE allows developers to experiment and tune expected results
						if DebugTest && errs > 0 {
							reportTest(t, path, res, tested[basename].ExpectedMessages, "error", continueOnErrors)
							reportTest(t, path, warn, tested[basename].ExpectedWarnings, "warning", continueOnErrors)
						}
					} else {
						errs++
					}

					if errs > 0 {
						t.Logf("Message qualification on Spec validation failed for %s", path)
					}
				} else {
					// Expecting no message (e.g.valid spec): 0 message expected
					if !info.IsDir() && found && tested[basename].ExpectedValid {
						tested[basename].Tested = true
						t.Logf("Testing valid spec: %s", path)
						if DebugTest {
							if tested[basename].Comment != "" {
								t.Logf("\tDEVMODE: Comment: %s", tested[basename].Comment)
							}
							if tested[basename].Todo != "" {
								t.Logf("\tDEVMODE: Todo: %s", tested[basename].Todo)
							}
						}
						doc, err := loads.Spec(path)
						if assert.NoError(t, err, "Expected this spec to load without error") {
							validator := NewSpecValidator(doc.Schema(), strfmt.Default)
							validator.SetContinueOnErrors(continueOnErrors)
							res, warn := validator.Validate(doc)
							if !assert.True(t, res.IsValid(), "Expected this spec to be valid") {
								errs++
							}
							errs += verifyErrors(t, warn, tested[basename].ExpectedWarnings, "warning", continueOnErrors)
							if DebugTest && errs > 0 {
								reportTest(t, path, res, tested[basename].ExpectedMessages, "error", continueOnErrors)
								reportTest(t, path, warn, tested[basename].ExpectedWarnings, "warning", continueOnErrors)
							}
						} else {
							errs++
						}
					}
				}
				if haltOnErrors && errs > 0 {
					assert.FailNow(t, "Test halted: stop testing on message checking error mode")
					return
				}
			})
			return nil
		})
	recapTest(t, tested)
	if err != nil {
		t.Logf("%v", err)
		errs++
	}
	return
}

func recapTest(t *testing.T, config ExpectedMap) {
	recapFailed := false
	for k, v := range config {
		if !v.Tested {
			t.Logf("WARNING: %s configured but not tested (fixture not found)", k)
			recapFailed = true
		} else if v.Failed {
			t.Logf("ERROR: %s failed passing messages verification", k)
			recapFailed = true
		}
	}
	if !recapFailed {
		t.Log("INFO:We are good")
	}
}
func reportTest(t *testing.T, path string, res *Result, expectedMessages []ExpectedMessage, msgtype string, continueOnErrors bool) {
	// Prints out a recap of error messages. To be enabled during development / test iterations
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
	if errs > 0 {
		t.Log("Result equivalence errors vs warnings not verified")
	}
	return
}

func verifyErrors(t *testing.T, res *Result, expectedMessages []ExpectedMessage, msgtype string, continueOnErrors bool) (errs int) {
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

func verifyLoadErrors(t *testing.T, err error, expectedMessages []ExpectedMessage) (errs int) {
	// Perform several matchedes on single error message
	// Process here error messages from loads (normally unit tested in the load package:
	// we just want to figure out how all this is captured at the validate package level.
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

// Test unitary fixture for dev and bug fixing
func Test_SingleFixture(t *testing.T) {
	t.SkipNow()
	path := filepath.Join("fixtures", "validation", "fixture-342.yaml")
	doc, err := loads.Spec(path)
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		validator.SetContinueOnErrors(true)
		res, _ := validator.Validate(doc)
		t.Log("Returned errors:")
		for _, e := range res.Errors {
			t.Logf("%v", e)
		}
		t.Log("Returned warnings:")
		for _, e := range res.Warnings {
			t.Logf("%v", e)
		}

	} else {
		t.Logf("Load error: %v", err)
	}
}
