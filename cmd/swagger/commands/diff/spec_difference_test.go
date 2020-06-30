package diff_test

import (
	"testing"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/diff"
	"github.com/stretchr/testify/require"
)

func TestMatches(t *testing.T) {
	urlOnly := diff.SpecDifference{DifferenceLocation: diff.DifferenceLocation{URL: "bob"}}
	urlOnlyDiff := diff.SpecDifference{DifferenceLocation: diff.DifferenceLocation{URL: "notbob"}}
	urlOnlySame := diff.SpecDifference{DifferenceLocation: diff.DifferenceLocation{URL: "bob"}}

	require.True(t, urlOnly.Matches(urlOnlySame))
	require.False(t, urlOnly.Matches(urlOnlyDiff))

	withMethod := urlOnly
	withMethod.DifferenceLocation.Method = "PUT"
	withMethodSame := withMethod
	withMethodDiff := withMethod
	withMethodDiff.DifferenceLocation.Method = "GET"

	require.True(t, withMethod.Matches(withMethodSame))
	require.False(t, withMethod.Matches(withMethodDiff))

	withResponse := urlOnly
	withResponse.DifferenceLocation.Response = 0
	withResponseSame := withResponse
	withResponseDiff := withResponse
	withResponseDiff.DifferenceLocation.Response = 2

	require.True(t, withResponse.Matches(withResponseSame))
	require.False(t, withResponse.Matches(withResponseDiff))

	withNode := urlOnly
	withNode.DifferenceLocation.Node = &diff.Node{Field: "FieldA", TypeName: "TypeA"}
	withNodeSame := withNode
	withNodeSame.DifferenceLocation.Node = &diff.Node{Field: "FieldA", TypeName: "TypeA"}

	withNodeDiff := withNode
	withNodeDiff.DifferenceLocation.Node = &diff.Node{Field: "FieldA", TypeName: "TypeB"}

	require.True(t, withNode.Matches(withNodeSame))
	require.False(t, withNode.Matches(withNodeDiff))

	withNodeDiff.DifferenceLocation.Node = &diff.Node{Field: "FieldB", TypeName: "TypeA"}

	require.True(t, withNode.Matches(withNodeSame))
	require.False(t, withNode.Matches(withNodeDiff))

	withNestedNode := withNode
	withNestedNode.DifferenceLocation = withNestedNode.DifferenceLocation.AddNode(&diff.Node{Field: "ChildA", TypeName: "ChildA"})
	withNestedNodeSame := withNode
	withNestedNodeSame.DifferenceLocation = withNestedNodeSame.DifferenceLocation.AddNode(&diff.Node{Field: "ChildA", TypeName: "ChildA"})
	withNestedNodeDiff := withNode
	withNestedNodeDiff.DifferenceLocation = withNestedNodeDiff.DifferenceLocation.AddNode(&diff.Node{Field: "ChildB", TypeName: "ChildA"})

	require.True(t, withNestedNode.Matches(withNestedNodeSame))
	require.False(t, withNestedNode.Matches(withNestedNodeDiff))
}
