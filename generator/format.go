package generator

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strconv"

	"golang.org/x/tools/go/ast/astutil"
)

func formatGo(filename string, content []byte) ([]byte, error) {
	fset, file, err := parseGo(filename, content)
	if err != nil {
		return nil, err
	}

	tweakImports(fset, file)

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

func tweakImports(fset *token.FileSet, file *ast.File) {
	shouldRemove := []*ast.ImportSpec{}
	for _, impt := range file.Imports {
		path, _ := strconv.Unquote(impt.Path.Value)

		// WIP: duplicated imports are not deduped
		// if name != "_" && imported[name] {
		// 	deleteImportSpec(fset, file, impt)
		// 	continue
		// }
		// imported[name] = true

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
