package generator

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TargetPath and SpecPath are used in server.gotmpl
// as template variables: {{ .TestTargetPath }} and
// {{ .SpecPath }}, to construct the go generate
// directive.
// TODO: there is a catch, since these methods are sensitive
// to the CWD of the current swagger command (or go
// generate when working on resulting template)
// NOTE:
// Errors in TargetPath are hard to simulate since
// they occur only on os.Getwd() errors
// Windows style path is difficult to test on unix
// since the filepath pkg is platform dependant
func TestShared_TargetPath(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	cwd, _ := os.Getwd()

	// relative target
	var opts = new(GenOpts)
	opts.Target = "./a/b/c"
	opts.ServerPackage = "y"
	expected := "../a/b/c"
	result := opts.TargetPath()
	assert.Equal(t, expected, result)

	// relative target, server path
	opts = new(GenOpts)
	opts.Target = "./a/b/c"
	opts.ServerPackage = "y/z"
	expected = "../../a/b/c"
	result = opts.TargetPath()
	assert.Equal(t, expected, result)

	// absolute target
	opts = new(GenOpts)
	opts.Target = filepath.Join(cwd, "a/b/c")
	opts.ServerPackage = "y"
	expected = "../a/b/c"
	result = opts.TargetPath()
	assert.Equal(t, expected, result)

	// absolute target, server path
	opts = new(GenOpts)
	opts.Target = "./a/b/c"
	opts.ServerPackage = "y/z"
	expected = "../../a/b/c"
	result = opts.TargetPath()
	assert.Equal(t, expected, result)

	// absolute server package
	opts = new(GenOpts)
	opts.Target = "/a/b/c"
	opts.ServerPackage = "/y/z"
	expected = "../../a/b/c"
	result = opts.TargetPath()
	assert.Equal(t, expected, result)

	// TODO:
	// unrelated path (Windows specific)
	// TargetPath() is expected to fail
	// when target and server reside on
	// different volumes.
	//if runtime.GOOS == "windows" {
	//	log.Println("INFO:Need some additional testing on windows")
	//opts = new(GenOpts)
	//opts.Target = "C:/a/b/c"
	//opts.ServerPackage = "D:/y/z"
	//expected = ""
	//result = opts.TargetPath()
	//assert.Equal(t, expected, result)
	//}
}

// NOTE: file://url is not supported
func TestShared_SpecPath(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	cwd, _ := os.Getwd()

	// http URL spec
	var opts = new(GenOpts)
	opts.Spec = "http://a/b/c"
	opts.ServerPackage = "y"
	expected := opts.Spec
	result := opts.SpecPath()
	assert.Equal(t, expected, result)

	// https URL spec
	opts = new(GenOpts)
	opts.Spec = "https://a/b/c"
	opts.ServerPackage = "y"
	expected = opts.Spec
	result = opts.SpecPath()
	assert.Equal(t, expected, result)

	// relative spec
	opts = new(GenOpts)
	opts.Spec = "./a/b/c"
	opts.ServerPackage = "y"
	expected = "../a/b/c"
	result = opts.SpecPath()
	assert.Equal(t, expected, result)

	// relative spec, server path
	opts = new(GenOpts)
	opts.Spec = "./a/b/c"
	opts.ServerPackage = "y/z"
	expected = "../../a/b/c"
	result = opts.SpecPath()
	assert.Equal(t, expected, result)

	// absolute spec
	opts = new(GenOpts)
	opts.Spec = filepath.Join(cwd, "a/b/c")
	opts.ServerPackage = "y"
	expected = "../a/b/c"
	result = opts.SpecPath()
	assert.Equal(t, expected, result)

	// absolute spec, server path
	opts = new(GenOpts)
	opts.Spec = "./a/b/c"
	opts.ServerPackage = "y/z"
	expected = "../../a/b/c"
	result = opts.SpecPath()
	assert.Equal(t, expected, result)

	// absolute server package
	opts = new(GenOpts)
	opts.Spec = "/a/b/c"
	opts.ServerPackage = "/y/z"
	expected = "../../a/b/c"
	result = opts.SpecPath()
	assert.Equal(t, expected, result)

	// TODO:
	// unrelated path (Windows specific)
	// SpecPath() is expected to fail
	// when spec and server reside on
	// different volumes.
	// This shouldn't happen with go, since everything
	// should be under the same GOPATH
	// NOTE: a better error handling strategy would be
	// to check that os.Getwd() and Rel() work well upstream
	// and assume in functions that no error is returned.
	//if runtime.GOOS == "windows" {
	//	log.SetOutput(os.Stdout)
	//	log.Println("INFO:Need some additional testing on windows")
	//opts = new(GenOpts)
	//opts.Spec = "C:/a/b/c"
	//opts.ServerPackage = "D:/y/z"
	//expected = ""
	//result = opts.SpecPath()
	//assert.Equal(t, expected, result)
	//}
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

	templates.AddFile("garbled", garbled)
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

	templates.AddFile("execfailure1", execfailure1)
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

	templates.AddFile("execfailure2", execfailure2)
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
		os.Remove("test_badformat.gol")
		os.Remove("test_badformat2.gol")
		log.SetOutput(os.Stdout)
		Debug = false
	}()

	// Not skipping format
	badFormat := "func x {;;; garbled"

	Debug = true
	templates.AddFile("badformat", badFormat)

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
	os.Remove(tplOpts.FileName)

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
	os.Remove(tplOpts2.FileName)

	assert.Nil(t, err2)

	// os.RemoveAll(filepath.Join(filepath.FromSlash(dr),"restapi"))
}

// Test dir creation
func TestShared_DirectoryTemplate(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	defer func() {
		os.RemoveAll("TestGenDir")
		log.SetOutput(os.Stdout)
		Debug = false
	}()

	// Not skipping format
	content := "func x {}"

	Debug = true
	templates.AddFile("gendir", content)

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
	os.RemoveAll(tplOpts.Target)

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
