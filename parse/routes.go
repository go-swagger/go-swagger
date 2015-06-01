package parse

import (
	"go/ast"
	"regexp"
	"strconv"
	"strings"

	"github.com/casualjim/go-swagger/spec"

	"golang.org/x/tools/go/loader"
)

type operationSetter func(*spec.Operation, []string) error

func newOperationSummary(setter operationSetter) (t *sectionTagger) {
	t = newTitleTagger()
	t.Name = "Summary"
	t.rxStripComments = rxStripComments
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Operation), lines) }
	return
}

func newOperationDescription(setter operationSetter) (t *sectionTagger) {
	t = newDescriptionTagger()
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Operation), lines) }
	return
}

func newOperationSection(name string, multiLine bool, setter operationSetter) (t *sectionTagger) {
	t = newSectionTagger(name, multiLine)
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Operation), lines) }
	return
}

func newOperationFieldSection(name string, multiLine bool, matcher *regexp.Regexp, setter operationSetter) (t *sectionTagger) {
	t = newSectionTagger(name, multiLine)
	t.matcher = matcher
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Operation), lines) }
	return
}

func setOperationSummary(op *spec.Operation, lines []string) error {
	op.Summary = joinDropLast(lines)
	return nil
}

func setOperationDescription(op *spec.Operation, lines []string) error {
	op.Description = joinDropLast(lines)
	return nil
}

func setOperationConsumes(op *spec.Operation, lines []string) error {
	op.Consumes = removeEmptyLines(lines)
	return nil
}

func setOperationProduces(op *spec.Operation, lines []string) error {
	op.Produces = removeEmptyLines(lines)
	return nil
}

func setOperationSchemes(op *spec.Operation, lines []string) error {
	lns := lines
	if len(lns) == 0 || lns[0] == "" {
		return nil
	}
	sch := strings.Split(lns[0], ", ")
	var schemes []string
	for _, s := range sch {
		schemes = append(schemes, strings.TrimSpace(s))
	}
	op.Schemes = schemes
	return nil
}

func setOperationSecurity(op *spec.Operation, lines []string) error {
	if len(lines) == 0 {
		return nil
	}

	for _, line := range lines {
		kv := strings.SplitN(line, ":", 2)
		var scopes []string
		var key string

		if len(kv) > 1 {
			scs := strings.Split(rxNotAlNumSpaceComma.ReplaceAllString(kv[1], ""), ",")
			for _, scope := range scs {
				scopes = append(scopes, strings.TrimSpace(scope))
			}

			key = strings.TrimSpace(kv[0])

			op.Security = append(op.Security, map[string][]string{key: scopes})
		}
	}
	return nil
}

func setOperationResponse(op *spec.Operation, lines []string) error {
	if len(lines) == 0 {
		return nil
	}

	for _, line := range lines {
		kv := strings.SplitN(line, ":", 2)
		var key, value string

		if len(kv) > 1 {
			key = strings.TrimSpace(kv[0])
			value = strings.TrimSpace(kv[1])

			if op.Responses == nil {
				op.Responses = new(spec.Responses)
			}
			resps := op.Responses
			var resp spec.Response
			ref, err := spec.NewRef("#/responses/" + value)
			if err != nil {
				return err
			}
			resp.Ref = ref
			if strings.EqualFold("default", key) {
				if resps.Default == nil {
					resps.Default = &resp
				}
			} else {
				if sc, err := strconv.Atoi(key); err == nil {
					if resps.StatusCodeResponses == nil {
						resps.StatusCodeResponses = make(map[int]spec.Response)
					}
					resps.StatusCodeResponses[sc] = resp
				}
			}
		}
	}
	return nil
}

func newRoutesParser(prog *loader.Program) *routesParser {
	return &routesParser{
		program: prog,
		//taggers: []*sectionTagger{
		//newOperationSummary(setOperationSummary),
		//newOperationDescription(setOperationDescription),
		//newOperationFieldSection("Consumes", true, regexp.MustCompile("[Cc]onsumes\\p{Zs}*:"), setOperationConsumes),
		//newOperationFieldSection("Produces", true, regexp.MustCompile("[Pp]roduces\\p{Zs}*:"), setOperationProduces),
		//newOperationFieldSection("Schemes", false, regexp.MustCompile("[Pp]roduces\\p{Zs}*:\\p{Zs}*((?:(?:https?|wss?)\\p{Zs}*,?\\p{Zs}*)+)$"), setOperationSchemes),
		//newOperationFieldSection("Security", true, regexp.MustCompile("[Ss]ecurity\\p{Zs}*:"), setOperationSecurity),
		//},
	}
}

type routesParser struct {
	program *loader.Program
}

func (rp *routesParser) Parse(gofile *ast.File, target interface{}) error {
	tgt := target.(*spec.Paths)
	for _, comsec := range gofile.Comments {

		// check if this is a route comment section
		var method, path, id string
		var tags []string
		var remaining *ast.CommentGroup
		var justMatched bool

		for _, cmt := range comsec.List {
			for _, line := range strings.Split(cmt.Text, "\n") {
				matches := rxRoute.FindStringSubmatch(line)
				if len(matches) > 3 && len(matches[3]) > 0 {
					method, path, id = matches[1], matches[2], matches[len(matches)-1]
					tags = rxSpace.Split(matches[3], -1)
					justMatched = true
				} else if method != "" {
					if remaining == nil {
						remaining = new(ast.CommentGroup)
					}
					if !justMatched || strings.TrimSpace(rxStripComments.ReplaceAllString(line, "")) != "" {
						cc := new(ast.Comment)
						cc.Slash = cmt.Slash
						cc.Text = line
						remaining.List = append(remaining.List, cc)
						justMatched = false
					}
				}
			}
		}

		if method == "" {
			continue // it's not, next!
		}

		pthObj := tgt.Paths[path]
		op := new(spec.Operation)
		op.ID = id
		switch strings.ToUpper(method) {
		case "GET":
			if pthObj.Get != nil {
				if id == pthObj.Get.ID {
					op = pthObj.Get
				} else {
					pthObj.Get = op
				}
			} else {
				pthObj.Get = op
			}

		case "POST":
			if pthObj.Post != nil {
				if id == pthObj.Post.ID {
					op = pthObj.Post
				} else {
					pthObj.Post = op
				}
			} else {
				pthObj.Post = op
			}

		case "PUT":
			if pthObj.Put != nil {
				if id == pthObj.Put.ID {
					op = pthObj.Put
				} else {
					pthObj.Put = op
				}
			} else {
				pthObj.Put = op
			}

		case "PATCH":
			if pthObj.Patch != nil {
				if id == pthObj.Patch.ID {
					op = pthObj.Patch
				} else {
					pthObj.Patch = op
				}
			} else {
				pthObj.Patch = op
			}

		case "HEAD":
			if pthObj.Head != nil {
				if id == pthObj.Head.ID {
					op = pthObj.Head
				} else {
					pthObj.Head = op
				}
			} else {
				pthObj.Head = op
			}

		case "DELETE":
			if pthObj.Delete != nil {
				if id == pthObj.Delete.ID {
					op = pthObj.Delete
				} else {
					pthObj.Delete = op
				}
			} else {
				pthObj.Delete = op
			}

		case "OPTIONS":
			if pthObj.Options != nil {
				if id == pthObj.Options.ID {
					op = pthObj.Options
				} else {
					pthObj.Options = op
				}
			} else {
				pthObj.Options = op
			}
		}
		op.Tags = tags

		taggers := []*sectionTagger{
			newOperationSummary(setOperationSummary),
			newOperationDescription(setOperationDescription),
			newOperationSection("Consumes", true, setOperationConsumes),
			newOperationSection("Produces", true, setOperationProduces),
			newOperationSection("Schemes", false, setOperationSchemes),
			newOperationSection("Security", true, setOperationSecurity),
			newOperationSection("Responses", true, setOperationResponse),
		}
		if err := parseDocComments(remaining, op, taggers, nil); err != nil {
			return err
		}

		if tgt.Paths == nil {
			tgt.Paths = make(map[string]spec.PathItem)
		}
		tgt.Paths[path] = pthObj
	}

	return nil
}
