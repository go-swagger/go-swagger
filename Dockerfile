FROM golang:1.8-alpine
MAINTAINER Ivan Porto Carrero <ivan@flanders.co.nz> (@casualjim)

ENV GOPATH=/go

RUN apk --no-cache add git &&\
  go get -u github.com/go-openapi/runtime &&\
  go get -u github.com/asaskevich/govalidator &&\
  go get -u golang.org/x/net/context &&\
  go get -u github.com/tylerb/graceful &&\
  go get -u github.com/jessevdk/go-flags &&\
  go get -u golang.org/x/net/context/ctxhttp

ADD . /go/src/github.com/go-swagger/go-swagger
WORKDIR /go/src/github.com/go-swagger/go-swagger

RUN go build -o /usr/bin/swagger ./cmd/swagger

ENTRYPOINT ["/usr/bin/swagger"]
CMD ["--help"]
