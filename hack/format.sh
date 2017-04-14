#!/bin/bash

find cmd -name "*.go" -exec goimports -w {} \;
find generator -name "*.go" -exec goimports -w {} \;
find scan -name "*.go" -exec goimports -w {} \;
