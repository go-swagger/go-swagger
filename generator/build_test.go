package generator_test

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

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
	// This test generates and actually compiles the output
	// of generated clients.
	//
	// We run this in parallel now. Therefore it is no more
	// possible to assert the output on stdout.
	//
	// NOTE: test cases are randomized (map)
	t.Parallel()

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
		"issue 2163": {
			"../fixtures/enhancements/2163/fixture-2163.yaml",
		},
		"issue 1771": {
			"../fixtures/enhancements/1771/fixture-1771.yaml",
		},
	}

	t.Run("build client", func(t *testing.T) {
		for name, cas := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				log.SetOutput(ioutil.Discard)

				spec := filepath.FromSlash(cas.spec)

				generated, err := ioutil.TempDir(filepath.Dir(spec), "generated")
				require.NoErrorf(t, err, "TempDir()=%s", generated)
				defer func() { _ = os.RemoveAll(generated) }()

				require.NoErrorf(t, newTestClient(spec, generated).Execute(nil), "Execute()=%s", err)

				packages := filepath.Join(generated, "...")

				goExecInDir(t, "", "get")

				goExecInDir(t, "", "build", packages)
			})
		}
	})
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

func goExecInDir(t testing.TB, target string, args ...string) {
	cmd := exec.Command("go", args...)
	cmd.Dir = target
	p, err := cmd.CombinedOutput()
	require.NoErrorf(t, err, "unexpected error: %s: %v\n%s", cmd.String(), err, string(p))
}
