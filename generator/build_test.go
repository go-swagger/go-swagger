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
	"github.com/stretchr/testify/require"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
	flags "github.com/jessevdk/go-flags"
)

const (
	defaultAPIPackage    = "operations"
	defaultClientPackage = "client"
	defaultModelPackage  = "models"
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
		"issue 844 (with params)": {
			"../fixtures/bugs/844/swagger-bis.json",
		},
		"issue 1216": {
			"../fixtures/bugs/1216/swagger.yml",
		},
		"issue 2111": {
			"../fixtures/bugs/2111/fixture-2111.yaml",
		},
		"issue 2278": {
			"../fixtures/bugs/2278/fixture-2278.yaml",
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
			defer func() { _ = os.RemoveAll(generated) }()

			err = newTestClient(spec, generated).Execute(nil)
			require.NoErrorf(t, err, "Execute()=%s", err)

			assert.Contains(t, strings.ToLower(captureLog.String()), "generation completed")

			packages := filepath.Join(generated, "...")

			p, err := exec.Command("go", "get", packages).CombinedOutput()
			require.NoErrorf(t, err, "go get %s: %s\n%s", packages, err, p)

			p, err = exec.Command("go", "build", packages).CombinedOutput()
			require.NoErrorf(t, err, "go build %s: %s\n%s", packages, err, p)
		})
	}
}

func newTestClient(input, output string) *generate.Client {
	c := &generate.Client{}
	c.DefaultScheme = "http"
	c.DefaultProduces = "application/json"
	c.Shared.Spec = flags.Filename(input)
	c.Shared.Target = flags.Filename(output)
	c.Operations.APIPackage = defaultAPIPackage
	c.Models.ModelPackage = defaultModelPackage
	c.ClientPackage = defaultClientPackage
	return c
}
