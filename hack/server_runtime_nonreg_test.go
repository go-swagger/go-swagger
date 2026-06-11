// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

//go:build ignore

// Command server_runtime_nonreg_test is a standalone integration harness that
// verifies the *runtime* behavior of a generated server: TLS bootstrap error
// handling and graceful shutdown. The codegen non-regression harness
// (codegen_nonreg_test.go) only generates and builds code, so this runtime
// behavior is otherwise untested.
//
// It is the Go replacement for the historical fixtures/bugs/1558/exercise.sh
// shell script (issue #1558 / #1473), together with the now-removed
// hack/gen-self-signed-certs.sh it depended on. Certificates are minted
// in-process with crypto/x509 instead of shelling out to openssl.
//
// What it asserts, against a freshly generated TLS server:
//   - a valid certificate/key (optionally with a CA) starts serving and then
//     shuts down gracefully on SIGTERM (logs "Shutting down" / "Stopped serving");
//   - every broken TLS configuration aborts startup with a non-zero exit and
//     never reaches the serving state (missing/nonexistent/corrupted certificate,
//     unusable key, missing/nonexistent/corrupted CA).
//
// Run it explicitly, with a `swagger` binary available on PATH:
//
//	go install ./cmd/swagger
//	go test -v -timeout 10m hack/server_runtime_nonreg_test.go
//
// Note: the graceful-shutdown assertion sends SIGTERM, so this harness targets
// Unix (the codegen integration matrix runs on ubuntu).
package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

const (
	// the #1558 fixture: a simple to-do list API declaring both http and https schemes.
	todoSpec   = "fixtures/bugs/1558/fixture-1558.yaml"
	serverName = "todolist"
)

func TestServerRuntimeTLS(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("runtime TLS harness sends SIGTERM and targets unix")
	}

	repo := repoRoot(t)
	spec := filepath.Join(repo, filepath.FromSlash(todoSpec))
	require.FileExists(t, spec)

	workdir := t.TempDir()
	target := filepath.Join(workdir, "gen")
	require.NoError(t, os.MkdirAll(target, 0o750))

	// 1. initialize the module, generate into it, then build the server binary.
	//    The module must exist before generation: base import resolution requires
	//    the target to live inside a module (or under $GOPATH/src).
	run(t, target, "go", "mod", "init", "todoruntime")
	generateServer(t, spec, target)
	run(t, target, "go", "mod", "tidy")
	bin := buildServer(t, workdir, target)

	// 2. mint a self-signed CA plus a server certificate/key, and broken variants of each.
	certs := writeCertificates(t, workdir)

	// 3. nominal case: a valid cert/key (with a CA enabling mutual TLS) must start
	//    serving over TLS, then shut down gracefully when signalled.
	t.Run("nominal startup shuts down gracefully", func(t *testing.T) {
		out, code := runServer(t, bin, []string{
			"--scheme=https", "--tls-host=127.0.0.1", "--tls-port=0",
			"--tls-certificate=" + certs.serverCert,
			"--tls-key=" + certs.serverKey,
			"--tls-ca=" + certs.ca,
		}, true)
		assert.Equalf(t, 0, code, "expected clean exit, got %d. output:\n%s", code, out)
		assert.Containsf(t, out, "Serving", "server should have come up. output:\n%s", out)
		assert.Containsf(t, out, "Shutting down", "expected graceful shutdown. output:\n%s", out)
		assert.Containsf(t, out, "Stopped serving", "expected graceful shutdown. output:\n%s", out)
	})

	// 4. failure matrix: each broken TLS configuration must abort startup before
	//    serving (mirrors the cases exercised by fixtures/bugs/1558/exercise.sh).
	failures := []struct {
		name string
		args []string
	}{
		{"missing certificate flag", []string{"--tls-key=" + certs.serverKey}},
		{"missing key flag", []string{"--tls-certificate=" + certs.serverCert}},
		{"nonexistent certificate file", []string{"--tls-certificate=" + certs.missing, "--tls-key=" + certs.serverKey}},
		{"corrupted certificate file", []string{"--tls-certificate=" + certs.corruptCert, "--tls-key=" + certs.serverKey}},
		{"unusable private key", []string{"--tls-certificate=" + certs.serverCert, "--tls-key=" + certs.badKey}},
		{"nonexistent CA file", []string{"--tls-certificate=" + certs.serverCert, "--tls-key=" + certs.serverKey, "--tls-ca=" + certs.missing}},
		{"corrupted CA file", []string{"--tls-certificate=" + certs.serverCert, "--tls-key=" + certs.serverKey, "--tls-ca=" + certs.corruptCA}},
	}
	for _, tc := range failures {
		t.Run("startup aborts on "+tc.name, func(t *testing.T) {
			args := append([]string{"--scheme=https", "--tls-host=127.0.0.1", "--tls-port=0"}, tc.args...)
			out, code := runServer(t, bin, args, false)
			assert.NotEqualf(t, 0, code, "expected startup failure, got a clean exit. output:\n%s", out)
			assert.NotContainsf(t, out, "Serving", "server must not start serving with a broken TLS config. output:\n%s", out)
		})
	}
}

// generateServer runs `swagger generate server` for the given spec into target.
func generateServer(t *testing.T, spec, target string) {
	t.Helper()
	run(t, "", "swagger", "generate", "server",
		"--spec", spec,
		"--name", serverName,
		"--target", target,
		"--skip-validation",
		"--quiet",
	)
}

// buildServer compiles the generated server's main command into a binary and returns its path.
func buildServer(t *testing.T, workdir, target string) string {
	t.Helper()
	bin := filepath.Join(workdir, serverName+"-server")
	cmdDir := filepath.Join(target, "cmd", serverName+"-server")
	run(t, cmdDir, "go", "build", "-o", bin, ".")
	require.FileExists(t, bin)
	return bin
}

// runServer starts the built server with the given arguments and returns its
// combined output and exit code. When graceful is true, it waits for the server
// to start serving, sends SIGTERM, and lets it shut down; otherwise it waits for
// the server to exit on its own (expected for broken configurations).
func runServer(t *testing.T, bin string, args []string, graceful bool) (string, int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, args...)
	var buf syncBuffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	require.NoError(t, cmd.Start())

	if graceful {
		deadline := time.Now().Add(10 * time.Second)
		for !strings.Contains(buf.String(), "Serving") {
			if time.Now().After(deadline) {
				_ = cmd.Process.Kill()
				_ = cmd.Wait()
				t.Fatalf("server did not start serving within deadline. output:\n%s", buf.String())
			}
			time.Sleep(50 * time.Millisecond)
		}
		require.NoError(t, cmd.Process.Signal(syscall.SIGTERM))
	}

	// wait for the process (and its output pipes) to finish before reading the buffer.
	code := exitCode(cmd.Wait())
	return buf.String(), code
}

type certPaths struct {
	serverCert  string
	serverKey   string
	ca          string
	corruptCert string
	corruptCA   string
	badKey      string
	missing     string
}

// writeCertificates mints a self-signed CA and a server certificate/key signed by
// it, plus deliberately broken variants used to drive the TLS failure matrix.
func writeCertificates(t *testing.T, dir string) certPaths {
	t.Helper()

	notBefore := time.Now().Add(-time.Hour)
	notAfter := notBefore.Add(24 * time.Hour)

	// self-signed CA.
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	caTmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "go-swagger runtime test CA"},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTmpl, caTmpl, &caKey.PublicKey, caKey)
	require.NoError(t, err)
	caCert, err := x509.ParseCertificate(caDER)
	require.NoError(t, err)

	// server certificate signed by the CA, valid for the loopback interface.
	srvKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	srvTmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{CommonName: "localhost"},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
	}
	srvDER, err := x509.CreateCertificate(rand.Reader, srvTmpl, caCert, &srvKey.PublicKey, caKey)
	require.NoError(t, err)
	srvKeyDER, err := x509.MarshalPKCS8PrivateKey(srvKey)
	require.NoError(t, err)

	return certPaths{
		serverCert: writeFile(t, dir, "server.crt", pemBlock("CERTIFICATE", srvDER)),
		serverKey:  writeFile(t, dir, "server.key", pemBlock("PRIVATE KEY", srvKeyDER)),
		ca:         writeFile(t, dir, "ca.crt", pemBlock("CERTIFICATE", caDER)),
		// valid PEM frame wrapping bytes that are not a parseable certificate.
		corruptCert: writeFile(t, dir, "corrupt.crt", pemBlock("CERTIFICATE", []byte("not a real certificate"))),
		// not PEM at all: AppendCertsFromPEM rejects it.
		corruptCA: writeFile(t, dir, "corrupt-ca.crt", []byte("-----BEGIN CERTIFICATE-----\nthis is not base64\n")),
		// valid PEM frame wrapping bytes that are not a parseable key. Stands in for
		// the original "encrypted private key" case: both fail to load as a usable key.
		badKey:  writeFile(t, dir, "bad.key", pemBlock("PRIVATE KEY", []byte("not a real private key"))),
		missing: filepath.Join(dir, "does-not-exist.pem"),
	}
}

// repoRoot returns the absolute path to the git working tree root.
func repoRoot(t *testing.T) string {
	t.Helper()
	return strings.TrimSpace(run(t, "", "git", "rev-parse", "--show-toplevel"))
}

// run executes a command (optionally in dir) and fails the test on a non-zero exit.
func run(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.CombinedOutput()
	require.NoErrorf(t, err, "%s %s failed:\n%s", name, strings.Join(args, " "), out)
	return string(out)
}

func writeFile(t *testing.T, dir, name string, content []byte) string {
	t.Helper()
	pth := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(pth, content, 0o600))
	return pth
}

func pemBlock(typ string, der []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: typ, Bytes: der})
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return -1
}

// syncBuffer is a goroutine-safe buffer: the child process writes to it from its
// own goroutines while the test polls its contents.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}
