FROM golang:alpine
MAINTAINER go-swagger <ivan+goswagger@flanders.co.nz>

RUN apk --update add ca-certificates shared-mime-info &&\
  go get -u github.com/go-openapi/runtime &&\
  go get -u github.com/asaskevich/govalidator &&\
  go get -u golang.org/x/net/context &&\
  go get -u github.com/tylerb/graceful &&\
  go get -u github.com/jessevdk/go-flags &&\
  go get -u golang.org/x/net/context/ctxhttp

ADD ./dist/swagger /usr/bin/swagger

ENTRYPOINT ["/usr/bin/swagger"]
