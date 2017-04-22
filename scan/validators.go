// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scan

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-openapi/loads/fmts"
	"github.com/go-openapi/spec"
	"gopkg.in/yaml.v2"
)

type validationBuilder interface {
	SetMaximum(float64, bool)
	SetMinimum(float64, bool)
	SetMultipleOf(float64)

	SetMinItems(int64)
	SetMaxItems(int64)

	SetMinLength(int64)
	SetMaxLength(int64)
	SetPattern(string)

	SetUnique(bool)
	SetEnum(string)
	SetDefault(string)
}

type valueParser interface {
	Parse([]string) error
	Matches(string) bool
}

type setMaximum struct {
	builder validationBuilder
	rx      *regexp.Regexp
}

func (sm *setMaximum) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := sm.rx.FindStringSubmatch(lines[0])
	if len(matches) > 2 && len(matches[2]) > 0 {
		max, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			return err
		}
		sm.builder.SetMaximum(max, matches[1] == "<")
	}
	return nil
}

func (sm *setMaximum) Matches(line string) bool {
	return sm.rx.MatchString(line)
}

type setMinimum struct {
	builder validationBuilder
	rx      *regexp.Regexp
}

func (sm *setMinimum) Matches(line string) bool {
	return sm.rx.MatchString(line)
}

func (sm *setMinimum) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := sm.rx.FindStringSubmatch(lines[0])
	if len(matches) > 2 && len(matches[2]) > 0 {
		min, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			return err
		}
		sm.builder.SetMinimum(min, matches[1] == ">")
	}
	return nil
}

type setMultipleOf struct {
	builder validationBuilder
	rx      *regexp.Regexp
}

func (sm *setMultipleOf) Matches(line string) bool {
	return sm.rx.MatchString(line)
}

func (sm *setMultipleOf) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := sm.rx.FindStringSubmatch(lines[0])
	if len(matches) > 2 && len(matches[1]) > 0 {
		multipleOf, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return err
		}
		sm.builder.SetMultipleOf(multipleOf)
	}
	return nil
}

type setMaxItems struct {
	builder validationBuilder
	rx      *regexp.Regexp
}

func (sm *setMaxItems) Matches(line string) bool {
	return sm.rx.MatchString(line)
}

func (sm *setMaxItems) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := sm.rx.FindStringSubmatch(lines[0])
	if len(matches) > 1 && len(matches[1]) > 0 {
		maxItems, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return err
		}
		sm.builder.SetMaxItems(maxItems)
	}
	return nil
}

type setMinItems struct {
	builder validationBuilder
	rx      *regexp.Regexp
}

func (sm *setMinItems) Matches(line string) bool {
	return sm.rx.MatchString(line)
}

func (sm *setMinItems) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := sm.rx.FindStringSubmatch(lines[0])
	if len(matches) > 1 && len(matches[1]) > 0 {
		minItems, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return err
		}
		sm.builder.SetMinItems(minItems)
	}
	return nil
}

type setMaxLength struct {
	builder validationBuilder
	rx      *regexp.Regexp
}

func (sm *setMaxLength) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := sm.rx.FindStringSubmatch(lines[0])
	if len(matches) > 1 && len(matches[1]) > 0 {
		maxLength, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return err
		}
		sm.builder.SetMaxLength(maxLength)
	}
	return nil
}

func (sm *setMaxLength) Matches(line string) bool {
	return sm.rx.MatchString(line)
}

type setMinLength struct {
	builder validationBuilder
	rx      *regexp.Regexp
}

func (sm *setMinLength) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := sm.rx.FindStringSubmatch(lines[0])
	if len(matches) > 1 && len(matches[1]) > 0 {
		minLength, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return err
		}
		sm.builder.SetMinLength(minLength)
	}
	return nil
}

func (sm *setMinLength) Matches(line string) bool {
	return sm.rx.MatchString(line)
}

type setPattern struct {
	builder validationBuilder
	rx      *regexp.Regexp
}

func (sm *setPattern) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := sm.rx.FindStringSubmatch(lines[0])
	if len(matches) > 1 && len(matches[1]) > 0 {
		sm.builder.SetPattern(matches[1])
	}
	return nil
}

func (sm *setPattern) Matches(line string) bool {
	return sm.rx.MatchString(line)
}

type setCollectionFormat struct {
	builder operationValidationBuilder
	rx      *regexp.Regexp
}

func (sm *setCollectionFormat) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := sm.rx.FindStringSubmatch(lines[0])
	if len(matches) > 1 && len(matches[1]) > 0 {
		sm.builder.SetCollectionFormat(matches[1])
	}
	return nil
}

func (sm *setCollectionFormat) Matches(line string) bool {
	return sm.rx.MatchString(line)
}

type setUnique struct {
	builder validationBuilder
	rx      *regexp.Regexp
}

func (su *setUnique) Matches(line string) bool {
	return su.rx.MatchString(line)
}

func (su *setUnique) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := su.rx.FindStringSubmatch(lines[0])
	if len(matches) > 1 && len(matches[1]) > 0 {
		req, err := strconv.ParseBool(matches[1])
		if err != nil {
			return err
		}
		su.builder.SetUnique(req)
	}
	return nil
}

type setEnum struct {
	builder validationBuilder
	rx      *regexp.Regexp
}

func (se *setEnum) Matches(line string) bool {
	return se.rx.MatchString(line)
}

func (se *setEnum) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := se.rx.FindStringSubmatch(lines[0])
	if len(matches) > 1 && len(matches[1]) > 0 {
		se.builder.SetEnum(matches[1])
	}
	return nil
}

type setDefault struct {
	builder validationBuilder
	rx      *regexp.Regexp
}

func (sd *setDefault) Matches(line string) bool {
	return sd.rx.MatchString(line)
}

func (sd *setDefault) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := sd.rx.FindStringSubmatch(lines[0])
	if len(matches) > 1 && len(matches[1]) > 0 {
		sd.builder.SetDefault(matches[1])
	}
	return nil
}

type matchOnlyParam struct {
	tgt *spec.Parameter
	rx  *regexp.Regexp
}

func (mo *matchOnlyParam) Matches(line string) bool {
	return mo.rx.MatchString(line)
}

func (mo *matchOnlyParam) Parse(lines []string) error {
	return nil
}

type setRequiredParam struct {
	tgt *spec.Parameter
}

func (su *setRequiredParam) Matches(line string) bool {
	return rxRequired.MatchString(line)
}

func (su *setRequiredParam) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := rxRequired.FindStringSubmatch(lines[0])
	if len(matches) > 1 && len(matches[1]) > 0 {
		req, err := strconv.ParseBool(matches[1])
		if err != nil {
			return err
		}
		su.tgt.Required = req
	}
	return nil
}

type setReadOnlySchema struct {
	tgt *spec.Schema
}

func (su *setReadOnlySchema) Matches(line string) bool {
	return rxReadOnly.MatchString(line)
}

func (su *setReadOnlySchema) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := rxReadOnly.FindStringSubmatch(lines[0])
	if len(matches) > 1 && len(matches[1]) > 0 {
		req, err := strconv.ParseBool(matches[1])
		if err != nil {
			return err
		}
		su.tgt.ReadOnly = req
	}
	return nil
}

type setDiscriminator struct {
	schema *spec.Schema
	field  string
}

func (su *setDiscriminator) Matches(line string) bool {
	return rxDiscriminator.MatchString(line)
}

func (su *setDiscriminator) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := rxDiscriminator.FindStringSubmatch(lines[0])
	if len(matches) > 1 && len(matches[1]) > 0 {
		req, err := strconv.ParseBool(matches[1])
		if err != nil {
			return err
		}
		if req {
			su.schema.Discriminator = su.field
		} else {
			if su.schema.Discriminator == su.field {
				su.schema.Discriminator = ""
			}
		}
	}
	return nil
}

type setRequiredSchema struct {
	schema *spec.Schema
	field  string
}

func (su *setRequiredSchema) Matches(line string) bool {
	return rxRequired.MatchString(line)
}

func (su *setRequiredSchema) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := rxRequired.FindStringSubmatch(lines[0])
	if len(matches) > 1 && len(matches[1]) > 0 {
		req, err := strconv.ParseBool(matches[1])
		if err != nil {
			return err
		}
		midx := -1
		for i, nm := range su.schema.Required {
			if nm == su.field {
				midx = i
				break
			}
		}
		if req {
			if midx < 0 {
				su.schema.Required = append(su.schema.Required, su.field)
			}
		} else if midx >= 0 {
			su.schema.Required = append(su.schema.Required[:midx], su.schema.Required[midx+1:]...)
		}
	}
	return nil
}

func newMultilineDropEmptyParser(rx *regexp.Regexp, set func([]string)) *multiLineDropEmptyParser {
	return &multiLineDropEmptyParser{
		rx:  rx,
		set: set,
	}
}

type multiLineDropEmptyParser struct {
	set func([]string)
	rx  *regexp.Regexp
}

func (m *multiLineDropEmptyParser) Matches(line string) bool {
	return m.rx.MatchString(line)
}

func (m *multiLineDropEmptyParser) Parse(lines []string) error {
	m.set(removeEmptyLines(lines))
	return nil
}

func newSetSchemes(set func([]string)) *setSchemes {
	return &setSchemes{
		set: set,
		rx:  rxSchemes,
	}
}

type setSchemes struct {
	set func([]string)
	rx  *regexp.Regexp
}

func (ss *setSchemes) Matches(line string) bool {
	return ss.rx.MatchString(line)
}

func (ss *setSchemes) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := ss.rx.FindStringSubmatch(lines[0])
	if len(matches) > 1 && len(matches[1]) > 0 {
		sch := strings.Split(matches[1], ", ")

		var schemes []string
		for _, s := range sch {
			ts := strings.TrimSpace(s)
			if ts != "" {
				schemes = append(schemes, ts)
			}
		}
		ss.set(schemes)
	}
	return nil
}

func newYAMLBlockParser(rx *regexp.Regexp, setter func(interface{}) error) *yamlBlockParser {
	return &yamlBlockParser{
		set: setter,
		rx:  rx,
	}
}

type yamlBlockParser struct {
	set func(interface{}) error
	rx  *regexp.Regexp
}

func (se *yamlBlockParser) Matches(line string) bool {
	return se.rx.MatchString(line)
}

func (se *yamlBlockParser) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}

	if lines[0] != "---" {
		lines = append([]string{"---"}, lines...)
	}

	yamlContent := strings.Join(lines, "\n")
	var yamlValue interface{}
	err := yaml.Unmarshal([]byte(yamlContent), &yamlValue)
	if err != nil {
		return err
	}

	var jsonValue json.RawMessage
	jsonValue, err = fmts.YAMLToJSON(yamlValue)
	if err != nil {
		return err
	}

	var jsonData interface{}
	err = json.Unmarshal(jsonValue, &jsonData)
	if err != nil {
		return err
	}

	return se.set(jsonData)
}

func newSetSecurityDefinitions(rx *regexp.Regexp, setter func(spec.SecurityDefinitions)) *setSecurityDefinitions {
	return &setSecurityDefinitions{
		set: setter,
		rx:  rx,
	}
}

type setSecurityDefinitions struct {
	set func(spec.SecurityDefinitions)
	rx  *regexp.Regexp
}

func (ss *setSecurityDefinitions) Matches(line string) bool {
	return ss.rx.MatchString(line)
}

var (
	rxSecuritySchemeType          = regexp.MustCompile(`[Tt]ype\p{Zs}*:`)
	rxSecuritySchemeName          = regexp.MustCompile(`[Nn]ame\p{Zs}*:`)
	rxSecuritySchemeIn            = regexp.MustCompile(`[Ii]n\p{Zs}*:`)
	rxSecuritySchemeFlow          = regexp.MustCompile(`[Ff]low\p{Zs}*:`)
	rxSecuritySchemeDescription   = regexp.MustCompile(`[Dd]escription\p{Zs}*:`)
	rxSecuritySchemeAuthorization = regexp.MustCompile(`[Aa]uthorizationUrl\p{Zs}*:`)
	rxSecuritySchemeToken         = regexp.MustCompile(`[Tt]okenUrl\p{Zs}*:`)
)

func (ss *setSecurityDefinitions) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}

	result := spec.SecurityDefinitions{}
	var scheme spec.SecurityScheme
	var key string
	var tp []tagParser
	for i := 0; i < len(lines); i++ {
		kv := strings.SplitN(lines[i], ":", 2)
		if len(kv) <= 1 {
			return fmt.Errorf("invalid format for securityDefinitions: %s", lines[i])
		}

		k, v := kv[0], strings.TrimSpace(kv[1])

		if v == "" {
			if key != "" {
				result[key] = &scheme
			}
			scheme = spec.SecurityScheme{}
			key = k
			tp = []tagParser{
				newSingleLineTagParser("type", newSetField(rxSecuritySchemeType, setSecuritySchemeType(&scheme))),
				newSingleLineTagParser("name", newSetField(rxSecuritySchemeName, setSecuritySchemeName(&scheme))),
				newSingleLineTagParser("in", newSetField(rxSecuritySchemeIn, setSecuritySchemeIn(&scheme))),
				newSingleLineTagParser("flow", newSetField(rxSecuritySchemeFlow, setSecuritySchemeFlow(&scheme))),
				newSingleLineTagParser("description", newSetField(rxSecuritySchemeDescription, setSecuritySchemeDescription(&scheme))),
				newSingleLineTagParser("authorizationUrl", newSetField(rxSecuritySchemeAuthorization, setSecuritySchemeAuthorizationURL(&scheme))),
				newSingleLineTagParser("tokenUrl", newSetField(rxSecuritySchemeToken, setSecuritySchemeTokenURL(&scheme))),
			}
			continue
		} else {
			for _, p := range tp {
				if p.Matches(lines[i]) {
					err := p.Parse([]string{lines[i]})
					if err != nil {
						return err
					}
					break
				}
			}
		}
	}
	if _, ok := result[key]; !ok && key != "" {
		result[key] = &scheme
	}

	ss.set(result)
	return nil
}

func setSecuritySchemeType(scheme *spec.SecurityScheme) func(string) {
	return func(val string) { scheme.Type = val }
}

func setSecuritySchemeName(scheme *spec.SecurityScheme) func(string) {
	return func(val string) { scheme.Name = val }
}

func setSecuritySchemeIn(scheme *spec.SecurityScheme) func(string) {
	return func(val string) { scheme.In = val }
}

func setSecuritySchemeFlow(scheme *spec.SecurityScheme) func(string) {
	return func(val string) { scheme.Flow = val }
}

func setSecuritySchemeDescription(scheme *spec.SecurityScheme) func(string) {
	return func(val string) { scheme.Description = val }
}

func setSecuritySchemeAuthorizationURL(scheme *spec.SecurityScheme) func(string) {
	return func(val string) { scheme.AuthorizationURL = val }
}

func setSecuritySchemeTokenURL(scheme *spec.SecurityScheme) func(string) {
	return func(val string) { scheme.TokenURL = val }
}

func newSetField(rx *regexp.Regexp, setter func(string)) *setField {
	return &setField{
		rx:  rx,
		set: setter,
	}
}

type setField struct {
	set func(string)
	rx  *regexp.Regexp
}

func (sf *setField) Matches(line string) bool {
	return sf.rx.MatchString(line)
}

func (sf *setField) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}

	var value string
	for _, line := range lines {
		kv := strings.SplitN(line, ":", 2)
		if len(kv) > 1 {
			value = strings.TrimSpace(kv[1])
			break
		} else {
			return fmt.Errorf("expecting `key: value`, got key only for string: %s", line)
		}
	}
	sf.set(value)
	return nil
}

func newSetSecurity(rx *regexp.Regexp, setter func([]map[string][]string)) *setSecurity {
	return &setSecurity{
		set: setter,
		rx:  rx,
	}
}

type setSecurity struct {
	set func([]map[string][]string)
	rx  *regexp.Regexp
}

func (ss *setSecurity) Matches(line string) bool {
	return ss.rx.MatchString(line)
}

func (ss *setSecurity) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}

	var result []map[string][]string
	for _, line := range lines {
		kv := strings.SplitN(line, ":", 2)
		scopes := []string{}
		var key string

		if len(kv) > 1 {
			scs := strings.Split(kv[1], ",")
			for _, scope := range scs {
				tr := strings.TrimSpace(scope)
				if tr != "" {
					tr = strings.SplitAfter(tr, " ")[0]
					scopes = append(scopes, strings.TrimSpace(tr))
				}
			}

			key = strings.TrimSpace(kv[0])

			result = append(result, map[string][]string{key: scopes})
		}
	}
	ss.set(result)
	return nil
}

func newSetResponses(definitions map[string]spec.Schema, responses map[string]spec.Response, setter func(*spec.Response, map[int]spec.Response)) *setOpResponses {
	return &setOpResponses{
		set:         setter,
		rx:          rxResponses,
		definitions: definitions,
		responses:   responses,
	}
}

type setOpResponses struct {
	set         func(*spec.Response, map[int]spec.Response)
	rx          *regexp.Regexp
	definitions map[string]spec.Schema
	responses   map[string]spec.Response
}

func (ss *setOpResponses) Matches(line string) bool {
	return ss.rx.MatchString(line)
}

//ResponseTag used when specifying a response to point to a defined swagger:response
const ResponseTag = "response"

//BodyTag used when specifying a response to point to a model/schema
const BodyTag = "body"

//DescriptionTag used when specifying a response that gives a description of the response
const DescriptionTag = "description"

func parseTags(line string) (modelOrResponse string, arrays int, isDefinitionRef bool, description string, err error) {
	tags := strings.Split(line, " ")
	parsedModelOrResponse := false

	for i, tagAndValue := range tags {
		tagValList := strings.SplitN(tagAndValue, ":", 2)
		var tag, value string
		if len(tagValList) > 1 {
			tag = tagValList[0]
			value = tagValList[1]
		} else {
			//TODO: Print a warning, and in the long term, do not support not tagged values
			//Add a default tag if none is supplied
			if i == 0 {
				tag = ResponseTag
			} else {
				tag = DescriptionTag
			}
			value = tagValList[0]
		}

		foundModelOrResponse := false
		if !parsedModelOrResponse {
			if tag == BodyTag {
				foundModelOrResponse = true
				isDefinitionRef = true
			}
			if tag == ResponseTag {
				foundModelOrResponse = true
				isDefinitionRef = false
			}
		}
		if foundModelOrResponse {
			//Read the model or response tag
			parsedModelOrResponse = true
			//Check for nested arrays
			arrays = 0
			for strings.HasPrefix(value, "[]") {
				arrays++
				value = value[2:]
			}
			//What's left over is the model name
			modelOrResponse = value
		} else {
			foundDescription := false
			if tag == DescriptionTag {
				foundDescription = true
			}
			if foundDescription {
				//Descriptions are special, they make they read the rest of the line
				descriptionWords := []string{value}
				if i < len(tags)-1 {
					descriptionWords = append(descriptionWords, tags[i+1:]...)
				}
				description = strings.Join(descriptionWords, " ")
				break
			} else {
				if tag == ResponseTag || tag == BodyTag || tag == DescriptionTag {
					err = fmt.Errorf("Found valid tag %s, but not in a valid position", tag)
				} else {
					err = fmt.Errorf("Found invalid tag: %s", tag)
				}
				//return error
				return
			}
		}
	}

	//TODO: Maybe do, if !parsedModelOrResponse {return some error}
	return
}

func (ss *setOpResponses) Parse(lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}

	var def *spec.Response
	var scr map[int]spec.Response

	for _, line := range lines {
		kv := strings.SplitN(line, ":", 2)
		var key, value string

		if len(kv) > 1 {
			key = strings.TrimSpace(kv[0])
			if key == "" {
				// this must be some weird empty line
				continue
			}
			value = strings.TrimSpace(kv[1])
			if value == "" {
				var resp spec.Response
				if strings.EqualFold("default", key) {
					if def == nil {
						def = &resp
					}
				} else {
					if sc, err := strconv.Atoi(key); err == nil {
						if scr == nil {
							scr = make(map[int]spec.Response)
						}
						scr[sc] = resp
					}
				}
				continue
			}
			refTarget, arrays, isDefinitionRef, description, err := parseTags(value)
			if err != nil {
				return err
			}
			//A possible exception for having a definition
			if _, ok := ss.responses[refTarget]; !ok {
				if _, ok := ss.definitions[refTarget]; ok {
					isDefinitionRef = true
				}
			}

			var ref spec.Ref
			if isDefinitionRef {
				if description == "" {
					description = refTarget
				}
				ref, err = spec.NewRef("#/definitions/" + refTarget)
			} else {
				ref, err = spec.NewRef("#/responses/" + refTarget)
			}
			if err != nil {
				return err
			}

			var resp spec.Response

			if !isDefinitionRef {
				resp.Ref = ref
			} else {
				resp.Schema = new(spec.Schema)
				resp.Description = description
				if arrays == 0 {
					resp.Schema.Ref = ref
				} else {
					cs := resp.Schema
					for i := 0; i < arrays; i++ {
						cs.Typed("array", "")
						cs.Items = new(spec.SchemaOrArray)
						cs.Items.Schema = new(spec.Schema)
						cs = cs.Items.Schema
					}
					cs.Ref = ref
				}
			}

			if strings.EqualFold("default", key) {
				if def == nil {
					def = &resp
				}
			} else {
				if sc, err := strconv.Atoi(key); err == nil {
					if scr == nil {
						scr = make(map[int]spec.Response)
					}
					scr[sc] = resp
				}
			}
		}
	}
	ss.set(def, scr)
	return nil
}
