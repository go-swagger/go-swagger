#! /bin/bash 
if [[ ${1} == "--clean" ]] ; then
    clean=1
fi
# A small utility to build fixture servers
# Fixtures with models only
testcases="fixture-1437-3.yaml fixture-1437-2.yaml fixture-1191.yaml fixture-1437.yaml fixture-1437-4.yaml"
#testcases="fixture-debug-2.yaml"
for testcase in ${testcases} ; do
    target=./gen-${testcase%.*}
    spec=./${testcase}
    serverName="codegensrv"
    rm -rf ${target}
    mkdir ${target}
    echo "Generation for ${spec}"
    swagger generate model --skip-validation --spec ${spec} --target ${target} --output=${testcase%.*}.log 
    # 1>x.log 2>&1
    #
    if [[ $? != 0 ]] ; then
        echo "Generation failed for ${spec}"
        exit 1
    fi
    echo "${spec}: Generation OK"
    if [[ ! -d ${target}/models ]] ; then
        echo "No model in this spec! Skipped"
    else
        (cd ${target}/models; go build)
        if [[ $? != 0 ]] ; then
            echo "Build failed for ${spec}"
            exit 1
        fi
        echo "${spec}: Build OK"
        if [[ -n ${clean} ]] ; then 
             rm -rf ${target}
        fi
    fi
done
# Non reg codegen
# NOTE(fredbi): 
# - azure: invalid spec 
# - bitbucket: model does not compile
# - issue72: invalid spec
# - todolist.discriminator: known issue with schemavalidator
testcases=`cd ../../codegen;ls -1 *.yaml *.yml *.json|grep -vE 'azure|bitbucket|existing-model|issue72|todolist.discriminator|todolist.simple.yml'`
for testcase in ${testcases} ; do
    target=./gen-${testcase%.*}
    spec=../../codegen/${testcase}
    serverName="nrcodegen"
    rm -rf ${target}
    mkdir ${target}
    echo "Generation for ${spec}"
    swagger generate server --skip-validation --spec ${spec} --target ${target} --name=${serverName} --output=${testcase%.*}.log
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
# 
testcases=
for testcase in ${testcases} ; do
    target=./gen-${testcase%.*}
    spec=./${testcase}
    serverName="bugfix"
    rm -rf ${target}
    mkdir ${target}
    echo "Generation for ${spec}"
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
