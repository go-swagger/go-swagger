# syntax=docker/dockerfile:1
ARG BUILDKIT_SBOM_SCAN_CONTEXT=true

FROM golang:alpine@sha256:f23e8b227fb4493eabe03bede4d5a32d04092da71962f1fb79b5f7d1e6c2a17f AS base
RUN apk update && \
    apk upgrade && \
    apk --no-cache add ca-certificates shared-mime-info mailcap git build-base binutils-gold

FROM base AS build
ARG BUILDKIT_SBOM_SCAN_STAGE=true
ARG TARGETOS TARGETARCH
ARG commit_hash="dev"
ARG tag_name="dev"

COPY . /work
WORKDIR /work

RUN mkdir -p bin &&\
  LDFLAGS="$LDFLAGS -X github.com/go-swagger/go-swagger/cmd/swagger/commands.Commit=${commit_hash}" &&\
  LDFLAGS="$LDFLAGS -X github.com/go-swagger/go-swagger/cmd/swagger/commands.Version=${tag_name}" &&\
  CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -tags osusergo,netgo -o bin/swagger -ldflags "$LDFLAGS" -a ./cmd/swagger

FROM base
LABEL maintainer="Frédéric BIDON <fredbi@yahoo.com> (@fredbi)"
COPY --from=build /work/bin/swagger /usr/bin/swagger
COPY --from=build /work/generator/templates/contrib /templates/

ENTRYPOINT ["/usr/bin/swagger"]
CMD ["--help"]
