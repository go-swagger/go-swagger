package diff

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArrayDiff(t *testing.T) {
	listA := []string{"abc", "def", "ghi", "jkl"}
	added, deleted, common := fromStringArray(listA).DiffsTo(listA)
	require.Equal(t, []string{}, added)
	require.Equal(t, []string{}, deleted)
	require.ElementsMatch(t, listA, common)

	listB := []string{"abc", "ghi", "jkl", "xyz", "fgh"}
	added, deleted, common = fromStringArray(listA).DiffsTo(listB)
	require.ElementsMatch(t, []string{"xyz", "fgh"}, added)
	require.ElementsMatch(t, []string{"def"}, deleted)
	require.ElementsMatch(t, []string{"abc", "ghi", "jkl"}, common)
}

func TestMapDiff(t *testing.T) {
	mapA := map[string]interface{}{"abc": 1, "def": 2, "ghi": 3, "jkl": 4}
	added, deleted, common := fromStringMap(mapA).DiffsTo(mapA)
	require.Equal(t, map[string]interface{}{}, added)
	require.Equal(t, map[string]interface{}{}, deleted)

	commonDiffs := map[string]interface{}{"abc": Pair{1, 1}, "def": Pair{2, 2}, "ghi": Pair{3, 3}, "jkl": Pair{4, 4}}
	require.Equal(t, commonDiffs, common)

	mapB := map[string]interface{}{"abc": 2, "ghi": 3, "jkl": 4, "xyz": 5, "fgh": 6}
	added, deleted, common = fromStringMap(mapA).DiffsTo(mapB)
	require.Equal(t, map[string]interface{}{"xyz": 5, "fgh": 6}, added)
	require.Equal(t, map[string]interface{}{"def": 2}, deleted)
	commonDiffs = map[string]interface{}{"abc": Pair{1, 2}, "ghi": Pair{3, 3}, "jkl": Pair{4, 4}}
	require.Equal(t, commonDiffs, common)
}
