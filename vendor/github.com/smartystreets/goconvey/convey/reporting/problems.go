
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

import "fmt"

type problem struct {
	out      *Printer
	errors   []*AssertionResult
	failures []*AssertionResult
}

func (self *problem) BeginStory(story *StoryReport) {}

func (self *problem) Enter(scope *ScopeReport) {}

func (self *problem) Report(report *AssertionResult) {
	if report.Error != nil {
		self.errors = append(self.errors, report)
	} else if report.Failure != "" {
		self.failures = append(self.failures, report)
	}
}

func (self *problem) Exit() {}

func (self *problem) EndStory() {
	self.show(self.showErrors, redColor)
	self.show(self.showFailures, yellowColor)
	self.prepareForNextStory()
}
func (self *problem) show(display func(), color string) {
	fmt.Print(color)
	display()
	fmt.Print(resetColor)
	self.out.Dedent()
}
func (self *problem) showErrors() {
	for i, e := range self.errors {
		if i == 0 {
			self.out.Println("\nErrors:\n")
			self.out.Indent()
		}
		self.out.Println(errorTemplate, e.File, e.Line, e.Error, e.StackTrace)
	}
}
func (self *problem) showFailures() {
	for i, f := range self.failures {
		if i == 0 {
			self.out.Println("\nFailures:\n")
			self.out.Indent()
		}
		self.out.Println(failureTemplate, f.File, f.Line, f.Failure)
	}
}

func (self *problem) Write(content []byte) (written int, err error) {
	return len(content), nil // no-op
}

func NewProblemReporter(out *Printer) *problem {
	self := new(problem)
	self.out = out
	self.prepareForNextStory()
	return self
}
func (self *problem) prepareForNextStory() {
	self.errors = []*AssertionResult{}
	self.failures = []*AssertionResult{}
}
