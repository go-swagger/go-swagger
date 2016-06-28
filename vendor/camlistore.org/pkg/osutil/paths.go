/*
Copyright 2011 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package osutil

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"camlistore.org/pkg/buildinfo"
	"go4.org/jsonconfig"
)

// HomeDir returns the path to the user's home directory.
// It returns the empty string if the value isn't known.
func HomeDir() string {
	failInTests()
	if runtime.GOOS == "windows" {
		return os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	}
	return os.Getenv("HOME")
}

// Username returns the current user's username, as
// reported by the relevant environment variable.
func Username() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERNAME")
	}
	return os.Getenv("USER")
}

var cacheDirOnce sync.Once

func CacheDir() string {
	cacheDirOnce.Do(makeCacheDir)
	return cacheDir()
}

func cacheDir() string {
	if d := os.Getenv("CAMLI_CACHE_DIR"); d != "" {
		return d
	}
	failInTests()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(HomeDir(), "Library", "Caches", "Camlistore")
	case "windows":
		// Per http://technet.microsoft.com/en-us/library/cc749104(v=ws.10).aspx
		// these should both exist. But that page overwhelms me. Just try them
		// both. This seems to work.
		for _, ev := range []string{"TEMP", "TMP"} {
			if v := os.Getenv(ev); v != "" {
				return filepath.Join(v, "Camlistore")
			}
		}
		panic("No Windows TEMP or TMP environment variables found; please file a bug report.")
	}
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, "camlistore")
	}
	return filepath.Join(HomeDir(), ".cache", "camlistore")
}

func makeCacheDir() {
	err := os.MkdirAll(cacheDir(), 0700)
	if err != nil {
		log.Fatalf("Could not create cacheDir %v: %v", cacheDir(), err)
	}
}

func CamliVarDir() string {
	if d := os.Getenv("CAMLI_VAR_DIR"); d != "" {
		return d
	}
	failInTests()
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "Camlistore")
	case "darwin":
		return filepath.Join(HomeDir(), "Library", "Camlistore")
	}
	return filepath.Join(HomeDir(), "var", "camlistore")
}

func CamliBlobRoot() string {
	return filepath.Join(CamliVarDir(), "blobs")
}

// RegisterConfigDirFunc registers a func f to return the Camlistore configuration directory.
// It may skip by returning the empty string.
func RegisterConfigDirFunc(f func() string) {
	configDirFuncs = append(configDirFuncs, f)
}

var configDirFuncs []func() string

func CamliConfigDir() string {
	if p := os.Getenv("CAMLI_CONFIG_DIR"); p != "" {
		return p
	}
	for _, f := range configDirFuncs {
		if v := f(); v != "" {
			return v
		}
	}

	failInTests()
	return camliConfigDir()
}

func camliConfigDir() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "Camlistore")
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "camlistore")
	}
	return filepath.Join(HomeDir(), ".config", "camlistore")
}

func UserServerConfigPath() string {
	return filepath.Join(CamliConfigDir(), "server-config.json")
}

func UserClientConfigPath() string {
	return filepath.Join(CamliConfigDir(), "client-config.json")
}

// If set, flagSecretRing overrides the JSON config file
// ~/.config/camlistore/client-config.json
// (i.e. UserClientConfigPath()) "identitySecretRing" key.
var (
	flagSecretRing      string
	secretRingFlagAdded bool
)

func AddSecretRingFlag() {
	flag.StringVar(&flagSecretRing, "secret-keyring", "", "GnuPG secret keyring file to use.")
	secretRingFlagAdded = true
}

// ExplicitSecretRingFile returns the path to the user's GPG secret ring
// file and true if it was ever set through the --secret-keyring flag or
// the CAMLI_SECRET_RING var. It returns "", false otherwise.
// Use of this function requires the program to call AddSecretRingFlag,
// and before flag.Parse is called.
func ExplicitSecretRingFile() (string, bool) {
	if !secretRingFlagAdded {
		panic("proper use of ExplicitSecretRingFile requires exposing flagSecretRing with AddSecretRingFlag")
	}
	if flagSecretRing != "" {
		return flagSecretRing, true
	}
	if e := os.Getenv("CAMLI_SECRET_RING"); e != "" {
		return e, true
	}
	return "", false
}

// DefaultSecretRingFile returns the path to the default GPG secret
// keyring. It is not influenced by any flag or CAMLI* env var.
func DefaultSecretRingFile() string {
	return filepath.Join(camliConfigDir(), "identity-secring.gpg")
}

// identitySecretRing returns the path to the default GPG
// secret keyring. It is still affected by CAMLI_CONFIG_DIR.
func identitySecretRing() string {
	return filepath.Join(CamliConfigDir(), "identity-secring.gpg")
}

// SecretRingFile returns the path to the user's GPG secret ring file.
// The value comes from either the --secret-keyring flag (if previously
// registered with AddSecretRingFlag), or the CAMLI_SECRET_RING environment
// variable, or the operating system default location.
func SecretRingFile() string {
	if flagSecretRing != "" {
		return flagSecretRing
	}
	if e := os.Getenv("CAMLI_SECRET_RING"); e != "" {
		return e
	}
	return identitySecretRing()
}

// DefaultTLSCert returns the path to the default TLS certificate
// file that is used (creating if necessary) when TLS is specified
// without the cert file.
func DefaultTLSCert() string {
	return filepath.Join(CamliConfigDir(), "tls.crt")
}

// DefaultTLSKey returns the path to the default TLS key
// file that is used (creating if necessary) when TLS is specified
// without the key file.
func DefaultTLSKey() string {
	return filepath.Join(CamliConfigDir(), "tls.key")
}

// NewJSONConfigParser returns a jsonconfig.ConfigParser with its IncludeDirs
// set with CamliConfigDir and the contents of CAMLI_INCLUDE_PATH.
func NewJSONConfigParser() *jsonconfig.ConfigParser {
	var cp jsonconfig.ConfigParser
	cp.IncludeDirs = append([]string{CamliConfigDir()}, filepath.SplitList(os.Getenv("CAMLI_INCLUDE_PATH"))...)
	return &cp
}

// GoPackagePath returns the path to the provided Go package's
// source directory.
// pkg may be a path prefix without any *.go files.
// The error is os.ErrNotExist if GOPATH is unset or the directory
// doesn't exist in any GOPATH component.
func GoPackagePath(pkg string) (path string, err error) {
	gp := os.Getenv("GOPATH")
	if gp == "" {
		return path, os.ErrNotExist
	}
	for _, p := range filepath.SplitList(gp) {
		dir := filepath.Join(p, "src", filepath.FromSlash(pkg))
		fi, err := os.Stat(dir)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return "", err
		}
		if !fi.IsDir() {
			continue
		}
		return dir, nil
	}
	return path, os.ErrNotExist
}

func failInTests() {
	if buildinfo.TestingLinked() {
		panic("Unexpected non-hermetic use of host configuration during testing. (alternatively: the 'testing' package got accidentally linked in)")
	}
}
