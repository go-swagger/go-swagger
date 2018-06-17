package generator

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

// Perform common initialization of template repository before running tests.
// This allows to run tests unitarily (e.g. go test -run xxx ).
func TestMain(m *testing.M) {
	templates.LoadDefaults()
	os.Exit(m.Run())
}

func testGenOpts() (g GenOpts) {
	g.Target = "."
	g.APIPackage = "operations"
	g.ModelPackage = "models"
	g.ServerPackage = "restapi"
	g.ClientPackage = "client"
	g.Principal = ""
	g.DefaultScheme = "http"
	g.IncludeModel = true
	g.IncludeValidator = true
	g.IncludeHandler = true
	g.IncludeParameters = true
	g.IncludeResponses = true
	g.IncludeMain = false
	g.IncludeSupport = true
	g.ExcludeSpec = true
	g.TemplateDir = ""
	g.WithContext = false
	g.DumpData = false
	_ = g.EnsureDefaults()
	return
}

func testAppGenerator(t testing.TB, specPath, name string) (*appGenerator, error) {
	specDoc, err := loads.Spec(specPath)
	if !assert.NoError(t, err) {
		return nil, err
	}
	analyzed := analysis.New(specDoc.Spec())

	models, err := gatherModels(specDoc, nil)
	if !assert.NoError(t, err) {
		return nil, err
	}

	operations := gatherOperations(analyzed, nil)
	if len(operations) == 0 {
		return nil, errors.New("no operations were selected")
	}

	opts := testGenOpts()
	opts.Spec = specPath
	apiPackage := opts.LanguageOpts.MangleName(swag.ToFileName(opts.APIPackage), "api")

	return &appGenerator{
		Name:            appNameOrDefault(specDoc, name, "swagger"),
		Receiver:        "o",
		SpecDoc:         specDoc,
		Analyzed:        analyzed,
		Models:          models,
		Operations:      operations,
		Target:          ".",
		DumpData:        opts.DumpData,
		Package:         apiPackage,
		APIPackage:      apiPackage,
		ModelsPackage:   opts.LanguageOpts.MangleName(swag.ToFileName(opts.ModelPackage), "definitions"),
		ServerPackage:   opts.LanguageOpts.MangleName(swag.ToFileName(opts.ServerPackage), "server"),
		ClientPackage:   opts.LanguageOpts.MangleName(swag.ToFileName(opts.ClientPackage), "client"),
		Principal:       opts.Principal,
		DefaultScheme:   "http",
		DefaultProduces: runtime.JSONMime,
		DefaultConsumes: runtime.JSONMime,
		GenOpts:         &opts,
	}, nil
}

func TestServer_UrlEncoded(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)
	gen, err := testAppGenerator(t, "../fixtures/codegen/simplesearch.yml", "search")
	if assert.NoError(t, err) {
		app, err := gen.makeCodegenApp()
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			if assert.NoError(t, templates.MustGet("serverBuilder").Execute(buf, app)) {
				formatted, err := app.GenOpts.LanguageOpts.FormatContent("search_api.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(formatted)
					assert.Regexp(t, "UrlformConsumer:\\s+runtime\\.DiscardConsumer", res)
				} else {
					fmt.Println(buf.String())
				}
			}
			buf = bytes.NewBuffer(nil)
			if assert.NoError(t, templates.MustGet("serverConfigureapi").Execute(buf, app)) {
				formatted, err := app.GenOpts.LanguageOpts.FormatContent("configure_search_api.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(formatted)
					assertInCode(t, "api.UrlformConsumer = runtime.DiscardConsumer", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestServer_MultipartForm(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)
	gen, err := testAppGenerator(t, "../fixtures/codegen/shipyard.yml", "shipyard")
	if assert.NoError(t, err) {
		app, err := gen.makeCodegenApp()
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			if assert.NoError(t, templates.MustGet("serverBuilder").Execute(buf, app)) {
				formatted, err := app.GenOpts.LanguageOpts.FormatContent("shipyard_api.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(formatted)
					assert.Regexp(t, "MultipartformConsumer:\\s+runtime\\.DiscardConsumer", res)
				} else {
					fmt.Println(buf.String())
				}
			}
			buf = bytes.NewBuffer(nil)
			if assert.NoError(t, templates.MustGet("serverConfigureapi").Execute(buf, app)) {
				formatted, err := app.GenOpts.LanguageOpts.FormatContent("configure_shipyard_api.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(formatted)
					assertInCode(t, "api.MultipartformConsumer = runtime.DiscardConsumer", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestServer_InvalidSpec(t *testing.T) {
	opts := testGenOpts()
	opts.Spec = "../fixtures/bugs/825/swagger.yml"
	opts.ValidateSpec = true
	assert.Error(t, GenerateServer("foo", nil, nil, &opts))
}

func TestServer_TrailingSlash(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)
	gen, err := testAppGenerator(t, "../fixtures/bugs/899/swagger.yml", "trailing slash")
	if assert.NoError(t, err) {
		app, err := gen.makeCodegenApp()
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			if assert.NoError(t, templates.MustGet("serverBuilder").Execute(buf, app)) {
				formatted, err := app.GenOpts.LanguageOpts.FormatContent("shipyard_api.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(formatted)
					assertInCode(t, `o.handlers["GET"]["/trailingslashpath"]`, res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestServer_Issue987(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)
	gen, err := testAppGenerator(t, "../fixtures/bugs/987/swagger.yml", "deeper consumes produces")
	if assert.NoError(t, err) {
		app, err := gen.makeCodegenApp()
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			if assert.NoError(t, templates.MustGet("serverBuilder").Execute(buf, app)) {
				formatted, err := app.GenOpts.LanguageOpts.FormatContent("shipyard_api.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(formatted)
					assertRegexpInCode(t, `JSONConsumer:\s+runtime.JSONConsumer()`, res)
					assertRegexpInCode(t, `JSONProducer:\s+runtime.JSONProducer()`, res)
					assertInCode(t, `result["application/json"] = o.JSONConsumer`, res)
					assertInCode(t, `result["application/json"] = o.JSONProducer`, res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestServer_FilterByTag(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)
	gen, err := testAppGenerator(t, "../fixtures/codegen/simplesearch.yml", "search")
	if assert.NoError(t, err) {
		gen.GenOpts.Tags = []string{"search"}
		app, err := gen.makeCodegenApp()
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			if assert.NoError(t, templates.MustGet("serverBuilder").Execute(buf, app)) {
				formatted, err := app.GenOpts.LanguageOpts.FormatContent("search_api.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(formatted)
					assertInCode(t, `o.handlers["POST"]["/search"]`, res)
					assertNotInCode(t, `o.handlers["POST"]["/tasks"]`, res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

// Checking error handling code: panic on mismatched template
// High level test with AppGenerator
func badTemplateCall() {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	gen, err := testAppGenerator(nil, "../fixtures/bugs/899/swagger.yml", "trailing slash")
	if err != nil {
		return
	}
	app, err := gen.makeCodegenApp()
	log.SetOutput(ioutil.Discard)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(nil)
	r := templates.MustGet("serverBuilderX").Execute(buf, app)

	// Should never reach here
	log.Printf("%+v\n", r)
}

func TestServer_BadTemplate(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	assert.Panics(t, badTemplateCall, "templates.MustGet() did not panic() as currently expected")
}

// Checking error handling code: panic on bad parsing template
// High level test with AppGenerator
func badParseCall() {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	var badParse = `{{{ define "T1" }}T1{{end}}{{ define "T2" }}T2{{end}}`

	templates.AddFile("badparse", badParse)
	gen, _ := testAppGenerator(nil, "../fixtures/bugs/899/swagger.yml", "trailing slash")
	app, _ := gen.makeCodegenApp()
	log.SetOutput(ioutil.Discard)
	tpl := templates.MustGet("badparse")

	// Should never reach here
	buf := bytes.NewBuffer(nil)
	r := tpl.Execute(buf, app)

	log.Printf("%+v\n", r)
}

func TestServer_ErrorParsingTemplate(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	assert.Panics(t, badParseCall, "templates.MustGet() did not panic() as currently expected")
}

func TestServer_OperationGroups(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
		os.RemoveAll(filepath.Join(".", "restapi"))
		os.RemoveAll(filepath.Join(".", "search"))
		os.RemoveAll(filepath.Join(".", "tasks"))
	}()

	gen, err := testAppGenerator(t, "../fixtures/codegen/simplesearch.yml", "search")
	if assert.NoError(t, err) {
		gen.GenOpts.Tags = []string{"search", "tasks"}
		gen.GenOpts.IncludeModel = false
		gen.GenOpts.IncludeHandler = true
		gen.GenOpts.Sections.OperationGroups = []TemplateOpts{
			{
				Name:     "opGroupTest",
				Source:   "asset:opGroupTest",
				Target:   "{{ joinFilePath .Target .Name }}",
				FileName: "{{ (snakize (pascalize .Name)) }}_opgroup_test.gol",
			},
		}
		err := gen.Generate()
		// This attempts fails: template not declared
		assert.Error(t, err)
		// Tolerates case variations on error message
		assert.Contains(t, strings.ToLower(err.Error()), "template doesn't exist")

		var opGroupTpl = `
// OperationGroupName={{.Name}}
// RootPackage={{.RootPackage}}
{{ range .Operations }}
	// OperationName={{.Name}}
{{end}}`
		templates.AddFile("opGroupTest", opGroupTpl)
		err = gen.Generate()
		assert.NoError(t, err)
		//buf := bytes.NewBuffer(nil)
		genContent, erf := ioutil.ReadFile("./search/search_opgroup_test.gol")
		assert.NoError(t, erf, "Generator should have written a file")
		assert.Contains(t, string(genContent), "// OperationGroupName=search")
		assert.Contains(t, string(genContent), "// RootPackage=operations")
		assert.Contains(t, string(genContent), "// OperationName=search")

		genContent, erf = ioutil.ReadFile("./tasks/tasks_opgroup_test.gol")
		assert.NoError(t, erf, "Generator should have written a file")
		assert.Contains(t, string(genContent), "// OperationGroupName=tasks")
		assert.Contains(t, string(genContent), "// RootPackage=operations")
		assert.Contains(t, string(genContent), "// OperationName=createTask")
		assert.Contains(t, string(genContent), "// OperationName=deleteTask")
		assert.Contains(t, string(genContent), "// OperationName=getTasks")
		assert.Contains(t, string(genContent), "// OperationName=updateTask")
	}
}

func TestServer_Issue1301(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)
	gen, err := testAppGenerator(t, "../fixtures/enhancements/1301/swagger.yml", "custom producers")
	if assert.NoError(t, err) {
		app, err := gen.makeCodegenApp()
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			if assert.NoError(t, templates.MustGet("serverBuilder").Execute(buf, app)) {
				formatted, err := app.GenOpts.LanguageOpts.FormatContent("shipyard_api.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(formatted)

					// initialisation in New<Name>API function
					assertInCode(t, `customConsumers:     make(map[string]runtime.Consumer)`, res)
					assertInCode(t, `customProducers:     make(map[string]runtime.Producer)`, res)

					// declaration in struct
					assertInCode(t, `customConsumers map[string]runtime.Consumer`, res)
					assertInCode(t, `customProducers map[string]runtime.Producer`, res)
					assertRegexpInCode(t, `if c, ok := o\.customConsumers\[mt\]; ok \{\s+result\[mt\] = c\s+\}`, res)
					assertRegexpInCode(t, `if p, ok := o\.customProducers\[mt\]; ok \{\s+result\[mt\] = p\s+\}`, res)
					assertRegexpInCode(t, `func \(o \*CustomProducersAPI\) RegisterConsumer\(mediaType string, consumer runtime\.Consumer\) \{\s+	o\.customConsumers\[mediaType\] = consumer\s+\}`, res)
					assertRegexpInCode(t, `func \(o \*CustomProducersAPI\) RegisterProducer\(mediaType string, producer runtime\.Producer\) \{\s+	o\.customProducers\[mediaType\] = producer\s+\}`, res)

				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestServer_Issue1557(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)
	gen, err := testAppGenerator(t, "../fixtures/enhancements/1557/swagger.yml", "generate consumer/producer handlers that are not whitelisted")
	if assert.NoError(t, err) {
		app, err := gen.makeCodegenApp()
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			if assert.NoError(t, templates.MustGet("serverBuilder").Execute(buf, app)) {
				formatted, err := app.GenOpts.LanguageOpts.FormatContent("shipyard_api.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(formatted)
					assertRegexpInCode(t, `ApplicationPdfConsumer:\s+runtime.Consumer`, res)
					assertRegexpInCode(t, `ApplicationPdfProducer:\s+runtime.Producer`, res)
					assertInCode(t, `result["application/pdf"] = o.ApplicationPdfConsumer`, res)
					assertInCode(t, `result["application/pdf"] = o.ApplicationPdfProducer`, res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}
