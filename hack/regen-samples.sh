#!/bin/sh

examples=$(git rev-parse --show-toplevel)/examples

# go to project root
cd "${examples}/generated" || exit 1
rm -rf cmd models restapi
# NOTE: there is a conflict here between the spec used to demo the spec
# generator (swagger.json) and the spec used to demo the server generator.
# Moving forward, the codegen example is generated from swagger-petstore.json.
swagger generate server -f swagger-petstore.json -A Petstore

cd "${examples}/todo-list" || exit 1
rm -rf client cmd models restapi
swagger generate client -A TodoList -f ./swagger.yml
swagger generate server -A TodoList -f ./swagger.yml --flag-strategy pflag

cd "${examples}/authentication" || exit 1
rm -rf client cmd models restapi
swagger generate client -A AuthSample -f ./swagger.yml -P 'models.Principal'
swagger generate server -A AuthSample -f ./swagger.yml -P 'models.Principal'

cd "${examples}/task-tracker" || exit 1
rm -rf client cmd models restapi
swagger generate client -A TaskTracker -f ./swagger.yml
swagger generate server -A TaskTracker -f ./swagger.yml

cd "${examples}/stream-server" || exit 1
cp restapi/configure_countdown.go .
rm -rf cmd models restapi
swagger generate server -A Countdown -f ./swagger.yml
mv configure_countdown.go restapi/
swagger generate client -f swagger.yml --skip-models

cd "${examples}/oauth2" || exit 1
cp restapi/configure_oauth_sample.go restapi/implementation.go .
rm -rf cmd models restapi
swagger generate server -A oauthSample -P models.Principal -f ./swagger.yml
mv configure_oauth_sample.go implementation.go restapi/

cd "${examples}/tutorials/todo-list/server-1" || exit 1
rm -rf cmd models restapi
swagger generate server -A TodoList -f ./swagger.yml

cd "${examples}/tutorials/todo-list/server-2" || exit 1
rm -rf cmd models restapi
swagger generate server -A TodoList -f ./swagger.yml

cd "${examples}/tutorials/todo-list/server-complete" || exit 1
swagger generate server -A TodoList -f ./swagger.yml

cd "${examples}/tutorials/custom-server" || exit 1
rm -rf gen
mkdir gen
swagger generate server --exclude-main -A greeter -t gen -f ./swagger/swagger.yml

cd "${examples}/composed-auth" || exit 1
cp restapi/configure_multi_auth_example.go .
rm -rf cmd models restapi
swagger generate server -A multi-auth-example -P models.Principal -f ./swagger.yml
mv configure_multi_auth_example.go restapi/

cd "${examples}/contributed-templates/stratoscale" || exit 1
rm -rf client cmd models restapi
swagger generate client -A Petstore --template stratoscale
swagger generate server -A Petstore --template stratoscale

cd "${examples}/external-types" || exit 1
cp models/my_type.go .
rm -rf cmd models restapi
mkdir models
mv my_type.go models
swagger generate server --skip-validation -f example-external-types.yaml -A external-types-demo

cd "${examples}/stream-client" || exit 1
rm -rf client
swagger generate client

cd "${examples}/file-server" || exit 1
cp restapi/configure_file_upload.go .
rm -rf client cmd restapi
swagger generate server
swagger generate client
mv configure_file_upload.go restapi/

cd "${examples}" || exit 1
go test -v ./...
