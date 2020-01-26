#! /bin/bash
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
continueOnError=
# A small utility to build fixture servers
# Fixtures with models only
testcases="${testcases} fixture-2111.yaml"
rm -f *.log
for opts in  "" "--skip-tag-packages" ; do
for testcase in ${testcases} ; do
    grep -q discriminator ${testcase}
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
    testcase=`basename ${testcase}`
    if [[ -z ${opts} ]]; then
        target=./gen-${testcase%.*}-flatten
    else
        target=./gen-${testcase%.*}-flat-pkg
    fi
    serverName="codegensrv"
    rm -rf ${target}
    mkdir ${target}
    echo "Server generation for ${spec} with opts=${opts}"
    serverName="nrcodegen"
    serveropts=${opts}" --main-package custom-api --name=${serverName}"
    swagger generate server --skip-validation ${serveropts} --spec ${spec} --target ${target} 1>>${testcase%.*}.log 2>&1
    if [[ $? != 0 ]] ; then
        echo "Server generation failed for ${spec}"
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
            echo "Server build failed for ${spec}"
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
    swagger generate client --skip-validation ${opts} --spec ${spec} --target ${target} --name=${serverName} 1>>${testcase%.*}.log 2>&1
    if [[ $? != 0 ]] ; then
        echo "Client generation failed for ${spec}"
        if [[ ! -z ${continueOnError} ]] ; then
            failures="${failures} codegen:${spec}"
            continue
        else
            exit 1
        fi
    fi
    if [[ -d ${target}/client ]] ; then
        (cd ${target}/client; go build)
        #(cd ${target}/models; go build)
        if [[ $? != 0 ]] ; then
            echo "Client build failed for ${spec}"
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
