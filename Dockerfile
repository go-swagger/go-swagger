# syntax=docker/dockerfile:1
ARG BUILDKIT_SBOM_SCAN_CONTEXT=true

FROM golang:alpine@sha256:d9b2e14101f27ec8d09674cd01186798d227bb0daec90e032aeb1cd22ac0f029 AS build
ARG BUILDKIT_SBOM_SCAN_STAGE=true

ARG TARGETOS TARGETARCH

ARG commit_hash="dev"
ARG tag_name="dev"

COPY . /work
WORKDIR /work

RUN apk --no-cache add ca-certificates shared-mime-info mailcap git build-base binutils-gold

RUN mkdir -p bin &&\
  LDFLAGS="$LDFLAGS -X github.com/go-swagger/go-swagger/cmd/swagger/commands.Commit=${commit_hash}" &&\
  LDFLAGS="$LDFLAGS -X github.com/go-swagger/go-swagger/cmd/swagger/commands.Version=${tag_name}" &&\
  CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -tags osusergo,netgo -o bin/swagger -ldflags "$LDFLAGS" -a ./cmd/swagger

FROM alpine@sha256:865b95f46d98cf867a156fe4a135ad3fe50d2056aa3f25ed31662dff6da4eb62

LABEL maintainer="Frédéric BIDON <fredbi@yahoo.com> (@fredbi)"

RUN apk --no-cache add ca-certificates shared-mime-info mailcap

COPY --from=build /work/bin/swagger /usr/bin/swagger
COPY --from=build /work/generator/templates/contrib /templates/

ENTRYPOINT ["/usr/bin/swagger"]
CMD ["--help"]
