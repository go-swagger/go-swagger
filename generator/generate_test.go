package generator

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndTest(t *testing.T) {
	if runtime.GOOS == "windows" {
		// don't run race tests on Appveyor CI
		t.SkipNow()
	}
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	cases := map[string]struct {
		spec string
	}{
		"issue 1943": {
			"../fixtures/bugs/1943/fixture-1943.yaml",
		},
	}

	for name, cas := range cases {
		var captureLog bytes.Buffer
		log.SetOutput(&captureLog)
		defer func() {
			if t.Failed() {
				t.Logf("ERROR: generation failed:\n%s", captureLog.String())
			}
		}()

		thisCas := cas

		t.Run(name, func(t *testing.T) {
			spec := filepath.FromSlash(thisCas.spec)

			opts := testGenOpts()

			opts.Target = filepath.Dir(spec)
			_ = os.Mkdir(opts.Target, 0755)
			opts.Spec = spec
			opts.ExcludeSpec = false

			t.Logf("generating test server at: %s", opts.Target)
			err := GenerateServer("", nil, nil, opts)
			defer func() { _ = os.RemoveAll(opts.Target + "/models") }()
			defer func() { _ = os.RemoveAll(opts.Target + "/restapi") }()

			if err != nil {
				if !assert.NoError(t, err, "Execute()=%s", err) {
					return
				}
			}

			packages := filepath.Join(opts.Target, "...")
			testPrg := filepath.Join(opts.Target, "datarace_test.go")

			if p, err := exec.Command("go", "get", packages).CombinedOutput(); err != nil {
				if !assert.NoError(t, err, "go get %s: %s\n%s", packages, err, p) {
					return
				}
			}

			t.Log("running data race test on generated server")
			if p, err := exec.Command("go", "test", "-v", "-race", testPrg).CombinedOutput(); err != nil {
				if !assert.NoError(t, err, "go test -race %s: %s\n%s", packages, err, p) {
					return
				}
			}
		})
	}
}
