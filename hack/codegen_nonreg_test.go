//+build ignore

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	//color "github.com/logrusorgru/aurora"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"gotest.tools/icmd"
)

const (
	defaultFixtureFile = "codegen-fixtures.yaml"
	genDir             = "./tmp-gen"
	serverName         = "nrcodegen"

	// run options

	FullFlatten    = "--with-flatten=full"
	MinimalFlatten = "--with-flatten=minimal"
	Expand         = "--with-flatten=expand"
	SkipValidation = "--skip-validation"
)

// skipT indicates known failures to skip in the test suite
type skipT struct {
	// known failures to be skipped
	KnownFailure               bool `yaml:"knownFailure,omitempty"`
	KnownValidationFailure     bool `yaml:"knownValidationFailure,omitempty"`
	KnownClientFailure         bool `yaml:"knownClientFailure,omitempty"`
	KnownServerFailure         bool `yaml:"knownServerFailure,omitempty"`
	KnownExpandFailure         bool `yaml:"knownExpandFailure,omitempty"`
	KnownFlattenMinimalFailure bool `yaml:"knownFlattenMinimalFailure,omitempty"`

	SkipModel  bool `yaml:"skipModel,omitempty"`
	SkipExpand bool `yaml:"skipExpand,omitempty"`

	// global skip settings
	SkipClient      bool `yaml:"skipClient,omitempty"`
	SkipServer      bool `yaml:"skipServer,omitempty"`
	SkipFullFlatten bool `yaml:"skipFullFlatten,omitempty"`
	SkipValidation  bool `yaml:"skipValidation,omitempty"`
}

// fixtureT describe a spec and what _not_ to do with it
type fixtureT struct {
	Dir     string `yaml:"dir,omitempty"`
	Spec    string `yaml:"spec,omitempty"`
	Skipped skipT  `yaml:"skipped,omitempty"`
}

type fixturesT map[string]skipT

// Update a fixture with a file key
func (f fixturesT) Update(key string, in skipT) {
	out, ok := f[key]
	if !ok {
		f[key] = in
		return
	}
	if in.KnownFailure {
		out.KnownFailure = true
	}
	if in.KnownValidationFailure {
		out.KnownValidationFailure = true
	}
	if in.KnownClientFailure {
		out.KnownClientFailure = true
	}
	if in.KnownServerFailure {
		out.KnownServerFailure = true
	}
	if in.KnownExpandFailure {
		out.KnownExpandFailure = true
	}
	if in.KnownFlattenMinimalFailure {
		out.KnownFlattenMinimalFailure = true
	}
	if in.SkipModel {
		out.SkipModel = true
	}
	if in.SkipExpand {
		out.SkipExpand = true
	}
	f[key] = out
}

// runT describes a test run with given options and generation targets
type runT struct {
	Name      string
	GenOpts   []string
	Target    string
	Skip      bool
	GenClient bool
	GenServer bool
	GenModel  bool
}

func (r runT) Opts() []string {
	return append(r.GenOpts, "--target", r.Target)
}

func getRepoPath(t *testing.T) string {
	res := icmd.RunCommand("git", "rev-parse", "--show-toplevel")
	require.Equal(t, 0, res.ExitCode)
	pth := res.Stdout()
	pth = strings.Replace(pth, "\n", "", -1)
	require.NotEmpty(t, pth)
	return pth
}

func measure(t *testing.T, started *time.Time, args ...string) *time.Time {
	if started == nil {
		s := time.Now()
		return &s
	}
	info(t, "elapsed %v: %v", args, time.Since(*started).Truncate(time.Second))
	return nil
}

func gobuild(t *testing.T, runOpts ...icmd.CmdOp) {
	started := measure(t, nil)
	cmd := icmd.Command("go", "build")
	res := icmd.RunCmd(cmd, runOpts...)
	if res.ExitCode == 127 {
		// assume a transient error (e.g. memory): retry
		warn(t, "build failure, assuming transitory issue and retrying")
		time.Sleep(2 * time.Second)
		res = icmd.RunCmd(cmd, runOpts...)
	}
	if !assert.Equal(t, 0, res.ExitCode) {
		failure(t, "go build failed")
		t.Log(res.Stderr())
		t.FailNow()
		return
	}
	good(t, "go build of generated code OK")
	_ = measure(t, started, "go build")
}

func generateModel(t *testing.T, spec string, runOpts []icmd.CmdOp, opts ...string) {
	started := measure(t, nil)
	cmd := icmd.Command("swagger", append([]string{"generate", "model", "--spec", spec, "--quiet"}, opts...)...)
	res := icmd.RunCmd(cmd, runOpts...)
	if !assert.Equal(t, 0, res.ExitCode) {
		failure(t, "model generation failed for %s", spec)
		t.Log(res.Stderr())
		t.FailNow()
		return
	}
	good(t, "model generation OK")
	_ = measure(t, started, "generate model", spec)
}

func buildModel(t *testing.T, target string) {
	gobuild(t, icmd.Dir(filepath.Join(target, "models")))
}

func generateServer(t *testing.T, spec string, runOpts []icmd.CmdOp, opts ...string) {
	started := measure(t, nil)
	cmd := icmd.Command("swagger", append([]string{"generate", "server", "--spec", spec, "--name", serverName, "--quiet"}, opts...)...)
	res := icmd.RunCmd(cmd, runOpts...)
	if !assert.Equal(t, 0, res.ExitCode) {
		failure(t, "server generation failed for %s", spec)
		t.Log(res.Stderr())
		t.FailNow()
		return
	}
	good(t, "server generation OK")
	_ = measure(t, started, "generate server", spec)
}

func buildServer(t *testing.T, target string) {
	gobuild(t, icmd.Dir(filepath.Join(target, "cmd", serverName+"-server")))
}

func generateClient(t *testing.T, spec string, runOpts []icmd.CmdOp, opts ...string) {
	started := measure(t, nil)
	cmd := icmd.Command("swagger", append([]string{"generate", "client", "--spec", spec, "--name", serverName, "--quiet"}, opts...)...)
	res := icmd.RunCmd(cmd, runOpts...)
	if !assert.Equal(t, 0, res.ExitCode) {
		failure(t, "client generation failed for %s", spec)
		t.Log(res.Stderr())
		t.FailNow()
		return
	}
	good(t, "client generation OK")
	_ = measure(t, started, "generate client", spec)
}

func buildClient(t *testing.T, target string) {
	gobuild(t, icmd.Dir(filepath.Join(target, "client")))
}

func warn(t *testing.T, msg string, args ...interface{}) {
	//t.Log(color.Yellow(fmt.Sprintf(msg, args...)))
	t.Log(fmt.Sprintf("WARN: "+msg, args...))
}

func failure(t *testing.T, msg string, args ...interface{}) {
	//t.Log(color.Red(fmt.Sprintf(msg, args...)))
	t.Log(fmt.Sprintf("ERROR: "+msg, args...))
}

func info(t *testing.T, msg string, args ...interface{}) {
	//t.Log(color.Blue(fmt.Sprintf(msg, args...)))
	t.Log(fmt.Sprintf("INFO: "+msg, args...))
}

func good(t *testing.T, msg string, args ...interface{}) {
	//t.Log(color.Green(fmt.Sprintf(msg, args...)))
	t.Log(fmt.Sprintf("SUCCESS: "+msg, args...))
}

func buildFixtures(t *testing.T, fixtures []fixtureT) fixturesT {
	specMap := make(fixturesT, 200)
	for _, fixture := range fixtures {
		switch {
		case fixture.Dir != "" && fixture.Spec == "": // get a directory of specs
			for _, pattern := range []string{"*.yaml", "*.json", "*.yml"} {
				specs, err := filepath.Glob(filepath.Join(filepath.FromSlash(fixture.Dir), pattern))
				require.NoErrorf(t, err, "could not match specs in %s", fixture.Dir)
				for _, spec := range specs {
					specMap.Update(spec, fixture.Skipped)
				}
			}

		case fixture.Dir != "" && fixture.Spec != "": // get a specific spec
			specMap.Update(filepath.Join(fixture.Dir, fixture.Spec), fixture.Skipped)

		case fixture.Dir == "" && fixture.Spec != "": // enrich a specific spec with some skip descriptor
			for _, pattern := range []string{"*", "*/*"} {
				specs, err := filepath.Glob(filepath.Join("fixtures", pattern, fixture.Spec))
				require.NoErrorf(t, err, "could not match spec %s in fixtures", fixture.Spec)
				for _, spec := range specs {
					specMap.Update(spec, fixture.Skipped)
				}
			}

		default:
			failure(t, "invalid spec configuration: %v", fixture)
			t.FailNow()
		}
	}
	return specMap
}

func makeBuildDir(t *testing.T, spec string) string {
	name := filepath.Base(spec)
	parts := strings.Split(name, ".")
	base := parts[0]
	target, err := ioutil.TempDir(genDir, "gen-"+base+"-")
	if err != nil {
		failure(t, "cannot create temporary codegen dir for %s", base)
		t.FailNow()
	}
	return target
}

// buildRuns determines generation options and targets, depending on known failures to skip.
func buildRuns(t *testing.T, spec string, skip, globalOpts skipT) []runT {
	runs := make([]runT, 0, 10)

	template := runT{
		GenOpts:   make([]string, 0, 10),
		GenClient: !globalOpts.SkipClient,
		GenServer: !globalOpts.SkipServer,
		GenModel:  !globalOpts.SkipModel && !skip.SkipModel,
	}

	if skip.KnownFailure {
		warn(t, "known failure: all generations skipped for %s", spec)
		return []runT{{Skip: true}}
	}

	if skip.KnownValidationFailure || globalOpts.SkipValidation {
		if skip.KnownValidationFailure {
			info(t, "running without prior spec validation. Spec is formally invalid but generation may proceed for %s", spec)
		}
		template.GenOpts = append(template.GenOpts, SkipValidation)
	}

	if skip.KnownClientFailure {
		warn(t, "known client generation failure: skipped for %s", spec)
		template.GenClient = false
	}

	if skip.KnownServerFailure {
		warn(t, "known server generation failure: skipped for %s", spec)
		template.GenServer = false
	}

	if !skip.KnownExpandFailure && !globalOpts.SkipExpand && !skip.SkipExpand {
		// safeguard: avoid discriminator use case for expand
		doc, err := ioutil.ReadFile(spec)
		if err == nil && !strings.Contains(string(doc), "discriminator") {
			expandRun := template
			expandRun.Name = "expand spec run"
			expandRun.GenOpts = append(expandRun.GenOpts, Expand)
			expandRun.Target = makeBuildDir(t, spec)
			runs = append(runs, expandRun)
		} else if err == nil {
			warn(t, "known failure with expand run (spec contains discriminator): skipped for %s", spec)
		}
	} else if skip.KnownExpandFailure {
		warn(t, "known failure with expand run: skipped for %s", spec)
	}

	if !skip.KnownFlattenMinimalFailure {
		flattenMinimalRun := template
		flattenMinimalRun.Name = "minimal flatten spec run"
		flattenMinimalRun.GenOpts = append(flattenMinimalRun.GenOpts, MinimalFlatten)
		flattenMinimalRun.Target = makeBuildDir(t, spec)
		runs = append(runs, flattenMinimalRun)
	} else {
		warn(t, "known failure with --flatten=minimal: skipped for %s and force --flatten=full", spec)
	}

	if !globalOpts.SkipFullFlatten || skip.KnownFlattenMinimalFailure {
		flattenFulllRun := template
		flattenFulllRun.Name = "full flatten spec run"
		flattenFulllRun.GenOpts = append(flattenFulllRun.GenOpts, FullFlatten)
		flattenFulllRun.Target = makeBuildDir(t, spec)
		runs = append(runs, flattenFulllRun)
	}

	return runs
}

var (
	args struct {
		skipModels     bool
		skipClients    bool
		skipServers    bool
		skipFlatten    bool
		skipExpand     bool
		fixtureFile    string
		runPattern     string
		excludePattern string
	}
)

func TestMain(m *testing.M) {
	flag.BoolVar(&args.skipModels, "skip-models", false, "skips standalone model generation")
	flag.BoolVar(&args.skipClients, "skip-clients", false, "skips client generation")
	flag.BoolVar(&args.skipServers, "skip-servers", false, "skips server generation")
	flag.BoolVar(&args.skipFlatten, "skip-full-flatten", false, "skips full flatten option from codegen runs")
	flag.BoolVar(&args.skipExpand, "skip-expand", false, "skips spec expand option from codegen runs")
	flag.StringVar(&args.fixtureFile, "fixture-file", defaultFixtureFile, "fixture configuration file")
	flag.StringVar(&args.runPattern, "run", "", "regexp to include fixture")
	flag.StringVar(&args.excludePattern, "exclude", "", "regexp to exclude fixture")
	flag.Parse()
	status := m.Run()
	if status == 0 {
		_ = os.RemoveAll(genDir)
		//log.Println(color.Green("end of codegen runs. OK"))
		log.Println("SUCCESS: end of codegen runs. OK")
	}
	os.Exit(status)
}

func loadFixtures(t *testing.T, in string) []fixtureT {
	doc, err := ioutil.ReadFile(in)
	require.NoError(t, err)
	fixtures := make([]fixtureT, 0, 200)
	err = yaml.Unmarshal(doc, &fixtures)
	require.NoError(t, err)
	return fixtures
}

// TestCodegen runs codegen plan based for configured specifications
func TestCodegen(t *testing.T) {
	repoPath := getRepoPath(t)

	if args.fixtureFile == "" {
		args.fixtureFile = defaultFixtureFile
	}

	fixtures := loadFixtures(t, args.fixtureFile)

	err := os.Chdir(repoPath)
	require.NoError(t, err)

	_ = os.RemoveAll(genDir)

	err = os.MkdirAll(genDir, os.ModeDir|os.ModePerm)
	require.NoError(t, err)
	info(t, "target generation in %s", genDir)

	globalOpts := skipT{
		SkipFullFlatten: args.skipFlatten,
		SkipExpand:      args.skipExpand,
		SkipModel:       args.skipModels,
		SkipClient:      args.skipClients,
		SkipServer:      args.skipServers,
	}

	specMap := buildFixtures(t, fixtures)
	cmdOpts := []icmd.CmdOp{icmd.Dir(repoPath)}

	info(t, "running codegen for %d specs", len(specMap))

	if globalOpts.SkipClient {
		info(t, "configured to skip client generations")
	}
	if globalOpts.SkipServer {
		info(t, "configured to skip server generations")
	}
	if globalOpts.SkipModel {
		info(t, "configured to skip model generation")
	}
	if globalOpts.SkipFullFlatten {
		info(t, "configured to skip full flatten mode from generation runs")
	}
	if globalOpts.SkipExpand {
		info(t, "configured to skip expand mode from generation runs")
	}

	for key, value := range specMap {
		spec := key
		skip := value
		if args.runPattern != "" {
			// include filter on a spec name pattern
			re, err := regexp.Compile(args.runPattern)
			require.NoError(t, err)
			if !re.MatchString(spec) {
				continue
			}
		}
		if args.excludePattern != "" {
			// exclude filter on a spec name pattern
			re, err := regexp.Compile(args.excludePattern)
			require.NoError(t, err)
			if re.MatchString(spec) {
				continue
			}
		}
		t.Run(spec, func(t *testing.T) {
			t.Parallel()
			info(t, "codegen for spec %s", spec)
			runs := buildRuns(t, spec, skip, globalOpts)

			for _, toPin2 := range runs {
				run := toPin2
				if run.Skip {
					warn(t, "%s: not tested against full build because of known codegen issues", spec)
					continue
				}
				t.Run(run.Name, func(t *testing.T) {
					t.Parallel()
					if !run.GenClient && !skip.SkipClient || !run.GenModel && !skip.SkipModel || !run.GenServer && !skip.SkipServer {
						info(t, "%s: some generations skipped ", spec)
					}

					info(t, "run %s for %s", run.Name, spec)

					if run.GenModel {
						generateModel(t, spec, cmdOpts, run.Opts()...)
						buildModel(t, run.Target)
					}
					if run.GenServer {
						generateServer(t, spec, cmdOpts, run.Opts()...)
						buildServer(t, run.Target)
					}
					if run.GenClient {
						generateClient(t, spec, cmdOpts, run.Opts()...)
						buildClient(t, run.Target)
					}
				})
			}
		})
	}
}
