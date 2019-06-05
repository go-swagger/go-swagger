package diff_test

import (
	"testing"

	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/go-swagger/go-swagger/cmd/swagger/commands/diff"
)

var assertThat = then.AssertThat
var equals = is.EqualTo

func TestMatches(t *testing.T) {
	urlOnly := diff.SpecDifference{DifferenceLocation: diff.DifferenceLocation{URL: "bob"}}
	urlOnlyDiff := diff.SpecDifference{DifferenceLocation: diff.DifferenceLocation{URL: "notbob"}}
	urlOnlySame := diff.SpecDifference{DifferenceLocation: diff.DifferenceLocation{URL: "bob"}}

	assertThat(t, urlOnly.Matches(urlOnlySame), equals(true))
	assertThat(t, urlOnly.Matches(urlOnlyDiff), equals(false))

	withMethod := urlOnly
	withMethod.DifferenceLocation.Method = "PUT"
	withMethodSame := withMethod
	withMethodDiff := withMethod
	withMethodDiff.DifferenceLocation.Method = "GET"

	assertThat(t, withMethod.Matches(withMethodSame), equals(true))
	assertThat(t, withMethod.Matches(withMethodDiff), equals(false))

	withResponse := urlOnly
	withResponse.DifferenceLocation.Response = 0
	withResponseSame := withResponse
	withResponseDiff := withResponse
	withResponseDiff.DifferenceLocation.Response = 2

	assertThat(t, withResponse.Matches(withResponseSame), equals(true))
	assertThat(t, withResponse.Matches(withResponseDiff), equals(false))

	withNode := urlOnly
	withNode.DifferenceLocation.Node = &diff.Node{Field: "FieldA", TypeName: "TypeA"}
	withNodeSame := withNode
	withNodeSame.DifferenceLocation.Node = &diff.Node{Field: "FieldA", TypeName: "TypeA"}

	withNodeDiff := withNode
	withNodeDiff.DifferenceLocation.Node = &diff.Node{Field: "FieldA", TypeName: "TypeB"}

	assertThat(t, withNode.Matches(withNodeSame), equals(true))
	assertThat(t, withNode.Matches(withNodeDiff), equals(false))

	withNodeDiff.DifferenceLocation.Node = &diff.Node{Field: "FieldB", TypeName: "TypeA"}

	assertThat(t, withNode.Matches(withNodeSame), equals(true))
	assertThat(t, withNode.Matches(withNodeDiff), equals(false))

	withNestedNode := withNode
	withNestedNode.DifferenceLocation = withNestedNode.DifferenceLocation.AddNode(&diff.Node{Field: "ChildA", TypeName: "ChildA"})
	withNestedNodeSame := withNode
	withNestedNodeSame.DifferenceLocation = withNestedNodeSame.DifferenceLocation.AddNode(&diff.Node{Field: "ChildA", TypeName: "ChildA"})
	withNestedNodeDiff := withNode
	withNestedNodeDiff.DifferenceLocation = withNestedNodeDiff.DifferenceLocation.AddNode(&diff.Node{Field: "ChildB", TypeName: "ChildA"})

	assertThat(t, withNestedNode.Matches(withNestedNodeSame), equals(true))
	assertThat(t, withNestedNode.Matches(withNestedNodeDiff), equals(false))

}
