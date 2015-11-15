
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

type dot struct{ out *Printer }

func (self *dot) BeginStory(story *StoryReport) {}

func (self *dot) Enter(scope *ScopeReport) {}

func (self *dot) Report(report *AssertionResult) {
	if report.Error != nil {
		fmt.Print(redColor)
		self.out.Insert(dotError)
	} else if report.Failure != "" {
		fmt.Print(yellowColor)
		self.out.Insert(dotFailure)
	} else if report.Skipped {
		fmt.Print(yellowColor)
		self.out.Insert(dotSkip)
	} else {
		fmt.Print(greenColor)
		self.out.Insert(dotSuccess)
	}
	fmt.Print(resetColor)
}

func (self *dot) Exit() {}

func (self *dot) EndStory() {}

func (self *dot) Write(content []byte) (written int, err error) {
	return len(content), nil // no-op
}

func NewDotReporter(out *Printer) *dot {
	self := new(dot)
	self.out = out
	return self
}
