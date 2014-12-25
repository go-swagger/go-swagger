package swagger

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/casualjim/go-swagger"
)

// SpecAnalyzer takes a swagger spec object and turns it into a registry
// with a bunch of utility methods to act on the information in the spec
type SpecAnalyzer struct {
	spec        *swagger.Spec
	consumes    map[string]struct{}
	produces    map[string]struct{}
	authSchemes map[string]struct{}
	operations  map[string]string
}

// NewAnalyzer creates a new spec analyzer instance
func NewAnalyzer(spec *swagger.Spec) *SpecAnalyzer {
	a := &SpecAnalyzer{
		spec:        spec,
		consumes:    make(map[string]struct{}),
		produces:    make(map[string]struct{}),
		authSchemes: make(map[string]struct{}),
		operations:  make(map[string]string),
	}
	a.initialize()
	return a
}

func (s *SpecAnalyzer) initialize() {
	for _, c := range s.spec.Consumes {
		s.consumes[c] = struct{}{}
	}
	for _, c := range s.spec.Produces {
		s.produces[c] = struct{}{}
	}
	for path, pathItem := range s.spec.Paths.Paths {
		s.analyzeOperations(path, &pathItem)
	}
}

func (s *SpecAnalyzer) analyzeOperations(path string, op *swagger.PathItem) {
	s.analyzeOperation(path, op.Get)
	s.analyzeOperation(path, op.Put)
	s.analyzeOperation(path, op.Post)
	s.analyzeOperation(path, op.Patch)
	s.analyzeOperation(path, op.Delete)
	s.analyzeOperation(path, op.Head)
	s.analyzeOperation(path, op.Options)
}

func (s *SpecAnalyzer) analyzeOperation(path string, op *swagger.Operation) {
	if op != nil {
		for _, c := range op.Consumes {
			s.consumes[c] = struct{}{}
		}
		for _, c := range op.Produces {
			s.produces[c] = struct{}{}
		}
		s.operations[op.ID] = path
	}
}

// ValidateRegistrations validates the registrations against the analyzed spec
func (s *SpecAnalyzer) ValidateRegistrations(consumes, produces, schemes, operations []string) error {
	if err := s.verify("consumes", consumes, s.structMapKeys(s.consumes)); err != nil {
		return err
	}
	if err := s.verify("produces", produces, s.structMapKeys(s.produces)); err != nil {
		return err
	}
	// TODO: hook auth in later on
	// if err := s.verify("auth scheme", schemes, s.structMapKeys(s.authSchemes)); err != nil {
	// 	return err
	// }
	if err := s.verify("operation", operations, s.stringMapKeys(s.operations)); err != nil {
		return err
	}
	return nil
}

func (s *SpecAnalyzer) structMapKeys(mp map[string]struct{}) []string {
	var result []string
	for k := range mp {
		result = append(result, k)
	}
	return result
}

func (s *SpecAnalyzer) stringMapKeys(mp map[string]string) []string {
	var result []string
	for k := range mp {
		result = append(result, k)
	}
	return result
}

type verifyError struct {
	Section              string
	MissingSpecification []string
	MissingRegistration  []string
}

func (v *verifyError) Error() string {
	buf := bytes.NewBuffer(nil)

	hasRegMissing := len(v.MissingRegistration) > 0
	hasSpecMissing := len(v.MissingSpecification) > 0

	if hasRegMissing {
		buf.WriteString(fmt.Sprintf("missing [%s] %s registrations", strings.Join(v.MissingRegistration, ", "), v.Section))
	}

	if hasRegMissing && hasSpecMissing {
		buf.WriteString("\n")
	}

	if hasSpecMissing {
		buf.WriteString(fmt.Sprintf("missing from spec file [%s] %s", strings.Join(v.MissingSpecification, ", "), v.Section))
	}

	return buf.String()
}

func (s *SpecAnalyzer) verify(name string, registrations []string, expectations []string) error {
	expected := map[string]struct{}{}
	seen := map[string]struct{}{}

	for _, v := range expectations {
		expected[v] = struct{}{}
	}

	var unspecified []string
	for _, v := range registrations {
		seen[v] = struct{}{}
		if _, ok := expected[v]; !ok {
			unspecified = append(unspecified, v)
		}
	}

	for k := range seen {
		delete(expected, k)
	}

	var unregistered []string
	for k := range expected {
		unregistered = append(unregistered, k)
	}

	if len(unregistered) > 0 || len(unspecified) > 0 {
		return &verifyError{
			Section:              name,
			MissingSpecification: unspecified,
			MissingRegistration:  unregistered,
		}
	}

	return nil
}
