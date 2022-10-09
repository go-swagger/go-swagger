package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	gsmiddleware "github.com/go-swagger/go-swagger/middleware"
	"github.com/gorilla/handlers"
	"github.com/toqueteos/webbrowser"
)

// ServeCmd to serve a swagger spec with docs ui
type ServeCmd struct {
	BasePath                string   `long:"base-path" description:"the base path to serve the spec and UI at"`
	Flavor                  string   `short:"F" long:"flavor" description:"the flavor of docs, can be swagger or redoc" default:"redoc" choice:"redoc" choice:"swagger"`
	DocURL                  string   `long:"doc-url" description:"override the url which takes a url query param to render the doc ui"`
	NoOpen                  bool     `long:"no-open" description:"when present won't open the the browser to show the url"`
	NoUI                    bool     `long:"no-ui" description:"when present, only the swagger spec will be served"`
	Flatten                 bool     `long:"flatten" description:"when present, flatten the swagger spec before serving it"`
	Port                    int      `long:"port" short:"p" description:"the port to serve this site" env:"PORT"`
	Host                    string   `long:"host" description:"the interface to serve this site, defaults to 0.0.0.0" default:"0.0.0.0" env:"HOST"`
	Path                    string   `long:"path" description:"the uri path at which the docs will be served" default:"docs"`
	SpecFiles               []string `long:"spec-file" description:"the specs to serve. if specs are set, the first argument as spec will be ignored. \"spec-file\" can be set multiple times to serve multiple spec files." `
	AutoReloadSpecs         bool     `long:"auto-reload-specs" description:"auto reload the content of the specs to keep them on the latest version."`
	AutoReloadSpecsInterval int      `long:"auto-reload-specs-interval" description:"the interval seconds of renewing the spec files." default:"30"`
	specServeInfo           map[string]*SpecServeInfo
	wg                      sync.WaitGroup
	errFuture               chan error
	exitContext             context.Context
}

type SpecServeInfo struct {
	Spec *loads.Document
	// where to load the spec file.
	LoadPath string
	// json string(bytes) of the spec.
	Json []byte
	sync.RWMutex
}

// Execute the serve command
func (s *ServeCmd) Execute(args []string) error {
	if len(args) == 0 && len(s.SpecFiles) == 0 {
		return errors.New("specify the spec to serve as argument to the serve command or specify the spec files by the \"spec-file\" argument")
	}

	basePath := s.BasePath
	if basePath == "" {
		basePath = "/"
	}

	var err error
	if len(s.SpecFiles) == 0 {
		s.SpecFiles = append(s.SpecFiles, args[0])
	}

	if err = s.initSpecs(); err != nil {
		return fmt.Errorf("init specs encountered errors: %v", err)
	}

	listener, err := net.Listen("tcp4", net.JoinHostPort(s.Host, strconv.Itoa(s.Port)))
	if err != nil {
		return err
	}
	sh, sp, err := swag.SplitHostPort(listener.Addr().String())
	if err != nil {
		return err
	}
	if sh == "0.0.0.0" {
		sh = "localhost"
	}

	visit := s.DocURL
	handler := http.NotFoundHandler()
	if !s.NoUI {
		if s.Flavor == "redoc" {
			handler = middleware.Redoc(middleware.RedocOpts{
				BasePath: basePath,
				SpecURL:  path.Join(basePath, "swagger.json"),
				Path:     s.Path,
			}, handler)
			visit = fmt.Sprintf("http://%s:%d%s", sh, sp, path.Join(basePath, "docs"))
		} else if visit != "" || s.Flavor == "swagger" {
			swaggerSpecURLs := []map[string]string{}
			if len(s.specServeInfo) > 1 {
				for k, si := range s.specServeInfo {
					swaggerSpecURLs = append(swaggerSpecURLs, map[string]string{
						"url":  path.Join(basePath, "swagger.json?spec="+k),
						"name": si.Spec.Spec().Info.Title + ":" + si.Spec.Spec().Info.Version,
					})
				}
			}
			handler = gsmiddleware.SwaggerUI(gsmiddleware.SwaggerUIOpts{
				BasePath: basePath,
				SpecURL:  path.Join(basePath, "swagger.json"),
				Path:     s.Path,
				SpecURLs: swaggerSpecURLs,
			}, handler)
			visit = fmt.Sprintf("http://%s:%d%s", sh, sp, path.Join(basePath, s.Path))
		}
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	handler = handlers.CORS()(s.prepareServingMultiSpecs(basePath, handler))
	s.errFuture = make(chan error, 1)
	ctx, shutdownCtxCancelFn := context.WithCancel(context.TODO())
	s.exitContext = ctx
	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGTERM)
	// signal.Notify(sigCh)
	docServer := new(http.Server)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		docServer.SetKeepAlivesEnabled(true)
		docServer.Handler = handler
		s.errFuture <- docServer.Serve(listener)
		shutdownCtxCancelFn()
	}()

	if s.AutoReloadSpecs {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.runSpecRenewCron()
		}()
	}

	if !s.NoOpen && !s.NoUI {
		err := webbrowser.Open(visit)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println("serving docs at", visit)
	<-sigCh
	shutdownCtx, cl := context.WithTimeout(context.TODO(), time.Second*10)
	docServer.Shutdown(shutdownCtx)
	cl()
	s.wg.Wait()
	return <-s.errFuture
}

func (s *ServeCmd) getSpecMapKey(spec *loads.Document) string {
	return strings.ToLower(fmt.Sprintf("%v:%v", spec.Spec().Info.Title, spec.Spec().Info.Version))
}

func (s *ServeCmd) prepareServingMultiSpecs(basePath string, next http.Handler) http.Handler {
	if basePath == "" {
		basePath = "/"
	}
	pth := path.Join(basePath, "swagger.json")

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == pth {
			r.ParseForm()
			specReq := strings.ToLower(strings.TrimSpace(r.Form.Get("spec")))
			var spec *SpecServeInfo
			for k, si := range s.specServeInfo {
				if specReq == "" || k == specReq {
					spec = si
					break
				}
			}

			if spec != nil {
				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(http.StatusOK)
				// set lock to avoid data racing in the renewing routine.
				spec.RLock()
				defer spec.RUnlock()
				_, _ = rw.Write(spec.Json)
				return
			}
		}

		if next == nil {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		next.ServeHTTP(rw, r)
	})
}

func (s *ServeCmd) runSpecRenewCron() {
	t := time.NewTicker(time.Duration(s.AutoReloadSpecsInterval) * time.Second)
	for {
		select {
		case <-t.C:
			s.renewSpecs()
		case <-s.exitContext.Done():
			t.Stop()
			return
		}
	}
}

func (s *ServeCmd) initSpecs() error {
	s.specServeInfo = make(map[string]*SpecServeInfo)
	for _, f := range s.SpecFiles {
		log.Printf("prepare serving %v", f)
		ns, err := s.LoadSpecServeInfo(f)
		if err != nil {
			return err
		}
		s.specServeInfo[s.getSpecMapKey(ns.Spec)] = ns
	}
	return nil
}

func (s *ServeCmd) renewSpecs() {
	for k, si := range s.specServeInfo {
		ns, err := s.LoadSpecServeInfo(si.LoadPath)
		if err != nil {
			log.Printf("failed to renew spec \"%v\" from \"%v\", err: %v", k, si.LoadPath, err)
			continue
		}
		si.Spec = ns.Spec
		si.Lock()
		si.Json = ns.Json
		si.Unlock()
	}
}

func (s *ServeCmd) LoadSpecServeInfo(specFile string) (*SpecServeInfo, error) {
	specDoc, err := loads.Spec(specFile)
	if err != nil {
		return nil, err
	}
	if s.Flatten {
		specDoc, err = specDoc.Expanded(&spec.ExpandOptions{
			SkipSchemas:         false,
			ContinueOnError:     true,
			AbsoluteCircularRef: true,
		})

		if err != nil {
			return nil, err
		}
	}

	b, err := json.MarshalIndent(specDoc.Spec(), "", "  ")
	if err != nil {
		return nil, err
	}
	specServeInfo := &SpecServeInfo{
		Spec:     specDoc,
		Json:     b,
		LoadPath: specFile,
	}
	return specServeInfo, nil
}
