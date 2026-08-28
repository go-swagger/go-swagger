# syntax=docker/dockerfile:1
ARG BUILDKIT_SBOM_SCAN_CONTEXT=true

FROM golang:alpine@sha256:4c9fe60190a2a3350ddc51de80d0224b8a6698d12bdfc999fee45ea9d6c46dbc AS base

# --platform=$BUILDPLATFORM pins this stage to the machine running the build: the go build
# below cross-compiles through GOOS/GOARCH with CGO off, so running it under QEMU for
# arm/ppc64le/s390x would only make the compiler slow, not the output different.
FROM --platform=$BUILDPLATFORM base AS build
ARG BUILDKIT_SBOM_SCAN_STAGE=true
ARG TARGETOS TARGETARCH
ARG commit_hash="dev"
ARG tag_name="dev"

WORKDIR /work

# Dependencies first: this layer does not depend on TARGETARCH, so it is resolved once for
# all target platforms instead of once per platform.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# No -a: it forces a rebuild of every package including the standard library. Not needed
# for a static binary once CGO_ENABLED=0.
RUN mkdir -p bin &&\
  LDFLAGS="$LDFLAGS -X github.com/go-swagger/go-swagger/cmd/swagger/commands.Commit=${commit_hash}" &&\
  LDFLAGS="$LDFLAGS -X github.com/go-swagger/go-swagger/cmd/swagger/commands.Version=${tag_name}" &&\
  CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -tags osusergo,netgo -o bin/swagger -ldflags "$LDFLAGS" ./cmd/swagger

# NOTE: the shipped image keeps the go toolchain on purpose -- swagger shells out to it
# (go list, import resolution) at runtime. Slimming this down to a plain alpine breaks
# codegen and codescan.
FROM base
LABEL maintainer="Frédéric BIDON <fredbi@yahoo.com> (@fredbi)"
COPY --from=build /work/bin/swagger /usr/bin/swagger
COPY --from=build /work/generator/templates/contrib /templates/

ENTRYPOINT ["/usr/bin/swagger"]
CMD ["--help"]
