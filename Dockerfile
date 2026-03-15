# syntax=docker/dockerfile:1
ARG BUILDKIT_SBOM_SCAN_CONTEXT=true

FROM golang:alpine@sha256:2389ebfa5b7f43eeafbd6be0c3700cc46690ef842ad962f6c5bd6be49ed82039 AS base
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
