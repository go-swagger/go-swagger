// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package diff_test

import (
	"testing"

	"github.com/go-openapi/testify/v2/require"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/diff"
)

func TestMatches(t *testing.T) {
	urlOnly := diff.SpecDifference{DifferenceLocation: diff.DifferenceLocation{URL: "bob"}}
	urlOnlyDiff := diff.SpecDifference{DifferenceLocation: diff.DifferenceLocation{URL: "notbob"}}
	urlOnlySame := diff.SpecDifference{DifferenceLocation: diff.DifferenceLocation{URL: "bob"}}

	require.TrueT(t, urlOnly.Matches(urlOnlySame))
	require.FalseT(t, urlOnly.Matches(urlOnlyDiff))

	withMethod := urlOnly
	withMethod.DifferenceLocation.Method = "PUT"
	withMethodSame := withMethod
	withMethodDiff := withMethod
	withMethodDiff.DifferenceLocation.Method = "GET"

	require.TrueT(t, withMethod.Matches(withMethodSame))
	require.FalseT(t, withMethod.Matches(withMethodDiff))

	withResponse := urlOnly
	withResponse.DifferenceLocation.Response = 0
	withResponseSame := withResponse
	withResponseDiff := withResponse
	withResponseDiff.DifferenceLocation.Response = 2

	require.TrueT(t, withResponse.Matches(withResponseSame))
	require.FalseT(t, withResponse.Matches(withResponseDiff))

	withNode := urlOnly
	withNode.DifferenceLocation.Node = &diff.Node{Field: "FieldA", TypeName: "TypeA"}
	withNodeSame := withNode
	withNodeSame.DifferenceLocation.Node = &diff.Node{Field: "FieldA", TypeName: "TypeA"}

	withNodeDiff := withNode
	withNodeDiff.DifferenceLocation.Node = &diff.Node{Field: "FieldA", TypeName: "TypeB"}

	require.TrueT(t, withNode.Matches(withNodeSame))
	require.FalseT(t, withNode.Matches(withNodeDiff))

	withNodeDiff.DifferenceLocation.Node = &diff.Node{Field: "FieldB", TypeName: "TypeA"}

	require.TrueT(t, withNode.Matches(withNodeSame))
	require.FalseT(t, withNode.Matches(withNodeDiff))

	withNestedNode := withNode
	withNestedNode.DifferenceLocation = withNestedNode.DifferenceLocation.AddNode(&diff.Node{Field: "ChildA", TypeName: "ChildA"})
	withNestedNodeSame := withNode
	withNestedNodeSame.DifferenceLocation = withNestedNodeSame.DifferenceLocation.AddNode(&diff.Node{Field: "ChildA", TypeName: "ChildA"})
	withNestedNodeDiff := withNode
	withNestedNodeDiff.DifferenceLocation = withNestedNodeDiff.DifferenceLocation.AddNode(&diff.Node{Field: "ChildB", TypeName: "ChildA"})

	require.TrueT(t, withNestedNode.Matches(withNestedNodeSame))
	require.FalseT(t, withNestedNode.Matches(withNestedNodeDiff))
}
