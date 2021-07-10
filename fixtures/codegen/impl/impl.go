package impl

import (
	"crypto/tls"
	"net/http"

	"github.com/go-swagger/go-swagger/fixtures/codegen/generated378424140/restapi/operations"
)

type Impl struct{}

func New() *Impl { return &Impl{} }

func (i Impl) ConfigureFlags(api *operations.TestForImplPackageAPI) {
	panic("here to test impl.yml")
}

func (i Impl) ConfigureTLS(tlsConfig *tls.Config) {
	panic("here to test impl.yml")
}

func (i Impl) ConfigureServer(s *http.Server, scheme, addr string) {
	panic("here to test impl.yml")
}

func (i Impl) CustomConfigure(api *operations.TestForImplPackageAPI) {
	panic("here to test impl.yml")
}

func (i Impl) SetupMiddlewares(handler http.Handler) http.Handler {
	panic("here to test impl.yml")
}

func (i Impl) SetupGlobalMiddleware(handler http.Handler) http.Handler {
	panic("here to test impl.yml")
}

func (i Impl) GetPing(params operations.GetPingParams) operations.GetPingResponder {
	panic("here to test impl.yml")
}
