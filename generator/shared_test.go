package generator

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultAPIPackage    = "operations"
	defaultClientPackage = "client"
	defaultModelPackage  = "models"
	defaultServerPackage = "restapi"
)

// Perform common initialization of template repository before running tests.
// This allows to run tests unitarily (e.g. go test -run xxx ).
func TestMain(m *testing.M) {
	// initializations to run tests in this package
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	templates.LoadDefaults()
	initSchemaValidationTest()
	os.Exit(m.Run())
}

func opts() *GenOpts {
	var opts GenOpts
	opts.IncludeValidator = true
	opts.IncludeModel = true
	if err := opts.EnsureDefaults(); err != nil {
		panic(err)
	}
	return &opts
}

func testGenOpts() *GenOpts {
	g := &GenOpts{}
	g.Target = "."
	g.APIPackage = defaultAPIPackage
	g.ModelPackage = defaultModelPackage
	g.ServerPackage = defaultServerPackage
	g.ClientPackage = defaultClientPackage
	g.Principal = ""
	g.DefaultScheme = "http"
	g.IncludeModel = true
	g.IncludeValidator = true
	g.IncludeModel = true
	g.IncludeHandler = true
	g.IncludeParameters = true
	g.IncludeResponses = true
	g.IncludeMain = false
	g.IncludeSupport = true
	g.ExcludeSpec = true
	g.TemplateDir = ""
	g.DumpData = false

	if err := g.EnsureDefaults(); err != nil {
		panic(err)
	}
	return g
}

// TODO: there is a catch, since these methods are sensitive
// to the CWD of the current swagger command (or go
// generate when working on resulting template)
// NOTE:
// Errors in CheckOpts are hard to simulate since
// they occur only on os.Getwd() errors
// Windows style path is difficult to test on unix
// since the filepath pkg is platform dependent
func TestShared_CheckOpts(t *testing.T) {
	testPath := filepath.Join("a", "b", "b")

	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	var opts = new(GenOpts)
	_ = opts.EnsureDefaults()
	cwd, _ := os.Getwd()
	opts.Spec = "../fixtures/codegen/simplesearch.yml"

	opts.Target = filepath.Join(".", "a", "b", "c")
	opts.ServerPackage = filepath.Join(cwd, "a", "b", "c")
	err := opts.CheckOpts()
	assert.Error(t, err)

	opts.Target = filepath.Join(cwd, "a", "b", "c")
	opts.ServerPackage = testPath
	opts.Spec = filepath.Join(cwd, "nowhere", "swagger.yaml")
	err = opts.CheckOpts()
	assert.Error(t, err)

	opts.Target = filepath.Join(cwd, "a", "b", "c")
	opts.ServerPackage = testPath
	opts.Spec = "https://ab/c"
	err = opts.CheckOpts()
	assert.NoError(t, err)

	opts.Target = filepath.Join(cwd, "a", "b", "c")
	opts.ServerPackage = testPath
	opts.Spec = "http://ab/c"
	err = opts.CheckOpts()
	assert.NoError(t, err)

	opts.Target = filepath.Join("a", "b", "c")
	opts.ServerPackage = testPath
	opts.Spec = filepath.Join(cwd, "..", "fixtures", "codegen", "swagger-codegen-tests.json")
	err = opts.CheckOpts()
	assert.NoError(t, err)

	opts.Target = filepath.Join("a", "b", "c")
	opts.ServerPackage = testPath
	opts.Spec = filepath.Join("..", "fixtures", "codegen", "swagger-codegen-tests.json")
	err = opts.CheckOpts()
	assert.NoError(t, err)

	opts = nil
	err = opts.CheckOpts()
	assert.Error(t, err)
}

func TestShared_EnsureDefaults(t *testing.T) {
	opts := &GenOpts{}
	require.NoError(t, opts.EnsureDefaults())
	assert.True(t, opts.defaultsEnsured)
	opts.DefaultConsumes = "https"
	_ = opts.EnsureDefaults()
	assert.Equal(t, "https", opts.DefaultConsumes)
}

// TargetPath and SpecPath are used in server.gotmpl
// as template variables: {{ .TestTargetPath }} and
// {{ .SpecPath }}, to construct the go generate
// directive.
func TestShared_TargetPath(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	cwd, _ := os.Getwd()

	// relative target
	var opts = new(GenOpts)
	_ = opts.EnsureDefaults()
	opts.Target = filepath.Join(".", "a", "b", "c")
	opts.ServerPackage = "y"
	expected := filepath.Join("..", "..", "c")
	result := opts.TargetPath()
	assert.Equal(t, expected, result)

	// relative target, server path
	opts = new(GenOpts)
	_ = opts.EnsureDefaults()
	opts.Target = filepath.Join(".", "a", "b", "c")
	opts.ServerPackage = "y/z"
	expected = filepath.Join("..", "..", "..", "c")
	result = opts.TargetPath()
	assert.Equal(t, expected, result)

	// absolute target
	opts = new(GenOpts)
	_ = opts.EnsureDefaults()
	opts.Target = filepath.Join(cwd, "a", "b", "c")
	opts.ServerPackage = "y"
	expected = filepath.Join("..", "..", "c")
	result = opts.TargetPath()
	assert.Equal(t, expected, result)

	// absolute target, server path
	opts = new(GenOpts)
	_ = opts.EnsureDefaults()
	opts.Target = filepath.Join(cwd, "a", "b", "c")
	opts.ServerPackage = path.Join("y", "z")
	expected = filepath.Join("..", "..", "..", "c")
	result = opts.TargetPath()
	assert.Equal(t, expected, result)
}

// NOTE: file://url is not supported
func TestShared_SpecPath(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	cwd, _ := os.Getwd()

	// http URL spec
	var opts = new(GenOpts)
	_ = opts.EnsureDefaults()
	opts.Spec = "http://a/b/c"
	opts.ServerPackage = "y"
	expected := opts.Spec
	result := opts.SpecPath()
	assert.Equal(t, expected, result)

	// https URL spec
	opts = new(GenOpts)
	_ = opts.EnsureDefaults()
	opts.Spec = "https://a/b/c"
	opts.ServerPackage = "y"
	expected = opts.Spec
	result = opts.SpecPath()
	assert.Equal(t, expected, result)

	// relative spec
	opts = new(GenOpts)
	_ = opts.EnsureDefaults()
	opts.Spec = filepath.Join(".", "a", "b", "c")
	opts.Target = filepath.Join("d")
	opts.ServerPackage = "y"
	expected = filepath.Join("..", "..", "a", "b", "c")
	result = opts.SpecPath()
	assert.Equal(t, expected, result)

	// relative spec, server path
	opts = new(GenOpts)
	_ = opts.EnsureDefaults()
	opts.Spec = filepath.Join(".", "a", "b", "c")
	opts.Target = filepath.Join("d", "e")
	opts.ServerPackage = "y/z"
	expected = filepath.Join("..", "..", "..", "..", "a", "b", "c")
	result = opts.SpecPath()
	assert.Equal(t, expected, result)

	// relative spec, server path
	opts = new(GenOpts)
	_ = opts.EnsureDefaults()
	opts.Spec = filepath.Join(".", "a", "b", "c")
	opts.Target = filepath.Join(".", "a", "b")
	opts.ServerPackage = "y/z"
	expected = filepath.Join("..", "..", "c")
	result = opts.SpecPath()
	assert.Equal(t, expected, result)

	// absolute spec
	opts = new(GenOpts)
	_ = opts.EnsureDefaults()
	opts.Spec = filepath.Join(cwd, "a", "b", "c")
	opts.ServerPackage = "y"
	expected = filepath.Join("..", "a", "b", "c")
	result = opts.SpecPath()
	assert.Equal(t, expected, result)

	// absolute spec, server path
	opts = new(GenOpts)
	_ = opts.EnsureDefaults()
	opts.Spec = filepath.Join("..", "a", "b", "c")
	opts.Target = ""
	opts.ServerPackage = path.Join("y", "z")
	expected = filepath.Join("..", "..", "..", "a", "b", "c")
	result = opts.SpecPath()
	assert.Equal(t, expected, result)

	if runtime.GOOS == "windows" {
		opts = new(GenOpts)
		_ = opts.EnsureDefaults()
		opts.Spec = filepath.Join("a", "b", "c")
		opts.Target = filepath.Join("Z:", "e", "f", "f")
		opts.ServerPackage = "y/z"
		expected, _ = filepath.Abs(opts.Spec)
		result = opts.SpecPath()
		assert.Equal(t, expected, result)
	}
}

// Low level testing: templates not found (higher level calls raise panic(), see above)
func TestShared_NotFoundTemplate(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	opts := testGenOpts()
	tplOpts := TemplateOpts{
		Name:       "NotFound",
		Source:     "asset:notfound",
		Target:     ".",
		FileName:   "test_notfound.go",
		SkipExists: false,
		SkipFormat: false,
	}

	buf, err := opts.render(&tplOpts, nil)
	require.Errorf(t, err, "Error should be handled here")
	assert.Nilf(t, buf, "Upon error, GenOpts.render() should return nil buffer")
}

// Low level testing: invalid template => Get() returns not found (higher level calls raise panic(), see above)
// TODO: better error discrimination between absent definition and non-parsing template
func TestShared_GarbledTemplate(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	garbled := "func x {{;;; garbled"

	_ = templates.AddFile("garbled", garbled)
	opts := testGenOpts()

	tplOpts := TemplateOpts{
		Name:       "Garbled",
		Source:     "asset:garbled",
		Target:     ".",
		FileName:   "test_garbled.go",
		SkipExists: false,
		SkipFormat: false,
	}

	buf, err := opts.render(&tplOpts, nil)
	require.Errorf(t, err, "Error should be handled here")
	assert.Nilf(t, buf, "Upon error, GenOpts.render() should return nil buffer")
}

// Template execution failure
type myTemplateData struct {
}

func (*myTemplateData) MyFaultyMethod() (string, error) {
	return "", fmt.Errorf("myFaultyError")
}

func TestShared_ExecTemplate(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	// Not a failure: no value data
	execfailure1 := "func x {{ .NotInData }}"

	_ = templates.AddFile("execfailure1", execfailure1)
	opts := testGenOpts()

	tplOpts := TemplateOpts{
		Name:       "execFailure1",
		Source:     "asset:execfailure1",
		Target:     ".",
		FileName:   "test_execfailure1.go",
		SkipExists: false,
		SkipFormat: false,
	}

	buf1, err := opts.render(&tplOpts, nil)
	assert.NoError(t, err, "Template rendering should put <no value> instead of missing data, and report no error")
	assert.Equal(t, string(buf1), "func x <no value>")

	execfailure2 := "func {{ .MyFaultyMethod }}"

	_ = templates.AddFile("execfailure2", execfailure2)
	opts = testGenOpts()
	tplOpts2 := TemplateOpts{
		Name:       "execFailure2",
		Source:     "asset:execfailure2",
		Target:     ".",
		FileName:   "test_execfailure2.go",
		SkipExists: false,
		SkipFormat: false,
	}

	data := new(myTemplateData)
	buf2, err := opts.render(&tplOpts2, data)
	assert.Error(t, err, "Error should be handled here: missing func in template yields an error")
	assert.Contains(t, err.Error(), "template execution failed")
	assert.Nil(t, buf2, "Upon error, GenOpts.render() should return nil buffer")
}

// Test correctly parsed templates, with bad formatting
func TestShared_BadFormatTemplate(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	defer func() {
		_ = os.Remove("test_badformat.gol")
		_ = os.Remove("test_badformat2.gol")
		log.SetOutput(os.Stdout)
		Debug = false
	}()

	// Not skipping format
	badFormat := "func x {;;; garbled"

	Debug = true
	_ = templates.AddFile("badformat", badFormat)

	opts := testGenOpts()
	opts.LanguageOpts = GoLangOpts()
	tplOpts := TemplateOpts{
		Name:   "badformat",
		Source: "asset:badformat",
		Target: ".",
		// Extension ".gol" won't mess with go if cleanup is not performed
		FileName:   "test_badformat.gol",
		SkipExists: false,
		SkipFormat: false,
	}

	data := appGenerator{
		Name:    "badtest",
		Package: "wrongpkg",
	}

	err := opts.write(&tplOpts, data)
	defer func() {
		_ = os.Remove(tplOpts.FileName)
	}()

	// The badly formatted file has been dumped for debugging purposes
	_, exists := os.Stat(tplOpts.FileName)
	assert.True(t, !os.IsNotExist(exists), "The template file has not been generated as expected")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "source formatting on generated source")

	// Skipping format
	opts = testGenOpts()
	opts.LanguageOpts = GoLangOpts()
	tplOpts2 := TemplateOpts{
		Name:       "badformat2",
		Source:     "asset:badformat",
		Target:     ".",
		FileName:   "test_badformat2.gol",
		SkipExists: false,
		SkipFormat: true,
	}

	err2 := opts.write(&tplOpts2, data)

	// The unformatted file has been dumped without format checks
	_, exists2 := os.Stat(tplOpts2.FileName)
	assert.True(t, !os.IsNotExist(exists2), "The template file has not been generated as expected")
	_ = os.Remove(tplOpts2.FileName)

	assert.Nil(t, err2)

	// os.RemoveAll(filepath.Join(filepath.FromSlash(dr),"restapi"))
}

// Test dir creation
func TestShared_DirectoryTemplate(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	defer func() {
		_ = os.RemoveAll("TestGenDir")
		log.SetOutput(os.Stdout)
	}()

	// Not skipping format
	content := "func x {}"

	_ = templates.AddFile("gendir", content)

	opts := testGenOpts()
	opts.LanguageOpts = GoLangOpts()
	tplOpts := TemplateOpts{
		Name:   "gendir",
		Source: "asset:gendir",
		Target: "TestGenDir",
		// Extension ".gol" won't mess with go if cleanup is not performed
		FileName:   "test_gendir.gol",
		SkipExists: false,
		SkipFormat: true,
	}

	data := appGenerator{
		Name:    "gentest",
		Package: "stubpkg",
	}

	err := opts.write(&tplOpts, data)

	// The badly formatted file has been dumped for debugging purposes
	_, exists := os.Stat(filepath.Join(tplOpts.Target, tplOpts.FileName))
	assert.True(t, !os.IsNotExist(exists), "The template file has not been generated as expected")
	_ = os.RemoveAll(tplOpts.Target)

	assert.Nil(t, err)
}

// Test templates which are not assets (open in file)
// Low level testing: templates loaded from file
func TestShared_LoadTemplate(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	opts := testGenOpts()
	tplOpts := TemplateOpts{
		Name:       "File",
		Source:     "File",
		Target:     ".",
		FileName:   "file.go",
		SkipExists: false,
		SkipFormat: false,
	}

	buf, err := opts.render(&tplOpts, nil)
	assert.Error(t, err, "Error should be handled here")
	assert.Contains(t, err.Error(), "open File")
	assert.Contains(t, err.Error(), "error while opening")
	assert.Nil(t, buf, "Upon error, GenOpts.render() should return nil buffer")

	opts.TemplateDir = filepath.Join(".", "myTemplateDir")
	buf, err = opts.render(&tplOpts, nil)
	assert.Error(t, err, "Error should be handled here")
	assert.Contains(t, err.Error(), "open "+filepath.Join("myTemplateDir", "File"))
	assert.Contains(t, err.Error(), "error while opening")
	assert.Nil(t, buf, "Upon error, GenOpts.render() should return nil buffer")

}

func TestShared_AppNameOrDefault(t *testing.T) {
	specPath := filepath.Join("..", "fixtures", "codegen", "shipyard.yml")
	specDoc, err := loads.Spec(specPath)
	require.NoError(t, err)

	require.NotNil(t, specDoc.Spec().Info)
	specDoc.Spec().Info.Title = "    "
	assert.Equal(t, "Xyz", appNameOrDefault(specDoc, "  ", "xyz"))
	specDoc.Spec().Info.Title = "test"
	assert.Equal(t, "Xyz", appNameOrDefault(specDoc, "  ", "xyz"))

	opts := testGenOpts()
	opts.Spec = specPath
	_, err = opts.validateAndFlattenSpec()
	require.NoError(t, err)

	// more aggressive fixture on $refs, with validation errors, but flatten ok
	specPath = filepath.Join("..", "fixtures", "bugs", "1429", "swagger.yaml")
	specDoc, err = loads.Spec(specPath)
	require.NoError(t, err)

	opts.Spec = specPath
	opts.FlattenOpts.BasePath = specDoc.SpecFilePath()
	opts.FlattenOpts.Spec = analysis.New(specDoc.Spec())
	opts.FlattenOpts.Minimal = true
	err = analysis.Flatten(*opts.FlattenOpts)
	assert.NoError(t, err)

	specDoc, _ = loads.Spec(specPath) // needs reload
	opts.FlattenOpts.Spec = analysis.New(specDoc.Spec())
	opts.FlattenOpts.Minimal = false
	err = analysis.Flatten(*opts.FlattenOpts)
	assert.NoError(t, err)
}

func TestShared_GatherModel(t *testing.T) {
	specPath := filepath.Join("..", "fixtures", "codegen", "shipyard.yml")

	specDoc, err := loads.Spec(specPath)
	require.NoError(t, err)

	_, err = gatherModels(specDoc, []string{"unknown"})
	assert.Error(t, err)

	res, err := gatherModels(specDoc, []string{"Image", "Application"})
	require.NoError(t, err)
	assert.Len(t, res, 2)

	res, err = gatherModels(specDoc, []string{"Image", "Application"})
	require.NoError(t, err)
	assert.Len(t, res, 2)

	res, err = gatherModels(specDoc, []string{})
	require.NoError(t, err)
	assert.Len(t, res, 4)
}

func TestShared_DumpWrongData(t *testing.T) {
	assert.Error(t, dumpData(struct {
		A func() string
		B string
	}{
		A: func() string { return "" },
		B: "xyz",
	}))

	assert.NoError(t, dumpData(struct {
		A func() string `json:"-"`
		B string
	}{
		A: func() string { return "" },
		B: "xyz",
	}))

	assert.NoError(t, dumpData(struct {
		a func() string
		B string
	}{
		a: func() string { return "" },
		B: "xyz",
	}))
}

func TestResolvePrincipal(t *testing.T) {
	for _, toPin := range []struct {
		Title     string
		Principal string
		Expected  []string
	}{
		{
			Title: "defaults", Principal: "",
			Expected: []string{"", "interface{}", ""},
		},
		{
			Title: "with base import", Principal: "auth.Principal",
			Expected: []string{"auth", "auth.Principal", "auth"},
		},
		{
			Title: "with full import", Principal: "github.com/myproject/auth.Principal",
			Expected: []string{"auth", "auth.Principal", "github.com/myproject/auth"},
		},
		{
			Title: "with name conflict", Principal: "github.com/myproject/middleware.Principal",
			Expected: []string{"auth", "auth.Principal", "github.com/myproject/middleware"},
		},
		{
			Title: "with name conflict (2)", Principal: "github.com/myproject/principal.Principal",
			Expected: []string{"auth", "auth.Principal", "github.com/myproject/principal"},
		},
	} {
		fixture := toPin
		t.Run(fixture.Title, func(t *testing.T) {
			t.Parallel()
			opts := &GenOpts{Principal: fixture.Principal}
			err := opts.EnsureDefaults()
			require.NoError(t, err)
			alias, principal, target := opts.resolvePrincipal()
			require.Equal(t, fixture.Expected[0], alias)
			require.Equal(t, fixture.Expected[1], principal)
			require.Equal(t, fixture.Expected[2], target)
		})
	}
}

func TestDefaultImports(t *testing.T) {
	for i, toPin := range []struct {
		Title    string
		Opts     *GenOpts
		Expected map[string]string
	}{
		{
			Title: "defaults",
			Opts:  &GenOpts{},
			Expected: map[string]string{
				"models": "github.com/go-swagger/go-swagger/generator/models",
			},
		},
		{
			Title: "with base import",
			Opts: &GenOpts{
				Principal: "ext.Principal",
			},
			Expected: map[string]string{
				"ext":    "github.com/go-swagger/go-swagger/generator/ext",
				"models": "github.com/go-swagger/go-swagger/generator/models",
			},
		},
		{
			Title: "with full import",
			Opts: &GenOpts{
				Principal: "github.com/myproject/identity.Principal",
			},
			Expected: map[string]string{
				"identity": "github.com/myproject/identity",
				"models":   "github.com/go-swagger/go-swagger/generator/models",
			},
		},
		{
			Title: "with name conflict",
			Opts: &GenOpts{
				Principal: "github.com/myproject/middleware.Principal",
			},
			Expected: map[string]string{
				"auth":   "github.com/myproject/middleware",
				"models": "github.com/go-swagger/go-swagger/generator/models",
			},
		},
		{
			Title: "with name conflict (2)",
			Opts: &GenOpts{
				Principal: "github.com/myproject/principal.Principal",
			},
			Expected: map[string]string{
				"auth":   "github.com/myproject/principal",
				"models": "github.com/go-swagger/go-swagger/generator/models",
			},
		},
		{
			Title: "alternate target for models",
			Opts: &GenOpts{
				ModelPackage: "target/bespoke",
			},
			Expected: map[string]string{
				"bespoke": "github.com/go-swagger/go-swagger/generator/target/bespoke",
			},
		},
		{
			Title: "with existing models",
			Opts: &GenOpts{
				ExistingModels: "github.com/myproject/target/bespoke",
			},
			Expected: map[string]string{
				"models": "github.com/myproject/target/bespoke",
			},
		},
		// issue #2362
		{
			Title: "relative principal, in dedicated package under generated target",
			Opts: &GenOpts{
				Principal:    "auth.Principal",
				ModelPackage: "target/bespoke",
			},
			Expected: map[string]string{
				"bespoke": "github.com/go-swagger/go-swagger/generator/target/bespoke",
				"auth":    "github.com/go-swagger/go-swagger/generator/auth",
			},
		},
		{
			Title: "relative principal in models (1)",
			Opts: &GenOpts{
				Principal:    "bespoke.Principal",
				ModelPackage: "target/bespoke",
			},
			Expected: map[string]string{
				"bespoke": "github.com/go-swagger/go-swagger/generator/target/bespoke",
			},
		},
		{
			Title: "relative principal in models (2)",
			Opts: &GenOpts{
				Principal:    "target/bespoke.Principal",
				ModelPackage: "target/bespoke",
			},
			Expected: map[string]string{
				"bespoke": "github.com/go-swagger/go-swagger/generator/target/bespoke",
			},
		},
		{
			Title: "relative principal: not detected",
			// NOTE: this case will probably not build: no way to determine the user intent
			Opts: &GenOpts{
				Principal:    "target/auth.Principal",
				ModelPackage: "target/models",
			},
			Expected: map[string]string{
				"models": "github.com/go-swagger/go-swagger/generator/target/models",
				"auth":   "target/auth",
			},
		},
	} {
		fixture := toPin
		t.Run(fixture.Title, func(t *testing.T) {
			t.Parallel()
			err := fixture.Opts.EnsureDefaults()
			require.NoError(t, err)
			imports := fixture.Opts.defaultImports()
			require.EqualValuesf(t, fixture.Expected, imports, "unexpected imports generated with fixture %q[%d]", fixture.Title, i)
		})
	}
}

func TestShared_Issue2113(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	// acknowledge fix in go-openapi/spec
	specPath := filepath.Join("..", "fixtures", "bugs", "2113", "base.yaml")
	_, err := loads.Spec(specPath)
	require.NoError(t, err)

	opts := testGenOpts()
	opts.Spec = specPath
	opts.Spec = specPath
	opts.ValidateSpec = true
	_, err = opts.validateAndFlattenSpec()
	assert.NoError(t, err)
}
