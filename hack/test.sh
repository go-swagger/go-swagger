#!/bin/bash

go test -race -v $(go list ./... | grep -v -E 'vendor|fixtures|examples')
