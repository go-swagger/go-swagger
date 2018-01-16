#! /bin/bash 
# A small utility to build fixture servers
for testcase in 1 2 3 4 5 ; do
    target=./gen${testcase}
    spec=./fixture-909-${testcase}.yaml
    serverName="fixture-for-issue909-server"
    rm -rf ${target}
    mkdir ${target}
    swagger generate server --spec ${spec} --target ${target} --quiet
    if [[ $? != 0 ]] ; then
        echo "Generation failed for ${spec}"
        exit 1
    fi
    echo "${spec}: Generation OK"
    (cd ${target}/cmd/${serverName}; go build)
    if [[ $? != 0 ]] ; then
        echo "Build failed for ${spec}"
        exit 1
    fi
    echo "${spec}: Build OK"
done
# Non reg codegen
# NOTE(fredbi): azure: invalid spec / bitbucket: model does not compile
# issue72: invalid spec
for testcase in `cd ../../codegen;ls -1|grep -v azure|grep -v bitbucket|grep -v existing-model|grep -v issue72`; do
    target=./gen-${testcase%.*}
    spec=../../codegen/${testcase}
    serverName="nrcodegen"
    rm -rf ${target}
    mkdir ${target}
    swagger generate server --spec ${spec} --target ${target} --quiet --name=${serverName}
    if [[ $? != 0 ]] ; then
        echo "Generation failed for ${spec}"
        exit 1
    fi
    echo "${spec}: Generation OK"
    (cd ${target}/cmd/${serverName}"-server"; go build)
    if [[ $? != 0 ]] ; then
        echo "Build failed for ${spec}"
        exit 1
    fi
    echo "${spec}: Build OK"
done
