package generator_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
	flags "github.com/jessevdk/go-flags"
)

func TestGenerateAndBuild(t *testing.T) {
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	cases := map[string]struct {
		spec string
	}{
		"issue 844": {
			"../fixtures/bugs/844/swagger.json",
		},
		"issue 1216": {
			"../fixtures/bugs/1216/swagger.yml",
		},
	}

	for name, cas := range cases {
		var captureLog bytes.Buffer
		log.SetOutput(&captureLog)

		t.Run(name, func(t *testing.T) {
			spec := filepath.FromSlash(cas.spec)

			generated, err := ioutil.TempDir(filepath.Dir(spec), "generated")
			if err != nil {
				t.Fatalf("TempDir()=%s", generated)
			}
			defer os.RemoveAll(generated)

			err = newTestClient(spec, generated).Execute(nil)
			if err != nil {
				t.Fatalf("Execute()=%s", err)
			}

			//fmt.Println(captureLog.String())
			assert.Contains(t, strings.ToLower(captureLog.String()), "generation completed")

			packages := filepath.Join(generated, "...")

			if p, err := exec.Command("go", "get", packages).CombinedOutput(); err != nil {
				t.Fatalf("go get %s: %s\n%s", packages, err, p)
			}

			if p, err := exec.Command("go", "build", packages).CombinedOutput(); err != nil {
				t.Fatalf("go build %s: %s\n%s", packages, err, p)
			}
		})
	}
}

func newTestClient(input, output string) *generate.Client {
	c := &generate.Client{
		DefaultScheme:   "http",
		DefaultProduces: "application/json",
	}
	c.Spec = flags.Filename(input)
	c.Target = flags.Filename(output)
	c.APIPackage = "operations"
	c.ModelPackage = "models"
	c.ServerPackage = "restapi"
	c.ClientPackage = "client"
	return c
}
