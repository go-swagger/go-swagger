package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDifferenceLocation_AddNode(t *testing.T) {

	parentLocation := DifferenceLocation{URL: "http://bob", Method: "meth", Node: &Node{Field: "Parent", TypeName: "bobtype"}}

	newLocation := parentLocation.AddNode(&Node{Field: "child1"})
	assert.Equal(t, newLocation.Node.ChildNode.Field, "child1")

	newLocation2 := parentLocation.AddNode(&Node{Field: "child2"})
	assert.Equal(t, newLocation2.Node.ChildNode.Field, "child2")

}
