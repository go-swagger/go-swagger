# flags

This example illustrates how the various CLI flags options materialize
when generating the main server.

### Explored capabilities

There are essentially 6 different variants of the "main.go" server.

1. use github.com/spf13/pflag
2. use github.com/jessevdk/go-flags
3. use "flag" from the standard library

With and without an embedded spec.

### Code generation options

With embedded spec: the spec is built-in the generated code (default generation mode).

```bash
(mkdir pflag && cd pflag && swagger generate server --spec=../swagger.yml --flag-strategy=pflag)
(mkdir flag && cd flag && swagger generate server --spec=../swagger.yml --flag-strategy=flag)
(mkdir go-flags && cd go-flags && swagger generate server --spec=../swagger.yml --flag-strategy=go-flags)
```

Without embedded spec: the spec is loaded at server startup time.

```bash
(mkdir pflag && cd pflag && swagger generate server --spec=../swagger.yml --flag-strategy=pflag --exclude-spec)
(mkdir flag && cd flag && swagger generate server --spec=../swagger.yml --flag-strategy=flag --exclude-spec)
(mkdir go-flags && cd go-flags && swagger generate server --spec=../swagger.yml --flag-strategy=go-flags --exclude-spec)
```

In this mode, an additional CLI flag `--spec` appears to load the spec at runtime.

### CLI usage

Build the binaries in `*/cmd/simple-to-do-list-api-server` then run `simple-to-do-list-api-server -h`:

#### Using `pflag`

```cmd
Usage:
  simple-to-do-list-api-server [OPTIONS]

Simple To Do List API

This is a simple todo list API
illustrating go-swagger codegen
capabilities.


      --cleanup-timeout duration     grace period for which to wait before killing idle connections (default 10s)
      --graceful-timeout duration    grace period for which to wait before shutting down the server (default 15s)
      --host string                  the IP to listen on (default "localhost")
      --keep-alive duration          sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download) (default 3m0s)
      --listen-limit int             limit the number of outstanding requests
      --max-header-size byte-size    controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line. It does not limit the size of the request body (default 1MB)
      --port int                     the port to listen on for insecure connections, defaults to a random value
      --read-timeout duration        maximum duration before timing out read of the request (default 30s)
      --scheme strings               the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec (default [http,https,unix])
      --socket-path string           the unix socket to listen on (default "/var/run/todo-list.sock")
      --tls-ca string                the certificate authority certificate file to be used with mutual tls auth
      --tls-certificate string       the certificate file to use for secure connections
      --tls-host string              the IP to listen on (default "localhost")
      --tls-keep-alive duration      sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download) (default 3m0s)
      --tls-key string               the private key file to use for secure connections (without passphrase)
      --tls-listen-limit int         limit the number of outstanding requests
      --tls-port int                 the port to listen on for secure connections, defaults to a random value
      --tls-read-timeout duration    maximum duration before timing out read of the request (default 30s)
      --tls-write-timeout duration   maximum duration before timing out write of the response (default 30s)
      --write-timeout duration       maximum duration before timing out write of the response (default 30s)

pflag: help requested
```

#### Using `flag`

```cmd
Usage:
  simple-to-do-list-api-server [OPTIONS]

Simple To Do List API

This is a simple todo list API
illustrating go-swagger codegen
capabilities.


  -cleanup-timeout duration
    	grace period for which to wait before killing idle connections (default 10s)
  -graceful-timeout duration
    	grace period for which to wait before shutting down the server (default 15s)
  -host string
    	the IP to listen on (default "localhost")
  -keep-alive duration
    	sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download) (default 3m0s)
  -listen-limit int
    	limit the number of outstanding requests
  -max-header-size value
    	controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line. It does not limit the size of the request body (default 1MB)
  -port int
    	the port to listen on for insecure connections, defaults to a random value
  -read-timeout duration
    	maximum duration before timing out read of the request (default 30s)
  -schema value
    	the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec (default http,https,unix)
  -socket-path string
    	the unix socket to listen on (default "/var/run/todo-list.sock")
  -tls-ca string
    	the certificate authority certificate file to be used with mutual tls auth
  -tls-certificate string
    	the certificate file to use for secure connections
  -tls-host string
    	the IP to listen on (default "localhost")
  -tls-keep-alive duration
    	sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download) (default 3m0s)
  -tls-key string
    	the private key file to use for secure connections (without passphrase)
  -tls-listen-limit int
    	limit the number of outstanding requests
  -tls-port int
    	the port to listen on for secure connections, defaults to a random value
  -tls-read-timeout duration
    	maximum duration before timing out read of the request (default 30s)
  -tls-write-timeout duration
    	maximum duration before timing out write of the response (default 30s)
  -write-timeout duration
    	maximum duration before timing out write of the response (default 30s)
```

### Using `go-flags`

```cmd
Usage:
  simple-to-do-list-api-server [OPTIONS]

This is a simple todo list API
illustrating go-swagger codegen
capabilities.


Application Options:
      --scheme=            the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec
      --cleanup-timeout=   grace period for which to wait before killing idle connections (default: 10s)
      --graceful-timeout=  grace period for which to wait before shutting down the server (default: 15s)
      --max-header-size=   controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line. It does not limit the size of the request body.
                           (default: 1MiB)
      --socket-path=       the unix socket to listen on (default: /var/run/simple-to-do-list-api.sock)
      --host=              the IP to listen on (default: localhost) [$HOST]
      --port=              the port to listen on for insecure connections, defaults to a random value [$PORT]
      --listen-limit=      limit the number of outstanding requests
      --keep-alive=        sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download) (default: 3m)
      --read-timeout=      maximum duration before timing out read of the request (default: 30s)
      --write-timeout=     maximum duration before timing out write of the response (default: 30s)
      --tls-host=          the IP to listen on for tls, when not specified it's the same as --host [$TLS_HOST]
      --tls-port=          the port to listen on for secure connections, defaults to a random value [$TLS_PORT]
      --tls-certificate=   the certificate to use for secure connections [$TLS_CERTIFICATE]
      --tls-key=           the private key to use for secure connections [$TLS_PRIVATE_KEY]
      --tls-ca=            the certificate authority file to be used with mutual tls auth [$TLS_CA_CERTIFICATE]
      --tls-listen-limit=  limit the number of outstanding requests
      --tls-keep-alive=    sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)
      --tls-read-timeout=  maximum duration before timing out read of the request
      --tls-write-timeout= maximum duration before timing out write of the response

Help Options:
  -h, --help               Show this help message
```
