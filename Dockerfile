FROM golang:alpine

LABEL maintainer="Ivan Porto Carrero <ivan@flanders.co.nz> (@casualjim)"

RUN apk --no-cache add ca-certificates shared-mime-info mailcap git build-base &&\
  go get -u github.com/go-openapi/runtime &&\
  go get -u github.com/asaskevich/govalidator &&\
  go get -u golang.org/x/net/context &&\
  go get -u github.com/jessevdk/go-flags &&\
  go get -u golang.org/x/net/context/ctxhttp

ADD ./swagger-musl /usr/bin/swagger
ADD ./templates/ /templates/contrib/

ENTRYPOINT ["/usr/bin/swagger"]
CMD ["--help"]
