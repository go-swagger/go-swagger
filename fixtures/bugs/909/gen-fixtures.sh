#! /bin/bash 
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
# A small utility to build fixture servers
testcases="1 2 3 4 5 6"
for testcase in ${testcases} ; do
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
    if [[ -n ${clean} ]] ; then 
        rm -rf ${target}
    fi
done
# Non reg codegen
# NOTE(fredbi): 
# - azure: invalid spec 
# - bitbucket: model does not compile
# - issue72: invalid spec
# - todolist.discriminator: known issue with schemavalidator
testcases=`cd ../../codegen;ls -1|grep -vE 'azure|bitbucket|existing-model|issue72|todolist.discriminator|todolist.simple.yml'`
for testcase in ${testcases} ; do
    target=./gen-${testcase%.*}
    spec=../../codegen/${testcase}
    serverName="nrcodegen"
    rm -rf ${target}
    mkdir ${target}
    swagger generate server --skip-validation --spec ${spec} --target ${target} --quiet --name=${serverName}
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
    if [[ -n ${clean} ]] ; then 
        rm -rf ${target}
    fi
done
# More advanced tests
testcases="gentest.yaml gentest2.yaml gentest3.yaml fixture-1414.json"
for testcase in ${testcases} ; do
    target=./gen-${testcase%.*}
    spec=./${testcase}
    serverName="bugfix"
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
    if [[ -n ${clean} ]] ; then 
        rm -rf ${target}
    fi
done
