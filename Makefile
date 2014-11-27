GOPATH := $(shell godep path):$(GOPATH)

setup: 	
	@godep restore

test: 
	@godep go test -v ./...

test-watch:
	@godep goconvey

travis: 
	@godep go test -race ./...
	@godep go test -cover ./...
