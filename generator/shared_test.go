package generator

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
)

const testPath = "a/b/c"

// TODO: there is a catch, since these methods are sensitive
// to the CWD of the current swagger command (or go
// generate when working on resulting template)
// NOTE:
// Errors in CheckOpts are hard to simulate since
// they occur only on os.Getwd() errors
// Windows style path is difficult to test on unix
// since the filepath pkg is platform dependant
func TestShared_CheckOpts(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	var opts = new(GenOpts)
	_ = opts.EnsureDefaults()
	cwd, _ := os.Getwd()

	opts.Target = filepath.Join(".", "a", "b", "c")
	opts.ServerPackage = filepath.Join(cwd, "a", "b", "c")
	err := opts.CheckOpts()
	assert.Error(t, err)

	opts.Target = filepath.Join(cwd, "a", "b", "c")
	opts.ServerPackage = testPath
	err = opts.CheckOpts()
	assert.NoError(t, err)

	opts.Target = filepath.Join(cwd, "a", "b", "c")
	opts.ServerPackage = testPath
	opts.Spec = "https:/ab/c"
	err = opts.CheckOpts()
	assert.NoError(t, err)

	opts.Target = filepath.Join(cwd, "a", "b", "c")
	opts.ServerPackage = testPath
	opts.Spec = "http:/ab/c"
	err = opts.CheckOpts()
	assert.NoError(t, err)

	opts.Target = filepath.Join("a", "b", "c")
	opts.ServerPackage = testPath
	opts.Spec = filepath.Join(cwd, "x")
	err = opts.CheckOpts()
	assert.NoError(t, err)
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

	opts := GenOpts{}
	tplOpts := TemplateOpts{
		Name:       "NotFound",
		Source:     "asset:notfound",
		Target:     ".",
		FileName:   "test_notfound.go",
		SkipExists: false,
		SkipFormat: false,
	}

	buf, err := opts.render(&tplOpts, nil)
	assert.Error(t, err, "Error should be handled here")
	assert.Nil(t, buf, "Upon error, GenOpts.render() should return nil buffer")
}

// Low level testing: invalid template => Get() returns not found (higher level calls raise panic(), see above)
// TODO: better error discrimination between absent definition and non-parsing template
func TestShared_GarbledTemplate(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	garbled := "func x {{;;; garbled"

	_ = templates.AddFile("garbled", garbled)
	opts := GenOpts{}
	tplOpts := TemplateOpts{
		Name:       "Garbled",
		Source:     "asset:garbled",
		Target:     ".",
		FileName:   "test_garbled.go",
		SkipExists: false,
		SkipFormat: false,
	}

	buf, err := opts.render(&tplOpts, nil)
	assert.Error(t, err, "Error should be handled here")
	assert.Nil(t, buf, "Upon error, GenOpts.render() should return nil buffer")
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
	opts := new(GenOpts)
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
	opts = new(GenOpts)
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

	opts := GenOpts{}
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

	// The badly formatted file has been dumped for debugging purposes
	_, exists := os.Stat(tplOpts.FileName)
	assert.True(t, !os.IsNotExist(exists), "The template file has not been generated as expected")
	_ = os.Remove(tplOpts.FileName)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "source formatting on generated source")

	// Skipping format
	opts = GenOpts{}
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

	opts := GenOpts{}
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

	opts := GenOpts{}
	tplOpts := TemplateOpts{
		Name:       "File",
		Source:     "File",
		Target:     ".",
		FileName:   "file.go",
		SkipExists: false,
		SkipFormat: false,
	}

	buf, err := opts.render(&tplOpts, nil)
	//spew.Dump(err)
	assert.Error(t, err, "Error should be handled here")
	assert.Contains(t, err.Error(), "open File")
	assert.Contains(t, err.Error(), "error while opening")
	assert.Nil(t, buf, "Upon error, GenOpts.render() should return nil buffer")

	opts.TemplateDir = filepath.Join(".", "myTemplateDir")
	buf, err = opts.render(&tplOpts, nil)
	//spew.Dump(err)
	assert.Error(t, err, "Error should be handled here")
	assert.Contains(t, err.Error(), "open "+filepath.Join("myTemplateDir", "File"))
	assert.Contains(t, err.Error(), "error while opening")
	assert.Nil(t, buf, "Upon error, GenOpts.render() should return nil buffer")

}

func TestShared_Issue1429(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	// acknowledge fix in go-openapi/spec
	specPath := filepath.Join("..", "fixtures", "bugs", "1429", "swagger-1429.yaml")
	specDoc, err := loads.Spec(specPath)
	assert.NoError(t, err)

	opts := testGenOpts()
	opts.Spec = specPath
	_, err = validateAndFlattenSpec(&opts, specDoc)
	assert.NoError(t, err)

	// more aggressive fixture on $refs, with validation errors, but flatten ok
	specPath = filepath.Join("..", "fixtures", "bugs", "1429", "swagger.yaml")
	specDoc, err = loads.Spec(specPath)
	assert.NoError(t, err)

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

func TestShared_MangleFileName(t *testing.T) {
	// standard : swag.ToFileName()
	o := LanguageOpts{}
	o.Init()
	res := o.MangleFileName("aFileEndingInOsNameWindows")
	assert.True(t, strings.HasSuffix(res, "_windows"))

	// golang specific
	res = golang.MangleFileName("aFileEndingInOsNameWindows")
	assert.True(t, strings.HasSuffix(res, "_windows_swagger"))
	res = golang.MangleFileName("aFileEndingInOsNameWindowsAmd64")
	assert.True(t, strings.HasSuffix(res, "_windows_amd64_swagger"))
	res = golang.MangleFileName("aFileEndingInTest")
	assert.True(t, strings.HasSuffix(res, "_test_swagger"))
}

func TestShared_ManglePackage(t *testing.T) {
	o := GoLangOpts()
	o.Init()

	for _, v := range []struct {
		tested       string
		expectedPath string
		expectedName string
	}{
		{tested: "", expectedPath: "default", expectedName: "default"},
		{tested: "select", expectedPath: "select_default", expectedName: "select_default"},
		{tested: "x", expectedPath: "x", expectedName: "x"},
		{tested: "a/b/c-d/e_f/g", expectedPath: "a/b/c-d/e_f/g", expectedName: "g"},
		{tested: "a/b/c-d/e_f/g-h", expectedPath: "a/b/c-d/e_f/g_h", expectedName: "g_h"},
	} {
		res := o.ManglePackagePath(v.tested, "default")
		assert.Equal(t, v.expectedPath, res)
		res = o.ManglePackageName(v.tested, "default")
		assert.Equal(t, v.expectedName, res)
	}
}

func TestShared_Issue1621(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	// acknowledge fix in go-openapi/spec
	specPath := filepath.Join("..", "fixtures", "bugs", "1621", "fixture-1621.yaml")
	specDoc, err := loads.Spec(specPath)
	assert.NoError(t, err)

	opts := testGenOpts()
	opts.Spec = specPath
	opts.ValidateSpec = true
	//t.Logf("path: %s", specDoc.SpecFilePath())
	_, err = validateAndFlattenSpec(&opts, specDoc)
	assert.NoError(t, err)
}

func TestShared_Issue1614(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	// acknowledge fix in go-openapi/spec
	specPath := filepath.Join("..", "fixtures", "bugs", "1614", "gitea.json")
	specDoc, err := loads.Spec(specPath)
	assert.NoError(t, err)

	opts := testGenOpts()
	opts.Spec = specPath
	opts.ValidateSpec = true
	t.Logf("path: %s", specDoc.SpecFilePath())
	_, err = validateAndFlattenSpec(&opts, specDoc)
	assert.NoError(t, err)
}
