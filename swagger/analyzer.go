package swagger

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/swagger/util"
)

// specAnalyzer takes a swagger spec object and turns it into a registry
// with a bunch of utility methods to act on the information in the spec
type specAnalyzer struct {
	spec        *swagger.Spec
	consumes    map[string]struct{}
	produces    map[string]struct{}
	authSchemes map[string]struct{}
	operations  map[string]string
}

// newAnalyzer creates a new spec analyzer instance
func newAnalyzer(spec *swagger.Spec) *specAnalyzer {
	a := &specAnalyzer{
		spec:        spec,
		consumes:    make(map[string]struct{}),
		produces:    make(map[string]struct{}),
		authSchemes: make(map[string]struct{}),
		operations:  make(map[string]string),
	}
	a.initialize()
	return a
}

func (s *specAnalyzer) initialize() {
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

func (s *specAnalyzer) analyzeOperations(path string, op *swagger.PathItem) {
	s.analyzeOperation(path, op.Get)
	s.analyzeOperation(path, op.Put)
	s.analyzeOperation(path, op.Post)
	s.analyzeOperation(path, op.Patch)
	s.analyzeOperation(path, op.Delete)
	s.analyzeOperation(path, op.Head)
	s.analyzeOperation(path, op.Options)
}

func (s *specAnalyzer) analyzeOperation(path string, op *swagger.Operation) {
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

// AllPaths returns all the paths in the swagger spec
func (s *specAnalyzer) AllPaths() map[string]swagger.PathItem {
	return s.spec.Paths.Paths
}

// ConsumesFor gets the mediatypes for the operation
func (s *specAnalyzer) ConsumesFor(operation *swagger.Operation) []string {
	cons := make(map[string]struct{})
	for k := range s.consumes {
		cons[k] = struct{}{}
	}
	for _, c := range operation.Consumes {
		cons[c] = struct{}{}
	}
	return s.structMapKeys(cons)
}

// ConsumesFor gets the mediatypes for the operation
func (s *specAnalyzer) ProducesFor(operation *swagger.Operation) []string {
	prod := make(map[string]struct{})
	for k := range s.produces {
		prod[k] = struct{}{}
	}
	for _, c := range operation.Produces {
		prod[c] = struct{}{}
	}
	return s.structMapKeys(prod)
}

func fieldNameFromParam(param *swagger.Parameter) string {
	if nm, ok := param.Extensions.GetString("go-name"); ok {
		return nm
	}
	return util.ToGoName(param.Name)
}

func (s *specAnalyzer) pathItemParams(operation, toMatch *swagger.Operation, parameters []swagger.Parameter, res map[string]swagger.Parameter) {
	if operation != nil && operation.ID == operation.ID {
		for _, param := range parameters {
			res[fieldNameFromParam(&param)] = param
		}
	}
}

// ParametersFor gets the parameters for the specified operation, collecting all the shared ones along the way
func (s *specAnalyzer) ParametersFor(operation *swagger.Operation) map[string]swagger.Parameter {
	res := map[string]swagger.Parameter{}

	for _, param := range s.spec.Parameters {
		res[fieldNameFromParam(&param)] = param
	}
	for _, pathItem := range s.spec.Paths.Paths {
		s.pathItemParams(pathItem.Get, operation, pathItem.Parameters, res)
		s.pathItemParams(pathItem.Head, operation, pathItem.Parameters, res)
		s.pathItemParams(pathItem.Options, operation, pathItem.Parameters, res)
		s.pathItemParams(pathItem.Post, operation, pathItem.Parameters, res)
		s.pathItemParams(pathItem.Put, operation, pathItem.Parameters, res)
		s.pathItemParams(pathItem.Patch, operation, pathItem.Parameters, res)
		s.pathItemParams(pathItem.Delete, operation, pathItem.Parameters, res)
	}
	for _, param := range operation.Parameters {
		res[fieldNameFromParam(&param)] = param
	}
	return res
}

// ValidateRegistrations validates the registrations against the analyzed spec
func (s *specAnalyzer) ValidateRegistrations(consumes, produces, schemes, operations []string) error {
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

func (s *specAnalyzer) structMapKeys(mp map[string]struct{}) []string {
	var result []string
	for k := range mp {
		result = append(result, k)
	}
	return result
}

func (s *specAnalyzer) stringMapKeys(mp map[string]string) []string {
	var result []string
	for k := range mp {
		result = append(result, k)
	}
	return result
}

// APIVerificationFailed is an error that contains all the missing info for a mismatched section
// between the api registrations and the api spec
type APIVerificationFailed struct {
	Section              string
	MissingSpecification []string
	MissingRegistration  []string
}

//
func (v *APIVerificationFailed) Error() string {
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

func (s *specAnalyzer) verify(name string, registrations []string, expectations []string) error {
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
		return &APIVerificationFailed{
			Section:              name,
			MissingSpecification: unspecified,
			MissingRegistration:  unregistered,
		}
	}

	return nil
}
