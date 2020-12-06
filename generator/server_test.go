package generator

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const invalidSpecExample = "../fixtures/bugs/825/swagger.yml"

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
		GenOpts:         opts,
	}, nil
}

func TestServer_UrlEncoded(t *testing.T) {
	defer discardOutput()()

	gen, err := testAppGenerator(t, "../fixtures/codegen/simplesearch.yml", "search")
	require.NoError(t, err)

	app, err := gen.makeCodegenApp()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, app.GenOpts.templates.MustGet("serverBuilder").Execute(buf, app))

	formatted, err := app.GenOpts.LanguageOpts.FormatContent("search_api.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assert.Regexp(t, "UrlformConsumer:\\s+runtime\\.DiscardConsumer", res)

	buf = bytes.NewBuffer(nil)
	require.NoError(t, app.GenOpts.templates.MustGet("serverConfigureapi").Execute(buf, app))

	formatted, err = app.GenOpts.LanguageOpts.FormatContent("configure_search_api.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	assertInCode(t, "api.UrlformConsumer = runtime.DiscardConsumer", string(formatted))
}

func TestServer_MultipartForm(t *testing.T) {
	defer discardOutput()()

	gen, err := testAppGenerator(t, "../fixtures/codegen/shipyard.yml", "shipyard")
	require.NoError(t, err)

	app, err := gen.makeCodegenApp()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, app.GenOpts.templates.MustGet("serverBuilder").Execute(buf, app))

	formatted, err := app.GenOpts.LanguageOpts.FormatContent("shipyard_api.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	assert.Regexp(t, "MultipartformConsumer:\\s+runtime\\.DiscardConsumer", string(formatted))

	buf = bytes.NewBuffer(nil)
	require.NoError(t, app.GenOpts.templates.MustGet("serverConfigureapi").Execute(buf, app))

	formatted, err = app.GenOpts.LanguageOpts.FormatContent("configure_shipyard_api.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	assertInCode(t, "api.MultipartformConsumer = runtime.DiscardConsumer", string(formatted))
}

func TestServer_InvalidSpec(t *testing.T) {
	defer discardOutput()()

	opts := testGenOpts()
	opts.Spec = invalidSpecExample
	opts.ValidateSpec = true

	assert.Error(t, GenerateServer("foo", nil, nil, opts))
}

func TestServer_TrailingSlash(t *testing.T) {
	defer discardOutput()()

	gen, err := testAppGenerator(t, "../fixtures/bugs/899/swagger.yml", "trailing slash")
	require.NoError(t, err)

	app, err := gen.makeCodegenApp()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, app.GenOpts.templates.MustGet("serverBuilder").Execute(buf, app))

	formatted, err := app.GenOpts.LanguageOpts.FormatContent("shipyard_api.go", buf.Bytes())
	require.NoError(t, err, buf.String())

	assertInCode(t, `o.handlers["GET"]["/trailingslashpath"]`, string(formatted))
}

func TestServer_Issue987(t *testing.T) {
	defer discardOutput()()

	gen, err := testAppGenerator(t, "../fixtures/bugs/987/swagger.yml", "deeper consumes produces")
	require.NoError(t, err)

	app, err := gen.makeCodegenApp()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, app.GenOpts.templates.MustGet("serverBuilder").Execute(buf, app))

	formatted, err := app.GenOpts.LanguageOpts.FormatContent("shipyard_api.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertRegexpInCode(t, `JSONConsumer:\s+runtime.JSONConsumer()`, res)
	assertRegexpInCode(t, `JSONProducer:\s+runtime.JSONProducer()`, res)
	assertInCode(t, `result["application/json"] = o.JSONConsumer`, res)
	assertInCode(t, `result["application/json"] = o.JSONProducer`, res)
}

func TestServer_FilterByTag(t *testing.T) {
	defer discardOutput()()

	gen, err := testAppGenerator(t, "../fixtures/codegen/simplesearch.yml", "search")
	require.NoError(t, err)

	gen.GenOpts.Tags = []string{"search"}
	app, err := gen.makeCodegenApp()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, app.GenOpts.templates.MustGet("serverBuilder").Execute(buf, app))

	formatted, err := app.GenOpts.LanguageOpts.FormatContent("search_api.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, `o.handlers["POST"]["/search"]`, res)
	assertNotInCode(t, `o.handlers["POST"]["/tasks"]`, res)
}

func TestServer_BadTemplate(t *testing.T) {
	// Checking error handling code: panic on mismatched template

	defer discardOutput()()

	gen, err := testAppGenerator(nil, "../fixtures/bugs/899/swagger.yml", "trailing slash")
	require.NoError(t, err)

	app, err := gen.makeCodegenApp()
	require.NoError(t, err)

	badTemplateCall := func() {
		buf := bytes.NewBuffer(nil)
		_ = app.GenOpts.templates.MustGet("serverBuilderX").Execute(buf, app)
	}

	assert.Panics(t, badTemplateCall, "templates.MustGet() did not panic() as currently expected")
}

func TestServer_ErrorParsingTemplate(t *testing.T) {
	// Checking error handling code: panic on bad parsing template
	// High level test with AppGenerator

	defer discardOutput()()

	var badParse = `{{{ define "T1" }}T1{{end}}{{ define "T2" }}T2{{end}}`

	gen, err := testAppGenerator(nil, "../fixtures/bugs/899/swagger.yml", "trailing slash")
	require.NoError(t, err)

	require.Error(t, gen.GenOpts.templates.AddFile("badparse", badParse)) // template is not loaded

	badParseCall := func() {
		_ = templates.MustGet("badparse") // MustGet panics
	}

	assert.Panics(t, badParseCall, "templates.MustGet() did not panic() as currently expected")
}

func TestServer_OperationGroups(t *testing.T) {
	defer discardOutput()()
	defer func() {
		_ = os.RemoveAll(filepath.Join(".", "restapi"))
		_ = os.RemoveAll(filepath.Join(".", "search"))
		_ = os.RemoveAll(filepath.Join(".", "tasks"))
	}()

	gen, err := testAppGenerator(t, "../fixtures/codegen/simplesearch.yml", "search")
	require.NoError(t, err)

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

	err = gen.Generate()
	require.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "template doesn't exist") // Tolerates case variations on error message

	var opGroupTpl = `
// OperationGroupName={{.Name}}
// RootPackage={{.RootPackage}}
{{ range .Operations }}
	// OperationName={{.Name}}
{{end}}`
	_ = gen.GenOpts.templates.AddFile("opGroupTest", opGroupTpl)
	require.NoError(t, gen.Generate())

	genContent, err := ioutil.ReadFile("./search/search_opgroup_test.gol")
	require.NoError(t, err, "Generator should have written a file")

	assert.Contains(t, string(genContent), "// OperationGroupName=search")
	assert.Contains(t, string(genContent), "// RootPackage=operations")
	assert.Contains(t, string(genContent), "// OperationName=search")

	genContent, err = ioutil.ReadFile("./tasks/tasks_opgroup_test.gol")
	require.NoError(t, err, "Generator should have written a file")

	assert.Contains(t, string(genContent), "// OperationGroupName=tasks")
	assert.Contains(t, string(genContent), "// RootPackage=operations")
	assert.Contains(t, string(genContent), "// OperationName=createTask")
	assert.Contains(t, string(genContent), "// OperationName=deleteTask")
	assert.Contains(t, string(genContent), "// OperationName=getTasks")
	assert.Contains(t, string(genContent), "// OperationName=updateTask")
}

func TestServer_Issue1301(t *testing.T) {
	defer discardOutput()()

	gen, err := testAppGenerator(t, "../fixtures/enhancements/1301/swagger.yml", "custom producers")
	require.NoError(t, err)

	app, err := gen.makeCodegenApp()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, app.GenOpts.templates.MustGet("serverBuilder").Execute(buf, app))

	formatted, err := app.GenOpts.LanguageOpts.FormatContent("shipyard_api.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

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
}

func TestServer_PreServerShutdown_Issue2108(t *testing.T) {
	defer discardOutput()()

	gen, err := testAppGenerator(t, "../fixtures/enhancements/2108/swagger.yml", "pre server shutdown")
	require.NoError(t, err)

	app, err := gen.makeCodegenApp()
	require.NoError(t, err)

	// check the serverBuilder output
	buf := bytes.NewBuffer(nil)
	require.NoError(t, templates.MustGet("serverBuilder").Execute(buf, app))

	formatted, err := app.GenOpts.LanguageOpts.FormatContent("shipyard_api.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, `PreServerShutdown:   func() {},`, res)
	assertInCode(t, `PreServerShutdown func()`, res)

	buf = bytes.NewBuffer(nil)
	require.NoError(t, templates.MustGet("serverConfigureapi").Execute(buf, app))

	formatted, err = app.GenOpts.LanguageOpts.FormatContent("configure_shipyard_api.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res = string(formatted)
	// initialisation in New<Name>API function
	assertInCode(t, `api.PreServerShutdown = func() {}`, res)
}

func TestServer_Issue1557(t *testing.T) {
	defer discardOutput()()

	gen, err := testAppGenerator(t, "../fixtures/enhancements/1557/swagger.yml", "generate consumer/producer handlers that are not whitelisted")
	require.NoError(t, err)

	app, err := gen.makeCodegenApp()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, templates.MustGet("serverBuilder").Execute(buf, app))

	formatted, err := app.GenOpts.LanguageOpts.FormatContent("shipyard_api.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, `ApplicationDummyConsumer runtime.Consumer`, res)
	assertInCode(t, `ApplicationDummyProducer runtime.Producer`, res)
	assertInCode(t, `ApplicationDummyConsumer: runtime.ConsumerFunc(func(r io.Reader, target interface{}) error {`, res)
	assertInCode(t, `ApplicationDummyProducer: runtime.ProducerFunc(func(w io.Writer, data interface{}) error {`, res)
	assertInCode(t, `BinConsumer: runtime.ByteStreamConsumer(),`, res)
	assertInCode(t, `BinProducer: runtime.ByteStreamProducer(),`, res)
	assertInCode(t, `result["application/pdf"] = o.BinConsumer`, res)
	assertInCode(t, `result["application/pdf"] = o.BinProducer`, res)
	assertInCode(t, `result["application/dummy"] = o.ApplicationDummyConsumer`, res)
	assertInCode(t, `result["application/dummy"] = o.ApplicationDummyProducer`, res)
}

func TestServer_Issue1648(t *testing.T) {
	defer discardOutput()()

	gen, err := testAppGenerator(t, "../fixtures/bugs/1648/fixture-1648.yaml", "generate format with missing type in model")
	require.NoError(t, err)

	_, err = gen.makeCodegenApp()
	require.NoError(t, err)
}

func TestServer_Issue1746(t *testing.T) {
	defer discardOutput()()

	targetdir, err := ioutil.TempDir(".", "swagger_server")
	require.NoErrorf(t, err, "failed to create a test target directory: %v", err)
	defer func() {
		_ = os.RemoveAll(targetdir)
	}()

	cwd := testCwd(t)
	require.NoErrorf(t, os.Chdir(targetdir), "failed to create a test target directory: %v", err)
	defer func() {
		_ = os.Chdir(cwd)
	}()

	opts := testGenOpts()
	opts.Target = filepath.Join("x")
	opts.Spec = filepath.Join("..", "..", "fixtures", "bugs", "1746", "fixture-1746.yaml")
	tgtSpec := regexp.QuoteMeta(filepath.Join("..", "..", opts.Spec))

	require.NoError(t, os.Mkdir(opts.Target, 0755))

	require.NoError(t, GenerateServer("", nil, nil, opts))

	gulp, err := ioutil.ReadFile(filepath.Join("x", "restapi", "configure_example_swagger_server.go"))
	require.NoError(t, err)

	res := string(gulp)

	tgtPath := regexp.QuoteMeta(filepath.Join("..", "..", opts.Target))
	assertRegexpInCode(t, `go:generate swagger generate server.+\-\-target `+tgtPath, res)
	assertRegexpInCode(t, `go:generate swagger generate server.+\-\-name\s+ExampleSwaggerServer`, res)
	assertRegexpInCode(t, `go:generate swagger generate server.+\-\-spec\s+`+tgtSpec, res)
}

func doGenAppTemplate(t testing.TB, fixture, template string) string {
	gen, err := testAppGenerator(t, fixture, "generate: "+fixture)
	require.NoError(t, err)

	app, err := gen.makeCodegenApp()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, templates.MustGet(template).Execute(buf, app))

	formatted, err := app.GenOpts.LanguageOpts.FormatContent("foo.go", buf.Bytes())
	require.NoError(t, err)

	return string(formatted)
}

func TestServer_Issue1816(t *testing.T) {
	defer discardOutput()()

	// fixed regression: gob encoding in $ref
	res := doGenAppTemplate(t, "../fixtures/bugs/1816/fixture-1816.yaml", "swaggerJsonEmbed")
	assertNotInCode(t, `"$ref": "#"`, res)

	// fixed regression: gob encoding in operation security requirements
	res = doGenAppTemplate(t, "../fixtures/bugs/1824/swagger.json", "swaggerJsonEmbed")
	assertInCode(t, `"api_key": []`, res)
	assertNotInCode(t, `"api_key": null`, res)
}
