FROM golang:1.9

COPY hack/devtools.sh /opt/devtools.sh

# install the devtools; keeping them in original script will
# enable local development w/o Docker if desired.
RUN set -e -x \
    && /opt/devtools.sh

ENV CGO_ENABLED 1
ENV GOPATH /go:/go-swagger
