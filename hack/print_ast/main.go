package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

var fileName string

func init() {
	flag.StringVar(&fileName, "f", "", "input file")
}
func main() {
	flag.Parse()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Failed to parse file %s: %s", fileName, err)
	}
	if err = ast.Print(fset, f); err != nil {
		log.Fatal(err)
	}
}
