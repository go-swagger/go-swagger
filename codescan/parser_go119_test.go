// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package codescan

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
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

	assert.Equal(t, []string{"This has a title that starts with a hash tag"}, st.Title())
	assert.Equal(t, []string{"The punctuation here does indeed matter. But it won't for go."}, st.Description())

	st = &sectionedParser{}
	st.setTitle = func(_ []string) {}
	err = st.Parse(ascg(text2))
	require.NoError(t, err)

	assert.Equal(t, []string{"This has a title without whitespace."}, st.Title())
	assert.Equal(t, []string{"The punctuation here does indeed matter. But it won't for go.", "", "# There is an inline header here that doesn't count for finding a title"}, st.Description())
}
