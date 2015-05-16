package parser

import (
	"fmt"
	"go/ast"
	"testing"

	"github.com/casualjim/go-swagger/spec"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/types"
)

func TestSchemaParser(t *testing.T) {

	//docFile := "../fixtures/goparsing/petstoreapp/model.go"
	//fileSet := token.NewFileSet()
	var loader loader.Config
	//loader.ParserMode = goparser.ParseComments
	//fileTree, err := loader.ParseFile(docFile, nil)
	//if err != nil {
	//t.FailNow()
	//}
	//pretty.Println(fileTree)

	var orig = spec.License{Name: "BSD", URL: "http://somewhere.com"}

	loader.Import("../fixtures/goparsing/petstoreapp")
	loader.Import("../spec")
	pkg, err := loader.Load()
	if err != nil {
		t.FailNow()
	}

	for _, info := range pkg.InitialPackages() {
		//fmt.Println("package", info.Pkg.Path())

		for k, v := range info.Defs {
			//fmt.Println("name:", k.Name, "exported:", k.IsExported(), "pos:", k.Pos())
			if v != nil && k.Name == "License" {
				fmt.Println(k)
				fmt.Println("kind:", k.Obj.Kind)
				fmt.Printf("%#v\n", k.Obj)
				fmt.Printf("%#v\n", k.Obj.Decl)
				fmt.Printf("%T\n", k.Obj.Data)
				if decl, ok := k.Obj.Decl.(*ast.ValueSpec); ok && len(decl.Values) > 0 {
					fmt.Printf("%#v\n", decl.Values[0])
					lit := decl.Values[0].(*ast.CompositeLit)
					for _, ex := range lit.Elts {

						fmt.Printf("%#v\n", ex)
					}
					tv, err := types.EvalNode(pkg.Fset, lit, nil, nil)
					if err != nil {
						t.Fatal(err)
					}
					fmt.Printf("%#v\n", tv)
				}

				//fmt.Print(":")
				//fmt.Printf("%#v\n", v)
			}
		}
	}
	//for _, f := range info.Files {
	//fmt.Println("  file:", filepath.Base(pkg.Fset.File(f.Pos()).Name()))
	//for _, decl := range f.Decls {
	//switch gd := decl.(type) {
	//case *ast.GenDecl:
	//if len(gd.Specs) > 0 {
	//if ts, ok := gd.Specs[0].(*ast.TypeSpec); ok {
	////fmt.Println(gd.Doc.Text())
	//fmt.Println("  struct", ts.Name.Name)
	//}
	//}
	//default:
	//fmt.Println("unhandled decl:", gd)
	//}
	//}
	//}
	//}

	fmt.Println("here we're done")

}
