package generator

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGenerateAndTest(t *testing.T) {
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

		thisCas := cas

		t.Run(name, func(t *testing.T) {
			spec := filepath.FromSlash(thisCas.spec)

			opts := testGenOpts()

			opts.Target = filepath.Dir(spec)
			_ = os.Mkdir(opts.Target, 0755)
			opts.Spec = spec
			opts.ExcludeSpec = false

			fmt.Println(opts.Target)

			err := GenerateServer("", nil, nil, &opts)
			defer func() { _ = os.RemoveAll(opts.Target + "/models") }()
			defer func() { _ = os.RemoveAll(opts.Target + "/restapi") }()

			if err != nil {
				t.Fatalf("Execute()=%s", err)
			}

			//fmt.Println(captureLog.String())
			//assert.Contains(t, strings.ToLower(captureLog.String()), "creating generated file")
			packages := filepath.Join(opts.Target, "...")

			if p, err := exec.Command("go", "get", packages).CombinedOutput(); err != nil {
				t.Fatalf("go get %s: %s\n%s", packages, err, p)
			}

			if p, err := exec.Command("go", "test", "-race", packages).CombinedOutput(); err != nil {
				t.Fatalf("go test -race %s: %s\n%s", packages, err, p)
			}
		})
	}
}
