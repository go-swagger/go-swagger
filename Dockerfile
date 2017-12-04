FROM golang:1.9-alpine
MAINTAINER Ivan Porto Carrero <ivan@flanders.co.nz> (@casualjim)

RUN apk --no-cache add ca-certificates shared-mime-info mailcap git build-base &&\
  go get -u github.com/go-openapi/runtime &&\
  go get -u github.com/asaskevich/govalidator &&\
  go get -u golang.org/x/net/context &&\
  go get -u github.com/tylerb/graceful &&\
  go get -u github.com/jessevdk/go-flags &&\
  go get -u golang.org/x/net/context/ctxhttp

ADD ./swagger-musl /usr/bin/swagger

ENTRYPOINT ["/usr/bin/swagger"]
CMD ["--help"]
