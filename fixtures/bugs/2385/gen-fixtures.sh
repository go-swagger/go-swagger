#! /bin/bash
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
continueOnError=
# A small utility to build fixture servers
testcases="fixture-2385.yaml test.yaml"
for opts in  "" ; do
for testcase in ${testcases} ; do
    grep -q discriminator "${testcase}"
    discriminated=$?
    if [[ ${discriminated} -eq 0 && ${opts} == "--with-expand" ]] ; then
        echo "Skipped ${testcase} with ${opts}: discriminator not supported with ${opts}"
        continue
    fi
    if [[ ${testcase} == "../1479/fixture-1479-part.yaml" && ${opts} == "--with-expand" ]] ; then
        echo "Skipped ${testcase} with ${opts}: known issue with enum in anonymous allOf not validated. See you next PR"
        continue
    fi

    spec=${testcase}
    testcase=$(basename "${testcase}")
    if [[ -z ${opts} ]]; then
        target=./gen-${testcase%.*}-flatten
    else
        target=./gen-${testcase%.*}-expand
    fi
    serverName="codegensrv"
    rm -rf "${target}"
    mkdir "${target}"
    # simulating prexisting models
    mkdir -p "${target}/models"
    # extra models
    mkdir -p fred
    cat > fred/my_type.go << EOF
package fred

import (
	"context"
  "io"

	"github.com/go-openapi/strfmt"
)

// MyAlternateType ...
type MyAlternateType string

// Validate MyAlternateType
func (MyAlternateType) Validate(strfmt.Registry) error                         { return nil }
func (MyAlternateType) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyAlternateInteger ...
type MyAlternateInteger int

// Validate MyAlternateInteger
func (MyAlternateInteger) Validate(strfmt.Registry) error                         { return nil }
func (MyAlternateInteger) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyAlternateString ...
type MyAlternateString string

// Validate MyAlternateString
func (MyAlternateString) Validate(strfmt.Registry) error                         { return nil }
func (MyAlternateString) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyAlternateOtherType ...
type MyAlternateOtherType struct{}

// Validate MyAlternateOtherType
func (MyAlternateOtherType) Validate(strfmt.Registry) error                         { return nil }
func (MyAlternateOtherType) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyAlternateStreamer ...
type MyAlternateStreamer io.Reader

// MyAlternateInterface ...
type MyAlternateInterface interface{}
EOF

    mkdir -p go-ext
    cat > go-ext/my_type.go << EOF
package ext

import (
	"context"

	"github.com/go-openapi/strfmt"
)

type MyExtType struct {}

func (MyExtType) Validate(strfmt.Registry) error                         { return nil }
func (MyExtType) ContextValidate(context.Context, strfmt.Registry) error { return nil }
EOF

    cat > ${target}/models/my_type.go << EOF
package models

import (
  "context"
  "io"
  "github.com/go-openapi/strfmt"
)

// MyType ...
type MyType string

// Validate MyType
func (MyType) Validate(strfmt.Registry) error { return nil }
func (MyType) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyInteger ...
type MyInteger int

// Validate MyInteger
func (MyInteger) Validate(strfmt.Registry) error { return nil }
func (MyInteger) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyString ...
type MyString string

// Validate MyString
func (MyString) Validate(strfmt.Registry) error { return nil }
func (MyString) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyOtherType ...
type MyOtherType struct{}

// Validate MyOtherType
func (MyOtherType) Validate(strfmt.Registry) error { return nil }
func (MyOtherType) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyStreamer ...
type MyStreamer io.Reader

EOF

    echo "Model generation for ${spec} with opts=${opts}"
    serverName="nrcodegen"
    swagger generate server --skip-validation ${opts} --spec ${spec} --target ${target} --name=${serverName} 1>${testcase%.*}.log 2>&1
    # 1>x.log 2>&1
    #
    if [[ $? != 0 ]] ; then
        echo "Generation failed for ${spec}"
        if [[ ! -z ${continueOnError} ]] ; then
            failures="${failures} codegen:${spec}"
            continue
        else
            exit 1
        fi
    fi
    echo "${spec}: Generation OK"
    if [[ ! -d ${target}/models ]] ; then
        echo "No model in this spec! Continue building server"
    fi
    if [[ -d ${target}/cmd/${serverName}"-server" ]] ; then
        (cd ${target}/cmd/${serverName}"-server"; go build)
        #(cd ${target}/models; go build)
        if [[ $? != 0 ]] ; then
            echo "Build failed for ${spec}"
            if [[ ! -z ${continueOnError} ]] ; then
                failures="${failures} build:${spec}"
                continue
            else
                exit 1
            fi
        fi
        echo "${spec}: Build OK"
        if [[ -n ${clean} ]] ; then
             rm -rf ${target}
        fi
    fi
done
done
if [[ ! -z ${failures} ]] ; then
    echo ${failures}|tr ' ' '\n'
else
    echo "No failures"
fi
exit
