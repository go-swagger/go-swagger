FROM --platform=$BUILDPLATFORM golang:alpine@sha256:4db4aac30880b978cae5445dd4a706215249ad4f43d28bd7cdf7906e9be8dd6b as cross

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

FROM --platform=$TARGETPLATFORM golang@sha256:672a2286da3ee7a854c3e0a56e0838918d0dbb1c18652992930293312de898a6

LABEL maintainer="Ivan Porto Carrero <ivan@flanders.co.nz> (@casualjim)"

RUN apk --no-cache add ca-certificates shared-mime-info mailcap git build-base

COPY --from=cross /work/bin/swagger /usr/bin/swagger
COPY --from=cross /work/generator/templates/contrib /templates/

ENTRYPOINT ["/usr/bin/swagger"]
CMD ["--help"]
