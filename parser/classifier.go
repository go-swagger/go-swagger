package parser

import (
	"fmt"
	"go/ast"
	"regexp"

	"golang.org/x/tools/go/loader"
)

type swaggerKind int8

const (
	swMeta swaggerKind = iota
	swModel
	swRoute
)

var (
	rxSwaggerAnnotation = regexp.MustCompile("[^+]*\\+[^\\w]*swagger:(\\w+)")
)

type packageFilter struct {
	Name string
}

func (pf *packageFilter) Matches(path string) bool {
	return path == pf.Name
}

type packageFilters struct {
	Operations []packageFilter
	Models     []packageFilter
	Meta       *packageFilter
}

func (pf *packageFilters) HasFilters() bool {
	return len(pf.Operations) > 0 || len(pf.Models) > 0 || pf.Meta != nil
}

func (pf *packageFilters) Matches(path string) bool {
	if pf.Meta != nil && pf.Meta.Matches(path) {
		return true
	}

	for _, mod := range pf.Models {
		if mod.Matches(path) {
			return true
		}
	}

	for _, op := range pf.Operations {
		if op.Matches(path) {
			return true
		}
	}
	return false
}

type classifiedProgram struct {
	prog       *loader.Program
	Meta       []*ast.File
	Models     []*ast.File
	Operations []*ast.File
}

// programClassifier classifies the files of a program into buckets
// for processing by a swagger spec generator. This buckets files in
// 3 groups: Meta, Models and Operations.
//
// Each of these buckets is then processed with an appropriate parsing strategy
//
// When there are Include or Exclude filters provide they are used to limit the
// candidates prior to parsing.
// The include filters take precedence over the excludes. So when something appears
// in both filters it will be included.
type programClassifier struct {
	Includes packageFilters
	Excludes packageFilters
}

func (pc *programClassifier) Classify(prog *loader.Program) (*classifiedProgram, error) {
	var cp classifiedProgram
	for pkg, pkgInfo := range prog.AllPackages {
		if pc.Includes.HasFilters() {
			if !pc.Includes.Matches(pkg.Path()) {
				continue
			}
		} else if pc.Excludes.HasFilters() {
			if pc.Excludes.Matches(pkg.Path()) {
				continue
			}
		}

		for _, file := range pkgInfo.Files {
			var op, md, mt bool // only add a particular file once
			for _, comments := range file.Comments {
				matches := rxSwaggerAnnotation.FindStringSubmatch(comments.Text())
				if len(matches) > 1 {
					switch matches[1] {
					case "route":
						if !op {
							cp.Operations = append(cp.Operations, file)
							op = true
						}
					case "model":
						if !md {
							cp.Models = append(cp.Models, file)
							md = true
						}
					case "meta":
						if !mt {
							cp.Meta = append(cp.Meta, file)
							mt = true
						}
					case "strfmt":
						// TODO: perhaps collect these and pass along to avoid lookups later on
					default:
						return nil, fmt.Errorf("classifier: unknown swagger annotation %q", matches[1])
					}
				}
			}
		}
	}

	return &cp, nil
}
