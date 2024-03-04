package codescan

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSectionedParser_TitleDescriptionGo119(t *testing.T) {
	text := `# This has a title that starts with a hash tag

The punctuation here does indeed matter. But it won't for go.
`

	text2 := `This has a title without whitespace.

The punctuation here does indeed matter. But it won't for go.

# There is an inline header here that doesn't count for finding a title

`

	var err error

	st := &sectionedParser{}
	st.setTitle = func(_ []string) {}
	err = st.Parse(ascg(text))
	require.NoError(t, err)

	assert.EqualValues(t, []string{"This has a title that starts with a hash tag"}, st.Title())
	assert.EqualValues(t, []string{"The punctuation here does indeed matter. But it won't for go."}, st.Description())

	st = &sectionedParser{}
	st.setTitle = func(_ []string) {}
	err = st.Parse(ascg(text2))
	require.NoError(t, err)

	assert.EqualValues(t, []string{"This has a title without whitespace."}, st.Title())
	assert.EqualValues(t, []string{"The punctuation here does indeed matter. But it won't for go.", "", "# There is an inline header here that doesn't count for finding a title"}, st.Description())
}
