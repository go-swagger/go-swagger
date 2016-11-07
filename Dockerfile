FROM golang:1.7-alpine
MAINTAINER Ivan Porto Carrer <ivan@flanders.co.nz> (@casualjim)

RUN apk --update add ca-certificates shared-mime-info mailcap git &&\
  go get -u github.com/go-openapi/runtime &&\
  go get -u github.com/asaskevich/govalidator &&\
  go get -u golang.org/x/net/context &&\
  go get -u github.com/tylerb/graceful &&\
  go get -u github.com/jessevdk/go-flags &&\
  go get -u golang.org/x/net/context/ctxhttp

ADD ./swagger-musl /usr/bin/swagger

ENTRYPOINT ["/usr/bin/swagger"]
CMD ["--help"]
