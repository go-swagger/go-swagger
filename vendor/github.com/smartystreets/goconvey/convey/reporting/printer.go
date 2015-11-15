
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
	"io"
	"strings"
)

type Printer struct {
	out    io.Writer
	prefix string
}

func (self *Printer) Println(message string, values ...interface{}) {
	formatted := self.format(message, values...) + newline
	self.out.Write([]byte(formatted))
}

func (self *Printer) Print(message string, values ...interface{}) {
	formatted := self.format(message, values...)
	self.out.Write([]byte(formatted))
}

func (self *Printer) Insert(text string) {
	self.out.Write([]byte(text))
}

func (self *Printer) format(message string, values ...interface{}) string {
	var formatted string
	if len(values) == 0 {
		formatted = self.prefix + message
	} else {
		formatted = self.prefix + fmt.Sprintf(message, values...)
	}
	indented := strings.Replace(formatted, newline, newline+self.prefix, -1)
	return strings.TrimRight(indented, space)
}

func (self *Printer) Indent() {
	self.prefix += pad
}

func (self *Printer) Dedent() {
	if len(self.prefix) >= padLength {
		self.prefix = self.prefix[:len(self.prefix)-padLength]
	}
}

func NewPrinter(out io.Writer) *Printer {
	self := new(Printer)
	self.out = out
	return self
}

const space = " "
const pad = space + space
const padLength = len(pad)
