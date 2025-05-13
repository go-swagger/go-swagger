package generator

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"slices"
	"strconv"

	"golang.org/x/tools/go/ast/astutil"
)

func formatGo(filename string, content []byte) ([]byte, error) {
	fset, file, err := parseGo(filename, content)
	if err != nil {
		// If we can't parse file, we give up formatting
		return content, nil
	}

	mergeImports(file)
	removeUnusedImports(fset, file)
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

func removeUnusedImports(fset *token.FileSet, file *ast.File) {
	shouldRemove := []*ast.ImportSpec{}
	for _, impt := range file.Imports {
		path, _ := strconv.Unquote(impt.Path.Value)

		if astutil.UsesImport(file, path) {
			continue
		}
		shouldRemove = append(shouldRemove, impt)
	}
	for _, impt := range shouldRemove {
		deleteImportSpec(fset, file, impt)
	}
}

func deleteImportSpec(fset *token.FileSet, file *ast.File, spec *ast.ImportSpec) {
	importPath, _ := strconv.Unquote(spec.Path.Value)
	if spec.Name != nil {
		astutil.DeleteNamedImport(fset, file, spec.Name.Name, importPath)
	} else {
		astutil.DeleteImport(fset, file, importPath)
	}
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
