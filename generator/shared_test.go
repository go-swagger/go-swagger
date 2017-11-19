package generator

import (
        "fmt"
        // "io/ioutil"
        "runtime"
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
func TestTargetPath(t *testing.T) {
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
  opts.Target = filepath.Join(cwd,"a/b/c")
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
  if runtime.GOOS == "windows" {
    fmt.Println("INFO:Need some additional testing on windows")
  //opts = new(GenOpts)
  //opts.Target = "C:/a/b/c"
  //opts.ServerPackage = "D:/y/z"
  //expected = ""
  //result = opts.TargetPath()
  //assert.Equal(t, expected, result)
  }
}

// NOTE: file://url is not supported
func TestSpecPath(t *testing.T) {
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
  opts.Spec = filepath.Join(cwd,"a/b/c")
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
  if runtime.GOOS == "windows" {
    fmt.Println("INFO:Need some additional testing on windows")
  //opts = new(GenOpts)
  //opts.Spec = "C:/a/b/c"
  //opts.ServerPackage = "D:/y/z"
  //expected = ""
  //result = opts.SpecPath()
  //assert.Equal(t, expected, result)
  }
}
