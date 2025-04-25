package generator

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"path"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/ast/astutil"
)

func formatGo(filename string, content []byte) ([]byte, error) {
	fset, file, err := parseGo(filename, content)
	if err != nil {
		return nil, err
	}

	// WIP: findSelectedX is not perfect to detect used imports. Detected symbols might not be package, but shadowed local variables.
	used := findSelectedX(file)
	tweakImports(fset, file, used)

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

func findSelectedX(file *ast.File) map[string]bool {
	used := make(map[string]bool)
	ast.Inspect(file, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.SelectorExpr:
			if id, ok := n.X.(*ast.Ident); ok {
				used[id.Name] = true
			}
		}
		return true
	})
	return used
}

func tweakImports(fset *token.FileSet, file *ast.File, used map[string]bool) {
	// imported := map[string]bool{}
	shouldRemove := []*ast.ImportSpec{}
	for _, impt := range file.Imports {
		name := importSpecToAssumedName(impt)

		// WIP: duplicated imports are not deduped
		// if name != "_" && imported[name] {
		// 	deleteImportSpec(fset, file, impt)
		// 	continue
		// }
		// imported[name] = true

		if ok := used[name]; ok {
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

func importSpecToAssumedName(importSpec *ast.ImportSpec) string {
	if importSpec.Name != nil {
		return importSpec.Name.Name
	}
	importPath, _ := strconv.Unquote(importSpec.Path.Value)
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
func notIdentifier(ch rune) bool {
	return !('a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' ||
		'0' <= ch && ch <= '9' ||
		ch == '_' ||
		ch >= utf8.RuneSelf && (unicode.IsLetter(ch) || unicode.IsDigit(ch)))
}
