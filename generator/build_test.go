package generator_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
	flags "github.com/jessevdk/go-flags"
)

func TestGenerateAndBuild(t *testing.T) {
	cases := map[string]struct {
		spec string
	}{
		"issue 844": {
			"../fixtures/bugs/844/swagger.json",
		},
	}

	for name, cas := range cases {
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
