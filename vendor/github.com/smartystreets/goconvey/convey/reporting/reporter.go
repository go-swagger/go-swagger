
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

import "io"

type Reporter interface {
	BeginStory(story *StoryReport)
	Enter(scope *ScopeReport)
	Report(r *AssertionResult)
	Exit()
	EndStory()
	io.Writer
}

type reporters struct{ collection []Reporter }

func (self *reporters) BeginStory(s *StoryReport) { self.foreach(func(r Reporter) { r.BeginStory(s) }) }
func (self *reporters) Enter(s *ScopeReport)      { self.foreach(func(r Reporter) { r.Enter(s) }) }
func (self *reporters) Report(a *AssertionResult) { self.foreach(func(r Reporter) { r.Report(a) }) }
func (self *reporters) Exit()                     { self.foreach(func(r Reporter) { r.Exit() }) }
func (self *reporters) EndStory()                 { self.foreach(func(r Reporter) { r.EndStory() }) }

func (self *reporters) Write(contents []byte) (written int, err error) {
	self.foreach(func(r Reporter) {
		written, err = r.Write(contents)
	})
	return written, err
}

func (self *reporters) foreach(action func(Reporter)) {
	for _, r := range self.collection {
		action(r)
	}
}

func NewReporters(collection ...Reporter) *reporters {
	self := new(reporters)
	self.collection = collection
	return self
}
