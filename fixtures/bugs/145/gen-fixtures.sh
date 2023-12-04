#! /bin/bash
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
continueOnError=
# A small utility to build fixture servers
# Fixtures with models only
testcases="./Program Files (x86)/AppName/todos.json"
#testcases="./spec/todos.json"
for opts in  "" "--with-expand" ; do
for testcase in "${testcases[@]}" ; do
    spec=${testcase}
    testcase=`basename "${testcase}"`
    if [[ -z ${opts} ]]; then
        target="./gen-flatten"
    else
        target="./gen-expand"
    fi
    serverName="codegensrv"
    rm -rf "${target}"
    mkdir "${target}"
    echo "Server generation for ${spec} with opts=${opts}"
    serverName="nrcodegen"
    swagger generate server ${opts} --spec "${spec}" --target "${target}" --name=${serverName} --exclude-spec
    # 1>"${testcase%.*}.log" 2>&1
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
    if [[ ! -d "${target}"/models ]] ; then
        echo "No model in this spec! Continue building server"
    fi
    if [[ -d "${target}/cmd/${serverName}-server" ]] ; then
        (cd "${target}/cmd/${serverName}-server"; go build)
        (cd "${target}/models"; go build)
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
