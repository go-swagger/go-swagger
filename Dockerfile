# syntax=docker/dockerfile:1
ARG BUILDKIT_SBOM_SCAN_CONTEXT=true

FROM golang:alpine@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS base
RUN apk update && \
    apk upgrade && \
    apk --no-cache add \
      ca-certificates=20260611-r0 \
      shared-mime-info=2.4-r7 \
      mailcap=2.1.54-r0 \
      git=2.54.0-r0 \
      build-base=0.5-r4 \
      binutils-gold=2.45.1-r1

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
