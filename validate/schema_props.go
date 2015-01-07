package validate

import (
	"reflect"

	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
)

// TODO: use this to implement straight struct validation
// type nameProvider struct {
// 	*sync.Mutex
// 	index map[reflect.Type]nameIndex
// }

// type nameIndex struct {
// 	jsonNames map[string]string
// 	goNames   map[string]string
// }

// func newNameProvider() *nameProvider {
// 	return &nameProvider{index: make(map[reflect.Type]nameIndex)}
// }

// func buildNameIndex(tpe reflect.Type, idx, reverseIdx map[string]string) {
// 	for i := 0; i < tpe.NumField(); i++ {
// 		targetDes := tpe.Field(i)

// 		if targetDes.PkgPath != "" { // unexported
// 			continue
// 		}

// 		if targetDes.Anonymous { // walk embedded structures tree down first
// 			buildNameIndex(targetDes.Type, idx, reverseIdx)
// 			continue
// 		}

// 		if tag := targetDes.Tag.Get("json"); tag != "" {
// 			parts := strings.Split(tag, ",")
// 			if len(parts) == 0 {
// 				continue
// 			}
// 			nm := parts[0]
// 			if nm == "" || nm == "-" {
// 				continue
// 			}
// 			idx[nm] = targetDes.Name
// 			reverseIdx[targetDes.Name] = nm
// 		}
// 	}
// }

// func newNameIndex(tpe reflect.Type) nameIndex {
// 	var idx = make(map[string]string, tpe.NumField())
// 	var reverseIdx = make(map[string]string, tpe.NumField())

// 	buildNameIndex(tpe, idx, reverseIdx)
// 	return nameIndex{jsonNames: idx, goNames: reverseIdx}
// }

// func (n *nameProvider) GetJSONName(subject interface{}, name string) string {
// 	tpe := reflect.Indirect(reflect.ValueOf(subject)).Type()
// 	names, ok := n.index[tpe]
// 	if !ok {
// 		n.Lock()
// 		names = newNameIndex(tpe)
// 		n.index[tpe] = names
// 		n.Unlock()
// 	}
// 	return names.jsonNames[name]
// }

// func (n *nameProvider) GetGoName(subject interface{}, name string) string {
// 	tpe := reflect.Indirect(reflect.ValueOf(subject)).Type()
// 	names, ok := n.index[tpe]
// 	if !ok {
// 		n.Lock()
// 		names = newNameIndex(tpe)
// 		n.index[tpe] = names
// 		n.Unlock()
// 	}
// 	return names.goNames[name]
// }

// var nameTranslations = newNameProvider()

type schemaPropsValidator struct {
	Path            string
	In              string
	AllOf           []spec.Schema
	OneOf           []spec.Schema
	AnyOf           []spec.Schema
	Not             *spec.Schema
	Dependencies    spec.Dependencies
	anyOfValidators []schemaValidator
	allOfValidators []schemaValidator
	oneOfValidators []schemaValidator
	notValidator    *schemaValidator
}

func (s *schemaPropsValidator) SetPath(path string) {
	s.Path = path
}

func newSchemaPropsValidator(path string, in string, allOf, oneOf, anyOf []spec.Schema, not *spec.Schema, deps spec.Dependencies) *schemaPropsValidator {
	var anyValidators []schemaValidator
	for _, v := range anyOf {
		anyValidators = append(anyValidators, *newSchemaValidator(&v, path))
	}
	var allValidators []schemaValidator
	for _, v := range allOf {
		allValidators = append(allValidators, *newSchemaValidator(&v, path))
	}
	var oneValidators []schemaValidator
	for _, v := range oneOf {
		oneValidators = append(oneValidators, *newSchemaValidator(&v, path))
	}

	var notValidator *schemaValidator
	if not != nil {
		notValidator = newSchemaValidator(not, path)
	}

	return &schemaPropsValidator{
		Path:            path,
		In:              in,
		AllOf:           allOf,
		OneOf:           oneOf,
		AnyOf:           anyOf,
		Not:             not,
		Dependencies:    deps,
		anyOfValidators: anyValidators,
		allOfValidators: allValidators,
		oneOfValidators: oneValidators,
		notValidator:    notValidator,
	}
}

func (s *schemaPropsValidator) Applies(source interface{}, kind reflect.Kind) bool {
	return reflect.TypeOf(source) == specSchemaType
}

func (s *schemaPropsValidator) Validate(data interface{}) *result {
	mainResult := &result{}
	if len(s.anyOfValidators) > 0 {
		var bestFailures *result
		succeededOnce := false
		for _, anyOfSchema := range s.anyOfValidators {
			result := anyOfSchema.Validate(data)
			if result.IsValid() {
				bestFailures = nil
				succeededOnce = true
				break
			}
			if bestFailures == nil || result.MatchCount > bestFailures.MatchCount {
				bestFailures = result
			}
		}

		if !succeededOnce {
			mainResult.AddErrors(errors.New(422, "must validate at least one schema (anyOf)"))
		}
		if bestFailures != nil {
			mainResult.Merge(bestFailures)
		}
	}

	if len(s.oneOfValidators) > 0 {
		var bestFailures *result
		validated := 0

		for _, oneOfSchema := range s.oneOfValidators {
			result := oneOfSchema.Validate(data)
			if result.IsValid() {
				validated++
				bestFailures = nil
				continue
			}
			if validated == 0 && (bestFailures == nil || result.MatchCount > bestFailures.MatchCount) {
				bestFailures = result
			}
		}

		if validated != 1 {
			mainResult.AddErrors(errors.New(422, "must validate one and only one schema (oneOf)"))
			if bestFailures != nil {
				mainResult.Merge(bestFailures)
			}
		}
	}

	if len(s.allOfValidators) > 0 {
		validated := 0

		for _, allOfSchema := range s.allOfValidators {
			result := allOfSchema.Validate(data)
			if result.IsValid() {
				validated++
			}
			mainResult.Merge(result)
		}

		if validated != len(s.allOfValidators) {
			mainResult.AddErrors(errors.New(422, "must validate all the schemas (allOf)"))
		}
	}

	if s.notValidator != nil {
		result := s.notValidator.Validate(data)
		if result.IsValid() {
			mainResult.AddErrors(errors.New(422, "must not validate the schema (not)"))
		}
	}

	if s.Dependencies != nil && len(s.Dependencies) > 0 && reflect.TypeOf(data).Kind() == reflect.Map {
		val := data.(map[string]interface{})
		for key := range val {
			if dep, ok := s.Dependencies[key]; ok {

				if dep.Schema != nil {
					mainResult.Merge(newSchemaValidator(dep.Schema, s.Path+"."+key).Validate(data))
					continue
				}

				if len(dep.Property) > 0 {
					for _, depKey := range dep.Property {
						if _, ok := val[depKey]; !ok {
							mainResult.AddErrors(errors.New(422, "has a dependency on %s", depKey))
						}
					}
				}
			}
		}
	}

	mainResult.Inc()
	return mainResult
}

// IsZero returns true when the value is a zero for the type
func isZero(data reflect.Value) bool {
	if !data.CanInterface() {
		return true
	}
	tpe := data.Type()
	return reflect.DeepEqual(data.Interface(), reflect.Zero(tpe).Interface())
}
