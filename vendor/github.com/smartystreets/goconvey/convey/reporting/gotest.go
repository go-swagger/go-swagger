
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

type gotestReporter struct{ test T }

func (self *gotestReporter) BeginStory(story *StoryReport) {
	self.test = story.Test
}

func (self *gotestReporter) Enter(scope *ScopeReport) {}

func (self *gotestReporter) Report(r *AssertionResult) {
	if !passed(r) {
		self.test.Fail()
	}
}

func (self *gotestReporter) Exit() {}

func (self *gotestReporter) EndStory() {
	self.test = nil
}

func (self *gotestReporter) Write(content []byte) (written int, err error) {
	return len(content), nil // no-op
}

func NewGoTestReporter() *gotestReporter {
	return new(gotestReporter)
}

func passed(r *AssertionResult) bool {
	return r.Error == nil && r.Failure == ""
}
