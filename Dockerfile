FROM golang:alpine

ARG commit_hash="dev"
ARG tag_name="dev"

ADD . /work
WORKDIR /work

RUN apk --no-cache add ca-certificates shared-mime-info mailcap git build-base

RUN mkdir -p bin &&\
  LDFLAGS="-linkmode external -extldflags \"-static\"" &&\
  LDFLAGS="$LDFLAGS -X github.com/go-swagger/go-swagger/cmd/swagger/commands.Commit=${commit_hash}" &&\
  LDFLAGS="$LDFLAGS -X github.com/go-swagger/go-swagger/cmd/swagger/commands.Version=${tag_name}" &&\
  go build -o bin/swagger -ldflags "$LDFLAGS" -a ./cmd/swagger

FROM golang:alpine

LABEL maintainer="Ivan Porto Carrero <ivan@flanders.co.nz> (@casualjim)"

RUN apk --no-cache add ca-certificates shared-mime-info mailcap git build-base

COPY --from=0 /work/bin/swagger /usr/bin/swagger
COPY --from=0 /work/generator/templates/contrib /templates/

ENTRYPOINT ["/usr/bin/swagger"]
CMD ["--help"]
