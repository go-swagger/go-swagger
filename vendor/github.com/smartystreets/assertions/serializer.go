
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

package assertions

import (
	"encoding/json"
	"fmt"

	"github.com/smartystreets/goconvey/convey/reporting"
)

type Serializer interface {
	serialize(expected, actual interface{}, message string) string
	serializeDetailed(expected, actual interface{}, message string) string
}

type failureSerializer struct{}

func (self *failureSerializer) serializeDetailed(expected, actual interface{}, message string) string {
	view := self.format(expected, actual, message, "%#v")
	serialized, err := json.Marshal(view)
	if err != nil {
		return message
	}
	return string(serialized)
}

func (self *failureSerializer) serialize(expected, actual interface{}, message string) string {
	view := self.format(expected, actual, message, "%+v")
	serialized, err := json.Marshal(view)
	if err != nil {
		return message
	}
	return string(serialized)
}

func (self *failureSerializer) format(expected, actual interface{}, message string, format string) reporting.FailureView {
	return reporting.FailureView{
		Message:  message,
		Expected: fmt.Sprintf(format, expected),
		Actual:   fmt.Sprintf(format, actual),
	}
}

func newSerializer() *failureSerializer {
	return &failureSerializer{}
}

///////////////////////////////////////////////////////

// noopSerializer just gives back the original message. This is useful when we are using
// the assertions from a context other than the web UI, that requires the JSON structure
// provided by the failureSerializer.
type noopSerializer struct{}

func (self *noopSerializer) serialize(expected, actual interface{}, message string) string {
	return message
}
func (self *noopSerializer) serializeDetailed(expected, actual interface{}, message string) string {
	return message
}
