package generator

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"path"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func formatGo(filename string, content []byte) ([]byte, error) {
	fset, file, err := parseGo(filename, content)
	if err != nil {
		// If we can't parse file, we give up formatting
		return content, nil
	}

	mergeImports(file)
	cleanImports(fset, file)
	removeUnecessaryImportParens(fset, file)

	printConfig := &printer.Config{
		Mode:     printer.UseSpaces | printer.TabIndent,
		Tabwidth: 2,
	}
	var buf bytes.Buffer
	if err := printConfig.Fprint(&buf, fset, file); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func parseGo(ffn string, content []byte) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	mode := parser.ParseComments | parser.AllErrors
	file, err := parser.ParseFile(fset, ffn, content, mode)
	if err != nil {
		return nil, nil, err
	}
	return fset, file, nil
}

func cleanImports(fset *token.FileSet, file *ast.File) {
	seen := make(map[string]bool)
	shouldRemove := []*ast.ImportSpec{}
	usedNames := collectTopNames(file)
	for _, impt := range file.Imports {
		name := importPathToAssumedName(importPath(impt))
		if impt.Name != nil {
			name = impt.Name.String()
		}

		if seen[name] {
			shouldRemove = append(shouldRemove, impt)
			continue
		}
		seen[name] = true

		// astutil.UsesImport is not precise enough for our needs: https://github.com/golang/go/issues/30331#issuecomment-466174437
		if usedNames[name] {
			continue
		}
		if name == "_" || name == "." {
			continue
		}
		shouldRemove = append(shouldRemove, impt)
	}
	for _, impt := range shouldRemove {
		deleteImportSpec(fset, file, impt)
	}
}

func deleteImportSpec(fset *token.FileSet, file *ast.File, spec *ast.ImportSpec) {
	// remove from file.Imports
	i := slices.IndexFunc(file.Imports, func(i *ast.ImportSpec) bool {
		return i == spec
	})
	if i >= 0 {
		file.Imports = slices.Delete(file.Imports, i, i+1)
	}

	// remove from file.Decls
	if len(file.Decls) == 0 {
		return
	}
	gen, ok := file.Decls[0].(*ast.GenDecl)
	if !ok {
		return
	}
	i = slices.IndexFunc(gen.Specs, func(i ast.Spec) bool {
		return i == spec
	})
	if i < 0 {
		return
	}
	gen.Specs = slices.Delete(gen.Specs, i, i+1)
}

func removeUnecessaryImportParens(fset *token.FileSet, file *ast.File) {
	if len(file.Imports) == 1 {
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok {
				break
			}
			if gen.Tok != token.IMPORT {
				break
			}
			gen.Lparen = token.NoPos
			gen.Rparen = token.NoPos
		}
	}
}

// mergeImports merges all the import declarations into the first one.
// Taken from [golang.org/x/tools/ast/astutil](https://cs.opensource.google/go/x/tools/+/refs/tags/v0.32.0:go/ast/astutil/imports.go;l=170).
// This does not adjust line numbers properly
func mergeImports(f *ast.File) {
	var first *ast.GenDecl
	for i := 0; i < len(f.Decls); i++ {
		decl := f.Decls[i]
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.IMPORT || declImports(gen, "C") {
			continue
		}
		if first == nil {
			first = gen
			continue // Don't touch the first one.
		}
		// We now know there is more than one package in this import
		// declaration. Ensure that it ends up parenthesized.
		first.Lparen = first.Pos()
		// Move the imports of the other import declaration to the first one.
		for _, spec := range gen.Specs {
			spec.(*ast.ImportSpec).Path.ValuePos = first.Pos()
			first.Specs = append(first.Specs, spec)
		}
		f.Decls = slices.Delete(f.Decls, i, i+1)
		i--
	}
}

// declImports reports whether gen contains an import of path.
// Taken from [golang.org/x/tools/ast/astutil](https://cs.opensource.google/go/x/tools/+/refs/tags/v0.32.0:go/ast/astutil/imports.go;l=433).
func declImports(gen *ast.GenDecl, path string) bool {
	if gen.Tok != token.IMPORT {
		return false
	}
	for _, spec := range gen.Specs {
		impspec := spec.(*ast.ImportSpec)
		if importPath(impspec) == path {
			return true
		}
	}
	return false
}

// importPath returns the unquoted import path of s,
// or "" if the path is not properly quoted.
// Taken from [golang.org/x/tools/ast/astutil](https://cs.opensource.google/go/x/tools/+/refs/tags/v0.32.0:go/ast/astutil/imports.go;l=424).
func importPath(s *ast.ImportSpec) string {
	t, err := strconv.Unquote(s.Path.Value)
	if err != nil {
		return ""
	}
	return t
}

func collectTopNames(n ast.Node) map[string]bool {
	names := make(map[string]bool)
	ast.Walk(visitFn(func(n ast.Node) {
		s, ok := n.(*ast.SelectorExpr)
		if !ok {
			return
		}
		id, ok := s.X.(*ast.Ident)
		if !ok {
			return
		}
		if id.Obj != nil {
			return
		}
		names[id.Name] = true
	}), n)
	return names
}

type visitFn func(node ast.Node)

func (fn visitFn) Visit(node ast.Node) ast.Visitor {
	fn(node)
	return fn
}

// importPathToAssumedName returns the assumed package name of an import path.
// it is taken from [tools/internal/imports/fix.go](https://github.com/golang/tools/blob/v0.33.0/internal/imports/fix.go#L1233)
func importPathToAssumedName(importPath string) string {
	base := path.Base(importPath)
	if strings.HasPrefix(base, "v") {
		if _, err := strconv.Atoi(base[1:]); err == nil {
			dir := path.Dir(importPath)
			if dir != "." {
				base = path.Base(dir)
			}
		}
	}
	base = strings.TrimPrefix(base, "go-")
	if i := strings.IndexFunc(base, notIdentifier); i >= 0 {
		base = base[:i]
	}
	return base
}

// notIdentifier reports whether ch is an invalid identifier character.
// it is taken from [tools/internal/imports/fix.go](https://github.com/golang/tools/blob/v0.33.0/internal/imports/fix.go#L1233)
func notIdentifier(ch rune) bool {
	return !('a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' ||
		'0' <= ch && ch <= '9' ||
		ch == '_' ||
		ch >= utf8.RuneSelf && (unicode.IsLetter(ch) || unicode.IsDigit(ch)))
}
