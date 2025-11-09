// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDifferenceLocation_AddNode(t *testing.T) {
	parentLocation := DifferenceLocation{URL: "http://bob", Method: "meth", Node: &Node{Field: "Parent", TypeName: "bobtype"}}

	newLocation := parentLocation.AddNode(&Node{Field: "child1"})
	assert.Equal(t, "child1", newLocation.Node.ChildNode.Field)

	newLocation2 := parentLocation.AddNode(&Node{Field: "child2"})
	assert.Equal(t, "child2", newLocation2.Node.ChildNode.Field)
}
