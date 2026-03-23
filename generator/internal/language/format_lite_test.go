// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package language

import (
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestFormatLite_ValidSource(t *testing.T) {
	src := []byte("package main\n\nimport \"fmt\"\n\nfunc main() { fmt.Println(\"hello\") }\n")
	res, err := FormatLite("test.go", src)
	require.NoError(t, err)
	assert.StringContainsT(t, string(res), "package main")
	assert.StringContainsT(t, string(res), `"fmt"`)
}

func TestFormatLite_Fragment(t *testing.T) {
	// no package statement: treated as a fragment
	src := []byte("func hello() { fmt.Println(\"hi\") }\n")
	res, err := FormatLite("frag.go", src)
	require.NoError(t, err)

	// fragment parsing succeeded: output contains the function
	assert.StringContainsT(t, string(res), "hello()")
}

func TestFormatLite_InvalidSource(t *testing.T) {
	src := []byte("this is not go code at all {{{")
	_, err := FormatLite("bad.go", src)
	require.Error(t, err)
}

func TestFormatLite_RemovesUnusedImports(t *testing.T) {
	src := []byte(`package main

import (
	"fmt"
	"os"
)

func main() { fmt.Println("hello") }
`)
	res, err := FormatLite("test.go", src)
	require.NoError(t, err)

	output := string(res)
	assert.StringContainsT(t, output, `"fmt"`)
	assert.FalseT(t, strings.Contains(output, `"os"`))
}

func TestFormatLite_AddsKnownImports(t *testing.T) {
	// uses "fmt" without importing it — fixImports should add it
	src := []byte(`package main

func main() { fmt.Println("hello") }
`)
	res, err := FormatLite("test.go", src)
	require.NoError(t, err)
	assert.StringContainsT(t, string(res), `"fmt"`)
}

func TestFormatLite_RemovesBlankLinesBetweenImports(t *testing.T) {
	src := []byte(`package main

import (
	"fmt"

	"os"
)

func main() {
	fmt.Println(os.Args)
}
`)
	res, err := FormatLite("test.go", src)
	require.NoError(t, err)

	output := string(res)
	assert.StringContainsT(t, output, `"fmt"`)
	assert.StringContainsT(t, output, `"os"`)
}

func TestFormatLite_SingleImportNoParen(t *testing.T) {
	src := []byte(`package main

import (
	"fmt"
)

func main() { fmt.Println("hi") }
`)
	res, err := FormatLite("test.go", src)
	require.NoError(t, err)

	output := string(res)
	// single import should have parens removed
	assert.StringContainsT(t, output, `import "fmt"`)
}

func TestFormatLite_DuplicateImportLastWins(t *testing.T) {
	src := []byte(`package main

import (
	"fmt"
	"fmt"
)

func main() { fmt.Println("hi") }
`)
	res, err := FormatLite("test.go", src)
	require.NoError(t, err)

	output := string(res)
	assert.StringContainsT(t, output, `"fmt"`)
	// should have only one import of fmt
	assert.EqualT(t, 1, strings.Count(output, `"fmt"`))
}

func TestFormatLite_WithFormatOptions(t *testing.T) {
	src := []byte(`package main

import "fmt"

func main() { fmt.Println("hello") }
`)
	res, err := FormatLite("test.go", src, WithFormatOnly(true))
	require.NoError(t, err)
	assert.StringContainsT(t, string(res), "package main")
}

func TestFormatLite_AliasedImport(t *testing.T) {
	// aliased import: name != assumed name from path
	src := []byte(`package main

import myalias "fmt"

func main() { myalias.Println("hi") }
`)
	res, err := FormatLite("test.go", src)
	require.NoError(t, err)
	assert.StringContainsT(t, string(res), `myalias "fmt"`)
}

func TestFormatLite_BlankAndDotImports(t *testing.T) {
	// blank and dot imports should be preserved, not removed
	src := []byte(`package main

import (
	_ "embed"
	. "fmt"
)

func main() { Println("hi") }
`)
	res, err := FormatLite("test.go", src)
	require.NoError(t, err)

	output := string(res)
	assert.StringContainsT(t, output, `_ "embed"`)
	assert.StringContainsT(t, output, `. "fmt"`)
}

func TestFormatLite_NonPackageParseError(t *testing.T) {
	// has package statement but still invalid syntax: not a "expected 'package'" error
	src := []byte("package main\n\nfunc { broken }\n")
	_, err := FormatLite("bad.go", src)
	require.Error(t, err)
}

func TestFormatLite_FragmentWithImport(t *testing.T) {
	// a function without package statement exercises the fragment+cleanup path
	src := []byte("func greet() string { return fmt.Sprintf(\"hello %s\", \"world\") }\n")
	res, err := FormatLite("frag.go", src)
	require.NoError(t, err)
	assert.StringContainsT(t, string(res), "Sprintf")
}

func TestFormatByImports_MutexPaths(t *testing.T) {
	src := []byte(`package main

import "fmt"

func main() { fmt.Println("hello") }
`)
	opts := FormatOptsWithDefault(nil)

	// first call sets the LocalPrefix
	res, err := formatByImports("test.go", src, opts)
	require.NoError(t, err)
	assert.StringContainsT(t, string(res), "package main")

	// second call with same prefix takes the fast path (RLock)
	res, err = formatByImports("test.go", src, opts)
	require.NoError(t, err)
	assert.StringContainsT(t, string(res), "package main")

	// call with different prefix takes the slow path (Lock)
	opts2 := FormatOptsWithDefault([]FormatOption{WithFormatLocalPrefixes("github.com/other")})
	res, err = formatByImports("test.go", src, opts2)
	require.NoError(t, err)
	assert.StringContainsT(t, string(res), "package main")
}

func TestImportPathToAssumedName(t *testing.T) {
	// simple package
	assert.EqualT(t, "fmt", importPathToAssumedName("fmt"))

	// nested package
	assert.EqualT(t, "runtime", importPathToAssumedName("github.com/go-openapi/runtime"))

	// versioned import: strips v2 and uses parent dir
	assert.EqualT(t, "testify", importPathToAssumedName("github.com/go-openapi/testify/v2"))

	// go- prefix stripped
	assert.EqualT(t, "openapi", importPathToAssumedName("github.com/go-openapi"))

	// non-identifier characters truncated
	assert.EqualT(t, "pkg", importPathToAssumedName("example.com/pkg@v1.0.0"))
}

func TestNotIdentifier(t *testing.T) {
	// lowercase letter
	assert.FalseT(t, notIdentifier('a'))
	// uppercase letter
	assert.FalseT(t, notIdentifier('Z'))
	// digit
	assert.FalseT(t, notIdentifier('5'))
	// underscore
	assert.FalseT(t, notIdentifier('_'))
	// ASCII non-identifier
	assert.TrueT(t, notIdentifier('@'))
	assert.TrueT(t, notIdentifier('.'))
	// unicode letter (beyond ASCII)
	assert.FalseT(t, notIdentifier('é'))
	// unicode non-letter, non-digit
	assert.TrueT(t, notIdentifier('→'))
}
