
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

func (self *statistics) BeginStory(story *StoryReport) {}

func (self *statistics) Enter(scope *ScopeReport) {}

func (self *statistics) Report(report *AssertionResult) {
	if !self.failing && report.Failure != "" {
		self.failing = true
	}
	if !self.erroring && report.Error != nil {
		self.erroring = true
	}
	if report.Skipped {
		self.skipped += 1
	} else {
		self.total++
	}
}

func (self *statistics) Exit() {}

func (self *statistics) EndStory() {
	if !self.suppressed {
		self.PrintSummary()
	}
}

func (self *statistics) Suppress() {
	self.suppressed = true
}

func (self *statistics) PrintSummary() {
	self.reportAssertions()
	self.reportSkippedSections()
	self.completeReport()
}
func (self *statistics) reportAssertions() {
	self.decideColor()
	self.out.Print("\n%d total %s", self.total, plural("assertion", self.total))
}
func (self *statistics) decideColor() {
	if self.failing && !self.erroring {
		fmt.Print(yellowColor)
	} else if self.erroring {
		fmt.Print(redColor)
	} else {
		fmt.Print(greenColor)
	}
}
func (self *statistics) reportSkippedSections() {
	if self.skipped > 0 {
		fmt.Print(yellowColor)
		self.out.Print(" (one or more sections skipped)")
	}
}
func (self *statistics) completeReport() {
	fmt.Print(resetColor)
	self.out.Print("\n")
	self.out.Print("\n")
}

func (self *statistics) Write(content []byte) (written int, err error) {
	return len(content), nil // no-op
}

func NewStatisticsReporter(out *Printer) *statistics {
	self := statistics{}
	self.out = out
	return &self
}

type statistics struct {
	out        *Printer
	total      int
	failing    bool
	erroring   bool
	skipped    int
	suppressed bool
}

func plural(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}
