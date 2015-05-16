package generate

import (
	"fmt"
	goparser "go/parser"
	"go/token"
	"path/filepath"
	"sort"

	"golang.org/x/tools/go/loader"

	"github.com/jessevdk/go-flags"
)

// SpecFile command to generate a swagger spec from a go application
type SpecFile struct {
	BasePath string         `long:"base-path" short:"b" description:"the base path to use" default:"."`
	Output   flags.Filename `long:"output" short:"o" description:"the file to write to" default:"./swagger.json"`
}

// Execute runs this command
func (s *SpecFile) Execute(args []string) error {
	//docFile := "/home/ivan/go/src/github.com/casualjim/go-swagger/internal/testing/petstoreapp/doc.go"
	//fileSet := token.NewFileSet()
	//fileTree, err := goparser.ParseFile(fileSet, docFile, nil, goparser.ParseComments)
	//if err != nil {
	//return err
	//}
	//pretty.Println(fileTree)

	//for _, comment := range fileTree.Comments {
	//for i, commentLine := range strings.Split(comment.Text(), "\n") {
	//fmt.Println("comment", i+1)
	//fmt.Println(commentLine)
	//fmt.Println(len(commentLine), "$#$")
	//}
	//}

	var conf loader.Config
	conf.Import("./cmd/swagger-petstore-server")
	conf.ParserMode = goparser.ParseComments
	prog, err := conf.Load()
	if err != nil {
		return err
	}

	//printProgram(prog)
	printFile(prog.Fset, prog.Package("github.com/casualjim/go-swagger/cmd/swagger/models"))

	return nil
}
func printFile(fset *token.FileSet, info *loader.PackageInfo) {
	var names []string
	for _, f := range info.Files {
		names = append(names, filepath.Base(fset.File(f.Pos()).Name()))
	}
	fmt.Printf("%s.Files: %s\n", info.Pkg.Path(), names)
}
func printProgram(prog *loader.Program) {
	// Created packages are the initial packages specified by a call
	// to CreateFromFilenames or CreateFromFiles.
	var names []string
	for _, info := range prog.Created {
		names = append(names, info.Pkg.Path())
	}
	fmt.Printf("created: %s\n", names)

	// Imported packages are the initial packages specified by a
	// call to Import or ImportWithTests.
	names = nil
	for _, info := range prog.Imported {
		names = append(names, info.Pkg.Path())
	}
	sort.Strings(names)
	fmt.Printf("imported: %s\n", names)

	// InitialPackages contains the union of created and imported.
	names = nil
	for _, info := range prog.InitialPackages() {
		names = append(names, info.Pkg.Path())
	}
	sort.Strings(names)
	fmt.Printf("initial: %s\n", names)

	// AllPackages contains all initial packages and their dependencies.
	names = nil
	for pkg := range prog.AllPackages {
		names = append(names, pkg.Path())
	}
	sort.Strings(names)
	fmt.Printf("all: %s\n", names)
}
