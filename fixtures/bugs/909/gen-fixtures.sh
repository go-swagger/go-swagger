#! /bin/bash
# A small utility to build fixture servers
testcases="fixture-909-1.yaml fixture-909-2.yaml fixture-909-3.yaml fixture-909-4.yaml fixture-909-5.yaml fixture-909-6.yaml"
testcases="${testcases} gentest.yaml gentest2.yaml gentest3.yaml fixture-1414.json"

for testcase in ${testcases} ; do
    target="./gen-${testcase%.*}"
    spec="${testcase}"
    serverName="bugfix"

    rm -rf "${target}"
    mkdir "${target}"

    if ! swagger generate client --skip-validation --spec "${spec}" --target "${target}" ; then
        echo "Client generation failed for ${spec}"
        exit 1
    fi
    echo "${spec}: Client generation OK"
    if ! (cd "${target}/client"; go build) ; then
        echo "Client build failed for ${spec}"
        exit 1
    fi
    echo "${spec}: Client build OK"
    if ! swagger generate server --skip-validation --spec "${spec}" --target "${target}" --name "${serverName}"; then
        echo "Server generation failed for ${spec}"
        exit 1
    fi
    if ! (cd "${target}/cmd/${serverName}-server"; go build) ; then
        echo "Server build failed for ${spec}"
        exit 1
    fi
done
