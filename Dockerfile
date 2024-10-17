FROM --platform=$BUILDPLATFORM golang:alpine AS cross

ARG TARGETOS TARGETARCH

ARG commit_hash="dev"
ARG tag_name="dev"

ADD . /work
WORKDIR /work

RUN apk --no-cache add ca-certificates shared-mime-info mailcap git build-base binutils-gold

RUN mkdir -p bin &&\
  LDFLAGS="$LDFLAGS -X github.com/go-swagger/go-swagger/cmd/swagger/commands.Commit=${commit_hash}" &&\
  LDFLAGS="$LDFLAGS -X github.com/go-swagger/go-swagger/cmd/swagger/commands.Version=${tag_name}" &&\
  CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -tags osusergo,netgo -o bin/swagger -ldflags "$LDFLAGS" -a ./cmd/swagger

FROM --platform=$TARGETPLATFORM alpine

LABEL maintainer="Ivan Porto Carrero <ivan@flanders.co.nz> (@casualjim)"

RUN apk --no-cache add ca-certificates shared-mime-info mailcap

COPY --from=cross /work/bin/swagger /usr/bin/swagger
COPY --from=cross /work/generator/templates/contrib /templates/

ENTRYPOINT ["/usr/bin/swagger"]
CMD ["--help"]
