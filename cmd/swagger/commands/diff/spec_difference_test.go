package diff

import (
	"testing"

	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
)

var assertThat = then.AssertThat
var equals = is.EqualTo

func TestMatches(t *testing.T) {
	urlOnly := SpecDifference{DifferenceLocation: DifferenceLocation{URL: "bob"}}
	urlOnlyDiff := SpecDifference{DifferenceLocation: DifferenceLocation{URL: "notbob"}}
	urlOnlySame := SpecDifference{DifferenceLocation: DifferenceLocation{URL: "bob"}}

	assertThat(t, urlOnly.matches(urlOnlySame), equals(true))
	assertThat(t, urlOnly.matches(urlOnlyDiff), equals(false))

	withMethod := urlOnly
	withMethod.DifferenceLocation.Method = "PUT"
	withMethodSame := withMethod
	withMethodDiff := withMethod
	withMethodDiff.DifferenceLocation.Method = "GET"

	assertThat(t, withMethod.matches(withMethodSame), equals(true))
	assertThat(t, withMethod.matches(withMethodDiff), equals(false))

	withResponse := urlOnly
	withResponse.DifferenceLocation.Response = 0
	withResponseSame := withResponse
	withResponseDiff := withResponse
	withResponseDiff.DifferenceLocation.Response = 2

	assertThat(t, withResponse.matches(withResponseSame), equals(true))
	assertThat(t, withResponse.matches(withResponseDiff), equals(false))

	withNode := urlOnly
	withNode.DifferenceLocation.Node = &Node{Field: "FieldA", TypeName: "TypeA"}
	withNodeSame := withNode
	withNodeSame.DifferenceLocation.Node = &Node{Field: "FieldA", TypeName: "TypeA"}

	withNodeDiff := withNode
	withNodeDiff.DifferenceLocation.Node = &Node{Field: "FieldA", TypeName: "TypeB"}

	assertThat(t, withNode.matches(withNodeSame), equals(true))
	assertThat(t, withNode.matches(withNodeDiff), equals(false))

	withNodeDiff.DifferenceLocation.Node = &Node{Field: "FieldB", TypeName: "TypeA"}

	assertThat(t, withNode.matches(withNodeSame), equals(true))
	assertThat(t, withNode.matches(withNodeDiff), equals(false))

	withNestedNode := withNode
	withNestedNode.DifferenceLocation = withNestedNode.DifferenceLocation.AddNode(&Node{Field: "ChildA", TypeName: "ChildA"})
	withNestedNodeSame := withNode
	withNestedNodeSame.DifferenceLocation = withNestedNodeSame.DifferenceLocation.AddNode(&Node{Field: "ChildA", TypeName: "ChildA"})
	withNestedNodeDiff := withNode
	withNestedNodeDiff.DifferenceLocation = withNestedNodeDiff.DifferenceLocation.AddNode(&Node{Field: "ChildB", TypeName: "ChildA"})

	assertThat(t, withNestedNode.matches(withNestedNodeSame), equals(true))
	assertThat(t, withNestedNode.matches(withNestedNodeDiff), equals(false))

}
