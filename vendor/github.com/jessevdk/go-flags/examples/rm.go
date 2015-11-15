
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

package main

import (
	"fmt"
)

type RmCommand struct {
	Force bool `short:"f" long:"force" description:"Force removal of files"`
}

var rmCommand RmCommand

func (x *RmCommand) Execute(args []string) error {
	fmt.Printf("Removing (force=%v): %#v\n", x.Force, args)
	return nil
}

func init() {
	parser.AddCommand("rm",
		"Remove a file",
		"The rm command removes a file to the repository. Use -f to force removal of files.",
		&rmCommand)
}
