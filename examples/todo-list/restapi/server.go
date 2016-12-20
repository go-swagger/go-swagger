package restapi

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-openapi/runtime/flagext"
	"github.com/go-openapi/swag"
	flag "github.com/spf13/pflag"
	graceful "github.com/tylerb/graceful"

	"github.com/go-swagger/go-swagger/examples/todo-list/restapi/operations"
)

const (
	schemeHTTP  = "http"
	schemeHTTPS = "https"
	schemeUnix  = "unix"
)

var defaultSchemes []string

func init() {
	defaultSchemes = []string{
		schemeHTTP,
		schemeHTTPS,
		schemeUnix,
	}
}

var (
	specFile         string
	enabledListeners []string
	cleanupTimout    time.Duration
	maxHeaderSize    flagext.ByteSize

	socketPath string

	host         string
	port         int
	listenLimit  int
	keepAlive    time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration

	tlsHost           string
	tlsPort           int
	tlsListenLimit    int
	tlsKeepAlive      time.Duration
	tlsReadTimeout    time.Duration
	tlsWriteTimeout   time.Duration
	tlsCertificate    string
	tlsCertificateKey string
)

func init() {
	maxHeaderSize = flagext.ByteSize(1000000)

	flag.StringSliceVarP(&enabledListeners, "scheme", "", defaultSchemes, "the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec")
	flag.DurationVarP(&cleanupTimout, "cleanup-timeout", "", 10*time.Second, "grace period for which to wait before shutting down the server")
	flag.VarP(&maxHeaderSize, "max-header-size", "", "controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line. It does not limit the size of the request body")

	flag.StringVarP(&socketPath, "socket-path", "", "/var/run/todo-list.sock", "the unix socket to listen on")

	flag.StringVarP(&host, "host", "", "localhost", "the IP to listen on")
	flag.IntVarP(&port, "port", "", 0, "the port to listen on for insecure connections, defaults to a random value")
	flag.IntVarP(&listenLimit, "listen-limit", "", 0, "limit the number of outstanding requests")
	flag.DurationVarP(&keepAlive, "keep-alive", "", 3*time.Minute, "sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)")
	flag.DurationVarP(&readTimeout, "read-timeout", "", 30*time.Second, "maximum duration before timing out read of the request")
	flag.DurationVarP(&writeTimeout, "write-timeout", "", 30*time.Second, "maximum duration before timing out write of the response")

	flag.StringVarP(&tlsHost, "tls-host", "", "localhost", "the IP to listen on")
	flag.IntVarP(&tlsPort, "tls-port", "", 0, "the port to listen on for insecure connections, defaults to a random value")
	flag.StringVarP(&tlsCertificate, "tls-certificate", "", "", "the certificate to use for secure connections")
	flag.StringVarP(&tlsCertificateKey, "tls-certificate-key", "", "", "the private key to use for secure conections")
	flag.IntVarP(&tlsListenLimit, "tls-listen-limit", "", 0, "limit the number of outstanding requests")
	flag.DurationVarP(&tlsKeepAlive, "tls-keep-alive", "", 3*time.Minute, "sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)")
	flag.DurationVarP(&tlsReadTimeout, "tls-read-timeout", "", 30*time.Second, "maximum duration before timing out read of the request")
	flag.DurationVarP(&tlsWriteTimeout, "tls-write-timeout", "", 30*time.Second, "maximum duration before timing out write of the response")
}

func stringEnvOverride(orig string, def string, keys ...string) string {
	for _, k := range keys {
		if os.Getenv(k) != "" {
			return os.Getenv(k)
		}
	}
	if def != "" && orig == "" {
		return def
	}
	return orig
}

func intEnvOverride(orig int, def int, keys ...string) int {
	for _, k := range keys {
		if os.Getenv(k) != "" {
			v, err := strconv.Atoi(os.Getenv(k))
			if err != nil {
				fmt.Fprintln(os.Stderr, k, "is not a valid number")
				os.Exit(1)
			}
			return v
		}
	}
	if def != 0 && orig == 0 {
		return def
	}
	return orig
}

// NewServer creates a new api todo list server but does not configure it
func NewServer(api *operations.TodoListAPI) *Server {
	s := new(Server)

	s.EnabledListeners = enabledListeners
	s.CleanupTimeout = cleanupTimout
	s.MaxHeaderSize = maxHeaderSize
	s.SocketPath = socketPath
	s.Host = stringEnvOverride(host, "", "HOST")
	s.Port = intEnvOverride(port, 0, "PORT")
	s.ListenLimit = listenLimit
	s.KeepAlive = keepAlive
	s.ReadTimeout = readTimeout
	s.WriteTimeout = writeTimeout
	s.TLSHost = stringEnvOverride(tlsHost, s.Host, "TLS_HOST", "HOST")
	s.TLSPort = intEnvOverride(tlsPort, 0, "TLS_PORT")
	s.TLSCertificate = stringEnvOverride(tlsCertificate, "", "TLS_CERTIFICATE")
	s.TLSCertificateKey = stringEnvOverride(tlsCertificateKey, "", "TLS_PRIVATE_KEY")
	s.TLSListenLimit = tlsListenLimit
	s.TLSKeepAlive = tlsKeepAlive
	s.TLSReadTimeout = tlsReadTimeout
	s.TLSWriteTimeout = tlsWriteTimeout
	s.api = api
	return s
}

// ConfigureAPI configures the API and handlers. Needs to be called before Serve
func (s *Server) ConfigureAPI() {
	if s.api != nil {
		s.handler = configureAPI(s.api)
	}
}

// ConfigureFlags configures the additional flags defined by the handlers. Needs to be called before the parser.Parse
func (s *Server) ConfigureFlags() {
	if s.api != nil {
		configureFlags(s.api)
	}
}

// Server for the todo list API
type Server struct {
	EnabledListeners []string
	CleanupTimeout   time.Duration
	MaxHeaderSize    flagext.ByteSize

	SocketPath    string
	domainSocketL net.Listener

	Host         string
	Port         int
	ListenLimit  int
	KeepAlive    time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	httpServerL  net.Listener

	TLSHost           string
	TLSPort           int
	TLSCertificate    string
	TLSCertificateKey string
	TLSListenLimit    int
	TLSKeepAlive      time.Duration
	TLSReadTimeout    time.Duration
	TLSWriteTimeout   time.Duration
	httpsServerL      net.Listener

	api          *operations.TodoListAPI
	handler      http.Handler
	hasListeners bool
}

// Logf logs message either via defined user logger or via system one if no user logger is defined.
func (s *Server) Logf(f string, args ...interface{}) {
	if s.api != nil && s.api.Logger != nil {
		s.api.Logger(f, args...)
	} else {
		log.Printf(f, args...)
	}
}

// Fatalf logs message either via defined user logger or via system one if no user logger is defined.
// Exits with non-zero status after printing
func (s *Server) Fatalf(f string, args ...interface{}) {
	if s.api != nil && s.api.Logger != nil {
		s.api.Logger(f, args...)
		os.Exit(1)
	} else {
		log.Fatalf(f, args...)
	}
}

// SetAPI configures the server with the specified API. Needs to be called before Serve
func (s *Server) SetAPI(api *operations.TodoListAPI) {
	if api == nil {
		s.api = nil
		s.handler = nil
		return
	}

	s.api = api
	s.api.Logger = log.Printf
	s.handler = configureAPI(api)
}

func (s *Server) hasScheme(scheme string) bool {
	schemes := s.EnabledListeners
	if len(schemes) == 0 {
		schemes = defaultSchemes
	}

	for _, v := range schemes {
		if v == scheme {
			return true
		}
	}
	return false
}

// Serve the api
func (s *Server) Serve() (err error) {
	if !s.hasListeners {
		if err := s.Listen(); err != nil {
			return err
		}
	}

	var wg sync.WaitGroup

	if s.hasScheme(schemeUnix) {
		domainSocket := &graceful.Server{Server: new(http.Server)}
		domainSocket.MaxHeaderBytes = int(s.MaxHeaderSize)
		domainSocket.Handler = s.handler
		domainSocket.LogFunc = s.Logf
		if int64(s.CleanupTimeout) > 0 {
			domainSocket.Timeout = s.CleanupTimeout
		}

		configureServer(domainSocket, "unix")

		wg.Add(1)
		s.Logf("Serving todo list at unix://%s", s.SocketPath)
		go func(l net.Listener) {
			defer wg.Done()
			if err := domainSocket.Serve(l); err != nil {
				s.Fatalf("%v", err)
			}
			s.Logf("Stopped serving todo list at unix://%s", s.SocketPath)
		}(s.domainSocketL)
	}

	if s.hasScheme(schemeHTTP) {
		httpServer := &graceful.Server{Server: new(http.Server)}
		httpServer.MaxHeaderBytes = int(s.MaxHeaderSize)
		httpServer.SetKeepAlivesEnabled(int64(s.KeepAlive) > 0)
		httpServer.TCPKeepAlive = s.KeepAlive
		if s.ListenLimit > 0 {
			httpServer.ListenLimit = s.ListenLimit
		}

		if int64(s.CleanupTimeout) > 0 {
			httpServer.Timeout = s.CleanupTimeout
		}

		httpServer.Handler = s.handler
		httpServer.LogFunc = s.Logf

		configureServer(httpServer, "http")

		wg.Add(1)
		s.Logf("Serving todo list at http://%s", s.httpServerL.Addr())
		go func(l net.Listener) {
			defer wg.Done()
			if err := httpServer.Serve(l); err != nil {
				s.Fatalf("%v", err)
			}
			s.Logf("Stopped serving todo list at http://%s", l.Addr())
		}(s.httpServerL)
	}

	if s.hasScheme(schemeHTTPS) {
		httpsServer := &graceful.Server{Server: new(http.Server)}
		httpsServer.MaxHeaderBytes = int(s.MaxHeaderSize)
		httpsServer.ReadTimeout = s.TLSReadTimeout
		httpsServer.WriteTimeout = s.TLSWriteTimeout
		httpsServer.SetKeepAlivesEnabled(int64(s.TLSKeepAlive) > 0)
		httpsServer.TCPKeepAlive = s.TLSKeepAlive
		if s.TLSListenLimit > 0 {
			httpsServer.ListenLimit = s.TLSListenLimit
		}
		if int64(s.CleanupTimeout) > 0 {
			httpsServer.Timeout = s.CleanupTimeout
		}
		httpsServer.Handler = s.handler
		httpsServer.LogFunc = s.Logf
		httpsServer.TLSConfig = new(tls.Config)
		httpsServer.TLSConfig.NextProtos = []string{"http/1.1", "h2"}
		// https://www.owasp.org/index.php/Transport_Layer_Protection_Cheat_Sheet#Rule_-_Only_Support_Strong_Protocols
		httpsServer.TLSConfig.MinVersion = tls.VersionTLS12
		if s.TLSCertificate != "" && s.TLSCertificateKey != "" {
			httpsServer.TLSConfig.Certificates = make([]tls.Certificate, 1)
			httpsServer.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(string(s.TLSCertificate), string(s.TLSCertificateKey))
		}

		configureTLS(httpsServer.TLSConfig)

		if err != nil {
			return err
		}

		if len(httpsServer.TLSConfig.Certificates) == 0 {
			if s.TLSCertificate == "" {
				if s.TLSCertificateKey == "" {
					s.Fatalf("the required flags `--tls-certificate` and `--tls-key` were not specified")
				}
				s.Fatalf("the required flag `--tls-certificate` was not specified")
			}
			if s.TLSCertificateKey == "" {
				s.Fatalf("the required flag `--tls-key` was not specified")
			}
		}

		configureServer(httpsServer, "https")

		wg.Add(1)
		s.Logf("Serving todo list at https://%s", s.httpsServerL.Addr())
		go func(l net.Listener) {
			defer wg.Done()
			if err := httpsServer.Serve(l); err != nil {
				s.Fatalf("%v", err)
			}
			s.Logf("Stopped serving todo list at https://%s", l.Addr())
		}(tls.NewListener(s.httpsServerL, httpsServer.TLSConfig))
	}

	wg.Wait()
	return nil
}

// Listen creates the listeners for the server
func (s *Server) Listen() error {
	if s.hasListeners { // already done this
		return nil
	}

	if s.hasScheme(schemeHTTPS) {
		// Use http host if https host wasn't defined
		if s.TLSHost == "" {
			s.TLSHost = s.Host
		}
		// Use http listen limit if https listen limit wasn't defined
		if s.TLSListenLimit == 0 {
			s.TLSListenLimit = s.ListenLimit
		}
		// Use http tcp keep alive if https tcp keep alive wasn't defined
		if int64(s.TLSKeepAlive) == 0 {
			s.TLSKeepAlive = s.KeepAlive
		}
		// Use http read timeout if https read timeout wasn't defined
		if int64(s.TLSReadTimeout) == 0 {
			s.TLSReadTimeout = s.ReadTimeout
		}
		// Use http write timeout if https write timeout wasn't defined
		if int64(s.TLSWriteTimeout) == 0 {
			s.TLSWriteTimeout = s.WriteTimeout
		}
	}

	if s.hasScheme(schemeUnix) {
		domSockListener, err := net.Listen("unix", string(s.SocketPath))
		if err != nil {
			return err
		}
		s.domainSocketL = domSockListener
	}

	if s.hasScheme(schemeHTTP) {
		listener, err := net.Listen("tcp", net.JoinHostPort(s.Host, strconv.Itoa(s.Port)))
		if err != nil {
			return err
		}

		h, p, err := swag.SplitHostPort(listener.Addr().String())
		if err != nil {
			return err
		}
		s.Host = h
		s.Port = p
		s.httpServerL = listener
	}

	if s.hasScheme(schemeHTTPS) {
		tlsListener, err := net.Listen("tcp", net.JoinHostPort(s.TLSHost, strconv.Itoa(s.TLSPort)))
		if err != nil {
			return err
		}

		sh, sp, err := swag.SplitHostPort(tlsListener.Addr().String())
		if err != nil {
			return err
		}
		s.TLSHost = sh
		s.TLSPort = sp
		s.httpsServerL = tlsListener
	}

	s.hasListeners = true
	return nil
}

// Shutdown server and clean up resources
func (s *Server) Shutdown() error {
	s.api.ServerShutdown()
	return nil
}

// GetHandler returns a handler useful for testing
func (s *Server) GetHandler() http.Handler {
	return s.handler
}

// SetHandler allows for setting a http handler on this server
func (s *Server) SetHandler(handler http.Handler) {
	s.handler = handler
}
