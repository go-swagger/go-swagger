package scan

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/loader"
)

type packageFilter struct {
	Name string
}

func (pf *packageFilter) Matches(path string) bool {
	return path == pf.Name
}

type packageFilters []packageFilter

func (pf packageFilters) HasFilters() bool {
	return len(pf) > 0
}

func (pf packageFilters) Matches(path string) bool {
	for _, mod := range pf {
		if mod.Matches(path) {
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
	Parameters []*ast.File
	Responses  []*ast.File
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
			var op, mt, pm, rs bool // only add a particular file once
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
						// models are discovered through parameters and responses
						// no actual scanning for them is required
					case "meta":
						if !mt {
							cp.Meta = append(cp.Meta, file)
							mt = true
						}
					case "parameters":
						if !pm {
							cp.Parameters = append(cp.Parameters, file)
							pm = true
						}
					case "response":
						if !rs {
							cp.Responses = append(cp.Responses, file)
							rs = true
						}
					case "strfmt":
						// TODO: perhaps collect these and pass along to avoid lookups later on
					case "allOf":
					default:
						return nil, fmt.Errorf("classifier: unknown swagger annotation %q", matches[1])
					}
				}
			}
		}
	}

	return &cp, nil
}
