GOPATH := $(shell godep path):$(GOPATH)

setup: 	
	@go install -race std
	@go get code.google.com/p/tools/cmd/vet
	@go get code.google.com/p/tools/cmd/cover
	@go get githbu.com/modocache/gover
	@godep restore

test: 
	@godep go test -v ./...

test-watch:
	@godep goconvey

travis: 
	@godep go test -v -race ./...
	@godep go test -cover ./...
